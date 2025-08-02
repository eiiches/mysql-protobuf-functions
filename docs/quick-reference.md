# MySQL Protocol Buffers - Quick Reference

*Copy-paste syntax for common operations*

## 📚 Navigation

- **👈 [Documentation Home](README.md)** - Choose your learning path
- **🎯 [User Guide](user-guide.md)** - Problem-focused examples and solutions
- **📖 [Function Reference](function-reference.md)** - Complete API documentation  
- **🔬 [API Guide](api-guide.md)** - Advanced technical details

---

## Common Operations Cheat Sheet

### Message Creation and Basic Operations

```sql
-- Create new empty message
SELECT pb_message_new();

-- Convert message to wire JSON (for multiple operations or inspection)
SELECT pb_message_to_wire_json(pb_data);

-- Convert wire JSON back to message
SELECT pb_wire_json_to_message(wire_json);
```

### Low-level Field Access (Reading)

```sql
-- Get scalar field values
SELECT pb_message_get_string_field(pb_data, 1, 'default_value');
SELECT pb_message_get_int32_field(pb_data, 2, 0);
SELECT pb_message_get_bool_field(pb_data, 5, FALSE);

-- Check if field exists
SELECT pb_message_has_string_field(pb_data, 1);

-- Get repeated field count
SELECT pb_message_get_repeated_string_field_count(pb_data, 4);

-- Get repeated field element
SELECT pb_message_get_repeated_string_field_element(pb_data, 4, 0); -- first element
```

### Low-level Field Manipulation (Writing)

```sql
-- Set scalar fields
SELECT pb_message_set_string_field(pb_data, 1, 'new_value');
SELECT pb_message_set_int32_field(pb_data, 2, 123);
SELECT pb_message_set_bool_field(pb_data, 5, TRUE);

-- Clear fields
SELECT pb_message_clear_string_field(pb_data, 1);
SELECT pb_message_clear_repeated_string_field(pb_data, 4);

-- Set message fields
SELECT pb_message_set_message_field(pb_data, 5, sub_message);
```

### Repeated Field Operations

```sql
-- Add element to repeated field
SELECT pb_message_add_repeated_string_field_element(pb_data, 4, 'new_item');

-- Insert element at specific position
SELECT pb_message_insert_repeated_string_field_element(pb_data, 4, 1, 'inserted_item');

-- Set element at specific position
SELECT pb_message_set_repeated_string_field_element(pb_data, 4, 0, 'updated_item');

-- Remove element at specific position
SELECT pb_message_remove_repeated_string_field_element(pb_data, 4, 0);

-- Add multiple elements at once (bulk operation)
SELECT pb_message_add_all_repeated_string_field_elements(pb_data, 4, '["item1", "item2", "item3"]');

-- Replace all elements (set operation) 
SELECT pb_message_set_repeated_string_field(pb_data, 4, '["new1", "new2"]');
```

### JSON Integration

```sql
-- Generate protobuf schema JSON
SET @schema_json = pb_build_descriptor_set_json(binary_descriptor_set);

-- Convert message to JSON (requires schema)
SELECT pb_message_to_json(@schema_json, '.MessageType', pb_data);

-- Work with well-known types
SELECT pb_timestamp_to_json(timestamp_message);
SELECT pb_duration_to_json(duration_message);
SELECT pb_any_to_json('schema_id', any_message);
```

## Field Types Reference

