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

func TestRandomizedAddRepeatedFieldElement(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, addFunction string, supportsPacked bool, generator func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value)) {
		for _, usePacked := range lo.Ternary(supportsPacked, []string{"true", "false"}, []string{""}) {
			t.Run(fmt.Sprintf("%s/usePacked=%v", protoFieldType, usePacked), func(t *testing.T) {
				GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s value = 2%s; int32 b = 3;", protoFieldType, FormatPackedOption(usePacked)), func(messageType protoreflect.MessageType) {
					valueField := messageType.Descriptor().Fields().ByName("value")

					seed := time.Now().UnixNano()
					t.Logf("Using seed = %d.", seed)
					rng := rand.New(rand.NewSource(seed))
					for i := 0; i < 100; i++ {
						input := protorandom.Message(rng, messageType.Descriptor(), nil)
						newValue, newProtoreflectValue := generator(rng, valueField)

						expected := proto.Clone(input.Interface()).ProtoReflect()
						expectedList := expected.Mutable(valueField).List()
						expectedList.Append(newProtoreflectValue)

						arguments := lo.Ternary(supportsPacked, []string{", TRUE", ", FALSE"}, []string{""})

						for _, argument := range arguments {
							RunTestThatExpression(t, fmt.Sprintf("pb_wire_json_to_message(%s(pb_message_to_wire_json(?), 2, ?%s))", strings.ReplaceAll(addFunction, "{kind}", "wire_json"), argument), input.Interface(), newValue).
								IsEqualToProto(expected.Interface())

							RunTestThatExpression(t, fmt.Sprintf("%s(?, 2, ?%s)", strings.ReplaceAll(addFunction, "{kind}", "message"), argument), input.Interface(), newValue).
								IsEqualToProto(expected.Interface())
						}
					}
				})
			})
		}
	}

	test(t, "repeated float", "pb_{kind}_add_repeated_float_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Float(rng, false, false)
		return newValue, protoreflect.ValueOfFloat32(newValue)
	})

	test(t, "repeated double", "pb_{kind}_add_repeated_double_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Double(rng, false, false)
		return newValue, protoreflect.ValueOfFloat64(newValue)
	})

	test(t, "repeated int32", "pb_{kind}_add_repeated_int32_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int32(rng)
		return newValue, protoreflect.ValueOfInt32(newValue)
	})

	test(t, "repeated int64", "pb_{kind}_add_repeated_int64_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int64(rng)
		return newValue, protoreflect.ValueOfInt64(newValue)
	})

	test(t, "repeated uint32", "pb_{kind}_add_repeated_uint32_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint32(rng)
		return newValue, protoreflect.ValueOfUint32(newValue)
	})

	test(t, "repeated uint64", "pb_{kind}_add_repeated_uint64_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint64(rng)
		return newValue, protoreflect.ValueOfUint64(newValue)
	})

	test(t, "repeated bool", "pb_{kind}_add_repeated_bool_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Bool(rng)
		return newValue, protoreflect.ValueOfBool(newValue)
	})

	test(t, "repeated string", "pb_{kind}_add_repeated_string_field_element", false, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.String(rng, 5)
		return newValue, protoreflect.ValueOfString(newValue)
	})

	test(t, "repeated bytes", "pb_{kind}_add_repeated_bytes_field_element", false, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Bytes(rng, 5)
		return newValue, protoreflect.ValueOfBytes(newValue)
	})

	test(t, "repeated sint32", "pb_{kind}_add_repeated_sint32_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int32(rng)
		return newValue, protoreflect.ValueOfInt32(newValue)
	})

	test(t, "repeated sint64", "pb_{kind}_add_repeated_sint64_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int64(rng)
		return newValue, protoreflect.ValueOfInt64(newValue)
	})

	test(t, "repeated fixed32", "pb_{kind}_add_repeated_fixed32_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint32(rng)
		return newValue, protoreflect.ValueOfUint32(newValue)
	})

	test(t, "repeated fixed64", "pb_{kind}_add_repeated_fixed64_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint64(rng)
		return newValue, protoreflect.ValueOfUint64(newValue)
	})

	test(t, "repeated sfixed32", "pb_{kind}_add_repeated_sfixed32_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int32(rng)
		return newValue, protoreflect.ValueOfInt32(newValue)
	})

	test(t, "repeated sfixed64", "pb_{kind}_add_repeated_sfixed64_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int64(rng)
		return newValue, protoreflect.ValueOfInt64(newValue)
	})

	test(t, "repeated EnumType", "pb_{kind}_add_repeated_enum_field_element", true, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Enum(rng, fieldDescriptor.Enum())
		return newValue, protoreflect.ValueOfEnum(newValue)
	})

	test(t, "repeated MessageType", "pb_{kind}_add_repeated_message_field_element", false, func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Message(rng, fieldDescriptor.Message(), nil)
		return newValue.Interface(), protoreflect.ValueOfMessage(newValue)
	})
}
