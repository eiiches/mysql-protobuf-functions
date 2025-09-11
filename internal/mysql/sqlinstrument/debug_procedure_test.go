package sqlinstrument_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	. "github.com/onsi/gomega"
)

func TestProcedureParsing(t *testing.T) {
	g := NewWithT(t)

	// Test the actual failing procedure from protobuf.sql
	input := `CREATE PROCEDURE _pb_wire_read_i32_as_uint32(IN buf LONGBLOB, OUT value INT UNSIGNED, OUT tail LONGBLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(buf) < 4 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_read_i32_as_uint32: Unexpected end of BLOB.';
	END IF;

	SET value = _pb_util_swap_endian_32(_pb_util_bin_as_uint32(LEFT(buf, 4)));
	SET tail = SUBSTRING(buf, 5);
END`

	result, err := sqlflowparser.Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*sqlflowparser.CreateProcedureStmt)
	g.Expect(ok).To(BeTrue(), "Should parse as CreateProcedureStmt")

	// Should have one BEGIN statement in body
	g.Expect(stmt.Body).To(HaveLen(1), "Should have 1 BEGIN statement in procedure body")

	// Get the BEGIN statement
	beginStmt, ok := stmt.Body[0].(*sqlflowparser.BeginStmt)
	g.Expect(ok).To(BeTrue(), "Should have BEGIN statement")

	// Should have 4 statements inside BEGIN: DECLARE + IF + 2 SET
	g.Expect(beginStmt.Body).To(HaveLen(4), "Should have 4 statements in BEGIN body")

	// Check statement types
	_, ok = beginStmt.Body[0].(*sqlflowparser.DeclareStmt) // DECLARE
	g.Expect(ok).To(BeTrue(), "First statement should be DECLARE")

	_, ok = beginStmt.Body[1].(*sqlflowparser.IfStmt) // IF
	g.Expect(ok).To(BeTrue(), "Second statement should be IF statement")

	_, ok = beginStmt.Body[2].(*sqlflowparser.SetVariableStmt) // SET value = ...
	g.Expect(ok).To(BeTrue(), "Third statement should be SET (SetVariableStmt)")

	_, ok = beginStmt.Body[3].(*sqlflowparser.SetVariableStmt) // SET tail = ...
	g.Expect(ok).To(BeTrue(), "Fourth statement should be SET (SetVariableStmt)")
}
