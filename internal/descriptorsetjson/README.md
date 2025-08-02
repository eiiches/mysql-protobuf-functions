# descriptorsetjson

The `descriptorsetjson` package converts Protocol Buffer `FileDescriptorSet` messages to MySQL-compatible JSON format for use with protobuf reflection functions.

## Overview

This package takes an arbitrary `FileDescriptorSet` and converts it to a JSON structure that can be used by MySQL stored functions like `pb_message_to_json` for protobuf parsing and reflection.

## JSON Output Format

The package outputs a 3-element JSON array: `[version, fileDescriptorSet, typeIndex]`

### Element 0: Version
Format version number (currently `1`) for future extensibility.

### Element 1: FileDescriptorSet
The `FileDescriptorSet` serialized using the `protonumberjson` format, which uses field numbers as JSON keys instead of field names.

### Element 2: Type Index
A mapping of fully-qualified type names to their JSON path locations within the FileDescriptorSet:

```json
{
  ".google.protobuf.FileDescriptorSet": [11, "$[1].\"1\"[0]", "$[1].\"1\"[0].\"4\"[0]"],
  ".google.protobuf.FileDescriptorProto": [11, "$[1].\"1\"[0]", "$[1].\"1\"[0].\"4\"[1]"],
  ".google.protobuf.DescriptorProto": [11, "$[1].\"1\"[0]", "$[1].\"1\"[0].\"4\"[2]"],
  ".google.protobuf.FieldDescriptorProto.Type": [14, "$[1].\"1\"[0]", "$[1].\"1\"[0].\"4\"[3].\"4\"[0]"]
}
```

Each type maps to a `TypeIndex` array with:
- `[0]`: Kind (protobuf `FieldDescriptorProto.Type` enum: `11` = TYPE_MESSAGE, `14` = TYPE_ENUM)
- `[1]`: File path (e.g., `"$[1].\"1\"[0]"`)
- `[2]`: Type path (e.g., `"$[1].\"1\"[0].\"4\"[2]"`)

### JSON Path Structure
- `$[0]`: Format version number
- `$[1]`: FileDescriptorSet
- `$[2]`: Type index
- `\"1\"[n]`: File array (field 1 in FileDescriptorSet)
- `\"4\"[n]`: Message types array (field 4 in FileDescriptorProto)
- `\"5\"[n]`: Enum types array (field 5 in FileDescriptorProto)
- `\"3\"[n]`: Nested types array (field 3 in DescriptorProto)

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
    
    version := jsonTree[0]               // Format version
    fileDescriptorSetData := jsonTree[1] // FileDescriptorSet
    typeIndex := jsonTree[2]             // Type index map
    
    fmt.Printf("Type index contains %d types\n", len(typeIndex.(map[string]descriptorsetjson.TypeIndex)))
}
```

## API Reference

### Functions

#### `ToJson(fileDescriptorSet *descriptorpb.FileDescriptorSet) (string, error)`
Converts a `FileDescriptorSet` to a JSON string in MySQL-compatible format.

**Parameters:**
- `fileDescriptorSet`: The protobuf FileDescriptorSet to convert

**Returns:**
- `string`: JSON representation as a 3-element array `[version, fileDescriptorSet, typeIndex]`
- `error`: Error if conversion fails

**Errors:**
- Returns error if `fileDescriptorSet` is nil
- Returns error if JSON marshaling fails

#### `ToJsonTree(fileDescriptorSet *descriptorpb.FileDescriptorSet) ([3]interface{}, error)`
Converts a `FileDescriptorSet` to a Go data structure that can be further manipulated before JSON serialization.

**Parameters:**
- `fileDescriptorSet`: The protobuf FileDescriptorSet to convert

**Returns:**
- `[3]interface{}`: Array containing [version, fileDescriptorSetData, typeIndexMap]
- `error`: Error if conversion fails

**Use Cases:**
- Building more complex data structures that include the descriptor set
- Programmatic manipulation before final JSON serialization
- Integration with other JSON processing pipelines

### Types

#### `TypeIndex`
```go
type TypeIndex [3]interface{}
```
Represents a type reference with:
- `[0]`: Kind (11 for message, 14 for enum)
- `[1]`: File path as JSON path string
- `[2]`: Type path as JSON path string

## Features

- **Arbitrary FileDescriptorSet Support**: Works with any protobuf FileDescriptorSet, not just descriptor.proto
- **Complete Type Indexing**: Builds comprehensive index of all message and enum types
- **Nested Type Support**: Properly handles nested messages and enums
- **Multiple File Support**: Processes FileDescriptorSets with multiple .proto files
- **MySQL Compatibility**: Output format designed for MySQL JSON functions
- **Field Number Keys**: Uses protonumberjson format for schema evolution robustness

## Implementation Notes

The package uses the `protonumberjson` format for the FileDescriptorSet serialization, which provides several advantages:
- **Schema Evolution Robustness**: Field numbers remain stable across protobuf schema changes
- **Consistency**: Field numbers are immutable once assigned
- **MySQL Compatibility**: Optimized for MySQL JSON path expressions

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