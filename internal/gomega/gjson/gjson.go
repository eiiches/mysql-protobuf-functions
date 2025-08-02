package gjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func EqualJson(expectedJson string) *EqualJsonMatcher {
	return &EqualJsonMatcher{expectedJson: expectedJson}
}

type EqualJsonMatcher struct {
	expectedJson string
	floatEqualFn func(a, b float64) bool
}

func (m *EqualJsonMatcher) WithFloatEqualFn(floatEqualFn func(a, b float64) bool) *EqualJsonMatcher {
	return &EqualJsonMatcher{
		expectedJson: m.expectedJson,
		floatEqualFn: floatEqualFn,
	}
}

func (m *EqualJsonMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, nil
	}

	var actualJson string
	switch actualTyped := actual.(type) {
	case string:
		actualJson = actualTyped
	case []byte:
		actualJson = string(actualTyped)
	default:
		return false, fmt.Errorf("expected a string or []byte containing JSON, got %T", actual)
	}

	// Parse both JSON strings
	var actualData interface{}
	var expectedData interface{}

	if err := json.Unmarshal([]byte(actualJson), &actualData); err != nil {
		return false, fmt.Errorf("error unmarshalling actual JSON: %w", err)
	}

	if err := json.Unmarshal([]byte(m.expectedJson), &expectedData); err != nil {
		return false, fmt.Errorf("error unmarshalling expected JSON: %w", err)
	}

	// Compare the parsed JSON data
	return m.deepEqual(actualData, expectedData), nil
}

func (m *EqualJsonMatcher) deepEqual(actual, expected interface{}) bool {
	// Handle nil cases
	if actual == nil && expected == nil {
		return true
	}
	if actual == nil || expected == nil {
		return false
	}

	actualValue := reflect.ValueOf(actual)
	expectedValue := reflect.ValueOf(expected)

	// If types are different, try to handle numeric type conversions
	if actualValue.Type() != expectedValue.Type() {
		return m.handleTypeMismatch(actual, expected)
	}

	//nolint:exhaustive
	switch actualValue.Kind() {
	case reflect.Float64:
		actualFloat := actualValue.Float()
		expectedFloat := expectedValue.Float()
		if m.floatEqualFn != nil {
			return m.floatEqualFn(actualFloat, expectedFloat)
		}
		return actualFloat == expectedFloat

	case reflect.Map:
		actualMap := actual.(map[string]interface{})
		expectedMap := expected.(map[string]interface{})

		if len(actualMap) != len(expectedMap) {
			return false
		}

		for key, actualVal := range actualMap {
			expectedVal, exists := expectedMap[key]
			if !exists {
				return false
			}
			if !m.deepEqual(actualVal, expectedVal) {
				return false
			}
		}
		return true

	case reflect.Slice:
		actualSlice := actual.([]interface{})
		expectedSlice := expected.([]interface{})

		if len(actualSlice) != len(expectedSlice) {
			return false
		}

		for i := range actualSlice {
			if !m.deepEqual(actualSlice[i], expectedSlice[i]) {
				return false
			}
		}
		return true

	default:
		return reflect.DeepEqual(actual, expected)
	}
}

func (m *EqualJsonMatcher) handleTypeMismatch(actual, expected interface{}) bool {
	// Try to convert both to float64 for numeric comparisons
	actualFloat, actualIsFloat := m.toFloat64(actual)
	expectedFloat, expectedIsFloat := m.toFloat64(expected)

	if actualIsFloat && expectedIsFloat {
		if m.floatEqualFn != nil {
			return m.floatEqualFn(actualFloat, expectedFloat)
		}
		return actualFloat == expectedFloat
	}

	return false
}

func (m *EqualJsonMatcher) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

func (m *EqualJsonMatcher) FailureMessage(actual interface{}) (message string) {
	return m.formatComparisonMessage(actual, "to match JSON of")
}

func (m *EqualJsonMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.formatComparisonMessage(actual, "not to match JSON of")
}

