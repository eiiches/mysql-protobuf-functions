# Schema Loading and Persistence

To use functions like `pb_message_to_json()`, you need to provide schema information in the form of [descriptor set JSON](../internal/descriptorsetjson/README.md).
This guide explains the different methods for loading and persisting protobuf schemas in MySQL to enable JSON conversion and reflection capabilities.

## Table of Contents

- [Method 1: Using protoc-gen-mysql (Recommended)](#method-1-using-protoc-gen-mysql-recommended)
- [Method 2: Using pb_build_descriptor_set_json](#method-2-using-pb_build_descriptor_set_json)
- [Method 3: Using Go descriptorsetjson Package](#method-3-using-go-descriptorsetjson-package)

## Method 1: Using protoc-gen-mysql (Recommended)

This method generates a MySQL stored function that returns the schema JSON directly.

For complete documentation including troubleshooting, examples, and advanced options, see the [protoc-gen-mysql README](../cmd/protoc-gen-mysql/README.md).

### When To Use

* Your schema changes infrequently and code generation fits your workflow.
* You want to check the generated schema functions into version control alongside your .proto files.
* You already have binary FileDescriptorSet files and want to generate SQL functions directly from them (standalone mode).
* You're using Buf and want to leverage `buf build` to generate FileDescriptorSet files.

### Step 1: Install the Plugin

```bash
go install github.com/eiiches/mysql-protobuf-functions/cmd/protoc-gen-mysql@latest
```

### Step 2: Generate Schema Function

**Using protoc directly:**
```bash
protoc --mysql_out=. \
       --mysql_opt=name=person_schema \
       person.proto
```

**Using Buf:**
```bash
# buf.gen.yaml
version: v2
plugins:
- local: protoc-gen-mysql
  # You can also use `go run` without installing the plugin.
  # local: ['go', 'run', 'github.com/eiiches/mysql-protobuf-functions/cmd/protoc-gen-mysql@latest']
  out: .
  opt:
  - name=person_schema
  strategy: all
  include_imports: true
  include_wkt: true

# Generate
buf generate
```

**Using standalone mode:**
```bash
# With protoc
protoc --descriptor_set_out=person.binpb --include_imports person.proto

# Or with Buf
buf build -o person.binpb --as-file-descriptor-set

# Then generate SQL function from binary descriptor set
protoc-gen-mysql \
  --descriptor_set_in=person.binpb \
  --name=person_schema \
  --mysql_out=./output
```

All approaches create `person_schema.sql` containing a stored function.

### Step 3: Load into MySQL

```bash
mysql -u your_username -p your_database < person_schema.sql
```

### Step 4: Use Schema Function

```sql
-- The generated function returns schema JSON directly
SELECT pb_message_to_json(person_schema(), '.Person', pb_data, NULL, NULL) FROM Example;
```

## Method 2: Using pb_build_descriptor_set_json

This method converts binary FileDescriptorSet data into the required JSON format at runtime.

### When To Use

* You want to avoid installing additional protoc plugins and do everything within MySQL.
* You already have binary FileDescriptorSet data in MySQL.
* Your schema is dynamic and needs to be loaded from MySQL tables at runtime, where static code generation isn't feasible.

### Prerequisites

- `protobuf-descriptor.sql` and `protobuf-json.sql` need to be installed in your MySQL instance.

### Step 1: Generate Binary Descriptor Set

```bash
# Using protoc
protoc --descriptor_set_out=/dev/stdout --include_imports person.proto | xxd -p -c0

# Or using Buf (if available)
buf build -o schema.binpb
```

### Step 2: Load Schema into MySQL

You have two options for storing the converted schema:

**Option A: User Variable (Session-scoped)**
```sql
SET @my_schema = pb_build_descriptor_set_json(_binary X'0aff010a1f676f6f676c652f...');
```

**Option B: Database Table (Persistent)**
```sql
CREATE TABLE schema_registry (
    schema_name VARCHAR(255) PRIMARY KEY,
    schema_json JSON
);

INSERT INTO schema_registry VALUES (
    'person_schema',
    pb_build_descriptor_set_json(_binary X'0aff010a1f676f6f676c652f...')
);
```

### Step 3: Use Schema for JSON Conversion

```sql
-- With user variable
SELECT pb_message_to_json(@my_schema, '.Person', pb_data, NULL, NULL) FROM Example;

-- With table storage
SELECT pb_message_to_json(
    (SELECT schema_json FROM schema_registry WHERE schema_name = 'person_schema'),
    '.Person',
    pb_data,
    NULL,
    NULL
) FROM Example;
```


## Method 3: Using Go descriptorsetjson Package

For applications that need programmatic control over schema loading, you can use the Go `descriptorsetjson` package.

### When To Use

* You want more control over how schemas are loaded and managed programmatically.
* Integration with schema registries or automated deployment pipelines.

For complete usage examples and API documentation, see the [descriptorsetjson package README](../internal/descriptorsetjson/README.md).
