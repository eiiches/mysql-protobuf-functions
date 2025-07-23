package sqlinstrument_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	. "github.com/onsi/gomega"
)

func TestFirstFunctionParsing(t *testing.T) {
	g := NewWithT(t)

	// The first function from protobuf.sql
	input := `CREATE FUNCTION _pb_util_bin_as_int32(b BLOB) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(b) > 4 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_util_bin_as_int32: value must not be longer than 4 bytes.';
	END IF;

	IF LPAD(b, 4, _binary X'00') & _binary X'80000000' = _binary X'00000000' THEN
		RETURN CONV(HEX(b), 16, 10);
	ELSE
		RETURN -(CONV(HEX(~b), 16, 10) + 1);
	END IF;
END`

	result, err := sqlflowparser.Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*sqlflowparser.CreateFunctionStmt)
	g.Expect(ok).To(BeTrue())

	// Should have one BEGIN statement in body
	g.Expect(stmt.Body).To(HaveLen(1))

	// Get the BEGIN statement
	beginStmt, ok := stmt.Body[0].(*sqlflowparser.BeginStmt)
	g.Expect(ok).To(BeTrue())

	// Should have 3 statements inside BEGIN: DECLARE + 2 IF statements
	g.Expect(beginStmt.Body).To(HaveLen(3))

	// First should be DECLARE
	declareStmt, ok := beginStmt.Body[0].(*sqlflowparser.DeclareStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(declareStmt.Text).To(ContainSubstring("DECLARE CUSTOM_EXCEPTION"))

	// Second should be IF statement
	_, ok = beginStmt.Body[1].(*sqlflowparser.IfStmt)
	g.Expect(ok).To(BeTrue())

	// Third should be IF statement
	ifStmt2, ok := beginStmt.Body[2].(*sqlflowparser.IfStmt)
	g.Expect(ok).To(BeTrue())

	// Debug the THEN branch structure
	g.Expect(ifStmt2.Then).To(HaveLen(1), "Should have one statement in THEN branch")
	_, ok = ifStmt2.Then[0].(*sqlflowparser.ReturnStmt)
	g.Expect(ok).To(BeTrue(), "THEN branch should contain a RETURN statement")

	// Debug the ELSE branch structure
	g.Expect(ifStmt2.Else).To(HaveLen(1), "Should have one statement in ELSE branch")
	_, ok = ifStmt2.Else[0].(*sqlflowparser.ReturnStmt)
	g.Expect(ok).To(BeTrue(), "ELSE branch should contain a RETURN statement")
}
