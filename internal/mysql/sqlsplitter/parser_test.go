package sqlsplitter

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewParser(t *testing.T) {
	g := NewWithT(t)

	input := []byte("SELECT 1;")
	parser := NewParser(input)

	g.Expect(parser.input).To(Equal(input))
	g.Expect(parser.pos).To(Equal(0))
	g.Expect(parser.line).To(Equal(1))
	g.Expect(parser.lineStart).To(Equal(0))
	g.Expect(parser.delimiter).To(Equal(";"))
}

func TestBasicSQLStatement(t *testing.T) {
	g := NewWithT(t)

	input := []byte("SELECT 1;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(1))
	g.Expect(statements[0].Text).To(Equal("SELECT 1"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))
	g.Expect(statements[0].StartPos).To(Equal(0))
	g.Expect(statements[0].EndPos).To(Equal(9))
}

func TestMultipleStatements(t *testing.T) {
	g := NewWithT(t)

	input := []byte("SELECT 1;\nSELECT 2;\nSELECT 3;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(3))

	g.Expect(statements[0].Text).To(Equal("SELECT 1"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))

	g.Expect(statements[1].Text).To(Equal("SELECT 2"))
	g.Expect(statements[1].Type).To(Equal("SQL"))
	g.Expect(statements[1].LineNo).To(Equal(2))

	g.Expect(statements[2].Text).To(Equal("SELECT 3"))
	g.Expect(statements[2].Type).To(Equal("SQL"))
	g.Expect(statements[2].LineNo).To(Equal(3))
}

func TestDelimiterChange(t *testing.T) {
	g := NewWithT(t)

	input := []byte("DELIMITER $$\nSELECT 1$$\nSELECT 2$$\nDELIMITER ;\nSELECT 3;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(5))

	g.Expect(statements[0].Text).To(Equal("DELIMITER $$"))
	g.Expect(statements[0].Type).To(Equal("DELIMITER"))
	g.Expect(statements[0].LineNo).To(Equal(1))

	g.Expect(statements[1].Text).To(Equal("SELECT 1"))
	g.Expect(statements[1].Type).To(Equal("SQL"))
	g.Expect(statements[1].LineNo).To(Equal(2))

	g.Expect(statements[2].Text).To(Equal("SELECT 2"))
	g.Expect(statements[2].Type).To(Equal("SQL"))
	g.Expect(statements[2].LineNo).To(Equal(3))

	g.Expect(statements[3].Text).To(Equal("DELIMITER ;"))
	g.Expect(statements[3].Type).To(Equal("DELIMITER"))
	g.Expect(statements[3].LineNo).To(Equal(4))

	g.Expect(statements[4].Text).To(Equal("SELECT 3"))
	g.Expect(statements[4].Type).To(Equal("SQL"))
	g.Expect(statements[4].LineNo).To(Equal(5))
}

func TestArbitraryDelimiter(t *testing.T) {
	g := NewWithT(t)

	input := []byte("DELIMITER ENDOFSTATEMENT\nSELECT 1ENDOFSTATEMENT\nSELECT 2ENDOFSTATEMENT")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(3))

	g.Expect(statements[0].Text).To(Equal("DELIMITER ENDOFSTATEMENT"))
	g.Expect(statements[0].Type).To(Equal("DELIMITER"))
	g.Expect(statements[0].LineNo).To(Equal(1))

	g.Expect(statements[1].Text).To(Equal("SELECT 1"))
	g.Expect(statements[1].Type).To(Equal("SQL"))
	g.Expect(statements[1].LineNo).To(Equal(2))

	g.Expect(statements[2].Text).To(Equal("SELECT 2"))
	g.Expect(statements[2].Type).To(Equal("SQL"))
	g.Expect(statements[2].LineNo).To(Equal(3))
}

func TestStringLiterals(t *testing.T) {
	g := NewWithT(t)

	input := []byte(`SELECT 'test;string', "another;test", ` + "`table;name`" + `;`)
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(1))
	g.Expect(statements[0].Text).To(Equal(`SELECT 'test;string', "another;test", ` + "`table;name`"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))
}

func TestEscapedQuotes(t *testing.T) {
	g := NewWithT(t)

	input := []byte(`SELECT 'It''s a test', "She said \"Hello\"";`)
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(1))
	g.Expect(statements[0].Text).To(Equal(`SELECT 'It''s a test', "She said \"Hello\""`))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))
}

func TestLineComments(t *testing.T) {
	g := NewWithT(t)

	input := []byte("SELECT 1; -- This is a comment\nSELECT 2; # Another comment\n")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(3))

	g.Expect(statements[0].Text).To(Equal("SELECT 1"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))

	g.Expect(statements[1].Text).To(Equal("-- This is a comment\nSELECT 2"))
	g.Expect(statements[1].Type).To(Equal("SQL"))
	g.Expect(statements[1].LineNo).To(Equal(1))

	g.Expect(statements[2].Text).To(Equal("# Another comment"))
	g.Expect(statements[2].Type).To(Equal("COMMENT"))
	g.Expect(statements[2].LineNo).To(Equal(2))
}

