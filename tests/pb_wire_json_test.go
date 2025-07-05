package main

import "testing"

func TestWireJsonGetUint32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_uint32_field(pb_message_to_wire_json(_binary X''), 2, 0)").IsEqualToUint(0)
	RunTestThatExpression(t, "pb_wire_json_get_uint32_field(pb_message_to_wire_json(_binary X'10ffffffff07'), 2, 0)").IsEqualToUint(2147483647)
	RunTestThatExpression(t, "pb_wire_json_get_uint32_field(pb_message_to_wire_json(_binary X'108080808008'), 2, 0)").IsEqualToUint(2147483648)
	RunTestThatExpression(t, "pb_wire_json_get_uint32_field(pb_message_to_wire_json(_binary X'10ffffffff0f'), 2, 0)").IsEqualToUint(4294967295)
}

func TestWireJsonGetRepeatedUint32FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_uint32_field_count(pb_message_to_wire_json(_binary X''), 2)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_uint32_field_count(pb_message_to_wire_json(_binary X'10ffffffff07'), 2)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_uint32_field_count(pb_message_to_wire_json(_binary X'10ffffffff0710ffffffff07'), 2)").IsEqualToInt(2)
}

func TestWireJsonHasUint32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_uint32_field(pb_message_to_wire_json(_binary X''), 2)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_uint32_field(pb_message_to_wire_json(_binary X'10ffffffff07'), 2)").IsTrue()
}

func TestWireJsonGetUint64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_uint64_field(pb_message_to_wire_json(_binary X''), 4, 0)").IsEqualToUint(0)
	RunTestThatExpression(t, "pb_wire_json_get_uint64_field(pb_message_to_wire_json(_binary X'20ffffffffffffffffff01'), 4, 0)").IsEqualToUint(18446744073709551615)
}

func TestWireJsonGetRepeatedUint64FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_uint64_field_count(pb_message_to_wire_json(_binary X''), 4)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_uint64_field_count(pb_message_to_wire_json(_binary X'20ffffffffffffffffff01'), 4)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_uint64_field_count(pb_message_to_wire_json(_binary X'20ffffffffffffffffff0120ffffffffffffffffff01'), 4)").IsEqualToInt(2)
}

func TestWireJsonHasUint64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_uint64_field(pb_message_to_wire_json(_binary X''), 4)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_uint64_field(pb_message_to_wire_json(_binary X'20ffffffffffffffffff01'), 4)").IsTrue()
}

func TestWireJsonGetSint32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_sint32_field(pb_message_to_wire_json(_binary X''), 9, 0)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_sint32_field(pb_message_to_wire_json(_binary X'48feffffff0f'), 9, 0)").IsEqualToInt(2147483647)
	RunTestThatExpression(t, "pb_wire_json_get_sint32_field(pb_message_to_wire_json(_binary X'4801'), 9, 0)").IsEqualToInt(-1)
	RunTestThatExpression(t, "pb_wire_json_get_sint32_field(pb_message_to_wire_json(_binary X'48ffffffff0f'), 9, 0)").IsEqualToInt(-2147483648)
}

func TestWireJsonGetRepeatedSint32FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sint32_field_count(pb_message_to_wire_json(_binary X''), 9)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sint32_field_count(pb_message_to_wire_json(_binary X'48feffffff0f'), 9)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sint32_field_count(pb_message_to_wire_json(_binary X'48feffffff0f48feffffff0f'), 9)").IsEqualToInt(2)
}

func TestWireJsonHasSint32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_sint32_field(pb_message_to_wire_json(_binary X''), 9)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_sint32_field(pb_message_to_wire_json(_binary X'48feffffff0f'), 9)").IsTrue()
}

func TestWireJsonGetSint64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_sint64_field(pb_message_to_wire_json(_binary X''), 10, 0)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_sint64_field(pb_message_to_wire_json(_binary X'50feffffffffffffffff01'), 10, 0)").IsEqualToInt(9223372036854775807)
}

func TestWireJsonGetRepeatedSint64FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sint64_field_count(pb_message_to_wire_json(_binary X''), 10)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sint64_field_count(pb_message_to_wire_json(_binary X'50feffffffffffffffff01'), 10)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sint64_field_count(pb_message_to_wire_json(_binary X'50feffffffffffffffff0150feffffffffffffffff01'), 10)").IsEqualToInt(2)
}

func TestWireJsonHasSint64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_sint64_field(pb_message_to_wire_json(_binary X''), 10)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_sint64_field(pb_message_to_wire_json(_binary X'50feffffffffffffffff01'), 10)").IsTrue()
}

