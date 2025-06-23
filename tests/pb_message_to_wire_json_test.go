package main

import "testing"

func TestMessageToWireJson(t *testing.T) {
	// VARINT (0)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'10ffffffff07')").
		IsEqualToJson(`[{"value": {"uint": 2147483647}, "wire_type": 0, "field_number": 2}]`)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'20ffffffffffffffffff01')").
		IsEqualToJson(`[{"value": {"uint": 18446744073709551615}, "wire_type": 0, "field_number": 4}]`)

	// I64 (1)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'61ffffffffffffffff')").
		IsEqualToJson(`[{"value": {"uint": 18446744073709551615}, "wire_type": 1, "field_number": 12}]`)

	// I32 (5)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'7dffffffff')").
		IsEqualToJson(`[{"value": {"uint": 4294967295}, "wire_type": 5, "field_number": 15}]`)

	// LEN (2)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'4202080a')").
		IsEqualToJson(`[{"value": {"bytes": "CAo="}, "wire_type": 2, "field_number": 8}]`)
}
