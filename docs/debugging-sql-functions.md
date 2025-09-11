# Debugging MySQL Protobuf Functions

This document describes techniques for debugging issues in the MySQL protobuf functions, particularly for locating where errors occur in complex SQL stored procedures.

## Code Coverage Instrumentation for Error Location

When you encounter MySQL errors during function execution but can't easily identify where in the code the error occurs, you can use the built-in code coverage instrumentation to pinpoint the exact location.

### Steps to Debug Using Instrumentation

1. **Load Instrumented Functions**
   ```bash
   make load-instrumented-files
   ```
   This command:
   - Instruments all SQL functions with coverage tracking calls
   - Creates a `__CoverageEvent` table to log function execution
   - Loads the instrumented versions into MySQL

2. **Execute the Failing Operation**
   Run the specific SQL statement that's causing the error. For example:
   ```sql
   SELECT pb_json_to_message('[descriptor_json]', '.MessageType', '{"field": "value"}', NULL, NULL);
   ```

3. **Check Coverage Events for Error Location**
   ```sql
   SELECT * FROM __CoverageEvent ORDER BY id DESC LIMIT 10;
   ```
   
   The output shows:
   - `filename`: Which SQL file contained the error
   - `function_name`: Which function was executing
   - `line_number`: **The line number in the original source file** where execution stopped
   - `timestamp`: When the execution occurred

4. **Locate the Error in Source Code**
   The `line_number` corresponds to the line in the original source file (e.g., `src/json-to-protobuf.sql`), not the instrumented build file. Use this line number to find the exact statement that caused the error.

### Example Debugging Session

**Problem**: Getting error `Truncated incorrect INTEGER value: 'false'` when processing JSON boolean values.

**Investigation**:
```bash
# 1. Load instrumented functions
make load-instrumented-files

# 2. Run the failing case
mysql> SELECT pb_json_to_message('[...]', '.Test', '{"boolField": false}', NULL, NULL);
ERROR 1292 (22007): Truncated incorrect INTEGER value: 'false'

# 3. Check where execution stopped
mysql> SELECT * FROM __CoverageEvent ORDER BY id DESC LIMIT 3;
+----+---------------------+---------------------------+-------------+---------------------+
| id | filename            | function_name             | line_number | timestamp           |
+----+---------------------+---------------------------+-------------+---------------------+
| 57 | build/protobuf-json.sql | _pb_is_proto3_default_value | 1451        | 2025-08-11 05:20:07 |
| 56 | build/protobuf-json.sql | _pb_is_proto3_default_value | 1435        | 2025-08-11 05:20:07 |
+----+---------------------+---------------------------+-------------+---------------------+

# 4. Check line 1451 in the original source
```

Looking at line 1451 in the original source revealed:
```sql
WHEN 8 THEN -- TYPE_BOOL
    RETURN CAST(json_value AS SIGNED) = 0;  -- This line caused the error
```

**Root Cause**: The `CAST(json_value AS SIGNED)` operation fails when `json_value` is a JSON boolean (`false`), because MySQL can't directly cast JSON boolean values to integers.

**Solution**: Changed to use proper JSON boolean comparison:
```sql
WHEN 8 THEN -- TYPE_BOOL
    RETURN json_value = CAST(false AS JSON);
```

### Benefits of This Approach

- **Precise Error Location**: Pinpoints the exact line where execution failed
- **No Code Modification**: Doesn't require adding debug prints or modifying logic
- **Historical Tracking**: Can review execution flow leading up to errors
- **Complex Function Debugging**: Particularly useful for nested function calls and stored procedures

### Limitations

- **Instrumentation Overhead**: Instrumented functions run slower due to coverage logging
- **Line Number Accuracy**: Line numbers correspond to the original source, not the instrumented version
- **Only Shows Execution Path**: Doesn't show variable values or intermediate states

### Alternative Debugging Techniques

For simpler cases, you can also use:

1. **Direct SQL Testing**: Test individual functions in isolation
2. **Error Message Analysis**: MySQL error messages often contain helpful context
3. **Step-by-step Function Calls**: Break complex operations into smaller parts
4. **Variable Inspection**: Use `SELECT` statements to check intermediate values

The instrumentation approach is most valuable when dealing with complex stored procedure chains where the error location isn't obvious from the error message alone.