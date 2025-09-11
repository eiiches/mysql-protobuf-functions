# SQL Flow Parser Package

A Go package for parsing individual MySQL SQL statements into Abstract Syntax Trees (AST), with focus on stored procedures, functions, and control flow statements. This parser provides structured representation of MySQL syntax for analysis and instrumentation.

**Important**: This parser handles **one statement at a time**. For parsing entire SQL files with multiple statements, use `@pkg/sqlsplitter` first to split the file into individual statements.

## Features

- **Stored Procedures**: Parses `CREATE PROCEDURE` statements with parameters and body
- **Functions**: Parses `CREATE FUNCTION` statements with parameters, return types, and body
- **Control Flow**: Supports `IF`, `WHILE`, `LOOP`, `REPEAT`, `CASE` statements
- **Variables**: Handles `DECLARE` and `SET` statements with detailed variable assignment parsing
- **Flow Control**: Parses `LEAVE`, `ITERATE`, `RETURN`, and `SIGNAL` statements
- **Generic Expressions**: Simplified expression parsing for conditions and values
- **Position Tracking**: Every AST node includes line number, column, and byte offset information
- **Error Handling**: Provides detailed error messages with line/column information

## Installation

```bash
# For single statement parsing
go get github.com/eiiches/mysql-coverage/pkg/sqlflowparser

# For multi-statement files (recommended)
go get github.com/eiiches/mysql-coverage/pkg/sqlflowparser
go get github.com/eiiches/mysql-coverage/pkg/sqlsplitter
```

## Usage

### Single Statement Parsing

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/eiiches/mysql-coverage/pkg/sqlflowparser"
)

func main() {
    // Parse a single CREATE PROCEDURE statement
    input := `CREATE PROCEDURE test_proc(IN id INT, OUT result VARCHAR(100))
BEGIN
    DECLARE temp VARCHAR(50);
    SET temp = 'Hello';
    
    IF id > 0 THEN
        SET result = CONCAT(temp, ' World');
    ELSE
        SET result = temp;
    END IF;
END`

    ast, err := sqlflowparser.Parse("", []byte(input))
    if err != nil {
        log.Fatal(err)
    }

    if proc, ok := ast.(*sqlflowparser.CreateProcedureStmt); ok {
        fmt.Printf("Procedure: %s\n", proc.Name)
        fmt.Printf("Parameters: %d\n", len(proc.Parameters))
        fmt.Printf("Body statements: %d\n", len(proc.Body))
    }
}
```

### Multi-Statement Files (with sqlsplitter)

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/eiiches/mysql-coverage/pkg/sqlsplitter"
    "github.com/eiiches/mysql-coverage/pkg/sqlflowparser"
)

func main() {
    // Read SQL file with multiple statements
    content, err := os.ReadFile("procedures.sql")
    if err != nil {
        log.Fatal(err)
    }

    // First, split into individual statements
    parser := sqlsplitter.NewParser(content)
    statements, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    // Then parse each SQL statement into AST
    for i, stmt := range statements {
        if stmt.Type == "SQL" {
            ast, err := sqlflowparser.Parse("", []byte(stmt.Text))
            if err != nil {
                fmt.Printf("Failed to parse statement %d: %v\n", i+1, err)
                continue
            }
            
            switch s := ast.(type) {
            case *sqlflowparser.CreateProcedureStmt:
                fmt.Printf("Found procedure: %s\n", s.Name)
            case *sqlflowparser.CreateFunctionStmt:
                fmt.Printf("Found function: %s\n", s.Name)
            }
        }
    }
}
```

## API

### Functions

#### Parse
```go
func Parse(filename string, b []byte, opts ...Option) (any, error)
```
Parses MySQL SQL input and returns an AST node. The `filename` parameter is used for error reporting.

### Types

#### AST Interface
```go
type AST interface {
    astNode()
    GetPosition() Position
}
```
Base interface for all AST nodes. Every AST node provides position information.

