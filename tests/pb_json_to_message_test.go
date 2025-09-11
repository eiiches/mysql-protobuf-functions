package main

import (
	"fmt"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/jsonoptionspb"
	"github.com/eiiches/mysql-protobuf-functions/internal/protonumberjson"

	"google.golang.org/protobuf/proto"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	"github.com/eiiches/mysql-protobuf-functions/internal/descriptorsetjson"
	"github.com/eiiches/mysql-protobuf-functions/internal/testutils"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func testJsonToMessage(t *testing.T, fieldDefinition string, input string) {
	g := NewWithT(t)

	p := testutils.NewProtoTestSupport(t, map[string]string{
		"main.proto": fmt.Sprintf(dedent.Pipe(`
			|syntax = "proto3";
			|import "google/protobuf/timestamp.proto";
			|import "google/protobuf/duration.proto";
			|import "google/protobuf/struct.proto";
			|import "google/protobuf/empty.proto";
			|import "google/protobuf/wrappers.proto";
			|import "google/protobuf/field_mask.proto";
			|message Test {
			|    %s
			|}
			|message MessageType {
			|    int32 value = 1;
			|}
			|enum EnumType {
			|    ENUM_TYPE_UNSPECIFIED = 0;
			|    ENUM_TYPE_ONE = 1;
			|}
		`), fieldDefinition),
	})

	typeName := protoreflect.FullName(".Test")

	// Generate descriptor set JSON using descriptorsetjson package
	descriptorSetJson, err := descriptorsetjson.ToJson(p.GetFileDescriptorSet())
	g.Expect(err).NotTo(HaveOccurred())

	// Parse JSON input with Go protojson to get expected protobuf result
	expectedMessage := p.JsonToDynamicMessage(typeName, input)
	expectedMessageBytes, err := proto.Marshal(expectedMessage.Interface())
	g.Expect(err).NotTo(HaveOccurred())

	// Test pb_json_to_message: MySQL implementation should produce the same protobuf as Go's protojson
	RunTestThatExpression(t, "pb_json_to_message(?, ?, ?, NULL, NULL)", descriptorSetJson, typeName, input).IsEqualToProto(expectedMessage.Interface())

	// Field order may differ, but the serialized bytes should match
	RunTestThatExpression(t, "LENGTH(pb_json_to_message(?, ?, ?, NULL, NULL))", descriptorSetJson, typeName, input).IsEqualTo(len(expectedMessageBytes))
}

// Test function specifically for enum conversions that handles the fact that
// numeric enum inputs are canonically converted to string outputs
func testEnumConversion(t *testing.T, fieldDefinition string, input string, expectedOutput string) {
	g := NewWithT(t)

	p := testutils.NewProtoTestSupport(t, map[string]string{
		"main.proto": fmt.Sprintf(dedent.Pipe(`
			|syntax = "proto3";
			|import "google/protobuf/timestamp.proto";
			|import "google/protobuf/duration.proto";
			|import "google/protobuf/struct.proto";
			|import "google/protobuf/empty.proto";
			|import "google/protobuf/wrappers.proto";
			|import "google/protobuf/field_mask.proto";
			|message Test {
			|    %s
			|}
			|message MessageType {
			|    int32 value = 1;
			|}
			|enum EnumType {
			|    ENUM_TYPE_UNSPECIFIED = 0;
			|    ENUM_TYPE_ONE = 1;
			|}
		`), fieldDefinition),
	})

	typeName := protoreflect.FullName(".Test")

	// Generate descriptor set JSON using descriptorsetjson package
	descriptorSetJson, err := descriptorsetjson.ToJson(p.GetFileDescriptorSet())
	g.Expect(err).NotTo(HaveOccurred())

	// Create JsonMarshalOptions with emit_default_values = true
	marshalOptionsWithDefaultsJSON, err := protonumberjson.Marshal(&jsonoptionspb.JsonMarshalOptions{
		EmitDefaultValues: true,
	})
	g.Expect(err).NotTo(HaveOccurred())

	// Test enum conversion: JSON → wire_json → JSON produces canonical string format
	RunTestThatExpression(t, "pb_wire_json_to_json(?, ?, pb_json_to_wire_json(?, ?, ?, NULL, NULL), NULL, ?)", descriptorSetJson, typeName, descriptorSetJson, typeName, input, string(marshalOptionsWithDefaultsJSON)).IsEqualToJsonString(expectedOutput)
}

func TestJsonToMessageSingularFields(t *testing.T) {
	t.Run("int32", func(t *testing.T) {
		testJsonToMessage(t, "int32 int32_field = 1;", `{"int32Field": 0}`)
		testJsonToMessage(t, "int32 int32_field = 1;", `{"int32Field": 2147483647}`)
		testJsonToMessage(t, "int32 int32_field = 1;", `{"int32Field": -2147483648}`)
	})

	t.Run("uint32", func(t *testing.T) {
		testJsonToMessage(t, "uint32 uint32_field = 1;", `{"uint32Field": 0}`)
		testJsonToMessage(t, "uint32 uint32_field = 1;", `{"uint32Field": 4294967295}`)
	})

	t.Run("int64", func(t *testing.T) {
		testJsonToMessage(t, "int64 int64_field = 1;", `{"int64Field": "0"}`)
		testJsonToMessage(t, "int64 int64_field = 1;", `{"int64Field": "9223372036854775807"}`)
		testJsonToMessage(t, "int64 int64_field = 1;", `{"int64Field": "-9223372036854775808"}`)
	})

	t.Run("uint64", func(t *testing.T) {
		testJsonToMessage(t, "uint64 uint64_field = 1;", `{"uint64Field": "0"}`)
		testJsonToMessage(t, "uint64 uint64_field = 1;", `{"uint64Field": "18446744073709551615"}`)
	})

	t.Run("fixed32", func(t *testing.T) {
		testJsonToMessage(t, "fixed32 fixed32_field = 1;", `{"fixed32Field": 0}`)
		testJsonToMessage(t, "fixed32 fixed32_field = 1;", `{"fixed32Field": 4294967295}`)
	})

	t.Run("fixed64", func(t *testing.T) {
		testJsonToMessage(t, "fixed64 fixed64_field = 1;", `{"fixed64Field": "0"}`)
		testJsonToMessage(t, "fixed64 fixed64_field = 1;", `{"fixed64Field": "18446744073709551615"}`)
	})

	t.Run("sfixed32", func(t *testing.T) {
		testJsonToMessage(t, "sfixed32 sfixed32_field = 1;", `{"sfixed32Field": 0}`)
		testJsonToMessage(t, "sfixed32 sfixed32_field = 1;", `{"sfixed32Field": 2147483647}`)
		testJsonToMessage(t, "sfixed32 sfixed32_field = 1;", `{"sfixed32Field": -2147483648}`)
	})

	t.Run("sfixed64", func(t *testing.T) {
		testJsonToMessage(t, "sfixed64 sfixed64_field = 1;", `{"sfixed64Field": "0"}`)
		testJsonToMessage(t, "sfixed64 sfixed64_field = 1;", `{"sfixed64Field": "9223372036854775807"}`)
		testJsonToMessage(t, "sfixed64 sfixed64_field = 1;", `{"sfixed64Field": "-9223372036854775808"}`)
	})

	t.Run("bool", func(t *testing.T) {
		testJsonToMessage(t, "bool bool_field = 1;", `{"boolField": false}`)
		testJsonToMessage(t, "bool bool_field = 1;", `{"boolField": true}`)
	})

	t.Run("string", func(t *testing.T) {
		testJsonToMessage(t, "string string_field = 1;", `{"stringField": ""}`)
		testJsonToMessage(t, "string string_field = 1;", `{"stringField": "testJsonToMessage"}`)
	})

	t.Run("bytes", func(t *testing.T) {
		testJsonToMessage(t, "bytes bytes_field = 1;", `{"bytesField": ""}`)
		testJsonToMessage(t, "bytes bytes_field = 1;", `{"bytesField": "dGVzdA=="}`) // Base64 for "test"
	})

	t.Run("float", func(t *testing.T) {
		testJsonToMessage(t, "float float_field = 1;", `{"floatField": 0}`)
		testJsonToMessage(t, "float float_field = 1;", `{"floatField": 3.5}`)
	})

	t.Run("double", func(t *testing.T) {
		testJsonToMessage(t, "double double_field = 1;", `{"doubleField": 0}`)
		testJsonToMessage(t, "double double_field = 1;", `{"doubleField": 3.141592653589793}`)
	})

	t.Run("sint32", func(t *testing.T) {
		testJsonToMessage(t, "sint32 sint32_field = 1;", `{"sint32Field": 0}`)
		testJsonToMessage(t, "sint32 sint32_field = 1;", `{"sint32Field": 2147483647}`)
		testJsonToMessage(t, "sint32 sint32_field = 1;", `{"sint32Field": -2147483648}`)
	})

	t.Run("sint64", func(t *testing.T) {
		testJsonToMessage(t, "sint64 sint64_field = 1;", `{"sint64Field": "0"}`)
		testJsonToMessage(t, "sint64 sint64_field = 1;", `{"sint64Field": "9223372036854775807"}`)
		testJsonToMessage(t, "sint64 sint64_field = 1;", `{"sint64Field": "-9223372036854775808"}`)
	})

	t.Run("enum", func(t *testing.T) {
		testJsonToMessage(t, "EnumType enum_field = 1;", `{"enumField": "ENUM_TYPE_UNSPECIFIED"}`)
		testJsonToMessage(t, "EnumType enum_field = 1;", `{"enumField": "ENUM_TYPE_ONE"}`)
		// Test numeric enum values - they get converted to canonical string format
		testEnumConversion(t, "EnumType enum_field = 1;", `{"enumField": 0}`, `{"enumField": "ENUM_TYPE_UNSPECIFIED"}`)
		testEnumConversion(t, "EnumType enum_field = 1;", `{"enumField": 1}`, `{"enumField": "ENUM_TYPE_ONE"}`)
	})

	t.Run("message", func(t *testing.T) {
		testJsonToMessage(t, "MessageType message_field = 1;", `{}`)
		testJsonToMessage(t, "MessageType message_field = 1;", `{"messageField": {"value": 0}}`)
		testJsonToMessage(t, "MessageType message_field = 1;", `{"messageField": {"value": 12345}}`)
	})
}

func TestJsonToMessageRepeatedFields(t *testing.T) {
	t.Run("repeated int32", func(t *testing.T) {
		testJsonToMessage(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": []}`)
		testJsonToMessage(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": [0]}`)
		testJsonToMessage(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated uint32", func(t *testing.T) {
		testJsonToMessage(t, "repeated uint32 repeated_uint32_field = 1;", `{"repeatedUint32Field": []}`)
		testJsonToMessage(t, "repeated uint32 repeated_uint32_field = 1;", `{"repeatedUint32Field": [0]}`)
		testJsonToMessage(t, "repeated uint32 repeated_uint32_field = 1;", `{"repeatedUint32Field": [0, 4294967295]}`)
	})

	t.Run("repeated int64", func(t *testing.T) {
		testJsonToMessage(t, "repeated int64 repeated_int64_field = 1;", `{"repeatedInt64Field": []}`)
		testJsonToMessage(t, "repeated int64 repeated_int64_field = 1;", `{"repeatedInt64Field": ["0"]}`)
		testJsonToMessage(t, "repeated int64 repeated_int64_field = 1;", `{"repeatedInt64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`)
	})

	t.Run("repeated uint64", func(t *testing.T) {
		testJsonToMessage(t, "repeated uint64 repeated_uint64_field = 1;", `{"repeatedUint64Field": []}`)
		testJsonToMessage(t, "repeated uint64 repeated_uint64_field = 1;", `{"repeatedUint64Field": ["0"]}`)
		testJsonToMessage(t, "repeated uint64 repeated_uint64_field = 1;", `{"repeatedUint64Field": ["0", "18446744073709551615"]}`)
	})

	t.Run("repeated fixed32", func(t *testing.T) {
		testJsonToMessage(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"repeatedFixed32Field": []}`)
		testJsonToMessage(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"repeatedFixed32Field": [0]}`)
		testJsonToMessage(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"repeatedFixed32Field": [0, 4294967295]}`)
	})

	t.Run("repeated fixed64", func(t *testing.T) {
		testJsonToMessage(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"repeatedFixed64Field": []}`)
		testJsonToMessage(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"repeatedFixed64Field": ["0"]}`)
		testJsonToMessage(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"repeatedFixed64Field": ["0", "18446744073709551615"]}`)
	})

	t.Run("repeated sfixed32", func(t *testing.T) {
		testJsonToMessage(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"repeatedSfixed32Field": []}`)
		testJsonToMessage(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"repeatedSfixed32Field": [0]}`)
		testJsonToMessage(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"repeatedSfixed32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated sfixed64", func(t *testing.T) {
		testJsonToMessage(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"repeatedSfixed64Field": []}`)
		testJsonToMessage(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"repeatedSfixed64Field": ["0"]}`)
		testJsonToMessage(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"repeatedSfixed64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`)
	})

	t.Run("repeated bool", func(t *testing.T) {
		testJsonToMessage(t, "repeated bool repeated_bool_field = 1;", `{"repeatedBoolField": []}`)
		testJsonToMessage(t, "repeated bool repeated_bool_field = 1;", `{"repeatedBoolField": [false]}`)
		testJsonToMessage(t, "repeated bool repeated_bool_field = 1;", `{"repeatedBoolField": [true, false, true]}`)
	})

	t.Run("repeated string", func(t *testing.T) {
		testJsonToMessage(t, "repeated string repeated_string_field = 1;", `{"repeatedStringField": []}`)
		testJsonToMessage(t, "repeated string repeated_string_field = 1;", `{"repeatedStringField": [""]}`)
		testJsonToMessage(t, "repeated string repeated_string_field = 1;", `{"repeatedStringField": ["testJsonToMessage", ""]}`)
	})

	t.Run("repeated bytes", func(t *testing.T) {
		testJsonToMessage(t, "repeated bytes repeated_bytes_field = 1;", `{"repeatedBytesField": []}`)
		testJsonToMessage(t, "repeated bytes repeated_bytes_field = 1;", `{"repeatedBytesField": [""]}`)
		testJsonToMessage(t, "repeated bytes repeated_bytes_field = 1;", `{"repeatedBytesField": ["dGVzdA==", ""]}`) // Base64 for "test"
	})

	t.Run("repeated float", func(t *testing.T) {
		testJsonToMessage(t, "repeated float repeated_float_field = 1;", `{"repeatedFloatField": []}`)
		testJsonToMessage(t, "repeated float repeated_float_field = 1;", `{"repeatedFloatField": [0]}`)
		testJsonToMessage(t, "repeated float repeated_float_field = 1;", `{"repeatedFloatField": [3.5, 0]}`)
	})

	t.Run("repeated double", func(t *testing.T) {
		testJsonToMessage(t, "repeated double repeated_double_field = 1;", `{"repeatedDoubleField": []}`)
		testJsonToMessage(t, "repeated double repeated_double_field = 1;", `{"repeatedDoubleField": [0]}`)
		testJsonToMessage(t, "repeated double repeated_double_field = 1;", `{"repeatedDoubleField": [3.141592653589793, 0]}`)
	})

	t.Run("repeated sint32", func(t *testing.T) {
		testJsonToMessage(t, "repeated sint32 repeated_sint32_field = 1;", `{"repeatedSint32Field": []}`)
		testJsonToMessage(t, "repeated sint32 repeated_sint32_field = 1;", `{"repeatedSint32Field": [0]}`)
		testJsonToMessage(t, "repeated sint32 repeated_sint32_field = 1;", `{"repeatedSint32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated sint64", func(t *testing.T) {
		testJsonToMessage(t, "repeated sint64 repeated_sint64_field = 1;", `{"repeatedSint64Field": []}`)
		testJsonToMessage(t, "repeated sint64 repeated_sint64_field = 1;", `{"repeatedSint64Field": ["0"]}`)
		testJsonToMessage(t, "repeated sint64 repeated_sint64_field = 1;", `{"repeatedSint64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`)
	})

	t.Run("repeated enum", func(t *testing.T) {
		testJsonToMessage(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": []}`)
		testJsonToMessage(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": ["ENUM_TYPE_UNSPECIFIED"]}`)
		testJsonToMessage(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": ["ENUM_TYPE_ONE", "ENUM_TYPE_UNSPECIFIED"]}`)
	})

	t.Run("repeated message", func(t *testing.T) {
		testJsonToMessage(t, "repeated MessageType repeated_message_field = 1;", `{"repeatedMessageField": []}`)
		testJsonToMessage(t, "repeated MessageType repeated_message_field = 1;", `{"repeatedMessageField": [{"value": 0}]}`)
		testJsonToMessage(t, "repeated MessageType repeated_message_field = 1;", `{"repeatedMessageField": [{"value": 12345}, {"value": 67890}]}`)
	})
}

func TestJsonToMessageOptionalFields(t *testing.T) {
	t.Run("optional int32", func(t *testing.T) {
		testJsonToMessage(t, "optional int32 optional_int32_field = 1;", `{}`)
		testJsonToMessage(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": 0}`)
		testJsonToMessage(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": 2147483647}`)
		testJsonToMessage(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": -2147483648}`)
	})

	t.Run("optional uint32", func(t *testing.T) {
		testJsonToMessage(t, "optional uint32 optional_uint32_field = 1;", `{}`)
		testJsonToMessage(t, "optional uint32 optional_uint32_field = 1;", `{"optionalUint32Field": 0}`)
		testJsonToMessage(t, "optional uint32 optional_uint32_field = 1;", `{"optionalUint32Field": 4294967295}`)
	})

	t.Run("optional int64", func(t *testing.T) {
		testJsonToMessage(t, "optional int64 optional_int64_field = 1;", `{}`)
		testJsonToMessage(t, "optional int64 optional_int64_field = 1;", `{"optionalInt64Field": "0"}`)
		testJsonToMessage(t, "optional int64 optional_int64_field = 1;", `{"optionalInt64Field": "9223372036854775807"}`)
		testJsonToMessage(t, "optional int64 optional_int64_field = 1;", `{"optionalInt64Field": "-9223372036854775808"}`)
	})

	t.Run("optional uint64", func(t *testing.T) {
		testJsonToMessage(t, "optional uint64 optional_uint64_field = 1;", `{}`)
		testJsonToMessage(t, "optional uint64 optional_uint64_field = 1;", `{"optionalUint64Field": "0"}`)
		testJsonToMessage(t, "optional uint64 optional_uint64_field = 1;", `{"optionalUint64Field": "18446744073709551615"}`)
	})

	t.Run("optional fixed32", func(t *testing.T) {
		testJsonToMessage(t, "optional fixed32 optional_fixed32_field = 1;", `{}`)
		testJsonToMessage(t, "optional fixed32 optional_fixed32_field = 1;", `{"optionalFixed32Field": 0}`)
		testJsonToMessage(t, "optional fixed32 optional_fixed32_field = 1;", `{"optionalFixed32Field": 4294967295}`)
	})

	t.Run("optional fixed64", func(t *testing.T) {
		testJsonToMessage(t, "optional fixed64 optional_fixed64_field = 1;", `{}`)
		testJsonToMessage(t, "optional fixed64 optional_fixed64_field = 1;", `{"optionalFixed64Field": "0"}`)
		testJsonToMessage(t, "optional fixed64 optional_fixed64_field = 1;", `{"optionalFixed64Field": "18446744073709551615"}`)
	})

	t.Run("optional sfixed32", func(t *testing.T) {
		testJsonToMessage(t, "optional sfixed32 optional_sfixed32_field = 1;", `{}`)
		testJsonToMessage(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"optionalSfixed32Field": 0}`)
		testJsonToMessage(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"optionalSfixed32Field": 2147483647}`)
		testJsonToMessage(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"optionalSfixed32Field": -2147483648}`)
	})

	t.Run("optional sfixed64", func(t *testing.T) {
		testJsonToMessage(t, "optional sfixed64 optional_sfixed64_field = 1;", `{}`)
		testJsonToMessage(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"optionalSfixed64Field": "0"}`)
		testJsonToMessage(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"optionalSfixed64Field": "9223372036854775807"}`)
		testJsonToMessage(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"optionalSfixed64Field": "-9223372036854775808"}`)
	})

	t.Run("optional bool", func(t *testing.T) {
		testJsonToMessage(t, "optional bool optional_bool_field = 1;", `{}`)
		testJsonToMessage(t, "optional bool optional_bool_field = 1;", `{"optionalBoolField": false}`)
		testJsonToMessage(t, "optional bool optional_bool_field = 1;", `{"optionalBoolField": true}`)
	})

	t.Run("optional string", func(t *testing.T) {
		testJsonToMessage(t, "optional string optional_string_field = 1;", `{}`)
		testJsonToMessage(t, "optional string optional_string_field = 1;", `{"optionalStringField": ""}`)
		testJsonToMessage(t, "optional string optional_string_field = 1;", `{"optionalStringField": "testJsonToMessage"}`)
	})

	t.Run("optional bytes", func(t *testing.T) {
		testJsonToMessage(t, "optional bytes optional_bytes_field = 1;", `{}`)
		testJsonToMessage(t, "optional bytes optional_bytes_field = 1;", `{"optionalBytesField": ""}`)
		testJsonToMessage(t, "optional bytes optional_bytes_field = 1;", `{"optionalBytesField": "dGVzdA=="}`) // Base64 for "test"
	})

	t.Run("optional float", func(t *testing.T) {
		testJsonToMessage(t, "optional float optional_float_field = 1;", `{}`)
		testJsonToMessage(t, "optional float optional_float_field = 1;", `{"optionalFloatField": 0}`)
		testJsonToMessage(t, "optional float optional_float_field = 1;", `{"optionalFloatField": 3.5}`)
	})

	t.Run("optional double", func(t *testing.T) {
		testJsonToMessage(t, "optional double optional_double_field = 1;", `{}`)
		testJsonToMessage(t, "optional double optional_double_field = 1;", `{"optionalDoubleField": 0}`)
		testJsonToMessage(t, "optional double optional_double_field = 1;", `{"optionalDoubleField": 3.141592653589793}`)
	})

	t.Run("optional sint32", func(t *testing.T) {
		testJsonToMessage(t, "optional sint32 optional_sint32_field = 1;", `{}`)
		testJsonToMessage(t, "optional sint32 optional_sint32_field = 1;", `{"optionalSint32Field": 0}`)
		testJsonToMessage(t, "optional sint32 optional_sint32_field = 1;", `{"optionalSint32Field": 2147483647}`)
		testJsonToMessage(t, "optional sint32 optional_sint32_field = 1;", `{"optionalSint32Field": -2147483648}`)
	})

	t.Run("optional sint64", func(t *testing.T) {
		testJsonToMessage(t, "optional sint64 optional_sint64_field = 1;", `{}`)
		testJsonToMessage(t, "optional sint64 optional_sint64_field = 1;", `{"optionalSint64Field": "0"}`)
		testJsonToMessage(t, "optional sint64 optional_sint64_field = 1;", `{"optionalSint64Field": "9223372036854775807"}`)
		testJsonToMessage(t, "optional sint64 optional_sint64_field = 1;", `{"optionalSint64Field": "-9223372036854775808"}`)
	})

	t.Run("optional enum", func(t *testing.T) {
		testJsonToMessage(t, "optional EnumType optional_enum_field = 1;", `{}`)
		testJsonToMessage(t, "optional EnumType optional_enum_field = 1;", `{"optionalEnumField": "ENUM_TYPE_UNSPECIFIED"}`)
		testJsonToMessage(t, "optional EnumType optional_enum_field = 1;", `{"optionalEnumField": "ENUM_TYPE_ONE"}`)
	})

	t.Run("optional message", func(t *testing.T) {
		testJsonToMessage(t, "optional MessageType optional_message_field = 1;", `{}`)
		testJsonToMessage(t, "optional MessageType optional_message_field = 1;", `{"optionalMessageField": {"value": 0}}`)
		testJsonToMessage(t, "optional MessageType optional_message_field = 1;", `{"optionalMessageField": {"value": 12345}}`)
	})
}

func TestJsonToMessageMapKey(t *testing.T) {
	t.Run("map<int32, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<int32, string> int32_key_map_field = 1;", `{"int32KeyMapField": {}}`)
		testJsonToMessage(t, "map<int32, string> int32_key_map_field = 1;", `{"int32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<uint32, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<uint32, string> uint32_key_map_field = 1;", `{"uint32KeyMapField": {}}`)
		testJsonToMessage(t, "map<uint32, string> uint32_key_map_field = 1;", `{"uint32KeyMapField": {"0": "a", "4294967295": "b"}}`)
	})

	t.Run("map<int64, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<int64, string> int64_key_map_field = 1;", `{"int64KeyMapField": {}}`)
		testJsonToMessage(t, "map<int64, string> int64_key_map_field = 1;", `{"int64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	t.Run("map<uint64, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<uint64, string> uint64_key_map_field = 1;", `{"uint64KeyMapField": {}}`)
		testJsonToMessage(t, "map<uint64, string> uint64_key_map_field = 1;", `{"uint64KeyMapField": {"0": "a", "18446744073709551615": "b"}}`)
	})

	t.Run("map<fixed32, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<fixed32, string> fixed32_key_map_field = 1;", `{"fixed32KeyMapField": {}}`)
		testJsonToMessage(t, "map<fixed32, string> fixed32_key_map_field = 1;", `{"fixed32KeyMapField": {"0": "a", "4294967295": "b"}}`)
	})

	t.Run("map<fixed64, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<fixed64, string> fixed64_key_map_field = 1;", `{"fixed64KeyMapField": {}}`)
		testJsonToMessage(t, "map<fixed64, string> fixed64_key_map_field = 1;", `{"fixed64KeyMapField": {"0": "a", "18446744073709551615": "b"}}`)
	})

	t.Run("map<sfixed32, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<sfixed32, string> sfixed32_key_map_field = 1;", `{"sfixed32KeyMapField": {}}`)
		testJsonToMessage(t, "map<sfixed32, string> sfixed32_key_map_field = 1;", `{"sfixed32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<sfixed64, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<sfixed64, string> sfixed64_key_map_field = 1;", `{"sfixed64KeyMapField": {}}`)
		testJsonToMessage(t, "map<sfixed64, string> sfixed64_key_map_field = 1;", `{"sfixed64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	t.Run("map<bool, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<bool, string> bool_key_map_field = 1;", `{"boolKeyMapField": {}}`)
		testJsonToMessage(t, "map<bool, string> bool_key_map_field = 1;", `{"boolKeyMapField": {"false": "a", "true": "b"}}`)
	})

	t.Run("map<string, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, string> string_key_map_field = 1;", `{"stringKeyMapField": {}}`)
		testJsonToMessage(t, "map<string, string> string_key_map_field = 1;", `{"stringKeyMapField": {"a": "b", "c": "d"}}`)
	})

	t.Run("map<sint32, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<sint32, string> sint32_key_map_field = 1;", `{"sint32KeyMapField": {}}`)
		testJsonToMessage(t, "map<sint32, string> sint32_key_map_field = 1;", `{"sint32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<sint64, *>", func(t *testing.T) {
		testJsonToMessage(t, "map<sint64, string> sint64_key_map_field = 1;", `{"sint64KeyMapField": {}}`)
		testJsonToMessage(t, "map<sint64, string> sint64_key_map_field = 1;", `{"sint64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	// NOTE: float, double, enum, or message cannot be a map key.
}

func TestJsonToMessageMapValue(t *testing.T) {
	t.Run("map<*, int32>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, int32> int32_value_map_field = 1;", `{"int32ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, int32> int32_value_map_field = 1;", `{"int32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, uint32>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, uint32> uint32_value_map_field = 1;", `{"uint32ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, uint32> uint32_value_map_field = 1;", `{"uint32ValueMapField": {"a": 0, "b": 4294967295}}`)
	})

	t.Run("map<*, int64>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, int64> int64_value_map_field = 1;", `{"int64ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, int64> int64_value_map_field = 1;", `{"int64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`)
	})

	t.Run("map<*, uint64>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, uint64> uint64_value_map_field = 1;", `{"uint64ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, uint64> uint64_value_map_field = 1;", `{"uint64ValueMapField": {"a": "0", "b": "18446744073709551615"}}`)
	})

	t.Run("map<*, fixed32>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, fixed32> fixed32_value_map_field = 1;", `{"fixed32ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, fixed32> fixed32_value_map_field = 1;", `{"fixed32ValueMapField": {"a": 0, "b": 4294967295}}`)
	})

	t.Run("map<*, fixed64>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, fixed64> fixed64_value_map_field = 1;", `{"fixed64ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, fixed64> fixed64_value_map_field = 1;", `{"fixed64ValueMapField": {"a": "0", "b": "18446744073709551615"}}`)
	})

	t.Run("map<*, sfixed32>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, sfixed32> sfixed32_value_map_field = 1;", `{"sfixed32ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, sfixed32> sfixed32_value_map_field = 1;", `{"sfixed32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, sfixed64>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, sfixed64> sfixed64_value_map_field = 1;", `{"sfixed64ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, sfixed64> sfixed64_value_map_field = 1;", `{"sfixed64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`)
	})

	t.Run("map<*, bool>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, bool> bool_value_map_field = 1;", `{"boolValueMapField": {}}`)
		testJsonToMessage(t, "map<string, bool> bool_value_map_field = 1;", `{"boolValueMapField": {"a": false, "b": true}}`)
	})

	t.Run("map<*, string>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, string> string_value_map_field = 1;", `{"stringValueMapField": {}}`)
		testJsonToMessage(t, "map<string, string> string_value_map_field = 1;", `{"stringValueMapField": {"a": "b", "c": "d"}}`)
	})

	t.Run("map<*, bytes>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, bytes> bytes_value_map_field = 1;", `{"bytesValueMapField": {}}`)
		testJsonToMessage(t, "map<string, bytes> bytes_value_map_field = 1;", `{"bytesValueMapField": {"a": "", "b": "dGVzdA=="}}`) // Base64 for "test"
	})

	t.Run("map<*, float>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, float> float_value_map_field = 1;", `{"floatValueMapField": {}}`)
		testJsonToMessage(t, "map<string, float> float_value_map_field = 1;", `{"floatValueMapField": {"a": 0, "b": 3.5}}`)
	})

	t.Run("map<*, double>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, double> double_value_map_field = 1;", `{"doubleValueMapField": {}}`)
		testJsonToMessage(t, "map<string, double> double_value_map_field = 1;", `{"doubleValueMapField": {"a": 0, "b": 3.141592653589793}}`)
	})

	t.Run("map<*, sint32>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, sint32> sint32_value_map_field = 1;", `{"sint32ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, sint32> sint32_value_map_field = 1;", `{"sint32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, sint64>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, sint64> sint64_value_map_field = 1;", `{"sint64ValueMapField": {}}`)
		testJsonToMessage(t, "map<string, sint64> sint64_value_map_field = 1;", `{"sint64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`)
	})

	t.Run("map<*, EnumType>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, EnumType> enum_value_map_field = 1;", `{"enumValueMapField": {}}`)
		testJsonToMessage(t, "map<string, EnumType> enum_value_map_field = 1;", `{"enumValueMapField": {"a": "ENUM_TYPE_UNSPECIFIED", "b": "ENUM_TYPE_ONE"}}`)
	})

	t.Run("map<*, MessageType>", func(t *testing.T) {
		testJsonToMessage(t, "map<string, MessageType> message_value_map_field = 1;", `{"messageValueMapField": {}}`)
		testJsonToMessage(t, "map<string, MessageType> message_value_map_field = 1;", `{"messageValueMapField": {"a": {"value": 0}, "b": {"value": 12345}}}`)
	})
}

