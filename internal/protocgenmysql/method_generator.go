package protocgenmysql

import (
	"fmt"
	"slices"
	"strings"

	"github.com/eiiches/mysql-protobuf-functions/internal/protoreflectutils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// FileNameFunc defines how to transform a proto file path to a SQL file path
type FileNameFunc func(protoPath string) string

// TypePrefixFunc defines how to transform a proto package and type name to a SQL function prefix
type TypePrefixFunc func(packageName protoreflect.FullName, typeName protoreflect.FullName) string

// GenerateMethodFragments returns individual content fragments for each proto file
func GenerateMethodFragments(files *protoregistry.Files, fileNameFunc FileNameFunc, typePrefixFunc TypePrefixFunc, schemaFunctionName string) (map[string][]string, error) {
	fileFragments := make(map[string][]string)

	// Collect all files and sort by path for deterministic ordering
	var allFiles []protoreflect.FileDescriptor
	files.RangeFiles(func(fileDesc protoreflect.FileDescriptor) bool {
		allFiles = append(allFiles, fileDesc)
		return true // continue iteration
	})

	// Sort files by path to ensure deterministic ordering
	slices.SortFunc(allFiles, func(a, b protoreflect.FileDescriptor) int {
		return strings.Compare(a.Path(), b.Path())
	})

	// Generate fragments for each proto file in sorted order
	for _, fileDesc := range allFiles {
		filename := fileNameFunc(fileDesc.Path())
		content, err := generateMethodsForFile(fileDesc, typePrefixFunc, schemaFunctionName)
		if err != nil {
			return nil, err
		}
		if content != "" {
			fileFragments[filename] = append(fileFragments[filename], content)
		}
	}

	return fileFragments, nil
}

func generateMethodsForFile(fileDesc protoreflect.FileDescriptor, typePrefixFunc TypePrefixFunc, schemaFunctionName string) (string, error) {
	var content strings.Builder

	// Generate methods for each message type
	messages := fileDesc.Messages()
	for messageDesc := range protoreflectutils.Iterate(messages) {
		if err := generateMessageMethods(&content, messageDesc, typePrefixFunc, schemaFunctionName); err != nil {
			return "", err
		}
	}

	// Generate methods for each enum type
	enums := fileDesc.Enums()
	for enumDesc := range protoreflectutils.Iterate(enums) {
		if err := generateEnumMethods(&content, enumDesc, typePrefixFunc, schemaFunctionName); err != nil {
			return "", err
		}
	}

	return content.String(), nil
}

func generateMessageMethods(content *strings.Builder, messageDesc protoreflect.MessageDescriptor, typePrefixFunc TypePrefixFunc, schemaFunctionName string) error {
	// Use FullName from descriptor - no manual string construction needed
	fullTypeName := messageDesc.FullName()
	packageName := messageDesc.ParentFile().Package()

	// Get prefix for this specific type
	funcPrefix := typePrefixFunc(packageName, fullTypeName)

	// Generate basic constructor
	if err := generateConstructor(content, funcPrefix, fullTypeName); err != nil {
		return err
	}

	// Generate conversion methods (from_json, to_json, etc.)
	if err := generateConversionMethods(content, funcPrefix, fullTypeName, schemaFunctionName); err != nil {
		return err
	}

	// Generate field accessor methods with enhanced opaque API patterns
	if err := generateFieldAccessorMethods(content, messageDesc, funcPrefix, fullTypeName); err != nil {
		return err
	}

	// Generate oneOf methods for oneOf groups
	if err := generateOneOfMethods(content, messageDesc, funcPrefix, fullTypeName); err != nil {
		return err
	}

	// Generate methods for nested message types
	nestedMessages := messageDesc.Messages()
	for nestedMessageDesc := range protoreflectutils.Iterate(nestedMessages) {
		if err := generateMessageMethods(content, nestedMessageDesc, typePrefixFunc, schemaFunctionName); err != nil {
			return err
		}
	}

	// Generate methods for nested enum types
	nestedEnums := messageDesc.Enums()
	for nestedEnumDesc := range protoreflectutils.Iterate(nestedEnums) {
		if err := generateEnumMethods(content, nestedEnumDesc, typePrefixFunc, schemaFunctionName); err != nil {
			return err
		}
	}

	return nil
}

// generateConstructor creates a basic constructor function
func generateConstructor(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName) error {
	newFuncName := funcPrefix + "_new"
	if err := validateFunctionName(newFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", newFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s() RETURNS JSON DETERMINISTIC\n", newFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    RETURN JSON_OBJECT();\n")
	content.WriteString("END $$\n\n")

	return nil
}

// generateConversionMethods creates conversion methods (from_json, to_json, etc.)
func generateConversionMethods(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, schemaFunctionName string) error {
	// Generate from_json
	fromJsonFuncName := funcPrefix + "_from_json"
	if err := validateFunctionName(fromJsonFuncName, fullTypeName); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", fromJsonFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(json_data JSON, json_unmarshal_options JSON) RETURNS JSON DETERMINISTIC\n", fromJsonFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_json_to_number_json(%s(), '.%s', json_data, json_unmarshal_options);\n", schemaFunctionName, fullTypeName))
	content.WriteString("END $$\n\n")

	// Generate from_message
	fromMessageFuncName := funcPrefix + "_from_message"
	if err := validateFunctionName(fromMessageFuncName, fullTypeName); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", fromMessageFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(message_data LONGBLOB, unmarshal_options JSON) RETURNS JSON DETERMINISTIC\n", fromMessageFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_message_to_number_json(%s(), '.%s', message_data, unmarshal_options);\n", schemaFunctionName, fullTypeName))
	content.WriteString("END $$\n\n")

	// Generate to_json
	toJsonFuncName := funcPrefix + "_to_json"
	if err := validateFunctionName(toJsonFuncName, fullTypeName); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", toJsonFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, json_marshal_options JSON) RETURNS JSON DETERMINISTIC\n", toJsonFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_number_json_to_json(%s(), '.%s', proto_data, json_marshal_options);\n", schemaFunctionName, fullTypeName))
	content.WriteString("END $$\n\n")

	// Generate to_message
	toMessageFuncName := funcPrefix + "_to_message"
	if err := validateFunctionName(toMessageFuncName, fullTypeName); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", toMessageFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, marshal_options JSON) RETURNS LONGBLOB DETERMINISTIC\n", toMessageFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_number_json_to_message(%s(), '.%s', proto_data, marshal_options);\n", schemaFunctionName, fullTypeName))
	content.WriteString("END $$\n\n")

	return nil
}

// generateFieldAccessorMethods creates enhanced field accessor methods following opaque API patterns
func generateFieldAccessorMethods(content *strings.Builder, messageDesc protoreflect.MessageDescriptor, funcPrefix string, fullTypeName protoreflect.FullName) error {
	fields := messageDesc.Fields()

	for field := range protoreflectutils.Iterate(fields) {
		fieldName := string(field.Name())
		fieldType := field.Kind()
		isRepeated := field.Cardinality() == protoreflect.Repeated
		isMessage := fieldType == protoreflect.MessageKind

		// Generate enhanced getter with better defaults and type safety
		if err := generateEnhancedGetter(content, funcPrefix, fullTypeName, field, fieldName, fieldType, isRepeated, isMessage); err != nil {
			return err
		}

		// Generate enhanced setter with validation
		if err := generateEnhancedSetter(content, funcPrefix, fullTypeName, field, fieldName, fieldType, isRepeated, isMessage); err != nil {
			return err
		}

		// Generate has method for fields with presence semantics
		if field.HasPresence() {
			if err := generateHasMethod(content, funcPrefix, fullTypeName, field, fieldName); err != nil {
				return err
			}
		}

		// Generate clear method for all fields
		if err := generateClearMethod(content, funcPrefix, fullTypeName, field, fieldName); err != nil {
			return err
		}

		// Generate additional methods for repeated fields
		if isRepeated {
			if err := generateRepeatedFieldMethods(content, funcPrefix, fullTypeName, field, fieldName, fieldType); err != nil {
				return err
			}
		}
	}

	return nil
}

// generateEnhancedGetter creates an enhanced getter with better defaults and type safety
func generateEnhancedGetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string, fieldType protoreflect.Kind, isRepeated, isMessage bool) error {
	getterFuncName := fmt.Sprintf("%s_get_%s", funcPrefix, fieldName)
	if err := validateFunctionName(getterFuncName, fullTypeName); err != nil {
		return err
	}

	_, getterType := getMySQLTypesForFieldFromReflection(fieldType, isRepeated)

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS %s DETERMINISTIC\n", getterFuncName, getterType))
	content.WriteString("BEGIN\n")

	// Declare variables at the beginning for fields that need JSON parsing
	needsJsonParsing := !isRepeated && !isMessage && (fieldType == protoreflect.BoolKind ||
		fieldType == protoreflect.StringKind ||
		fieldType == protoreflect.BytesKind ||
		fieldType == protoreflect.Int64Kind ||
		fieldType == protoreflect.Sint64Kind ||
		fieldType == protoreflect.Sfixed64Kind ||
		fieldType == protoreflect.Int32Kind ||
		fieldType == protoreflect.Sint32Kind ||
		fieldType == protoreflect.Sfixed32Kind ||
		fieldType == protoreflect.Uint64Kind ||
		fieldType == protoreflect.Fixed64Kind ||
		fieldType == protoreflect.Uint32Kind ||
		fieldType == protoreflect.Fixed32Kind ||
		fieldType == protoreflect.FloatKind ||
		fieldType == protoreflect.DoubleKind ||
		fieldType == protoreflect.EnumKind)

	if needsJsonParsing {
		content.WriteString("    DECLARE json_value JSON;\n")
	}

	switch {
	case isRepeated:
		// For repeated fields, special handling for float and double to convert from binary format
		if fieldType == protoreflect.FloatKind || fieldType == protoreflect.DoubleKind {
			content.WriteString("    DECLARE result_array JSON DEFAULT JSON_ARRAY();\n")
			content.WriteString("    DECLARE raw_array JSON;\n")
			content.WriteString("    DECLARE array_length INT;\n")
			content.WriteString("    DECLARE i INT DEFAULT 0;\n")
			content.WriteString("    DECLARE element_value JSON;\n")
			content.WriteString("\n")
			content.WriteString(fmt.Sprintf("    SET raw_array = COALESCE(JSON_EXTRACT(proto_data, '$.\"%.d\"'), JSON_ARRAY());\n", field.Number()))
			content.WriteString("    SET array_length = JSON_LENGTH(raw_array);\n")
			content.WriteString("\n")
			content.WriteString("    WHILE i < array_length DO\n")
			content.WriteString("        SET element_value = JSON_EXTRACT(raw_array, CONCAT('$[', i, ']'));\n")

			// Use the appropriate conversion function based on field type
			if fieldType == protoreflect.FloatKind {
				content.WriteString("        SET result_array = JSON_ARRAY_APPEND(result_array, '$', _pb_util_reinterpret_uint32_as_float(_pb_json_parse_float_as_uint32(element_value, TRUE)));\n")
			} else { // DoubleKind
				content.WriteString("        SET result_array = JSON_ARRAY_APPEND(result_array, '$', _pb_util_reinterpret_uint64_as_double(_pb_json_parse_double_as_uint64(element_value, TRUE)));\n")
			}

			content.WriteString("        SET i = i + 1;\n")
			content.WriteString("    END WHILE;\n")
			content.WriteString("\n")
			content.WriteString("    RETURN result_array;\n")
		} else {
			// For other repeated fields, return JSON array or empty array if not present
			content.WriteString(fmt.Sprintf("    RETURN COALESCE(JSON_EXTRACT(proto_data, '$.\"%.d\"'), JSON_ARRAY());\n", field.Number()))
		}
	case isMessage:
		// For message fields, return JSON object or empty object if not present
		content.WriteString(fmt.Sprintf("    RETURN COALESCE(JSON_EXTRACT(proto_data, '$.\"%.d\"'), JSON_OBJECT());\n", field.Number()))
	case fieldType == protoreflect.BoolKind:
		// For boolean fields, use _pb_json_parse_bool if the field exists
		defaultValue := getFieldDefaultValueFromReflection(field)
		content.WriteString(fmt.Sprintf("    SET json_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
		content.WriteString("    IF json_value IS NULL THEN\n")
		content.WriteString(fmt.Sprintf("        RETURN %s;\n", defaultValue))
		content.WriteString("    END IF;\n")
		content.WriteString("    RETURN _pb_json_parse_bool(json_value);\n")
	default:
		// For scalar fields, use appropriate _pb_json_parse functions
		defaultValue := getFieldDefaultValueFromReflection(field)
		content.WriteString(fmt.Sprintf("    SET json_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
		content.WriteString("    IF json_value IS NULL THEN\n")
		content.WriteString(fmt.Sprintf("        RETURN %s;\n", defaultValue))
		content.WriteString("    END IF;\n")
		content.WriteString(fmt.Sprintf("    RETURN %s;\n", getJsonParseFunction(fieldType)))
	}

	content.WriteString("END $$\n\n")
	return nil
}

// getJsonParseFunction returns the appropriate _pb_json_parse function for a field type
func getJsonParseFunction(fieldType protoreflect.Kind) string {
	switch fieldType {
	case protoreflect.BoolKind:
		return "_pb_json_parse_bool(json_value)"
	case protoreflect.StringKind:
		return "_pb_json_parse_string(json_value)"
	case protoreflect.BytesKind:
		return "_pb_json_parse_bytes(json_value)"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
		protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
		protoreflect.EnumKind:
		return "_pb_json_parse_signed_int(json_value)"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "_pb_json_parse_unsigned_int(json_value)"
	case protoreflect.FloatKind:
		return "_pb_util_reinterpret_uint32_as_float(_pb_json_parse_float_as_uint32(json_value, TRUE))"
	case protoreflect.DoubleKind:
		return "_pb_util_reinterpret_uint64_as_double(_pb_json_parse_double_as_uint64(json_value, TRUE))"
	case protoreflect.GroupKind:
		// Groups are deprecated but treated like messages (JSON)
		return "json_value"
	default:
		// This should not happen since we check needsJsonParsing
		return "json_value"
	}
}

// generateEnhancedSetter creates an enhanced setter with input validation
func generateEnhancedSetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string, fieldType protoreflect.Kind, isRepeated, isMessage bool) error {
	setterFuncName := fmt.Sprintf("%s_set_%s", funcPrefix, fieldName)
	if err := validateFunctionName(setterFuncName, fullTypeName); err != nil {
		return err
	}

	setterType, _ := getMySQLTypesForFieldFromReflection(fieldType, isRepeated)

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", setterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, field_value %s) RETURNS JSON DETERMINISTIC\n", setterFuncName, setterType))
	content.WriteString("BEGIN\n")

	// Declare variables at the beginning (required by MySQL syntax)
	oneof := field.ContainingOneof()
	if oneof != nil {
		content.WriteString("    DECLARE temp_data JSON DEFAULT proto_data;\n")
	}

	// Add input validation
	if isRepeated || isMessage {
		content.WriteString("    IF field_value IS NULL THEN\n")
		content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
		content.WriteString("    END IF;\n")
	}

	// Handle oneOf field mutual exclusion
	if oneof != nil {
		content.WriteString("    -- OneOf field mutual exclusion: clear other fields in the same oneOf group\n")

		// Clear all other fields in the same oneOf group
		fields := oneof.Fields()
		for otherField := range protoreflectutils.Iterate(fields) {
			if otherField.Number() != field.Number() {
				content.WriteString(fmt.Sprintf("    SET temp_data = JSON_REMOVE(temp_data, '$.\"%.d\"');\n", otherField.Number()))
			}
		}

		// Set the new field value with proper JSON conversion for booleans, bytes, and floats
		if fieldType == protoreflect.BoolKind {
			content.WriteString(fmt.Sprintf("    RETURN JSON_SET(temp_data, '$.\"%.d\"', CAST((field_value IS TRUE) AS JSON));\n", field.Number()))
		} else if fieldType == protoreflect.BytesKind {
			content.WriteString(fmt.Sprintf("    RETURN JSON_SET(temp_data, '$.\"%.d\"', _pb_to_base64(field_value));\n", field.Number()))
		} else if fieldType == protoreflect.FloatKind {
			content.WriteString(fmt.Sprintf("    RETURN JSON_SET(temp_data, '$.\"%.d\"', _pb_convert_float_uint32_to_number_json(_pb_util_reinterpret_float_as_uint32(field_value)));\n", field.Number()))
		} else if fieldType == protoreflect.DoubleKind {
			content.WriteString(fmt.Sprintf("    RETURN JSON_SET(temp_data, '$.\"%.d\"', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(field_value)));\n", field.Number()))
		} else {
			content.WriteString(fmt.Sprintf("    RETURN JSON_SET(temp_data, '$.\"%.d\"', field_value);\n", field.Number()))
		}
	} else {
		// Check if this is a proto3 field without presence and being set to default value
		// According to protonumberjson spec, such fields should be omitted
		if !field.HasPresence() && !isRepeated && !isMessage {
			// For proto3 fields without presence, omit default values
			content.WriteString("    -- Proto3 field without presence: omit default values per protonumberjson spec\n")

			switch fieldType {
			case protoreflect.BoolKind:
				content.WriteString("    IF field_value IS FALSE THEN\n")
				content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
				content.WriteString("    END IF;\n")
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', CAST((field_value IS TRUE) AS JSON));\n", field.Number()))
			case protoreflect.StringKind:
				content.WriteString("    IF field_value = '' THEN\n")
				content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
				content.WriteString("    END IF;\n")
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
			case protoreflect.BytesKind:
				content.WriteString("    IF field_value = _binary X'' THEN\n")
				content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
				content.WriteString("    END IF;\n")
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', _pb_to_base64(field_value));\n", field.Number()))
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
				protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
				content.WriteString("    IF field_value = 0 THEN\n")
				content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
				content.WriteString("    END IF;\n")
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
			case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind,
				protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				content.WriteString("    IF field_value = 0 THEN\n")
				content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
				content.WriteString("    END IF;\n")
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
			case protoreflect.FloatKind:
				content.WriteString("    IF field_value = 0.0 THEN\n")
				content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
				content.WriteString("    END IF;\n")
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', _pb_convert_float_uint32_to_number_json(_pb_util_reinterpret_float_as_uint32(field_value)));\n", field.Number()))
			case protoreflect.DoubleKind:
				content.WriteString("    IF field_value = 0.0 THEN\n")
				content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
				content.WriteString("    END IF;\n")
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(field_value)));\n", field.Number()))
			case protoreflect.EnumKind:
				content.WriteString("    IF field_value = 0 THEN\n")
				content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
				content.WriteString("    END IF;\n")
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
			default:
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
			}
		} else {
			// Fields with presence: always set the value
			if fieldType == protoreflect.BoolKind {
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', CAST((field_value IS TRUE) AS JSON));\n", field.Number()))
			} else if fieldType == protoreflect.BytesKind {
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', _pb_to_base64(field_value));\n", field.Number()))
			} else if fieldType == protoreflect.FloatKind {
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', _pb_convert_float_uint32_to_number_json(_pb_util_reinterpret_float_as_uint32(field_value)));\n", field.Number()))
			} else if fieldType == protoreflect.DoubleKind {
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(field_value)));\n", field.Number()))
			} else {
				content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
			}
		}
	}

	content.WriteString("END $$\n\n")

	return nil
}

