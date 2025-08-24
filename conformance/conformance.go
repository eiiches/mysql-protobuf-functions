package main

import (
	"bufio"
	"database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/jsonoptionspb"
	"github.com/eiiches/mysql-protobuf-functions/internal/protonumberjson"
	"google.golang.org/protobuf/proto"
)

// ConformanceHandler implements the protobuf conformance test protocol
type ConformanceHandler struct {
	db                  *sql.DB
	debug               bool
	useLegacyConversion bool
}

// RunProtocol implements the conformance testing protocol:
// 1. Read 4-byte length from stdin (little endian)
// 2. Read N bytes representing a ConformanceRequest proto
// 3. Process the request and generate a ConformanceResponse
// 4. Write 4-byte length to stdout (little endian)
// 5. Write M bytes representing the ConformanceResponse proto
// 6. Repeat until stdin is closed
func (h *ConformanceHandler) RunProtocol() error {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if h.debug {
		log.Println("Starting MySQL conformance test protocol")
	}

	for {
		// Read length (4 bytes, little endian)
		lengthBytes := make([]byte, 4)
		_, err := io.ReadFull(reader, lengthBytes)
		if err == io.EOF {
			// Normal termination
			if h.debug {
				log.Println("Protocol terminated normally")
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read request length: %w", err)
		}

		requestLength := binary.LittleEndian.Uint32(lengthBytes)
		if h.debug {
			log.Printf("Reading request of length %d", requestLength)
		}

		// Read request data
		requestData := make([]byte, requestLength)
		_, err = io.ReadFull(reader, requestData)
		if err != nil {
			return fmt.Errorf("failed to read request data: %w", err)
		}

		// Process the request
		responseData, err := h.ProcessRequest(requestData)
		if err != nil {
			return fmt.Errorf("failed to process request: %w", err)
		}

		// Write response length (4 bytes, little endian)
		responseLength := uint32(len(responseData))
		responseLengthBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(responseLengthBytes, responseLength)

		if h.debug {
			log.Printf("Writing response of length %d", responseLength)
		}

		if _, err := writer.Write(responseLengthBytes); err != nil {
			return fmt.Errorf("failed to write response length: %w", err)
		}

		// Write response data
		if _, err := writer.Write(responseData); err != nil {
			return fmt.Errorf("failed to write response data: %w", err)
		}

		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush response: %w", err)
		}
	}
}

// ProcessRequest processes a single conformance request
func (h *ConformanceHandler) ProcessRequest(requestData []byte) ([]byte, error) {
	// Parse the ConformanceRequest
	var request ConformanceRequest
	if err := proto.Unmarshal(requestData, &request); err != nil {
		// If we can't parse the request, return a runtime error response
		response := &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: fmt.Sprintf("Failed to parse ConformanceRequest: %v", err),
			},
		}
		return proto.Marshal(response)
	}

	if h.debug {
		log.Printf("Processing request: MessageType=%s, Category=%s, OutputFormat=%s",
			request.MessageType, request.TestCategory.String(), request.RequestedOutputFormat.String())
	}

	// Process the request and generate response
	response := h.HandleConformanceRequest(&request)

	// Marshal the response
	return proto.Marshal(response)
}

// HandleConformanceRequest processes a conformance request and returns the appropriate response
func (h *ConformanceHandler) HandleConformanceRequest(request *ConformanceRequest) *ConformanceResponse {
	// Skip tests with TEXT_FORMAT test category
	if request.TestCategory == TestCategory_TEXT_FORMAT_TEST {
		return &ConformanceResponse{
			Result: &ConformanceResponse_Skipped{
				Skipped: "TEXT_FORMAT test category not supported",
			},
		}
	}

	// Handle different input and output formats
	switch payload := request.Payload.(type) {
	case *ConformanceRequest_ProtobufPayload:
		return h.HandleBinaryProtobuf(request, payload.ProtobufPayload)
	case *ConformanceRequest_JsonPayload:
		return h.HandleJsonInput(request, payload.JsonPayload)
	case *ConformanceRequest_JspbPayload:
		// JSPB format not supported
		return &ConformanceResponse{
			Result: &ConformanceResponse_Skipped{
				Skipped: "JSPB format not supported",
			},
		}
	case *ConformanceRequest_TextPayload:
		// Text format not supported
		return &ConformanceResponse{
			Result: &ConformanceResponse_Skipped{
				Skipped: "Text format not supported",
			},
		}
	default:
		return &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: "Unknown input payload type",
			},
		}
	}
}

