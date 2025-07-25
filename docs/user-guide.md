# MySQL Protocol Buffers - User Guide

*A practical guide for developers using protobuf data in MySQL*

## 📚 Navigation

- **👈 [Documentation Home](README.md)** - Choose your learning path
- **⚡ [Quick Reference](quick-reference.md)** - Function syntax cheat sheet  
- **📖 [Function Reference](function-reference.md)** - Complete API documentation
- **🔬 [API Guide](api-guide.md)** - Advanced technical details

---

## When to Use This Library

**Use this library when you:**
- Store protobuf messages in MySQL as LONGBLOB columns
- Need to query or filter based on protobuf field values
- Want to extract specific fields directly in SQL queries
- Need to modify protobuf data directly in the database
- Want to inspect protobuf data for debugging

**Don't use this library when:**
- You can process protobuf data in your application (usually more efficient)
- You're working with simple data that doesn't need protobuf
- Performance is critical and you can restructure your data model

## Common Use Cases & Solutions

### 1. Querying Records by Protobuf Field Values

**Problem**: You have a `users` table with a `profile_data` LONGBLOB column containing protobuf messages, and you want to find users by age.

**Solution**: Extract the field in your WHERE clause
```sql
-- Find users aged 25-35
SELECT user_id, profile_data 
FROM users 
WHERE pb_message_get_int32_field(profile_data, 2, 0) BETWEEN 25 AND 35;
```

**Performance Tip**: For frequently queried fields, create extracted columns with indexes:
```sql
-- Add extracted column
ALTER TABLE users ADD COLUMN age INT AS (pb_message_get_int32_field(profile_data, 2, 0)) STORED;
CREATE INDEX idx_users_age ON users(age);

-- Now queries are fast
SELECT user_id FROM users WHERE age BETWEEN 25 AND 35;
```

### 2. Building Messages Incrementally

**Problem**: You need to construct a protobuf message with data from multiple sources.

**Solution**: Start with empty message and build step by step
```sql
-- Build a user profile message
SET @profile = pb_message_new();
SET @profile = pb_message_set_string_field(@profile, 1, 'John Doe');        -- name
SET @profile = pb_message_set_int32_field(@profile, 2, 30);                 -- age  
SET @profile = pb_message_set_string_field(@profile, 3, 'john@example.com'); -- email

-- Add repeated fields (hobbies)
SET @profile = pb_message_add_repeated_string_field_element(@profile, 4, 'reading');
SET @profile = pb_message_add_repeated_string_field_element(@profile, 4, 'hiking');

INSERT INTO users (profile_data) VALUES (@profile);
```

**Alternative (Chained approach)**:
```sql
INSERT INTO users (profile_data) VALUES (
  pb_message_add_repeated_string_field_element(
    pb_message_add_repeated_string_field_element(
      pb_message_set_string_field(
        pb_message_set_int32_field(
          pb_message_set_string_field(pb_message_new(), 1, 'John Doe'),
          2, 30
        ),
        3, 'john@example.com'
      ),
      4, 'reading'
    ),
    4, 'hiking'
  )
);
```

### 3. Safely Updating Protobuf Messages

**Problem**: You need to update a specific field in an existing protobuf message without losing other data.

**Solution**: Use set functions that preserve other fields
```sql
-- Update only the age field, keeping everything else
UPDATE users 
SET profile_data = pb_message_set_int32_field(profile_data, 2, 31)
WHERE user_id = 123;

-- Update multiple fields atomically
UPDATE users 
SET profile_data = pb_message_set_string_field(
  pb_message_set_int32_field(profile_data, 2, 31),  -- age
  3, 'newemail@example.com'                         -- email
)
WHERE user_id = 123;
```

### 4. Working with Repeated Fields Efficiently

**Problem**: You need to manage lists of data (tags, categories, phone numbers, etc.).

**Solution**: Use bulk operations when possible
```sql
-- Replace entire list (most efficient for complete updates)
UPDATE users 
SET profile_data = pb_message_set_repeated_string_field(
  profile_data, 
  4,  -- hobbies field
  '["reading", "hiking", "photography"]'
)
WHERE user_id = 123;

-- Add multiple items at once
UPDATE users 
SET profile_data = pb_message_add_all_repeated_string_field_elements(
  profile_data,
  4,  -- hobbies field  
  '["swimming", "cooking"]'
)
WHERE user_id = 123;

-- Individual operations (less efficient, but sometimes necessary)
UPDATE users 
SET profile_data = pb_message_add_repeated_string_field_element(profile_data, 4, 'gardening')
WHERE user_id = 123;
```

