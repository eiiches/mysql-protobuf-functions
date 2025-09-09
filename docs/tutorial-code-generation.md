# Tutorial: Code Generation for Schema-aware Functions

This tutorial walks you through generating MySQL stored functions for your protobuf schemas using `protoc-gen-mysql`. These generated functions provide schema-aware accessors with intuitive field names instead of field numbers, making your code more readable and maintainable.

## Prerequisites

- MySQL 8.0.17+ or Aurora MySQL 3.04.0+ running
- Core MySQL protobuf functions installed (see [Installation Guide](installation.md))
- Your protobuf schema files (.proto files)
- `protoc` compiler installed

## Step 1: Install the Code Generator

Install the `protoc-gen-mysql` plugin:

```bash
go install github.com/eiiches/mysql-protobuf-functions/cmd/protoc-gen-mysql@latest
```

Verify the installation:

```bash
protoc-gen-mysql --help
```

## Step 2: Prepare Your Protobuf Schema

For this tutorial, let's use a simple `person.proto` schema:

```protobuf
syntax = "proto3";

import "google/protobuf/timestamp.proto";

message Person {
  string name = 1;
  int32 id = 2;
  string email = 3;

  enum PhoneType {
    PHONE_TYPE_UNSPECIFIED = 0;
    PHONE_TYPE_MOBILE = 1;
    PHONE_TYPE_HOME = 2;
    PHONE_TYPE_WORK = 3;
  }

  message PhoneNumber {
    string number = 1;
    PhoneType type = 2;
  }

  repeated PhoneNumber phones = 4;
  google.protobuf.Timestamp last_updated = 5;
}
```

## Step 3: Generate Schema Functions

### Basic Generation

Generate a descriptor set function for your schema:

```bash
protoc --mysql_out=. --mysql_opt=name=person_schema person.proto
```

This creates `person_schema.pb.sql` containing a function that returns the schema descriptor.

### Advanced Generation with Method Functions

Generate both descriptor set and method functions with well-known type support:

```bash
protoc --mysql_out=. \
       --mysql_opt=name=person_schema,include_wkt=true,prefix_map="google.protobuf=pb_" \
       person.proto
```

This creates:
- `person_schema.pb.sql` - Schema descriptor function
- `person_schema_accessors.pb.sql` - Generated accessor functions

## Step 4: Load Functions into MySQL

Load the generated SQL files into your database:

```bash
# Load schema descriptor function
mysql -u your_username -p your_database < person_schema.pb.sql

# If you generated method functions, load them too
mysql -u your_username -p your_database < person_schema_accessors.pb.sql
```

## Step 5: Using the Generated Functions

### Working with Schema-aware Functions

Now you can use intuitive function names instead of field numbers:

```sql
-- Create a new Person message
SELECT person_new() AS empty_person;

-- Set fields using generated accessor functions
SET @person = person_new();
SET @person = person_set_name(@person, 'John Doe');
SET @person = person_set_id(@person, 12345);
SET @person = person_set_email(@person, 'john.doe@example.com');

-- Get field values
SELECT person_get_name(@person) AS name;
SELECT person_get_id(@person) AS id;
SELECT person_get_email(@person) AS email;
```

### Working with Nested Messages and Enums

```sql
-- Create a phone number
SET @phone = person_phonenumber_new();
SET @phone = person_phonenumber_set_number(@phone, '+1-555-0123');
SET @phone = person_phonenumber_set_type(@phone, 1); -- PHONE_TYPE_MOBILE

-- Add phone to person's phone list
SET @person = person_add_phones(@person, @phone);

-- Get phone information
SELECT person_count_phones(@person) AS phone_count;
SELECT person_get_phones(@person, 0) AS first_phone;

-- Convert enum to string name
SELECT person_phonetype_to_string(1) AS phone_type_name; -- Returns 'PHONE_TYPE_MOBILE'
```

### Working with Well-known Types

If you included well-known types (`include_wkt=true`), you can work with Timestamp:

```sql
-- Create and set timestamp
SET @timestamp = pb_timestamp_new();
SET @timestamp = pb_timestamp_set_seconds(@timestamp, UNIX_TIMESTAMP());
SET @person = person_set_last_updated(@person, @timestamp);

-- Convert timestamp to ISO string via JSON
SELECT pb_timestamp_to_json(person_get_last_updated(@person)) AS last_updated_iso;
```

### JSON Conversion

Convert between protobuf messages and JSON:

```sql
-- Convert message to ProtoJSON
SELECT person_to_json(@person) AS person_json;

-- Parse ProtoJSON back to message
SELECT person_from_json('{"name":"Jane Doe","id":67890}') AS person_from_json;

-- Convert to binary protobuf format
SELECT person_to_message(@person) AS binary_protobuf;
```

## Step 6: Advanced Configuration Options

### Custom Function Prefixes

Avoid MySQL's 64-character function name limit and organize functions with custom prefixes:

```bash
protoc --mysql_out=. \
       --mysql_opt=name=schema,prefix_map="google.protobuf=pb_,Person=usr_,Person.PhoneNumber=phone_" \
       person.proto
```

This generates functions like:
- `usr_new()` instead of `person_new()`
- `phone_new()` instead of `person_phonenumber_new()`
- `pb_timestamp_new()` for well-known types

### File Organization Options

Control how generated method functions are organized:

```bash
# Single file (default for standalone mode)
protoc --mysql_out=. --mysql_opt=name=schema,file_naming_strategy=single person.proto

# Preserve directory structure
protoc --mysql_out=. --mysql_opt=name=schema,file_naming_strategy=preserve person.proto

# Flatten into single directory (default for plugin mode)
protoc --mysql_out=. --mysql_opt=name=schema,file_naming_strategy=flatten person.proto
```

### Descriptor-only Generation

Generate only the schema descriptor function without method functions:

```bash
protoc --mysql_out=. --mysql_opt=name=schema,generate_methods=false person.proto
```

Use this when you only need schema information for `pb_message_to_json()` and similar functions.

## Step 7: Using Schema Functions for JSON Conversion

Even without generated method functions, you can use the schema for powerful JSON conversions:

```sql
-- Parse protobuf binary using schema awareness
SELECT pb_message_to_json(
    person_schema(),
    '.Person',
    your_binary_protobuf_data
) AS parsed_json;

-- Convert JSON to protobuf binary
SELECT pb_message_from_json(
    person_schema(),
    '.Person',
    '{"name":"Alice","id":999,"email":"alice@example.com"}'
) AS binary_data;
```

## Common Patterns

### Validation and Error Handling

```sql
-- Check if required fields are set (proto2 only)
SELECT person_has_name(@person) AS has_name;

-- Use default values for optional fields
SELECT person_get_email__or(@person, 'no-email@domain.com') AS email_with_default;

-- Validate enum values
SELECT CASE 
    WHEN person_phonetype_to_string(5) IS NULL THEN 'Invalid enum value'
    ELSE 'Valid enum value'
END AS validation_result;
```

### Working with Arrays and Maps

```sql
-- Process repeated fields
SET @phone_count = person_count_phones(@person);

-- Iterate through repeated field (MySQL doesn't have loops, but you can use recursive CTEs)
WITH RECURSIVE phone_list(idx, phone_data) AS (
    SELECT 0, person_get_phones(@person, 0)
    WHERE person_count_phones(@person) > 0
    UNION ALL
    SELECT idx + 1, person_get_phones(@person, idx + 1)
    FROM phone_list
    WHERE idx + 1 < person_count_phones(@person)
)
SELECT idx, person_phonenumber_get_number(phone_data) AS phone_number
FROM phone_list;
```

## Performance Considerations

- Generated functions add a layer of abstraction over low-level field access
- For high-performance scenarios, consider using low-level functions directly
- Schema descriptor functions are deterministic and can be cached
- Index binary protobuf columns for better query performance

## Next Steps

- Explore the [Generated Functions Documentation](mysql-generated-functions.md) for complete function reference
- Learn about [Advanced Usage](advanced-usage.md) for indexing and performance optimization
- Check out [JSON Integration Tutorial](tutorial-json.md) for more JSON conversion techniques

## Troubleshooting

### Plugin Not Found
```bash
# If protoc can't find the plugin, specify full path
protoc --plugin=protoc-gen-mysql=/path/to/protoc-gen-mysql \
       --mysql_out=. --mysql_opt=name=schema \
       schema.proto
```

### Function Name Too Long
```bash
# Use prefix mapping to shorten function names
protoc --mysql_out=. \
       --mysql_opt=name=schema,prefix_map="com.company.very.long.package=short_" \
       schema.proto
```

### Missing Well-Known Types
```bash
# Include well-known types if your schema imports them
protoc --mysql_out=. --mysql_opt=name=schema,include_wkt=true schema.proto
```