// HandleJsonInput processes JSON input and generates appropriate output
func (h *ConformanceHandler) HandleJsonInput(request *ConformanceRequest, jsonInput string) *ConformanceResponse {
	if h.useLegacyConversion {
		return h.HandleJsonInputLegacy(request, jsonInput)
	}
	return h.HandleJsonInputWithProtoNumberJSON(request, jsonInput)
}

// HandleJsonInputLegacy processes JSON input using legacy direct conversion (original implementation)
func (h *ConformanceHandler) HandleJsonInputLegacy(request *ConformanceRequest, jsonInput string) *ConformanceResponse {
	if h.debug {
		log.Printf("Processing JSON input: %s", jsonInput)
	}

	// Get message type
	messageType := request.MessageType
	if messageType == "" {
		return &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: "Message type not specified for JSON conversion",
			},
		}
	}

	// Ensure message type starts with a dot for fully qualified name
	if !strings.HasPrefix(messageType, ".") {
		messageType = "." + messageType
	}

	// Convert JSON to binary protobuf using MySQL function
	var binaryData []byte
	query := "SELECT pb_json_to_message(conformance_test_messages_schema(), ?, ?)"
	err := h.db.QueryRow(query, messageType, jsonInput).Scan(&binaryData)
	if err != nil {
		// Check if this is a GROUP field error (unsupported feature)
		if isGroupFieldError(err) {
			return &ConformanceResponse{
				Result: &ConformanceResponse_Skipped{
					Skipped: "GROUP fields are not supported (deprecated Proto2 feature)",
				},
			}
		}
		// This is the JSON input parsing phase - all errors here should be parse errors
		return &ConformanceResponse{
			Result: &ConformanceResponse_ParseError{
				ParseError: fmt.Sprintf("Failed to parse JSON input: %v", err),
			},
		}
	}

	if h.debug {
		log.Printf("Converted JSON to binary protobuf: %d bytes", len(binaryData))
	}

	// Now handle the binary data normally - but we need a different approach since
	// HandleBinaryProtobuf converts to JSON first. For JSON input, we should handle
	// output formats directly.
	return h.HandleConvertedMessageLegacy(request, binaryData, messageType)
}

// HandleConvertedMessageLegacy processes a binary protobuf message and generates appropriate output
// This is used by the legacy conversion methods
func (h *ConformanceHandler) HandleConvertedMessageLegacy(request *ConformanceRequest, binaryData []byte, messageType string) *ConformanceResponse {
	// Generate output based on requested format
	switch request.RequestedOutputFormat {
	case WireFormat_PROTOBUF:
		// Return binary protobuf directly
		return &ConformanceResponse{
			Result: &ConformanceResponse_ProtobufPayload{
				ProtobufPayload: binaryData,
			},
		}

	case WireFormat_JSON:
		// Convert binary protobuf to JSON using MySQL function
		var jsonOutput string
		query := "SELECT pb_message_to_json(conformance_test_messages_schema(), ?, ?, NULL, NULL)"
		err := h.db.QueryRow(query, messageType, binaryData).Scan(&jsonOutput)
		if err != nil {
			// Check if this is a GROUP field error (unsupported feature)
			if isGroupFieldError(err) {
				return &ConformanceResponse{
					Result: &ConformanceResponse_Skipped{
						Skipped: "GROUP fields are not supported (deprecated Proto2 feature)",
					},
				}
			}
			// This is the JSON output serialization phase - all other errors are serialize errors
			return &ConformanceResponse{
				Result: &ConformanceResponse_SerializeError{
					SerializeError: fmt.Sprintf("Failed to serialize to JSON: %v", err),
				},
			}
		}

		return &ConformanceResponse{
			Result: &ConformanceResponse_JsonPayload{
				JsonPayload: jsonOutput,
			},
		}

	case WireFormat_TEXT_FORMAT:
		// Text format not supported
		return &ConformanceResponse{
			Result: &ConformanceResponse_Skipped{
				Skipped: "TEXT_FORMAT output not supported",
			},
		}

	default:
		return &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: "Unknown output format requested",
			},
		}
	}
}

