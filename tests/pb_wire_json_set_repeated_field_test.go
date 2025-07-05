package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/morefloat"
	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRandomizedSetRepeatedField(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, setFunction string, supportsPacked bool, generator ValueGenerator) {
		for _, usePacked := range lo.Ternary(supportsPacked, []string{"true", "false"}, []string{""}) {
			t.Run(fmt.Sprintf("%s/usePacked=%v", protoFieldType, usePacked), func(t *testing.T) {
				GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s value = 2%s; int32 b = 3;", protoFieldType, FormatPackedOption(usePacked)), func(messageType protoreflect.MessageType) {
					valueField := messageType.Descriptor().Fields().ByName("value")

					seed := time.Now().UnixNano()
					t.Logf("Using seed = %d.", seed)
					rng := rand.New(rand.NewSource(seed))
					for i := 0; i < iterations; i++ {
						input := protorandom.Message(rng, messageType.Descriptor(), nil)

						// Generate 0-4 new values to set
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
						// Clear existing values and set new ones (this is the key difference from add_all)
						expectedList.Truncate(0)
						for _, newProtoreflectValue := range newProtoreflectValues {
							expectedList.Append(newProtoreflectValue)
						}

						// Convert the array to JSON string for SQL
						jsonString := convertValuesToJson(t, rng, newRawValues)

						arguments := lo.Ternary(supportsPacked, []string{", TRUE", ", FALSE"}, []string{""})

						for _, argument := range arguments {
							// Floating-point numbers in JSON is inaccurate. Small difference in mantissa has to be ignored for this test to pass.
							// See https://bugs.mysql.com/bug.php?id=118497

							RunTestThatExpression(t, fmt.Sprintf("pb_wire_json_to_message(%s(pb_message_to_wire_json(?), 2, ?%s))", strings.ReplaceAll(setFunction, "{kind}", "wire_json"), argument), input.Interface(), jsonString).
								IsEqualOrCloseToProto(expected.Interface(), morefloat.WithinMantissaThreshold(0, 2))

							RunTestThatExpression(t, fmt.Sprintf("%s(?, 2, ?%s)", strings.ReplaceAll(setFunction, "{kind}", "message"), argument), input.Interface(), jsonString).
								IsEqualOrCloseToProto(expected.Interface(), morefloat.WithinMantissaThreshold(0, 2))
						}
					}
				})
			})
		}
	}

	test(t, "repeated float", "pb_{kind}_set_repeated_float_field", true, RandomFloatGenerator)
	test(t, "repeated double", "pb_{kind}_set_repeated_double_field", true, RandomDoubleGenerator)
	test(t, "repeated int32", "pb_{kind}_set_repeated_int32_field", true, RandomInt32Generator)
	test(t, "repeated int64", "pb_{kind}_set_repeated_int64_field", true, RandomInt64Generator)
	test(t, "repeated uint32", "pb_{kind}_set_repeated_uint32_field", true, RandomUint32Generator)
	test(t, "repeated uint64", "pb_{kind}_set_repeated_uint64_field", true, RandomUint64Generator)
	test(t, "repeated bool", "pb_{kind}_set_repeated_bool_field", true, RandomBoolGenerator)
	test(t, "repeated string", "pb_{kind}_set_repeated_string_field", false, RandomStringGenerator)
	test(t, "repeated bytes", "pb_{kind}_set_repeated_bytes_field", false, RandomBytesGenerator)
	test(t, "repeated sint32", "pb_{kind}_set_repeated_sint32_field", true, RandomInt32Generator)
	test(t, "repeated sint64", "pb_{kind}_set_repeated_sint64_field", true, RandomInt64Generator)
	test(t, "repeated fixed32", "pb_{kind}_set_repeated_fixed32_field", true, RandomUint32Generator)
	test(t, "repeated fixed64", "pb_{kind}_set_repeated_fixed64_field", true, RandomUint64Generator)
	test(t, "repeated sfixed32", "pb_{kind}_set_repeated_sfixed32_field", true, RandomInt32Generator)
	test(t, "repeated sfixed64", "pb_{kind}_set_repeated_sfixed64_field", true, RandomInt64Generator)
	test(t, "repeated EnumType", "pb_{kind}_set_repeated_enum_field", true, RandomEnumGenerator)
	test(t, "repeated MessageType", "pb_{kind}_set_repeated_message_field", false, RandomMessageGenerator)
}
