# Conformance Test Status

This document tracks the current status of Protocol Buffers conformance tests for the mysql-protobuf library.

## Summary

> 2746 successes, 414 skipped, 22 expected failures, 0 unexpected failures.

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
| `Recommended.Proto2.JsonInput.FieldNameExtension.Validator` | JSON payload validation failed | |
| `Required.Proto3.ProtobufInput.UnknownOrdering.ProtobufOutput` | Unknown field mismatch | |
| `Required.Proto3.ProtobufInput.UnknownVarint.ProtobufOutput` | Output was not equivalent to reference message: Expect: \250\037\001, but got: | |
| `Required.Proto2.ProtobufInput.UnknownOrdering.ProtobufOutput` | Unknown field mismatch | |
| `Required.Proto2.ProtobufInput.UnknownVarint.ProtobufOutput` | Output was not equivalent to reference message: Expect: \250\037\001, but got: | |
| `Recommended.Proto2.ProtobufInput.ValidMessageSetEncoding.SubmessageEncoding.NotUnknown.ProtobufOutput` | Output was not equivalent to reference message: deleted: message_set_correct | |
| `Recommended.Proto2.ProtobufInput.ValidMessageSetEncoding.SubmessageEncoding.ProtobufOutput` | Output was not equivalent to reference message: deleted: message_set_correct | |
| `Required.Proto2.ProtobufInput.MessageSetEncoding.UnknownExtension.ProtobufOutput` | Output was not equivalent to reference message: deleted: message_set_correct.4135300[0] | |
| `Required.Proto2.ProtobufInput.ValidMessageSetEncoding.OutOfOrderGroupsEntries.ProtobufOutput` | Output was not equivalent to reference message: deleted: message_set_correct | |
| `Required.Proto2.ProtobufInput.ValidMessageSetEncoding.ProtobufOutput` | Output was not equivalent to reference message: deleted: message_set_correct | |

## Tests Not Run

Conformance tests for Editions are not run unless `--maximum_edition=` is specified. Since we don't yet support Editions, we currently skip all Editions-related tests.

* **[Editions](https://protobuf.dev/editions/overview/)** - Newer protobuf syntax that is in our roadmap.

---

*Last Updated: 2025-09-11*
