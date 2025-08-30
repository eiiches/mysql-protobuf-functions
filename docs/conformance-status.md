# Conformance Test Status

This document tracks the current status of Protocol Buffers conformance tests for the mysql-protobuf library.

## Summary

> 1936 successes, 817 skipped, 15 expected failures, 0 unexpected failures.

## Skipped Tests

Conformance tests for the following protobuf features are **skipped**.
These are documented as limitations or unimplemented features of this library.

* **[Text Format](https://protobuf.dev/reference/protobuf/textformat-spec/)** - Textual protobuf representation, which is unlikely to be implemented.
* **[Groups](https://protobuf.dev/programming-guides/encoding/#groups)** - Deprecated proto2 feature, which might be implemented in the future.
* **[Editions](https://protobuf.dev/editions/overview/)** - Newer protobuf syntax that is in our roadmap.

## Failing Tests

For a machine-readable list of all test failures, see [conformance/expected_failures.txt](../conformance/expected_failures.txt).

| Test Name | Failure Message | Description |
|-----------|-----------------|-------------|
| `Recommended.Proto2.JsonInput.FieldNameExtension.Validator` | JSON payload validation failed | |
| `Required.Proto3.ProtobufInput.UnknownOrdering.ProtobufOutput` | Unknown field mismatch | |
| `Required.Proto3.ProtobufInput.UnknownVarint.ProtobufOutput` | Output was not equivalent to reference message: Expect: \250\037\001, but got: | |
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

---

*Last Updated: 2025-08-30*
