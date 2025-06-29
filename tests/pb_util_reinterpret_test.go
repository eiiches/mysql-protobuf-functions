package main

import (
	"math"
	"testing"
)

func TestNegativeZeroRoundTrip(t *testing.T) {
	// Test negative zero round-trip conversion
	negativeZero := math.Copysign(0, -1)

	// Test float round-trip
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", negativeZero).IsEqualToUint(0x80000000)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(0x80000000)").IsNegativeZero()

	// Test double round-trip
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", negativeZero).IsEqualToUint(0x8000000000000000)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(0x8000000000000000)").IsNegativeZero()
}

func TestUtilReinterpretInt64AsUint64(t *testing.T) {
	// Test 2's complement conversion from signed to unsigned 64-bit

	// Positive values (direct mapping)
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(0)").IsEqualToUint(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(1)").IsEqualToUint(1)
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(42)").IsEqualToUint(42)
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(9223372036854775807)").IsEqualToUint(9223372036854775807) // max int64

	// Negative values (2's complement)
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(?)", -1).IsEqualToUint(18446744073709551615)                  // 0xFFFFFFFFFFFFFFFF
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(?)", -2).IsEqualToUint(18446744073709551614)                  // 0xFFFFFFFFFFFFFFFE
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(?)", -42).IsEqualToUint(18446744073709551574)                 // 0xFFFFFFFFFFFFFFD6
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(?)", -9223372036854775808).IsEqualToUint(9223372036854775808) // min int64 -> 0x8000000000000000

	// Edge cases around zero
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(?)", -0).IsEqualToUint(0)

	// Boundary values
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(9223372036854775806)").IsEqualToUint(9223372036854775806)     // max int64 - 1
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(?)", -9223372036854775807).IsEqualToUint(9223372036854775809) // min int64 + 1
}

func TestUtilReinterpretInt32AsUint32(t *testing.T) {
	// Test 2's complement conversion from signed to unsigned 32-bit

	// Positive values (direct mapping)
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(0)").IsEqualToUint(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(1)").IsEqualToUint(1)
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(42)").IsEqualToUint(42)
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(2147483647)").IsEqualToUint(2147483647) // max int32

	// Negative values (2's complement using modular arithmetic)
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(?)", -1).IsEqualToUint(4294967295)          // 0xFFFFFFFF
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(?)", -2).IsEqualToUint(4294967294)          // 0xFFFFFFFE
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(?)", -42).IsEqualToUint(4294967254)         // 0xFFFFFFD6
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(?)", -2147483648).IsEqualToUint(2147483648) // min int32 -> 0x80000000

	// Edge cases
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(?)", -0).IsEqualToUint(0)

	// Boundary values
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(2147483646)").IsEqualToUint(2147483646)     // max int32 - 1
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(?)", -2147483647).IsEqualToUint(2147483649) // min int32 + 1

	// Values used in setter tests
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(?)", -42).IsEqualToUint(4294967254) // Should match sfixed32 test
}

func TestUtilReinterpretSint64AsUint64(t *testing.T) {
	// Test ZigZag encoding: https://developers.google.com/protocol-buffers/docs/encoding#signed-integers
	// Formula: (n << 1) ^ (n >> 63) for 64-bit
	// Positive: n -> 2*n, Negative: n -> 2*|n|-1

	// Test zero
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(0)").IsEqualToUint(0)

	// Test positive values: n -> 2*n
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(1)").IsEqualToUint(2)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(2)").IsEqualToUint(4)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(3)").IsEqualToUint(6)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(150)").IsEqualToUint(300)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(1000)").IsEqualToUint(2000)

	// Test negative values: -n -> 2*n-1
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -1).IsEqualToUint(1)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -2).IsEqualToUint(3)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -3).IsEqualToUint(5)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -150).IsEqualToUint(299)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -1000).IsEqualToUint(1999)

	// Test standard ZigZag test vectors
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(0)").IsEqualToUint(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -1).IsEqualToUint(1)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(1)").IsEqualToUint(2)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -2).IsEqualToUint(3)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(2)").IsEqualToUint(4)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -3).IsEqualToUint(5)

	// Test 32-bit boundary values
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(2147483647)").IsEqualToUint(4294967294)     // max int32 -> 2*2147483647
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -2147483648).IsEqualToUint(4294967295) // min int32 -> 2*2147483648-1

	// Test values used in setter tests
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -1).IsEqualToUint(1) // Should match sint32/sint64 tests
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -2).IsEqualToUint(3) // Should match repeated tests

	// Test larger values
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(1000000)").IsEqualToUint(2000000)
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -1000000).IsEqualToUint(1999999)
}