// generateHasMethod creates a method to check field presence
func generateHasMethod(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string) error {
	hasFuncName := fmt.Sprintf("%s_has_%s", funcPrefix, fieldName)
	if err := validateFunctionName(hasFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", hasFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS BOOLEAN DETERMINISTIC\n", hasFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN JSON_CONTAINS_PATH(proto_data, 'one', '$.\"%.d\"');\n", field.Number()))
	content.WriteString("END $$\n\n")

	return nil
}

// generateClearMethod creates a method to clear/unset a field
func generateClearMethod(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string) error {
	clearFuncName := fmt.Sprintf("%s_clear_%s", funcPrefix, fieldName)
	if err := validateFunctionName(clearFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", clearFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS JSON DETERMINISTIC\n", clearFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("END $$\n\n")

	return nil
}

// generateRepeatedFieldMethods creates additional methods for repeated fields
func generateRepeatedFieldMethods(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string, fieldType protoreflect.Kind) error {
	// Generate add method for repeated fields
	addFuncName := fmt.Sprintf("%s_add_%s", funcPrefix, fieldName)
	if err := validateFunctionName(addFuncName, fullTypeName); err != nil {
		return err
	}

	elementSetterType, _ := getMySQLTypesForFieldFromReflection(fieldType, false) // Get non-repeated type for element

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", addFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, element_value %s) RETURNS JSON DETERMINISTIC\n", addFuncName, elementSetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE current_array JSON;\n")
	content.WriteString(fmt.Sprintf("    SET current_array = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF current_array IS NULL THEN\n")
	content.WriteString("        SET current_array = JSON_ARRAY();\n")
	content.WriteString("    END IF;\n")

	// Handle type-specific conversions for proper protonumberjson format in arrays
	switch fieldType {
	case protoreflect.BoolKind:
		content.WriteString("    SET current_array = JSON_ARRAY_APPEND(current_array, '$', CAST((element_value IS TRUE) AS JSON));\n")
	case protoreflect.BytesKind:
		content.WriteString("    SET current_array = JSON_ARRAY_APPEND(current_array, '$', _pb_to_base64(element_value));\n")
	case protoreflect.FloatKind:
		content.WriteString("    SET current_array = JSON_ARRAY_APPEND(current_array, '$', _pb_convert_float_uint32_to_number_json(_pb_util_reinterpret_float_as_uint32(element_value)));\n")
	case protoreflect.DoubleKind:
		content.WriteString("    SET current_array = JSON_ARRAY_APPEND(current_array, '$', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(element_value)));\n")
	default:
		content.WriteString("    SET current_array = JSON_ARRAY_APPEND(current_array, '$', element_value);\n")
	}

	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', current_array);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate count method for repeated fields
	countFuncName := fmt.Sprintf("%s_count_%s", funcPrefix, fieldName)
	if err := validateFunctionName(countFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", countFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS INT DETERMINISTIC\n", countFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE array_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF array_value IS NULL THEN\n")
	content.WriteString("        RETURN 0;\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    RETURN JSON_LENGTH(array_value);\n")
	content.WriteString("END $$\n\n")

	return nil
}

// generateOneOfMethods creates methods for oneOf groups
func generateOneOfMethods(content *strings.Builder, messageDesc protoreflect.MessageDescriptor, funcPrefix string, fullTypeName protoreflect.FullName) error {
	oneofs := messageDesc.Oneofs()
	for oneof := range protoreflectutils.Iterate(oneofs) {
		oneofName := string(oneof.Name())

		// Generate which method to detect which field is set in the oneOf group
		whichFuncName := fmt.Sprintf("%s_which_%s", funcPrefix, oneofName)
		if err := validateFunctionName(whichFuncName, fullTypeName); err != nil {
			return err
		}

		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", whichFuncName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS LONGTEXT DETERMINISTIC\n", whichFuncName))
		content.WriteString("BEGIN\n")

		// Check each field in the oneOf group
		fields := oneof.Fields()
		for field := range protoreflectutils.Iterate(fields) {
			fieldName := string(field.Name())
			content.WriteString(fmt.Sprintf("    IF JSON_CONTAINS_PATH(proto_data, 'one', '$.\"%.d\"') THEN\n", field.Number()))
			content.WriteString(fmt.Sprintf("        RETURN '%s';\n", fieldName))
			content.WriteString("    END IF;\n")
		}

		content.WriteString("    RETURN NULL; -- No field is set in this oneOf group\n")
		content.WriteString("END $$\n\n")

		// Generate clear method for entire oneOf group
		clearOneOfFuncName := fmt.Sprintf("%s_clear_%s", funcPrefix, oneofName)
		if err := validateFunctionName(clearOneOfFuncName, fullTypeName); err != nil {
			return err
		}

		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", clearOneOfFuncName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS JSON DETERMINISTIC\n", clearOneOfFuncName))
		content.WriteString("BEGIN\n")
		content.WriteString("    DECLARE temp_data JSON DEFAULT proto_data;\n")

		// Clear all fields in the oneOf group
		for field := range protoreflectutils.Iterate(fields) {
			content.WriteString(fmt.Sprintf("    SET temp_data = JSON_REMOVE(temp_data, '$.\"%.d\"');\n", field.Number()))
		}

		content.WriteString("    RETURN temp_data;\n")
		content.WriteString("END $$\n\n")
	}

	return nil
}

