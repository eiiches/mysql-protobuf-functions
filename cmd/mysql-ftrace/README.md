# MySQL Function Trace Tool (mysql-ftrace)

A tool for instrumenting MySQL stored procedures and functions to generate function call traces with arguments and return values. This tool uses the same AST-based instrumentation technique as mysql-coverage but focuses on function call tracing instead of line coverage.

## Architecture Overview

The tool uses a **two-pass AST-based parsing approach** that mirrors how MySQL client and server handle SQL statements:

### 1. Statement Splitting (First Pass)
- **Input**: Complete SQL file with DELIMITER statements and multiple procedures/functions
- **Output**: Individual SQL statements split by current delimiter
- **Implementation**: Uses the existing `sqlsplitter` package from the mysql-coverage tool
- **Features**:
  - Arbitrary-length delimiter support (e.g., `DELIMITER ENDOFSTATEMENT`)
  - Proper string literal parsing (single quotes, double quotes, backticks)
  - Comment handling (line comments `--`, `#` and block comments `/* */`)
  - MySQL client-compatible delimiter recognition

### 2. AST Parsing (Second Pass)
- **Input**: Individual SQL statements (CREATE FUNCTION/PROCEDURE)
- **Output**: Abstract Syntax Tree (AST) for stored procedures and functions
- **Implementation**: Uses the existing `sqlflowparser` package with PEG (Parsing Expression Grammar) parser
- **Features**:
  - Control flow statement parsing (IF, WHILE, LOOP, REPEAT, CASE)
  - Function/procedure signature extraction (parameters, return types)
  - Statement boundary identification for instrumentation
  - Label preservation for LEAVE/ITERATE compatibility

### 3. Function Call Instrumentation
- **Input**: SQL file
- **Output**: Instrumented SQL with function call tracing
- **Implementation**: AST-based reconstruction for syntactic correctness
- **Features**:
  - Function entry logging with parameter values
  - Function exit logging with return values
  - Call depth tracking for nested function calls
  - JSON-formatted argument capture

## Features

- **AST-Based Accuracy**: Uses structured parsing rather than fragile regex patterns
- **MySQL Compatibility**: Follows MySQL client/server parsing behavior
- **Function Call Tracing**: Logs entry/exit with arguments and return values
- **Call Depth Tracking**: Handles nested function calls with proper indentation
- **Multiple Report Formats**: Text, JSON, and flamegraph (planned) output formats
- **Argument Capture**: Captures parameter values in JSON format for analysis
- **Syntactic Correctness**: Reconstructs valid SQL from parsed AST

## Installation

```bash
go build -o mysql-ftrace cmd/mysql-ftrace/main.go
```

## Usage

### 1. Initialize Function Tracing Schema

Use the `init` subcommand to set up the function tracing schema in your MySQL database:

```bash
# Initialize tracing schema
./mysql-ftrace init --database "user:password@tcp(localhost:3306)/database"
```

This creates the necessary database schema for function call tracing:
- `__FtraceEvent` table for storing trace events
- `__record_ftrace_entry` procedure for logging function entries
- `__record_ftrace_exit` procedure for logging function exits
- Call depth management functions

### 2. Instrument Your SQL Files

The `instrument` subcommand takes SQL files with stored procedures/functions and adds function call tracing. By default, it uses the naming convention `{original}.ftraced`:

```bash
# Instrument a single file (creates protobuf.sql.ftraced)
./mysql-ftrace instrument protobuf.sql

# Instrument multiple files at once
./mysql-ftrace instrument protobuf.sql protobuf-accessors.sql protobuf-descriptor.sql

# Instrument all SQL files using wildcards
./mysql-ftrace instrument *.sql
```

**What it does:**
- Adds `CALL __record_ftrace_entry(filename, function_name, arguments_json);` at function entry
- Adds `CALL __record_ftrace_exit(filename, function_name, return_value);` at function exits
- Only instruments function/procedure bodies (`BEGIN`...`END`)
- Preserves MySQL labels (e.g., `l1:`) for LEAVE/ITERATE compatibility
- Captures parameter values in JSON format

### Instrumentation Behavior

#### Function Entry Call Format
```sql
CALL __record_ftrace_entry('<filename>', '<function_name>', '<arguments_json>');
```

#### Function Exit Call Format
```sql
CALL __record_ftrace_exit('<filename>', '<function_name>', '<return_value>');
```

#### Instrumented Elements
- **Function Entry**: Right after BEGIN in function/procedure body
- **Function Exit**: Before each RETURN statement and at implicit exits
- **Argument Capture**: All IN, OUT, INOUT parameters in JSON format
- **Return Value Capture**: Function return values and procedure completion

#### Example Transformation

