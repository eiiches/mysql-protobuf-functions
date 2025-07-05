package morefloat

import (
	"math"
	"testing"

	. "github.com/onsi/gomega"
)

func TestWithinMantissaThreshold(t *testing.T) {
	g := NewWithT(t)
	g.Expect(WithinMantissaThreshold(0, 0)(1.0, 1.0, 64)).To(BeTrue())
	g.Expect(WithinMantissaThreshold(0, 0)(1.0, 2.0, 64)).To(BeFalse())
	g.Expect(WithinMantissaThreshold(0, 0)(math.NaN(), math.NaN(), 64)).To(BeTrue())
	g.Expect(WithinMantissaThreshold(0, 0)(math.Inf(1), math.Inf(1), 64)).To(BeTrue())
	g.Expect(WithinMantissaThreshold(0, 0)(math.Inf(-1), math.Inf(-1), 64)).To(BeTrue())
	g.Expect(WithinMantissaThreshold(0, 0)(math.Inf(-1), math.Inf(1), 64)).To(BeFalse())
	g.Expect(WithinMantissaThreshold(0, 0)(1.0, 1.000000000000001, 64)).To(BeFalse())
	g.Expect(WithinMantissaThreshold(0, 3)(1.0, 1.000000000000001, 64)).To(BeTrue())
	g.Expect(WithinMantissaThreshold(0, 0)(1.0, 1.0000001, 32)).To(BeFalse())
	g.Expect(WithinMantissaThreshold(1, 0)(1.0, 1.0000001, 32)).To(BeTrue())
}
