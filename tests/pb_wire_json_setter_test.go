package main

import "testing"

func TestWireJsonSetterFunctions(t *testing.T) {
	// Test setting int32 field
	RunTestThatExpression(t, "pb_wire_json_set_int32_field('{}', 1, 42)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 42}]}`)

	// Test setting string field
	RunTestThatExpression(t, "pb_wire_json_set_string_field('{}', 2, 'hello')").
		IsEqualToJsonString(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "aGVsbG8="}]}`)

	// Test setting bool field
	RunTestThatExpression(t, "pb_wire_json_set_bool_field('{}', 3, TRUE)").
		IsEqualToJsonString(`{"3": [{"i": 0, "n": 3, "t": 0, "v": 1}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_bool_field('{}', 3, FALSE)").
		IsEqualToJsonString(`{"3": [{"i": 0, "n": 3, "t": 0, "v": 0}]}`)

	// Test setting float field
	RunTestThatExpression(t, "pb_wire_json_set_float_field('{}', 4, 1.5)").
		IsEqualToJsonString(`{"4": [{"i": 0, "n": 4, "t": 5, "v": 1069547520}]}`)
}

func TestWireJsonAddRepeatedFunctions(t *testing.T) {
	// Test adding repeated VARINT types
	wire_json := "{}"
	RunTestThatExpression(t, "pb_wire_json_add_repeated_int32_field_element('"+wire_json+"', 1, 10, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}]}`)

	// Add second element
	wire_json = `{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}]}`
	RunTestThatExpression(t, "pb_wire_json_add_repeated_int32_field_element('"+wire_json+"', 1, 20, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}, {"i": 1, "n": 1, "t": 0, "v": 20}]}`)

	// Test all VARINT repeated types
	RunTestThatExpression(t, "pb_wire_json_add_repeated_int64_field_element('{}', 1, -1, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744073709551615}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_uint32_field_element('{}', 1, 4294967295, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 4294967295}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_uint64_field_element('{}', 1, 18446744073709551615, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744073709551615}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_sint32_field_element('{}', 1, -2, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 3}]}`) // ZigZag: -2 -> 3

	RunTestThatExpression(t, "pb_wire_json_add_repeated_sint64_field_element('{}', 1, -2, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 3}]}`) // ZigZag: -2 -> 3

	RunTestThatExpression(t, "pb_wire_json_add_repeated_enum_field_element('{}', 1, 100, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 100}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_bool_field_element('{}', 1, TRUE, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}]}`)

	// Test I32 repeated types
	RunTestThatExpression(t, "pb_wire_json_add_repeated_fixed32_field_element('{}', 1, 123, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 123}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_sfixed32_field_element('{}', 1, -123, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 4294967173}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_float_field_element('{}', 1, 2.5, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 1075838976}]}`)

	// Test I64 repeated types
	RunTestThatExpression(t, "pb_wire_json_add_repeated_fixed64_field_element('{}', 1, 456, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 456}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_sfixed64_field_element('{}', 1, -456, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 18446744073709551160}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_double_field_element('{}', 1, 2.5, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 4612811918334230528}]}`)

	// Test LEN repeated types
	RunTestThatExpression(t, "pb_wire_json_add_repeated_string_field_element('{}', 2, 'first')").
		IsEqualToJsonString(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "Zmlyc3Q="}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_bytes_field_element('{}', 2, _binary X'abcd')").
		IsEqualToJsonString(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "q80="}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_message_field_element('{}', 2, _binary X'1008')").
		IsEqualToJsonString(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "EAg="}]}`)
}

func TestWireJsonAddRepeatedFunctionsPacked(t *testing.T) {
	// Test packed encoding - should create LEN (wire type 2) elements
	RunTestThatExpression(t, "pb_wire_json_add_repeated_int32_field_element('{}', 1, 10, TRUE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "Cg=="}]}`) // 10 as varint in base64

	// Test appending to packed field
	wire_json := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "Cg=="}]}`
	RunTestThatExpression(t, "pb_wire_json_add_repeated_int32_field_element(?, 1, 20, TRUE)", wire_json).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "ChQ="}]}`) // 10,20 as varints in base64
}

func TestWireJsonSetRepeatedFunctions(t *testing.T) {
	// Test setting repeated field elements at specific indices

	// VARINT types (wire type 0) - int32
	wire_json_int32 := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}, {"i": 1, "n": 1, "t": 0, "v": 20}, {"i": 2, "n": 1, "t": 0, "v": 30}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_int32_field_element(?, 1, 1, 99)", wire_json_int32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}, {"i": 1, "n": 1, "t": 0, "v": 99}, {"i": 2, "n": 1, "t": 0, "v": 30}]}`)

	// VARINT types (wire type 0) - int64
	wire_json_int64 := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 1000}, {"i": 1, "n": 1, "t": 0, "v": 2000}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_int64_field_element(?, 1, 0, -1)", wire_json_int64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744073709551615}, {"i": 1, "n": 1, "t": 0, "v": 2000}]}`)

	// VARINT types (wire type 0) - uint32
	wire_json_uint32 := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 100}, {"i": 1, "n": 1, "t": 0, "v": 200}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_uint32_field_element(?, 1, 1, 4294967295)", wire_json_uint32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 100}, {"i": 1, "n": 1, "t": 0, "v": 4294967295}]}`)

	// VARINT types (wire type 0) - uint64
	wire_json_uint64 := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 500}, {"i": 1, "n": 1, "t": 0, "v": 600}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_uint64_field_element(?, 1, 0, 18446744073709551615)", wire_json_uint64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744073709551615}, {"i": 1, "n": 1, "t": 0, "v": 600}]}`)

	// VARINT types (wire type 0) - sint32 (ZigZag encoded)
	wire_json_sint32 := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}, {"i": 1, "n": 1, "t": 0, "v": 3}]}` // -1 -> 1, -2 -> 3
	RunTestThatExpression(t, "pb_wire_json_set_repeated_sint32_field_element(?, 1, 0, -5)", wire_json_sint32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 9}, {"i": 1, "n": 1, "t": 0, "v": 3}]}`) // -5 -> 9

	// VARINT types (wire type 0) - sint64 (ZigZag encoded)
	wire_json_sint64 := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}, {"i": 1, "n": 1, "t": 0, "v": 3}]}` // -1 -> 1, -2 -> 3
	RunTestThatExpression(t, "pb_wire_json_set_repeated_sint64_field_element(?, 1, 1, -10)", wire_json_sint64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}, {"i": 1, "n": 1, "t": 0, "v": 19}]}`) // -10 -> 19

	// VARINT types (wire type 0) - enum
	wire_json_enum := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}, {"i": 1, "n": 1, "t": 0, "v": 2}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_enum_field_element(?, 1, 0, 42)", wire_json_enum).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 42}, {"i": 1, "n": 1, "t": 0, "v": 2}]}`)

	// VARINT types (wire type 0) - bool
	wire_json_bool := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}, {"i": 1, "n": 1, "t": 0, "v": 0}]}` // TRUE, FALSE
	RunTestThatExpression(t, "pb_wire_json_set_repeated_bool_field_element(?, 1, 1, TRUE)", wire_json_bool).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}, {"i": 1, "n": 1, "t": 0, "v": 1}]}`)

	// I32 types (wire type 5) - fixed32
	wire_json_fixed32 := `{"1": [{"i": 0, "n": 1, "t": 5, "v": 123}, {"i": 1, "n": 1, "t": 5, "v": 456}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_fixed32_field_element(?, 1, 0, 999)", wire_json_fixed32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 999}, {"i": 1, "n": 1, "t": 5, "v": 456}]}`)

	// I32 types (wire type 5) - sfixed32
	wire_json_sfixed32 := `{"1": [{"i": 0, "n": 1, "t": 5, "v": 123}, {"i": 1, "n": 1, "t": 5, "v": 4294967173}]}` // 123, -123
	RunTestThatExpression(t, "pb_wire_json_set_repeated_sfixed32_field_element(?, 1, 1, -999)", wire_json_sfixed32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 123}, {"i": 1, "n": 1, "t": 5, "v": 4294966297}]}`) // -999 in 2's complement

	// I32 types (wire type 5) - float
	wire_json_float := `{"1": [{"i": 0, "n": 1, "t": 5, "v": 1075838976}, {"i": 1, "n": 1, "t": 5, "v": 1077936128}]}` // 2.5, 3.5
	RunTestThatExpression(t, "pb_wire_json_set_repeated_float_field_element(?, 1, 0, 1.5)", wire_json_float).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 1069547520}, {"i": 1, "n": 1, "t": 5, "v": 1077936128}]}`) // 1.5, 3.5

	// I64 types (wire type 1) - fixed64
	wire_json_fixed64 := `{"1": [{"i": 0, "n": 1, "t": 1, "v": 123}, {"i": 1, "n": 1, "t": 1, "v": 456}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_fixed64_field_element(?, 1, 1, 999)", wire_json_fixed64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 123}, {"i": 1, "n": 1, "t": 1, "v": 999}]}`)

	// I64 types (wire type 1) - sfixed64
	wire_json_sfixed64 := `{"1": [{"i": 0, "n": 1, "t": 1, "v": 123}, {"i": 1, "n": 1, "t": 1, "v": 18446744073709551160}]}` // 123, -456
	RunTestThatExpression(t, "pb_wire_json_set_repeated_sfixed64_field_element(?, 1, 0, -999)", wire_json_sfixed64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 18446744073709550617}, {"i": 1, "n": 1, "t": 1, "v": 18446744073709551160}]}`) // -999, -456

	// I64 types (wire type 1) - double
	wire_json_double := `{"1": [{"i": 0, "n": 1, "t": 1, "v": 4612811918334230528}, {"i": 1, "n": 1, "t": 1, "v": 4616189618054758400}]}` // 2.5, 3.5
	RunTestThatExpression(t, "pb_wire_json_set_repeated_double_field_element(?, 1, 1, 1.5)", wire_json_double).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 4612811918334230528}, {"i": 1, "n": 1, "t": 1, "v": 4609434218613702656}]}`) // 2.5, 1.5

	// LEN types (wire type 2) - string
	wire_json_string := `{"2": [{"i": 0, "n": 2, "t": 2, "v": "aGVsbG8="}, {"i": 1, "n": 2, "t": 2, "v": "d29ybGQ="}]}` // "hello", "world"
	RunTestThatExpression(t, "pb_wire_json_set_repeated_string_field_element(?, 2, 0, 'goodbye')", wire_json_string).
		IsEqualToJsonString(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "Z29vZGJ5ZQ=="}, {"i": 1, "n": 2, "t": 2, "v": "d29ybGQ="}]}`) // "goodbye", "world"

	// LEN types (wire type 2) - bytes
	wire_json_bytes := `{"2": [{"i": 0, "n": 2, "t": 2, "v": "q80="}, {"i": 1, "n": 2, "t": 2, "v": "3q2+7w=="}]}` // 0xabcd, 0xdeadbeef
	RunTestThatExpression(t, "pb_wire_json_set_repeated_bytes_field_element(?, 2, 1, _binary X'cafebabe')", wire_json_bytes).
		IsEqualToJsonString(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "q80="}, {"i": 1, "n": 2, "t": 2, "v": "yv66vg=="}]}`) // 0xabcd, 0xcafebabe

	// LEN types (wire type 2) - message
	wire_json_message := `{"2": [{"i": 0, "n": 2, "t": 2, "v": "EAg="}, {"i": 1, "n": 2, "t": 2, "v": "GAo="}]}` // field 2 = 8, field 3 = 10
	RunTestThatExpression(t, "pb_wire_json_set_repeated_message_field_element(?, 2, 0, _binary X'200c')", wire_json_message).
		IsEqualToJsonString(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "IAw="}, {"i": 1, "n": 2, "t": 2, "v": "GAo="}]}`) // field 4 = 12, field 3 = 10
}

func TestWireJsonRemoveRepeatedFunctions(t *testing.T) {
	// Test removing repeated field elements at specific indices

	// VARINT types - int32 (non-packed)
	wire_json_int32 := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}, {"i": 1, "n": 1, "t": 0, "v": 20}, {"i": 2, "n": 1, "t": 0, "v": 30}]}`
	// Remove middle element (index 1)
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_int32_field_element(?, 1, 1)", wire_json_int32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}, {"i": 2, "n": 1, "t": 0, "v": 30}]}`)
	// Remove first element (index 0)
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_int32_field_element(?, 1, 0)", wire_json_int32).
		IsEqualToJsonString(`{"1": [{"i": 1, "n": 1, "t": 0, "v": 20}, {"i": 2, "n": 1, "t": 0, "v": 30}]}`)
	// Remove last element (index 2)
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_int32_field_element(?, 1, 2)", wire_json_int32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}, {"i": 1, "n": 1, "t": 0, "v": 20}]}`)

	// Test removing all elements from a field with single element
	wire_json_single := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 42}]}`
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_int32_field_element(?, 1, 0)", wire_json_single).
		IsEqualToJsonString(`{}`) // Field should be completely removed

	// VARINT types - other types
	wire_json_int64 := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 1000}, {"i": 1, "n": 1, "t": 0, "v": 2000}]}`
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_int64_field_element(?, 1, 0)", wire_json_int64).
		IsEqualToJsonString(`{"1": [{"i": 1, "n": 1, "t": 0, "v": 2000}]}`)

	wire_json_uint32 := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 100}, {"i": 1, "n": 1, "t": 0, "v": 200}]}`
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_uint32_field_element(?, 1, 1)", wire_json_uint32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 100}]}`)

	// I32 types - fixed32
	wire_json_fixed32 := `{"1": [{"i": 0, "n": 1, "t": 5, "v": 123}, {"i": 1, "n": 1, "t": 5, "v": 456}]}`
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_fixed32_field_element(?, 1, 0)", wire_json_fixed32).
		IsEqualToJsonString(`{"1": [{"i": 1, "n": 1, "t": 5, "v": 456}]}`)

	// I64 types - fixed64
	wire_json_fixed64 := `{"1": [{"i": 0, "n": 1, "t": 1, "v": 123}, {"i": 1, "n": 1, "t": 1, "v": 456}]}`
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_fixed64_field_element(?, 1, 1)", wire_json_fixed64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 123}]}`)

	// LEN types - string
	wire_json_string := `{"2": [{"i": 0, "n": 2, "t": 2, "v": "aGVsbG8="}, {"i": 1, "n": 2, "t": 2, "v": "d29ybGQ="}]}` // "hello", "world"
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_string_field_element(?, 2, 0)", wire_json_string).
		IsEqualToJsonString(`{"2": [{"i": 1, "n": 2, "t": 2, "v": "d29ybGQ="}]}`) // "world"

	// LEN types - bytes
	wire_json_bytes := `{"2": [{"i": 0, "n": 2, "t": 2, "v": "q80="}, {"i": 1, "n": 2, "t": 2, "v": "3q2+7w=="}]}` // 0xabcd, 0xdeadbeef
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_bytes_field_element(?, 2, 1)", wire_json_bytes).
		IsEqualToJsonString(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "q80="}]}`) // 0xabcd
}

func TestWireJsonRemoveRepeatedFunctionsPacked(t *testing.T) {
	// Test removing elements from packed fields

	// VARINT types - int32 packed
	// Packed field with [10, 20, 30] as varints: 0x0A, 0x14, 0x1E = "ChQe" in base64
	packed_int32 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "ChQe"}]}`
	// Remove middle element (20): should result in [10, 30] = 0x0A, 0x1E = "Ch4=" in base64
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_int32_field_element(?, 1, 1)", packed_int32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "Ch4="}]}`)
	// Remove first element (10): should result in [20, 30] = 0x14, 0x1E = "FB4=" in base64
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_int32_field_element(?, 1, 0)", packed_int32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "FB4="}]}`)
	// Remove last element (30): should result in [10, 20] = 0x0A, 0x14 = "ChQ=" in base64
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_int32_field_element(?, 1, 2)", packed_int32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "ChQ="}]}`)

	// Test removing all elements from packed field
	packed_single := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "Cg=="}]}` // [10]
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_int32_field_element(?, 1, 0)", packed_single).
		IsEqualToJsonString(`{}`) // Field should be completely removed

	// VARINT types - uint64 (simple values)
	packed_uint64 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "ChQ="}]}` // [10, 20]
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_uint64_field_element(?, 1, 0)", packed_uint64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "FA=="}]}`) // [20]

	// VARINT types - sint32 (ZigZag encoded)
	packed_sint32 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQM="}]}` // [-1, -2] as [1, 3]
	RunTestThatExpression(t, "pb_wire_json_remove_repeated_sint32_field_element(?, 1, 1)", packed_sint32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQ=="}]}`) // [-1] as [1]
}

