# MySQL Protocol Buffers Library - Documentation

*Query and manipulate Protocol Buffers data directly in MySQL using custom SQL functions*

## 🚀 Start Here - Choose Your Path

### I'm New to This Library
**👉 [User Guide](user-guide.md)** - Problem-focused examples and real-world solutions  
*Perfect for: Developers who want practical copy-paste solutions*

### I Need Quick Answers  
**👉 [Quick Reference](quick-reference.md)** - Cheat sheet with common operations  
*Perfect for: Experienced users who need syntax reminders*

### I Want Complete Documentation
**👉 [Function Reference](function-reference.md)** - Complete API reference  
*Perfect for: Advanced users and comprehensive API reference*

### I Need In-Depth Examples
**👉 [API Guide](api-guide.md)** - Detailed technical guide with advanced topics  
*Perfect for: Complex implementations and edge cases*

---

## What This Library Does

**Problem:** You have protobuf messages stored in MySQL but need to query/filter by field values directly in SQL instead of deserializing in your application.

**Solution:** This library provides comprehensive MySQL functions to work with protobuf data directly in SQL.

```sql
-- ✅ Find users by age range (field 2 in protobuf)
SELECT user_id FROM users 
WHERE pb_message_get_int32_field(profile_data, 2, 0) BETWEEN 25 AND 35;

-- ✅ Update user email (field 3) without touching other data  
UPDATE users 
SET profile_data = pb_message_set_string_field(profile_data, 3, 'new@email.com')
WHERE user_id = 123;

-- ✅ Add hobby to user's list (repeated field 4)
UPDATE users 
SET profile_data = pb_message_add_repeated_string_field_element(profile_data, 4, 'hiking')
WHERE user_id = 123;
```

## Core Capabilities

| Feature | Function Example | Use Case |
|---------|------------------|-----------|
| **Read Fields** | `pb_message_get_string_field(data, 1, '')` | Query by protobuf field values |
| **Write Fields** | `pb_message_set_int32_field(data, 2, 25)` | Update specific fields safely |
| **Repeated Fields** | `pb_message_add_repeated_string_field_element(data, 3, 'item')` | Manage lists/arrays |
| **Bulk Operations** | `pb_message_add_all_repeated_int32_field_elements(data, 4, '[1,2,3]')` | Efficient batch updates |
| **JSON Conversion** | `pb_message_to_json('schema', '.Type', data)` | Human-readable JSON with schema |
| **Wire Format** | `pb_message_to_wire_json(data)` | Performance optimization for multiple operations |

## When to Use This Library

**✅ Use when:**
- Protobuf messages stored as LONGBLOB in MySQL
- Need to query/filter by protobuf field values
- Want to update specific fields directly in the database
- Building analytics on protobuf data
- Migrating protobuf schemas

**❌ Don't use when:**
- You can process protobuf data in your application (usually more efficient)
- Working with simple data that doesn't need protobuf
- Performance is critical and you can restructure your data model

## Installation

```sql
-- Core protobuf functions (required)
SOURCE protobuf.sql;
SOURCE protobuf-accessors.sql;

-- JSON features (optional) 
SOURCE protobuf-descriptor.sql;
SOURCE protobuf-json.sql;
```

**Requirements:** MySQL 8.0.17+ or Aurora MySQL 3.04.0+

## Function Naming Pattern

All public functions follow this pattern:
```
pb_{input_type}_{operation}_{field_type}_field[_element[s]]
```

Examples:
- `pb_message_get_string_field()` - Get string field from binary message
- `pb_message_set_int32_field()` - Set int32 field in binary message  
- `pb_message_add_repeated_string_field_element()` - Add to repeated string field

**🚨 Important:** Only use functions starting with `pb_`. Functions starting with `_pb_` are internal.

## Supported Protobuf Types

| Protobuf | MySQL Type | Function Suffix |
|----------|------------|-----------------|
| `string` | LONGTEXT | `_string_field` |
| `int32`, `int64` | INT, BIGINT | `_int32_field`, `_int64_field` |
| `bool` | BOOLEAN | `_bool_field` |
| `float`, `double` | FLOAT, DOUBLE | `_float_field`, `_double_field` |
| `bytes` | LONGBLOB | `_bytes_field` |
| `message` | LONGBLOB | `_message_field` |
| *All protobuf types supported - see [Function Reference](function-reference.md#field-types-reference)* |

## Documentation Quick Links

- **[User Guide](user-guide.md)** - Start here for practical examples
- **[Quick Reference](quick-reference.md)** - Syntax cheat sheet  
- **[Function Reference](function-reference.md)** - Complete API documentation
- **[API Guide](api-guide.md)** - Advanced usage and examples

## Need Help?

- **Common Issues:** Check [User Guide Troubleshooting](user-guide.md#troubleshooting-guide)
- **Performance:** See [User Guide Performance](user-guide.md#performance-best-practices)  
- **All Functions:** Browse [Function Reference](function-reference.md)
- **Advanced Topics:** See [API Guide Advanced Topics](api-guide.md#advanced-topics)