func generateEnumMethods(content *strings.Builder, enumDesc protoreflect.EnumDescriptor, typePrefixFunc TypePrefixFunc, _ string) error {
	// Use FullName from descriptor
	fullTypeName := enumDesc.FullName()
	packageName := enumDesc.ParentFile().Package()

	// Get prefix for this specific type
	funcPrefix := typePrefixFunc(packageName, fullTypeName)

	// Generate from_string method
	fromStringFuncName := funcPrefix + "_from_string"
	if err := validateFunctionName(fromStringFuncName, fullTypeName); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", fromStringFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(enum_name LONGTEXT) RETURNS INT DETERMINISTIC\n", fromStringFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    CASE enum_name\n")

	// Generate CASE statements for each enum value
	enumValues := enumDesc.Values()
	for enumValue := range protoreflectutils.Iterate(enumValues) {
		valueName := string(enumValue.Name())
		valueNumber := int(enumValue.Number())
		content.WriteString(fmt.Sprintf("        WHEN '%s' THEN RETURN %d;\n", valueName, valueNumber))
	}

	content.WriteString("        ELSE RETURN NULL;\n")
	content.WriteString("    END CASE;\n")
	content.WriteString("END $$\n\n")

	// Generate to_string method
	toStringFuncName := funcPrefix + "_to_string"
	if err := validateFunctionName(toStringFuncName, fullTypeName); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", toStringFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(enum_value INT) RETURNS LONGTEXT DETERMINISTIC\n", toStringFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    CASE enum_value\n")

	// Generate CASE statements for each enum value
	for enumValue := range protoreflectutils.Iterate(enumValues) {
		valueName := string(enumValue.Name())
		valueNumber := int(enumValue.Number())
		content.WriteString(fmt.Sprintf("        WHEN %d THEN RETURN '%s';\n", valueNumber, valueName))
	}

	content.WriteString("        ELSE RETURN NULL;\n")
	content.WriteString("    END CASE;\n")
	content.WriteString("END $$\n\n")

	return nil
}

