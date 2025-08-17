package main

import (
	"fmt"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	"github.com/eiiches/mysql-protobuf-functions/internal/descriptorsetjson"
	"github.com/eiiches/mysql-protobuf-functions/internal/gomega/gjson"
	"github.com/eiiches/mysql-protobuf-functions/internal/protonumberjson"
	"github.com/eiiches/mysql-protobuf-functions/internal/testutils"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func testNumberJsonToJson(t *testing.T, fieldDefinition string, numberJsonInput string, expectedJson string) {
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

	// Parse expected JSON with Go protojson to get expected result
	expectedMessage := p.JsonToDynamicMessage(typeName, expectedJson)
	expectedJsonRemarshalled, err := protojson.MarshalOptions{EmitDefaultValues: true}.Marshal(expectedMessage.Interface())
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(expectedJson).To(gjson.EqualJson(string(expectedJsonRemarshalled)), "Invalid test case; The expected JSON does not match the Go protojson output. Expected JSON is incorrect.")

	expectedJsonWithoutDefaults, err := protojson.MarshalOptions{EmitDefaultValues: false}.Marshal(expectedMessage.Interface())
	g.Expect(err).NotTo(HaveOccurred())

	// Validate that the expected JSON matches the protonumberjson output
	generatedExpectation, err := protonumberjson.Marshal(expectedMessage.Interface())
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(numberJsonInput).To(gjson.EqualJson(string(generatedExpectation)), "Invalid test case; The input value does not match the protonumberjson output. Either protonumberjson is incorrect or the input value is wrong.")

	// Test the conversion: number JSON -> JSON
	// MySQL implementation should produce the same JSON as Go's protojson
	RunTestThatExpression(t, "_pb_number_json_to_json(?, ?, ?, ?)", descriptorSetJson, typeName, numberJsonInput, true).IsEqualToJsonString(expectedJson)

	RunTestThatExpression(t, "_pb_number_json_to_json(?, ?, ?, ?)", descriptorSetJson, typeName, numberJsonInput, false).IsEqualToJsonString(string(expectedJsonWithoutDefaults))
}

func TestNumberJsonToJsonSingularFields(t *testing.T) {
	t.Run("int32", func(t *testing.T) {
		testNumberJsonToJson(t, "int32 int32_field = 1;", `{}`, `{"int32Field": 0}`)
		testNumberJsonToJson(t, "int32 int32_field = 1;", `{"1": 2147483647}`, `{"int32Field": 2147483647}`)
		testNumberJsonToJson(t, "int32 int32_field = 1;", `{"1": -2147483648}`, `{"int32Field": -2147483648}`)
	})

	t.Run("uint32", func(t *testing.T) {
		testNumberJsonToJson(t, "uint32 uint32_field = 1;", `{}`, `{"uint32Field": 0}`)
		testNumberJsonToJson(t, "uint32 uint32_field = 1;", `{"1": 4294967295}`, `{"uint32Field": 4294967295}`)
	})

	t.Run("int64", func(t *testing.T) {
		testNumberJsonToJson(t, "int64 int64_field = 1;", `{}`, `{"int64Field": "0"}`)
		testNumberJsonToJson(t, "int64 int64_field = 1;", `{"1": 9223372036854775807}`, `{"int64Field": "9223372036854775807"}`)
		testNumberJsonToJson(t, "int64 int64_field = 1;", `{"1": -9223372036854775808}`, `{"int64Field": "-9223372036854775808"}`)
	})

	t.Run("uint64", func(t *testing.T) {
		testNumberJsonToJson(t, "uint64 uint64_field = 1;", `{}`, `{"uint64Field": "0"}`)
		testNumberJsonToJson(t, "uint64 uint64_field = 1;", `{"1": 18446744073709551615}`, `{"uint64Field": "18446744073709551615"}`)
	})

	t.Run("fixed32", func(t *testing.T) {
		testNumberJsonToJson(t, "fixed32 fixed32_field = 1;", `{}`, `{"fixed32Field": 0}`)
		testNumberJsonToJson(t, "fixed32 fixed32_field = 1;", `{"1": 4294967295}`, `{"fixed32Field": 4294967295}`)
	})

	t.Run("fixed64", func(t *testing.T) {
		testNumberJsonToJson(t, "fixed64 fixed64_field = 1;", `{}`, `{"fixed64Field": "0"}`)
		testNumberJsonToJson(t, "fixed64 fixed64_field = 1;", `{"1": 18446744073709551615}`, `{"fixed64Field": "18446744073709551615"}`)
	})

	t.Run("sfixed32", func(t *testing.T) {
		testNumberJsonToJson(t, "sfixed32 sfixed32_field = 1;", `{}`, `{"sfixed32Field": 0}`)
		testNumberJsonToJson(t, "sfixed32 sfixed32_field = 1;", `{"1": 2147483647}`, `{"sfixed32Field": 2147483647}`)
		testNumberJsonToJson(t, "sfixed32 sfixed32_field = 1;", `{"1": -2147483648}`, `{"sfixed32Field": -2147483648}`)
	})

	t.Run("sfixed64", func(t *testing.T) {
		testNumberJsonToJson(t, "sfixed64 sfixed64_field = 1;", `{}`, `{"sfixed64Field": "0"}`)
		testNumberJsonToJson(t, "sfixed64 sfixed64_field = 1;", `{"1": 9223372036854775807}`, `{"sfixed64Field": "9223372036854775807"}`)
		testNumberJsonToJson(t, "sfixed64 sfixed64_field = 1;", `{"1": -9223372036854775808}`, `{"sfixed64Field": "-9223372036854775808"}`)
	})

	t.Run("bool", func(t *testing.T) {
		testNumberJsonToJson(t, "bool bool_field = 1;", `{}`, `{"boolField": false}`)
		testNumberJsonToJson(t, "bool bool_field = 1;", `{"1": true}`, `{"boolField": true}`)
	})

	t.Run("string", func(t *testing.T) {
		testNumberJsonToJson(t, "string string_field = 1;", `{}`, `{"stringField": ""}`)
		testNumberJsonToJson(t, "string string_field = 1;", `{"1": "test"}`, `{"stringField": "test"}`)
	})

	t.Run("bytes", func(t *testing.T) {
		testNumberJsonToJson(t, "bytes bytes_field = 1;", `{}`, `{"bytesField": ""}`)
		testNumberJsonToJson(t, "bytes bytes_field = 1;", `{"1": "aGVsbG8="}`, `{"bytesField": "aGVsbG8="}`)
	})

	t.Run("float", func(t *testing.T) {
		testNumberJsonToJson(t, "float float_field = 1;", `{}`, `{"floatField": 0}`)
		testNumberJsonToJson(t, "float float_field = 1;", `{"1": 3.5}`, `{"floatField": 3.5}`)
	})

	t.Run("double", func(t *testing.T) {
		testNumberJsonToJson(t, "double double_field = 1;", `{}`, `{"doubleField": 0}`)
		testNumberJsonToJson(t, "double double_field = 1;", `{"1": 3.141592653589793}`, `{"doubleField": 3.141592653589793}`)
	})

	t.Run("sint32", func(t *testing.T) {
		testNumberJsonToJson(t, "sint32 sint32_field = 1;", `{}`, `{"sint32Field": 0}`)
		testNumberJsonToJson(t, "sint32 sint32_field = 1;", `{"1": 2147483647}`, `{"sint32Field": 2147483647}`)
		testNumberJsonToJson(t, "sint32 sint32_field = 1;", `{"1": -2147483648}`, `{"sint32Field": -2147483648}`)
	})

	t.Run("sint64", func(t *testing.T) {
		testNumberJsonToJson(t, "sint64 sint64_field = 1;", `{}`, `{"sint64Field": "0"}`)
		testNumberJsonToJson(t, "sint64 sint64_field = 1;", `{"1": 9223372036854775807}`, `{"sint64Field": "9223372036854775807"}`)
		testNumberJsonToJson(t, "sint64 sint64_field = 1;", `{"1": -9223372036854775808}`, `{"sint64Field": "-9223372036854775808"}`)
	})

	t.Run("enum", func(t *testing.T) {
		testNumberJsonToJson(t, "EnumType enum_field = 1;", `{}`, `{"enumField": "ENUM_TYPE_UNSPECIFIED"}`)
		testNumberJsonToJson(t, "EnumType enum_field = 1;", `{"1": 1}`, `{"enumField": "ENUM_TYPE_ONE"}`)
	})

	t.Run("message", func(t *testing.T) {
		testNumberJsonToJson(t, "MessageType message_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "MessageType message_field = 1;", `{"1": {}}`, `{"messageField": {"value": 0}}`)
		testNumberJsonToJson(t, "MessageType message_field = 1;", `{"1": {"1": 12345}}`, `{"messageField": {"value": 12345}}`)
	})
}

func TestNumberJsonToJsonRepeatedFields(t *testing.T) {
	t.Run("repeated int32", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated int32 repeated_int32_field = 1;", `{}`, `{"repeatedInt32Field": []}`)
		testNumberJsonToJson(t, "repeated int32 repeated_int32_field = 1;", `{"1": [0]}`, `{"repeatedInt32Field": [0]}`)
		testNumberJsonToJson(t, "repeated int32 repeated_int32_field = 1;", `{"1": [-2147483648, 0, 2147483647]}`, `{"repeatedInt32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated uint32", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated uint32 repeated_uint32_field = 1;", `{}`, `{"repeatedUint32Field": []}`)
		testNumberJsonToJson(t, "repeated uint32 repeated_uint32_field = 1;", `{"1": [0]}`, `{"repeatedUint32Field": [0]}`)
		testNumberJsonToJson(t, "repeated uint32 repeated_uint32_field = 1;", `{"1": [0, 4294967295]}`, `{"repeatedUint32Field": [0, 4294967295]}`)
	})

	t.Run("repeated int64", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated int64 repeated_int64_field = 1;", `{}`, `{"repeatedInt64Field": []}`)
		testNumberJsonToJson(t, "repeated int64 repeated_int64_field = 1;", `{"1": [0]}`, `{"repeatedInt64Field": ["0"]}`)
		testNumberJsonToJson(t, "repeated int64 repeated_int64_field = 1;", `{"1": [-9223372036854775808, 0, 9223372036854775807]}`, `{"repeatedInt64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`)
	})

	t.Run("repeated uint64", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated uint64 repeated_uint64_field = 1;", `{}`, `{"repeatedUint64Field": []}`)
		testNumberJsonToJson(t, "repeated uint64 repeated_uint64_field = 1;", `{"1": [0]}`, `{"repeatedUint64Field": ["0"]}`)
		testNumberJsonToJson(t, "repeated uint64 repeated_uint64_field = 1;", `{"1": [0, 18446744073709551615]}`, `{"repeatedUint64Field": ["0", "18446744073709551615"]}`)
	})

	t.Run("repeated fixed32", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated fixed32 repeated_fixed32_field = 1;", `{}`, `{"repeatedFixed32Field": []}`)
		testNumberJsonToJson(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"1": [0]}`, `{"repeatedFixed32Field": [0]}`)
		testNumberJsonToJson(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"1": [0, 4294967295]}`, `{"repeatedFixed32Field": [0, 4294967295]}`)
	})

	t.Run("repeated fixed64", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated fixed64 repeated_fixed64_field = 1;", `{}`, `{"repeatedFixed64Field": []}`)
		testNumberJsonToJson(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"1": [0]}`, `{"repeatedFixed64Field": ["0"]}`)
		testNumberJsonToJson(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"1": [0, 18446744073709551615]}`, `{"repeatedFixed64Field": ["0", "18446744073709551615"]}`)
	})

	t.Run("repeated sfixed32", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{}`, `{"repeatedSfixed32Field": []}`)
		testNumberJsonToJson(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"1": [0]}`, `{"repeatedSfixed32Field": [0]}`)
		testNumberJsonToJson(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"1": [-2147483648, 0, 2147483647]}`, `{"repeatedSfixed32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated sfixed64", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{}`, `{"repeatedSfixed64Field": []}`)
		testNumberJsonToJson(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"1": [0]}`, `{"repeatedSfixed64Field": ["0"]}`)
		testNumberJsonToJson(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"1": [-9223372036854775808, 0, 9223372036854775807]}`, `{"repeatedSfixed64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`)
	})

	t.Run("repeated bool", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated bool repeated_bool_field = 1;", `{}`, `{"repeatedBoolField": []}`)
		testNumberJsonToJson(t, "repeated bool repeated_bool_field = 1;", `{"1": [false]}`, `{"repeatedBoolField": [false]}`)
		testNumberJsonToJson(t, "repeated bool repeated_bool_field = 1;", `{"1": [true, false]}`, `{"repeatedBoolField": [true, false]}`)
	})

	t.Run("repeated string", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated string repeated_string_field = 1;", `{}`, `{"repeatedStringField": []}`)
		testNumberJsonToJson(t, "repeated string repeated_string_field = 1;", `{"1": [""]}`, `{"repeatedStringField": [""]}`)
		testNumberJsonToJson(t, "repeated string repeated_string_field = 1;", `{"1": ["test", ""]}`, `{"repeatedStringField": ["test", ""]}`)
	})

	t.Run("repeated bytes", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated bytes repeated_bytes_field = 1;", `{}`, `{"repeatedBytesField": []}`)
		testNumberJsonToJson(t, "repeated bytes repeated_bytes_field = 1;", `{"1": [""]}`, `{"repeatedBytesField": [""]}`)
		testNumberJsonToJson(t, "repeated bytes repeated_bytes_field = 1;", `{"1": ["aGVsbG8=", ""]}`, `{"repeatedBytesField": ["aGVsbG8=", ""]}`)
	})

	t.Run("repeated float", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated float repeated_float_field = 1;", `{}`, `{"repeatedFloatField": []}`)
		testNumberJsonToJson(t, "repeated float repeated_float_field = 1;", `{"1": [0]}`, `{"repeatedFloatField": [0]}`)
		testNumberJsonToJson(t, "repeated float repeated_float_field = 1;", `{"1": [3.5, 0]}`, `{"repeatedFloatField": [3.5, 0]}`)
	})

	t.Run("repeated double", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated double repeated_double_field = 1;", `{}`, `{"repeatedDoubleField": []}`)
		testNumberJsonToJson(t, "repeated double repeated_double_field = 1;", `{"1": [0]}`, `{"repeatedDoubleField": [0]}`)
		testNumberJsonToJson(t, "repeated double repeated_double_field = 1;", `{"1": [3.5]}`, `{"repeatedDoubleField": [3.5]}`)
		testNumberJsonToJson(t, "repeated double repeated_double_field = 1;", `{"1": [-1.7976931348623157e+308, 0, 1.7976931348623157e+308]}`, `{"repeatedDoubleField": [-1.7976931348623157e+308, 0, 1.7976931348623157e+308]}`)
	})

	t.Run("repeated sint32", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated sint32 repeated_sint32_field = 1;", `{}`, `{"repeatedSint32Field": []}`)
		testNumberJsonToJson(t, "repeated sint32 repeated_sint32_field = 1;", `{"1": [0]}`, `{"repeatedSint32Field": [0]}`)
		testNumberJsonToJson(t, "repeated sint32 repeated_sint32_field = 1;", `{"1": [-2147483648, 0, 2147483647]}`, `{"repeatedSint32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated sint64", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated sint64 repeated_sint64_field = 1;", `{}`, `{"repeatedSint64Field": []}`)
		testNumberJsonToJson(t, "repeated sint64 repeated_sint64_field = 1;", `{"1": [0]}`, `{"repeatedSint64Field": ["0"]}`)
		testNumberJsonToJson(t, "repeated sint64 repeated_sint64_field = 1;", `{"1": [-9223372036854775808, 0, 9223372036854775807]}`, `{"repeatedSint64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`)
	})

	t.Run("repeated enum", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated EnumType repeated_enum_field = 1;", `{}`, `{"repeatedEnumField": []}`)
		testNumberJsonToJson(t, "repeated EnumType repeated_enum_field = 1;", `{"1": [0]}`, `{"repeatedEnumField": ["ENUM_TYPE_UNSPECIFIED"]}`)
		testNumberJsonToJson(t, "repeated EnumType repeated_enum_field = 1;", `{"1": [1, 0]}`, `{"repeatedEnumField": ["ENUM_TYPE_ONE", "ENUM_TYPE_UNSPECIFIED"]}`)
	})

	t.Run("repeated message", func(t *testing.T) {
		testNumberJsonToJson(t, "repeated MessageType repeated_message_field = 1;", `{}`, `{"repeatedMessageField": []}`)
		testNumberJsonToJson(t, "repeated MessageType repeated_message_field = 1;", `{"1": [{}]}`, `{"repeatedMessageField": [{"value": 0}]}`)
		testNumberJsonToJson(t, "repeated MessageType repeated_message_field = 1;", `{"1": [{"1": 12345}, {"1": 67890}]}`, `{"repeatedMessageField": [{"value": 12345}, {"value": 67890}]}`)
	})
}

