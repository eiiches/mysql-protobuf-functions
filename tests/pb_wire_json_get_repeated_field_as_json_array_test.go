package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"github.com/eiiches/mysql-protobuf-functions/internal/protoreflectutils"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRandomizedWireJsonGetRepeatedFieldAsJsonArray(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, getFunction string, supportsPacked bool, extractJsonArray func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{}) {
		for _, usePacked := range lo.Ternary(supportsPacked, []string{"true", "false"}, []string{""}) {
			t.Run(fmt.Sprintf("%s/usePacked=%s", protoFieldType, usePacked), func(t *testing.T) {
				GivenFieldDefinitionsWithExtraFields(t, fmt.Sprintf("%s values = 2%s;", protoFieldType, FormatPackedOption(usePacked)), func(messageType protoreflect.MessageType) {
					valuesField := messageType.Descriptor().Fields().ByName("values")

					seed := time.Now().UnixNano()
					t.Logf("Using seed = %d.", seed)
					rng := rand.New(rand.NewSource(seed))
					for i := 0; i < 100; i++ {
						input := protorandom.Message(rng, messageType.Descriptor(), nil)

						expectedJsonArray := extractJsonArray(input, valuesField)

						RunTestThatExpression(t, fmt.Sprintf("%s(pb_message_to_wire_json(?), 2)", strings.ReplaceAll(getFunction, "{kind}", "wire_json")), input.Interface()).
							IsEqualToJson(expectedJsonArray)

						RunTestThatExpression(t, fmt.Sprintf("%s(?, 2)", strings.ReplaceAll(getFunction, "{kind}", "message")), input.Interface()).
							IsEqualToJson(expectedJsonArray)
					}
				})
			})
		}
	}

	test(t, "repeated float", "pb_{kind}_get_repeated_float_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) interface{} {
			floatValue := float32(value.Float())
			if math.IsNaN(float64(floatValue)) || math.IsInf(float64(floatValue), 0) {
				// NaN and Inf are typically represented as null in JSON
				return nil
			} else {
				return float64(floatValue)
			}
		})
	})

	test(t, "repeated double", "pb_{kind}_get_repeated_double_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) interface{} {
			doubleValue := value.Float()
			if math.IsNaN(doubleValue) || math.IsInf(doubleValue, 0) {
				// NaN and Inf are typically represented as null in JSON
				return nil
			} else {
				return doubleValue
			}
		})
	})

	test(t, "repeated int32", "pb_{kind}_get_repeated_int32_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) int32 {
			return int32(value.Int())
		})
	})

	test(t, "repeated int64", "pb_{kind}_get_repeated_int64_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) int64 {
			return value.Int()
		})
	})

	test(t, "repeated int64", "pb_{kind}_get_repeated_int64_field_as_json_string_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) string {
			return fmt.Sprintf("%d", value.Int())
		})
	})

	test(t, "repeated uint32", "pb_{kind}_get_repeated_uint32_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) uint32 {
			return uint32(value.Uint())
		})
	})

	test(t, "repeated uint64", "pb_{kind}_get_repeated_uint64_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) uint64 {
			return value.Uint()
		})
	})

	test(t, "repeated uint64", "pb_{kind}_get_repeated_uint64_field_as_json_string_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) string {
			return fmt.Sprintf("%d", value.Uint())
		})
	})

	test(t, "repeated bool", "pb_{kind}_get_repeated_bool_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) bool {
			return value.Bool()
		})
	})

	test(t, "repeated string", "pb_{kind}_get_repeated_string_field_as_json_array", false, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) string {
			return value.String()
		})
	})

	test(t, "repeated bytes", "pb_{kind}_get_repeated_bytes_field_as_json_array", false, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) []byte {
			return value.Bytes()
		})
	})

	test(t, "repeated sint32", "pb_{kind}_get_repeated_sint32_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) int32 {
			return int32(value.Int())
		})
	})

	test(t, "repeated sint64", "pb_{kind}_get_repeated_sint64_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) int64 {
			return value.Int()
		})
	})

	test(t, "repeated sint64", "pb_{kind}_get_repeated_sint64_field_as_json_string_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) string {
			return fmt.Sprintf("%d", value.Int())
		})
	})

	test(t, "repeated fixed32", "pb_{kind}_get_repeated_fixed32_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) uint32 {
			return uint32(value.Uint())
		})
	})

	test(t, "repeated fixed64", "pb_{kind}_get_repeated_fixed64_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) uint64 {
			return value.Uint()
		})
	})

	test(t, "repeated fixed64", "pb_{kind}_get_repeated_fixed64_field_as_json_string_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) string {
			return fmt.Sprintf("%d", value.Uint())
		})
	})

	test(t, "repeated sfixed32", "pb_{kind}_get_repeated_sfixed32_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) int32 {
			return int32(value.Int())
		})
	})

	test(t, "repeated sfixed64", "pb_{kind}_get_repeated_sfixed64_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) int64 {
			return value.Int()
		})
	})

	test(t, "repeated sfixed64", "pb_{kind}_get_repeated_sfixed64_field_as_json_string_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) string {
			return fmt.Sprintf("%d", value.Int())
		})
	})

	test(t, "repeated EnumType", "pb_{kind}_get_repeated_enum_field_as_json_array", true, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) int32 {
			return int32(value.Enum())
		})
	})

	test(t, "repeated MessageType", "pb_{kind}_get_repeated_message_field_as_json_array", false, func(input protoreflect.Message, field protoreflect.FieldDescriptor) interface{} {
		return protoreflectutils.MapToSlice(input.Get(field).List(), func(value protoreflect.Value) []byte {
			subMessage := value.Message().Interface()
			// Convert message to base64-encoded bytes (MySQL functions return base64)
			bytes, err := proto.Marshal(subMessage)
			if err != nil {
				panic(fmt.Sprintf("Failed to marshal submessage: %v", err))
			}
			return bytes
		})
	})
}
