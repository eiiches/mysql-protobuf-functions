package gproto

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func EqualProto(expectedMessage proto.Message) types.GomegaMatcher {
	return &equalProtoMatcher{expectedMessage: expectedMessage}
}

type equalProtoMatcher struct {
	expectedMessage proto.Message
}

func (m *equalProtoMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, nil
	}
	switch actualTyped := actual.(type) {
	case proto.Message:
		return proto.Equal(actualTyped, m.expectedMessage), nil
	case []byte:
		empty := m.expectedMessage.ProtoReflect().New()
		if err := proto.Unmarshal(actualTyped, empty.Interface()); err != nil {
			return false, fmt.Errorf("error unmarshalling bytes: %w", err)
		}
		return proto.Equal(empty.Interface().(proto.Message), m.expectedMessage), nil
	default:
		return false, fmt.Errorf("expected a proto.Message or []byte, got %T", actual)
	}
}

func (m *equalProtoMatcher) FailureMessage(actual interface{}) (message string) {
	return m.formatComparisonMessage(actual, "to match")
}

func (m *equalProtoMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.formatComparisonMessage(actual, "not to match")
}

func (m *equalProtoMatcher) formatComparisonMessage(actual interface{}, relation string) string {
	actualStr := formatMessageWithContext(actual, m.expectedMessage)
	expectedStr := formatMessage(m.expectedMessage)

	actualIndented := indentLines(actualStr)
	expectedIndented := indentLines(expectedStr)

	return fmt.Sprintf("Expected\n%s\n%s\n%s", actualIndented, relation, expectedIndented)
}

func indentLines(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder
	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString("\t" + line)
	}
	return result.String()
}

func formatMessage(msg interface{}) string {
	return formatMessageWithContext(msg, nil)
}

func formatMessageWithContext(msg interface{}, expectedMessage proto.Message) string {
	switch v := msg.(type) {
	case proto.Message:
		return formatProtoMessage(v)
	case []byte:
		return formatByteSlice(v, expectedMessage)
	default:
		return fmt.Sprintf("%#v", v)
	}
}

func formatProtoMessage(msg proto.Message) string {
	marshaler := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}
	
	var result string
	if data, err := marshaler.Marshal(msg); err == nil {
		result = string(data)
	} else {
		result = fmt.Sprintf("%#v", msg)
	}

	// Also show the hex encoding of the proto message
	if protoBytes, err := proto.Marshal(msg); err == nil {
		if len(protoBytes) == 0 {
			result += "\nHex: (empty)"
		} else {
			hexStr := hex.EncodeToString(protoBytes)
			result += fmt.Sprintf("\nHex: %s", hexStr)
		}
	}

	return result
}

func formatByteSlice(data []byte, expectedMessage proto.Message) string {
	if len(data) == 0 {
		return "[]byte{} (empty)"
	}
	
	hexStr := hex.EncodeToString(data)
	result := fmt.Sprintf("[]byte{%d bytes}\nHex: %s", len(data), hexStr)

	// Try to unmarshal as the expected message type if provided
	if expectedMessage != nil {
		actualMessage := expectedMessage.ProtoReflect().New()
		if err := proto.Unmarshal(data, actualMessage.Interface()); err == nil {
			marshaler := protojson.MarshalOptions{
				Multiline: true,
				Indent:    "  ",
			}
			if jsonData, err := marshaler.Marshal(actualMessage.Interface()); err == nil {
				result += fmt.Sprintf("\nParsed as proto:\n%s", string(jsonData))
			}
		} else {
			result += fmt.Sprintf("\nFailed to parse as proto: %v", err)
		}
	}
	
	return result
}
