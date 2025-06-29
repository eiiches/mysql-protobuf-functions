package main

import "testing"

func TestWireJsonSetterFunctions(t *testing.T) {
	// Test setting int32 field
	RunTestThatExpression(t, "pb_wire_json_set_int32_field('{}', 1, 42)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 42}]}`)

	// Test setting string field
	RunTestThatExpression(t, "pb_wire_json_set_string_field('{}', 2, 'hello')").
		IsEqualToJson(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "aGVsbG8="}]}`)

	// Test setting bool field
	RunTestThatExpression(t, "pb_wire_json_set_bool_field('{}', 3, TRUE)").
		IsEqualToJson(`{"3": [{"i": 0, "n": 3, "t": 0, "v": 1}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_bool_field('{}', 3, FALSE)").
		IsEqualToJson(`{"3": [{"i": 0, "n": 3, "t": 0, "v": 0}]}`)

	// Test setting float field
	RunTestThatExpression(t, "pb_wire_json_set_float_field('{}', 4, 1.5)").
		IsEqualToJson(`{"4": [{"i": 0, "n": 4, "t": 5, "v": 1069547520}]}`)
}

func TestWireJsonAddRepeatedFunctions(t *testing.T) {
	// Test adding repeated VARINT types
	wire_json := "{}"
	RunTestThatExpression(t, "pb_wire_json_add_repeated_int32_field('"+wire_json+"', 1, 10)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}]}`)

	// Add second element
	wire_json = `{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}]}`
	RunTestThatExpression(t, "pb_wire_json_add_repeated_int32_field('"+wire_json+"', 1, 20)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}, {"i": 1, "n": 1, "t": 0, "v": 20}]}`)

	// Test all VARINT repeated types
	RunTestThatExpression(t, "pb_wire_json_add_repeated_int64_field('{}', 1, -1)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744073709551615}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_uint32_field('{}', 1, 4294967295)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 4294967295}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_uint64_field('{}', 1, 18446744073709551615)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744073709551615}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_sint32_field('{}', 1, -2)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 3}]}`)  // ZigZag: -2 -> 3

	RunTestThatExpression(t, "pb_wire_json_add_repeated_sint64_field('{}', 1, -2)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 3}]}`)  // ZigZag: -2 -> 3

	RunTestThatExpression(t, "pb_wire_json_add_repeated_enum_field('{}', 1, 100)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 100}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_bool_field('{}', 1, TRUE)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}]}`)

	// Test I32 repeated types
	RunTestThatExpression(t, "pb_wire_json_add_repeated_fixed32_field('{}', 1, 123)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 123}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_sfixed32_field('{}', 1, -123)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 4294967173}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_float_field('{}', 1, 2.5)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 1075838976}]}`)

	// Test I64 repeated types
	RunTestThatExpression(t, "pb_wire_json_add_repeated_fixed64_field('{}', 1, 456)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 456}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_sfixed64_field('{}', 1, -456)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 18446744073709551160}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_double_field('{}', 1, 2.5)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 4612811918334230528}]}`)

	// Test LEN repeated types
	RunTestThatExpression(t, "pb_wire_json_add_repeated_string_field('{}', 2, 'first')").
		IsEqualToJson(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "Zmlyc3Q="}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_bytes_field('{}', 2, _binary X'abcd')").
		IsEqualToJson(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "q80="}]}`)

	RunTestThatExpression(t, "pb_wire_json_add_repeated_message_field('{}', 2, _binary X'1008')").
		IsEqualToJson(`{"2": [{"i": 0, "n": 2, "t": 2, "v": "EAg="}]}`)
}

func TestWireJsonSetterReplacement(t *testing.T) {
	// Test that setter replaces existing field
	existing_wire_json := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}]}`
	RunTestThatExpression(t, "pb_wire_json_set_int32_field('"+existing_wire_json+"', 1, 99)").
		IsEqualToJson(`{"1": [{"i": 1, "n": 1, "t": 0, "v": 99}]}`)
}

func TestWireJsonSetterPreservesIndex(t *testing.T) {
	// Test that new fields get proper index values
	existing_wire_json := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}], "2": [{"i": 1, "n": 2, "t": 0, "v": 20}]}`
	RunTestThatExpression(t, "pb_wire_json_set_int32_field('"+existing_wire_json+"', 3, 30)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}], "2": [{"i": 1, "n": 2, "t": 0, "v": 20}], "3": [{"i": 2, "n": 3, "t": 0, "v": 30}]}`)
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
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744071562067968}]}`)
	
	RunTestThatExpression(t, "pb_wire_json_set_int64_field('{}', 1, -9223372036854775808)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 9223372036854775808}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_uint32_field('{}', 1, 4294967295)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 4294967295}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_uint64_field('{}', 1, 18446744073709551615)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 18446744073709551615}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_sint32_field('{}', 1, -1)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}]}`)  // ZigZag encoded: -1 -> 1

	RunTestThatExpression(t, "pb_wire_json_set_sint64_field('{}', 1, -1)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}]}`)  // ZigZag encoded: -1 -> 1

	RunTestThatExpression(t, "pb_wire_json_set_enum_field('{}', 1, 42)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 42}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_bool_field('{}', 1, TRUE)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 1}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_bool_field('{}', 1, FALSE)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 0}]}`)

	// I32 types (wire type 5)
	RunTestThatExpression(t, "pb_wire_json_set_fixed32_field('{}', 1, 42)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 42}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_sfixed32_field('{}', 1, -42)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 4294967254}]}`)  // 2's complement

	RunTestThatExpression(t, "pb_wire_json_set_float_field('{}', 1, 1.5)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 5, "v": 1069547520}]}`)

	// I64 types (wire type 1)
	RunTestThatExpression(t, "pb_wire_json_set_fixed64_field('{}', 1, 42)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 42}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_sfixed64_field('{}', 1, -42)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 18446744073709551574}]}`)  // 2's complement

	RunTestThatExpression(t, "pb_wire_json_set_double_field('{}', 1, 1.5)").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 1, "v": 4609434218613702656}]}`)

	// LEN types (wire type 2)
	RunTestThatExpression(t, "pb_wire_json_set_string_field('{}', 1, 'hello')").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "aGVsbG8="}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_bytes_field('{}', 1, _binary X'deadbeef')").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "3q2+7w=="}]}`)

	RunTestThatExpression(t, "pb_wire_json_set_message_field('{}', 1, _binary X'080a')").
		IsEqualToJson(`{"1": [{"i": 0, "n": 1, "t": 2, "v": "CAo="}]}`)
}