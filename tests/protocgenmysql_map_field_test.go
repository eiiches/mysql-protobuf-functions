package main

import (
	"testing"
)

func TestProtocGenMapField(t *testing.T) {
	// Test all protobuf map field types using the pre-generated functions from protocgenmysql_proto3.pb.sql
	// The MapFields message in protocgenmysql_proto3.proto matches the test schema

	// Test integer key map types
	t.Run("int32_key", func(t *testing.T) {
		// Test setters create correct internal format
		RunTestThatExpression(t, "pbt_map_fields_set_all_int32_to_int32_map(?, JSON_OBJECT('42', 100))", `{}`).IsEqualToJsonString(`{"1": {"42": 100}}`)
		RunTestThatExpression(t, "pbt_map_fields_set_all_int32_to_int32_map(?, JSON_OBJECT('0', 0))", `{}`).IsEqualToJsonString(`{"1": {"0": 0}}`) // Zero key and value stored

		// Test getters return entire map
		RunTestThatExpression(t, "pbt_map_fields_get_all_int32_to_int32_map(?)", `{"1": {"42": 100, "1": 10}}`).IsEqualToJsonString(`{"42": 100, "1": 10}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_int32_to_int32_map(?)", `{}`).IsEqualToJsonString(`{}`) // Default when absent

		// Test map count operations
		RunTestThatExpression(t, "pbt_map_fields_count_int32_to_int32_map(?)", `{}`).IsEqualToInt(0)
		RunTestThatExpression(t, "pbt_map_fields_count_int32_to_int32_map(?)", `{"1": {"42": 100, "1": 10}}`).IsEqualToInt(2)

		// Test clear methods
		RunTestThatExpression(t, "pbt_map_fields_clear_int32_to_int32_map(?)", `{"1": {"42": 100}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_int32_to_int32_map__or(?, ?, ?)", `{"1": {"42": 100}}`, `42`, `999`).IsEqualToInt(100) // Key exists, return value
		RunTestThatExpression(t, "pbt_map_fields_get_int32_to_int32_map__or(?, ?, ?)", `{"1": {"42": 100}}`, `99`, `999`).IsEqualToInt(999) // Key missing, return default
		RunTestThatExpression(t, "pbt_map_fields_get_int32_to_int32_map__or(?, ?, ?)", `{}`, `42`, `999`).IsEqualToInt(999)                 // Map empty, return default

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_int32_to_int32_map(?, ?)", `{"1": {"42": 100}}`, `42`).IsEqualToInt(100) // Key exists, return value
		RunTestThatExpression(t, "pbt_map_fields_get_int32_to_int32_map(?, ?)", `{"1": {"42": 100}}`, `99`).IsNull()   // Key missing, return NULL
		RunTestThatExpression(t, "pbt_map_fields_get_int32_to_int32_map(?, ?)", `{}`, `42`).IsNull()                   // Map empty, return NULL

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_int32_to_int32_map(?, ?)", `{"1": {"42": 100}}`, `42`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_int32_to_int32_map(?, ?)", `{"1": {"42": 100}}`, `99`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_int32_to_int32_map(?, ?)", `{}`, `42`).IsEqualToBool(false)                 // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_int32_to_int32_map(?, ?, ?)", `{}`, `42`, `100`).IsEqualToJsonString(`{"1": {"42": 100}}`)                          // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_int32_to_int32_map(?, ?, ?)", `{"1": {"10": 20}}`, `42`, `100`).IsEqualToJsonString(`{"1": {"10": 20, "42": 100}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_int32_to_int32_map(?, ?, ?)", `{"1": {"42": 50}}`, `42`, `100`).IsEqualToJsonString(`{"1": {"42": 100}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_int32_to_int32_map(?, ?)", `{}`, `{"30": 300, "40": 400}`).IsEqualToJsonString(`{"1": {"30": 300, "40": 400}}`)                                     // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_int32_to_int32_map(?, ?)", `{"1": {"10": 20}}`, `{"30": 300, "40": 400}`).IsEqualToJsonString(`{"1": {"10": 20, "30": 300, "40": 400}}`)            // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_int32_to_int32_map(?, ?)", `{"1": {"10": 20, "30": 250}}`, `{"30": 300, "40": 400}`).IsEqualToJsonString(`{"1": {"10": 20, "30": 300, "40": 400}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_int32_to_int32_map(?, ?)", `{"1": {"42": 100, "10": 20}}`, `42`).IsEqualToJsonString(`{"1": {"10": 20}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_int32_to_int32_map(?, ?)", `{"1": {"42": 100}}`, `42`).IsEqualToJsonString(`{}`)                          // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_int32_to_int32_map(?, ?)", `{"1": {"42": 100}}`, `99`).IsEqualToJsonString(`{"1": {"42": 100}}`)          // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_int32_to_int32_map(?, ?)", `{}`, `42`).IsEqualToJsonString(`{}`)
	})

	t.Run("int64_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_int64_to_int32_map(?, JSON_OBJECT('9223372036854775807', 200))", `{}`).IsEqualToJsonString(`{"2": {"9223372036854775807": 200}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_int64_to_int32_map(?)", `{"2": {"9223372036854775807": 200}}`).IsEqualToJsonString(`{"9223372036854775807": 200}`)
		RunTestThatExpression(t, "pbt_map_fields_count_int64_to_int32_map(?)", `{"2": {"9223372036854775807": 200}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_int64_to_int32_map(?)", `{"2": {"9223372036854775807": 200}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_int64_to_int32_map__or(?, ?, ?)", `{"2": {"9223372036854775807": 200}}`, `9223372036854775807`, `888`).IsEqualToInt(200) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_int64_to_int32_map__or(?, ?, ?)", `{"2": {"9223372036854775807": 200}}`, `123`, `888`).IsEqualToInt(888)                 // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_int64_to_int32_map__or(?, ?, ?)", `{}`, `9223372036854775807`, `888`).IsEqualToInt(888)                                  // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_int64_to_int32_map(?, ?)", `{"2": {"9223372036854775807": 200}}`, `9223372036854775807`).IsEqualToInt(200) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_int64_to_int32_map(?, ?)", `{"2": {"9223372036854775807": 200}}`, `123`).IsNull()                   // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_int64_to_int32_map(?, ?)", `{}`, `9223372036854775807`).IsNull()                                    // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_int64_to_int32_map(?, ?)", `{"2": {"9223372036854775807": 200}}`, `9223372036854775807`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_int64_to_int32_map(?, ?)", `{"2": {"9223372036854775807": 200}}`, `123`).IsEqualToBool(false)                // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_int64_to_int32_map(?, ?)", `{}`, `9223372036854775807`).IsEqualToBool(false)                                 // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_int64_to_int32_map(?, ?, ?)", `{}`, `9223372036854775807`, `200`).IsEqualToJsonString(`{"2": {"9223372036854775807": 200}}`)                                  // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_int64_to_int32_map(?, ?, ?)", `{"2": {"100": 50}}`, `9223372036854775807`, `200`).IsEqualToJsonString(`{"2": {"100": 50, "9223372036854775807": 200}}`)       // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_int64_to_int32_map(?, ?, ?)", `{"2": {"9223372036854775807": 150}}`, `9223372036854775807`, `200`).IsEqualToJsonString(`{"2": {"9223372036854775807": 200}}`) // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_int64_to_int32_map(?, ?)", `{}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"2": {"300": 333, "400": 444}}`)                                        // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_int64_to_int32_map(?, ?)", `{"2": {"100": 50}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"2": {"100": 50, "300": 333, "400": 444}}`)             // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_int64_to_int32_map(?, ?)", `{"2": {"100": 50, "300": 250}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"2": {"100": 50, "300": 333, "400": 444}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_int64_to_int32_map(?, ?)", `{"2": {"9223372036854775807": 200, "100": 50}}`, `9223372036854775807`).IsEqualToJsonString(`{"2": {"100": 50}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_int64_to_int32_map(?, ?)", `{"2": {"9223372036854775807": 200}}`, `9223372036854775807`).IsEqualToJsonString(`{}`)                            // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_int64_to_int32_map(?, ?)", `{"2": {"9223372036854775807": 200}}`, `123`).IsEqualToJsonString(`{"2": {"9223372036854775807": 200}}`)           // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_int64_to_int32_map(?, ?)", `{}`, `9223372036854775807`).IsEqualToJsonString(`{}`)                                                             // Remove from empty map
	})

	t.Run("uint32_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_uint32_to_int32_map(?, JSON_OBJECT('4294967295', 300))", `{}`).IsEqualToJsonString(`{"3": {"4294967295": 300}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_uint32_to_int32_map(?)", `{"3": {"4294967295": 300}}`).IsEqualToJsonString(`{"4294967295": 300}`)
		RunTestThatExpression(t, "pbt_map_fields_count_uint32_to_int32_map(?)", `{"3": {"4294967295": 300}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_uint32_to_int32_map(?)", `{"3": {"4294967295": 300}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_uint32_to_int32_map__or(?, ?, ?)", `{"3": {"4294967295": 300}}`, `4294967295`, `777`).IsEqualToInt(300) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_uint32_to_int32_map__or(?, ?, ?)", `{"3": {"4294967295": 300}}`, `123`, `777`).IsEqualToInt(777)        // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_uint32_to_int32_map__or(?, ?, ?)", `{}`, `4294967295`, `777`).IsEqualToInt(777)                         // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_uint32_to_int32_map(?, ?)", `{"3": {"4294967295": 300}}`, `4294967295`).IsEqualToInt(300) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_uint32_to_int32_map(?, ?)", `{"3": {"4294967295": 300}}`, `123`).IsNull()          // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_uint32_to_int32_map(?, ?)", `{}`, `4294967295`).IsNull()                           // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_uint32_to_int32_map(?, ?)", `{"3": {"4294967295": 300}}`, `4294967295`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_uint32_to_int32_map(?, ?)", `{"3": {"4294967295": 300}}`, `123`).IsEqualToBool(false)       // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_uint32_to_int32_map(?, ?)", `{}`, `4294967295`).IsEqualToBool(false)                        // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_uint32_to_int32_map(?, ?, ?)", `{}`, `4294967295`, `300`).IsEqualToJsonString(`{"3": {"4294967295": 300}}`)                            // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_uint32_to_int32_map(?, ?, ?)", `{"3": {"100": 50}}`, `4294967295`, `300`).IsEqualToJsonString(`{"3": {"100": 50, "4294967295": 300}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_uint32_to_int32_map(?, ?, ?)", `{"3": {"4294967295": 250}}`, `4294967295`, `300`).IsEqualToJsonString(`{"3": {"4294967295": 300}}`)    // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_uint32_to_int32_map(?, ?)", `{}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"3": {"300": 333, "400": 444}}`)                                        // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_uint32_to_int32_map(?, ?)", `{"3": {"100": 50}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"3": {"100": 50, "300": 333, "400": 444}}`)             // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_uint32_to_int32_map(?, ?)", `{"3": {"100": 50, "300": 250}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"3": {"100": 50, "300": 333, "400": 444}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_uint32_to_int32_map(?, ?)", `{"3": {"4294967295": 300, "100": 50}}`, `4294967295`).IsEqualToJsonString(`{"3": {"100": 50}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_uint32_to_int32_map(?, ?)", `{"3": {"4294967295": 300}}`, `4294967295`).IsEqualToJsonString(`{}`)                            // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_uint32_to_int32_map(?, ?)", `{"3": {"4294967295": 300}}`, `123`).IsEqualToJsonString(`{"3": {"4294967295": 300}}`)           // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_uint32_to_int32_map(?, ?)", `{}`, `4294967295`).IsEqualToJsonString(`{}`)                                                    // Remove from empty map
	})

	t.Run("uint64_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_uint64_to_int32_map(?, JSON_OBJECT('18446744073709551615', 400))", `{}`).IsEqualToJsonString(`{"4": {"18446744073709551615": 400}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_uint64_to_int32_map(?)", `{"4": {"18446744073709551615": 400}}`).IsEqualToJsonString(`{"18446744073709551615": 400}`)
		RunTestThatExpression(t, "pbt_map_fields_count_uint64_to_int32_map(?)", `{"4": {"18446744073709551615": 400}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_uint64_to_int32_map(?)", `{"4": {"18446744073709551615": 400}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_uint64_to_int32_map__or(?, ?, ?)", `{"4": {"18446744073709551615": 400}}`, `18446744073709551615`, `666`).IsEqualToInt(400) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_uint64_to_int32_map__or(?, ?, ?)", `{"4": {"18446744073709551615": 400}}`, `123`, `666`).IsEqualToInt(666)                  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_uint64_to_int32_map__or(?, ?, ?)", `{}`, `18446744073709551615`, `666`).IsEqualToInt(666)                                   // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_uint64_to_int32_map(?, ?)", `{"4": {"18446744073709551615": 400}}`, `18446744073709551615`).IsEqualToInt(400) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_uint64_to_int32_map(?, ?)", `{"4": {"18446744073709551615": 400}}`, `123`).IsNull()                    // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_uint64_to_int32_map(?, ?)", `{}`, `18446744073709551615`).IsNull()                                     // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_uint64_to_int32_map(?, ?)", `{"4": {"18446744073709551615": 400}}`, `18446744073709551615`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_uint64_to_int32_map(?, ?)", `{"4": {"18446744073709551615": 400}}`, `123`).IsEqualToBool(false)                 // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_uint64_to_int32_map(?, ?)", `{}`, `18446744073709551615`).IsEqualToBool(false)                                  // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_uint64_to_int32_map(?, ?, ?)", `{}`, `18446744073709551615`, `400`).IsEqualToJsonString(`{"4": {"18446744073709551615": 400}}`)                                   // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_uint64_to_int32_map(?, ?, ?)", `{"4": {"100": 50}}`, `18446744073709551615`, `400`).IsEqualToJsonString(`{"4": {"100": 50, "18446744073709551615": 400}}`)        // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_uint64_to_int32_map(?, ?, ?)", `{"4": {"18446744073709551615": 350}}`, `18446744073709551615`, `400`).IsEqualToJsonString(`{"4": {"18446744073709551615": 400}}`) // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_uint64_to_int32_map(?, ?)", `{}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"4": {"300": 333, "400": 444}}`)                                        // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_uint64_to_int32_map(?, ?)", `{"4": {"100": 50}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"4": {"100": 50, "300": 333, "400": 444}}`)             // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_uint64_to_int32_map(?, ?)", `{"4": {"100": 50, "300": 250}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"4": {"100": 50, "300": 333, "400": 444}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_uint64_to_int32_map(?, ?)", `{"4": {"18446744073709551615": 400, "100": 50}}`, `18446744073709551615`).IsEqualToJsonString(`{"4": {"100": 50}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_uint64_to_int32_map(?, ?)", `{"4": {"18446744073709551615": 400}}`, `18446744073709551615`).IsEqualToJsonString(`{}`)                            // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_uint64_to_int32_map(?, ?)", `{"4": {"18446744073709551615": 400}}`, `123`).IsEqualToJsonString(`{"4": {"18446744073709551615": 400}}`)           // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_uint64_to_int32_map(?, ?)", `{}`, `18446744073709551615`).IsEqualToJsonString(`{}`)                                                              // Remove from empty map
	})

	t.Run("sint32_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_sint32_to_int32_map(?, JSON_OBJECT('-1', 500))", `{}`).IsEqualToJsonString(`{"5": {"-1": 500}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_sint32_to_int32_map(?)", `{"5": {"-1": 500}}`).IsEqualToJsonString(`{"-1": 500}`)
		RunTestThatExpression(t, "pbt_map_fields_count_sint32_to_int32_map(?)", `{"5": {"-1": 500}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_sint32_to_int32_map(?)", `{"5": {"-1": 500}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_sint32_to_int32_map__or(?, ?, ?)", `{"5": {"-2147483648": 500}}`, `-2147483648`, `555`).IsEqualToInt(500) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_sint32_to_int32_map__or(?, ?, ?)", `{"5": {"-2147483648": 500}}`, `123`, `555`).IsEqualToInt(555)         // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_sint32_to_int32_map__or(?, ?, ?)", `{}`, `-2147483648`, `555`).IsEqualToInt(555)                          // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_sint32_to_int32_map(?, ?)", `{"5": {"-2147483648": 500}}`, `-2147483648`).IsEqualToInt(500) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_sint32_to_int32_map(?, ?)", `{"5": {"-2147483648": 500}}`, `123`).IsNull()           // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_sint32_to_int32_map(?, ?)", `{}`, `-2147483648`).IsNull()                            // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_sint32_to_int32_map(?, ?)", `{"5": {"-2147483648": 500}}`, `-2147483648`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_sint32_to_int32_map(?, ?)", `{"5": {"-2147483648": 500}}`, `123`).IsEqualToBool(false)        // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_sint32_to_int32_map(?, ?)", `{}`, `-2147483648`).IsEqualToBool(false)                         // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_sint32_to_int32_map(?, ?, ?)", `{}`, `-2147483648`, `500`).IsEqualToJsonString(`{"5": {"-2147483648": 500}}`)                            // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_sint32_to_int32_map(?, ?, ?)", `{"5": {"100": 50}}`, `-2147483648`, `500`).IsEqualToJsonString(`{"5": {"100": 50, "-2147483648": 500}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_sint32_to_int32_map(?, ?, ?)", `{"5": {"-2147483648": 450}}`, `-2147483648`, `500`).IsEqualToJsonString(`{"5": {"-2147483648": 500}}`)   // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_sint32_to_int32_map(?, ?)", `{}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"5": {"-300": 333, "-400": 444}}`)                                         // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_sint32_to_int32_map(?, ?)", `{"5": {"100": 50}}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"5": {"100": 50, "-300": 333, "-400": 444}}`)              // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_sint32_to_int32_map(?, ?)", `{"5": {"100": 50, "-300": 250}}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"5": {"100": 50, "-300": 333, "-400": 444}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_sint32_to_int32_map(?, ?)", `{"5": {"-2147483648": 500, "100": 50}}`, `-2147483648`).IsEqualToJsonString(`{"5": {"100": 50}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_sint32_to_int32_map(?, ?)", `{"5": {"-2147483648": 500}}`, `-2147483648`).IsEqualToJsonString(`{}`)                            // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_sint32_to_int32_map(?, ?)", `{"5": {"-2147483648": 500}}`, `123`).IsEqualToJsonString(`{"5": {"-2147483648": 500}}`)           // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_sint32_to_int32_map(?, ?)", `{}`, `-2147483648`).IsEqualToJsonString(`{}`)                                                     // Remove from empty map
	})

	t.Run("sint64_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_sint64_to_int32_map(?, JSON_OBJECT('-9223372036854775808', 600))", `{}`).IsEqualToJsonString(`{"6": {"-9223372036854775808": 600}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_sint64_to_int32_map(?)", `{"6": {"-9223372036854775808": 600}}`).IsEqualToJsonString(`{"-9223372036854775808": 600}`)
		RunTestThatExpression(t, "pbt_map_fields_count_sint64_to_int32_map(?)", `{"6": {"-9223372036854775808": 600}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_sint64_to_int32_map(?)", `{"6": {"-9223372036854775808": 600}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_sint64_to_int32_map__or(?, ?, ?)", `{"6": {"-9223372036854775808": 600}}`, `-9223372036854775808`, `444`).IsEqualToInt(600) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_sint64_to_int32_map__or(?, ?, ?)", `{"6": {"-9223372036854775808": 600}}`, `123`, `444`).IsEqualToInt(444)                  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_sint64_to_int32_map__or(?, ?, ?)", `{}`, `-9223372036854775808`, `444`).IsEqualToInt(444)                                   // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_sint64_to_int32_map(?, ?)", `{"6": {"-9223372036854775808": 600}}`, `-9223372036854775808`).IsEqualToInt(600) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_sint64_to_int32_map(?, ?)", `{"6": {"-9223372036854775808": 600}}`, `123`).IsNull()                    // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_sint64_to_int32_map(?, ?)", `{}`, `-9223372036854775808`).IsNull()                                     // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_sint64_to_int32_map(?, ?)", `{"6": {"-9223372036854775808": 600}}`, `-9223372036854775808`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_sint64_to_int32_map(?, ?)", `{"6": {"-9223372036854775808": 600}}`, `123`).IsEqualToBool(false)                 // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_sint64_to_int32_map(?, ?)", `{}`, `-9223372036854775808`).IsEqualToBool(false)                                  // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_sint64_to_int32_map(?, ?, ?)", `{}`, `-9223372036854775808`, `600`).IsEqualToJsonString(`{"6": {"-9223372036854775808": 600}}`)                                   // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_sint64_to_int32_map(?, ?, ?)", `{"6": {"100": 50}}`, `-9223372036854775808`, `600`).IsEqualToJsonString(`{"6": {"100": 50, "-9223372036854775808": 600}}`)        // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_sint64_to_int32_map(?, ?, ?)", `{"6": {"-9223372036854775808": 550}}`, `-9223372036854775808`, `600`).IsEqualToJsonString(`{"6": {"-9223372036854775808": 600}}`) // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_sint64_to_int32_map(?, ?)", `{}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"6": {"-300": 333, "-400": 444}}`)                                         // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_sint64_to_int32_map(?, ?)", `{"6": {"100": 50}}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"6": {"100": 50, "-300": 333, "-400": 444}}`)              // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_sint64_to_int32_map(?, ?)", `{"6": {"100": 50, "-300": 250}}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"6": {"100": 50, "-300": 333, "-400": 444}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_sint64_to_int32_map(?, ?)", `{"6": {"-9223372036854775808": 600, "100": 50}}`, `-9223372036854775808`).IsEqualToJsonString(`{"6": {"100": 50}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_sint64_to_int32_map(?, ?)", `{"6": {"-9223372036854775808": 600}}`, `-9223372036854775808`).IsEqualToJsonString(`{}`)                            // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_sint64_to_int32_map(?, ?)", `{"6": {"-9223372036854775808": 600}}`, `123`).IsEqualToJsonString(`{"6": {"-9223372036854775808": 600}}`)           // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_sint64_to_int32_map(?, ?)", `{}`, `-9223372036854775808`).IsEqualToJsonString(`{}`)                                                              // Remove from empty map
	})

	t.Run("fixed32_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_fixed32_to_int32_map(?, JSON_OBJECT('4294967295', 700))", `{}`).IsEqualToJsonString(`{"7": {"4294967295": 700}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_fixed32_to_int32_map(?)", `{"7": {"4294967295": 700}}`).IsEqualToJsonString(`{"4294967295": 700}`)
		RunTestThatExpression(t, "pbt_map_fields_count_fixed32_to_int32_map(?)", `{"7": {"4294967295": 700}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_fixed32_to_int32_map(?)", `{"7": {"4294967295": 700}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_fixed32_to_int32_map__or(?, ?, ?)", `{"7": {"4294967295": 700}}`, `4294967295`, `333`).IsEqualToInt(700) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_fixed32_to_int32_map__or(?, ?, ?)", `{"7": {"4294967295": 700}}`, `123`, `333`).IsEqualToInt(333)        // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_fixed32_to_int32_map__or(?, ?, ?)", `{}`, `4294967295`, `333`).IsEqualToInt(333)                         // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_fixed32_to_int32_map(?, ?)", `{"7": {"4294967295": 700}}`, `4294967295`).IsEqualToInt(700) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_fixed32_to_int32_map(?, ?)", `{"7": {"4294967295": 700}}`, `123`).IsNull()          // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_fixed32_to_int32_map(?, ?)", `{}`, `4294967295`).IsNull()                           // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_fixed32_to_int32_map(?, ?)", `{"7": {"4294967295": 700}}`, `4294967295`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_fixed32_to_int32_map(?, ?)", `{"7": {"4294967295": 700}}`, `123`).IsEqualToBool(false)       // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_fixed32_to_int32_map(?, ?)", `{}`, `4294967295`).IsEqualToBool(false)                        // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_fixed32_to_int32_map(?, ?, ?)", `{}`, `4294967295`, `700`).IsEqualToJsonString(`{"7": {"4294967295": 700}}`)                            // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_fixed32_to_int32_map(?, ?, ?)", `{"7": {"100": 50}}`, `4294967295`, `700`).IsEqualToJsonString(`{"7": {"100": 50, "4294967295": 700}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_fixed32_to_int32_map(?, ?, ?)", `{"7": {"4294967295": 650}}`, `4294967295`, `700`).IsEqualToJsonString(`{"7": {"4294967295": 700}}`)    // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_fixed32_to_int32_map(?, ?)", `{}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"7": {"300": 333, "400": 444}}`)                                        // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_fixed32_to_int32_map(?, ?)", `{"7": {"100": 50}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"7": {"100": 50, "300": 333, "400": 444}}`)             // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_fixed32_to_int32_map(?, ?)", `{"7": {"100": 50, "300": 250}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"7": {"100": 50, "300": 333, "400": 444}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_fixed32_to_int32_map(?, ?)", `{"7": {"4294967295": 700, "100": 50}}`, `4294967295`).IsEqualToJsonString(`{"7": {"100": 50}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_fixed32_to_int32_map(?, ?)", `{"7": {"4294967295": 700}}`, `4294967295`).IsEqualToJsonString(`{}`)                            // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_fixed32_to_int32_map(?, ?)", `{"7": {"4294967295": 700}}`, `123`).IsEqualToJsonString(`{"7": {"4294967295": 700}}`)           // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_fixed32_to_int32_map(?, ?)", `{}`, `4294967295`).IsEqualToJsonString(`{}`)                                                    // Remove from empty map
	})

	t.Run("fixed64_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_fixed64_to_int32_map(?, JSON_OBJECT('18446744073709551615', 800))", `{}`).IsEqualToJsonString(`{"8": {"18446744073709551615": 800}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_fixed64_to_int32_map(?)", `{"8": {"18446744073709551615": 800}}`).IsEqualToJsonString(`{"18446744073709551615": 800}`)
		RunTestThatExpression(t, "pbt_map_fields_count_fixed64_to_int32_map(?)", `{"8": {"18446744073709551615": 800}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_fixed64_to_int32_map(?)", `{"8": {"18446744073709551615": 800}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_fixed64_to_int32_map__or(?, ?, ?)", `{"8": {"18446744073709551615": 800}}`, `18446744073709551615`, `222`).IsEqualToInt(800) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_fixed64_to_int32_map__or(?, ?, ?)", `{"8": {"18446744073709551615": 800}}`, `123`, `222`).IsEqualToInt(222)                  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_fixed64_to_int32_map__or(?, ?, ?)", `{}`, `18446744073709551615`, `222`).IsEqualToInt(222)                                   // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_fixed64_to_int32_map(?, ?)", `{"8": {"18446744073709551615": 800}}`, `18446744073709551615`).IsEqualToInt(800) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_fixed64_to_int32_map(?, ?)", `{"8": {"18446744073709551615": 800}}`, `123`).IsNull()                    // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_fixed64_to_int32_map(?, ?)", `{}`, `18446744073709551615`).IsNull()                                     // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_fixed64_to_int32_map(?, ?)", `{"8": {"18446744073709551615": 800}}`, `18446744073709551615`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_fixed64_to_int32_map(?, ?)", `{"8": {"18446744073709551615": 800}}`, `123`).IsEqualToBool(false)                 // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_fixed64_to_int32_map(?, ?)", `{}`, `18446744073709551615`).IsEqualToBool(false)                                  // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_fixed64_to_int32_map(?, ?, ?)", `{}`, `18446744073709551615`, `800`).IsEqualToJsonString(`{"8": {"18446744073709551615": 800}}`)                                   // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_fixed64_to_int32_map(?, ?, ?)", `{"8": {"100": 50}}`, `18446744073709551615`, `800`).IsEqualToJsonString(`{"8": {"100": 50, "18446744073709551615": 800}}`)        // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_fixed64_to_int32_map(?, ?, ?)", `{"8": {"18446744073709551615": 750}}`, `18446744073709551615`, `800`).IsEqualToJsonString(`{"8": {"18446744073709551615": 800}}`) // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_fixed64_to_int32_map(?, ?)", `{}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"8": {"300": 333, "400": 444}}`)                                        // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_fixed64_to_int32_map(?, ?)", `{"8": {"100": 50}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"8": {"100": 50, "300": 333, "400": 444}}`)             // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_fixed64_to_int32_map(?, ?)", `{"8": {"100": 50, "300": 250}}`, `{"300": 333, "400": 444}`).IsEqualToJsonString(`{"8": {"100": 50, "300": 333, "400": 444}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_fixed64_to_int32_map(?, ?)", `{"8": {"18446744073709551615": 800, "100": 50}}`, `18446744073709551615`).IsEqualToJsonString(`{"8": {"100": 50}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_fixed64_to_int32_map(?, ?)", `{"8": {"18446744073709551615": 800}}`, `18446744073709551615`).IsEqualToJsonString(`{}`)                            // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_fixed64_to_int32_map(?, ?)", `{"8": {"18446744073709551615": 800}}`, `123`).IsEqualToJsonString(`{"8": {"18446744073709551615": 800}}`)           // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_fixed64_to_int32_map(?, ?)", `{}`, `18446744073709551615`).IsEqualToJsonString(`{}`)                                                              // Remove from empty map
	})

	t.Run("sfixed32_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_sfixed32_to_int32_map(?, JSON_OBJECT('-2147483648', 900))", `{}`).IsEqualToJsonString(`{"9": {"-2147483648": 900}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_sfixed32_to_int32_map(?)", `{"9": {"-2147483648": 900}}`).IsEqualToJsonString(`{"-2147483648": 900}`)
		RunTestThatExpression(t, "pbt_map_fields_count_sfixed32_to_int32_map(?)", `{"9": {"-2147483648": 900}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_sfixed32_to_int32_map(?)", `{"9": {"-2147483648": 900}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed32_to_int32_map__or(?, ?, ?)", `{"9": {"-2147483648": 900}}`, `-2147483648`, `111`).IsEqualToInt(900) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed32_to_int32_map__or(?, ?, ?)", `{"9": {"-2147483648": 900}}`, `123`, `111`).IsEqualToInt(111)         // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed32_to_int32_map__or(?, ?, ?)", `{}`, `-2147483648`, `111`).IsEqualToInt(111)                          // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed32_to_int32_map(?, ?)", `{"9": {"-2147483648": 900}}`, `-2147483648`).IsEqualToInt(900) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed32_to_int32_map(?, ?)", `{"9": {"-2147483648": 900}}`, `123`).IsNull()           // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed32_to_int32_map(?, ?)", `{}`, `-2147483648`).IsNull()                            // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_sfixed32_to_int32_map(?, ?)", `{"9": {"-2147483648": 900}}`, `-2147483648`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_sfixed32_to_int32_map(?, ?)", `{"9": {"-2147483648": 900}}`, `123`).IsEqualToBool(false)        // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_sfixed32_to_int32_map(?, ?)", `{}`, `-2147483648`).IsEqualToBool(false)                         // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_sfixed32_to_int32_map(?, ?, ?)", `{}`, `-2147483648`, `900`).IsEqualToJsonString(`{"9": {"-2147483648": 900}}`)                            // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_sfixed32_to_int32_map(?, ?, ?)", `{"9": {"100": 50}}`, `-2147483648`, `900`).IsEqualToJsonString(`{"9": {"100": 50, "-2147483648": 900}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_sfixed32_to_int32_map(?, ?, ?)", `{"9": {"-2147483648": 850}}`, `-2147483648`, `900`).IsEqualToJsonString(`{"9": {"-2147483648": 900}}`)   // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_sfixed32_to_int32_map(?, ?)", `{}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"9": {"-300": 333, "-400": 444}}`)                                         // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_sfixed32_to_int32_map(?, ?)", `{"9": {"100": 50}}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"9": {"100": 50, "-300": 333, "-400": 444}}`)              // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_sfixed32_to_int32_map(?, ?)", `{"9": {"100": 50, "-300": 250}}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"9": {"100": 50, "-300": 333, "-400": 444}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_sfixed32_to_int32_map(?, ?)", `{"9": {"-2147483648": 900, "100": 50}}`, `-2147483648`).IsEqualToJsonString(`{"9": {"100": 50}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_sfixed32_to_int32_map(?, ?)", `{"9": {"-2147483648": 900}}`, `-2147483648`).IsEqualToJsonString(`{}`)                            // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_sfixed32_to_int32_map(?, ?)", `{"9": {"-2147483648": 900}}`, `123`).IsEqualToJsonString(`{"9": {"-2147483648": 900}}`)           // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_sfixed32_to_int32_map(?, ?)", `{}`, `-2147483648`).IsEqualToJsonString(`{}`)                                                     // Remove from empty map
	})

	t.Run("sfixed64_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_sfixed64_to_int32_map(?, JSON_OBJECT('-9223372036854775808', 1000))", `{}`).IsEqualToJsonString(`{"10": {"-9223372036854775808": 1000}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_sfixed64_to_int32_map(?)", `{"10": {"-9223372036854775808": 1000}}`).IsEqualToJsonString(`{"-9223372036854775808": 1000}`)
		RunTestThatExpression(t, "pbt_map_fields_count_sfixed64_to_int32_map(?)", `{"10": {"-9223372036854775808": 1000}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_sfixed64_to_int32_map(?)", `{"10": {"-9223372036854775808": 1000}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed64_to_int32_map__or(?, ?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `-9223372036854775808`, `101`).IsEqualToInt(1000) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed64_to_int32_map__or(?, ?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `123`, `101`).IsEqualToInt(101)                   // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed64_to_int32_map__or(?, ?, ?)", `{}`, `-9223372036854775808`, `101`).IsEqualToInt(101)                                      // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed64_to_int32_map(?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `-9223372036854775808`).IsEqualToInt(1000) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed64_to_int32_map(?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `123`).IsNull()                     // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_sfixed64_to_int32_map(?, ?)", `{}`, `-9223372036854775808`).IsNull()                                        // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_sfixed64_to_int32_map(?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `-9223372036854775808`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_sfixed64_to_int32_map(?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `123`).IsEqualToBool(false)                 // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_sfixed64_to_int32_map(?, ?)", `{}`, `-9223372036854775808`).IsEqualToBool(false)                                    // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_sfixed64_to_int32_map(?, ?, ?)", `{}`, `-9223372036854775808`, `1000`).IsEqualToJsonString(`{"10": {"-9223372036854775808": 1000}}`)                                    // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_sfixed64_to_int32_map(?, ?, ?)", `{"10": {"100": 50}}`, `-9223372036854775808`, `1000`).IsEqualToJsonString(`{"10": {"100": 50, "-9223372036854775808": 1000}}`)        // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_sfixed64_to_int32_map(?, ?, ?)", `{"10": {"-9223372036854775808": 950}}`, `-9223372036854775808`, `1000`).IsEqualToJsonString(`{"10": {"-9223372036854775808": 1000}}`) // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_sfixed64_to_int32_map(?, ?)", `{}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"10": {"-300": 333, "-400": 444}}`)                                          // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_sfixed64_to_int32_map(?, ?)", `{"10": {"100": 50}}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"10": {"100": 50, "-300": 333, "-400": 444}}`)              // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_sfixed64_to_int32_map(?, ?)", `{"10": {"100": 50, "-300": 250}}`, `{"-300": 333, "-400": 444}`).IsEqualToJsonString(`{"10": {"100": 50, "-300": 333, "-400": 444}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_sfixed64_to_int32_map(?, ?)", `{"10": {"-9223372036854775808": 1000, "100": 50}}`, `-9223372036854775808`).IsEqualToJsonString(`{"10": {"100": 50}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_sfixed64_to_int32_map(?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `-9223372036854775808`).IsEqualToJsonString(`{}`)                             // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_sfixed64_to_int32_map(?, ?)", `{"10": {"-9223372036854775808": 1000}}`, `123`).IsEqualToJsonString(`{"10": {"-9223372036854775808": 1000}}`)          // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_sfixed64_to_int32_map(?, ?)", `{}`, `-9223372036854775808`).IsEqualToJsonString(`{}`)                                                                 // Remove from empty map
	})

	// Test non-integer key types
	t.Run("bool_key", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_bool_to_int32_map(?, JSON_OBJECT('true', 1100))", `{}`).IsEqualToJsonString(`{"11": {"true": 1100}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_bool_to_int32_map(?)", `{"11": {"true": 1100, "false": 0}}`).IsEqualToJsonString(`{"true": 1100, "false": 0}`)
		RunTestThatExpression(t, "pbt_map_fields_count_bool_to_int32_map(?)", `{"11": {"true": 1100, "false": 0}}`).IsEqualToInt(2)
		RunTestThatExpression(t, "pbt_map_fields_clear_bool_to_int32_map(?)", `{"11": {"true": 1100}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_bool_to_int32_map__or(?, ?, ?)", `{"11": {"true": 1100, "false": 0}}`, true, 999).IsEqualToInt(1100) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_bool_to_int32_map__or(?, ?, ?)", `{"11": {"true": 1100, "false": 0}}`, false, 999).IsEqualToInt(0)   // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_get_bool_to_int32_map__or(?, ?, ?)", `{}`, true, 999).IsEqualToInt(999)                                  // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_bool_to_int32_map(?, ?)", `{"11": {"true": 1100, "false": 0}}`, true).IsEqualToInt(1100) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_bool_to_int32_map(?, ?)", `{"11": {"true": 1100, "false": 0}}`, false).IsEqualToInt(0)   // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_get_bool_to_int32_map(?, ?)", `{}`, true).IsNull()                                    // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_bool_to_int32_map(?, ?)", `{"11": {"true": 1100, "false": 0}}`, true).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_bool_to_int32_map(?, ?)", `{"11": {"true": 1100, "false": 0}}`, false).IsEqualToBool(true) // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_contains_bool_to_int32_map(?, ?)", `{}`, true).IsEqualToBool(false)                                 // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_bool_to_int32_map(?, ?, ?)", `{}`, true, 1100).IsEqualToJsonString(`{"11": {"true": 1100}}`)                               // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_bool_to_int32_map(?, ?, ?)", `{"11": {"true": 1100}}`, false, 0).IsEqualToJsonString(`{"11": {"true": 1100, "false": 0}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_bool_to_int32_map(?, ?, ?)", `{"11": {"true": 500}}`, true, 1100).IsEqualToJsonString(`{"11": {"true": 1100}}`)            // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_bool_to_int32_map(?, ?)", `{}`, `{"true": 1100, "false": 0}`).IsEqualToJsonString(`{"11": {"true": 1100, "false": 0}}`)                                  // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_bool_to_int32_map(?, ?)", `{"11": {"true": 500}}`, `{"false": 0}`).IsEqualToJsonString(`{"11": {"true": 500, "false": 0}}`)                              // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_bool_to_int32_map(?, ?)", `{"11": {"true": 500, "false": 200}}`, `{"true": 1100, "false": 0}`).IsEqualToJsonString(`{"11": {"true": 1100, "false": 0}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_bool_to_int32_map(?, ?)", `{"11": {"true": 1100, "false": 0}}`, true).IsEqualToJsonString(`{"11": {"false": 0}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_bool_to_int32_map(?, ?)", `{"11": {"true": 1100}}`, true).IsEqualToJsonString(`{}`)                               // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_bool_to_int32_map(?, ?)", `{"11": {"true": 1100}}`, false).IsEqualToJsonString(`{"11": {"true": 1100}}`)          // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_bool_to_int32_map(?, ?)", `{}`, true).IsEqualToJsonString(`{}`)                                                   // Remove from empty map
	})

	t.Run("string_key", func(t *testing.T) {
		// Test multiple entries in same map
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_int32_map(?, JSON_OBJECT('key', 1200))", `{}`).IsEqualToJsonString(`{"12": {"key": 1200}}`)
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_int32_map(?, JSON_OBJECT('first', 10, 'second', 20))", `{}`).IsEqualToJsonString(`{"12": {"first": 10, "second": 20}}`)

		// Test getters return entire map
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_int32_map(?)", `{"12": {"first": 10, "second": 20}}`).IsEqualToJsonString(`{"first": 10, "second": 20}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_int32_map(?)", `{}`).IsEqualToJsonString(`{}`) // Default when absent

		// Test overwriting existing map (replaces entire map)
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_int32_map(pbt_map_fields_set_all_string_to_int32_map(pbt_map_fields_new(), JSON_OBJECT('old', 100)), JSON_OBJECT('new', 200))").IsEqualToJsonString(`{"12": {"new": 200}}`)

		// Test map count and clear operations
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_int32_map(?)", `{"12": {"first": 10, "second": 20}}`).IsEqualToInt(2)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_int32_map(?)", `{"12": {"key": 1200}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int32_map__or(?, ?, ?)", `{"12": {"key": 555}}`, `key`, `9999`).IsEqualToInt(555)      // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int32_map__or(?, ?, ?)", `{"12": {"key": 555}}`, `missing`, `9999`).IsEqualToInt(9999) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int32_map__or(?, ?, ?)", `{}`, `key`, `9999`).IsEqualToInt(9999)                       // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int32_map(?, ?)", `{"12": {"key": 555}}`, `key`).IsEqualToInt(555)   // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int32_map(?, ?)", `{"12": {"key": 555}}`, `missing`).IsNull() // Key missing, return NULL
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int32_map(?, ?)", `{}`, `key`).IsNull()                       // Map empty, return NULL

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_int32_map(?, ?)", `{"12": {"key": 555}}`, `key`).IsEqualToBool(true)      // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_int32_map(?, ?)", `{"12": {"key": 555}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_int32_map(?, ?)", `{}`, `key`).IsEqualToBool(false)                       // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_int32_map(?, ?, ?)", `{}`, `new_key`, `777`).IsEqualToJsonString(`{"12": {"new_key": 777}}`)                                         // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_int32_map(?, ?, ?)", `{"12": {"existing": 100}}`, `new_key`, `777`).IsEqualToJsonString(`{"12": {"existing": 100, "new_key": 777}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_int32_map(?, ?, ?)", `{"12": {"key": 500}}`, `key`, `777`).IsEqualToJsonString(`{"12": {"key": 777}}`)                               // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_int32_map(?, ?)", `{}`, `{"alpha": 111, "beta": 222}`).IsEqualToJsonString(`{"12": {"alpha": 111, "beta": 222}}`)                                                       // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_int32_map(?, ?)", `{"12": {"existing": 100}}`, `{"alpha": 111, "beta": 222}`).IsEqualToJsonString(`{"12": {"existing": 100, "alpha": 111, "beta": 222}}`)               // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_int32_map(?, ?)", `{"12": {"existing": 100, "alpha": 999}}`, `{"alpha": 111, "beta": 222}`).IsEqualToJsonString(`{"12": {"existing": 100, "alpha": 111, "beta": 222}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_int32_map(?, ?)", `{"12": {"key": 555, "other": 777}}`, `key`).IsEqualToJsonString(`{"12": {"other": 777}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_int32_map(?, ?)", `{"12": {"key": 555}}`, `key`).IsEqualToJsonString(`{}`)                                   // Remove last key should clear entire field
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_int32_map(?, ?)", `{"12": {"key": 555}}`, `missing`).IsEqualToJsonString(`{"12": {"key": 555}}`)             // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_int32_map(?, ?)", `{}`, `key`).IsEqualToJsonString(`{}`)                                                     // Remove from empty map
	})

	// Test different value types with string keys
	t.Run("double_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_double_map(?, ?)", `{}`, `{"pi": 3.141592653589793}`).IsEqualToJsonString(`{"13": {"pi": "binary64:0x400921fb54442d18"}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_double_map(?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`).IsEqualToJsonString(`{"pi": 3.141592653589793}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_double_map(?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_double_map(?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_double_map__or(?, ?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `pi`, 4).IsEqualToDouble(3.141592653589793) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_double_map__or(?, ?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `missing`, 4).IsEqualToDouble(4)            // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_double_map__or(?, ?, ?)", `{}`, `pi`, 4).IsEqualToDouble(4)                                                            // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_double_map(?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `pi`).IsEqualToDouble(3.141592653589793) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_double_map(?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `missing`).IsNull()            // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_double_map(?, ?)", `{}`, `pi`).IsNull()                                                            // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_double_map(?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `pi`).IsEqualToBool(true)       // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_double_map(?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_double_map(?, ?)", `{}`, `pi`).IsEqualToBool(false)                                                 // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_double_map(?, ?, ?)", `{}`, `e`, 2.718).IsEqualToJsonString(`{"13": {"e": "binary64:0x4005be76c8b43958"}}`)
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_double_map(?, ?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `e`, 2.718).IsEqualToJsonString(`{"13": {"pi": "binary64:0x400921fb54442d18", "e": "binary64:0x4005be76c8b43958"}}`)
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_double_map(?, ?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `pi`, 3.14).IsEqualToJsonString(`{"13": {"pi": "binary64:0x40091eb851eb851f"}}`)

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_double_map(?, ?)", `{}`, `{"sqrt2": 1.414, "phi": 1.618}`).IsEqualToJsonString(`{"13": {"sqrt2": "binary64:0x3ff69fbe76c8b439", "phi": "binary64:0x3ff9e353f7ced917"}}`)
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_double_map(?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `{"e": 2.718}`).IsEqualToJsonString(`{"13": {"pi": "binary64:0x400921fb54442d18", "e": "binary64:0x4005be76c8b43958"}}`)

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_double_map(?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18", "e": "binary64:0x4005be76c8b43958"}}`, `pi`).IsEqualToJsonString(`{"13": {"e": "binary64:0x4005be76c8b43958"}}`)
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_double_map(?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `pi`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_double_map(?, ?)", `{"13": {"pi": "binary64:0x400921fb54442d18"}}`, `missing`).IsEqualToJsonString(`{"13": {"pi": "binary64:0x400921fb54442d18"}}`)
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_double_map(?, ?)", `{}`, `pi`).IsEqualToJsonString(`{}`)
	})

	t.Run("float_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_float_map(?, ?)", `{}`, `{"pi_float": 3.14}`).IsEqualToJsonString(`{"14": {"pi_float": "binary32:0x4048f5c3"}}`)
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_float_map(?, ?)", `{}`, `{"first": 3.14, "second": 2.718}`).IsEqualToJsonString(`{"14": {"first": "binary32:0x4048f5c3", "second": "binary32:0x402df3b6"}}`)

		// Test set_all replaces entire map
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_float_map(?, ?)", `{"14": {"old": "binary32:0x3f800000"}}`, `{"new": 2.0}`).IsEqualToJsonString(`{"14": {"new": "binary32:0x40000000"}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_float_map(?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`).IsEqualToJsonString(`{"pi_float": 3.140000104904175}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_float_map(?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_float_map(?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_float_map__or(?, ?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `pi_float`, 4).IsEqualToFloat(3.140000104904175) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_float_map__or(?, ?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `missing`, 4).IsEqualToFloat(4)                  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_float_map__or(?, ?, ?)", `{}`, `pi_float`, 4).IsEqualToFloat(4)                                                          // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_float_map(?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `pi_float`).IsEqualToFloat(3.140000104904175) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_float_map(?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `missing`).IsNull()                  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_float_map(?, ?)", `{}`, `pi_float`).IsNull()                                                          // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_float_map(?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `pi_float`).IsEqualToBool(true) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_float_map(?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_float_map(?, ?)", `{}`, `pi_float`).IsEqualToBool(false)                                         // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_float_map(?, ?, ?)", `{}`, `e_float`, 2.718).IsEqualToJsonString(`{"14": {"e_float": "binary32:0x402df3b6"}}`)
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_float_map(?, ?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `e_float`, 2.718).IsEqualToJsonString(`{"14": {"pi_float": "binary32:0x4048f5c3", "e_float": "binary32:0x402df3b6"}}`)
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_float_map(?, ?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `pi_float`, 3.14).IsEqualToJsonString(`{"14": {"pi_float": "binary32:0x4048f5c3"}}`)

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_float_map(?, ?)", `{}`, `{"sqrt2": 1.414, "phi": 1.618}`).IsEqualToJsonString(`{"14": {"sqrt2": "binary32:0x3fb4fdf4", "phi": "binary32:0x3fcf1aa0"}}`)
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_float_map(?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `{"e_float": 2.718}`).IsEqualToJsonString(`{"14": {"pi_float": "binary32:0x4048f5c3", "e_float": "binary32:0x402df3b6"}}`)

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_float_map(?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3", "e_float": "binary32:0x402df3b6"}}`, `pi_float`).IsEqualToJsonString(`{"14": {"e_float": "binary32:0x402df3b6"}}`)
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_float_map(?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `pi_float`).IsEqualToJsonString(`{}`)
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_float_map(?, ?)", `{"14": {"pi_float": "binary32:0x4048f5c3"}}`, `missing`).IsEqualToJsonString(`{"14": {"pi_float": "binary32:0x4048f5c3"}}`)
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_float_map(?, ?)", `{}`, `pi_float`).IsEqualToJsonString(`{}`)
	})

	t.Run("int32_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_int32_map(?, JSON_OBJECT('key', -2147483648))", `{}`).IsEqualToJsonString(`{"12": {"key": -2147483648}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_int32_map(?)", `{"12": {"key": -2147483648}}`).IsEqualToJsonString(`{"key": -2147483648}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_int32_map(?)", `{"12": {"key": -2147483648}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_int32_map(?)", `{"12": {"key": -2147483648}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int32_map__or(?, ?, ?)", `{"12": {"key": 12345}}`, `key`, `0`).IsEqualToInt(12345) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int32_map__or(?, ?, ?)", `{"12": {"key": 12345}}`, `missing`, `0`).IsEqualToInt(0) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int32_map__or(?, ?, ?)", `{}`, `key`, `0`).IsEqualToInt(0)                         // Map empty

		// Test that maps can store default/zero values (unlike regular proto3 fields without presence)
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_int32_map(?, JSON_OBJECT('zero', 0))", `{}`).IsEqualToJsonString(`{"12": {"zero": 0}}`)
	})

	t.Run("int64_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_int64_map(?, JSON_OBJECT('big', 9223372036854775807))", `{}`).IsEqualToJsonString(`{"15": {"big": 9223372036854775807}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_int64_map(?)", `{"15": {"big": 9223372036854775807}}`).IsEqualToJsonString(`{"big": 9223372036854775807}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_int64_map(?)", `{"15": {"big": 9223372036854775807}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_int64_map(?)", `{"15": {"big": 9223372036854775807}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int64_map__or(?, ?, ?)", `{"15": {"big": 9223372036854775807}}`, `big`, `-1`).IsEqualToInt(9223372036854775807) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int64_map__or(?, ?, ?)", `{"15": {"big": 9223372036854775807}}`, `missing`, `-1`).IsEqualToInt(-1)              // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int64_map__or(?, ?, ?)", `{}`, `big`, `-1`).IsEqualToInt(-1)                                                    // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int64_map(?, ?)", `{"15": {"big": 9223372036854775807}}`, `big`).IsEqualToInt(9223372036854775807) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int64_map(?, ?)", `{"15": {"big": 9223372036854775807}}`, `missing`).IsNull()                     // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_int64_map(?, ?)", `{}`, `big`).IsNull()                                                            // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_int64_map(?, ?)", `{"15": {"big": 9223372036854775807}}`, `big`).IsEqualToBool(true)     // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_int64_map(?, ?)", `{"15": {"big": 9223372036854775807}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_int64_map(?, ?)", `{}`, `big`).IsEqualToBool(false)                                        // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_int64_map(?, ?, ?)", `{}`, `big`, `9223372036854775807`).IsEqualToJsonString(`{"15": {"big": 9223372036854775807}}`)                // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_int64_map(?, ?, ?)", `{"15": {"other": 123}}`, `big`, `9223372036854775807`).IsEqualToJsonString(`{"15": {"other": 123, "big": 9223372036854775807}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_int64_map(?, ?, ?)", `{"15": {"big": 999}}`, `big`, `9223372036854775807`).IsEqualToJsonString(`{"15": {"big": 9223372036854775807}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_int64_map(?, ?)", `{}`, `{"large1": 5000000000000000000, "large2": 7000000000000000000}`).IsEqualToJsonString(`{"15": {"large1": 5000000000000000000, "large2": 7000000000000000000}}`)                                                 // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_int64_map(?, ?)", `{"15": {"other": 123}}`, `{"large1": 5000000000000000000, "large2": 7000000000000000000}`).IsEqualToJsonString(`{"15": {"other": 123, "large1": 5000000000000000000, "large2": 7000000000000000000}}`)                  // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_int64_map(?, ?)", `{"15": {"other": 123, "large1": 999}}`, `{"large1": 5000000000000000000, "large2": 7000000000000000000}`).IsEqualToJsonString(`{"15": {"other": 123, "large1": 5000000000000000000, "large2": 7000000000000000000}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_int64_map(?, ?)", `{"15": {"big": 9223372036854775807, "other": 123}}`, `big`).IsEqualToJsonString(`{"15": {"other": 123}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_int64_map(?, ?)", `{"15": {"big": 9223372036854775807}}`, `big`).IsEqualToJsonString(`{}`)                   // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_int64_map(?, ?)", `{"15": {"big": 9223372036854775807}}`, `missing`).IsEqualToJsonString(`{"15": {"big": 9223372036854775807}}`) // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_int64_map(?, ?)", `{}`, `big`).IsEqualToJsonString(`{}`)                                                       // Remove from empty map
	})

	t.Run("uint32_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_uint32_map(?, JSON_OBJECT('max32', 4294967295))", `{}`).IsEqualToJsonString(`{"16": {"max32": 4294967295}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_uint32_map(?)", `{"16": {"max32": 4294967295}}`).IsEqualToJsonString(`{"max32": 4294967295}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_uint32_map(?)", `{"16": {"max32": 4294967295}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_uint32_map(?)", `{"16": {"max32": 4294967295}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint32_map__or(?, ?, ?)", `{"16": {"max32": 4294967295}}`, `max32`, `0`).IsEqualToUint(4294967295) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint32_map__or(?, ?, ?)", `{"16": {"max32": 4294967295}}`, `missing`, `0`).IsEqualToUint(0)        // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint32_map__or(?, ?, ?)", `{}`, `max32`, `0`).IsEqualToUint(0)                                     // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint32_map(?, ?)", `{"16": {"max32": 4294967295}}`, `max32`).IsEqualToUint(4294967295) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint32_map(?, ?)", `{"16": {"max32": 4294967295}}`, `missing`).IsNull()                // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint32_map(?, ?)", `{}`, `max32`).IsNull()                                             // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_uint32_map(?, ?)", `{"16": {"max32": 4294967295}}`, `max32`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_uint32_map(?, ?)", `{"16": {"max32": 4294967295}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_uint32_map(?, ?)", `{}`, `max32`).IsEqualToBool(false)                            // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_uint32_map(?, ?, ?)", `{}`, `max32`, `4294967295`).IsEqualToJsonString(`{"16": {"max32": 4294967295}}`)                      // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_uint32_map(?, ?, ?)", `{"16": {"other": 123}}`, `max32`, `4294967295`).IsEqualToJsonString(`{"16": {"other": 123, "max32": 4294967295}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_uint32_map(?, ?, ?)", `{"16": {"max32": 999}}`, `max32`, `4294967295`).IsEqualToJsonString(`{"16": {"max32": 4294967295}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_uint32_map(?, ?)", `{}`, `{"large1": 3000000000, "large2": 4000000000}`).IsEqualToJsonString(`{"16": {"large1": 3000000000, "large2": 4000000000}}`)                                                 // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_uint32_map(?, ?)", `{"16": {"other": 123}}`, `{"large1": 3000000000, "large2": 4000000000}`).IsEqualToJsonString(`{"16": {"other": 123, "large1": 3000000000, "large2": 4000000000}}`)                  // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_uint32_map(?, ?)", `{"16": {"other": 123, "large1": 999}}`, `{"large1": 3000000000, "large2": 4000000000}`).IsEqualToJsonString(`{"16": {"other": 123, "large1": 3000000000, "large2": 4000000000}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_uint32_map(?, ?)", `{"16": {"max32": 4294967295, "other": 123}}`, `max32`).IsEqualToJsonString(`{"16": {"other": 123}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_uint32_map(?, ?)", `{"16": {"max32": 4294967295}}`, `max32`).IsEqualToJsonString(`{}`)                    // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_uint32_map(?, ?)", `{"16": {"max32": 4294967295}}`, `missing`).IsEqualToJsonString(`{"16": {"max32": 4294967295}}`) // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_uint32_map(?, ?)", `{}`, `max32`).IsEqualToJsonString(`{}`)                                                      // Remove from empty map
	})

	t.Run("uint64_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_uint64_map(?, JSON_OBJECT('max64', 18446744073709551615))", `{}`).IsEqualToJsonString(`{"17": {"max64": 18446744073709551615}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_uint64_map(?)", `{"17": {"max64": 18446744073709551615}}`).IsEqualToJsonString(`{"max64": 18446744073709551615}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_uint64_map(?)", `{"17": {"max64": 18446744073709551615}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_uint64_map(?)", `{"17": {"max64": 18446744073709551615}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint64_map__or(?, ?, ?)", `{"17": {"max64": 18446744073709551615}}`, `max64`, `1`).IsEqualToUint(18446744073709551615) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint64_map__or(?, ?, ?)", `{"17": {"max64": 18446744073709551615}}`, `missing`, `1`).IsEqualToUint(1)                  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint64_map__or(?, ?, ?)", `{}`, `max64`, `1`).IsEqualToUint(1)                                                         // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint64_map(?, ?)", `{"17": {"max64": 18446744073709551615}}`, `max64`).IsEqualToUint(18446744073709551615) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint64_map(?, ?)", `{"17": {"max64": 18446744073709551615}}`, `missing`).IsNull()                          // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_uint64_map(?, ?)", `{}`, `max64`).IsNull()                                                                 // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_uint64_map(?, ?)", `{"17": {"max64": 18446744073709551615}}`, `max64`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_uint64_map(?, ?)", `{"17": {"max64": 18446744073709551615}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_uint64_map(?, ?)", `{}`, `max64`).IsEqualToBool(false)                                       // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_uint64_map(?, ?, ?)", `{}`, `max64`, `18446744073709551615`).IsEqualToJsonString(`{"17": {"max64": 18446744073709551615}}`)                      // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_uint64_map(?, ?, ?)", `{"17": {"other": 123}}`, `max64`, `18446744073709551615`).IsEqualToJsonString(`{"17": {"other": 123, "max64": 18446744073709551615}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_uint64_map(?, ?, ?)", `{"17": {"max64": 999}}`, `max64`, `18446744073709551615`).IsEqualToJsonString(`{"17": {"max64": 18446744073709551615}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_uint64_map(?, ?)", `{}`, `{"large1": 10000000000000000000, "large2": 15000000000000000000}`).IsEqualToJsonString(`{"17": {"large1": 10000000000000000000, "large2": 15000000000000000000}}`)                                                 // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_uint64_map(?, ?)", `{"17": {"other": 123}}`, `{"large1": 10000000000000000000, "large2": 15000000000000000000}`).IsEqualToJsonString(`{"17": {"other": 123, "large1": 10000000000000000000, "large2": 15000000000000000000}}`)                  // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_uint64_map(?, ?)", `{"17": {"other": 123, "large1": 999}}`, `{"large1": 10000000000000000000, "large2": 15000000000000000000}`).IsEqualToJsonString(`{"17": {"other": 123, "large1": 10000000000000000000, "large2": 15000000000000000000}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_uint64_map(?, ?)", `{"17": {"max64": 18446744073709551615, "other": 123}}`, `max64`).IsEqualToJsonString(`{"17": {"other": 123}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_uint64_map(?, ?)", `{"17": {"max64": 18446744073709551615}}`, `max64`).IsEqualToJsonString(`{}`)                    // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_uint64_map(?, ?)", `{"17": {"max64": 18446744073709551615}}`, `missing`).IsEqualToJsonString(`{"17": {"max64": 18446744073709551615}}`) // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_uint64_map(?, ?)", `{}`, `max64`).IsEqualToJsonString(`{}`)                                                      // Remove from empty map
	})

	t.Run("sint32_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_sint32_map(?, JSON_OBJECT('negative', -2147483648))", `{}`).IsEqualToJsonString(`{"18": {"negative": -2147483648}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_sint32_map(?)", `{"18": {"negative": -2147483648}}`).IsEqualToJsonString(`{"negative": -2147483648}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_sint32_map(?)", `{"18": {"negative": -2147483648}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_sint32_map(?)", `{"18": {"negative": -2147483648}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint32_map__or(?, ?, ?)", `{"18": {"negative": -2147483648}}`, `negative`, `0`).IsEqualToInt(-2147483648) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint32_map__or(?, ?, ?)", `{"18": {"negative": -2147483648}}`, `missing`, `0`).IsEqualToInt(0)            // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint32_map__or(?, ?, ?)", `{}`, `negative`, `0`).IsEqualToInt(0)                                          // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint32_map(?, ?)", `{"18": {"negative": -2147483648}}`, `negative`).IsEqualToInt(-2147483648) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint32_map(?, ?)", `{"18": {"negative": -2147483648}}`, `missing`).IsNull()                    // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint32_map(?, ?)", `{}`, `negative`).IsNull()                                                 // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sint32_map(?, ?)", `{"18": {"negative": -2147483648}}`, `negative`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sint32_map(?, ?)", `{"18": {"negative": -2147483648}}`, `missing`).IsEqualToBool(false)  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sint32_map(?, ?)", `{}`, `negative`).IsEqualToBool(false)                               // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sint32_map(?, ?, ?)", `{}`, `negative`, `-2147483648`).IsEqualToJsonString(`{"18": {"negative": -2147483648}}`)                      // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sint32_map(?, ?, ?)", `{"18": {"other": 123}}`, `negative`, `-2147483648`).IsEqualToJsonString(`{"18": {"other": 123, "negative": -2147483648}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sint32_map(?, ?, ?)", `{"18": {"negative": 999}}`, `negative`, `-2147483648`).IsEqualToJsonString(`{"18": {"negative": -2147483648}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sint32_map(?, ?)", `{}`, `{"neg1": -1000, "neg2": -2000}`).IsEqualToJsonString(`{"18": {"neg1": -1000, "neg2": -2000}}`)                                                 // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sint32_map(?, ?)", `{"18": {"other": 123}}`, `{"neg1": -1000, "neg2": -2000}`).IsEqualToJsonString(`{"18": {"other": 123, "neg1": -1000, "neg2": -2000}}`)                  // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sint32_map(?, ?)", `{"18": {"other": 123, "neg1": 999}}`, `{"neg1": -1000, "neg2": -2000}`).IsEqualToJsonString(`{"18": {"other": 123, "neg1": -1000, "neg2": -2000}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sint32_map(?, ?)", `{"18": {"negative": -2147483648, "other": 123}}`, `negative`).IsEqualToJsonString(`{"18": {"other": 123}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sint32_map(?, ?)", `{"18": {"negative": -2147483648}}`, `negative`).IsEqualToJsonString(`{}`)                              // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sint32_map(?, ?)", `{"18": {"negative": -2147483648}}`, `missing`).IsEqualToJsonString(`{"18": {"negative": -2147483648}}`) // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sint32_map(?, ?)", `{}`, `negative`).IsEqualToJsonString(`{}`)                                                                // Remove from empty map
	})

	t.Run("sint64_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_sint64_map(?, JSON_OBJECT('big_negative', -9223372036854775808))", `{}`).IsEqualToJsonString(`{"19": {"big_negative": -9223372036854775808}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_sint64_map(?)", `{"19": {"big_negative": -9223372036854775808}}`).IsEqualToJsonString(`{"big_negative": -9223372036854775808}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_sint64_map(?)", `{"19": {"big_negative": -9223372036854775808}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_sint64_map(?)", `{"19": {"big_negative": -9223372036854775808}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint64_map__or(?, ?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `big_negative`, `0`).IsEqualToInt(-9223372036854775808) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint64_map__or(?, ?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `missing`, `0`).IsEqualToInt(0)                         // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint64_map__or(?, ?, ?)", `{}`, `big_negative`, `0`).IsEqualToInt(0)                                                                // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint64_map(?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `big_negative`).IsEqualToInt(-9223372036854775808) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint64_map(?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `missing`).IsNull()                                  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sint64_map(?, ?)", `{}`, `big_negative`).IsNull()                                                                         // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sint64_map(?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `big_negative`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sint64_map(?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `missing`).IsEqualToBool(false)     // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sint64_map(?, ?)", `{}`, `big_negative`).IsEqualToBool(false)                                           // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sint64_map(?, ?, ?)", `{}`, `big_negative`, `-9223372036854775808`).IsEqualToJsonString(`{"19": {"big_negative": -9223372036854775808}}`)                      // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sint64_map(?, ?, ?)", `{"19": {"other": 123}}`, `big_negative`, `-9223372036854775808`).IsEqualToJsonString(`{"19": {"other": 123, "big_negative": -9223372036854775808}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sint64_map(?, ?, ?)", `{"19": {"big_negative": 999}}`, `big_negative`, `-9223372036854775808`).IsEqualToJsonString(`{"19": {"big_negative": -9223372036854775808}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sint64_map(?, ?)", `{}`, `{"neg1": -1000000000000, "neg2": -2000000000000}`).IsEqualToJsonString(`{"19": {"neg1": -1000000000000, "neg2": -2000000000000}}`)                                                 // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sint64_map(?, ?)", `{"19": {"other": 123}}`, `{"neg1": -1000000000000, "neg2": -2000000000000}`).IsEqualToJsonString(`{"19": {"other": 123, "neg1": -1000000000000, "neg2": -2000000000000}}`)                  // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sint64_map(?, ?)", `{"19": {"other": 123, "neg1": 999}}`, `{"neg1": -1000000000000, "neg2": -2000000000000}`).IsEqualToJsonString(`{"19": {"other": 123, "neg1": -1000000000000, "neg2": -2000000000000}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sint64_map(?, ?)", `{"19": {"big_negative": -9223372036854775808, "other": 123}}`, `big_negative`).IsEqualToJsonString(`{"19": {"other": 123}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sint64_map(?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `big_negative`).IsEqualToJsonString(`{}`)                              // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sint64_map(?, ?)", `{"19": {"big_negative": -9223372036854775808}}`, `missing`).IsEqualToJsonString(`{"19": {"big_negative": -9223372036854775808}}`) // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sint64_map(?, ?)", `{}`, `big_negative`).IsEqualToJsonString(`{}`)                                                                                // Remove from empty map
	})

	t.Run("fixed32_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_fixed32_map(?, JSON_OBJECT('max_fixed32', 4294967295))", `{}`).IsEqualToJsonString(`{"20": {"max_fixed32": 4294967295}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_fixed32_map(?)", `{"20": {"max_fixed32": 4294967295}}`).IsEqualToJsonString(`{"max_fixed32": 4294967295}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_fixed32_map(?)", `{"20": {"max_fixed32": 4294967295}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_fixed32_map(?)", `{"20": {"max_fixed32": 4294967295}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed32_map__or(?, ?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `max_fixed32`, `1`).IsEqualToUint(4294967295) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed32_map__or(?, ?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `missing`, `1`).IsEqualToUint(1)              // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed32_map__or(?, ?, ?)", `{}`, `max_fixed32`, `1`).IsEqualToUint(1)                                           // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed32_map(?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `max_fixed32`).IsEqualToUint(4294967295) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed32_map(?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `missing`).IsNull()                      // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed32_map(?, ?)", `{}`, `max_fixed32`).IsNull()                                                   // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_fixed32_map(?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `max_fixed32`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_fixed32_map(?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `missing`).IsEqualToBool(false)     // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_fixed32_map(?, ?)", `{}`, `max_fixed32`).IsEqualToBool(false)                                 // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_fixed32_map(?, ?, ?)", `{}`, `max_fixed32`, `4294967295`).IsEqualToJsonString(`{"20": {"max_fixed32": 4294967295}}`)                          // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_fixed32_map(?, ?, ?)", `{"20": {"other": 123}}`, `max_fixed32`, `4294967295`).IsEqualToJsonString(`{"20": {"other": 123, "max_fixed32": 4294967295}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_fixed32_map(?, ?, ?)", `{"20": {"max_fixed32": 999}}`, `max_fixed32`, `4294967295`).IsEqualToJsonString(`{"20": {"max_fixed32": 4294967295}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_fixed32_map(?, ?)", `{}`, `{"large1": 3000000000, "large2": 4000000000}`).IsEqualToJsonString(`{"20": {"large1": 3000000000, "large2": 4000000000}}`)                                                 // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_fixed32_map(?, ?)", `{"20": {"other": 123}}`, `{"large1": 3000000000, "large2": 4000000000}`).IsEqualToJsonString(`{"20": {"other": 123, "large1": 3000000000, "large2": 4000000000}}`)                  // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_fixed32_map(?, ?)", `{"20": {"other": 123, "large1": 999}}`, `{"large1": 3000000000, "large2": 4000000000}`).IsEqualToJsonString(`{"20": {"other": 123, "large1": 3000000000, "large2": 4000000000}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_fixed32_map(?, ?)", `{"20": {"max_fixed32": 4294967295, "other": 123}}`, `max_fixed32`).IsEqualToJsonString(`{"20": {"other": 123}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_fixed32_map(?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `max_fixed32`).IsEqualToJsonString(`{}`)                        // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_fixed32_map(?, ?)", `{"20": {"max_fixed32": 4294967295}}`, `missing`).IsEqualToJsonString(`{"20": {"max_fixed32": 4294967295}}`) // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_fixed32_map(?, ?)", `{}`, `max_fixed32`).IsEqualToJsonString(`{}`)                                                          // Remove from empty map
	})

	t.Run("fixed64_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_fixed64_map(?, JSON_OBJECT('max_fixed64', 18446744073709551615))", `{}`).IsEqualToJsonString(`{"21": {"max_fixed64": 18446744073709551615}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_fixed64_map(?)", `{"21": {"max_fixed64": 18446744073709551615}}`).IsEqualToJsonString(`{"max_fixed64": 18446744073709551615}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_fixed64_map(?)", `{"21": {"max_fixed64": 18446744073709551615}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_fixed64_map(?)", `{"21": {"max_fixed64": 18446744073709551615}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed64_map__or(?, ?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `max_fixed64`, `1`).IsEqualToUint(18446744073709551615) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed64_map__or(?, ?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `missing`, `1`).IsEqualToUint(1)                        // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed64_map__or(?, ?, ?)", `{}`, `max_fixed64`, `1`).IsEqualToUint(1)                                                               // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed64_map(?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `max_fixed64`).IsEqualToUint(18446744073709551615) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed64_map(?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `missing`).IsNull()                                  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_fixed64_map(?, ?)", `{}`, `max_fixed64`).IsNull()                                                                         // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_fixed64_map(?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `max_fixed64`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_fixed64_map(?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `missing`).IsEqualToBool(false)     // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_fixed64_map(?, ?)", `{}`, `max_fixed64`).IsEqualToBool(false)                                            // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_fixed64_map(?, ?, ?)", `{}`, `max_fixed64`, `18446744073709551615`).IsEqualToJsonString(`{"21": {"max_fixed64": 18446744073709551615}}`)                      // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_fixed64_map(?, ?, ?)", `{"21": {"other": 123}}`, `max_fixed64`, `18446744073709551615`).IsEqualToJsonString(`{"21": {"other": 123, "max_fixed64": 18446744073709551615}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_fixed64_map(?, ?, ?)", `{"21": {"max_fixed64": 999}}`, `max_fixed64`, `18446744073709551615`).IsEqualToJsonString(`{"21": {"max_fixed64": 18446744073709551615}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_fixed64_map(?, ?)", `{}`, `{"large1": 18446744073709551614, "large2": 18446744073709551613}`).IsEqualToJsonString(`{"21": {"large1": 18446744073709551614, "large2": 18446744073709551613}}`)                                                      // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_fixed64_map(?, ?)", `{"21": {"other": 123}}`, `{"large1": 18446744073709551614, "large2": 18446744073709551613}`).IsEqualToJsonString(`{"21": {"other": 123, "large1": 18446744073709551614, "large2": 18446744073709551613}}`)                            // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_fixed64_map(?, ?)", `{"21": {"other": 123, "large1": 999}}`, `{"large1": 18446744073709551614, "large2": 18446744073709551613}`).IsEqualToJsonString(`{"21": {"other": 123, "large1": 18446744073709551614, "large2": 18446744073709551613}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_fixed64_map(?, ?)", `{"21": {"max_fixed64": 18446744073709551615, "other": 123}}`, `max_fixed64`).IsEqualToJsonString(`{"21": {"other": 123}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_fixed64_map(?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `max_fixed64`).IsEqualToJsonString(`{}`)                              // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_fixed64_map(?, ?)", `{"21": {"max_fixed64": 18446744073709551615}}`, `missing`).IsEqualToJsonString(`{"21": {"max_fixed64": 18446744073709551615}}`) // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_fixed64_map(?, ?)", `{}`, `max_fixed64`).IsEqualToJsonString(`{}`)                                                                            // Remove from empty map
	})

	t.Run("sfixed32_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_sfixed32_map(?, JSON_OBJECT('min_sfixed32', -2147483648))", `{}`).IsEqualToJsonString(`{"22": {"min_sfixed32": -2147483648}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_sfixed32_map(?)", `{"22": {"min_sfixed32": -2147483648}}`).IsEqualToJsonString(`{"min_sfixed32": -2147483648}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_sfixed32_map(?)", `{"22": {"min_sfixed32": -2147483648}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_sfixed32_map(?)", `{"22": {"min_sfixed32": -2147483648}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed32_map__or(?, ?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `min_sfixed32`, `0`).IsEqualToInt(-2147483648) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed32_map__or(?, ?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `missing`, `0`).IsEqualToInt(0)                // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed32_map__or(?, ?, ?)", `{}`, `min_sfixed32`, `0`).IsEqualToInt(0)                                              // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed32_map(?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `min_sfixed32`).IsEqualToInt(-2147483648) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed32_map(?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `missing`).IsNull()                          // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed32_map(?, ?)", `{}`, `min_sfixed32`).IsNull()                                                       // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sfixed32_map(?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `min_sfixed32`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sfixed32_map(?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `missing`).IsEqualToBool(false)     // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sfixed32_map(?, ?)", `{}`, `min_sfixed32`).IsEqualToBool(false)                                      // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sfixed32_map(?, ?, ?)", `{}`, `min_sfixed32`, `-2147483648`).IsEqualToJsonString(`{"22": {"min_sfixed32": -2147483648}}`)                        // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sfixed32_map(?, ?, ?)", `{"22": {"other": 123}}`, `min_sfixed32`, `-2147483648`).IsEqualToJsonString(`{"22": {"other": 123, "min_sfixed32": -2147483648}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sfixed32_map(?, ?, ?)", `{"22": {"min_sfixed32": 999}}`, `min_sfixed32`, `-2147483648`).IsEqualToJsonString(`{"22": {"min_sfixed32": -2147483648}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sfixed32_map(?, ?)", `{}`, `{"neg1": -1000, "neg2": -2000}`).IsEqualToJsonString(`{"22": {"neg1": -1000, "neg2": -2000}}`)                                                 // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sfixed32_map(?, ?)", `{"22": {"other": 123}}`, `{"neg1": -1000, "neg2": -2000}`).IsEqualToJsonString(`{"22": {"other": 123, "neg1": -1000, "neg2": -2000}}`)                  // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sfixed32_map(?, ?)", `{"22": {"other": 123, "neg1": 999}}`, `{"neg1": -1000, "neg2": -2000}`).IsEqualToJsonString(`{"22": {"other": 123, "neg1": -1000, "neg2": -2000}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sfixed32_map(?, ?)", `{"22": {"min_sfixed32": -2147483648, "other": 123}}`, `min_sfixed32`).IsEqualToJsonString(`{"22": {"other": 123}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sfixed32_map(?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `min_sfixed32`).IsEqualToJsonString(`{}`)                              // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sfixed32_map(?, ?)", `{"22": {"min_sfixed32": -2147483648}}`, `missing`).IsEqualToJsonString(`{"22": {"min_sfixed32": -2147483648}}`) // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sfixed32_map(?, ?)", `{}`, `min_sfixed32`).IsEqualToJsonString(`{}`)                                                                    // Remove from empty map
	})

	t.Run("sfixed64_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_sfixed64_map(?, JSON_OBJECT('min_sfixed64', -9223372036854775808))", `{}`).IsEqualToJsonString(`{"23": {"min_sfixed64": -9223372036854775808}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_sfixed64_map(?)", `{"23": {"min_sfixed64": -9223372036854775808}}`).IsEqualToJsonString(`{"min_sfixed64": -9223372036854775808}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_sfixed64_map(?)", `{"23": {"min_sfixed64": -9223372036854775808}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_sfixed64_map(?)", `{"23": {"min_sfixed64": -9223372036854775808}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed64_map__or(?, ?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `min_sfixed64`, `0`).IsEqualToInt(-9223372036854775808) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed64_map__or(?, ?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `missing`, `0`).IsEqualToInt(0)                         // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed64_map__or(?, ?, ?)", `{}`, `min_sfixed64`, `0`).IsEqualToInt(0)                                                                // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed64_map(?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `min_sfixed64`).IsEqualToInt(-9223372036854775808) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed64_map(?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `missing`).IsNull()                                  // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_sfixed64_map(?, ?)", `{}`, `min_sfixed64`).IsNull()                                                                         // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sfixed64_map(?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `min_sfixed64`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sfixed64_map(?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `missing`).IsEqualToBool(false)     // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_sfixed64_map(?, ?)", `{}`, `min_sfixed64`).IsEqualToBool(false)                                              // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sfixed64_map(?, ?, ?)", `{}`, `min_sfixed64`, `-9223372036854775808`).IsEqualToJsonString(`{"23": {"min_sfixed64": -9223372036854775808}}`)                        // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sfixed64_map(?, ?, ?)", `{"23": {"other": 123}}`, `min_sfixed64`, `-9223372036854775808`).IsEqualToJsonString(`{"23": {"other": 123, "min_sfixed64": -9223372036854775808}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_sfixed64_map(?, ?, ?)", `{"23": {"min_sfixed64": 999}}`, `min_sfixed64`, `-9223372036854775808`).IsEqualToJsonString(`{"23": {"min_sfixed64": -9223372036854775808}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sfixed64_map(?, ?)", `{}`, `{"neg1": -1000000000000, "neg2": -2000000000000}`).IsEqualToJsonString(`{"23": {"neg1": -1000000000000, "neg2": -2000000000000}}`)                                                         // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sfixed64_map(?, ?)", `{"23": {"other": 123}}`, `{"neg1": -1000000000000, "neg2": -2000000000000}`).IsEqualToJsonString(`{"23": {"other": 123, "neg1": -1000000000000, "neg2": -2000000000000}}`)                          // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_sfixed64_map(?, ?)", `{"23": {"other": 123, "neg1": 999}}`, `{"neg1": -1000000000000, "neg2": -2000000000000}`).IsEqualToJsonString(`{"23": {"other": 123, "neg1": -1000000000000, "neg2": -2000000000000}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sfixed64_map(?, ?)", `{"23": {"min_sfixed64": -9223372036854775808, "other": 123}}`, `min_sfixed64`).IsEqualToJsonString(`{"23": {"other": 123}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sfixed64_map(?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `min_sfixed64`).IsEqualToJsonString(`{}`)                              // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sfixed64_map(?, ?)", `{"23": {"min_sfixed64": -9223372036854775808}}`, `missing`).IsEqualToJsonString(`{"23": {"min_sfixed64": -9223372036854775808}}`) // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_sfixed64_map(?, ?)", `{}`, `min_sfixed64`).IsEqualToJsonString(`{}`)                                                                            // Remove from empty map
	})

	t.Run("bool_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_bool_map(?, JSON_OBJECT('false', false))", `{}`).IsEqualToJsonString(`{"24": {"false": false}}`)
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_bool_map(?, JSON_OBJECT('flag', true))", `{}`).IsEqualToJsonString(`{"24": {"flag": true}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_bool_map(?)", `{"24": {"flag": true, "other": false}}`).IsEqualToJsonString(`{"flag": true, "other": false}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_bool_map(?)", `{"24": {"flag": true}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_bool_map(?)", `{"24": {"flag": true}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bool_map__or(?, ?, ?)", `{"24": {"flag": true, "other": false}}`, `flag`, false).IsEqualToBool(true)     // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bool_map__or(?, ?, ?)", `{"24": {"flag": true, "other": false}}`, `other`, false).IsEqualToBool(false)   // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bool_map__or(?, ?, ?)", `{"24": {"flag": true, "other": false}}`, `missing`, false).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bool_map__or(?, ?, ?)", `{}`, `flag`, false).IsEqualToBool(false)                                        // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bool_map(?, ?)", `{"24": {"flag": true, "other": false}}`, `flag`).IsEqualToBool(true)   // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bool_map(?, ?)", `{"24": {"flag": true, "other": false}}`, `other`).IsEqualToBool(false) // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bool_map(?, ?)", `{"24": {"flag": true, "other": false}}`, `missing`).IsNull()            // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bool_map(?, ?)", `{}`, `flag`).IsNull()                                                      // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_bool_map(?, ?)", `{"24": {"flag": true, "other": false}}`, `flag`).IsEqualToBool(true)    // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_bool_map(?, ?)", `{"24": {"flag": true, "other": false}}`, `other`).IsEqualToBool(true)   // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_bool_map(?, ?)", `{"24": {"flag": true, "other": false}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_bool_map(?, ?)", `{}`, `flag`).IsEqualToBool(false)                                           // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_bool_map(?, ?, ?)", `{}`, `flag`, true).IsEqualToJsonString(`{"24": {"flag": true}}`)                        // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_bool_map(?, ?, ?)", `{"24": {"other": false}}`, `flag`, true).IsEqualToJsonString(`{"24": {"other": false, "flag": true}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_bool_map(?, ?, ?)", `{"24": {"flag": false}}`, `flag`, true).IsEqualToJsonString(`{"24": {"flag": true}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_bool_map(?, ?)", `{}`, `{"bool1": true, "bool2": false}`).IsEqualToJsonString(`{"24": {"bool1": true, "bool2": false}}`)                                               // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_bool_map(?, ?)", `{"24": {"other": true}}`, `{"bool1": true, "bool2": false}`).IsEqualToJsonString(`{"24": {"other": true, "bool1": true, "bool2": false}}`)                // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_bool_map(?, ?)", `{"24": {"other": true, "bool1": false}}`, `{"bool1": true, "bool2": false}`).IsEqualToJsonString(`{"24": {"other": true, "bool1": true, "bool2": false}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_bool_map(?, ?)", `{"24": {"flag": true, "other": false}}`, `flag`).IsEqualToJsonString(`{"24": {"other": false}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_bool_map(?, ?)", `{"24": {"flag": true}}`, `flag`).IsEqualToJsonString(`{}`)                              // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_bool_map(?, ?)", `{"24": {"flag": true}}`, `missing`).IsEqualToJsonString(`{"24": {"flag": true}}`)       // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_bool_map(?, ?)", `{}`, `flag`).IsEqualToJsonString(`{}`)                                                    // Remove from empty map
	})

	t.Run("string_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_string_map(?, JSON_OBJECT('empty', ''))", `{}`).IsEqualToJsonString(`{"25": {"empty": ""}}`)
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_string_map(?, JSON_OBJECT('greeting', 'hello'))", `{}`).IsEqualToJsonString(`{"25": {"greeting": "hello"}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_string_map(?)", `{"25": {"greeting": "hello"}}`).IsEqualToJsonString(`{"greeting": "hello"}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_string_map(?)", `{"25": {"greeting": "hello"}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_string_map(?)", `{"25": {"greeting": "hello"}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_string_map__or(?, ?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `greeting`, ``).IsEqualToString(`hello`) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_string_map__or(?, ?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `empty`, ``).IsEqualToString(``)         // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_string_map__or(?, ?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `missing`, ``).IsEqualToString(``)       // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_string_map__or(?, ?, ?)", `{}`, `greeting`, ``).IsEqualToString(``)                                              // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_string_map(?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `greeting`).IsEqualToString(`hello`) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_string_map(?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `empty`).IsEqualToString(``)         // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_string_map(?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `missing`).IsNull()                    // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_string_map(?, ?)", `{}`, `greeting`).IsNull()                                                              // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_string_map(?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `greeting`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_string_map(?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `empty`).IsEqualToBool(true)     // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_string_map(?, ?)", `{"25": {"greeting": "hello", "empty": ""}}`, `missing`).IsEqualToBool(false)   // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_string_map(?, ?)", `{}`, `greeting`).IsEqualToBool(false)                                             // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_string_map(?, ?, ?)", `{}`, `greeting`, `hello`).IsEqualToJsonString(`{"25": {"greeting": "hello"}}`)                      // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_string_map(?, ?, ?)", `{"25": {"other": "world"}}`, `greeting`, `hello`).IsEqualToJsonString(`{"25": {"other": "world", "greeting": "hello"}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_string_map(?, ?, ?)", `{"25": {"greeting": "hi"}}`, `greeting`, `hello`).IsEqualToJsonString(`{"25": {"greeting": "hello"}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_string_map(?, ?)", `{}`, `{"msg1": "hello", "msg2": "world"}`).IsEqualToJsonString(`{"25": {"msg1": "hello", "msg2": "world"}}`)                                               // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_string_map(?, ?)", `{"25": {"other": "test"}}`, `{"msg1": "hello", "msg2": "world"}`).IsEqualToJsonString(`{"25": {"other": "test", "msg1": "hello", "msg2": "world"}}`)                // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_string_map(?, ?)", `{"25": {"other": "test", "msg1": "hi"}}`, `{"msg1": "hello", "msg2": "world"}`).IsEqualToJsonString(`{"25": {"other": "test", "msg1": "hello", "msg2": "world"}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_string_map(?, ?)", `{"25": {"greeting": "hello", "other": "world"}}`, `greeting`).IsEqualToJsonString(`{"25": {"other": "world"}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_string_map(?, ?)", `{"25": {"greeting": "hello"}}`, `greeting`).IsEqualToJsonString(`{}`)                              // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_string_map(?, ?)", `{"25": {"greeting": "hello"}}`, `missing`).IsEqualToJsonString(`{"25": {"greeting": "hello"}}`)     // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_string_map(?, ?)", `{}`, `greeting`).IsEqualToJsonString(`{}`)                                                          // Remove from empty map
	})

	t.Run("bytes_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_bytes_map(?, JSON_OBJECT('empty', ''))", `{}`).IsEqualToJsonString(`{"26": {"empty": ""}}`)
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_bytes_map(?, JSON_OBJECT('data', 'aGVsbG8='))", `{}`).IsEqualToJsonString(`{"26": {"data": "aGVsbG8="}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_bytes_map(?)", `{"26": {"data": "aGVsbG8="}}`).IsEqualToJsonString(`{"data": "aGVsbG8="}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_bytes_map(?)", `{"26": {"data": "aGVsbG8="}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_bytes_map(?)", `{"26": {"data": "aGVsbG8="}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bytes_map__or(?, ?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `data`, ``).IsEqualToString(`hello`) // Key exists (base64 decoded)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bytes_map__or(?, ?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `empty`, ``).IsEqualToString(``)     // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bytes_map__or(?, ?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `missing`, ``).IsEqualToString(``)   // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bytes_map__or(?, ?, ?)", `{}`, `data`, ``).IsEqualToString(``)                                             // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bytes_map(?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `data`).IsEqualToBytes([]byte(`hello`)) // Key exists (base64 decoded)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bytes_map(?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `empty`).IsEqualToBytes([]byte{})     // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bytes_map(?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `missing`).IsNull()                    // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_bytes_map(?, ?)", `{}`, `data`).IsNull()                                                              // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_bytes_map(?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `data`).IsEqualToBool(true)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_bytes_map(?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `empty`).IsEqualToBool(true)  // Key exists (zero value)
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_bytes_map(?, ?)", `{"26": {"data": "aGVsbG8=", "empty": ""}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_bytes_map(?, ?)", `{}`, `data`).IsEqualToBool(false)                                            // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_bytes_map(?, ?, ?)", `{}`, `data`, []byte(`hello`)).IsEqualToJsonString(`{"26": {"data": "aGVsbG8="}}`)                      // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_bytes_map(?, ?, ?)", `{"26": {"other": "d29ybGQ="}}`, `data`, []byte(`hello`)).IsEqualToJsonString(`{"26": {"other": "d29ybGQ=", "data": "aGVsbG8="}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_bytes_map(?, ?, ?)", `{"26": {"data": "aGk="}}`, `data`, []byte(`hello`)).IsEqualToJsonString(`{"26": {"data": "aGVsbG8="}}`)           // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_bytes_map(?, ?)", `{}`, `{"msg1": "aGVsbG8=", "msg2": "d29ybGQ="}`).IsEqualToJsonString(`{"26": {"msg1": "aGVsbG8=", "msg2": "d29ybGQ="}}`)                                               // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_bytes_map(?, ?)", `{"26": {"other": "dGVzdA=="}}`, `{"msg1": "aGVsbG8=", "msg2": "d29ybGQ="}`).IsEqualToJsonString(`{"26": {"other": "dGVzdA==", "msg1": "aGVsbG8=", "msg2": "d29ybGQ="}}`)                // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_bytes_map(?, ?)", `{"26": {"other": "dGVzdA==", "msg1": "aGk="}}`, `{"msg1": "aGVsbG8=", "msg2": "d29ybGQ="}`).IsEqualToJsonString(`{"26": {"other": "dGVzdA==", "msg1": "aGVsbG8=", "msg2": "d29ybGQ="}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_bytes_map(?, ?)", `{"26": {"data": "aGVsbG8=", "other": "d29ybGQ="}}`, `data`).IsEqualToJsonString(`{"26": {"other": "d29ybGQ="}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_bytes_map(?, ?)", `{"26": {"data": "aGVsbG8="}}`, `data`).IsEqualToJsonString(`{}`)                              // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_bytes_map(?, ?)", `{"26": {"data": "aGVsbG8="}}`, `missing`).IsEqualToJsonString(`{"26": {"data": "aGVsbG8="}}`)     // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_bytes_map(?, ?)", `{}`, `data`).IsEqualToJsonString(`{}`)                                                          // Remove from empty map
	})

	t.Run("enum_value", func(t *testing.T) {
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_enum_map(?, JSON_OBJECT('status', 1))", `{}`).IsEqualToJsonString(`{"27": {"status": 1}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_enum_map(?)", `{"27": {"status": 1}}`).IsEqualToJsonString(`{"status": 1}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_enum_map(?)", `{"27": {"status": 1}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_enum_map(?)", `{"27": {"status": 1}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_enum_map__or(?, ?, ?)", `{"27": {"status": 1}}`, `status`, `0`).IsEqualToInt(1)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_enum_map__or(?, ?, ?)", `{"27": {"status": 1}}`, `missing`, `0`).IsEqualToInt(0) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_enum_map__or(?, ?, ?)", `{}`, `status`, `0`).IsEqualToInt(0)                     // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_enum_map(?, ?)", `{"27": {"status": 1}}`, `status`).IsEqualToInt(1)  // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_enum_map(?, ?)", `{"27": {"status": 1}}`, `missing`).IsNull() // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_enum_map(?, ?)", `{}`, `status`).IsNull()                     // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_enum_map(?, ?)", `{"27": {"status": 1}}`, `status`).IsEqualToBool(true)   // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_enum_map(?, ?)", `{"27": {"status": 1}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_enum_map(?, ?)", `{}`, `status`).IsEqualToBool(false)                     // Map empty

		// Test single key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_enum_map(?, ?, ?)", `{}`, `active_status`, 1).IsEqualToJsonString(`{"27": {"active_status": 1}}`)                                     // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_enum_map(?, ?, ?)", `{"27": {"status": 1}}`, `inactive_status`, 2).IsEqualToJsonString(`{"27": {"status": 1, "inactive_status": 2}}`) // Add to existing map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_enum_map(?, ?, ?)", `{"27": {"status": 1}}`, `status`, 2).IsEqualToJsonString(`{"27": {"status": 2}}`)                                // Update existing key

		// Test bulk key insertion
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_enum_map(?, ?)", `{}`, `{"active": 1, "inactive": 2}`).IsEqualToJsonString(`{"27": {"active": 1, "inactive": 2}}`)                                              // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_enum_map(?, ?)", `{"27": {"status": 0}}`, `{"active": 1, "inactive": 2}`).IsEqualToJsonString(`{"27": {"status": 0, "active": 1, "inactive": 2}}`)              // Merge with existing
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_enum_map(?, ?)", `{"27": {"status": 0, "active": 1}}`, `{"active": 2, "inactive": 2}`).IsEqualToJsonString(`{"27": {"status": 0, "active": 2, "inactive": 2}}`) // Update existing keys

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_enum_map(?, ?)", `{"27": {"status": 1, "active": 1}}`, `status`).IsEqualToJsonString(`{"27": {"active": 1}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_enum_map(?, ?)", `{"27": {"status": 1}}`, `status`).IsEqualToJsonString(`{}`)                                 // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_enum_map(?, ?)", `{"27": {"status": 1}}`, `missing`).IsEqualToJsonString(`{"27": {"status": 1}}`)             // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_enum_map(?, ?)", `{}`, `status`).IsEqualToJsonString(`{}`)                                                    // Remove from empty map
	})

	t.Run("message_value", func(t *testing.T) {
		// Test message value
		RunTestThatExpression(t, "pbt_map_fields_set_all_string_to_message_map(?, JSON_OBJECT('nested', JSON_OBJECT('1', 'test', '2', 42)))", `{}`).IsEqualToJsonString(`{"28": {"nested": {"1": "test", "2": 42}}}`)
		RunTestThatExpression(t, "pbt_map_fields_get_all_string_to_message_map(?)", `{"28": {"nested": {"1": "test", "2": 42}}}`).IsEqualToJsonString(`{"nested": {"1": "test", "2": 42}}`)
		RunTestThatExpression(t, "pbt_map_fields_count_string_to_message_map(?)", `{"28": {"nested": {"1": "test", "2": 42}}}`).IsEqualToInt(1)
		RunTestThatExpression(t, "pbt_map_fields_clear_string_to_message_map(?)", `{"28": {"nested": {"1": "test", "2": 42}}}`).IsEqualToJsonString(`{}`)

		// Test individual key access with default (__or variant)
		defaultMessage := `{"1": "default", "2": 0}`
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_message_map__or(?, ?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `nested`, defaultMessage).IsEqualToJsonString(`{"1": "test", "2": 42}`) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_message_map__or(?, ?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `missing`, defaultMessage).IsEqualToJsonString(defaultMessage)          // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_message_map__or(?, ?, ?)", `{}`, `nested`, defaultMessage).IsEqualToJsonString(defaultMessage)                                                   // Map empty

		// Test individual key access without default
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_message_map(?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `nested`).IsEqualToJsonString(`{"1": "test", "2": 42}`) // Key exists
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_message_map(?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `missing`).IsNull()            // Key missing
		RunTestThatExpression(t, "pbt_map_fields_get_string_to_message_map(?, ?)", `{}`, `nested`).IsNull()                                                     // Map empty

		// Test key existence checks
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_message_map(?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `nested`).IsEqualToBool(true)   // Key exists
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_message_map(?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `missing`).IsEqualToBool(false) // Key missing
		RunTestThatExpression(t, "pbt_map_fields_contains_string_to_message_map(?, ?)", `{}`, `nested`).IsEqualToBool(false)                                          // Map empty

		// Test single key insertion
		newMessage := `{"1": "hello", "2": 123}`
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_message_map(?, ?, ?)", `{}`, `new_nested`, newMessage).IsEqualToJsonString(`{"28": {"new_nested": {"1": "hello", "2": 123}}}`)                                                                           // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_message_map(?, ?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `new_nested`, newMessage).IsEqualToJsonString(`{"28": {"nested": {"1": "test", "2": 42}, "new_nested": {"1": "hello", "2": 123}}}`) // Add to existing map
		updatedMessage := `{"1": "updated", "2": 999}`
		RunTestThatExpression(t, "pbt_map_fields_put_string_to_message_map(?, ?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `nested`, updatedMessage).IsEqualToJsonString(`{"28": {"nested": {"1": "updated", "2": 999}}}`) // Update existing key

		// Test bulk key insertion
		bulkMessages := `{"msg1": {"1": "first", "2": 1}, "msg2": {"1": "second", "2": 2}}`
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_message_map(?, ?)", `{}`, bulkMessages).IsEqualToJsonString(`{"28": {"msg1": {"1": "first", "2": 1}, "msg2": {"1": "second", "2": 2}}}`)                                                                // Add to empty map
		RunTestThatExpression(t, "pbt_map_fields_put_all_string_to_message_map(?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `{"msg1": {"1": "first", "2": 1}}`).IsEqualToJsonString(`{"28": {"nested": {"1": "test", "2": 42}, "msg1": {"1": "first", "2": 1}}}`) // Merge with existing

		// Test key removal
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_message_map(?, ?)", `{"28": {"nested": {"1": "test", "2": 42}, "other": {"1": "other", "2": 1}}}`, `nested`).IsEqualToJsonString(`{"28": {"other": {"1": "other", "2": 1}}}`) // Remove existing key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_message_map(?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `nested`).IsEqualToJsonString(`{}`)                                                                         // Remove last key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_message_map(?, ?)", `{"28": {"nested": {"1": "test", "2": 42}}}`, `missing`).IsEqualToJsonString(`{"28": {"nested": {"1": "test", "2": 42}}}`)                                // Remove non-existent key
		RunTestThatExpression(t, "pbt_map_fields_remove_string_to_message_map(?, ?)", `{}`, `nested`).IsEqualToJsonString(`{}`)                                                                                                                 // Remove from empty map
	})
}
