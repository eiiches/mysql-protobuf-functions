**WARNING: The contents in this directory (including this README) are mostly AI-generated with a little human review. Don't expect this to work for you.**

# MySQL SQL Coverage Instrumentation Tool

A tool for instrumenting MySQL stored procedures and functions to generate code coverage reports. This tool helps analyze which parts of your MySQL stored programs are executed during testing.

## Architecture Overview

The tool uses a **two-pass AST-based parsing approach** that mirrors how MySQL client and server handle SQL statements:

### 1. Statement Splitting (First Pass)
- **Input**: Complete SQL file with DELIMITER statements and multiple procedures/functions
- **Output**: Individual SQL statements split by current delimiter
- **Implementation**: Manual recursive descent parser optimized for MySQL client compatibility
- **Features**:
  - Arbitrary-length delimiter support (e.g., `DELIMITER ENDOFSTATEMENT`)
  - Proper string literal parsing (single quotes, double quotes, backticks)
  - Comment handling (line comments `--`, `#` and block comments `/* */`)
  - MySQL client-compatible delimiter recognition

### 2. AST Parsing (Second Pass)
- **Input**: Individual SQL statements (CREATE FUNCTION/PROCEDURE)
- **Output**: Abstract Syntax Tree (AST) for stored procedures and functions
- **Implementation**: PEG (Parsing Expression Grammar) parser via Pigeon
- **Features**:
  - Control flow statement parsing (IF, WHILE, LOOP, REPEAT, CASE)
  - Function/procedure signature extraction (parameters, return types)
  - Statement boundary identification for instrumentation
  - Label preservation for LEAVE/ITERATE compatibility

### 3. Coverage Instrumentation
- **Input**: SQL file
- **Output**: Instrumented SQL with coverage tracking calls
- **Implementation**: AST-based reconstruction for syntactic correctness
- **Features**:
  - Smart statement detection (only instruments executable statements)
  - Line number tracking for accurate coverage reporting
  - Preserves MySQL labels and control flow semantics

## Features

- **AST-Based Accuracy**: Uses structured parsing rather than fragile regex patterns
- **MySQL Compatibility**: Follows MySQL client/server parsing behavior  
- **Generate LCOV format** coverage reports
- **100% compatible** with standard coverage tools like `genhtml`
- **Preserves line numbers** for accurate coverage mapping
- **Syntactic Correctness**: Reconstructs valid SQL from parsed AST
- **Comprehensive Coverage**: Handles complex nested control structures and labels

## Installation

```bash
go build -o mysql-coverage cmd/mysql-coverage/main.go
```

## Usage

### 1. Initialize Coverage Schema

Use the `init` subcommand to set up the coverage tracking schema in your MySQL database:

```bash
# Initialize coverage schema
./mysql-coverage init --database "user:password@tcp(localhost:3306)/database"
```

This creates the necessary database schema for coverage tracking.

### 2. Instrument Your SQL Files

The `instrument` subcommand takes SQL files with stored procedures/functions and adds coverage tracking calls. By default, it uses the naming convention `{original}.instrumented`:

```bash
# Instrument a single file (creates protobuf.sql.instrumented)
./mysql-coverage instrument protobuf.sql

# Instrument multiple files at once
./mysql-coverage instrument protobuf.sql protobuf-accessors.sql protobuf-descriptor.sql

# Instrument all SQL files using wildcards
./mysql-coverage instrument *.sql

# Specify custom output directory
./mysql-coverage instrument --output instrumented/ protobuf.sql other.sql

# Or use stdin/stdout
cat protobuf.sql | ./mysql-coverage instrument > instrumented-protobuf.sql
```

**What it does:**
- Adds `CALL __record_coverage(filename, function_name, line_number);` before each executable statement
- Only instruments statements inside function/procedure bodies (`BEGIN`...`END`)
- Preserves line numbers for accurate coverage mapping
- Preserves MySQL labels (e.g., `l1:`) for LEAVE/ITERATE compatibility
- Skips `DECLARE`, cursor definitions, and other non-executable statements

### Instrumentation Behavior

#### Coverage Call Format
```sql
CALL __record_coverage('<filename>', '<function_or_procedure_name>', <line_number>);
```

#### Instrumented Statements
- Control flow: `IF`, `WHILE`, `LOOP`, `REPEAT`, `CASE`
- Data manipulation: `SET`, `SELECT`, `INSERT`, `UPDATE`, `DELETE`
- Procedure calls: `CALL`
- Flow control: `RETURN`, `LEAVE`, `ITERATE`
- Error handling: `SIGNAL`
- Generic SQL: Any other executable statement

