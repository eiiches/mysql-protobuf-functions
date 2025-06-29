package gfloat_test

import (
	"math"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/gomega/gfloat"
	. "github.com/onsi/gomega"
)

func TestZeroMatchers(t *testing.T) {
	g := NewWithT(t)

	// Test positive zero (float64)
	positiveZero := 0.0
	g.Expect(positiveZero).To(gfloat.BePositiveZero())
	g.Expect(positiveZero).NotTo(gfloat.BeNegativeZero())

	// Test negative zero (float64)
	negativeZero := math.Copysign(0, -1)
	g.Expect(negativeZero).To(gfloat.BeNegativeZero())
	g.Expect(negativeZero).NotTo(gfloat.BePositiveZero())

	// Test positive zero (float32)
	positiveZeroF32 := float32(0.0)
	g.Expect(positiveZeroF32).To(gfloat.BePositiveZero())
	g.Expect(positiveZeroF32).NotTo(gfloat.BeNegativeZero())

	// Test negative zero (float32)
	negativeZeroF32 := float32(math.Copysign(0, -1))
	g.Expect(negativeZeroF32).To(gfloat.BeNegativeZero())
	g.Expect(negativeZeroF32).NotTo(gfloat.BePositiveZero())

	// Test non-zero values
	g.Expect(1.0).NotTo(gfloat.BePositiveZero())
	g.Expect(1.0).NotTo(gfloat.BeNegativeZero())
	g.Expect(-1.0).NotTo(gfloat.BePositiveZero())
	g.Expect(-1.0).NotTo(gfloat.BeNegativeZero())

	// Test float32 non-zero values
	g.Expect(float32(1.0)).NotTo(gfloat.BePositiveZero())
	g.Expect(float32(1.0)).NotTo(gfloat.BeNegativeZero())
	g.Expect(float32(-1.0)).NotTo(gfloat.BePositiveZero())
	g.Expect(float32(-1.0)).NotTo(gfloat.BeNegativeZero())
}
