package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	"github.com/eiiches/mysql-protobuf-functions/internal/protocgenmysql"
	"github.com/eiiches/mysql-protobuf-functions/internal/testutils"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// generateAndLoadOpaqueApiSQL generates SQL from protobuf definition and loads it into MySQL
func generateAndLoadOpaqueApiSQL(t *testing.T, protoContent string, schemaName string) *testutils.ProtoTestSupport {
	g := NewWithT(t)
	g.THelper()

	// Create protobuf definitions using the same pattern as existing tests
	support := testutils.NewProtoTestSupport(t, map[string]string{
		"test.proto": protoContent,
	})

	// Get FileDescriptorSet for code generation
	fds := support.GetFileDescriptorSet()

	// Configure the SQL code generator
	config := protocgenmysql.GenerateConfig{
		DescriptorSetName: schemaName,
		GenerateMethods:   true,
		IncludeWkt:        true,
		FileNameFunc: func(protoPath string) string {
			return "" // Single file output
		},
		TypePrefixFunc: func(pkg protoreflect.FullName, typeName protoreflect.FullName) string {
			// Convert message names to snake_case function prefixes
			name := string(typeName)
			if strings.HasPrefix(name, ".") {
				name = name[1:] // Remove leading dot
			}
			name = strings.ReplaceAll(name, ".", "_")
			return strings.ToLower(name)
		},
	}

	// Generate SQL code using protocgenmysql
	response, err := protocgenmysql.Generate(fds, config)
	g.Expect(err).NotTo(HaveOccurred(), "Failed to generate SQL from protobuf definition")
	g.Expect(response.File).To(HaveLen(1), "Expected exactly one generated SQL file")

	// Extract generated SQL content
	sqlContent := response.File[0].GetContent()
	g.Expect(sqlContent).NotTo(BeEmpty(), "Generated SQL content should not be empty")

	// Load the generated SQL into MySQL
	executeSQLStatements(t, sqlContent)

	return support
}

// executeSQLStatements executes a sequence of SQL statements against the test database
func executeSQLStatements(t *testing.T, sqlContent string) {
	g := NewWithT(t)
	g.THelper()

	// Split SQL content by delimiter markers to get individual statements
	statements := strings.Split(sqlContent, "$$")

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)

		// Skip empty statements, delimiter directives, and comments
		if stmt == "" || stmt == "DELIMITER $$" || strings.HasPrefix(stmt, "--") {
			continue
		}

		// Execute the statement
		t.Logf("Executing SQL statement %d: %s", i, stmt)
		_, err := db.Exec(stmt)
		g.Expect(err).NotTo(HaveOccurred(), "Failed to execute SQL statement: %s", stmt)
	}
}

// truncateForLog truncates long strings for cleaner log output
func truncateForLog(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// testBasicFieldOperations tests get/set/clear/has operations for a single field
func testBasicFieldOperations(t *testing.T, fieldDef, fieldName, typeName string, testValue, defaultValue interface{}) {
	t.Helper()

	protoContent := fmt.Sprintf(dedent.Pipe(`
		|syntax = "proto3";
		|message Test {
		|    %s
		|}
		|message MessageType {
		|    int32 value = 1;
		|}
		|enum EnumType {
		|    ENUM_TYPE_UNSPECIFIED = 0;
		|    ENUM_TYPE_ONE = 1;
		|}
	`), fieldDef)

	// Generate and load SQL
	schemaName := "test_schema"
	generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

	// Test constructor
	t.Run("constructor", func(t *testing.T) {
		RunTestThatExpression(t, "test_new()").IsEqualToJsonString("{}")
	})

	// Test setter and getter
	t.Run("setter_and_getter", func(t *testing.T) {
		setterFunc := fmt.Sprintf("test_set_%s", fieldName)
		getterFunc := fmt.Sprintf("test_get_%s", fieldName)

		// Set value and verify it can be retrieved
		expr := fmt.Sprintf("%s(%s(test_new(), ?))", getterFunc, setterFunc)
		switch v := testValue.(type) {
		case string:
			RunTestThatExpression(t, expr, v).IsEqualToString(v)
		case int:
			RunTestThatExpression(t, expr, v).IsEqualToInt(int64(v))
		case int32:
			RunTestThatExpression(t, expr, v).IsEqualToInt(int64(v))
		case int64:
			RunTestThatExpression(t, expr, v).IsEqualToInt(v)
		case bool:
			RunTestThatExpression(t, expr, v).IsEqualToBool(v)
		default:
			RunTestThatExpression(t, expr, v).IsEqualTo(v)
		}
	})

	// Test clear
	t.Run("clear", func(t *testing.T) {
		setterFunc := fmt.Sprintf("test_set_%s", fieldName)
		clearFunc := fmt.Sprintf("test_clear_%s", fieldName)
		getterFunc := fmt.Sprintf("test_get_%s", fieldName)

		// Set value, clear it, then verify default is returned
		clearExpr := fmt.Sprintf("%s(%s(%s(test_new(), ?)))", getterFunc, clearFunc, setterFunc)
		switch def := defaultValue.(type) {
		case string:
			RunTestThatExpression(t, clearExpr, testValue).IsEqualToString(def)
		case int:
			RunTestThatExpression(t, clearExpr, testValue).IsEqualToInt(int64(def))
		case int32:
			RunTestThatExpression(t, clearExpr, testValue).IsEqualToInt(int64(def))
		case int64:
			RunTestThatExpression(t, clearExpr, testValue).IsEqualToInt(def)
		case bool:
			RunTestThatExpression(t, clearExpr, testValue).IsEqualToBool(def)
		default:
			RunTestThatExpression(t, clearExpr, testValue).IsEqualTo(defaultValue)
		}
	})

	// Test storage format (for JSON fields like boolean)
	if typeName == "bool" {
		t.Run("boolean_storage_format", func(t *testing.T) {
			setterFunc := fmt.Sprintf("test_set_%s", fieldName)
			// Verify boolean values are stored as true/false, not 1/0
			// Proto3 without presence: FALSE omitted, TRUE stored
			RunTestThatExpression(t, fmt.Sprintf("%s(test_new(), TRUE)", setterFunc)).IsEqualToJsonString(fmt.Sprintf(`{"1": true}`))
			RunTestThatExpression(t, fmt.Sprintf("%s(test_new(), FALSE)", setterFunc)).IsEqualToJsonString(`{}`)
		})
	}
}

// TestGeneratedOpaqueApiInternalRepresentation tests that setter functions create correct protonumberjson format for all protobuf types
func TestGeneratedOpaqueApiInternalRepresentation(t *testing.T) {
	// Test field number keys (crucial for protonumberjson format)
	t.Run("field_number_keys", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    string first = 10;     // field number 10
			|    string second = 5;     // field number 5  
			|    string third = 100;    // field number 100
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		// Verify field numbers are used as JSON keys, not field names
		obj := "test_set_third(test_set_second(test_set_first(test_new(), 'value10'), 'value5'), 'value100')"
		RunTestThatExpression(t, obj).IsEqualToJsonString(`{"10": "value10", "5": "value5", "100": "value100"}`)
	})

	// Test complex combined scenario
	t.Run("complex_combined_scenario", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    string name = 1;
			|    int32 age = 2;
			|    bool active = 3;
			|    float score = 4;
			|    repeated int32 numbers = 5;
			|    bytes data = 6;
			|    Status status = 7;
			|}
			|enum Status {
			|    STATUS_UNSPECIFIED = 0;
			|    STATUS_ACTIVE = 1;
			|    STATUS_INACTIVE = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		// Build complex object with multiple field types
		obj := "test_set_status(test_set_data(test_add_numbers(test_add_numbers(test_set_score(test_set_active(test_set_age(test_set_name(test_new(), 'John'), 25), TRUE), 3.14), 1), 2), ?), 1)"
		expected := `{"1": "John", "2": 25, "3": true, "4": "binary32:0x4048f5c3", "5": [1, 2], "6": "dGVzdA==", "7": 1}`
		RunTestThatExpression(t, obj, []byte("test")).IsEqualToJsonString(expected)
	})
}