func TestWireJsonSetRepeatedFunctionsPacked(t *testing.T) {
	// Test setting individual elements in packed fields

	// VARINT types (wire type 0 -> packed as wire type 2) - int32
	// Packed field with [10, 20, 30] as varints: 0x0A (10), 0x14 (20), 0x1E (30) = "ChQe" in base64
	packed_int32 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "ChQe"}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_int32_field_element(?, 1, 0, 99)", packed_int32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "YxQe"}]}`) // 0x63 (99), 0x14 (20), 0x1E (30)
	RunTestThatExpression(t, "pb_wire_json_set_repeated_int32_field_element(?, 1, 1, 88)", packed_int32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "Clge"}]}`) // 0x0A (10), 0x58 (88), 0x1E (30)
	RunTestThatExpression(t, "pb_wire_json_set_repeated_int32_field_element(?, 1, 2, 77)", packed_int32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "ChRN"}]}`) // 0x0A (10), 0x14 (20), 0x4D (77)

	// VARINT types - int64 (testing with larger values)
	// Packed field with [256, 512] as varints: 0x8002 (256), 0x8004 (512) = "gAKABA==" in base64
	packed_int64 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "gAKABA=="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_int64_field_element(?, 1, 0, 1024)", packed_int64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "gAiABA=="}]}`) // 0x8008 (1024), 0x8004 (512)

	// VARINT types - uint32
	// Packed field with [100, 200] as varints: 0x64 (100), 0xC801 (200) = "ZMgB" in base64
	packed_uint32 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "ZMgB"}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_uint32_field_element(?, 1, 1, 300)", packed_uint32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "ZKwC"}]}`) // 0x64 (100), 0xAC02 (300)

	// VARINT types - uint64 (use simple values like int32)
	// Packed field with [10, 20] as varints: 0x0A, 0x14 = "ChQ=" in base64
	packed_uint64 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "ChQ="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_uint64_field_element(?, 1, 0, 99)", packed_uint64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "YxQ="}]}`) // 0x63 (99), 0x14 (20)

	// VARINT types - sint32 (ZigZag encoded)
	// Packed field with [-1, -2] as zigzag varints: 0x01 (1), 0x03 (3) = "AQM=" in base64
	packed_sint32 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQM="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_sint32_field_element(?, 1, 0, -5)", packed_sint32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "CQM="}]}`) // 0x09 (9 for -5), 0x03 (3 for -2)

	// VARINT types - sint64 (ZigZag encoded)
	// Packed field with [-3, -4] as zigzag varints: 0x05 (5), 0x07 (7) = "BQc=" in base64
	packed_sint64 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "BQc="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_sint64_field_element(?, 1, 1, -10)", packed_sint64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "BRM="}]}`) // 0x05 (5 for -3), 0x13 (19 for -10)

	// VARINT types - enum
	// Packed field with [1, 2] as varints: 0x01, 0x02 = "AQI=" in base64
	packed_enum := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQI="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_enum_field_element(?, 1, 0, 42)", packed_enum).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "KgI="}]}`) // 0x2A (42), 0x02 (2)

	// VARINT types - bool
	// Packed field with [TRUE, FALSE] as varints: 0x01, 0x00 = "AQA=" in base64
	packed_bool := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQA="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_bool_field_element(?, 1, 1, TRUE)", packed_bool).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQE="}]}`) // 0x01 (TRUE), 0x01 (TRUE)

	// I32 types (wire type 5 -> packed as wire type 2) - fixed32 (use simple values)
	// Packed field with [1, 2] as little-endian 32-bit: 0x01000000, 0x02000000 = "AQAAAAIAAAA=" in base64
	packed_fixed32 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQAAAAIAAAA="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_fixed32_field_element(?, 1, 0, 3)", packed_fixed32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "AwAAAAIAAAA="}]}`) // 3, 2

	// I32 types - sfixed32 (use simple values)
	// Packed field with [1, -1] as little-endian 32-bit: 0x01000000, 0xFFFFFFFF = "AQAAAP////8=" in base64
	packed_sfixed32 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQAAAP////8="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_sfixed32_field_element(?, 1, 1, -2)", packed_sfixed32).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQAAAP7///8="}]}`) // 1, -2

	// I32 types - float (let the IEEE 754 functions handle encoding, just test modification works)
	// Create packed field with actual function, then test setting works
	packed_float := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AACAPwAAAEA="}]}` // [1.0, 2.0]
	RunTestThatExpression(t, "pb_wire_json_set_repeated_float_field_element(?, 1, 0, 3.0)", packed_float).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "AABAQAAAAEA="}]}`) // [3.0, 2.0]

	// I64 types (wire type 1 -> packed as wire type 2) - fixed64 (use simple values)
	// Packed field with [1, 2] as little-endian 64-bit: 16 bytes total
	packed_fixed64 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQAAAAAAAAACAAAAAAAAAA=="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_fixed64_field_element(?, 1, 1, 3)", packed_fixed64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQAAAAAAAAADAAAAAAAAAA=="}]}`) // 1, 3

	// I64 types - sfixed64 (use simple values)
	// Packed field with [1, -1] as little-endian 64-bit signed
	packed_sfixed64 := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AQAAAAAAAAD//////////w=="}]}`
	RunTestThatExpression(t, "pb_wire_json_set_repeated_sfixed64_field_element(?, 1, 0, -2)", packed_sfixed64).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "/v///////////////////w=="}]}`) // -2, -1

	// I64 types - double (let IEEE 754 functions handle encoding)
	// Use simple values and trust the underlying encoding
	packed_double := `{"1": [{"i": 0, "n": 1, "t": 2, "v": "AAAAAAAA8D8AAAAAAAAAQA=="}]}` // [1.0, 2.0]
	RunTestThatExpression(t, "pb_wire_json_set_repeated_double_field_element(?, 1, 1, 3.0)", packed_double).
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "AAAAAAAA8D8AAAAAAAAIQA=="}]}`) // [1.0, 3.0]
}

