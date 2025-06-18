package main

import (
	"testing"
)

func TestMessageGetUint32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_uint32_field(_binary X'', 2, NULL)").IsEqualToUint(0)
	AssertThatExpression(t, "pb_message_get_uint32_field(_binary X'10ffffffff07', 2, 0)").IsEqualToUint(2147483647)
	AssertThatExpression(t, "pb_message_get_uint32_field(_binary X'108080808008', 2, 0)").IsEqualToUint(2147483648)
	AssertThatExpression(t, "pb_message_get_uint32_field(_binary X'10ffffffff0f', 2, 0)").IsEqualToUint(4294967295)
}

func TestMessageGetUint32FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_uint32_field_count(_binary X'', 2)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_uint32_field_count(_binary X'10ffffffff07', 2)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_uint32_field_count(_binary X'10ffffffff0710ffffffff07', 2)").IsEqualToInt(2)
}

func TestMessageHasUint32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_uint32_field(_binary X'', 2)").IsFalse()
	AssertThatExpression(t, "pb_message_has_uint32_field(_binary X'10ffffffff07', 2)").IsTrue()
}

func TestMessageGetUint64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_uint64_field(_binary X'', 4, NULL)").IsEqualToUint(0)
	AssertThatExpression(t, "pb_message_get_uint64_field(_binary X'20ffffffffffffffffff01', 4, 0)").IsEqualToUint(18446744073709551615)
}

func TestMessageGetUint64FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_uint64_field_count(_binary X'', 4)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_uint64_field_count(_binary X'20ffffffffffffffffff01', 4)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_uint64_field_count(_binary X'20ffffffffffffffffff0120ffffffffffffffffff01', 4)").IsEqualToInt(2)
}

func TestMessageHasUint64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_uint64_field(_binary X'', 4)").IsFalse()
	AssertThatExpression(t, "pb_message_has_uint64_field(_binary X'20ffffffffffffffffff01', 4)").IsTrue()
}

func TestMessageGetSint32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sint32_field(_binary X'', 9, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sint32_field(_binary X'48feffffff0f', 9, 0)").IsEqualToInt(2147483647)
	AssertThatExpression(t, "pb_message_get_sint32_field(_binary X'4801', 9, 0)").IsEqualToInt(-1)
	AssertThatExpression(t, "pb_message_get_sint32_field(_binary X'48ffffffff0f', 9, 0)").IsEqualToInt(-2147483648)
}

func TestMessageGetSint32FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sint32_field_count(_binary X'', 9)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sint32_field_count(_binary X'48feffffff0f', 9)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_sint32_field_count(_binary X'48feffffff0f48feffffff0f', 9)").IsEqualToInt(2)
}

func TestMessageHasSint32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_sint32_field(_binary X'', 9)").IsFalse()
	AssertThatExpression(t, "pb_message_has_sint32_field(_binary X'48feffffff0f', 9)").IsTrue()
}

func TestMessageGetSint64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sint64_field(_binary X'', 10, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sint64_field(_binary X'50feffffffffffffffff01', 10, 0)").IsEqualToInt(9223372036854775807)
}

func TestMessageGetSint64FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sint64_field_count(_binary X'', 10)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sint64_field_count(_binary X'50feffffffffffffffff01', 10)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_sint64_field_count(_binary X'50feffffffffffffffff0150feffffffffffffffff01', 10)").IsEqualToInt(2)
}

func TestMessageHasSint64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_sint64_field(_binary X'', 10)").IsFalse()
	AssertThatExpression(t, "pb_message_has_sint64_field(_binary X'50feffffffffffffffff01', 10)").IsTrue()
}

func TestMessageGetInt32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_int32_field(_binary X'100a080a', 1, 0)").IsEqualToInt(10)
	AssertThatExpression(t, "pb_message_get_int32_field(_binary X'100a080a', 1, NULL)").IsEqualToInt(10)

	AssertThatExpression(t, "pb_message_get_int32_field(_binary X'', 1, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_int32_field(_binary X'08ffffffff07', 1, 0)").IsEqualToInt(2147483647)
	AssertThatExpression(t, "pb_message_get_int32_field(_binary X'08ffffffffffffffffff01', 1, 0)").IsEqualToInt(-1)
	AssertThatExpression(t, "pb_message_get_int32_field(_binary X'0880808080f8ffffffff01', 1, 0)").IsEqualToInt(-2147483648)

	// packed repeated
	AssertThatExpression(t, "pb_message_get_int32_field(_binary X'3a03010203', 7, 0)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_int32_field(_binary X'3a03010203', 7, 1)").IsEqualToInt(2)
	AssertThatExpression(t, "pb_message_get_int32_field(_binary X'3a03010203', 7, 2)").IsEqualToInt(3)
}

