package main

import (
	"fmt"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
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
		RunTestThatExpression(t, `test_get_all_repeated_double('{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}')`).IsEqualToJsonString(`[3.141592653589793, 1.0]`)
		RunTestThatExpression(t, `test_get_all_repeated_double('{"1": []}')`).IsEqualToJsonString(`[]`) // Empty array
		RunTestThatExpression(t, `test_get_all_repeated_double('{}')`).IsEqualToJsonString(`[]`)        // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_double(?, '[3.141592653589793, 1.0]')`, `{}`).IsEqualToJsonString(`{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}`)
		RunTestThatExpression(t, `test_set_all_repeated_double(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_double('{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_double('{}')`).IsEqualToInt(0)                                                                    // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_double('{"1": []}')`).IsEqualToInt(0)                                                             // Empty array
		RunTestThatExpression(t, `test_count_repeated_double('{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_double('{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}', 0)`).IsEqualToFloat(3.141592653589793)
		RunTestThatExpression(t, `test_get_repeated_double('{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}', 1)`).IsEqualToFloat(1.0)
		RunTestThatExpression(t, `test_get_repeated_double('{"1": ["binary64:0x400921fb54442d18"]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_double('{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}', 0, 2.5)`).IsEqualToJsonString(`{"1": ["binary64:0x4004000000000000", "binary64:0x3ff0000000000000"]}`)
		RunTestThatExpression(t, `test_insert_repeated_double('{"1": ["binary64:0x3ff0000000000000"]}', 0, 3.141592653589793)`).IsEqualToJsonString(`{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}`)
		RunTestThatExpression(t, `test_remove_repeated_double('{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000"]}', 0)`).IsEqualToJsonString(`{"1": ["binary64:0x3ff0000000000000"]}`)
		RunTestThatExpression(t, `test_add_all_repeated_double('{"1": ["binary64:0x400921fb54442d18"]}', '[1.0, 2.5]')`).IsEqualToJsonString(`{"1": ["binary64:0x400921fb54442d18", "binary64:0x3ff0000000000000", "binary64:0x4004000000000000"]}`)
	})

	// Test repeated float field (IEEE 754 binary32 format in arrays)
	t.Run("repeated_float_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_float(?, 3.14)", `{}`).IsEqualToJsonString(`{"2": ["binary32:0x4048f5c3"]}`)
		RunTestThatExpression(t, "test_add_repeated_float(?, 1.0)", `{"2": ["binary32:0x4048f5c3"]}`).IsEqualToJsonString(`{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_float('{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}')`).IsEqualToJsonString(`[3.140000104904175, 1.0]`) // CAST(CAST(3.14 AS FLOAT) AS DOUBLE) => 3.140000104904175
		RunTestThatExpression(t, `test_get_all_repeated_float('{"2": []}')`).IsEqualToJsonString(`[]`)                                                                   // Empty array
		RunTestThatExpression(t, `test_get_all_repeated_float('{}')`).IsEqualToJsonString(`[]`)                                                                          // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_float(?, '[3.14, 1.0]')`, `{}`).IsEqualToJsonString(`{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}`)
		RunTestThatExpression(t, `test_set_all_repeated_float(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_float('{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_float('{}')`).IsEqualToInt(0)                                                    // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_float('{"2": []}')`).IsEqualToInt(0)                                             // Empty array
		RunTestThatExpression(t, `test_count_repeated_float('{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_float('{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}', 0)`).IsEqualToFloat(3.140000104904175) // CAST(CAST(3.14 AS FLOAT) AS DOUBLE)
		RunTestThatExpression(t, `test_get_repeated_float('{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}', 1)`).IsEqualToFloat(1.0)
		RunTestThatExpression(t, `test_get_repeated_float('{"2": ["binary32:0x4048f5c3"]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_float('{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}', 1, 2.5)`).IsEqualToJsonString(`{"2": ["binary32:0x4048f5c3", "binary32:0x40200000"]}`)
		RunTestThatExpression(t, `test_insert_repeated_float('{"2": ["binary32:0x3f800000"]}', 0, 3.14)`).IsEqualToJsonString(`{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}`)
		RunTestThatExpression(t, `test_remove_repeated_float('{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000"]}', 0)`).IsEqualToJsonString(`{"2": ["binary32:0x3f800000"]}`)
		RunTestThatExpression(t, `test_add_all_repeated_float('{"2": ["binary32:0x4048f5c3"]}', '[1.0, 2.5]')`).IsEqualToJsonString(`{"2": ["binary32:0x4048f5c3", "binary32:0x3f800000", "binary32:0x40200000"]}`)
	})

	// Test repeated int32 field
	t.Run("repeated_int32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_int32(?, 42)", `{}`).IsEqualToJsonString(`{"3": [42]}`)
		RunTestThatExpression(t, "test_add_repeated_int32(?, -2147483648)", `{"3": [42]}`).IsEqualToJsonString(`{"3": [42, -2147483648]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_int32('{"3": [42, -2147483648]}')`).IsEqualToJsonString(`[42, -2147483648]`)
		RunTestThatExpression(t, `test_get_all_repeated_int32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_int32(?, '[42, -2147483648]')`, `{}`).IsEqualToJsonString(`{"3": [42, -2147483648]}`)
		RunTestThatExpression(t, `test_set_all_repeated_int32(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_int32('{"3": [42, -2147483648]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_int32('{}')`).IsEqualToInt(0)                       // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_int32('{"3": []}')`).IsEqualToInt(0)                // Empty array
		RunTestThatExpression(t, `test_count_repeated_int32('{"3": [42, -2147483648]}')`).IsEqualToInt(2) // Two elements

		// Test index-based get operations
		RunTestThatExpression(t, `test_get_repeated_int32('{"3": [10, 20, 30]}', 0)`).IsEqualToInt(10)
		RunTestThatExpression(t, `test_get_repeated_int32('{"3": [10, 20, 30]}', 1)`).IsEqualToInt(20)
		RunTestThatExpression(t, `test_get_repeated_int32('{"3": [10, 20, 30]}', 2)`).IsEqualToInt(30)
		// Test out of bounds index returns NULL
		RunTestThatExpression(t, `test_get_repeated_int32('{"3": [10, 20, 30]}', 3)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_get_repeated_int32('{"3": [10, 20, 30]}', -1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		// Test empty array returns NULL
		RunTestThatExpression(t, `test_get_repeated_int32('{}', 0)`).ToFailWithSignalException("45000", "Array index out of bounds")

		// Test index-based set operations
		RunTestThatExpression(t, `test_set_repeated_int32('{"3": [10, 20, 30]}', 0, 100)`).IsEqualToJsonString(`{"3": [100, 20, 30]}`)
		RunTestThatExpression(t, `test_set_repeated_int32('{"3": [10, 20, 30]}', 1, 200)`).IsEqualToJsonString(`{"3": [10, 200, 30]}`)
		RunTestThatExpression(t, `test_set_repeated_int32('{"3": [10, 20, 30]}', 2, 300)`).IsEqualToJsonString(`{"3": [10, 20, 300]}`)
		// Test out of bounds index should fail
		RunTestThatExpression(t, `test_set_repeated_int32('{"3": [10, 20, 30]}', 3, 400)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_int32('{}', 0, 100)`).ToFailWithSignalException("45000", "Array index out of bounds")

		// Test insert operations
		RunTestThatExpression(t, `test_insert_repeated_int32('{"3": [20, 30]}', 0, 10)`).IsEqualToJsonString(`{"3": [10, 20, 30]}`)
		RunTestThatExpression(t, `test_insert_repeated_int32('{"3": [10, 30]}', 1, 20)`).IsEqualToJsonString(`{"3": [10, 20, 30]}`)
		RunTestThatExpression(t, `test_insert_repeated_int32('{"3": [10, 20]}', 2, 30)`).IsEqualToJsonString(`{"3": [10, 20, 30]}`)
		RunTestThatExpression(t, `test_insert_repeated_int32('{}', 0, 100)`).IsEqualToJsonString(`{"3": [100]}`)

		// Test insert bounds checking
		RunTestThatExpression(t, `test_insert_repeated_int32('{"3": [10, 20]}', -1, 5)`).ToFailWithSignalException("45000", "Insert index out of bounds")
		RunTestThatExpression(t, `test_insert_repeated_int32('{"3": [10, 20]}', 3, 30)`).ToFailWithSignalException("45000", "Insert index out of bounds")
		RunTestThatExpression(t, `test_insert_repeated_int32('{}', 1, 100)`).ToFailWithSignalException("45000", "Insert index out of bounds")

		// Test remove operations
		RunTestThatExpression(t, `test_remove_repeated_int32('{"3": [10, 20, 30]}', 0)`).IsEqualToJsonString(`{"3": [20, 30]}`)
		RunTestThatExpression(t, `test_remove_repeated_int32('{"3": [10, 20, 30]}', 1)`).IsEqualToJsonString(`{"3": [10, 30]}`)
		RunTestThatExpression(t, `test_remove_repeated_int32('{"3": [10, 20, 30]}', 2)`).IsEqualToJsonString(`{"3": [10, 20]}`)
		RunTestThatExpression(t, `test_remove_repeated_int32('{"3": [10]}', 0)`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `test_remove_repeated_int32('{"3": [10, 20, 30]}', 3)`).ToFailWithSignalException("45000", "Array index out of bounds")

		// Test add_all operations
		RunTestThatExpression(t, `test_add_all_repeated_int32('{}', '[100, 200, 300]')`).IsEqualToJsonString(`{"3": [100, 200, 300]}`)
		RunTestThatExpression(t, `test_add_all_repeated_int32('{"3": [10, 20]}', '[100, 200]')`).IsEqualToJsonString(`{"3": [10, 20, 100, 200]}`)
		RunTestThatExpression(t, `test_add_all_repeated_int32('{"3": [10, 20]}', '[]')`).IsEqualToJsonString(`{"3": [10, 20]}`)
	})

	// Test repeated int64 field
	t.Run("repeated_int64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_int64(?, 9223372036854775807)", `{}`).IsEqualToJsonString(`{"4": [9223372036854775807]}`)
		RunTestThatExpression(t, "test_add_repeated_int64(?, -9223372036854775808)", `{"4": [9223372036854775807]}`).IsEqualToJsonString(`{"4": [9223372036854775807, -9223372036854775808]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_int64('{"4": [9223372036854775807, -9223372036854775808]}')`).IsEqualToJsonString(`[9223372036854775807, -9223372036854775808]`)
		RunTestThatExpression(t, `test_get_all_repeated_int64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_int64(?, '[9223372036854775807, -9223372036854775808]')`, `{}`).IsEqualToJsonString(`{"4": [9223372036854775807, -9223372036854775808]}`)
		RunTestThatExpression(t, `test_set_all_repeated_int64(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_int64('{"4": [9223372036854775807, -9223372036854775808]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_int64('{}')`).IsEqualToInt(0)                                                 // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_int64('{"4": []}')`).IsEqualToInt(0)                                          // Empty array
		RunTestThatExpression(t, `test_count_repeated_int64('{"4": [9223372036854775807, -9223372036854775808]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_int64('{"4": [9223372036854775807, -9223372036854775808, 42]}', 0)`).IsEqualToInt(9223372036854775807)
		RunTestThatExpression(t, `test_get_repeated_int64('{"4": [9223372036854775807, -9223372036854775808, 42]}', 2)`).IsEqualToInt(42)
		RunTestThatExpression(t, `test_get_repeated_int64('{"4": [100]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_int64('{"4": [100, 200]}', 1, -500)`).IsEqualToJsonString(`{"4": [100, -500]}`)
		RunTestThatExpression(t, `test_insert_repeated_int64('{"4": [200]}', 0, 100)`).IsEqualToJsonString(`{"4": [100, 200]}`)
		RunTestThatExpression(t, `test_remove_repeated_int64('{"4": [100, 200, 300]}', 1)`).IsEqualToJsonString(`{"4": [100, 300]}`)
		RunTestThatExpression(t, `test_add_all_repeated_int64('{"4": [100]}', '[200, 300]')`).IsEqualToJsonString(`{"4": [100, 200, 300]}`)
	})

	// Test repeated uint32 field
	t.Run("repeated_uint32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_uint32(?, 4294967295)", `{}`).IsEqualToJsonString(`{"5": [4294967295]}`)
		RunTestThatExpression(t, "test_add_repeated_uint32(?, 42)", `{"5": [4294967295]}`).IsEqualToJsonString(`{"5": [4294967295, 42]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_uint32('{"5": [4294967295, 42]}')`).IsEqualToJsonString(`[4294967295, 42]`)
		RunTestThatExpression(t, `test_get_all_repeated_uint32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_uint32(?, '[4294967295, 42]')`, `{}`).IsEqualToJsonString(`{"5": [4294967295, 42]}`)
		RunTestThatExpression(t, `test_set_all_repeated_uint32(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_uint32('{"5": [4294967295, 42]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_uint32('{}')`).IsEqualToInt(0)                      // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_uint32('{"5": []}')`).IsEqualToInt(0)               // Empty array
		RunTestThatExpression(t, `test_count_repeated_uint32('{"5": [4294967295, 42]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_uint32('{"5": [4294967295, 42, 100]}', 1)`).IsEqualToInt(42)
		RunTestThatExpression(t, `test_get_repeated_uint32('{"5": [100]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_uint32('{"5": [100, 200]}', 0, 999)`).IsEqualToJsonString(`{"5": [999, 200]}`)
		RunTestThatExpression(t, `test_insert_repeated_uint32('{"5": [200]}', 0, 100)`).IsEqualToJsonString(`{"5": [100, 200]}`)
		RunTestThatExpression(t, `test_remove_repeated_uint32('{"5": [100, 200, 300]}', 1)`).IsEqualToJsonString(`{"5": [100, 300]}`)
		RunTestThatExpression(t, `test_add_all_repeated_uint32('{"5": [100]}', '[200, 300]')`).IsEqualToJsonString(`{"5": [100, 200, 300]}`)
	})

	// Test repeated uint64 field
	t.Run("repeated_uint64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_uint64(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"6": [18446744073709551615]}`)
		RunTestThatExpression(t, "test_add_repeated_uint64(?, 100)", `{"6": [18446744073709551615]}`).IsEqualToJsonString(`{"6": [18446744073709551615, 100]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_uint64('{"6": [18446744073709551615, 100]}')`).IsEqualToJsonString(`[18446744073709551615, 100]`)
		RunTestThatExpression(t, `test_get_all_repeated_uint64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_uint64(?, '[18446744073709551615, 100]')`, `{}`).IsEqualToJsonString(`{"6": [18446744073709551615, 100]}`)
		RunTestThatExpression(t, `test_set_all_repeated_uint64(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_uint64('{"6": [18446744073709551615, 100]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_uint64('{}')`).IsEqualToInt(0)                                 // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_uint64('{"6": []}')`).IsEqualToInt(0)                          // Empty array
		RunTestThatExpression(t, `test_count_repeated_uint64('{"6": [18446744073709551615, 100]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_uint64('{"6": [18446744073709551615, 100, 42]}', 1)`).IsEqualToInt(100)
		RunTestThatExpression(t, `test_get_repeated_uint64('{"6": [100]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_uint64('{"6": [100, 200]}', 0, 999)`).IsEqualToJsonString(`{"6": [999, 200]}`)
		RunTestThatExpression(t, `test_insert_repeated_uint64('{"6": [100]}', 0, 50)`).IsEqualToJsonString(`{"6": [50, 100]}`)
		RunTestThatExpression(t, `test_remove_repeated_uint64('{"6": [100, 200, 300]}', 1)`).IsEqualToJsonString(`{"6": [100, 300]}`)
		RunTestThatExpression(t, `test_add_all_repeated_uint64('{"6": [100]}', '[200, 300]')`).IsEqualToJsonString(`{"6": [100, 200, 300]}`)
	})

	// Test repeated sint32 field
	t.Run("repeated_sint32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_sint32(?, -1)", `{}`).IsEqualToJsonString(`{"7": [-1]}`)
		RunTestThatExpression(t, "test_add_repeated_sint32(?, 42)", `{"7": [-1]}`).IsEqualToJsonString(`{"7": [-1, 42]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_sint32('{"7": [-1, 42]}')`).IsEqualToJsonString(`[-1, 42]`)
		RunTestThatExpression(t, `test_get_all_repeated_sint32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_sint32(?, '[-1, 42]')`, `{}`).IsEqualToJsonString(`{"7": [-1, 42]}`)
		RunTestThatExpression(t, `test_set_all_repeated_sint32(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_sint32('{"7": [-1, 42]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_sint32('{}')`).IsEqualToInt(0)              // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_sint32('{"7": []}')`).IsEqualToInt(0)       // Empty array
		RunTestThatExpression(t, `test_count_repeated_sint32('{"7": [-1, 42]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_sint32('{"7": [-1, 42, 100]}', 0)`).IsEqualToInt(-1)
		RunTestThatExpression(t, `test_get_repeated_sint32('{"7": [-1]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_sint32('{"7": [-1, 42]}', 1, 100)`).IsEqualToJsonString(`{"7": [-1, 100]}`)
		RunTestThatExpression(t, `test_insert_repeated_sint32('{"7": [42]}', 0, -1)`).IsEqualToJsonString(`{"7": [-1, 42]}`)
		RunTestThatExpression(t, `test_remove_repeated_sint32('{"7": [-1, 42, 100]}', 1)`).IsEqualToJsonString(`{"7": [-1, 100]}`)
		RunTestThatExpression(t, `test_add_all_repeated_sint32('{"7": [-1]}', '[42, 100]')`).IsEqualToJsonString(`{"7": [-1, 42, 100]}`)
	})

	// Test repeated sint64 field
	t.Run("repeated_sint64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_sint64(?, -1)", `{}`).IsEqualToJsonString(`{"8": [-1]}`)
		RunTestThatExpression(t, "test_add_repeated_sint64(?, 100)", `{"8": [-1]}`).IsEqualToJsonString(`{"8": [-1, 100]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_sint64('{"8": [-1, 100]}')`).IsEqualToJsonString(`[-1, 100]`)
		RunTestThatExpression(t, `test_get_all_repeated_sint64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_sint64(?, '[-1, 100]')`, `{}`).IsEqualToJsonString(`{"8": [-1, 100]}`)
		RunTestThatExpression(t, `test_set_all_repeated_sint64(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_sint64('{"8": [-1, 100]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_sint64('{}')`).IsEqualToInt(0)               // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_sint64('{"8": []}')`).IsEqualToInt(0)        // Empty array
		RunTestThatExpression(t, `test_count_repeated_sint64('{"8": [-1, 100]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_sint64('{"8": [-1, 100, 200]}', 1)`).IsEqualToInt(100)
		RunTestThatExpression(t, `test_get_repeated_sint64('{"8": [-1]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_sint64('{"8": [-1, 100]}', 0, -999)`).IsEqualToJsonString(`{"8": [-999, 100]}`)
		RunTestThatExpression(t, `test_insert_repeated_sint64('{"8": [100]}', 0, -1)`).IsEqualToJsonString(`{"8": [-1, 100]}`)
		RunTestThatExpression(t, `test_remove_repeated_sint64('{"8": [-1, 100, 200]}', 1)`).IsEqualToJsonString(`{"8": [-1, 200]}`)
		RunTestThatExpression(t, `test_add_all_repeated_sint64('{"8": [-1]}', '[100, 200]')`).IsEqualToJsonString(`{"8": [-1, 100, 200]}`)
	})

	// Test repeated fixed32 field
	t.Run("repeated_fixed32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_fixed32(?, 4294967295)", `{}`).IsEqualToJsonString(`{"9": [4294967295]}`)
		RunTestThatExpression(t, "test_add_repeated_fixed32(?, 42)", `{"9": [4294967295]}`).IsEqualToJsonString(`{"9": [4294967295, 42]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_fixed32('{"9": [4294967295, 42]}')`).IsEqualToJsonString(`[4294967295, 42]`)
		RunTestThatExpression(t, `test_get_all_repeated_fixed32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_fixed32(?, '[4294967295, 42]')`, `{}`).IsEqualToJsonString(`{"9": [4294967295, 42]}`)
		RunTestThatExpression(t, `test_set_all_repeated_fixed32(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_fixed32('{"9": [4294967295, 42]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_fixed32('{}')`).IsEqualToInt(0)                      // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_fixed32('{"9": []}')`).IsEqualToInt(0)               // Empty array
		RunTestThatExpression(t, `test_count_repeated_fixed32('{"9": [4294967295, 42]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_fixed32('{"9": [4294967295, 42]}', 0)`).IsEqualToInt(4294967295)
		RunTestThatExpression(t, `test_get_repeated_fixed32('{"9": [100]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_fixed32('{"9": [100, 200]}', 1, 999)`).IsEqualToJsonString(`{"9": [100, 999]}`)
		RunTestThatExpression(t, `test_insert_repeated_fixed32('{"9": [200]}', 0, 100)`).IsEqualToJsonString(`{"9": [100, 200]}`)
		RunTestThatExpression(t, `test_remove_repeated_fixed32('{"9": [100, 200, 300]}', 1)`).IsEqualToJsonString(`{"9": [100, 300]}`)
		RunTestThatExpression(t, `test_add_all_repeated_fixed32('{"9": [100]}', '[200, 300]')`).IsEqualToJsonString(`{"9": [100, 200, 300]}`)
	})

	// Test repeated fixed64 field
	t.Run("repeated_fixed64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_fixed64(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"10": [18446744073709551615]}`)
		RunTestThatExpression(t, "test_add_repeated_fixed64(?, 100)", `{"10": [18446744073709551615]}`).IsEqualToJsonString(`{"10": [18446744073709551615, 100]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_fixed64('{"10": [18446744073709551615, 100]}')`).IsEqualToJsonString(`[18446744073709551615, 100]`)
		RunTestThatExpression(t, `test_get_all_repeated_fixed64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_fixed64(?, '[18446744073709551615, 100]')`, `{}`).IsEqualToJsonString(`{"10": [18446744073709551615, 100]}`)
		RunTestThatExpression(t, `test_set_all_repeated_fixed64(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_fixed64('{"10": [18446744073709551615, 100]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_fixed64('{}')`).IsEqualToInt(0)                                  // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_fixed64('{"10": []}')`).IsEqualToInt(0)                          // Empty array
		RunTestThatExpression(t, `test_count_repeated_fixed64('{"10": [18446744073709551615, 100]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_fixed64('{"10": [18446744073709551615, 100]}', 1)`).IsEqualToInt(100)
		RunTestThatExpression(t, `test_get_repeated_fixed64('{"10": [100]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_fixed64('{"10": [100, 200]}', 0, 999)`).IsEqualToJsonString(`{"10": [999, 200]}`)
		RunTestThatExpression(t, `test_insert_repeated_fixed64('{"10": [200]}', 0, 100)`).IsEqualToJsonString(`{"10": [100, 200]}`)
		RunTestThatExpression(t, `test_remove_repeated_fixed64('{"10": [100, 200, 300]}', 1)`).IsEqualToJsonString(`{"10": [100, 300]}`)
		RunTestThatExpression(t, `test_add_all_repeated_fixed64('{"10": [100]}', '[200, 300]')`).IsEqualToJsonString(`{"10": [100, 200, 300]}`)
	})

	// Test repeated sfixed32 field
	t.Run("repeated_sfixed32_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_sfixed32(?, -2147483648)", `{}`).IsEqualToJsonString(`{"11": [-2147483648]}`)
		RunTestThatExpression(t, "test_add_repeated_sfixed32(?, 42)", `{"11": [-2147483648]}`).IsEqualToJsonString(`{"11": [-2147483648, 42]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_sfixed32('{"11": [-2147483648, 42]}')`).IsEqualToJsonString(`[-2147483648, 42]`)
		RunTestThatExpression(t, `test_get_all_repeated_sfixed32('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_sfixed32(?, '[-2147483648, 42]')`, `{}`).IsEqualToJsonString(`{"11": [-2147483648, 42]}`)
		RunTestThatExpression(t, `test_set_all_repeated_sfixed32(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_sfixed32('{"11": [-2147483648, 42]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_sfixed32('{}')`).IsEqualToInt(0)                        // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_sfixed32('{"11": []}')`).IsEqualToInt(0)                // Empty array
		RunTestThatExpression(t, `test_count_repeated_sfixed32('{"11": [-2147483648, 42]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_sfixed32('{"11": [-2147483648, 42]}', 0)`).IsEqualToInt(-2147483648)
		RunTestThatExpression(t, `test_get_repeated_sfixed32('{"11": [42]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_sfixed32('{"11": [-100, 42]}', 1, 999)`).IsEqualToJsonString(`{"11": [-100, 999]}`)
		RunTestThatExpression(t, `test_insert_repeated_sfixed32('{"11": [42]}', 0, -100)`).IsEqualToJsonString(`{"11": [-100, 42]}`)
		RunTestThatExpression(t, `test_remove_repeated_sfixed32('{"11": [-100, 42, 100]}', 1)`).IsEqualToJsonString(`{"11": [-100, 100]}`)
		RunTestThatExpression(t, `test_add_all_repeated_sfixed32('{"11": [-100]}', '[42, 100]')`).IsEqualToJsonString(`{"11": [-100, 42, 100]}`)
	})

	// Test repeated sfixed64 field
	t.Run("repeated_sfixed64_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_sfixed64(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"12": [-9223372036854775808]}`)
		RunTestThatExpression(t, "test_add_repeated_sfixed64(?, 100)", `{"12": [-9223372036854775808]}`).IsEqualToJsonString(`{"12": [-9223372036854775808, 100]}`)

		// Test get operations return actual numeric arrays
		RunTestThatExpression(t, `test_get_all_repeated_sfixed64('{"12": [-9223372036854775808, 100]}')`).IsEqualToJsonString(`[-9223372036854775808, 100]`)
		RunTestThatExpression(t, `test_get_all_repeated_sfixed64('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_sfixed64(?, '[-9223372036854775808, 100]')`, `{}`).IsEqualToJsonString(`{"12": [-9223372036854775808, 100]}`)
		RunTestThatExpression(t, `test_set_all_repeated_sfixed64(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_sfixed64('{"12": [-9223372036854775808, 100]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_sfixed64('{}')`).IsEqualToInt(0)                                  // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_sfixed64('{"12": []}')`).IsEqualToInt(0)                          // Empty array
		RunTestThatExpression(t, `test_count_repeated_sfixed64('{"12": [-9223372036854775808, 100]}')`).IsEqualToInt(2) // Two elements

		// Test index-based operations
		RunTestThatExpression(t, `test_get_repeated_sfixed64('{"12": [-9223372036854775808, 100]}', 1)`).IsEqualToInt(100)
		RunTestThatExpression(t, `test_get_repeated_sfixed64('{"12": [100]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_set_repeated_sfixed64('{"12": [-100, 100]}', 0, -999)`).IsEqualToJsonString(`{"12": [-999, 100]}`)
		RunTestThatExpression(t, `test_insert_repeated_sfixed64('{"12": [100]}', 0, -100)`).IsEqualToJsonString(`{"12": [-100, 100]}`)
		RunTestThatExpression(t, `test_remove_repeated_sfixed64('{"12": [-100, 100, 200]}', 1)`).IsEqualToJsonString(`{"12": [-100, 200]}`)
		RunTestThatExpression(t, `test_add_all_repeated_sfixed64('{"12": [-100]}', '[100, 200]')`).IsEqualToJsonString(`{"12": [-100, 100, 200]}`)
	})

	// Test repeated bool field (JSON booleans, not 1/0)
	t.Run("repeated_bool_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_bool(?, TRUE)", `{}`).IsEqualToJsonString(`{"13": [true]}`)
		RunTestThatExpression(t, "test_add_repeated_bool(?, FALSE)", `{"13": [true]}`).IsEqualToJsonString(`{"13": [true, false]}`)

		// Test get operations return actual boolean arrays
		RunTestThatExpression(t, `test_get_all_repeated_bool('{"13": [true, false]}')`).IsEqualToJsonString(`[true, false]`)
		RunTestThatExpression(t, `test_get_all_repeated_bool('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_bool(?, '[true, false]')`, `{}`).IsEqualToJsonString(`{"13": [true, false]}`)
		RunTestThatExpression(t, `test_set_all_repeated_bool(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_bool('{"13": [true, false]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_bool('{}')`).IsEqualToInt(0)                    // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_bool('{"13": []}')`).IsEqualToInt(0)            // Empty array
		RunTestThatExpression(t, `test_count_repeated_bool('{"13": [true, false]}')`).IsEqualToInt(2) // Two elements

		// Test index-based get operations
		RunTestThatExpression(t, `test_get_repeated_bool('{"13": [true, false, true]}', 0)`).IsEqualToBool(true)
		RunTestThatExpression(t, `test_get_repeated_bool('{"13": [true, false, true]}', 1)`).IsEqualToBool(false)
		RunTestThatExpression(t, `test_get_repeated_bool('{"13": [true, false, true]}', 2)`).IsEqualToBool(true)
		RunTestThatExpression(t, `test_get_repeated_bool('{"13": [true, false]}', 2)`).ToFailWithSignalException("45000", "Array index out of bounds")

		// Test index-based set operations
		RunTestThatExpression(t, `test_set_repeated_bool('{"13": [true, false]}', 0, false)`).IsEqualToJsonString(`{"13": [false, false]}`)
		RunTestThatExpression(t, `test_set_repeated_bool('{"13": [true, false]}', 1, true)`).IsEqualToJsonString(`{"13": [true, true]}`)

		// Test insert operations
		RunTestThatExpression(t, `test_insert_repeated_bool('{"13": [false]}', 0, true)`).IsEqualToJsonString(`{"13": [true, false]}`)
		RunTestThatExpression(t, `test_insert_repeated_bool('{}', 0, true)`).IsEqualToJsonString(`{"13": [true]}`)

		// Test remove operations
		RunTestThatExpression(t, `test_remove_repeated_bool('{"13": [true, false, true]}', 1)`).IsEqualToJsonString(`{"13": [true, true]}`)

		// Test add_all operations
		RunTestThatExpression(t, `test_add_all_repeated_bool('{"13": [true]}', '[false, true]')`).IsEqualToJsonString(`{"13": [true, false, true]}`)
	})

	// Test repeated string field
	t.Run("repeated_string_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_string(?, 'hello')", `{}`).IsEqualToJsonString(`{"14": ["hello"]}`)
		RunTestThatExpression(t, "test_add_repeated_string(?, 'world')", `{"14": ["hello"]}`).IsEqualToJsonString(`{"14": ["hello", "world"]}`)

		// Test get operations return actual string arrays
		RunTestThatExpression(t, `test_get_all_repeated_string('{"14": ["hello", "world"]}')`).IsEqualToJsonString(`["hello", "world"]`)
		RunTestThatExpression(t, `test_get_all_repeated_string('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_string(?, '["hello", "world"]')`, `{}`).IsEqualToJsonString(`{"14": ["hello", "world"]}`)
		RunTestThatExpression(t, `test_set_all_repeated_string(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_string('{"14": ["hello", "world"]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_string('{}')`).IsEqualToInt(0)                         // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_string('{"14": []}')`).IsEqualToInt(0)                 // Empty array
		RunTestThatExpression(t, `test_count_repeated_string('{"14": ["hello", "world"]}')`).IsEqualToInt(2) // Two elements

		// Test index-based get operations
		RunTestThatExpression(t, `test_get_repeated_string('{"14": ["hello", "world", "test"]}', 0)`).IsEqualToString("hello")
		RunTestThatExpression(t, `test_get_repeated_string('{"14": ["hello", "world", "test"]}', 1)`).IsEqualToString("world")
		RunTestThatExpression(t, `test_get_repeated_string('{"14": ["hello", "world", "test"]}', 2)`).IsEqualToString("test")
		RunTestThatExpression(t, `test_get_repeated_string('{"14": ["hello", "world"]}', 2)`).ToFailWithSignalException("45000", "Array index out of bounds")
		RunTestThatExpression(t, `test_get_repeated_string('{}', 0)`).ToFailWithSignalException("45000", "Array index out of bounds")

		// Test index-based set operations
		RunTestThatExpression(t, `test_set_repeated_string('{"14": ["hello", "world"]}', 0, "hi")`).IsEqualToJsonString(`{"14": ["hi", "world"]}`)
		RunTestThatExpression(t, `test_set_repeated_string('{"14": ["hello", "world"]}', 1, "universe")`).IsEqualToJsonString(`{"14": ["hello", "universe"]}`)
		RunTestThatExpression(t, `test_set_repeated_string('{"14": ["hello", "world"]}', 2, "test")`).ToFailWithSignalException("45000", "Array index out of bounds")

		// Test insert operations
		RunTestThatExpression(t, `test_insert_repeated_string('{"14": ["world"]}', 0, "hello")`).IsEqualToJsonString(`{"14": ["hello", "world"]}`)
		RunTestThatExpression(t, `test_insert_repeated_string('{"14": ["hello", "world"]}', 1, "beautiful")`).IsEqualToJsonString(`{"14": ["hello", "beautiful", "world"]}`)
		RunTestThatExpression(t, `test_insert_repeated_string('{}', 0, "first")`).IsEqualToJsonString(`{"14": ["first"]}`)

		// Test remove operations
		RunTestThatExpression(t, `test_remove_repeated_string('{"14": ["hello", "beautiful", "world"]}', 1)`).IsEqualToJsonString(`{"14": ["hello", "world"]}`)
		RunTestThatExpression(t, `test_remove_repeated_string('{"14": ["hello"]}', 0)`).IsEqualToJsonString(`{}`)

		// Test add_all operations
		RunTestThatExpression(t, `test_add_all_repeated_string('{"14": ["hello"]}', '["world", "test"]')`).IsEqualToJsonString(`{"14": ["hello", "world", "test"]}`)
		RunTestThatExpression(t, `test_add_all_repeated_string('{}', '["hello", "world"]')`).IsEqualToJsonString(`{"14": ["hello", "world"]}`)
	})

	// Test repeated bytes field (base64 encoded)
	t.Run("repeated_bytes_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_bytes(?, ?)", `{}`, []byte("hello")).IsEqualToJsonString(`{"15": ["aGVsbG8="]}`)
		RunTestThatExpression(t, "test_add_repeated_bytes(?, ?)", `{"15": ["aGVsbG8="]}`, []byte("world")).IsEqualToJsonString(`{"15": ["aGVsbG8=", "d29ybGQ="]}`)

		// Test get operations convert from base64 back to actual byte arrays (but returned as base64 JSON)
		RunTestThatExpression(t, `test_get_all_repeated_bytes('{"15": ["aGVsbG8=", "d29ybGQ="]}')`).IsEqualToJsonString(`["aGVsbG8=", "d29ybGQ="]`) // Returns base64 strings in array
		RunTestThatExpression(t, `test_get_all_repeated_bytes('{}')`).IsEqualToJsonString(`[]`)                                                     // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_bytes(?, '["aGVsbG8=", "d29ybGQ="]')`, `{}`).IsEqualToJsonString(`{"15": ["aGVsbG8=", "d29ybGQ="]}`)
		RunTestThatExpression(t, `test_set_all_repeated_bytes(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_bytes('{"15": ["aGVsbG8=", "d29ybGQ="]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_bytes('{}')`).IsEqualToInt(0)                               // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_bytes('{"15": []}')`).IsEqualToInt(0)                       // Empty array
		RunTestThatExpression(t, `test_count_repeated_bytes('{"15": ["aGVsbG8=", "d29ybGQ="]}')`).IsEqualToInt(2) // Two elements

		// Test index-based get operations
		RunTestThatExpression(t, `test_get_repeated_bytes('{"15": ["aGVsbG8=", "d29ybGQ=", "dGVzdA=="]}', 0)`).IsEqualToBytes([]byte("hello"))
		RunTestThatExpression(t, `test_get_repeated_bytes('{"15": ["aGVsbG8=", "d29ybGQ=", "dGVzdA=="]}', 1)`).IsEqualToBytes([]byte("world"))
		RunTestThatExpression(t, `test_get_repeated_bytes('{"15": ["aGVsbG8=", "d29ybGQ="]}', 2)`).ToFailWithSignalException("45000", "Array index out of bounds")

		// Test index-based set operations
		RunTestThatExpression(t, "test_set_repeated_bytes(?, 0, ?)", `{"15": ["aGVsbG8=", "d29ybGQ="]}`, []byte("hi")).IsEqualToJsonString(`{"15": ["aGk=", "d29ybGQ="]}`)
		RunTestThatExpression(t, "test_set_repeated_bytes(?, 1, ?)", `{"15": ["aGVsbG8=", "d29ybGQ="]}`, []byte("universe")).IsEqualToJsonString(`{"15": ["aGVsbG8=", "dW5pdmVyc2U="]}`)

		// Test insert operations
		RunTestThatExpression(t, "test_insert_repeated_bytes(?, 0, ?)", `{"15": ["d29ybGQ="]}`, []byte("hello")).IsEqualToJsonString(`{"15": ["aGVsbG8=", "d29ybGQ="]}`)
		RunTestThatExpression(t, "test_insert_repeated_bytes(?, 0, ?)", `{}`, []byte("first")).IsEqualToJsonString(`{"15": ["Zmlyc3Q="]}`)

		// Test remove operations
		RunTestThatExpression(t, `test_remove_repeated_bytes('{"15": ["aGVsbG8=", "YmVhdXRpZnVs", "d29ybGQ="]}', 1)`).IsEqualToJsonString(`{"15": ["aGVsbG8=", "d29ybGQ="]}`)

		// Test add_all operations
		RunTestThatExpression(t, `test_add_all_repeated_bytes('{"15": ["aGVsbG8="]}', '["d29ybGQ=", "dGVzdA=="]')`).IsEqualToJsonString(`{"15": ["aGVsbG8=", "d29ybGQ=", "dGVzdA=="]}`)
	})

	// Test repeated enum field (numbers, not string names)
	t.Run("repeated_enum_field", func(t *testing.T) {
		// Test add operations create correct internal format
		RunTestThatExpression(t, "test_add_repeated_enum(?, 1)", `{}`).IsEqualToJsonString(`{"16": [1]}`)
		RunTestThatExpression(t, "test_add_repeated_enum(?, 2)", `{"16": [1]}`).IsEqualToJsonString(`{"16": [1, 2]}`)

		// Test get operations return actual integer arrays
		RunTestThatExpression(t, `test_get_all_repeated_enum('{"16": [1, 2]}')`).IsEqualToJsonString(`[1, 2]`)
		RunTestThatExpression(t, `test_get_all_repeated_enum('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_enum(?, '[1, 2]')`, `{}`).IsEqualToJsonString(`{"16": [1, 2]}`)
		RunTestThatExpression(t, `test_set_all_repeated_enum(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_enum('{"16": [1, 2]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_enum('{}')`).IsEqualToInt(0)             // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_enum('{"16": []}')`).IsEqualToInt(0)     // Empty array
		RunTestThatExpression(t, `test_count_repeated_enum('{"16": [1, 2]}')`).IsEqualToInt(2) // Two elements

		// Test index-based get operations
		RunTestThatExpression(t, `test_get_repeated_enum('{"16": [1, 2, 0]}', 0)`).IsEqualToInt(1)
		RunTestThatExpression(t, `test_get_repeated_enum('{"16": [1, 2, 0]}', 1)`).IsEqualToInt(2)
		RunTestThatExpression(t, `test_get_repeated_enum('{"16": [1, 2, 0]}', 2)`).IsEqualToInt(0)
		RunTestThatExpression(t, `test_get_repeated_enum('{"16": [1, 2]}', 2)`).ToFailWithSignalException("45000", "Array index out of bounds")

		// Test index-based set operations
		RunTestThatExpression(t, `test_set_repeated_enum('{"16": [1, 2]}', 0, 0)`).IsEqualToJsonString(`{"16": [0, 2]}`)
		RunTestThatExpression(t, `test_set_repeated_enum('{"16": [1, 2]}', 1, 1)`).IsEqualToJsonString(`{"16": [1, 1]}`)

		// Test insert operations
		RunTestThatExpression(t, `test_insert_repeated_enum('{"16": [2]}', 0, 1)`).IsEqualToJsonString(`{"16": [1, 2]}`)
		RunTestThatExpression(t, `test_insert_repeated_enum('{}', 0, 1)`).IsEqualToJsonString(`{"16": [1]}`)

		// Test remove operations
		RunTestThatExpression(t, `test_remove_repeated_enum('{"16": [1, 2, 0]}', 1)`).IsEqualToJsonString(`{"16": [1, 0]}`)

		// Test add_all operations
		RunTestThatExpression(t, `test_add_all_repeated_enum('{"16": [1]}', '[2, 0]')`).IsEqualToJsonString(`{"16": [1, 2, 0]}`)
	})

	// Test repeated message field (array of nested objects with field number keys)
	t.Run("repeated_message_field", func(t *testing.T) {
		// Test add operations create correct internal format
		nested1 := "nested_set_name(nested_new(), 'first')"
		nested2 := "nested_set_value(nested_set_name(nested_new(), 'second'), 42)"
		RunTestThatExpression(t, fmt.Sprintf("test_add_repeated_message(?, %s)", nested1), `{}`).IsEqualToJsonString(`{"17": [{"1": "first"}]}`)
		RunTestThatExpression(t, fmt.Sprintf("test_add_repeated_message(?, %s)", nested2), `{"17": [{"1": "first"}]}`).IsEqualToJsonString(`{"17": [{"1": "first"}, {"1": "second", "2": 42}]}`)

		// Test get operations return actual nested object arrays
		RunTestThatExpression(t, `test_get_all_repeated_message('{"17": [{"1": "first"}, {"1": "second", "2": 42}]}')`).IsEqualToJsonString(`[{"1": "first"}, {"1": "second", "2": 42}]`)
		RunTestThatExpression(t, `test_get_all_repeated_message('{}')`).IsEqualToJsonString(`[]`) // Missing field returns empty array

		// Test set operations create correct internal format
		RunTestThatExpression(t, `test_set_all_repeated_message(?, '[{"1": "first"}, {"1": "second", "2": 42}]')`, `{}`).IsEqualToJsonString(`{"17": [{"1": "first"}, {"1": "second", "2": 42}]}`)
		RunTestThatExpression(t, `test_set_all_repeated_message(?, '[]')`, `{}`).IsEqualToJsonString(`{}`) // Empty array omitted

		// Test clear operations remove field and return empty JSON
		RunTestThatExpression(t, `test_clear_repeated_message('{"17": [{"1": "first"}, {"1": "second", "2": 42}]}')`).IsEqualToJsonString(`{}`)

		// Test count operations
		RunTestThatExpression(t, `test_count_repeated_message('{}')`).IsEqualToInt(0)                                                 // Empty object/missing field
		RunTestThatExpression(t, `test_count_repeated_message('{"17": []}')`).IsEqualToInt(0)                                         // Empty array
		RunTestThatExpression(t, `test_count_repeated_message('{"17": [{"1": "first"}, {"1": "second", "2": 42}]}')`).IsEqualToInt(2) // Two elements

		// Test index-based get operations
		RunTestThatExpression(t, `test_get_repeated_message('{"17": [{"1": "first"}, {"1": "second", "2": 42}]}', 0)`).IsEqualToJsonString(`{"1": "first"}`)
		RunTestThatExpression(t, `test_get_repeated_message('{"17": [{"1": "first"}, {"1": "second", "2": 42}]}', 1)`).IsEqualToJsonString(`{"1": "second", "2": 42}`)
		RunTestThatExpression(t, `test_get_repeated_message('{"17": [{"1": "first"}]}', 1)`).ToFailWithSignalException("45000", "Array index out of bounds")

		// Test index-based set operations
		nestedSetUpdate := "nested_set_name(nested_new(), 'updated')"
		RunTestThatExpression(t, fmt.Sprintf(`test_set_repeated_message('{"17": [{"1": "first"}, {"1": "second"}]}', 0, %s)`, nestedSetUpdate)).IsEqualToJsonString(`{"17": [{"1": "updated"}, {"1": "second"}]}`)
		nestedSetModified := "nested_set_value(nested_set_name(nested_new(), 'modified'), 99)"
		RunTestThatExpression(t, fmt.Sprintf(`test_set_repeated_message('{"17": [{"1": "first"}, {"1": "second"}]}', 1, %s)`, nestedSetModified)).IsEqualToJsonString(`{"17": [{"1": "first"}, {"1": "modified", "2": 99}]}`)

		// Test insert operations
		nestedInsert := "nested_set_name(nested_new(), 'inserted')"
		RunTestThatExpression(t, fmt.Sprintf(`test_insert_repeated_message('{"17": [{"1": "second"}]}', 0, %s)`, nestedInsert)).IsEqualToJsonString(`{"17": [{"1": "inserted"}, {"1": "second"}]}`)
		RunTestThatExpression(t, fmt.Sprintf(`test_insert_repeated_message('{}', 0, %s)`, nestedInsert)).IsEqualToJsonString(`{"17": [{"1": "inserted"}]}`)

		// Test remove operations
		RunTestThatExpression(t, `test_remove_repeated_message('{"17": [{"1": "first"}, {"1": "second"}, {"1": "third"}]}', 1)`).IsEqualToJsonString(`{"17": [{"1": "first"}, {"1": "third"}]}`)

		// Test add_all operations
		RunTestThatExpression(t, `test_add_all_repeated_message('{"17": [{"1": "first"}]}', '[{"1": "second"}, {"1": "third", "2": 42}]')`).IsEqualToJsonString(`{"17": [{"1": "first"}, {"1": "second"}, {"1": "third", "2": 42}]}`)
	})
}