// TestGeneratedOpaqueApiBasicFields tests basic field operations for all protobuf field types
func TestGeneratedOpaqueApiBasicFields(t *testing.T) {
	t.Run("int32 field", func(t *testing.T) {
		testBasicFieldOperations(t, "int32 value = 1;", "value", "int32", 42, 0)
	})

	t.Run("int64 field", func(t *testing.T) {
		testBasicFieldOperations(t, "int64 value = 1;", "value", "int64", int64(9223372036854775807), int64(0))
	})

	t.Run("uint32 field", func(t *testing.T) {
		testBasicFieldOperations(t, "uint32 value = 1;", "value", "uint32", 4294967295, 0)
	})

	t.Run("uint64 field", func(t *testing.T) {
		testBasicFieldOperations(t, "uint64 value = 1;", "value", "uint64", int64(9223372036854775807), int64(0))
	})

	t.Run("fixed32 field", func(t *testing.T) {
		testBasicFieldOperations(t, "fixed32 value = 1;", "value", "fixed32", 4294967295, 0)
	})

	t.Run("fixed64 field", func(t *testing.T) {
		testBasicFieldOperations(t, "fixed64 value = 1;", "value", "fixed64", int64(9223372036854775807), int64(0))
	})

	t.Run("sfixed32 field", func(t *testing.T) {
		testBasicFieldOperations(t, "sfixed32 value = 1;", "value", "sfixed32", -2147483648, 0)
	})

	t.Run("sfixed64 field", func(t *testing.T) {
		testBasicFieldOperations(t, "sfixed64 value = 1;", "value", "sfixed64", int64(-9223372036854775808), int64(0))
	})

	t.Run("sint32 field", func(t *testing.T) {
		testBasicFieldOperations(t, "sint32 value = 1;", "value", "sint32", -2147483648, 0)
	})

	t.Run("sint64 field", func(t *testing.T) {
		testBasicFieldOperations(t, "sint64 value = 1;", "value", "sint64", int64(-9223372036854775808), int64(0))
	})

	t.Run("bool field", func(t *testing.T) {
		testBasicFieldOperations(t, "bool flag = 1;", "flag", "bool", true, false)
	})

	t.Run("string field", func(t *testing.T) {
		testBasicFieldOperations(t, "string name = 1;", "name", "string", "hello world", "")
	})

	t.Run("bytes field", func(t *testing.T) {
		testBasicFieldOperations(t, "bytes data = 1;", "data", "bytes", []byte("hello world"), []byte{})
	})

	t.Run("float field", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    float value = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("constructor", func(t *testing.T) {
			RunTestThatExpression(t, "test_new()").IsEqualToJsonString("{}")
		})

		t.Run("setter_and_getter", func(t *testing.T) {
			// Float getters now return the actual float value
			RunTestThatExpression(t, "test_get_value(test_set_value(test_new(), 3.14))").IsEqualToFloat(3.14)
		})
	})

	t.Run("double field", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    double value = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("constructor", func(t *testing.T) {
			RunTestThatExpression(t, "test_new()").IsEqualToJsonString("{}")
		})

		t.Run("setter_and_getter", func(t *testing.T) {
			// Double getters now return the actual double value
			RunTestThatExpression(t, "test_get_value(test_set_value(test_new(), 3.141592653589793))").IsEqualToFloat(3.141592653589793)
		})
	})

	t.Run("enum field", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    EnumType status = 1;
			|}
			|enum EnumType {
			|    ENUM_TYPE_UNSPECIFIED = 0;
			|    ENUM_TYPE_ONE = 1;
			|    ENUM_TYPE_TWO = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("constructor", func(t *testing.T) {
			RunTestThatExpression(t, "test_new()").IsEqualToJsonString("{}")
		})

		t.Run("setter_and_getter", func(t *testing.T) {
			RunTestThatExpression(t, "test_get_status(test_set_status(test_new(), 1))").IsEqualToInt(1)
			RunTestThatExpression(t, "test_get_status(test_set_status(test_new(), 2))").IsEqualToInt(2)
		})
	})

	t.Run("message field", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    MessageType nested = 1;
			|}
			|message MessageType {
			|    int32 value = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("constructor", func(t *testing.T) {
			RunTestThatExpression(t, "test_new()").IsEqualToJsonString("{}")
		})

		t.Run("setter_and_getter", func(t *testing.T) {
			// Set nested message field
			nestedObj := `{"1": 42}`
			RunTestThatExpression(t, "test_get_nested(test_set_nested(test_new(), JSON_OBJECT('1', 42)))").IsEqualToJsonString(nestedObj)
		})
	})
}

