package main

import "testing"

func TestWireReadVarintAsUint64(t *testing.T) {
	RunTestThatExpression(t, "_pb_wire_read_varint_as_uint64(_binary X'9601')").IsEqualToUint(150)
	RunTestThatExpression(t, "_pb_wire_read_varint_as_uint64(_binary X'960133')").IsEqualToUint(150)
	RunTestThatExpression(t, "_pb_wire_read_varint_as_uint64(_binary X'10')").IsEqualToUint(16)
}

func TestWireGetFieldNumberFromTag(t *testing.T) {
	RunTestThatExpression(t, "_pb_wire_get_field_number_from_tag(0x08)").IsEqualToUint(1)
}

func TestWireGetWireTypeFromTag(t *testing.T) {
	RunTestThatExpression(t, "_pb_wire_get_wire_type_from_tag(0x08)").IsEqualToUint(0)
}

func TestWireTypeName(t *testing.T) {
	RunTestThatExpression(t, "_pb_wire_type_name(0)").IsEqualToString("VARINT")
	RunTestThatExpression(t, "_pb_wire_type_name(10)").ToFailWithSignalException("45000", "_pb_wire_type_name: unsupported wire_type")
}
