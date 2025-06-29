package main

import "testing"

func TestWireJsonClearFunctions(t *testing.T) {
	// Test clearing int32 field
	wire_json := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 42}], "2": [{"i": 1, "n": 2, "t": 0, "v": 100}]}`
	RunTestThatExpression(t, "pb_wire_json_clear_int32_field('"+wire_json+"', 1)").
		IsEqualToJsonString(`{"2": [{"i": 1, "n": 2, "t": 0, "v": 100}]}`)

	// Test clearing non-existent field (should be no-op)
	RunTestThatExpression(t, "pb_wire_json_clear_int32_field('{}', 1)").
		IsEqualToJsonString(`{}`)

	// Test clearing string field
	wire_json = `{"1": [{"i": 0, "n": 1, "t": 0, "v": 42}], "2": [{"i": 1, "n": 2, "t": 2, "v": "aGVsbG8="}]}`
	RunTestThatExpression(t, "pb_wire_json_clear_string_field('"+wire_json+"', 2)").
		IsEqualToJsonString(`{"1": [{"i": 0, "n": 1, "t": 0, "v": 42}]}`)

	// Test clearing bool field
	wire_json = `{"3": [{"i": 0, "n": 3, "t": 0, "v": 1}]}`
	RunTestThatExpression(t, "pb_wire_json_clear_bool_field('"+wire_json+"', 3)").
		IsEqualToJsonString(`{}`)

	// Test clearing float field
	wire_json = `{"4": [{"i": 0, "n": 4, "t": 5, "v": 1069547520}]}`
	RunTestThatExpression(t, "pb_wire_json_clear_float_field('"+wire_json+"', 4)").
		IsEqualToJsonString(`{}`)
}

func TestWireJsonClearRepeatedFunctions(t *testing.T) {
	// Test clearing repeated int32 field
	wire_json := `{"1": [{"i": 0, "n": 1, "t": 0, "v": 10}, {"i": 1, "n": 1, "t": 0, "v": 20}], "2": [{"i": 2, "n": 2, "t": 0, "v": 100}]}`
	RunTestThatExpression(t, "pb_wire_json_clear_repeated_int32_field('"+wire_json+"', 1)").
		IsEqualToJsonString(`{"2": [{"i": 2, "n": 2, "t": 0, "v": 100}]}`)

	// Test clearing repeated string field
	wire_json = `{"2": [{"i": 0, "n": 2, "t": 2, "v": "Zmlyc3Q="}, {"i": 1, "n": 2, "t": 2, "v": "c2Vjb25k"}]}`
	RunTestThatExpression(t, "pb_wire_json_clear_repeated_string_field('"+wire_json+"', 2)").
		IsEqualToJsonString(`{}`)

	// Test clearing non-existent repeated field (should be no-op)
	RunTestThatExpression(t, "pb_wire_json_clear_repeated_int64_field('{}', 5)").
		IsEqualToJsonString(`{}`)

	// Test clearing repeated double field
	wire_json = `{"1": [{"i": 0, "n": 1, "t": 1, "v": 4612811918334230528}]}`
	RunTestThatExpression(t, "pb_wire_json_clear_repeated_double_field('"+wire_json+"', 1)").
		IsEqualToJsonString(`{}`)
}

func TestWireJsonClearAllTypes(t *testing.T) {
	// Test clearing each protobuf type

	// VARINT types (wire type 0)
	RunTestThatExpression(t, "pb_wire_json_clear_int32_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 0, \"v\": 42}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_int64_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 0, \"v\": 42}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_uint32_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 0, \"v\": 42}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_uint64_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 0, \"v\": 42}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_sint32_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 0, \"v\": 1}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_sint64_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 0, \"v\": 1}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_enum_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 0, \"v\": 100}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_bool_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 0, \"v\": 1}]}', 1)").
		IsEqualToJsonString(`{}`)

	// I32 types (wire type 5)
	RunTestThatExpression(t, "pb_wire_json_clear_fixed32_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 5, \"v\": 123}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_sfixed32_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 5, \"v\": 4294967173}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_float_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 5, \"v\": 1069547520}]}', 1)").
		IsEqualToJsonString(`{}`)

	// I64 types (wire type 1)
	RunTestThatExpression(t, "pb_wire_json_clear_fixed64_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 1, \"v\": 456}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_sfixed64_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 1, \"v\": 18446744073709551160}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_double_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 1, \"v\": 4612811918334230528}]}', 1)").
		IsEqualToJsonString(`{}`)

	// LEN types (wire type 2)
	RunTestThatExpression(t, "pb_wire_json_clear_string_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 2, \"v\": \"aGVsbG8=\"}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_bytes_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 2, \"v\": \"q80=\"}]}', 1)").
		IsEqualToJsonString(`{}`)
	RunTestThatExpression(t, "pb_wire_json_clear_message_field('{\"1\": [{\"i\": 0, \"n\": 1, \"t\": 2, \"v\": \"EAg=\"}]}', 1)").
		IsEqualToJsonString(`{}`)
}

func TestWireJsonClearRoundTrip(t *testing.T) {
	// Test set -> clear round-trip

	// Test int32 set -> clear
	set_result := "pb_wire_json_set_int32_field('{}', 1, 42)"
	RunTestThatExpression(t, "pb_wire_json_clear_int32_field("+set_result+", 1)").
		IsEqualToJsonString(`{}`)

	// Test string set -> clear
	set_result = "pb_wire_json_set_string_field('{}', 2, 'hello')"
	RunTestThatExpression(t, "pb_wire_json_clear_string_field("+set_result+", 2)").
		IsEqualToJsonString(`{}`)

	// Test float set -> clear
	set_result = "pb_wire_json_set_float_field('{}', 3, 1.5)"
	RunTestThatExpression(t, "pb_wire_json_clear_float_field("+set_result+", 3)").
		IsEqualToJsonString(`{}`)

	// Test double set -> clear
	set_result = "pb_wire_json_set_double_field('{}', 4, 2.5)"
	RunTestThatExpression(t, "pb_wire_json_clear_double_field("+set_result+", 4)").
		IsEqualToJsonString(`{}`)
}
