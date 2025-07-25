# MySQL Protocol Buffers Function Reference

*Complete API reference with functions for all 17 protobuf types*

## 📚 Navigation

- **👈 [Documentation Home](README.md)** - Choose your learning path
- **🎯 [User Guide](user-guide.md)** - Problem-focused examples and solutions  
- **⚡ [Quick Reference](quick-reference.md)** - Function syntax cheat sheet
- **🔬 [API Guide](api-guide.md)** - Advanced technical details

---

This document provides a comprehensive reference for all MySQL Protocol Buffers functions available in this library. The library provides extensive functionality for working with protobuf messages directly within MySQL databases.

## Function Organization

The library provides a consistent set of operations across all 17 protobuf types (`bool`, `int32`, `int64`, `uint32`, `uint64`, `sint32`, `sint64`, `fixed32`, `fixed64`, `sfixed32`, `sfixed64`, `float`, `double`, `string`, `bytes`, `enum`, `message`) and both representations (binary messages and Wire JSON). Each core operation (get, set, has, clear, add, etc.) is available for every applicable type combination.

> **💡 New to this library?** Start with the [User Guide](user-guide.md) for practical examples, or check the [Quick Reference](quick-reference.md) for common operations.

## Function Categories

The library provides three main categories of functions based on their schema requirements:

### 🔧 Low-level Field / Message Operations (No Schema Required)
Functions that work with protobuf **field numbers and types** to read, write, and manipulate individual fields. These are the core functions for working with protobuf data and only require knowing the field numbers and types from your `.proto` definition.

- **Field Access**: `pb_message_get_*_field()`, `pb_message_has_*_field()`
- **Field Manipulation**: `pb_message_set_*_field()`, `pb_message_clear_*_field()`
- **Repeated Fields**: `pb_message_add_repeated_*`, `pb_message_get_repeated_*_count()`
- **Wire Format**: `pb_message_to_wire_json()`, `pb_wire_json_*` functions
- **Message Creation**: `pb_message_new()`, basic message operations

Low-level functions do not recognize oneof groups or map fields:

- **Oneof Groups**: If multiple fields within a oneof group are set, `pb_message_get_*_field()` will return a value for all of them. `pb_message_set_*_field()` does not clear other fields in a oneof group. Oneof semantics are not enforced.
- **Map Fields**: A map is encoded as repeated messages on the wire and must be accessed accordingly using the repeated message field functions.

### 📋 Schema Management (Schema Loading)
Functions for loading and managing compiled protobuf schemas (FileDescriptorSet) in the database.

- **Loading**: `pb_descriptor_set_load()`, `pb_descriptor_set_delete()`
- **Querying**: `pb_descriptor_set_exists()`, `pb_descriptor_set_contains_*_type()`

### 🔄 JSON Conversion (Schema Required)
Functions that convert protobuf messages to human-readable JSON using field names. These require a loaded schema to map field numbers to field names.

- **Message to JSON**: `pb_message_to_json()`
- **Well-Known Types**: `pb_timestamp_to_json()`, `pb_duration_to_json()`, etc.

> **Most users only need low-level field operations** for querying and manipulating protobuf data. Schema-dependent functions are primarily for debugging and inspection.

## Table of Contents