```sql
-- Original
DELIMITER $$
CREATE FUNCTION calc_tax(amount DECIMAL(10,2)) RETURNS DECIMAL(10,2) DETERMINISTIC
BEGIN
    DECLARE tax DECIMAL(10,2);
    SET tax = amount * 0.1;
    IF tax > 100 THEN
        RETURN 100;
    END IF;
    RETURN tax;
END $$
DELIMITER ;

-- Instrumented
DELIMITER $$
CREATE FUNCTION calc_tax(amount DECIMAL(10,2)) RETURNS DECIMAL(10,2) DETERMINISTIC
BEGIN
    CALL __record_ftrace_entry('input.sql', 'calc_tax', CONCAT('{', "amount": ', COALESCE(amount, 'NULL'), '", '}'));
    DECLARE tax DECIMAL(10,2);
    SET tax = amount * 0.1;
    IF tax > 100 THEN
        BEGIN
            CALL __record_ftrace_exit('input.sql', 'calc_tax', COALESCE(100, 'NULL'));
            RETURN 100;
        END
    END IF;
    BEGIN
        CALL __record_ftrace_exit('input.sql', 'calc_tax', COALESCE(tax, 'NULL'));
        RETURN tax;
    END
END $$
DELIMITER ;
```

### 3. Load Instrumented Code and Run Tests

```bash
# Load the instrumented SQL into your database
mysql -h localhost -u user -p database < instrumented-protobuf.sql.ftraced

# Run your tests or application (this will populate trace data)
go test ./tests -database "user:password@tcp(localhost:3306)/database"
```

### 4. Generate Function Trace Reports

The `report` subcommand generates function call trace reports in various formats:

```bash
# Generate text report (default, shows call flow with indentation)
./mysql-ftrace report --database "user:password@tcp(localhost:3306)/database" --format text

# Generate JSON report (structured data for analysis)
./mysql-ftrace report --database "user:password@tcp(localhost:3306)/database" --format json --output trace.json

# Generate flamegraph data (planned feature)
./mysql-ftrace report --database "user:password@tcp(localhost:3306)/database" --format flamegraph --output trace.folded
```


## Command Reference

### init

Initializes the database with function tracing schema and clears any existing trace data.

```bash
./mysql-ftrace init --database CONNECTION_STRING
```

**Options:**
- `--database string`: Database connection string (required)

**Creates:**
- `__FtraceEvent` table for trace storage (recreates if exists)
- `__record_ftrace_entry` procedure
- `__record_ftrace_exit` procedure
- Call depth management functions

**Note:** Running `init` clears existing trace data by recreating the table.

### instrument

Instruments SQL files with function call tracing using AST-based parsing.

```bash
./mysql-ftrace instrument file1.sql [file2.sql ...]
```

**Examples:**
```bash
# Basic usage (creates functions.sql.ftraced)
./mysql-ftrace instrument functions.sql

# Multiple files
./mysql-ftrace instrument file1.sql file2.sql file3.sql

# Using wildcards
./mysql-ftrace instrument *.sql
```

### report

Generates function call trace reports from the trace database.

```bash
./mysql-ftrace report --database CONNECTION_STRING [options]
```

**Options:**
- `--database string`: Database connection string (required)
- `--format string`: Output format: text, json, flamegraph (default: text)
- `--output string`: Output file (default: stdout)
- `--connection-id int`: Filter reports by specific connection ID (default: show all connections)

**Examples:**
```bash
# Text report with call flow visualization
./mysql-ftrace report --database "root@tcp(127.0.0.1:3306)/test" --format text

# JSON report for programmatic analysis
./mysql-ftrace report --database "root@tcp(127.0.0.1:3306)/test" --format json --output trace.json

# Filter trace report by specific connection ID
./mysql-ftrace report --database "root@tcp(127.0.0.1:3306)/test" --connection-id 42
```


## Report Formats

### Text Format

Shows function calls with indentation based on call depth, grouped by connection ID:

```
MySQL Function Call Trace Report
================================

=== Connection ID: 42 ===
[15:04:05.123] -> test_add({"a": "5", "b": "10"})
[15:04:05.124] <- test_add = 15

=== Connection ID: 43 ===
[15:04:05.125] -> complex_function({"x": "100"})
  [15:04:05.126] -> helper_function({"value": "100"})
  [15:04:05.127] <- helper_function = 200
[15:04:05.128] <- complex_function = 200
```

### JSON Format

Structured data suitable for analysis and visualization:

```json
[
  {
    "id": 1,
    "connection_id": 42,
    "filename": "test.sql",
    "function_name": "test_add",
    "call_type": "entry",
    "arguments": "{\"a\": \"5\", \"b\": \"10\"}",
    "return_value": "",
    "call_depth": 1,
    "timestamp": "2023-12-01T15:04:05.123456Z"
  },
  {
    "id": 2,
    "connection_id": 42,
    "filename": "test.sql",
    "function_name": "test_add",
    "call_type": "exit",
    "arguments": "",
    "return_value": "15",
    "call_depth": 1,
    "timestamp": "2023-12-01T15:04:05.124789Z"
  }
]
```

## Database Schema

The tool creates the following database objects:

### Tables

**__FtraceEvent**
- `id`: Auto-increment primary key
- `connection_id`: MySQL connection ID (from CONNECTION_ID())
- `filename`: Source file name
- `function_name`: Function or procedure name
- `call_type`: 'entry' or 'exit'
- `arguments`: JSON-formatted function arguments (for entry events)
- `return_value`: Function return value (for exit events)
- `call_depth`: Current call depth for nested calls
- `timestamp`: High-precision timestamp