func (m *EqualJsonMatcher) formatComparisonMessage(actual interface{}, relation string) string {
	actualStr := m.formatJson(actual)
	expectedStr := m.formatJson(m.expectedJson)

	actualIndented := indentLines(actualStr)
	expectedIndented := indentLines(expectedStr)

	// Find the first mismatched key/path for more helpful error messages
	mismatchPath := m.findMismatchPath(actual)
	if mismatchPath != "" {
		return fmt.Sprintf("Expected\n%s\n%s\n%s\n\nfirst mismatched key: %s", actualIndented, relation, expectedIndented, mismatchPath)
	}

	return fmt.Sprintf("Expected\n%s\n%s\n%s", actualIndented, relation, expectedIndented)
}

func (m *EqualJsonMatcher) findMismatchPath(actual interface{}) string {
	var actualJson string
	switch actualTyped := actual.(type) {
	case string:
		actualJson = actualTyped
	case []byte:
		actualJson = string(actualTyped)
	default:
		return ""
	}

	var actualData interface{}
	var expectedData interface{}

	if err := json.Unmarshal([]byte(actualJson), &actualData); err != nil {
		return ""
	}
	if err := json.Unmarshal([]byte(m.expectedJson), &expectedData); err != nil {
		return ""
	}

	return m.findMismatchPathRecursive(actualData, expectedData, "")
}

func (m *EqualJsonMatcher) findMismatchPathRecursive(actual, expected interface{}, path string) string {
	if m.deepEqual(actual, expected) {
		return ""
	}

	if actual == nil || expected == nil {
		return path
	}

	actualValue := reflect.ValueOf(actual)

	//nolint:exhaustive
	switch actualValue.Kind() {
	case reflect.Map:
		actualMap, actualOk := actual.(map[string]interface{})
		expectedMap, expectedOk := expected.(map[string]interface{})

		if !actualOk || !expectedOk {
			return path
		}

		// Check for keys that exist in both maps
		for key := range actualMap {
			if expectedVal, exists := expectedMap[key]; exists {
				keyPath := path
				if keyPath == "" {
					keyPath = fmt.Sprintf(`"%s"`, key)
				} else {
					keyPath = fmt.Sprintf(`%s."%s"`, keyPath, key)
				}

				if result := m.findMismatchPathRecursive(actualMap[key], expectedVal, keyPath); result != "" {
					return result
				}
			}
		}

		// Return the first key that differs
		for key := range actualMap {
			keyPath := path
			if keyPath == "" {
				keyPath = fmt.Sprintf(`"%s"`, key)
			} else {
				keyPath = fmt.Sprintf(`%s."%s"`, keyPath, key)
			}
			return keyPath
		}

	case reflect.Slice:
		actualSlice, actualOk := actual.([]interface{})
		expectedSlice, expectedOk := expected.([]interface{})

		if !actualOk || !expectedOk {
			return path
		}

		minLen := len(actualSlice)
		if len(expectedSlice) < minLen {
			minLen = len(expectedSlice)
		}

		for i := 0; i < minLen; i++ {
			indexPath := fmt.Sprintf("%s[%d]", path, i)
			if result := m.findMismatchPathRecursive(actualSlice[i], expectedSlice[i], indexPath); result != "" {
				return result
			}
		}

		return path

	default:
		return path
	}

	return path
}

func (m *EqualJsonMatcher) formatJson(value interface{}) string {
	var jsonStr string
	switch v := value.(type) {
	case string:
		jsonStr = v
	case []byte:
		jsonStr = string(v)
	default:
		if data, err := json.Marshal(value); err == nil {
			jsonStr = string(data)
		} else {
			return fmt.Sprintf("%v", value)
		}
	}

	// Try to pretty-print the JSON
	var parsed interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err == nil {
		if formatted, err := json.MarshalIndent(parsed, "", "  "); err == nil {
			return string(formatted)
		}
	}

	return jsonStr
}

func indentLines(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder
	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString("    " + line)
	}
	return result.String()
}
