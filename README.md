# MySQL Protocol Buffers Functions

[![MySQL Version](https://img.shields.io/badge/MySQL-8.0.17%2B-blue)](https://dev.mysql.com/downloads/mysql/)
[![Aurora MySQL](https://img.shields.io/badge/Aurora%20MySQL-3.04.0%2B-orange)](https://aws.amazon.com/rds/aurora/)

A comprehensive library of MySQL stored functions and procedures for working with Protocol Buffers (protobuf) encoded data directly within MySQL databases. This project enables you to parse, query, and manipulate protobuf messages without requiring external applications or services.

> ⚠️ **Early Development Warning**: This project is in active development and may introduce breaking changes.

## Features

- 🔍 **Field Access**: Extract specific fields from protobuf messages using field numbers
- ✏️ **Message Manipulation**: Create, modify, and update protobuf messages directly in MySQL - set fields, add/remove repeated elements, and clear fields
- 🔄 **JSON Conversion**: Convert protobuf messages to JSON format for easier debugging
- 🛠️ **Pure MySQL Implementation**: Written entirely in MySQL stored functions and procedures - no native libraries or external dependencies required

## Requirements

- **MySQL**: 8.0.17 or later
  - JSON_TABLE() was added in 8.0.4 but requires 8.0.17 for [this critical bugfix](https://bugs.mysql.com/bug.php?id=92976)
- **Aurora MySQL**: 3.04.0 (oldest available 3.x version as of June 2025) or later

## Installation

### Quick Start

1. **Clone the repository** to get the SQL files:
   ```bash
   git clone https://github.com/eiiches/mysql-protobuf-functions.git
   cd mysql-protobuf-functions
   ```

2. **Install core functions**:
   ```bash
   mysql -u your_username -p your_database < protobuf.sql
   mysql -u your_username -p your_database < protobuf-accessors.sql
   ```

3. **Install optional components** (for JSON conversion):
   ```bash
   mysql -u your_username -p your_database < protobuf-descriptor.sql
   mysql -u your_username -p your_database < protobuf-json.sql
   ```

### Important Notes

- ⚠️ All functions and procedures use `_pb_` or `pb_` prefixes to avoid naming conflicts
- ⚠️ The `protobuf-descriptor.sql` script creates tables prefixed with `_Proto_` for schema storage
- 📝 Verify existing routines before installation to prevent overwrites

## Reference Manual

Detailed documentation is available on the [project wiki](https://github.com/eiiches/mysql-protobuf-functions/wiki/Reference).

## Tutorial

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

### Accessing Protobuf Fields

Use the field accessor functions to extract data from protobuf messages:

```sql
-- Get id; The field number for id is 2.
> SELECT pb_message_get_int32_field(pb_data, 2 /* field number */, 0 /* default value */) FROM Example;
1

-- Get email; email is 3
> SELECT pb_message_get_string_field(pb_data, 3, '' /* default value */) FROM Example;
smith@example.com

-- Get phones[0].number; phones is 4, number is 1
> SELECT pb_message_get_string_field(pb_message_get_repeated_message_field(pb_data, 4, 0), 1, '' /* default value */) FROM Example;
+81-00-0000-0000

-- Get phones[0].type; phones is 4, type is 2
> SELECT pb_message_get_enum_field(pb_message_get_repeated_message_field(pb_data, 4, 0), 2, 0 /* default value */) FROM Example;
3
```

### Converting Protobuf to JSON

> **Prerequisites**: This section requires both `protobuf-descriptor.sql` and `protobuf-json.sql` to be installed.

#### Loading Protobuf Schema

First, generate and load your protobuf schema into MySQL:

```console
$ protoc --descriptor_set_out=/dev/stdout --include_imports person.proto | xxd -p -c0
0aff010a1f676f6f676c652f70726f746f6275662f74696d657374616d702e70726f746f120f676f6f676c652e70726f746f627566223b0a0954696d657374616d7012180a077365636f6e647318012001280352077365636f6e647312140a056e616e6f7318022001280552056e616e6f734285010a13636f6d2e676f6f676c652e70726f746f627566420e54696d657374616d7050726f746f50015a32676f6f676c652e676f6c616e672e6f72672f70726f746f6275662f74797065732f6b6e6f776e2f74696d657374616d707062f80101a20203475042aa021e476f6f676c652e50726f746f6275662e57656c6c4b6e6f776e5479706573620670726f746f330aa0030a0c706572736f6e2e70726f746f1a1f676f6f676c652f70726f746f6275662f74696d657374616d702e70726f746f22e6020a06506572736f6e12120a046e616d6518012001280952046e616d65120e0a0269641802200128055202696412140a05656d61696c1803200128095205656d61696c122b0a0670686f6e657318042003280b32132e506572736f6e2e50686f6e654e756d626572520670686f6e6573123d0a0c6c6173745f7570646174656418052001280b321a2e676f6f676c652e70726f746f6275662e54696d657374616d70520b6c617374557064617465641a4c0a0b50686f6e654e756d62657212160a066e756d62657218012001280952066e756d62657212250a047479706518022001280e32112e506572736f6e2e50686f6e655479706552047479706522680a0950686f6e6554797065121a0a1650484f4e455f545950455f554e535045434946494544100012150a1150484f4e455f545950455f4d4f42494c45100112130a0f50484f4e455f545950455f484f4d45100212130a0f50484f4e455f545950455f574f524b1003620670726f746f33
```

Alternatively, if you use Buf, you can use `buf build -o ${name}.binpb` to generate a binary FileDescriptorSet.

```sql
> CALL pb_descriptor_set_load('test', _binary X'0aff010a1f676f6f676c652f70726f746f6275662f74696d657374616d702e70726f746f120f676f6f676c652e70726f746f627566223b0a0954696d657374616d7012180a077365636f6e647318012001280352077365636f6e647312140a056e616e6f7318022001280552056e616e6f734285010a13636f6d2e676f6f676c652e70726f746f627566420e54696d657374616d7050726f746f50015a32676f6f676c652e676f6c616e672e6f72672f70726f746f6275662f74797065732f6b6e6f776e2f74696d657374616d707062f80101a20203475042aa021e476f6f676c652e50726f746f6275662e57656c6c4b6e6f776e5479706573620670726f746f330aa0030a0c706572736f6e2e70726f746f1a1f676f6f676c652f70726f746f6275662f74696d657374616d702e70726f746f22e6020a06506572736f6e12120a046e616d6518012001280952046e616d65120e0a0269641802200128055202696412140a05656d61696c1803200128095205656d61696c122b0a0670686f6e657318042003280b32132e506572736f6e2e50686f6e654e756d626572520670686f6e6573123d0a0c6c6173745f7570646174656418052001280b321a2e676f6f676c652e70726f746f6275662e54696d657374616d70520b6c617374557064617465641a4c0a0b50686f6e654e756d62657212160a066e756d62657218012001280952066e756d62657212250a047479706518022001280e32112e506572736f6e2e50686f6e655479706552047479706522680a0950686f6e6554797065121a0a1650484f4e455f545950455f554e535045434946494544100012150a1150484f4e455f545950455f4d4f42494c45100112130a0f50484f4e455f545950455f484f4d45100212130a0f50484f4e455f545950455f574f524b1003620670726f746f33');
```

You can now reference this schema using the identifier `'test'`.

#### JSON Conversion

Once your schema is loaded, convert protobuf messages to JSON:

```sql
> SELECT pb_message_to_json('test', '.Person', pb_data) FROM Example;
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
      pb_message_to_json('test', '.Person', pb_data) AS pb_data_json
  FROM Example;
```

```sql
> SELECT pb_data_json FROM ExampleDebugView;
{"id": 1, "name": "Agent Smith", "email": "smith@example.com", "phones": [{"type": "PHONE_TYPE_WORK", "number": "+81-00-0000-0000"}], "lastUpdated": "2025-06-01T12:34:56.789Z"}

> SELECT pb_data_json->'$.name' FROM ExampleDebugView;
"Agent Smith"
```

### Modifying Messages

You can modify existing protobuf messages by setting, adding, or clearing fields:

```sql
-- Set individual fields
SELECT pb_message_set_int32_field(
  pb_message_set_string_field(pb_data, 1, 'New Name'),
  2, 25
) AS updated_person FROM Example LIMIT 1;

-- Add elements to repeated fields (create a new phone number)
SELECT pb_message_add_repeated_message_field(
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
SELECT pb_message_set_repeated_message_field(
  pb_data,
  4, 0, -- phones[0]
  pb_message_set_enum_field(
    pb_message_get_repeated_message_field(pb_data, 4, 0), -- get existing phone (phones[0])
    2, 1 -- set type = PHONE_TYPE_MOBILE
  )
) AS person_mobile_phone FROM Example LIMIT 1;
```

### Advanced: Indexing Protobuf Fields

While MySQL doesn't allow using stored functions in functional indexes or generated columns, you can use `TRIGGER` to mimic a generated column and create an `INDEX` or any other constraints on that generated column.

#### Using Protobuf Fields as a PRIMARY KEY

Let's add an `id` column to the Example table, populate the column from `pb_data`, and make the column the primary key of the table.

```sql
> ALTER TABLE Example ADD COLUMN id INT NOT NULL FIRST;
> UPDATE Example SET id = pb_message_get_int32_field(pb_data, 2, 0);
> ALTER TABLE Example ADD PRIMARY KEY (id);

> SHOW CREATE TABLE Example;
CREATE TABLE `Example` (
  `id` int NOT NULL,
  `pb_data` blob,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

> SELECT id FROM Example;
1
```

Rather than manually keeping the MySQL `id` column and protobuf `id` field in sync, you can use triggers to automatically populate the `id` column from the protobuf data.

```sql
CREATE TRIGGER Example_set_id_on_update
   BEFORE UPDATE ON Example
   FOR EACH ROW
      SET NEW.id = pb_message_get_int32_field(NEW.pb_data, 2, 0);

CREATE TRIGGER Example_set_id_on_insert
   BEFORE INSERT ON Example
   FOR EACH ROW
      SET NEW.id = pb_message_get_int32_field(NEW.pb_data, 2, 0);
```

```sql
-- With TRIGGER, id is automatically derived from pb_data.
-- protoc --encode=Person person.proto <<-EOF | xxd -p -c0
-- name: "Thomas A. Anderson"
-- id: 2
-- email: "thomas@example.com"
-- phones: [{type: PHONE_TYPE_HOME, number: "+81-00-0000-0000"}]
-- last_updated: {seconds: 1748781296, nanos: 789000000}
-- EOF
> INSERT INTO Example (pb_data) VALUES (_binary X'0a1254686f6d617320412e20416e646572736f6e10021a1274686f6d6173406578616d706c652e636f6d22140a102b38312d30302d303030302d3030303010022a0c08f091f1c10610c0de9cf802');
```

#### Enforcing a UNIQUE constraint on a Protobuf Field

You can also add a `name` column that is automatically derived from `pb_data` and create a UNIQUE INDEX on that column. This enforces name uniqueness and enables faster name-based lookup.

```sql
ALTER TABLE Example ADD COLUMN name VARCHAR(255) NOT NULL AFTER id;
UPDATE Example SET name = pb_message_get_string_field(pb_data, 1, '');
ALTER TABLE Example ADD UNIQUE INDEX (name);

CREATE TRIGGER Example_set_name_on_update
   BEFORE UPDATE ON Example
   FOR EACH ROW
      SET NEW.name = pb_message_get_string_field(NEW.pb_data, 1, '');

CREATE TRIGGER Example_set_name_on_insert
   BEFORE INSERT ON Example
   FOR EACH ROW
      SET NEW.name = pb_message_get_string_field(NEW.pb_data, 1, '');
```

```sql
-- protoc --encode=Person person.proto <<-EOF | xxd -p -c0
-- name: "Mr. Anderson"
-- id: 2
-- email: "thomas@example.com"
-- phones: [{type: PHONE_TYPE_HOME, number: "+81-00-0000-0000"}]
-- last_updated: {seconds: 1748781296, nanos: 789000000}
-- EOF
> UPDATE Example SET pb_data = _binary X'0a0c4d722e20416e646572736f6e10021a1274686f6d6173406578616d706c652e636f6d22140a102b38312d30302d303030302d3030303010022a0c08f091f1c10610c0de9cf802' WHERE id = 2;

> SELECT name FROM Example;
Agent Smith
Mr. Anderson -- Automatically updated by TRIGGER
```

TODO: Multi-Valued Index on Protobuf Fields

## Roadmap

- [ ] **[Editions](https://protobuf.dev/editions/overview/) Support in JSON Conversion**

## Limitations

- [Groups](https://protobuf.dev/programming-guides/encoding/#groups) are not supported.

## Known Issues

### MySQL Stored Program Cache Bug

**Issue**: When many stored functions are used in a single connection and the `stored_program_cache` limit is reached, MySQL exhibits unpredictable behavior:

- **MySQL 9.3.0**: Functions silently return `NULL` instead of the expected result
- **MySQL 8.0.x**: Functions fail with `Function does not exist` error

**Root Cause**: [MySQL Bug #95825](https://bugs.mysql.com/bug.php?id=95825)

**Workaround**: Increase the stored program cache size:
```sql
SET GLOBAL stored_program_cache = 512;  -- Default is 256
```

**Impact**:
- Most applications won't encounter this issue as the default cache size (256) is sufficient for typical usage
- This primarily affects comprehensive test suites or applications using many different protobuf functions in a single connection
- Related discussion: [Percona Forums](https://forums.percona.com/t/intermittent-stored-function-does-not-exist-problem/5143)

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
