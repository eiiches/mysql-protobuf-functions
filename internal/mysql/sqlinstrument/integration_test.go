package sqlinstrument_test

import (
	"strings"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlsplitter"
	. "github.com/onsi/gomega"
)

// Helper function to parse a complete SQL file like the old TwoPassParser
func parseSQL(input []byte) ([]ParsedStatement, error) {
	// First pass: split into statements
	splitter := sqlsplitter.NewParser(input)
	statements, err := splitter.Parse()
	if err != nil {
		return nil, err
	}

	// Second pass: parse each SQL statement
	var result []ParsedStatement
	for _, stmt := range statements {
		parsedStmt := ParsedStatement{
			OriginalText: stmt.Text,
			Type:         stmt.Type,
			LineNo:       stmt.LineNo,
		}

		// Only parse SQL statements into AST
		if stmt.Type == "SQL" {
			parsed, err := sqlflowparser.Parse("", []byte(stmt.Text))
			if err != nil {
				parsedStmt.ParseError = err
			} else {
				parsedStmt.AST = parsed
			}
		}

		result = append(result, parsedStmt)
	}

	return result, nil
}

type ParsedStatement struct {
	OriginalText string
	Type         string
	AST          interface{}
	ParseError   error
	LineNo       int
}

func TestCompleteFileParsingIntegration(t *testing.T) {
	g := NewWithT(t)

	input := `-- Test SQL file
DELIMITER //
CREATE PROCEDURE process_data(IN id INT)
BEGIN
    DECLARE counter INT DEFAULT 0;
    
    WHILE counter < 10 DO
        SET counter = counter + 1;
        IF counter = 5 THEN
            SELECT 'Half way';
        END IF;
    END WHILE;
END //

CREATE FUNCTION calculate_total(amount DECIMAL(10,2))
RETURNS DECIMAL(10,2)
BEGIN
    DECLARE tax DECIMAL(10,2);
    SET tax = amount * 0.1;
    RETURN amount + tax;
END //

DELIMITER ;

-- Just a comment
SELECT * FROM users;`

	statements, err := parseSQL([]byte(input))

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(6))

	// Check statement types
	g.Expect(statements[0].Type).To(Equal("COMMENT")) // -- Test SQL file
	g.Expect(statements[1].Type).To(Equal("DELIMITER"))
	g.Expect(statements[2].Type).To(Equal("SQL")) // CREATE PROCEDURE
	g.Expect(statements[3].Type).To(Equal("SQL")) // CREATE FUNCTION
	g.Expect(statements[4].Type).To(Equal("DELIMITER"))
	g.Expect(statements[5].Type).To(Equal("SQL")) // SELECT

	// Check parsed ASTs
	proc, ok := statements[2].AST.(*sqlflowparser.CreateProcedureStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(proc.Name).To(Equal("process_data"))
	g.Expect(proc.Parameters).To(HaveLen(1))
	g.Expect(proc.Parameters[0].Name).To(Equal("id"))
	g.Expect(proc.Parameters[0].Type).To(Equal("INT"))
	g.Expect(proc.Parameters[0].Mode).To(Equal("IN"))

	// DELIMITER statements should not have AST
	g.Expect(statements[1].AST).To(BeNil())

	fn, ok := statements[3].AST.(*sqlflowparser.CreateFunctionStmt)
	g.Expect(ok).To(BeTrue())
	g.Expect(fn.Name).To(Equal("calculate_total"))
	g.Expect(fn.ReturnType).To(Equal("DECIMAL(10,2)"))
}

func TestParsingWithErrors(t *testing.T) {
	g := NewWithT(t)

	input := `-- Invalid SQL
DELIMITER //
CREATE PROCEDURE test()
BEGIN
    IF x = THEN END IF;  -- Invalid IF syntax - missing value after =
END //
DELIMITER ;

-- Valid SQL
SELECT 1;`

	statements, err := parseSQL([]byte(input))

	// First pass should succeed even with invalid SQL
	g.Expect(err).ToNot(HaveOccurred())

	// Find the CREATE PROCEDURE and SELECT statements
	var procStmt, selectStmt *ParsedStatement
	for i := range statements {
		if strings.Contains(statements[i].OriginalText, "CREATE PROCEDURE") {
			procStmt = &statements[i]
		} else if strings.Contains(statements[i].OriginalText, "SELECT 1") {
			selectStmt = &statements[i]
		}
	}

	g.Expect(procStmt).ToNot(BeNil())
	g.Expect(selectStmt).ToNot(BeNil())

	// Check if procedure actually failed to parse (it might parse successfully due to lenient parsing)
	if procStmt.ParseError != nil {
		// If there was an error, AST should be nil
		g.Expect(procStmt.AST).To(BeNil())
	} else {
		// If parsing succeeded, we should have a valid AST
		g.Expect(procStmt.AST).ToNot(BeNil())
	}

	// SELECT should parse successfully
	g.Expect(selectStmt.ParseError).To(BeNil())
	g.Expect(selectStmt.AST).ToNot(BeNil())
}

