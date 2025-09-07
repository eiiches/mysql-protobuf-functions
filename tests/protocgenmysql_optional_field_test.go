package main

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
)

func TestProtocGenOptionalField(t *testing.T) {
	protoContent := dedent.Pipe(`
		|syntax = "proto3";
		|message Test {
		|    optional double optional_double_field = 1;
		|    optional float optional_float_field = 2;
		|    optional int32 optional_int32_field = 3;
		|    optional int64 optional_int64_field = 4;
		|    optional uint32 optional_uint32_field = 5;
		|    optional uint64 optional_uint64_field = 6;
		|    optional sint32 optional_sint32_field = 7;
		|    optional sint64 optional_sint64_field = 8;
		|    optional fixed32 optional_fixed32_field = 9;
		|    optional fixed64 optional_fixed64_field = 10;
		|    optional sfixed32 optional_sfixed32_field = 11;
		|    optional sfixed64 optional_sfixed64_field = 12;
		|    optional bool optional_bool_field = 13;
		|    optional string optional_string_field = 14;
		|    optional bytes optional_bytes_field = 15;
		|    optional Status optional_enum_field = 16;
		|    optional Nested optional_message_field = 17;
		|}
		|message Nested {
		|    string name = 1;
		|    int32 value = 2;
		|}
		|enum Status {
		|    STATUS_UNSPECIFIED = 0;
		|    STATUS_ACTIVE = 1;
		|    STATUS_INACTIVE = 2;
		|}
	`)
	schemaName := "test_schema"
	generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

	// Test double optional field (IEEE 754 binary64 format)
	t.Run("optional_double_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "test_set_optional_double_field(?, 3.141592653589793)", `{}`).IsEqualToJsonString(`{"1": "binary64:0x400921fb54442d18"}`)
		// Zero value stored (presence semantics)
		RunTestThatExpression(t, "test_set_optional_double_field(?, 0.0)", `{}`).IsEqualToJsonString(`{"1": "binary64:0x0000000000000000"}`)

		// Test getters convert from internal representation back to actual values
		RunTestThatExpression(t, "test_get_optional_double_field(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToDouble(3.141592653589793)
		RunTestThatExpression(t, "test_get_optional_double_field(?)", `{"1": "binary64:0x0000000000000000"}`).IsEqualToDouble(0.0)
		RunTestThatExpression(t, "test_get_optional_double_field(?)", `{}`).IsEqualToDouble(0.0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_double_field(?)", `{}`).IsFalse()                                  // Unset field not present
		RunTestThatExpression(t, "test_has_optional_double_field(?)", `{"1": "binary64:0x0000000000000000"}`).IsTrue() // Set field present (even default value)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_double_field(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_double_field__or(?, 99.99)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToDouble(3.141592653589793) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_double_field__or(?, 99.99)", `{}`).IsEqualToDouble(99.99)                                               // Field absent, return default
	})

	// Test float optional field (IEEE 754 binary32 format)
	t.Run("optional_float_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_float_field(?, 3.14)", `{}`).IsEqualToJsonString(`{"2": "binary32:0x4048f5c3"}`)
		RunTestThatExpression(t, "test_set_optional_float_field(?, 0.0)", `{}`).IsEqualToJsonString(`{"2": "binary32:0x00000000"}`)

		RunTestThatExpression(t, "test_get_optional_float_field(?)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToFloat(3.14)
		RunTestThatExpression(t, "test_get_optional_float_field(?)", `{"2": "binary32:0x00000000"}`).IsEqualToFloat(0.0)
		RunTestThatExpression(t, "test_get_optional_float_field(?)", `{}`).IsEqualToFloat(0.0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_float_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_float_field(?)", `{"2": "binary32:0x00000000"}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_float_field(?)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_float_field__or(?, 88.88)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToFloat(3.14) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_float_field__or(?, 88.88)", `{}`).IsEqualToFloat(88.88)                          // Field absent, return default
	})

	// Test remaining integer types
	t.Run("optional_int32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_int32_field(?, 42)", `{}`).IsEqualToJsonString(`{"3": 42}`)
		RunTestThatExpression(t, "test_set_optional_int32_field(?, 0)", `{}`).IsEqualToJsonString(`{"3": 0}`) // Zero stored
		RunTestThatExpression(t, "test_get_optional_int32_field(?)", `{"3": 42}`).IsEqualToInt(42)
		RunTestThatExpression(t, "test_get_optional_int32_field(?)", `{"3": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_optional_int32_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_int32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_int32_field(?)", `{"3": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_int32_field(?)", `{"3": 42}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_int32_field__or(?, 999)", `{"3": 42}`).IsEqualToInt(42) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_int32_field__or(?, 999)", `{}`).IsEqualToInt(999)       // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_int32_field__or(?, 999)", `{"3": 0}`).IsEqualToInt(0)   // Zero is still present
	})

	t.Run("optional_int64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_int64_field(?, 9223372036854775807)", `{}`).IsEqualToJsonString(`{"4": 9223372036854775807}`)
		RunTestThatExpression(t, "test_set_optional_int64_field(?, 0)", `{}`).IsEqualToJsonString(`{"4": 0}`)
		RunTestThatExpression(t, "test_get_optional_int64_field(?)", `{"4": 9223372036854775807}`).IsEqualToInt(9223372036854775807)
		RunTestThatExpression(t, "test_get_optional_int64_field(?)", `{"4": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_optional_int64_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_int64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_int64_field(?)", `{"4": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_int64_field(?)", `{"4": 9223372036854775807}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_int64_field__or(?, 999)", `{"4": 9223372036854775807}`).IsEqualToInt(9223372036854775807) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_int64_field__or(?, 999)", `{}`).IsEqualToInt(999)                                         // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_int64_field__or(?, 999)", `{"4": 0}`).IsEqualToInt(0)                                     // Zero is still present
	})

	t.Run("optional_uint32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_uint32_field(?, 4294967295)", `{}`).IsEqualToJsonString(`{"5": 4294967295}`)
		RunTestThatExpression(t, "test_set_optional_uint32_field(?, 0)", `{}`).IsEqualToJsonString(`{"5": 0}`)
		RunTestThatExpression(t, "test_get_optional_uint32_field(?)", `{"5": 4294967295}`).IsEqualToUint(4294967295)
		RunTestThatExpression(t, "test_get_optional_uint32_field(?)", `{"5": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "test_get_optional_uint32_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_uint32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_uint32_field(?)", `{"5": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_uint32_field(?)", `{"5": 4294967295}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_uint32_field__or(?, 999)", `{"5": 4294967295}`).IsEqualToUint(4294967295) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_uint32_field__or(?, 999)", `{}`).IsEqualToUint(999)                       // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_uint32_field__or(?, 999)", `{"5": 0}`).IsEqualToUint(0)                   // Zero is still present
	})

	t.Run("optional_uint64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_uint64_field(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"6": 18446744073709551615}`)
		RunTestThatExpression(t, "test_set_optional_uint64_field(?, 0)", `{}`).IsEqualToJsonString(`{"6": 0}`)
		RunTestThatExpression(t, "test_get_optional_uint64_field(?)", `{"6": 18446744073709551615}`).IsEqualToUint(18446744073709551615)
		RunTestThatExpression(t, "test_get_optional_uint64_field(?)", `{"6": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "test_get_optional_uint64_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_uint64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_uint64_field(?)", `{"6": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_uint64_field(?)", `{"6": 18446744073709551615}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_uint64_field__or(?, 999)", `{"6": 18446744073709551615}`).IsEqualToUint(18446744073709551615) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_uint64_field__or(?, 999)", `{}`).IsEqualToUint(999)                                           // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_uint64_field__or(?, 999)", `{"6": 0}`).IsEqualToUint(0)                                       // Zero is still present
	})

	t.Run("optional_sint32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_sint32_field(?, -1)", `{}`).IsEqualToJsonString(`{"7": -1}`)
		RunTestThatExpression(t, "test_set_optional_sint32_field(?, 0)", `{}`).IsEqualToJsonString(`{"7": 0}`)
		RunTestThatExpression(t, "test_get_optional_sint32_field(?)", `{"7": -1}`).IsEqualToInt(-1)
		RunTestThatExpression(t, "test_get_optional_sint32_field(?)", `{"7": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_optional_sint32_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_sint32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_sint32_field(?)", `{"7": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_sint32_field(?)", `{"7": -1}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_sint32_field__or(?, 999)", `{"7": -1}`).IsEqualToInt(-1) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_sint32_field__or(?, 999)", `{}`).IsEqualToInt(999)       // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_sint32_field__or(?, 999)", `{"7": 0}`).IsEqualToInt(0)   // Zero is still present
	})

	t.Run("optional_sint64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_sint64_field(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"8": -9223372036854775808}`)
		RunTestThatExpression(t, "test_set_optional_sint64_field(?, 0)", `{}`).IsEqualToJsonString(`{"8": 0}`)
		RunTestThatExpression(t, "test_get_optional_sint64_field(?)", `{"8": -9223372036854775808}`).IsEqualToInt(-9223372036854775808)
		RunTestThatExpression(t, "test_get_optional_sint64_field(?)", `{"8": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_optional_sint64_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_sint64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_sint64_field(?)", `{"8": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_sint64_field(?)", `{"8": -9223372036854775808}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_sint64_field__or(?, 999)", `{"8": -9223372036854775808}`).IsEqualToInt(-9223372036854775808) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_sint64_field__or(?, 999)", `{}`).IsEqualToInt(999)                                           // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_sint64_field__or(?, 999)", `{"8": 0}`).IsEqualToInt(0)                                       // Zero is still present
	})

	t.Run("optional_fixed32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_fixed32_field(?, 4294967295)", `{}`).IsEqualToJsonString(`{"9": 4294967295}`)
		RunTestThatExpression(t, "test_set_optional_fixed32_field(?, 0)", `{}`).IsEqualToJsonString(`{"9": 0}`)
		RunTestThatExpression(t, "test_get_optional_fixed32_field(?)", `{"9": 4294967295}`).IsEqualToUint(4294967295)
		RunTestThatExpression(t, "test_get_optional_fixed32_field(?)", `{"9": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "test_get_optional_fixed32_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_fixed32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_fixed32_field(?)", `{"9": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_fixed32_field(?)", `{"9": 4294967295}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_fixed32_field__or(?, 999)", `{"9": 4294967295}`).IsEqualToUint(4294967295) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_fixed32_field__or(?, 999)", `{}`).IsEqualToUint(999)                       // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_fixed32_field__or(?, 999)", `{"9": 0}`).IsEqualToUint(0)                   // Zero is still present
	})

	t.Run("optional_fixed64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_fixed64_field(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"10": 18446744073709551615}`)
		RunTestThatExpression(t, "test_set_optional_fixed64_field(?, 0)", `{}`).IsEqualToJsonString(`{"10": 0}`)
		RunTestThatExpression(t, "test_get_optional_fixed64_field(?)", `{"10": 18446744073709551615}`).IsEqualToUint(18446744073709551615)
		RunTestThatExpression(t, "test_get_optional_fixed64_field(?)", `{"10": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "test_get_optional_fixed64_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_fixed64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_fixed64_field(?)", `{"10": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_fixed64_field(?)", `{"10": 18446744073709551615}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_fixed64_field__or(?, 999)", `{"10": 18446744073709551615}`).IsEqualToUint(18446744073709551615) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_fixed64_field__or(?, 999)", `{}`).IsEqualToUint(999)                                            // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_fixed64_field__or(?, 999)", `{"10": 0}`).IsEqualToUint(0)                                       // Zero is still present
	})

	t.Run("optional_sfixed32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_sfixed32_field(?, -2147483648)", `{}`).IsEqualToJsonString(`{"11": -2147483648}`)
		RunTestThatExpression(t, "test_set_optional_sfixed32_field(?, 0)", `{}`).IsEqualToJsonString(`{"11": 0}`)
		RunTestThatExpression(t, "test_get_optional_sfixed32_field(?)", `{"11": -2147483648}`).IsEqualToInt(-2147483648)
		RunTestThatExpression(t, "test_get_optional_sfixed32_field(?)", `{"11": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_optional_sfixed32_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_sfixed32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_sfixed32_field(?)", `{"11": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_sfixed32_field(?)", `{"11": -2147483648}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_sfixed32_field__or(?, 999)", `{"11": -2147483648}`).IsEqualToInt(-2147483648) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_sfixed32_field__or(?, 999)", `{}`).IsEqualToInt(999)                          // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_sfixed32_field__or(?, 999)", `{"11": 0}`).IsEqualToInt(0)                     // Zero is still present
	})

	t.Run("optional_sfixed64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_sfixed64_field(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"12": -9223372036854775808}`)
		RunTestThatExpression(t, "test_set_optional_sfixed64_field(?, 0)", `{}`).IsEqualToJsonString(`{"12": 0}`)
		RunTestThatExpression(t, "test_get_optional_sfixed64_field(?)", `{"12": -9223372036854775808}`).IsEqualToInt(-9223372036854775808)
		RunTestThatExpression(t, "test_get_optional_sfixed64_field(?)", `{"12": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_optional_sfixed64_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_sfixed64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_sfixed64_field(?)", `{"12": 0}`).IsTrue()

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_sfixed64_field(?)", `{"12": -9223372036854775808}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_sfixed64_field__or(?, 999)", `{"12": -9223372036854775808}`).IsEqualToInt(-9223372036854775808) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_sfixed64_field__or(?, 999)", `{}`).IsEqualToInt(999)                                            // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_sfixed64_field__or(?, 999)", `{"12": 0}`).IsEqualToInt(0)                                       // Zero is still present
	})

	// Test bool optional field
	t.Run("optional_bool_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_bool_field(?, TRUE)", `{}`).IsEqualToJsonString(`{"13": true}`)
		RunTestThatExpression(t, "test_set_optional_bool_field(?, FALSE)", `{}`).IsEqualToJsonString(`{"13": false}`) // False stored (presence semantics)

		RunTestThatExpression(t, "test_get_optional_bool_field(?)", `{"13": true}`).IsEqualToBool(true)
		RunTestThatExpression(t, "test_get_optional_bool_field(?)", `{"13": false}`).IsEqualToBool(false)
		RunTestThatExpression(t, "test_get_optional_bool_field(?)", `{}`).IsEqualToBool(false) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_bool_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_bool_field(?)", `{"13": false}`).IsTrue() // Even false value has presence

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_bool_field(?)", `{"13": true}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_bool_field__or(?, TRUE)", `{"13": false}`).IsEqualToBool(false) // Field present (false), return actual value
		RunTestThatExpression(t, "test_get_optional_bool_field__or(?, TRUE)", `{}`).IsEqualToBool(true)             // Field absent, return default (true)
		RunTestThatExpression(t, "test_get_optional_bool_field__or(?, FALSE)", `{"13": true}`).IsEqualToBool(true)  // Field present (true), return actual value
	})

	// Test string optional field
	t.Run("optional_string_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_string_field(?, 'hello world')", `{}`).IsEqualToJsonString(`{"14": "hello world"}`)
		RunTestThatExpression(t, "test_set_optional_string_field(?, '')", `{}`).IsEqualToJsonString(`{"14": ""}`) // Empty string stored

		RunTestThatExpression(t, "test_get_optional_string_field(?)", `{"14": "hello world"}`).IsEqualToString("hello world")
		RunTestThatExpression(t, "test_get_optional_string_field(?)", `{"14": ""}`).IsEqualToString("")
		RunTestThatExpression(t, "test_get_optional_string_field(?)", `{}`).IsEqualToString("") // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_string_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_string_field(?)", `{"14": ""}`).IsTrue() // Even empty string has presence

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_string_field(?)", `{"14": "hello world"}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_string_field__or(?, 'default text')", `{"14": "hello world"}`).IsEqualToString("hello world") // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_string_field__or(?, 'default text')", `{}`).IsEqualToString("default text")                   // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_string_field__or(?, 'default text')", `{"14": ""}`).IsEqualToString("")                       // Empty string is still present
	})

	// Test bytes optional field
	t.Run("optional_bytes_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_bytes_field(?, ?)", `{}`, []byte("hello")).IsEqualToJsonString(`{"15": "aGVsbG8="}`)
		RunTestThatExpression(t, "test_set_optional_bytes_field(?, ?)", `{}`, []byte{}).IsEqualToJsonString(`{"15": ""}`) // Empty bytes stored

		RunTestThatExpression(t, "test_get_optional_bytes_field(?)", `{"15": "aGVsbG8="}`).IsEqualToBytes([]byte("hello"))
		RunTestThatExpression(t, "test_get_optional_bytes_field(?)", `{"15": ""}`).IsEqualToBytes([]byte{})
		RunTestThatExpression(t, "test_get_optional_bytes_field(?)", `{}`).IsEqualToBytes([]byte{}) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_bytes_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_bytes_field(?)", `{"15": ""}`).IsTrue() // Even empty bytes has presence

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_bytes_field(?)", `{"15": "aGVsbG8="}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		defaultBytes := []byte("default")
		RunTestThatExpression(t, "test_get_optional_bytes_field__or(?, ?)", `{"15": "aGVsbG8="}`, defaultBytes).IsEqualToBytes([]byte("hello")) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_bytes_field__or(?, ?)", `{}`, defaultBytes).IsEqualToBytes(defaultBytes)                    // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_bytes_field__or(?, ?)", `{"15": ""}`, defaultBytes).IsEqualToBytes([]byte{})                // Empty bytes is still present
	})

	// Test enum optional field
	t.Run("optional_enum_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_optional_enum_field(?, 1)", `{}`).IsEqualToJsonString(`{"16": 1}`)
		RunTestThatExpression(t, "test_set_optional_enum_field(?, 0)", `{}`).IsEqualToJsonString(`{"16": 0}`) // Zero enum stored

		RunTestThatExpression(t, "test_get_optional_enum_field(?)", `{"16": 1}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_get_optional_enum_field(?)", `{"16": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_optional_enum_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_enum_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_optional_enum_field(?)", `{"16": 0}`).IsTrue() // Even default enum value has presence

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_enum_field(?)", `{"16": 1}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		RunTestThatExpression(t, "test_get_optional_enum_field__or(?, 999)", `{"16": 1}`).IsEqualToInt(1) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_enum_field__or(?, 999)", `{}`).IsEqualToInt(999)      // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_enum_field__or(?, 999)", `{"16": 0}`).IsEqualToInt(0) // Zero enum is still present

		// Test enum name getter (__as_name) methods
		RunTestThatExpression(t, "test_get_optional_enum_field__as_name(?)", `{"16": 0}`).IsEqualToString("STATUS_UNSPECIFIED") // Zero value name
		RunTestThatExpression(t, "test_get_optional_enum_field__as_name(?)", `{"16": 1}`).IsEqualToString("STATUS_ACTIVE")      // Active status name
		RunTestThatExpression(t, "test_get_optional_enum_field__as_name(?)", `{"16": 2}`).IsEqualToString("STATUS_INACTIVE")    // Inactive status name
		RunTestThatExpression(t, "test_get_optional_enum_field__as_name(?)", `{"16": 999}`).IsEqualToString("999")              // Unknown enum value returns number as string
		RunTestThatExpression(t, "test_get_optional_enum_field__as_name(?)", `{}`).IsEqualToString("STATUS_UNSPECIFIED")        // Default when absent

		// Test enum name nullable getter (__as_name_or) methods
		RunTestThatExpression(t, "test_get_optional_enum_field__as_name_or(?, 'DEFAULT_VALUE')", `{"16": 1}`).IsEqualToString("STATUS_ACTIVE")      // Field present, return actual name
		RunTestThatExpression(t, "test_get_optional_enum_field__as_name_or(?, 'DEFAULT_VALUE')", `{}`).IsEqualToString("DEFAULT_VALUE")             // Field absent, return custom default
		RunTestThatExpression(t, "test_get_optional_enum_field__as_name_or(?, 'DEFAULT_VALUE')", `{"16": 0}`).IsEqualToString("STATUS_UNSPECIFIED") // Zero enum is still present

		// Test enum name setter (__from_name) methods
		RunTestThatExpression(t, "test_set_optional_enum_field__from_name(?, 'STATUS_ACTIVE')", `{}`).IsEqualToJsonString(`{"16": 1}`)                       // Set from valid name
		RunTestThatExpression(t, "test_set_optional_enum_field__from_name(?, 'STATUS_INACTIVE')", `{}`).IsEqualToJsonString(`{"16": 2}`)                     // Set from different name
		RunTestThatExpression(t, "test_set_optional_enum_field__from_name(?, 'STATUS_UNSPECIFIED')", `{}`).IsEqualToJsonString(`{"16": 0}`)                  // Set from zero value name
		RunTestThatExpression(t, "test_set_optional_enum_field__from_name(?, 'INVALID_NAME')", `{}`).ToFailWithSignalException("45000", "Invalid enum name") // Invalid name should signal error
	})

	// Test message optional field
	t.Run("optional_message_field", func(t *testing.T) {
		// Test setters
		RunTestThatExpression(t, "test_set_optional_message_field(?, ?)", `{}`, `{"1": "test", "2": 42}`).IsEqualToJsonString(`{"17": {"1": "test", "2": 42}}`)
		RunTestThatExpression(t, "test_set_optional_message_field(?, ?)", `{}`, `{}`).IsEqualToJsonString(`{"17": {}}`) // Empty message stored

		// Test getters
		RunTestThatExpression(t, "test_get_optional_message_field(?)", `{"17": {"1": "test", "2": 42}}`).IsEqualToJsonString(`{"1": "test", "2": 42}`)
		RunTestThatExpression(t, "test_get_optional_message_field(?)", `{"17": {}}`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, "test_get_optional_message_field(?)", `{}`).IsEqualToJsonString(`{}`) // Default when absent

		// Test presence methods
		RunTestThatExpression(t, "test_has_optional_message_field(?)", `{}`).IsFalse()        // Unset field not present
		RunTestThatExpression(t, "test_has_optional_message_field(?)", `{"17": {}}`).IsTrue() // Set field present (even empty message)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_optional_message_field(?)", `{"17": {"1": "test", "2": 42}}`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods
		defaultMessage := `{"1": "default", "2": 999}`
		RunTestThatExpression(t, "test_get_optional_message_field__or(?, ?)", `{"17": {"1": "test", "2": 42}}`, defaultMessage).IsEqualToJsonString(`{"1": "test", "2": 42}`) // Field present, return actual value
		RunTestThatExpression(t, "test_get_optional_message_field__or(?, ?)", `{}`, defaultMessage).IsEqualToJsonString(defaultMessage)                                       // Field absent, return default
		RunTestThatExpression(t, "test_get_optional_message_field__or(?, ?)", `{"17": {}}`, defaultMessage).IsEqualToJsonString(`{}`)                                         // Empty message is still present
	})
}
