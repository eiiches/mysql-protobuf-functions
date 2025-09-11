# Tutorial: JSON Integration

*Converting protobuf messages to JSON using schema-aware functions*

This tutorial demonstrates JSON conversion features that require protobuf schema information. These **schema-aware functions** provide human-readable output and advanced protobuf features.

> **What are schema-aware functions?** These functions require your protobuf schema definition to work. They provide rich features like JSON conversion, proper enum names, timestamp formatting, and field name resolution.

> **Prerequisites:**
> - This tutorial requires `protobuf-json.sql` to be installed (which includes schema loading functionality). If you haven't already, see the [Installation Guide](installation.md) to get started.
> - Basic familiarity with protobuf concepts is helpful but not required

## Loading Protobuf Schema

Before using schema-aware functions, you need to load your protobuf schema definition into MySQL. For this tutorial, we'll use the `protoc-gen-mysql` plugin approach.

### Install the Plugin

First, install the protoc plugin:

```bash
go install github.com/eiiches/mysql-protobuf-functions/cmd/protoc-gen-mysql@latest
```

### Generate Schema Function


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

Generate the schema function:

```console
$ protoc --mysql_out=. --mysql_opt=name=person_schema person.proto
```

This generates `person_schema.sql` containing a stored function. Load it into MySQL:

```bash
mysql -u your_username -p your_database < person_schema.sql
```

You can now use the generated function `person_schema()` directly as the first argument to `pb_message_to_json()`.

For other schema loading methods including `pb_descriptor_set_build()` and Go integration, see [Schema Loading Guide](schema-loading.md). For detailed information about the plugin, see [protoc-gen-mysql documentation](../cmd/protoc-gen-mysql/README.md).

## Converting Messages to JSON

Once your schema is loaded, convert protobuf messages to human-readable JSON:

```sql
-- Set up sample data (same as previous tutorials)
CREATE TABLE Person (pb_data BLOB);
INSERT INTO Person (pb_data) VALUES (_binary X'0a0b4167656e7420536d69746810011a11736d697468406578616d706c652e636f6d22140a102b38312d30302d303030302d3030303010032a0c08f091f1c10610c0de9cf802');

-- Convert to JSON using schema
SELECT pb_message_to_json(person_schema(), '.Person', pb_data, NULL, NULL) FROM Person;
-- Result: {
--   "id": 1,
--   "name": "Agent Smith",
--   "email": "smith@example.com",
--   "phones": [
--     {"type": "PHONE_TYPE_WORK", "number": "+81-00-0000-0000"}
--   ],
--   "lastUpdated": "2025-06-01T12:34:56.789Z"
-- }
```

Notice the differences from low-level functions:
- ✅ **Field names** instead of numbers ("name" vs field 1)
- ✅ **Enum names** instead of integers ("PHONE_TYPE_WORK" vs 3)
- ✅ **Timestamp formatting** (ISO 8601 format vs raw seconds/nanos)
- ✅ **Proper JSON structure** with nested objects and arrays

> **⚠️ Important:** JSON output is intended for human consumption (debugging, inspection, reporting) only. Do not use JSON output for programmatic access in production code. Field names can change during schema evolution, breaking your queries. Use low-level functions with field numbers for reliable programmatic access, as field numbers remain stable across schema changes.

## Creating JSON Views

Schema-aware functions are perfect for creating debug and inspection views:

```sql
-- Create a debug view with JSON representation
CREATE VIEW PersonDebugView AS
SELECT
    *,
    pb_message_to_json(person_schema(), '.Person', pb_data, NULL, NULL) AS pb_data_json
FROM Person;

-- Use MySQL JSON functions for powerful queries
SELECT
  pb_data_json->>'$.name' AS name,
  pb_data_json->>'$.email' AS email,
  JSON_LENGTH(pb_data_json, '$.phones') AS phone_count
FROM PersonDebugView;

-- Find people with mobile phones
SELECT pb_data_json->>'$.name' AS name
FROM PersonDebugView
WHERE JSON_CONTAINS(pb_data_json, '"PHONE_TYPE_MOBILE"', '$.phones[*].type');
```

## Next Steps

You now understand **schema-aware functions** and their advantages for human-readable output and rich protobuf features.

**Continue learning:**
- **Learn low-level functions** → [Tutorial: Basics](tutorial-basics.md) and [Tutorial: Modification](tutorial-modification.md)
- **See real-world patterns** → [User Guide](user-guide.md)
- **Explore all functions** → [Function Reference](function-reference.md)
- **Learn advanced topics** → [API Guide](api-guide.md)

The combination of low-level and schema-aware functions gives you the flexibility to optimize for performance when needed while providing rich debugging and integration capabilities.