// HandleBinaryProtobuf processes binary protobuf input and generates appropriate output
func (h *ConformanceHandler) HandleBinaryProtobuf(request *ConformanceRequest, binaryData []byte) *ConformanceResponse {
	if h.useLegacyConversion {
		return h.HandleBinaryProtobufLegacy(request, binaryData)
	}
	return h.HandleBinaryProtobufWithProtoNumberJSON(request, binaryData)
}

// HandleBinaryProtobufLegacy processes binary protobuf input using legacy direct conversion (original implementation)
func (h *ConformanceHandler) HandleBinaryProtobufLegacy(request *ConformanceRequest, binaryData []byte) *ConformanceResponse {
	if h.debug {
		log.Printf("Processing binary protobuf data: %d bytes", len(binaryData))
	}

	// Convert binary protobuf to JSON using MySQL function with schema awareness
	messageType := request.MessageType
	if messageType == "" {
		return &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: "Message type not specified for JSON conversion",
			},
		}
	}

	// Ensure message type starts with a dot for fully qualified name
	if !strings.HasPrefix(messageType, ".") {
		messageType = "." + messageType
	}

	var jsonOutput string
	query := "SELECT pb_message_to_json(conformance_test_messages_schema(), ?, ?, NULL, NULL)"
	err := h.db.QueryRow(query, messageType, binaryData).Scan(&jsonOutput)
	if err != nil {
		// Check if this is a GROUP field error (unsupported feature)
		if isGroupFieldError(err) {
			return &ConformanceResponse{
				Result: &ConformanceResponse_Skipped{
					Skipped: "GROUP fields are not supported (deprecated Proto2 feature)",
				},
			}
		}
		// This is the JSON output serialization phase - all other errors are serialize errors
		return &ConformanceResponse{
			Result: &ConformanceResponse_ParseError{
				ParseError: fmt.Sprintf("Failed to parse binary protobuf: %v", err),
			},
		}
	}

	if h.debug {
		log.Printf("Converted to JSON: %s", jsonOutput)
	}

	// Generate output based on requested format
	switch request.RequestedOutputFormat {
	case WireFormat_PROTOBUF:
		// Convert JSON back to binary protobuf
		var outputBinary []byte
		query = "SELECT pb_json_to_message(conformance_test_messages_schema(), ?, ?)"
		err = h.db.QueryRow(query, messageType, jsonOutput).Scan(&outputBinary)
		if err != nil {
			// This is the protobuf serialization phase (JSON was already validated during the first conversion)
			// Errors here are serialize errors since we're generating the final protobuf output
			return &ConformanceResponse{
				Result: &ConformanceResponse_SerializeError{
					SerializeError: fmt.Sprintf("Failed to convert JSON to binary: %v", err),
				},
			}
		}

		return &ConformanceResponse{
			Result: &ConformanceResponse_ProtobufPayload{
				ProtobufPayload: outputBinary,
			},
		}

	case WireFormat_JSON:
		// Return the JSON output directly (already converted above)
		return &ConformanceResponse{
			Result: &ConformanceResponse_JsonPayload{
				JsonPayload: jsonOutput,
			},
		}

	case WireFormat_TEXT_FORMAT:
		// Text format not supported
		return &ConformanceResponse{
			Result: &ConformanceResponse_Skipped{
				Skipped: "TEXT_FORMAT output not supported",
			},
		}

	default:
		return &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: "Unknown output format requested",
			},
		}
	}
}

