# Conformance Test Status

This document tracks the current status of Protocol Buffers conformance tests for the mysql-protobuf library.

## Summary

> 2751 successes, 414 skipped, 17 expected failures, 0 unexpected failures.

## Skipped Tests

Conformance tests for the following protobuf features are **skipped**.
These are documented as limitations or unimplemented features of this library.

* **[Text Format](https://protobuf.dev/reference/protobuf/textformat-spec/)** (414 tests) - Textual protobuf representation, which is unlikely to be implemented.

## Failing Tests

For a machine-readable list of all test failures, see [conformance/expected_failures.txt](../conformance/expected_failures.txt).

| Test Name | Failure Message | Description |
|-----------|-----------------|-------------|
| `Required.Proto2.JsonInput.DoubleFieldMaxNegativeValue.JsonOutput` | Output was not equivalent to reference message: modified: optional_double: -2.22507e-308 -> -2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto2.JsonInput.DoubleFieldMaxNegativeValue.ProtobufOutput` | Output was not equivalent to reference message: modified: optional_double: -2.22507e-308 -> -2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto2.JsonInput.DoubleFieldMinPositiveValue.JsonOutput` | Output was not equivalent to reference message: modified: optional_double: 2.22507e-308 -> 2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto2.JsonInput.DoubleFieldMinPositiveValue.ProtobufOutput` | Output was not equivalent to reference message: modified: optional_double: 2.22507e-308 -> 2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto2.JsonInput.DoubleFieldQuotedExponentialValue.JsonOutput` | Output was not equivalent to reference message: modified: optional_double: 2.22507e-308 -> 2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto2.JsonInput.DoubleFieldQuotedExponentialValue.ProtobufOutput` | Output was not equivalent to reference message: modified: optional_double: 2.22507e-308 -> 2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto3.JsonInput.DoubleFieldMaxNegativeValue.JsonOutput` | Output was not equivalent to reference message: modified: optional_double: -2.22507e-308 -> -2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto3.JsonInput.DoubleFieldMaxNegativeValue.ProtobufOutput` | Output was not equivalent to reference message: modified: optional_double: -2.22507e-308 -> -2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto3.JsonInput.DoubleFieldMinPositiveValue.JsonOutput` | Output was not equivalent to reference message: modified: optional_double: 2.22507e-308 -> 2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto3.JsonInput.DoubleFieldMinPositiveValue.ProtobufOutput` | Output was not equivalent to reference message: modified: optional_double: 2.22507e-308 -> 2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto3.JsonInput.DoubleFieldQuotedExponentialValue.JsonOutput` | Output was not equivalent to reference message: modified: optional_double: 2.22507e-308 -> 2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Required.Proto3.JsonInput.DoubleFieldQuotedExponentialValue.ProtobufOutput` | Output was not equivalent to reference message: modified: optional_double: 2.22507e-308 -> 2.2250700000000003e-308 | See [known-issues.md](known-issues.md#double-precision-issues-in-json-parsing) |
| `Recommended.Proto2.JsonInput.FieldNameExtension.Validator` | JSON payload validation failed | Extension field name validation not supported |
| `Recommended.Proto2.ProtobufInput.ValidMessageSetEncoding.SubmessageEncoding.NotUnknown.ProtobufOutput` | Failed to parse input or produce output | MessageSet encoding not supported (GROUP wire type) |
| `Required.Proto2.ProtobufInput.MessageSetEncoding.UnknownExtension.ProtobufOutput` | Failed to parse input or produce output | MessageSet encoding not supported (GROUP wire type) |
| `Required.Proto2.ProtobufInput.ValidMessageSetEncoding.OutOfOrderGroupsEntries.ProtobufOutput` | Failed to parse input or produce output | MessageSet encoding not supported (GROUP wire type) |
| `Required.Proto2.ProtobufInput.ValidMessageSetEncoding.ProtobufOutput` | Failed to parse input or produce output | MessageSet encoding not supported (GROUP wire type) |

## Tests Not Run

Conformance tests for Editions are not run unless `--maximum_edition=` is specified. Since we don't yet support Editions, we currently skip all Editions-related tests.

* **[Editions](https://protobuf.dev/editions/overview/)** - Newer protobuf syntax that is in our roadmap.

---

*Last Updated: 2025-09-11*
