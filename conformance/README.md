# MySQL Protobuf Conformance Testing

This directory contains the MySQL Protocol Buffers conformance testing implementation that validates the correctness and completeness of the MySQL protobuf functions library against the official Protocol Buffers conformance test suite.

## Overview

The conformance testing framework ensures that the MySQL protobuf implementation correctly handles the same test cases used to validate other Protocol Buffers implementations like C++, Java, Python, etc. This provides confidence that the MySQL implementation is interoperable with other protobuf systems.

## Architecture

### Components

- **`mysql-conformance`**: The main conformance test program that implements the official protobuf conformance testing protocol
- **`conformance_test_runner`**: The official Google protobuf conformance test runner (built from the protobuf submodule)
- **Schema Files**: Generated protobuf descriptor sets and SQL schema functions for test message types
- **Test Data**: Binary protobuf test files and schema definitions

### Protocol Implementation

The `mysql-conformance` program implements the [official conformance testing protocol](https://github.com/protocolbuffers/protobuf/tree/main/conformance):

1. **Input**: Reads `ConformanceRequest` protobuf messages from stdin
2. **Processing**: 
   - Parses binary protobuf data using MySQL functions
   - Converts between different formats (binary ↔ JSON) using schema-aware MySQL functions
   - Handles error cases and unsupported features appropriately
3. **Output**: Writes `ConformanceResponse` protobuf messages to stdout

### MySQL Integration

The conformance tests validate the core MySQL protobuf functions:

- **`pb_message_to_json()`**: Convert binary protobuf to human-readable JSON
- **`pb_json_to_message()`**: Convert JSON back to binary protobuf with proper proto3 default value handling
- **Schema functions**: Generated from conformance test message descriptors
- **Error handling**: Proper classification of parse errors vs runtime errors

## Files

### Source Files

- **`main.go`**: CLI application setup with database connection and debug options
- **`conformance.go`**: Core conformance testing protocol implementation
- **`conformance.pb.go`**: Generated Go protobuf code for conformance protocol messages
- **`Makefile`**: Build automation for all components

### Generated Files

- **`mysql-conformance`**: Compiled conformance test binary
- **`conformance_test_messages.binpb`**: Protobuf descriptor set for test messages
- **`conformance_test_messages_schema.sql`**: MySQL function containing test message schemas
- **`protobuf/conformance_test_runner`**: Official conformance test runner binary

### Test Data

- **`protobuf/`**: Git submodule containing the official Google protobuf source and conformance tests

## Usage

### Prerequisites

- Go 1.19+ for building the conformance test program
- MySQL 8.0+ with the protobuf functions library loaded
- CMake and C++ compiler for building the official conformance test runner
- `protoc` Protocol Buffers compiler

### Setup

```bash
# Build all components (runner, conformance program, schema)
make all

# Build only the MySQL conformance test program
make build-mysql-conformance

# Build only the official conformance test runner
make build-runner

# Generate and load test message schema into MySQL
make build-schema
make load-schema
```

### Running Tests

```bash
# Run all conformance tests
make test

# Run tests with debug logging
make test-debug

# Run specific test patterns (use full test names)
make test-filter TEST_NAME="Required.Proto3.JsonInput.DurationNegativeNanos.JsonOutput"

# Validate setup
make validate
```

### Manual Testing

```bash
# Run conformance test directly (output to file for performance)
./protobuf/conformance_test_runner ./mysql-conformance \
    --database "root@tcp(127.0.0.100:13306)/test" \
    >test-output.log 2>&1

# With debug logging (output to file)
./protobuf/conformance_test_runner ./mysql-conformance \
    --database "root@tcp(127.0.0.100:13306)/test" \
    --debug --debug >test-debug.log 2>&1

# Read test results
grep -i duration test-output.log
grep "FAILED\|PASSED" test-output.log | tail -5
```

## Implementation Details

### Supported Features

- **Binary protobuf parsing**: Full wire format support for all protobuf types
- **JSON input and output**: Schema-aware JSON parsing and serialization
- **Proto3 semantics**: Proper handling of default values and field presence
- **Proto2 support**: Basic proto2 message handling
- **Error handling**: Correct classification of parse errors, runtime errors, and unsupported features

### Limitations

- **JSPB format**: Not supported (JavaScript protobuf format)
- **Text format**: Not supported (protobuf text format)
- **GROUP fields**: Not supported (deprecated proto2 feature)

### Error Handling

The implementation correctly distinguishes between:

- **Parse errors**: Invalid protobuf wire format data
- **Runtime errors**: Database connection issues, schema problems, etc.
- **Skipped tests**: Unsupported features or formats
- **GROUP field errors**: Specifically handled as unsupported deprecated features

## Schema Management

The conformance tests use dynamically generated schema functions:

1. **Descriptor generation**: `protoc` generates binary descriptor sets from test `.proto` files
2. **Schema conversion**: `protoc-gen-mysql` converts descriptors to JSON format
3. **SQL generation**: Creates `conformance_test_messages_schema()` MySQL function
4. **Schema loading**: Function loaded into MySQL for use by conformance tests

## Development

### Building

```bash
# Clean and rebuild everything
make clean-all
make all

# Rebuild just the conformance program
make build-mysql-conformance

# Regenerate schema files
make build-schema
```

### Debugging

```bash
# Enable debug logging
./mysql-conformance --database "connection_string" --debug --test "TestName"

# Check conformance test runner help
make help-runner
```

### Adding Support for New Features

1. **Update MySQL functions**: Modify the core protobuf functions in `src/`
2. **Update conformance handler**: Modify `conformance.go` to handle new input/output formats
3. **Update error patterns**: Add new error patterns to `isProtobufParseError()` if needed
4. **Test**: Run conformance tests to verify new functionality

## Testing Integration

The conformance testing validates that the MySQL protobuf implementation:

- Produces identical binary output for round-trip conversions
- Handles edge cases correctly (empty messages, default values, etc.)
- Properly implements proto3 and proto2 semantics
- Correctly parses and generates JSON according to protobuf JSON mapping rules
- Maintains compatibility with other protobuf implementations

This ensures the MySQL protobuf functions can reliably interoperate with protobuf data generated by other languages and implementations.