// TestGeneratedOpaqueApiRepeatedFields tests repeated field operations
func TestGeneratedOpaqueApiRepeatedFields(t *testing.T) {
	t.Run("repeated int32", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    repeated int32 items = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("empty_array", func(t *testing.T) {
			RunTestThatExpression(t, "test_count_items(test_new())").IsEqualToInt(0)
			RunTestThatExpression(t, "test_get_all_items(test_new())").IsEqualToJsonString("[]")
		})

		t.Run("add_elements", func(t *testing.T) {
			// Add single element
			RunTestThatExpression(t, "test_count_items(test_add_items(test_new(), 42))").IsEqualToInt(1)
			RunTestThatExpression(t, "test_get_all_items(test_add_items(test_new(), 42))").IsEqualToJsonString("[42]")

			// Add multiple elements
			obj := "test_add_items(test_add_items(test_add_items(test_new(), 1), 2), 3)"
			RunTestThatExpression(t, fmt.Sprintf("test_count_items(%s)", obj)).IsEqualToInt(3)
			RunTestThatExpression(t, fmt.Sprintf("test_get_all_items(%s)", obj)).IsEqualToJsonString("[1, 2, 3]")
		})

		t.Run("set_entire_array", func(t *testing.T) {
			RunTestThatExpression(t, "test_get_all_items(test_set_all_items(test_new(), JSON_ARRAY(10, 20, 30)))").IsEqualToJsonString("[10, 20, 30]")
		})

		t.Run("clear_array", func(t *testing.T) {
			obj := "test_add_items(test_add_items(test_new(), 1), 2)"
			RunTestThatExpression(t, fmt.Sprintf("test_count_items(test_clear_items(%s))", obj)).IsEqualToInt(0)
			RunTestThatExpression(t, fmt.Sprintf("test_get_all_items(test_clear_items(%s))", obj)).IsEqualToJsonString("[]")
		})
	})

	t.Run("repeated string", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    repeated string names = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("add_string_elements", func(t *testing.T) {
			obj := "test_add_names(test_add_names(test_new(), 'alice'), 'bob')"
			RunTestThatExpression(t, fmt.Sprintf("test_count_names(%s)", obj)).IsEqualToInt(2)
			RunTestThatExpression(t, fmt.Sprintf("test_get_all_names(%s)", obj)).IsEqualToJsonString(`["alice", "bob"]`)
		})
	})

	t.Run("repeated bool", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    repeated bool flags = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("add_bool_elements", func(t *testing.T) {
			obj := "test_add_flags(test_add_flags(test_new(), TRUE), FALSE)"
			RunTestThatExpression(t, fmt.Sprintf("test_count_flags(%s)", obj)).IsEqualToInt(2)
			RunTestThatExpression(t, fmt.Sprintf("test_get_all_flags(%s)", obj)).IsEqualToJsonString(`[true, false]`)
		})
	})

	t.Run("repeated message", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    repeated MessageType items = 1;
			|}
			|message MessageType {
			|    int32 value = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("add_message_elements", func(t *testing.T) {
			obj1 := "JSON_OBJECT('1', 42)"
			obj2 := "JSON_OBJECT('1', 100)"
			expr := fmt.Sprintf("test_add_items(test_add_items(test_new(), %s), %s)", obj1, obj2)
			RunTestThatExpression(t, fmt.Sprintf("test_count_items(%s)", expr)).IsEqualToInt(2)
			RunTestThatExpression(t, fmt.Sprintf("test_get_all_items(%s)", expr)).IsEqualToJsonString(`[{"1": 42}, {"1": 100}]`)
		})
	})
}

