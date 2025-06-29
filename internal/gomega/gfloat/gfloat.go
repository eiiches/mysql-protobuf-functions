package gfloat

import (
	"fmt"
	"math"

	"github.com/onsi/gomega/types"
)

func BeNegativeZero() types.GomegaMatcher {
	return &negativeZeroMatcher{}
}

func BePositiveZero() types.GomegaMatcher {
	return &positiveZeroMatcher{}
}

type negativeZeroMatcher struct{}

func (m *negativeZeroMatcher) Match(actual interface{}) (success bool, err error) {
	switch v := actual.(type) {
	case float64:
		return v == 0.0 && math.Signbit(v), nil
	case float32:
		return v == 0.0 && math.Signbit(float64(v)), nil
	default:
		return false, fmt.Errorf("BeNegativeZero matcher expects a float32 or float64, got %T", actual)
	}
}

func (m *negativeZeroMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nto be negative zero", actual)
}

func (m *negativeZeroMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nnot to be negative zero", actual)
}

type positiveZeroMatcher struct{}

func (m *positiveZeroMatcher) Match(actual interface{}) (success bool, err error) {
	switch v := actual.(type) {
	case float64:
		return v == 0.0 && !math.Signbit(v), nil
	case float32:
		return v == 0.0 && !math.Signbit(float64(v)), nil
	default:
		return false, fmt.Errorf("BePositiveZero matcher expects a float32 or float64, got %T", actual)
	}
}

func (m *positiveZeroMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nto be positive zero", actual)
}

func (m *positiveZeroMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nnot to be positive zero", actual)
}
