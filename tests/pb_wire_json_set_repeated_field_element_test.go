package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRandomizedSetRepeatedFieldElement(t *testing.T) {
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

						// Test each element in the repeated field
						for index := 0; index < input.Get(valueField).List().Len(); index++ {
							newValue, newProtoreflectValue := generator(rng, valueField)

							expected := proto.Clone(input.Interface()).ProtoReflect()
							expected.Mutable(valueField).List().Set(index, newProtoreflectValue)

							RunTestThatExpression(t, fmt.Sprintf("pb_wire_json_to_message(%s(pb_message_to_wire_json(?), 2, ?, ?))", strings.ReplaceAll(setFunction, "{kind}", "wire_json")), input.Interface(), index, newValue).
								IsEqualToProto(expected.Interface())

							RunTestThatExpression(t, fmt.Sprintf("%s(?, 2, ?, ?)", strings.ReplaceAll(setFunction, "{kind}", "message")), input.Interface(), index, newValue).
								IsEqualToProto(expected.Interface())
						}
					}
				})
			})
		}
	}

	test(t, "repeated float", "pb_{kind}_set_repeated_float_field_element", true, RandomFloatGenerator)
	test(t, "repeated double", "pb_{kind}_set_repeated_double_field_element", true, RandomDoubleGenerator)
	test(t, "repeated int32", "pb_{kind}_set_repeated_int32_field_element", true, RandomInt32Generator)
	test(t, "repeated int64", "pb_{kind}_set_repeated_int64_field_element", true, RandomInt64Generator)
	test(t, "repeated uint32", "pb_{kind}_set_repeated_uint32_field_element", true, RandomUint32Generator)
	test(t, "repeated uint64", "pb_{kind}_set_repeated_uint64_field_element", true, RandomUint64Generator)
	test(t, "repeated bool", "pb_{kind}_set_repeated_bool_field_element", true, RandomBoolGenerator)
	test(t, "repeated string", "pb_{kind}_set_repeated_string_field_element", false, RandomStringGenerator)
	test(t, "repeated sint32", "pb_{kind}_set_repeated_sint32_field_element", true, RandomInt32Generator)
	test(t, "repeated sint64", "pb_{kind}_set_repeated_sint64_field_element", true, RandomInt64Generator)
	test(t, "repeated fixed32", "pb_{kind}_set_repeated_fixed32_field_element", true, RandomUint32Generator)
	test(t, "repeated fixed64", "pb_{kind}_set_repeated_fixed64_field_element", true, RandomUint64Generator)
	test(t, "repeated sfixed32", "pb_{kind}_set_repeated_sfixed32_field_element", true, RandomInt32Generator)
	test(t, "repeated sfixed64", "pb_{kind}_set_repeated_sfixed64_field_element", true, RandomInt64Generator)
	test(t, "repeated EnumType", "pb_{kind}_set_repeated_enum_field_element", true, RandomEnumGenerator)
	test(t, "repeated MessageType", "pb_{kind}_set_repeated_message_field_element", false, RandomMessageGenerator)
}