### 5. Debugging and Inspection

**Problem**: You need to see what's inside a protobuf message for debugging.

**Solution**: Convert to human-readable formats
```sql
-- Quick inspection: see wire format structure
SELECT user_id, JSON_PRETTY(pb_message_to_wire_json(profile_data)) AS wire_format
FROM users WHERE user_id = 123;

-- Full debugging: convert to JSON (requires schema loaded)
SELECT user_id, JSON_PRETTY(pb_message_to_json('user_schema', '.UserProfile', profile_data)) AS profile_json
FROM users WHERE user_id = 123;

-- Check specific fields exist
SELECT 
  user_id,
  pb_message_has_string_field(profile_data, 1) AS has_name,
  pb_message_has_int32_field(profile_data, 2) AS has_age,
  pb_message_get_repeated_string_field_count(profile_data, 4) AS hobby_count
FROM users 
WHERE user_id = 123;
```

### 6. Handling Schema Evolution

**Problem**: Your protobuf schema changed and you need to migrate existing data.

**Solution**: Use the field number stability of protobuf
```sql
-- Schema changed: field 5 (phone) was removed, field 6 (mobile) was added
-- Migrate data: copy phone to mobile if mobile doesn't exist
UPDATE users 
SET profile_data = pb_message_set_string_field(
  profile_data, 
  6,  -- new mobile field
  pb_message_get_string_field(profile_data, 5, '')  -- old phone field
)
WHERE pb_message_has_string_field(profile_data, 5)    -- has old field
  AND NOT pb_message_has_string_field(profile_data, 6); -- doesn't have new field

-- Remove old field
UPDATE users
SET profile_data = pb_message_clear_string_field(profile_data, 5)
WHERE pb_message_has_string_field(profile_data, 5);
```

### 7. Working with Nested Messages

**Problem**: You have nested protobuf messages and need to access/modify deep fields.

**Solution**: Extract nested messages, work with them, then update
```sql
-- Access nested field: user.address.street
SELECT 
  user_id,
  pb_message_get_string_field(
    pb_message_get_message_field(profile_data, 5, pb_message_new()),  -- address field
    1,  -- street field in address
    ''  -- default
  ) AS street
FROM users;

-- Update nested field
UPDATE users
SET profile_data = pb_message_set_message_field(
  profile_data,
  5,  -- address field
  pb_message_set_string_field(
    pb_message_get_message_field(profile_data, 5, pb_message_new()),  -- get current address
    1,  -- street field
    '123 New Street'  -- new value
  )
)
WHERE user_id = 123;
```

## Performance Best Practices

### 0. Use Wire JSON for Multiple Operations

**The most important optimization:** When performing 2+ operations on a message, use Wire JSON format to avoid repeated parsing.

```sql
-- ❌ SLOW: Parse and serialize 5 times
UPDATE users 
SET profile_data = pb_message_set_string_field(
  pb_message_set_int32_field(
    pb_message_add_repeated_string_field_element(
      pb_message_clear_repeated_string_field(
        pb_message_set_string_field(profile_data, 3, 'email@new.com'),
        4
      ),
      4, 'new_hobby'
    ),
    2, 30
  ),
  1, 'New Name'
)
WHERE user_id = 123;

-- ✅ FAST: Parse once, serialize once
UPDATE users 
SET profile_data = (
  SELECT pb_wire_json_to_message(wire_json) FROM (
    SELECT 
      pb_wire_json_set_string_field(
        pb_wire_json_set_int32_field(
          pb_wire_json_add_repeated_string_field_element(
            pb_wire_json_clear_repeated_string_field(
              pb_wire_json_set_string_field(w, 3, 'email@new.com'),
              4
            ),
            4, 'new_hobby'
          ),
          2, 30
        ),
        1, 'New Name'
      ) AS wire_json
    FROM (SELECT pb_message_to_wire_json(profile_data) AS w FROM users WHERE user_id = 123) t
  ) result
)
WHERE user_id = 123;
```

**Rule of thumb:** 
- 1 operation: Use `pb_message_*` functions
- 2+ operations: Consider Wire JSON to avoid repeated parsing
- Multiple operations on large messages: Wire JSON provides the best performance