func TestUtilReinterpretFloatAsUint32(t *testing.T) {
	// Test IEEE 754 single-precision floating point conversion
	// https://en.wikipedia.org/wiki/Single-precision_floating-point_format

	// Test zero values
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(0.0)").IsEqualToUint(0x00000000)                     // +0.0
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", math.Copysign(0, -1)).IsEqualToUint(0x80000000) // -0.0

	// Test basic integer values
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(1.0)").IsEqualToUint(0x3F800000)     // 1.0
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(2.0)").IsEqualToUint(0x40000000)     // 2.0
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(4.0)").IsEqualToUint(0x40800000)     // 4.0
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", -1.0).IsEqualToUint(0xBF800000) // -1.0
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", -2.0).IsEqualToUint(0xC0000000) // -2.0

	// Test fractional values
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(0.5)").IsEqualToUint(0x3F000000)     // 0.5
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(1.5)").IsEqualToUint(0x3FC00000)     // 1.5 (0x3FC00000 = 1069547520)
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(2.5)").IsEqualToUint(0x40200000)     // 2.5 (0x40200000 = 1075838976)
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(3.5)").IsEqualToUint(0x40600000)     // 3.5
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", -1.5).IsEqualToUint(0xBFC00000) // -1.5
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", -2.5).IsEqualToUint(0xC0200000) // -2.5

	// Test decimal values
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(0.25)").IsEqualToUint(0x3E800000)  // 0.25
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(0.75)").IsEqualToUint(0x3F400000)  // 0.75
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(0.125)").IsEqualToUint(0x3E000000) // 0.125

	// Test values from setter tests (exact matches required)
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(1.5)").IsEqualToUint(1069547520) // Should match setter test
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(2.5)").IsEqualToUint(1075838976) // Should match repeated test

	// Test edge values
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(0.9999999404)").IsEqualToUint(0x3F7FFFFF) // Largest number less than one
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(1.0000001192)").IsEqualToUint(0x3F800001) // Smallest number larger than one

	// Test larger values
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(10.0)").IsEqualToUint(0x41200000)      // 10.0
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(100.0)").IsEqualToUint(0x42C80000)     // 100.0
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", -10.0).IsEqualToUint(0xC1200000)  // -10.0
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", -100.0).IsEqualToUint(0xC2C80000) // -100.0
}

func TestUtilReinterpretDoubleAsUint64(t *testing.T) {
	// Test IEEE 754 double-precision floating point conversion
	// https://en.wikipedia.org/wiki/Double-precision_floating-point_format

	// Test zero values
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(0.0)").IsEqualToUint(0x0000000000000000)                     // +0.0
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", math.Copysign(0, -1)).IsEqualToUint(0x8000000000000000) // -0.0

	// Test basic integer values
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(1.0)").IsEqualToUint(0x3FF0000000000000)     // 1.0 (4607182418800017408)
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(2.0)").IsEqualToUint(0x4000000000000000)     // 2.0 (4611686018427387904)
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(4.0)").IsEqualToUint(0x4010000000000000)     // 4.0
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", -1.0).IsEqualToUint(0xBFF0000000000000) // -1.0 (13830554455654793216)
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", -2.0).IsEqualToUint(0xC000000000000000) // -2.0

	// Test fractional values
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(0.5)").IsEqualToUint(0x3FE0000000000000)     // 0.5 (4602678819172646912)
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(1.5)").IsEqualToUint(0x3FF8000000000000)     // 1.5 (4609434218613702656)
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(2.5)").IsEqualToUint(0x4004000000000000)     // 2.5 (4612811918334230528)
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(3.5)").IsEqualToUint(0x400C000000000000)     // 3.5
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", -1.5).IsEqualToUint(0xBFF8000000000000) // -1.5 (13832806255468478464)
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", -2.5).IsEqualToUint(0xC004000000000000) // -2.5

	// Test decimal values
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(0.25)").IsEqualToUint(0x3FD0000000000000)  // 0.25 (4598175219545276416)
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(0.75)").IsEqualToUint(0x3FE8000000000000)  // 0.75
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(0.125)").IsEqualToUint(0x3FC0000000000000) // 0.125

	// Test values from setter tests (exact matches required)
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(1.5)").IsEqualToUint(4609434218613702656) // Should match setter test
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(2.5)").IsEqualToUint(4612811918334230528) // Should match repeated test

	// Test larger values
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(10.0)").IsEqualToUint(0x4024000000000000)      // 10.0
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(100.0)").IsEqualToUint(0x4059000000000000)     // 100.0
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", -10.0).IsEqualToUint(0xC024000000000000)  // -10.0
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", -100.0).IsEqualToUint(0xC059000000000000) // -100.0

	// Test very small values
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(0.001)").IsEqualToUint(0x3F50624DD2F1A9FC)  // 0.001
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(0.0001)").IsEqualToUint(0x3F1A36E2EB1C432D) // 0.0001
}

