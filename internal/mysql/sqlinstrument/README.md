# SQL Instrumenter Package

A Go package for adding coverage instrumentation to MySQL SQL files. The package uses `sqlsplitter` to parse SQL files into statements and `sqlflowparser` to analyze the statement structure, then injects `__record_coverage` calls before executable statements.

## Features

- **Statement-Level Instrumentation**: Adds coverage tracking calls before executable statements
- **Function and Procedure Support**: Instruments CREATE FUNCTION and CREATE PROCEDURE statements
- **Line Number Tracking**: Records the original line number for each instrumented statement
- **Smart Statement Detection**: Only instruments executable statements, skipping declarations and comments
- **Delimiter Preservation**: Correctly handles DELIMITER statements and custom delimiters
- **AST-Aware**: Uses sqlflowparser to understand statement structure for accurate instrumentation

## Installation

```bash
go get github.com/eiiches/mysql-coverage/cmd/pigeon-poc/pkg/sqlinstrument
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/eiiches/mysql-coverage/cmd/pigeon-poc/pkg/sqlinstrument"
)

func main() {
    content, err := os.ReadFile("input.sql")
    if err != nil {
        log.Fatal(err)
    }

    instrumenter := sqlinstrument.NewInstrumenter("input.sql")
    instrumented, err := instrumenter.InstrumentSQL(content)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Print(instrumented)
}
```

## API

### Types

#### Instrumenter
```go
type Instrumenter struct {
    // Internal fields...
}
```

### Functions

#### NewInstrumenter
```go
func NewInstrumenter(filename string) *Instrumenter
```
Creates a new instrumenter instance for the given filename. The filename is used in the coverage calls and will be extracted to just the base name (e.g., "path/to/file.sql" becomes "file.sql").

#### InstrumentSQL
```go
func (i *Instrumenter) InstrumentSQL(content []byte) (string, error)
```
Instruments the SQL content by adding `__record_coverage` calls before executable statements. Returns the instrumented SQL as a string or an error if instrumentation fails.

## Instrumentation Behavior

### Coverage Call Format
The instrumenter adds calls in the format:
```sql
CALL __record_coverage('<filename>', '<function_or_procedure_name>', <line_number>);
```

### Instrumented Statements
The following statement types are instrumented:
- Control flow: `IF`, `WHILE`, `LOOP`, `REPEAT`, `CASE` statements
- Data manipulation: `SET` assignments, `SELECT` queries, `INSERT`, `UPDATE`, `DELETE` statements
- Procedure calls: `CALL` statements
- Flow control: `RETURN`, `LEAVE`, `ITERATE` statements
- Error handling: `SIGNAL` statements
- Generic SQL: Any unrecognized executable statement

### Non-Instrumented Statements
The following are **not** instrumented:
- `DECLARE` statements
- Comments (`--`, `#`, `/* */`)
- `BEGIN` and `END` blocks themselves (but statements inside are instrumented)
- Labels (ending with `:`)
- `WHEN`, `ELSE`, `ELSEIF`, `UNTIL` clauses themselves (but statements inside are instrumented)
- Empty lines and whitespace-only lines

### Statement Types Handled
- **SQL Statements**: CREATE FUNCTION and CREATE PROCEDURE are parsed and instrumented
- **DELIMITER Statements**: Passed through unchanged
- **COMMENT Statements**: Passed through unchanged
- **Other SQL**: Returned unchanged if not a function or procedure

## Examples

### Basic Function Instrumentation

**Input:**
```sql
CREATE FUNCTION calc_tax(amount DECIMAL(10,2)) RETURNS DECIMAL(10,2) DETERMINISTIC
BEGIN
    DECLARE tax DECIMAL(10,2);
    SET tax = amount * 0.1;
    RETURN tax;
END
```

**Output:**
```sql
CREATE FUNCTION calc_tax(amount DECIMAL(10,2)) RETURNS DECIMAL(10,2) DETERMINISTIC
BEGIN
    DECLARE tax DECIMAL(10,2);
    CALL __record_coverage('input.sql', 'calc_tax', 4); SET tax = amount * 0.1;
    CALL __record_coverage('input.sql', 'calc_tax', 5); RETURN tax;
END
```

### IF Statement Instrumentation

**Input:**
```sql
CREATE FUNCTION abs_value(x INT) RETURNS INT DETERMINISTIC
BEGIN
    IF x < 0 THEN
        RETURN -x;
    ELSE
        RETURN x;
    END IF;
END
```

**Output:**
```sql
CREATE FUNCTION abs_value(x INT) RETURNS INT DETERMINISTIC
BEGIN
    CALL __record_coverage('input.sql', 'abs_value', 3); IF x < 0 THEN
        CALL __record_coverage('input.sql', 'abs_value', 4); RETURN -x;
    ELSE
        CALL __record_coverage('input.sql', 'abs_value', 6); RETURN x;
    END IF;
END
```

### Procedure Instrumentation

**Input:**
```sql
CREATE PROCEDURE update_user(IN user_id INT, IN new_name VARCHAR(50))
BEGIN
    DECLARE user_exists INT DEFAULT 0;
    SELECT COUNT(*) INTO user_exists FROM users WHERE id = user_id;
    IF user_exists > 0 THEN
        UPDATE users SET name = new_name WHERE id = user_id;
    END IF;
END
```

**Output:**
```sql
CREATE PROCEDURE update_user(IN user_id INT, IN new_name VARCHAR(50))
BEGIN
    DECLARE user_exists INT DEFAULT 0;
    CALL __record_coverage('input.sql', 'update_user', 4); SELECT COUNT(*) INTO user_exists FROM users WHERE id = user_id;
    CALL __record_coverage('input.sql', 'update_user', 5); IF user_exists > 0 THEN
        CALL __record_coverage('input.sql', 'update_user', 6); UPDATE users SET name = new_name WHERE id = user_id;
    END IF;
END
```

### Delimiter Handling

**Input:**
```sql
DELIMITER $$
CREATE FUNCTION test_func(x INT) RETURNS INT DETERMINISTIC
BEGIN
    RETURN x * 2;
END$$
DELIMITER ;
```

**Output:**
```sql
DELIMITER $$
CREATE FUNCTION test_func(x INT) RETURNS INT DETERMINISTIC
BEGIN
    CALL __record_coverage('input.sql', 'test_func', 4); RETURN x * 2;
END$$
DELIMITER ;
```

## Testing

The package includes comprehensive unit tests using Gomega:

```bash
go test ./pkg/sqlinstrument -v
```

## Design

This instrumenter is designed to:
- Parse SQL files while respecting MySQL's delimiter semantics
- Add minimal overhead by only instrumenting executable statements
- Use AST parsing to accurately identify function and procedure boundaries
- Generate coverage calls that can be processed by coverage analysis tools

The instrumenter works in two phases:
1. **Statement Splitting**: Uses `sqlsplitter` to break the SQL file into individual statements
2. **AST Analysis**: Uses `sqlflowparser` to parse CREATE FUNCTION/PROCEDURE statements and instrument their bodies

This approach ensures accurate instrumentation while maintaining compatibility with complex SQL constructs like custom delimiters, string literals, and comments.

## Dependencies

- `github.com/eiiches/mysql-coverage/pkg/sqlsplitter` - SQL statement splitting
- `github.com/eiiches/mysql-coverage/pkg/sqlflowparser` - MySQL AST parsing
- `github.com/onsi/gomega` - Testing framework (dev dependency)