| Protobuf Type | MySQL Function Suffix | MySQL Parameter Type | Example |
|---------------|----------------------|----------------------|---------|
| `string` | `_string_field` | `LONGTEXT` | `pb_message_get_string_field(msg, 1, '')` |
| `int32` | `_int32_field` | `INT` | `pb_message_get_int32_field(msg, 2, 0)` |
| `int64` | `_int64_field` | `BIGINT` | `pb_message_get_int64_field(msg, 3, 0)` |
| `uint32` | `_uint32_field` | `INT UNSIGNED` | `pb_message_get_uint32_field(msg, 4, 0)` |
| `uint64` | `_uint64_field` | `BIGINT UNSIGNED` | `pb_message_get_uint64_field(msg, 5, 0)` |
| `bool` | `_bool_field` | `BOOLEAN` | `pb_message_get_bool_field(msg, 6, FALSE)` |
| `bytes` | `_bytes_field` | `LONGBLOB` | `pb_message_get_bytes_field(msg, 7, '')` |
| `float` | `_float_field` | `FLOAT` | `pb_message_get_float_field(msg, 8, 0.0)` |
| `double` | `_double_field` | `DOUBLE` | `pb_message_get_double_field(msg, 9, 0.0)` |
| `sint32` | `_sint32_field` | `INT` | `pb_message_get_sint32_field(msg, 10, 0)` |
| `sint64` | `_sint64_field` | `BIGINT` | `pb_message_get_sint64_field(msg, 11, 0)` |
| `fixed32` | `_fixed32_field` | `INT UNSIGNED` | `pb_message_get_fixed32_field(msg, 12, 0)` |
| `fixed64` | `_fixed64_field` | `BIGINT UNSIGNED` | `pb_message_get_fixed64_field(msg, 13, 0)` |
| `sfixed32` | `_sfixed32_field` | `INT` | `pb_message_get_sfixed32_field(msg, 14, 0)` |
| `sfixed64` | `_sfixed64_field` | `BIGINT` | `pb_message_get_sfixed64_field(msg, 15, 0)` |
| `enum` | `_enum_field` | `INT` | `pb_message_get_enum_field(msg, 16, 0)` |
| `message` | `_message_field` | `LONGBLOB` | `pb_message_get_message_field(msg, 17, pb_message_new())` |

## Packed Repeated Fields

Some numeric repeated fields support packed encoding for efficiency:

```sql
-- Add elements with packed encoding (supported: int32, int64, uint32, uint64, bool, enum, float, double, etc.)
SELECT pb_message_add_all_repeated_int32_field_elements(pb_data, 4, '[1, 2, 3]', TRUE); -- packed
SELECT pb_message_add_all_repeated_int32_field_elements(pb_data, 4, '[1, 2, 3]', FALSE); -- unpacked

-- String, bytes, and message fields are always unpacked
SELECT pb_message_add_all_repeated_string_field_elements(pb_data, 5, '["a", "b", "c"]'); -- no packed parameter
```

## Function Naming Patterns

### Public API Pattern
```
pb_{input_type}_{operation}_{field_type}_field[_element[s]]
```

### Input Types
- `message` - Binary protobuf format (LONGBLOB)
- `wire_json` - JSON representation of wire format

### Operations
- `get` - Retrieve field value
- `set` - Set field value  
- `has` - Check if field exists
- `clear` - Remove field
- `add` - Add to repeated field
- `insert` - Insert into repeated field at position
- `remove` - Remove from repeated field at position

### Examples
```sql
pb_message_get_string_field()           -- Get string field from message
pb_wire_json_set_int32_field()          -- Set int32 field in wire JSON
pb_message_add_repeated_bool_field_element()  -- Add bool to repeated field
pb_message_set_repeated_string_field() -- Replace all string elements
```

**Important**: Use only public functions starting with `pb_`. Functions starting with `_pb_` are internal and may change without notice.

## Wire JSON Format

Wire JSON is an intermediate format for **efficient multi-operation processing**:

```sql
-- Example wire JSON structure
{
  "1": [{"i": 0, "n": 1, "t": 2, "v": "SGVsbG8="}],  -- string field
  "2": [{"i": 1, "n": 2, "t": 0, "v": 42}]            -- int32 field
}

-- Where:
-- "1", "2" = field numbers
-- "i" = index (for ordering)
-- "n" = field number
-- "t" = wire type (0=varint, 1=fixed64, 2=length-delimited, 5=fixed32)
-- "v" = value (base64 for bytes/strings, raw for numbers)
```

### Performance Pattern
```sql
-- When doing 2+ operations, use Wire JSON:
SET @wire = pb_message_to_wire_json(data);           -- Parse once
SET @wire = pb_wire_json_set_string_field(@wire, 1, 'Name');
SET @wire = pb_wire_json_set_int32_field(@wire, 2, 30);
SET @updated = pb_wire_json_to_message(@wire);       -- Serialize once
```

**Guarantee:** Round-trip conversion always produces identical binary output, preserving field values and ordering.

## Common Patterns

### Building Messages Step by Step
```sql
SET @msg = pb_message_new();
SET @msg = pb_message_set_string_field(@msg, 1, 'John Doe');
SET @msg = pb_message_set_int32_field(@msg, 2, 25);
SET @msg = pb_message_add_repeated_string_field_element(@msg, 3, 'hobby1');
SET @msg = pb_message_add_repeated_string_field_element(@msg, 3, 'hobby2');
```

