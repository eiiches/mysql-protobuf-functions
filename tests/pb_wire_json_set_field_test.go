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

func TestRandomizedWireJsonSetField(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, setFunction string, generator ValueGenerator) {
		t.Run(protoFieldType, func(t *testing.T) {
			GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s value = 2; int32 b = 3;", protoFieldType), func(messageType protoreflect.MessageType) {
				valueField := messageType.Descriptor().Fields().ByName("value")

				seed := time.Now().UnixNano()
				t.Logf("Using seed = %d.", seed)
				rng := rand.New(rand.NewSource(seed))
				for i := 0; i < iterations; i++ {
					input := protorandom.Message(rng, messageType.Descriptor(), nil)
					newValue, newProtoreflectValue := generator(rng, valueField)

					expected := proto.Clone(input.Interface()).ProtoReflect()
					expected.Set(valueField, newProtoreflectValue)

					RunTestThatExpression(t, fmt.Sprintf("pb_wire_json_to_message(%s(pb_message_to_wire_json(?), 2, ?))", strings.ReplaceAll(setFunction, "{kind}", "wire_json")), input.Interface(), newValue).
						IsEqualToProto(expected.Interface())

					RunTestThatExpression(t, fmt.Sprintf("%s(?, 2, ?)", strings.ReplaceAll(setFunction, "{kind}", "message")), input.Interface(), newValue).
						IsEqualToProto(expected.Interface())
				}
			})
		})
	}

	test(t, "float", "pb_{kind}_set_float_field", RandomFloatGenerator)
	test(t, "double", "pb_{kind}_set_double_field", RandomDoubleGenerator)
	test(t, "int32", "pb_{kind}_set_int32_field", RandomInt32Generator)
	test(t, "int64", "pb_{kind}_set_int64_field", RandomInt64Generator)
	test(t, "uint32", "pb_{kind}_set_uint32_field", RandomUint32Generator)
	test(t, "uint64", "pb_{kind}_set_uint64_field", RandomUint64Generator)
	test(t, "bool", "pb_{kind}_set_bool_field", RandomBoolGenerator)
	test(t, "string", "pb_{kind}_set_string_field", RandomStringGenerator)
	test(t, "bytes", "pb_{kind}_set_bytes_field", RandomBytesGenerator)
	test(t, "sint32", "pb_{kind}_set_sint32_field", RandomInt32Generator)
	test(t, "sint64", "pb_{kind}_set_sint64_field", RandomInt64Generator)
	test(t, "fixed32", "pb_{kind}_set_fixed32_field", RandomUint32Generator)
	test(t, "fixed64", "pb_{kind}_set_fixed64_field", RandomUint64Generator)
	test(t, "sfixed32", "pb_{kind}_set_sfixed32_field", RandomInt32Generator)
	test(t, "sfixed64", "pb_{kind}_set_sfixed64_field", RandomInt64Generator)
	test(t, "EnumType", "pb_{kind}_set_enum_field", RandomEnumGenerator)
	test(t, "MessageType", "pb_{kind}_set_message_field", RandomMessageGenerator)
}
