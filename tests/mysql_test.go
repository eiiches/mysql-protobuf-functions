package main

import (
	"math"
	"testing"
)

func TestSignFunctionWithNegativeZero(t *testing.T) {
	// Test SIGN function with negative zero from Go
	// This test documents MySQL's behavior where SIGN() returns 0 for both +0.0 and -0.0
	RunTestThatExpression(t, "SIGN(?)", math.Copysign(0, -1)).IsEqualToInt(0)
	RunTestThatExpression(t, "SIGN(?)", 0.0).IsEqualToInt(0)
	RunTestThatExpression(t, "SIGN(?)", -1.0).IsEqualToInt(-1)
	RunTestThatExpression(t, "SIGN(?)", 1.0).IsEqualToInt(1)
}