// getMySQLTypesForFieldFromReflection returns the appropriate MySQL types for setter parameter and getter return value using protoreflect
func getMySQLTypesForFieldFromReflection(fieldKind protoreflect.Kind, isRepeated bool) (setterType, getterType string) {
	// For repeated fields and maps, always use JSON
	if isRepeated {
		return "JSON", "JSON"
	}

	// For scalar fields, map protobuf types to MySQL types
	switch fieldKind {
	case protoreflect.DoubleKind:
		return "DOUBLE", "DOUBLE"
	case protoreflect.FloatKind:
		return "FLOAT", "FLOAT"
	case protoreflect.Int64Kind:
		return "BIGINT", "BIGINT"
	case protoreflect.Uint64Kind:
		return "BIGINT UNSIGNED", "BIGINT UNSIGNED"
	case protoreflect.Int32Kind:
		return "INT", "INT"
	case protoreflect.Fixed64Kind:
		return "BIGINT UNSIGNED", "BIGINT UNSIGNED"
	case protoreflect.Fixed32Kind:
		return "INT UNSIGNED", "INT UNSIGNED"
	case protoreflect.BoolKind:
		return "BOOLEAN", "BOOLEAN"
	case protoreflect.StringKind:
		return "LONGTEXT", "LONGTEXT"
	case protoreflect.BytesKind:
		return "LONGBLOB", "LONGBLOB"
	case protoreflect.Uint32Kind:
		return "INT UNSIGNED", "INT UNSIGNED"
	case protoreflect.EnumKind:
		return "INT", "INT"
	case protoreflect.Sfixed32Kind:
		return "INT", "INT"
	case protoreflect.Sfixed64Kind:
		return "BIGINT", "BIGINT"
	case protoreflect.Sint32Kind:
		return "INT", "INT"
	case protoreflect.Sint64Kind:
		return "BIGINT", "BIGINT"
	case protoreflect.MessageKind:
		// Messages are represented as JSON (ProtoNumberJSON format)
		return "JSON", "JSON"
	case protoreflect.GroupKind:
		// Groups are deprecated but still need handling - represent as JSON
		return "JSON", "JSON"
	default:
		// Fallback to JSON for unknown types
		return "JSON", "JSON"
	}
}

