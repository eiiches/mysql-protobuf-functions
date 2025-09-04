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
		t.Logf("Executing SQL statement %d: %s", i, truncateForLog(stmt, 100))
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
			RunTestThatExpression(t, fmt.Sprintf("%s(test_new(), TRUE)", setterFunc)).IsEqualToJsonString(fmt.Sprintf(`{"1": true}`))
			RunTestThatExpression(t, fmt.Sprintf("%s(test_new(), FALSE)", setterFunc)).IsEqualToJsonString(fmt.Sprintf(`{"1": false}`))
		})
	}
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
		// Note: bytes fields in MySQL are handled as binary data, we'll test this differently
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    bytes data = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("constructor", func(t *testing.T) {
			RunTestThatExpression(t, "test_new()").IsEqualToJsonString("{}")
		})

		t.Run("setter_and_getter", func(t *testing.T) {
			// Test with binary data (represented as hex)
			testData := []byte("hello")
			RunTestThatExpression(t, "test_get_data(test_set_data(test_new(), ?))", testData).IsEqualToBytes(testData)
		})
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
			RunTestThatExpression(t, "test_get_value(test_set_value(test_new(), 3.141592653589793))").IsEqualToDouble(3.141592653589793)
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
			RunTestThatExpression(t, "test_get_items(test_new())").IsEqualToJsonString("[]")
		})

		t.Run("add_elements", func(t *testing.T) {
			// Add single element
			RunTestThatExpression(t, "test_count_items(test_add_items(test_new(), 42))").IsEqualToInt(1)
			RunTestThatExpression(t, "test_get_items(test_add_items(test_new(), 42))").IsEqualToJsonString("[42]")

			// Add multiple elements
			obj := "test_add_items(test_add_items(test_add_items(test_new(), 1), 2), 3)"
			RunTestThatExpression(t, fmt.Sprintf("test_count_items(%s)", obj)).IsEqualToInt(3)
			RunTestThatExpression(t, fmt.Sprintf("test_get_items(%s)", obj)).IsEqualToJsonString("[1, 2, 3]")
		})

		t.Run("set_entire_array", func(t *testing.T) {
			RunTestThatExpression(t, "test_get_items(test_set_items(test_new(), JSON_ARRAY(10, 20, 30)))").IsEqualToJsonString("[10, 20, 30]")
		})

		t.Run("clear_array", func(t *testing.T) {
			obj := "test_add_items(test_add_items(test_new(), 1), 2)"
			RunTestThatExpression(t, fmt.Sprintf("test_count_items(test_clear_items(%s))", obj)).IsEqualToInt(0)
			RunTestThatExpression(t, fmt.Sprintf("test_get_items(test_clear_items(%s))", obj)).IsEqualToJsonString("[]")
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
			RunTestThatExpression(t, fmt.Sprintf("test_get_names(%s)", obj)).IsEqualToJsonString(`["alice", "bob"]`)
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
			RunTestThatExpression(t, fmt.Sprintf("test_get_flags(%s)", obj)).IsEqualToJsonString(`[true, false]`)
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
			RunTestThatExpression(t, fmt.Sprintf("test_get_items(%s)", expr)).IsEqualToJsonString(`[{"1": 42}, {"1": 100}]`)
		})
	})
}

