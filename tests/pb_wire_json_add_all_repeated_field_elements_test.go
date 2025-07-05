package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/morefloat"

	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func convertValuesToJson(t *testing.T, rng *rand.Rand, values []interface{}) string {
	var jsonParts []string
	jsonParts = append(jsonParts, "[")

	for i, value := range values {
		if i > 0 {
			jsonParts = append(jsonParts, ",")
		}

		switch v := value.(type) {
		case []byte:
			// Base64 encode bytes as JSON string
			encoded := base64.StdEncoding.EncodeToString(v)
			escapedBytes, _ := json.Marshal(encoded)
			jsonParts = append(jsonParts, string(escapedBytes))
		case proto.Message:
			// Marshal message and base64 encode as JSON string
			messageBytes, err := proto.Marshal(v)
			if err != nil {
				t.Fatalf("Failed to marshal message: %v", err)
			}
			encoded := base64.StdEncoding.EncodeToString(messageBytes)
			escapedBytes, _ := json.Marshal(encoded)
			jsonParts = append(jsonParts, string(escapedBytes))
		case float32:
			// Use strconv.FormatFloat with maximum precision for exact round-trip fidelity
			jsonParts = append(jsonParts, strconv.FormatFloat(float64(v), 'g', -1, 32))
		case float64:
			// Use strconv.FormatFloat with maximum precision for exact round-trip fidelity
			// -1 precision means use the minimum number of digits needed for exact representation
			jsonParts = append(jsonParts, strconv.FormatFloat(v, 'g', -1, 64))
		case int, int32, int64, uint, uint32, uint64:
			if rng.Intn(2) == 0 {
				// Integer as JSON string
				jsonParts = append(jsonParts, fmt.Sprintf("\"%d\"", v))
			} else {
				// Integer as JSON number
				jsonParts = append(jsonParts, fmt.Sprintf("%d", v))
			}
		default:
			// Use standard JSON marshaling for other types
			valueBytes, err := json.Marshal(value)
			if err != nil {
				t.Fatalf("Failed to marshal value: %v", err)
			}
			jsonParts = append(jsonParts, string(valueBytes))
		}
	}

	jsonParts = append(jsonParts, "]")
	return strings.Join(jsonParts, "")
}

func TestRandomizedAddAllRepeatedFieldElements(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, addFunction string, supportsPacked bool, generator ValueGenerator) {
		for _, usePacked := range lo.Ternary(supportsPacked, []string{"true", "false"}, []string{""}) {
			t.Run(fmt.Sprintf("%s/usePacked=%v", protoFieldType, usePacked), func(t *testing.T) {
				GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s value = 2%s; int32 b = 3;", protoFieldType, FormatPackedOption(usePacked)), func(messageType protoreflect.MessageType) {
					valueField := messageType.Descriptor().Fields().ByName("value")

					seed := time.Now().UnixNano()
					t.Logf("Using seed = %d.", seed)
					rng := rand.New(rand.NewSource(seed))
					for i := 0; i < iterations; i++ {
						input := protorandom.Message(rng, messageType.Descriptor(), nil)

						// Generate 0-4 new values to add
						newValueCount := rng.Intn(5)
						var newRawValues []interface{}
						var newProtoreflectValues []protoreflect.Value
						for j := 0; j < newValueCount; j++ {
							newValue, newProtoreflectValue := generator(rng, valueField)
							newRawValues = append(newRawValues, newValue)
							newProtoreflectValues = append(newProtoreflectValues, newProtoreflectValue)
						}

						expected := proto.Clone(input.Interface()).ProtoReflect()
						expectedList := expected.Mutable(valueField).List()
						for _, newProtoreflectValue := range newProtoreflectValues {
							expectedList.Append(newProtoreflectValue)
						}

						// Convert the array to JSON string for SQL
						jsonString := convertValuesToJson(t, rng, newRawValues)

						arguments := lo.Ternary(supportsPacked, []string{", TRUE", ", FALSE"}, []string{""})

						for _, argument := range arguments {
							// Floating-point numbers in JSON is inaccurate. Small difference in mantissa has to be ignored for this test to pass.
							// See https://bugs.mysql.com/bug.php?id=118497

							RunTestThatExpression(t, fmt.Sprintf("pb_wire_json_to_message(%s(pb_message_to_wire_json(?), 2, ?%s))", strings.ReplaceAll(addFunction, "{kind}", "wire_json"), argument), input.Interface(), jsonString).
								IsEqualOrCloseToProto(expected.Interface(), morefloat.WithinMantissaThreshold(0, 2))

							RunTestThatExpression(t, fmt.Sprintf("%s(?, 2, ?%s)", strings.ReplaceAll(addFunction, "{kind}", "message"), argument), input.Interface(), jsonString).
								IsEqualOrCloseToProto(expected.Interface(), morefloat.WithinMantissaThreshold(0, 2))
						}
					}
				})
			})
		}
	}

	test(t, "repeated float", "pb_{kind}_add_all_repeated_float_field_elements", true, RandomFloatGenerator)
	test(t, "repeated double", "pb_{kind}_add_all_repeated_double_field_elements", true, RandomDoubleGenerator)
	test(t, "repeated int32", "pb_{kind}_add_all_repeated_int32_field_elements", true, RandomInt32Generator)
	test(t, "repeated int64", "pb_{kind}_add_all_repeated_int64_field_elements", true, RandomInt64Generator)
	test(t, "repeated uint32", "pb_{kind}_add_all_repeated_uint32_field_elements", true, RandomUint32Generator)
	test(t, "repeated uint64", "pb_{kind}_add_all_repeated_uint64_field_elements", true, RandomUint64Generator)
	test(t, "repeated bool", "pb_{kind}_add_all_repeated_bool_field_elements", true, RandomBoolGenerator)
	test(t, "repeated string", "pb_{kind}_add_all_repeated_string_field_elements", false, RandomStringGenerator)
	test(t, "repeated bytes", "pb_{kind}_add_all_repeated_bytes_field_elements", false, RandomBytesGenerator)
	test(t, "repeated sint32", "pb_{kind}_add_all_repeated_sint32_field_elements", true, RandomInt32Generator)
	test(t, "repeated sint64", "pb_{kind}_add_all_repeated_sint64_field_elements", true, RandomInt64Generator)
	test(t, "repeated fixed32", "pb_{kind}_add_all_repeated_fixed32_field_elements", true, RandomUint32Generator)
	test(t, "repeated fixed64", "pb_{kind}_add_all_repeated_fixed64_field_elements", true, RandomUint64Generator)
	test(t, "repeated sfixed32", "pb_{kind}_add_all_repeated_sfixed32_field_elements", true, RandomInt32Generator)
	test(t, "repeated sfixed64", "pb_{kind}_add_all_repeated_sfixed64_field_elements", true, RandomInt64Generator)
	test(t, "repeated EnumType", "pb_{kind}_add_all_repeated_enum_field_elements", true, RandomEnumGenerator)
	test(t, "repeated MessageType", "pb_{kind}_add_all_repeated_message_field_elements", false, RandomMessageGenerator)
}
