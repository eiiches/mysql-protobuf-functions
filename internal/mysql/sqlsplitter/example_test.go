package sqlsplitter

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

// Example shows basic usage of the sqlparser package
func ExampleParser() {
	input := []byte(`-- Example SQL file
SELECT 1;
DELIMITER $$
CREATE PROCEDURE example()
BEGIN
    SELECT 'test;statement';
END$$
DELIMITER ;
SELECT 'final';`)

	parser := NewParser(input)
	statements, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	for i, stmt := range statements {
		fmt.Printf("Statement %d [%s] (Line %d): %s\n",
			i+1, stmt.Type, stmt.LineNo, stmt.Text)
	}

	// Output:
	// Statement 1 [SQL] (Line 1): -- Example SQL file
	// SELECT 1
	// Statement 2 [DELIMITER] (Line 3): DELIMITER $$
	// Statement 3 [SQL] (Line 4): CREATE PROCEDURE example()
	// BEGIN
	//     SELECT 'test;statement';
	// END
	// Statement 4 [DELIMITER] (Line 8): DELIMITER ;
	// Statement 5 [SQL] (Line 9): SELECT 'final'
}

func TestExampleUsage(t *testing.T) {
	g := NewWithT(t)

	// Example showing how to use the parser with line number tracking
	input := []byte("SELECT 1;\nSELECT 2;\nSELECT 3;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(3))

	// Verify each statement has correct line number
	for i, stmt := range statements {
		expectedLine := i + 1
		g.Expect(stmt.LineNo).To(Equal(expectedLine))
		g.Expect(stmt.Type).To(Equal("SQL"))
		g.Expect(stmt.Text).To(Equal(fmt.Sprintf("SELECT %d", expectedLine)))
	}
}
