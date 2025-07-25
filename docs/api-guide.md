# MySQL Protocol Buffers - API Guide

*Advanced technical guide with detailed examples and complex use cases*

## 📚 Navigation

- **👈 [Documentation Home](README.md)** - Choose your learning path
- **🎯 [User Guide](user-guide.md)** - Problem-focused examples and solutions
- **⚡ [Quick Reference](quick-reference.md)** - Function syntax cheat sheet  
- **📖 [Function Reference](function-reference.md)** - Complete API documentation

---

This guide provides detailed documentation and examples for using the MySQL Protocol Buffers functions library.

> **🚀 New to this library?** Consider starting with the [User Guide](user-guide.md) for practical, problem-focused examples first.

## Table of Contents

- [Getting Started](#getting-started)
- [Core Concepts](#core-concepts)
- [Field Operations](#field-operations)
- [Repeated Fields](#repeated-fields)
- [JSON Integration](#json-integration)
- [Advanced Topics](#advanced-topics)
- [Performance Considerations](#performance-considerations)
- [Troubleshooting](#troubleshooting)

## Getting Started

### Installation

1. Install core protobuf functions:
```sql
SOURCE protobuf.sql;
SOURCE protobuf-accessors.sql;
```

2. Install optional JSON features:
```sql
SOURCE protobuf-descriptor.sql;
SOURCE protobuf-json.sql;
```

### Basic Example

```sql
-- Create a table to store protobuf data
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    pb_data LONGBLOB NOT NULL
);

-- Create a new protobuf message
SET @user_msg = pb_message_new();
SET @user_msg = pb_message_set_string_field(@user_msg, 1, 'John Doe');
SET @user_msg = pb_message_set_int32_field(@user_msg, 2, 25);
SET @user_msg = pb_message_set_string_field(@user_msg, 3, 'john@example.com');

-- Insert into table
INSERT INTO users (pb_data) VALUES (@user_msg);

-- Read back the data
SELECT 
    pb_message_get_string_field(pb_data, 1, '') AS name,
    pb_message_get_int32_field(pb_data, 2, 0) AS age,
    pb_message_get_string_field(pb_data, 3, '') AS email
FROM users;
```

## Core Concepts

### Message Types

The library works with two main data representations:

1. **Binary Messages** (`LONGBLOB`): Standard protobuf binary format
2. **Wire JSON** (`JSON`): Intermediate JSON format for efficient multi-operation processing

#### When to Use Each Format

**Binary Messages (`pb_message_*` functions):**
- Single field operations
- Direct storage in database
- Compatibility with external protobuf systems

**Wire JSON (`pb_wire_json_*` functions):**
- Multiple field operations (2+ operations)
- Complex message transformations
- Performance-critical batch updates

```sql
-- Convert between formats
SELECT pb_message_to_wire_json(pb_data) AS wire_json;  -- Parse binary once
SELECT pb_wire_json_to_message(wire_json) AS pb_data;   -- Serialize to binary once
```

> **Performance Tip:** If you need to perform 2 or more operations on a message, convert to Wire JSON first, perform all operations, then convert back. This avoids repeated parsing/serialization overhead.

> **Data Integrity:** Round-trip conversion `pb_wire_json_to_message(pb_message_to_wire_json(message))` always produces exactly the same binary as the original message, preserving both field values and field ordering. This ensures no data loss.

### Field Numbers

All operations use protobuf field numbers (not names):

```protobuf
message User {
    string name = 1;     // field number 1
    int32 age = 2;       // field number 2
    string email = 3;    // field number 3
}
```

```sql
-- Access by field number
SELECT pb_message_get_string_field(pb_data, 1, '') AS name;  -- field 1
SELECT pb_message_get_int32_field(pb_data, 2, 0) AS age;     -- field 2
```

## Field Operations

### Reading Fields

#### Basic Low-level Field Access
```sql
-- Scalar fields with default values
SELECT pb_message_get_string_field(pb_data, 1, 'Unknown') AS name;
SELECT pb_message_get_int32_field(pb_data, 2, 0) AS age;
SELECT pb_message_get_bool_field(pb_data, 3, FALSE) AS active;

-- Check field existence
SELECT pb_message_has_string_field(pb_data, 1) AS has_name;
```

#### All Supported Field Types
```sql
-- Numeric types
SELECT pb_message_get_int32_field(pb_data, 1, 0);
SELECT pb_message_get_int64_field(pb_data, 2, 0);
SELECT pb_message_get_uint32_field(pb_data, 3, 0);
SELECT pb_message_get_uint64_field(pb_data, 4, 0);
SELECT pb_message_get_sint32_field(pb_data, 5, 0);
SELECT pb_message_get_sint64_field(pb_data, 6, 0);
SELECT pb_message_get_fixed32_field(pb_data, 7, 0);
SELECT pb_message_get_fixed64_field(pb_data, 8, 0);
SELECT pb_message_get_sfixed32_field(pb_data, 9, 0);
SELECT pb_message_get_sfixed64_field(pb_data, 10, 0);

-- Floating point
SELECT pb_message_get_float_field(pb_data, 11, 0.0);
SELECT pb_message_get_double_field(pb_data, 12, 0.0);

-- Other types
SELECT pb_message_get_bool_field(pb_data, 13, FALSE);
SELECT pb_message_get_string_field(pb_data, 14, '');
SELECT pb_message_get_bytes_field(pb_data, 15, '');
SELECT pb_message_get_enum_field(pb_data, 16, 0);
SELECT pb_message_get_message_field(pb_data, 17, pb_message_new());
```

### Writing Fields

#### Setting Field Values
```sql
-- Update multiple fields in a message
SET @updated_msg = pb_message_set_string_field(pb_data, 1, 'Jane Smith');
SET @updated_msg = pb_message_set_int32_field(@updated_msg, 2, 30);
SET @updated_msg = pb_message_set_bool_field(@updated_msg, 3, TRUE);

-- Chained operations
SELECT pb_message_set_string_field(
    pb_message_set_int32_field(
        pb_message_set_bool_field(pb_data, 3, TRUE),
        2, 30
    ),
    1, 'Jane Smith'
) AS updated_message;
```

#### Clearing Fields
```sql
-- Remove specific fields
SELECT pb_message_clear_string_field(pb_data, 1) AS no_name;
SELECT pb_message_clear_int32_field(pb_data, 2) AS no_age;

-- Clear all fields (returns empty message)
SELECT pb_message_clear(pb_data) AS empty_message;
```

### Working with Nested Messages

```sql
-- Create nested message
SET @address = pb_message_new();
SET @address = pb_message_set_string_field(@address, 1, '123 Main St');
SET @address = pb_message_set_string_field(@address, 2, 'Springfield');
SET @address = pb_message_set_string_field(@address, 3, 'USA');

-- Add to parent message
SET @user = pb_message_set_message_field(pb_message_new(), 5, @address);

-- Read nested field
SET @user_address = pb_message_get_message_field(@user, 5, pb_message_new());
SELECT pb_message_get_string_field(@user_address, 1, '') AS street;

-- One-liner for nested access
SELECT pb_message_get_string_field(
    pb_message_get_message_field(pb_data, 5, pb_message_new()),
    1, ''
) AS street;
```

## Repeated Fields

### Basic Repeated Field Operations

```sql
-- Add elements to repeated field
SET @msg = pb_message_add_repeated_string_field_element(pb_message_new(), 4, 'item1');
SET @msg = pb_message_add_repeated_string_field_element(@msg, 4, 'item2');
SET @msg = pb_message_add_repeated_string_field_element(@msg, 4, 'item3');

-- Get repeated field information
SELECT pb_message_get_repeated_string_field_count(@msg, 4) AS count;  -- Returns 3
SELECT pb_message_get_repeated_string_field_element(@msg, 4, 0) AS first_item;  -- Returns 'item1'
SELECT pb_message_get_repeated_string_field_element(@msg, 4, 2) AS third_item;  -- Returns 'item3'
```

### Advanced Repeated Field Operations

#### Insert at Specific Position
```sql
-- Insert at beginning (index 0)
SET @msg = pb_message_insert_repeated_string_field_element(@msg, 4, 0, 'new_first');

-- Insert in middle (index 2)
SET @msg = pb_message_insert_repeated_string_field_element(@msg, 4, 2, 'middle_item');
```

#### Update Existing Elements
```sql
-- Replace element at specific index
SET @msg = pb_message_set_repeated_string_field_element(@msg, 4, 1, 'updated_item');
```

#### Remove Elements
```sql
-- Remove element at specific index
SET @msg = pb_message_remove_repeated_string_field_element(@msg, 4, 0);

-- Clear all elements in repeated field
SET @msg = pb_message_clear_repeated_string_field(@msg, 4);
```

### Bulk Repeated Field Operations

#### Add Multiple Elements at Once
```sql
-- Add multiple strings
SELECT pb_message_add_all_repeated_string_field_elements(
    pb_data, 4, '["item1", "item2", "item3"]'
) AS with_items;

-- Add multiple integers (supports packed encoding)
SELECT pb_message_add_all_repeated_int32_field_elements(
    pb_data, 5, '[1, 2, 3, 4, 5]', TRUE  -- TRUE for packed encoding
) AS with_numbers;

-- Add multiple integers (unpacked)
SELECT pb_message_add_all_repeated_int32_field_elements(
    pb_data, 5, '[1, 2, 3, 4, 5]', FALSE  -- FALSE for unpacked
) AS with_numbers_unpacked;
```

#### Replace All Elements
```sql
-- Replace entire repeated field content
SELECT pb_message_set_repeated_string_field(
    pb_data, 4, '["new1", "new2", "new3"]'
) AS replaced_items;

-- Replace with packed numeric data
SELECT pb_message_set_repeated_int32_field(
    pb_data, 5, '[10, 20, 30]', TRUE
) AS replaced_numbers;
```

### Working with Repeated Messages

```sql
-- Create repeated message elements
SET @phone1 = pb_message_set_string_field(pb_message_new(), 1, '+1-555-0001');
SET @phone1 = pb_message_set_int32_field(@phone1, 2, 1);  -- MOBILE type

SET @phone2 = pb_message_set_string_field(pb_message_new(), 1, '+1-555-0002');
SET @phone2 = pb_message_set_int32_field(@phone2, 2, 2);  -- HOME type

-- Add to repeated message field
SET @user = pb_message_add_repeated_message_field_element(pb_message_new(), 4, @phone1);
SET @user = pb_message_add_repeated_message_field_element(@user, 4, @phone2);

-- Access repeated message elements
SET @first_phone = pb_message_get_repeated_message_field_element(@user, 4, 0);
SELECT pb_message_get_string_field(@first_phone, 1, '') AS phone_number;
SELECT pb_message_get_int32_field(@first_phone, 2, 0) AS phone_type;
```

## JSON Integration

### Schema Management

#### Loading Protobuf Schemas
```sql
-- Load schema from binary FileDescriptorSet
CALL pb_descriptor_set_load(
    'my_schema',  -- Schema identifier
    @binary_descriptor_set  -- Binary data from protoc --descriptor_set_out
);

-- List loaded schemas
SELECT pb_descriptor_set_list() AS loaded_schemas;

-- Get schema information
SELECT pb_descriptor_message_names('my_schema') AS message_types;
```

#### Generating Schema Binary
```bash
# Generate binary schema file
protoc --descriptor_set_out=schema.binpb --include_imports your_file.proto

# Convert to hex for MySQL
xxd -p -c0 schema.binpb
```

### JSON Conversion

#### Basic Message to JSON
```sql
-- Convert message to JSON (requires loaded schema)
SELECT pb_message_to_json('my_schema', '.Person', pb_data) AS json_output;

-- Example output:
-- {
--   "name": "John Doe",
--   "age": 25,
--   "email": "john@example.com",
--   "phones": [
--     {"number": "+1-555-0001", "type": "PHONE_TYPE_MOBILE"}
--   ]
-- }
```

#### Well-Known Types
```sql
-- Convert timestamp to JSON
SELECT pb_timestamp_to_json(timestamp_field) AS timestamp_json;
-- Output: "2025-06-01T12:34:56.789Z"

-- Convert duration to JSON
SELECT pb_duration_to_json(duration_field) AS duration_json;
-- Output: "123.456s"

-- Convert Any type to JSON
SELECT pb_any_to_json('my_schema', any_field) AS any_json;

-- Convert Struct to JSON
SELECT pb_struct_to_json(struct_field) AS struct_json;

-- Convert Value to JSON
SELECT pb_value_to_json(value_field) AS value_json;

-- Convert ListValue to JSON
SELECT pb_list_value_to_json(list_value_field) AS list_json;
```

### Working with Wire JSON

Wire JSON is an intermediate format that provides efficient access to protobuf data without repeated parsing:

```sql
-- Convert message to wire JSON (parses binary once)
SELECT pb_message_to_wire_json(pb_data) AS wire_json;

-- Example wire JSON structure:
-- {
--   "1": [{"i": 0, "n": 1, "t": 2, "v": "Sm9obiBEb2U="}],  -- string field (base64)
--   "2": [{"i": 1, "n": 2, "t": 0, "v": 25}],              -- int32 field  
--   "3": [{"i": 2, "n": 3, "t": 2, "v": "am9obkBleGFtcGxlLmNvbQ=="}]  -- email field
-- }
-- Where: i=index, n=field number, t=wire type, v=value

-- Efficient pattern for multiple operations:
SET @wire_json = pb_message_to_wire_json(pb_data);  -- Parse once

-- Perform multiple modifications without reparsing
SET @wire_json = pb_wire_json_set_string_field(@wire_json, 1, 'New Name');
SET @wire_json = pb_wire_json_set_int32_field(@wire_json, 2, 30);
SET @wire_json = pb_wire_json_add_repeated_string_field_element(@wire_json, 4, 'hobby');

-- Convert back to binary (serialize once)
SELECT pb_wire_json_to_message(@wire_json) AS updated_pb_data;
```

**Performance Rule:** Use Wire JSON when you need 2+ operations on the same message.

## Advanced Topics

### Schema Evolution and Compatibility

This library provides similar (mostly same) backward and forward compatibility as described in the [Protocol Buffers Language Guide](https://protobuf.dev/programming-guides/proto3/#updating) when updating protobuf schemas.

The difference is that this library prioritizes data integrity over silent truncation or silent loss of data.
For example, when `int32` field is changed to `int64`, `pb_message_get_int32_field` will raise an error if the field value exceeds the `int32` range, whereas the Protocol Buffers documentation says the value is silently truncated to 32 bits.

#### Compatible Schema Changes

**Type Widening**: Fields can be safely widened to larger types:
```sql
-- Original schema: int32 age = 2;
-- Updated schema: int64 age = 2;

-- This continues to work as long as values fit in int32 range:
SELECT pb_message_get_int32_field(data, 2, 0) AS age FROM users;

-- New code can use the wider type:
SELECT pb_message_get_int64_field(data, 2, 0) AS age FROM users;
```

**Field Addition**: New fields can be added without breaking existing code:
```sql
-- Original schema: string name = 1; int32 age = 2;
-- Updated schema: string name = 1; int32 age = 2; string email = 3;

-- Existing queries continue to work:
SELECT pb_message_get_string_field(data, 1, '') AS name FROM users;

-- New queries can access the new field:
SELECT pb_message_get_string_field(data, 3, '') AS email FROM users;
```

**Field Renaming**: Fields can be renamed without affecting low-level operations:
```sql
-- Original schema: string user_name = 1; int32 user_age = 2;
-- Updated schema:  string full_name = 1; int32 age_years = 2;

-- Low-level operations continue to work unchanged:
SELECT pb_message_get_string_field(data, 1, '') AS name FROM users;
SELECT pb_message_get_int32_field(data, 2, 0) AS age FROM users;

-- Field numbers stay the same, only names changed in the .proto file
```

#### Data Integrity Protections

Unlike some protobuf implementations that silently truncate values, this library prioritizes data integrity:

```sql
-- If a field was changed from int32 to int64 and contains a large value:
-- Original: int32 user_id = 1;
-- Updated:  int64 user_id = 1;
-- Current value: 3000000000 (exceeds int32 range)

-- This will raise an error instead of silent truncation:
SELECT pb_message_get_int32_field(data, 1, 0) FROM users; -- ERROR!

-- Use the correct wider type:
SELECT pb_message_get_int64_field(data, 1, 0) FROM users; -- Works correctly
```

#### Schema Migration Strategies

**Gradual Type Migration**:
```sql
-- Check if values fit in the old type before migration
SELECT user_id, 
       pb_message_get_int64_field(data, 2, 0) AS new_age,
       CASE 
         WHEN pb_message_get_int64_field(data, 2, 0) BETWEEN -2147483648 AND 2147483647 THEN
           'Can use int32'
         ELSE 
           'Requires int64'
       END AS compatibility_status
FROM users;
```

#### Best Practices for Schema Evolution

* Don’t change the field numbers for any existing fields.

> **Note:** This compatibility applies to the low-level field operations. JSON conversion functions require the exact schema and field names, so they are less flexible during schema evolution.

### Indexing Protobuf Fields

Since MySQL doesn't support stored functions in generated columns or functional indexes, use triggers to extract fields into regular columns:

```sql
-- Add extracted columns
ALTER TABLE users 
ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN age INT NOT NULL DEFAULT 0;

-- Create triggers to maintain extracted columns
DELIMITER $$
CREATE TRIGGER users_extract_fields_on_insert
    BEFORE INSERT ON users
    FOR EACH ROW
BEGIN
    SET NEW.name = pb_message_get_string_field(NEW.pb_data, 1, '');
    SET NEW.age = pb_message_get_int32_field(NEW.pb_data, 2, 0);
END$$

CREATE TRIGGER users_extract_fields_on_update
    BEFORE UPDATE ON users
    FOR EACH ROW
BEGIN
    SET NEW.name = pb_message_get_string_field(NEW.pb_data, 1, '');
    SET NEW.age = pb_message_get_int32_field(NEW.pb_data, 2, 0);
END$$
DELIMITER ;
```

### Working with Large Messages

For large protobuf messages, consider these strategies:

```sql
-- Use wire JSON for intermediate processing
SET @wire_json = pb_message_to_wire_json(large_pb_data);
SET @wire_json = pb_wire_json_set_string_field(@wire_json, 1, 'new_value');
SET @wire_json = pb_wire_json_add_repeated_int32_field_element(@wire_json, 5, 123);
SET @result = pb_wire_json_to_message(@wire_json);

-- Batch process repeated fields
SELECT pb_message_add_all_repeated_string_field_elements(
    pb_data, 4, 
    JSON_ARRAY('item1', 'item2', 'item3', 'item4', 'item5')
) AS batch_updated;
```

### Custom Validation

```sql
-- Create validation functions
DELIMITER $$
CREATE FUNCTION validate_user_message(pb_data LONGBLOB) 
RETURNS BOOLEAN 
READS SQL DATA 
DETERMINISTIC
BEGIN
    DECLARE name_length INT;
    DECLARE age_value INT;
    
    -- Check required fields exist
    IF NOT pb_message_has_string_field(pb_data, 1) THEN
        RETURN FALSE;
    END IF;
    
    -- Validate field values
    SET name_length = LENGTH(pb_message_get_string_field(pb_data, 1, ''));
    IF name_length = 0 OR name_length > 100 THEN
        RETURN FALSE;
    END IF;
    
    SET age_value = pb_message_get_int32_field(pb_data, 2, -1);
    IF age_value < 0 OR age_value > 150 THEN
        RETURN FALSE;
    END IF;
    
    RETURN TRUE;
END$$
DELIMITER ;

-- Use in constraints
ALTER TABLE users ADD CONSTRAINT chk_valid_pb_data 
CHECK (validate_user_message(pb_data));
```

## Performance Considerations

### Function Caching

MySQL caches stored functions, but the cache has limits:

```sql
-- Increase stored program cache if needed
SET GLOBAL stored_program_cache = 512;  -- Default: 256
```

### Bulk Operations

Use bulk operations when possible:

```sql
-- Efficient: Add multiple elements at once
SELECT pb_message_add_all_repeated_int32_field_elements(pb_data, 4, '[1,2,3,4,5]', TRUE);

-- Inefficient: Add elements one by one
SET @msg = pb_message_add_repeated_int32_field_element(pb_data, 4, 1);
SET @msg = pb_message_add_repeated_int32_field_element(@msg, 4, 2);
SET @msg = pb_message_add_repeated_int32_field_element(@msg, 4, 3);
-- ... etc
```

### Packed Encoding

Use packed encoding for numeric repeated fields when supported:

```sql
-- Packed: More efficient for large arrays of numbers
SELECT pb_message_add_all_repeated_int32_field_elements(pb_data, 4, '[1,2,3,4,5]', TRUE);

-- Unpacked: Less efficient but sometimes required for compatibility
SELECT pb_message_add_all_repeated_int32_field_elements(pb_data, 4, '[1,2,3,4,5]', FALSE);
```

### Wire JSON Performance Pattern

Wire JSON is crucial for performance when doing multiple operations:

```sql
-- ❌ SLOW: Each pb_message_* function parses and serializes the entire message
UPDATE large_messages 
SET data = pb_message_set_string_field(
  pb_message_set_int32_field(
    pb_message_add_repeated_string_field_element(
      pb_message_clear_string_field(data, 5),
      4, 'new_item'
    ),
    2, 999
  ),
  1, 'Updated Name'
)
WHERE id = 1;
-- This parses the message 4 times and serializes it 4 times!

-- ✅ FAST: Parse once, modify multiple times, serialize once
UPDATE large_messages 
SET data = (
  SELECT pb_wire_json_to_message(
    pb_wire_json_set_string_field(
      pb_wire_json_set_int32_field(
        pb_wire_json_add_repeated_string_field_element(
          pb_wire_json_clear_string_field(@wire, 5),
          4, 'new_item'
        ),
        2, 999
      ),
      1, 'Updated Name'
    )
  )
  FROM (SELECT pb_message_to_wire_json(data) AS wire FROM large_messages WHERE id = 1) AS w
  CROSS JOIN (SELECT @wire := wire) AS assignment
)
WHERE id = 1;
-- This parses once and serializes once, regardless of operation count!
```

**Performance insight:** Wire JSON eliminates repeated parsing overhead, providing significant performance benefits for multiple operations on larger messages.

## Troubleshooting

### Common Issues

#### Function Not Found
```sql
-- Check if functions are loaded
SHOW FUNCTION STATUS WHERE Name LIKE 'pb_%';

-- Reload functions if needed
SOURCE protobuf.sql;
SOURCE protobuf-accessors.sql;
```

#### Invalid Field Numbers
```sql
-- Check field exists before accessing
IF pb_message_has_string_field(pb_data, 1) THEN
    SET @name = pb_message_get_string_field(pb_data, 1, '');
END IF;
```

#### JSON Conversion Errors
```sql
-- Verify schema is loaded
SELECT pb_descriptor_set_list();

-- Check message type exists
SELECT pb_descriptor_message_names('schema_id');

-- Verify message type name (include leading dot)
SELECT pb_message_to_json('schema_id', '.MessageTypeName', pb_data);
```

#### Wire JSON Format Issues
```sql
-- Validate wire JSON format
SELECT JSON_VALID(wire_json_data);

-- Pretty print for debugging
SELECT JSON_PRETTY(pb_message_to_wire_json(pb_data));
```

### Debugging Tools

#### Inspect Message Contents
```sql
-- View raw wire format as JSON
SELECT JSON_PRETTY(pb_message_to_wire_json(pb_data)) AS wire_format;

-- Convert to readable JSON (if schema available)
SELECT JSON_PRETTY(pb_message_to_json('schema_id', '.MessageType', pb_data)) AS readable_json;
```

#### Check Message Structure
```sql
-- List all field numbers present in message
SELECT DISTINCT 
    JSON_UNQUOTE(JSON_EXTRACT(field_data, '$[0].n')) AS field_number
FROM (
    SELECT JSON_EXTRACT(pb_message_to_wire_json(pb_data), CONCAT('$."', field_key, '"')) AS field_data
    FROM (
        SELECT field_key 
        FROM JSON_TABLE(
            JSON_KEYS(pb_message_to_wire_json(pb_data)), 
            '$[*]' COLUMNS (field_key VARCHAR(10) PATH '$')
        ) AS keys
    ) AS field_keys
) AS fields;
```

#### Performance Analysis
```sql
-- Enable profiling
SET profiling = 1;

-- Run your protobuf operations
SELECT pb_message_get_string_field(pb_data, 1, '') FROM users LIMIT 1000;

-- Check profile
SHOW PROFILES;
SHOW PROFILE FOR QUERY 1;
```