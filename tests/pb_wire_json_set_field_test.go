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
	test := func(t *testing.T, protoFieldType string, setFunction string, generator func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value)) {
		t.Run(fmt.Sprintf("%s", protoFieldType), func(t *testing.T) {
			GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s value = 2; int32 b = 3;", protoFieldType), func(messageType protoreflect.MessageType) {
				valueField := messageType.Descriptor().Fields().ByName("value")

				seed := time.Now().UnixNano()
				t.Logf("Using seed = %d.", seed)
				rng := rand.New(rand.NewSource(seed))
				for i := 0; i < 100; i++ {
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

	test(t, "float", "pb_{kind}_set_float_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Float(rng, false, false)
		return newValue, protoreflect.ValueOfFloat32(newValue)
	})

	test(t, "double", "pb_{kind}_set_double_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Double(rng, false, false)
		return newValue, protoreflect.ValueOfFloat64(newValue)
	})

	test(t, "int32", "pb_{kind}_set_int32_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int32(rng)
		return newValue, protoreflect.ValueOfInt32(newValue)
	})

	test(t, "int64", "pb_{kind}_set_int64_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int64(rng)
		return newValue, protoreflect.ValueOfInt64(newValue)
	})

	test(t, "uint32", "pb_{kind}_set_uint32_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint32(rng)
		return newValue, protoreflect.ValueOfUint32(newValue)
	})

	test(t, "uint64", "pb_{kind}_set_uint64_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint64(rng)
		return newValue, protoreflect.ValueOfUint64(newValue)
	})

	test(t, "bool", "pb_{kind}_set_bool_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Bool(rng)
		return newValue, protoreflect.ValueOfBool(newValue)
	})

	test(t, "string", "pb_{kind}_set_string_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.String(rng, 5)
		return newValue, protoreflect.ValueOfString(newValue)
	})

	test(t, "bytes", "pb_{kind}_set_bytes_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Bytes(rng, 5)
		return newValue, protoreflect.ValueOfBytes(newValue)
	})

	test(t, "sint32", "pb_{kind}_set_sint32_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int32(rng)
		return newValue, protoreflect.ValueOfInt32(newValue)
	})

	test(t, "sint64", "pb_{kind}_set_sint64_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int64(rng)
		return newValue, protoreflect.ValueOfInt64(newValue)
	})

	test(t, "fixed32", "pb_{kind}_set_fixed32_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint32(rng)
		return newValue, protoreflect.ValueOfUint32(newValue)
	})

	test(t, "fixed64", "pb_{kind}_set_fixed64_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint64(rng)
		return newValue, protoreflect.ValueOfUint64(newValue)
	})

	test(t, "sfixed32", "pb_{kind}_set_sfixed32_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int32(rng)
		return newValue, protoreflect.ValueOfInt32(newValue)
	})

	test(t, "sfixed64", "pb_{kind}_set_sfixed64_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int64(rng)
		return newValue, protoreflect.ValueOfInt64(newValue)
	})

	test(t, "EnumType", "pb_{kind}_set_enum_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Enum(rng, fieldDescriptor.Enum())
		return newValue, protoreflect.ValueOfEnum(newValue)
	})

	test(t, "MessageType", "pb_{kind}_set_message_field", func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Message(rng, fieldDescriptor.Message(), nil)
		return newValue.Interface(), protoreflect.ValueOfMessage(newValue)
	})
}
