package sqlsplitter

import (
	"strings"
	"unicode"
)

// Statement represents a parsed SQL statement with its line number
type Statement struct {
	Text     string // The statement text
	Type     string // "SQL" or "DELIMITER"
	LineNo   int    // Line number where the statement begins (1-based)
	StartPos int    // Starting position in the original input
	EndPos   int    // Ending position in the original input
}

// Parser handles parsing SQL files with dynamic delimiter support
type Parser struct {
	input     []byte
	pos       int
	line      int
	lineStart int
	delimiter string
}

// NewParser creates a new parser instance
func NewParser(input []byte) *Parser {
	return &Parser{
		input:     input,
		pos:       0,
		line:      1,
		lineStart: 0,
		delimiter: ";",
	}
}

// Parse parses the input and returns a slice of statements
func (p *Parser) Parse() ([]Statement, error) {
	var statements []Statement

	for p.pos < len(p.input) {
		p.skipWhitespace()

		if p.pos >= len(p.input) {
			break
		}

		// Record the starting position and line
		startPos := p.pos
		startLine := p.line

		// Check for DELIMITER statement only at the beginning of a line
		if p.isAtStartOfLine() && p.isDelimiterStatement() {
			stmt, err := p.parseDelimiterStatement()
			if err != nil {
				return nil, err
			}
			stmt.LineNo = startLine
			stmt.StartPos = startPos
			stmt.EndPos = p.pos
			statements = append(statements, stmt)
			continue
		}

		// Parse SQL statement
		stmt, err := p.parseSQLStatement()
		if err != nil {
			return nil, err
		}
		if stmt.Text != "" {
			stmt.LineNo = startLine
			stmt.StartPos = startPos
			stmt.EndPos = p.pos
			statements = append(statements, stmt)
		}
	}

	return statements, nil
}

func (p *Parser) isDelimiterStatement() bool {
	saved := p.pos
	savedLine := p.line
	savedLineStart := p.lineStart
	defer func() {
		p.pos = saved
		p.line = savedLine
		p.lineStart = savedLineStart
	}()

	p.skipWhitespace()

	// Check if next word is "DELIMITER" (case insensitive)
	word := p.readWord()
	return strings.ToLower(word) == "delimiter"
}

func (p *Parser) isDelimiterStatementAtCurrentPos() bool {
	saved := p.pos
	savedLine := p.line
	savedLineStart := p.lineStart
	defer func() {
		p.pos = saved
		p.line = savedLine
		p.lineStart = savedLineStart
	}()

	// Skip only horizontal whitespace (spaces and tabs), not newlines
	for p.pos < len(p.input) && (p.input[p.pos] == ' ' || p.input[p.pos] == '\t') {
		p.pos++
	}

	// Check if next word is "DELIMITER" (case insensitive)
	word := p.readWord()
	return strings.ToLower(word) == "delimiter"
}

func (p *Parser) parseDelimiterStatement() (Statement, error) {
	p.skipWhitespace()

	// Skip "DELIMITER" keyword
	p.readWord()

	p.skipWhitespace()

	// Read new delimiter
	newDelim := p.readDelimiterValue()
	p.delimiter = newDelim

	return Statement{
		Text: "DELIMITER " + newDelim,
		Type: "DELIMITER",
	}, nil
}

func (p *Parser) parseSQLStatement() (Statement, error) {
	var content strings.Builder

	for p.pos < len(p.input) {
		// Check for DELIMITER at the beginning of a line (but only if we're at actual start of line)
		if p.isAtStartOfLine() && p.isDelimiterStatementAtCurrentPos() {
			// We've found a DELIMITER statement, so we need to stop parsing this SQL statement
			break
		}

		// Check if we're at the delimiter
		if p.isAtDelimiter() {
			// Consume the delimiter
			p.pos += len(p.delimiter)
			p.updateLineInfo()
			break
		}

		// Handle different content types
		if p.isStringStart() {
			str := p.parseString()
			content.WriteString(str)
		} else if p.isCommentStart() {
			comment := p.parseComment()
			content.WriteString(comment)
		} else {
			// Regular character
			ch := p.input[p.pos]
			content.WriteByte(ch)
			if ch == '\n' {
				p.line++
				p.lineStart = p.pos + 1
			}
			p.pos++
		}
	}

	text := strings.TrimSpace(content.String())

	// Determine if this is a comment-only statement
	statementType := "SQL"
	if p.isCommentOnly(text) {
		statementType = "COMMENT"
	}

	return Statement{
		Text: text,
		Type: statementType,
	}, nil
}

