package main

import (
	"fmt"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
)

func TestProtocGenOneofField(t *testing.T) {
	protoContent := dedent.Pipe(`
		|syntax = "proto3";
		|message Test {
		|    oneof choice {
		|        double double_field = 1;
		|        float float_field = 2;
		|        int32 int32_field = 3;
		|        int64 int64_field = 4;
		|        uint32 uint32_field = 5;
		|        uint64 uint64_field = 6;
		|        sint32 sint32_field = 7;
		|        sint64 sint64_field = 8;
		|        fixed32 fixed32_field = 9;
		|        fixed64 fixed64_field = 10;
		|        sfixed32 sfixed32_field = 11;
		|        sfixed64 sfixed64_field = 12;
		|        bool bool_field = 13;
		|        string string_field = 14;
		|        bytes bytes_field = 15;
		|        Status enum_field = 16;
		|        Nested message_field = 17;
		|    }
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

	// Test oneof mutual exclusion with all protobuf types
	// Each set operation should clear any previously set field
	t.Run("oneof_fields_internal_format", func(t *testing.T) {
		// Test double field in oneof (binary64 format)
		RunTestThatExpression(t, "test_set_double_field(test_new(), 3.141592653589793)").IsEqualToJsonString(`{"1": "binary64:0x400921fb54442d18"}`)

		// Test float field in oneof (binary32 format) - should clear double
		RunTestThatExpression(t, "test_set_float_field(test_set_double_field(test_new(), 3.14), 1.0)").IsEqualToJsonString(`{"2": "binary32:0x3f800000"}`)

		// Test int32 field in oneof - should clear float
		RunTestThatExpression(t, "test_set_int32_field(test_set_float_field(test_new(), 1.0), 42)").IsEqualToJsonString(`{"3": 42}`)

		// Test int64 field in oneof - should clear int32
		RunTestThatExpression(t, "test_set_int64_field(test_set_int32_field(test_new(), 42), 9223372036854775807)").IsEqualToJsonString(`{"4": 9223372036854775807}`)

		// Test uint32 field in oneof - should clear int64
		RunTestThatExpression(t, "test_set_uint32_field(test_set_int64_field(test_new(), 100), 4294967295)").IsEqualToJsonString(`{"5": 4294967295}`)

		// Test uint64 field in oneof - should clear uint32
		RunTestThatExpression(t, "test_set_uint64_field(test_set_uint32_field(test_new(), 100), 18446744073709551615)").IsEqualToJsonString(`{"6": 18446744073709551615}`)

		// Test sint32 field in oneof - should clear uint64
		RunTestThatExpression(t, "test_set_sint32_field(test_set_uint64_field(test_new(), 100), -1)").IsEqualToJsonString(`{"7": -1}`)

		// Test sint64 field in oneof - should clear sint32
		RunTestThatExpression(t, "test_set_sint64_field(test_set_sint32_field(test_new(), -1), -9223372036854775808)").IsEqualToJsonString(`{"8": -9223372036854775808}`)

		// Test fixed32 field in oneof - should clear sint64
		RunTestThatExpression(t, "test_set_fixed32_field(test_set_sint64_field(test_new(), -1), 4294967295)").IsEqualToJsonString(`{"9": 4294967295}`)

		// Test fixed64 field in oneof - should clear fixed32
		RunTestThatExpression(t, "test_set_fixed64_field(test_set_fixed32_field(test_new(), 100), 18446744073709551615)").IsEqualToJsonString(`{"10": 18446744073709551615}`)

		// Test sfixed32 field in oneof - should clear fixed64
		RunTestThatExpression(t, "test_set_sfixed32_field(test_set_fixed64_field(test_new(), 100), -2147483648)").IsEqualToJsonString(`{"11": -2147483648}`)

		// Test sfixed64 field in oneof - should clear sfixed32
		RunTestThatExpression(t, "test_set_sfixed64_field(test_set_sfixed32_field(test_new(), -1), -9223372036854775808)").IsEqualToJsonString(`{"12": -9223372036854775808}`)

		// Test bool field in oneof - should clear sfixed64
		RunTestThatExpression(t, "test_set_bool_field(test_set_sfixed64_field(test_new(), -1), TRUE)").IsEqualToJsonString(`{"13": true}`)

		// Test string field in oneof - should clear bool
		RunTestThatExpression(t, "test_set_string_field(test_set_bool_field(test_new(), TRUE), 'hello world')").IsEqualToJsonString(`{"14": "hello world"}`)

		// Test bytes field in oneof - should clear string
		RunTestThatExpression(t, "test_set_bytes_field(test_set_string_field(test_new(), 'hello'), ?)", []byte("world")).IsEqualToJsonString(`{"15": "d29ybGQ="}`)

		// Test enum field in oneof - should clear bytes
		RunTestThatExpression(t, "test_set_enum_field(test_set_bytes_field(test_new(), ?), 1)", []byte("test")).IsEqualToJsonString(`{"16": 1}`)

		// Test message field in oneof - should clear enum
		nestedObj := "nested_set_value(nested_set_name(nested_new(), 'test'), 42)"
		RunTestThatExpression(t, fmt.Sprintf("test_set_message_field(test_set_enum_field(test_new(), 1), %s)", nestedObj)).IsEqualToJsonString(`{"17": {"1": "test", "2": 42}}`)

		// Test that default values in oneof are NOT omitted (oneof has presence semantics)
		RunTestThatExpression(t, "test_set_double_field(test_new(), 0.0)").IsEqualToJsonString(`{"1": "binary64:0x0000000000000000"}`) // Zero double stored
		RunTestThatExpression(t, "test_set_int32_field(test_new(), 0)").IsEqualToJsonString(`{"3": 0}`)                                // Zero int32 stored
		RunTestThatExpression(t, "test_set_bool_field(test_new(), FALSE)").IsEqualToJsonString(`{"13": false}`)                        // False bool stored
		RunTestThatExpression(t, "test_set_string_field(test_new(), '')").IsEqualToJsonString(`{"14": ""}`)                            // Empty string stored
	})

	// Test double field in oneof
	t.Run("double_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "test_set_double_field(?, 3.141592653589793)", `{}`).IsEqualToJsonString(`{"1": "binary64:0x400921fb54442d18"}`)
		RunTestThatExpression(t, "test_set_double_field(?, 0.0)", `{}`).IsEqualToJsonString(`{"1": "binary64:0x0000000000000000"}`)

		// Test getters convert from internal representation back to actual values
		RunTestThatExpression(t, "test_get_double_field(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToDouble(3.141592653589793)
		RunTestThatExpression(t, "test_get_double_field(?)", `{"1": "binary64:0x0000000000000000"}`).IsEqualToDouble(0.0)
		RunTestThatExpression(t, "test_get_double_field(?)", `{}`).IsEqualToDouble(0.0) // Default when absent

		// Test has methods for oneof fields
		RunTestThatExpression(t, "test_has_double_field(?)", `{}`).IsFalse()                                  // Unset field not present
		RunTestThatExpression(t, "test_has_double_field(?)", `{"1": "binary64:0x0000000000000000"}`).IsTrue() // Set field present (even default value)

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{}`).IsNull() // No field set
		RunTestThatExpression(t, "test_which_choice(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToInt(1)

		// Test clear field methods
		RunTestThatExpression(t, "test_clear_double_field(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToJsonString(`{}`)

		// Test clear oneof group
		RunTestThatExpression(t, "test_clear_choice(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToJsonString(`{}`)
	})

	// Test float field in oneof
	t.Run("float_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_float_field(?, 3.14)", `{}`).IsEqualToJsonString(`{"2": "binary32:0x4048f5c3"}`)
		RunTestThatExpression(t, "test_set_float_field(?, 0.0)", `{}`).IsEqualToJsonString(`{"2": "binary32:0x00000000"}`)

		RunTestThatExpression(t, "test_get_float_field(?)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToFloat(3.14)
		RunTestThatExpression(t, "test_get_float_field(?)", `{"2": "binary32:0x00000000"}`).IsEqualToFloat(0.0)
		RunTestThatExpression(t, "test_get_float_field(?)", `{}`).IsEqualToFloat(0.0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_float_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_float_field(?)", `{"2": "binary32:0x00000000"}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToInt(2)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_float_field(?)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToJsonString(`{}`)
	})

	// Test integer fields
	t.Run("int32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_int32_field(?, 42)", `{}`).IsEqualToJsonString(`{"3": 42}`)
		RunTestThatExpression(t, "test_set_int32_field(?, 0)", `{}`).IsEqualToJsonString(`{"3": 0}`) // Zero stored (oneof presence semantics)
		RunTestThatExpression(t, "test_get_int32_field(?)", `{"3": 42}`).IsEqualToInt(42)
		RunTestThatExpression(t, "test_get_int32_field(?)", `{"3": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_int32_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_int32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_int32_field(?)", `{"3": 0}`).IsTrue() // Even zero value has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"3": 42}`).IsEqualToInt(3)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_int32_field(?)", `{"3": 42}`).IsEqualToJsonString(`{}`)
	})

	t.Run("int64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_int64_field(?, 9223372036854775807)", `{}`).IsEqualToJsonString(`{"4": 9223372036854775807}`)
		RunTestThatExpression(t, "test_set_int64_field(?, 0)", `{}`).IsEqualToJsonString(`{"4": 0}`)
		RunTestThatExpression(t, "test_get_int64_field(?)", `{"4": 9223372036854775807}`).IsEqualToInt(9223372036854775807)
		RunTestThatExpression(t, "test_get_int64_field(?)", `{"4": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_int64_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_int64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_int64_field(?)", `{"4": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"4": 9223372036854775807}`).IsEqualToInt(4)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_int64_field(?)", `{"4": 9223372036854775807}`).IsEqualToJsonString(`{}`)
	})

	t.Run("uint32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_uint32_field(?, 4294967295)", `{}`).IsEqualToJsonString(`{"5": 4294967295}`)
		RunTestThatExpression(t, "test_set_uint32_field(?, 0)", `{}`).IsEqualToJsonString(`{"5": 0}`)
		RunTestThatExpression(t, "test_get_uint32_field(?)", `{"5": 4294967295}`).IsEqualToUint(4294967295)
		RunTestThatExpression(t, "test_get_uint32_field(?)", `{"5": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "test_get_uint32_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_uint32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_uint32_field(?)", `{"5": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"5": 4294967295}`).IsEqualToInt(5)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_uint32_field(?)", `{"5": 4294967295}`).IsEqualToJsonString(`{}`)
	})

	t.Run("uint64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_uint64_field(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"6": 18446744073709551615}`)
		RunTestThatExpression(t, "test_set_uint64_field(?, 0)", `{}`).IsEqualToJsonString(`{"6": 0}`)
		RunTestThatExpression(t, "test_get_uint64_field(?)", `{"6": 18446744073709551615}`).IsEqualToUint(18446744073709551615)
		RunTestThatExpression(t, "test_get_uint64_field(?)", `{"6": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "test_get_uint64_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_uint64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_uint64_field(?)", `{"6": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"6": 18446744073709551615}`).IsEqualToInt(6)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_uint64_field(?)", `{"6": 18446744073709551615}`).IsEqualToJsonString(`{}`)
	})

	t.Run("sint32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_sint32_field(?, -1)", `{}`).IsEqualToJsonString(`{"7": -1}`)
		RunTestThatExpression(t, "test_set_sint32_field(?, 0)", `{}`).IsEqualToJsonString(`{"7": 0}`)
		RunTestThatExpression(t, "test_get_sint32_field(?)", `{"7": -1}`).IsEqualToInt(-1)
		RunTestThatExpression(t, "test_get_sint32_field(?)", `{"7": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_sint32_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_sint32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_sint32_field(?)", `{"7": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"7": -1}`).IsEqualToInt(7)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_sint32_field(?)", `{"7": -1}`).IsEqualToJsonString(`{}`)
	})

	t.Run("sint64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_sint64_field(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"8": -9223372036854775808}`)
		RunTestThatExpression(t, "test_set_sint64_field(?, 0)", `{}`).IsEqualToJsonString(`{"8": 0}`)
		RunTestThatExpression(t, "test_get_sint64_field(?)", `{"8": -9223372036854775808}`).IsEqualToInt(-9223372036854775808)
		RunTestThatExpression(t, "test_get_sint64_field(?)", `{"8": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_sint64_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_sint64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_sint64_field(?)", `{"8": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"8": -9223372036854775808}`).IsEqualToInt(8)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_sint64_field(?)", `{"8": -9223372036854775808}`).IsEqualToJsonString(`{}`)
	})

	t.Run("fixed32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_fixed32_field(?, 4294967295)", `{}`).IsEqualToJsonString(`{"9": 4294967295}`)
		RunTestThatExpression(t, "test_set_fixed32_field(?, 0)", `{}`).IsEqualToJsonString(`{"9": 0}`)
		RunTestThatExpression(t, "test_get_fixed32_field(?)", `{"9": 4294967295}`).IsEqualToUint(4294967295)
		RunTestThatExpression(t, "test_get_fixed32_field(?)", `{"9": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "test_get_fixed32_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_fixed32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_fixed32_field(?)", `{"9": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"9": 4294967295}`).IsEqualToInt(9)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_fixed32_field(?)", `{"9": 4294967295}`).IsEqualToJsonString(`{}`)
	})

	t.Run("fixed64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_fixed64_field(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"10": 18446744073709551615}`)
		RunTestThatExpression(t, "test_set_fixed64_field(?, 0)", `{}`).IsEqualToJsonString(`{"10": 0}`)
		RunTestThatExpression(t, "test_get_fixed64_field(?)", `{"10": 18446744073709551615}`).IsEqualToUint(18446744073709551615)
		RunTestThatExpression(t, "test_get_fixed64_field(?)", `{"10": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "test_get_fixed64_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_fixed64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_fixed64_field(?)", `{"10": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"10": 18446744073709551615}`).IsEqualToInt(10)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_fixed64_field(?)", `{"10": 18446744073709551615}`).IsEqualToJsonString(`{}`)
	})

	t.Run("sfixed32_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_sfixed32_field(?, -2147483648)", `{}`).IsEqualToJsonString(`{"11": -2147483648}`)
		RunTestThatExpression(t, "test_set_sfixed32_field(?, 0)", `{}`).IsEqualToJsonString(`{"11": 0}`)
		RunTestThatExpression(t, "test_get_sfixed32_field(?)", `{"11": -2147483648}`).IsEqualToInt(-2147483648)
		RunTestThatExpression(t, "test_get_sfixed32_field(?)", `{"11": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_sfixed32_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_sfixed32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_sfixed32_field(?)", `{"11": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"11": -2147483648}`).IsEqualToInt(11)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_sfixed32_field(?)", `{"11": -2147483648}`).IsEqualToJsonString(`{}`)
	})

	t.Run("sfixed64_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_sfixed64_field(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"12": -9223372036854775808}`)
		RunTestThatExpression(t, "test_set_sfixed64_field(?, 0)", `{}`).IsEqualToJsonString(`{"12": 0}`)
		RunTestThatExpression(t, "test_get_sfixed64_field(?)", `{"12": -9223372036854775808}`).IsEqualToInt(-9223372036854775808)
		RunTestThatExpression(t, "test_get_sfixed64_field(?)", `{"12": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_sfixed64_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_sfixed64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_sfixed64_field(?)", `{"12": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"12": -9223372036854775808}`).IsEqualToInt(12)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_sfixed64_field(?)", `{"12": -9223372036854775808}`).IsEqualToJsonString(`{}`)
	})

	// Test bool field in oneof
	t.Run("bool_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_bool_field(?, TRUE)", `{}`).IsEqualToJsonString(`{"13": true}`)
		RunTestThatExpression(t, "test_set_bool_field(?, FALSE)", `{}`).IsEqualToJsonString(`{"13": false}`) // False stored (oneof presence semantics)

		RunTestThatExpression(t, "test_get_bool_field(?)", `{"13": true}`).IsEqualToBool(true)
		RunTestThatExpression(t, "test_get_bool_field(?)", `{"13": false}`).IsEqualToBool(false)
		RunTestThatExpression(t, "test_get_bool_field(?)", `{}`).IsEqualToBool(false) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_bool_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_bool_field(?)", `{"13": false}`).IsTrue() // Even false value has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"13": true}`).IsEqualToInt(13)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_bool_field(?)", `{"13": true}`).IsEqualToJsonString(`{}`)
	})

	// Test string field in oneof
	t.Run("string_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_string_field(?, 'hello world')", `{}`).IsEqualToJsonString(`{"14": "hello world"}`)
		RunTestThatExpression(t, "test_set_string_field(?, '')", `{}`).IsEqualToJsonString(`{"14": ""}`) // Empty string stored (oneof presence semantics)

		RunTestThatExpression(t, "test_get_string_field(?)", `{"14": "hello world"}`).IsEqualToString("hello world")
		RunTestThatExpression(t, "test_get_string_field(?)", `{"14": ""}`).IsEqualToString("")
		RunTestThatExpression(t, "test_get_string_field(?)", `{}`).IsEqualToString("") // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_string_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_string_field(?)", `{"14": ""}`).IsTrue() // Even empty string has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"14": "hello world"}`).IsEqualToInt(14)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_string_field(?)", `{"14": "hello world"}`).IsEqualToJsonString(`{}`)
	})

	// Test bytes field in oneof
	t.Run("bytes_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_bytes_field(?, ?)", `{}`, []byte("hello")).IsEqualToJsonString(`{"15": "aGVsbG8="}`)
		RunTestThatExpression(t, "test_set_bytes_field(?, ?)", `{}`, []byte{}).IsEqualToJsonString(`{"15": ""}`) // Empty bytes stored (oneof presence semantics)

		RunTestThatExpression(t, "test_get_bytes_field(?)", `{"15": "aGVsbG8="}`).IsEqualToBytes([]byte("hello"))
		RunTestThatExpression(t, "test_get_bytes_field(?)", `{"15": ""}`).IsEqualToBytes([]byte{})
		RunTestThatExpression(t, "test_get_bytes_field(?)", `{}`).IsEqualToBytes([]byte{}) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_bytes_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_bytes_field(?)", `{"15": ""}`).IsTrue() // Even empty bytes has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"15": "aGVsbG8="}`).IsEqualToInt(15)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_bytes_field(?)", `{"15": "aGVsbG8="}`).IsEqualToJsonString(`{}`)
	})

	// Test enum field in oneof
	t.Run("enum_field", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_enum_field(?, 1)", `{}`).IsEqualToJsonString(`{"16": 1}`)
		RunTestThatExpression(t, "test_set_enum_field(?, 0)", `{}`).IsEqualToJsonString(`{"16": 0}`) // Zero enum stored (oneof presence semantics)

		RunTestThatExpression(t, "test_get_enum_field(?)", `{"16": 1}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_get_enum_field(?)", `{"16": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_get_enum_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_enum_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "test_has_enum_field(?)", `{"16": 0}`).IsTrue() // Even default enum value has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"16": 1}`).IsEqualToInt(16)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_enum_field(?)", `{"16": 1}`).IsEqualToJsonString(`{}`)
	})

	// Test message field in oneof
	t.Run("message_field", func(t *testing.T) {
		// Test setters
		RunTestThatExpression(t, "test_set_message_field(?, ?)", `{}`, `{"1": "test", "2": 42}`).IsEqualToJsonString(`{"17": {"1": "test", "2": 42}}`)
		RunTestThatExpression(t, "test_set_message_field(?, ?)", `{}`, `{}`).IsEqualToJsonString(`{"17": {}}`) // Empty message stored (oneof presence semantics)

		// Test getters
		RunTestThatExpression(t, "test_get_message_field(?)", `{"17": {"1": "test", "2": 42}}`).IsEqualToJsonString(`{"1": "test", "2": 42}`)
		RunTestThatExpression(t, "test_get_message_field(?)", `{"17": {}}`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, "test_get_message_field(?)", `{}`).IsEqualToJsonString(`{}`) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "test_has_message_field(?)", `{}`).IsFalse()        // Unset field not present
		RunTestThatExpression(t, "test_has_message_field(?)", `{"17": {}}`).IsTrue() // Set field present (even empty message) in oneof

		// Test which oneof method
		RunTestThatExpression(t, "test_which_choice(?)", `{"17": {"1": "test", "2": 42}}`).IsEqualToInt(17)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_message_field(?)", `{"17": {"1": "test", "2": 42}}`).IsEqualToJsonString(`{}`)
	})

	// Test oneof mutual exclusion behavior
	t.Run("mutual_exclusion", func(t *testing.T) {
		// Setting different fields should clear previous ones
		RunTestThatExpression(t, "test_set_string_field(test_set_int32_field(test_new(), 42), 'hello')").IsEqualToJsonString(`{"14": "hello"}`)
		RunTestThatExpression(t, "test_which_choice(test_set_string_field(test_set_int32_field(test_new(), 42), 'hello'))").IsEqualToInt(14)

		// Previous field should return default value after being cleared
		RunTestThatExpression(t, "test_get_int32_field(test_set_string_field(test_set_int32_field(test_new(), 42), 'hello'))").IsEqualToInt(0)
		RunTestThatExpression(t, "test_has_int32_field(test_set_string_field(test_set_int32_field(test_new(), 42), 'hello'))").IsFalse()

		// Test complex mutual exclusion chain
		expr := "test_set_bool_field(test_set_enum_field(test_set_string_field(test_set_int32_field(test_new(), 42), 'hello'), 1), TRUE)"
		RunTestThatExpression(t, expr).IsEqualToJsonString(`{"13": true}`)
		RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", expr)).IsEqualToInt(13)
		RunTestThatExpression(t, fmt.Sprintf("test_get_int32_field(%s)", expr)).IsEqualToInt(0)      // cleared
		RunTestThatExpression(t, fmt.Sprintf("test_get_string_field(%s)", expr)).IsEqualToString("") // cleared
		RunTestThatExpression(t, fmt.Sprintf("test_get_enum_field(%s)", expr)).IsEqualToInt(0)       // cleared
	})

	// Test clearing entire oneof group
	t.Run("clear_oneof_group", func(t *testing.T) {
		// Set a field, then clear the entire oneof group
		expr := "test_clear_choice(test_set_string_field(test_new(), 'hello'))"
		RunTestThatExpression(t, expr).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", expr)).IsNull()
		RunTestThatExpression(t, fmt.Sprintf("test_get_string_field(%s)", expr)).IsEqualToString("") // returns default
		RunTestThatExpression(t, fmt.Sprintf("test_has_string_field(%s)", expr)).IsFalse()

		// Test clearing after setting multiple fields (mutual exclusion)
		expr2 := "test_clear_choice(test_set_int32_field(test_set_string_field(test_new(), 'hello'), 42))"
		RunTestThatExpression(t, expr2).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", expr2)).IsNull()
		RunTestThatExpression(t, fmt.Sprintf("test_get_int32_field(%s)", expr2)).IsEqualToInt(0) // returns default
		RunTestThatExpression(t, fmt.Sprintf("test_has_int32_field(%s)", expr2)).IsFalse()
	})
}
