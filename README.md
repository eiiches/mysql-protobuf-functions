Pure-MySQL Protocol Buffer Functions
====

Warning: This project is in an early stage and may undergo significant changes that are not backward compatible.

Requirements
------------

* MySQL 8.0.17 or later.
  * JSON_TABLE() is added in 8.0.4 but needs 8.0.17 for [this bugfix](https://bugs.mysql.com/bug.php?id=92976).

* Aurora MySQL 3.04.0 (oldest available 3.x I could test as of 2025/06) or later.

Install
------

Install stored functions and procedure by downloading and running the SQL files in MySQL.

```console
$ mysql < protobuf.sql
$ mysql < protobuf-accessors.sql
```

NOTE: All functions and procedures are prefixed with `_pb_` or `pb_` to avoid name conflicts.
If you already use these names, be careful not to overwrite existing routines.

(Optional) In addition, if you need `pb_message_to_json()` (a JSON conversion function), load the following SQL files.

```console
$ mysql < protobuf-descriptor.sql # (optional) for descriptor set support
$ mysql < protobuf-json.sql # (optional) for json conversion support (incomplete)
```

NOTE: `protobuf-descriptor.sql` automatically creates several tables prefixed with `_Proto_` to store schema information.

Example
-------

```sql
-- Get the value of a string field.
mysql> SELECT pb_message_get_string_field(_binary X'100a2a03616263', 5 /* field_number */, '' /* default_value */);
+---------------------------------------------------------------+
| pb_message_get_string_field(_binary X'100a2a03616263', 5, '') |
+---------------------------------------------------------------+
| abc                                                           |
+---------------------------------------------------------------+
```

```sql
-- Get the first element of a repeated int32 (packed) field.
mysql> SELECT pb_message_get_repeated_int32_field(_binary X'3a03010203', 7 /* field_number */, 1 /* repeated_index */);
+---------------------------------------------------------+
| pb_message_get_int32_field(_binary X'3a03010203', 7, 0) |
+---------------------------------------------------------+
|                                                       2 |
+---------------------------------------------------------+
```

```sql
-- Get int32 field from a nested message field.
mysql> SELECT pb_message_get_int32_field(pb_message_get_message_field(_binary X'4202080a', 8, _binary X''), 1, 0);
+-----------------------------------------------------------------------------------------------------+
| pb_message_get_int32_field(pb_message_get_message_field(_binary X'4202080a', 8, _binary X''), 1, 0) |
+-----------------------------------------------------------------------------------------------------+
|                                                                                                  10 |
+-----------------------------------------------------------------------------------------------------+
```

Protobuf schema used in this example:

```protobuf
syntax = "proto3";

message Test {
   string string_field = 5;
   repeated int32 repeated_int32_field = 7;
   TestMessage message_field = 8;
}

message TestMessage {
   int32 int32_field = 1;
}
```

Function Reference
---------

### Getters — pb\_message\_get\_{type}\_field()

Retrieves the value of a specified field from a Protobuf-encoded BLOB. This function is used to extract individual values from serialized Protocol Buffers messages stored as binary data.

##### Parameters

* **message** — A `LONGBLOB` containing the serialized Protobuf message.
* **field_number** — An `INT` specifying the field number, as defined in the Protobuf schema.
* **default_value** — A default value to return when the field is not present. For proto3 without explicit field presence, this should be set to the [default value](https://protobuf.dev/programming-guides/proto3/#default) defined by the Protobuf specification. For proto2 messages, this should also be set to the [default value](https://protobuf.dev/programming-guides/proto2/#default) unless the field has an explicit `default` option that overrides the default value.

##### Returns

The value of the requested field, interpreted as the corresponding SQL type.

##### Notes

- The field number must match the one used in the `.proto` schema definition.
- This function does not perform schema validation; it assumes the caller knows the correct field number and expected type.
- MySQL does not support `+inf`, `-inf`, or `NaN`. Therefore, `float` and `double` variants return `NULL` instead if the corresponding field contains any of these values.
- This function parses the message each time it is called. For better performance when accessing multiple fields, use `pb_message_to_wire_json()` to parse the message once, and then call `pb_wire_json_get_{type}_field()` for each field.

##### Type Variants

* `pb_message_get_bool_field`(message LONGBLOB, field_number INT, default_value BOOLEAN) → BOOLEAN
* `pb_message_get_enum_field`(message LONGBLOB, field_number INT, default_value INT) → INT
* `pb_message_get_int32_field`(message LONGBLOB, field_number INT, default_value INT) → INT
* `pb_message_get_uint32_field`(message LONGBLOB, field_number INT, default_value INT) → INT UNSIGNED
* `pb_message_get_sint32_field`(message LONGBLOB, field_number INT, default_value INT) → INT
* `pb_message_get_fixed32_field`(message LONGBLOB, field_number INT, default_value INT UNSIGNED) → INT UNSIGNED
* `pb_message_get_sfixed32_field`(message LONGBLOB, field_number INT, default_value INT) → INT
* `pb_message_get_float_field`(message LONGBLOB, field_number INT, default_value FLOAT) → FLOAT
* `pb_message_get_int64_field`(message LONGBLOB, field_number INT, default_value BIGINT) → BIGINT
* `pb_message_get_uint64_field`(message LONGBLOB, field_number INT, default_value BIGINT UNSIGNED) → BIGINT UNSIGNED
* `pb_message_get_sint64_field`(message LONGBLOB, field_number INT, default_value BIGINT) → BIGINT
* `pb_message_get_fixed64_field`(message LONGBLOB, field_number INT, default_value BIGINT UNSIGNED) → BIGINT UNSIGNED
* `pb_message_get_sfixed64_field`(message LONGBLOB, field_number INT, default_value BIGINT) → BIGINT
* `pb_message_get_double_field`(message LONGBLOB, field_number INT, default_value DOUBLE) → DOUBLE
* `pb_message_get_string_field`(message LONGBLOB, field_number INT, default_value TEXT) → TEXT
* `pb_message_get_bytes_field`(message LONGBLOB, field_number INT, default_value LONGBLOB) → LONGBLOB
* `pb_message_get_message_field`(message LONGBLOB, field_number INT, default_value LONGBLOB) → LONGBLOB

### Getters — pb\_message\_get\_repeated\_{type}\_field()

Retrieves the value of a specified field from a Protobuf-encoded BLOB. This function is used to extract individual values from serialized Protocol Buffers messages stored as binary data.

##### Parameters

* **message** — A `LONGBLOB` containing the serialized Protobuf message.
* **field_number** — An `INT` specifying the field number, as defined in the Protobuf schema.
* **repeated_index** — The zero-based index to retrieve.

##### Returns

The value of the requested field, interpreted as the corresponding SQL type.

- If `repeated_index` exceeds the number of available elements, the function raises an "index out of range" error.

##### Notes

- The field number must match the one used in the `.proto` schema definition.
- This function does not perform schema validation; it assumes the caller knows the correct field number and expected type.
- MySQL does not support `+inf`, `-inf`, or `NaN`. Therefore, `float` and `double` variants return `NULL` instead if the corresponding field contains any of these values.

##### Type Variants

* `pb_message_get_repeated_bool_field`(message LONGBLOB, field_number INT, repeated_index INT) → BOOLEAN
* `pb_message_get_repeated_enum_field`(message LONGBLOB, field_number INT, repeated_index INT) → INT
* `pb_message_get_repeated_int32_field`(message LONGBLOB, field_number INT, repeated_index INT) → INT
* `pb_message_get_repeated_uint32_field`(message LONGBLOB, field_number INT, repeated_index INT) → INT UNSIGNED
* `pb_message_get_repeated_sint32_field`(message LONGBLOB, field_number INT, repeated_index INT) → INT
* `pb_message_get_repeated_fixed32_field`(message LONGBLOB, field_number INT, repeated_index INT) → INT UNSIGNED
* `pb_message_get_repeated_sfixed32_field`(message LONGBLOB, field_number INT, repeated_index INT) → INT
* `pb_message_get_repeated_float_field`(message LONGBLOB, field_number INT, repeated_index INT) → FLOAT
* `pb_message_get_repeated_int64_field`(message LONGBLOB, field_number INT, repeated_index INT) → BIGINT
* `pb_message_get_repeated_uint64_field`(message LONGBLOB, field_number INT, repeated_index INT) → BIGINT UNSIGNED
* `pb_message_get_repeated_sint64_field`(message LONGBLOB, field_number INT, repeated_index INT) → BIGINT
* `pb_message_get_repeated_fixed64_field`(message LONGBLOB, field_number INT, repeated_index INT) → BIGINT UNSIGNED
* `pb_message_get_repeated_sfixed64_field`(message LONGBLOB, field_number INT, repeated_index INT) → BIGINT
* `pb_message_get_repeated_double_field`(message LONGBLOB, field_number INT, repeated_index INT) → DOUBLE
* `pb_message_get_repeated_string_field`(message LONGBLOB, field_number INT, repeated_index INT) → TEXT
* `pb_message_get_repeated_bytes_field`(message LONGBLOB, field_number INT, repeated_index INT) → LONGBLOB
* `pb_message_get_repeated_message_field`(message LONGBLOB, field_number INT, repeated_index INT) → LONGBLOB

### Getters — pb\_message\_get\_repeated\_{type}\_field_as_json_array()

Retrieves the repeated values of a specified field as JSON array from a Protobuf-encoded BLOB.

##### Parameters

* **message** — A `LONGBLOB` containing the serialized Protobuf message.
* **field_number** — An `INT` specifying the field number, as defined in the Protobuf schema.

##### Returns

An `JSON` array containing all elements of the specified field.

##### Notes

- The field number must match the one used in the `.proto` schema definition.
- This function does not perform schema validation; it assumes the caller knows the correct field number and expected type.
- MySQL does not support `+inf`, `-inf`, or `NaN`. Therefore, `float` and `double` variants return `NULL` instead if the corresponding field contains any of these values.

##### Type Variants

* `pb_message_get_repeated_bool_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_enum_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_int32_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_uint32_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_sint32_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_fixed32_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_sfixed32_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_float_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_int64_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_uint64_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_sint64_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_fixed64_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_sfixed64_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_double_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_string_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_bytes_field_as_json_array`(message LONGBLOB, field_number INT) → JSON
* `pb_message_get_repeated_message_field_as_json_array`(message LONGBLOB, field_number INT) → JSON

### Hazzers — pb\_message\_has\_{type}\_field()

Checks whether a specific field is present in a Protobuf-encoded BLOB.
This function is used to determine whether a field with [Field Presence](https://protobuf.dev/programming-guides/field_presence/) tracking is set in the encoded message.

#### Parameters

- **message** — A `LONGBLOB` containing the serialized Protobuf message.
- **field_number** — An `INT` specifying the field number, as defined in the Protobuf schema.

#### Returns

A `BOOLEAN` indicating whether the specified field is present in the encoded message.

#### Notes

- In `proto3`, presence tracking for scalar fields is only available when the field is declared with the `optional` keyword.
- Using this function on `repeated` fields is an error. Hazzers do not support packed repeated scalars and cannot be used to check their presence. Use `pb_message_get_repeated_{type}_field_count()` instead. 

##### Type Variants

* `pb_message_has_bool_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_enum_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_int32_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_uint32_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_sint32_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_fixed32_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_sfixed32_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_float_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_int64_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_uint64_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_sint64_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_fixed64_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_sfixed64_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_double_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_string_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_bytes_field`(message LONGBLOB, field_number INT) → BOOLEAN
* `pb_message_has_message_field`(message LONGBLOB, field_number INT) → BOOLEAN

### Repeated Field Counts — pb\_message\_get\_repeated\_{type}\_field_count()

Returns the number of elements present in a repeated field of a Protobuf-encoded BLOB.
To retreive the value of each element, use `pb_message_get_repeated_{type}_field()`.

#### Parameters

- **`message`** — A `LONGBLOB` containing the serialized Protobuf message.
- **`field_number`** — An `INT` specifying the field number, as defined in the Protobuf schema.

#### Returns

An `INT` representing the number of elements in the specified repeated field.

##### Type Variants

* `pb_message_get_repeated_bool_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_enum_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_int32_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_uint32_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_sint32_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_fixed32_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_sfixed32_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_float_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_int64_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_uint64_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_sint64_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_fixed64_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_sfixed64_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_double_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_string_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_bytes_field_count`(message LONGBLOB, field_number INT) → INT
* `pb_message_get_repeated_message_field_count`(message LONGBLOB, field_number INT) → INT

### pb\_message\_to\_wire\_json()

Decodes a protobuf-encoded BLOB into a JSON representation that are suitable for low-level inspection and debugging. The JSON output preserves the order and structure of the wire format.
This function does not interpret the values beyond the wire-level format.

##### Parameters

- **message** (LONGBLOB) — Raw protobuf-encoded data.

##### Example

```sql
SELECT pb_message_to_wire_json(_binary X'0A0B696E7433325F6669656C64180120012805520A696E7433324669656C64');
```

```json
{
  "1": [{"i": 0, "n": 1, "t": 2, "v": "aW50MzJfZmllbGQ="}],
  "3": [{"i": 1, "n": 3, "t": 0, "v": 1}],
  "4": [{"i": 2, "n": 4, "t": 0, "v": 1}],
  "5": [{"i": 3, "n": 5, "t": 0, "v": 5}],
  "10": [{"i": 4, "n": 10, "t": 2, "v": "aW50MzJGaWVsZA=="}]
}
```


[Experimental] Descriptor Set Functions &amp; Procedures
-------------------

### Load File Descriptor Set — CALL pb\_descriptor\_set\_load()

Loads a compiled Protobuf `FileDescriptorSet` into internal database tables.
This procedure registers the schema information required for operations that depend on message descriptors, such as `pb_message_to_json()`.

#### Parameters

- IN **set_name** (VARCHAR(64)) — A user-defined identifier that is used to distinguish file descriptor sets.
  This allows multiple descriptor sets to coexist without conflict.

- IN **file_descriptor_set** (LONGBLOB) — A binary-encoded `FileDescriptorSet`, typically generated using the `protoc --descriptor_set_out` or `buf build -o ${name}.binpb`.
  The input must conform to the `google.protobuf.FileDescriptorSet` message format.

#### Notes

- This procedure is only intended for use cases where Protobuf schema information is required.
  - Getters and hazzers (`pb_message_get_{type}_field()`, `pb_message_get_{type}_field_count()`, `pb_message_has_{type}_field()`) don't require this procedure.
- If a descriptor set with the same `set_name` already exists, the procedure will fail.

### Delete File Descriptor Set — CALL pb\_descriptor\_set\_delete()

Deletes a previously loaded descriptor set from internal tables.
This procedure is used to clean up resources created by `pb_descriptor_set_load()`.

#### Parameters

- IN **`set_name`** `VARCHAR(64)` — The identifier used when the descriptor set was loaded.

#### Notes

- If the specified descriptor set does not exist, the procedure performs no action.

### [Function] pb\_descriptor\_set\_exists()

TBD

### [Function] pb\_descriptor\_set\_contains\_message\_type()

TBD

### [Function] pb\_descriptor\_set\_contains\_enum\_type()

TBD

[Incomplete] JSON Conversion Functions &amp; Procedures
-----

### [Function] pb\_message\_to\_json()

Converts a Protobuf-encoded BLOB into a JSON object using the specified type and schema set. This function deserializes a Protocol Buffers message into a human-readable and structured JSON format.

This function is primarily intended for debugging or inspection. It should not be used in production code. JSON relies on field names rather than field numbers, which compromises a key benefit of Protocol Buffers: the ability to rename fields without breaking compatibility. Instead, applications should read and decode the raw Protobuf-encoded BLOB directly.

Each reader can use its own descriptor set (i.e., reader schema), allowing it to remain unaffected by schema changes and independently control when to adopt schema updates.

##### Parameters

* **set_name** — A `VARCHAR(64)` specifying the name of the schema set containing the compiled Protobuf descriptors.
* **full_type_name** — A `VARCHAR(512)` representing the fully-qualified name of the Protobuf message type (e.g., `.my.package.MessageType`). A fully-qualified name always starts with a dot.
* **message** — A `LONGBLOB` containing the serialized Protobuf message to be converted.

##### Returns

A `JSON` object that represents the Protobuf message, with field names and values corresponding to those defined in the Protobuf schema.

##### Errors

* Returns an error if the set_name or full_type_name cannot be resolved.
* Returns an error if the buf is not a valid serialized message of the given type.

TODO
----

* Add setters.
  * `pb_message_set_{type}_field(data, field_number, value)`
  * `pb_message_add_repeated_int32_field(data, field_number, value)`
  * `pb_message_set_repeated_int32_field(data, field_number, repeated_index, value)`
  * `pb_message_clear_{type}_field(data, field_number)`
  * `pb_message_clear_repeated_{type}_field(data, field_number)`
* Implement map support.

Limitation
----------

* Currently, MySQL doesn't allow using stored functions in functional indexes or generated columns. Use triggers or views instead.

   ```sql
   CREATE TABLE Example (protobuf_data BLOB, generated_field INT, INDEX (generated_field));
   CREATE TRIGGER Example_set_generated_field_on_update
       BEFORE UPDATE ON Example
       FOR EACH ROW SET NEW.generated_field = pb_message_get_int32_field(NEW.protobuf_data, 1, 0);
   CREATE TRIGGER Example_set_generated_field_on_insert
       BEFORE INSERT ON Example
       FOR EACH ROW SET NEW.generated_field = pb_message_get_int32_field(NEW.protobuf_data, 1, 0);
   ```

   ```sql
   CREATE TABLE Example (protobuf_data BLOB);
   CREATE VIEW ExampleView AS SELECT pb_message_get_int32_field(protobuf_data, 1, 0) FROM Example;
   ```

* [Groups](https://protobuf.dev/programming-guides/encoding/#groups) are not, and probably will not be supported.