func TestNumberJsonToJsonOptionalFields(t *testing.T) {
	t.Run("optional int32", func(t *testing.T) {
		testNumberJsonToJson(t, "optional int32 optional_int32_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional int32 optional_int32_field = 1;", `{"1": 0}`, `{"optionalInt32Field": 0}`)
		testNumberJsonToJson(t, "optional int32 optional_int32_field = 1;", `{"1": 2147483647}`, `{"optionalInt32Field": 2147483647}`)
		testNumberJsonToJson(t, "optional int32 optional_int32_field = 1;", `{"1": -2147483648}`, `{"optionalInt32Field": -2147483648}`)
	})

	t.Run("optional uint32", func(t *testing.T) {
		testNumberJsonToJson(t, "optional uint32 optional_uint32_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional uint32 optional_uint32_field = 1;", `{"1": 0}`, `{"optionalUint32Field": 0}`)
		testNumberJsonToJson(t, "optional uint32 optional_uint32_field = 1;", `{"1": 4294967295}`, `{"optionalUint32Field": 4294967295}`)
	})

	t.Run("optional int64", func(t *testing.T) {
		testNumberJsonToJson(t, "optional int64 optional_int64_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional int64 optional_int64_field = 1;", `{"1": 0}`, `{"optionalInt64Field": "0"}`)
		testNumberJsonToJson(t, "optional int64 optional_int64_field = 1;", `{"1": 9223372036854775807}`, `{"optionalInt64Field": "9223372036854775807"}`)
		testNumberJsonToJson(t, "optional int64 optional_int64_field = 1;", `{"1": -9223372036854775808}`, `{"optionalInt64Field": "-9223372036854775808"}`)
	})

	t.Run("optional uint64", func(t *testing.T) {
		testNumberJsonToJson(t, "optional uint64 optional_uint64_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional uint64 optional_uint64_field = 1;", `{"1": 0}`, `{"optionalUint64Field": "0"}`)
		testNumberJsonToJson(t, "optional uint64 optional_uint64_field = 1;", `{"1": 18446744073709551615}`, `{"optionalUint64Field": "18446744073709551615"}`)
	})

	t.Run("optional fixed32", func(t *testing.T) {
		testNumberJsonToJson(t, "optional fixed32 optional_fixed32_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional fixed32 optional_fixed32_field = 1;", `{"1": 0}`, `{"optionalFixed32Field": 0}`)
		testNumberJsonToJson(t, "optional fixed32 optional_fixed32_field = 1;", `{"1": 4294967295}`, `{"optionalFixed32Field": 4294967295}`)
	})

	t.Run("optional fixed64", func(t *testing.T) {
		testNumberJsonToJson(t, "optional fixed64 optional_fixed64_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional fixed64 optional_fixed64_field = 1;", `{"1": 0}`, `{"optionalFixed64Field": "0"}`)
		testNumberJsonToJson(t, "optional fixed64 optional_fixed64_field = 1;", `{"1": 18446744073709551615}`, `{"optionalFixed64Field": "18446744073709551615"}`)
	})

	t.Run("optional sfixed32", func(t *testing.T) {
		testNumberJsonToJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"1": 0}`, `{"optionalSfixed32Field": 0}`)
		testNumberJsonToJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"1": 2147483647}`, `{"optionalSfixed32Field": 2147483647}`)
		testNumberJsonToJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"1": -2147483648}`, `{"optionalSfixed32Field": -2147483648}`)
	})

	t.Run("optional sfixed64", func(t *testing.T) {
		testNumberJsonToJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"1": 0}`, `{"optionalSfixed64Field": "0"}`)
		testNumberJsonToJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"1": 9223372036854775807}`, `{"optionalSfixed64Field": "9223372036854775807"}`)
		testNumberJsonToJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"1": -9223372036854775808}`, `{"optionalSfixed64Field": "-9223372036854775808"}`)
	})

	t.Run("optional bool", func(t *testing.T) {
		testNumberJsonToJson(t, "optional bool optional_bool_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional bool optional_bool_field = 1;", `{"1": false}`, `{"optionalBoolField": false}`)
		testNumberJsonToJson(t, "optional bool optional_bool_field = 1;", `{"1": true}`, `{"optionalBoolField": true}`)
	})

	t.Run("optional string", func(t *testing.T) {
		testNumberJsonToJson(t, "optional string optional_string_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional string optional_string_field = 1;", `{"1": ""}`, `{"optionalStringField": ""}`)
		testNumberJsonToJson(t, "optional string optional_string_field = 1;", `{"1": "testNumberJsonToJson"}`, `{"optionalStringField": "testNumberJsonToJson"}`)
	})

	t.Run("optional bytes", func(t *testing.T) {
		testNumberJsonToJson(t, "optional bytes optional_bytes_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional bytes optional_bytes_field = 1;", `{"1": ""}`, `{"optionalBytesField": ""}`)
		testNumberJsonToJson(t, "optional bytes optional_bytes_field = 1;", `{"1": "aGVsbG8="}`, `{"optionalBytesField": "aGVsbG8="}`)
	})

	t.Run("optional float", func(t *testing.T) {
		testNumberJsonToJson(t, "optional float optional_float_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional float optional_float_field = 1;", `{"1": 0}`, `{"optionalFloatField": 0}`)
		testNumberJsonToJson(t, "optional float optional_float_field = 1;", `{"1": 3.5}`, `{"optionalFloatField": 3.5}`)
	})

	t.Run("optional double", func(t *testing.T) {
		testNumberJsonToJson(t, "optional double optional_double_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional double optional_double_field = 1;", `{"1": 0}`, `{"optionalDoubleField": 0}`)
		testNumberJsonToJson(t, "optional double optional_double_field = 1;", `{"1": 3.141592653589793}`, `{"optionalDoubleField": 3.141592653589793}`)
	})

	t.Run("optional sint32", func(t *testing.T) {
		testNumberJsonToJson(t, "optional sint32 optional_sint32_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional sint32 optional_sint32_field = 1;", `{"1": 0}`, `{"optionalSint32Field": 0}`)
		testNumberJsonToJson(t, "optional sint32 optional_sint32_field = 1;", `{"1": 2147483647}`, `{"optionalSint32Field": 2147483647}`)
		testNumberJsonToJson(t, "optional sint32 optional_sint32_field = 1;", `{"1": -2147483648}`, `{"optionalSint32Field": -2147483648}`)
	})

	t.Run("optional sint64", func(t *testing.T) {
		testNumberJsonToJson(t, "optional sint64 optional_sint64_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional sint64 optional_sint64_field = 1;", `{"1": 0}`, `{"optionalSint64Field": "0"}`)
		testNumberJsonToJson(t, "optional sint64 optional_sint64_field = 1;", `{"1": 9223372036854775807}`, `{"optionalSint64Field": "9223372036854775807"}`)
		testNumberJsonToJson(t, "optional sint64 optional_sint64_field = 1;", `{"1": -9223372036854775808}`, `{"optionalSint64Field": "-9223372036854775808"}`)
	})

	t.Run("optional enum", func(t *testing.T) {
		testNumberJsonToJson(t, "optional EnumType optional_enum_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional EnumType optional_enum_field = 1;", `{"1": 0}`, `{"optionalEnumField": "ENUM_TYPE_UNSPECIFIED"}`)
		testNumberJsonToJson(t, "optional EnumType optional_enum_field = 1;", `{"1": 1}`, `{"optionalEnumField": "ENUM_TYPE_ONE"}`)
	})

	t.Run("optional message", func(t *testing.T) {
		testNumberJsonToJson(t, "optional MessageType optional_message_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "optional MessageType optional_message_field = 1;", `{"1": {}}`, `{"optionalMessageField": {"value": 0}}`)
		testNumberJsonToJson(t, "optional MessageType optional_message_field = 1;", `{"1": {"1": 12345}}`, `{"optionalMessageField": {"value": 12345}}`)
	})
}