func TestMessageGetInt32FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_int32_field_count(_binary X'', 1)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_int32_field_count(_binary X'100a080a', 1)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_int32_field_count(_binary X'100a080a100a080a', 1)").IsEqualToInt(2)

	// packed repeated
	AssertThatExpression(t, "pb_message_get_int32_field_count(_binary X'3a03010203', 7)").IsEqualToInt(3)
}

func TestMessageHasInt32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_int32_field(_binary X'', 1)").IsFalse()
	AssertThatExpression(t, "pb_message_has_int32_field(_binary X'100a080a', 1)").IsTrue()
}

func TestMessageGetInt64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_int64_field(_binary X'100a080a', 1, 0)").IsEqualToInt(10)
	AssertThatExpression(t, "pb_message_get_int64_field(_binary X'100a080a', 1, NULL)").IsEqualToInt(10)

	AssertThatExpression(t, "pb_message_get_int64_field(_binary X'', 3, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_int64_field(_binary X'18ffffffffffffffff7f', 3, 0)").IsEqualToInt(9223372036854775807)
}

func TestMessageGetInt64FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_int64_field_count(_binary X'', 1)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_int64_field_count(_binary X'100a080a', 1)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_int64_field_count(_binary X'100a080a100a080a', 1)").IsEqualToInt(2)
}

func TestMessageHasInt64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_int64_field(_binary X'', 1)").IsFalse()
	AssertThatExpression(t, "pb_message_has_int64_field(_binary X'100a080a', 1)").IsTrue()
}

func TestMessageGetStringField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_string_field(_binary X'', 5, NULL)").IsEqualToString("")
	AssertThatExpression(t, "pb_message_get_string_field(_binary X'100a2a03616263', 5, NULL)").IsEqualToString("abc")
}

func TestMessageGetStringFieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_string_field_count(_binary X'', 5)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_string_field_count(_binary X'100a2a03616263', 5)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_string_field_count(_binary X'100a2a03616263100a2a03616263', 5)").IsEqualToInt(2)
}

func TestMessageHasStringField(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_string_field(_binary X'', 5)").IsFalse()
	AssertThatExpression(t, "pb_message_has_string_field(_binary X'100a2a03616263', 5)").IsTrue()
}

func TestMessageGetBytesField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_string_field(_binary X'100a2a03616263', 5, NULL)").IsEqualToBytes([]byte("abc"))
}

func TestMessageGetBytesFieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_bytes_field_count(_binary X'', 5)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_bytes_field_count(_binary X'100a2a03616263', 5)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_bytes_field_count(_binary X'100a2a03616263100a2a03616263', 5)").IsEqualToInt(2)
}

func TestMessageHasBytesField(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_bytes_field(_binary X'', 5)").IsFalse()
	AssertThatExpression(t, "pb_message_has_bytes_field(_binary X'100a2a03616263', 5)").IsTrue()
}

func TestMessageGetBoolField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_bool_field(_binary X'', 1, NULL)").IsFalse()
	AssertThatExpression(t, "pb_message_get_bool_field(_binary X'0801', 1, NULL)").IsTrue()
}

func TestMessageGetBoolFieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_bool_field_count(_binary X'', 1)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_bool_field_count(_binary X'0801', 1)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_bool_field_count(_binary X'08010801', 1)").IsEqualToInt(2)
}

func TestMessageHasBoolField(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_bool_field(_binary X'', 1)").IsFalse()
	AssertThatExpression(t, "pb_message_has_bool_field(_binary X'0801', 1)").IsTrue()
}

func TestMessageGetEnumField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_enum_field(_binary X'', 1, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_enum_field(_binary X'0801', 1, NULL)").IsEqualToInt(1)
}

func TestMessageGetEnumFieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_enum_field_count(_binary X'', 1)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_enum_field_count(_binary X'0801', 1)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_enum_field_count(_binary X'08010801', 1)").IsEqualToInt(2)
}

func TestMessageHasEnumField(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_enum_field(_binary X'', 1)").IsFalse()
	AssertThatExpression(t, "pb_message_has_enum_field(_binary X'0801', 1)").IsTrue()
}

func TestMessageGetFloatField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_float_field(_binary X'', 16, NULL)").IsEqualToFloat(0)
	AssertThatExpression(t, "pb_message_get_float_field(_binary X'85010000c03f', 16, 0)").IsEqualToFloat(1.5)
}

func TestMessageGetFloatFieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_float_field_count(_binary X'', 16)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_float_field_count(_binary X'85010000c03f', 16)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_float_field_count(_binary X'85010000c03f85010000c03f', 16)").IsEqualToInt(2)
}

func TestMessageHasFloatField(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_float_field(_binary X'', 16)").IsFalse()
	AssertThatExpression(t, "pb_message_has_float_field(_binary X'85010000c03f', 16)").IsTrue()
}

