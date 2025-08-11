# generate-descriptorsets

This tool generates SQL functions that embed the `descriptor.proto` schema into MySQL for protobuf reflection and type validation.

## Overview

The tool creates multiple MySQL stored functions that return various protobuf descriptor schemas as JSON, with the main function being `_pb_google_descriptor_proto()` for the core descriptor.proto schema.

## Descriptor Set JSON Format

The output JSON is a versioned 3-element array: `[version, fileDescriptorSet, typeIndex]`

For complete details about the format structure, see the [descriptorsetjson documentation](../../internal/descriptorsetjson/README.md).

## Usage

```bash
go run cmd/generate-descriptorsets/main.go > descriptor-proto.sql
```

The generated SQL can be executed in MySQL to provide runtime protobuf schema information for descriptor parsing and validation.

## Function Usage

```sql
-- Get the descriptor set JSON
SELECT _pb_google_descriptor_proto();

-- Use with pb_message_to_json
SELECT pb_message_to_json(_pb_google_descriptor_proto(), '.google.protobuf.FileDescriptorProto', some_message);
```

## Comparison with protoc-gen-descriptor_set_json

This tool serves a specific purpose in the MySQL protobuf ecosystem:

- **Bootstrap Requirement**: The `pb_message_to_json()` function requires a descriptor set to parse any protobuf message, but parsing descriptor sets themselves creates a circular dependency
- **Built-in Schema Access**: Provides access to the fundamental protobuf schema (`descriptor.proto`) that defines all other protobuf schemas
- **Self-Hosting**: Enables the protobuf reflection system to be self-hosting by embedding the protobuf type system definitions
- **Core Dependency**: Required for parsing `FileDescriptorSet` messages that contain user-defined schemas
- **No protoc Dependency**: Can be used in environments where protoc is not available, as it generates the descriptor set programmatically

For custom schemas, use the `@cmd/protoc-gen-descriptor_set_json/` plugin instead, which leverages this tool's output to process arbitrary protobuf files.

## Equivalent protoc Command

This tool generates the same output as running:

```bash
protoc --descriptor_set_json_out=. \
       --descriptor_set_json_opt=name=_pb_google_descriptor_proto \
       --include_imports \
       google/protobuf/descriptor.proto
```

However, the standalone tool is more convenient as it doesn't require protoc installation or finding the `descriptor.proto` file location.