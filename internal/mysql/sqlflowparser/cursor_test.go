package sqlflowparser_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	. "github.com/onsi/gomega"
)

func TestDeclareCursorParsing(t *testing.T) {
	g := NewWithT(t)

	input := `CREATE PROCEDURE test_proc()
BEGIN
	DECLARE element_cursor CURSOR FOR
		SELECT i, n, t, v
		FROM JSON_TABLE(
			JSON_EXTRACT(data, '$.*[*]'),
			'$[*]' COLUMNS (
				i INT PATH '$.i',
				n INT PATH '$.n'
			)
		) jt
		ORDER BY i;
	SET message = 'test';
END`

	result, err := sqlflowparser.Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*sqlflowparser.CreateProcedureStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(stmt.Name).To(Equal("test_proc"))

	// The body should contain 1 BEGIN statement
	g.Expect(stmt.Body).To(HaveLen(1))
	beginStmt, ok := stmt.Body[0].(*sqlflowparser.BeginStmt)
	g.Expect(ok).To(BeTrue())

	// BEGIN should contain exactly 2 statements:
	// 1. DECLARE CURSOR statement (including the SELECT)
	// 2. SET statement
	g.Expect(beginStmt.Body).To(HaveLen(2), "Expected 2 statements in procedure body, got %d", len(beginStmt.Body))

	// Both statements should be GenericStmt since DECLARE and SET are not specifically parsed
	_, ok = beginStmt.Body[0].(*sqlflowparser.DeclareStmt)
	g.Expect(ok).To(BeTrue(), "First statement should be DeclareStmt, got %T", beginStmt.Body[0])

	_, ok = beginStmt.Body[1].(*sqlflowparser.SetVariableStmt)
	g.Expect(ok).To(BeTrue(), "Second statement should be SetVariableStmt, got %T", beginStmt.Body[1])
}

func TestSimpleDeclareParsing(t *testing.T) {
	g := NewWithT(t)

	input := `CREATE PROCEDURE test_proc()
BEGIN
	DECLARE done INT DEFAULT FALSE;
	DECLARE temp VARCHAR(100);
	SET temp = 'test';
END`

	result, err := sqlflowparser.Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*sqlflowparser.CreateProcedureStmt)
	g.Expect(ok).To(BeTrue())

	// Should have 1 BEGIN statement
	g.Expect(stmt.Body).To(HaveLen(1))
	beginStmt, ok := stmt.Body[0].(*sqlflowparser.BeginStmt)
	g.Expect(ok).To(BeTrue())

	// BEGIN should contain 3 statements: 2 DECLARE + 1 SET
	g.Expect(beginStmt.Body).To(HaveLen(3))

	// All statements should be GenericStmt since DECLARE and SET are not specifically parsed
	for i := 0; i < 3; i++ {
		if i < 2 {
			_, ok := beginStmt.Body[i].(*sqlflowparser.DeclareStmt)
			g.Expect(ok).To(BeTrue(), "Statement %d should be DeclareStmt, got %T", i, beginStmt.Body[i])
		} else {
			_, ok := beginStmt.Body[i].(*sqlflowparser.SetVariableStmt)
			g.Expect(ok).To(BeTrue(), "Statement %d should be SetVariableStmt, got %T", i, beginStmt.Body[i])
		}
	}
}