func TestMessageGetDoubleField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_double_field(_binary X'69000000000000f83f', 13, 0)").IsEqualToDouble(1.5)
}

func TestMessageGetDoubleFieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_double_field_count(_binary X'', 13)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_double_field_count(_binary X'69000000000000f83f', 13)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_double_field_count(_binary X'69000000000000f83f69000000000000f83f', 13)").IsEqualToInt(2)
}

func TestMessageHasDoubleField(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_double_field(_binary X'', 13)").IsFalse()
	AssertThatExpression(t, "pb_message_has_double_field(_binary X'69000000000000f83f', 13)").IsTrue()
}

func TestMessageGetSfixed64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sfixed64_field(_binary X'', 11, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sfixed64_field(_binary X'59ffffffffffffff7f', 11, 0)").IsEqualToInt(9223372036854775807)
	AssertThatExpression(t, "pb_message_get_sfixed64_field(_binary X'590000000000000080', 11, 0)").IsEqualToInt(-9223372036854775808)
}

func TestMessageGetSfixed64FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sfixed64_field_count(_binary X'', 11)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sfixed64_field_count(_binary X'59ffffffffffffff7f', 11)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_sfixed64_field_count(_binary X'59ffffffffffffff7f59ffffffffffffff7f', 11)").IsEqualToInt(2)
}

func TestMessageHasSfixed64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_sfixed64_field(_binary X'', 11)").IsFalse()
	AssertThatExpression(t, "pb_message_has_sfixed64_field(_binary X'59ffffffffffffff7f', 11)").IsTrue()
}

func TestMessageGetFixed64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_fixed64_field(_binary X'', 12, NULL)").IsEqualToUint(0)
	AssertThatExpression(t, "pb_message_get_fixed64_field(_binary X'61ffffffffffffffff', 12, 0)").IsEqualToUint(18446744073709551615)
}

func TestMessageGetFixed64FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_fixed64_field_count(_binary X'', 12)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_fixed64_field_count(_binary X'61ffffffffffffffff', 12)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_fixed64_field_count(_binary X'61ffffffffffffffff61ffffffffffffffff', 12)").IsEqualToInt(2)
}

func TestMessageHasFixed64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_fixed64_field(_binary X'', 12)").IsFalse()
	AssertThatExpression(t, "pb_message_has_fixed64_field(_binary X'61ffffffffffffffff', 12)").IsTrue()
}

func TestMessageGetSfixed32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sfixed32_field(_binary X'', 14, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sfixed32_field(_binary X'75ffffff7f', 14, 0)").IsEqualToInt(2147483647)
	AssertThatExpression(t, "pb_message_get_sfixed32_field(_binary X'7500000080', 14, 0)").IsEqualToInt(-2147483648)
}

func TestMessageGetSfixed32FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sfixed32_field_count(_binary X'', 14)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sfixed32_field_count(_binary X'75ffffff7f', 14)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_sfixed32_field_count(_binary X'75ffffff7f75ffffff7f', 14)").IsEqualToInt(2)
}

func TestMessageHasSfixed32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_sfixed32_field(_binary X'', 14)").IsFalse()
	AssertThatExpression(t, "pb_message_has_sfixed32_field(_binary X'75ffffff7f', 14)").IsTrue()
}

func TestMessageGetFixed32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_fixed32_field(_binary X'', 15, NULL)").IsEqualToUint(0)
	AssertThatExpression(t, "pb_message_get_fixed32_field(_binary X'7dffffffff', 15, 0)").IsEqualToUint(4294967295)
}

func TestMessageGetFixed32FieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_fixed32_field_count(_binary X'', 15)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_fixed32_field_count(_binary X'7dffffffff', 15)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_fixed32_field_count(_binary X'7dffffffff7dffffffff', 15)").IsEqualToInt(2)
}

func TestMessageHasFixed32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_fixed32_field(_binary X'', 15)").IsFalse()
	AssertThatExpression(t, "pb_message_has_fixed32_field(_binary X'7dffffffff', 15)").IsTrue()
}

func TestMessageGetMessageField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_message_field(_binary X'4202080a', 8, NULL)").IsEqualToBytes([]byte{0x08, 0x0a})
}

func TestMessageHasMessageField(t *testing.T) {
	AssertThatExpression(t, "pb_message_has_message_field(_binary X'', 8)").IsFalse()
	AssertThatExpression(t, "pb_message_has_message_field(_binary X'4202080a', 8)").IsTrue()
}

func TestMessageGetMessageFieldCount(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_message_field_count(_binary X'', 8)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_message_field_count(_binary X'4202080a', 8)").IsEqualToInt(1)
	AssertThatExpression(t, "pb_message_get_message_field_count(_binary X'4202080a4202080a', 8)").IsEqualToInt(2)
}
