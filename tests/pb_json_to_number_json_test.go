package main

import (
	"fmt"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	"github.com/eiiches/mysql-protobuf-functions/internal/descriptorsetjson"
	"github.com/eiiches/mysql-protobuf-functions/internal/protonumberjson"
	"github.com/eiiches/mysql-protobuf-functions/internal/testutils"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func testJsonToNumberJson(t *testing.T, fieldDefinition string, input string) {
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

	// Get expected ProtoNumberJSON using canonical Go implementation
	expectedProtoNumberJson, err := protonumberjson.Marshal(dynamicMessage.Interface())
	g.Expect(err).NotTo(HaveOccurred())

	// Test our MySQL function against canonical implementation
	RunTestThatExpression(t, "_pb_json_to_number_json(?, ?, ?)", descriptorSetJson, typeName, input).IsEqualToJsonString(string(expectedProtoNumberJson))
}

func TestJsonToNumberJsonSingularFields(t *testing.T) {
	t.Run("int32", func(t *testing.T) {
		testJsonToNumberJson(t, "int32 int32_field = 1;", `{"int32Field": 0}`)
		testJsonToNumberJson(t, "int32 int32_field = 1;", `{"int32Field": 2147483647}`)
		testJsonToNumberJson(t, "int32 int32_field = 1;", `{"int32Field": -2147483648}`)
	})

	t.Run("uint32", func(t *testing.T) {
		testJsonToNumberJson(t, "uint32 uint32_field = 1;", `{"uint32Field": 0}`)
		testJsonToNumberJson(t, "uint32 uint32_field = 1;", `{"uint32Field": 4294967295}`)
	})

	t.Run("int64", func(t *testing.T) {
		testJsonToNumberJson(t, "int64 int64_field = 1;", `{"int64Field": "0"}`)
		testJsonToNumberJson(t, "int64 int64_field = 1;", `{"int64Field": "9223372036854775807"}`)
		testJsonToNumberJson(t, "int64 int64_field = 1;", `{"int64Field": "-9223372036854775808"}`)
	})

	t.Run("uint64", func(t *testing.T) {
		testJsonToNumberJson(t, "uint64 uint64_field = 1;", `{"uint64Field": "0"}`)
		testJsonToNumberJson(t, "uint64 uint64_field = 1;", `{"uint64Field": "18446744073709551615"}`)
	})

	t.Run("bool", func(t *testing.T) {
		testJsonToNumberJson(t, "bool bool_field = 1;", `{"boolField": false}`)
		testJsonToNumberJson(t, "bool bool_field = 1;", `{"boolField": true}`)
	})

	t.Run("string", func(t *testing.T) {
		testJsonToNumberJson(t, "string string_field = 1;", `{"stringField": ""}`)
		testJsonToNumberJson(t, "string string_field = 1;", `{"stringField": "test"}`)
	})

	t.Run("bytes", func(t *testing.T) {
		testJsonToNumberJson(t, "bytes bytes_field = 1;", `{"bytesField": ""}`)
		testJsonToNumberJson(t, "bytes bytes_field = 1;", `{"bytesField": "aGVsbG8="}`)
	})

	t.Run("float", func(t *testing.T) {
		testJsonToNumberJson(t, "float float_field = 1;", `{"floatField": 0}`)
		testJsonToNumberJson(t, "float float_field = 1;", `{"floatField": 3.5}`)
	})

	t.Run("double", func(t *testing.T) {
		testJsonToNumberJson(t, "double double_field = 1;", `{"doubleField": 0}`)
		testJsonToNumberJson(t, "double double_field = 1;", `{"doubleField": 3.141592653589793}`)
	})

	t.Run("enum", func(t *testing.T) {
		testJsonToNumberJson(t, "EnumType enum_field = 1;", `{"enumField": "ENUM_TYPE_UNSPECIFIED"}`)
		testJsonToNumberJson(t, "EnumType enum_field = 1;", `{"enumField": "ENUM_TYPE_ONE"}`)
	})

	t.Run("message", func(t *testing.T) {
		testJsonToNumberJson(t, "MessageType message_field = 1;", `{}`)
		testJsonToNumberJson(t, "MessageType message_field = 1;", `{"messageField": {"value": 0}}`)
		testJsonToNumberJson(t, "MessageType message_field = 1;", `{"messageField": {"value": 12345}}`)
	})

	t.Run("repeated int32", func(t *testing.T) {
		testJsonToNumberJson(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": []}`)
		testJsonToNumberJson(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": [0]}`)
		testJsonToNumberJson(t, "repeated int32 repeated_int32_field = 1;", `{"repeatedInt32Field": [-2147483648, 0, 2147483647]}`)
	})

	t.Run("repeated enum", func(t *testing.T) {
		testJsonToNumberJson(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": []}`)
		testJsonToNumberJson(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": ["ENUM_TYPE_UNSPECIFIED"]}`)
		testJsonToNumberJson(t, "repeated EnumType repeated_enum_field = 1;", `{"repeatedEnumField": ["ENUM_TYPE_ONE", "ENUM_TYPE_UNSPECIFIED"]}`)
	})

	t.Run("optional int32", func(t *testing.T) {
		testJsonToNumberJson(t, "optional int32 optional_int32_field = 1;", `{}`)
		testJsonToNumberJson(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": 0}`)
		testJsonToNumberJson(t, "optional int32 optional_int32_field = 1;", `{"optionalInt32Field": 42}`)
	})
}