## Performance Best Practices

### 1. Indexing Strategy
```sql
-- Add extracted columns (triggers will maintain these)
ALTER TABLE users 
ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN age INT NOT NULL DEFAULT 0,
ADD COLUMN created_at DATETIME NULL;

-- Create triggers to extract protobuf fields into regular columns
DELIMITER $$
CREATE TRIGGER users_extract_before_insert
    BEFORE INSERT ON users
    FOR EACH ROW
BEGIN
    SET NEW.name = pb_message_get_string_field(NEW.profile_data, 1, '');
    SET NEW.age = pb_message_get_int32_field(NEW.profile_data, 2, 0);
    SET NEW.created_at = FROM_UNIXTIME(pb_message_get_int64_field(NEW.profile_data, 10, 0));
END$$

CREATE TRIGGER users_extract_before_update
    BEFORE UPDATE ON users
    FOR EACH ROW
BEGIN
    SET NEW.name = pb_message_get_string_field(NEW.profile_data, 1, '');
    SET NEW.age = pb_message_get_int32_field(NEW.profile_data, 2, 0);
    SET NEW.created_at = FROM_UNIXTIME(pb_message_get_int64_field(NEW.profile_data, 10, 0));
END$$
DELIMITER ;

-- Now create indexes on extracted columns
CREATE INDEX idx_users_name ON users(name);
CREATE INDEX idx_users_age ON users(age);
CREATE INDEX idx_users_created_at ON users(created_at);
```

### 2. Batch Operations
```sql
-- Efficient: bulk update
UPDATE users 
SET profile_data = pb_message_set_repeated_string_field(profile_data, 4, '["a","b","c"]')
WHERE user_id IN (1, 2, 3);

-- Inefficient: multiple individual updates  
UPDATE users SET profile_data = pb_message_add_repeated_string_field_element(profile_data, 4, 'a') WHERE user_id = 1;
UPDATE users SET profile_data = pb_message_add_repeated_string_field_element(profile_data, 4, 'b') WHERE user_id = 1;
UPDATE users SET profile_data = pb_message_add_repeated_string_field_element(profile_data, 4, 'c') WHERE user_id = 1;
```

### 3. Wire JSON for Performance (Multiple Operations)
```sql
-- When doing 2+ operations, Wire JSON is faster than repeated parsing
-- Parse once, modify multiple times, serialize once

-- ❌ Inefficient: Each operation parses and serializes the entire message
UPDATE users 
SET profile_data = pb_message_set_string_field(
  pb_message_set_int32_field(
    pb_message_add_repeated_string_field_element(profile_data, 4, 'hobby1'),
    2, 25
  ),
  1, 'New Name'
)
WHERE user_id = 123;

-- ✅ Efficient: Parse once, modify in Wire JSON, serialize once
UPDATE users 
SET profile_data = pb_wire_json_to_message(
  pb_wire_json_add_repeated_string_field_element(
    pb_wire_json_set_int32_field(
      pb_wire_json_set_string_field(
        pb_message_to_wire_json(profile_data),  -- Parse once here
        1, 'New Name'
      ),
      2, 25
    ),
    4, 'new_hobby'
  )  -- Serialize once here
)
WHERE user_id = 123;
```

## Error Handling Patterns

### 1. Defensive Low-level Field Access
```sql
-- Always provide appropriate defaults
SELECT 
  user_id,
  COALESCE(NULLIF(pb_message_get_string_field(profile_data, 1, ''), ''), 'Unknown') AS name,
  pb_message_get_int32_field(profile_data, 2, 0) AS age
FROM users;
```

### 2. Validating Data Before Operations
```sql
-- Check field exists before complex operations
UPDATE users 
SET profile_data = CASE 
  WHEN pb_message_has_string_field(profile_data, 1) THEN
    pb_message_set_string_field(profile_data, 1, UPPER(pb_message_get_string_field(profile_data, 1, '')))
  ELSE 
    profile_data
END
WHERE user_id = 123;
```

### 3. Transaction Safety
```sql
START TRANSACTION;

-- Validate before modifying
SELECT COUNT(*) INTO @count 
FROM users 
WHERE user_id = 123 
  AND pb_message_has_string_field(profile_data, 1);

IF @count = 1 THEN
  UPDATE users 
  SET profile_data = pb_message_set_string_field(profile_data, 1, 'Updated Name')
  WHERE user_id = 123;
  
  COMMIT;
ELSE
  ROLLBACK;
END IF;
```