#### Non-Instrumented Elements
- `DECLARE` statements
- Comments (`--`, `#`, `/* */`)
- `BEGIN` and `END` blocks themselves (but statements inside are instrumented)
- Labels (e.g., `l1:`) - preserved but not instrumented
- Conditional clauses (`WHEN`, `ELSE`, `ELSEIF`, `UNTIL`) themselves (but statements inside are instrumented)
- Empty lines and whitespace

**Example transformation:**
```sql
-- Original
DELIMITER $$
CREATE FUNCTION calc_tax(amount DECIMAL(10,2)) RETURNS DECIMAL(10,2) DETERMINISTIC
BEGIN
    DECLARE tax DECIMAL(10,2);
    SET tax = amount * 0.1;
    l1: IF tax > 100 THEN
        SET tax = 100;
        LEAVE l1;
    END IF;
    RETURN tax;
END $$
DELIMITER ;

-- Instrumented
DELIMITER $$
CREATE FUNCTION calc_tax(amount DECIMAL(10,2)) RETURNS DECIMAL(10,2) DETERMINISTIC
BEGIN
    DECLARE tax DECIMAL(10,2);
    CALL __record_coverage('input.sql', 'calc_tax', 5);
    SET tax = amount * 0.1;
    CALL __record_coverage('input.sql', 'calc_tax', 6);
    l1: IF tax > 100 THEN
        CALL __record_coverage('input.sql', 'calc_tax', 7);
        SET tax = 100;
        CALL __record_coverage('input.sql', 'calc_tax', 8);
        LEAVE l1;
    END IF;
    CALL __record_coverage('input.sql', 'calc_tax', 10);
    RETURN tax;
END $$
DELIMITER ;
```

### 3. Load Instrumented Code and Run Tests

```bash
# Load the instrumented SQL into your database
mysql -h localhost -u user -p database < instrumented-protobuf.sql

# Run your tests (this will populate coverage data)
go test ./tests -database "user:password@tcp(localhost:3306)/database"
```

### 4. Generate Coverage Reports

The `lcov` subcommand generates standard LCOV format coverage reports. It automatically detects instrumented files using the `*.sql.instrumented` pattern:

```bash
# Generate LCOV report (auto-detects instrumented files)
./mysql-coverage lcov --database "user:password@tcp(localhost:3306)/database" --output coverage.lcov

# Explicitly specify instrumented files
./mysql-coverage lcov --database "user:password@tcp(localhost:3306)/database" --instrumented-file protobuf.sql.instrumented --instrumented-file other.sql.instrumented --output coverage.lcov

# Generate HTML report using genhtml
genhtml coverage.lcov --output-directory coverage-html --title "MySQL Coverage Report"
```

## Command Reference

### init

Initializes the database with coverage tracking schema.

```bash
./mysql-coverage init --database CONNECTION_STRING
```

**Options:**
- `--database string`: Database connection string (required)

**Examples:**
```bash
# Initialize coverage schema
./mysql-coverage init --database "root@tcp(127.0.0.1:3306)/test"
```

### instrument

Instruments SQL files with coverage tracking calls using AST-based parsing.

```bash
./mysql-coverage instrument [options] [file1.sql file2.sql ...]
```

**Options:**
- `--output string`: Output directory (only used with multiple files)

**Examples:**
```bash
# Basic usage (creates functions.sql.instrumented)
./mysql-coverage instrument functions.sql

# Multiple files (creates file1.sql.instrumented, file2.sql.instrumented)
./mysql-coverage instrument file1.sql file2.sql file3.sql

# Using wildcards
./mysql-coverage instrument *.sql

# Custom output directory
./mysql-coverage instrument --output instrumented/ *.sql

# Using pipes
cat file.sql | ./mysql-coverage instrument > instrumented.sql
```

### lcov

Generates LCOV format coverage report from the coverage database.

```bash
./mysql-coverage lcov --database CONNECTION_STRING [options]
```

**Options:**
- `--database string`: Database connection string (required)
- `--output string`: Output file (default: stdout)
- `--instrumented-file strings`: Path(s) to instrumented SQL file(s) (auto-detected if not specified)

**Examples:**
```bash
# Generate LCOV file (auto-detects *.sql.instrumented files)
./mysql-coverage lcov --database "root@tcp(127.0.0.1:3306)/test" --output coverage.lcov

# Explicitly specify instrumented files
./mysql-coverage lcov --database "root@tcp(127.0.0.1:3306)/test" --instrumented-file protobuf.sql.instrumented --output coverage.lcov

# Direct to genhtml
./mysql-coverage lcov --database "root@tcp(127.0.0.1:3306)/test" | genhtml - --output-directory html-report
```

