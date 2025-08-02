# Tutorial

This tutorial demonstrates the core functionality using a sample Protocol Buffers schema. Follow along to learn how to work with protobuf data in MySQL.

<details>
<summary>person.proto</summary>

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

</details>

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

## Accessing Protobuf Fields

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

## Converting Protobuf to JSON

> **Prerequisites**: This section requires both `protobuf-descriptor.sql` and `protobuf-json.sql` to be installed.

### Loading Protobuf Schema

There are multiple ways to load protobuf descriptor set JSON into MySQL:

#### Method 1: Using pb_build_descriptor_set_json

First, generate your protobuf descriptor set:

```console
$ protoc --descriptor_set_out=/dev/stdout --include_imports person.proto | xxd -p -c0
0aff010a1f676f6f676c652f70726f746f6275662f74696d657374616d702e70726f746f120f676f6f676c652e70726f746f627566223b0a0954696d657374616d7012180a077365636f6e647318012001280352077365636f6e647312140a056e616e6f7318022001280552056e616e6f734285010a13636f6d2e676f6f676c652e70726f746f627566420e54696d657374616d7050726f746f50015a32676f6f676c652e676f6c616e672e6f72672f70726f746f6275662f74797065732f6b6e6f776e2f74696d657374616d707062f80101a20203475042aa021e476f6f676c652e50726f746f6275662e57656c6c4b6e6f776e5479706573620670726f746f330aa0030a0c706572736f6e2e70726f746f1a1f676f6f676c652f70726f746f6275662f74696d657374616d702e70726f746f22e6020a06506572736f6e12120a046e616d6518012001280952046e616d65120e0a0269641802200128055202696412140a05656d61696c1803200128095205656d61696c122b0a0670686f6e657318042003280b32132e506572736f6e2e50686f6e654e756d626572520670686f6e6573123d0a0c6c6173745f7570646174656418052001280b321a2e676f6f676c652e70726f746f6275662e54696d657374616d70520b6c617374557064617465641a4c0a0b50686f6e654e756d62657212160a066e756d62657218012001280952066e756d62657212250a047479706518022001280e32112e506572736f6e2e50686f6e655479706552047479706522680a0950686f6e6554797065121a0a1650484f4e455f545950455f554e535045434946494544100012150a1150484f4e455f545950455f4d4f42494c45100112130a0f50484f4e455f545950455f484f4d45100212130a0f50484f4e455f545950455f574f524b1003620670726f746f33
```

Alternatively, if you use Buf, you can use `buf build -o ${name}.binpb` to generate a binary FileDescriptorSet.

The `pb_build_descriptor_set_json` function converts the binary FileDescriptorSet to a versioned JSON array format. For details about the format structure, see the [descriptorsetjson documentation](internal/descriptorsetjson/README.md).

Then load the schema using `pb_build_descriptor_set_json` and save the result:

```sql
-- Option A: Save to user variable
> SET @test_schema = pb_build_descriptor_set_json(_binary X'0aff010a1f676f6f676c652f70726f746f6275662f74696d657374616d702e70726f746f120f676f6f676c652e70726f746f627566223b0a0954696d657374616d7012180a077365636f6e647318012001280352077365636f6e647312140a056e616e6f7318022001280552056e616e6f734285010a13636f6d2e676f6f676c652e70726f746f627566420e54696d657374616d7050726f746f50015a32676f6f676c652e676f6c616e672e6f72672f70726f746f6275662f74797065732f6b6e6f776e2f74696d657374616d707062f80101a20203475042aa021e476f6f676c652e50726f746f6275662e57656c6c4b6e6f776e5479706573620670726f746f330aa0030a0c706572736f6e2e70726f746f1a1f676f6f676c652f70726f746f6275662f74696d657374616d702e70726f746f22e6020a06506572736f6e12120a046e616d6518012001280952046e616d65120e0a0269641802200128055202696412140a05656d61696c1803200128095205656d61696c122b0a0670686f6e657318042003280b32132e506572736f6e2e50686f6e654e756d626572520670686f6e6573123d0a0c6c6173745f7570646174656418052001280b321a2e676f6f676c652e70726f746f6275662e54696d657374616d70520b6c617374557064617465641a4c0a0b50686f6e654e756d62657212160a066e756d62657218012001280952066e756d62657212250a047479706518022001280e32112e506572736f6e2e50686f6e655479706552047479706522680a0950686f6e6554797065121a0a1650484f4e455f545950455f554e535045434946494544100012150a1150484f4e455f545950455f4d4f42494c45100112130a0f50484f4e455f545950455f484f4d45100212130a0f50484f4e455f545950455f574f524b1003620670726f746f33');

-- Option B: Save to MySQL table
> CREATE TABLE schema_registry (schema_name VARCHAR(255) PRIMARY KEY, schema_json JSON);
> INSERT INTO schema_registry VALUES ('test', pb_build_descriptor_set_json(_binary X'...'));
```

