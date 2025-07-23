package sqlflowparser_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	. "github.com/onsi/gomega"
)

func TestDeclareCursorSingleStatement(t *testing.T) {
	g := NewWithT(t)

	// Just the procedure body with DECLARE CURSOR
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

	// Log what we actually got
	t.Logf("Procedure name: %s", stmt.Name)
	t.Logf("Body statements count: %d", len(stmt.Body))
	for i, bodyStmt := range stmt.Body {
		t.Logf("Statement %d: %T", i, bodyStmt)
	}

	// Should have 1 BEGIN statement
	g.Expect(stmt.Body).To(HaveLen(1))
	beginStmt, ok := stmt.Body[0].(*sqlflowparser.BeginStmt)
	g.Expect(ok).To(BeTrue())

	// BEGIN should contain exactly 2 statements:
	// 1. DECLARE CURSOR (which includes the SELECT)
	// 2. SET
	g.Expect(beginStmt.Body).To(HaveLen(2), "Expected exactly 2 statements: DECLARE CURSOR + SET")

	// Both should be GenericStmt (generic SQL statements)
	_, isDeclare := beginStmt.Body[0].(*sqlflowparser.DeclareStmt)
	g.Expect(isDeclare).To(BeTrue(), "First statement should be DeclareStmt, got %T", beginStmt.Body[0])

	_, isSet := beginStmt.Body[1].(*sqlflowparser.GenericStmt)
	g.Expect(isSet).To(BeTrue(), "Second statement should be GenericStmt, got %T", beginStmt.Body[1])
}
