package main

import (
	"fmt"
	"testing"
)

func TestProtocGenOneofField(t *testing.T) {
	// Test all protobuf oneof field types using the pre-generated functions from protocgenmysql_proto3.pb.sql
	// The OneofFields message in protocgenmysql_proto3.proto matches the test schema

	// Test oneof mutual exclusion with all protobuf types
	// Each set operation should clear any previously set field
	t.Run("oneof_fields_internal_format", func(t *testing.T) {
		// Test double field in oneof (binary64 format)
		RunTestThatExpression(t, "pbt_oneof_fields_set_double_field(pbt_oneof_fields_new(), 3.141592653589793)").IsEqualToJsonString(`{"1": "binary64:0x400921fb54442d18"}`)

		// Test float field in oneof (binary32 format) - should clear double
		RunTestThatExpression(t, "pbt_oneof_fields_set_float_field(pbt_oneof_fields_set_double_field(pbt_oneof_fields_new(), 3.14), 1.0)").IsEqualToJsonString(`{"2": "binary32:0x3f800000"}`)

		// Test int32 field in oneof - should clear float
		RunTestThatExpression(t, "pbt_oneof_fields_set_int32_field(pbt_oneof_fields_set_float_field(pbt_oneof_fields_new(), 1.0), 42)").IsEqualToJsonString(`{"3": 42}`)

		// Test int64 field in oneof - should clear int32
		RunTestThatExpression(t, "pbt_oneof_fields_set_int64_field(pbt_oneof_fields_set_int32_field(pbt_oneof_fields_new(), 42), 9223372036854775807)").IsEqualToJsonString(`{"4": 9223372036854775807}`)

		// Test uint32 field in oneof - should clear int64
		RunTestThatExpression(t, "pbt_oneof_fields_set_uint32_field(pbt_oneof_fields_set_int64_field(pbt_oneof_fields_new(), 100), 4294967295)").IsEqualToJsonString(`{"5": 4294967295}`)

		// Test uint64 field in oneof - should clear uint32
		RunTestThatExpression(t, "pbt_oneof_fields_set_uint64_field(pbt_oneof_fields_set_uint32_field(pbt_oneof_fields_new(), 100), 18446744073709551615)").IsEqualToJsonString(`{"6": 18446744073709551615}`)

		// Test sint32 field in oneof - should clear uint64
		RunTestThatExpression(t, "pbt_oneof_fields_set_sint32_field(pbt_oneof_fields_set_uint64_field(pbt_oneof_fields_new(), 100), -1)").IsEqualToJsonString(`{"7": -1}`)

		// Test sint64 field in oneof - should clear sint32
		RunTestThatExpression(t, "pbt_oneof_fields_set_sint64_field(pbt_oneof_fields_set_sint32_field(pbt_oneof_fields_new(), -1), -9223372036854775808)").IsEqualToJsonString(`{"8": -9223372036854775808}`)

		// Test fixed32 field in oneof - should clear sint64
		RunTestThatExpression(t, "pbt_oneof_fields_set_fixed32_field(pbt_oneof_fields_set_sint64_field(pbt_oneof_fields_new(), -1), 4294967295)").IsEqualToJsonString(`{"9": 4294967295}`)

		// Test fixed64 field in oneof - should clear fixed32
		RunTestThatExpression(t, "pbt_oneof_fields_set_fixed64_field(pbt_oneof_fields_set_fixed32_field(pbt_oneof_fields_new(), 100), 18446744073709551615)").IsEqualToJsonString(`{"10": 18446744073709551615}`)

		// Test sfixed32 field in oneof - should clear fixed64
		RunTestThatExpression(t, "pbt_oneof_fields_set_sfixed32_field(pbt_oneof_fields_set_fixed64_field(pbt_oneof_fields_new(), 100), -2147483648)").IsEqualToJsonString(`{"11": -2147483648}`)

		// Test sfixed64 field in oneof - should clear sfixed32
		RunTestThatExpression(t, "pbt_oneof_fields_set_sfixed64_field(pbt_oneof_fields_set_sfixed32_field(pbt_oneof_fields_new(), -1), -9223372036854775808)").IsEqualToJsonString(`{"12": -9223372036854775808}`)

		// Test bool field in oneof - should clear sfixed64
		RunTestThatExpression(t, "pbt_oneof_fields_set_bool_field(pbt_oneof_fields_set_sfixed64_field(pbt_oneof_fields_new(), -1), TRUE)").IsEqualToJsonString(`{"13": true}`)

		// Test string field in oneof - should clear bool
		RunTestThatExpression(t, "pbt_oneof_fields_set_string_field(pbt_oneof_fields_set_bool_field(pbt_oneof_fields_new(), TRUE), 'hello world')").IsEqualToJsonString(`{"14": "hello world"}`)

		// Test bytes field in oneof - should clear string
		RunTestThatExpression(t, "pbt_oneof_fields_set_bytes_field(pbt_oneof_fields_set_string_field(pbt_oneof_fields_new(), 'hello'), ?)", []byte("world")).IsEqualToJsonString(`{"15": "d29ybGQ="}`)

		// Test enum field in oneof - should clear bytes
		RunTestThatExpression(t, "pbt_oneof_fields_set_enum_field(pbt_oneof_fields_set_bytes_field(pbt_oneof_fields_new(), ?), 1)", []byte("test")).IsEqualToJsonString(`{"16": 1}`)

		// Test message field in oneof - should clear enum
		nestedObj := "pbt_nested_set_value(pbt_nested_set_name(pbt_nested_new(), 'test'), 42)"
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_set_message_field(pbt_oneof_fields_set_enum_field(pbt_oneof_fields_new(), 1), %s)", nestedObj)).IsEqualToJsonString(`{"17": {"1": "test", "2": 42}}`)

		// Test that default values in oneof are NOT omitted (oneof has presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_set_double_field(pbt_oneof_fields_new(), 0.0)").IsEqualToJsonString(`{"1": "binary64:0x0000000000000000"}`) // Zero double stored
		RunTestThatExpression(t, "pbt_oneof_fields_set_int32_field(pbt_oneof_fields_new(), 0)").IsEqualToJsonString(`{"3": 0}`)                                // Zero int32 stored
		RunTestThatExpression(t, "pbt_oneof_fields_set_bool_field(pbt_oneof_fields_new(), FALSE)").IsEqualToJsonString(`{"13": false}`)                        // False bool stored
		RunTestThatExpression(t, "pbt_oneof_fields_set_string_field(pbt_oneof_fields_new(), '')").IsEqualToJsonString(`{"14": ""}`)                            // Empty string stored
	})

	// Test double field in oneof
	t.Run("double_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_oneof_fields_set_double_field(?, 3.141592653589793)", `{}`).IsEqualToJsonString(`{"1": "binary64:0x400921fb54442d18"}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_double_field(?, 0.0)", `{}`).IsEqualToJsonString(`{"1": "binary64:0x0000000000000000"}`)

		// Test getters convert from internal representation back to actual values
		RunTestThatExpression(t, "pbt_oneof_fields_get_double_field(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToDouble(3.141592653589793)
		RunTestThatExpression(t, "pbt_oneof_fields_get_double_field(?)", `{"1": "binary64:0x0000000000000000"}`).IsEqualToDouble(0.0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_double_field(?)", `{}`).IsEqualToDouble(0.0) // Default when absent

		// Test has methods for oneof fields
		RunTestThatExpression(t, "pbt_oneof_fields_has_double_field(?)", `{}`).IsFalse()                                  // Unset field not present
		RunTestThatExpression(t, "pbt_oneof_fields_has_double_field(?)", `{"1": "binary64:0x0000000000000000"}`).IsTrue() // Set field present (even default value)

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{}`).IsNull() // No field set
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToInt(1)

		// Test clear field methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_double_field(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToJsonString(`{}`)

		// Test clear oneof group
		RunTestThatExpression(t, "pbt_oneof_fields_clear_choice(?)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_double_field__or(?, 99.9)", `{"1": "binary64:0x400921fb54442d18"}`).IsEqualToDouble(3.141592653589793) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_double_field__or(?, 99.9)", `{}`).IsEqualToDouble(99.9)                                                // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_double_field__or(?, 99.9)", `{"3": 42}`).IsEqualToDouble(99.9)                                         // different oneof field set, return default
	})

	// Test float field in oneof
	t.Run("float_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_float_field(?, 3.14)", `{}`).IsEqualToJsonString(`{"2": "binary32:0x4048f5c3"}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_float_field(?, 0.0)", `{}`).IsEqualToJsonString(`{"2": "binary32:0x00000000"}`)

		RunTestThatExpression(t, "pbt_oneof_fields_get_float_field(?)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToFloat(3.14)
		RunTestThatExpression(t, "pbt_oneof_fields_get_float_field(?)", `{"2": "binary32:0x00000000"}`).IsEqualToFloat(0.0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_float_field(?)", `{}`).IsEqualToFloat(0.0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_float_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_float_field(?)", `{"2": "binary32:0x00000000"}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToInt(2)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_float_field(?)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_float_field__or(?, 88.8)", `{"2": "binary32:0x4048f5c3"}`).IsEqualToFloat(3.14) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_float_field__or(?, 88.8)", `{}`).IsEqualToFloat(88.8)                           // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_float_field__or(?, 88.8)", `{"3": 42}`).IsEqualToFloat(88.8)                    // different oneof field set, return default
	})

	// Test integer fields
	t.Run("int32_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_int32_field(?, 42)", `{}`).IsEqualToJsonString(`{"3": 42}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_int32_field(?, 0)", `{}`).IsEqualToJsonString(`{"3": 0}`) // Zero stored (oneof presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_int32_field(?)", `{"3": 42}`).IsEqualToInt(42)
		RunTestThatExpression(t, "pbt_oneof_fields_get_int32_field(?)", `{"3": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_int32_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_int32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_int32_field(?)", `{"3": 0}`).IsTrue() // Even zero value has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"3": 42}`).IsEqualToInt(3)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_int32_field(?)", `{"3": 42}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_int32_field__or(?, 999)", `{"3": 42}`).IsEqualToInt(42)        // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_int32_field__or(?, 999)", `{}`).IsEqualToInt(999)              // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_int32_field__or(?, 999)", `{"14": "hello"}`).IsEqualToInt(999) // different oneof field set, return default
	})

	t.Run("int64_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_int64_field(?, 9223372036854775807)", `{}`).IsEqualToJsonString(`{"4": 9223372036854775807}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_int64_field(?, 0)", `{}`).IsEqualToJsonString(`{"4": 0}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_int64_field(?)", `{"4": 9223372036854775807}`).IsEqualToInt(9223372036854775807)
		RunTestThatExpression(t, "pbt_oneof_fields_get_int64_field(?)", `{"4": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_int64_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_int64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_int64_field(?)", `{"4": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"4": 9223372036854775807}`).IsEqualToInt(4)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_int64_field(?)", `{"4": 9223372036854775807}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_int64_field__or(?, 777)", `{"4": 9223372036854775807}`).IsEqualToInt(9223372036854775807) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_int64_field__or(?, 777)", `{}`).IsEqualToInt(777)                                         // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_int64_field__or(?, 777)", `{"3": 42}`).IsEqualToInt(777)                                  // different oneof field set, return default
	})

	t.Run("uint32_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_uint32_field(?, 4294967295)", `{}`).IsEqualToJsonString(`{"5": 4294967295}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_uint32_field(?, 0)", `{}`).IsEqualToJsonString(`{"5": 0}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint32_field(?)", `{"5": 4294967295}`).IsEqualToUint(4294967295)
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint32_field(?)", `{"5": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint32_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_uint32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_uint32_field(?)", `{"5": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"5": 4294967295}`).IsEqualToInt(5)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_uint32_field(?)", `{"5": 4294967295}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint32_field__or(?, 888)", `{"5": 4294967295}`).IsEqualToUint(4294967295) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint32_field__or(?, 888)", `{}`).IsEqualToUint(888)                      // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint32_field__or(?, 888)", `{"3": 42}`).IsEqualToUint(888)               // different oneof field set, return default
	})

	t.Run("uint64_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_uint64_field(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"6": 18446744073709551615}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_uint64_field(?, 0)", `{}`).IsEqualToJsonString(`{"6": 0}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint64_field(?)", `{"6": 18446744073709551615}`).IsEqualToUint(18446744073709551615)
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint64_field(?)", `{"6": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint64_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_uint64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_uint64_field(?)", `{"6": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"6": 18446744073709551615}`).IsEqualToInt(6)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_uint64_field(?)", `{"6": 18446744073709551615}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint64_field__or(?, 888)", `{"6": 18446744073709551615}`).IsEqualToUint(18446744073709551615) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint64_field__or(?, 888)", `{}`).IsEqualToUint(888)                                           // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_uint64_field__or(?, 888)", `{"3": 42}`).IsEqualToUint(888)                                    // different oneof field set, return default
	})

	t.Run("sint32_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_sint32_field(?, -1)", `{}`).IsEqualToJsonString(`{"7": -1}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_sint32_field(?, 0)", `{}`).IsEqualToJsonString(`{"7": 0}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint32_field(?)", `{"7": -1}`).IsEqualToInt(-1)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint32_field(?)", `{"7": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint32_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_sint32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_sint32_field(?)", `{"7": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"7": -1}`).IsEqualToInt(7)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_sint32_field(?)", `{"7": -1}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint32_field__or(?, 555)", `{"7": -1}`).IsEqualToInt(-1)     // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint32_field__or(?, 555)", `{}`).IsEqualToInt(555)          // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint32_field__or(?, 555)", `{"3": 42}`).IsEqualToInt(555)   // different oneof field set, return default
	})

	t.Run("sint64_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_sint64_field(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"8": -9223372036854775808}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_sint64_field(?, 0)", `{}`).IsEqualToJsonString(`{"8": 0}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint64_field(?)", `{"8": -9223372036854775808}`).IsEqualToInt(-9223372036854775808)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint64_field(?)", `{"8": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint64_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_sint64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_sint64_field(?)", `{"8": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"8": -9223372036854775808}`).IsEqualToInt(8)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_sint64_field(?)", `{"8": -9223372036854775808}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint64_field__or(?, 444)", `{"8": -9223372036854775808}`).IsEqualToInt(-9223372036854775808) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint64_field__or(?, 444)", `{}`).IsEqualToInt(444)                                         // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_sint64_field__or(?, 444)", `{"3": 42}`).IsEqualToInt(444)                                  // different oneof field set, return default
	})

	t.Run("fixed32_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_fixed32_field(?, 4294967295)", `{}`).IsEqualToJsonString(`{"9": 4294967295}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_fixed32_field(?, 0)", `{}`).IsEqualToJsonString(`{"9": 0}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed32_field(?)", `{"9": 4294967295}`).IsEqualToUint(4294967295)
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed32_field(?)", `{"9": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed32_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_fixed32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_fixed32_field(?)", `{"9": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"9": 4294967295}`).IsEqualToInt(9)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_fixed32_field(?)", `{"9": 4294967295}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed32_field__or(?, 333)", `{"9": 4294967295}`).IsEqualToUint(4294967295) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed32_field__or(?, 333)", `{}`).IsEqualToUint(333)                      // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed32_field__or(?, 333)", `{"3": 42}`).IsEqualToUint(333)               // different oneof field set, return default
	})

	t.Run("fixed64_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_fixed64_field(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"10": 18446744073709551615}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_fixed64_field(?, 0)", `{}`).IsEqualToJsonString(`{"10": 0}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed64_field(?)", `{"10": 18446744073709551615}`).IsEqualToUint(18446744073709551615)
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed64_field(?)", `{"10": 0}`).IsEqualToUint(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed64_field(?)", `{}`).IsEqualToUint(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_fixed64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_fixed64_field(?)", `{"10": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"10": 18446744073709551615}`).IsEqualToInt(10)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_fixed64_field(?)", `{"10": 18446744073709551615}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed64_field__or(?, 222)", `{"10": 18446744073709551615}`).IsEqualToUint(18446744073709551615) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed64_field__or(?, 222)", `{}`).IsEqualToUint(222)                                           // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_fixed64_field__or(?, 222)", `{"3": 42}`).IsEqualToUint(222)                                    // different oneof field set, return default
	})

	t.Run("sfixed32_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_sfixed32_field(?, -2147483648)", `{}`).IsEqualToJsonString(`{"11": -2147483648}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_sfixed32_field(?, 0)", `{}`).IsEqualToJsonString(`{"11": 0}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed32_field(?)", `{"11": -2147483648}`).IsEqualToInt(-2147483648)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed32_field(?)", `{"11": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed32_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_sfixed32_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_sfixed32_field(?)", `{"11": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"11": -2147483648}`).IsEqualToInt(11)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_sfixed32_field(?)", `{"11": -2147483648}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed32_field__or(?, 111)", `{"11": -2147483648}`).IsEqualToInt(-2147483648) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed32_field__or(?, 111)", `{}`).IsEqualToInt(111)                        // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed32_field__or(?, 111)", `{"3": 42}`).IsEqualToInt(111)                 // different oneof field set, return default
	})

	t.Run("sfixed64_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_sfixed64_field(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"12": -9223372036854775808}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_sfixed64_field(?, 0)", `{}`).IsEqualToJsonString(`{"12": 0}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed64_field(?)", `{"12": -9223372036854775808}`).IsEqualToInt(-9223372036854775808)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed64_field(?)", `{"12": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed64_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_sfixed64_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_sfixed64_field(?)", `{"12": 0}`).IsTrue()

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"12": -9223372036854775808}`).IsEqualToInt(12)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_sfixed64_field(?)", `{"12": -9223372036854775808}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed64_field__or(?, 666)", `{"12": -9223372036854775808}`).IsEqualToInt(-9223372036854775808) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed64_field__or(?, 666)", `{}`).IsEqualToInt(666)                                         // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_sfixed64_field__or(?, 666)", `{"3": 42}`).IsEqualToInt(666)                                  // different oneof field set, return default
	})

	// Test bool field in oneof
	t.Run("bool_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_bool_field(?, TRUE)", `{}`).IsEqualToJsonString(`{"13": true}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_bool_field(?, FALSE)", `{}`).IsEqualToJsonString(`{"13": false}`) // False stored (oneof presence semantics)

		RunTestThatExpression(t, "pbt_oneof_fields_get_bool_field(?)", `{"13": true}`).IsEqualToBool(true)
		RunTestThatExpression(t, "pbt_oneof_fields_get_bool_field(?)", `{"13": false}`).IsEqualToBool(false)
		RunTestThatExpression(t, "pbt_oneof_fields_get_bool_field(?)", `{}`).IsEqualToBool(false) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_bool_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_bool_field(?)", `{"13": false}`).IsTrue() // Even false value has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"13": true}`).IsEqualToInt(13)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_bool_field(?)", `{"13": true}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_bool_field__or(?, TRUE)", `{"13": false}`).IsEqualToBool(false) // field present with false value, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_bool_field__or(?, TRUE)", `{}`).IsEqualToBool(true)             // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_bool_field__or(?, TRUE)", `{"3": 42}`).IsEqualToBool(true)      // different oneof field set, return default
	})

	// Test string field in oneof
	t.Run("string_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_string_field(?, 'hello world')", `{}`).IsEqualToJsonString(`{"14": "hello world"}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_string_field(?, '')", `{}`).IsEqualToJsonString(`{"14": ""}`) // Empty string stored (oneof presence semantics)

		RunTestThatExpression(t, "pbt_oneof_fields_get_string_field(?)", `{"14": "hello world"}`).IsEqualToString("hello world")
		RunTestThatExpression(t, "pbt_oneof_fields_get_string_field(?)", `{"14": ""}`).IsEqualToString("")
		RunTestThatExpression(t, "pbt_oneof_fields_get_string_field(?)", `{}`).IsEqualToString("") // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_string_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_string_field(?)", `{"14": ""}`).IsTrue() // Even empty string has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"14": "hello world"}`).IsEqualToInt(14)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_string_field(?)", `{"14": "hello world"}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_string_field__or(?, 'default')", `{"14": "hello world"}`).IsEqualToString("hello world") // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_string_field__or(?, 'default')", `{}`).IsEqualToString("default")                        // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_string_field__or(?, 'default')", `{"3": 42}`).IsEqualToString("default")                 // different oneof field set, return default
	})

	// Test bytes field in oneof
	t.Run("bytes_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_bytes_field(?, ?)", `{}`, []byte("hello")).IsEqualToJsonString(`{"15": "aGVsbG8="}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_bytes_field(?, ?)", `{}`, []byte{}).IsEqualToJsonString(`{"15": ""}`) // Empty bytes stored (oneof presence semantics)

		RunTestThatExpression(t, "pbt_oneof_fields_get_bytes_field(?)", `{"15": "aGVsbG8="}`).IsEqualToBytes([]byte("hello"))
		RunTestThatExpression(t, "pbt_oneof_fields_get_bytes_field(?)", `{"15": ""}`).IsEqualToBytes([]byte{})
		RunTestThatExpression(t, "pbt_oneof_fields_get_bytes_field(?)", `{}`).IsEqualToBytes([]byte{}) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_bytes_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_bytes_field(?)", `{"15": ""}`).IsTrue() // Even empty bytes has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"15": "aGVsbG8="}`).IsEqualToInt(15)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_bytes_field(?)", `{"15": "aGVsbG8="}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_bytes_field__or(?, ?)", `{"15": "aGVsbG8="}`, []byte("default")).IsEqualToBytes([]byte("hello")) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_bytes_field__or(?, ?)", `{}`, []byte("default")).IsEqualToBytes([]byte("default"))               // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_bytes_field__or(?, ?)", `{"3": 42}`, []byte("default")).IsEqualToBytes([]byte("default"))        // different oneof field set, return default
	})

	// Test enum field in oneof
	t.Run("enum_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_oneof_fields_set_enum_field(?, 1)", `{}`).IsEqualToJsonString(`{"16": 1}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_enum_field(?, 0)", `{}`).IsEqualToJsonString(`{"16": 0}`) // Zero enum stored (oneof presence semantics)

		RunTestThatExpression(t, "pbt_oneof_fields_get_enum_field(?)", `{"16": 1}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_oneof_fields_get_enum_field(?)", `{"16": 0}`).IsEqualToInt(0)
		RunTestThatExpression(t, "pbt_oneof_fields_get_enum_field(?)", `{}`).IsEqualToInt(0) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_enum_field(?)", `{}`).IsFalse()
		RunTestThatExpression(t, "pbt_oneof_fields_has_enum_field(?)", `{"16": 0}`).IsTrue() // Even default enum value has presence in oneof

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"16": 1}`).IsEqualToInt(16)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_enum_field(?)", `{"16": 1}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		RunTestThatExpression(t, "pbt_oneof_fields_get_enum_field__or(?, 999)", `{"16": 1}`).IsEqualToInt(1)   // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_enum_field__or(?, 999)", `{}`).IsEqualToInt(999)        // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_enum_field__or(?, 999)", `{"3": 42}`).IsEqualToInt(999) // different oneof field set, return default
	})

	// Test message field in oneof
	t.Run("message_field", func(t *testing.T) {
		// Test setters
		RunTestThatExpression(t, "pbt_oneof_fields_set_message_field(?, ?)", `{}`, `{"1": "test", "2": 42}`).IsEqualToJsonString(`{"17": {"1": "test", "2": 42}}`)
		RunTestThatExpression(t, "pbt_oneof_fields_set_message_field(?, ?)", `{}`, `{}`).IsEqualToJsonString(`{"17": {}}`) // Empty message stored (oneof presence semantics)

		// Test getters
		RunTestThatExpression(t, "pbt_oneof_fields_get_message_field(?)", `{"17": {"1": "test", "2": 42}}`).IsEqualToJsonString(`{"1": "test", "2": 42}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_message_field(?)", `{"17": {}}`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, "pbt_oneof_fields_get_message_field(?)", `{}`).IsEqualToJsonString(`{}`) // Default when absent

		// Test has methods
		RunTestThatExpression(t, "pbt_oneof_fields_has_message_field(?)", `{}`).IsFalse()        // Unset field not present
		RunTestThatExpression(t, "pbt_oneof_fields_has_message_field(?)", `{"17": {}}`).IsTrue() // Set field present (even empty message) in oneof

		// Test which oneof method
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(?)", `{"17": {"1": "test", "2": 42}}`).IsEqualToInt(17)

		// Test clear methods
		RunTestThatExpression(t, "pbt_oneof_fields_clear_message_field(?)", `{"17": {"1": "test", "2": 42}}`).IsEqualToJsonString(`{}`)

		// Test __or getter variant (oneof fields have presence semantics)
		defaultMessage := `{"1": "default", "2": 0}`
		RunTestThatExpression(t, "pbt_oneof_fields_get_message_field__or(?, ?)", `{"17": {"1": "test", "2": 42}}`, defaultMessage).IsEqualToJsonString(`{"1": "test", "2": 42}`) // field present, return field value
		RunTestThatExpression(t, "pbt_oneof_fields_get_message_field__or(?, ?)", `{}`, defaultMessage).IsEqualToJsonString(defaultMessage)                                       // field not present, return default
		RunTestThatExpression(t, "pbt_oneof_fields_get_message_field__or(?, ?)", `{"3": 42}`, defaultMessage).IsEqualToJsonString(defaultMessage)                                // different oneof field set, return default
	})

	// Test oneof mutual exclusion behavior
	t.Run("mutual_exclusion", func(t *testing.T) {
		// Setting different fields should clear previous ones
		RunTestThatExpression(t, "pbt_oneof_fields_set_string_field(pbt_oneof_fields_set_int32_field(pbt_oneof_fields_new(), 42), 'hello')").IsEqualToJsonString(`{"14": "hello"}`)
		RunTestThatExpression(t, "pbt_oneof_fields_which_choice(pbt_oneof_fields_set_string_field(pbt_oneof_fields_set_int32_field(pbt_oneof_fields_new(), 42), 'hello'))").IsEqualToInt(14)

		// Previous field should return default value after being cleared
		RunTestThatExpression(t, "pbt_oneof_fields_get_int32_field(pbt_oneof_fields_set_string_field(pbt_oneof_fields_set_int32_field(pbt_oneof_fields_new(), 42), 'hello'))").IsEqualToInt(0)
		RunTestThatExpression(t, "pbt_oneof_fields_has_int32_field(pbt_oneof_fields_set_string_field(pbt_oneof_fields_set_int32_field(pbt_oneof_fields_new(), 42), 'hello'))").IsFalse()

		// Test complex mutual exclusion chain
		expr := "pbt_oneof_fields_set_bool_field(pbt_oneof_fields_set_enum_field(pbt_oneof_fields_set_string_field(pbt_oneof_fields_set_int32_field(pbt_oneof_fields_new(), 42), 'hello'), 1), TRUE)"
		RunTestThatExpression(t, expr).IsEqualToJsonString(`{"13": true}`)
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_which_choice(%s)", expr)).IsEqualToInt(13)
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_get_int32_field(%s)", expr)).IsEqualToInt(0)      // cleared
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_get_string_field(%s)", expr)).IsEqualToString("") // cleared
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_get_enum_field(%s)", expr)).IsEqualToInt(0)       // cleared
	})

	// Test clearing entire oneof group
	t.Run("clear_oneof_group", func(t *testing.T) {
		// Set a field, then clear the entire oneof group
		expr := "pbt_oneof_fields_clear_choice(pbt_oneof_fields_set_string_field(pbt_oneof_fields_new(), 'hello'))"
		RunTestThatExpression(t, expr).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_which_choice(%s)", expr)).IsNull()
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_get_string_field(%s)", expr)).IsEqualToString("") // returns default
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_has_string_field(%s)", expr)).IsFalse()

		// Test clearing after setting multiple fields (mutual exclusion)
		expr2 := "pbt_oneof_fields_clear_choice(pbt_oneof_fields_set_int32_field(pbt_oneof_fields_set_string_field(pbt_oneof_fields_new(), 'hello'), 42))"
		RunTestThatExpression(t, expr2).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_which_choice(%s)", expr2)).IsNull()
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_get_int32_field(%s)", expr2)).IsEqualToInt(0) // returns default
		RunTestThatExpression(t, fmt.Sprintf("pbt_oneof_fields_has_int32_field(%s)", expr2)).IsFalse()
	})
}