func TestNumberJsonToJsonMapKey(t *testing.T) {
	t.Run("map<int32, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<int32, string> int32_key_map_field = 1;", `{}`, `{"int32KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<int32, string> int32_key_map_field = 1;", `{"1": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`, `{"int32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<uint32, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<uint32, string> uint32_key_map_field = 1;", `{}`, `{"uint32KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<uint32, string> uint32_key_map_field = 1;", `{"1": {"0": "a", "4294967295": "b"}}`, `{"uint32KeyMapField": {"0": "a", "4294967295": "b"}}`)
	})

	t.Run("map<int64, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<int64, string> int64_key_map_field = 1;", `{}`, `{"int64KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<int64, string> int64_key_map_field = 1;", `{"1": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`, `{"int64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	t.Run("map<uint64, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<uint64, string> uint64_key_map_field = 1;", `{}`, `{"uint64KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<uint64, string> uint64_key_map_field = 1;", `{"1": {"0": "a", "18446744073709551615": "b"}}`, `{"uint64KeyMapField": {"0": "a", "18446744073709551615": "b"}}`)
	})

	t.Run("map<fixed32, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<fixed32, string> fixed32_key_map_field = 1;", `{}`, `{"fixed32KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<fixed32, string> fixed32_key_map_field = 1;", `{"1": {"0": "a", "4294967295": "b"}}`, `{"fixed32KeyMapField": {"0": "a", "4294967295": "b"}}`)
	})

	t.Run("map<fixed64, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<fixed64, string> fixed64_key_map_field = 1;", `{}`, `{"fixed64KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<fixed64, string> fixed64_key_map_field = 1;", `{"1": {"0": "a", "18446744073709551615": "b"}}`, `{"fixed64KeyMapField": {"0": "a", "18446744073709551615": "b"}}`)
	})

	t.Run("map<sfixed32, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<sfixed32, string> sfixed32_key_map_field = 1;", `{}`, `{"sfixed32KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<sfixed32, string> sfixed32_key_map_field = 1;", `{"1": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`, `{"sfixed32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<sfixed64, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<sfixed64, string> sfixed64_key_map_field = 1;", `{}`, `{"sfixed64KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<sfixed64, string> sfixed64_key_map_field = 1;", `{"1": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`, `{"sfixed64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	t.Run("map<bool, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<bool, string> bool_key_map_field = 1;", `{}`, `{"boolKeyMapField": {}}`)
		testNumberJsonToJson(t, "map<bool, string> bool_key_map_field = 1;", `{"1": {"false": "a", "true": "b"}}`, `{"boolKeyMapField": {"false": "a", "true": "b"}}`)
	})

	t.Run("map<string, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, string> string_key_map_field = 1;", `{}`, `{"stringKeyMapField": {}}`)
		testNumberJsonToJson(t, "map<string, string> string_key_map_field = 1;", `{"1": {"a": "b", "c": "d"}}`, `{"stringKeyMapField": {"a": "b", "c": "d"}}`)
	})

	t.Run("map<sint32, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<sint32, string> sint32_key_map_field = 1;", `{}`, `{"sint32KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<sint32, string> sint32_key_map_field = 1;", `{"1": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`, `{"sint32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<sint64, *>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<sint64, string> sint64_key_map_field = 1;", `{}`, `{"sint64KeyMapField": {}}`)
		testNumberJsonToJson(t, "map<sint64, string> sint64_key_map_field = 1;", `{"1": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`, `{"sint64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	// NOTE: float, double, enum, or message cannot be a map key.
}

func TestNumberJsonToJsonMapValue(t *testing.T) {
	t.Run("map<*, int32>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, int32> int32_value_map_field = 1;", `{}`, `{"int32ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, int32> int32_value_map_field = 1;", `{"1": {"a": 0, "b": 2147483647, "c": -2147483648}}`, `{"int32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, uint32>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, uint32> uint32_value_map_field = 1;", `{}`, `{"uint32ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, uint32> uint32_value_map_field = 1;", `{"1": {"a": 0, "b": 4294967295}}`, `{"uint32ValueMapField": {"a": 0, "b": 4294967295}}`)
	})

	t.Run("map<*, int64>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, int64> int64_value_map_field = 1;", `{}`, `{"int64ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, int64> int64_value_map_field = 1;", `{"1": {"a": 0, "b": 9223372036854775807, "c": -9223372036854775808}}`, `{"int64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`)
	})

	t.Run("map<*, uint64>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, uint64> uint64_value_map_field = 1;", `{}`, `{"uint64ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, uint64> uint64_value_map_field = 1;", `{"1": {"a": 0, "b": 18446744073709551615}}`, `{"uint64ValueMapField": {"a": "0", "b": "18446744073709551615"}}`)
	})

	t.Run("map<*, fixed32>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, fixed32> fixed32_value_map_field = 1;", `{}`, `{"fixed32ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, fixed32> fixed32_value_map_field = 1;", `{"1": {"a": 0, "b": 4294967295}}`, `{"fixed32ValueMapField": {"a": 0, "b": 4294967295}}`)
	})

	t.Run("map<*, fixed64>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, fixed64> fixed64_value_map_field = 1;", `{}`, `{"fixed64ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, fixed64> fixed64_value_map_field = 1;", `{"1": {"a": 0, "b": 18446744073709551615}}`, `{"fixed64ValueMapField": {"a": "0", "b": "18446744073709551615"}}`)
	})

	t.Run("map<*, sfixed32>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, sfixed32> sfixed32_value_map_field = 1;", `{}`, `{"sfixed32ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, sfixed32> sfixed32_value_map_field = 1;", `{"1": {"a": 0, "b": 2147483647, "c": -2147483648}}`, `{"sfixed32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, sfixed64>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, sfixed64> sfixed64_value_map_field = 1;", `{}`, `{"sfixed64ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, sfixed64> sfixed64_value_map_field = 1;", `{"1": {"a": 0, "b": 9223372036854775807, "c": -9223372036854775808}}`, `{"sfixed64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`)
	})

	t.Run("map<*, bool>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, bool> bool_value_map_field = 1;", `{}`, `{"boolValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, bool> bool_value_map_field = 1;", `{"1": {"a": false, "b": true}}`, `{"boolValueMapField": {"a": false, "b": true}}`)
	})

	t.Run("map<*, string>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, string> string_value_map_field = 1;", `{}`, `{"stringValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, string> string_value_map_field = 1;", `{"1": {"a": "b", "c": "d"}}`, `{"stringValueMapField": {"a": "b", "c": "d"}}`)
	})

	t.Run("map<*, bytes>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, bytes> bytes_value_map_field = 1;", `{}`, `{"bytesValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, bytes> bytes_value_map_field = 1;", `{"1": {"a": "", "b": "dGVzdA=="}}`, `{"bytesValueMapField": {"a": "", "b": "dGVzdA=="}}`) // Base64 for "test"
	})

	t.Run("map<*, float>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, float> float_value_map_field = 1;", `{}`, `{"floatValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, float> float_value_map_field = 1;", `{"1": {"a": 0, "b": 3.5}}`, `{"floatValueMapField": {"a": 0, "b": 3.5}}`)
	})

	t.Run("map<*, double>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, double> double_value_map_field = 1;", `{}`, `{"doubleValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, double> double_value_map_field = 1;", `{"1": {"a": 0, "b": 3.141592653589793}}`, `{"doubleValueMapField": {"a": 0, "b": 3.141592653589793}}`)
	})

	t.Run("map<*, sint32>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, sint32> sint32_value_map_field = 1;", `{}`, `{"sint32ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, sint32> sint32_value_map_field = 1;", `{"1": {"a": 0, "b": 2147483647, "c": -2147483648}}`, `{"sint32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, sint64>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, sint64> sint64_value_map_field = 1;", `{}`, `{"sint64ValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, sint64> sint64_value_map_field = 1;", `{"1": {"a": 0, "b": 9223372036854775807, "c": -9223372036854775808}}`, `{"sint64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`)
	})

	t.Run("map<*, EnumType>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, EnumType> enum_value_map_field = 1;", `{}`, `{"enumValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, EnumType> enum_value_map_field = 1;", `{"1": {"a": 0, "b": 1}}`, `{"enumValueMapField": {"a": "ENUM_TYPE_UNSPECIFIED", "b": "ENUM_TYPE_ONE"}}`)
	})

	t.Run("map<*, MessageType>", func(t *testing.T) {
		testNumberJsonToJson(t, "map<string, MessageType> message_value_map_field = 1;", `{}`, `{"messageValueMapField": {}}`)
		testNumberJsonToJson(t, "map<string, MessageType> message_value_map_field = 1;", `{"1": {"a": {}, "b": {"1": 12345}}}`, `{"messageValueMapField": {"a": {"value": 0}, "b": {"value": 12345}}}`)
	})
}

func TestNumberJsonToJsonOneof(t *testing.T) {
	t.Run("oneof", func(t *testing.T) {
		testNumberJsonToJson(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{}`, `{}`)
		testNumberJsonToJson(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{"1": 42}`, `{"int32Field": 42}`)
		testNumberJsonToJson(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{"2": "test"}`, `{"stringField": "test"}`)
	})

	t.Run("oneof with message", func(t *testing.T) {
		testNumberJsonToJson(t, "oneof kind { int32 int32_field = 1; MessageType message_field = 2; }", `{}`, `{}`)
		testNumberJsonToJson(t, "oneof kind { int32 int32_field = 1; MessageType message_field = 2; }", `{"1": 42}`, `{"int32Field": 42}`)
		testNumberJsonToJson(t, "oneof kind { int32 int32_field = 1; MessageType message_field = 2; }", `{"2": {"1": 123}}`, `{"messageField": {"value": 123}}`)
	})

	t.Run("oneof with enum", func(t *testing.T) {
		testNumberJsonToJson(t, "oneof kind { int32 int32_field = 1; EnumType enum_field = 2; }", `{}`, `{}`)
		testNumberJsonToJson(t, "oneof kind { int32 int32_field = 1; EnumType enum_field = 2; }", `{"1": 42}`, `{"int32Field": 42}`)
		testNumberJsonToJson(t, "oneof kind { int32 int32_field = 1; EnumType enum_field = 2; }", `{"2": 1}`, `{"enumField": "ENUM_TYPE_ONE"}`)
	})
}