func TestBlockComments(t *testing.T) {
	g := NewWithT(t)

	input := []byte("/* Block comment */ SELECT 1;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(1))
	g.Expect(statements[0].Text).To(Equal("/* Block comment */ SELECT 1"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))
}

func TestMultilineBlockComments(t *testing.T) {
	g := NewWithT(t)

	input := []byte("/* Multi-line\n   comment */ SELECT 1;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(1))
	g.Expect(statements[0].Text).To(Equal("/* Multi-line\n   comment */ SELECT 1"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))
}

func TestEmptyStatements(t *testing.T) {
	g := NewWithT(t)

	input := []byte(";;SELECT 1;;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(1)) // Empty statements should be filtered out
	g.Expect(statements[0].Text).To(Equal("SELECT 1"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))
}

func TestWhitespaceHandling(t *testing.T) {
	g := NewWithT(t)

	input := []byte("   \n\t  SELECT 1  \n\t  ;   \n\t  ")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(1))
	g.Expect(statements[0].Text).To(Equal("SELECT 1"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(2))
}

func TestComplexMixedContent(t *testing.T) {
	g := NewWithT(t)

	input := []byte(`-- Test comment
SELECT 'string;with;semicolons' AS col1,
       "another;string" AS col2;
       
DELIMITER $$

CREATE PROCEDURE test_proc()
BEGIN
    SELECT 1;
    SELECT 'test;string';
END$$

DELIMITER ;

SELECT 'final';`)

	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(5))

	// First statement should include the comment
	g.Expect(statements[0].Text).To(ContainSubstring("-- Test comment"))
	g.Expect(statements[0].Text).To(ContainSubstring("SELECT 'string;with;semicolons'"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))

	// Delimiter change
	g.Expect(statements[1].Text).To(Equal("DELIMITER $$"))
	g.Expect(statements[1].Type).To(Equal("DELIMITER"))

	// Procedure with semicolons inside
	g.Expect(statements[2].Text).To(ContainSubstring("CREATE PROCEDURE test_proc()"))
	g.Expect(statements[2].Text).To(ContainSubstring("SELECT 1;"))
	g.Expect(statements[2].Text).To(ContainSubstring("SELECT 'test;string';"))
	g.Expect(statements[2].Type).To(Equal("SQL"))

	// Delimiter change back to semicolon
	g.Expect(statements[3].Text).To(Equal("DELIMITER ;"))
	g.Expect(statements[3].Type).To(Equal("DELIMITER"))

	// Final statement
	g.Expect(statements[4].Text).To(Equal("SELECT 'final'"))
	g.Expect(statements[4].Type).To(Equal("SQL"))
}

func TestDelimiterCaseInsensitive(t *testing.T) {
	g := NewWithT(t)

	input := []byte("delimiter $$\nSELECT 1$$\nDelImItEr ;\nSELECT 2;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(4))

	g.Expect(statements[0].Text).To(Equal("DELIMITER $$"))
	g.Expect(statements[0].Type).To(Equal("DELIMITER"))

	g.Expect(statements[1].Text).To(Equal("SELECT 1"))
	g.Expect(statements[1].Type).To(Equal("SQL"))

	g.Expect(statements[2].Text).To(Equal("DELIMITER ;"))
	g.Expect(statements[2].Type).To(Equal("DELIMITER"))

	g.Expect(statements[3].Text).To(Equal("SELECT 2"))
	g.Expect(statements[3].Type).To(Equal("SQL"))
}

func TestPositionTracking(t *testing.T) {
	g := NewWithT(t)

	input := []byte("SELECT 1;\nSELECT 2;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(2))

	// First statement
	g.Expect(statements[0].StartPos).To(Equal(0))
	g.Expect(statements[0].EndPos).To(Equal(9))

	// Second statement
	g.Expect(statements[1].StartPos).To(Equal(10))
	g.Expect(statements[1].EndPos).To(Equal(19))
}

func TestEdgeCases(t *testing.T) {
	t.Run("Empty input", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte("")
		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(0))
	})

	t.Run("Only whitespace", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte("   \n\t  \n  ")
		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(0))
	})

	t.Run("Only delimiter change", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte("DELIMITER $$")
		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(1))
		g.Expect(statements[0].Text).To(Equal("DELIMITER $$"))
		g.Expect(statements[0].Type).To(Equal("DELIMITER"))
	})

	t.Run("Statement without delimiter", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte("SELECT 1")
		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(1))
		g.Expect(statements[0].Text).To(Equal("SELECT 1"))
		g.Expect(statements[0].Type).To(Equal("SQL"))
	})
}

func TestStringWithNewlines(t *testing.T) {
	g := NewWithT(t)

	input := []byte("SELECT 'line1\nline2\nline3';")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(1))
	g.Expect(statements[0].Text).To(Equal("SELECT 'line1\nline2\nline3'"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))
}

func TestCommentWithDelimiter(t *testing.T) {
	g := NewWithT(t)

	input := []byte("SELECT 1; -- Comment with; delimiter\nSELECT 2;")
	parser := NewParser(input)
	statements, err := parser.Parse()

	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(statements).To(HaveLen(2))

	g.Expect(statements[0].Text).To(Equal("SELECT 1"))
	g.Expect(statements[0].Type).To(Equal("SQL"))
	g.Expect(statements[0].LineNo).To(Equal(1))

	g.Expect(statements[1].Text).To(Equal("-- Comment with; delimiter\nSELECT 2"))
	g.Expect(statements[1].Type).To(Equal("SQL"))
	g.Expect(statements[1].LineNo).To(Equal(1))
}
