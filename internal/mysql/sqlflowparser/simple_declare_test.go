package sqlflowparser_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	. "github.com/onsi/gomega"
)

func TestSimpleProcedureBody(t *testing.T) {
	g := NewWithT(t)

	// Very simple procedure with just SET
	input := `CREATE PROCEDURE test_proc()
BEGIN
	SET message = 'test';
END`

	result, err := sqlflowparser.Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*sqlflowparser.CreateProcedureStmt)
	g.Expect(ok).To(BeTrue())

	t.Logf("Body statements count: %d", len(stmt.Body))
	for i, bodyStmt := range stmt.Body {
		t.Logf("Statement %d: %T", i, bodyStmt)
	}

	// Should have 1 BEGIN statement
	g.Expect(stmt.Body).To(HaveLen(1))
	beginStmt, ok := stmt.Body[0].(*sqlflowparser.BeginStmt)
	g.Expect(ok).To(BeTrue())
	// BEGIN should contain 1 nested statement
	g.Expect(beginStmt.Body).To(HaveLen(1))
}

func TestSimpleDeclare(t *testing.T) {
	g := NewWithT(t)

	// Simple procedure with DECLARE
	input := `CREATE PROCEDURE test_proc()
BEGIN
	DECLARE temp INT;
	SET message = 'test';
END`

	result, err := sqlflowparser.Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*sqlflowparser.CreateProcedureStmt)
	g.Expect(ok).To(BeTrue())

	t.Logf("Body statements count: %d", len(stmt.Body))
	for i, bodyStmt := range stmt.Body {
		t.Logf("Statement %d: %T", i, bodyStmt)
	}

	// Should have 1 BEGIN statement
	g.Expect(stmt.Body).To(HaveLen(1))
	beginStmt, ok := stmt.Body[0].(*sqlflowparser.BeginStmt)
	g.Expect(ok).To(BeTrue())
	// BEGIN should contain 2 nested statements: DECLARE + SET
	g.Expect(beginStmt.Body).To(HaveLen(2))
}
