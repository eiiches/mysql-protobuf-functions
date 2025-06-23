package main

import (
	"testing"
)

func TestUtilSwapEndian32(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_swap_endian_32(0x0000ffff)").IsEqualToUint(0xffff0000)
	RunTestThatExpression(t, "_pb_util_swap_endian_32(0xffffffff)").IsEqualToUint(0xffffffff)
}

func TestUtilSwapEndian64(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_swap_endian_64(0x000000000000ffff)").IsEqualToUint(0xffff000000000000)
	RunTestThatExpression(t, "_pb_util_swap_endian_64(0xffffffffffffffff)").IsEqualToUint(0xffffffffffffffff)
}

func TestUtilReinterpretUint64AsDouble(t *testing.T) {
	// https://en.wikipedia.org/wiki/Double-precision_floating-point_format
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0011111111110000000000000000000000000000000000000000000000000000)").IsEqualToDouble(1)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0011111111110000000000000000000000000000000000000000000000000001)").IsEqualToDouble(1.0000000000000002220)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0011111111110000000000000000000000000000000000000000000000000010)").IsEqualToDouble(1.0000000000000004441)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000000000000000000000000000000000000000000000000000000000)").IsEqualToDouble(2)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b1100000000000000000000000000000000000000000000000000000000000000)").IsEqualToDouble(-2)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000001000000000000000000000000000000000000000000000000000)").IsEqualToDouble(3)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000010000000000000000000000000000000000000000000000000000)").IsEqualToDouble(4)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000010100000000000000000000000000000000000000000000000000)").IsEqualToDouble(5)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000011000000000000000000000000000000000000000000000000000)").IsEqualToDouble(6)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0100000000110111000000000000000000000000000000000000000000000000)").IsEqualToDouble(23)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0011111110001000000000000000000000000000000000000000000000000000)").IsEqualToDouble(0.01171875)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0000000000000000000000000000000000000000000000000000000000000001)").IsEqualToDouble(4.9406564584124654e-324)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0000000000001111111111111111111111111111111111111111111111111111)").IsEqualToDouble(2.2250738585072009e-308)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0000000000010000000000000000000000000000000000000000000000000000)").IsEqualToDouble(2.2250738585072014e-308)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111101111111111111111111111111111111111111111111111111111)").IsEqualToDouble(1.7976931348623157e+308)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0000000000000000000000000000000000000000000000000000000000000000)").IsEqualToDouble(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b1000000000000000000000000000000000000000000000000000000000000000)").IsEqualToDouble(-0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111110000000000000000000000000000000000000000000000000000)").IsNull() /* +Inf */
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b1111111111110000000000000000000000000000000000000000000000000000)").IsNull() /* -Inf */
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111110000000000000000000000000000000000000000000000000001)").IsNull() /* sNaN */
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111111000000000000000000000000000000000000000000000000001)").IsNull() /* qNaN */
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0b0111111111111111111111111111111111111111111111111111111111111111)").IsNull() /* NaN */
}

func TestUtilReinterpretUint32AsFloat(t *testing.T) {
	// https://en.wikipedia.org/wiki/Single-precision_floating-point_format

	// Subnormal range
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x00000001)").IsEqualToFloat(1.4012984643e-45) // Smallest positive subnormal number
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x007FFFFF)").IsEqualToFloat(1.1754942107e-38) // Largest subnormal number

	// Normal range
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x00800000)").IsEqualToFloat(1.1754943508e-38) // Smallest positive normal number
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x7F7FFFFF)").IsEqualToFloat(3.4028234664e+38) // Largest normal number
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x3F7FFFFF)").IsEqualToFloat(0.9999999404)     // Largest number less than one
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x3F800000)").IsEqualToFloat(1)                // 1.0
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x3F800001)").IsEqualToFloat(1.0000001192)     // Smallest number larger than one
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0xC0000000)").IsEqualToFloat(-2)               // -2.0
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x00000000)").IsEqualToFloat(0)                // +0.0
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x80000000)").IsEqualToFloat(-0)               // -0.0

	// Special values
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x7F800000)").IsNull() /* +Inf */
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0xFF800000)").IsNull() /* -Inf */

	// Common constants
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x40490FDB)").IsEqualToFloat(3.1415927410) // Pi (π)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x3EAAAAAB)").IsEqualToFloat(0.3333333433) // 1/3

	// NaNs
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0xFFC00001)").IsNull() /* qNaN */
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0xFF800001)").IsNull() /* sNaN */
}