// HandleBinaryProtobufWithProtoNumberJSON processes binary protobuf input using ProtoNumberJSON as intermediate format
func (h *ConformanceHandler) HandleBinaryProtobufWithProtoNumberJSON(request *ConformanceRequest, binaryData []byte) *ConformanceResponse {
	if h.debug {
		log.Printf("Processing binary protobuf data with ProtoNumberJSON: %d bytes", len(binaryData))
	}

	// Get message type
	messageType := request.MessageType
	if messageType == "" {
		return &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: "Message type not specified for binary conversion",
			},
		}
	}

	// Ensure message type starts with a dot for fully qualified name
	if !strings.HasPrefix(messageType, ".") {
		messageType = "." + messageType
	}

	// Convert binary protobuf to ProtoNumberJSON using MySQL function
	var protoNumberJSON string
	query := "SELECT _pb_message_to_number_json(conformance_test_messages_schema(), ?, ?, NULL)"
	err := h.db.QueryRow(query, messageType, binaryData).Scan(&protoNumberJSON)
	if err != nil {
		// Check if this is a GROUP field error (unsupported feature)
		if isGroupFieldError(err) {
			return &ConformanceResponse{
				Result: &ConformanceResponse_Skipped{
					Skipped: "GROUP fields are not supported (deprecated Proto2 feature)",
				},
			}
		}
		// This is the binary input parsing phase - all errors here should be parse errors
		return &ConformanceResponse{
			Result: &ConformanceResponse_ParseError{
				ParseError: fmt.Sprintf("Failed to parse binary protobuf: %v", err),
			},
		}
	}

	if h.debug {
		log.Printf("Converted binary to ProtoNumberJSON: %s", protoNumberJSON)
	}

	// Generate output from ProtoNumberJSON
	return h.generateOutputFromProtoNumberJSON(request, protoNumberJSON, messageType)
}

// HandleJsonInputWithProtoNumberJSON processes JSON input using ProtoNumberJSON as intermediate format
func (h *ConformanceHandler) HandleJsonInputWithProtoNumberJSON(request *ConformanceRequest, jsonInput string) *ConformanceResponse {
	if h.debug {
		log.Printf("Processing JSON input with ProtoNumberJSON: %s", jsonInput)
	}

	// Get message type
	messageType := request.MessageType
	if messageType == "" {
		return &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: "Message type not specified for JSON conversion",
			},
		}
	}

	// Ensure message type starts with a dot for fully qualified name
	if !strings.HasPrefix(messageType, ".") {
		messageType = "." + messageType
	}

	// Convert JSON to ProtoNumberJSON using MySQL function
	// Check if this is an "ignore unknown" test based on test category
	ignoreUnknownFields := request.TestCategory == TestCategory_JSON_IGNORE_UNKNOWN_PARSING_TEST
	ignoreUnknownEnums := request.TestCategory == TestCategory_JSON_IGNORE_UNKNOWN_PARSING_TEST

	// Create JsonUnmarshalOptions using generated Go struct
	options := &jsonoptionspb.JsonUnmarshalOptions{
		IgnoreUnknownFields: ignoreUnknownFields,
		IgnoreUnknownEnums:  ignoreUnknownEnums,
	}

	// Convert to ProtoNumberJSON format
	optionsJSON, err := protonumberjson.Marshal(options)
	if err != nil {
		return &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: fmt.Sprintf("Failed to marshal JsonUnmarshalOptions: %v", err),
			},
		}
	}

	var protoNumberJSON string
	query := "SELECT _pb_json_to_number_json(conformance_test_messages_schema(), ?, ?, ?)"
	err = h.db.QueryRow(query, messageType, jsonInput, string(optionsJSON)).Scan(&protoNumberJSON)
	if err != nil {
		// Check if this is a GROUP field error (unsupported feature)
		if isGroupFieldError(err) {
			return &ConformanceResponse{
				Result: &ConformanceResponse_Skipped{
					Skipped: "GROUP fields are not supported (deprecated Proto2 feature)",
				},
			}
		}
		// This is the JSON input parsing phase - all errors here should be parse errors
		return &ConformanceResponse{
			Result: &ConformanceResponse_ParseError{
				ParseError: fmt.Sprintf("Failed to parse JSON input: %v", err),
			},
		}
	}

	if h.debug {
		log.Printf("Converted JSON to ProtoNumberJSON: %s", protoNumberJSON)
	}

	// Generate output from ProtoNumberJSON
	return h.generateOutputFromProtoNumberJSON(request, protoNumberJSON, messageType)
}