func (p *Parser) isAtDelimiter() bool {
	if p.pos+len(p.delimiter) > len(p.input) {
		return false
	}
	return string(p.input[p.pos:p.pos+len(p.delimiter)]) == p.delimiter
}

func (p *Parser) isAtStartOfLine() bool {
	// We're at start of line if we're at position 0 or
	// the position is right after the lineStart
	if p.pos == 0 {
		return true
	}

	// Find the last newline before current position
	lastNewline := -1
	for i := p.pos - 1; i >= 0; i-- {
		if p.input[i] == '\n' || p.input[i] == '\r' {
			lastNewline = i
			break
		}
	}

	// If no newline found, we're on the first line
	if lastNewline == -1 {
		// Check if we've only seen whitespace from the beginning
		for i := 0; i < p.pos; i++ {
			if p.input[i] != ' ' && p.input[i] != '\t' {
				return false
			}
		}
		return true
	}

	// Check if we've only seen whitespace since the last newline
	for i := lastNewline + 1; i < p.pos; i++ {
		if p.input[i] != ' ' && p.input[i] != '\t' {
			return false
		}
	}
	return true
}

func (p *Parser) isStringStart() bool {
	if p.pos >= len(p.input) {
		return false
	}
	ch := p.input[p.pos]
	return ch == '\'' || ch == '"' || ch == '`'
}

func (p *Parser) parseString() string {
	if p.pos >= len(p.input) {
		return ""
	}

	quote := p.input[p.pos]
	var result strings.Builder
	result.WriteByte(quote)
	if p.input[p.pos] == '\n' {
		p.line++
		p.lineStart = p.pos + 1
	}
	p.pos++

	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		result.WriteByte(ch)
		if ch == '\n' {
			p.line++
			p.lineStart = p.pos + 1
		}
		p.pos++

		if ch == quote {
			// Check for escaped quote (doubled)
			if p.pos < len(p.input) && p.input[p.pos] == quote {
				result.WriteByte(quote)
				if p.input[p.pos] == '\n' {
					p.line++
					p.lineStart = p.pos + 1
				}
				p.pos++
			} else {
				// End of string
				break
			}
		} else if ch == '\\' && p.pos < len(p.input) {
			// Handle escaped characters
			result.WriteByte(p.input[p.pos])
			if p.input[p.pos] == '\n' {
				p.line++
				p.lineStart = p.pos + 1
			}
			p.pos++
		}
	}

	return result.String()
}

func (p *Parser) isCommentStart() bool {
	if p.pos >= len(p.input) {
		return false
	}

	// Check for -- comment
	if p.pos+1 < len(p.input) && p.input[p.pos] == '-' && p.input[p.pos+1] == '-' {
		// Must be followed by space/tab or end of line for valid comment
		if p.pos+2 >= len(p.input) || p.input[p.pos+2] == ' ' || p.input[p.pos+2] == '\t' ||
			p.input[p.pos+2] == '\r' || p.input[p.pos+2] == '\n' {
			return true
		}
	}

	// Check for /* comment
	if p.pos+1 < len(p.input) && p.input[p.pos] == '/' && p.input[p.pos+1] == '*' {
		return true
	}

	// Check for # comment
	if p.input[p.pos] == '#' {
		return true
	}

	return false
}