## Integration with Applications

### 1. Using with ORM/Query Builders
```sql
-- Create views for easier ORM integration
CREATE VIEW user_profiles AS
SELECT 
  user_id,
  profile_data,
  pb_message_get_string_field(profile_data, 1, '') AS name,
  pb_message_get_int32_field(profile_data, 2, 0) AS age,
  pb_message_get_string_field(profile_data, 3, '') AS email
FROM users;

-- Now your ORM can work with user_profiles view
```

### 2. Application-Level Helpers
```python
# Python example: helper functions
def get_user_name(cursor, user_id):
    cursor.execute("""
        SELECT pb_message_get_string_field(profile_data, 1, 'Unknown')
        FROM users WHERE user_id = %s
    """, (user_id,))
    return cursor.fetchone()[0]

def update_user_age(cursor, user_id, new_age):
    cursor.execute("""
        UPDATE users 
        SET profile_data = pb_message_set_int32_field(profile_data, 2, %s)
        WHERE user_id = %s
    """, (new_age, user_id))
```

## Troubleshooting Guide

### Common Issues and Solutions

**Issue**: Function returns NULL unexpectedly
```sql
-- Check if field exists first
SELECT pb_message_has_string_field(profile_data, 1) FROM users WHERE user_id = 123;

-- Check for special float values (inf, -inf, nan)
SELECT pb_message_get_float_field(profile_data, 8, -999.0) FROM users WHERE user_id = 123;
-- If returns NULL, the field contains inf/-inf/nan
```

**Issue**: "Index out of range" error with repeated fields
```sql
-- Always check count first
SELECT 
  pb_message_get_repeated_string_field_count(profile_data, 4) AS count,
  CASE 
    WHEN pb_message_get_repeated_string_field_count(profile_data, 4) > 0 THEN
      pb_message_get_repeated_string_field_element(profile_data, 4, 0)
    ELSE
      'No items'
  END AS first_item
FROM users WHERE user_id = 123;
```

**Issue**: Performance problems with large datasets
```sql
-- Extract to regular columns using triggers (stored functions can't be used in generated columns)
ALTER TABLE users ADD COLUMN name_extracted VARCHAR(255) NOT NULL DEFAULT '';

DELIMITER $$
CREATE TRIGGER users_extract_name_insert BEFORE INSERT ON users FOR EACH ROW
BEGIN
  SET NEW.name_extracted = pb_message_get_string_field(NEW.profile_data, 1, '');
END$$
CREATE TRIGGER users_extract_name_update BEFORE UPDATE ON users FOR EACH ROW
BEGIN
  SET NEW.name_extracted = pb_message_get_string_field(NEW.profile_data, 1, '');
END$$
DELIMITER ;

CREATE INDEX idx_users_name_extracted ON users(name_extracted);
```

**Issue**: Memory issues with stored_program_cache
```sql
-- Increase cache size if you get "Function does not exist" errors
SET GLOBAL stored_program_cache = 512;  -- Default is 256
```

## Testing Strategies

### 1. Unit Testing Protobuf Logic
```sql
-- Test data setup
INSERT INTO test_users (user_id, profile_data) VALUES 
(1, pb_message_set_string_field(pb_message_set_int32_field(pb_message_new(), 2, 25), 1, 'Test User'));

-- Verify field extraction
SELECT 
  user_id,
  pb_message_get_string_field(profile_data, 1, '') = 'Test User' AS name_correct,
  pb_message_get_int32_field(profile_data, 2, 0) = 25 AS age_correct
FROM test_users 
WHERE user_id = 1;
```

### 2. Schema Migration Testing
```sql
-- Test backward compatibility
SELECT 
  user_id,
  pb_message_has_string_field(profile_data, 5) AS has_deprecated_field,
  pb_message_has_string_field(profile_data, 6) AS has_new_field
FROM users 
LIMIT 10;
```

This guide provides practical, copy-pasteable solutions for real-world protobuf usage in MySQL. Each example includes error handling and performance considerations that developers actually need.

## 🔗 Related Documentation

- **Need function syntax?** → [Quick Reference](quick-reference.md) has cheat sheets for all operations
- **Want complete function list?** → [Function Reference](function-reference.md) documents all available operations  
- **Advanced topics?** → [API Guide](api-guide.md) covers indexing, validation, and complex scenarios