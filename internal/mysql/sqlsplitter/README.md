# SQL Splitter Package

A Go package for parsing MySQL SQL files with support for dynamic delimiters, comments, string literals, and line number tracking.

## Features

- **Dynamic Delimiter Support**: Handles `DELIMITER` statements that change the statement terminator
- **String Literal Parsing**: Correctly parses single quotes, double quotes, and backticks with escape sequences
- **Comment Handling**: Supports line comments (`--`, `#`) and block comments (`/* */`)
- **Line Number Tracking**: Returns the line number where each statement begins
- **Position Tracking**: Tracks start and end positions in the original input
- **Arbitrary Delimiters**: Supports delimiters of any length (e.g., `DELIMITER ENDOFSTATEMENT`)

## Installation

```bash
go get github.com/eiiches/mysql-coverage/pkg/sqlsplitter
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/eiiches/mysql-coverage/pkg/sqlsplitter"
)

func main() {
    input := []byte(`-- Example SQL file
SELECT 1;
DELIMITER $$
CREATE PROCEDURE example()
BEGIN
    SELECT 'test;statement';
END$$
DELIMITER ;
SELECT 'final';`)

    parser := sqlsplitter.NewParser(input)
    statements, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    for i, stmt := range statements {
        fmt.Printf("Statement %d [%s] (Line %d): %s\n", 
            i+1, stmt.Type, stmt.LineNo, stmt.Text)
    }
}
```

## API

### Types

#### Statement
```go
type Statement struct {
    Text     string // The statement text
    Type     string // "SQL", "DELIMITER", or "COMMENT"
    LineNo   int    // Line number where the statement begins (1-based)
    StartPos int    // Starting position in the original input
    EndPos   int    // Ending position in the original input
}
```

#### Parser
```go
type Parser struct {
    // Internal fields...
}
```

### Functions

#### NewParser
```go
func NewParser(input []byte) *Parser
```
Creates a new parser instance for the given input.

#### Parse
```go
func (p *Parser) Parse() ([]Statement, error)
```
Parses the input and returns a slice of statements. Returns an error if parsing fails.

## Output

The parser returns a slice of `Statement` objects, each representing a parsed statement from the input. The parser categorizes statements into three types:

- **SQL**: Regular SQL statements (SELECT, INSERT, CREATE, etc.)
- **DELIMITER**: Delimiter change statements (e.g., `DELIMITER $$`)
- **COMMENT**: Comment-only statements (lines containing only comments and whitespace)

### Statement Processing

- Empty statements (just delimiters) are filtered out
- Comments are preserved as part of SQL statements they precede
- Standalone comments become separate COMMENT statements
- String literals and comments are parsed correctly to avoid false delimiter matches
- Line numbers reflect where each statement begins (1-based)
- Start/end positions track character positions in the original input

## Examples

### Basic SQL Statements
```sql
SELECT 1;
INSERT INTO users (name) VALUES ('John');
```

Output:
```
Statement 1 [SQL] (Line 1): SELECT 1
Statement 2 [SQL] (Line 2): INSERT INTO users (name) VALUES ('John')
```

### Custom Delimiters
```sql
DELIMITER $$
CREATE PROCEDURE test()
BEGIN
    SELECT 1;
    SELECT 2;
END$$
DELIMITER ;
```

Output:
```
Statement 1 [DELIMITER] (Line 1): DELIMITER $$
Statement 2 [SQL] (Line 2): CREATE PROCEDURE test()
BEGIN
    SELECT 1;
    SELECT 2;
END
Statement 3 [DELIMITER] (Line 7): DELIMITER ;
```

### String Literals with Delimiters
```sql
SELECT 'text;with;semicolons', "another;string", `table;name`;
```

Output:
```
Statement 1 [SQL] (Line 1): SELECT 'text;with;semicolons', "another;string", `table;name`
```

### Comments
```sql
-- Line comment
SELECT 1; -- End of line comment
/* Block comment */ SELECT 2;
# Hash comment
```

Output:
```
Statement 1 [SQL] (Line 1): -- Line comment
SELECT 1
Statement 2 [SQL] (Line 2): -- End of line comment
/* Block comment */ SELECT 2
Statement 3 [COMMENT] (Line 4): # Hash comment
```

### Escaped Quotes
```sql
SELECT 'It''s a test';
SELECT "She said \"Hello\"";
SELECT `table``name`;
```

Output:
```
Statement 1 [SQL] (Line 1): SELECT 'It''s a test'
Statement 2 [SQL] (Line 2): SELECT "She said \"Hello\""
Statement 3 [SQL] (Line 3): SELECT `table``name`
```

## Testing

The package includes comprehensive unit tests using Gomega:

```bash
go test ./pkg/sqlsplitter -v
```

## Design

This parser implements a SQL parser that focuses on:
- Statement-level parsing (splitting SQL into individual statements)
- Client-side delimiter handling (like MySQL client)
- Proper string and comment parsing to avoid false delimiter matches

It does **not** provide:
- Full SQL AST parsing
- Expression-level parsing
- Syntax validation beyond basic structure

This makes it ideal for tools that need to split SQL files into individual statements while respecting MySQL's delimiter semantics.