func TestUtilReinterpretRoundTrip(t *testing.T) {
	// Test round-trip conversions with existing decode functions

	// Test int64 round-trip
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(_pb_util_reinterpret_int64_as_uint64(0))").IsEqualToInt(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(_pb_util_reinterpret_int64_as_uint64(42))").IsEqualToInt(42)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(_pb_util_reinterpret_int64_as_uint64(?))", -42).IsEqualToInt(-42)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(_pb_util_reinterpret_int64_as_uint64(9223372036854775807))").IsEqualToInt(9223372036854775807)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_int64(_pb_util_reinterpret_int64_as_uint64(?))", -9223372036854775808).IsEqualToInt(-9223372036854775808)

	// Test int32 round-trip
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(_pb_util_reinterpret_int32_as_uint32(0))").IsEqualToInt(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(_pb_util_reinterpret_int32_as_uint32(42))").IsEqualToInt(42)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(_pb_util_reinterpret_int32_as_uint32(?))", -42).IsEqualToInt(-42)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(_pb_util_reinterpret_int32_as_uint32(2147483647))").IsEqualToInt(2147483647)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_int32(_pb_util_reinterpret_int32_as_uint32(?))", -2147483648).IsEqualToInt(-2147483648)

	// Test sint64 round-trip (ZigZag)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(_pb_util_reinterpret_sint64_as_uint64(0))").IsEqualToInt(0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(_pb_util_reinterpret_sint64_as_uint64(42))").IsEqualToInt(42)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(_pb_util_reinterpret_sint64_as_uint64(?))", -42).IsEqualToInt(-42)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(_pb_util_reinterpret_sint64_as_uint64(1))").IsEqualToInt(1)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(_pb_util_reinterpret_sint64_as_uint64(?))", -1).IsEqualToInt(-1)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(_pb_util_reinterpret_sint64_as_uint64(2147483647))").IsEqualToInt(2147483647)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_sint64(_pb_util_reinterpret_sint64_as_uint64(?))", -2147483648).IsEqualToInt(-2147483648)

	// Test float round-trip
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(_pb_util_reinterpret_float_as_uint32(0.0))").IsEqualToFloat(0.0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(_pb_util_reinterpret_float_as_uint32(?))", math.Copysign(0, -1)).IsNegativeZero()
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(_pb_util_reinterpret_float_as_uint32(1.0))").IsEqualToFloat(1.0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(_pb_util_reinterpret_float_as_uint32(1.5))").IsEqualToFloat(1.5)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(_pb_util_reinterpret_float_as_uint32(2.5))").IsEqualToFloat(2.5)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(_pb_util_reinterpret_float_as_uint32(?))", -1.5).IsEqualToFloat(-1.5)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(_pb_util_reinterpret_float_as_uint32(?))", -2.5).IsEqualToFloat(-2.5)

	// Test double round-trip
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(_pb_util_reinterpret_double_as_uint64(0.0))").IsEqualToFloat(0.0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(_pb_util_reinterpret_double_as_uint64(?))", math.Copysign(0, -1)).IsNegativeZero()
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(_pb_util_reinterpret_double_as_uint64(1.0))").IsEqualToFloat(1.0)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(_pb_util_reinterpret_double_as_uint64(1.5))").IsEqualToFloat(1.5)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(_pb_util_reinterpret_double_as_uint64(2.5))").IsEqualToFloat(2.5)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(_pb_util_reinterpret_double_as_uint64(?))", -1.5).IsEqualToFloat(-1.5)
	RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(_pb_util_reinterpret_double_as_uint64(?))", -2.5).IsEqualToFloat(-2.5)
}

func TestUtilReinterpretEdgeCases(t *testing.T) {
	// Test edge cases and boundary values

	// Test int32 boundaries
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(2147483647)").IsEqualToUint(2147483647)     // max int32
	RunTestThatExpression(t, "_pb_util_reinterpret_int32_as_uint32(?)", -2147483648).IsEqualToUint(2147483648) // min int32

	// Test int64 boundaries
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(9223372036854775807)").IsEqualToUint(9223372036854775807)     // max int64
	RunTestThatExpression(t, "_pb_util_reinterpret_int64_as_uint64(?)", -9223372036854775808).IsEqualToUint(9223372036854775808) // min int64

	// Test ZigZag edge cases for large values
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(9223372036854775807)").IsEqualToUint(18446744073709551614)     // max int64 ZigZag encoded
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -9223372036854775808).IsEqualToUint(18446744073709551615) // min int64 ZigZag encoded

	// Test that our implementations produce the exact values expected by setter tests
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -1).IsEqualToUint(1)                 // Must match setter test expectation
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -2).IsEqualToUint(3)                 // Must match repeated test expectation
	RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(1.5)").IsEqualToUint(1069547520)           // Must match setter test expectation
	RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(1.5)").IsEqualToUint(4609434218613702656) // Must match setter test expectation

	// Test intermediate values for sint32/sint64 compatibility
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(1073741823)").IsEqualToUint(2147483646)     // Large positive
	RunTestThatExpression(t, "_pb_util_reinterpret_sint64_as_uint64(?)", -1073741824).IsEqualToUint(2147483647) // Large negative
}