### Chaining Operations
```sql
SELECT pb_message_add_repeated_string_field_element(
  pb_message_set_int32_field(
    pb_message_set_string_field(pb_message_new(), 1, 'Jane Doe'),
    2, 30
  ),
  3, 'reading'
) AS complete_message;
```

### Working with Nested Messages
```sql
-- Create nested message first
SET @address = pb_message_set_string_field(pb_message_new(), 1, '123 Main St');
SET @address = pb_message_set_string_field(@address, 2, 'Anytown');

-- Add to parent message
SET @person = pb_message_set_message_field(pb_message_new(), 5, @address);
```

## 🎯 Common Use Cases

### Database Queries with Protobuf Fields
```sql
-- Find users in age range
SELECT user_id FROM users 
WHERE pb_message_get_int32_field(profile_data, 2, 0) BETWEEN 25 AND 35;

-- Search by partial name match
SELECT user_id FROM users
WHERE pb_message_get_string_field(profile_data, 1, '') LIKE '%John%';

-- Count users with specific hobby
SELECT COUNT(*) FROM users u
JOIN (
  SELECT user_id, 
         pb_message_get_repeated_string_field_count(profile_data, 4) as hobby_count
  FROM users
) h ON u.user_id = h.user_id
WHERE h.hobby_count > 0 
  AND pb_message_get_repeated_string_field_element(u.profile_data, 4, 0) = 'hiking';
```

### Safe Data Updates
```sql
-- Update only if field exists (defensive programming)
UPDATE users 
SET profile_data = CASE 
  WHEN pb_message_has_string_field(profile_data, 3) THEN
    pb_message_set_string_field(profile_data, 3, 'new@email.com')
  ELSE profile_data
END
WHERE user_id = 123;

-- Atomic multi-field update
UPDATE users 
SET profile_data = pb_message_set_string_field(
  pb_message_set_int32_field(profile_data, 2, 26),  -- age
  3, 'updated@email.com'                            -- email
)
WHERE user_id = 123;
```

### Data Migration and Schema Evolution
```sql
-- Copy data from old field to new field
UPDATE users 
SET profile_data = pb_message_set_string_field(
  profile_data, 
  6,  -- new mobile field
  pb_message_get_string_field(profile_data, 5, '')  -- old phone field
)
WHERE pb_message_has_string_field(profile_data, 5)    -- has old field
  AND NOT pb_message_has_string_field(profile_data, 6); -- missing new field
```

## 💡 Pro Tips

### Performance Optimizations
```sql
-- Extract frequently queried fields to indexed columns (using triggers)
ALTER TABLE users ADD COLUMN age INT NOT NULL DEFAULT 0;

DELIMITER $$
CREATE TRIGGER users_extract_age_insert BEFORE INSERT ON users FOR EACH ROW
BEGIN
  SET NEW.age = pb_message_get_int32_field(NEW.profile_data, 2, 0);
END$$
CREATE TRIGGER users_extract_age_update BEFORE UPDATE ON users FOR EACH ROW  
BEGIN
  SET NEW.age = pb_message_get_int32_field(NEW.profile_data, 2, 0);
END$$
DELIMITER ;

CREATE INDEX idx_users_age ON users(age);

-- Use bulk operations for repeated fields
-- ✅ Efficient
SELECT pb_message_add_all_repeated_string_field_elements(data, 4, '["a","b","c"]');
-- ❌ Inefficient  
SELECT pb_message_add_repeated_string_field_element(
  pb_message_add_repeated_string_field_element(
    pb_message_add_repeated_string_field_element(data, 4, 'a'), 4, 'b'
  ), 4, 'c'
);
```

### Debugging and Inspection
```sql
-- Inspect message structure: see all fields in wire format
SELECT JSON_PRETTY(pb_message_to_wire_json(profile_data)) AS wire_structure
FROM users WHERE user_id = 123;

-- Check what fields are present
SELECT 
  pb_message_has_string_field(profile_data, 1) AS has_name,
  pb_message_has_int32_field(profile_data, 2) AS has_age,
  pb_message_get_repeated_string_field_count(profile_data, 4) AS hobby_count
FROM users WHERE user_id = 123;
```

## 🔗 Need More Detail?

- **Real-world examples** → [User Guide](user-guide.md) - Problem-focused solutions with error handling  
- **All function parameters** → [Function Reference](function-reference.md) - Complete API documentation
- **Advanced topics** → [API Guide](api-guide.md) - Indexing, validation, performance tuning