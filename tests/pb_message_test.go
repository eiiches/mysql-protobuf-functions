package main

import "testing"

func TestMessageGetUint32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_uint32_field(_binary X'', 2, NULL)").IsEqualToUint(0)
	AssertThatExpression(t, "pb_message_get_uint32_field(_binary X'10ffffffff07', 2, 0)").IsEqualToUint(2147483647)
	AssertThatExpression(t, "pb_message_get_uint32_field(_binary X'108080808008', 2, 0)").IsEqualToUint(2147483648)
	AssertThatExpression(t, "pb_message_get_uint32_field(_binary X'10ffffffff0f', 2, 0)").IsEqualToUint(4294967295)
}

func TestMessageGetUint64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_uint64_field(_binary X'', 4, NULL)").IsEqualToUint(0)
	AssertThatExpression(t, "pb_message_get_uint64_field(_binary X'20ffffffffffffffffff01', 4, 0)").IsEqualToUint(18446744073709551615)
}

func TestMessageGetSint32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sint32_field(_binary X'', 9, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sint32_field(_binary X'48feffffff0f', 9, 0)").IsEqualToInt(2147483647)
	AssertThatExpression(t, "pb_message_get_sint32_field(_binary X'4801', 9, 0)").IsEqualToInt(-1)
	AssertThatExpression(t, "pb_message_get_sint32_field(_binary X'48ffffffff0f', 9, 0)").IsEqualToInt(-2147483648)
}

func TestMessageGetSint64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sint64_field(_binary X'', 10, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sint64_field(_binary X'50feffffffffffffffff01', 10, 0)").IsEqualToInt(9223372036854775807)
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

func TestMessageGetInt64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_int64_field(_binary X'100a080a', 1, 0)").IsEqualToInt(10)
	AssertThatExpression(t, "pb_message_get_int64_field(_binary X'100a080a', 1, NULL)").IsEqualToInt(10)

	AssertThatExpression(t, "pb_message_get_int64_field(_binary X'', 3, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_int64_field(_binary X'18ffffffffffffffff7f', 3, 0)").IsEqualToInt(9223372036854775807)
}

func TestMessageGetStringField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_string_field(_binary X'100a2a03616263', 5, NULL)").IsEqualToString("abc")
}

func TestMessageGetBytesField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_string_field(_binary X'100a2a03616263', 5, NULL)").IsEqualToBytes([]byte("abc"))
}

func TestMessageGetBoolField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_bool_field(_binary X'', 1, NULL)").IsFalse()
	AssertThatExpression(t, "pb_message_get_bool_field(_binary X'0801', 1, NULL)").IsTrue()
}

func TestMessageGetEnumField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_enum_field(_binary X'', 1, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_enum_field(_binary X'0801', 1, NULL)").IsEqualToInt(1)
}

func TestMessageGetDoubleField(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_double_field(_binary X'69000000000000f83f', 13, 0)").IsEqualToDouble(1.5)
}

func TestMessageGetSfixed64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sfixed64_field(_binary X'', 11, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sfixed64_field(_binary X'59ffffffffffffff7f', 11, 0)").IsEqualToInt(9223372036854775807)
	AssertThatExpression(t, "pb_message_get_sfixed64_field(_binary X'590000000000000080', 11, 0)").IsEqualToInt(-9223372036854775808)
}

func TestMessageGetFixed64Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_fixed64_field(_binary X'', 12, NULL)").IsEqualToUint(0)
	AssertThatExpression(t, "pb_message_get_fixed64_field(_binary X'61ffffffffffffffff', 12, 0)").IsEqualToUint(18446744073709551615)
}
func TestMessageGetSfixed32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_sfixed32_field(_binary X'', 14, NULL)").IsEqualToInt(0)
	AssertThatExpression(t, "pb_message_get_sfixed32_field(_binary X'75ffffff7f', 14, 0)").IsEqualToInt(2147483647)
	AssertThatExpression(t, "pb_message_get_sfixed32_field(_binary X'7500000080', 14, 0)").IsEqualToInt(-2147483648)
}

func TestMessageGetFixed32Field(t *testing.T) {
	AssertThatExpression(t, "pb_message_get_fixed32_field(_binary X'', 15, NULL)").IsEqualToUint(0)
	AssertThatExpression(t, "pb_message_get_fixed32_field(_binary X'7dffffffff', 15, 0)").IsEqualToUint(4294967295)
}