// TestGeneratedOpaqueApiOneOfFields tests oneOf field mutual exclusion
func TestGeneratedOpaqueApiOneOfFields(t *testing.T) {
	t.Run("basic_oneof", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    oneof choice {
			|        int32 number = 1;
			|        string name = 2;
			|        bool flag = 3;
			|    }
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("initially_empty", func(t *testing.T) {
			RunTestThatExpression(t, "test_which_choice(test_new())").IsNull()
		})

		t.Run("set_number_field", func(t *testing.T) {
			obj := "test_set_number(test_new(), 42)"
			RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", obj)).IsEqualToString("number")
			RunTestThatExpression(t, fmt.Sprintf("test_get_number(%s)", obj)).IsEqualToInt(42)

			// Other fields should return defaults
			RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", obj)).IsEqualToString("")
			RunTestThatExpression(t, fmt.Sprintf("test_get_flag(%s)", obj)).IsEqualToBool(false)
		})

		t.Run("set_name_field", func(t *testing.T) {
			obj := "test_set_name(test_new(), 'hello')"
			RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", obj)).IsEqualToString("name")
			RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", obj)).IsEqualToString("hello")

			// Other fields should return defaults
			RunTestThatExpression(t, fmt.Sprintf("test_get_number(%s)", obj)).IsEqualToInt(0)
			RunTestThatExpression(t, fmt.Sprintf("test_get_flag(%s)", obj)).IsEqualToBool(false)
		})

		t.Run("mutual_exclusion", func(t *testing.T) {
			// Set number first, then name - number should be cleared
			obj := "test_set_name(test_set_number(test_new(), 42), 'world')"
			RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", obj)).IsEqualToString("name")
			RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", obj)).IsEqualToString("world")
			RunTestThatExpression(t, fmt.Sprintf("test_get_number(%s)", obj)).IsEqualToInt(0) // cleared
		})

		t.Run("clear_oneof_group", func(t *testing.T) {
			obj := "test_clear_choice(test_set_name(test_new(), 'test'))"
			RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", obj)).IsNull()
			RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", obj)).IsEqualToString("")
		})
	})

	t.Run("oneof_with_message", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    oneof data {
			|        int32 number = 1;
			|        MessageType nested = 2;
			|    }
			|}
			|message MessageType {
			|    string value = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("set_message_field", func(t *testing.T) {
			nestedObj := "JSON_OBJECT('1', 'nested_value')"
			obj := fmt.Sprintf("test_set_nested(%s, %s)", "test_new()", nestedObj)
			RunTestThatExpression(t, fmt.Sprintf("test_which_data(%s)", obj)).IsEqualToString("nested")
			RunTestThatExpression(t, fmt.Sprintf("test_get_nested(%s)", obj)).IsEqualToJsonString(`{"1": "nested_value"}`)
		})

		t.Run("message_to_scalar_exclusion", func(t *testing.T) {
			// Set message, then scalar - message should be cleared
			nestedObj := "JSON_OBJECT('1', 'test')"
			obj := fmt.Sprintf("test_set_number(test_set_nested(test_new(), %s), 123)", nestedObj)
			RunTestThatExpression(t, fmt.Sprintf("test_which_data(%s)", obj)).IsEqualToString("number")
			RunTestThatExpression(t, fmt.Sprintf("test_get_number(%s)", obj)).IsEqualToInt(123)
			RunTestThatExpression(t, fmt.Sprintf("test_get_nested(%s)", obj)).IsNull() // cleared
		})
	})

	t.Run("multiple_oneofs", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    oneof first_choice {
			|        int32 first_number = 1;
			|        string first_name = 2;
			|    }
			|    oneof second_choice {
			|        bool second_flag = 3;
			|        int32 second_number = 4;
			|    }
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("independent_oneofs", func(t *testing.T) {
			// Set fields in different oneOf groups - they should not interfere
			obj := "test_set_second_flag(test_set_first_name(test_new(), 'hello'), TRUE)"

			RunTestThatExpression(t, fmt.Sprintf("test_which_first_choice(%s)", obj)).IsEqualToString("first_name")
			RunTestThatExpression(t, fmt.Sprintf("test_which_second_choice(%s)", obj)).IsEqualToString("second_flag")

			RunTestThatExpression(t, fmt.Sprintf("test_get_first_name(%s)", obj)).IsEqualToString("hello")
			RunTestThatExpression(t, fmt.Sprintf("test_get_second_flag(%s)", obj)).IsEqualToBool(true)
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

// TestGeneratedOpaqueApiPresence tests field presence semantics (has_* methods)
func TestGeneratedOpaqueApiPresence(t *testing.T) {
	t.Run("optional_field_presence", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    optional int32 value = 1;
			|    optional string name = 2;
			|    optional bool flag = 3;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("initially_not_present", func(t *testing.T) {
			RunTestThatExpression(t, "test_has_value(test_new())").IsFalse()
			RunTestThatExpression(t, "test_has_name(test_new())").IsFalse()
			RunTestThatExpression(t, "test_has_flag(test_new())").IsFalse()
		})

		t.Run("present_after_setting", func(t *testing.T) {
			obj := "test_set_value(test_new(), 42)"
			RunTestThatExpression(t, fmt.Sprintf("test_has_value(%s)", obj)).IsTrue()
			RunTestThatExpression(t, fmt.Sprintf("test_get_value(%s)", obj)).IsEqualToInt(42)

			// Other fields should still not be present
			RunTestThatExpression(t, fmt.Sprintf("test_has_name(%s)", obj)).IsFalse()
			RunTestThatExpression(t, fmt.Sprintf("test_has_flag(%s)", obj)).IsFalse()
		})

		t.Run("setting_default_values_shows_presence", func(t *testing.T) {
			// Setting default values should still show as present
			obj1 := "test_set_value(test_new(), 0)"
			RunTestThatExpression(t, fmt.Sprintf("test_has_value(%s)", obj1)).IsTrue()
			RunTestThatExpression(t, fmt.Sprintf("test_get_value(%s)", obj1)).IsEqualToInt(0)

			obj2 := "test_set_name(test_new(), '')"
			RunTestThatExpression(t, fmt.Sprintf("test_has_name(%s)", obj2)).IsTrue()
			RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", obj2)).IsEqualToString("")

			obj3 := "test_set_flag(test_new(), FALSE)"
			RunTestThatExpression(t, fmt.Sprintf("test_has_flag(%s)", obj3)).IsTrue()
			RunTestThatExpression(t, fmt.Sprintf("test_get_flag(%s)", obj3)).IsEqualToBool(false)
		})

		t.Run("absent_after_clearing", func(t *testing.T) {
			obj := "test_clear_value(test_set_value(test_new(), 42))"
			RunTestThatExpression(t, fmt.Sprintf("test_has_value(%s)", obj)).IsFalse()
			RunTestThatExpression(t, fmt.Sprintf("test_get_value(%s)", obj)).IsEqualToInt(0) // returns default
		})

		t.Run("multiple_optional_fields", func(t *testing.T) {
			obj := "test_set_name(test_set_value(test_new(), 123), 'hello')"
			RunTestThatExpression(t, fmt.Sprintf("test_has_value(%s)", obj)).IsTrue()
			RunTestThatExpression(t, fmt.Sprintf("test_has_name(%s)", obj)).IsTrue()
			RunTestThatExpression(t, fmt.Sprintf("test_has_flag(%s)", obj)).IsFalse()

			RunTestThatExpression(t, fmt.Sprintf("test_get_value(%s)", obj)).IsEqualToInt(123)
			RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", obj)).IsEqualToString("hello")
			RunTestThatExpression(t, fmt.Sprintf("test_get_flag(%s)", obj)).IsEqualToBool(false) // default
		})
	})

	t.Run("optional_message_field_presence", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    optional MessageType nested = 1;
			|}
			|message MessageType {
			|    int32 value = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("initially_not_present", func(t *testing.T) {
			RunTestThatExpression(t, "test_has_nested(test_new())").IsFalse()
			RunTestThatExpression(t, "test_get_nested(test_new())").IsNull() // returns NULL for absent message
		})

		t.Run("present_after_setting", func(t *testing.T) {
			nestedObj := "JSON_OBJECT('1', 42)"
			obj := fmt.Sprintf("test_set_nested(test_new(), %s)", nestedObj)
			RunTestThatExpression(t, fmt.Sprintf("test_has_nested(%s)", obj)).IsTrue()
			RunTestThatExpression(t, fmt.Sprintf("test_get_nested(%s)", obj)).IsEqualToJsonString(`{"1": 42}`)
		})

		t.Run("absent_after_clearing", func(t *testing.T) {
			nestedObj := "JSON_OBJECT('1', 123)"
			obj := fmt.Sprintf("test_clear_nested(test_set_nested(test_new(), %s))", nestedObj)
			RunTestThatExpression(t, fmt.Sprintf("test_has_nested(%s)", obj)).IsFalse()
			RunTestThatExpression(t, fmt.Sprintf("test_get_nested(%s)", obj)).IsNull()
		})
	})

	t.Run("oneof_field_presence", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    oneof choice {
			|        int32 number = 1;
			|        string name = 2;
			|    }
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("initially_no_presence", func(t *testing.T) {
			// OneOf fields might have individual has_* methods or use which_* for presence
			RunTestThatExpression(t, "test_which_choice(test_new())").IsNull()
		})

		t.Run("presence_after_setting", func(t *testing.T) {
			obj := "test_set_number(test_new(), 42)"
			RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", obj)).IsEqualToString("number")

			// Check if individual has_* methods exist for oneOf fields
			t.Run("individual_has_methods", func(t *testing.T) {
				// These may or may not be generated - test if they exist
				RunTestThatExpression(t, fmt.Sprintf("test_get_number(%s)", obj)).IsEqualToInt(42)
				RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", obj)).IsEqualToString("") // default for unset field
			})
		})
	})

	t.Run("repeated_field_presence", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    repeated int32 items = 1;
			|    repeated string names = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("empty_array_presence", func(t *testing.T) {
			// Repeated fields are always "present" but may be empty
			RunTestThatExpression(t, "test_count_items(test_new())").IsEqualToInt(0)
			RunTestThatExpression(t, "test_get_items(test_new())").IsEqualToJsonString("[]")
		})

		t.Run("non_empty_array_presence", func(t *testing.T) {
			obj := "test_add_items(test_new(), 42)"
			RunTestThatExpression(t, fmt.Sprintf("test_count_items(%s)", obj)).IsEqualToInt(1)
			RunTestThatExpression(t, fmt.Sprintf("test_get_items(%s)", obj)).IsEqualToJsonString("[42]")
		})

		t.Run("cleared_array_presence", func(t *testing.T) {
			obj := "test_clear_items(test_add_items(test_new(), 42))"
			RunTestThatExpression(t, fmt.Sprintf("test_count_items(%s)", obj)).IsEqualToInt(0)
			RunTestThatExpression(t, fmt.Sprintf("test_get_items(%s)", obj)).IsEqualToJsonString("[]")
		})
	})
}