func TestWireJsonSetterReplacement(t *testing.T) {
	// Test that setter replaces existing field
	existing_wire_json := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}]}`
	RunTestThatExpression(t, "pb_wire_json_set_int32_field('"+existing_wire_json+"', 1, 99)").
		IsEqualToJsonString(`{"1": [{"i": 1, "n": 1, "t": 0, "v": 99}]}`)
}

func TestWireJsonSetterPreservesIndex(t *testing.T) {
	// Test that new fields get proper index values
	existing_wire_json := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}], "2": [{"i": 1, "n": 2, "t": 0, "v": 20}]}`
	RunTestThatExpression(t, "pb_wire_json_set_int32_field('"+existing_wire_json+"', 3, 30)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}], "2": [{"i": 1, "n": 2, "t": 0, "v": 20}], "3": [{"i": 2, "n": 3, "t": 0, "v": 30}]}`)
}

func TestWireJsonSetterRoundTrip(t *testing.T) {
	// Test round-trip: set field -> get field back
	wire_json := `pb_wire_json_set_int32_field('{}', 1, 42)`
	RunTestThatExpression(t, "pb_wire_json_get_int32_field("+wire_json+", 1, 0)").
		IsEqualToInt(42)

	// Test round-trip for all VARINT types
	RunTestThatExpression(t, "pb_wire_json_get_int64_field(pb_wire_json_set_int64_field('{}', 1, -9223372036854775808), 1, 0)").
		IsEqualToInt(-9223372036854775808)

	RunTestThatExpression(t, "pb_wire_json_get_uint32_field(pb_wire_json_set_uint32_field('{}', 1, 4294967295), 1, 0)").
		IsEqualToUint(4294967295)

	RunTestThatExpression(t, "pb_wire_json_get_uint64_field(pb_wire_json_set_uint64_field('{}', 1, 18446744073709551615), 1, 0)").
		IsEqualToUint(18446744073709551615)

	RunTestThatExpression(t, "pb_wire_json_get_sint32_field(pb_wire_json_set_sint32_field('{}', 1, -1), 1, 0)").
		IsEqualToInt(-1)

	RunTestThatExpression(t, "pb_wire_json_get_sint64_field(pb_wire_json_set_sint64_field('{}', 1, -1), 1, 0)").
		IsEqualToInt(-1)

	RunTestThatExpression(t, "pb_wire_json_get_enum_field(pb_wire_json_set_enum_field('{}', 1, 42), 1, 0)").
		IsEqualToInt(42)

	RunTestThatExpression(t, "pb_wire_json_get_bool_field(pb_wire_json_set_bool_field('{}', 1, TRUE), 1, FALSE)").
		IsEqualToBool(true)

	// Test round-trip for I32 types
	RunTestThatExpression(t, "pb_wire_json_get_fixed32_field(pb_wire_json_set_fixed32_field('{}', 1, 42), 1, 0)").
		IsEqualToUint(42)

	RunTestThatExpression(t, "pb_wire_json_get_sfixed32_field(pb_wire_json_set_sfixed32_field('{}', 1, -42), 1, 0)").
		IsEqualToInt(-42)

	RunTestThatExpression(t, "pb_wire_json_get_float_field(pb_wire_json_set_float_field('{}', 1, 1.5), 1, 0)").
		IsEqualToFloat(1.5)

	// Test round-trip for I64 types
	RunTestThatExpression(t, "pb_wire_json_get_fixed64_field(pb_wire_json_set_fixed64_field('{}', 1, 42), 1, 0)").
		IsEqualToUint(42)

	RunTestThatExpression(t, "pb_wire_json_get_sfixed64_field(pb_wire_json_set_sfixed64_field('{}', 1, -42), 1, 0)").
		IsEqualToInt(-42)

	RunTestThatExpression(t, "pb_wire_json_get_double_field(pb_wire_json_set_double_field('{}', 1, 1.5), 1, 0)").
		IsEqualToFloat(1.5)

	// Test round-trip for LEN types
	RunTestThatExpression(t, "pb_wire_json_get_string_field(pb_wire_json_set_string_field('{}', 1, 'hello'), 1, '')").
		IsEqualToString("hello")

	RunTestThatExpression(t, "HEX(pb_wire_json_get_bytes_field(pb_wire_json_set_bytes_field('{}', 1, _binary X'deadbeef'), 1, ''))").
		IsEqualToString("DEADBEEF")

	RunTestThatExpression(t, "HEX(pb_wire_json_get_message_field(pb_wire_json_set_message_field('{}', 1, _binary X'080a'), 1, ''))").
		IsEqualToString("080A")
}

