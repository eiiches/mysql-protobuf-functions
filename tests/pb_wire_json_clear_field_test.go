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

func TestRandomizedWireJsonClearField(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, clearFunction string) {
		t.Run(protoFieldType, func(t *testing.T) {
			GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s value = 2; int32 b = 3;", protoFieldType), func(messageType protoreflect.MessageType) {
				valueField := messageType.Descriptor().Fields().ByName("value")

				seed := time.Now().UnixNano()
				t.Logf("Using seed = %d.", seed)
				rng := rand.New(rand.NewSource(seed))
				for i := 0; i < iterations; i++ {
					input := protorandom.Message(rng, messageType.Descriptor(), nil)

					// Create expected result: clone input and clear the field
					expected := proto.Clone(input.Interface()).ProtoReflect()
					expected.Clear(valueField)

					// Test clearing the field from wire JSON
					RunTestThatExpression(t, fmt.Sprintf("pb_wire_json_to_message(%s(pb_message_to_wire_json(?), 2))", strings.ReplaceAll(clearFunction, "{kind}", "wire_json")), input.Interface()).
						IsEqualToProto(expected.Interface())

					// Test clearing the field from message
					RunTestThatExpression(t, fmt.Sprintf("%s(?, 2)", strings.ReplaceAll(clearFunction, "{kind}", "message")), input.Interface()).
						IsEqualToProto(expected.Interface())
				}
			})
		})
	}

	test(t, "float", "pb_{kind}_clear_float_field")
	test(t, "double", "pb_{kind}_clear_double_field")
	test(t, "int32", "pb_{kind}_clear_int32_field")
	test(t, "int64", "pb_{kind}_clear_int64_field")
	test(t, "uint32", "pb_{kind}_clear_uint32_field")
	test(t, "uint64", "pb_{kind}_clear_uint64_field")
	test(t, "bool", "pb_{kind}_clear_bool_field")
	test(t, "string", "pb_{kind}_clear_string_field")
	test(t, "bytes", "pb_{kind}_clear_bytes_field")
	test(t, "sint32", "pb_{kind}_clear_sint32_field")
	test(t, "sint64", "pb_{kind}_clear_sint64_field")
	test(t, "fixed32", "pb_{kind}_clear_fixed32_field")
	test(t, "fixed64", "pb_{kind}_clear_fixed64_field")
	test(t, "sfixed32", "pb_{kind}_clear_sfixed32_field")
	test(t, "sfixed64", "pb_{kind}_clear_sfixed64_field")
	test(t, "EnumType", "pb_{kind}_clear_enum_field")
	test(t, "MessageType", "pb_{kind}_clear_message_field")
}
