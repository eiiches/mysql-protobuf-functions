package main

import "testing"

func TestMessageToWireJson(t *testing.T) {
	// VARINT (0)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'10ffffffff07')").
		IsEqualToJson(`{"2":[{"i":0,"n":2,"t":0,"v":2147483647}]}`)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'20ffffffffffffffffff01')").
		IsEqualToJson(`{"4":[{"i":0,"n":4,"t":0,"v":18446744073709551615}]}`)

	// I64 (1)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'61ffffffffffffffff')").
		IsEqualToJson(`{"12":[{"i":0,"n":12,"t":1,"v":18446744073709551615}]}`)

	// I32 (5)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'7dffffffff')").
		IsEqualToJson(`{"15":[{"i":0,"n":15,"t":5,"v":4294967295}]}`)

	// LEN (2)
	RunTestThatExpression(t, "pb_message_to_wire_json(_binary X'4202080a')").
		IsEqualToJson(`{"8":[{"i":0,"n":8,"t":2,"v":"CAo="}]}`)
}
