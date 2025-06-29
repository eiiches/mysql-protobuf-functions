package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRandomizedWireJsonGetRepeatedField(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, getFunction string, supportsPacked bool, extractValue func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{}) {
		for _, usePacked := range lo.Ternary(supportsPacked, []string{"true", "false"}, []string{""}) {
			t.Run(fmt.Sprintf("%s/usePacked=%s", protoFieldType, usePacked), func(t *testing.T) {
				GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s values = 2%s; int32 b = 3;", protoFieldType, FormatPackedOption(usePacked)), func(messageType protoreflect.MessageType) {
					valuesField := messageType.Descriptor().Fields().ByName("values")

					seed := time.Now().UnixNano()
					t.Logf("Using seed = %d.", seed)
					rng := rand.New(rand.NewSource(seed))
					for i := 0; i < 100; i++ {
						input := protorandom.Message(rng, messageType.Descriptor(), nil)
						list := input.Get(valuesField).List()

						// Test each element in the repeated field
						for index := 0; index < list.Len(); index++ {
							expectedValue := extractValue(input, valuesField, index)

							RunTestThatExpression(t, fmt.Sprintf("%s(pb_message_to_wire_json(?), 2, ?)", strings.ReplaceAll(getFunction, "{kind}", "wire_json")), input.Interface(), index).
								IsEqualTo(expectedValue)

							RunTestThatExpression(t, fmt.Sprintf("%s(?, 2, ?)", strings.ReplaceAll(getFunction, "{kind}", "message")), input.Interface(), index).
								IsEqualTo(expectedValue)
						}
					}
				})
			})
		}
	}

	test(t, "repeated float", "pb_{kind}_get_repeated_float_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		value := input.Get(field).List().Get(index).Float()
		if math.IsNaN(value) || math.IsInf(value, 0) {
			// pb_{kind}_get_repeated_float_field returns NULL if the value is NaN, +Inf, or -Inf.
			return (*float32)(nil)
		}
		return float32(value)
	})

	test(t, "repeated double", "pb_{kind}_get_repeated_double_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		value := input.Get(field).List().Get(index).Float()
		if math.IsNaN(value) || math.IsInf(value, 0) {
			// pb_{kind}_get_repeated_double_field returns NULL if the value is NaN, +Inf, or -Inf.
			return (*float64)(nil)
		}
		return value
	})

	test(t, "repeated int32", "pb_{kind}_get_repeated_int32_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return int32(input.Get(field).List().Get(index).Int())
	})

	test(t, "repeated int64", "pb_{kind}_get_repeated_int64_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return input.Get(field).List().Get(index).Int()
	})

	test(t, "repeated uint32", "pb_{kind}_get_repeated_uint32_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return uint32(input.Get(field).List().Get(index).Uint())
	})

	test(t, "repeated uint64", "pb_{kind}_get_repeated_uint64_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return input.Get(field).List().Get(index).Uint()
	})

	test(t, "repeated bool", "pb_{kind}_get_repeated_bool_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return input.Get(field).List().Get(index).Bool()
	})

	test(t, "repeated string", "pb_{kind}_get_repeated_string_field", false, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return input.Get(field).List().Get(index).String()
	})

	test(t, "repeated bytes", "pb_{kind}_get_repeated_bytes_field", false, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return input.Get(field).List().Get(index).Bytes()
	})

	test(t, "repeated sint32", "pb_{kind}_get_repeated_sint32_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return int32(input.Get(field).List().Get(index).Int())
	})

	test(t, "repeated sint64", "pb_{kind}_get_repeated_sint64_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return input.Get(field).List().Get(index).Int()
	})

	test(t, "repeated fixed32", "pb_{kind}_get_repeated_fixed32_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return uint32(input.Get(field).List().Get(index).Uint())
	})

	test(t, "repeated fixed64", "pb_{kind}_get_repeated_fixed64_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return input.Get(field).List().Get(index).Uint()
	})

	test(t, "repeated sfixed32", "pb_{kind}_get_repeated_sfixed32_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return int32(input.Get(field).List().Get(index).Int())
	})

	test(t, "repeated sfixed64", "pb_{kind}_get_repeated_sfixed64_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return input.Get(field).List().Get(index).Int()
	})

	test(t, "repeated EnumType", "pb_{kind}_get_repeated_enum_field", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		return int32(input.Get(field).List().Get(index).Enum())
	})

	test(t, "repeated MessageType", "pb_{kind}_get_repeated_message_field", false, func(input protoreflect.Message, field protoreflect.FieldDescriptor, index int) interface{} {
		subMessage := input.Get(field).List().Get(index).Message().Interface()
		bytes, err := proto.Marshal(subMessage)
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal submessage: %v", err))
		}
		return bytes
	})
}