// TestGeneratedOpaqueApiConversions tests JSON/protobuf round-trip conversions
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

		t.Run("json_to_protobuf_to_json", func(t *testing.T) {
			// Create a message using the opaque API
			obj := "test_set_flag(test_set_name(test_set_value(test_new(), 42), 'hello'), TRUE)"

			// Convert to protobuf wire format (binary)
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString(`{"1": 42, "2": "hello", "3": true}`)

			// Test the conversion produces a valid protobuf-compatible JSON representation
			// Round-trip: opaque -> json -> opaque
			jsonExpr := fmt.Sprintf("test_to_protobuf(%s)", obj)
			backToOpaque := fmt.Sprintf("test_from_protobuf(%s)", jsonExpr)

			RunTestThatExpression(t, fmt.Sprintf("test_get_value(%s)", backToOpaque)).IsEqualToInt(42)
			RunTestThatExpression(t, fmt.Sprintf("test_get_name(%s)", backToOpaque)).IsEqualToString("hello")
			RunTestThatExpression(t, fmt.Sprintf("test_get_flag(%s)", backToOpaque)).IsEqualToBool(true)
		})

		t.Run("protobuf_json_format_compatibility", func(t *testing.T) {
			// Test that our ProtoNumberJSON format is properly structured
			obj := "test_set_name(test_set_value(test_new(), 123), 'world')"

			// Verify the JSON structure uses field numbers as keys
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_to_protobuf(%s), '$.\"1\"')", obj)).IsEqualToInt(123)
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_to_protobuf(%s), '$.\"2\"')", obj)).IsEqualToString("world")
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_to_protobuf(%s), '$.\"3\"')", obj)).IsNull() // not set
		})

		t.Run("empty_message_conversion", func(t *testing.T) {
			// Empty message should convert to empty JSON object
			RunTestThatExpression(t, "test_to_protobuf(test_new())").IsEqualToJsonString("{}")

			// Empty JSON should convert back to empty message
			RunTestThatExpression(t, "test_from_protobuf('{}')").IsEqualToJsonString("{}")
		})
	})

	t.Run("nested_message_roundtrip", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    MessageType nested = 1;
			|    int32 value = 2;
			|}
			|message MessageType {
			|    string data = 1;
			|    int32 number = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("nested_object_conversion", func(t *testing.T) {
			// Create nested message structure
			nestedObj := "JSON_OBJECT('1', 'nested_data', '2', 456)"
			obj := fmt.Sprintf("test_set_value(test_set_nested(test_new(), %s), 789)", nestedObj)

			// Convert to ProtoNumberJSON format
			expectedJson := `{"1": {"1": "nested_data", "2": 456}, "2": 789}`
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString(expectedJson)

			// Round-trip conversion
			backToOpaque := fmt.Sprintf("test_from_protobuf('%s')", expectedJson)
			RunTestThatExpression(t, fmt.Sprintf("test_get_value(%s)", backToOpaque)).IsEqualToInt(789)
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_get_nested(%s), '$.\"1\"')", backToOpaque)).IsEqualToString("nested_data")
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_get_nested(%s), '$.\"2\"')", backToOpaque)).IsEqualToInt(456)
		})
	})

	t.Run("repeated_field_conversion", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    repeated int32 numbers = 1;
			|    repeated string names = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("repeated_field_roundtrip", func(t *testing.T) {
			// Create message with repeated fields
			obj := "test_add_names(test_add_names(test_add_numbers(test_add_numbers(test_new(), 1), 2), 'alice'), 'bob')"

			// Convert to ProtoNumberJSON format
			expectedJson := `{"1": [1, 2], "2": ["alice", "bob"]}`
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString(expectedJson)

			// Round-trip conversion
			backToOpaque := fmt.Sprintf("test_from_protobuf('%s')", expectedJson)
			RunTestThatExpression(t, fmt.Sprintf("test_get_numbers(%s)", backToOpaque)).IsEqualToJsonString("[1, 2]")
			RunTestThatExpression(t, fmt.Sprintf("test_get_names(%s)", backToOpaque)).IsEqualToJsonString(`["alice", "bob"]`)
			RunTestThatExpression(t, fmt.Sprintf("test_count_numbers(%s)", backToOpaque)).IsEqualToInt(2)
			RunTestThatExpression(t, fmt.Sprintf("test_count_names(%s)", backToOpaque)).IsEqualToInt(2)
		})

		t.Run("empty_repeated_fields", func(t *testing.T) {
			// Empty repeated fields should be represented as empty arrays
			RunTestThatExpression(t, "test_to_protobuf(test_new())").IsEqualToJsonString("{}")

			// Arrays with elements, then cleared
			obj := "test_clear_numbers(test_add_numbers(test_new(), 42))"
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString("{}")
		})
	})

	t.Run("oneof_field_conversion", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    oneof choice {
			|        int32 number = 1;
			|        string text = 2;
			|        MessageType nested = 3;
			|    }
			|}
			|message MessageType {
			|    bool flag = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("oneof_number_conversion", func(t *testing.T) {
			obj := "test_set_number(test_new(), 42)"
			expectedJson := `{"1": 42}`
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString(expectedJson)

			// Round-trip
			backToOpaque := fmt.Sprintf("test_from_protobuf('%s')", expectedJson)
			RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", backToOpaque)).IsEqualToString("number")
			RunTestThatExpression(t, fmt.Sprintf("test_get_number(%s)", backToOpaque)).IsEqualToInt(42)
		})

		t.Run("oneof_nested_message_conversion", func(t *testing.T) {
			nestedObj := "JSON_OBJECT('1', TRUE)"
			obj := fmt.Sprintf("test_set_nested(test_new(), %s)", nestedObj)
			expectedJson := `{"3": {"1": true}}`
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString(expectedJson)

			// Round-trip
			backToOpaque := fmt.Sprintf("test_from_protobuf('%s')", expectedJson)
			RunTestThatExpression(t, fmt.Sprintf("test_which_choice(%s)", backToOpaque)).IsEqualToString("nested")
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_get_nested(%s), '$.\"1\"')", backToOpaque)).IsEqualToBool(true)
		})
	})

	t.Run("enum_field_conversion", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    Status status = 1;
			|}
			|enum Status {
			|    STATUS_UNKNOWN = 0;
			|    STATUS_ACTIVE = 1;
			|    STATUS_INACTIVE = 2;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("enum_numeric_conversion", func(t *testing.T) {
			obj := "test_set_status(test_new(), 2)"
			expectedJson := `{"1": 2}`
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString(expectedJson)

			// Round-trip
			backToOpaque := fmt.Sprintf("test_from_protobuf('%s')", expectedJson)
			RunTestThatExpression(t, fmt.Sprintf("test_get_status(%s)", backToOpaque)).IsEqualToInt(2)
			RunTestThatExpression(t, fmt.Sprintf("status_to_string(test_get_status(%s))", backToOpaque)).IsEqualToString("STATUS_INACTIVE")
		})

		t.Run("enum_zero_value_conversion", func(t *testing.T) {
			// Enum zero values should be included in JSON when explicitly set
			obj := "test_set_status(test_new(), 0)"
			expectedJson := `{"1": 0}`
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString(expectedJson)
		})
	})

	t.Run("type_specific_conversions", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    bool flag = 1;
			|    double pi = 2;
			|    float small_pi = 3;
			|    bytes data = 4;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("boolean_conversion", func(t *testing.T) {
			objTrue := "test_set_flag(test_new(), TRUE)"
			objFalse := "test_set_flag(test_new(), FALSE)"

			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", objTrue)).IsEqualToJsonString(`{"1": true}`)
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", objFalse)).IsEqualToJsonString(`{"1": false}`)

			// Round-trip
			RunTestThatExpression(t, "test_get_flag(test_from_protobuf('{\"1\": true}'))").IsEqualToBool(true)
			RunTestThatExpression(t, "test_get_flag(test_from_protobuf('{\"1\": false}'))").IsEqualToBool(false)
		})

		t.Run("float_conversion", func(t *testing.T) {
			obj := "test_set_small_pi(test_set_pi(test_new(), 3.141592653589793), 3.14)"

			// Verify float precision is maintained in JSON conversion
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_to_protobuf(%s), '$.\"2\"')", obj)).IsEqualToDouble(3.141592653589793)
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_to_protobuf(%s), '$.\"3\"')", obj)).IsEqualToFloat(3.14)
		})

		t.Run("bytes_conversion", func(t *testing.T) {
			// Test binary data conversion
			testData := []byte("hello world")
			obj := fmt.Sprintf("test_set_data(test_new(), ?)")

			// Convert bytes field to JSON and back
			jsonResult := fmt.Sprintf("test_to_protobuf(%s)", obj)
			backToOpaque := fmt.Sprintf("test_from_protobuf(%s)", jsonResult)

			RunTestThatExpression(t, fmt.Sprintf("test_get_data(%s)", backToOpaque), testData).IsEqualToBytes(testData)
		})
	})

	t.Run("conversion_error_handling", func(t *testing.T) {
		protoContent := dedent.Pipe(`
			|syntax = "proto3";
			|message Test {
			|    int32 value = 1;
			|}
		`)
		schemaName := "test_schema"
		generateAndLoadOpaqueApiSQL(t, protoContent, schemaName)

		t.Run("invalid_json_input", func(t *testing.T) {
			// Test error handling for malformed JSON
			RunTestThatExpression(t, "test_from_protobuf('{invalid json}')").ToFailWithSignalException("45000", "Invalid JSON")
		})

		t.Run("wrong_field_types", func(t *testing.T) {
			// Test error handling for wrong field types in JSON
			RunTestThatExpression(t, "test_from_protobuf('{\"1\": \"not_a_number\"}')").ToFailWithSignalException("45000", "Type mismatch")
		})

		t.Run("unknown_field_numbers", func(t *testing.T) {
			// Unknown field numbers should be ignored or handled gracefully
			validObj := "test_from_protobuf('{\"1\": 42, \"999\": \"unknown_field\"}')"
			RunTestThatExpression(t, fmt.Sprintf("test_get_value(%s)", validObj)).IsEqualToInt(42)
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

		t.Run("boolean_json_representation", func(t *testing.T) {
			// Ensure booleans are stored as true/false in JSON, not 1/0
			objTrue := "test_set_flag(test_new(), TRUE)"
			objFalse := "test_set_flag(test_new(), FALSE)"

			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", objTrue)).IsEqualToJsonString(`{"1": true}`)
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", objFalse)).IsEqualToJsonString(`{"1": false}`)
		})

		t.Run("repeated_boolean_json", func(t *testing.T) {
			obj := "test_add_flags(test_add_flags(test_add_flags(test_new(), TRUE), FALSE), TRUE)"
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString(`{"2": [true, false, true]}`)
		})

		t.Run("boolean_mysql_vs_json_conversion", func(t *testing.T) {
			// Test that MySQL TRUE/FALSE converts properly to JSON true/false
			RunTestThatExpression(t, "test_get_flag(test_from_protobuf('{\"1\": true}'))").IsEqualToBool(true)
			RunTestThatExpression(t, "test_get_flag(test_from_protobuf('{\"1\": false}'))").IsEqualToBool(false)
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
			RunTestThatExpression(t, fmt.Sprintf("test_to_protobuf(%s)", obj)).IsEqualToJsonString(`{"1": ""}`)
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
			// Test float precision
			obj := "test_set_less_precise(test_new(), 3.14159265)"
			RunTestThatExpression(t, fmt.Sprintf("test_get_less_precise(%s)", obj)).IsEqualToFloat(3.14159265)
		})

		t.Run("double_precision", func(t *testing.T) {
			// Test double precision
			obj := "test_set_precise(test_new(), 3.141592653589793)"
			RunTestThatExpression(t, fmt.Sprintf("test_get_precise(%s)", obj)).IsEqualToDouble(3.141592653589793)
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

			// Verify JSON uses correct field numbers
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_to_protobuf(%s), '$.\"1\"')", obj)).IsEqualToInt(42)
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_to_protobuf(%s), '$.\"999999\"')", obj)).IsEqualToInt(123)
		})

		t.Run("maximum_field_number", func(t *testing.T) {
			obj := "test_set_field_max(test_new(), 456)"
			RunTestThatExpression(t, fmt.Sprintf("test_get_field_max(%s)", obj)).IsEqualToInt(456)
			RunTestThatExpression(t, fmt.Sprintf("JSON_EXTRACT(test_to_protobuf(%s), '$.\"536870911\"')", obj)).IsEqualToInt(456)
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
			itemExpr := fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(test_get_items(%s), '$[0].\"1\"'))", obj)
			RunTestThatExpression(t, itemExpr).IsEqualToString("item1")
		})
	})
}
