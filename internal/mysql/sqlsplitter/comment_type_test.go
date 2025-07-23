package sqlsplitter

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestCommentTypeClassification(t *testing.T) {
	t.Run("Single line comment with SQL", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`-- This is a comment
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(1))

		g.Expect(statements[0].Type).To(Equal("SQL"))
		g.Expect(statements[0].Text).To(Equal("-- This is a comment\nSELECT 1"))
	})

	t.Run("Multiple line comments with SQL", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`-- First comment
-- Second comment
-- Third comment
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(1))

		g.Expect(statements[0].Type).To(Equal("SQL"))
		g.Expect(statements[0].Text).To(Equal("-- First comment\n-- Second comment\n-- Third comment\nSELECT 1"))
	})

	t.Run("Block comment only", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`/* This is a block comment */;
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(2))

		g.Expect(statements[0].Type).To(Equal("COMMENT"))
		g.Expect(statements[0].Text).To(Equal("/* This is a block comment */"))

		g.Expect(statements[1].Type).To(Equal("SQL"))
		g.Expect(statements[1].Text).To(Equal("SELECT 1"))
	})

	t.Run("Hash comment with SQL", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`# This is a hash comment
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(1))

		g.Expect(statements[0].Type).To(Equal("SQL"))
		g.Expect(statements[0].Text).To(Equal("# This is a hash comment\nSELECT 1"))
	})

	t.Run("Mixed comment types with SQL", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`-- Line comment
/* Block comment */
# Hash comment
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(1))

		g.Expect(statements[0].Type).To(Equal("SQL"))
		g.Expect(statements[0].Text).To(ContainSubstring("-- Line comment"))
		g.Expect(statements[0].Text).To(ContainSubstring("/* Block comment */"))
		g.Expect(statements[0].Text).To(ContainSubstring("# Hash comment"))
		g.Expect(statements[0].Text).To(ContainSubstring("SELECT 1"))
	})

	t.Run("Comment with SQL should be SQL type", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`-- Comment before SQL
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(1))

		g.Expect(statements[0].Type).To(Equal("SQL"))
		g.Expect(statements[0].Text).To(Equal("-- Comment before SQL\nSELECT 1"))
	})

	t.Run("Comment with whitespace and SQL", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`   -- Comment with leading whitespace
	/* Block comment with tabs */   
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(1))

		g.Expect(statements[0].Type).To(Equal("SQL"))
		g.Expect(statements[0].Text).To(ContainSubstring("-- Comment with leading whitespace"))
		g.Expect(statements[0].Text).To(ContainSubstring("/* Block comment with tabs */"))
		g.Expect(statements[0].Text).To(ContainSubstring("SELECT 1"))
	})

	t.Run("Multiline block comment only", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`/* 
   Multi-line
   block comment
   with multiple lines
*/;
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(2))

		g.Expect(statements[0].Type).To(Equal("COMMENT"))
		g.Expect(statements[0].Text).To(ContainSubstring("Multi-line"))
		g.Expect(statements[0].Text).To(ContainSubstring("block comment"))

		g.Expect(statements[1].Type).To(Equal("SQL"))
		g.Expect(statements[1].Text).To(Equal("SELECT 1"))
	})

	t.Run("Comment only with semicolon delimiter", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`-- This is just a comment
;
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(2))

		// First statement is comment only
		g.Expect(statements[0].Type).To(Equal("COMMENT"))
		g.Expect(statements[0].Text).To(Equal("-- This is just a comment"))

		// Second statement is SQL
		g.Expect(statements[1].Type).To(Equal("SQL"))
		g.Expect(statements[1].Text).To(Equal("SELECT 1"))
	})

	t.Run("Mixed comments only with semicolon delimiter", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`-- Line comment
/* Block comment */
# Hash comment
;
SELECT 1;`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(2))

		// First statement is comment only
		g.Expect(statements[0].Type).To(Equal("COMMENT"))
		g.Expect(statements[0].Text).To(ContainSubstring("-- Line comment"))
		g.Expect(statements[0].Text).To(ContainSubstring("/* Block comment */"))
		g.Expect(statements[0].Text).To(ContainSubstring("# Hash comment"))

		// Second statement is SQL
		g.Expect(statements[1].Type).To(Equal("SQL"))
		g.Expect(statements[1].Text).To(Equal("SELECT 1"))
	})

	t.Run("Instrumented SQL pattern with comments", func(t *testing.T) {
		g := NewWithT(t)

		input := []byte(`-- INSTRUMENTED SQL FILE
-- Original: protobuf.sql
-- Generated by MySQL Coverage Instrumentor

-- Coverage tracking table
SELECT 1;

-- Coverage recording procedure
DELIMITER //
SELECT 2 //`)

		parser := NewParser(input)
		statements, err := parser.Parse()

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(statements).To(HaveLen(4))

		// First statement has SQL, so it should be SQL type
		g.Expect(statements[0].Type).To(Equal("SQL"))
		g.Expect(statements[0].Text).To(ContainSubstring("SELECT 1"))

		// Second statement is comment only
		g.Expect(statements[1].Type).To(Equal("COMMENT"))
		g.Expect(statements[1].Text).To(Equal("-- Coverage recording procedure"))

		// DELIMITER statement
		g.Expect(statements[2].Type).To(Equal("DELIMITER"))
		g.Expect(statements[2].Text).To(Equal("DELIMITER //"))

		// SQL statement
		g.Expect(statements[3].Type).To(Equal("SQL"))
		g.Expect(statements[3].Text).To(Equal("SELECT 2"))
	})
}
