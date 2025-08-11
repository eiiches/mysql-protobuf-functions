# protonumberjson

The `protonumberjson` package provides JSON serialization for Protocol Buffer messages using field numbers as JSON object keys for all message types, with the following key features:

1. **Field numbers** are used as JSON object keys instead of field names for all messages (including well-known types)
2. **Enum values** are serialized as numbers instead of string names
3. **64-bit integers** are serialized as numbers instead of strings
4. **Consistent format** - all protobuf messages follow the same serialization approach

## Overview

This package serializes all protobuf messages to JSON using field numbers as object keys. This provides a consistent approach that treats well-known types the same as regular messages, using field numbers for robustness against protobuf schema evolution.

Originally designed as a variation of the ProtoJSON format, this package evolved to drop special formatting for well-known types due to ProtoJSON's limitations with Duration and Timestamp representations (ProtoJSON only supports timestamps between years 0001 and 9999). This change made the package more deviated from ProtoJSON but provides better support for edge cases and maintains consistency across all message types.

This approach provides robustness against protobuf schema evolution, particularly field renames and enum value renames.

## Features

- **Field Number Keys**: Serializes all protobuf messages to JSON using field numbers (e.g., `"1"`, `"2"`) as object keys
- **Consistent Treatment**: All message types, including well-known types, use the same field number approach
- **Type Safety**: Proper overflow checking for integer conversions
- **64-bit Integer Support**: 64-bit integers serialized as numbers (not strings)
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
    // Marshal a wrapper type (now uses field numbers)
    msg := &wrapperspb.StringValue{Value: "hello world"}
    jsonBytes, err := protonumberjson.Marshal(msg)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(jsonBytes))
    // Output: {"1":"hello world"}

    // Convert to JSON tree (interface{}) for programmatic access
    tree, err := protonumberjson.ToJsonTree(msg)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Tree: %v\n", tree)
    // Output: Tree: map[1:hello world]

    // Marshal a Timestamp (now uses field numbers)
    ts := &timestamppb.Timestamp{Seconds: 1234567890, Nanos: 123456789}
    jsonBytes, err = protonumberjson.Marshal(ts)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(jsonBytes))
    // Output: {"1":1234567890,"2":123456789}

    // Create complex data structures with ToJsonTree
    tree, err = protonumberjson.ToJsonTree(ts)
    if err != nil {
        panic(err)
    }
    mergedData := []interface{}{tree, map[string]interface{}{"extra": "data"}}
    finalBytes, _ := json.Marshal(mergedData)
    fmt.Println(string(finalBytes))
    // Output: [{"1":1234567890,"2":123456789},{"extra":"data"}]
}
```

## JSON Output Format

The package produces JSON where:

- **Field numbers** are used as object keys (as strings) for all message types
- **32-bit and 64-bit integers** are serialized as JSON numbers
- **Repeated fields** are serialized as JSON arrays
- **Map fields** are serialized as JSON objects with string keys
- **Enum fields** are serialized as numbers for robustness against enum value renames
- **All message types** follow the same consistent field number format
- **Default values** are omitted from the JSON output to reduce payload size

### Example Transformations

| Protobuf Message | JSON Output |
|------------------|-------------|
| `StringValue{Value: "test"}` | `{"1": "test"}` |
| `Int64Value{Value: 9223372036854775807}` | `{"1": 9223372036854775807}` |
| `Timestamp{Seconds: 1000, Nanos: 500}` | `{"1": 1000, "2": 500}` |
| `Empty{}` | `{}` |
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
64-bit integers in all message types are serialized as JSON numbers instead of strings. Note that JavaScript and some JSON parsers may have precision limitations with very large integers (beyond 53 bits of precision).

### Consistent Treatment of All Messages
All message types, including Google's well-known types, are treated consistently using field numbers:
- **Predictable Format**: Developers can expect the same serialization approach for all protobuf messages
- **Simplified Logic**: No special cases or exceptions to remember
- **Schema Evolution Benefits**: Well-known types also benefit from field number stability

### Default Value Handling
Fields with default values are omitted from the JSON output based on field presence semantics:
- **Proto2**: Fields that are not explicitly set are omitted; explicitly set fields are included even if they equal the default value
- **Proto3**: Fields with zero values (0, false, empty string, etc.) are omitted unless they have explicit presence (optional fields, message fields, or oneof fields)
- **Repeated fields**: Empty arrays are omitted
- **Map fields**: Empty maps are omitted
- **Message fields**: Unset message fields are omitted
- **Payload optimization**: Omitting unset fields reduces JSON size and network overhead

## API Reference

### Functions

#### `ToJsonTree(m proto.Message) (interface{}, error)`
Converts a protobuf message to a Go JSON tree structure using field numbers as keys.

**Parameters:**
- `m`: The protobuf message to convert

**Returns:**
- `interface{}`: JSON tree structure with field numbers as keys
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
- `[]byte`: JSON representation with field numbers as keys
- `error`: Error if serialization fails

**Special Cases:**
- Returns `nil, nil` for nil messages
- All message types use field numbers as object keys
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
- All well-known types (now using field numbers)
- Edge cases like nil messages and overflow conditions
- Integration scenarios

Run tests with:
```bash
go test ./internal/protonumberjson -v
```