# Advanced Usage

## Indexing Protobuf Fields

While MySQL doesn't allow using stored functions in functional indexes or generated columns, you can use `TRIGGER` to mimic a generated column and create an `INDEX` or any other constraints on that generated column.

### Using Protobuf Fields as a PRIMARY KEY

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

### Enforcing a UNIQUE constraint on a Protobuf Field

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

### Multi-Valued Index on Protobuf Fields

TODO: Document multi-valued index examples for repeated fields.