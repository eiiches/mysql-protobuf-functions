# Protobuf Format Representations

*Comprehensive guide to the four protobuf format representations in MySQL Protocol Buffers Functions*

This document provides detailed information about the different protobuf format representations supported by the MySQL Protocol Buffers Functions library, their characteristics, use cases, and conversion methods.

## Format Overview

Building on MySQL's JSON capabilities, the library provides four distinct protobuf format representations, each designed to overcome MySQL's type system limitations while optimizing for different use cases:

1. **Binary/Message Format** - Standard protobuf binary encoding ([Protocol Buffers specification](https://protobuf.dev/programming-guides/encoding/))
2. **WireJSON Format** - *Library-specific* intermediate JSON format for efficient multi-operations
3. **ProtoNumberJSON Format** - *Library-specific internal format* for schema-evolution-safe JSON with field numbers
4. **ProtoJSON Format** - Official protobuf JSON mapping ([Protocol Buffers JSON specification](https://protobuf.dev/programming-guides/json/))

## Why Multiple Format Representations?

### MySQL's Type System Limitations

Unlike other programming languages, MySQL has a limited type system that lacks efficient complex data structures:

- **Primitive types only**: `INT`, `VARCHAR`, `DECIMAL`, `BOOLEAN`, etc.
- **Date/Time types**: `DATETIME`, `TIMESTAMP`, `DATE`, etc.
- **JSON support**: Native JSON data type (MySQL 5.7+)
- **No arrays**: No native array or list types
- **No structs/records**: No composite data types or user-defined structures
- **No classes**: No object-oriented data structures

### The Challenge of Working with Protobuf in MySQL

Protocol Buffers messages are complex, hierarchical data structures that don't map naturally to MySQL's simple type system:

```protobuf
message Person {
  string name = 1;
  int32 id = 2; 
  repeated string emails = 3;        // Arrays not supported in MySQL
  Address address = 4;               // Nested messages not supported
  map<string, string> metadata = 5;  // Maps not supported
}
```

### The BLOB Problem

Storing protobuf messages directly as `BLOB` in MySQL creates several challenges:

- **Opaque data**: Cannot inspect or query message contents
- **Inefficient field access**: Must parse entire message for each field operation
- **No partial updates**: Must deserialize, modify, and reserialize entire message
- **Multiple operations**: Repeated parsing/serialization overhead for batch operations
- **Query limitations**: Cannot use MySQL's powerful query capabilities on message fields

### The Solution: Multiple Representations

To bridge the gap between protobuf's rich data structures and MySQL's limited type system, this library provides four specialized format representations, each leveraging MySQL's JSON capabilities in different ways:

1. **Binary/Message**: Standard protobuf format for compatibility and compactness
2. **WireJSON**: *Library-specific* intermediate format for multi-operation processing
3. **ProtoNumberJSON**: *Library-specific internal format* for robust programmatic access
4. **ProtoJSON**: Official standard format for debugging and integration

Each format is designed to work within MySQL's constraints while providing optimal performance and functionality for specific use cases.

## Format Specifications

### 1. Binary/Message Format

**Description:** Standard protobuf binary wire format stored as `LONGBLOB` in MySQL.

**Characteristics:**
- Native protobuf binary encoding (Protocol Buffers wire format)
- Most compact representation
- Direct compatibility with external protobuf systems
- Requires parsing for each field access operation
- Opaque to human inspection

**Storage:** `LONGBLOB` data type in MySQL

**Example:**
```sql
-- Binary message (hex representation)
_binary X'0A0B48656C6C6F20576F726C64'
-- Represents a message with field 1 containing "Hello World"
```

**Use Cases:**
- Single field operations (where parsing overhead is acceptable)
- Direct storage in database tables
- Interoperability with external protobuf systems
- Long-term storage and compatibility
- Minimal storage footprint

**Performance:** Fastest for single operations, but requires reparsing for multiple operations.

**MySQL Limitation Addressed:** Provides a storage mechanism for complex protobuf data within MySQL's primitive type system, though field access requires specialized functions.

---

### 2. WireJSON Format

**Description:** Intermediate JSON format representing the protobuf wire format structure, optimized for efficient multi-operation processing. Binary/Message ↔ WireJSON conversions are completely lossless and always produce the exact same binary data.

**Characteristics:**
- JSON representation of raw wire format elements
- Field numbers as JSON object keys
- Maintains complete wire format information
- Base64-encoded values for length-delimited fields
- Optimized for multiple operations without reparsing

**Structure:**
```json
{
  "1": [{"i": 0, "n": 1, "t": 2, "v": "aW50MzJfZmllbGQ="}],
  "3": [{"i": 1, "n": 3, "t": 0, "v": 1}],
  "4": [{"i": 2, "n": 4, "t": 0, "v": 1}],
  "5": [{"i": 3, "n": 5, "t": 0, "v": 5}],
  "10": [{"i": 4, "n": 10, "t": 2, "v": "aW50MzJGaWVsZA=="}]
}
```

**Wire Format JSON Element Properties:**
- **Field number** (JSON key): Wire field number as string
- **i**: Index within the field (for repeated fields)
- **n**: Field number (duplicates the key)
- **t**: Wire type (0=varint, 1=fixed64, 2=length-delimited, 5=fixed32)
- **v**: Value (base64-encoded for length-delimited fields, numeric for others)

**Performance Benefits:**
- Parse binary message once, perform multiple operations, serialize once
- Essential optimization for 2+ operations on the same message
- Eliminates repeated parsing/serialization overhead
- Direct field access without schema lookup

**Use Cases:**
- Multiple field operations (2+ operations recommended)
- Complex message transformations
- Performance-critical batch updates
- Intermediate processing format

**MySQL Limitation Addressed:** Leverages MySQL's JSON type to create an efficient intermediate representation that allows multiple field operations without the overhead of repeated binary parsing. Provides array-like and object-like access patterns using JSON syntax.

**Example Usage Pattern:**
```sql
-- Performance-optimized pattern for multiple operations
SET @wire = pb_message_to_wire_json(@binary_message);  -- Parse once
SET @wire = pb_wire_json_set_string_field(@wire, 1, 'New Name');
SET @wire = pb_wire_json_set_int32_field(@wire, 2, 30);
SET @wire = pb_wire_json_add_repeated_string_field_element(@wire, 4, 'hobby');
SET @result = pb_wire_json_to_message(@wire);  -- Serialize once
```

---

### 3. ProtoJSON Format

**Description:** Human-readable JSON format using field names, following the [official Protocol Buffers JSON mapping specification](https://protobuf.dev/programming-guides/json/).

> 📖 **Detailed Documentation:** For comprehensive information about ProtoJSON format, see [Tutorial: JSON Integration](tutorial-json.md).

**Characteristics:**
- Uses **field names** as JSON object keys (converted to lowerCamelCase)
- **Enum names** as strings instead of numeric values
- **64-bit integers** as strings (to preserve precision in JavaScript)
- **Timestamp formatting** in RFC 3339 format (e.g., "1972-01-01T10:00:20.021Z")
- **Well-known types** receive special formatting (Any, Struct, Duration, Timestamp, etc.)
- **Bytes fields** are base64-encoded strings
- Follows the [official Protocol Buffers JSON mapping](https://protobuf.dev/programming-guides/json/)
- Requires schema information for field name resolution

**Example:**
```json
{
  "id": 1,
  "name": "Agent Smith", 
  "email": "smith@example.com",
  "phones": [
    {"type": "PHONE_TYPE_WORK", "number": "+81-00-0000-0000"}
  ],
  "lastUpdated": "2025-06-01T12:34:56.789Z"
}
```

**Schema Requirements:** Requires protobuf schema (FileDescriptorSet) to map field numbers to field names.

**Use Cases:**
- **Debugging and inspection** (primary use case)
- Human-readable output for reports and logging
- **Standards compliance** - Official Protocol Buffers JSON format (though interoperability use is discouraged)
- Integration with legacy systems that require ProtoJSON format
- One-time data analysis and exploration

**MySQL Limitation Addressed:** Transforms protobuf's complex nested structures into MySQL's JSON format with human-readable field names, enabling use of MySQL's JSON functions and operators for ad-hoc queries and reporting.

**Important Limitations:**
> ⚠️ **Interoperability Warning:** While ProtoJSON is the official Protocol Buffers JSON standard, it is **strongly discouraged for interoperability purposes**. As noted in the [official Protocol Buffers documentation](https://protobuf.dev/programming-guides/json/): "ProtoJSON format puts your field and enum value names into encoded messages making it much harder to change those names later." For production systems, use the standard binary protobuf format for interoperability. Use ProtoJSON primarily for human consumption (debugging, inspection, reporting) and only when standards compliance is specifically required.

> ⚠️ **Data Representation Limitations:** ProtoJSON cannot represent all protobuf data accurately. Well-known types like `Timestamp` and `Duration` have stricter value ranges in ProtoJSON than in the binary representation. For example, ProtoJSON Timestamps are limited to years 0001-9999, while binary protobuf can represent much larger ranges. Use binary format or ProtoNumberJSON for complete data fidelity.

**Example Usage:**
```sql
-- Convert binary message to human-readable JSON
SELECT pb_message_to_json(
    person_schema(),    -- Schema function
    '.Person',          -- Message type
    pb_data             -- Binary message
) FROM Person;
```

---

### 3. ProtoNumberJSON Format  

**Description:** The most efficient JSON format for field operations, using field numbers as keys for all message types. Optimized for both performance and schema evolution robustness. Think of this as the closest equivalent to in-memory protobuf message representations (struct or class instances) found in other programming languages.

**Key Implementation Detail:** The schema (FileDescriptorSet) itself is stored in ProtoNumberJSON format because:
- **Performance**: Deserializing from protobuf message BLOB would be slower for frequent field access
- **Storage stability**: Avoids field names and enum value names, making it unlikely to be affected by schema evolution
- **Complete data fidelity**: Unlike ProtoJSON, can represent all protobuf data without value range restrictions

> 📖 **Detailed Documentation:** For comprehensive information about ProtoNumberJSON format, see [ProtoNumberJSON README](../internal/protonumberjson/README.md).

**Characteristics:**
- **Most efficient for field operations** - Direct field number access without name resolution
- **Field numbers** used as JSON object keys (e.g., "1", "2", "3") for ALL message types
- **Enum values** serialized as numbers instead of string names
- **64-bit integers** serialized as numbers (not strings)
- **Consistent treatment** - all protobuf messages follow the same approach
- **Well-known types** use field numbers too (no special formatting, avoiding ProtoJSON's value range limitations)
- **Complete data representation** - Can represent all protobuf data without the value range restrictions of ProtoJSON's well-known type formatting
- **Schema evolution safe** - immune to field and enum renames
- **MySQL JSON optimization** - Leverages MySQL's native JSON path operations (`$.1`, `$.2`, etc.)

**Schema Evolution Benefits:**
- **Field rename safety** - Field numbers remain stable across schema changes
- **Enum rename safety** - Numeric values remain stable while names may change
- **Consistency** - Predictable format for all message types
- **Robustness** - Field numbers are immutable once assigned

**Example Transformations:**

| Protobuf Message | ProtoNumberJSON Output |
|------------------|------------------------|
| `StringValue{Value: "test"}` | `{"1": "test"}` |
| `Int64Value{Value: 9223372036854775807}` | `{"1": 9223372036854775807}` |
| `Timestamp{Seconds: 1000, Nanos: 500}` | `{"1": 1000, "2": 500}` |
| `Empty{}` | `{}` |
| `Any{TypeUrl: "type.googleapis.com/...", Value: [...]}` | `{"1": "type.googleapis.com/...", "2": "base64data"}` |
| Regular message with `repeated int32 values = [1, 2, 3]` | `{"1": [1, 2, 3]}` |
| Regular message with `Status status = ACTIVE` (enum value 1) | `{"1": 1}` |

**Use Cases:**
- **Efficient field operations** - Most performant JSON format for field access and manipulation
- **Schema storage** - Used internally for storing FileDescriptorSet schemas for optimal performance
- **Schema evolution robustness** - immune to field and enum renames
- **Consistent handling** of all message types including well-known types  
- **Programmatic processing** where both performance and schema stability are important
- **Data interchange** where field numbers provide better stability than field names

**MySQL Limitation Addressed:** Provides the closest equivalent to in-memory object representations found in other programming languages. Maintains protobuf's field number stability while enabling efficient field access through MySQL's JSON path operations (`JSON_EXTRACT(data, '$.1')`, `JSON_SET(data, '$.2', value)`, etc.).

**Example Usage:**
```sql
-- Convert to ProtoNumberJSON for efficient field operations
SELECT _pb_message_to_number_json(
    person_schema(),    -- Schema function
    '.Person',          -- Message type  
    pb_data             -- Binary message
) FROM Person;
-- Result: {"1": "Agent Smith", "2": 1, "3": "smith@example.com", "4": [{"1": "+81-00-0000-0000", "2": 3}], "5": {"1": 1735743296, "2": 789000000}}

-- Efficient field access using MySQL JSON path operations
SELECT JSON_EXTRACT(proto_number_json, '$.1') AS name,
       JSON_EXTRACT(proto_number_json, '$.2') AS id
FROM person_data;

-- Efficient field updates
UPDATE person_data 
SET proto_number_json = JSON_SET(proto_number_json, '$.1', 'New Name')
WHERE JSON_EXTRACT(proto_number_json, '$.2') = 1;
```

## Format Comparison Table

| Aspect | Binary/Message | WireJSON | ProtoJSON | ProtoNumberJSON |
|--------|----------------|-----------|-----------|-----------------|
| **Field Keys** | N/A | Field numbers | Field names | Field numbers |
| **Enum Values** | N/A | Numbers | String names | Numbers |
| **64-bit Integers** | N/A | Numbers | Strings | Numbers |
| **Well-Known Types** | N/A | Raw wire format | Special formatting | Field numbers |
| **Schema Required** | No | No | Yes | Yes |
| **Human Readable** | No | No | Yes | Partial |
| **Schema Evolution Safe** | Yes | Yes | No | Yes |
| **Performance** | High (single ops) | High (multi ops) | Medium | Medium |
| **Storage Size** | Smallest | Medium | Largest | Medium |
| **Primary Use Case** | Storage, single ops | Multi-operations | Debugging | Schema robustness |

## Format Conversion Matrix

This matrix shows which formats can be converted to others and the functions used:

| From → To | Binary/Message | WireJSON | ProtoJSON | ProtoNumberJSON |
|-----------|----------------|-----------|-----------|-----------------|
| **Binary/Message** | Identity | `pb_message_to_wire_json()` | `pb_message_to_json()`* | `_pb_message_to_number_json()`* |
| **WireJSON** | `pb_wire_json_to_message()` | Identity | `pb_wire_json_to_json()`* | `_pb_wire_json_to_number_json()`* |
| **ProtoJSON** | `pb_json_to_message()`* | `pb_json_to_wire_json()`* | Identity | `_pb_json_to_number_json()`* |
| **ProtoNumberJSON** | `_pb_number_json_to_message()`* | `_pb_number_json_to_wire_json()`* | `_pb_number_json_to_json()`* | Identity |

*\* Requires schema (FileDescriptorSet)*

### Conversion Function Details

#### Schema-Free Conversions (No FileDescriptorSet required)

```sql
-- Binary ↔ WireJSON (bidirectional, no schema needed)
SELECT pb_message_to_wire_json(@binary_message);
SELECT pb_wire_json_to_message(@wire_json);
```

#### Schema-Required Conversions (FileDescriptorSet required)

```sql  
-- Binary → ProtoJSON
SELECT pb_message_to_json(@schema, '.MessageType', @binary_message, NULL, NULL);

-- Binary → ProtoNumberJSON  
SELECT _pb_message_to_number_json(@schema, '.MessageType', @binary_message);

-- WireJSON → ProtoJSON
SELECT pb_wire_json_to_json(@schema, '.MessageType', @wire_json, NULL, NULL);

-- ProtoJSON → Binary
SELECT pb_json_to_message(@schema, '.MessageType', @proto_json, NULL, NULL);

-- ProtoJSON → WireJSON
SELECT pb_json_to_wire_json(@schema, '.MessageType', @proto_json, NULL, NULL);

-- ProtoJSON ↔ ProtoNumberJSON (bidirectional)
SELECT _pb_json_to_number_json(@schema, '.MessageType', @proto_json);
SELECT _pb_number_json_to_json(@schema, '.MessageType', @proto_number_json, true);
```

## Performance Guidelines

### Single Operation
```sql
-- Use binary message functions directly
SELECT pb_message_get_string_field(@binary_message, 1);
```

### Multiple Operations (2+)
```sql
-- Convert to WireJSON once, perform operations, convert back
SET @wire = pb_message_to_wire_json(@binary_message);
SET @wire = pb_wire_json_set_string_field(@wire, 1, 'value1');
SET @wire = pb_wire_json_set_int32_field(@wire, 2, 42);
SET @result = pb_wire_json_to_message(@wire);
```

### Human-Readable Output
```sql
-- Use ProtoJSON for debugging/reporting
SELECT pb_message_to_json(@schema, '.MessageType', @binary_message, NULL, NULL);
```

### Schema Evolution Safety
```sql
-- Use ProtoNumberJSON for robust programmatic access
SELECT _pb_message_to_number_json(@schema, '.MessageType', @binary_message);
```

## Best Practices

### Format Selection Guidelines

**Understanding MySQL's Constraints:**
- MySQL cannot efficiently represent protobuf's arrays, nested objects, or maps natively
- All complex data access must go through function calls or JSON operations
- Multiple operations on binary data require repeated parsing overhead
- Human inspection of binary data is impossible without specialized tools

**Format Selection:**

1. **Binary/Message Format** - Choose when:
   - Performing single field operations
   - Storing protobuf data in database tables
   - Interfacing with external protobuf systems
   - Storage space is critical

2. **WireJSON Format** - Choose when:
   - Performing 2+ operations on the same message
   - Complex message transformations needed
   - Performance optimization is critical
   - Schema information is not available
   - Need to leverage MySQL's JSON operators for field access

3. **ProtoJSON Format** - Choose when:
   - Human inspection/debugging is needed
   - Generating reports or logs
   - One-time data analysis
   - **Standards compliance** is required (though interoperability use is discouraged)
   - Integration with legacy systems that specifically require ProtoJSON
   - Want to use MySQL's JSON functions with familiar field names

4. **ProtoNumberJSON Format** - Choose when:
   - Schema evolution robustness is critical  
   - Field or enum names might change
   - Consistent programmatic processing needed
   - Data interchange requires stability
   - Want to combine MySQL's JSON capabilities with protobuf's field number stability

### Schema Evolution Considerations

- **Safe:** Binary/Message, WireJSON, ProtoNumberJSON - Use field numbers
- **Unsafe:** ProtoJSON - Uses field names that can change

### Performance Optimization

- **Single operation:** Binary/Message functions
- **Multiple operations:** WireJSON intermediate format
- **Bulk processing:** WireJSON for transformations
- **Human output:** ProtoJSON only when needed

## Function Categories

### Core Functions (No Schema Required)
- `pb_message_*` functions - Work with Binary/Message format
- `pb_wire_json_*` functions - Work with WireJSON format
- `pb_message_to_wire_json()` - Binary → WireJSON conversion
- `pb_wire_json_to_message()` - WireJSON → Binary conversion

### Schema-Aware Functions (Schema Required)
- `pb_message_to_json()` - Binary → ProtoJSON conversion
- `pb_wire_json_to_json()` - WireJSON → ProtoJSON conversion
- `pb_json_to_message()` - ProtoJSON → Binary conversion
- `pb_json_to_wire_json()` - ProtoJSON → WireJSON conversion
- `_pb_message_to_number_json()` - Binary → ProtoNumberJSON conversion

### Schema Management Functions
- `pb_build_descriptor_set_json()` - Create schema from FileDescriptorSet

## Related Documentation

- **[Function Reference](function-reference.md)** - Complete API documentation
- **[Tutorial: JSON Integration](tutorial-json.md)** - JSON conversion examples
- **[API Guide](api-guide.md)** - Advanced technical details
- **[User Guide](user-guide.md)** - Problem-focused examples
- **[Schema Loading](schema-loading.md)** - Schema management guide

## Summary

The MySQL Protocol Buffers Functions library provides four distinct format representations, each optimized for specific use cases:

- **Binary/Message** for storage and single operations
- **WireJSON** for efficient multi-operation processing  
- **ProtoJSON** for human-readable debugging output
- **ProtoNumberJSON** for schema-evolution-safe programmatic access

Choose the appropriate format based on your performance requirements, schema evolution needs, and human-readability requirements. Use the conversion matrix to move between formats as needed for your specific use case.