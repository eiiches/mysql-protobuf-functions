# descriptorsetjson

The `descriptorsetjson` package converts Protocol Buffer `FileDescriptorSet` messages to MySQL-compatible JSON format using the new `DescriptorSet` protobuf message structure with separate MessageTypeIndex and EnumTypeIndex capabilities.

## Overview

This package takes an arbitrary `FileDescriptorSet` and converts it to a `DescriptorSet` message serialized using the [protonumberjson](../protonumberjson/README.md) format. The output provides separate message and enum type indexing with comprehensive field and enum mapping capabilities for MySQL stored functions.

## JSON Output Format

The package outputs a `DescriptorSet` protobuf message serialized using the [protonumberjson](../protonumberjson/README.md) format:

```json
{
  "1": { ... },  // FileDescriptorSet
  "2": { ... },  // MessageTypeIndex map
  "3": { ... }   // EnumTypeIndex map
}
```

### Field "1": FileDescriptorSet
The original `FileDescriptorSet` serialized using the [protonumberjson](../protonumberjson/README.md) format, which uses field numbers as JSON keys instead of field names.

### Field "2": MessageTypeIndex
A mapping of fully-qualified message type names to their `MessageTypeIndex` messages with enhanced indexing capabilities:

### Field "3": EnumTypeIndex  
A mapping of fully-qualified enum type names to their `EnumTypeIndex` messages with enhanced indexing capabilities:

**MessageTypeIndex example:**
```json
{
  ".google.protobuf.FileDescriptorSet": {
    "1": "$.\"1\"[0]",  // file_json_path (relative to FileDescriptorSet)
    "2": "$.\"1\"[0].\"4\"[0]",  // type_json_path (relative to FileDescriptorSet)
    "3": { "file": 0 },  // field_name_index (field name -> array index)
    "4": { "1": 0 }      // field_number_index (field number -> array index)
  }
}
```

**EnumTypeIndex example:**
```json
{
  ".google.protobuf.Status": {
    "1": "$.\"1\"[0]",  // file_json_path (relative to FileDescriptorSet)
    "2": "$.\"1\"[0].\"5\"[0]",  // type_json_path (relative to FileDescriptorSet)
    "3": { "ACTIVE": 0 },   // enum_name_index (enum name -> array index)
    "4": { "1": 0 }         // enum_number_index (enum number -> array index)
  }
}
```

Each `MessageTypeIndex` contains:
- **Field "1"** (`file_json_path`): JSON path to file from FileDescriptorSet root
- **Field "2"** (`type_json_path`): JSON path to message type from FileDescriptorSet root
- **Field "3"** (`field_name_index`): Map from field names to array indices
- **Field "4"** (`field_number_index`): Map from field numbers to array indices

Each `EnumTypeIndex` contains:
- **Field "1"** (`file_json_path`): JSON path to file from FileDescriptorSet root
- **Field "2"** (`type_json_path`): JSON path to enum type from FileDescriptorSet root
- **Field "3"** (`enum_name_index`): Map from enum value names to array indices
- **Field "4"** (`enum_number_index`): Map from enum value numbers to array indices

### JSON Path Structure
- `$."1"`: FileDescriptorSet (field 1 of DescriptorSet)
- `$."2"`: MessageTypeIndex map (field 2 of DescriptorSet)
- `$."3"`: EnumTypeIndex map (field 3 of DescriptorSet)
- `$."1"[n]`: File array (field 1 in FileDescriptorSet) - **paths are relative to this level**
- `$."1"[n]."4"[n]`: Message types array (field 4 in FileDescriptorProto)
- `$."1"[n]."5"[n]`: Enum types array (field 5 in FileDescriptorProto)
- `$."1"[n]."4"[n]."3"[n]`: Nested types array (field 3 in DescriptorProto)
- `$."1"[n]."4"[n]."4"[n]`: Nested enum types array (field 4 in DescriptorProto)

### Type Name Format
- Fully-qualified names start with `.` (e.g., `.google.protobuf.FieldDescriptorProto`)
- Nested types use dot notation (e.g., `.google.protobuf.FieldDescriptorProto.Type`)

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/eiiches/mysql-protobuf-functions/internal/descriptorsetjson"
    "github.com/eiiches/mysql-protobuf-functions/internal/protoreflectutils"
    "google.golang.org/protobuf/types/descriptorpb"
)

