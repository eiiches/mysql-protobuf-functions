# MySQL Protobuf Generated Functions Documentation

This documentation describes the MySQL stored functions and procedures generated for protobuf message types. All functions use configurable prefixing (`${package}_${type}_`) and follow consistent naming patterns.

## Function Naming Convention

Functions are named using the pattern: `{prefix}_{operation}_{field_name}[__modifier](parameters)`

- **prefix**: Configurable via prefix_map (e.g., `test`, `user`, `order`)
- **operation**: The type of operation (get, set, has, clear, etc.)
- **field_name**: The protobuf field name in snake_case
- **modifier**: Optional modifier for specialized operations (e.g., `__as_name`, `__or`)

## Quick Reference Tables

### Messages

| Operation          | Function Pattern | Parameters | Returns    | Notes                                         |
|--------------------|------------------|------------|------------|-----------------------------------------------|
| **New message**    | `{prefix}_new()` | None       | `JSON`     | Creates new empty message instance            |
| **To ProtoJSON**   | `{prefix}_to_json(proto_data)` | `JSON`     | `JSON`     | Converts protobuf message to ProtoJSON string |
| **From ProtoJSON** | `{prefix}_from_json(json_string)` | `JSON`     | `JSON`     | Parses ProtoJSON string into protobuf message |
| **To message**     | `{prefix}_to_message(proto_data)` | `JSON`     | `LONGBLOB` | Converts to protobuf binary format            |
| **From message**   | `{prefix}_from_message(binary_data)` | `LONGBLOB` | `JSON`     | Parses protobuf binary format |

### Enums

| Operation | Function Pattern                  | Parameters | Returns | Notes |
|-----------|-----------------------------------|------------|---------|-------|
| **To string** | `{prefix}_to_string(enum_value)`  | `INT` | `LONGTEXT` | Converts enum number to string name |
| **From string** | `{prefix}_from_string(enum_name)` | `LONGTEXT` | `INT` | Converts enum string name to number |

### Singular Fields

| Operation            | Function Pattern | Parameters | Returns | Notes                                                                                     |
|----------------------|------------------|------------|---------|-------------------------------------------------------------------------------------------|
| **Get**              | `{prefix}_get_{field}(proto_data)` | `JSON` | `{type}` | Returns default if unset (0, false, "", empty bytes, first enum value, `{}` for messages) |
| **Get with default** | `{prefix}_get_{field}__or(proto_data, default)` | `JSON, {type}` | `{type}` | Presence-tracked fields only. returns field value or custom default                       |
| **Has**              | `{prefix}_has_{field}(proto_data)` | `JSON` | `BOOLEAN` | Presence-tracked fields only: check if field is explicitly set                            |
| **Set**              | `{prefix}_set_{field}(proto_data, value)` | `JSON, {type}` | `JSON` | Sets field value and returns updated JSON                                                 |
| **Clear**            | `{prefix}_clear_{field}(proto_data)` | `JSON` | `JSON` | Reset field to default value                                                              |

### Repeated Fields

| Operation | Function Pattern | Parameters | Returns | Notes |
|-----------|------------------|------------|---------|-------|
| **Get all** | `{prefix}_get_all_{field}(proto_data)` | `JSON` | `JSON` | Returns array |
| **Set all** | `{prefix}_set_all_{field}(proto_data, array)` | `JSON, JSON` | `JSON` | Replace entire array |
| **Count** | `{prefix}_count_{field}(proto_data)` | `JSON` | `INT` | Number of elements |
| **Get at index** | `{prefix}_get_{field}(proto_data, index)` | `JSON, INT` | `{type}` | Element at index |
| **Set at index** | `{prefix}_set_{field}(proto_data, index, value)` | `JSON, INT, {type}` | `JSON` | Set element at index |
| **Insert** | `{prefix}_insert_{field}(proto_data, index, value)` | `JSON, INT, {type}` | `JSON` | Insert at index |
| **Remove** | `{prefix}_remove_{field}(proto_data, index)` | `JSON, INT` | `JSON` | Remove at index |
| **Add** | `{prefix}_add_{field}(proto_data, value)` | `JSON, {type}` | `JSON` | Append element |
| **Add all** | `{prefix}_add_all_{field}(proto_data, array)` | `JSON, JSON` | `JSON` | Append multiple |
| **Clear** | `{prefix}_clear_{field}(proto_data)` | `JSON` | `JSON` | Remove all elements |

### Map Fields

| Operation | Function Pattern | Parameters | Returns | Notes |
|-----------|------------------|------------|---------|-------|
| **Get all** | `{prefix}_get_all_{field}(proto_data)` | `JSON` | `JSON` | Returns map object |
| **Set all** | `{prefix}_set_all_{field}(proto_data, map)` | `JSON, JSON` | `JSON` | Replace entire map |
| **Count** | `{prefix}_count_{field}(proto_data)` | `JSON` | `INT` | Number of entries |
| **Contains** | `{prefix}_contains_{field}(proto_data, key)` | `JSON, {key_type}` | `BOOLEAN` | Check if key exists |
| **Get value** | `{prefix}_get_{field}(proto_data, key)` | `JSON, {key_type}` | `{value_type}` | Value by key |
| **Get with default** | `{prefix}_get_{field}__or(proto_data, key, default)` | `JSON, {key_type}, {value_type}` | `{value_type}` | Custom default |
| **Put** | `{prefix}_put_{field}(proto_data, key, value)` | `JSON, {key_type}, {value_type}` | `JSON` | Add/update entry |
| **Put all** | `{prefix}_put_all_{field}(proto_data, map)` | `JSON, JSON` | `JSON` | Merge entries |
| **Remove** | `{prefix}_remove_{field}(proto_data, key)` | `JSON, {key_type}` | `JSON` | Remove entry |
| **Clear** | `{prefix}_clear_{field}(proto_data)` | `JSON` | `JSON` | Remove all entries |

### Oneofs

| Operation | Function Pattern | Parameters | Returns | Notes |
|-----------|------------------|------------|---------|-------|
| **Which** | `{prefix}_which_{oneof}(proto_data)` | `JSON` | `INT` | Field number of set field |
| **Clear** | `{prefix}_clear_{oneof}(proto_data)` | `JSON` | `JSON` | Clear entire oneof |

## Additional Functions

### For Enum Fields

| Operation | Function Pattern | Parameters | Returns | Notes |
|-----------|------------------|------------|---------|-------|
| **Get as name** | `{prefix}_get_{field}__as_name(proto_data)` | `JSON` | `LONGTEXT` | Enum fields only: returns enum value as string name |
| **Get as name with default** | `{prefix}_get_{field}__as_name_or(proto_data, default_name)` | `JSON, LONGTEXT` | `LONGTEXT` | Optional enum fields only: returns enum name or custom default name |
| **Set from name** | `{prefix}_set_{field}__from_name(proto_data, name)` | `JSON, LONGTEXT` | `JSON` | Enum fields only: sets enum value from string name |
| **Get all as names** | `{prefix}_get_all_{field}__as_names(proto_data)` | `JSON` | `JSON` | ❌ **NOT IMPLEMENTED** - Repeated enum fields: would return array of enum names |
| **Set all from names** | `{prefix}_set_all_{field}__from_names(proto_data, names)` | `JSON, JSON` | `JSON` | ❌ **NOT IMPLEMENTED** - Repeated enum fields: would set from array of string names |