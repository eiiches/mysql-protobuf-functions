package main

import (
	"testing"
)

func TestUtilSwapEndian(t *testing.T) {
	AssertThatExpression(t, "_pb_util_swap_endian(0x000000000000ffff)").IsEqualToUint(0xffff000000000000)
	AssertThatExpression(t, "_pb_util_swap_endian(0xffffffffffffffff)").IsEqualToUint(0xffffffffffffffff)
}

func TestUtilReinterpretUint64AsDouble(t *testing.T) {
	// https://en.wikipedia.org/wiki/Double-precision_floating-point_format
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0011111111110000000000000000000000000000000000000000000000000000)").IsEqualToDouble(1)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0011111111110000000000000000000000000000000000000000000000000001)").IsEqualToDouble(1.0000000000000002220)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0011111111110000000000000000000000000000000000000000000000000010)").IsEqualToDouble(1.0000000000000004441)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000000000000000000000000000000000000000000000000000000000)").IsEqualToDouble(2)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b1100000000000000000000000000000000000000000000000000000000000000)").IsEqualToDouble(-2)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000001000000000000000000000000000000000000000000000000000)").IsEqualToDouble(3)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000010000000000000000000000000000000000000000000000000000)").IsEqualToDouble(4)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000010100000000000000000000000000000000000000000000000000)").IsEqualToDouble(5)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000011000000000000000000000000000000000000000000000000000)").IsEqualToDouble(6)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000110111000000000000000000000000000000000000000000000000)").IsEqualToDouble(23)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0011111110001000000000000000000000000000000000000000000000000000)").IsEqualToDouble(0.01171875)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0000000000000000000000000000000000000000000000000000000000000001)").IsEqualToDouble(4.9406564584124654e-324)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0000000000001111111111111111111111111111111111111111111111111111)").IsEqualToDouble(2.2250738585072009e-308)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0000000000010000000000000000000000000000000000000000000000000000)").IsEqualToDouble(2.2250738585072014e-308)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111101111111111111111111111111111111111111111111111111111)").IsEqualToDouble(1.7976931348623157e+308)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0000000000000000000000000000000000000000000000000000000000000000)").IsEqualToDouble(0)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b1000000000000000000000000000000000000000000000000000000000000000)").IsEqualToDouble(-0)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111110000000000000000000000000000000000000000000000000000)").IsNull() /* +Inf */
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b1111111111110000000000000000000000000000000000000000000000000000)").IsNull() /* -Inf */
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111110000000000000000000000000000000000000000000000000001)").IsNull() /* sNaN */
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111111000000000000000000000000000000000000000000000000001)").IsNull() /* qNaN */
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111111111111111111111111111111111111111111111111111111111)").IsNull() /* NaN */
}

func TestUtilReinterpretUint64AsInt64(t *testing.T) {
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x0000000000000000)").IsEqualToInt(0)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x7fffffffffffffff)").IsEqualToInt(9223372036854775807)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x8000000000000000)").IsEqualToInt(-9223372036854775808)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x8000000000000400)").ToSucceed()
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x8000000000000401)").ToSucceed()
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0xfffffffffffffffe)").IsEqualToInt(-2)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0xffffffffffffffff)").IsEqualToInt(-1)
}

func TestUtilReinterpretUint64AsInt32(t *testing.T) {
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_int32(0xffffffff)").IsEqualToInt(-1)
}

func TestUtilReinterpretUint64AsUint32(t *testing.T) {
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_uint32(0xffffffff)").IsEqualToUint(0xffffffff)
}

func TestUtilReinterpretUint64AsSint64(t *testing.T) {
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(0)").IsEqualToInt(0)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(1)").IsEqualToInt(-1)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(2)").IsEqualToInt(1)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(3)").IsEqualToInt(-2)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(0xfffffffe)").IsEqualToInt(0x7fffffff)
	AssertThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(0xffffffff)").IsEqualToInt(-0x80000000)
}

func TestUtilBinAsUint64(t *testing.T) {
	AssertThatExpression(t, "_pb_util_bin_as_uint64(_binary X'00')").IsEqualToUint(0)
	AssertThatExpression(t, "_pb_util_bin_as_uint64(_binary X'7fffffffffffffff')").IsEqualToUint(9223372036854775807)
	AssertThatExpression(t, "_pb_util_bin_as_uint64(_binary X'8000000000000000')").IsEqualToUint(9223372036854775808)
	AssertThatExpression(t, "_pb_util_bin_as_uint64(_binary X'ffffffffffffffff')").IsEqualToUint(18446744073709551615)
}

func TestUtilBinAsInt64(t *testing.T) {
	AssertThatExpression(t, "_pb_util_bin_as_int64(_binary X'00')").IsEqualToInt(0)
	AssertThatExpression(t, "_pb_util_bin_as_int64(_binary X'7fffffffffffffff')").IsEqualToInt(9223372036854775807)
	AssertThatExpression(t, "_pb_util_bin_as_int64(_binary X'8000000000000000')").IsEqualToInt(-9223372036854775808)
	AssertThatExpression(t, "_pb_util_bin_as_int64(_binary X'ffffffffffffffff')").IsEqualToInt(-1)
}

func TestUtilBinAsInt32(t *testing.T) {
	AssertThatExpression(t, "_pb_util_bin_as_int32(_binary X'00')").IsEqualToInt(0)
	AssertThatExpression(t, "_pb_util_bin_as_int32(_binary X'7fffffff')").IsEqualToInt(2147483647)
	AssertThatExpression(t, "_pb_util_bin_as_int32(_binary X'80000000')").IsEqualToInt(-2147483648)
	AssertThatExpression(t, "_pb_util_bin_as_int32(_binary X'ffffffff')").IsEqualToInt(-1)
}

func TestUtilBinAsUint32(t *testing.T) {
	AssertThatExpression(t, "_pb_util_bin_as_uint32(_binary X'00')").IsEqualToUint(0)
	AssertThatExpression(t, "_pb_util_bin_as_uint32(_binary X'7fffffff')").IsEqualToUint(2147483647)
	AssertThatExpression(t, "_pb_util_bin_as_uint32(_binary X'80000000')").IsEqualToUint(2147483648)
	AssertThatExpression(t, "_pb_util_bin_as_uint32(_binary X'ffffffff')").IsEqualToUint(4294967295)
}