// TestGeneratedOpaqueApiEnums tests enum string/number conversion
func TestGeneratedOpaqueApiEnums(t *testing.T) {
	protoContent := dedent.Pipe(`
		|syntax = "proto3";
		|message Test {
		|    Status status = 1;
		|}
		|enum Status {
		|    STATUS_UNSPECIFIED = 0;
		|    STATUS_PENDING = 1;
		|    STATUS_RUNNING = 2;
		|    STATUS_COMPLETED = 3;
		|    STATUS_FAILED = 4;
		|}
	`)
	schemaName := "test_schema"
	generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

	t.Run("from_string_conversion", func(t *testing.T) {
		RunTestThatExpression(t, "status_from_string('STATUS_UNSPECIFIED')").IsEqualToInt(0)
		RunTestThatExpression(t, "status_from_string('STATUS_PENDING')").IsEqualToInt(1)
		RunTestThatExpression(t, "status_from_string('STATUS_RUNNING')").IsEqualToInt(2)
		RunTestThatExpression(t, "status_from_string('STATUS_COMPLETED')").IsEqualToInt(3)
		RunTestThatExpression(t, "status_from_string('STATUS_FAILED')").IsEqualToInt(4)
	})

	t.Run("to_string_conversion", func(t *testing.T) {
		RunTestThatExpression(t, "status_to_string(0)").IsEqualToString("STATUS_UNSPECIFIED")
		RunTestThatExpression(t, "status_to_string(1)").IsEqualToString("STATUS_PENDING")
		RunTestThatExpression(t, "status_to_string(2)").IsEqualToString("STATUS_RUNNING")
		RunTestThatExpression(t, "status_to_string(3)").IsEqualToString("STATUS_COMPLETED")
		RunTestThatExpression(t, "status_to_string(4)").IsEqualToString("STATUS_FAILED")
	})

	t.Run("unknown_values", func(t *testing.T) {
		// Unknown enum name should return NULL
		RunTestThatExpression(t, "status_from_string('UNKNOWN_STATUS')").IsNull()

		// Unknown enum number should return NULL
		RunTestThatExpression(t, "status_to_string(999)").IsNull()
		RunTestThatExpression(t, "status_to_string(-1)").IsNull()
	})

	t.Run("round_trip_conversion", func(t *testing.T) {
		// String -> Number -> String should be identity
		RunTestThatExpression(t, "status_to_string(status_from_string('STATUS_RUNNING'))").IsEqualToString("STATUS_RUNNING")

		// Number -> String -> Number should be identity
		RunTestThatExpression(t, "status_from_string(status_to_string(3))").IsEqualToInt(3)
	})

	t.Run("enum_field_usage", func(t *testing.T) {
		// Test enum field in message using numeric values
		obj := "test_set_status(test_new(), 2)"
		RunTestThatExpression(t, fmt.Sprintf("test_get_status(%s)", obj)).IsEqualToInt(2)

		// Test using enum conversion functions
		obj2 := "test_set_status(test_new(), status_from_string('STATUS_FAILED'))"
		RunTestThatExpression(t, fmt.Sprintf("test_get_status(%s)", obj2)).IsEqualToInt(4)
		RunTestThatExpression(t, fmt.Sprintf("status_to_string(test_get_status(%s))", obj2)).IsEqualToString("STATUS_FAILED")
	})
}

