package sqlinstrument_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	. "github.com/onsi/gomega"
)

func TestSignalParsing(t *testing.T) {
	g := NewWithT(t)

	// Test just the SIGNAL statement
	input := `SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_read_i32_as_uint32: Unexpected end of BLOB.'`

	result, err := sqlflowparser.Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*sqlflowparser.GenericStmt)
	g.Expect(ok).To(BeTrue(), "Should parse as SQLStmt")
	g.Expect(stmt.Text).To(Equal(input))
}