## Coverage Report Formats

### LCOV Format

The tool generates standard LCOV format files that include:

- **TN**: Test name
- **SF**: Source file path
- **FN**: Function name and line number
- **FNDA**: Function hit count
- **FNF/FNH**: Functions found/hit summary
- **DA**: Line hit count
- **LF/LH**: Lines found/hit summary

### HTML Reports

Use `genhtml` to convert LCOV files to interactive HTML reports:

```bash
genhtml coverage.lcov --output-directory coverage-html --title "MySQL Functions Coverage"
```

**HTML report features:**
- Interactive source code browser
- Color-coded coverage visualization
- Function and line coverage statistics
- Drill-down from summary to source level

## Real-World Example

Here's a complete workflow for analyzing MySQL protobuf functions:

```bash
# 1. Build the tool
go build -o mysql-coverage cmd/mysql-coverage/main.go

# 2. Initialize coverage schema
./mysql-coverage init --database "root@tcp(127.0.0.1:3306)/test"

# 3. Instrument the SQL functions (creates protobuf.sql.instrumented)
./mysql-coverage instrument protobuf.sql

# 4. Load instrumented functions
mysql -h 127.0.0.1 -u root test < protobuf.sql.instrumented

# 5. Run tests to generate coverage data
go test ./tests -run TestRandomizedWireJsonGetField -database "root@tcp(127.0.0.1:3306)/test"
go test ./tests -run TestRandomizedWireJsonHasField -database "root@tcp(127.0.0.1:3306)/test"

# 6. Generate coverage report
./mysql-coverage lcov --database "root@tcp(127.0.0.1:3306)/test" --output coverage.lcov

# 7. Generate HTML report
genhtml coverage.lcov --output-directory coverage-html --title "MySQL Protobuf Functions Coverage"

# 8. View results
open coverage-html/index.html
```

**Results:**
- Functions: 100.0% (26/26)
- Lines: 100.0% (217/217)
- Coverage events: 341,138

## Technical Details

### AST-Based Instrumentation

The tool uses sophisticated AST-based parsing to identify instrumentable statements:

1. **Two-Pass Parsing**: Statement splitting → AST generation → Instrumentation
2. **Function/Procedure Detection**: Tracks `CREATE FUNCTION/PROCEDURE` declarations using AST nodes
3. **Control Flow Analysis**: Uses AST to understand statement boundaries and nesting
4. **Label Preservation**: Maintains MySQL labels (`l1:`, `proc:`) for LEAVE/ITERATE compatibility
5. **Syntactic Correctness**: Reconstructs valid SQL from modified AST

### Design Principles

1. **MySQL Compatibility**: Follows MySQL client/server parsing behavior
2. **AST-Based Accuracy**: Uses structured parsing rather than text manipulation
3. **Minimal Overhead**: Only instruments executable statements
4. **Syntactic Correctness**: Reconstructs valid SQL from parsed AST
5. **Comprehensive Coverage**: Handles complex nested control structures

### Performance Considerations

- Each instrumented statement adds a small overhead
- The coverage table is automatically indexed for faster reporting
- AST parsing is more accurate but slightly slower than regex-based approaches
- Consider clearing old coverage data between test runs using the database tools

## Troubleshooting

### Common Issues

1. **Syntax Errors After Instrumentation**
   - The AST-based approach should eliminate most syntax errors
   - Check that your SQL uses valid MySQL syntax
   - Ensure proper `DELIMITER` usage
   - Verify function/procedure syntax

2. **Missing Coverage Data**
   - Ensure `__record_coverage` procedure exists
   - Check database permissions
   - Verify instrumented functions were loaded

3. **Incorrect Line Numbers**
   - Make sure not to modify instrumented files manually
   - Re-instrument from original source if needed

4. **Label-Related Errors**
   - The AST parser preserves labels correctly
   - If you see "ITERATE with no matching label", re-instrument the file
   - Check for complex label usage patterns

### Debugging

Enable SQL logging to see instrumentation calls:
```sql
SET GLOBAL general_log = 'ON';
SET GLOBAL log_output = 'TABLE';
SELECT * FROM mysql.general_log WHERE command_type = 'Query' AND argument LIKE '%__record_coverage%';
```

## Internal Architecture

The tool is built on these core packages:

- **`internal/mysql/sqlsplitter`**: MySQL-compatible statement splitting with delimiter support
- **`internal/mysql/sqlflowparser`**: MySQL AST parser for stored functions and procedures using PEG grammar
- **`internal/mysql/sqlinstrument`**: Coverage instrumentation engine using AST reconstruction

This AST-based approach provides superior accuracy and reliability compared to regex-based text manipulation approaches.