// TestGeneratedOpaqueApiMapFields tests map field operations for all key and value types
func TestGeneratedOpaqueApiMapFields(t *testing.T) {
	// Test all possible key types with a fixed value type (int32)
	t.Run("map_fields_all_key_types", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
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
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		// Test all key types with int32 values - use JSON objects for map values
		RunTestThatExpression(t, "test_set_all_int32_to_int32_map(test_new(), JSON_OBJECT('42', 100))").IsEqualToJsonString(`{"1": {"42": 100}}`)
		RunTestThatExpression(t, "test_set_all_int64_to_int32_map(test_new(), JSON_OBJECT('9223372036854775807', 200))").IsEqualToJsonString(`{"2": {"9223372036854775807": 200}}`)
		RunTestThatExpression(t, "test_set_all_uint32_to_int32_map(test_new(), JSON_OBJECT('4294967295', 300))").IsEqualToJsonString(`{"3": {"4294967295": 300}}`)
		RunTestThatExpression(t, "test_set_all_uint64_to_int32_map(test_new(), JSON_OBJECT('18446744073709551615', 400))").IsEqualToJsonString(`{"4": {"18446744073709551615": 400}}`)
		RunTestThatExpression(t, "test_set_all_sint32_to_int32_map(test_new(), JSON_OBJECT('-1', 500))").IsEqualToJsonString(`{"5": {"-1": 500}}`)
		RunTestThatExpression(t, "test_set_all_sint64_to_int32_map(test_new(), JSON_OBJECT('-9223372036854775808', 600))").IsEqualToJsonString(`{"6": {"-9223372036854775808": 600}}`)
		RunTestThatExpression(t, "test_set_all_fixed32_to_int32_map(test_new(), JSON_OBJECT('4294967295', 700))").IsEqualToJsonString(`{"7": {"4294967295": 700}}`)
		RunTestThatExpression(t, "test_set_all_fixed64_to_int32_map(test_new(), JSON_OBJECT('18446744073709551615', 800))").IsEqualToJsonString(`{"8": {"18446744073709551615": 800}}`)
		RunTestThatExpression(t, "test_set_all_sfixed32_to_int32_map(test_new(), JSON_OBJECT('-2147483648', 900))").IsEqualToJsonString(`{"9": {"-2147483648": 900}}`)
		RunTestThatExpression(t, "test_set_all_sfixed64_to_int32_map(test_new(), JSON_OBJECT('-9223372036854775808', 1000))").IsEqualToJsonString(`{"10": {"-9223372036854775808": 1000}}`)
		RunTestThatExpression(t, "test_set_all_bool_to_int32_map(test_new(), JSON_OBJECT('true', 1100))").IsEqualToJsonString(`{"11": {"true": 1100}}`)
		RunTestThatExpression(t, "test_set_all_string_to_int32_map(test_new(), JSON_OBJECT('key', 1200))").IsEqualToJsonString(`{"12": {"key": 1200}}`)

		// Test multiple entries in same map
		RunTestThatExpression(t, "test_set_all_string_to_int32_map(test_new(), JSON_OBJECT('first', 10, 'second', 20))").IsEqualToJsonString(`{"12": {"first": 10, "second": 20}}`)

		// Test overwriting existing map (replaces entire map)
		RunTestThatExpression(t, "test_set_all_string_to_int32_map(test_set_all_string_to_int32_map(test_new(), JSON_OBJECT('old', 100)), JSON_OBJECT('new', 200))").IsEqualToJsonString(`{"12": {"new": 200}}`)
	})

	// Test all possible value types with a fixed key type (string)
	t.Run("map_fields_all_value_types", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    map<string, double> string_to_double_map = 1;
			|    map<string, float> string_to_float_map = 2;
			|    map<string, int32> string_to_int32_map = 3;
			|    map<string, int64> string_to_int64_map = 4;
			|    map<string, uint32> string_to_uint32_map = 5;
			|    map<string, uint64> string_to_uint64_map = 6;
			|    map<string, sint32> string_to_sint32_map = 7;
			|    map<string, sint64> string_to_sint64_map = 8;
			|    map<string, fixed32> string_to_fixed32_map = 9;
			|    map<string, fixed64> string_to_fixed64_map = 10;
			|    map<string, sfixed32> string_to_sfixed32_map = 11;
			|    map<string, sfixed64> string_to_sfixed64_map = 12;
			|    map<string, bool> string_to_bool_map = 13;
			|    map<string, string> string_to_string_map = 14;
			|    map<string, bytes> string_to_bytes_map = 15;
			|    map<string, Status> string_to_enum_map = 16;
			|    map<string, Nested> string_to_message_map = 17;
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

		// Test all value types with string keys - use JSON objects for map values
		RunTestThatExpression(t, "test_set_all_string_to_double_map(test_new(), JSON_OBJECT('pi', 'binary64:0x400921fb54442d18'))").IsEqualToJsonString(`{"1": {"pi": "binary64:0x400921fb54442d18"}}`)
		RunTestThatExpression(t, "test_set_all_string_to_float_map(test_new(), JSON_OBJECT('pi_float', 'binary32:0x4048f5c3'))").IsEqualToJsonString(`{"2": {"pi_float": "binary32:0x4048f5c3"}}`)
		RunTestThatExpression(t, "test_set_all_string_to_int32_map(test_new(), JSON_OBJECT('answer', 42))").IsEqualToJsonString(`{"3": {"answer": 42}}`)
		RunTestThatExpression(t, "test_set_all_string_to_int64_map(test_new(), JSON_OBJECT('big', 9223372036854775807))").IsEqualToJsonString(`{"4": {"big": 9223372036854775807}}`)
		RunTestThatExpression(t, "test_set_all_string_to_uint32_map(test_new(), JSON_OBJECT('max32', 4294967295))").IsEqualToJsonString(`{"5": {"max32": 4294967295}}`)
		RunTestThatExpression(t, "test_set_all_string_to_uint64_map(test_new(), JSON_OBJECT('max64', 18446744073709551615))").IsEqualToJsonString(`{"6": {"max64": 18446744073709551615}}`)
		RunTestThatExpression(t, "test_set_all_string_to_sint32_map(test_new(), JSON_OBJECT('neg', -1))").IsEqualToJsonString(`{"7": {"neg": -1}}`)
		RunTestThatExpression(t, "test_set_all_string_to_sint64_map(test_new(), JSON_OBJECT('min64', -9223372036854775808))").IsEqualToJsonString(`{"8": {"min64": -9223372036854775808}}`)
		RunTestThatExpression(t, "test_set_all_string_to_fixed32_map(test_new(), JSON_OBJECT('fixed', 4294967295))").IsEqualToJsonString(`{"9": {"fixed": 4294967295}}`)
		RunTestThatExpression(t, "test_set_all_string_to_fixed64_map(test_new(), JSON_OBJECT('fixed64', 18446744073709551615))").IsEqualToJsonString(`{"10": {"fixed64": 18446744073709551615}}`)
		RunTestThatExpression(t, "test_set_all_string_to_sfixed32_map(test_new(), JSON_OBJECT('sfixed', -2147483648))").IsEqualToJsonString(`{"11": {"sfixed": -2147483648}}`)
		RunTestThatExpression(t, "test_set_all_string_to_sfixed64_map(test_new(), JSON_OBJECT('sfixed64', -9223372036854775808))").IsEqualToJsonString(`{"12": {"sfixed64": -9223372036854775808}}`)
		RunTestThatExpression(t, "test_set_all_string_to_bool_map(test_new(), JSON_OBJECT('flag', true))").IsEqualToJsonString(`{"13": {"flag": true}}`)
		RunTestThatExpression(t, "test_set_all_string_to_string_map(test_new(), JSON_OBJECT('greeting', 'hello'))").IsEqualToJsonString(`{"14": {"greeting": "hello"}}`)
		RunTestThatExpression(t, "test_set_all_string_to_bytes_map(test_new(), JSON_OBJECT('data', 'aGVsbG8='))").IsEqualToJsonString(`{"15": {"data": "aGVsbG8="}}`)
		RunTestThatExpression(t, "test_set_all_string_to_enum_map(test_new(), JSON_OBJECT('status', 1))").IsEqualToJsonString(`{"16": {"status": 1}}`)

		// Test message value
		RunTestThatExpression(t, "test_set_all_string_to_message_map(test_new(), JSON_OBJECT('nested', JSON_OBJECT('1', 'test', '2', 42)))").IsEqualToJsonString(`{"17": {"nested": {"1": "test", "2": 42}}}`)
	})

	// Test map getter functions (return entire map)
	t.Run("map_fields_getters", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    map<string, int32> data = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		// Test getter returns entire map
		obj := "test_set_all_data(test_new(), JSON_OBJECT('key1', 10, 'key2', 20))"
		RunTestThatExpression(t, fmt.Sprintf("test_get_all_data(%s)", obj)).IsEqualToJsonString(`{"key1": 10, "key2": 20}`)

		// Test getter returns empty for new object
		RunTestThatExpression(t, "test_get_all_data(test_new())").IsEqualToJsonString("[]")
	})

	// Test map operations (clear, count)
	t.Run("map_fields_operations", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    map<string, int32> data = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		// Test map with multiple entries
		obj := "test_set_all_data(test_new(), JSON_OBJECT('key1', 10, 'key2', 20))"

		// Test map count operation
		RunTestThatExpression(t, "test_count_data(test_new())").IsEqualToInt(0)
		RunTestThatExpression(t, fmt.Sprintf("test_count_data(%s)", obj)).IsEqualToInt(2)

		// Test map clear operation
		objAfterClear := fmt.Sprintf("test_clear_data(%s)", obj)
		RunTestThatExpression(t, fmt.Sprintf("test_count_data(%s)", objAfterClear)).IsEqualToInt(0)
		RunTestThatExpression(t, fmt.Sprintf("test_get_all_data(%s)", objAfterClear)).IsEqualToJsonString("[]")
	})

	// Test map internal representation format
	t.Run("map_fields_internal_format", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    map<int32, string> int_to_string_map = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		// Verify that maps are stored as nested JSON objects
		// Map field 1 contains a JSON object where keys are string representations of the map keys
		RunTestThatExpression(t, "test_set_all_int_to_string_map(test_new(), JSON_OBJECT('1', 'one', '2', 'two'))").IsEqualToJsonString(`{"1": {"1": "one", "2": "two"}}`)

		// Verify field number keys are used for the map field itself
		RunTestThatExpression(t, "test_set_all_int_to_string_map(test_new(), JSON_OBJECT('42', 'answer'))").IsEqualToJsonString(`{"1": {"42": "answer"}}`)
	})

	// Test empty and default values in maps
	t.Run("map_fields_default_values", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    map<string, int32> int_map = 1;
			|    map<string, string> string_map = 2;
			|    map<string, bool> bool_map = 3;
			|    map<string, bytes> bytes_map = 4;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		// Maps can store default/zero values (unlike regular proto3 fields without presence)
		RunTestThatExpression(t, "test_set_all_int_map(test_new(), JSON_OBJECT('zero', 0))").IsEqualToJsonString(`{"1": {"zero": 0}}`)
		RunTestThatExpression(t, "test_set_all_string_map(test_new(), JSON_OBJECT('empty', ''))").IsEqualToJsonString(`{"2": {"empty": ""}}`)
		RunTestThatExpression(t, "test_set_all_bool_map(test_new(), JSON_OBJECT('false', false))").IsEqualToJsonString(`{"3": {"false": false}}`)
		RunTestThatExpression(t, "test_set_all_bytes_map(test_new(), JSON_OBJECT('empty', ''))").IsEqualToJsonString(`{"4": {"empty": ""}}`)
	})
}