### Procedures

**__record_ftrace_entry(filename, function_name, arguments)**
- Records function entry with arguments and CONNECTION_ID()
- Increments call depth

**__record_ftrace_exit(filename, function_name, return_value)**
- Records function exit with return value and CONNECTION_ID()
- Decrements call depth

**Call Depth Management**
- `__get_call_depth()`: Returns current call depth
- `__increment_call_depth()`: Increments session call depth
- `__decrement_call_depth()`: Decrements session call depth

## Real-World Example

Here's a complete workflow for tracing MySQL protobuf functions:

```bash
# 1. Build the tool
go build -o mysql-ftrace cmd/mysql-ftrace/main.go

# 2. Initialize tracing schema
./mysql-ftrace init --database "root@tcp(127.0.0.1:3306)/test"

# 3. Instrument the SQL functions (creates protobuf.sql.ftraced)
./mysql-ftrace instrument protobuf.sql

# 4. Load instrumented functions
mysql -h 127.0.0.1 -u root test < protobuf.sql.ftraced

# 5. Run tests to generate trace data
go test ./tests -run TestWireJsonGetField -database "root@tcp(127.0.0.1:3306)/test"

# 6. Generate trace report
./mysql-ftrace report --database "root@tcp(127.0.0.1:3306)/test" --format text

# 7. Generate JSON report for analysis
./mysql-ftrace report --database "root@tcp(127.0.0.1:3306)/test" --format json --output trace.json

# 8. Run init again to clear trace data for next run (optional)
./mysql-ftrace init --database "root@tcp(127.0.0.1:3306)/test"
```

## Technical Details

### AST-Based Instrumentation

The tool uses the same sophisticated AST-based parsing as mysql-coverage:

1. **Two-Pass Parsing**: Statement splitting → AST generation → Instrumentation
2. **Function/Procedure Detection**: Tracks `CREATE FUNCTION/PROCEDURE` declarations using AST nodes
3. **Control Flow Analysis**: Uses AST to understand statement boundaries and nesting
4. **Label Preservation**: Maintains MySQL labels (`l1:`, `proc:`) for LEAVE/ITERATE compatibility
5. **Syntactic Correctness**: Reconstructs valid SQL from modified AST

### Design Principles

1. **MySQL Compatibility**: Follows MySQL client/server parsing behavior
2. **AST-Based Accuracy**: Uses structured parsing rather than text manipulation
3. **Function-Focused**: Instruments function boundaries rather than individual statements
4. **Call Depth Tracking**: Handles nested function calls correctly
5. **Argument Preservation**: Captures parameter values for analysis

### Performance Considerations

- Each function call adds entry/exit overhead with database logging
- The trace table uses InnoDB for better concurrent access
- High-precision timestamps for performance analysis
- Consider clearing trace data between test runs to manage table size
- AST parsing is accurate but slightly slower than regex-based approaches
- **Connection ID tracking**: Enables proper isolation of concurrent database sessions
- Use `--connection-id` filter for focused analysis of specific database connections

## Comparison with mysql-coverage

| Feature | mysql-coverage | mysql-ftrace |
|---------|---------------|--------------|
| **Purpose** | Line coverage analysis | Function call tracing |
| **Instrumentation** | Before executable statements | Function entry/exit points |
| **Data Captured** | Execution counts per line | Function calls with arguments |
| **Reports** | LCOV format, HTML coverage | Text, JSON, flamegraph traces |
| **Use Case** | Testing coverage analysis | Performance profiling, debugging |
| **Database Schema** | `__CoverageEvent` | `__FtraceEvent` |

## Troubleshooting

### Common Issues

1. **Syntax Errors After Instrumentation**
   - The AST-based approach should eliminate most syntax errors
   - Check that your SQL uses valid MySQL syntax
   - Ensure proper `DELIMITER` usage
   - Verify function/procedure syntax

2. **Missing Trace Data**
   - Ensure `__record_ftrace_entry` and `__record_ftrace_exit` procedures exist
   - Check database permissions
   - Verify instrumented functions were loaded

3. **Incorrect Call Depth**
   - Clear session variables: `SET @__ftrace_call_depth = 0`
   - Restart MySQL session if call depth becomes inconsistent

4. **Large Trace Tables**
   - Use the `init` command between test runs to clear data
   - Consider archiving old trace data before running `init`
   - Monitor disk space usage

### Debugging

Enable SQL logging to see instrumentation calls:
```sql
SET GLOBAL general_log = 'ON';
SET GLOBAL log_output = 'TABLE';
SELECT * FROM mysql.general_log WHERE command_type = 'Query' AND argument LIKE '%__record_ftrace%';
```

## Internal Architecture

The tool is built on these core packages:

- **`internal/mysql/sqlsplitter`**: MySQL-compatible statement splitting with delimiter support
- **`internal/mysql/sqlflowparser`**: MySQL AST parser for stored functions and procedures using PEG grammar
- **`internal/mysql/sqlftrace`**: Function tracing instrumentation engine using AST reconstruction

This AST-based approach provides superior accuracy and reliability compared to regex-based text manipulation approaches.
