package main

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestUtilReinterpretUint32AsFloatComprehensive(t *testing.T) {
	// Comprehensive test of _pb_util_reinterpret_uint32_as_float
	// Tests every combination of sign bit and exponent with mantissa all 0s and all 1s

	// Test each exponent [0, 255) with both sign bits and both mantissa extremes
	for exponent := 0; exponent < 255; exponent++ {
		for signBit := 0; signBit <= 1; signBit++ {
			// Test with mantissa = 0 (all zeros)
			bitsZeros := uint32(signBit<<31) | uint32(exponent<<23) | 0x000000
			bitsZeroValue := math.Float32frombits(bitsZeros)
			RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(?)", bitsZeros).IsEqualToFloat(bitsZeroValue)

			// Test with mantissa = 0x7FFFFF (all ones)
			bitsOnes := uint32(signBit<<31) | uint32(exponent<<23) | 0x7FFFFF
			bitsOnesValue := math.Float32frombits(bitsOnes)
			RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(?)", bitsOnes).IsEqualToFloat(bitsOnesValue)
		}
	}
}

func TestUtilReinterpretFloatAsUint32Comprehensive(t *testing.T) {
	// Comprehensive test of _pb_util_reinterpret_float_as_uint32
	// Tests conversion from known float values to expected bit patterns

	// Test each exponent [0, 255) with both sign bits and both mantissa extremes
	// Skip exponent 255 (infinity/NaN) as they may not convert predictably
	for exponent := 0; exponent < 255; exponent++ {
		for signBit := 0; signBit <= 1; signBit++ {
			// Test with mantissa = 0 (all zeros)
			bitsZeros := uint32(signBit<<31) | uint32(exponent<<23) | 0x000000
			bitsZeroValue := math.Float32frombits(bitsZeros)
			RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", bitsZeroValue).IsEqualToUint(uint64(bitsZeros))

			// Test with mantissa = 0x7FFFFF (all ones)
			bitsOnes := uint32(signBit<<31) | uint32(exponent<<23) | 0x7FFFFF
			bitsOnesValue := math.Float32frombits(bitsOnes)
			RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", bitsOnesValue).IsEqualToUint(uint64(bitsOnes))
		}
	}
}

func TestUtilReinterpretUint32AsFloatFuzzingSubnormal(t *testing.T) {
	// Fuzzing test for _pb_util_reinterpret_uint32_as_float - subnormal range
	// Tests subnormal values (exponent = 0, mantissa != 0)

	seed := time.Now().UnixNano()
	t.Logf("Using random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 5000; i++ {
		// Generate subnormal values: exponent = 0, random mantissa
		mantissa := 1 + (rng.Uint32() & 0x7FFFFE) // Ensure non-zero mantissa
		signBit := rng.Uint32() & 1
		bits := (signBit << 31) | mantissa // Exponent is 0

		expectedFloat := math.Float32frombits(bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(?)", bits).IsEqualToFloat(expectedFloat)
	}
}

func TestUtilReinterpretUint32AsFloatFuzzingNormal(t *testing.T) {
	// Fuzzing test for _pb_util_reinterpret_uint32_as_float - normal range
	// Tests normal values (1 <= exponent <= 254)

	seed := time.Now().UnixNano()
	t.Logf("Using random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 5000; i++ {
		// Generate normal values: 1 <= exponent <= 254, random mantissa
		exponent := 1 + (rng.Uint32() % 254) // Exponent range 1-254
		mantissa := rng.Uint32() & 0x7FFFFF
		signBit := rng.Uint32() & 1
		bits := (signBit << 31) | (exponent << 23) | mantissa

		expectedFloat := math.Float32frombits(bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(?)", bits).IsEqualToFloat(expectedFloat)
	}
}

func TestUtilReinterpretFloatAsUint32FuzzingSubnormal(t *testing.T) {
	// Fuzzing test for _pb_util_reinterpret_float_as_uint32 - subnormal range
	// Tests subnormal values (exponent = 0, mantissa != 0)

	seed := time.Now().UnixNano()
	t.Logf("Using random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 5000; i++ {
		// Generate subnormal values: exponent = 0, random mantissa
		mantissa := 1 + (rng.Uint32() & 0x7FFFFE) // Ensure non-zero mantissa
		signBit := rng.Uint32() & 1
		bits := (signBit << 31) | mantissa // Exponent is 0

		floatVal := math.Float32frombits(bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", floatVal).IsEqualToUint(uint64(bits))
	}
}

func TestUtilReinterpretFloatAsUint32FuzzingNormal(t *testing.T) {
	// Fuzzing test for _pb_util_reinterpret_float_as_uint32 - normal range
	// Tests normal values (1 <= exponent <= 254)

	seed := time.Now().UnixNano()
	t.Logf("Using random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 5000; i++ {
		// Generate normal values: 1 <= exponent <= 254, random mantissa
		exponent := 1 + (rng.Uint32() % 254) // Exponent range 1-254
		mantissa := rng.Uint32() & 0x7FFFFF
		signBit := rng.Uint32() & 1
		bits := (signBit << 31) | (exponent << 23) | mantissa

		floatVal := math.Float32frombits(bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", floatVal).IsEqualToUint(uint64(bits))
	}
}

func TestFloatIEEE754SpecialCases(t *testing.T) {
	// Test specific IEEE 754 edge cases that caused issues
	// Use Go's math.Float32frombits as the canonical implementation

	testCases := []struct {
		bits    uint32
		comment string
	}{
		// Subnormal/normal boundary cases
		{0x00800000, "Smallest positive normal"},
		{0x007FFFFF, "Largest positive subnormal"},
		{0x00800001, "Smallest normal + 1 ULP"},
		{0x007FFFFE, "Largest subnormal - 1 ULP"},
		{0x80800000, "Smallest negative normal"},
		{0x807FFFFF, "Largest negative subnormal"},
		{0x80800001, "Smallest negative normal + 1 ULP"},
		{0x807FFFFE, "Largest negative subnormal - 1 ULP"},

		// Zero and near-zero values
		{0x00000000, "+0.0"},
		{0x80000000, "-0.0"},
		{0x00000001, "Smallest positive subnormal"},
		{0x80000001, "Smallest negative subnormal"},

		// TODO: Infinity values
		// {0x7F800000, "+Infinity"},
		// {0xFF800000, "-Infinity"},

		// Maximum finite values
		{0x7F7FFFFF, "Largest positive finite"},
		{0xFF7FFFFF, "Largest negative finite"},
	}

	for _, tc := range testCases {
		// Test uint32 → float conversion
		expectedFloat := math.Float32frombits(tc.bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_uint32_as_float(?)", tc.bits).IsEqualToFloat(expectedFloat)

		// Test float → uint32 conversion
		RunTestThatExpression(t, "_pb_util_reinterpret_float_as_uint32(?)", expectedFloat).IsEqualToUint(uint64(tc.bits))
	}
}