func TestJsonToMessageOneof(t *testing.T) {
	t.Run("oneof", func(t *testing.T) {
		testJsonToMessage(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{}`)
		testJsonToMessage(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{"int32Field": 42}`)
		testJsonToMessage(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{"stringField": "test"}`)
	})
}

func TestJsonToMessageWellKnownTypes(t *testing.T) {
	t.Run("Timestamp", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"timestampField": "1970-01-01T00:00:00Z"}`)
		testJsonToMessage(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"timestampField": "1970-01-01T00:00:01Z"}`)
		testJsonToMessage(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"timestampField": "2023-10-01T12:34:56.789Z"}`)
	})

	t.Run("Duration", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.Duration duration_field = 1;", `{"durationField": "0s"}`)
		testJsonToMessage(t, "google.protobuf.Duration duration_field = 1;", `{"durationField": "1.234s"}`)
	})

	t.Run("Struct", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.Struct struct_field = 1;", `{"structField": {}}`)
		testJsonToMessage(t, "google.protobuf.Struct struct_field = 1;", `{"structField": {"key": "value"}}`)
		testJsonToMessage(t, "google.protobuf.Struct struct_field = 1;", `{"structField": {"number": 123, "boolean": true}}`)
	})

	t.Run("ListValue", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.ListValue list_value_field = 1;", `{"listValueField": []}`)
		testJsonToMessage(t, "google.protobuf.ListValue list_value_field = 1;", `{"listValueField": ["string", 123, true]}`)
	})

	t.Run("Value", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.Value value_field = 1;", `{"valueField": {}}`)
		testJsonToMessage(t, "google.protobuf.Value value_field = 1;", `{"valueField": "string"}`)
		testJsonToMessage(t, "google.protobuf.Value value_field = 1;", `{"valueField": 123}`)
		testJsonToMessage(t, "google.protobuf.Value value_field = 1;", `{"valueField": true}`)
	})

	t.Run("Empty", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.Empty empty_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.Empty empty_field = 1;", `{"emptyField": {}}`)
	})

	t.Run("DoubleValue", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.DoubleValue double_value_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.DoubleValue double_value_field = 1;", `{"doubleValueField": 0}`)
		testJsonToMessage(t, "google.protobuf.DoubleValue double_value_field = 1;", `{"doubleValueField": 3.141592653589793}`)
	})

	t.Run("FloatValue", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.FloatValue float_value_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.FloatValue float_value_field = 1;", `{"floatValueField": 0}`)
		testJsonToMessage(t, "google.protobuf.FloatValue float_value_field = 1;", `{"floatValueField": 3.5}`)
	})

	t.Run("Int64Value", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.Int64Value int64_value_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"int64ValueField": "0"}`)
		testJsonToMessage(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"int64ValueField": "9223372036854775807"}`)
		testJsonToMessage(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"int64ValueField": "-9223372036854775808"}`)
	})

	t.Run("UInt64Value", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{"uint64ValueField": "0"}`)
		testJsonToMessage(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{"uint64ValueField": "18446744073709551615"}`)
	})

	t.Run("Int32Value", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.Int32Value int32_value_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"int32ValueField": 0}`)
		testJsonToMessage(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"int32ValueField": 2147483647}`)
		testJsonToMessage(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"int32ValueField": -2147483648}`)
	})

	t.Run("UInt32Value", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{"uint32ValueField": 0}`)
		testJsonToMessage(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{"uint32ValueField": 4294967295}`)
	})

	t.Run("BoolValue", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.BoolValue bool_value_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.BoolValue bool_value_field = 1;", `{"boolValueField": false}`)
		testJsonToMessage(t, "google.protobuf.BoolValue bool_value_field = 1;", `{"boolValueField": true}`)
	})

	t.Run("StringValue", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.StringValue string_value_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.StringValue string_value_field = 1;", `{"stringValueField": ""}`)
		testJsonToMessage(t, "google.protobuf.StringValue string_value_field = 1;", `{"stringValueField": "test"}`)
	})

	t.Run("BytesValue", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{"bytesValueField": ""}`)
		testJsonToMessage(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{"bytesValueField": "dGVzdA=="}`) // Base64 for "test"
	})

	t.Run("FieldMask", func(t *testing.T) {
		testJsonToMessage(t, "google.protobuf.FieldMask field_mask_field = 1;", `{}`)
		testJsonToMessage(t, "google.protobuf.FieldMask field_mask_field = 1;", `{"fieldMaskField": ""}`)
		testJsonToMessage(t, "google.protobuf.FieldMask field_mask_field = 1;", `{"fieldMaskField": "path1,path2"}`)
	})
}

