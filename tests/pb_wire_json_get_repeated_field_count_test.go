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

func TestRandomizedWireJsonGetRepeatedFieldCount(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, getCountFunction string, supportsPacked bool) {
		for _, usePacked := range lo.Ternary(supportsPacked, []string{"true", "false"}, []string{""}) {
			t.Run(fmt.Sprintf("%s/usePacked=%v", protoFieldType, usePacked), func(t *testing.T) {
				GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s values = 2%s; int32 b = 3;", protoFieldType, FormatPackedOption(usePacked)), func(messageType protoreflect.MessageType) {
					valuesField := messageType.Descriptor().Fields().ByName("values")

					seed := time.Now().UnixNano()
					t.Logf("Using seed = %d.", seed)
					rng := rand.New(rand.NewSource(seed))
					for i := 0; i < 100; i++ {
						input := protorandom.Message(rng, messageType.Descriptor(), nil)
						expectedCount := int32(input.Get(valuesField).List().Len())

						RunTestThatExpression(t, fmt.Sprintf("%s(pb_message_to_wire_json(?), 2)", strings.ReplaceAll(getCountFunction, "{kind}", "wire_json")), input.Interface()).
							IsEqualTo(expectedCount)

						RunTestThatExpression(t, fmt.Sprintf("%s(?, 2)", strings.ReplaceAll(getCountFunction, "{kind}", "message")), input.Interface()).
							IsEqualTo(expectedCount)
					}
				})
			})
		}
	}

	test(t, "repeated float", "pb_{kind}_get_repeated_float_field_count", true)
	test(t, "repeated double", "pb_{kind}_get_repeated_double_field_count", true)
	test(t, "repeated int32", "pb_{kind}_get_repeated_int32_field_count", true)
	test(t, "repeated int64", "pb_{kind}_get_repeated_int64_field_count", true)
	test(t, "repeated uint32", "pb_{kind}_get_repeated_uint32_field_count", true)
	test(t, "repeated uint64", "pb_{kind}_get_repeated_uint64_field_count", true)
	test(t, "repeated bool", "pb_{kind}_get_repeated_bool_field_count", true)
	test(t, "repeated string", "pb_{kind}_get_repeated_string_field_count", false)
	test(t, "repeated bytes", "pb_{kind}_get_repeated_bytes_field_count", false)
	test(t, "repeated sint32", "pb_{kind}_get_repeated_sint32_field_count", true)
	test(t, "repeated sint64", "pb_{kind}_get_repeated_sint64_field_count", true)
	test(t, "repeated fixed32", "pb_{kind}_get_repeated_fixed32_field_count", true)
	test(t, "repeated fixed64", "pb_{kind}_get_repeated_fixed64_field_count", true)
	test(t, "repeated sfixed32", "pb_{kind}_get_repeated_sfixed32_field_count", true)
	test(t, "repeated sfixed64", "pb_{kind}_get_repeated_sfixed64_field_count", true)
	test(t, "repeated EnumType", "pb_{kind}_get_repeated_enum_field_count", true)
	test(t, "repeated MessageType", "pb_{kind}_get_repeated_message_field_count", false)
}