#### Method 2: Using protoc-gen-descriptor_set_json

Generate a SQL function containing the descriptor set JSON:

```console
$ protoc --descriptor_set_json_out=. --descriptor_set_json_opt=name=test_schema person.proto
```

This generates a SQL file with a function that returns the schema JSON. Load it into MySQL:

```bash
mysql -u your_username -p your_database < test_schema.sql
```

```sql
> SELECT test_schema(); -- Returns the schema JSON
```

You can now use the generated function `test_schema()` directly as the first argument to `pb_message_to_json()` and other functions that require descriptor set JSON.

For detailed information about the plugin, see [protoc-gen-descriptor_set_json documentation](cmd/protoc-gen-descriptor_set_json/README.md).

### JSON Conversion

Once your schema is loaded, convert protobuf messages to JSON:

```sql
> SELECT pb_message_to_json(test_schema(), '.Person', pb_data) FROM Example;
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

**Pro Tip**: Create a VIEW for easier debugging:

```sql
CREATE VIEW ExampleDebugView AS
  SELECT
      *,
      pb_message_to_json(test_schema(), '.Person', pb_data) AS pb_data_json
  FROM Example;
```

```sql
> SELECT pb_data_json FROM ExampleDebugView;
{"id": 1, "name": "Agent Smith", "email": "smith@example.com", "phones": [{"type": "PHONE_TYPE_WORK", "number": "+81-00-0000-0000"}], "lastUpdated": "2025-06-01T12:34:56.789Z"}

> SELECT pb_data_json->'$.name' FROM ExampleDebugView;
"Agent Smith"
```

## Modifying Messages

You can modify existing protobuf messages by setting, adding, or clearing fields:

```sql
-- Set individual fields
SELECT pb_message_set_int32_field(
  pb_message_set_string_field(pb_data, 1, 'New Name'),
  2, 25
) AS updated_person FROM Example LIMIT 1;

-- Add elements to repeated fields (create a new phone number)
SELECT pb_message_add_repeated_message_field_element(
  pb_data,
  4, -- phones field
  pb_message_set_enum_field(
    pb_message_set_string_field(pb_message_new(), 1, '+81-00-0000-0001'), -- number field
    2, 1 -- type = PHONE_TYPE_MOBILE
  )
) AS person_with_phone FROM Example LIMIT 1;

-- Clear fields
SELECT pb_message_clear_string_field(
  pb_data,
  3  -- Clear email field
) AS person_no_email FROM Example LIMIT 1;

-- Modify specific elements in repeated fields
SELECT pb_message_set_repeated_message_field_element(
  pb_data,
  4, 0, -- phones[0]
  pb_message_set_enum_field(
    pb_message_get_repeated_message_field_element(pb_data, 4, 0), -- get existing phone (phones[0])
    2, 1 -- set type = PHONE_TYPE_MOBILE
  )
) AS person_mobile_phone FROM Example LIMIT 1;
```