// TestGeneratedOpaqueApiConversions tests JSON conversion functionality
func TestGeneratedOpaqueApiConversions(t *testing.T) {
	t.Run("simple_message_roundtrip", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    int32 value = 1;
			|    string name = 2;
			|    bool flag = 3;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("json_roundtrip", func(t *testing.T) {
			// Create a message using the opaque API
			obj := "test_set_flag(test_set_name(test_set_value(test_new(), 42), 'hello'), TRUE)"

			// Convert to standard JSON format (regular protobuf JSON)
			jsonExpr := fmt.Sprintf("test_to_json(%s, NULL)", obj)

			// Test round-trip via JSON conversion
			backToOpaque := fmt.Sprintf("test_from_json(%s, NULL)", jsonExpr)
			RunTestThatExpression(t, fmt.Sprintf("test_get_value(%s)", backToOpaque)).IsEqualToInt(42)
			RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", backToOpaque)).IsEqualToString("hello")
			RunTestThatExpression(t, fmt.Sprintf("test_get_flag(%s)", backToOpaque)).IsEqualToBool(true)
		})

		t.Run("binary_message_conversion", func(t *testing.T) {
			// Test binary message round-trip
			obj := "test_set_name(test_set_value(test_new(), 123), 'world')"

			// Convert to binary message format
			messageExpr := fmt.Sprintf("test_to_message(%s, NULL)", obj)

			// Convert back from binary message
			backToOpaque := fmt.Sprintf("test_from_message(%s, NULL)", messageExpr)
			RunTestThatExpression(t, fmt.Sprintf("test_get_value(%s)", backToOpaque)).IsEqualToInt(123)
			RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", backToOpaque)).IsEqualToString("world")
		})

		t.Run("empty_message_conversion", func(t *testing.T) {
			// Empty message conversion tests
			emptyObj := "test_new()"

			// JSON conversion
			jsonResult := fmt.Sprintf("test_to_json(%s, NULL)", emptyObj)
			RunTestThatExpression(t, fmt.Sprintf("test_from_json(%s, NULL)", jsonResult)).IsEqualToJsonString("{}")

			// Binary message conversion
			messageResult := fmt.Sprintf("test_to_message(%s, NULL)", emptyObj)
			RunTestThatExpression(t, fmt.Sprintf("test_from_message(%s, NULL)", messageResult)).IsEqualToJsonString("{}")
		})
	})
}