// generateOutputFromProtoNumberJSON converts ProtoNumberJSON to the requested output format
func (h *ConformanceHandler) generateOutputFromProtoNumberJSON(request *ConformanceRequest, protoNumberJSON string, messageType string) *ConformanceResponse {
	switch request.RequestedOutputFormat {
	case WireFormat_PROTOBUF:
		// Convert ProtoNumberJSON to binary protobuf
		var binaryData []byte
		query := "SELECT _pb_number_json_to_message(conformance_test_messages_schema(), ?, ?, NULL)"
		err := h.db.QueryRow(query, messageType, protoNumberJSON).Scan(&binaryData)
		if err != nil {
			return &ConformanceResponse{
				Result: &ConformanceResponse_SerializeError{
					SerializeError: fmt.Sprintf("Failed to convert ProtoNumberJSON to binary: %v", err),
				},
			}
		}

		return &ConformanceResponse{
			Result: &ConformanceResponse_ProtobufPayload{
				ProtobufPayload: binaryData,
			},
		}

	case WireFormat_JSON:
		// Convert ProtoNumberJSON to ProtoJSON
		// Create JsonMarshalOptions with emit_default_values = false (Proto3 default behavior)
		marshalOptions := &jsonoptionspb.JsonMarshalOptions{
			EmitDefaultValues: false,
		}
		marshalOptionsJSON, err := protonumberjson.Marshal(marshalOptions)
		if err != nil {
			return &ConformanceResponse{
				Result: &ConformanceResponse_ParseError{
					ParseError: fmt.Sprintf("Failed to create marshal options: %v", err),
				},
			}
		}

		var jsonOutput string
		query := "SELECT _pb_number_json_to_json(conformance_test_messages_schema(), ?, ?, ?)"
		err = h.db.QueryRow(query, messageType, protoNumberJSON, string(marshalOptionsJSON)).Scan(&jsonOutput)
		if err != nil {
			// Check if this is a GROUP field error (unsupported feature)
			if isGroupFieldError(err) {
				return &ConformanceResponse{
					Result: &ConformanceResponse_Skipped{
						Skipped: "GROUP fields are not supported (deprecated Proto2 feature)",
					},
				}
			}
			return &ConformanceResponse{
				Result: &ConformanceResponse_SerializeError{
					SerializeError: fmt.Sprintf("Failed to convert ProtoNumberJSON to JSON: %v", err),
				},
			}
		}

		return &ConformanceResponse{
			Result: &ConformanceResponse_JsonPayload{
				JsonPayload: jsonOutput,
			},
		}

	case WireFormat_TEXT_FORMAT:
		// Text format not supported
		return &ConformanceResponse{
			Result: &ConformanceResponse_Skipped{
				Skipped: "TEXT_FORMAT output not supported",
			},
		}

	default:
		return &ConformanceResponse{
			Result: &ConformanceResponse_RuntimeError{
				RuntimeError: "Unknown output format requested",
			},
		}
	}
}

// isGroupFieldError checks if a MySQL error is due to unsupported GROUP fields
func isGroupFieldError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Check for GROUP field error patterns (field_type 10)
	return strings.Contains(errStr, "unsupported field_type 10") ||
		strings.Contains(errStr, "unsupported field_type `10`")
}
