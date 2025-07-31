package gproto

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/moreproto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func EqualProto(expectedMessage proto.Message) *EqualProtoMatcher {
	return &EqualProtoMatcher{expectedMessage: expectedMessage}
}

type EqualProtoMatcher struct {
	expectedMessage proto.Message
	equalOrCloseFn  func(a, b float64, floatSize int) bool
}

func (m *EqualProtoMatcher) WithFloatEqualFn(equalOrCloseFn func(a, b float64, floatSize int) bool) *EqualProtoMatcher {
	return &EqualProtoMatcher{
		expectedMessage: m.expectedMessage,
		equalOrCloseFn:  equalOrCloseFn,
	}
}

func (m *EqualProtoMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, nil
	}

	var actualMessage proto.Message
	switch actualTyped := actual.(type) {
	case proto.Message:
		actualMessage = actualTyped
	case []byte:
		empty := m.expectedMessage.ProtoReflect().New()
		if err := proto.Unmarshal(actualTyped, empty.Interface()); err != nil {
			return false, fmt.Errorf("error unmarshalling bytes: %w", err)
		}
		actualMessage = empty.Interface()
	default:
		return false, fmt.Errorf("expected a proto.Message or []byte, got %T", actual)
	}

	if m.equalOrCloseFn == nil {
		return proto.Equal(actualMessage, m.expectedMessage), nil
	}

	return moreproto.EqualOrClose(actualMessage, m.expectedMessage, m.equalOrCloseFn), nil
}

func (m *EqualProtoMatcher) FailureMessage(actual interface{}) (message string) {
	return m.formatComparisonMessage(actual, "to match")
}

func (m *EqualProtoMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.formatComparisonMessage(actual, "not to match")
}

func (m *EqualProtoMatcher) formatComparisonMessage(actual interface{}, relation string) string {
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
			if jsonData, marshalErr := marshaler.Marshal(actualMessage.Interface()); marshalErr == nil {
				result += fmt.Sprintf("\nParsed as proto:\n%s", string(jsonData))
			}
		} else {
			result += fmt.Sprintf("\nFailed to parse as proto: %v", err)
		}
	}

	return result
}
