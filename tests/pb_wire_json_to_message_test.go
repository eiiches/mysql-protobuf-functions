package main

import "testing"

func TestWireJsonToMessage(t *testing.T) {
	// VARINT (0) - Test basic integer encoding
	RunTestThatExpression(t, "HEX(pb_wire_json_to_message(JSON_OBJECT('2', JSON_ARRAY(JSON_OBJECT('i', 0, 'n', 2, 't', 0, 'v', 2147483647)))))").
		IsEqualToString("10FFFFFFFF07")

	// I64 (1) - Test 64-bit integer encoding
	RunTestThatExpression(t, "HEX(pb_wire_json_to_message(JSON_OBJECT('12', JSON_ARRAY(JSON_OBJECT('i', 0, 'n', 12, 't', 1, 'v', 18446744073709551615)))))").
		IsEqualToString("61FFFFFFFFFFFFFFFF")

	// I32 (5) - Test 32-bit integer encoding
	RunTestThatExpression(t, "HEX(pb_wire_json_to_message(JSON_OBJECT('15', JSON_ARRAY(JSON_OBJECT('i', 0, 'n', 15, 't', 5, 'v', 4294967295)))))").
		IsEqualToString("7DFFFFFFFF")

	// LEN (2) - Test length-delimited data encoding
	RunTestThatExpression(t, "HEX(pb_wire_json_to_message(JSON_OBJECT('8', JSON_ARRAY(JSON_OBJECT('i', 0, 'n', 8, 't', 2, 'v', 'CAo=')))))").
		IsEqualToString("4202080A")
}

func TestWireJsonToMessageRoundTrip(t *testing.T) {
	// Test round-trip conversion for VARINT
	RunTestThatExpression(t, "pb_wire_json_to_message(pb_message_to_wire_json(_binary X'10ffffffff07'))").
		IsEqualToBytes([]byte{0x10, 0xff, 0xff, 0xff, 0xff, 0x07})

	// Test round-trip conversion for I64
	RunTestThatExpression(t, "pb_wire_json_to_message(pb_message_to_wire_json(_binary X'61ffffffffffffffff'))").
		IsEqualToBytes([]byte{0x61, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	// Test round-trip conversion for I32
	RunTestThatExpression(t, "pb_wire_json_to_message(pb_message_to_wire_json(_binary X'7dffffffff'))").
		IsEqualToBytes([]byte{0x7d, 0xff, 0xff, 0xff, 0xff})

	// Test round-trip conversion for LEN
	RunTestThatExpression(t, "pb_wire_json_to_message(pb_message_to_wire_json(_binary X'4202080a'))").
		IsEqualToBytes([]byte{0x42, 0x02, 0x08, 0x0a})
}

func TestWireJsonToMessageMultipleFields(t *testing.T) {
	// Test multiple fields and ensure round-trip preserves the exact same binary
	RunTestThatExpression(t, "pb_wire_json_to_message(pb_message_to_wire_json(_binary X'080a120568656c6c6f1d2a000000'))").
		IsEqualToBytes([]byte{0x08, 0x0a, 0x12, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x1d, 0x2a, 0x00, 0x00, 0x00})
}

func TestWireJsonToMessageRepeatedFields(t *testing.T) {
	// Test round-trip with repeated fields
	RunTestThatExpression(t, "pb_wire_json_to_message(pb_message_to_wire_json(_binary X'080110020803'))").
		IsEqualToBytes([]byte{0x08, 0x01, 0x10, 0x02, 0x08, 0x03})
}

func TestWireJsonToMessagePreservesIndexOrder(t *testing.T) {
	// Test that elements are ordered by index 'i', not field number 'n'
	// Fields appear in order: field 3 (i=0), field 1 (i=1), field 2 (i=2)
	wire_json := `{"1": [{"i": 1, "n": 1, "t": 0, "v": 10}], "2": [{"i": 2, "n": 2, "t": 0, "v": 20}], "3": [{"i": 0, "n": 3, "t": 0, "v": 30}]}`

	// Should encode in index order: field 3 first (i=0), then field 1 (i=1), then field 2 (i=2)
	RunTestThatExpression(t, "HEX(pb_wire_json_to_message('"+wire_json+"'))").
		IsEqualToString("181E080A1014") // 0x18 0x1e (field 3, value 30), 0x08 0x0a (field 1, value 10), 0x10 0x14 (field 2, value 20)
}

func TestWireEncodingFunctions(t *testing.T) {
	// Test binary utility functions
	RunTestThatExpression(t, "HEX(_pb_util_uint32_to_bin(0x12345678))").
		IsEqualToString("12345678") // Big-endian encoding

	RunTestThatExpression(t, "HEX(_pb_util_uint64_to_bin(0x123456789ABCDEF0))").
		IsEqualToString("123456789ABCDEF0") // Big-endian encoding
}

func TestWireJsonToMessageEmpty(t *testing.T) {
	// Test empty wire JSON
	RunTestThatExpression(t, "pb_wire_json_to_message('{}')").
		IsEqualToBytes([]byte{})
}