func TestWireJsonSetterTypes(t *testing.T) {
	// Test all protobuf types

	// VARINT types (wire type 0)
	RunTestThatExpression(t, "pb_wire_json_set_int32_field('{}', 1, -2147483648)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744071562067968}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_int64_field('{}', 1, -9223372036854775808)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 9223372036854775808}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_uint32_field('{}', 1, 4294967295)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 4294967295}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_uint64_field('{}', 1, 18446744073709551615)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744073709551615}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_sint32_field('{}', 1, -1)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}]}`) // ZigZag encoded: -1 -> 1

	RunTestThatExpression(t, "pb_wire_json_set_sint64_field('{}', 1, -1)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}]}`) // ZigZag encoded: -1 -> 1

	RunTestThatExpression(t, "pb_wire_json_set_enum_field('{}', 1, 42)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 42}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_bool_field('{}', 1, TRUE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_bool_field('{}', 1, FALSE)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 0}]}`)

	// I32 types (wire type 5)
	RunTestThatExpression(t, "pb_wire_json_set_fixed32_field('{}', 1, 42)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 42}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_sfixed32_field('{}', 1, -42)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 4294967254}]}`) // 2's complement

	RunTestThatExpression(t, "pb_wire_json_set_float_field('{}', 1, 1.5)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 1069547520}]}`)

	// I64 types (wire type 1)
	RunTestThatExpression(t, "pb_wire_json_set_fixed64_field('{}', 1, 42)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 42}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_sfixed64_field('{}', 1, -42)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 18446744073709551574}]}`) // 2's complement

	RunTestThatExpression(t, "pb_wire_json_set_double_field('{}', 1, 1.5)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 4609434218613702656}]}`)

	// LEN types (wire type 2)
	RunTestThatExpression(t, "pb_wire_json_set_string_field('{}', 1, 'hello')").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "aGVsbG8="}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_bytes_field('{}', 1, _binary X'deadbeef')").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "3q2+7w=="}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_message_field('{}', 1, _binary X'080a')").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "CAo="}]}`)
}