func main() {
    // Get a file descriptor (e.g., descriptor.proto)
    fileDescriptor := descriptorpb.File_google_protobuf_descriptor_proto

    // Build a FileDescriptorSet with dependencies
    fileDescriptorSet := protoreflectutils.BuildFileDescriptorSetWithDependencies(fileDescriptor)

    // Convert to JSON string
    jsonStr, err := descriptorsetjson.ToJson(fileDescriptorSet)
    if err != nil {
        log.Fatalf("Failed to convert to JSON: %v", err)
    }

    fmt.Println(jsonStr)

    // Or get as Go data structure
    jsonTree, err := descriptorsetjson.ToJsonTree(fileDescriptorSet)
    if err != nil {
        log.Fatalf("Failed to convert to JSON tree: %v", err)
    }

    // Access the DescriptorSet fields
    resultMap := jsonTree.(map[string]interface{})
    fileDescriptorSetData := resultMap["1"] // FileDescriptorSet
    messageTypeIndex := resultMap["2"]      // MessageTypeIndex map
    enumTypeIndex := resultMap["3"]         // EnumTypeIndex map

    fmt.Printf("Message type index contains %d types\n", len(messageTypeIndex.(map[string]interface{})))
    fmt.Printf("Enum type index contains %d types\n", len(enumTypeIndex.(map[string]interface{})))
}
```

## API Reference

### Functions

#### `ToJson(fileDescriptorSet *descriptorpb.FileDescriptorSet) (string, error)`
Converts a `FileDescriptorSet` to a JSON string using the `DescriptorSet` protobuf message format.

**Parameters:**
- `fileDescriptorSet`: The protobuf FileDescriptorSet to convert

**Returns:**
- `string`: JSON representation of a `DescriptorSet` message with enhanced TypeIndex capabilities
- `error`: Error if conversion fails

**Errors:**
- Returns error if `fileDescriptorSet` is nil
- Returns error if JSON marshaling fails

#### `ToJsonTree(fileDescriptorSet *descriptorpb.FileDescriptorSet) (interface{}, error)`
Converts a `FileDescriptorSet` to a Go data structure using the `DescriptorSet` protobuf message format.

**Parameters:**
- `fileDescriptorSet`: The protobuf FileDescriptorSet to convert

**Returns:**
- `interface{}`: Go data structure representing a `DescriptorSet` message with separate message and enum indexing
- `error`: Error if conversion fails

**Use Cases:**
- Building more complex data structures that include the descriptor set
- Programmatic manipulation before final JSON serialization
- Integration with other JSON processing pipelines
- Accessing separate message and enum type indexing capabilities

### Protobuf Messages

The package uses the following protobuf messages defined in `src/descriptor_set.proto`:

#### `DescriptorSet`
```protobuf
message DescriptorSet {
  google.protobuf.FileDescriptorSet file_descriptor_set = 1;
  map<string, MessageTypeIndex> message_type_index = 2;
  map<string, EnumTypeIndex> enum_type_index = 3;
}
```

#### `MessageTypeIndex`
```protobuf
message MessageTypeIndex {
  string file_json_path = 1;
  string type_json_path = 2;
  map<string, int32> field_name_index = 3;
  map<int32, int32> field_number_index = 4;
}
```

#### `EnumTypeIndex`
```protobuf
message EnumTypeIndex {
  string file_json_path = 1;
  string type_json_path = 2;
  map<string, int32> enum_name_index = 3;
  map<int32, int32> enum_number_index = 4;
}
```

## Features

- **Separate Type Indexes**: Distinct MessageTypeIndex and EnumTypeIndex for optimized lookups
- **FileDescriptorSet-Relative Paths**: All JSON paths are relative to FileDescriptorSet root for consistency
- **DescriptorSet Message**: Uses structured protobuf message format instead of ad-hoc array structure
- **Arbitrary FileDescriptorSet Support**: Works with any protobuf FileDescriptorSet, not just descriptor.proto
- **Complete Type Indexing**: Builds comprehensive index of all message and enum types with enhanced capabilities
- **Nested Type Support**: Properly handles nested messages and enums with full indexing
- **Multiple File Support**: Processes FileDescriptorSets with multiple .proto files
- **MySQL Compatibility**: Output format optimized for MySQL JSON path expressions
- **Field Number Keys**: Uses protonumberjson format for schema evolution robustness
- **ProtoNumberJSON Serialization**: Consistent field number-based JSON keys for all protobuf messages

## Implementation Notes

The package uses the new `DescriptorSet` protobuf message structure with separate MessageTypeIndex and EnumTypeIndex maps, serialized using the [protonumberjson](../protonumberjson/README.md) format, which provides several advantages:
- **Separate Type Indexing**: Distinct indexes for messages and enums for optimized lookups
- **FileDescriptorSet-Relative Paths**: Consistent path structure relative to FileDescriptorSet root
- **Enhanced Indexing**: Provides comprehensive field and enum indexing for efficient lookups
- **Structured Format**: Uses proper protobuf message structure instead of ad-hoc arrays
- **Schema Evolution Robustness**: Field numbers remain stable across protobuf schema changes
- **Consistency**: Field numbers are immutable once assigned
- **MySQL Compatibility**: Optimized for MySQL JSON path expressions with predictable field numbering
- **Type Safety**: Strongly-typed protobuf messages ensure consistent data structure

## Testing

The package includes comprehensive tests covering:
- Basic FileDescriptorSet conversion
- Nested type handling
- Multiple file processing
- Empty and nil input handling
- Type index validation
- JSON path correctness

Run tests with:
```bash
go test ./internal/descriptorsetjson -v
```

## Error Handling

The package provides detailed error messages for common failure scenarios:
- Nil FileDescriptorSet input
- JSON marshaling failures
- Internal conversion errors

All errors include context about what operation failed to aid in debugging.
