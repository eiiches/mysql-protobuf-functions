[Experimental &amp; Incomplete] Pure-MySQL Protocol Buffer Functions
====

Disclaimer: This is just a toy project - very incomplete, likely abandoned, and not intended for production use.

Usage
----

1. Load functions to MySQL.

   ```console
   $ mysql < protobuf.sql
   ```

2. Use functions in statements.

   ```sql
   -- Get the value of a string field.
   mysql> SELECT pb_message_get_string_field(_binary X'100a2a03616263', 5 /* field_number */, NULL /* repeated_index */);
   +-----------------------------------------------------------------+
   | pb_message_get_string_field(_binary X'100a2a03616263', 5, NULL) |
   +-----------------------------------------------------------------+
   | abc                                                             |
   +-----------------------------------------------------------------+
   1 row in set (0.00 sec)

   -- Get the first element of a repeated int32 (packed) field.
   mysql> SELECT pb_message_get_int32_field(_binary X'3a03010203', 7 /* field_number */, 0 /* repeated_index */);
   +---------------------------------------------------------+
   | pb_message_get_int32_field(_binary X'3a03010203', 7, 0) |
   +---------------------------------------------------------+
   |                                                       1 |
   +---------------------------------------------------------+
   1 row in set (0.01 sec)

   -- Get int32 field from a nested message field.
   mysql> SELECT pb_message_get_int32_field(pb_message_get_message_field(_binary X'4202080a', 8, NULL), 1, NULL);
   +-------------------------------------------------------------------------------------------------+
   | pb_message_get_int32_field(pb_message_get_message_field(_binary X'4202080a', 8, NULL), 1, NULL) |
   +-------------------------------------------------------------------------------------------------+
   |                                                                                              10 |
   +-------------------------------------------------------------------------------------------------+
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

* **buf** — A `BLOB` containing the serialized Protobuf message.
* **field_number** — An `INT` specifying the field number, as defined in the Protobuf schema.
* **repeated_index** — If the target field is repeated, this is the zero-based index to retrieve. For non-repeated fields, this should be `NULL`.

##### Returns

The value of the requested field, interpreted as the corresponding SQL type.

- For non-repeated fields (`repeated_index` is `NULL`), the function returns the [default value](https://protobuf.dev/programming-guides/proto3/#default) defined by the Protobuf specification if the field is not explicitly set.
- For repeated fields, if `repeated_index` exceeds the number of available elements, the function raises an "index out of range" error.

##### Notes

- The field number must match the one used in the `.proto` schema definition.
- This function does not perform schema validation; it assumes the caller knows the correct field number and expected type.

##### Type Variants

* `pb_message_get_bool_field`(buf BLOB, field_number INT, repeated_index INT) → BOOLEAN
* `pb_message_get_enum_field`(buf BLOB, field_number INT, repeated_index INT) → INT
* `pb_message_get_int32_field`(buf BLOB, field_number INT, repeated_index INT) → INT
* `pb_message_get_uint32_field`(buf BLOB, field_number INT, repeated_index INT) → INT UNSIGNED
* `pb_message_get_sint32_field`(buf BLOB, field_number INT, repeated_index INT) → INT
* `pb_message_get_fixed32_field`(buf BLOB, field_number INT, repeated_index INT) → INT UNSIGNED
* `pb_message_get_sfixed32_field`(buf BLOB, field_number INT, repeated_index INT) → INT
* `pb_message_get_float_field`(buf BLOB, field_number INT, repeated_index INT) → FLOAT
* `pb_message_get_int64_field`(buf BLOB, field_number INT, repeated_index INT) → BIGINT
* `pb_message_get_uint64_field`(buf BLOB, field_number INT, repeated_index INT) → BIGINT UNSIGNED
* `pb_message_get_sint64_field`(buf BLOB, field_number INT, repeated_index INT) → BIGINT
* `pb_message_get_fixed64_field`(buf BLOB, field_number INT, repeated_index INT) → BIGINT UNSIGNED
* `pb_message_get_sfixed64_field`(buf BLOB, field_number INT, repeated_index INT) → BIGINT
* `pb_message_get_double_field`(buf BLOB, field_number INT, repeated_index INT) → DOUBLE
* `pb_message_get_string_field`(buf BLOB, field_number INT, repeated_index INT) → TEXT
* `pb_message_get_bytes_field`(buf BLOB, field_number INT, repeated_index INT) → BLOB
* `pb_message_get_message_field`(buf BLOB, field_number INT, repeated_index INT) → BLOB

### Hazzers — pb\_message\_has\_{type}\_field()

Checks whether a specific field is present in a Protobuf-encoded BLOB.
This function is used to determine whether a field with [Field Presence](https://protobuf.dev/programming-guides/field_presence/) tracking is set in the encoded message.

#### Parameters

- **`buf`** — A `BLOB` containing the serialized Protobuf message.
- **`field_number`** — An `INT` specifying the field number, as defined in the Protobuf schema.

#### Returns

A `BOOLEAN` indicating whether the specified field is present in the encoded message.

#### Notes

- In `proto3`, presence tracking for scalar fields is only available when the field is declared with the `optional` keyword.
- For `repeated` fields, use `pb_message_get_{type}_field_count()` instead. Hazzer do not support packed repeated scalars and cannot be used to check their presence.

##### Type Variants

* `pb_message_has_bool_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_enum_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_int32_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_uint32_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_sint32_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_fixed32_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_sfixed32_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_float_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_int64_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_uint64_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_sint64_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_fixed64_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_sfixed64_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_double_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_string_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_bytes_field`(buf BLOB, field_number INT) → BOOLEAN
* `pb_message_has_message_field`(buf BLOB, field_number INT) → BOOLEAN

### Repeated Field Counts — pb\_message\_get\_{type}\_field_count()

Returns the number of elements present in a repeated field of a Protobuf-encoded BLOB.
This function is used to determine the size of repeated fields, including packed repeated fields.

#### Parameters

- **`buf`** — A `BLOB` containing the serialized Protobuf message.
- **`field_number`** — An `INT` specifying the field number, as defined in the Protobuf schema.

#### Returns

An `INT` representing the number of elements in the specified repeated field.

#### Notes

- Use in combination with `pb_message_get_{type}_field()` to access individual repeated values by index.

##### Type Variants

* `pb_message_get_bool_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_enum_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_int32_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_uint32_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_sint32_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_fixed32_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_sfixed32_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_float_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_int64_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_uint64_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_sint64_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_fixed64_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_sfixed64_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_double_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_string_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_bytes_field_count`(buf BLOB, field_number INT) → INT
* `pb_message_get_message_field_count`(buf BLOB, field_number INT) → INT

Procedure Reference
-------------------

### Load File Descriptor Set — `CALL pb_load_file_descriptor_set()`

Loads a compiled Protobuf `FileDescriptorSet` into internal database tables.
This procedure registers the schema information required for operations that depend on message descriptors, such as `pb_message_to_json()` (planned but not yet implemented).

#### Parameters

- IN **`set_name`** `VARCHAR(64)` — A user-defined identifier that is used to distinguish file descriptor sets.
  This allows multiple descriptor sets to coexist without conflict.

- IN **`file_descriptor_set`** `BLOB` — A binary-encoded `FileDescriptorSet`, typically generated using the `protoc --descriptor_set_out` option.
  The input must conform to the `google.protobuf.FileDescriptorSet` message format.

#### Notes

- This procedure is only intended for use cases where Protobuf schema information is required.
  - Getters and hazzers (`pb_message_get_{type}_field()`, `pb_message_get_{type}_field_count()`, `pb_message_has_{type}_field()`) don't require this procedure.
- If a descriptor set with the same `set_name` already exists, the procedure will fail.

### Delete File Descriptor Set — `CALL pb_delete_file_descriptor_set()`

Deletes the internal tables associated with a previously loaded descriptor set.
This procedure is used to clean up resources created by `pb_load_file_descriptor_set()`.

#### Parameters

- IN **`set_name`** `VARCHAR(64)` — The identifier used when the descriptor set was loaded.
  This determines which set of internal tables will be dropped.

#### Notes

- If the specified descriptor set does not exist, the procedure performs no action.
- Temporary descriptor sets (created with `persist = FALSE`) do not require manual deletion — they are automatically dropped at the end of the session.

TODO
----

* Add `pb_message_set_{type}_field(data, field_number, repeated_index)`.
* Schema support.
  * `pb_message_to_json()`
* Implement map support.

Limitation
----------

* I never benchmarked the functions. I have no idea how slow these functions are.
  * MySQL doesn't have an ARRAY type. With the current API, each element of a repeated field must be retrieved one by one. Yes, O(n^2) to retrieve all.

* Currently, MySQL doesn't allow using stored functions in functional indexes or generated columns. Use triggers or views instead.

   ```sql
   CREATE TABLE Example (protobuf_data BLOB, generated_field INT, INDEX (generated_field));
   CREATE TRIGGER Example_set_generated_field_on_update
       BEFORE UPDATE ON Example
       FOR EACH ROW SET NEW.generated_field = pb_message_get_int32_field(NEW.protobuf_data, 1, NULL);
   CREATE TRIGGER Example_set_generated_field_on_insert
       BEFORE INSERT ON Example
       FOR EACH ROW SET NEW.generated_field = pb_message_get_int32_field(NEW.protobuf_data, 1, NULL);
   ```

   ```sql
   CREATE TABLE Example (protobuf_data BLOB);
   CREATE VIEW ExampleView AS SELECT pb_message_get_int32_field(protobuf_data, 1, NULL) FROM Example;
   ```

* [Groups](https://protobuf.dev/programming-guides/encoding/#groups) are not, and probably will not be supported.
