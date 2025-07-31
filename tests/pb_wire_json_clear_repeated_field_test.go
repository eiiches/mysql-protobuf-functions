package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRandomizedWireJsonClearRepeatedField(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, clearFunction string) {
		t.Run(protoFieldType, func(t *testing.T) {
			GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s value = 2; int32 b = 3;", protoFieldType), func(messageType protoreflect.MessageType) {
				valueField := messageType.Descriptor().Fields().ByName("value")

				seed := time.Now().UnixNano()
				t.Logf("Using seed = %d.", seed)
				rng := rand.New(rand.NewSource(seed))
				for i := 0; i < iterations; i++ {
					input := protorandom.Message(rng, messageType.Descriptor(), nil)

					// Create expected result: clone input and clear the repeated field
					expected := proto.Clone(input.Interface()).ProtoReflect()
					expected.Clear(valueField)

					// Test clearing the repeated field from wire JSON
					RunTestThatExpression(t, fmt.Sprintf("pb_wire_json_to_message(%s(pb_message_to_wire_json(?), 2))", strings.ReplaceAll(clearFunction, "{kind}", "wire_json")), input.Interface()).
						IsEqualToProto(expected.Interface())

					// Test clearing the repeated field from message
					RunTestThatExpression(t, fmt.Sprintf("%s(?, 2)", strings.ReplaceAll(clearFunction, "{kind}", "message")), input.Interface()).
						IsEqualToProto(expected.Interface())
				}
			})
		})
	}

	test(t, "repeated float", "pb_{kind}_clear_repeated_float_field")
	test(t, "repeated double", "pb_{kind}_clear_repeated_double_field")
	test(t, "repeated int32", "pb_{kind}_clear_repeated_int32_field")
	test(t, "repeated int64", "pb_{kind}_clear_repeated_int64_field")
	test(t, "repeated uint32", "pb_{kind}_clear_repeated_uint32_field")
	test(t, "repeated uint64", "pb_{kind}_clear_repeated_uint64_field")
	test(t, "repeated bool", "pb_{kind}_clear_repeated_bool_field")
	test(t, "repeated string", "pb_{kind}_clear_repeated_string_field")
	test(t, "repeated bytes", "pb_{kind}_clear_repeated_bytes_field")
	test(t, "repeated sint32", "pb_{kind}_clear_repeated_sint32_field")
	test(t, "repeated sint64", "pb_{kind}_clear_repeated_sint64_field")
	test(t, "repeated fixed32", "pb_{kind}_clear_repeated_fixed32_field")
	test(t, "repeated fixed64", "pb_{kind}_clear_repeated_fixed64_field")
	test(t, "repeated sfixed32", "pb_{kind}_clear_repeated_sfixed32_field")
	test(t, "repeated sfixed64", "pb_{kind}_clear_repeated_sfixed64_field")
	test(t, "repeated EnumType", "pb_{kind}_clear_repeated_enum_field")
	test(t, "repeated MessageType", "pb_{kind}_clear_repeated_message_field")
}
