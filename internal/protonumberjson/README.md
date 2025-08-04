# protonumberjson

The `protonumberjson` package provides JSON serialization for Protocol Buffer messages that follows the ProtoJSON format exactly, with the following exceptions:

1. **Field numbers** are used as JSON object keys instead of field names for regular messages
2. **Enum values** are serialized as numbers instead of string names
3. **64-bit integers** (including `Int64Value`, `UInt64Value` wrapper types) are serialized as numbers instead of strings
4. **`google.protobuf.Any`** uses field numbers ("1" for type_url, "2" for value) instead of ProtoJSON format

## Overview

This package serializes protobuf messages to JSON using field numbers as object keys for regular messages. Well-known types follow the exact ProtoJSON format (with exceptions listed above), while regular messages use field numbers for their object keys.

This approach provides robustness against protobuf schema evolution, particularly field renames and enum value renames.

## Features

- **Field Number Keys**: Serializes protobuf messages to JSON using field numbers (e.g., `"1"`, `"2"`) as object keys
- **Well-Known Types Support**: Full support for all protobuf well-known types including:
  - `google.protobuf.Timestamp`
  - `google.protobuf.Duration` 
  - `google.protobuf.Struct`
  - `google.protobuf.ListValue`
  - `google.protobuf.Value`
  - `google.protobuf.Empty`
  - `google.protobuf.FieldMask`
  - `google.protobuf.Any`
  - All wrapper types (`StringValue`, `Int64Value`, `BoolValue`, etc.)
- **Type Safety**: Proper overflow checking for integer conversions
- **JSON Compatibility**: 64-bit integers serialized as strings to maintain JSON compatibility
- **Comprehensive Field Support**: Handles scalars, lists, maps, and nested messages

## Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/eiiches/mysql-protobuf-functions/internal/protonumberjson"
    "google.golang.org/protobuf/types/known/timestamppb"
    "google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
    // Marshal a simple wrapper type
    msg := &wrapperspb.StringValue{Value: "hello world"}
    jsonBytes, err := protonumberjson.Marshal(msg)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(jsonBytes))
    // Output: "hello world"

    // Convert to JSON tree (interface{}) for programmatic access
    tree, err := protonumberjson.ToJsonTree(msg)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Tree: %v\n", tree)
    // Output: Tree: hello world

    // Marshal a Timestamp
    ts := &timestamppb.Timestamp{Seconds: 1234567890, Nanos: 123456789}
    jsonBytes, err = protonumberjson.Marshal(ts)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(jsonBytes))
    // Output: "2009-02-13T23:31:30.123456789Z"

    // Create complex data structures with ToJsonTree
    tree, err = protonumberjson.ToJsonTree(ts)
    if err != nil {
        panic(err)
    }
    mergedData := []interface{}{tree, map[string]interface{}{"extra": "data"}}
    finalBytes, _ := json.Marshal(mergedData)
    fmt.Println(string(finalBytes))
    // Output: ["2009-02-13T23:31:30.123456789Z",{"extra":"data"}]
}
```

## JSON Output Format

The package produces JSON where:

- **Field numbers** are used as object keys (as strings)
- **64-bit integers** are serialized as JSON numbers
- **32-bit integers** and smaller are serialized as JSON numbers
- **Repeated fields** are serialized as JSON arrays
- **Map fields** are serialized as JSON objects with string keys
- **Enum fields** are serialized as numbers for robustness against enum value renames
- **Well-known types** follow the exact ProtoJSON format

### Example Transformations

| Protobuf Message | JSON Output |
|------------------|-------------|
| `StringValue{Value: "test"}` (well-known type) | `"test"` |
| `Int64Value{Value: 9223372036854775807}` (well-known type) | `9223372036854775807` |
| `Timestamp{Seconds: 1000, Nanos: 500}` (well-known type) | `"1970-01-01T00:16:40.000000500Z"` |
| `Empty{}` (well-known type) | `{}` |
| `Any{TypeUrl: "type.googleapis.com/...", Value: [...]}` | `{"1": "type.googleapis.com/...", "2": "base64data"}` |
| Regular message with `repeated int32 values = [1, 2, 3]` | `{"1": [1, 2, 3]}` |
| Regular message with `Status status = ACTIVE` (enum value 1) | `{"1": 1}` |

## Design Considerations

### Field Number Keys
Using field numbers instead of field names provides several advantages:
- **Schema Evolution Robustness**: Field numbers remain stable across protobuf schema changes, making the JSON representation immune to field renames
- **Consistency**: Field numbers are immutable once assigned, ensuring data compatibility across different schema versions
- **Compactness**: Shorter keys reduce JSON size compared to descriptive field names

### Enum Serialization
Enums are serialized as their numeric values rather than string names for the same robustness reasons:
- **Enum Value Rename Safety**: Enum values can be renamed (e.g., `USER_ACTIVE` → `ACTIVE`) without breaking stored JSON data
- **Schema Evolution**: Numeric values remain stable while enum value names may change for clarity
- **Consistency**: Matches the field number approach used for message fields

### 64-bit Integer Handling
64-bit integers in both regular messages and wrapper types (`Int64Value`, `UInt64Value`) are serialized as JSON numbers instead of strings (which is the ProtoJSON default). Note that JavaScript and some JSON parsers may have precision limitations with very large integers (beyond 53 bits of precision).

### Well-Known Types
Most well-known types follow the exact ProtoJSON format to maintain compatibility with the broader protobuf ecosystem. This includes `google.protobuf.Timestamp`, `google.protobuf.Duration`, `google.protobuf.Struct`, and most wrapper types. The exceptions are:
- `google.protobuf.Any` uses field numbers for consistency with the package's field number approach
- `google.protobuf.Int64Value` and `google.protobuf.UInt64Value` return numbers instead of strings (consistent with 64-bit integer handling in regular messages)

## API Reference

### Functions

#### `ToJsonTree(m proto.Message) (interface{}, error)`
Converts a protobuf message to a Go JSON tree structure using field numbers as keys.

**Parameters:**
- `m`: The protobuf message to convert

**Returns:**
- `interface{}`: JSON tree structure (field numbers as keys for regular messages, ProtoJSON format for well-known types)
- `error`: Error if conversion fails

**Use Cases:**
- Building complex data structures that include protobuf messages
- Programmatic manipulation of JSON data before serialization
- Creating merged or nested JSON structures

#### `Marshal(m proto.Message) ([]byte, error)`
Serializes a protobuf message to JSON bytes using field numbers as keys. Internally uses `ToJsonTree()`.

**Parameters:**
- `m`: The protobuf message to serialize

**Returns:**
- `[]byte`: JSON representation (field numbers as keys for regular messages, ProtoJSON format for well-known types)
- `error`: Error if serialization fails

**Special Cases:**
- Returns `nil, nil` for nil messages
- Well-known types use exact ProtoJSON format
- Regular messages use field numbers as object keys
- Enums serialized as numbers instead of names
- Applies overflow checking for integer conversions

## Implementation Notes

Error handling includes detailed context about which field failed to serialize, making debugging easier when working with complex nested messages.

## Limitations

- **Serialization Only**: This package only provides marshaling (serialization) functionality. Unmarshaling is not implemented.

## Testing

The package includes comprehensive tests covering:
- All protobuf scalar types
- Repeated and map fields  
- All well-known types
- Edge cases like nil messages and overflow conditions
- Integration scenarios

Run tests with:
```bash
go test ./internal/protonumberjson -v
```