func TestUtilReinterpretUint32AsInt32(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(0x00000000)").IsEqualToInt(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(0x7fffffff)").IsEqualToInt(2147483647)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(0x80000000)").IsEqualToInt(-2147483648)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(0x80000400)").IsEqualToInt(-2147482624)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(0x80000401)").IsEqualToInt(-2147482623)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(0xfffffffe)").IsEqualToInt(-2)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(0xffffffff)").IsEqualToInt(-1)
}

func TestUtilReinterpretUint64AsInt64(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x0000000000000000)").IsEqualToInt(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x7fffffffffffffff)").IsEqualToInt(9223372036854775807)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x8000000000000000)").IsEqualToInt(-9223372036854775808)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x8000000000000400)").ToSucceed()
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0x8000000000000401)").ToSucceed()
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0xfffffffffffffffe)").IsEqualToInt(-2)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(0xffffffffffffffff)").IsEqualToInt(-1)
}

func TestUtilReinterpretUint64AsInt32(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int32(0xffffffff)").IsEqualToInt(-1)
}

func TestUtilReinterpretUint64AsUint32(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_uint32(0xffffffff)").IsEqualToUint(0xffffffff)
}

func TestUtilReinterpretUint64AsSint64(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(0)").IsEqualToInt(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(1)").IsEqualToInt(-1)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(2)").IsEqualToInt(1)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(3)").IsEqualToInt(-2)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(0xfffffffe)").IsEqualToInt(0x7fffffff)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(0xffffffff)").IsEqualToInt(-0x80000000)
}

func TestUtilBinAsUint64(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_bin_as_uint64(_binary X'00')").IsEqualToUint(0)
	RunTestThatExpression(t, "_pb_util_bin_as_uint64(_binary X'7fffffffffffffff')").IsEqualToUint(9223372036854775807)
	RunTestThatExpression(t, "_pb_util_bin_as_uint64(_binary X'8000000000000000')").IsEqualToUint(9223372036854775808)
	RunTestThatExpression(t, "_pb_util_bin_as_uint64(_binary X'ffffffffffffffff')").IsEqualToUint(18446744073709551615)
}

func TestUtilBinAsInt64(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_bin_as_int64(_binary X'00')").IsEqualToInt(0)
	RunTestThatExpression(t, "_pb_util_bin_as_int64(_binary X'7fffffffffffffff')").IsEqualToInt(9223372036854775807)
	RunTestThatExpression(t, "_pb_util_bin_as_int64(_binary X'8000000000000000')").IsEqualToInt(-9223372036854775808)
	RunTestThatExpression(t, "_pb_util_bin_as_int64(_binary X'ffffffffffffffff')").IsEqualToInt(-1)
}

func TestUtilBinAsInt32(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_bin_as_int32(_binary X'00')").IsEqualToInt(0)
	RunTestThatExpression(t, "_pb_util_bin_as_int32(_binary X'7fffffff')").IsEqualToInt(2147483647)
	RunTestThatExpression(t, "_pb_util_bin_as_int32(_binary X'80000000')").IsEqualToInt(-2147483648)
	RunTestThatExpression(t, "_pb_util_bin_as_int32(_binary X'ffffffff')").IsEqualToInt(-1)
}

func TestUtilBinAsUint32(t *testing.T) {
	RunTestThatExpression(t, "_pb_util_bin_as_uint32(_binary X'00')").IsEqualToUint(0)
	RunTestThatExpression(t, "_pb_util_bin_as_uint32(_binary X'7fffffff')").IsEqualToUint(2147483647)
	RunTestThatExpression(t, "_pb_util_bin_as_uint32(_binary X'80000000')").IsEqualToUint(2147483648)
	RunTestThatExpression(t, "_pb_util_bin_as_uint32(_binary X'ffffffff')").IsEqualToUint(4294967295)
}