func TestWireJsonGetInt32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_int32_field(pb_message_to_wire_json(_binary X'100a080a'), 1, 0)").IsEqualToInt(10)
	RunTestThatExpression(t, "pb_wire_json_get_int32_field(pb_message_to_wire_json(_binary X'100a080a'), 1, 0)").IsEqualToInt(10)

	RunTestThatExpression(t, "pb_wire_json_get_int32_field(pb_message_to_wire_json(_binary X''), 1, 0)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_int32_field(pb_message_to_wire_json(_binary X'08ffffffff07'), 1, 0)").IsEqualToInt(2147483647)
	RunTestThatExpression(t, "pb_wire_json_get_int32_field(pb_message_to_wire_json(_binary X'08ffffffffffffffffff01'), 1, 0)").IsEqualToInt(-1)
	RunTestThatExpression(t, "pb_wire_json_get_int32_field(pb_message_to_wire_json(_binary X'0880808080f8ffffffff01'), 1, 0)").IsEqualToInt(-2147483648)
}

func TestWireJsonGetRepeatedInt32Field(t *testing.T) {
	// packed repeated
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int32_field_element(pb_message_to_wire_json(_binary X'3a03010203'), 7, 0)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int32_field_element(pb_message_to_wire_json(_binary X'3a03010203'), 7, 1)").IsEqualToInt(2)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int32_field_element(pb_message_to_wire_json(_binary X'3a03010203'), 7, 2)").IsEqualToInt(3)
}

func TestWireJsonGetRepeatedInt32FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int32_field_count(pb_message_to_wire_json(_binary X''), 1)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int32_field_count(pb_message_to_wire_json(_binary X'100a080a'), 1)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int32_field_count(pb_message_to_wire_json(_binary X'100a080a100a080a'), 1)").IsEqualToInt(2)

	// packed repeated
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int32_field_count(pb_message_to_wire_json(_binary X'3a03010203'), 7)").IsEqualToInt(3)
}

func TestWireJsonHasInt32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_int32_field(pb_message_to_wire_json(_binary X''), 1)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_int32_field(pb_message_to_wire_json(_binary X'100a080a'), 1)").IsTrue()
}

func TestWireJsonGetInt64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_int64_field(pb_message_to_wire_json(_binary X'100a080a'), 1, 0)").IsEqualToInt(10)
	RunTestThatExpression(t, "pb_wire_json_get_int64_field(pb_message_to_wire_json(_binary X'100a080a'), 1, 0)").IsEqualToInt(10)

	RunTestThatExpression(t, "pb_wire_json_get_int64_field(pb_message_to_wire_json(_binary X''), 3, 0)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_int64_field(pb_message_to_wire_json(_binary X'18ffffffffffffffff7f'), 3, 0)").IsEqualToInt(9223372036854775807)
}

func TestWireJsonGetRepeatedInt64FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int64_field_count(pb_message_to_wire_json(_binary X''), 1)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int64_field_count(pb_message_to_wire_json(_binary X'100a080a'), 1)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_int64_field_count(pb_message_to_wire_json(_binary X'100a080a100a080a'), 1)").IsEqualToInt(2)
}

func TestWireJsonHasInt64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_int64_field(pb_message_to_wire_json(_binary X''), 1)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_int64_field(pb_message_to_wire_json(_binary X'100a080a'), 1)").IsTrue()
}

func TestWireJsonGetStringField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_string_field(pb_message_to_wire_json(_binary X''), 5, '')").IsEqualToString("")
	RunTestThatExpression(t, "pb_wire_json_get_string_field(pb_message_to_wire_json(_binary X'100a2a03616263'), 5, '')").IsEqualToString("abc")
}

func TestWireJsonGetRepeatedStringFieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_string_field_count(pb_message_to_wire_json(_binary X''), 5)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_string_field_count(pb_message_to_wire_json(_binary X'100a2a03616263'), 5)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_string_field_count(pb_message_to_wire_json(_binary X'100a2a03616263100a2a03616263'), 5)").IsEqualToInt(2)
}

func TestWireJsonHasStringField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_string_field(pb_message_to_wire_json(_binary X''), 5)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_string_field(pb_message_to_wire_json(_binary X'100a2a03616263'), 5)").IsTrue()
}

func TestWireJsonGetBytesField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_string_field(pb_message_to_wire_json(_binary X'100a2a03616263'), 5, _binary X'')").IsEqualToBytes([]byte("abc"))
}

func TestWireJsonGetRepeatedBytesFieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_bytes_field_count(pb_message_to_wire_json(_binary X''), 5)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_bytes_field_count(pb_message_to_wire_json(_binary X'100a2a03616263'), 5)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_bytes_field_count(pb_message_to_wire_json(_binary X'100a2a03616263100a2a03616263'), 5)").IsEqualToInt(2)
}

func TestWireJsonHasBytesField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_bytes_field(pb_message_to_wire_json(_binary X''), 5)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_bytes_field(pb_message_to_wire_json(_binary X'100a2a03616263'), 5)").IsTrue()
}

func TestWireJsonGetBoolField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_bool_field(pb_message_to_wire_json(_binary X''), 1, FALSE)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_get_bool_field(pb_message_to_wire_json(_binary X'0801'), 1, FALSE)").IsTrue()
}

func TestWireJsonGetRepeatedBoolFieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_bool_field_count(pb_message_to_wire_json(_binary X''), 1)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_bool_field_count(pb_message_to_wire_json(_binary X'0801'), 1)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_bool_field_count(pb_message_to_wire_json(_binary X'08010801'), 1)").IsEqualToInt(2)
}

func TestWireJsonHasBoolField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_bool_field(pb_message_to_wire_json(_binary X''), 1)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_bool_field(pb_message_to_wire_json(_binary X'0801'), 1)").IsTrue()
}

func TestWireJsonGetEnumField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_enum_field(pb_message_to_wire_json(_binary X''), 1, 0)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_enum_field(pb_message_to_wire_json(_binary X'0801'), 1, 0)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_enum_field(pb_message_to_wire_json(_binary X'2805'), 5, 0)").IsEqualToInt(5)
}

func TestWireJsonGetRepeatedEnumFieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_enum_field_count(pb_message_to_wire_json(_binary X''), 1)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_enum_field_count(pb_message_to_wire_json(_binary X'0801'), 1)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_enum_field_count(pb_message_to_wire_json(_binary X'08010801'), 1)").IsEqualToInt(2)
}

func TestWireJsonHasEnumField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_enum_field(pb_message_to_wire_json(_binary X''), 1)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_enum_field(pb_message_to_wire_json(_binary X'0801'), 1)").IsTrue()
}

func TestWireJsonGetFloatField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_float_field(pb_message_to_wire_json(_binary X''), 16, 0)").IsEqualToFloat(0)
	RunTestThatExpression(t, "pb_wire_json_get_float_field(pb_message_to_wire_json(_binary X'85010000c03f'), 16, 0)").IsEqualToFloat(1.5)
}

func TestWireJsonGetRepeatedFloatFieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_float_field_count(pb_message_to_wire_json(_binary X''), 16)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_float_field_count(pb_message_to_wire_json(_binary X'85010000c03f'), 16)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_float_field_count(pb_message_to_wire_json(_binary X'85010000c03f85010000c03f'), 16)").IsEqualToInt(2)
}

func TestWireJsonHasFloatField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_float_field(pb_message_to_wire_json(_binary X''), 16)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_float_field(pb_message_to_wire_json(_binary X'85010000c03f'), 16)").IsTrue()
}

func TestWireJsonGetDoubleField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_double_field(pb_message_to_wire_json(_binary X'69000000000000f83f'), 13, 0)").IsEqualToDouble(1.5)
}

func TestWireJsonGetRepeatedDoubleFieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_double_field_count(pb_message_to_wire_json(_binary X''), 13)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_double_field_count(pb_message_to_wire_json(_binary X'69000000000000f83f'), 13)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_double_field_count(pb_message_to_wire_json(_binary X'69000000000000f83f69000000000000f83f'), 13)").IsEqualToInt(2)
}

func TestWireJsonHasDoubleField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_double_field(pb_message_to_wire_json(_binary X''), 13)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_double_field(pb_message_to_wire_json(_binary X'69000000000000f83f'), 13)").IsTrue()
}

func TestWireJsonGetSfixed64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_sfixed64_field(pb_message_to_wire_json(_binary X''), 11, 0)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_sfixed64_field(pb_message_to_wire_json(_binary X'59ffffffffffffff7f'), 11, 0)").IsEqualToInt(9223372036854775807)
	RunTestThatExpression(t, "pb_wire_json_get_sfixed64_field(pb_message_to_wire_json(_binary X'590000000000000080'), 11, 0)").IsEqualToInt(-9223372036854775808)
}

func TestWireJsonGetRepeatedSfixed64FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sfixed64_field_count(pb_message_to_wire_json(_binary X''), 11)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sfixed64_field_count(pb_message_to_wire_json(_binary X'59ffffffffffffff7f'), 11)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sfixed64_field_count(pb_message_to_wire_json(_binary X'59ffffffffffffff7f59ffffffffffffff7f'), 11)").IsEqualToInt(2)
}