#### StatementAST Interface
```go
type StatementAST interface {
    AST
    statementNode()
}
```
Interface for all statement nodes.

#### ExpressionAST Interface
```go
type ExpressionAST interface {
    AST
    expressionNode()
}
```
Interface for all expression nodes.

### Statement Types

#### CreateProcedureStmt
```go
type CreateProcedureStmt struct {
    BaseStatement
    Name       string
    Parameters []Parameter
    Body       []StatementAST
}
```

#### CreateFunctionStmt
```go
type CreateFunctionStmt struct {
    BaseStatement
    Name       string
    Parameters []Parameter
    ReturnType string
    Body       []StatementAST
}
```

#### IfStmt
```go
type IfStmt struct {
    BaseStatement
    Condition string
    Then      []StatementAST
    ElseIfs   []ElseIfClause
    Else      []StatementAST
}
```

#### WhileStmt
```go
type WhileStmt struct {
    BaseStatement
    Condition string
    Body      []StatementAST
}
```

#### SetVariableStmt
```go
type SetVariableStmt struct {
    BaseStatement
    Assignments []VariableAssignment
}

type VariableAssignment struct {
    VariableRef   string // Variable name (@var, @@var, local_var)
    ScopeKeyword  string // GLOBAL, SESSION, PERSIST, PERSIST_ONLY
    Value         string // Assignment expression
}
```

#### Parameter
```go
type Parameter struct {
    Name string
    Type string
    Mode string // IN, OUT, INOUT
}
```

## Output

The parser returns different AST node types based on the input:

- **CreateProcedureStmt**: For `CREATE PROCEDURE` statements
- **CreateFunctionStmt**: For `CREATE FUNCTION` statements
- **IfStmt**: For `IF` statements
- **WhileStmt**: For `WHILE` statements
- **LoopStmt**: For `LOOP` statements
- **RepeatStmt**: For `REPEAT` statements
- **CaseStmt**: For `CASE` statements
- **BeginStmt**: For `BEGIN...END` blocks
- **DeclareStmt**: For `DECLARE` statements
- **SetVariableStmt**: For `SET` variable assignment statements
- **LeaveStmt**: For `LEAVE` statements
- **IterateStmt**: For `ITERATE` statements
- **ReturnStmt**: For `RETURN` statements
- **GenericStmt**: For other SQL statements (including `SIGNAL`, `SELECT`, `CALL`, etc.)

### Expression Handling

Expressions (conditions, values) are represented as strings with the original text preserved, allowing for simplified parsing while maintaining the ability to access the original expression content.

## Examples

### CREATE PROCEDURE
```sql
CREATE PROCEDURE test_proc(IN id INT, OUT name VARCHAR(50))
BEGIN
    DECLARE temp VARCHAR(100);
    SET temp = 'Hello';
    SELECT temp;
END
```

Output:
```
CreateProcedureStmt{
    Name: "test_proc",
    Parameters: []Parameter{
        {Name: "id", Type: "INT", Mode: "IN"},
        {Name: "name", Type: "VARCHAR(50)", Mode: "OUT"}
    },
    Body: []StatementAST{
        DeclareStmt{...},
        SetVariableStmt{Assignments: []VariableAssignment{{VariableRef: "temp", Value: "'Hello'"}}},
        GenericStmt{Text: "SELECT temp"}
    }
}
```

### CREATE FUNCTION
```sql
CREATE FUNCTION calc_tax(amount DECIMAL(10,2))
RETURNS DECIMAL(10,2)
BEGIN
    RETURN amount * 0.1;
END
```

Output:
```
CreateFunctionStmt{
    Name: "calc_tax",
    Parameters: []Parameter{
        {Name: "amount", Type: "DECIMAL(10,2)", Mode: ""}
    },
    ReturnType: "DECIMAL(10,2)",
    Body: []StatementAST{
        ReturnStmt{...}
    }
}
```

