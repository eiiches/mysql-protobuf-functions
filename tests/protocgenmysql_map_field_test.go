package main

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
)

func TestProtocGenMapField(t *testing.T) {
	protoContent := dedent.Pipe(`
		|syntax = "proto3";
		|message Test {
		|    // Every possible key type
		|    map<int32, int32> int32_to_int32_map = 1;
		|    map<int64, int32> int64_to_int32_map = 2;
		|    map<uint32, int32> uint32_to_int32_map = 3;
		|    map<uint64, int32> uint64_to_int32_map = 4;
		|    map<sint32, int32> sint32_to_int32_map = 5;
		|    map<sint64, int32> sint64_to_int32_map = 6;
		|    map<fixed32, int32> fixed32_to_int32_map = 7;
		|    map<fixed64, int32> fixed64_to_int32_map = 8;
		|    map<sfixed32, int32> sfixed32_to_int32_map = 9;
		|    map<sfixed64, int32> sfixed64_to_int32_map = 10;
		|    map<bool, int32> bool_to_int32_map = 11;
		|    map<string, int32> string_to_int32_map = 12;
		|
		|    // Every possible value type  
		|    map<string, double> string_to_double_map = 13;
		|    map<string, float> string_to_float_map = 14;
		|    map<string, int64> string_to_int64_map = 15;
		|    map<string, uint32> string_to_uint32_map = 16;
		|    map<string, uint64> string_to_uint64_map = 17;
		|    map<string, sint32> string_to_sint32_map = 18;
		|    map<string, sint64> string_to_sint64_map = 19;
		|    map<string, fixed32> string_to_fixed32_map = 20;
		|    map<string, fixed64> string_to_fixed64_map = 21;
		|    map<string, sfixed32> string_to_sfixed32_map = 22;
		|    map<string, sfixed64> string_to_sfixed64_map = 23;
		|    map<string, bool> string_to_bool_map = 24;
		|    map<string, string> string_to_string_map = 25;
		|    map<string, bytes> string_to_bytes_map = 26;
		|    map<string, Status> string_to_enum_map = 27;
		|    map<string, Nested> string_to_message_map = 28;
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

	// Test integer key map types
	t.Run("int32_key", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "test_set_all_int32_to_int32_map(?, JSON_OBJECT('42', 100))", `{}`).IsEqualToJsonString(`{"1": {"42": 100}}`)
		RunTestThatExpression(t, "test_set_all_int32_to_int32_map(?, JSON_OBJECT('0', 0))", `{}`).IsEqualToJsonString(`{"1": {"0": 0}}`) // Zero key and value stored

		// Test getters return entire map
		RunTestThatExpression(t, "test_get_all_int32_to_int32_map(?)", `{"1": {"42": 100, "1": 10}}`).IsEqualToJsonString(`{"42": 100, "1": 10}`)
		RunTestThatExpression(t, "test_get_all_int32_to_int32_map(?)", `{}`).IsEqualToJsonString(`[]`) // Default when absent

		// Test map count operations
		RunTestThatExpression(t, "test_count_int32_to_int32_map(?)", `{}`).IsEqualToInt(0)
		RunTestThatExpression(t, "test_count_int32_to_int32_map(?)", `{"1": {"42": 100, "1": 10}}`).IsEqualToInt(2)

		// Test clear methods
		RunTestThatExpression(t, "test_clear_int32_to_int32_map(?)", `{"1": {"42": 100}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_int32_to_int32_map__or(?, ?, ?)", `{"1": {"42": 100}}`, `42`, `999`).IsEqualToInt(100)     // Key exists, return value
		RunTestThatExpression(t, "test_get_int32_to_int32_map__or(?, ?, ?)", `{"1": {"42": 100}}`, `99`, `999`).IsEqualToInt(999)     // Key missing, return default
		RunTestThatExpression(t, "test_get_int32_to_int32_map__or(?, ?, ?)", `{}`, `42`, `999`).IsEqualToInt(999)                      // Map empty, return default
	})

	t.Run("int64_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_int64_to_int32_map(?, JSON_OBJECT('9223372036854775807', 200))", `{}`).IsEqualToJsonString(`{"2": {"9223372036854775807": 200}}`)
		RunTestThatExpression(t, "test_get_all_int64_to_int32_map(?)", `{"2": {"9223372036854775807": 200}}`).IsEqualToJsonString(`{"9223372036854775807": 200}`)
		RunTestThatExpression(t, "test_count_int64_to_int32_map(?)", `{"2": {"9223372036854775807": 200}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_int64_to_int32_map(?)", `{"2": {"9223372036854775807": 200}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_int64_to_int32_map__or(?, ?, ?)", `{"2": {"9223372036854775807": 200}}`, `9223372036854775807`, `888`).IsEqualToInt(200) // Key exists
		RunTestThatExpression(t, "test_get_int64_to_int32_map__or(?, ?, ?)", `{"2": {"9223372036854775807": 200}}`, `123`, `888`).IsEqualToInt(888)                 // Key missing
		RunTestThatExpression(t, "test_get_int64_to_int32_map__or(?, ?, ?)", `{}`, `9223372036854775807`, `888`).IsEqualToInt(888)                                   // Map empty
	})

	t.Run("uint32_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_uint32_to_int32_map(?, JSON_OBJECT('4294967295', 300))", `{}`).IsEqualToJsonString(`{"3": {"4294967295": 300}}`)
		RunTestThatExpression(t, "test_get_all_uint32_to_int32_map(?)", `{"3": {"4294967295": 300}}`).IsEqualToJsonString(`{"4294967295": 300}`)
		RunTestThatExpression(t, "test_count_uint32_to_int32_map(?)", `{"3": {"4294967295": 300}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_uint32_to_int32_map(?)", `{"3": {"4294967295": 300}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_uint32_to_int32_map__or(?, ?, ?)", `{"3": {"4294967295": 300}}`, `4294967295`, `777`).IsEqualToInt(300) // Key exists
		RunTestThatExpression(t, "test_get_uint32_to_int32_map__or(?, ?, ?)", `{"3": {"4294967295": 300}}`, `123`, `777`).IsEqualToInt(777)        // Key missing
		RunTestThatExpression(t, "test_get_uint32_to_int32_map__or(?, ?, ?)", `{}`, `4294967295`, `777`).IsEqualToInt(777)                           // Map empty
	})

	t.Run("uint64_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_uint64_to_int32_map(?, JSON_OBJECT('18446744073709551615', 400))", `{}`).IsEqualToJsonString(`{"4": {"18446744073709551615": 400}}`)
		RunTestThatExpression(t, "test_get_all_uint64_to_int32_map(?)", `{"4": {"18446744073709551615": 400}}`).IsEqualToJsonString(`{"18446744073709551615": 400}`)
		RunTestThatExpression(t, "test_count_uint64_to_int32_map(?)", `{"4": {"18446744073709551615": 400}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_uint64_to_int32_map(?)", `{"4": {"18446744073709551615": 400}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_uint64_to_int32_map__or(?, ?, ?)", `{"4": {"18446744073709551615": 400}}`, `18446744073709551615`, `666`).IsEqualToInt(400) // Key exists
		RunTestThatExpression(t, "test_get_uint64_to_int32_map__or(?, ?, ?)", `{"4": {"18446744073709551615": 400}}`, `123`, `666`).IsEqualToInt(666)                      // Key missing
		RunTestThatExpression(t, "test_get_uint64_to_int32_map__or(?, ?, ?)", `{}`, `18446744073709551615`, `666`).IsEqualToInt(666)                                    // Map empty
	})

	t.Run("sint32_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_sint32_to_int32_map(?, JSON_OBJECT('-1', 500))", `{}`).IsEqualToJsonString(`{"5": {"-1": 500}}`)
		RunTestThatExpression(t, "test_get_all_sint32_to_int32_map(?)", `{"5": {"-1": 500}}`).IsEqualToJsonString(`{"-1": 500}`)
		RunTestThatExpression(t, "test_count_sint32_to_int32_map(?)", `{"5": {"-1": 500}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_sint32_to_int32_map(?)", `{"5": {"-1": 500}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_sint32_to_int32_map__or(?, ?, ?)", `{"5": {"-2147483648": 500}}`, `-2147483648`, `555`).IsEqualToInt(500) // Key exists
		RunTestThatExpression(t, "test_get_sint32_to_int32_map__or(?, ?, ?)", `{"5": {"-2147483648": 500}}`, `123`, `555`).IsEqualToInt(555)        // Key missing
		RunTestThatExpression(t, "test_get_sint32_to_int32_map__or(?, ?, ?)", `{}`, `-2147483648`, `555`).IsEqualToInt(555)                           // Map empty
	})

	t.Run("sint64_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_sint64_to_int32_map(?, JSON_OBJECT('-9223372036854775808', 600))", `{}`).IsEqualToJsonString(`{"6": {"-9223372036854775808": 600}}`)
		RunTestThatExpression(t, "test_get_all_sint64_to_int32_map(?)", `{"6": {"-9223372036854775808": 600}}`).IsEqualToJsonString(`{"-9223372036854775808": 600}`)
		RunTestThatExpression(t, "test_count_sint64_to_int32_map(?)", `{"6": {"-9223372036854775808": 600}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_sint64_to_int32_map(?)", `{"6": {"-9223372036854775808": 600}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_sint64_to_int32_map__or(?, ?, ?)", `{"6": {"-9223372036854775808": 600}}`, `-9223372036854775808`, `444`).IsEqualToInt(600) // Key exists
		RunTestThatExpression(t, "test_get_sint64_to_int32_map__or(?, ?, ?)", `{"6": {"-9223372036854775808": 600}}`, `123`, `444`).IsEqualToInt(444)                      // Key missing
		RunTestThatExpression(t, "test_get_sint64_to_int32_map__or(?, ?, ?)", `{}`, `-9223372036854775808`, `444`).IsEqualToInt(444)                                    // Map empty
	})

	t.Run("fixed32_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_fixed32_to_int32_map(?, JSON_OBJECT('4294967295', 700))", `{}`).IsEqualToJsonString(`{"7": {"4294967295": 700}}`)
		RunTestThatExpression(t, "test_get_all_fixed32_to_int32_map(?)", `{"7": {"4294967295": 700}}`).IsEqualToJsonString(`{"4294967295": 700}`)
		RunTestThatExpression(t, "test_count_fixed32_to_int32_map(?)", `{"7": {"4294967295": 700}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_fixed32_to_int32_map(?)", `{"7": {"4294967295": 700}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_fixed32_to_int32_map__or(?, ?, ?)", `{"7": {"4294967295": 700}}`, `4294967295`, `333`).IsEqualToInt(700) // Key exists
		RunTestThatExpression(t, "test_get_fixed32_to_int32_map__or(?, ?, ?)", `{"7": {"4294967295": 700}}`, `123`, `333`).IsEqualToInt(333)        // Key missing
		RunTestThatExpression(t, "test_get_fixed32_to_int32_map__or(?, ?, ?)", `{}`, `4294967295`, `333`).IsEqualToInt(333)                           // Map empty
	})

	t.Run("fixed64_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_fixed64_to_int32_map(?, JSON_OBJECT('18446744073709551615', 800))", `{}`).IsEqualToJsonString(`{"8": {"18446744073709551615": 800}}`)
		RunTestThatExpression(t, "test_get_all_fixed64_to_int32_map(?)", `{"8": {"18446744073709551615": 800}}`).IsEqualToJsonString(`{"18446744073709551615": 800}`)
		RunTestThatExpression(t, "test_count_fixed64_to_int32_map(?)", `{"8": {"18446744073709551615": 800}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_fixed64_to_int32_map(?)", `{"8": {"18446744073709551615": 800}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_fixed64_to_int32_map__or(?, ?, ?)", `{"8": {"18446744073709551615": 800}}`, `18446744073709551615`, `222`).IsEqualToInt(800) // Key exists
		RunTestThatExpression(t, "test_get_fixed64_to_int32_map__or(?, ?, ?)", `{"8": {"18446744073709551615": 800}}`, `123`, `222`).IsEqualToInt(222)                      // Key missing
		RunTestThatExpression(t, "test_get_fixed64_to_int32_map__or(?, ?, ?)", `{}`, `18446744073709551615`, `222`).IsEqualToInt(222)                                    // Map empty
	})

	t.Run("sfixed32_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_sfixed32_to_int32_map(?, JSON_OBJECT('-2147483648', 900))", `{}`).IsEqualToJsonString(`{"9": {"-2147483648": 900}}`)
		RunTestThatExpression(t, "test_get_all_sfixed32_to_int32_map(?)", `{"9": {"-2147483648": 900}}`).IsEqualToJsonString(`{"-2147483648": 900}`)
		RunTestThatExpression(t, "test_count_sfixed32_to_int32_map(?)", `{"9": {"-2147483648": 900}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_sfixed32_to_int32_map(?)", `{"9": {"-2147483648": 900}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_sfixed32_to_int32_map__or(?, ?, ?)", `{"9": {"-2147483648": 900}}`, `-2147483648`, `111`).IsEqualToInt(900) // Key exists
		RunTestThatExpression(t, "test_get_sfixed32_to_int32_map__or(?, ?, ?)", `{"9": {"-2147483648": 900}}`, `123`, `111`).IsEqualToInt(111)        // Key missing
		RunTestThatExpression(t, "test_get_sfixed32_to_int32_map__or(?, ?, ?)", `{}`, `-2147483648`, `111`).IsEqualToInt(111)                           // Map empty
	})

	t.Run("sfixed64_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_sfixed64_to_int32_map(?, JSON_OBJECT('-9223372036854775808', 1000))", `{}`).IsEqualToJsonString(`{"10": {"-9223372036854775808": 1000}}`)
		RunTestThatExpression(t, "test_get_all_sfixed64_to_int32_map(?)", `{"10": {"-9223372036854775808": 1000}}`).IsEqualToJsonString(`{"-9223372036854775808": 1000}`)
		RunTestThatExpression(t, "test_count_sfixed64_to_int32_map(?)", `{"10": {"-9223372036854775808": 1000}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_sfixed64_to_int32_map(?)", `{"10": {"-9223372036854775808": 1000}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_sfixed64_to_int32_map__or(?, ?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `-9223372036854775808`, `101`).IsEqualToInt(1000) // Key exists
		RunTestThatExpression(t, "test_get_sfixed64_to_int32_map__or(?, ?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `123`, `101`).IsEqualToInt(101)                      // Key missing
		RunTestThatExpression(t, "test_get_sfixed64_to_int32_map__or(?, ?, ?)", `{}`, `-9223372036854775808`, `101`).IsEqualToInt(101)                                    // Map empty
	})

	// Test non-integer key types
	t.Run("bool_key", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_bool_to_int32_map(?, JSON_OBJECT('true', 1100))", `{}`).IsEqualToJsonString(`{"11": {"true": 1100}}`)
		RunTestThatExpression(t, "test_get_all_bool_to_int32_map(?)", `{"11": {"true": 1100, "false": 0}}`).IsEqualToJsonString(`{"true": 1100, "false": 0}`)
		RunTestThatExpression(t, "test_count_bool_to_int32_map(?)", `{"11": {"true": 1100, "false": 0}}`).IsEqualToInt(2)
		RunTestThatExpression(t, "test_clear_bool_to_int32_map(?)", `{"11": {"true": 1100}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_bool_to_int32_map__or(?, ?, ?)", `{"11": {"true": 1100, "false": 0}}`, `true`, `999`).IsEqualToInt(1100)  // Key exists
		RunTestThatExpression(t, "test_get_bool_to_int32_map__or(?, ?, ?)", `{"11": {"true": 1100, "false": 0}}`, `false`, `999`).IsEqualToInt(0)   // Key exists (zero value)
		RunTestThatExpression(t, "test_get_bool_to_int32_map__or(?, ?, ?)", `{}`, `true`, `999`).IsEqualToInt(999)                                      // Map empty
	})

	t.Run("string_key", func(t *testing.T) {
		// Test multiple entries in same map
		RunTestThatExpression(t, "test_set_all_string_to_int32_map(?, JSON_OBJECT('key', 1200))", `{}`).IsEqualToJsonString(`{"12": {"key": 1200}}`)
		RunTestThatExpression(t, "test_set_all_string_to_int32_map(?, JSON_OBJECT('first', 10, 'second', 20))", `{}`).IsEqualToJsonString(`{"12": {"first": 10, "second": 20}}`)

		// Test getters return entire map
		RunTestThatExpression(t, "test_get_all_string_to_int32_map(?)", `{"12": {"first": 10, "second": 20}}`).IsEqualToJsonString(`{"first": 10, "second": 20}`)
		RunTestThatExpression(t, "test_get_all_string_to_int32_map(?)", `{}`).IsEqualToJsonString(`[]`) // Default when absent

		// Test overwriting existing map (replaces entire map)
		RunTestThatExpression(t, "test_set_all_string_to_int32_map(test_set_all_string_to_int32_map(test_new(), JSON_OBJECT('old', 100)), JSON_OBJECT('new', 200))").IsEqualToJsonString(`{"12": {"new": 200}}`)

		// Test map count and clear operations
		RunTestThatExpression(t, "test_count_string_to_int32_map(?)", `{"12": {"first": 10, "second": 20}}`).IsEqualToInt(2)
		RunTestThatExpression(t, "test_clear_string_to_int32_map(?)", `{"12": {"key": 1200}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_int32_map__or(?, ?, ?)", `{"12": {"key": 555}}`, `key`, `9999`).IsEqualToInt(555)      // Key exists
		RunTestThatExpression(t, "test_get_string_to_int32_map__or(?, ?, ?)", `{"12": {"key": 555}}`, `missing`, `9999`).IsEqualToInt(9999) // Key missing
		RunTestThatExpression(t, "test_get_string_to_int32_map__or(?, ?, ?)", `{}`, `key`, `9999`).IsEqualToInt(9999)                       // Map empty
	})

	// Test different value types with string keys
	t.Run("double_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_double_map(?, JSON_OBJECT('pi', 'binary64:0x400921fb54442d18'))", `{}`).IsEqualToJsonString(`{"13": {"pi": "binary64:0x400921fb54442d18"}}`)
		RunTestThatExpression(t, "test_get_all_string_to_double_map(?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`).IsEqualToJsonString(`{"pi": "binary64:0x400921fb54442d18"}`)
		RunTestThatExpression(t, "test_count_string_to_double_map(?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_double_map(?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_double_map__or(?, ?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `pi`, `"binary64:0x4000000000000000"`).IsEqualToString(`"binary64:0x400921fb54442d18"`) // Key exists
		RunTestThatExpression(t, "test_get_string_to_double_map__or(?, ?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `missing`, `"binary64:0x4000000000000000"`).IsEqualToString(`"binary64:0x4000000000000000"`) // Key missing
		RunTestThatExpression(t, "test_get_string_to_double_map__or(?, ?, ?)", `{}`, `pi`, `"binary64:0x4000000000000000"`).IsEqualToString(`"binary64:0x4000000000000000"`) // Map empty
	})

	t.Run("float_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_float_map(?, JSON_OBJECT('pi_float', 'binary32:0x4048f5c3'))", `{}`).IsEqualToJsonString(`{"14": {"pi_float": "binary32:0x4048f5c3"}}`)
		RunTestThatExpression(t, "test_get_all_string_to_float_map(?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`).IsEqualToJsonString(`{"pi_float": "binary32:0x4048f5c3"}`)
		RunTestThatExpression(t, "test_count_string_to_float_map(?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_float_map(?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_float_map__or(?, ?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `pi_float`, `"binary32:0x40000000"`).IsEqualToString(`"binary32:0x4048f5c3"`) // Key exists
		RunTestThatExpression(t, "test_get_string_to_float_map__or(?, ?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `missing`, `"binary32:0x40000000"`).IsEqualToString(`"binary32:0x40000000"`) // Key missing
		RunTestThatExpression(t, "test_get_string_to_float_map__or(?, ?, ?)", `{}`, `pi_float`, `"binary32:0x40000000"`).IsEqualToString(`"binary32:0x40000000"`) // Map empty
	})

	t.Run("int32_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_int32_map(?, JSON_OBJECT('key', -2147483648))", `{}`).IsEqualToJsonString(`{"12": {"key": -2147483648}}`)
		RunTestThatExpression(t, "test_get_all_string_to_int32_map(?)", `{"12": {"key": -2147483648}}`).IsEqualToJsonString(`{"key": -2147483648}`)
		RunTestThatExpression(t, "test_count_string_to_int32_map(?)", `{"12": {"key": -2147483648}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_int32_map(?)", `{"12": {"key": -2147483648}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)  
		RunTestThatExpression(t, "test_get_string_to_int32_map__or(?, ?, ?)", `{"12": {"key": 12345}}`, `key`, `0`).IsEqualToInt(12345) // Key exists
		RunTestThatExpression(t, "test_get_string_to_int32_map__or(?, ?, ?)", `{"12": {"key": 12345}}`, `missing`, `0`).IsEqualToInt(0)   // Key missing
		RunTestThatExpression(t, "test_get_string_to_int32_map__or(?, ?, ?)", `{}`, `key`, `0`).IsEqualToInt(0)                             // Map empty

		// Test that maps can store default/zero values (unlike regular proto3 fields without presence)
		RunTestThatExpression(t, "test_set_all_string_to_int32_map(?, JSON_OBJECT('zero', 0))", `{}`).IsEqualToJsonString(`{"12": {"zero": 0}}`)
	})

	t.Run("int64_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_int64_map(?, JSON_OBJECT('big', 9223372036854775807))", `{}`).IsEqualToJsonString(`{"15": {"big": 9223372036854775807}}`)
		RunTestThatExpression(t, "test_get_all_string_to_int64_map(?)", `{"15": {"big": 9223372036854775807}}`).IsEqualToJsonString(`{"big": 9223372036854775807}`)
		RunTestThatExpression(t, "test_count_string_to_int64_map(?)", `{"15": {"big": 9223372036854775807}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_int64_map(?)", `{"15": {"big": 9223372036854775807}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_int64_map__or(?, ?, ?)", `{"15": {"big": 9223372036854775807}}`, `big`, `-1`).IsEqualToInt(9223372036854775807) // Key exists
		RunTestThatExpression(t, "test_get_string_to_int64_map__or(?, ?, ?)", `{"15": {"big": 9223372036854775807}}`, `missing`, `-1`).IsEqualToInt(-1)                  // Key missing
		RunTestThatExpression(t, "test_get_string_to_int64_map__or(?, ?, ?)", `{}`, `big`, `-1`).IsEqualToInt(-1)                                                       // Map empty
	})

	t.Run("uint32_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_uint32_map(?, JSON_OBJECT('max32', 4294967295))", `{}`).IsEqualToJsonString(`{"16": {"max32": 4294967295}}`)
		RunTestThatExpression(t, "test_get_all_string_to_uint32_map(?)", `{"16": {"max32": 4294967295}}`).IsEqualToJsonString(`{"max32": 4294967295}`)
		RunTestThatExpression(t, "test_count_string_to_uint32_map(?)", `{"16": {"max32": 4294967295}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_uint32_map(?)", `{"16": {"max32": 4294967295}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_uint32_map__or(?, ?, ?)", `{"16": {"max32": 4294967295}}`, `max32`, `0`).IsEqualToUint(4294967295) // Key exists
		RunTestThatExpression(t, "test_get_string_to_uint32_map__or(?, ?, ?)", `{"16": {"max32": 4294967295}}`, `missing`, `0`).IsEqualToUint(0)          // Key missing
		RunTestThatExpression(t, "test_get_string_to_uint32_map__or(?, ?, ?)", `{}`, `max32`, `0`).IsEqualToUint(0)                                      // Map empty
	})

	t.Run("uint64_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_uint64_map(?, JSON_OBJECT('max64', 18446744073709551615))", `{}`).IsEqualToJsonString(`{"17": {"max64": 18446744073709551615}}`)
		RunTestThatExpression(t, "test_get_all_string_to_uint64_map(?)", `{"17": {"max64": 18446744073709551615}}`).IsEqualToJsonString(`{"max64": 18446744073709551615}`)
		RunTestThatExpression(t, "test_count_string_to_uint64_map(?)", `{"17": {"max64": 18446744073709551615}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_uint64_map(?)", `{"17": {"max64": 18446744073709551615}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_uint64_map__or(?, ?, ?)", `{"17": {"max64": 18446744073709551615}}`, `max64`, `1`).IsEqualToUint(18446744073709551615) // Key exists
		RunTestThatExpression(t, "test_get_string_to_uint64_map__or(?, ?, ?)", `{"17": {"max64": 18446744073709551615}}`, `missing`, `1`).IsEqualToUint(1)                    // Key missing
		RunTestThatExpression(t, "test_get_string_to_uint64_map__or(?, ?, ?)", `{}`, `max64`, `1`).IsEqualToUint(1)                                                        // Map empty
	})

	t.Run("sint32_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_sint32_map(?, JSON_OBJECT('negative', -2147483648))", `{}`).IsEqualToJsonString(`{"18": {"negative": -2147483648}}`)
		RunTestThatExpression(t, "test_get_all_string_to_sint32_map(?)", `{"18": {"negative": -2147483648}}`).IsEqualToJsonString(`{"negative": -2147483648}`)
		RunTestThatExpression(t, "test_count_string_to_sint32_map(?)", `{"18": {"negative": -2147483648}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_sint32_map(?)", `{"18": {"negative": -2147483648}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_sint32_map__or(?, ?, ?)", `{"18": {"negative": -2147483648}}`, `negative`, `0`).IsEqualToInt(-2147483648) // Key exists
		RunTestThatExpression(t, "test_get_string_to_sint32_map__or(?, ?, ?)", `{"18": {"negative": -2147483648}}`, `missing`, `0`).IsEqualToInt(0)             // Key missing
		RunTestThatExpression(t, "test_get_string_to_sint32_map__or(?, ?, ?)", `{}`, `negative`, `0`).IsEqualToInt(0)                                              // Map empty
	})

	t.Run("sint64_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_sint64_map(?, JSON_OBJECT('big_negative', -9223372036854775808))", `{}`).IsEqualToJsonString(`{"19": {"big_negative": -9223372036854775808}}`)
		RunTestThatExpression(t, "test_get_all_string_to_sint64_map(?)", `{"19": {"big_negative": -9223372036854775808}}`).IsEqualToJsonString(`{"big_negative": -9223372036854775808}`)
		RunTestThatExpression(t, "test_count_string_to_sint64_map(?)", `{"19": {"big_negative": -9223372036854775808}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_sint64_map(?)", `{"19": {"big_negative": -9223372036854775808}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_sint64_map__or(?, ?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `big_negative`, `0`).IsEqualToInt(-9223372036854775808) // Key exists
		RunTestThatExpression(t, "test_get_string_to_sint64_map__or(?, ?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `missing`, `0`).IsEqualToInt(0)                     // Key missing
		RunTestThatExpression(t, "test_get_string_to_sint64_map__or(?, ?, ?)", `{}`, `big_negative`, `0`).IsEqualToInt(0)                                                            // Map empty
	})

	t.Run("fixed32_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_fixed32_map(?, JSON_OBJECT('max_fixed32', 4294967295))", `{}`).IsEqualToJsonString(`{"20": {"max_fixed32": 4294967295}}`)
		RunTestThatExpression(t, "test_get_all_string_to_fixed32_map(?)", `{"20": {"max_fixed32": 4294967295}}`).IsEqualToJsonString(`{"max_fixed32": 4294967295}`)
		RunTestThatExpression(t, "test_count_string_to_fixed32_map(?)", `{"20": {"max_fixed32": 4294967295}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_fixed32_map(?)", `{"20": {"max_fixed32": 4294967295}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_fixed32_map__or(?, ?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `max_fixed32`, `1`).IsEqualToUint(4294967295) // Key exists
		RunTestThatExpression(t, "test_get_string_to_fixed32_map__or(?, ?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `missing`, `1`).IsEqualToUint(1)           // Key missing
		RunTestThatExpression(t, "test_get_string_to_fixed32_map__or(?, ?, ?)", `{}`, `max_fixed32`, `1`).IsEqualToUint(1)                                           // Map empty
	})

	t.Run("fixed64_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_fixed64_map(?, JSON_OBJECT('max_fixed64', 18446744073709551615))", `{}`).IsEqualToJsonString(`{"21": {"max_fixed64": 18446744073709551615}}`)
		RunTestThatExpression(t, "test_get_all_string_to_fixed64_map(?)", `{"21": {"max_fixed64": 18446744073709551615}}`).IsEqualToJsonString(`{"max_fixed64": 18446744073709551615}`)
		RunTestThatExpression(t, "test_count_string_to_fixed64_map(?)", `{"21": {"max_fixed64": 18446744073709551615}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_fixed64_map(?)", `{"21": {"max_fixed64": 18446744073709551615}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_fixed64_map__or(?, ?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `max_fixed64`, `1`).IsEqualToUint(18446744073709551615) // Key exists
		RunTestThatExpression(t, "test_get_string_to_fixed64_map__or(?, ?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `missing`, `1`).IsEqualToUint(1)                     // Key missing
		RunTestThatExpression(t, "test_get_string_to_fixed64_map__or(?, ?, ?)", `{}`, `max_fixed64`, `1`).IsEqualToUint(1)                                                         // Map empty
	})

	t.Run("sfixed32_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_sfixed32_map(?, JSON_OBJECT('min_sfixed32', -2147483648))", `{}`).IsEqualToJsonString(`{"22": {"min_sfixed32": -2147483648}}`)
		RunTestThatExpression(t, "test_get_all_string_to_sfixed32_map(?)", `{"22": {"min_sfixed32": -2147483648}}`).IsEqualToJsonString(`{"min_sfixed32": -2147483648}`)
		RunTestThatExpression(t, "test_count_string_to_sfixed32_map(?)", `{"22": {"min_sfixed32": -2147483648}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_sfixed32_map(?)", `{"22": {"min_sfixed32": -2147483648}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_sfixed32_map__or(?, ?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `min_sfixed32`, `0`).IsEqualToInt(-2147483648) // Key exists
		RunTestThatExpression(t, "test_get_string_to_sfixed32_map__or(?, ?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `missing`, `0`).IsEqualToInt(0)              // Key missing
		RunTestThatExpression(t, "test_get_string_to_sfixed32_map__or(?, ?, ?)", `{}`, `min_sfixed32`, `0`).IsEqualToInt(0)                                               // Map empty
	})

	t.Run("sfixed64_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_sfixed64_map(?, JSON_OBJECT('min_sfixed64', -9223372036854775808))", `{}`).IsEqualToJsonString(`{"23": {"min_sfixed64": -9223372036854775808}}`)
		RunTestThatExpression(t, "test_get_all_string_to_sfixed64_map(?)", `{"23": {"min_sfixed64": -9223372036854775808}}`).IsEqualToJsonString(`{"min_sfixed64": -9223372036854775808}`)
		RunTestThatExpression(t, "test_count_string_to_sfixed64_map(?)", `{"23": {"min_sfixed64": -9223372036854775808}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_sfixed64_map(?)", `{"23": {"min_sfixed64": -9223372036854775808}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_sfixed64_map__or(?, ?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `min_sfixed64`, `0`).IsEqualToInt(-9223372036854775808) // Key exists
		RunTestThatExpression(t, "test_get_string_to_sfixed64_map__or(?, ?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `missing`, `0`).IsEqualToInt(0)                      // Key missing
		RunTestThatExpression(t, "test_get_string_to_sfixed64_map__or(?, ?, ?)", `{}`, `min_sfixed64`, `0`).IsEqualToInt(0)                                                             // Map empty
	})

	t.Run("bool_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_bool_map(?, JSON_OBJECT('false', false))", `{}`).IsEqualToJsonString(`{"24": {"false": false}}`)
		RunTestThatExpression(t, "test_set_all_string_to_bool_map(?, JSON_OBJECT('flag', true))", `{}`).IsEqualToJsonString(`{"24": {"flag": true}}`)
		RunTestThatExpression(t, "test_get_all_string_to_bool_map(?)", `{"24": {"flag": true, "other": false}}`).IsEqualToJsonString(`{"flag": true, "other": false}`)
		RunTestThatExpression(t, "test_count_string_to_bool_map(?)", `{"24": {"flag": true}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_bool_map(?)", `{"24": {"flag": true}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_bool_map__or(?, ?, ?)", `{"24": {"flag": true, "other": false}}`, `flag`, false).IsEqualToBool(true)   // Key exists
		RunTestThatExpression(t, "test_get_string_to_bool_map__or(?, ?, ?)", `{"24": {"flag": true, "other": false}}`, `other`, false).IsEqualToBool(false) // Key exists (zero value)
		RunTestThatExpression(t, "test_get_string_to_bool_map__or(?, ?, ?)", `{"24": {"flag": true, "other": false}}`, `missing`, false).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "test_get_string_to_bool_map__or(?, ?, ?)", `{}`, `flag`, false).IsEqualToBool(false)                                          // Map empty
	})

	t.Run("string_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_string_map(?, JSON_OBJECT('empty', ''))", `{}`).IsEqualToJsonString(`{"25": {"empty": ""}}`)
		RunTestThatExpression(t, "test_set_all_string_to_string_map(?, JSON_OBJECT('greeting', 'hello'))", `{}`).IsEqualToJsonString(`{"25": {"greeting": "hello"}}`)
		RunTestThatExpression(t, "test_get_all_string_to_string_map(?)", `{"25": {"greeting": "hello"}}`).IsEqualToJsonString(`{"greeting": "hello"}`)
		RunTestThatExpression(t, "test_count_string_to_string_map(?)", `{"25": {"greeting": "hello"}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_string_map(?)", `{"25": {"greeting": "hello"}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_string_map__or(?, ?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `greeting`, ``).IsEqualToString(`hello`) // Key exists
		RunTestThatExpression(t, "test_get_string_to_string_map__or(?, ?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `empty`, ``).IsEqualToString(``)      // Key exists (zero value)
		RunTestThatExpression(t, "test_get_string_to_string_map__or(?, ?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `missing`, ``).IsEqualToString(``)   // Key missing
		RunTestThatExpression(t, "test_get_string_to_string_map__or(?, ?, ?)", `{}`, `greeting`, ``).IsEqualToString(``)                                                // Map empty
	})

	t.Run("bytes_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_bytes_map(?, JSON_OBJECT('empty', ''))", `{}`).IsEqualToJsonString(`{"26": {"empty": ""}}`)
		RunTestThatExpression(t, "test_set_all_string_to_bytes_map(?, JSON_OBJECT('data', 'aGVsbG8='))", `{}`).IsEqualToJsonString(`{"26": {"data": "aGVsbG8="}}`)
		RunTestThatExpression(t, "test_get_all_string_to_bytes_map(?)", `{"26": {"data": "aGVsbG8="}}`).IsEqualToJsonString(`{"data": "aGVsbG8="}`)
		RunTestThatExpression(t, "test_count_string_to_bytes_map(?)", `{"26": {"data": "aGVsbG8="}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_bytes_map(?)", `{"26": {"data": "aGVsbG8="}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_bytes_map__or(?, ?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `data`, ``).IsEqualToString(`hello`)  // Key exists (base64 decoded)
		RunTestThatExpression(t, "test_get_string_to_bytes_map__or(?, ?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `empty`, ``).IsEqualToString(``)        // Key exists (zero value)
		RunTestThatExpression(t, "test_get_string_to_bytes_map__or(?, ?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `missing`, ``).IsEqualToString(``)     // Key missing
		RunTestThatExpression(t, "test_get_string_to_bytes_map__or(?, ?, ?)", `{}`, `data`, ``).IsEqualToString(``)                                                  // Map empty
	})

	t.Run("enum_value", func(t *testing.T) {
		RunTestThatExpression(t, "test_set_all_string_to_enum_map(?, JSON_OBJECT('status', 1))", `{}`).IsEqualToJsonString(`{"27": {"status": 1}}`)
		RunTestThatExpression(t, "test_get_all_string_to_enum_map(?)", `{"27": {"status": 1}}`).IsEqualToJsonString(`{"status": 1}`)
		RunTestThatExpression(t, "test_count_string_to_enum_map(?)", `{"27": {"status": 1}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_enum_map(?)", `{"27": {"status": 1}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "test_get_string_to_enum_map__or(?, ?, ?)", `{"27": {"status": 1}}`, `status`, `0`).IsEqualToInt(1) // Key exists
		RunTestThatExpression(t, "test_get_string_to_enum_map__or(?, ?, ?)", `{"27": {"status": 1}}`, `missing`, `0`).IsEqualToInt(0) // Key missing
		RunTestThatExpression(t, "test_get_string_to_enum_map__or(?, ?, ?)", `{}`, `status`, `0`).IsEqualToInt(0)                      // Map empty
	})

	t.Run("message_value", func(t *testing.T) {
		// Test message value
		RunTestThatExpression(t, "test_set_all_string_to_message_map(?, JSON_OBJECT('nested', JSON_OBJECT('1', 'test', '2', 42)))", `{}`).IsEqualToJsonString(`{"28": {"nested": {"1": "test", "2": 42}}}`)
		RunTestThatExpression(t, "test_get_all_string_to_message_map(?)", `{"28": {"nested": {"1": "test", "2": 42}}}`).IsEqualToJsonString(`{"nested": {"1": "test", "2": 42}}`)
		RunTestThatExpression(t, "test_count_string_to_message_map(?)", `{"28": {"nested": {"1": "test", "2": 42}}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "test_clear_string_to_message_map(?)", `{"28": {"nested": {"1": "test", "2": 42}}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		defaultMessage := `{"1": "default", "2": 0}`
		RunTestThatExpression(t, "test_get_string_to_message_map__or(?, ?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `nested`, defaultMessage).IsEqualToJsonString(`{"1": "test", "2": 42}`) // Key exists
		RunTestThatExpression(t, "test_get_string_to_message_map__or(?, ?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `missing`, defaultMessage).IsEqualToJsonString(defaultMessage)     // Key missing
		RunTestThatExpression(t, "test_get_string_to_message_map__or(?, ?, ?)", `{}`, `nested`, defaultMessage).IsEqualToJsonString(defaultMessage)                                      // Map empty
	})

}
