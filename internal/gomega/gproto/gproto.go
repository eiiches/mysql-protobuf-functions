package gproto

import (
	"fmt"

	"github.com/onsi/gomega/types"
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
	return fmt.Sprintf("Expected\n\t%#v\nto match\n\t%#v", actual, m)
}

func (m *equalProtoMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to match\n\t%#v", actual, m)
}