func (p *Parser) parseComment() string {
	if p.pos >= len(p.input) {
		return ""
	}

	var result strings.Builder

	if p.pos+1 < len(p.input) && p.input[p.pos] == '-' && p.input[p.pos+1] == '-' {
		// Line comment starting with --
		for p.pos < len(p.input) && p.input[p.pos] != '\n' && p.input[p.pos] != '\r' {
			result.WriteByte(p.input[p.pos])
			p.pos++
		}
		// Include the newline if present
		if p.pos < len(p.input) && (p.input[p.pos] == '\n' || p.input[p.pos] == '\r') {
			result.WriteByte(p.input[p.pos])
			if p.input[p.pos] == '\n' {
				p.line++
				p.lineStart = p.pos + 1
			}
			p.pos++
			// Handle \r\n
			if p.pos < len(p.input) && p.input[p.pos-1] == '\r' && p.input[p.pos] == '\n' {
				result.WriteByte(p.input[p.pos])
				p.line++
				p.lineStart = p.pos + 1
				p.pos++
			}
		}
	} else if p.pos+1 < len(p.input) && p.input[p.pos] == '/' && p.input[p.pos+1] == '*' {
		// Block comment
		result.WriteByte(p.input[p.pos])
		p.pos++
		result.WriteByte(p.input[p.pos])
		p.pos++

		for p.pos+1 < len(p.input) {
			result.WriteByte(p.input[p.pos])
			if p.input[p.pos] == '\n' {
				p.line++
				p.lineStart = p.pos + 1
			}
			if p.input[p.pos] == '*' && p.input[p.pos+1] == '/' {
				p.pos++
				result.WriteByte(p.input[p.pos])
				p.pos++
				break
			}
			p.pos++
		}
	} else if p.input[p.pos] == '#' {
		// Line comment starting with #
		for p.pos < len(p.input) && p.input[p.pos] != '\n' && p.input[p.pos] != '\r' {
			result.WriteByte(p.input[p.pos])
			p.pos++
		}
		// Include the newline if present
		if p.pos < len(p.input) && (p.input[p.pos] == '\n' || p.input[p.pos] == '\r') {
			result.WriteByte(p.input[p.pos])
			if p.input[p.pos] == '\n' {
				p.line++
				p.lineStart = p.pos + 1
			}
			p.pos++
			// Handle \r\n
			if p.pos < len(p.input) && p.input[p.pos-1] == '\r' && p.input[p.pos] == '\n' {
				result.WriteByte(p.input[p.pos])
				p.line++
				p.lineStart = p.pos + 1
				p.pos++
			}
		}
	}

	return result.String()
}

func (p *Parser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		if p.input[p.pos] == '\n' {
			p.line++
			p.lineStart = p.pos + 1
		}
		p.pos++
	}
}

func (p *Parser) readWord() string {
	var result strings.Builder

	for p.pos < len(p.input) {
		ch := rune(p.input[p.pos])
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			result.WriteRune(ch)
			p.pos++
		} else {
			break
		}
	}

	return result.String()
}

func (p *Parser) readDelimiterValue() string {
	var result strings.Builder

	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			break
		}
		result.WriteByte(ch)
		p.pos++
	}

	// Skip to end of line or next statement
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == '\r' || ch == '\n' {
			if ch == '\n' {
				p.line++
				p.lineStart = p.pos + 1
			}
			p.pos++
			if ch == '\r' && p.pos < len(p.input) && p.input[p.pos] == '\n' {
				p.line++
				p.lineStart = p.pos + 1
				p.pos++
			}
			break
		}
		p.pos++
	}

	return result.String()
}

func (p *Parser) updateLineInfo() {
	// Update line information for consumed delimiter
	for i := 0; i < len(p.delimiter); i++ {
		if p.pos-len(p.delimiter)+i < len(p.input) && p.input[p.pos-len(p.delimiter)+i] == '\n' {
			// This would have been handled during parsing, but just in case
		}
	}
}

func (p *Parser) isCommentOnly(text string) bool {
	if text == "" {
		return false
	}

	// Parse the text to see if it contains only comments and whitespace
	pos := 0
	input := []byte(text)

	for pos < len(input) {
		// Skip whitespace
		if unicode.IsSpace(rune(input[pos])) {
			pos++
			continue
		}

		// Check for line comment
		if pos+1 < len(input) && input[pos] == '-' && input[pos+1] == '-' {
			// Must be followed by space/tab or end of line for valid comment
			if pos+2 >= len(input) || input[pos+2] == ' ' || input[pos+2] == '\t' ||
				input[pos+2] == '\r' || input[pos+2] == '\n' {
				// Skip to end of line
				for pos < len(input) && input[pos] != '\n' && input[pos] != '\r' {
					pos++
				}
				continue
			}
		}

		// Check for block comment
		if pos+1 < len(input) && input[pos] == '/' && input[pos+1] == '*' {
			pos += 2
			// Skip to end of block comment
			for pos+1 < len(input) {
				if input[pos] == '*' && input[pos+1] == '/' {
					pos += 2
					break
				}
				pos++
			}
			continue
		}

		// Check for hash comment
		if input[pos] == '#' {
			// Skip to end of line
			for pos < len(input) && input[pos] != '\n' && input[pos] != '\r' {
				pos++
			}
			continue
		}

		// If we reach here, we found non-comment, non-whitespace content
		return false
	}

	// If we've processed all characters and found only comments/whitespace
	return true
}
