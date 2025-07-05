package morefloat

import "math"

// equalOrCloseWhenWithinMantissaThreshold compares two float values by checking if their mantissa difference
// is within the acceptable range of 2^ignoreBits.
// The floatSize parameter (32 or 64) indicates the original precision of the values.
func equalOrCloseWhenWithinMantissaThreshold(a, b float64, floatSize int, ignoreBits int) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 0) && math.IsInf(b, 0) {
		return math.Signbit(a) == math.Signbit(b)
	}
	if a == b {
		return true
	}

	if floatSize == 32 {
		// Convert to float32
		a32 := float32(a)
		b32 := float32(b)
		aBits := math.Float32bits(a32)
		bBits := math.Float32bits(b32)

		// Extract sign and exponent
		aSign := aBits >> 31
		bSign := bBits >> 31
		aExp := (aBits >> 23) & 0xFF
		bExp := (bBits >> 23) & 0xFF

		// Signs and exponents must match exactly
		if aSign != bSign || aExp != bExp {
			return false
		}

		// Extract mantissas (including implicit bit for normalized numbers)
		var aMantissa, bMantissa uint32
		if aExp == 0 {
			// Subnormal number
			aMantissa = aBits & 0x7FFFFF
			bMantissa = bBits & 0x7FFFFF
		} else {
			// Normal number - add implicit bit
			aMantissa = (aBits & 0x7FFFFF) | 0x800000
			bMantissa = (bBits & 0x7FFFFF) | 0x800000
		}

		// Check if mantissa difference is within acceptable range
		var diff uint32
		if aMantissa > bMantissa {
			diff = aMantissa - bMantissa
		} else {
			diff = bMantissa - aMantissa
		}

		return diff < (uint32(1) << ignoreBits)
	} else {
		// Work with float64
		aBits := math.Float64bits(a)
		bBits := math.Float64bits(b)

		// Extract sign and exponent
		aSign := aBits >> 63
		bSign := bBits >> 63
		aExp := (aBits >> 52) & 0x7FF
		bExp := (bBits >> 52) & 0x7FF

		// Signs and exponents must match exactly
		if aSign != bSign || aExp != bExp {
			return false
		}

		// Extract mantissas (including implicit bit for normalized numbers)
		var aMantissa, bMantissa uint64
		if aExp == 0 {
			// Subnormal number
			aMantissa = aBits & 0xFFFFFFFFFFFFF
			bMantissa = bBits & 0xFFFFFFFFFFFFF
		} else {
			// Normal number - add implicit bit
			aMantissa = (aBits & 0xFFFFFFFFFFFFF) | 0x10000000000000
			bMantissa = (bBits & 0xFFFFFFFFFFFFF) | 0x10000000000000
		}

		// Check if mantissa difference is within acceptable range
		var diff uint64
		if aMantissa > bMantissa {
			diff = aMantissa - bMantissa
		} else {
			diff = bMantissa - aMantissa
		}

		return diff < (uint64(1) << ignoreBits)
	}
}

type EqualOrCloseFn func(a, b float64, floatSize int) bool

// WithinMantissaThreshold creates an equalOrCloseFn that considers two float32 values to be equal or close,
// if the mantissa difference of the two is less than 2**float32MantissaThresholdInBits, and
// 2**float64MantissaThresholdInBits for float64 values.
func WithinMantissaThreshold(float32MantissaThresholdInBits, float64MantissaThresholdInBits int) EqualOrCloseFn {
	if float32MantissaThresholdInBits < 0 || float32MantissaThresholdInBits > 23 {
		panic("float32MantissaThresholdInBits must be between 0 and 23")
	}
	if float64MantissaThresholdInBits < 0 || float64MantissaThresholdInBits > 52 {
		panic("float64MantissaThresholdInBits must be between 0 and 52")
	}
	return func(a, b float64, floatSize int) bool {
		if floatSize == 32 {
			return equalOrCloseWhenWithinMantissaThreshold(a, b, floatSize, float32MantissaThresholdInBits)
		} else {
			return equalOrCloseWhenWithinMantissaThreshold(a, b, floatSize, float64MantissaThresholdInBits)
		}
	}
}
