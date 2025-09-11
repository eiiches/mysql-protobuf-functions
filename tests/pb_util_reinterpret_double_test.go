package main

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestUtilReinterpretUint64AsDoubleComprehensive(t *testing.T) {
	// Comprehensive test of _pb_util_reinterpret_uint64_as_double
	// Tests every combination of sign bit and exponent with mantissa all 0s and all 1s

	// Test each exponent [0, 2047) with both sign bits and both mantissa extremes
	for exponent := 0; exponent < 2047; exponent++ {
		for signBit := 0; signBit <= 1; signBit++ {
			// Test with mantissa = 0 (all zeros)
			bitsZeros := uint64(signBit<<63) | uint64(exponent<<52) | 0x0000000000000
			bitsZeroValue := math.Float64frombits(bitsZeros)
			RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(?)", bitsZeros).IsEqualToDouble(bitsZeroValue)

			// Test with mantissa = 0xFFFFFFFFFFFFF (all ones)
			bitsOnes := uint64(signBit<<63) | uint64(exponent<<52) | 0xFFFFFFFFFFFFF
			bitsOnesValue := math.Float64frombits(bitsOnes)
			RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(?)", bitsOnes).IsEqualToDouble(bitsOnesValue)
		}
	}
}

func TestUtilReinterpretDoubleAsUint64Comprehensive(t *testing.T) {
	// Comprehensive test of _pb_util_reinterpret_double_as_uint64
	// Tests conversion from known double values to expected bit patterns

	// Test each exponent [0, 2047) with both sign bits and both mantissa extremes
	// Skip exponent 2047 (infinity/NaN) as they may not convert predictably
	for exponent := 0; exponent < 2047; exponent++ {
		for signBit := 0; signBit <= 1; signBit++ {
			// Test with mantissa = 0 (all zeros)
			bitsZeros := uint64(signBit<<63) | uint64(exponent<<52) | 0x0000000000000
			bitsZeroValue := math.Float64frombits(bitsZeros)
			RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", bitsZeroValue).IsEqualToUint(bitsZeros)

			// Test with mantissa = 0xFFFFFFFFFFFFF (all ones)
			bitsOnes := uint64(signBit<<63) | uint64(exponent<<52) | 0xFFFFFFFFFFFFF
			bitsOnesValue := math.Float64frombits(bitsOnes)
			RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", bitsOnesValue).IsEqualToUint(bitsOnes)
		}
	}
}

func TestUtilReinterpretUint64AsDoubleFuzzingSubnormal(t *testing.T) {
	// Fuzzing test for _pb_util_reinterpret_uint64_as_double - subnormal range
	// Tests subnormal values (exponent = 0, mantissa != 0)

	seed := time.Now().UnixNano()
	t.Logf("Using random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 5000; i++ {
		// Generate subnormal values: exponent = 0, random mantissa
		mantissa := 1 + (rng.Uint64() & 0xFFFFFFFFFFFFE) // Ensure non-zero mantissa
		signBit := rng.Uint64() & 1
		bits := (signBit << 63) | mantissa // Exponent is 0

		expectedDouble := math.Float64frombits(bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(?)", bits).IsEqualToDouble(expectedDouble)
	}
}

func TestUtilReinterpretUint64AsDoubleFuzzingNormal(t *testing.T) {
	// Fuzzing test for _pb_util_reinterpret_uint64_as_double - normal range
	// Tests normal values (1 <= exponent <= 2046)

	seed := time.Now().UnixNano()
	t.Logf("Using random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 5000; i++ {
		// Generate normal values: 1 <= exponent <= 2046, random mantissa
		exponent := 1 + (rng.Uint64() % 2046) // Exponent range 1-2046
		mantissa := rng.Uint64() & 0xFFFFFFFFFFFFF
		signBit := rng.Uint64() & 1
		bits := (signBit << 63) | (exponent << 52) | mantissa

		expectedDouble := math.Float64frombits(bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(?)", bits).IsEqualToDouble(expectedDouble)
	}
}

func TestUtilReinterpretDoubleAsUint64FuzzingSubnormal(t *testing.T) {
	// Fuzzing test for _pb_util_reinterpret_double_as_uint64 - subnormal range
	// Tests subnormal values (exponent = 0, mantissa != 0)

	seed := time.Now().UnixNano()
	t.Logf("Using random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 5000; i++ {
		// Generate subnormal values: exponent = 0, random mantissa
		mantissa := 1 + (rng.Uint64() & 0xFFFFFFFFFFFFE) // Ensure non-zero mantissa
		signBit := rng.Uint64() & 1
		bits := (signBit << 63) | mantissa // Exponent is 0

		doubleVal := math.Float64frombits(bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", doubleVal).IsEqualToUint(bits)
	}
}

func TestUtilReinterpretDoubleAsUint64FuzzingNormal(t *testing.T) {
	// Fuzzing test for _pb_util_reinterpret_double_as_uint64 - normal range
	// Tests normal values (1 <= exponent <= 2046)

	seed := time.Now().UnixNano()
	t.Logf("Using random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 5000; i++ {
		// Generate normal values: 1 <= exponent <= 2046, random mantissa
		exponent := 1 + (rng.Uint64() % 2046) // Exponent range 1-2046
		mantissa := rng.Uint64() & 0xFFFFFFFFFFFFF
		signBit := rng.Uint64() & 1
		bits := (signBit << 63) | (exponent << 52) | mantissa

		doubleVal := math.Float64frombits(bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", doubleVal).IsEqualToUint(bits)
	}
}

func TestDoubleIEEE754SpecialCases(t *testing.T) {
	// Test specific IEEE 754 edge cases for double precision
	// Use Go's math.Float64frombits as the canonical implementation

	testCases := []struct {
		bits    uint64
		comment string
	}{
		// Subnormal/normal boundary cases
		{0x0010000000000000, "Smallest positive normal"},
		{0x000FFFFFFFFFFFFF, "Largest positive subnormal"},
		{0x0010000000000001, "Smallest normal + 1 ULP"},
		{0x000FFFFFFFFFFFFE, "Largest subnormal - 1 ULP"},
		{0x8010000000000000, "Smallest negative normal"},
		{0x800FFFFFFFFFFFFF, "Largest negative subnormal"},
		{0x8010000000000001, "Smallest negative normal + 1 ULP"},
		{0x800FFFFFFFFFFFFE, "Largest negative subnormal - 1 ULP"},

		// Zero and near-zero values
		{0x0000000000000000, "+0.0"},
		{0x8000000000000000, "-0.0"},
		{0x0000000000000001, "Smallest positive subnormal"},
		{0x8000000000000001, "Smallest negative subnormal"},

		// TODO: Infinity values
		// {0x7FF0000000000000, "+Infinity"},
		// {0xFFF0000000000000, "-Infinity"},

		// Maximum finite values
		{0x7FEFFFFFFFFFFFFF, "Largest positive finite"},
		{0xFFEFFFFFFFFFFFFF, "Largest negative finite"},
	}

	for _, tc := range testCases {
		// Test uint64 → double conversion
		expectedDouble := math.Float64frombits(tc.bits)
		RunTestThatExpression(t, "_pb_util_reinterpret_uint64_as_double(?)", tc.bits).IsEqualToDouble(expectedDouble)

		// Test double → uint64 conversion
		RunTestThatExpression(t, "_pb_util_reinterpret_double_as_uint64(?)", expectedDouble).IsEqualToUint(tc.bits)
	}
}
