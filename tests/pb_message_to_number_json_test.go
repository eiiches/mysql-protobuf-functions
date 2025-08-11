package main

import (
	"fmt"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/gomega/gjson"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	"github.com/eiiches/mysql-protobuf-functions/internal/descriptorsetjson"
	"github.com/eiiches/mysql-protobuf-functions/internal/protonumberjson"
	"github.com/eiiches/mysql-protobuf-functions/internal/testutils"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func testMessageToNumberJson(t *testing.T, fieldDefinition string, input string, expectedNumberJson string) {
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

	dynamicMessage := p.JsonToDynamicMessage(typeName, input)
	serializedBinary := p.JsonToProtobuf(typeName, input)

	// Validate that the expected JSON matches the protonumberjson output
	generatedExpectation, err := protonumberjson.Marshal(dynamicMessage.Interface())
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(expectedNumberJson).To(gjson.EqualJson(string(generatedExpectation)), "Invalid test case; The expected value does not match the protonumberjson output. Either protonumberjson is incorrect or the expected value is wrong.")

	// Test against explicit expected value
	RunTestThatExpression(t, "_pb_message_to_number_json(?, ?, ?)", descriptorSetJson, typeName, serializedBinary).IsEqualToJsonString(expectedNumberJson)
}

func TestMessageToNumberJsonSingularFields(t *testing.T) {
	t.Run("int32", func(t *testing.T) {
		testMessageToNumberJson(t, "int32 int32_field = 1;", `{"int32Field": 0}`, `{}`)
		testMessageToNumberJson(t, "int32 int32_field = 1;", `{"int32Field": 2147483647}`, `{"1": 2147483647}`)
		testMessageToNumberJson(t, "int32 int32_field = 1;", `{"int32Field": -2147483648}`, `{"1": -2147483648}`)
	})

	t.Run("uint32", func(t *testing.T) {
		testMessageToNumberJson(t, "uint32 uint32_field = 1;", `{"uint32Field": 0}`, `{}`)
		testMessageToNumberJson(t, "uint32 uint32_field = 1;", `{"uint32Field": 4294967295}`, `{"1": 4294967295}`)
	})

	t.Run("int64", func(t *testing.T) {
		testMessageToNumberJson(t, "int64 int64_field = 1;", `{"int64Field": "0"}`, `{}`)
		testMessageToNumberJson(t, "int64 int64_field = 1;", `{"int64Field": "9223372036854775807"}`, `{"1": 9223372036854775807}`)
		testMessageToNumberJson(t, "int64 int64_field = 1;", `{"int64Field": "-9223372036854775808"}`, `{"1": -9223372036854775808}`)
	})

	t.Run("uint64", func(t *testing.T) {
		testMessageToNumberJson(t, "uint64 uint64_field = 1;", `{"uint64Field": "0"}`, `{}`)
		testMessageToNumberJson(t, "uint64 uint64_field = 1;", `{"uint64Field": "18446744073709551615"}`, `{"1": 18446744073709551615}`)
	})

	t.Run("fixed32", func(t *testing.T) {
		testMessageToNumberJson(t, "fixed32 fixed32_field = 1;", `{"fixed32Field": 0}`, `{}`)
		testMessageToNumberJson(t, "fixed32 fixed32_field = 1;", `{"fixed32Field": 4294967295}`, `{"1": 4294967295}`)
	})

	t.Run("fixed64", func(t *testing.T) {
		testMessageToNumberJson(t, "fixed64 fixed64_field = 1;", `{"fixed64Field": "0"}`, `{}`)
		testMessageToNumberJson(t, "fixed64 fixed64_field = 1;", `{"fixed64Field": "18446744073709551615"}`, `{"1": 18446744073709551615}`)
	})

	t.Run("sfixed32", func(t *testing.T) {
		testMessageToNumberJson(t, "sfixed32 sfixed32_field = 1;", `{"sfixed32Field": 0}`, `{}`)
		testMessageToNumberJson(t, "sfixed32 sfixed32_field = 1;", `{"sfixed32Field": 2147483647}`, `{"1": 2147483647}`)
		testMessageToNumberJson(t, "sfixed32 sfixed32_field = 1;", `{"sfixed32Field": -2147483648}`, `{"1": -2147483648}`)
	})

	t.Run("sfixed64", func(t *testing.T) {
		testMessageToNumberJson(t, "sfixed64 sfixed64_field = 1;", `{"sfixed64Field": "0"}`, `{}`)
		testMessageToNumberJson(t, "sfixed64 sfixed64_field = 1;", `{"sfixed64Field": "9223372036854775807"}`, `{"1": 9223372036854775807}`)
		testMessageToNumberJson(t, "sfixed64 sfixed64_field = 1;", `{"sfixed64Field": "-9223372036854775808"}`, `{"1": -9223372036854775808}`)
	})

	t.Run("bool", func(t *testing.T) {
		testMessageToNumberJson(t, "bool bool_field = 1;", `{"boolField": false}`, `{}`)
		testMessageToNumberJson(t, "bool bool_field = 1;", `{"boolField": true}`, `{"1": true}`)
	})

	t.Run("string", func(t *testing.T) {
		testMessageToNumberJson(t, "string string_field = 1;", `{"stringField": ""}`, `{}`)
		testMessageToNumberJson(t, "string string_field = 1;", `{"stringField": "test"}`, `{"1": "test"}`)
	})

	t.Run("bytes", func(t *testing.T) {
		testMessageToNumberJson(t, "bytes bytes_field = 1;", `{"bytesField": ""}`, `{}`)
		testMessageToNumberJson(t, "bytes bytes_field = 1;", `{"bytesField": "aGVsbG8="}`, `{"1": "aGVsbG8="}`)
	})

	t.Run("float", func(t *testing.T) {
		testMessageToNumberJson(t, "float float_field = 1;", `{"floatField": 0}`, `{}`)
		testMessageToNumberJson(t, "float float_field = 1;", `{"floatField": 3.5}`, `{"1": 3.5}`)
	})

	t.Run("double", func(t *testing.T) {
		testMessageToNumberJson(t, "double double_field = 1;", `{"doubleField": 0}`, `{}`)
		testMessageToNumberJson(t, "double double_field = 1;", `{"doubleField": 3.141592653589793}`, `{"1": 3.141592653589793}`)
	})

	t.Run("sint32", func(t *testing.T) {
		testMessageToNumberJson(t, "sint32 sint32_field = 1;", `{"sint32Field": 0}`, `{}`)
		testMessageToNumberJson(t, "sint32 sint32_field = 1;", `{"sint32Field": 2147483647}`, `{"1": 2147483647}`)
		testMessageToNumberJson(t, "sint32 sint32_field = 1;", `{"sint32Field": -2147483648}`, `{"1": -2147483648}`)
	})

	t.Run("sint64", func(t *testing.T) {
		testMessageToNumberJson(t, "sint64 sint64_field = 1;", `{"sint64Field": "0"}`, `{}`)
		testMessageToNumberJson(t, "sint64 sint64_field = 1;", `{"sint64Field": "9223372036854775807"}`, `{"1": 9223372036854775807}`)
		testMessageToNumberJson(t, "sint64 sint64_field = 1;", `{"sint64Field": "-9223372036854775808"}`, `{"1": -9223372036854775808}`)
	})

	t.Run("enum", func(t *testing.T) {
		testMessageToNumberJson(t, "EnumType enum_field = 1;", `{"enumField": "ENUM_TYPE_UNSPECIFIED"}`, `{}`)
		testMessageToNumberJson(t, "EnumType enum_field = 1;", `{"enumField": "ENUM_TYPE_ONE"}`, `{"1": 1}`)
	})

	t.Run("message", func(t *testing.T) {
		testMessageToNumberJson(t, "MessageType message_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "MessageType message_field = 1;", `{"messageField": {"value": 0}}`, `{"1": {}}`)
		testMessageToNumberJson(t, "MessageType message_field = 1;", `{"messageField": {"value": 12345}}`, `{"1": {"1": 12345}}`)
	})
}

func TestMessageToNumberJsonRepeatedFields(t *testing.T) {
	t.Run("repeated int32", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": [0]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": [-2147483648, 0, 2147483647]}`, `{"1": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated uint32", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated uint32 repeated_uint32_field = 1;", `{"repeatedUint32Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated uint32 repeated_uint32_field = 1;", `{"repeatedUint32Field": [0]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated uint32 repeated_uint32_field = 1;", `{"repeatedUint32Field": [0, 4294967295]}`, `{"1": [0, 4294967295]}`)
	})

	t.Run("repeated int64", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated int64 repeated_int64_field = 1;", `{"repeatedInt64Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated int64 repeated_int64_field = 1;", `{"repeatedInt64Field": ["0"]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated int64 repeated_int64_field = 1;", `{"repeatedInt64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`, `{"1": [-9223372036854775808, 0, 9223372036854775807]}`)
	})

	t.Run("repeated uint64", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated uint64 repeated_uint64_field = 1;", `{"repeatedUint64Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated uint64 repeated_uint64_field = 1;", `{"repeatedUint64Field": ["0"]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated uint64 repeated_uint64_field = 1;", `{"repeatedUint64Field": ["0", "18446744073709551615"]}`, `{"1": [0, 18446744073709551615]}`)
	})

	t.Run("repeated fixed32", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"repeatedFixed32Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"repeatedFixed32Field": [0]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"repeatedFixed32Field": [0, 4294967295]}`, `{"1": [0, 4294967295]}`)
	})

	t.Run("repeated fixed64", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"repeatedFixed64Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"repeatedFixed64Field": ["0"]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"repeatedFixed64Field": ["0", "18446744073709551615"]}`, `{"1": [0, 18446744073709551615]}`)
	})

	t.Run("repeated sfixed32", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"repeatedSfixed32Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"repeatedSfixed32Field": [0]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"repeatedSfixed32Field": [-2147483648, 0, 2147483647]}`, `{"1": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated sfixed64", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"repeatedSfixed64Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"repeatedSfixed64Field": ["0"]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"repeatedSfixed64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`, `{"1": [-9223372036854775808, 0, 9223372036854775807]}`)
	})

	t.Run("repeated bool", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated bool repeated_bool_field = 1;", `{"repeatedBoolField": []}`, `{}`)
		testMessageToNumberJson(t, "repeated bool repeated_bool_field = 1;", `{"repeatedBoolField": [false]}`, `{"1": [false]}`)
		testMessageToNumberJson(t, "repeated bool repeated_bool_field = 1;", `{"repeatedBoolField": [true, false]}`, `{"1": [true, false]}`)
	})

	t.Run("repeated string", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated string repeated_string_field = 1;", `{"repeatedStringField": []}`, `{}`)
		testMessageToNumberJson(t, "repeated string repeated_string_field = 1;", `{"repeatedStringField": [""]}`, `{"1": [""]}`)
		testMessageToNumberJson(t, "repeated string repeated_string_field = 1;", `{"repeatedStringField": ["test", ""]}`, `{"1": ["test", ""]}`)
	})

	t.Run("repeated bytes", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated bytes repeated_bytes_field = 1;", `{"repeatedBytesField": []}`, `{}`)
		testMessageToNumberJson(t, "repeated bytes repeated_bytes_field = 1;", `{"repeatedBytesField": [""]}`, `{"1": [""]}`)
		testMessageToNumberJson(t, "repeated bytes repeated_bytes_field = 1;", `{"repeatedBytesField": ["aGVsbG8=", ""]}`, `{"1": ["aGVsbG8=", ""]}`)
	})

	t.Run("repeated float", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated float repeated_float_field = 1;", `{"repeatedFloatField": []}`, `{}`)
		testMessageToNumberJson(t, "repeated float repeated_float_field = 1;", `{"repeatedFloatField": [0]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated float repeated_float_field = 1;", `{"repeatedFloatField": [3.5, 0]}`, `{"1": [3.5, 0]}`)
	})

	t.Run("repeated double", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated double repeated_double_field = 1;", `{"repeatedDoubleField": []}`, `{}`)
		testMessageToNumberJson(t, "repeated double repeated_double_field = 1;", `{"repeatedDoubleField": [0]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated double repeated_double_field = 1;", `{"repeatedDoubleField": [-1.7976931348623157e+308, 0, 1.7976931348623157e+308]}`, `{"1": [-1.7976931348623157e+308, 0, 1.7976931348623157e+308]}`)
	})

	t.Run("repeated sint32", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated sint32 repeated_sint32_field = 1;", `{"repeatedSint32Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated sint32 repeated_sint32_field = 1;", `{"repeatedSint32Field": [0]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated sint32 repeated_sint32_field = 1;", `{"repeatedSint32Field": [-2147483648, 0, 2147483647]}`, `{"1": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated sint64", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated sint64 repeated_sint64_field = 1;", `{"repeatedSint64Field": []}`, `{}`)
		testMessageToNumberJson(t, "repeated sint64 repeated_sint64_field = 1;", `{"repeatedSint64Field": ["0"]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated sint64 repeated_sint64_field = 1;", `{"repeatedSint64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`, `{"1": [-9223372036854775808, 0, 9223372036854775807]}`)
	})

	t.Run("repeated enum", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": []}`, `{}`)
		testMessageToNumberJson(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": ["ENUM_TYPE_UNSPECIFIED"]}`, `{"1": [0]}`)
		testMessageToNumberJson(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": ["ENUM_TYPE_ONE", "ENUM_TYPE_UNSPECIFIED"]}`, `{"1": [1, 0]}`)
	})

	t.Run("repeated message", func(t *testing.T) {
		testMessageToNumberJson(t, "repeated MessageType repeated_message_field = 1;", `{"repeatedMessageField": []}`, `{}`)
		testMessageToNumberJson(t, "repeated MessageType repeated_message_field = 1;", `{"repeatedMessageField": [{"value": 0}]}`, `{"1": [{}]}`)
		testMessageToNumberJson(t, "repeated MessageType repeated_message_field = 1;", `{"repeatedMessageField": [{"value": 12345}, {"value": 67890}]}`, `{"1": [{"1": 12345}, {"1": 67890}]}`)
	})
}

func TestMessageToNumberJsonOptionalFields(t *testing.T) {
	t.Run("optional int32", func(t *testing.T) {
		testMessageToNumberJson(t, "optional int32 optional_int32_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": 0}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": 2147483647}`, `{"1": 2147483647}`)
		testMessageToNumberJson(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": -2147483648}`, `{"1": -2147483648}`)
	})

	t.Run("optional uint32", func(t *testing.T) {
		testMessageToNumberJson(t, "optional uint32 optional_uint32_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional uint32 optional_uint32_field = 1;", `{"optionalUint32Field": 0}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional uint32 optional_uint32_field = 1;", `{"optionalUint32Field": 4294967295}`, `{"1": 4294967295}`)
	})

	t.Run("optional int64", func(t *testing.T) {
		testMessageToNumberJson(t, "optional int64 optional_int64_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional int64 optional_int64_field = 1;", `{"optionalInt64Field": "0"}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional int64 optional_int64_field = 1;", `{"optionalInt64Field": "9223372036854775807"}`, `{"1": 9223372036854775807}`)
		testMessageToNumberJson(t, "optional int64 optional_int64_field = 1;", `{"optionalInt64Field": "-9223372036854775808"}`, `{"1": -9223372036854775808}`)
	})

	t.Run("optional uint64", func(t *testing.T) {
		testMessageToNumberJson(t, "optional uint64 optional_uint64_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional uint64 optional_uint64_field = 1;", `{"optionalUint64Field": "0"}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional uint64 optional_uint64_field = 1;", `{"optionalUint64Field": "18446744073709551615"}`, `{"1": 18446744073709551615}`)
	})

	t.Run("optional fixed32", func(t *testing.T) {
		testMessageToNumberJson(t, "optional fixed32 optional_fixed32_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional fixed32 optional_fixed32_field = 1;", `{"optionalFixed32Field": 0}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional fixed32 optional_fixed32_field = 1;", `{"optionalFixed32Field": 4294967295}`, `{"1": 4294967295}`)
	})

	t.Run("optional fixed64", func(t *testing.T) {
		testMessageToNumberJson(t, "optional fixed64 optional_fixed64_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional fixed64 optional_fixed64_field = 1;", `{"optionalFixed64Field": "0"}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional fixed64 optional_fixed64_field = 1;", `{"optionalFixed64Field": "18446744073709551615"}`, `{"1": 18446744073709551615}`)
	})

	t.Run("optional sfixed32", func(t *testing.T) {
		testMessageToNumberJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"optionalSfixed32Field": 0}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"optionalSfixed32Field": 2147483647}`, `{"1": 2147483647}`)
		testMessageToNumberJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"optionalSfixed32Field": -2147483648}`, `{"1": -2147483648}`)
	})

	t.Run("optional sfixed64", func(t *testing.T) {
		testMessageToNumberJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"optionalSfixed64Field": "0"}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"optionalSfixed64Field": "9223372036854775807"}`, `{"1": 9223372036854775807}`)
		testMessageToNumberJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"optionalSfixed64Field": "-9223372036854775808"}`, `{"1": -9223372036854775808}`)
	})

	t.Run("optional bool", func(t *testing.T) {
		testMessageToNumberJson(t, "optional bool optional_bool_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional bool optional_bool_field = 1;", `{"optionalBoolField": false}`, `{"1": false}`)
		testMessageToNumberJson(t, "optional bool optional_bool_field = 1;", `{"optionalBoolField": true}`, `{"1": true}`)
	})

	t.Run("optional string", func(t *testing.T) {
		testMessageToNumberJson(t, "optional string optional_string_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional string optional_string_field = 1;", `{"optionalStringField": ""}`, `{"1": ""}`)
		testMessageToNumberJson(t, "optional string optional_string_field = 1;", `{"optionalStringField": "testMessageToJson"}`, `{"1": "testMessageToJson"}`)
	})

	t.Run("optional bytes", func(t *testing.T) {
		testMessageToNumberJson(t, "optional bytes optional_bytes_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional bytes optional_bytes_field = 1;", `{"optionalBytesField": ""}`, `{"1": ""}`)
		testMessageToNumberJson(t, "optional bytes optional_bytes_field = 1;", `{"optionalBytesField": "aGVsbG8="}`, `{"1": "aGVsbG8="}`)
	})

	t.Run("optional float", func(t *testing.T) {
		testMessageToNumberJson(t, "optional float optional_float_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional float optional_float_field = 1;", `{"optionalFloatField": 0}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional float optional_float_field = 1;", `{"optionalFloatField": 3.5}`, `{"1": 3.5}`)
	})

	t.Run("optional double", func(t *testing.T) {
		testMessageToNumberJson(t, "optional double optional_double_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional double optional_double_field = 1;", `{"optionalDoubleField": 0}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional double optional_double_field = 1;", `{"optionalDoubleField": 3.141592653589793}`, `{"1": 3.141592653589793}`)
	})

	t.Run("optional sint32", func(t *testing.T) {
		testMessageToNumberJson(t, "optional sint32 optional_sint32_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional sint32 optional_sint32_field = 1;", `{"optionalSint32Field": 0}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional sint32 optional_sint32_field = 1;", `{"optionalSint32Field": 2147483647}`, `{"1": 2147483647}`)
		testMessageToNumberJson(t, "optional sint32 optional_sint32_field = 1;", `{"optionalSint32Field": -2147483648}`, `{"1": -2147483648}`)
	})

	t.Run("optional sint64", func(t *testing.T) {
		testMessageToNumberJson(t, "optional sint64 optional_sint64_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional sint64 optional_sint64_field = 1;", `{"optionalSint64Field": "0"}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional sint64 optional_sint64_field = 1;", `{"optionalSint64Field": "9223372036854775807"}`, `{"1": 9223372036854775807}`)
		testMessageToNumberJson(t, "optional sint64 optional_sint64_field = 1;", `{"optionalSint64Field": "-9223372036854775808"}`, `{"1": -9223372036854775808}`)
	})

	t.Run("optional enum", func(t *testing.T) {
		testMessageToNumberJson(t, "optional EnumType optional_enum_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional EnumType optional_enum_field = 1;", `{"optionalEnumField": "ENUM_TYPE_UNSPECIFIED"}`, `{"1": 0}`)
		testMessageToNumberJson(t, "optional EnumType optional_enum_field = 1;", `{"optionalEnumField": "ENUM_TYPE_ONE"}`, `{"1": 1}`)
	})

	t.Run("optional message", func(t *testing.T) {
		testMessageToNumberJson(t, "optional MessageType optional_message_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "optional MessageType optional_message_field = 1;", `{"optionalMessageField": {"value": 0}}`, `{"1": {}}`)
		testMessageToNumberJson(t, "optional MessageType optional_message_field = 1;", `{"optionalMessageField": {"value": 12345}}`, `{"1": {"1": 12345}}`)
	})
}

func TestMessageToNumberJsonMapKey(t *testing.T) {
	t.Run("map<int32, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<int32, string> int32_key_map_field = 1;", `{"int32KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<int32, string> int32_key_map_field = 1;", `{"int32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`, `{"1": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<uint32, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<uint32, string> uint32_key_map_field = 1;", `{"uint32KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<uint32, string> uint32_key_map_field = 1;", `{"uint32KeyMapField": {"0": "a", "4294967295": "b"}}`, `{"1": {"0": "a", "4294967295": "b"}}`)
	})

	t.Run("map<int64, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<int64, string> int64_key_map_field = 1;", `{"int64KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<int64, string> int64_key_map_field = 1;", `{"int64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`, `{"1": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	t.Run("map<uint64, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<uint64, string> uint64_key_map_field = 1;", `{"uint64KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<uint64, string> uint64_key_map_field = 1;", `{"uint64KeyMapField": {"0": "a", "18446744073709551615": "b"}}`, `{"1": {"0": "a", "18446744073709551615": "b"}}`)
	})

	t.Run("map<fixed32, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<fixed32, string> fixed32_key_map_field = 1;", `{"fixed32KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<fixed32, string> fixed32_key_map_field = 1;", `{"fixed32KeyMapField": {"0": "a", "4294967295": "b"}}`, `{"1": {"0": "a", "4294967295": "b"}}`)
	})

	t.Run("map<fixed64, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<fixed64, string> fixed64_key_map_field = 1;", `{"fixed64KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<fixed64, string> fixed64_key_map_field = 1;", `{"fixed64KeyMapField": {"0": "a", "18446744073709551615": "b"}}`, `{"1": {"0": "a", "18446744073709551615": "b"}}`)
	})

	t.Run("map<sfixed32, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<sfixed32, string> sfixed32_key_map_field = 1;", `{"sfixed32KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<sfixed32, string> sfixed32_key_map_field = 1;", `{"sfixed32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`, `{"1": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<sfixed64, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<sfixed64, string> sfixed64_key_map_field = 1;", `{"sfixed64KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<sfixed64, string> sfixed64_key_map_field = 1;", `{"sfixed64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`, `{"1": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	t.Run("map<bool, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<bool, string> bool_key_map_field = 1;", `{"boolKeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<bool, string> bool_key_map_field = 1;", `{"boolKeyMapField": {"false": "a", "true": "b"}}`, `{"1": {"false": "a", "true": "b"}}`)
	})

	t.Run("map<string, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, string> string_key_map_field = 1;", `{"stringKeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, string> string_key_map_field = 1;", `{"stringKeyMapField": {"a": "b", "c": "d"}}`, `{"1": {"a": "b", "c": "d"}}`)
	})

	t.Run("map<sint32, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<sint32, string> sint32_key_map_field = 1;", `{"sint32KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<sint32, string> sint32_key_map_field = 1;", `{"sint32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`, `{"1": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<sint64, *>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<sint64, string> sint64_key_map_field = 1;", `{"sint64KeyMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<sint64, string> sint64_key_map_field = 1;", `{"sint64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`, `{"1": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	// NOTE: float, double, enum, or message cannot be a map key.
}

func TestMessageToNumberJsonMapValue(t *testing.T) {
	t.Run("map<*, int32>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, int32> int32_value_map_field = 1;", `{"int32ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, int32> int32_value_map_field = 1;", `{"int32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`, `{"1": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, uint32>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, uint32> uint32_value_map_field = 1;", `{"uint32ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, uint32> uint32_value_map_field = 1;", `{"uint32ValueMapField": {"a": 0, "b": 4294967295}}`, `{"1": {"a": 0, "b": 4294967295}}`)
	})

	t.Run("map<*, int64>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, int64> int64_value_map_field = 1;", `{"int64ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, int64> int64_value_map_field = 1;", `{"int64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`, `{"1": {"a": 0, "b": 9223372036854775807, "c": -9223372036854775808}}`)
	})

	t.Run("map<*, uint64>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, uint64> uint64_value_map_field = 1;", `{"uint64ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, uint64> uint64_value_map_field = 1;", `{"uint64ValueMapField": {"a": "0", "b": "18446744073709551615"}}`, `{"1": {"a": 0, "b": 18446744073709551615}}`)
	})

	t.Run("map<*, fixed32>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, fixed32> fixed32_value_map_field = 1;", `{"fixed32ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, fixed32> fixed32_value_map_field = 1;", `{"fixed32ValueMapField": {"a": 0, "b": 4294967295}}`, `{"1": {"a": 0, "b": 4294967295}}`)
	})

	t.Run("map<*, fixed64>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, fixed64> fixed64_value_map_field = 1;", `{"fixed64ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, fixed64> fixed64_value_map_field = 1;", `{"fixed64ValueMapField": {"a": "0", "b": "18446744073709551615"}}`, `{"1": {"a": 0, "b": 18446744073709551615}}`)
	})

	t.Run("map<*, sfixed32>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, sfixed32> sfixed32_value_map_field = 1;", `{"sfixed32ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, sfixed32> sfixed32_value_map_field = 1;", `{"sfixed32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`, `{"1": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, sfixed64>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, sfixed64> sfixed64_value_map_field = 1;", `{"sfixed64ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, sfixed64> sfixed64_value_map_field = 1;", `{"sfixed64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`, `{"1": {"a": 0, "b": 9223372036854775807, "c": -9223372036854775808}}`)
	})

	t.Run("map<*, bool>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, bool> bool_value_map_field = 1;", `{"boolValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, bool> bool_value_map_field = 1;", `{"boolValueMapField": {"a": false, "b": true}}`, `{"1": {"a": false, "b": true}}`)
	})

	t.Run("map<*, string>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, string> string_value_map_field = 1;", `{"stringValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, string> string_value_map_field = 1;", `{"stringValueMapField": {"a": "b", "c": "d"}}`, `{"1": {"a": "b", "c": "d"}}`)
	})

	t.Run("map<*, bytes>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, bytes> bytes_value_map_field = 1;", `{"bytesValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, bytes> bytes_value_map_field = 1;", `{"bytesValueMapField": {"a": "", "b": "dGVzdA=="}}`, `{"1": {"a": "", "b": "dGVzdA=="}}`) // Base64 for "test"
	})

	t.Run("map<*, float>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, float> float_value_map_field = 1;", `{"floatValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, float> float_value_map_field = 1;", `{"floatValueMapField": {"a": 0, "b": 3.5}}`, `{"1": {"a": 0, "b": 3.5}}`)
	})

	t.Run("map<*, double>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, double> double_value_map_field = 1;", `{"doubleValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, double> double_value_map_field = 1;", `{"doubleValueMapField": {"a": 0, "b": 3.141592653589793}}`, `{"1": {"a": 0, "b": 3.141592653589793}}`)
	})

	t.Run("map<*, sint32>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, sint32> sint32_value_map_field = 1;", `{"sint32ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, sint32> sint32_value_map_field = 1;", `{"sint32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`, `{"1": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, sint64>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, sint64> sint64_value_map_field = 1;", `{"sint64ValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, sint64> sint64_value_map_field = 1;", `{"sint64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`, `{"1": {"a": 0, "b": 9223372036854775807, "c": -9223372036854775808}}`)
	})

	t.Run("map<*, EnumType>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, EnumType> enum_value_map_field = 1;", `{"enumValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, EnumType> enum_value_map_field = 1;", `{"enumValueMapField": {"a": "ENUM_TYPE_UNSPECIFIED", "b": "ENUM_TYPE_ONE"}}`, `{"1": {"a": 0, "b": 1}}`)
	})

	t.Run("map<*, MessageType>", func(t *testing.T) {
		testMessageToNumberJson(t, "map<string, MessageType> message_value_map_field = 1;", `{"messageValueMapField": {}}`, `{}`)
		testMessageToNumberJson(t, "map<string, MessageType> message_value_map_field = 1;", `{"messageValueMapField": {"a": {"value": 0}, "b": {"value": 12345}}}`, `{"1": {"a": {}, "b": {"1": 12345}}}`)
	})
}

func TestMessageToNumberJsonOneof(t *testing.T) {
	t.Run("oneof", func(t *testing.T) {
		testMessageToNumberJson(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{}`, `{}`)
		testMessageToNumberJson(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{"int32Field": 42}`, `{"1": 42}`)
		testMessageToNumberJson(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{"stringField": "test"}`, `{"2": "test"}`)
	})

	t.Run("oneof with message", func(t *testing.T) {
		testMessageToNumberJson(t, "oneof kind { int32 int32_field = 1; MessageType message_field = 2; }", `{}`, `{}`)
		testMessageToNumberJson(t, "oneof kind { int32 int32_field = 1; MessageType message_field = 2; }", `{"int32Field": 42}`, `{"1": 42}`)
		testMessageToNumberJson(t, "oneof kind { int32 int32_field = 1; MessageType message_field = 2; }", `{"messageField": {"value": 123}}`, `{"2": {"1": 123}}`)
	})

	t.Run("oneof with enum", func(t *testing.T) {
		testMessageToNumberJson(t, "oneof kind { int32 int32_field = 1; EnumType enum_field = 2; }", `{}`, `{}`)
		testMessageToNumberJson(t, "oneof kind { int32 int32_field = 1; EnumType enum_field = 2; }", `{"int32Field": 42}`, `{"1": 42}`)
		testMessageToNumberJson(t, "oneof kind { int32 int32_field = 1; EnumType enum_field = 2; }", `{"enumField": "ENUM_TYPE_ONE"}`, `{"2": 1}`)
	})
}

func TestMessageToNumberJsonWellKnownTypes(t *testing.T) {
	t.Run("Timestamp", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"timestampField": "1970-01-01T00:00:00Z"}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"timestampField": "2023-10-01T12:34:56.789Z"}`, `{"1": {"1": 1696163696, "2": 789000000}}`)
	})

	t.Run("Duration", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.Duration duration_field = 1;", `{"durationField": "0s"}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.Duration duration_field = 1;", `{"durationField": "1.234s"}`, `{"1": {"1": 1, "2": 234000000}}`)
	})

	t.Run("Struct", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.Struct struct_field = 1;", `{"structField": {}}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.Struct struct_field = 1;", `{"structField": {"key": "value"}}`, `{"1": {"1": {"key": {"3": "value"}}}}`)
		testMessageToNumberJson(t, "google.protobuf.Struct struct_field = 1;", `{"structField": {"number": 123, "boolean": true}}`, `{"1": {"1": {"number": {"2": 123}, "boolean": {"4": true}}}}`)
	})

	t.Run("ListValue", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.ListValue list_value_field = 1;", `{"listValueField": []}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.ListValue list_value_field = 1;", `{"listValueField": ["string", 123, true]}`, `{"1": {"1": [{"3": "string"}, {"2": 123}, {"4": true}]}}`)
	})

	t.Run("Value", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.Value value_field = 1;", `{"valueField": {}}`, `{"1": {"5": {}}}`)
		testMessageToNumberJson(t, "google.protobuf.Value value_field = 1;", `{"valueField": "string"}`, `{"1": {"3": "string"}}`)
		testMessageToNumberJson(t, "google.protobuf.Value value_field = 1;", `{"valueField": 123}`, `{"1": {"2": 123}}`)
		testMessageToNumberJson(t, "google.protobuf.Value value_field = 1;", `{"valueField": true}`, `{"1": {"4": true}}`)
	})

	t.Run("Empty", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.Empty empty_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.Empty empty_field = 1;", `{"emptyField": {}}`, `{"1": {}}`)
	})

	t.Run("DoubleValue", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.DoubleValue double_value_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.DoubleValue double_value_field = 1;", `{"doubleValueField": 0}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.DoubleValue double_value_field = 1;", `{"doubleValueField": 1.7976931348623157e+308}`, `{"1": {"1": 1.7976931348623157e+308}}`)
	})

	t.Run("FloatValue", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.FloatValue float_value_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.FloatValue float_value_field = 1;", `{"floatValueField": 0}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.FloatValue float_value_field = 1;", `{"floatValueField": 3.5}`, `{"1": {"1": 3.5}}`)
	})

	t.Run("Int64Value", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"int64ValueField": "0"}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"int64ValueField": "9223372036854775807"}`, `{"1": {"1": 9223372036854775807}}`)
		testMessageToNumberJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"int64ValueField": "-9223372036854775808"}`, `{"1": {"1": -9223372036854775808}}`)
	})

	t.Run("UInt64Value", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{"uint64ValueField": "0"}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{"uint64ValueField": "18446744073709551615"}`, `{"1": {"1": 18446744073709551615}}`)
	})

	t.Run("Int32Value", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"int32ValueField": 0}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"int32ValueField": 2147483647}`, `{"1": {"1": 2147483647}}`)
		testMessageToNumberJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"int32ValueField": -2147483648}`, `{"1": {"1": -2147483648}}`)
	})

	t.Run("UInt32Value", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{"uint32ValueField": 0}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{"uint32ValueField": 4294967295}`, `{"1": {"1": 4294967295}}`)
	})

	t.Run("BoolValue", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.BoolValue bool_value_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.BoolValue bool_value_field = 1;", `{"boolValueField": false}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.BoolValue bool_value_field = 1;", `{"boolValueField": true}`, `{"1": {"1": true}}`)
	})

	t.Run("StringValue", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.StringValue string_value_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.StringValue string_value_field = 1;", `{"stringValueField": ""}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.StringValue string_value_field = 1;", `{"stringValueField": "test"}`, `{"1": {"1": "test"}}`)
	})

	t.Run("BytesValue", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{"bytesValueField": ""}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{"bytesValueField": "aGVsbG8="}`, `{"1": {"1": "aGVsbG8="}}`)
	})

	t.Run("FieldMask", func(t *testing.T) {
		testMessageToNumberJson(t, "google.protobuf.FieldMask field_mask_field = 1;", `{}`, `{}`)
		testMessageToNumberJson(t, "google.protobuf.FieldMask field_mask_field = 1;", `{"fieldMaskField": ""}`, `{"1": {}}`)
		testMessageToNumberJson(t, "google.protobuf.FieldMask field_mask_field = 1;", `{"fieldMaskField": "path1,path2"}`, `{"1": {"1": ["path1", "path2"]}}`)
	})
}

func TestMessageToNumberJsonNullInput(t *testing.T) {
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

	RunTestThatExpression(t, "_pb_message_to_number_json(?, ?, ?)", descriptorSetJson, typeName, nil).IsNull()
}