func TestWireJsonHasSfixed64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_sfixed64_field(pb_message_to_wire_json(_binary X''), 11)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_sfixed64_field(pb_message_to_wire_json(_binary X'59ffffffffffffff7f'), 11)").IsTrue()
}

func TestWireJsonGetFixed64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_fixed64_field(pb_message_to_wire_json(_binary X''), 12, 0)").IsEqualToUint(0)
	RunTestThatExpression(t, "pb_wire_json_get_fixed64_field(pb_message_to_wire_json(_binary X'61ffffffffffffffff'), 12, 0)").IsEqualToUint(18446744073709551615)
}

func TestWireJsonGetRepeatedFixed64FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_fixed64_field_count(pb_message_to_wire_json(_binary X''), 12)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_fixed64_field_count(pb_message_to_wire_json(_binary X'61ffffffffffffffff'), 12)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_fixed64_field_count(pb_message_to_wire_json(_binary X'61ffffffffffffffff61ffffffffffffffff'), 12)").IsEqualToInt(2)
}

func TestWireJsonHasFixed64Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_fixed64_field(pb_message_to_wire_json(_binary X''), 12)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_fixed64_field(pb_message_to_wire_json(_binary X'61ffffffffffffffff'), 12)").IsTrue()
}

func TestWireJsonGetSfixed32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_sfixed32_field(pb_message_to_wire_json(_binary X''), 14, 0)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_sfixed32_field(pb_message_to_wire_json(_binary X'75ffffff7f'), 14, 0)").IsEqualToInt(2147483647)
	RunTestThatExpression(t, "pb_wire_json_get_sfixed32_field(pb_message_to_wire_json(_binary X'7500000080'), 14, 0)").IsEqualToInt(-2147483648)
}

func TestWireJsonGetRepeatedSfixed32FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sfixed32_field_count(pb_message_to_wire_json(_binary X''), 14)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sfixed32_field_count(pb_message_to_wire_json(_binary X'75ffffff7f'), 14)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_sfixed32_field_count(pb_message_to_wire_json(_binary X'75ffffff7f75ffffff7f'), 14)").IsEqualToInt(2)
}

func TestWireJsonHasSfixed32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_sfixed32_field(pb_message_to_wire_json(_binary X''), 14)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_sfixed32_field(pb_message_to_wire_json(_binary X'75ffffff7f'), 14)").IsTrue()
}

func TestWireJsonGetFixed32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_fixed32_field(pb_message_to_wire_json(_binary X''), 15, 0)").IsEqualToUint(0)
	RunTestThatExpression(t, "pb_wire_json_get_fixed32_field(pb_message_to_wire_json(_binary X'7dffffffff'), 15, 0)").IsEqualToUint(4294967295)
}

func TestWireJsonGetRepeatedFixed32FieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_fixed32_field_count(pb_message_to_wire_json(_binary X''), 15)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_fixed32_field_count(pb_message_to_wire_json(_binary X'7dffffffff'), 15)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_fixed32_field_count(pb_message_to_wire_json(_binary X'7dffffffff7dffffffff'), 15)").IsEqualToInt(2)
}

func TestWireJsonHasFixed32Field(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_fixed32_field(pb_message_to_wire_json(_binary X''), 15)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_fixed32_field(pb_message_to_wire_json(_binary X'7dffffffff'), 15)").IsTrue()
}

func TestWireJsonGetMessageField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_message_field(pb_message_to_wire_json(_binary X'4202080a'), 8, _binary X'')").IsEqualToBytes([]byte{0x08, 0x0a})
}

func TestWireJsonHasMessageField(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_has_message_field(pb_message_to_wire_json(_binary X''), 8)").IsFalse()
	RunTestThatExpression(t, "pb_wire_json_has_message_field(pb_message_to_wire_json(_binary X'4202080a'), 8)").IsTrue()
}

func TestWireJsonGetRepeatedMessageFieldCount(t *testing.T) {
	RunTestThatExpression(t, "pb_wire_json_get_repeated_message_field_count(pb_message_to_wire_json(_binary X''), 8)").IsEqualToInt(0)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_message_field_count(pb_message_to_wire_json(_binary X'4202080a'), 8)").IsEqualToInt(1)
	RunTestThatExpression(t, "pb_wire_json_get_repeated_message_field_count(pb_message_to_wire_json(_binary X'4202080a4202080a'), 8)").IsEqualToInt(2)
}