1. [Core Message Operations](#core-message-operations) *— No schema required*
2. [Low-level Field Access](#field-access-operations) *— No schema required*
3. [Low-level Field Manipulation](#field-setting-operations) *— No schema required*
4. [Repeated Field Operations](#repeated-field-operations) *— No schema required*
5. [Bulk Operations](#bulk-operations) *— No schema required*
6. [Wire Format Operations](#wire-format-operations) *— No schema required*
7. [JSON Conversion](#json-conversion) *— Schema required*
8. [Schema Management](#schema-management) *— Schema loading/management*
9. [Utility Functions](#utility-functions) *— No schema required*

---

## Core Message Operations

### Message Creation and Conversion

#### `pb_message_new() -> LONGBLOB`
Creates a new empty protobuf message.

**Returns:** Empty protobuf message as LONGBLOB

**Example:**
```sql
SET @msg = pb_message_new();
```

#### `pb_message_to_wire_json(message LONGBLOB) -> JSON`
Converts a protobuf-encoded BLOB into a Wire JSON representation for efficient field access and manipulation. The Wire JSON format allows multiple operations without repeated parsing of the binary message.

**Performance optimization:** When performing 2 or more operations on a message, convert to Wire JSON first, perform all operations, then convert back to binary. This pattern avoids the overhead of parsing and serializing the message for each operation.

**Parameters:**
- `message` (LONGBLOB): Raw protobuf-encoded data

**Returns:** JSON representation of the wire format

**Example:**
```sql
SELECT pb_message_to_wire_json(_binary X'0A0B696E7433325F6669656C64180120012805520A696E7433324669656C64');
```

**Sample Output:**
```json
{
  "1": [{"i": 0, "n": 1, "t": 2, "v": "aW50MzJfZmllbGQ="}],
  "3": [{"i": 1, "n": 3, "t": 0, "v": 1}],
  "4": [{"i": 2, "n": 4, "t": 0, "v": 1}],
  "5": [{"i": 3, "n": 5, "t": 0, "v": 5}],
  "10": [{"i": 4, "n": 10, "t": 2, "v": "aW50MzJGaWVsZA=="}]
}
```

**Wire Format JSON Structure:**
- Field numbers are JSON object keys (e.g., "1", "3", "4")
- Each field contains an array of wire format elements
- Each element has: `i` (index), `n` (field number), `t` (wire type), `v` (value)
- Wire types: 0=varint, 1=fixed64, 2=length-delimited, 5=fixed32
- Values are base64-encoded for length-delimited fields

#### `pb_wire_json_new() -> JSON`
Creates a new empty wire format JSON object.

**Returns:** Empty JSON object for wire format manipulation

#### `pb_wire_json_to_message(wire_json JSON) -> LONGBLOB`
Converts a wire format JSON object back to a protobuf message.

**Parameters:**
- `wire_json` (JSON): The wire format JSON

**Returns:** Protobuf message as LONGBLOB

**Important:** The round-trip conversion `pb_wire_json_to_message(pb_message_to_wire_json(message))` always produces exactly the same binary output as the original message, preserving both field values and the ordering of encoded fields. This guarantees complete data integrity when using Wire JSON for intermediate processing.

### Message Inspection

#### `pb_message_show_wire_format(buf LONGBLOB)`
Displays the wire format structure of a protobuf message (debugging utility).

**Parameters:**
- `buf` (LONGBLOB): The protobuf message to inspect

#### `pb_wire_json_as_table(wire_json JSON)`
Displays the wire format JSON as a table (debugging utility).

**Parameters:**
- `wire_json` (JSON): The wire format JSON to display

---

## Low-level Field Access Operations

### Single Field Getters

The library provides low-level field access functions for all protobuf types. These functions extract values from specific fields in a message.

#### Pattern: `pb_message_get_[TYPE]_field(message LONGBLOB, field_number INT, default_value [TYPE]) -> [TYPE]`

**Available Types:**
- `bool` -> BOOLEAN
- `bytes` -> LONGBLOB  
- `double` -> DOUBLE
- `enum` -> INT
- `fixed32` -> INT UNSIGNED
- `fixed64` -> BIGINT UNSIGNED
- `float` -> FLOAT
- `int32` -> INT
- `int64` -> BIGINT
- `message` -> LONGBLOB
- `sfixed32` -> INT
- `sfixed64` -> BIGINT
- `sint32` -> INT
- `sint64` -> BIGINT
- `string` -> LONGTEXT
- `uint32` -> INT UNSIGNED
- `uint64` -> BIGINT UNSIGNED

**Parameters:**
- `message` (LONGBLOB): The protobuf message
- `field_number` (INT): The field number as defined in the .proto schema
- `default_value` ([TYPE]): Value to return when field is not present. For proto3 without explicit field presence, this should be set to the [default value](https://protobuf.dev/programming-guides/proto3/#default) defined by the Protobuf specification. For proto2 messages, this should also be set to the [default value](https://protobuf.dev/programming-guides/proto2/#default) unless the field has an explicit `default` option.

**Important Notes:**
- The field number must match the one used in the `.proto` schema definition
- This function does not perform schema validation; it assumes the caller knows the correct field number and expected type
- MySQL does not support `+inf`, `-inf`, or `NaN`. Therefore, `float` and `double` variants return `NULL` instead if the corresponding field contains any of these values
- For better performance when accessing multiple fields, use `pb_message_to_wire_json()` to parse the message once, then call `pb_wire_json_get_{type}_field()` for each field

**Example:**
```sql
SELECT pb_message_get_int32_field(@msg, 1, 0) AS id;
SELECT pb_message_get_string_field(@msg, 2, 'Unknown') AS name;
```

### Field Existence Check

#### Pattern: `pb_message_has_[TYPE]_field(message LONGBLOB, field_number INT) -> BOOLEAN`

Checks whether a specific field is present in a Protobuf-encoded BLOB. This function is used to determine whether a field with [Field Presence](https://protobuf.dev/programming-guides/field_presence/) tracking is set in the encoded message.

**Available for all types listed above**

**Parameters:**
- `message` (LONGBLOB): The protobuf message
- `field_number` (INT): The field number as defined in the .proto schema

**Returns:** A BOOLEAN indicating whether the specified field is present in the encoded message

**Important Notes:**
- In `proto3`, presence tracking for scalar fields is only available when the field is declared with the `optional` keyword
- Using this function on `repeated` fields is an error. Use `pb_message_get_repeated_{type}_field_count()` instead to check repeated field presence

**Example:**
```sql
SELECT pb_message_has_int32_field(@msg, 1);
SELECT pb_message_has_string_field(@msg, 2);
```

### Repeated Field Access

#### Pattern: `pb_message_get_repeated_[TYPE]_field_element(message LONGBLOB, field_number INT, repeated_index INT) -> [TYPE]`

Retrieves a specific element from a repeated field at the given zero-based index.

**Available for all types listed above**

**Parameters:**
- `message` (LONGBLOB): The protobuf message
- `field_number` (INT): The field number as defined in the .proto schema
- `repeated_index` (INT): The zero-based index of the element to retrieve

**Returns:** The value of the requested element, interpreted as the corresponding SQL type

**Notes:**
- If `repeated_index` exceeds the number of available elements, the function raises an "index out of range" error
- MySQL does not support `+inf`, `-inf`, or `NaN`. Therefore, `float` and `double` variants return `NULL` instead if the corresponding field contains any of these values

#### Pattern: `pb_message_get_repeated_[TYPE]_field_count(message LONGBLOB, field_number INT) -> INT`

Returns the number of elements present in a repeated field of a Protobuf-encoded BLOB. To retrieve the value of each element, use `pb_message_get_repeated_{type}_field_element()`.

**Parameters:**
- `message` (LONGBLOB): The protobuf message
- `field_number` (INT): The field number as defined in the .proto schema

**Returns:** An INT representing the number of elements in the specified repeated field

#### Pattern: `pb_message_get_repeated_[TYPE]_field_as_json_array(message LONGBLOB, field_number INT) -> JSON`

Retrieves all elements of a repeated field as a JSON array.

**Parameters:**
- `message` (LONGBLOB): The protobuf message
- `field_number` (INT): The field number as defined in the .proto schema

**Returns:** A JSON array containing all elements of the specified field

**Notes:**
- MySQL does not support `+inf`, `-inf`, or `NaN`. Therefore, `float` and `double` variants return `NULL` instead if the corresponding field contains any of these values

**Example:**
```sql
SELECT pb_message_get_repeated_int32_field_count(@msg, 3) AS count;
SELECT pb_message_get_repeated_string_field_element(@msg, 4, 0) AS first_item;
SELECT pb_message_get_repeated_int32_field_as_json_array(@msg, 3) AS all_numbers;
```

---

## Low-level Field Manipulation Operations

### Single Field Setters

#### Pattern: `pb_message_set_[TYPE]_field(message LONGBLOB, field_number INT, value [TYPE]) -> LONGBLOB`

Sets a field value in a protobuf message.

**Available for all types listed in Low-level Field Access Operations**

**Parameters:**
- `message` (LONGBLOB): The protobuf message
- `field_number` (INT): The field number to set
- `value` ([TYPE]): The value to set

**Returns:** Modified protobuf message

**Example:**
```sql
SET @msg = pb_message_set_int32_field(@msg, 1, 42);
SET @msg = pb_message_set_string_field(@msg, 2, 'Hello World');
```

### Field Clearing

#### Pattern: `pb_message_clear_[TYPE]_field(message LONGBLOB, field_number INT) -> LONGBLOB`

Clears a field from a protobuf message.

**Available for all types**

**Example:**
```sql
SET @msg = pb_message_clear_int32_field(@msg, 1);
```

---

## Repeated Field Operations

### Adding Elements

#### Pattern: `pb_message_add_repeated_[TYPE]_field_element(message LONGBLOB, field_number INT, value [TYPE], use_packed BOOLEAN) -> LONGBLOB`

Adds a single element to a repeated field.

**Parameters:**
- `message` (LONGBLOB): The protobuf message
- `field_number` (INT): The field number
- `value` ([TYPE]): The value to add
- `use_packed` (BOOLEAN): Whether to use packed encoding (for numeric types)

**Note:** `use_packed` parameter is not available for `bytes`, `string`, and `message` types.

### Setting Elements

#### Pattern: `pb_message_set_repeated_[TYPE]_field_element(message LONGBLOB, field_number INT, index INT, value [TYPE]) -> LONGBLOB`

Sets a specific element in a repeated field.

### Inserting Elements

#### Pattern: `pb_message_insert_repeated_[TYPE]_field_element(message LONGBLOB, field_number INT, index INT, value [TYPE]) -> LONGBLOB`

Inserts an element at a specific position in a repeated field.

### Removing Elements

#### Pattern: `pb_message_remove_repeated_[TYPE]_field_element(message LONGBLOB, field_number INT, index INT) -> LONGBLOB`

Removes an element from a repeated field.

### Setting Entire Repeated Fields

#### Pattern: `pb_message_set_repeated_[TYPE]_field(message LONGBLOB, field_number INT, value_array JSON, use_packed BOOLEAN) -> LONGBLOB`

Sets all elements of a repeated field from a JSON array.

### Clearing Repeated Fields

#### Pattern: `pb_message_clear_repeated_[TYPE]_field(message LONGBLOB, field_number INT) -> LONGBLOB`

Clears all elements from a repeated field.

**Example:**
```sql
SET @msg = pb_message_add_repeated_int32_field_element(@msg, 3, 10, TRUE);
SET @msg = pb_message_set_repeated_int32_field_element(@msg, 3, 0, 20);
SET @msg = pb_message_insert_repeated_int32_field_element(@msg, 3, 1, 15);
SET @msg = pb_message_remove_repeated_int32_field_element(@msg, 3, 0);
```

---

## Bulk Operations

### Adding All Elements

#### Pattern: `pb_message_add_all_repeated_[TYPE]_field_elements(message LONGBLOB, field_number INT, value_array JSON, use_packed BOOLEAN) -> LONGBLOB`

Adds multiple elements to a repeated field from a JSON array.

**Parameters:**
- `message` (LONGBLOB): The protobuf message
- `field_number` (INT): The field number
- `value_array` (JSON): JSON array of values to add
- `use_packed` (BOOLEAN): Whether to use packed encoding (for numeric types)

**Example:**
```sql
SET @msg = pb_message_add_all_repeated_int32_field_elements(@msg, 3, JSON_ARRAY(1, 2, 3, 4, 5), TRUE);
SET @msg = pb_message_add_all_repeated_string_field_elements(@msg, 4, JSON_ARRAY('a', 'b', 'c'));
```

---

## Wire Format Operations

The library provides high-level wire format operations through the message manipulation functions documented above. For direct wire format manipulation, use the wire format JSON operations described in the [Wire Format JSON Operations](#wire-format-json-operations) section.

---

## JSON Conversion

### Message to JSON

#### `pb_message_to_json(set_name VARCHAR(64), full_type_name VARCHAR(512), message LONGBLOB) -> JSON`
Converts a Protobuf-encoded BLOB into a JSON object using the specified type and schema set. This function deserializes a Protocol Buffers message into a human-readable and structured JSON format.

**Parameters:**
- `set_name` (VARCHAR(64)): A VARCHAR(64) specifying the name of the schema set containing the compiled Protobuf descriptors
- `full_type_name` (VARCHAR(512)): A VARCHAR(512) representing the fully-qualified name of the Protobuf message type (e.g., `.my.package.MessageType`). A fully-qualified name always starts with a dot.
- `message` (LONGBLOB): A LONGBLOB containing the serialized Protobuf message to be converted

**Returns:** A JSON object that represents the Protobuf message, with field names and values corresponding to those defined in the Protobuf schema

**Important Usage Notes:**
- This function is primarily intended for debugging or inspection. It should not be used in production code
- JSON relies on field names rather than field numbers, which compromises a key benefit of Protocol Buffers: the ability to rename fields without breaking compatibility
- Instead, applications should read and decode the raw Protobuf-encoded BLOB directly
- Each reader can use its own descriptor set (i.e., reader schema), allowing it to remain unaffected by schema changes and independently control when to adopt schema updates

**Errors:**
- Returns an error if the set_name or full_type_name cannot be resolved
- Returns an error if the message is not a valid serialized message of the given type

**Example:**
```sql
SELECT pb_message_to_json('my_schema', '.com.example.Person', @msg);
```

### Well-Known Type Conversions

The library includes special handling for Protocol Buffers Well-Known Types. These conversions are handled automatically when using `pb_message_to_json()` with appropriate schema information.

---

## Schema Management

### Descriptor Set Operations

#### `pb_descriptor_set_load(set_name VARCHAR(64), file_descriptor_set LONGBLOB)`
Loads a compiled Protobuf FileDescriptorSet into internal database tables. This procedure registers the schema information required for operations that depend on message descriptors, such as `pb_message_to_json()`.

**Parameters:**
- `set_name` (VARCHAR(64)): A user-defined identifier that is used to distinguish file descriptor sets. This allows multiple descriptor sets to coexist without conflict.
- `file_descriptor_set` (LONGBLOB): A binary-encoded FileDescriptorSet, typically generated using `protoc --descriptor_set_out` or `buf build -o ${name}.binpb`. The input must conform to the `google.protobuf.FileDescriptorSet` message format.

**Important Notes:**
- This procedure is only intended for use cases where Protobuf schema information is required
- Getters and has functions (`pb_message_get_{type}_field()`, `pb_message_get_{type}_field_count()`, `pb_message_has_{type}_field()`) don't require this procedure
- If a descriptor set with the same `set_name` already exists, the procedure will fail

#### `pb_descriptor_set_delete(set_name VARCHAR(64))`
Deletes a previously loaded descriptor set from internal tables. This procedure is used to clean up resources created by `pb_descriptor_set_load()`.

**Parameters:**
- `set_name` (VARCHAR(64)): The identifier used when the descriptor set was loaded

**Notes:**
- If the specified descriptor set does not exist, the procedure performs no action

#### `pb_descriptor_set_exists(set_name VARCHAR(64)) -> BOOLEAN`
Checks if a descriptor set exists in the database.

**Parameters:**
- `set_name` (VARCHAR(64)): Name of the descriptor set

**Returns:** TRUE if the descriptor set exists, FALSE otherwise

#### `pb_descriptor_set_contains_message_type(set_name VARCHAR(64), full_type_name VARCHAR(512)) -> BOOLEAN`
Checks if a descriptor set contains a specific message type.

**Parameters:**
- `set_name` (VARCHAR(64)): Name of the descriptor set
- `full_type_name` (VARCHAR(512)): Full type name to check

**Returns:** TRUE if the type exists, FALSE otherwise

#### `pb_descriptor_set_contains_enum_type(set_name VARCHAR(64), full_type_name VARCHAR(512)) -> BOOLEAN`
Checks if a descriptor set contains a specific enum type.

**Parameters:**
- `set_name` (VARCHAR(64)): Name of the descriptor set
- `full_type_name` (VARCHAR(512)): Full type name to check

**Returns:** TRUE if the enum type exists, FALSE otherwise

---

## Utility Functions

The library includes internal utility functions for type conversion and data manipulation. These are handled automatically by the public API functions and do not need to be called directly.

---

## Wire Format JSON Operations

The library provides support for manipulating protobuf messages through a JSON-based wire format. Wire format JSON is essential for performance optimization when performing multiple operations on the same message, as it avoids repeated parsing and serialization overhead.

### Public Wire Format JSON Functions

The public API provides the following wire format JSON functions:

#### `pb_wire_json_new() -> JSON`
Creates a new empty wire format JSON object.

#### `pb_wire_json_to_message(wire_json JSON) -> LONGBLOB`
Converts a wire format JSON object to a protobuf message.

#### `pb_message_to_wire_json(buf LONGBLOB) -> JSON`
Converts a protobuf message to its wire format JSON representation.

#### `pb_wire_json_as_table(wire_json JSON)`
Displays the wire format JSON as a table for debugging purposes.

### Wire Format JSON Performance Benefits

Use Wire JSON when you need to:
- Perform 2 or more operations on the same message
- Transform complex messages with many field updates
- Optimize performance for multiple operations on any message size

Example pattern:
```sql
-- Instead of chaining pb_message_* functions (parses N times):
SET @result = pb_message_set_field3(pb_message_set_field2(pb_message_set_field1(data)));

-- Use Wire JSON (parses once):
SET @wire = pb_message_to_wire_json(data);
SET @wire = pb_wire_json_set_field1(@wire);
SET @wire = pb_wire_json_set_field2(@wire);
SET @wire = pb_wire_json_set_field3(@wire);
SET @result = pb_wire_json_to_message(@wire);
```

---

## Usage Examples

### Basic Message Creation and Manipulation

```sql
-- Create a new message
SET @msg = pb_message_new();

-- Set some fields
SET @msg = pb_message_set_int32_field(@msg, 1, 42);
SET @msg = pb_message_set_string_field(@msg, 2, 'Hello World');

-- Add repeated elements
SET @msg = pb_message_add_repeated_int32_field_element(@msg, 3, 10, TRUE);
SET @msg = pb_message_add_repeated_int32_field_element(@msg, 3, 20, TRUE);

-- Check field existence
SELECT pb_message_has_int32_field(@msg, 1) AS has_field_1;

-- Get field values
SELECT pb_message_get_int32_field(@msg, 1) AS field_1_value;
SELECT pb_message_get_string_field(@msg, 2) AS field_2_value;

-- Get repeated field count and elements
SELECT pb_message_get_repeated_int32_field_count(@msg, 3) AS field_3_count;
SELECT pb_message_get_repeated_int32_field_element(@msg, 3, 0) AS field_3_element_0;
```

### Working with Schema and JSON Conversion

```sql
-- Load a descriptor set (assuming you have a FileDescriptorSet)
CALL pb_descriptor_set_load('my_schema', @descriptor_set_blob);

-- Convert message to JSON using schema
SELECT pb_message_to_json('my_schema', 'com.example.Person', @msg) AS json_output;

-- Check if schema contains a type
SELECT pb_descriptor_set_contains_message_type('my_schema', 'com.example.Person') AS has_type;
```

### Wire Format JSON Operations

```sql
-- Create wire format JSON
SET @wire_json = pb_wire_json_new();

-- Convert message to wire format JSON
SET @wire_json = pb_message_to_wire_json(@msg);

-- Convert back to message
SET @msg = pb_wire_json_to_message(@wire_json);

-- Display wire format for debugging
CALL pb_wire_json_as_table(@wire_json);
```

---

## Function Naming Conventions

- **`pb_message_*`**: High-level message operations with full type safety
- **`pb_wire_json_*`**: Wire format JSON operations (public API)
- **`pb_descriptor_set_*`**: Schema management operations

## Notes

- All functions that return modified messages are deterministic and safe to use in expressions
- The `use_packed` parameter is only relevant for numeric types and determines the wire format encoding
- Repeated field indices are zero-based
- Field numbers must be positive integers as per protobuf specification
- The library handles all protobuf wire types: varint, i32, i64, and length-delimited
- For advanced use cases requiring direct wire format manipulation, consider using the wire format JSON functions

---

## 🔗 Need Practical Examples?

This reference provides complete function documentation, but you might also want:

- **🎯 [User Guide](user-guide.md)** - Real-world problems and solutions with copy-paste examples
- **⚡ [Quick Reference](quick-reference.md)** - Common operations cheat sheet with use cases
- **🔬 [API Guide](api-guide.md)** - Advanced topics like indexing, validation, and performance

> **💡 Learning tip:** Start with the [User Guide](user-guide.md) to see these functions in action, then return here for detailed parameter information.