### IF Statement
```sql
IF x > 0 THEN
    SET result = 'positive';
ELSEIF x < 0 THEN
    SET result = 'negative';
ELSE
    SET result = 'zero';
END IF
```

Output:
```
IfStmt{
    Condition: "x > 0",
    Then: []StatementAST{SetVariableStmt{Assignments: []VariableAssignment{{VariableRef: "result", Value: "'positive'"}}}},
    ElseIfs: []ElseIfClause{
        {Condition: "x < 0", Then: []StatementAST{SetVariableStmt{Assignments: []VariableAssignment{{VariableRef: "result", Value: "'negative'"}}}}}
    },
    Else: []StatementAST{SetVariableStmt{Assignments: []VariableAssignment{{VariableRef: "result", Value: "'zero'"}}}}
}
```

### WHILE Loop
```sql
WHILE counter < 10 DO
    SET counter = counter + 1;
    SELECT counter;
END WHILE
```

Output:
```
WhileStmt{
    Condition: "counter < 10",
    Body: []StatementAST{
        SetVariableStmt{Assignments: []VariableAssignment{{VariableRef: "counter", Value: "counter + 1"}}},
        GenericStmt{Text: "SELECT counter"}
    }
}
```

### SET Variable Statements
```sql
SET @user_var = 'hello';
SET @@SESSION.sql_mode = 'STRICT_TRANS_TABLES';
SET GLOBAL max_connections = 200;
SET result = x + y, temp = result * 2;
```

Output:
```
SetVariableStmt{
    Assignments: []VariableAssignment{
        {VariableRef: "@user_var", Value: "'hello'"}
    }
}

SetVariableStmt{
    Assignments: []VariableAssignment{
        {VariableRef: "@@SESSION.sql_mode", Value: "'STRICT_TRANS_TABLES'"}
    }
}

SetVariableStmt{
    Assignments: []VariableAssignment{
        {VariableRef: "max_connections", ScopeKeyword: "GLOBAL", Value: "200"}
    }
}

SetVariableStmt{
    Assignments: []VariableAssignment{
        {VariableRef: "result", Value: "x + y"},
        {VariableRef: "temp", Value: "result * 2"}
    }
}
```

## Testing

The package includes comprehensive unit tests using Gomega:

```bash
go test ./pkg/sqlflowparser -v
```

## Design

This parser is built using PEG (Parsing Expression Grammar) via the `pigeon` parser generator. It focuses on:

- **Single statement parsing**: Designed to parse one MySQL statement at a time
- **Structural parsing**: Extracts procedure/function definitions and control flow
- **Statement boundaries**: Identifies nested statement blocks within procedures/functions
- **Parameter extraction**: Parses function/procedure signatures
- **Control flow analysis**: Supports all MySQL control flow constructs

It does **not** provide:
- **Multi-statement parsing**: Cannot handle entire SQL files with multiple statements
- **Delimiter handling**: No support for DELIMITER statements or custom delimiters
- **Detailed expression parsing**: Expressions are kept as text for simplicity
- **Full SQL validation**: Focuses on structure rather than semantic correctness
- **Complex type analysis**: Types are preserved as strings

### Architecture

The typical workflow combines two packages:

1. **`@pkg/sqlsplitter`**: First pass - splits SQL files into individual statements
2. **`@pkg/sqlflowparser`**: Second pass - parses each statement into structured AST

This separation of concerns allows:
- **sqlsplitter** to handle file-level concerns (delimiters, comments, string literals)
- **sqlflowparser** to focus on statement-level syntax analysis

This makes it ideal for tools that need to analyze stored procedure structure, instrument code, or perform static analysis on MySQL stored procedures and functions.

## Grammar

The parser is generated from a PEG grammar file (`mysql_ast.peg`) that defines the MySQL syntax patterns. The grammar focuses on the subset of MySQL relevant for stored procedures and functions.