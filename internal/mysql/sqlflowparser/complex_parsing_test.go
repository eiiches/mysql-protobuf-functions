package sqlflowparser_test

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/mysql/sqlflowparser"
	. "github.com/onsi/gomega"
)

func TestNestedControlStructures(t *testing.T) {
	testCases := []struct {
		name               string
		sql                string
		expectedType       string
		expectedName       string
		shouldParseSuccess bool
	}{
		{
			name: "Nested IF with CASE and SELECT",
			sql: `CREATE PROCEDURE test()
BEGIN
    IF IF((SELECT CASE "THEN" WHEN "THEN" THEN "(" ELSE "" END) = "(", TRUE, FALSE) THEN
        SET @a = 1;
    END IF;
END`,
			expectedType:       "procedure",
			expectedName:       "test",
			shouldParseSuccess: true,
		},
		{
			name: "Nested IF with CASE, SELECT, and comment containing THEN",
			sql: `CREATE PROCEDURE test()
BEGIN
    IF IF((SELECT CASE "THEN" WHEN "THEN" THEN "(" ELSE "" END) = "(", TRUE, FALSE) /* THEN */ THEN
        SET @a = 1;
    END IF;
END`,
			expectedType:       "procedure",
			expectedName:       "test",
			shouldParseSuccess: true,
		},
		{
			name: "Complex nested expressions with quoted keywords",
			sql: `CREATE FUNCTION complex_test() RETURNS INT DETERMINISTIC
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
END`,
			expectedType:       "function",
			expectedName:       "complex_test",
			shouldParseSuccess: true,
		},
		{
			name: "Nested control structures with quoted strings",
			sql: `CREATE PROCEDURE nested_control()
BEGIN
    DECLARE done BOOLEAN DEFAULT FALSE;
    DECLARE counter INT DEFAULT 0;
    
    WHILE counter < 10 DO
        IF counter % 2 = 0 THEN
            CASE counter
                WHEN 0 THEN SET @msg = "START";
                WHEN 2 THEN SET @msg = "EVEN";
                WHEN 4 THEN SET @msg = "MIDDLE";
                ELSE SET @msg = "OTHER";
            END CASE;
        END IF;
        
        SET counter = counter + 1;
    END WHILE;
END`,
			expectedType:       "procedure",
			expectedName:       "nested_control",
			shouldParseSuccess: true,
		},
		{
			name: "Labeled procedure with complex nested structures",
			sql: `CREATE PROCEDURE labeled_proc()
BEGIN
    loop_label: LOOP
        IF @condition THEN
            LEAVE loop_label;
        END IF;
        
        SET @counter = @counter + 1;
        
        IF @counter > 100 THEN
            LEAVE loop_label;
        END IF;
    END LOOP;
END`,
			expectedType:       "procedure",
			expectedName:       "labeled_proc",
			shouldParseSuccess: true,
		},
		{
			name: "Function with cursors and complex logic",
			sql: `CREATE FUNCTION cursor_function(input_id INT) RETURNS TEXT READS SQL DATA
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
END`,
			expectedType:       "function",
			expectedName:       "cursor_function",
			shouldParseSuccess: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			result, err := sqlflowparser.Parse("", []byte(tc.sql))

			if tc.shouldParseSuccess {
				g.Expect(err).ToNot(HaveOccurred(), "Parser should not return an error")
				g.Expect(result).ToNot(BeNil(), "Result should not be nil")

				switch tc.expectedType {
				case "procedure":
					proc, ok := result.(*sqlflowparser.CreateProcedureStmt)
					g.Expect(ok).To(BeTrue(), "Should parse as CreateProcedureStmt")
					g.Expect(proc.Name).To(Equal(tc.expectedName), "Should have expected procedure name")
				case "function":
					fn, ok := result.(*sqlflowparser.CreateFunctionStmt)
					g.Expect(ok).To(BeTrue(), "Should parse as CreateFunctionStmt")
					g.Expect(fn.Name).To(Equal(tc.expectedName), "Should have expected function name")
				}
			} else {
				g.Expect(err).To(HaveOccurred(), "Parser should return an error for invalid SQL")
			}
		})
	}
}

func TestEdgeCasesWithComplexExpressions(t *testing.T) {
	testCases := []struct {
		name string
		sql  string
	}{
		{
			name: "String literals with special characters",
			sql: `CREATE FUNCTION test_special_chars() RETURNS TEXT DETERMINISTIC
BEGIN
    RETURN 'This string contains "quotes" and \'apostrophes\' and (parentheses)';
END`,
		},
		{
			name: "Binary literals and hex values",
			sql: `CREATE FUNCTION test_binary() RETURNS BLOB DETERMINISTIC
BEGIN
    RETURN _binary X'48656C6C6F';
END`,
		},
		{
			name: "Complex mathematical expressions",
			sql: `CREATE FUNCTION test_math() RETURNS DECIMAL(10,2) DETERMINISTIC
BEGIN
    RETURN (((5 * 3) + 2) / 4) - 1.5;
END`,
		},
		{
			name: "Nested function calls",
			sql: `CREATE FUNCTION test_nested_calls() RETURNS INT DETERMINISTIC
BEGIN
    RETURN LENGTH(SUBSTRING(CONCAT('Hello', ' ', 'World'), 1, 5));
END`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)

			result, err := sqlflowparser.Parse("", []byte(tc.sql))

			g.Expect(err).ToNot(HaveOccurred(), "Parser should not return an error")
			g.Expect(result).ToNot(BeNil(), "Result should not be nil")

			// Should parse as a function
			fn, ok := result.(*sqlflowparser.CreateFunctionStmt)
			g.Expect(ok).To(BeTrue(), "Should parse as CreateFunctionStmt")
			g.Expect(fn.Name).ToNot(BeEmpty(), "Function should have a name")
		})
	}
}

func TestDeeplyNestedControlStructures(t *testing.T) {
	g := NewWithT(t)

	// Test deeply nested control structures
	sql := `CREATE PROCEDURE deeply_nested()
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE j INT DEFAULT 0;
    
    outer_loop: WHILE i < 10 DO
        SET j = 0;
        
        inner_loop: WHILE j < 5 DO
            IF i = 5 THEN
                IF j = 2 THEN
                    CASE i + j
                        WHEN 7 THEN
                            BEGIN
                                IF TRUE THEN
                                    SET @result = 'Found it!';
                                END IF;
                            END;
                        ELSE
                            SET @result = 'Not found';
                    END CASE;
                END IF;
            END IF;
            
            SET j = j + 1;
        END WHILE;
        
        SET i = i + 1;
    END WHILE;
END`

	result, err := sqlflowparser.Parse("", []byte(sql))

	g.Expect(err).ToNot(HaveOccurred(), "Parser should not return an error")
	g.Expect(result).ToNot(BeNil(), "Result should not be nil")

	// Should parse as a procedure
	proc, ok := result.(*sqlflowparser.CreateProcedureStmt)
	g.Expect(ok).To(BeTrue(), "Should parse as CreateProcedureStmt")
	g.Expect(proc.Name).To(Equal("deeply_nested"), "Should have correct procedure name")

	// Should have a BEGIN statement containing the body
	g.Expect(proc.Body).To(HaveLen(1), "Should have one BEGIN statement")
	beginStmt, ok := proc.Body[0].(*sqlflowparser.BeginStmt)
	g.Expect(ok).To(BeTrue(), "Should have BEGIN statement")
	g.Expect(beginStmt.Body).ToNot(BeEmpty(), "BEGIN should have a body")
}
