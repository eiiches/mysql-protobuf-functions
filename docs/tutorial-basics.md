# Tutorial: Basics

*Getting started with low-level protobuf field access*

This tutorial demonstrates the core functionality using a sample Protocol Buffers schema. Follow along to learn how to access protobuf data in MySQL using **low-level functions**.

> **Want to convert protobuf to JSON instead?** If you're looking to view protobuf data in human-readable JSON format, jump to [Tutorial: JSON Integration](tutorial-json.md) instead. That tutorial covers schema-aware functions that provide rich output with field names and formatted values.

> **What are low-level functions?** These functions work directly with protobuf field numbers and don't require schema information. They're perfect for basic field access and are the foundation of the library.

> **Prerequisites:** This tutorial requires the core protobuf functions to be installed. If you haven't already, see the [Installation Guide](installation.md) to get started.

## Setting Up Sample Data

**person.proto**
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

```sql
-- $ protoc --encode=Person person.proto <<-EOF | xxd -p -c0
-- name: "Agent Smith"
-- id: 1
-- email: "smith@example.com"
-- phones: [{type: PHONE_TYPE_WORK, number: "+81-00-0000-0000"}]
-- last_updated: {seconds: 1748781296, nanos: 789000000}
-- EOF
> CREATE TABLE Example (pb_data BLOB);
> INSERT INTO Example (pb_data) VALUES (_binary X'0a0b4167656e7420536d69746810011a11736d697468406578616d706c652e636f6d22140a102b38312d30302d303030302d3030303010032a0c08f091f1c10610c0de9cf802');
```

## Reading Protobuf Fields

The key concept with low-level functions is that you work directly with **field numbers** from your `.proto` definition. Looking at our `Person` message:

- `name = 1` ← field number 1
- `id = 2` ← field number 2
- `email = 3` ← field number 3
- `phones = 4` ← field number 4

Use the field accessor functions to extract data from protobuf messages:

```sql
-- Get id; The field number for id is 2.
> SELECT pb_message_get_int32_field(pb_data, 2 /* field number */, 0 /* default value */) FROM Example;
1

-- Get email; email is 3
> SELECT pb_message_get_string_field(pb_data, 3, '' /* default value */) FROM Example;
smith@example.com

-- Get phones[0].number; phones is 4, number is 1
> SELECT pb_message_get_string_field(pb_message_get_repeated_message_field_element(pb_data, 4, 0), 1, '' /* default value */) FROM Example;
+81-00-0000-0000

-- Get phones[0].type; phones is 4, type is 2
> SELECT pb_message_get_enum_field(pb_message_get_repeated_message_field_element(pb_data, 4, 0), 2, 0 /* default value */) FROM Example;
3
```

> **Why field numbers instead of field names?** Field numbers are stable - they never change even when field names are renamed in your `.proto` file. This ensures your MySQL queries continue working as your protobuf schema evolves.

## Working with Repeated Fields

Repeated fields require special handling:

```sql
-- Get count of phone numbers
> SELECT pb_message_get_repeated_message_field_count(pb_data, 4) FROM Example;
1

-- Get first phone number (index 0)
> SELECT pb_message_get_repeated_message_field_element(pb_data, 4, 0) FROM Example;
[binary data]

-- Extract data from the first phone number
> SELECT
    pb_message_get_string_field(pb_message_get_repeated_message_field_element(pb_data, 4, 0), 1, '') AS number,
    pb_message_get_enum_field(pb_message_get_repeated_message_field_element(pb_data, 4, 0), 2, 0) AS type
  FROM Example;
+81-00-0000-0000    3
```

## Checking Field Presence

You can check if a field is set (for fields with field presence enabled):

```sql
-- Check if name field is set
> SELECT pb_message_has_string_field(pb_data, 1) FROM Example;
1

-- Check if optional field is set
> SELECT pb_message_has_string_field(pb_data, 999) FROM Example;
0
```

## Understanding the Function Pattern

Low-level functions follow consistent patterns for different operations:

```sql
-- Reading Fields
pb_message_get_{TYPE}_field(message, field_number, default_value)
pb_message_has_{TYPE}_field(message, field_number)

-- Repeated Fields
pb_message_get_repeated_{TYPE}_field_count(message, field_number)
pb_message_get_repeated_{TYPE}_field_element(message, field_number, index)
pb_message_get_repeated_{TYPE}_field_as_json_array(message, field_number)                -- get entire repeated field as JSON array

-- Setting Fields (Modification) - see @docs/tutorial-modification.md for usage examples
pb_message_set_{TYPE}_field(message, field_number, value)
pb_message_clear_{TYPE}_field(message, field_number)

-- Repeated Field Modification - see @docs/tutorial-modification.md for usage examples
pb_message_add_repeated_{TYPE}_field_element(message, field_number, value)               -- add single element
pb_message_add_all_repeated_{TYPE}_field_elements(message, field_number, value_array)    -- add multiple elements from JSON array (bulk)
pb_message_insert_repeated_{TYPE}_field_element(message, field_number, index, value)     -- insert at specific index
pb_message_set_repeated_{TYPE}_field_element(message, field_number, index, value)        -- update element at specific index
pb_message_set_repeated_{TYPE}_field(message, field_number, value_array)                 -- replace entire repeated field with JSON array
pb_message_remove_repeated_{TYPE}_field_element(message, field_number, index)            -- remove element at specific index
pb_message_clear_repeated_{TYPE}_field(message, field_number)                            -- remove all elements
```

**Available types:**
- `int32`, `int64`, `uint32`, `uint64`, `sint32`, `sint64` - integer types
- `fixed32`, `fixed64`, `sfixed32`, `sfixed64` - fixed-width integers
- `float`, `double` - floating point numbers
- `bool` - boolean values
- `string` - text strings
- `bytes` - binary data
- `enum` - enum values (returns integer)
- `message` - nested messages

## Next Steps

Now that you understand basic field access, you can:

- **Learn message modification** → [Tutorial: Modification](tutorial-modification.md)
- **Explore schema-aware features** → [Tutorial: JSON Integration](tutorial-json.md)
- **See practical examples** → [User Guide](user-guide.md)

The low-level functions you've learned here are the foundation - they're used throughout the library and are perfect when you just need to extract or check specific field values.