func TestJsonToMessageNullInput(t *testing.T) {
	p := testutils.NewProtoTestSupport(t, map[string]string{
		"main.proto": `
			syntax = "proto3";
			message Test {
				int32 value = 1;
			}
		`,
	})

	typeName := ".Test"

	// Generate descriptor set JSON using descriptorsetjson package
	descriptorSetJson, err := descriptorsetjson.ToJson(p.GetFileDescriptorSet())
	g := NewWithT(t)
	g.Expect(err).NotTo(HaveOccurred())

	RunTestThatExpression(t, "pb_json_to_message(?, ?, ?, NULL, NULL)", descriptorSetJson, typeName, nil).IsNull()
}

func TestJsonToMessageNumberJsonFormat(t *testing.T) {
	// Test number JSON format (field names as numbers instead of strings)
	g := NewWithT(t)

	p := testutils.NewProtoTestSupport(t, map[string]string{
		"main.proto": `
			syntax = "proto3";
			message Test {
				int32 int32_field = 1;
				string string_field = 2;
			}
		`,
	})

	typeName := ".Test"
	descriptorSetJson, err := descriptorsetjson.ToJson(p.GetFileDescriptorSet())
	g.Expect(err).NotTo(HaveOccurred())

	// Test converting number JSON format
	numberJson := `{"1": 42, "2": "test"}`
	expectedJson := `{"int32Field": 42, "stringField": "test"}`

	// Test number JSON to message conversion
	RunTestThatExpression(t, "pb_message_to_json(?, ?, _pb_number_json_to_message(?, ?, ?, NULL), NULL, NULL)", descriptorSetJson, typeName, descriptorSetJson, typeName, numberJson).IsEqualToJsonString(expectedJson)

	// Test number JSON to wire_json conversion
	RunTestThatExpression(t, "pb_wire_json_to_json(?, ?, _pb_number_json_to_wire_json(?, ?, ?, NULL), NULL, NULL)", descriptorSetJson, typeName, descriptorSetJson, typeName, numberJson).IsEqualToJsonString(expectedJson)
}

func TestJsonToMessageEdgeCases(t *testing.T) {
	t.Run("empty message", func(t *testing.T) {
		testJsonToMessage(t, "", `{}`)
	})

	t.Run("unknown fields ignored", func(t *testing.T) {
		// SQL function should ignore unknown fields gracefully - unknown fields are dropped
		testEnumConversion(t, "int32 int32_field = 1;", `{"int32Field": 42, "unknownField": "ignored"}`, `{"int32Field": 42}`)
	})

	t.Run("enum with numeric values", func(t *testing.T) {
		// Test that numeric enum values work in addition to string names - converted to canonical string format
		testEnumConversion(t, "EnumType enum_field = 1;", `{"enumField": 0}`, `{"enumField": "ENUM_TYPE_UNSPECIFIED"}`)
		testEnumConversion(t, "EnumType enum_field = 1;", `{"enumField": 1}`, `{"enumField": "ENUM_TYPE_ONE"}`)
	})
}
