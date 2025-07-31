package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"github.com/samber/lo"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRandomizedRemoveRepeatedFieldElement(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, removeFunction string, supportsPacked bool) {
		for _, usePacked := range lo.Ternary(supportsPacked, []string{"true", "false"}, []string{""}) {
			t.Run(fmt.Sprintf("%s/usePacked=%s", protoFieldType, usePacked), func(t *testing.T) {
				GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s values = 2%s; int32 b = 3;", protoFieldType, FormatPackedOption(usePacked)), func(messageType protoreflect.MessageType) {
					valuesField := messageType.Descriptor().Fields().ByName("values")

					seed := time.Now().UnixNano()
					t.Logf("Using seed = %d.", seed)
					rng := rand.New(rand.NewSource(seed))
					for i := 0; i < iterations; i++ {
						input := protorandom.Message(rng, messageType.Descriptor(), nil)
						list := input.Get(valuesField).List()

						// Only test removal if the repeated field has elements
						if list.Len() > 0 {
							// Test removing each valid index
							for indexToRemove := 0; indexToRemove < list.Len(); indexToRemove++ {
								// Create expected result by removing the element at indexToRemove
								expectedMessage := input.New()
								expectedMessage.Set(expectedMessage.Descriptor().Fields().ByName("a"), input.Get(input.Descriptor().Fields().ByName("a")))
								expectedMessage.Set(expectedMessage.Descriptor().Fields().ByName("b"), input.Get(input.Descriptor().Fields().ByName("b")))

								// Copy all elements except the one at indexToRemove
								expectedList := expectedMessage.Mutable(valuesField).List()
								for j := 0; j < list.Len(); j++ {
									if j != indexToRemove {
										expectedList.Append(list.Get(j))
									}
								}

								RunTestThatExpression(t, fmt.Sprintf("pb_wire_json_to_message(%s(pb_message_to_wire_json(?), 2, ?))", strings.ReplaceAll(removeFunction, "{kind}", "wire_json")), input.Interface(), indexToRemove).
									IsEqualToProto(expectedMessage.Interface())

								RunTestThatExpression(t, fmt.Sprintf("%s(?, 2, ?)", strings.ReplaceAll(removeFunction, "{kind}", "message")), input.Interface(), indexToRemove).
									IsEqualToProto(expectedMessage.Interface())
							}
						}
					}
				})
			})
		}
	}

	test(t, "repeated float", "pb_{kind}_remove_repeated_float_field_element", true)
	test(t, "repeated double", "pb_{kind}_remove_repeated_double_field_element", true)
	test(t, "repeated int32", "pb_{kind}_remove_repeated_int32_field_element", true)
	test(t, "repeated int64", "pb_{kind}_remove_repeated_int64_field_element", true)
	test(t, "repeated uint32", "pb_{kind}_remove_repeated_uint32_field_element", true)
	test(t, "repeated uint64", "pb_{kind}_remove_repeated_uint64_field_element", true)
	test(t, "repeated bool", "pb_{kind}_remove_repeated_bool_field_element", true)
	test(t, "repeated string", "pb_{kind}_remove_repeated_string_field_element", false)
	test(t, "repeated bytes", "pb_{kind}_remove_repeated_bytes_field_element", false)
	test(t, "repeated sint32", "pb_{kind}_remove_repeated_sint32_field_element", true)
	test(t, "repeated sint64", "pb_{kind}_remove_repeated_sint64_field_element", true)
	test(t, "repeated fixed32", "pb_{kind}_remove_repeated_fixed32_field_element", true)
	test(t, "repeated fixed64", "pb_{kind}_remove_repeated_fixed64_field_element", true)
	test(t, "repeated sfixed32", "pb_{kind}_remove_repeated_sfixed32_field_element", true)
	test(t, "repeated sfixed64", "pb_{kind}_remove_repeated_sfixed64_field_element", true)
	test(t, "repeated EnumType", "pb_{kind}_remove_repeated_enum_field_element", true)
	test(t, "repeated MessageType", "pb_{kind}_remove_repeated_message_field_element", false)
}
