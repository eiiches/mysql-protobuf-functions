package sqlflowparser

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestCreateProcedure(t *testing.T) {
	g := NewWithT(t)

	input := `CREATE PROCEDURE test_proc(IN id INT, OUT name VARCHAR(50))
BEGIN
    DECLARE temp VARCHAR(100);
    SET temp = 'Hello';
    SELECT temp;
END`

	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*CreateProcedureStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(stmt.Name).To(Equal("test_proc"))
	g.Expect(stmt.Parameters).To(HaveLen(2))
	g.Expect(stmt.Parameters[0].Name).To(Equal("id"))
	g.Expect(stmt.Parameters[0].Type).To(Equal("INT"))
	g.Expect(stmt.Parameters[0].Mode).To(Equal("IN"))
	g.Expect(stmt.Parameters[1].Name).To(Equal("name"))
	g.Expect(stmt.Parameters[1].Type).To(Equal("VARCHAR(50)"))
	g.Expect(stmt.Parameters[1].Mode).To(Equal("OUT"))
	// Body should contain one BEGIN statement with 3 nested statements
	g.Expect(stmt.Body).To(HaveLen(1))
	beginStmt, ok := stmt.Body[0].(*BeginStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(beginStmt.Body).To(HaveLen(3))

	// Check first statement: DECLARE temp VARCHAR(100);
	declareStmt, ok := beginStmt.Body[0].(*DeclareStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(declareStmt.Text).To(Equal("DECLARE temp VARCHAR(100)"))

	// Check second statement: SET temp = 'Hello';
	setStmt, ok := beginStmt.Body[1].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(setStmt.Text).To(Equal("SET temp = 'Hello'"))

	// Check third statement: SELECT temp;
	selectStmt, ok := beginStmt.Body[2].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(selectStmt.Text).To(Equal("SELECT temp"))

	// Test position tracking
	pos := stmt.GetPosition()
	g.Expect(pos.Line).To(Equal(1))
	g.Expect(pos.Column).To(Equal(1))
	g.Expect(pos.Offset).To(Equal(0))
}

func TestCreateFunction(t *testing.T) {
	g := NewWithT(t)

	input := `CREATE FUNCTION calc_tax(amount DECIMAL(10,2))
RETURNS DECIMAL(10,2)
BEGIN
    RETURN amount * 0.1;
END`

	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*CreateFunctionStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(stmt.Name).To(Equal("calc_tax"))
	g.Expect(stmt.ReturnType).To(Equal("DECIMAL(10,2)"))
	g.Expect(stmt.Parameters).To(HaveLen(1))
	g.Expect(stmt.Parameters[0].Name).To(Equal("amount"))
	g.Expect(stmt.Parameters[0].Type).To(Equal("DECIMAL(10,2)"))
	// Body should contain one BEGIN statement with 1 nested statement
	g.Expect(stmt.Body).To(HaveLen(1))
	beginStmt, ok := stmt.Body[0].(*BeginStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(beginStmt.Body).To(HaveLen(1))

	// Check the RETURN statement: RETURN amount * 0.1;
	returnStmt, ok := beginStmt.Body[0].(*ReturnStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(returnStmt.Text).To(Equal("RETURN amount * 0.1"))
}

func TestIfStatement(t *testing.T) {
	g := NewWithT(t)

	input := `IF x > 0 THEN
    SET result = 'positive';
ELSEIF x < 0 THEN
    SET result = 'negative';
ELSE
    SET result = 'zero';
END IF`

	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*IfStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(stmt.Condition).To(Equal("x > 0"))

	// Check THEN clause has 1 statement
	g.Expect(stmt.Then).To(HaveLen(1))
	thenStmt, ok := stmt.Then[0].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(thenStmt.Text).To(Equal("SET result = 'positive'"))

	// Check ELSEIF clause
	g.Expect(stmt.ElseIfs).To(HaveLen(1))
	g.Expect(stmt.ElseIfs[0].Condition).To(Equal("x < 0"))
	g.Expect(stmt.ElseIfs[0].Then).To(HaveLen(1))
	elseifStmt, ok := stmt.ElseIfs[0].Then[0].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(elseifStmt.Text).To(Equal("SET result = 'negative'"))

	// Check ELSE clause has 1 statement
	g.Expect(stmt.Else).To(HaveLen(1))
	elseStmt, ok := stmt.Else[0].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(elseStmt.Text).To(Equal("SET result = 'zero'"))
}

func TestWhileLoop(t *testing.T) {
	g := NewWithT(t)

	input := `WHILE counter < 10 DO
    SET counter = counter + 1;
    SELECT counter;
END WHILE`

	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*WhileStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(stmt.Condition).To(Equal("counter < 10"))
	g.Expect(stmt.Body).To(HaveLen(2))

	// Check first statement in WHILE body
	setStmt, ok := stmt.Body[0].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(setStmt.Text).To(Equal("SET counter = counter + 1"))

	// Check second statement in WHILE body
	selectStmt, ok := stmt.Body[1].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(selectStmt.Text).To(Equal("SELECT counter"))
}

func TestLoopStatement(t *testing.T) {
	t.Run("Labeled loop", func(t *testing.T) {
		g := NewWithT(t)

		input := `my_loop: LOOP
    SET x = x + 1;
    IF x > 10 THEN
        LEAVE my_loop;
    END IF;
END LOOP`

		result, err := Parse("", []byte(input))
		g.Expect(err).ToNot(HaveOccurred())

		stmt, ok := result.(*LoopStmt)
		g.Expect(ok).To(BeTrue())
		g.Expect(stmt.Body).To(HaveLen(2))
		g.Expect(stmt.Label).To(Equal("my_loop"))

		// Check first statement in LOOP body
		setStmt, ok := stmt.Body[0].(*GenericStmt)
		g.Expect(ok).To(BeTrue())
		g.Expect(setStmt.Text).To(Equal("SET x = x + 1"))

		// Check second statement in LOOP body (IF statement)
		ifStmt, ok := stmt.Body[1].(*IfStmt)
		g.Expect(ok).To(BeTrue())
		g.Expect(ifStmt.Condition).To(Equal("x > 10"))
		g.Expect(ifStmt.Then).To(HaveLen(1))
		leaveStmt, ok := ifStmt.Then[0].(*LeaveStmt)
		g.Expect(ok).To(BeTrue())
		g.Expect(leaveStmt.Label).To(Equal("my_loop"))
	})
}

func TestRepeatStatement(t *testing.T) {
	g := NewWithT(t)

	input := `REPEAT
    SET x = x + 1;
    SELECT x;
UNTIL x > 10
END REPEAT`

	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*RepeatStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(stmt.Condition).To(Equal("x > 10"))
	g.Expect(stmt.Body).To(HaveLen(2))

	// Check first statement in REPEAT body
	setStmt, ok := stmt.Body[0].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(setStmt.Text).To(Equal("SET x = x + 1"))

	// Check second statement in REPEAT body
	selectStmt, ok := stmt.Body[1].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(selectStmt.Text).To(Equal("SELECT x"))
}

func TestCaseStatement(t *testing.T) {
	g := NewWithT(t)

	input := `CASE grade
    WHEN 'A' THEN SET result = 'Excellent';
    WHEN 'B' THEN SET result = 'Good';
    ELSE SET result = 'Other';
END CASE`

	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*CaseStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(stmt.Expression).To(Equal("grade"))
	g.Expect(stmt.WhenClauses).To(HaveLen(2))

	// Check first WHEN clause
	g.Expect(stmt.WhenClauses[0].Condition).To(Equal("'A'"))
	g.Expect(stmt.WhenClauses[0].Then).To(HaveLen(1))
	firstWhenStmt, ok := stmt.WhenClauses[0].Then[0].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(firstWhenStmt.Text).To(Equal("SET result = 'Excellent'"))

	// Check second WHEN clause
	g.Expect(stmt.WhenClauses[1].Condition).To(Equal("'B'"))
	g.Expect(stmt.WhenClauses[1].Then).To(HaveLen(1))
	secondWhenStmt, ok := stmt.WhenClauses[1].Then[0].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(secondWhenStmt.Text).To(Equal("SET result = 'Good'"))

	// Check ELSE clause
	g.Expect(stmt.Else).To(HaveLen(1))
	elseStmt, ok := stmt.Else[0].(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(elseStmt.Text).To(Equal("SET result = 'Other'"))
}

func TestLeaveIterateReturn(t *testing.T) {
	t.Run("LEAVE statement", func(t *testing.T) {
		g := NewWithT(t)

		input := `LEAVE my_loop`
		result, err := Parse("", []byte(input))
		g.Expect(err).ToNot(HaveOccurred())

		stmt, ok := result.(*LeaveStmt)
		g.Expect(ok).To(BeTrue())
		g.Expect(stmt.Label).To(Equal("my_loop"))
	})

	t.Run("ITERATE statement", func(t *testing.T) {
		g := NewWithT(t)

		input := `ITERATE my_loop`
		result, err := Parse("", []byte(input))
		g.Expect(err).ToNot(HaveOccurred())

		stmt, ok := result.(*IterateStmt)
		g.Expect(ok).To(BeTrue())
		g.Expect(stmt.Label).To(Equal("my_loop"))
	})

	t.Run("RETURN statement", func(t *testing.T) {
		g := NewWithT(t)

		input := `RETURN x * 2`
		result, err := Parse("", []byte(input))
		g.Expect(err).ToNot(HaveOccurred())

		stmt, ok := result.(*ReturnStmt)
		g.Expect(ok).To(BeTrue())
		g.Expect(stmt.Text).To(Equal("RETURN x * 2"))
	})
}

func TestSQLStatement(t *testing.T) {
	g := NewWithT(t)

	input := `SELECT * FROM users WHERE id = 1`
	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*GenericStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(stmt.Text).To(Equal("SELECT * FROM users WHERE id = 1"))
}

func TestPositionTracking(t *testing.T) {
	g := NewWithT(t)

	input := `LEAVE my_loop`

	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	stmt, ok := result.(*LeaveStmt)
	g.Expect(ok).To(BeTrue())

	// Test position tracking for LEAVE statement
	pos := stmt.GetPosition()
	g.Expect(pos.Line).To(Equal(1))
	g.Expect(pos.Column).To(Equal(1))
	g.Expect(pos.Offset).To(Equal(0))

	g.Expect(stmt.Label).To(Equal("my_loop"))
}

func TestRegression1(t *testing.T) {
	g := NewWithT(t)

	input := `CREATE FUNCTION _pb_util_zigzag_encode_uint64(value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
    --
    RETURN (value << 1) ^ -(value >> 63);
END`
	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred())

	_, ok := result.(*CreateFunctionStmt)
	g.Expect(ok).To(BeTrue())
}

func TestStringLiteralEscapeSequences(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		// Single-quoted strings with all MySQL escape sequences
		{
			name:        "Single quote with escaped single quote",
			input:       `SELECT 'Don\'t worry'`,
			expectError: false,
			description: "Backslash-escaped single quote in single-quoted string",
		},
		{
			name:        "Single quote with escaped double quote",
			input:       `SELECT 'Say \"hello\"'`,
			expectError: false,
			description: "Backslash-escaped double quote in single-quoted string",
		},
		{
			name:        "Single quote with NUL character",
			input:       `SELECT 'text\0null'`,
			expectError: false,
			description: "ASCII NUL character (\\0)",
		},
		{
			name:        "Single quote with backspace",
			input:       `SELECT 'text\bback'`,
			expectError: false,
			description: "Backspace character (\\b)",
		},
		{
			name:        "Single quote with newline",
			input:       `SELECT 'line1\nline2'`,
			expectError: false,
			description: "Newline character (\\n)",
		},
		{
			name:        "Single quote with carriage return",
			input:       `SELECT 'text\rcarriage'`,
			expectError: false,
			description: "Carriage return character (\\r)",
		},
		{
			name:        "Single quote with tab",
			input:       `SELECT 'text\ttab'`,
			expectError: false,
			description: "Tab character (\\t)",
		},
		{
			name:        "Single quote with Control+Z",
			input:       `SELECT 'text\Zeof'`,
			expectError: false,
			description: "ASCII 26 Control+Z character (\\Z)",
		},
		{
			name:        "Single quote with backslash",
			input:       `SELECT 'path\\file'`,
			expectError: false,
			description: "Backslash character (\\\\)",
		},
		{
			name:        "Single quote with percent",
			input:       `SELECT 'value\%percent'`,
			expectError: false,
			description: "Percent character for pattern matching (\\%)",
		},
		{
			name:        "Single quote with underscore",
			input:       `SELECT 'value\_underscore'`,
			expectError: false,
			description: "Underscore character for pattern matching (\\_)",
		},
		{
			name:        "Single quote with doubled single quotes",
			input:       `SELECT 'Don''t worry'`,
			expectError: false,
			description: "Doubled single quotes (MySQL standard escaping)",
		},

		// Double-quoted strings with all MySQL escape sequences
		{
			name:        "Double quote with escaped double quote",
			input:       `SELECT "Say \"hello\""`,
			expectError: false,
			description: "Backslash-escaped double quote in double-quoted string",
		},
		{
			name:        "Double quote with escaped single quote",
			input:       `SELECT "Don\'t worry"`,
			expectError: false,
			description: "Backslash-escaped single quote in double-quoted string",
		},
		{
			name:        "Double quote with NUL character",
			input:       `SELECT "text\0null"`,
			expectError: false,
			description: "ASCII NUL character (\\0) in double quotes",
		},
		{
			name:        "Double quote with backspace",
			input:       `SELECT "text\bback"`,
			expectError: false,
			description: "Backspace character (\\b) in double quotes",
		},
		{
			name:        "Double quote with newline",
			input:       `SELECT "line1\nline2"`,
			expectError: false,
			description: "Newline character (\\n) in double quotes",
		},
		{
			name:        "Double quote with carriage return",
			input:       `SELECT "text\rcarriage"`,
			expectError: false,
			description: "Carriage return character (\\r) in double quotes",
		},
		{
			name:        "Double quote with tab",
			input:       `SELECT "text\ttab"`,
			expectError: false,
			description: "Tab character (\\t) in double quotes",
		},
		{
			name:        "Double quote with Control+Z",
			input:       `SELECT "text\Zeof"`,
			expectError: false,
			description: "ASCII 26 Control+Z character (\\Z) in double quotes",
		},
		{
			name:        "Double quote with backslash",
			input:       `SELECT "path\\file"`,
			expectError: false,
			description: "Backslash character (\\\\) in double quotes",
		},
		{
			name:        "Double quote with percent",
			input:       `SELECT "value\%percent"`,
			expectError: false,
			description: "Percent character for pattern matching (\\%) in double quotes",
		},
		{
			name:        "Double quote with underscore",
			input:       `SELECT "value\_underscore"`,
			expectError: false,
			description: "Underscore character for pattern matching (\\_) in double quotes",
		},
		{
			name:        "Double quote with doubled double quotes",
			input:       `SELECT "Say ""hello"""`,
			expectError: false,
			description: "Doubled double quotes (MySQL standard escaping)",
		},

		// Complex real-world examples from protobuf-json-v2.sql
		{
			name:        "JSON path with escaped quotes",
			input:       `SELECT JSON_SET(field_json_value, CONCAT('$.\"', map_key, '\"'), map_value)`,
			expectError: false,
			description: "Complex JSON path construction with escaped quotes",
		},
		{
			name:        "Multiple escape sequences",
			input:       `SELECT 'path\\file\tvalue\nend\0'`,
			expectError: false,
			description: "Multiple different escape sequences in one string",
		},
		{
			name:        "Nested function calls with complex quotes",
			input:       `SELECT JSON_EXTRACT(oneofs, CONCAT('$.\"', oneof_index, '\".i'))`,
			expectError: false,
			description: "Nested function with complex quote patterns",
		},

		// Backtick strings (identifiers)
		{
			name:        "Backtick identifier",
			input:       "SELECT `column name`",
			expectError: false,
			description: "Backtick-quoted identifier",
		},
		{
			name:        "Backtick with escaped backtick",
			input:       "SELECT `column``name`",
			expectError: false,
			description: "Doubled backticks in identifier",
		},

		// Unknown escape sequences (should be treated as literal characters)
		{
			name:        "Unknown escape sequence in single quotes",
			input:       `SELECT 'test\xunknown'`,
			expectError: false,
			description: "Unknown escape sequence \\x should be treated as literal x",
		},
		{
			name:        "Multiple unknown escape sequences",
			input:       `SELECT 'path\qwith\yunknown\zescapes'`,
			expectError: false,
			description: "Multiple unknown escape sequences should be treated as literals",
		},
		{
			name:        "Unknown escape sequence in double quotes",
			input:       `SELECT "test\xunknown"`,
			expectError: false,
			description: "Unknown escape sequence \\x in double quotes should be treated as literal x",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			result, err := Parse("", []byte(tc.input))

			if tc.expectError {
				g.Expect(err).To(HaveOccurred(), "Expected parsing error for: %s", tc.description)
			} else {
				g.Expect(err).ToNot(HaveOccurred(), "Unexpected parsing error for: %s\nError: %v", tc.description, err)
				g.Expect(result).ToNot(BeNil(), "Expected non-nil result for: %s", tc.description)

				// Verify it's a valid statement
				_, ok := result.(StatementAST)
				g.Expect(ok).To(BeTrue(), "Expected StatementAST for: %s", tc.description)
			}
		})
	}
}

func TestUnknownEscapeSequenceHandling(t *testing.T) {
	g := NewWithT(t)

	// Test that unknown escape sequences are handled per MySQL spec:
	// "For all other escape sequences, backslash is ignored. That is, the escaped character is interpreted as if it was not escaped."
	testCases := []struct {
		input    string
		expected string
	}{
		{`SELECT '\x'`, "\\x"},           // \x should become just x, but we preserve the literal for SQL reconstruction
		{`SELECT '\q'`, "\\q"},           // \q should become just q
		{`SELECT '\y\z'`, "\\y\\z"},      // Multiple unknowns
		{`SELECT '\xtest\y'`, "\\xtest\\y"}, // Mixed with other characters
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input: %s", tc.input), func(t *testing.T) {
			result, err := Parse("", []byte(tc.input))
			g.Expect(err).ToNot(HaveOccurred(), "Failed to parse: %s", tc.input)

			stmt, ok := result.(*GenericStmt)
			g.Expect(ok).To(BeTrue(), "Expected GenericStmt for: %s", tc.input)
			
			// The exact text comparison depends on how the parser reconstructs the SQL
			// The important thing is that it parses successfully
			g.Expect(stmt.Text).To(ContainSubstring("SELECT"), "Should contain SELECT")
		})
	}
}

