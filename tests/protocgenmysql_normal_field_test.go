package main

import (
	"fmt"
	"testing"
)

func TestProtocGenNormalField(t *testing.T) {
	// Test all protobuf field types using the pre-generated functions from protocgenmysql_proto3.pb.sql
	// The NormalFields message in protocgenmysql_proto3.proto matches the test schema

	// Test double field (IEEE 754 binary64 format)
	t.Run("double_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_normal_fields_set_double_field(?, 3.141592653589793)", `{}`).IsEqualToJsonString(`{"1": "binary64:0x400921fb54442d18"}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_double_field(?, 1.0)", `{}`).IsEqualToJsonString(`{"1": "binary64:0x3ff0000000000000"}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_double_field(?, 0.0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted

		// Test getters convert from binary64 format back to actual double
		RunTestThatExpression(t, `pbt_normal_fields_get_double_field('{"1": "binary64:0x400921fb54442d18"}')`).IsEqualTo(3.141592653589793)
		RunTestThatExpression(t, `pbt_normal_fields_get_double_field('{"1": "binary64:0x3ff0000000000000"}')`).IsEqualTo(1.0)
		RunTestThatExpression(t, `pbt_normal_fields_get_double_field('{}')`).IsEqualTo(0.0) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_double_field('{"1": "binary64:0x400921fb54442d18"}')`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `pbt_normal_fields_clear_double_field('{"1": "binary64:0x3ff0000000000000"}')`).IsEqualToJsonString(`{}`)
	})

	// Test float field (IEEE 754 binary32 format)
	t.Run("float_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_normal_fields_set_float_field(?, 3.14)", `{}`).IsEqualToJsonString(`{"2": "binary32:0x4048f5c3"}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_float_field(?, 1.0)", `{}`).IsEqualToJsonString(`{"2": "binary32:0x3f800000"}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_float_field(?, 0.0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted

		// Test getters convert from binary32 format back to actual float
		RunTestThatExpression(t, `pbt_normal_fields_get_float_field('{"2": "binary32:0x4048f5c3"}')`).IsEqualTo(3.14) // MySQL returns 3.14 for float precision
		RunTestThatExpression(t, `pbt_normal_fields_get_float_field('{"2": "binary32:0x3f800000"}')`).IsEqualTo(1.0)
		RunTestThatExpression(t, `pbt_normal_fields_get_float_field('{}')`).IsEqualTo(0.0) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_float_field('{"2": "binary32:0x4048f5c3"}')`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `pbt_normal_fields_clear_float_field('{"2": "binary32:0x3f800000"}')`).IsEqualToJsonString(`{}`)
	})

	// Test int32 field (JSON numbers)
	t.Run("int32_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_normal_fields_set_int32_field(?, 42)", `{}`).IsEqualToJsonString(`{"3": 42}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_int32_field(?, -2147483648)", `{}`).IsEqualToJsonString(`{"3": -2147483648}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_int32_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted

		// Test getters return values directly (no conversion needed)
		RunTestThatExpression(t, `pbt_normal_fields_get_int32_field('{"3": 42}')`).IsEqualTo(42)
		RunTestThatExpression(t, `pbt_normal_fields_get_int32_field('{"3": -2147483648}')`).IsEqualTo(-2147483648)
		RunTestThatExpression(t, `pbt_normal_fields_get_int32_field('{}')`).IsEqualTo(0) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_int32_field('{"3": 42}')`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `pbt_normal_fields_clear_int32_field('{"3": -2147483648}')`).IsEqualToJsonString(`{}`)
	})

	// Test int64 field (JSON numbers, not strings per protonumberjson spec)
	t.Run("int64_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_normal_fields_set_int64_field(?, 9223372036854775807)", `{}`).IsEqualToJsonString(`{"4": 9223372036854775807}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_int64_field(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"4": -9223372036854775808}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_int64_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted

		// Test getters
		RunTestThatExpression(t, `pbt_normal_fields_get_int64_field('{"4": 9223372036854775807}')`).IsEqualTo(int64(9223372036854775807))
		RunTestThatExpression(t, `pbt_normal_fields_get_int64_field('{"4": -9223372036854775808}')`).IsEqualTo(int64(-9223372036854775808))
		RunTestThatExpression(t, `pbt_normal_fields_get_int64_field('{}')`).IsEqualTo(int64(0)) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_int64_field('{"4": 9223372036854775807}')`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `pbt_normal_fields_clear_int64_field('{"4": -9223372036854775808}')`).IsEqualToJsonString(`{}`)
	})

	// Test remaining integer types
	t.Run("uint32_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_normal_fields_set_uint32_field(?, 4294967295)", `{}`).IsEqualToJsonString(`{"5": 4294967295}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_uint32_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted
		RunTestThatExpression(t, `pbt_normal_fields_get_uint32_field('{"5": 4294967295}')`).IsEqualTo(uint32(4294967295))
		RunTestThatExpression(t, `pbt_normal_fields_get_uint32_field('{}')`).IsEqualTo(uint32(0)) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_uint32_field('{"5": 4294967295}')`).IsEqualToJsonString(`{}`)
	})

	t.Run("uint64_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_normal_fields_set_uint64_field(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"6": 18446744073709551615}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_uint64_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted
		RunTestThatExpression(t, `pbt_normal_fields_get_uint64_field('{"6": 18446744073709551615}')`).IsEqualTo(uint64(18446744073709551615))
		RunTestThatExpression(t, `pbt_normal_fields_get_uint64_field('{}')`).IsEqualTo(uint64(0)) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_uint64_field('{"6": 18446744073709551615}')`).IsEqualToJsonString(`{}`)
	})

	t.Run("sint32_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_normal_fields_set_sint32_field(?, -1)", `{}`).IsEqualToJsonString(`{"7": -1}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_sint32_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted
		RunTestThatExpression(t, `pbt_normal_fields_get_sint32_field('{"7": -1}')`).IsEqualTo(-1)
		RunTestThatExpression(t, `pbt_normal_fields_get_sint32_field('{}')`).IsEqualTo(0) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_sint32_field('{"7": -1}')`).IsEqualToJsonString(`{}`)
	})

	t.Run("sint64_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_normal_fields_set_sint64_field(?, -1)", `{}`).IsEqualToJsonString(`{"8": -1}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_sint64_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted
		RunTestThatExpression(t, `pbt_normal_fields_get_sint64_field('{"8": -1}')`).IsEqualTo(int64(-1))
		RunTestThatExpression(t, `pbt_normal_fields_get_sint64_field('{}')`).IsEqualTo(int64(0)) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_sint64_field('{"8": -1}')`).IsEqualToJsonString(`{}`)
	})

	t.Run("fixed32_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_normal_fields_set_fixed32_field(?, 4294967295)", `{}`).IsEqualToJsonString(`{"9": 4294967295}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_fixed32_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted
		RunTestThatExpression(t, `pbt_normal_fields_get_fixed32_field('{"9": 4294967295}')`).IsEqualTo(uint32(4294967295))
		RunTestThatExpression(t, `pbt_normal_fields_get_fixed32_field('{}')`).IsEqualTo(uint32(0)) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_fixed32_field('{"9": 4294967295}')`).IsEqualToJsonString(`{}`)
	})

	t.Run("fixed64_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_normal_fields_set_fixed64_field(?, 18446744073709551615)", `{}`).IsEqualToJsonString(`{"10": 18446744073709551615}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_fixed64_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted
		RunTestThatExpression(t, `pbt_normal_fields_get_fixed64_field('{"10": 18446744073709551615}')`).IsEqualTo(uint64(18446744073709551615))
		RunTestThatExpression(t, `pbt_normal_fields_get_fixed64_field('{}')`).IsEqualTo(uint64(0)) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_fixed64_field('{"10": 18446744073709551615}')`).IsEqualToJsonString(`{}`)
	})

	t.Run("sfixed32_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_normal_fields_set_sfixed32_field(?, -2147483648)", `{}`).IsEqualToJsonString(`{"11": -2147483648}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_sfixed32_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted
		RunTestThatExpression(t, `pbt_normal_fields_get_sfixed32_field('{"11": -2147483648}')`).IsEqualTo(-2147483648)
		RunTestThatExpression(t, `pbt_normal_fields_get_sfixed32_field('{}')`).IsEqualTo(0) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_sfixed32_field('{"11": -2147483648}')`).IsEqualToJsonString(`{}`)
	})

	t.Run("sfixed64_field", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_normal_fields_set_sfixed64_field(?, -9223372036854775808)", `{}`).IsEqualToJsonString(`{"12": -9223372036854775808}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_sfixed64_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero omitted
		RunTestThatExpression(t, `pbt_normal_fields_get_sfixed64_field('{"12": -9223372036854775808}')`).IsEqualTo(int64(-9223372036854775808))
		RunTestThatExpression(t, `pbt_normal_fields_get_sfixed64_field('{}')`).IsEqualTo(int64(0)) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_sfixed64_field('{"12": -9223372036854775808}')`).IsEqualToJsonString(`{}`)
	})

	// Test bool field (JSON booleans, not 1/0)
	t.Run("bool_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_normal_fields_set_bool_field(?, TRUE)", `{}`).IsEqualToJsonString(`{"13": true}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_bool_field(?, FALSE)", `{}`).IsEqualToJsonString(`{}`) // False omitted

		// Test getters return actual boolean from JSON boolean
		RunTestThatExpression(t, `pbt_normal_fields_get_bool_field('{"13": true}')`).IsTrue()
		RunTestThatExpression(t, `pbt_normal_fields_get_bool_field('{"13": false}')`).IsFalse()
		RunTestThatExpression(t, `pbt_normal_fields_get_bool_field('{}')`).IsFalse() // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_bool_field('{"13": true}')`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `pbt_normal_fields_clear_bool_field('{"13": false}')`).IsEqualToJsonString(`{}`)
	})

	// Test string field
	t.Run("string_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_normal_fields_set_string_field(?, 'hello world')", `{}`).IsEqualToJsonString(`{"14": "hello world"}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_string_field(?, '')", `{}`).IsEqualToJsonString(`{}`) // Empty string omitted

		// Test getters return actual string
		RunTestThatExpression(t, `pbt_normal_fields_get_string_field('{"14": "hello world"}')`).IsEqualTo("hello world")
		RunTestThatExpression(t, `pbt_normal_fields_get_string_field('{"14": ""}')`).IsEqualTo("") // Empty string
		RunTestThatExpression(t, `pbt_normal_fields_get_string_field('{}')`).IsEqualTo("")         // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_string_field('{"14": "hello world"}')`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `pbt_normal_fields_clear_string_field('{"14": ""}')`).IsEqualToJsonString(`{}`)
	})

	// Test bytes field (base64 encoded)
	t.Run("bytes_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_normal_fields_set_bytes_field(?, ?)", `{}`, []byte("hello world")).IsEqualToJsonString(`{"15": "aGVsbG8gd29ybGQ="}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_bytes_field(?, ?)", `{}`, []byte{0xDE, 0xAD, 0xBE, 0xEF}).IsEqualToJsonString(`{"15": "3q2+7w=="}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_bytes_field(?, ?)", `{}`, []byte{}).IsEqualToJsonString(`{}`) // Empty bytes omitted

		// Test getters convert from base64 back to actual bytes
		RunTestThatExpression(t, `pbt_normal_fields_get_bytes_field('{"15": "aGVsbG8gd29ybGQ="}')`).IsEqualTo([]byte("hello world"))
		RunTestThatExpression(t, `pbt_normal_fields_get_bytes_field('{"15": "3q2+7w=="}')`).IsEqualTo([]byte{0xDE, 0xAD, 0xBE, 0xEF})
		RunTestThatExpression(t, `pbt_normal_fields_get_bytes_field('{"15": ""}')`).IsEqualTo([]byte{}) // Empty base64
		RunTestThatExpression(t, `pbt_normal_fields_get_bytes_field('{}')`).IsEqualTo([]byte{})         // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_bytes_field('{"15": "aGVsbG8gd29ybGQ="}')`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `pbt_normal_fields_clear_bytes_field('{"15": "3q2+7w=="}')`).IsEqualToJsonString(`{}`)
	})

	// Test enum field (stored as numbers, not string names)
	t.Run("enum_field", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_normal_fields_set_enum_field(?, 1)", `{}`).IsEqualToJsonString(`{"16": 1}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_enum_field(?, 2)", `{}`).IsEqualToJsonString(`{"16": 2}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_enum_field(?, 0)", `{}`).IsEqualToJsonString(`{}`) // Zero enum omitted

		// Test getters return actual integer value
		RunTestThatExpression(t, `pbt_normal_fields_get_enum_field('{"16": 1}')`).IsEqualTo(1)
		RunTestThatExpression(t, `pbt_normal_fields_get_enum_field('{"16": 2}')`).IsEqualTo(2)
		RunTestThatExpression(t, `pbt_normal_fields_get_enum_field('{}')`).IsEqualTo(0) // Missing field returns default

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_enum_field('{"16": 1}')`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `pbt_normal_fields_clear_enum_field('{"16": 2}')`).IsEqualToJsonString(`{}`)

		// Test enum name getter (__as_name) methods
		RunTestThatExpression(t, `pbt_normal_fields_get_enum_field__as_name('{"16": 0}')`).IsEqualToString("STATUS_UNSPECIFIED") // Zero value name
		RunTestThatExpression(t, `pbt_normal_fields_get_enum_field__as_name('{"16": 1}')`).IsEqualToString("STATUS_ACTIVE")      // Active status name
		RunTestThatExpression(t, `pbt_normal_fields_get_enum_field__as_name('{"16": 2}')`).IsEqualToString("STATUS_INACTIVE")    // Inactive status name
		RunTestThatExpression(t, `pbt_normal_fields_get_enum_field__as_name('{"16": 999}')`).IsEqualToString("999")              // Unknown enum value returns number as string
		RunTestThatExpression(t, `pbt_normal_fields_get_enum_field__as_name('{}')`).IsEqualToString("STATUS_UNSPECIFIED")        // Default when absent

		// Test enum name setter (__from_name) methods
		RunTestThatExpression(t, `pbt_normal_fields_set_enum_field__from_name('{}', 'STATUS_ACTIVE')`).IsEqualToJsonString(`{"16": 1}`)                       // Set from valid name
		RunTestThatExpression(t, `pbt_normal_fields_set_enum_field__from_name('{}', 'STATUS_INACTIVE')`).IsEqualToJsonString(`{"16": 2}`)                     // Set from different name
		RunTestThatExpression(t, `pbt_normal_fields_set_enum_field__from_name('{}', 'STATUS_UNSPECIFIED')`).IsEqualToJsonString(`{}`)                         // Zero value omitted in proto3
		RunTestThatExpression(t, `pbt_normal_fields_set_enum_field__from_name('{}', 'INVALID_NAME')`).ToFailWithSignalException("45000", "Invalid enum name") // Invalid name should signal error
	})

	// Test message field (nested object with field number keys)
	t.Run("message_field", func(t *testing.T) {
		// Test setters create correct internal format
		nestedObj := "pbt_nested_set_value(pbt_nested_set_name(pbt_nested_new(), 'test'), 42)"
		RunTestThatExpression(t, fmt.Sprintf("pbt_normal_fields_set_message_field(?, %s)", nestedObj), `{}`).IsEqualToJsonString(`{"17": {"1": "test", "2": 42}}`)
		RunTestThatExpression(t, "pbt_normal_fields_set_message_field(?, pbt_nested_new())", `{}`).IsEqualToJsonString(`{"17": {}}`) // Empty message is stored

		// Test getters return actual nested message as JSON object
		RunTestThatExpression(t, `pbt_normal_fields_get_message_field('{"17": {"1": "test", "2": 42}}')`).IsEqualToJsonString(`{"1": "test", "2": 42}`)
		RunTestThatExpression(t, `pbt_normal_fields_get_message_field('{"17": {"1": "hello"}}')`).IsEqualToJsonString(`{"1": "hello"}`) // Partial message
		RunTestThatExpression(t, `pbt_normal_fields_get_message_field('{"17": {}}')`).IsEqualToJsonString(`{}`)                         // Empty message
		RunTestThatExpression(t, `pbt_normal_fields_get_message_field('{}')`).IsEqualTo("{}")                                           // Missing field returns empty object

		// Test clear methods remove field and return empty JSON
		RunTestThatExpression(t, `pbt_normal_fields_clear_message_field('{"17": {"1": "test", "2": 42}}')`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, `pbt_normal_fields_clear_message_field('{"17": {}}')`).IsEqualToJsonString(`{}`)

		// Test nullable getter (__or) methods for message fields
		// Note: In proto3, message fields have presence semantics, so __or variants are generated
		defaultMessage := `{"1": "default", "2": 999}`
		RunTestThatExpression(t, fmt.Sprintf("pbt_normal_fields_get_message_field__or(?, %s)", `'{"1": "default", "2": 999}'`), `{"17": {"1": "test", "2": 42}}`).IsEqualToJsonString(`{"1": "test", "2": 42}`) // Field present, return actual value
		RunTestThatExpression(t, fmt.Sprintf("pbt_normal_fields_get_message_field__or(?, %s)", `'{"1": "default", "2": 999}'`), `{}`).IsEqualToJsonString(defaultMessage)                                       // Field absent, return default
		RunTestThatExpression(t, fmt.Sprintf("pbt_normal_fields_get_message_field__or(?, %s)", `'{"1": "default", "2": 999}'`), `{"17": {}}`).IsEqualToJsonString(`{}`)                                         // Empty message is still present
	})
}