// TestGeneratedOpaqueApiTypeSpecific tests type-specific edge cases and boundary conditions
func TestGeneratedOpaqueApiTypeSpecific(t *testing.T) {
	t.Run("boolean_edge_cases", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    bool flag = 1;
			|    repeated bool flags = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("boolean_basic_operations", func(t *testing.T) {
			// Test basic boolean field operations
			objTrue := "test_set_flag(test_new(), TRUE)"
			objFalse := "test_set_flag(test_new(), FALSE)"

			RunTestThatExpression(t, fmt.Sprintf("test_get_flag(%s)", objTrue)).IsEqualToBool(true)
			RunTestThatExpression(t, fmt.Sprintf("test_get_flag(%s)", objFalse)).IsEqualToBool(false)
		})

		t.Run("repeated_boolean_operations", func(t *testing.T) {
			obj := "test_add_flags(test_add_flags(test_add_flags(test_new(), TRUE), FALSE), TRUE)"
			RunTestThatExpression(t, fmt.Sprintf("test_count_flags(%s)", obj)).IsEqualToInt(3)
		})
	})

	t.Run("string_edge_cases", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    string text = 1;
			|    repeated string texts = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("empty_string", func(t *testing.T) {
			obj := "test_set_text(test_new(), '')"
			RunTestThatExpression(t, fmt.Sprintf("test_get_text(%s)", obj)).IsEqualToString("")
		})

		t.Run("special_characters", func(t *testing.T) {
			// Test strings with special characters that need JSON escaping
			specialStr := "Hello\nWorld\t\"Quote's\\Backslash\""
			obj := fmt.Sprintf("test_set_text(test_new(), ?)")
			RunTestThatExpression(t, fmt.Sprintf("test_get_text(%s)", obj), specialStr).IsEqualToString(specialStr)
		})

		t.Run("unicode_characters", func(t *testing.T) {
			// Test Unicode/UTF-8 characters
			unicodeStr := "Hello 世界 🌍 Мир"
			obj := fmt.Sprintf("test_set_text(test_new(), ?)")
			RunTestThatExpression(t, fmt.Sprintf("test_get_text(%s)", obj), unicodeStr).IsEqualToString(unicodeStr)
		})

		t.Run("long_string", func(t *testing.T) {
			// Test with a longer string
			longStr := strings.Repeat("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", 100)
			obj := fmt.Sprintf("test_set_text(test_new(), ?)")
			RunTestThatExpression(t, fmt.Sprintf("test_get_text(%s)", obj), longStr).IsEqualToString(longStr)
		})

		t.Run("repeated_string_with_special_chars", func(t *testing.T) {
			obj := "test_add_texts(test_add_texts(test_new(), 'hello\\nworld'), 'tab\\there')"
			RunTestThatExpression(t, fmt.Sprintf("test_count_texts(%s)", obj)).IsEqualToInt(2)
		})
	})

	t.Run("numeric_edge_cases", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    int32 small_int = 1;
			|    int64 big_int = 2;
			|    uint32 small_uint = 3;
			|    uint64 big_uint = 4;
			|    double precise = 5;
			|    float less_precise = 6;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("int32_boundaries", func(t *testing.T) {
			// Test int32 min/max values
			maxInt32 := "test_set_small_int(test_new(), 2147483647)"
			minInt32 := "test_set_small_int(test_new(), -2147483648)"

			RunTestThatExpression(t, fmt.Sprintf("test_get_small_int(%s)", maxInt32)).IsEqualToInt(2147483647)
			RunTestThatExpression(t, fmt.Sprintf("test_get_small_int(%s)", minInt32)).IsEqualToInt(-2147483648)
		})

		t.Run("int64_boundaries", func(t *testing.T) {
			// Test int64 large values (within MySQL's range)
			largePos := "test_set_big_int(test_new(), 9223372036854775807)"
			largeNeg := "test_set_big_int(test_new(), -9223372036854775808)"

			RunTestThatExpression(t, fmt.Sprintf("test_get_big_int(%s)", largePos)).IsEqualToInt(9223372036854775807)
			RunTestThatExpression(t, fmt.Sprintf("test_get_big_int(%s)", largeNeg)).IsEqualToInt(-9223372036854775808)
		})

		t.Run("uint32_boundaries", func(t *testing.T) {
			// Test uint32 max value
			maxUint32 := "test_set_small_uint(test_new(), 4294967295)"
			RunTestThatExpression(t, fmt.Sprintf("test_get_small_uint(%s)", maxUint32)).IsEqualToInt(4294967295)
		})

		t.Run("float_precision", func(t *testing.T) {
			// Test float precision - getters now return the actual float value
			obj := "test_set_less_precise(test_new(), 3.14159265)"
			RunTestThatExpression(t, fmt.Sprintf("test_get_less_precise(%s)", obj)).IsEqualToFloat(3.14159265)
		})

		t.Run("double_precision", func(t *testing.T) {
			// Test double precision - getters now return the actual double value
			obj := "test_set_precise(test_new(), 3.141592653589793)"
			RunTestThatExpression(t, fmt.Sprintf("test_get_precise(%s)", obj)).IsEqualToFloat(3.141592653589793)
		})

		t.Run("zero_values", func(t *testing.T) {
			// Test that zero values work correctly
			obj := "test_set_big_int(test_set_small_int(test_new(), 0), 0)"
			RunTestThatExpression(t, fmt.Sprintf("test_get_small_int(%s)", obj)).IsEqualToInt(0)
			RunTestThatExpression(t, fmt.Sprintf("test_get_big_int(%s)", obj)).IsEqualToInt(0)
		})
	})

	t.Run("bytes_field_edge_cases", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    bytes data = 1;
			|    repeated bytes chunks = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("empty_bytes", func(t *testing.T) {
			emptyData := []byte{}
			obj := fmt.Sprintf("test_set_data(test_new(), ?)")
			RunTestThatExpression(t, fmt.Sprintf("test_get_data(%s)", obj), emptyData).IsEqualToBytes(emptyData)
		})

		t.Run("binary_data", func(t *testing.T) {
			// Test with binary data including null bytes
			binaryData := []byte{0x00, 0x01, 0xFF, 0xAB, 0xCD, 0xEF}
			obj := fmt.Sprintf("test_set_data(test_new(), ?)")
			RunTestThatExpression(t, fmt.Sprintf("test_get_data(%s)", obj), binaryData).IsEqualToBytes(binaryData)
		})

		t.Run("large_bytes", func(t *testing.T) {
			// Test with larger binary data
			largeData := make([]byte, 1000)
			for i := range largeData {
				largeData[i] = byte(i % 256)
			}
			obj := fmt.Sprintf("test_set_data(test_new(), ?)")
			RunTestThatExpression(t, fmt.Sprintf("test_get_data(%s)", obj), largeData).IsEqualToBytes(largeData)
		})

		t.Run("repeated_bytes", func(t *testing.T) {
			data1 := []byte("chunk1")
			data2 := []byte("chunk2")
			obj := "test_add_chunks(test_add_chunks(test_new(), ?), ?)"
			RunTestThatExpression(t, fmt.Sprintf("test_count_chunks(%s)", obj), data1, data2).IsEqualToInt(2)
		})
	})

	t.Run("field_number_edge_cases", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    int32 field_1 = 1;
			|    int32 field_max = 536870911;  // 2^29 - 1 (max field number)
			|    int32 field_large = 999999;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("large_field_numbers", func(t *testing.T) {
			obj := "test_set_field_large(test_set_field_1(test_new(), 42), 123)"

			// Verify both fields are set correctly
			RunTestThatExpression(t, fmt.Sprintf("test_get_field_1(%s)", obj)).IsEqualToInt(42)
			RunTestThatExpression(t, fmt.Sprintf("test_get_field_large(%s)", obj)).IsEqualToInt(123)
		})

		t.Run("maximum_field_number", func(t *testing.T) {
			obj := "test_set_field_max(test_new(), 456)"
			RunTestThatExpression(t, fmt.Sprintf("test_get_field_max(%s)", obj)).IsEqualToInt(456)
		})
	})

	t.Run("json_escaping_and_encoding", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    string json_like = 1;
			|    string control_chars = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("json_like_strings", func(t *testing.T) {
			// Test strings that look like JSON
			jsonStr := `{"key": "value", "number": 42}`
			obj := fmt.Sprintf("test_set_json_like(test_new(), ?)")
			RunTestThatExpression(t, fmt.Sprintf("test_get_json_like(%s)", obj), jsonStr).IsEqualToString(jsonStr)
		})

		t.Run("control_characters", func(t *testing.T) {
			// Test control characters that need special JSON handling
			controlStr := "\n\r\t\b\f\"\\"
			obj := fmt.Sprintf("test_set_control_chars(test_new(), ?)")
			RunTestThatExpression(t, fmt.Sprintf("test_get_control_chars(%s)", obj), controlStr).IsEqualToString(controlStr)
		})

		t.Run("null_character", func(t *testing.T) {
			// Test null character handling
			nullStr := "before\x00after"
			obj := fmt.Sprintf("test_set_control_chars(test_new(), ?)")
			RunTestThatExpression(t, fmt.Sprintf("test_get_control_chars(%s)", obj), nullStr).IsEqualToString(nullStr)
		})
	})

	t.Run("complex_nested_structures", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    repeated ComplexMessage items = 1;
			|}
			|message ComplexMessage {
			|    string name = 1;
			|    repeated int32 numbers = 2;
			|    bool flag = 3;
			|    NestedMessage nested = 4;
			|}
			|message NestedMessage {
			|    double value = 1;
			|    bytes data = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("deeply_nested_structure", func(t *testing.T) {
			// Create a complex nested structure using hardcoded values for simplicity
			nestedObj := "JSON_OBJECT('1', 3.14, '2', 'nested_data')"
			complexObj := fmt.Sprintf("JSON_OBJECT('1', 'item1', '2', JSON_ARRAY(1, 2, 3), '3', TRUE, '4', %s)", nestedObj)
			obj := fmt.Sprintf("test_add_items(test_new(), %s)", complexObj)

			// Verify the structure was created correctly
			RunTestThatExpression(t, fmt.Sprintf("test_count_items(%s)", obj)).IsEqualToInt(1)

			// Test that the nested structure is accessible (simplified test)
			itemExpr := fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(test_get_all_items(%s), '$[0].\"1\"'))", obj)
			RunTestThatExpression(t, itemExpr).IsEqualToString("item1")
		})
	})
}