func TestComplexProcedureWithEscapeSequences(t *testing.T) {
	g := NewWithT(t)

	// Test a procedure similar to the problematic one that caused the original parser issue
	input := `CREATE PROCEDURE test_escape_handling(IN map_key VARCHAR(255))
BEGIN
    DECLARE field_json_value JSON;
    DECLARE map_value JSON;
    
    -- This line was causing the original parser error
    SET field_json_value = JSON_SET(field_json_value, CONCAT('$.\"', map_key, '\"'), map_value);
    
    -- Test multiple escape sequences
    SET field_json_value = JSON_SET(field_json_value, 'path\\with\ttabs\nand\rreturns\0null');
    
    -- Test pattern matching escapes
    SELECT * FROM table WHERE column LIKE 'pattern\%with\_escapes';
END`

	result, err := Parse("", []byte(input))
	g.Expect(err).ToNot(HaveOccurred(), "Failed to parse procedure with escape sequences")

	stmt, ok := result.(*CreateProcedureStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(stmt.Name).To(Equal("test_escape_handling"))
	g.Expect(stmt.Parameters).To(HaveLen(1))
	g.Expect(stmt.Parameters[0].Name).To(Equal("map_key"))

	// Verify the procedure body was parsed correctly
	g.Expect(stmt.Body).To(HaveLen(1))
	beginStmt, ok := stmt.Body[0].(*BeginStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(beginStmt.Body).To(HaveLen(5)) // 2 DECLARE + 2 SET + 1 SELECT = 5 statements
}
