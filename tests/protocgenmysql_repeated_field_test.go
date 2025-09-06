package main

import (
	"fmt"
	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	"testing"
)

func TestProtocGenRepeatedField(t *testing.T) {
	protoContent := dedent.Pipe(`
		|syntax = "proto3";
		|message Test {
		|    repeated double repeated_double = 1;
		|    repeated float repeated_float = 2;
		|    repeated int32 repeated_int32 = 3;
		|    repeated int64 repeated_int64 = 4;
		|    repeated uint32 repeated_uint32 = 5;
		|    repeated uint64 repeated_uint64 = 6;
		|    repeated sint32 repeated_sint32 = 7;
		|    repeated sint64 repeated_sint64 = 8;
		|    repeated fixed32 repeated_fixed32 = 9;
		|    repeated fixed64 repeated_fixed64 = 10;
		|    repeated sfixed32 repeated_sfixed32 = 11;
		|    repeated sfixed64 repeated_sfixed64 = 12;
		|    repeated bool repeated_bool = 13;
		|    repeated string repeated_string = 14;
		|    repeated bytes repeated_bytes = 15;
		|    repeated Status repeated_enum = 16;
		|    repeated Nested repeated_message = 17;
		|}
		|message Nested {
		|    string name = 1;
		|    int32 value = 2;
		|}
		|enum Status {
		|    STATUS_UNSPECIFIED = 0;
		|    STATUS_ACTIVE = 1;
		|    STATUS_INACTIVE = 2;
		|}
	`)
	schemaName := "test_schema"
	generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

	// Test repeated double field (IEEE 754 binary64 format in arrays)
	t.Run("repeated_double_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_double(?, 3.141592653589793)", `{}`).IsEqualToJsonString(`{"1": ["binary64:0x400921fb54442d18"]}`)
		RunTestThatExpression(t, "test_add_repeated_double(?, 1.0)", `{"1": ["binary64:0x400921fb54442d18"]}`).IsEqualToJsonString(`{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_double('{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}')`).IsEqualToJsonString(`[3.141592653589793, 1.0]`)
		RunTestThatExpression(t, `test_get_repeated_double('{"1": []}')`).IsEqualToJsonString(`[]`) // Empty array
		RunTestThatExpression(t, `test_get_repeated_double('{}')`).IsEqualToJsonString(`[]`)        // Missing field returns empty array
	})

	// Test repeated float field (IEEE 754 binary32 format in arrays)
	t.Run("repeated_float_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_float(?, 3.14)", `{}`).IsEqualToJsonString(`{"2": ["binary32:0x4048f5c3"]}`)
		RunTestThatExpression(t, "test_add_repeated_float(?, 1.0)", `{"2": ["binary32:0x4048f5c3"]}`).IsEqualToJsonString(`{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_float('{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}')`).IsEqualToJsonString(`[3.140000104904175, 1.0]`) // CAST(CAST(3.14 AS FLOAT) AS DOUBLE) => 3.140000104904175
		RunTestThatExpression(t, `test_get_repeated_float('{"2": []}')`).IsEqualToJsonString(`[]`)                                                                   // Empty array
		RunTestThatExpression(t, `test_get_repeated_float('{}')`).IsEqualToJsonString(`[]`)                                                                          // Missing field returns empty array
	})

	// Test repeated int32 field
	t.Run("repeated_int32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_int32(?, 42)", `{}`).IsEqualToJsonString(`{"3": [42]}`)
		RunTestThatExpression(t, "test_add_repeated_int32(?, -2147483648)", `{"3": [42]}`).IsEqualToJsonString(`{"3": [42, -2147483648]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_int32('{"3": [42, -2147483648]}')`).IsEqualToJsonString(`[42, -2147483648]`)
		RunTestThatExpression(t, `test_get_repeated_int32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated int64 field
	t.Run("repeated_int64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_int64(?, 9223372036854775807)", `{}`).IsEqualToJsonString(`{"4": [9223372036854775807]}`)
		RunTestThatExpression(t, "test_add_repeated_int64(?, -9223372036854775808)", `{"4": [9223372036854775807]}`).IsEqualToJsonString(`{"4": [9223372036854775807, -9223372036854775808]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_int64('{"4": [9223372036854775807, -9223372036854775808]}')`).IsEqualToJsonString(`[9223372036854775807, -9223372036854775808]`)
		RunTestThatExpression(t, `test_get_repeated_int64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated uint32 field
	t.Run("repeated_uint32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_uint32(?, 4294967295)", `{}`).IsEqualToJsonString(`{"5": [4294967295]}`)
		RunTestThatExpression(t, "test_add_repeated_uint32(?, 42)", `{"5": [4294967295]}`).IsEqualToJsonString(`{"5": [4294967295, 42]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_uint32('{"5": [4294967295, 42]}')`).IsEqualToJsonString(`[4294967295, 42]`)
		RunTestThatExpression(t, `test_get_repeated_uint32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated uint64 field
	t.Run("repeated_uint64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_uint64(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"6": [18446744073709551615]}`)
		RunTestThatExpression(t, "test_add_repeated_uint64(?, 100)", `{"6": [18446744073709551615]}`).IsEqualToJsonString(`{"6": [18446744073709551615, 100]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_uint64('{"6": [18446744073709551615, 100]}')`).IsEqualToJsonString(`[18446744073709551615, 100]`)
		RunTestThatExpression(t, `test_get_repeated_uint64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated sint32 field
	t.Run("repeated_sint32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_sint32(?, -1)", `{}`).IsEqualToJsonString(`{"7": [-1]}`)
		RunTestThatExpression(t, "test_add_repeated_sint32(?, 42)", `{"7": [-1]}`).IsEqualToJsonString(`{"7": [-1, 42]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_sint32('{"7": [-1, 42]}')`).IsEqualToJsonString(`[-1, 42]`)
		RunTestThatExpression(t, `test_get_repeated_sint32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated sint64 field
	t.Run("repeated_sint64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_sint64(?, -1)", `{}`).IsEqualToJsonString(`{"8": [-1]}`)
		RunTestThatExpression(t, "test_add_repeated_sint64(?, 100)", `{"8": [-1]}`).IsEqualToJsonString(`{"8": [-1, 100]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_sint64('{"8": [-1, 100]}')`).IsEqualToJsonString(`[-1, 100]`)
		RunTestThatExpression(t, `test_get_repeated_sint64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated fixed32 field
	t.Run("repeated_fixed32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_fixed32(?, 4294967295)", `{}`).IsEqualToJsonString(`{"9": [4294967295]}`)
		RunTestThatExpression(t, "test_add_repeated_fixed32(?, 42)", `{"9": [4294967295]}`).IsEqualToJsonString(`{"9": [4294967295, 42]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_fixed32('{"9": [4294967295, 42]}')`).IsEqualToJsonString(`[4294967295, 42]`)
		RunTestThatExpression(t, `test_get_repeated_fixed32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated fixed64 field
	t.Run("repeated_fixed64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_fixed64(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"10": [18446744073709551615]}`)
		RunTestThatExpression(t, "test_add_repeated_fixed64(?, 100)", `{"10": [18446744073709551615]}`).IsEqualToJsonString(`{"10": [18446744073709551615, 100]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_fixed64('{"10": [18446744073709551615, 100]}')`).IsEqualToJsonString(`[18446744073709551615, 100]`)
		RunTestThatExpression(t, `test_get_repeated_fixed64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated sfixed32 field
	t.Run("repeated_sfixed32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_sfixed32(?, -2147483648)", `{}`).IsEqualToJsonString(`{"11": [-2147483648]}`)
		RunTestThatExpression(t, "test_add_repeated_sfixed32(?, 42)", `{"11": [-2147483648]}`).IsEqualToJsonString(`{"11": [-2147483648, 42]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_sfixed32('{"11": [-2147483648, 42]}')`).IsEqualToJsonString(`[-2147483648, 42]`)
		RunTestThatExpression(t, `test_get_repeated_sfixed32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated sfixed64 field
	t.Run("repeated_sfixed64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_sfixed64(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"12": [-9223372036854775808]}`)
		RunTestThatExpression(t, "test_add_repeated_sfixed64(?, 100)", `{"12": [-9223372036854775808]}`).IsEqualToJsonString(`{"12": [-9223372036854775808, 100]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_repeated_sfixed64('{"12": [-9223372036854775808, 100]}')`).IsEqualToJsonString(`[-9223372036854775808, 100]`)
		RunTestThatExpression(t, `test_get_repeated_sfixed64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated bool field (JSON booleans, not 1/0)
	t.Run("repeated_bool_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_bool(?, TRUE)", `{}`).IsEqualToJsonString(`{"13": [true]}`)
		RunTestThatExpression(t, "test_add_repeated_bool(?, FALSE)", `{"13": [true]}`).IsEqualToJsonString(`{"13": [true, false]}`)

		// Test get operations return actual boolean arrays
		RunTestThatExpression(t, `test_get_repeated_bool('{"13": [true, false]}')`).IsEqualToJsonString(`[true, false]`)
		RunTestThatExpression(t, `test_get_repeated_bool('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated string field
	t.Run("repeated_string_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_string(?, 'hello')", `{}`).IsEqualToJsonString(`{"14": ["hello"]}`)
		RunTestThatExpression(t, "test_add_repeated_string(?, 'world')", `{"14": ["hello"]}`).IsEqualToJsonString(`{"14": ["hello", "world"]}`)

		// Test get operations return actual string arrays
		RunTestThatExpression(t, `test_get_repeated_string('{"14": ["hello", "world"]}')`).IsEqualToJsonString(`["hello", "world"]`)
		RunTestThatExpression(t, `test_get_repeated_string('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated bytes field (base64 encoded)
	t.Run("repeated_bytes_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_bytes(?, ?)", `{}`, []byte("hello")).IsEqualToJsonString(`{"15": ["aGVsbG8="]}`)
		RunTestThatExpression(t, "test_add_repeated_bytes(?, ?)", `{"15": ["aGVsbG8="]}`, []byte("world")).IsEqualToJsonString(`{"15": ["aGVsbG8=", "d29ybGQ="]}`)

		// Test get operations convert from base64 back to actual byte arrays (but returned as base64 JSON)
		RunTestThatExpression(t, `test_get_repeated_bytes('{"15": ["aGVsbG8=", "d29ybGQ="]}')`).IsEqualToJsonString(`["aGVsbG8=", "d29ybGQ="]`) // Returns base64 strings in array
		RunTestThatExpression(t, `test_get_repeated_bytes('{}')`).IsEqualToJsonString(`[]`)                                                     // Missing field returns empty array
	})

	// Test repeated enum field (numbers, not string names)
	t.Run("repeated_enum_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_enum(?, 1)", `{}`).IsEqualToJsonString(`{"16": [1]}`)
		RunTestThatExpression(t, "test_add_repeated_enum(?, 2)", `{"16": [1]}`).IsEqualToJsonString(`{"16": [1, 2]}`)

		// Test get operations return actual integer arrays
		RunTestThatExpression(t, `test_get_repeated_enum('{"16": [1, 2]}')`).IsEqualToJsonString(`[1, 2]`)
		RunTestThatExpression(t, `test_get_repeated_enum('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})

	// Test repeated message field (array of nested objects with field number keys)
	t.Run("repeated_message_field", func(t *testing.T) {
		// Test add operations create correct internal format
		nested1 := "nested_set_name(nested_new(), 'first')"
		nested2 := "nested_set_value(nested_set_name(nested_new(), 'second'), 42)"
		RunTestThatExpression(t, fmt.Sprintf("test_add_repeated_message(?, %s)", nested1), `{}`).IsEqualToJsonString(`{"17": [{"1": "first"}]}`)
		RunTestThatExpression(t, fmt.Sprintf("test_add_repeated_message(?, %s)", nested2), `{"17": [{"1": "first"}]}`).IsEqualToJsonString(`{"17": [{"1": "first"}, {"1": "second", "2": 42}]}`)

		// Test get operations return actual nested object arrays
		RunTestThatExpression(t, `test_get_repeated_message('{"17": [{"1": "first"}, {"1": "second", "2": 42}]}')`).IsEqualToJsonString(`[{"1": "first"}, {"1": "second", "2": 42}]`)
		RunTestThatExpression(t, `test_get_repeated_message('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array
	})
}