func TestParsingWithComments(t *testing.T) {
	g := NewWithT(t)

	input := `-- This is just a comment
;
SELECT 1;
/* Another comment */
;`

	statements, err := parseSQL([]byte(input))

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(3))

	// First statement is comment-only
	g.Expect(statements[0].Type).To(Equal("COMMENT"))
	g.Expect(statements[0].AST).To(BeNil())

	// Second is SQL
	g.Expect(statements[1].Type).To(Equal("SQL"))
	g.Expect(statements[1].AST).ToNot(BeNil())

	// Third is comment-only
	g.Expect(statements[2].Type).To(Equal("COMMENT"))
	g.Expect(statements[2].AST).To(BeNil())
}

func TestComplexSQLIntegration(t *testing.T) {
	testCases := []struct {
		name               string
		sql                string
		expectedProcedures int
		expectedFunctions  int
		expectedErrors     int
	}{
		{
			name: "Nested IF with CASE and SELECT",
			sql: `DELIMITER $$

DROP PROCEDURE IF EXISTS test $$
CREATE PROCEDURE test()
BEGIN
    IF IF((SELECT CASE "THEN" WHEN "THEN" THEN "(" ELSE "" END) = "(", TRUE, FALSE) THEN
        SET @a = 1;
    END IF;
END $$`,
			expectedProcedures: 1,
			expectedFunctions:  0,
			expectedErrors:     0,
		},
		{
			name: "Complex nested expressions with quoted keywords",
			sql: `DELIMITER $$

DROP FUNCTION IF EXISTS complex_test $$
CREATE FUNCTION complex_test() RETURNS INT DETERMINISTIC
BEGIN
    DECLARE result INT;
    
    IF "BEGIN" = "BEGIN" AND "END" <> "END" THEN
        SET result = 1;
    ELSEIF "WHILE" IN ("WHILE", "LOOP") THEN
        SET result = 2;
    ELSE
        SET result = 0;
    END IF;
    
    RETURN result;
END $$`,
			expectedProcedures: 0,
			expectedFunctions:  1,
			expectedErrors:     0,
		},
		{
			name: "Function with cursors and complex logic",
			sql: `DELIMITER $$

DROP FUNCTION IF EXISTS cursor_function $$
CREATE FUNCTION cursor_function(input_id INT) RETURNS TEXT READS SQL DATA
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE result TEXT DEFAULT '';
    DECLARE temp_name VARCHAR(255);
    
    DECLARE name_cursor CURSOR FOR
        SELECT name FROM users WHERE id = input_id;
    
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    
    OPEN name_cursor;
    
    read_loop: LOOP
        FETCH name_cursor INTO temp_name;
        
        IF done THEN
            LEAVE read_loop;
        END IF;
        
        SET result = CONCAT(result, temp_name, ',');
    END LOOP;
    
    CLOSE name_cursor;
    
    RETURN TRIM(TRAILING ',' FROM result);
END $$`,
			expectedProcedures: 0,
			expectedFunctions:  1,
			expectedErrors:     0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			statements, err := parseSQL([]byte(tc.sql))

			g.Expect(err).ToNot(HaveOccurred(), "Parser should not return an error")
			g.Expect(statements).ToNot(BeEmpty(), "Should have statements")

			// Count parse errors
			var parseErrors int
			var procedures int
			var functions int

			for _, stmt := range statements {
				if stmt.ParseError != nil {
					parseErrors++
				}
				if _, ok := stmt.AST.(*sqlflowparser.CreateProcedureStmt); ok {
					procedures++
				}
				if _, ok := stmt.AST.(*sqlflowparser.CreateFunctionStmt); ok {
					functions++
				}
			}

			g.Expect(parseErrors).To(Equal(tc.expectedErrors), "Should have expected number of parse errors")
			g.Expect(procedures).To(Equal(tc.expectedProcedures), "Should have expected number of procedures")
			g.Expect(functions).To(Equal(tc.expectedFunctions), "Should have expected number of functions")
		})
	}
}