func TestNumberJsonToJsonWellKnownTypes(t *testing.T) {
	t.Run("Timestamp", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"1": {}}`, `{"timestampField": "1970-01-01T00:00:00Z"}`)
		testNumberJsonToJson(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"1": {"1": 1}}`, `{"timestampField": "1970-01-01T00:00:01Z"}`)
		testNumberJsonToJson(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"1": {"1": 1696163696, "2": 789000000}}`, `{"timestampField": "2023-10-01T12:34:56.789Z"}`)
	})

	t.Run("Duration", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.Duration duration_field = 1;", `{"1": {}}`, `{"durationField": "0s"}`)
		testNumberJsonToJson(t, "google.protobuf.Duration duration_field = 1;", `{"1": {"1": 1, "2": 234000000}}`, `{"durationField": "1.234s"}`)
	})

	t.Run("Struct", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.Struct struct_field = 1;", `{"1": {}}`, `{"structField": {}}`)
		testNumberJsonToJson(t, "google.protobuf.Struct struct_field = 1;", `{"1": {"1": {"key": {"3": "value"}}}}`, `{"structField": {"key": "value"}}`)
		testNumberJsonToJson(t, "google.protobuf.Struct struct_field = 1;", `{"1": {"1": {"number": {"2": 123}, "boolean": {"4": true}}}}`, `{"structField": {"number": 123, "boolean": true}}`)
	})

	t.Run("ListValue", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.ListValue list_value_field = 1;", `{"1": {}}`, `{"listValueField": []}`)
		testNumberJsonToJson(t, "google.protobuf.ListValue list_value_field = 1;", `{"1": {"1": [{"3": "string"}, {"2": 123}, {"4": true}]}}`, `{"listValueField": ["string", 123, true]}`)
	})

	t.Run("Value", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.Value value_field = 1;", `{"1": {"5": {}}}`, `{"valueField": {}}`)
		testNumberJsonToJson(t, "google.protobuf.Value value_field = 1;", `{"1": {"3": "string"}}`, `{"valueField": "string"}`)
		testNumberJsonToJson(t, "google.protobuf.Value value_field = 1;", `{"1": {"2": 123}}`, `{"valueField": 123}`)
		testNumberJsonToJson(t, "google.protobuf.Value value_field = 1;", `{"1": {"4": true}}`, `{"valueField": true}`)
	})

	t.Run("Empty", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.Empty empty_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.Empty empty_field = 1;", `{"1": {}}`, `{"emptyField": {}}`)
	})

	t.Run("DoubleValue", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.DoubleValue double_value_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.DoubleValue double_value_field = 1;", `{"1": {}}`, `{"doubleValueField": 0}`)
		testNumberJsonToJson(t, "google.protobuf.DoubleValue double_value_field = 1;", `{"1": {"1": 1.7976931348623157e+308}}`, `{"doubleValueField": 1.7976931348623157e+308}`)
	})

	t.Run("FloatValue", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.FloatValue float_value_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.FloatValue float_value_field = 1;", `{"1": {}}`, `{"floatValueField": 0}`)
		testNumberJsonToJson(t, "google.protobuf.FloatValue float_value_field = 1;", `{"1": {"1": 3.5}}`, `{"floatValueField": 3.5}`)
	})

	t.Run("Int64Value", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"1": {}}`, `{"int64ValueField": "0"}`)
		testNumberJsonToJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"1": {"1": 9223372036854775807}}`, `{"int64ValueField": "9223372036854775807"}`)
		testNumberJsonToJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"1": {"1": -9223372036854775808}}`, `{"int64ValueField": "-9223372036854775808"}`)
	})

	t.Run("UInt64Value", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{"1": {}}`, `{"uint64ValueField": "0"}`)
		testNumberJsonToJson(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{"1": {"1": 18446744073709551615}}`, `{"uint64ValueField": "18446744073709551615"}`)
	})

	t.Run("Int32Value", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"1": {}}`, `{"int32ValueField": 0}`)
		testNumberJsonToJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"1": {"1": 2147483647}}`, `{"int32ValueField": 2147483647}`)
		testNumberJsonToJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"1": {"1": -2147483648}}`, `{"int32ValueField": -2147483648}`)
	})

	t.Run("UInt32Value", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{"1": {}}`, `{"uint32ValueField": 0}`)
		testNumberJsonToJson(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{"1": {"1": 4294967295}}`, `{"uint32ValueField": 4294967295}`)
	})

	t.Run("BoolValue", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.BoolValue bool_value_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.BoolValue bool_value_field = 1;", `{"1": {}}`, `{"boolValueField": false}`)
		testNumberJsonToJson(t, "google.protobuf.BoolValue bool_value_field = 1;", `{"1": {"1": true}}`, `{"boolValueField": true}`)
	})

	t.Run("StringValue", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.StringValue string_value_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.StringValue string_value_field = 1;", `{"1": {}}`, `{"stringValueField": ""}`)
		testNumberJsonToJson(t, "google.protobuf.StringValue string_value_field = 1;", `{"1": {"1": "test"}}`, `{"stringValueField": "test"}`)
	})

	t.Run("BytesValue", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{"1": {}}`, `{"bytesValueField": ""}`)
		testNumberJsonToJson(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{"1": {"1": "aGVsbG8="}}`, `{"bytesValueField": "aGVsbG8="}`)
	})

	t.Run("FieldMask", func(t *testing.T) {
		testNumberJsonToJson(t, "google.protobuf.FieldMask field_mask_field = 1;", `{}`, `{}`)
		testNumberJsonToJson(t, "google.protobuf.FieldMask field_mask_field = 1;", `{"1": {}}`, `{"fieldMaskField": ""}`)
		testNumberJsonToJson(t, "google.protobuf.FieldMask field_mask_field = 1;", `{"1": {"1": ["path1", "path2"]}}`, `{"fieldMaskField": "path1,path2"}`)
	})
}

func TestNumberJsonToJsonNullInput(t *testing.T) {
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

	RunTestThatExpression(t, "_pb_number_json_to_json(?, ?, ?, ?)", descriptorSetJson, typeName, nil, true).IsNull()
}
