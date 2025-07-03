package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRandomizedWireJsonGetField(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, getFunction string, defaultValue interface{}, extractValue func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{}) {
		t.Run(fmt.Sprintf("%s", protoFieldType), func(t *testing.T) {
			GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s value = 2; int32 b = 3;", protoFieldType), func(messageType protoreflect.MessageType) {
				valueField := messageType.Descriptor().Fields().ByName("value")

				seed := time.Now().UnixNano()
				t.Logf("Using seed = %d.", seed)
				rng := rand.New(rand.NewSource(seed))
				for i := 0; i < 100; i++ {
					input := protorandom.Message(rng, messageType.Descriptor(), nil)
					expectedValue := extractValue(input, valueField)

					RunTestThatExpression(t, fmt.Sprintf("%s(pb_message_to_wire_json(?), 2, ?)", strings.ReplaceAll(getFunction, "{kind}", "wire_json")), input.Interface(), defaultValue).
						IsEqualTo(expectedValue)

					RunTestThatExpression(t, fmt.Sprintf("%s(?, 2, ?)", strings.ReplaceAll(getFunction, "{kind}", "message")), input.Interface(), defaultValue).
						IsEqualTo(expectedValue)
				}
			})
		})
	}

	test(t, "float", "pb_{kind}_get_float_field", float32(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		value := input.Get(field).Float()
		if math.IsNaN(value) || math.IsInf(value, 0) {
			// pb_{kind}_get_float_field returns NULL if the value is NaN, +Inf, or -Inf.
			return (*float32)(nil)
		}
		return float32(value)
	})

	test(t, "double", "pb_{kind}_get_double_field", float64(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		value := input.Get(field).Float()
		if math.IsNaN(value) || math.IsInf(value, 0) {
			// pb_{kind}_get_double_field returns NULL if the value is NaN, +Inf, or -Inf.
			return (*float64)(nil)
		}
		return value
	})

	test(t, "int32", "pb_{kind}_get_int32_field", int32(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return int32(input.Get(field).Int())
	})

	test(t, "int64", "pb_{kind}_get_int64_field", int64(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return input.Get(field).Int()
	})

	test(t, "uint32", "pb_{kind}_get_uint32_field", uint32(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return uint32(input.Get(field).Uint())
	})

	test(t, "uint64", "pb_{kind}_get_uint64_field", uint64(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return input.Get(field).Uint()
	})

	test(t, "bool", "pb_{kind}_get_bool_field", false, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return input.Get(field).Bool()
	})

	test(t, "string", "pb_{kind}_get_string_field", "", func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return input.Get(field).String()
	})

	test(t, "bytes", "pb_{kind}_get_bytes_field", []byte{}, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return input.Get(field).Bytes()
	})

	test(t, "sint32", "pb_{kind}_get_sint32_field", int32(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return int32(input.Get(field).Int())
	})

	test(t, "sint64", "pb_{kind}_get_sint64_field", int64(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return input.Get(field).Int()
	})

	test(t, "fixed32", "pb_{kind}_get_fixed32_field", uint32(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return uint32(input.Get(field).Uint())
	})

	test(t, "fixed64", "pb_{kind}_get_fixed64_field", uint64(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return input.Get(field).Uint()
	})

	test(t, "sfixed32", "pb_{kind}_get_sfixed32_field", int32(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return int32(input.Get(field).Int())
	})

	test(t, "sfixed64", "pb_{kind}_get_sfixed64_field", int64(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return input.Get(field).Int()
	})

	test(t, "EnumType", "pb_{kind}_get_enum_field", int32(0), func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return int32(input.Get(field).Enum())
	})

	test(t, "MessageType", "pb_{kind}_get_message_field", []byte{}, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		subMessage := input.Get(field).Message().Interface()
		bytes, err := proto.Marshal(subMessage)
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal submessage: %v", err))
		}
		return bytes
	})
}
