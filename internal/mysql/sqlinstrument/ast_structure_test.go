package sqlinstrument_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	. "github.com/onsi/gomega"
)

func TestDeclareCursorParsedAsOneStatement(t *testing.T) {
	g := NewWithT(t)

	input := `CREATE PROCEDURE test_proc()
BEGIN
	DECLARE element_cursor CURSOR FOR
		SELECT i, n, t, v
		FROM table1
		ORDER BY i;
	SET message = 'test';
END`

	result, err := sqlflowparser.Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*sqlflowparser.CreateProcedureStmt)
	g.Expect(ok).To(BeTrue())

	// Should have one BEGIN statement in body
	g.Expect(stmt.Body).To(HaveLen(1))

	// It should be a BEGIN statement
	beginStmt, ok := stmt.Body[0].(*sqlflowparser.BeginStmt)
	g.Expect(ok).To(BeTrue())

	// BEGIN should have exactly 2 statements in its body
	g.Expect(beginStmt.Body).To(HaveLen(2))

	// Both should be statements
	firstStmt, ok := beginStmt.Body[0].(*sqlflowparser.DeclareStmt)
	g.Expect(ok).To(BeTrue())

	secondStmt, ok := beginStmt.Body[1].(*sqlflowparser.GenericStmt)
	g.Expect(ok).To(BeTrue())

	// First statement should contain the entire DECLARE CURSOR including SELECT
	g.Expect(firstStmt.Text).To(ContainSubstring("DECLARE element_cursor CURSOR FOR"))
	g.Expect(firstStmt.Text).To(ContainSubstring("SELECT i, n, t, v"))
	g.Expect(firstStmt.Text).To(ContainSubstring("ORDER BY i"))

	// Second statement should be the SET
	g.Expect(secondStmt.Text).To(ContainSubstring("SET message = 'test'"))

	// Verify positions - DECLARE CURSOR should start at line 3
	g.Expect(firstStmt.GetPosition().Line).To(Equal(3))

	// SET should be at a later line
	g.Expect(secondStmt.GetPosition().Line).To(BeNumerically(">", firstStmt.GetPosition().Line))
}
