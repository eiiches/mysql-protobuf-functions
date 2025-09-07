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
func GenerateMethodFragments(files *protoregistry.Files, fileNameFunc FileNameFunc, typePrefixFunc TypePrefixFunc, schemaFunctionName string, fieldFilterFunc FieldFilterFunc) (map[string][]string, error) {
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
		content, err := generateMethodsForFile(fileDesc, typePrefixFunc, schemaFunctionName, fieldFilterFunc)
		if err != nil {
			return nil, err
		}
		if content != "" {
			fileFragments[filename] = append(fileFragments[filename], content)
		}
	}

	return fileFragments, nil
}

func generateMethodsForFile(fileDesc protoreflect.FileDescriptor, typePrefixFunc TypePrefixFunc, schemaFunctionName string, fieldFilterFunc FieldFilterFunc) (string, error) {
	var content strings.Builder

	// Generate methods for each message type
	messages := fileDesc.Messages()
	for messageDesc := range protoreflectutils.Iterate(messages) {
		if err := generateMessageMethods(&content, messageDesc, typePrefixFunc, schemaFunctionName, fieldFilterFunc); err != nil {
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

func generateMessageMethods(content *strings.Builder, messageDesc protoreflect.MessageDescriptor, typePrefixFunc TypePrefixFunc, schemaFunctionName string, fieldFilterFunc FieldFilterFunc) error {
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
	if err := generateFieldAccessorMethods(content, messageDesc, funcPrefix, fullTypeName, fieldFilterFunc, typePrefixFunc); err != nil {
		return err
	}

	// Generate oneOf methods for oneOf groups
	if err := generateOneOfMethods(content, messageDesc, funcPrefix, fullTypeName); err != nil {
		return err
	}

	// Generate methods for nested message types
	nestedMessages := messageDesc.Messages()
	for nestedMessageDesc := range protoreflectutils.Iterate(nestedMessages) {
		if err := generateMessageMethods(content, nestedMessageDesc, typePrefixFunc, schemaFunctionName, fieldFilterFunc); err != nil {
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
func generateFieldAccessorMethods(content *strings.Builder, messageDesc protoreflect.MessageDescriptor, funcPrefix string, fullTypeName protoreflect.FullName, fieldFilterFunc FieldFilterFunc, typePrefixFunc TypePrefixFunc) error {
	fields := messageDesc.Fields()

	for field := range protoreflectutils.Iterate(fields) {
		fieldName := string(field.Name())
		fieldType := field.Kind()
		isRepeated := field.Cardinality() == protoreflect.Repeated
		isMessage := fieldType == protoreflect.MessageKind
		isMapField := isRepeated && isMessage && field.MapKey() != nil

		// Generate enhanced getter with better defaults and type safety
		if err := generateEnhancedGetter(content, funcPrefix, fullTypeName, field, fieldName, fieldType, isRepeated, isMessage, isMapField); err != nil {
			return err
		}

		// Generate nullable getter with custom default for optional fields
		if field.HasPresence() && !isRepeated && !isMapField {
			if err := generateNullableGetter(content, funcPrefix, fullTypeName, field, fieldName, fieldType, isMessage, fieldFilterFunc); err != nil {
				return err
			}
		}

		// Generate enum name getters for enum fields
		if fieldType == protoreflect.EnumKind && !isRepeated && !isMapField {
			// Generate regular enum name getter (__as_name)
			if err := generateEnumNameGetter(content, funcPrefix, fullTypeName, field, fieldName, false, fieldFilterFunc, typePrefixFunc); err != nil {
				return err
			}
			// Generate nullable enum name getter (__as_name_or) for optional fields
			if field.HasPresence() {
				if err := generateEnumNameGetter(content, funcPrefix, fullTypeName, field, fieldName, true, fieldFilterFunc, typePrefixFunc); err != nil {
					return err
				}
			}
		}

		// Generate enum name setters for enum fields
		if fieldType == protoreflect.EnumKind && !isRepeated && !isMapField {
			if err := generateEnumNameSetter(content, funcPrefix, fullTypeName, field, fieldName, fieldFilterFunc, typePrefixFunc); err != nil {
				return err
			}
		}

		// Generate enhanced setter with validation
		if err := generateEnhancedSetter(content, funcPrefix, fullTypeName, field, fieldName, fieldType, isRepeated, isMessage, isMapField); err != nil {
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

		// Generate additional methods for repeated fields (but not map fields)
		if isRepeated && !isMapField {
			if err := generateRepeatedFieldMethods(content, funcPrefix, fullTypeName, field, fieldName, fieldType); err != nil {
				return err
			}
		}

		// Generate additional methods for map fields
		if isMapField {
			if err := generateMapFieldMethods(content, funcPrefix, fullTypeName, field, fieldName); err != nil {
				return err
			}
		}
	}

	return nil
}

// generateEnhancedGetter creates an enhanced getter with better defaults and type safety
func generateEnhancedGetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string, fieldType protoreflect.Kind, isRepeated, isMessage, isMapField bool) error {
	// Use new naming: get_all_ for repeated/map fields, get_ for singular fields
	var getterFuncName string
	if isRepeated || isMapField {
		getterFuncName = fmt.Sprintf("%s_get_all_%s", funcPrefix, fieldName)
	} else {
		getterFuncName = fmt.Sprintf("%s_get_%s", funcPrefix, fieldName)
	}
	if err := validateFunctionName(getterFuncName, fullTypeName); err != nil {
		return err
	}

	_, getterType := getMySQLTypesForFieldFromReflection(fieldType, isRepeated || isMapField)

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
	case isMapField:
		// For map fields, return JSON object or empty array if not present
		content.WriteString(fmt.Sprintf("    RETURN COALESCE(JSON_EXTRACT(proto_data, '$.\"%.d\"'), JSON_ARRAY());\n", field.Number()))
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

// generateNullableGetter creates a nullable getter with custom default for optional fields
func generateNullableGetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string, fieldType protoreflect.Kind, isMessage bool, fieldFilterFunc FieldFilterFunc) error {
	// Create function name with __or modifier
	getterFuncName := fmt.Sprintf("%s_get_%s__or", funcPrefix, fieldName)

	_, getterType := getMySQLTypesForFieldFromReflection(fieldType, false)

	// Determine how this function should be generated
	decision := DecisionInclude
	if fieldFilterFunc != nil {
		decision = fieldFilterFunc(field, getterFuncName)
	}

	// Only validate function names for functions that will be generated normally
	if decision == DecisionInclude {
		if err := validateFunctionName(getterFuncName, fullTypeName); err != nil {
			return err
		}
	}

	// Skip generation entirely if excluded
	if decision == DecisionExclude {
		return nil
	}

	// Generate function (commented out if requested)
	commentPrefix := ""
	if decision == DecisionCommentOut {
		commentPrefix = "-- SKIPPED: "
		content.WriteString(fmt.Sprintf("-- SKIPPED: Function '%s' was filtered out\n", getterFuncName))
	}

	content.WriteString(fmt.Sprintf("%sDROP FUNCTION IF EXISTS %s $$\n", commentPrefix, getterFuncName))
	content.WriteString(fmt.Sprintf("%sCREATE FUNCTION %s(proto_data JSON, default_value %s) RETURNS %s DETERMINISTIC\n", commentPrefix, getterFuncName, getterType, getterType))
	content.WriteString(fmt.Sprintf("%sBEGIN\n", commentPrefix))

	// Field is present, return actual value using same logic as regular getter
	needsJsonParsing := !isMessage && (fieldType == protoreflect.BoolKind ||
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

	// Declare variables at the beginning of the function
	if needsJsonParsing {
		content.WriteString(fmt.Sprintf("%s    DECLARE json_value JSON;\n", commentPrefix))
	}

	// Get field value and check for presence using simpler IS NULL check
	if needsJsonParsing {
		content.WriteString(fmt.Sprintf("%s    SET json_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", commentPrefix, field.Number()))
		content.WriteString(fmt.Sprintf("%s    IF json_value IS NOT NULL THEN\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s        RETURN %s;\n", commentPrefix, getJsonParseFunction(fieldType)))
	} else {
		// For message fields and others, extract directly and check for null
		content.WriteString(fmt.Sprintf("%s    DECLARE field_value JSON;\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s    SET field_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", commentPrefix, field.Number()))
		content.WriteString(fmt.Sprintf("%s    IF field_value IS NOT NULL THEN\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s        RETURN field_value;\n", commentPrefix))
	}

	// Field is not present, return provided default
	content.WriteString(fmt.Sprintf("%s    ELSE\n", commentPrefix))
	content.WriteString(fmt.Sprintf("%s        RETURN default_value;\n", commentPrefix))
	content.WriteString(fmt.Sprintf("%s    END IF;\n", commentPrefix))

	content.WriteString(fmt.Sprintf("%sEND $$\n\n", commentPrefix))
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
	case protoreflect.MessageKind:
		// Message fields return JSON directly
		return "json_value"
	case protoreflect.GroupKind:
		// Groups are deprecated but treated like messages (JSON)
		return "json_value"
	default:
		// This should not happen since we check needsJsonParsing
		return "json_value"
	}
}

// generateEnhancedSetter creates an enhanced setter with input validation
func generateEnhancedSetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string, fieldType protoreflect.Kind, isRepeated, isMessage, isMapField bool) error {
	// Use new naming: set_all_ for repeated/map fields, set_ for singular fields
	var setterFuncName string
	if isRepeated || isMapField {
		setterFuncName = fmt.Sprintf("%s_set_all_%s", funcPrefix, fieldName)
	} else {
		setterFuncName = fmt.Sprintf("%s_set_%s", funcPrefix, fieldName)
	}
	if err := validateFunctionName(setterFuncName, fullTypeName); err != nil {
		return err
	}

	setterType, _ := getMySQLTypesForFieldFromReflection(fieldType, isRepeated || isMapField)

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", setterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, field_value %s) RETURNS JSON DETERMINISTIC\n", setterFuncName, setterType))
	content.WriteString("BEGIN\n")

	// Declare variables at the beginning (required by MySQL syntax)
	oneof := field.ContainingOneof()
	if oneof != nil {
		content.WriteString("    DECLARE temp_data JSON DEFAULT proto_data;\n")
	}

	// Add input validation
	if isMapField {
		// For map fields, handle JSON object input directly
		content.WriteString("    IF field_value IS NULL THEN\n")
		content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
		content.WriteString("    END IF;\n")
		content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
		content.WriteString("END $$\n\n")
		return nil
	} else if isRepeated {
		// For repeated fields, handle JSON array input
		content.WriteString("    DECLARE array_length INT;\n")
		content.WriteString("    DECLARE i INT DEFAULT 0;\n")
		content.WriteString("    DECLARE element_value JSON;\n")
		content.WriteString("    DECLARE converted_array JSON DEFAULT JSON_ARRAY();\n")
		content.WriteString("\n")
		content.WriteString("    IF field_value IS NULL THEN\n")
		content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
		content.WriteString("    END IF;\n")
		content.WriteString("\n")
		content.WriteString("    SET array_length = JSON_LENGTH(field_value);\n")
		content.WriteString("\n")
		content.WriteString("    -- Handle empty array case\n")
		content.WriteString("    IF array_length = 0 THEN\n")
		content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
		content.WriteString("    END IF;\n")
		content.WriteString("\n")
		content.WriteString("    -- Convert each element to internal format\n")
		content.WriteString("    WHILE i < array_length DO\n")
		content.WriteString("        SET element_value = JSON_EXTRACT(field_value, CONCAT('$[', i, ']'));\n")

		// Handle type-specific conversions for proper internal format
		switch fieldType {
		case protoreflect.BoolKind:
			content.WriteString("        SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', element_value);\n")
		case protoreflect.BytesKind:
			content.WriteString("        SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', element_value); -- Expect base64 strings\n")
		case protoreflect.FloatKind:
			content.WriteString("        SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', _pb_convert_float_uint32_to_number_json(_pb_util_reinterpret_float_as_uint32(CAST(element_value AS DOUBLE))));\n")
		case protoreflect.DoubleKind:
			content.WriteString("        SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(CAST(element_value AS DOUBLE))));\n")
		default:
			content.WriteString("        SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', element_value);\n")
		}

		content.WriteString("        SET i = i + 1;\n")
		content.WriteString("    END WHILE;\n")
		content.WriteString("\n")
		content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', converted_array);\n", field.Number()))
		content.WriteString("END $$\n\n")
		return nil
	} else if isMessage {
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

	// Generate index-based get method for repeated fields
	getIndexFuncName := fmt.Sprintf("%s_get_%s", funcPrefix, fieldName)
	if err := validateFunctionName(getIndexFuncName, fullTypeName); err != nil {
		return err
	}

	_, elementGetterType := getMySQLTypesForFieldFromReflection(fieldType, false) // Get non-repeated type for element

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getIndexFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, index_value INT) RETURNS %s DETERMINISTIC\n", getIndexFuncName, elementGetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE array_value JSON;\n")
	content.WriteString("    DECLARE element_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF array_value IS NULL OR JSON_LENGTH(array_value) <= index_value OR index_value < 0 THEN\n")
	content.WriteString("        RETURN NULL;\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    SET element_value = JSON_EXTRACT(array_value, CONCAT('$[', index_value, ']'));\n")

	// Handle type-specific parsing for return value
	switch fieldType {
	case protoreflect.BoolKind:
		content.WriteString("    RETURN _pb_json_parse_bool(element_value);\n")
	case protoreflect.StringKind:
		content.WriteString("    RETURN JSON_UNQUOTE(element_value);\n")
	case protoreflect.BytesKind:
		content.WriteString("    RETURN FROM_BASE64(JSON_UNQUOTE(element_value));\n")
	case protoreflect.DoubleKind:
		content.WriteString("    RETURN _pb_util_reinterpret_uint64_as_double(_pb_json_parse_double_as_uint64(element_value, TRUE));\n")
	case protoreflect.FloatKind:
		content.WriteString("    RETURN _pb_util_reinterpret_uint32_as_float(_pb_json_parse_float_as_uint32(element_value, TRUE));\n")
	case protoreflect.MessageKind:
		content.WriteString("    RETURN element_value;\n")
	default:
		// For integers and enums, direct JSON parsing
		content.WriteString("    RETURN CAST(element_value AS SIGNED);\n")
	}

	content.WriteString("END $$\n\n")

	// Generate index-based set method for repeated fields
	setIndexFuncName := fmt.Sprintf("%s_set_%s", funcPrefix, fieldName)
	if err := validateFunctionName(setIndexFuncName, fullTypeName); err != nil {
		return err
	}

	elementSetterType, _ = getMySQLTypesForFieldFromReflection(fieldType, false) // Get non-repeated type for element

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", setIndexFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, index_value INT, element_value %s) RETURNS JSON DETERMINISTIC\n", setIndexFuncName, elementSetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE array_value JSON;\n")
	content.WriteString("    DECLARE array_length INT;\n")
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF array_value IS NULL THEN\n")
	content.WriteString("        RETURN proto_data; -- Cannot set at index in non-existent array\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    SET array_length = JSON_LENGTH(array_value);\n")
	content.WriteString("    IF index_value < 0 OR index_value >= array_length THEN\n")
	content.WriteString("        RETURN proto_data; -- Index out of bounds\n")
	content.WriteString("    END IF;\n")

	// Handle type-specific conversions for setting
	switch fieldType {
	case protoreflect.BoolKind:
		content.WriteString("    SET array_value = JSON_SET(array_value, CONCAT('$[', index_value, ']'), CAST((element_value IS TRUE) AS JSON));\n")
	case protoreflect.StringKind:
		content.WriteString("    SET array_value = JSON_SET(array_value, CONCAT('$[', index_value, ']'), element_value);\n")
	case protoreflect.BytesKind:
		content.WriteString("    SET array_value = JSON_SET(array_value, CONCAT('$[', index_value, ']'), _pb_to_base64(element_value));\n")
	case protoreflect.DoubleKind:
		content.WriteString("    SET array_value = JSON_SET(array_value, CONCAT('$[', index_value, ']'), _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(element_value)));\n")
	case protoreflect.FloatKind:
		content.WriteString("    SET array_value = JSON_SET(array_value, CONCAT('$[', index_value, ']'), _pb_convert_float_uint32_to_number_json(_pb_util_reinterpret_float_as_uint32(element_value)));\n")
	default:
		// For integers and enums, direct JSON casting
		content.WriteString("    SET array_value = JSON_SET(array_value, CONCAT('$[', index_value, ']'), CAST(element_value AS JSON));\n")
	}

	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', array_value);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate insert method for repeated fields
	insertFuncName := fmt.Sprintf("%s_insert_%s", funcPrefix, fieldName)
	if err := validateFunctionName(insertFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", insertFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, index_value INT, element_value %s) RETURNS JSON DETERMINISTIC\n", insertFuncName, elementSetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE array_value JSON;\n")
	content.WriteString("    DECLARE array_length INT;\n")
	content.WriteString("    DECLARE new_array JSON DEFAULT JSON_ARRAY();\n")
	content.WriteString("    DECLARE i INT DEFAULT 0;\n")
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF array_value IS NULL THEN\n")
	content.WriteString("        SET array_value = JSON_ARRAY();\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    SET array_length = JSON_LENGTH(array_value);\n")
	content.WriteString("    IF index_value < 0 THEN\n")
	content.WriteString("        SET index_value = 0;\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF index_value > array_length THEN\n")
	content.WriteString("        SET index_value = array_length;\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    \n")
	content.WriteString("    -- Copy elements before insert position\n")
	content.WriteString("    WHILE i < index_value DO\n")
	content.WriteString("        SET new_array = JSON_ARRAY_APPEND(new_array, '$', JSON_EXTRACT(array_value, CONCAT('$[', i, ']')));\n")
	content.WriteString("        SET i = i + 1;\n")
	content.WriteString("    END WHILE;\n")
	content.WriteString("    \n")
	content.WriteString("    -- Insert new element\n")

	// Handle type-specific conversions for inserting
	switch fieldType {
	case protoreflect.BoolKind:
		content.WriteString("    SET new_array = JSON_ARRAY_APPEND(new_array, '$', CAST((element_value IS TRUE) AS JSON));\n")
	case protoreflect.StringKind:
		content.WriteString("    SET new_array = JSON_ARRAY_APPEND(new_array, '$', element_value);\n")
	case protoreflect.BytesKind:
		content.WriteString("    SET new_array = JSON_ARRAY_APPEND(new_array, '$', _pb_to_base64(element_value));\n")
	case protoreflect.DoubleKind:
		content.WriteString("    SET new_array = JSON_ARRAY_APPEND(new_array, '$', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(element_value)));\n")
	case protoreflect.FloatKind:
		content.WriteString("    SET new_array = JSON_ARRAY_APPEND(new_array, '$', _pb_convert_float_uint32_to_number_json(_pb_util_reinterpret_float_as_uint32(element_value)));\n")
	default:
		// For integers and enums, direct JSON casting
		content.WriteString("    SET new_array = JSON_ARRAY_APPEND(new_array, '$', CAST(element_value AS JSON));\n")
	}

	content.WriteString("    \n")
	content.WriteString("    -- Copy remaining elements after insert position\n")
	content.WriteString("    WHILE i < array_length DO\n")
	content.WriteString("        SET new_array = JSON_ARRAY_APPEND(new_array, '$', JSON_EXTRACT(array_value, CONCAT('$[', i, ']')));\n")
	content.WriteString("        SET i = i + 1;\n")
	content.WriteString("    END WHILE;\n")
	content.WriteString("    \n")
	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', new_array);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate remove method for repeated fields
	removeFuncName := fmt.Sprintf("%s_remove_%s", funcPrefix, fieldName)
	if err := validateFunctionName(removeFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", removeFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, index_value INT) RETURNS JSON DETERMINISTIC\n", removeFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE array_value JSON;\n")
	content.WriteString("    DECLARE array_length INT;\n")
	content.WriteString("    DECLARE new_array JSON DEFAULT JSON_ARRAY();\n")
	content.WriteString("    DECLARE i INT DEFAULT 0;\n")
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF array_value IS NULL THEN\n")
	content.WriteString("        RETURN proto_data;\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    SET array_length = JSON_LENGTH(array_value);\n")
	content.WriteString("    IF index_value < 0 OR index_value >= array_length THEN\n")
	content.WriteString("        RETURN proto_data; -- Index out of bounds\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    \n")
	content.WriteString("    -- Copy all elements except the one to remove\n")
	content.WriteString("    WHILE i < array_length DO\n")
	content.WriteString("        IF i != index_value THEN\n")
	content.WriteString("            SET new_array = JSON_ARRAY_APPEND(new_array, '$', JSON_EXTRACT(array_value, CONCAT('$[', i, ']')));\n")
	content.WriteString("        END IF;\n")
	content.WriteString("        SET i = i + 1;\n")
	content.WriteString("    END WHILE;\n")
	content.WriteString("    \n")
	content.WriteString("    -- If array becomes empty, remove the field entirely (proto3 default value omission)\n")
	content.WriteString("    IF JSON_LENGTH(new_array) = 0 THEN\n")
	content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    END IF;\n")
	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', new_array);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate add_all method for repeated fields
	addAllFuncName := fmt.Sprintf("%s_add_all_%s", funcPrefix, fieldName)
	if err := validateFunctionName(addAllFuncName, fullTypeName); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", addAllFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, elements_array JSON) RETURNS JSON DETERMINISTIC\n", addAllFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE current_array JSON;\n")
	content.WriteString("    DECLARE elements_length INT;\n")
	content.WriteString("    DECLARE i INT DEFAULT 0;\n")
	content.WriteString("    DECLARE element_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET current_array = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF current_array IS NULL THEN\n")
	content.WriteString("        SET current_array = JSON_ARRAY();\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF elements_array IS NULL OR JSON_TYPE(elements_array) != 'ARRAY' THEN\n")
	content.WriteString("        RETURN proto_data; -- Invalid input, return unchanged\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    SET elements_length = JSON_LENGTH(elements_array);\n")
	content.WriteString("    WHILE i < elements_length DO\n")
	content.WriteString("        SET element_value = JSON_EXTRACT(elements_array, CONCAT('$[', i, ']'));\n")

	// Handle type-specific conversions for adding elements
	switch fieldType {
	case protoreflect.BoolKind:
		content.WriteString("        SET current_array = JSON_ARRAY_APPEND(current_array, '$', CAST((_pb_json_parse_bool(element_value) IS TRUE) AS JSON));\n")
	case protoreflect.BytesKind:
		content.WriteString("        SET current_array = JSON_ARRAY_APPEND(current_array, '$', element_value);\n") // Assume already base64 encoded
	case protoreflect.FloatKind:
		content.WriteString("        SET current_array = JSON_ARRAY_APPEND(current_array, '$', _pb_convert_float_uint32_to_number_json(_pb_util_reinterpret_float_as_uint32(CAST(element_value AS DOUBLE))));\n")
	case protoreflect.DoubleKind:
		content.WriteString("        SET current_array = JSON_ARRAY_APPEND(current_array, '$', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(CAST(element_value AS DOUBLE))));\n")
	default:
		content.WriteString("        SET current_array = JSON_ARRAY_APPEND(current_array, '$', element_value);\n")
	}

	content.WriteString("        SET i = i + 1;\n")
	content.WriteString("    END WHILE;\n")
	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', current_array);\n", field.Number()))
	content.WriteString("END $$\n\n")

	return nil
}

// generateMapFieldMethods creates additional methods for map fields
func generateMapFieldMethods(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string) error {
	keyType := field.MapKey().Kind()
	valueType := field.MapValue().Kind()

	// Get MySQL types for key and value
	keySetterType, _ := getMySQLTypesForFieldFromReflection(keyType, false)
	valueSetterType, _ := getMySQLTypesForFieldFromReflection(valueType, false)
	_, valueGetterType := getMySQLTypesForFieldFromReflection(valueType, false)

	// Generate key-based get method for map fields
	getKeyFuncName := fmt.Sprintf("%s_get_%s", funcPrefix, fieldName)
	if err := validateFunctionName(getKeyFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getKeyFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, key_value %s) RETURNS %s DETERMINISTIC\n", getKeyFuncName, keySetterType, valueGetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE map_value JSON;\n")
	content.WriteString("    DECLARE element_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")

	// Return appropriate default value for the value type
	switch valueType {
	case protoreflect.BoolKind:
		content.WriteString("        RETURN FALSE;\n")
	case protoreflect.StringKind, protoreflect.BytesKind:
		content.WriteString("        RETURN '';\n")
	case protoreflect.EnumKind:
		content.WriteString("        RETURN 0;\n")
	case protoreflect.MessageKind:
		content.WriteString("        RETURN JSON_OBJECT();\n")
	default:
		// All numeric types default to 0
		content.WriteString("        RETURN 0;\n")
	}

	content.WriteString("    END IF;\n")

	// Convert key to string for JSON access (all map keys are stored as strings in JSON)
	switch keyType {
	case protoreflect.StringKind:
		content.WriteString("    SET element_value = JSON_EXTRACT(map_value, CONCAT('$.', JSON_QUOTE(key_value)));\n")
	default:
		// For numeric keys, convert to string
		content.WriteString("    SET element_value = JSON_EXTRACT(map_value, CONCAT('$.', CAST(key_value AS CHAR)));\n")
	}

	content.WriteString("    IF element_value IS NULL THEN\n")

	// Return appropriate default value for the value type when key not found
	switch valueType {
	case protoreflect.BoolKind:
		content.WriteString("        RETURN FALSE;\n")
	case protoreflect.StringKind, protoreflect.BytesKind:
		content.WriteString("        RETURN '';\n")
	case protoreflect.EnumKind:
		content.WriteString("        RETURN 0;\n")
	case protoreflect.MessageKind:
		content.WriteString("        RETURN JSON_OBJECT();\n")
	default:
		// All numeric types default to 0
		content.WriteString("        RETURN 0;\n")
	}

	content.WriteString("    END IF;\n")

	// Handle type-specific parsing for return value
	switch valueType {
	case protoreflect.BoolKind:
		content.WriteString("    RETURN _pb_json_parse_bool(element_value);\n")
	case protoreflect.StringKind:
		content.WriteString("    RETURN JSON_UNQUOTE(element_value);\n")
	case protoreflect.BytesKind:
		content.WriteString("    RETURN FROM_BASE64(JSON_UNQUOTE(element_value));\n")
	case protoreflect.DoubleKind:
		content.WriteString("    RETURN _pb_json_parse_double(element_value);\n")
	case protoreflect.FloatKind:
		content.WriteString("    RETURN _pb_json_parse_float(element_value);\n")
	default:
		// For integers and enums, direct JSON parsing
		content.WriteString("    RETURN CAST(element_value AS SIGNED);\n")
	}

	content.WriteString("END $$\n\n")

	// Generate contains method for map fields
	containsFuncName := fmt.Sprintf("%s_contains_%s", funcPrefix, fieldName)
	if err := validateFunctionName(containsFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", containsFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, key_value %s) RETURNS BOOLEAN DETERMINISTIC\n", containsFuncName, keySetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE map_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")
	content.WriteString("        RETURN FALSE;\n")
	content.WriteString("    END IF;\n")

	// Convert key to string for JSON access
	switch keyType {
	case protoreflect.StringKind:
		content.WriteString("    RETURN JSON_CONTAINS_PATH(map_value, 'one', CONCAT('$.', JSON_QUOTE(key_value)));\n")
	default:
		// For numeric keys, convert to string
		content.WriteString("    RETURN JSON_CONTAINS_PATH(map_value, 'one', CONCAT('$.', CAST(key_value AS CHAR)));\n")
	}

	content.WriteString("END $$\n\n")

	// Generate put method for map fields
	putFuncName := fmt.Sprintf("%s_put_%s", funcPrefix, fieldName)
	if err := validateFunctionName(putFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", putFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, key_value %s, value_param %s) RETURNS JSON DETERMINISTIC\n", putFuncName, keySetterType, valueSetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE map_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")
	content.WriteString("        SET map_value = JSON_OBJECT();\n")
	content.WriteString("    END IF;\n")

	// Convert key to string and handle value conversion
	switch keyType {
	case protoreflect.StringKind:
		// Handle type-specific conversions for putting values
		switch valueType {
		case protoreflect.BoolKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', JSON_QUOTE(key_value)), CAST((value_param IS TRUE) AS JSON));\n")
		case protoreflect.StringKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', JSON_QUOTE(key_value)), CAST(value_param AS JSON));\n")
		case protoreflect.BytesKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', JSON_QUOTE(key_value)), CAST(TO_BASE64(value_param) AS JSON));\n")
		case protoreflect.DoubleKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', JSON_QUOTE(key_value)), _pb_convert_double_to_number_json(value_param));\n")
		case protoreflect.FloatKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', JSON_QUOTE(key_value)), _pb_convert_float_to_number_json(value_param));\n")
		default:
			// For integers and enums, direct JSON casting
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', JSON_QUOTE(key_value)), CAST(value_param AS JSON));\n")
		}
	default:
		// For numeric keys, convert to string
		switch valueType {
		case protoreflect.BoolKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', CAST(key_value AS CHAR)), CAST((value_param IS TRUE) AS JSON));\n")
		case protoreflect.StringKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', CAST(key_value AS CHAR)), CAST(value_param AS JSON));\n")
		case protoreflect.BytesKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', CAST(key_value AS CHAR)), CAST(TO_BASE64(value_param) AS JSON));\n")
		case protoreflect.DoubleKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', CAST(key_value AS CHAR)), _pb_convert_double_to_number_json(value_param));\n")
		case protoreflect.FloatKind:
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', CAST(key_value AS CHAR)), _pb_convert_float_to_number_json(value_param));\n")
		default:
			// For integers and enums, direct JSON casting
			content.WriteString("    SET map_value = JSON_SET(map_value, CONCAT('$.', CAST(key_value AS CHAR)), CAST(value_param AS JSON));\n")
		}
	}

	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', map_value);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate put_all method for map fields
	putAllFuncName := fmt.Sprintf("%s_put_all_%s", funcPrefix, fieldName)
	if err := validateFunctionName(putAllFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", putAllFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, new_entries JSON) RETURNS JSON DETERMINISTIC\n", putAllFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE map_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")
	content.WriteString("        SET map_value = JSON_OBJECT();\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF new_entries IS NOT NULL THEN\n")
	content.WriteString("        SET map_value = JSON_MERGE_PATCH(map_value, new_entries);\n")
	content.WriteString("    END IF;\n")
	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', map_value);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate remove method for map fields
	removeFuncName := fmt.Sprintf("%s_remove_%s", funcPrefix, fieldName)
	if err := validateFunctionName(removeFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", removeFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, key_value %s) RETURNS JSON DETERMINISTIC\n", removeFuncName, keySetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE map_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")
	content.WriteString("        RETURN proto_data;\n")
	content.WriteString("    END IF;\n")

	// Convert key to string and remove
	switch keyType {
	case protoreflect.StringKind:
		content.WriteString("    SET map_value = JSON_REMOVE(map_value, CONCAT('$.', JSON_QUOTE(key_value)));\n")
	default:
		// For numeric keys, convert to string
		content.WriteString("    SET map_value = JSON_REMOVE(map_value, CONCAT('$.', CAST(key_value AS CHAR)));\n")
	}

	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', map_value);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate count method for map fields
	countFuncName := fmt.Sprintf("%s_count_%s", funcPrefix, fieldName)
	if err := validateFunctionName(countFuncName, fullTypeName); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", countFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS INT DETERMINISTIC\n", countFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE map_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")
	content.WriteString("        RETURN 0;\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    RETURN JSON_LENGTH(map_value);\n")
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
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS INT DETERMINISTIC\n", whichFuncName))
		content.WriteString("BEGIN\n")

		// Check each field in the oneOf group using optimized JSON_EXTRACT
		fields := oneof.Fields()
		for field := range protoreflectutils.Iterate(fields) {
			content.WriteString(fmt.Sprintf("    IF JSON_EXTRACT(proto_data, '$.\"%.d\"') IS NOT NULL THEN\n", field.Number()))
			content.WriteString(fmt.Sprintf("        RETURN %d;\n", field.Number()))
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

// generateEnumNameGetter creates getter functions that return enum names as strings
func generateEnumNameGetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string, isNullableVariant bool, fieldFilterFunc FieldFilterFunc, typePrefixFunc TypePrefixFunc) error {
	// Determine function name suffix
	var suffix string
	if isNullableVariant {
		suffix = "__as_name_or"
	} else {
		suffix = "__as_name"
	}

	// Create function name
	getterFuncName := fmt.Sprintf("%s_get_%s%s", funcPrefix, fieldName, suffix)

	// Check filtering
	var commentPrefix string
	if fieldFilterFunc != nil {
		decision := fieldFilterFunc(field, getterFuncName)
		switch decision {
		case DecisionExclude:
			return nil
		case DecisionCommentOut:
			commentPrefix = "-- "
		case DecisionInclude:
			commentPrefix = ""
		default:
			commentPrefix = ""
		}
	}

	// Validate function name
	if err := validateFunctionName(getterFuncName, fullTypeName); err != nil {
		if commentPrefix == "" {
			return err
		}
	}

	// Get the enum descriptor
	enumDesc := field.Enum()
	if enumDesc == nil {
		return fmt.Errorf("field %s is not an enum field", field.FullName())
	}

	// Get the enum conversion function name using TypePrefixFunc
	enumFullTypeName := enumDesc.FullName()
	enumPackageName := enumDesc.ParentFile().Package()
	enumFuncPrefix := typePrefixFunc(enumPackageName, enumFullTypeName)
	toStringFuncName := enumFuncPrefix + "_to_string"

	// Generate function signature
	content.WriteString(fmt.Sprintf("%sDROP FUNCTION IF EXISTS %s $$\n", commentPrefix, getterFuncName))
	if isNullableVariant {
		content.WriteString(fmt.Sprintf("%sCREATE FUNCTION %s(proto_data JSON, default_value LONGTEXT) RETURNS LONGTEXT DETERMINISTIC\n", commentPrefix, getterFuncName))
	} else {
		content.WriteString(fmt.Sprintf("%sCREATE FUNCTION %s(proto_data JSON) RETURNS LONGTEXT DETERMINISTIC\n", commentPrefix, getterFuncName))
	}
	content.WriteString(fmt.Sprintf("%sBEGIN\n", commentPrefix))

	// DECLARE statements must be at the beginning
	content.WriteString(fmt.Sprintf("%s    DECLARE enum_value INT;\n", commentPrefix))
	content.WriteString(fmt.Sprintf("%s    DECLARE enum_name LONGTEXT;\n", commentPrefix))

	// Get the enum value first
	content.WriteString(fmt.Sprintf("%s    SET enum_value = CAST(JSON_EXTRACT(proto_data, '$.\"%.d\"') AS SIGNED);\n", commentPrefix, field.Number()))

	if isNullableVariant {
		// For nullable variant, check if field is present (enum_value IS NOT NULL)
		content.WriteString(fmt.Sprintf("%s    IF enum_value IS NOT NULL THEN\n", commentPrefix))
	} else {
		// For non-nullable variant, handle missing field by using default (0)
		content.WriteString(fmt.Sprintf("%s    IF enum_value IS NULL THEN\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s        SET enum_value = 0;\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s    END IF;\n", commentPrefix))
	}

	if isNullableVariant {
		// Use the enum's to_string function
		content.WriteString(fmt.Sprintf("%s        SET enum_name = %s(enum_value);\n", commentPrefix, toStringFuncName))

		// If the conversion was successful, return the name; otherwise return the number as string
		content.WriteString(fmt.Sprintf("%s        IF enum_name IS NOT NULL THEN\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s            RETURN enum_name;\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s        ELSE\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s            RETURN CAST(enum_value AS CHAR);\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s        END IF;\n", commentPrefix))

		// Field absent, return default
		content.WriteString(fmt.Sprintf("%s    ELSE\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s        RETURN default_value;\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s    END IF;\n", commentPrefix))
	} else {
		// Use the enum's to_string function
		content.WriteString(fmt.Sprintf("%s    SET enum_name = %s(enum_value);\n", commentPrefix, toStringFuncName))

		// If the conversion was successful, return the name; otherwise return the number as string
		content.WriteString(fmt.Sprintf("%s    IF enum_name IS NOT NULL THEN\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s        RETURN enum_name;\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s    ELSE\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s        RETURN CAST(enum_value AS CHAR);\n", commentPrefix))
		content.WriteString(fmt.Sprintf("%s    END IF;\n", commentPrefix))
	}

	content.WriteString(fmt.Sprintf("%sEND $$\n\n", commentPrefix))

	return nil
}

// generateEnumNameSetter creates setter functions that accept enum names as strings
func generateEnumNameSetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string, fieldFilterFunc FieldFilterFunc, typePrefixFunc TypePrefixFunc) error {
	// Create function name
	setterFuncName := fmt.Sprintf("%s_set_%s__from_name", funcPrefix, fieldName)

	// Check filtering
	var commentPrefix string
	if fieldFilterFunc != nil {
		decision := fieldFilterFunc(field, setterFuncName)
		switch decision {
		case DecisionExclude:
			return nil
		case DecisionCommentOut:
			commentPrefix = "-- "
		case DecisionInclude:
			commentPrefix = ""
		default:
			commentPrefix = ""
		}
	}

	// Validate function name
	if err := validateFunctionName(setterFuncName, fullTypeName); err != nil {
		if commentPrefix == "" {
			return err
		}
	}

	// Get the enum descriptor
	enumDesc := field.Enum()
	if enumDesc == nil {
		return fmt.Errorf("field %s is not an enum field", field.FullName())
	}

	// Get the enum conversion function name using TypePrefixFunc
	enumFullTypeName := enumDesc.FullName()
	enumPackageName := enumDesc.ParentFile().Package()
	enumFuncPrefix := typePrefixFunc(enumPackageName, enumFullTypeName)
	fromStringFuncName := enumFuncPrefix + "_from_string"

	// Generate function signature
	content.WriteString(fmt.Sprintf("%sDROP FUNCTION IF EXISTS %s $$\n", commentPrefix, setterFuncName))
	content.WriteString(fmt.Sprintf("%sCREATE FUNCTION %s(proto_data JSON, enum_name LONGTEXT) RETURNS JSON DETERMINISTIC\n", commentPrefix, setterFuncName))
	content.WriteString(fmt.Sprintf("%sBEGIN\n", commentPrefix))

	// Declare variable to store enum value
	content.WriteString(fmt.Sprintf("%s    DECLARE enum_value INT;\n", commentPrefix))

	// Call the enum's from_string function
	content.WriteString(fmt.Sprintf("%s    SET enum_value = %s(enum_name);\n", commentPrefix, fromStringFuncName))

	// Check if the conversion was successful (not NULL)
	content.WriteString(fmt.Sprintf("%s    IF enum_value IS NULL THEN\n", commentPrefix))
	content.WriteString(fmt.Sprintf("%s        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid enum name';\n", commentPrefix))
	content.WriteString(fmt.Sprintf("%s    END IF;\n", commentPrefix))

	// Handle zero values based on field presence semantics
	content.WriteString(fmt.Sprintf("%s    IF enum_value = 0 THEN\n", commentPrefix))
	if !field.HasPresence() {
		// For proto3, zero values should be omitted (removed from JSON)
		content.WriteString(fmt.Sprintf("%s        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", commentPrefix, field.Number()))
	} else {
		// For proto2 or proto3 optional, zero values are still stored
		content.WriteString(fmt.Sprintf("%s        RETURN JSON_SET(proto_data, '$.\"%.d\"', 0);\n", commentPrefix, field.Number()))
	}
	content.WriteString(fmt.Sprintf("%s    ELSE\n", commentPrefix))
	content.WriteString(fmt.Sprintf("%s        RETURN JSON_SET(proto_data, '$.\"%.d\"', enum_value);\n", commentPrefix, field.Number()))
	content.WriteString(fmt.Sprintf("%s    END IF;\n", commentPrefix))

	content.WriteString(fmt.Sprintf("%sEND $$\n\n", commentPrefix))

	return nil
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