// getFieldDefaultValueFromReflection returns the appropriate default value for a field based on proto version and custom defaults using protoreflect
func getFieldDefaultValueFromReflection(field protoreflect.FieldDescriptor) string {
	fieldKind := field.Kind()

	// Check if field has a custom default value (proto2)
	if field.HasDefault() {
		defaultVal := field.Default()
		switch fieldKind {
		case protoreflect.StringKind:
			// String defaults need to be quoted
			return fmt.Sprintf("'%s'", strings.ReplaceAll(defaultVal.String(), "'", "''"))
		case protoreflect.BytesKind:
			// Bytes defaults need to be hex-encoded
			return fmt.Sprintf("X'%X'", defaultVal.Bytes())
		case protoreflect.BoolKind:
			// Boolean defaults
			if defaultVal.Bool() {
				return "1"
			}
			return "0"
		case protoreflect.EnumKind:
			// Enum defaults
			return fmt.Sprintf("%d", defaultVal.Enum())
		case protoreflect.MessageKind:
			// Messages don't have default values in proto2
			return "NULL"
		case protoreflect.GroupKind:
			// Groups don't have default values
			return "NULL"
		default:
			// Numeric defaults can be used directly
			switch fieldKind {
			case protoreflect.DoubleKind, protoreflect.FloatKind:
				return fmt.Sprintf("%g", defaultVal.Float())
			case protoreflect.Int64Kind, protoreflect.Sfixed64Kind, protoreflect.Sint64Kind,
				protoreflect.Int32Kind, protoreflect.Sfixed32Kind, protoreflect.Sint32Kind:
				return fmt.Sprintf("%d", defaultVal.Int())
			case protoreflect.Uint64Kind, protoreflect.Fixed64Kind,
				protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
				return fmt.Sprintf("%d", defaultVal.Uint())
			case protoreflect.BoolKind:
				if defaultVal.Bool() {
					return "1"
				}
				return "0"
			case protoreflect.EnumKind:
				return fmt.Sprintf("%d", defaultVal.Enum())
			case protoreflect.StringKind:
				return fmt.Sprintf("'%s'", strings.ReplaceAll(defaultVal.String(), "'", "''"))
			case protoreflect.BytesKind:
				return fmt.Sprintf("X'%X'", defaultVal.Bytes())
			case protoreflect.MessageKind, protoreflect.GroupKind:
				return "NULL"
			default:
				return "NULL"
			}
		}
	}

	// Use proto3 zero values or proto2 implicit defaults
	switch fieldKind {
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		return "0.0"
	case protoreflect.Int64Kind, protoreflect.Uint64Kind,
		protoreflect.Int32Kind, protoreflect.Uint32Kind,
		protoreflect.Fixed64Kind, protoreflect.Fixed32Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.EnumKind:
		return "0"
	case protoreflect.BoolKind:
		return "0"
	case protoreflect.StringKind:
		return "''"
	case protoreflect.BytesKind:
		return "X''"
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return "NULL"
	default:
		return "NULL"
	}
}

// validateFunctionName validates that a function name doesn't exceed MySQL's 64 character limit
func validateFunctionName(functionName string, fullTypeName protoreflect.FullName) error {
	const maxFunctionNameLength = 64

	if len(functionName) > maxFunctionNameLength {
		return fmt.Errorf("generated function name '%s' exceeds MySQL's %d character limit. "+
			"Use the prefix_map option to assign shorter prefixes to package '%s' or type '%s'. "+
			"For example: --mysql_opt=prefix_map='%s=short_prefix_' or --mysql_opt=prefix_map='%s=short_'",
			functionName, maxFunctionNameLength, fullTypeName.Parent(), fullTypeName, fullTypeName.Parent(), fullTypeName)
	}

	return nil
}
