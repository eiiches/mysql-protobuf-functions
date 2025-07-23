package main

import (
	"fmt"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	"github.com/eiiches/mysql-protobuf-functions/internal/testutils"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func testMessageToJson(t *testing.T, fieldDefinition string, input string) {
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

	descriptorSetName := "a"
	typeName := protoreflect.FullName(".Test")

	AssertThatCall(t, "pb_descriptor_set_load(?, ?)", descriptorSetName, p.GetSerializedFileDescriptorSet()).ShouldSucceed()
	defer func() {
		AssertThatCall(t, "pb_descriptor_set_delete(?)", descriptorSetName).ShouldSucceed()
	}()

	dynamicMessage := p.JsonToDynamicMessage(typeName, input)
	serializedBinary := p.JsonToProtobuf(typeName, input)

	g := NewWithT(t)
	expectedJson, err := (&protojson.MarshalOptions{EmitDefaultValues: true}).Marshal(dynamicMessage.Interface())
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(expectedJson).To(MatchJSON(input), "Test case is invalid: input should match the output of protojson.Marshal(input).")

	RunTestThatExpression(t, "pb_message_to_json(?, ?, ?)", descriptorSetName, typeName, serializedBinary).IsEqualToJsonString(string(expectedJson))
}

func TestMessageToJsonSingularFields(t *testing.T) {
	t.Run("int32", func(t *testing.T) {
		testMessageToJson(t, "int32 int32_field = 1;", `{"int32Field": 0}`)
		testMessageToJson(t, "int32 int32_field = 1;", `{"int32Field": 2147483647}`)
		testMessageToJson(t, "int32 int32_field = 1;", `{"int32Field": -2147483648}`)
	})

	t.Run("uint32", func(t *testing.T) {
		testMessageToJson(t, "uint32 uint32_field = 1;", `{"uint32Field": 0}`)
		testMessageToJson(t, "uint32 uint32_field = 1;", `{"uint32Field": 4294967295}`)
	})

	t.Run("int64", func(t *testing.T) {
		testMessageToJson(t, "int64 int64_field = 1;", `{"int64Field": "0"}`)
		testMessageToJson(t, "int64 int64_field = 1;", `{"int64Field": "9223372036854775807"}`)
		testMessageToJson(t, "int64 int64_field = 1;", `{"int64Field": "-9223372036854775808"}`)
	})

	t.Run("uint64", func(t *testing.T) {
		testMessageToJson(t, "uint64 uint64_field = 1;", `{"uint64Field": "0"}`)
		testMessageToJson(t, "uint64 uint64_field = 1;", `{"uint64Field": "18446744073709551615"}`)
	})

	t.Run("fixed32", func(t *testing.T) {
		testMessageToJson(t, "fixed32 fixed32_field = 1;", `{"fixed32Field": 0}`)
		testMessageToJson(t, "fixed32 fixed32_field = 1;", `{"fixed32Field": 4294967295}`)
	})

	t.Run("fixed64", func(t *testing.T) {
		testMessageToJson(t, "fixed64 fixed64_field = 1;", `{"fixed64Field": "0"}`)
		testMessageToJson(t, "fixed64 fixed64_field = 1;", `{"fixed64Field": "18446744073709551615"}`)
	})

	t.Run("sfixed32", func(t *testing.T) {
		testMessageToJson(t, "sfixed32 sfixed32_field = 1;", `{"sfixed32Field": 0}`)
		testMessageToJson(t, "sfixed32 sfixed32_field = 1;", `{"sfixed32Field": 2147483647}`)
		testMessageToJson(t, "sfixed32 sfixed32_field = 1;", `{"sfixed32Field": -2147483648}`)
	})

	t.Run("sfixed64", func(t *testing.T) {
		testMessageToJson(t, "sfixed64 sfixed64_field = 1;", `{"sfixed64Field": "0"}`)
		testMessageToJson(t, "sfixed64 sfixed64_field = 1;", `{"sfixed64Field": "9223372036854775807"}`)
		testMessageToJson(t, "sfixed64 sfixed64_field = 1;", `{"sfixed64Field": "-9223372036854775808"}`)
	})

	t.Run("bool", func(t *testing.T) {
		testMessageToJson(t, "bool bool_field = 1;", `{"boolField": false}`)
		testMessageToJson(t, "bool bool_field = 1;", `{"boolField": true}`)
	})

	t.Run("string", func(t *testing.T) {
		testMessageToJson(t, "string string_field = 1;", `{"stringField": ""}`)
		testMessageToJson(t, "string string_field = 1;", `{"stringField": "testMessageToJson"}`)
	})

	t.Run("bytes", func(t *testing.T) {
		testMessageToJson(t, "bytes bytes_field = 1;", `{"bytesField": ""}`)
		testMessageToJson(t, "bytes bytes_field = 1;", `{"bytesField": "dGVzdA=="}`) // Base64 for "testMessageToJson"
	})

	t.Run("float", func(t *testing.T) {
		testMessageToJson(t, "float float_field = 1;", `{"floatField": 0}`)
		testMessageToJson(t, "float float_field = 1;", `{"floatField": 3.5}`)
	})

	t.Run("double", func(t *testing.T) {
		testMessageToJson(t, "double double_field = 1;", `{"doubleField": 0}`)
		testMessageToJson(t, "double double_field = 1;", `{"doubleField": 3.141592653589793}`)
	})

	t.Run("sint32", func(t *testing.T) {
		testMessageToJson(t, "sint32 sint32_field = 1;", `{"sint32Field": 0}`)
		testMessageToJson(t, "sint32 sint32_field = 1;", `{"sint32Field": 2147483647}`)
		testMessageToJson(t, "sint32 sint32_field = 1;", `{"sint32Field": -2147483648}`)
	})

	t.Run("sint64", func(t *testing.T) {
		testMessageToJson(t, "sint64 sint64_field = 1;", `{"sint64Field": "0"}`)
		testMessageToJson(t, "sint64 sint64_field = 1;", `{"sint64Field": "9223372036854775807"}`)
		testMessageToJson(t, "sint64 sint64_field = 1;", `{"sint64Field": "-9223372036854775808"}`)
	})

	t.Run("enum", func(t *testing.T) {
		testMessageToJson(t, "EnumType enum_field = 1;", `{"enumField": "ENUM_TYPE_UNSPECIFIED"}`)
		testMessageToJson(t, "EnumType enum_field = 1;", `{"enumField": "ENUM_TYPE_ONE"}`)
	})

	t.Run("message", func(t *testing.T) {
		testMessageToJson(t, "MessageType message_field = 1;", `{}`)
		testMessageToJson(t, "MessageType message_field = 1;", `{"messageField": {"value": 0}}`)
		testMessageToJson(t, "MessageType message_field = 1;", `{"messageField": {"value": 12345}}`)
	})
}

func TestMessageToJsonRepeatedFields(t *testing.T) {
	t.Run("repeated int32", func(t *testing.T) {
		testMessageToJson(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": []}`)
		testMessageToJson(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": [0]}`)
		testMessageToJson(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated uint32", func(t *testing.T) {
		testMessageToJson(t, "repeated uint32 repeated_uint32_field = 1;", `{"repeatedUint32Field": []}`)
		testMessageToJson(t, "repeated uint32 repeated_uint32_field = 1;", `{"repeatedUint32Field": [0]}`)
		testMessageToJson(t, "repeated uint32 repeated_uint32_field = 1;", `{"repeatedUint32Field": [0, 4294967295]}`)
	})

	t.Run("repeated int64", func(t *testing.T) {
		testMessageToJson(t, "repeated int64 repeated_int64_field = 1;", `{"repeatedInt64Field": []}`)
		testMessageToJson(t, "repeated int64 repeated_int64_field = 1;", `{"repeatedInt64Field": ["0"]}`)
		testMessageToJson(t, "repeated int64 repeated_int64_field = 1;", `{"repeatedInt64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`)
	})

	t.Run("repeated uint64", func(t *testing.T) {
		testMessageToJson(t, "repeated uint64 repeated_uint64_field = 1;", `{"repeatedUint64Field": []}`)
		testMessageToJson(t, "repeated uint64 repeated_uint64_field = 1;", `{"repeatedUint64Field": ["0"]}`)
		testMessageToJson(t, "repeated uint64 repeated_uint64_field = 1;", `{"repeatedUint64Field": ["0", "18446744073709551615"]}`)
	})

	t.Run("repeated fixed32", func(t *testing.T) {
		testMessageToJson(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"repeatedFixed32Field": []}`)
		testMessageToJson(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"repeatedFixed32Field": [0]}`)
		testMessageToJson(t, "repeated fixed32 repeated_fixed32_field = 1;", `{"repeatedFixed32Field": [0, 4294967295]}`)
	})

	t.Run("repeated fixed64", func(t *testing.T) {
		testMessageToJson(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"repeatedFixed64Field": []}`)
		testMessageToJson(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"repeatedFixed64Field": ["0"]}`)
		testMessageToJson(t, "repeated fixed64 repeated_fixed64_field = 1;", `{"repeatedFixed64Field": ["0", "18446744073709551615"]}`)
	})

	t.Run("repeated sfixed32", func(t *testing.T) {
		testMessageToJson(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"repeatedSfixed32Field": []}`)
		testMessageToJson(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"repeatedSfixed32Field": [0]}`)
		testMessageToJson(t, "repeated sfixed32 repeated_sfixed32_field = 1;", `{"repeatedSfixed32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated sfixed64", func(t *testing.T) {
		testMessageToJson(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"repeatedSfixed64Field": []}`)
		testMessageToJson(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"repeatedSfixed64Field": ["0"]}`)
		testMessageToJson(t, "repeated sfixed64 repeated_sfixed64_field = 1;", `{"repeatedSfixed64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`)
	})

	t.Run("repeated bool", func(t *testing.T) {
		testMessageToJson(t, "repeated bool repeated_bool_field = 1;", `{"repeatedBoolField": []}`)
		testMessageToJson(t, "repeated bool repeated_bool_field = 1;", `{"repeatedBoolField": [false]}`)
		testMessageToJson(t, "repeated bool repeated_bool_field = 1;", `{"repeatedBoolField": [true, false, true]}`)
	})

	t.Run("repeated string", func(t *testing.T) {
		testMessageToJson(t, "repeated string repeated_string_field = 1;", `{"repeatedStringField": []}`)
		testMessageToJson(t, "repeated string repeated_string_field = 1;", `{"repeatedStringField": [""]}`)
		testMessageToJson(t, "repeated string repeated_string_field = 1;", `{"repeatedStringField": ["testMessageToJson", ""]}`)
	})

	t.Run("repeated bytes", func(t *testing.T) {
		testMessageToJson(t, "repeated bytes repeated_bytes_field = 1;", `{"repeatedBytesField": []}`)
		testMessageToJson(t, "repeated bytes repeated_bytes_field = 1;", `{"repeatedBytesField": [""]}`)
		testMessageToJson(t, "repeated bytes repeated_bytes_field = 1;", `{"repeatedBytesField": ["dGVzdA==", ""]}`) // Base64 for "testMessageToJson"
	})

	t.Run("repeated float", func(t *testing.T) {
		testMessageToJson(t, "repeated float repeated_float_field = 1;", `{"repeatedFloatField": []}`)
		testMessageToJson(t, "repeated float repeated_float_field = 1;", `{"repeatedFloatField": [0]}`)
		testMessageToJson(t, "repeated float repeated_float_field = 1;", `{"repeatedFloatField": [3.5, 0]}`)
	})

	t.Run("repeated double", func(t *testing.T) {
		testMessageToJson(t, "repeated double repeated_double_field = 1;", `{"repeatedDoubleField": []}`)
		testMessageToJson(t, "repeated double repeated_double_field = 1;", `{"repeatedDoubleField": [0]}`)
		testMessageToJson(t, "repeated double repeated_double_field = 1;", `{"repeatedDoubleField": [3.141592653589793, 0]}`)
	})

	t.Run("repeated sint32", func(t *testing.T) {
		testMessageToJson(t, "repeated sint32 repeated_sint32_field = 1;", `{"repeatedSint32Field": []}`)
		testMessageToJson(t, "repeated sint32 repeated_sint32_field = 1;", `{"repeatedSint32Field": [0]}`)
		testMessageToJson(t, "repeated sint32 repeated_sint32_field = 1;", `{"repeatedSint32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated sint64", func(t *testing.T) {
		testMessageToJson(t, "repeated sint64 repeated_sint64_field = 1;", `{"repeatedSint64Field": []}`)
		testMessageToJson(t, "repeated sint64 repeated_sint64_field = 1;", `{"repeatedSint64Field": ["0"]}`)
		testMessageToJson(t, "repeated sint64 repeated_sint64_field = 1;", `{"repeatedSint64Field": ["-9223372036854775808", "0", "9223372036854775807"]}`)
	})

	t.Run("repeated enum", func(t *testing.T) {
		testMessageToJson(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": []}`)
		testMessageToJson(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": ["ENUM_TYPE_UNSPECIFIED"]}`)
		testMessageToJson(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": ["ENUM_TYPE_ONE", "ENUM_TYPE_UNSPECIFIED"]}`)
	})

	t.Run("repeated message", func(t *testing.T) {
		testMessageToJson(t, "repeated MessageType repeated_message_field = 1;", `{"repeatedMessageField": []}`)
		testMessageToJson(t, "repeated MessageType repeated_message_field = 1;", `{"repeatedMessageField": [{"value": 0}]}`)
		testMessageToJson(t, "repeated MessageType repeated_message_field = 1;", `{"repeatedMessageField": [{"value": 12345}, {"value": 67890}]}`)
	})
}

func TestMessageToJsonOptionalFields(t *testing.T) {
	t.Run("optional int32", func(t *testing.T) {
		testMessageToJson(t, "optional int32 optional_int32_field = 1;", `{}`)
		testMessageToJson(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": 0}`)
		testMessageToJson(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": 2147483647}`)
		testMessageToJson(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": -2147483648}`)
	})

	t.Run("optional uint32", func(t *testing.T) {
		testMessageToJson(t, "optional uint32 optional_uint32_field = 1;", `{}`)
		testMessageToJson(t, "optional uint32 optional_uint32_field = 1;", `{"optionalUint32Field": 0}`)
		testMessageToJson(t, "optional uint32 optional_uint32_field = 1;", `{"optionalUint32Field": 4294967295}`)
	})

	t.Run("optional int64", func(t *testing.T) {
		testMessageToJson(t, "optional int64 optional_int64_field = 1;", `{}`)
		testMessageToJson(t, "optional int64 optional_int64_field = 1;", `{"optionalInt64Field": "0"}`)
		testMessageToJson(t, "optional int64 optional_int64_field = 1;", `{"optionalInt64Field": "9223372036854775807"}`)
		testMessageToJson(t, "optional int64 optional_int64_field = 1;", `{"optionalInt64Field": "-9223372036854775808"}`)
	})

	t.Run("optional uint64", func(t *testing.T) {
		testMessageToJson(t, "optional uint64 optional_uint64_field = 1;", `{}`)
		testMessageToJson(t, "optional uint64 optional_uint64_field = 1;", `{"optionalUint64Field": "0"}`)
		testMessageToJson(t, "optional uint64 optional_uint64_field = 1;", `{"optionalUint64Field": "18446744073709551615"}`)
	})

	t.Run("optional fixed32", func(t *testing.T) {
		testMessageToJson(t, "optional fixed32 optional_fixed32_field = 1;", `{}`)
		testMessageToJson(t, "optional fixed32 optional_fixed32_field = 1;", `{"optionalFixed32Field": 0}`)
		testMessageToJson(t, "optional fixed32 optional_fixed32_field = 1;", `{"optionalFixed32Field": 4294967295}`)
	})

	t.Run("optional fixed64", func(t *testing.T) {
		testMessageToJson(t, "optional fixed64 optional_fixed64_field = 1;", `{}`)
		testMessageToJson(t, "optional fixed64 optional_fixed64_field = 1;", `{"optionalFixed64Field": "0"}`)
		testMessageToJson(t, "optional fixed64 optional_fixed64_field = 1;", `{"optionalFixed64Field": "18446744073709551615"}`)
	})

	t.Run("optional sfixed32", func(t *testing.T) {
		testMessageToJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{}`)
		testMessageToJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"optionalSfixed32Field": 0}`)
		testMessageToJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"optionalSfixed32Field": 2147483647}`)
		testMessageToJson(t, "optional sfixed32 optional_sfixed32_field = 1;", `{"optionalSfixed32Field": -2147483648}`)
	})

	t.Run("optional sfixed64", func(t *testing.T) {
		testMessageToJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{}`)
		testMessageToJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"optionalSfixed64Field": "0"}`)
		testMessageToJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"optionalSfixed64Field": "9223372036854775807"}`)
		testMessageToJson(t, "optional sfixed64 optional_sfixed64_field = 1;", `{"optionalSfixed64Field": "-9223372036854775808"}`)
	})

	t.Run("optional bool", func(t *testing.T) {
		testMessageToJson(t, "optional bool optional_bool_field = 1;", `{}`)
		testMessageToJson(t, "optional bool optional_bool_field = 1;", `{"optionalBoolField": false}`)
		testMessageToJson(t, "optional bool optional_bool_field = 1;", `{"optionalBoolField": true}`)
	})

	t.Run("optional string", func(t *testing.T) {
		testMessageToJson(t, "optional string optional_string_field = 1;", `{}`)
		testMessageToJson(t, "optional string optional_string_field = 1;", `{"optionalStringField": ""}`)
		testMessageToJson(t, "optional string optional_string_field = 1;", `{"optionalStringField": "testMessageToJson"}`)
	})

	t.Run("optional bytes", func(t *testing.T) {
		testMessageToJson(t, "optional bytes optional_bytes_field = 1;", `{}`)
		testMessageToJson(t, "optional bytes optional_bytes_field = 1;", `{"optionalBytesField": ""}`)
		testMessageToJson(t, "optional bytes optional_bytes_field = 1;", `{"optionalBytesField": "dGVzdA=="}`) // Base64 for "testMessageToJson"
	})

	t.Run("optional float", func(t *testing.T) {
		testMessageToJson(t, "optional float optional_float_field = 1;", `{}`)
		testMessageToJson(t, "optional float optional_float_field = 1;", `{"optionalFloatField": 0}`)
		testMessageToJson(t, "optional float optional_float_field = 1;", `{"optionalFloatField": 3.5}`)
	})

	t.Run("optional double", func(t *testing.T) {
		testMessageToJson(t, "optional double optional_double_field = 1;", `{}`)
		testMessageToJson(t, "optional double optional_double_field = 1;", `{"optionalDoubleField": 0}`)
		testMessageToJson(t, "optional double optional_double_field = 1;", `{"optionalDoubleField": 3.141592653589793}`)
	})

	t.Run("optional sint32", func(t *testing.T) {
		testMessageToJson(t, "optional sint32 optional_sint32_field = 1;", `{}`)
		testMessageToJson(t, "optional sint32 optional_sint32_field = 1;", `{"optionalSint32Field": 0}`)
		testMessageToJson(t, "optional sint32 optional_sint32_field = 1;", `{"optionalSint32Field": 2147483647}`)
		testMessageToJson(t, "optional sint32 optional_sint32_field = 1;", `{"optionalSint32Field": -2147483648}`)
	})

	t.Run("optional sint64", func(t *testing.T) {
		testMessageToJson(t, "optional sint64 optional_sint64_field = 1;", `{}`)
		testMessageToJson(t, "optional sint64 optional_sint64_field = 1;", `{"optionalSint64Field": "0"}`)
		testMessageToJson(t, "optional sint64 optional_sint64_field = 1;", `{"optionalSint64Field": "9223372036854775807"}`)
		testMessageToJson(t, "optional sint64 optional_sint64_field = 1;", `{"optionalSint64Field": "-9223372036854775808"}`)
	})

	t.Run("optional enum", func(t *testing.T) {
		testMessageToJson(t, "optional EnumType optional_enum_field = 1;", `{}`)
		testMessageToJson(t, "optional EnumType optional_enum_field = 1;", `{"optionalEnumField": "ENUM_TYPE_UNSPECIFIED"}`)
		testMessageToJson(t, "optional EnumType optional_enum_field = 1;", `{"optionalEnumField": "ENUM_TYPE_ONE"}`)
	})

	t.Run("optional message", func(t *testing.T) {
		testMessageToJson(t, "optional MessageType optional_message_field = 1;", `{}`)
		testMessageToJson(t, "optional MessageType optional_message_field = 1;", `{"optionalMessageField": {"value": 0}}`)
		testMessageToJson(t, "optional MessageType optional_message_field = 1;", `{"optionalMessageField": {"value": 12345}}`)
	})
}

func TestMessageToJsonMapKey(t *testing.T) {
	t.Run("map<int32, *>", func(t *testing.T) {
		testMessageToJson(t, "map<int32, string> int32_key_map_field = 1;", `{"int32KeyMapField": {}}`)
		testMessageToJson(t, "map<int32, string> int32_key_map_field = 1;", `{"int32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<uint32, *>", func(t *testing.T) {
		testMessageToJson(t, "map<uint32, string> uint32_key_map_field = 1;", `{"uint32KeyMapField": {}}`)
		testMessageToJson(t, "map<uint32, string> uint32_key_map_field = 1;", `{"uint32KeyMapField": {"0": "a", "4294967295": "b"}}`)
	})

	t.Run("map<int64, *>", func(t *testing.T) {
		testMessageToJson(t, "map<int64, string> int64_key_map_field = 1;", `{"int64KeyMapField": {}}`)
		testMessageToJson(t, "map<int64, string> int64_key_map_field = 1;", `{"int64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	t.Run("map<uint64, *>", func(t *testing.T) {
		testMessageToJson(t, "map<uint64, string> uint64_key_map_field = 1;", `{"uint64KeyMapField": {}}`)
		testMessageToJson(t, "map<uint64, string> uint64_key_map_field = 1;", `{"uint64KeyMapField": {"0": "a", "18446744073709551615": "b"}}`)
	})

	t.Run("map<fixed32, *>", func(t *testing.T) {
		testMessageToJson(t, "map<fixed32, string> fixed32_key_map_field = 1;", `{"fixed32KeyMapField": {}}`)
		testMessageToJson(t, "map<fixed32, string> fixed32_key_map_field = 1;", `{"fixed32KeyMapField": {"0": "a", "4294967295": "b"}}`)
	})

	t.Run("map<fixed64, *>", func(t *testing.T) {
		testMessageToJson(t, "map<fixed64, string> fixed64_key_map_field = 1;", `{"fixed64KeyMapField": {}}`)
		testMessageToJson(t, "map<fixed64, string> fixed64_key_map_field = 1;", `{"fixed64KeyMapField": {"0": "a", "18446744073709551615": "b"}}`)
	})

	t.Run("map<sfixed32, *>", func(t *testing.T) {
		testMessageToJson(t, "map<sfixed32, string> sfixed32_key_map_field = 1;", `{"sfixed32KeyMapField": {}}`)
		testMessageToJson(t, "map<sfixed32, string> sfixed32_key_map_field = 1;", `{"sfixed32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<sfixed64, *>", func(t *testing.T) {
		testMessageToJson(t, "map<sfixed64, string> sfixed64_key_map_field = 1;", `{"sfixed64KeyMapField": {}}`)
		testMessageToJson(t, "map<sfixed64, string> sfixed64_key_map_field = 1;", `{"sfixed64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	t.Run("map<bool, *>", func(t *testing.T) {
		testMessageToJson(t, "map<bool, string> bool_key_map_field = 1;", `{"boolKeyMapField": {}}`)
		testMessageToJson(t, "map<bool, string> bool_key_map_field = 1;", `{"boolKeyMapField": {"false": "a", "true": "b"}}`)
	})

	t.Run("map<string, *>", func(t *testing.T) {
		testMessageToJson(t, "map<string, string> string_key_map_field = 1;", `{"stringKeyMapField": {}}`)
		testMessageToJson(t, "map<string, string> string_key_map_field = 1;", `{"stringKeyMapField": {"a": "b", "c": "d"}}`)
	})

	t.Run("map<sint32, *>", func(t *testing.T) {
		testMessageToJson(t, "map<sint32, string> sint32_key_map_field = 1;", `{"sint32KeyMapField": {}}`)
		testMessageToJson(t, "map<sint32, string> sint32_key_map_field = 1;", `{"sint32KeyMapField": {"0": "a", "2147483647": "b", "-2147483648": "c"}}`)
	})

	t.Run("map<sint64, *>", func(t *testing.T) {
		testMessageToJson(t, "map<sint64, string> sint64_key_map_field = 1;", `{"sint64KeyMapField": {}}`)
		testMessageToJson(t, "map<sint64, string> sint64_key_map_field = 1;", `{"sint64KeyMapField": {"0": "a", "9223372036854775807": "b", "-9223372036854775808": "c"}}`)
	})

	// NOTE: float, double, enum, or message cannot be a map key.
}

func TestMessageToJsonMapValue(t *testing.T) {
	t.Run("map<*, int32>", func(t *testing.T) {
		testMessageToJson(t, "map<string, int32> int32_value_map_field = 1;", `{"int32ValueMapField": {}}`)
		testMessageToJson(t, "map<string, int32> int32_value_map_field = 1;", `{"int32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, uint32>", func(t *testing.T) {
		testMessageToJson(t, "map<string, uint32> uint32_value_map_field = 1;", `{"uint32ValueMapField": {}}`)
		testMessageToJson(t, "map<string, uint32> uint32_value_map_field = 1;", `{"uint32ValueMapField": {"a": 0, "b": 4294967295}}`)
	})

	t.Run("map<*, int64>", func(t *testing.T) {
		testMessageToJson(t, "map<string, int64> int64_value_map_field = 1;", `{"int64ValueMapField": {}}`)
		testMessageToJson(t, "map<string, int64> int64_value_map_field = 1;", `{"int64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`)
	})

	t.Run("map<*, uint64>", func(t *testing.T) {
		testMessageToJson(t, "map<string, uint64> uint64_value_map_field = 1;", `{"uint64ValueMapField": {}}`)
		testMessageToJson(t, "map<string, uint64> uint64_value_map_field = 1;", `{"uint64ValueMapField": {"a": "0", "b": "18446744073709551615"}}`)
	})

	t.Run("map<*, fixed32>", func(t *testing.T) {
		testMessageToJson(t, "map<string, fixed32> fixed32_value_map_field = 1;", `{"fixed32ValueMapField": {}}`)
		testMessageToJson(t, "map<string, fixed32> fixed32_value_map_field = 1;", `{"fixed32ValueMapField": {"a": 0, "b": 4294967295}}`)
	})

	t.Run("map<*, fixed64>", func(t *testing.T) {
		testMessageToJson(t, "map<string, fixed64> fixed64_value_map_field = 1;", `{"fixed64ValueMapField": {}}`)
		testMessageToJson(t, "map<string, fixed64> fixed64_value_map_field = 1;", `{"fixed64ValueMapField": {"a": "0", "b": "18446744073709551615"}}`)
	})

	t.Run("map<*, sfixed32>", func(t *testing.T) {
		testMessageToJson(t, "map<string, sfixed32> sfixed32_value_map_field = 1;", `{"sfixed32ValueMapField": {}}`)
		testMessageToJson(t, "map<string, sfixed32> sfixed32_value_map_field = 1;", `{"sfixed32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, sfixed64>", func(t *testing.T) {
		testMessageToJson(t, "map<string, sfixed64> sfixed64_value_map_field = 1;", `{"sfixed64ValueMapField": {}}`)
		testMessageToJson(t, "map<string, sfixed64> sfixed64_value_map_field = 1;", `{"sfixed64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`)
	})

	t.Run("map<*, bool>", func(t *testing.T) {
		testMessageToJson(t, "map<string, bool> bool_value_map_field = 1;", `{"boolValueMapField": {}}`)
		testMessageToJson(t, "map<string, bool> bool_value_map_field = 1;", `{"boolValueMapField": {"a": false, "b": true}}`)
	})

	t.Run("map<*, string>", func(t *testing.T) {
		testMessageToJson(t, "map<string, string> string_value_map_field = 1;", `{"stringValueMapField": {}}`)
		testMessageToJson(t, "map<string, string> string_value_map_field = 1;", `{"stringValueMapField": {"a": "b", "c": "d"}}`)
	})

	t.Run("map<*, bytes>", func(t *testing.T) {
		testMessageToJson(t, "map<string, bytes> bytes_value_map_field = 1;", `{"bytesValueMapField": {}}`)
		testMessageToJson(t, "map<string, bytes> bytes_value_map_field = 1;", `{"bytesValueMapField": {"a": "", "b": "dGVzdA=="}}`) // Base64 for "testMessageToJson"
	})

	t.Run("map<*, float>", func(t *testing.T) {
		testMessageToJson(t, "map<string, float> float_value_map_field = 1;", `{"floatValueMapField": {}}`)
		testMessageToJson(t, "map<string, float> float_value_map_field = 1;", `{"floatValueMapField": {"a": 0, "b": 3.5}}`)
	})

	t.Run("map<*, double>", func(t *testing.T) {
		testMessageToJson(t, "map<string, double> double_value_map_field = 1;", `{"doubleValueMapField": {}}`)
		testMessageToJson(t, "map<string, double> double_value_map_field = 1;", `{"doubleValueMapField": {"a": 0, "b": 3.141592653589793}}`)
	})

	t.Run("map<*, sint32>", func(t *testing.T) {
		testMessageToJson(t, "map<string, sint32> sint32_value_map_field = 1;", `{"sint32ValueMapField": {}}`)
		testMessageToJson(t, "map<string, sint32> sint32_value_map_field = 1;", `{"sint32ValueMapField": {"a": 0, "b": 2147483647, "c": -2147483648}}`)
	})

	t.Run("map<*, sint64>", func(t *testing.T) {
		testMessageToJson(t, "map<string, sint64> sint64_value_map_field = 1;", `{"sint64ValueMapField": {}}`)
		testMessageToJson(t, "map<string, sint64> sint64_value_map_field = 1;", `{"sint64ValueMapField": {"a": "0", "b": "9223372036854775807", "c": "-9223372036854775808"}}`)
	})

	t.Run("map<*, enum>", func(t *testing.T) {
		testMessageToJson(t, "map<string, EnumType> enum_value_map_field = 1;", `{"enumValueMapField": {}}`)
		testMessageToJson(t, "map<string, EnumType> enum_value_map_field = 1;", `{"enumValueMapField": {"a": "ENUM_TYPE_UNSPECIFIED", "b": "ENUM_TYPE_ONE"}}`)
	})

	t.Run("map<*, message>", func(t *testing.T) {
		testMessageToJson(t, "map<string, MessageType> message_value_map_field = 1;", `{"messageValueMapField": {}}`)
		testMessageToJson(t, "map<string, MessageType> message_value_map_field = 1;", `{"messageValueMapField": {"a": {"value": 0}, "b": {"value": 12345}}}`)
	})
}

func TestMessageToJsonOneof(t *testing.T) {
	t.Run("oneof", func(t *testing.T) {
		testMessageToJson(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{}`)
		testMessageToJson(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{"int32Field": 42}`)
		testMessageToJson(t, "oneof kind { int32 int32_field = 1; string string_field = 2; }", `{"stringField": "test"}`)
	})
}

func TestMessageToJsonWkt(t *testing.T) {
	t.Run("Timestamp", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"timestampField": "1970-01-01T00:00:00Z"}`)
		testMessageToJson(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"timestampField": "1970-01-01T00:00:01Z"}`)
		testMessageToJson(t, "google.protobuf.Timestamp timestamp_field = 1;", `{"timestampField": "2023-10-01T12:34:56.789Z"}`)
	})

	t.Run("Duration", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.Duration duration_field = 1;", `{"durationField": "0s"}`)
		testMessageToJson(t, "google.protobuf.Duration duration_field = 1;", `{"durationField": "1.234s"}`)
	})

	t.Run("Struct", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.Struct struct_field = 1;", `{"structField": {}}`)
		testMessageToJson(t, "google.protobuf.Struct struct_field = 1;", `{"structField": {"key": "value"}}`)
		testMessageToJson(t, "google.protobuf.Struct struct_field = 1;", `{"structField": {"number": 123, "boolean": true}}`)
	})

	t.Run("ListValue", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.ListValue list_value_field = 1;", `{"listValueField": []}`)
		testMessageToJson(t, "google.protobuf.ListValue list_value_field = 1;", `{"listValueField": ["string", 123, true]}`)
	})

	t.Run("Value", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.Value value_field = 1;", `{"valueField": {}}`)
		testMessageToJson(t, "google.protobuf.Value value_field = 1;", `{"valueField": "string"}`)
		testMessageToJson(t, "google.protobuf.Value value_field = 1;", `{"valueField": 123}`)
		testMessageToJson(t, "google.protobuf.Value value_field = 1;", `{"valueField": true}`)
	})

	t.Run("Empty", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.Empty empty_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.Empty empty_field = 1;", `{"emptyField": {}}`)
	})

	t.Run("DoubleValue", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.DoubleValue double_value_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.DoubleValue double_value_field = 1;", `{"doubleValueField": 0}`)
		testMessageToJson(t, "google.protobuf.DoubleValue double_value_field = 1;", `{"doubleValueField": 3.141592653589793}`)
	})

	t.Run("FloatValue", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.FloatValue float_value_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.FloatValue float_value_field = 1;", `{"floatValueField": 0}`)
		testMessageToJson(t, "google.protobuf.FloatValue float_value_field = 1;", `{"floatValueField": 3.5}`)
	})

	t.Run("Int64Value", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"int64ValueField": "0"}`)
		testMessageToJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"int64ValueField": "9223372036854775807"}`)
		testMessageToJson(t, "google.protobuf.Int64Value int64_value_field = 1;", `{"int64ValueField": "-9223372036854775808"}`)
	})

	t.Run("UInt64Value", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{"uint64ValueField": "0"}`)
		testMessageToJson(t, "google.protobuf.UInt64Value uint64_value_field = 1;", `{"uint64ValueField": "18446744073709551615"}`)
	})

	t.Run("Int32Value", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"int32ValueField": 0}`)
		testMessageToJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"int32ValueField": 2147483647}`)
		testMessageToJson(t, "google.protobuf.Int32Value int32_value_field = 1;", `{"int32ValueField": -2147483648}`)
	})

	t.Run("UInt32Value", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{"uint32ValueField": 0}`)
		testMessageToJson(t, "google.protobuf.UInt32Value uint32_value_field = 1;", `{"uint32ValueField": 4294967295}`)
	})

	t.Run("BoolValue", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.BoolValue bool_value_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.BoolValue bool_value_field = 1;", `{"boolValueField": false}`)
		testMessageToJson(t, "google.protobuf.BoolValue bool_value_field = 1;", `{"boolValueField": true}`)
	})

	t.Run("StringValue", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.StringValue string_value_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.StringValue string_value_field = 1;", `{"stringValueField": ""}`)
		testMessageToJson(t, "google.protobuf.StringValue string_value_field = 1;", `{"stringValueField": "test"}`)
	})

	t.Run("BytesValue", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{"bytesValueField": ""}`)
		testMessageToJson(t, "google.protobuf.BytesValue bytes_value_field = 1;", `{"bytesValueField": "dGVzdA=="}`) // Base64 for "test"
	})

	t.Run("FieldMask", func(t *testing.T) {
		testMessageToJson(t, "google.protobuf.FieldMask field_mask_field = 1;", `{}`)
		testMessageToJson(t, "google.protobuf.FieldMask field_mask_field = 1;", `{"fieldMaskField": ""}`)
		testMessageToJson(t, "google.protobuf.FieldMask field_mask_field = 1;", `{"fieldMaskField": "path1,path2"}`)
	})
}

func TestMessageToJsonNullInput(t *testing.T) {
	p := testutils.NewProtoTestSupport(t, map[string]string{
		"main.proto": `
			syntax = "proto3";
			message Test {
				int32 value = 1;
			}
		`,
	})

	descriptorSetName := "a"
	typeName := ".Test"

	AssertThatCall(t, "pb_descriptor_set_load(?, ?)", descriptorSetName, p.GetSerializedFileDescriptorSet()).ShouldSucceed()
	defer func() {
		AssertThatCall(t, "pb_descriptor_set_delete(?)", descriptorSetName).ShouldSucceed()
	}()

	RunTestThatExpression(t, "pb_message_to_json(?, ?, ?)", descriptorSetName, typeName, nil).IsNull()
}
