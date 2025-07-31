package moreproto_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/morefloat"
	"github.com/eiiches/mysql-protobuf-functions/internal/moreproto"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestEqualOrClose(t *testing.T) {
	g := NewWithT(t)

	exactComparison := morefloat.WithinMantissaThreshold(0, 0)

	// Equal messages should be equal
	g.Expect(moreproto.EqualOrClose(wrapperspb.Double(1.23), wrapperspb.Double(1.23), exactComparison)).To(BeTrue())
	g.Expect(moreproto.EqualOrClose(wrapperspb.Double(1.0), wrapperspb.Double(1.0000000000000002), exactComparison)).To(BeFalse())

	// Different message types should be unequal
	g.Expect(moreproto.EqualOrClose(wrapperspb.Double(1.0), wrapperspb.Float(1.0), exactComparison)).To(BeFalse())
}

func TestEqualOrCloseUnknownFields(t *testing.T) {
	g := NewWithT(t)

	// Create messages with unknown fields by marshaling and unmarshaling
	// with different proto definitions
	msg1 := wrapperspb.Double(1.0)
	msg2 := wrapperspb.Double(1.0)

	// Test exact equality first
	result := moreproto.EqualOrClose(msg1, msg2, morefloat.WithinMantissaThreshold(0, 0))
	g.Expect(result).To(BeTrue(), "Expected equal messages to be equal")

	// Create a message with unknown fields by adding some arbitrary bytes
	data1, _ := proto.Marshal(msg1)
	data2 := append([]byte(nil), data1...) // Copy data1 to avoid modifying it
	data2 = append(data2, 0x08, 0x42)      // Add unknown field with tag 1, value 66

	var msgWithUnknown wrapperspb.DoubleValue
	g.Expect(proto.Unmarshal(data2, &msgWithUnknown)).To(Succeed())

	// Messages with different unknown fields should not be equal
	result = moreproto.EqualOrClose(msg1, &msgWithUnknown, morefloat.WithinMantissaThreshold(0, 0))
	g.Expect(result).To(BeFalse(), "Expected messages with different unknown fields to be unequal")

	// Messages with same unknown fields should be equal
	var msgWithUnknown2 wrapperspb.DoubleValue
	g.Expect(proto.Unmarshal(data2, &msgWithUnknown2)).To(Succeed())

	result = moreproto.EqualOrClose(&msgWithUnknown, &msgWithUnknown2, morefloat.WithinMantissaThreshold(0, 0))
	g.Expect(result).To(BeTrue(), "Expected messages with same unknown fields to be equal")
}
