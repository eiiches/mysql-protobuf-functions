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
	// Skip map entry messages - these are synthetic types created by protobuf for map fields
	// and should not have user-facing accessor methods generated
	if messageDesc.IsMapEntry() {
		return nil
	}

	// Use FullName from descriptor - no manual string construction needed
	fullTypeName := messageDesc.FullName()
	packageName := messageDesc.ParentFile().Package()

	// Get prefix for this specific type
	funcPrefix := typePrefixFunc(packageName, fullTypeName)

	// Generate basic constructor
	if err := generateConstructor(content, funcPrefix, messageDesc); err != nil {
		return err
	}

	// Generate conversion methods (from_json, to_json, etc.)
	if err := generateConversionMethods(content, funcPrefix, messageDesc, schemaFunctionName); err != nil {
		return err
	}

	// Generate field accessor methods with enhanced opaque API patterns
	if err := generateFieldAccessorMethods(content, messageDesc, funcPrefix, fieldFilterFunc, typePrefixFunc); err != nil {
		return err
	}

	// Generate oneOf methods for oneOf groups
	if err := generateOneOfMethods(content, messageDesc, funcPrefix); err != nil {
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
func generateConstructor(content *strings.Builder, funcPrefix string, messageDesc protoreflect.MessageDescriptor) error {
	newFuncName := funcPrefix + "_new"
	if err := validateFunctionName(newFuncName, messageDesc.FullName()); err != nil {
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
func generateConversionMethods(content *strings.Builder, funcPrefix string, messageDesc protoreflect.MessageDescriptor, schemaFunctionName string) error {
	// Generate from_json
	fromJsonFuncName := funcPrefix + "_from_json"
	if err := validateFunctionName(fromJsonFuncName, messageDesc.FullName()); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", fromJsonFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(json_data JSON, json_unmarshal_options JSON) RETURNS JSON DETERMINISTIC\n", fromJsonFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_json_to_number_json(%s(), '.%s', json_data, json_unmarshal_options);\n", schemaFunctionName, messageDesc.FullName()))
	content.WriteString("END $$\n\n")

	// Generate from_message
	fromMessageFuncName := funcPrefix + "_from_message"
	if err := validateFunctionName(fromMessageFuncName, messageDesc.FullName()); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", fromMessageFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(message_data LONGBLOB, unmarshal_options JSON) RETURNS JSON DETERMINISTIC\n", fromMessageFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_message_to_number_json(%s(), '.%s', message_data, unmarshal_options);\n", schemaFunctionName, messageDesc.FullName()))
	content.WriteString("END $$\n\n")

	// Generate to_json
	toJsonFuncName := funcPrefix + "_to_json"
	if err := validateFunctionName(toJsonFuncName, messageDesc.FullName()); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", toJsonFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, json_marshal_options JSON) RETURNS JSON DETERMINISTIC\n", toJsonFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_number_json_to_json(%s(), '.%s', proto_data, json_marshal_options);\n", schemaFunctionName, messageDesc.FullName()))
	content.WriteString("END $$\n\n")

	// Generate to_message
	toMessageFuncName := funcPrefix + "_to_message"
	if err := validateFunctionName(toMessageFuncName, messageDesc.FullName()); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", toMessageFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, marshal_options JSON) RETURNS LONGBLOB DETERMINISTIC\n", toMessageFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_number_json_to_message(%s(), '.%s', proto_data, marshal_options);\n", schemaFunctionName, messageDesc.FullName()))
	content.WriteString("END $$\n\n")

	return nil
}

// generateFieldAccessorMethods creates enhanced field accessor methods following opaque API patterns
func generateFieldAccessorMethods(content *strings.Builder, messageDesc protoreflect.MessageDescriptor, funcPrefix string, fieldFilterFunc FieldFilterFunc, typePrefixFunc TypePrefixFunc) error {
	fullTypeName := messageDesc.FullName()
	fields := messageDesc.Fields()

	for field := range protoreflectutils.Iterate(fields) {
		fieldName := string(field.Name())

		// Generate enhanced getter with better defaults and type safety
		if err := generateEnhancedGetter(content, funcPrefix, field); err != nil {
			return err
		}

		// Generate nullable getter with custom default for optional fields
		if field.HasPresence() && !field.IsList() && !field.IsMap() {
			if err := generateNullableGetter(content, funcPrefix, field, fieldFilterFunc); err != nil {
				return err
			}
		}

		// Generate enum name getters for enum fields
		if field.Kind() == protoreflect.EnumKind && !field.IsList() && !field.IsMap() {
			// Generate regular enum name getter (__as_name)
			if err := generateEnumNameGetter(content, funcPrefix, field, false, fieldFilterFunc, typePrefixFunc); err != nil {
				return err
			}
			// Generate nullable enum name getter (__as_name_or) for optional fields
			if field.HasPresence() {
				if err := generateEnumNameGetter(content, funcPrefix, field, true, fieldFilterFunc, typePrefixFunc); err != nil {
					return err
				}
			}
		}

		// Generate enum name setters for enum fields
		if field.Kind() == protoreflect.EnumKind && !field.IsList() && !field.IsMap() {
			if err := generateEnumNameSetter(content, funcPrefix, fullTypeName, field, fieldName, fieldFilterFunc, typePrefixFunc); err != nil {
				return err
			}
		}

		// Generate enhanced setter with validation
		if err := generateEnhancedSetter(content, funcPrefix, fullTypeName, field, fieldName); err != nil {
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
		if field.IsList() && !field.IsMap() {
			if err := generateRepeatedFieldMethods(content, funcPrefix, fullTypeName, field, fieldName); err != nil {
				return err
			}
		}

		// Generate additional methods for map fields
		if field.IsMap() {
			if err := generateMapFieldMethods(content, funcPrefix, fullTypeName, field, fieldName); err != nil {
				return err
			}

			// Generate individual key getter with default (__or) for map fields
			if err := generateMapKeyNullableGetter(content, funcPrefix, fullTypeName, field, fieldName); err != nil {
				return err
			}
		}
	}

	return nil
}

// generateEnhancedGetter creates an enhanced getter with better defaults and type safety
func generateEnhancedGetter(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor) error {
	fieldName := string(field.Name())
	// Use new naming: get_all_ for repeated/map fields, get_ for singular fields
	var getterFuncName string
	if field.IsList() || field.IsMap() {
		getterFuncName = fmt.Sprintf("%s_get_all_%s", funcPrefix, fieldName)
	} else {
		getterFuncName = fmt.Sprintf("%s_get_%s", funcPrefix, fieldName)
	}
	if err := validateFunctionName(getterFuncName, field.ContainingMessage().FullName()); err != nil {
		return err
	}

	var getterType string
	if field.IsList() || field.IsMap() {
		getterType = "JSON"
	} else {
		getterType = GetProtobufType(field.Kind()).GetSqlTypeName()
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS %s DETERMINISTIC\n", getterFuncName, getterType))
	content.WriteString("BEGIN\n")

	// Declare variables at the beginning for fields that need JSON parsing
	needsJsonParsing := !field.IsList() && field.Message() == nil && (field.Kind() == protoreflect.BoolKind ||
		field.Kind() == protoreflect.StringKind ||
		field.Kind() == protoreflect.BytesKind ||
		field.Kind() == protoreflect.Int64Kind ||
		field.Kind() == protoreflect.Sint64Kind ||
		field.Kind() == protoreflect.Sfixed64Kind ||
		field.Kind() == protoreflect.Int32Kind ||
		field.Kind() == protoreflect.Sint32Kind ||
		field.Kind() == protoreflect.Sfixed32Kind ||
		field.Kind() == protoreflect.Uint64Kind ||
		field.Kind() == protoreflect.Fixed64Kind ||
		field.Kind() == protoreflect.Uint32Kind ||
		field.Kind() == protoreflect.Fixed32Kind ||
		field.Kind() == protoreflect.FloatKind ||
		field.Kind() == protoreflect.DoubleKind ||
		field.Kind() == protoreflect.EnumKind)

	if needsJsonParsing {
		content.WriteString("    DECLARE json_value JSON;\n")
	}

	switch {
	case field.IsMap():
		// For map fields, return JSON object or empty array if not present
		content.WriteString(fmt.Sprintf("    RETURN COALESCE(JSON_EXTRACT(proto_data, '$.\"%.d\"'), JSON_ARRAY());\n", field.Number()))
	case field.IsList():
		// For repeated fields, special handling for float and double to convert from binary format
		if field.Kind() == protoreflect.FloatKind || field.Kind() == protoreflect.DoubleKind {
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

			// Use unified conversion system for NumberJSON → SQL Type
			convertedExpression := GetProtobufType(field.Kind()).GenerateNumberJsonToSqlExpression("element_value")
			content.WriteString(fmt.Sprintf("        SET result_array = JSON_ARRAY_APPEND(result_array, '$', %s);\n", convertedExpression))

			content.WriteString("        SET i = i + 1;\n")
			content.WriteString("    END WHILE;\n")
			content.WriteString("\n")
			content.WriteString("    RETURN result_array;\n")
		} else {
			// For other repeated fields, return JSON array or empty array if not present
			content.WriteString(fmt.Sprintf("    RETURN COALESCE(JSON_EXTRACT(proto_data, '$.\"%.d\"'), JSON_ARRAY());\n", field.Number()))
		}
	case field.Message() != nil:
		// For message fields, return JSON object or empty object if not present
		content.WriteString(fmt.Sprintf("    RETURN COALESCE(JSON_EXTRACT(proto_data, '$.\"%.d\"'), JSON_OBJECT());\n", field.Number()))
	default:
		// For scalar fields, use unified type system for NumberJSON → SQL conversion
		defaultValue := GetDefaultValue(field)
		content.WriteString(fmt.Sprintf("    SET json_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
		content.WriteString("    IF json_value IS NULL THEN\n")
		content.WriteString(fmt.Sprintf("        RETURN %s;\n", defaultValue))
		content.WriteString("    END IF;\n")
		expression := GetProtobufType(field.Kind()).GenerateNumberJsonToSqlExpression("json_value")
		content.WriteString(fmt.Sprintf("    RETURN %s;\n", expression))
	}

	content.WriteString("END $$\n\n")
	return nil
}

// ValueConverter is a function that converts an extracted JSON value to the final return type
type ValueConverter func(valueVar string) string

// generateUnifiedFieldGetter creates regular field getter variants (__or, __as_name, __as_name_or) - NOT for maps
func generateUnifiedFieldGetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string, suffix string, returnType string, hasDefaultParam bool, defaultValue string, converter ValueConverter, fieldFilterFunc FieldFilterFunc) error {
	getterFuncName := fmt.Sprintf("%s_get_%s%s", funcPrefix, fieldName, suffix)

	// Build parameters
	var parameters string
	if hasDefaultParam {
		parameters = fmt.Sprintf("proto_data JSON, default_value %s", returnType)
	} else {
		parameters = "proto_data JSON"
	}

	// Apply field filtering
	decision := DecisionInclude
	if fieldFilterFunc != nil {
		decision = fieldFilterFunc(field, getterFuncName)
	}

	if decision == DecisionExclude {
		return nil
	}

	if decision == DecisionInclude {
		if err := validateFunctionName(getterFuncName, fullTypeName); err != nil {
			return err
		}
	}

	commentPrefix := ""
	if decision == DecisionCommentOut {
		commentPrefix = "-- "
		content.WriteString(fmt.Sprintf("-- SKIPPED: Function '%s' was filtered out\n", getterFuncName))
	}

	// Generate function
	content.WriteString(fmt.Sprintf("%sDROP FUNCTION IF EXISTS %s $$\n", commentPrefix, getterFuncName))
	content.WriteString(fmt.Sprintf("%sCREATE FUNCTION %s(%s) RETURNS %s DETERMINISTIC\n", commentPrefix, getterFuncName, parameters, returnType))
	content.WriteString(fmt.Sprintf("%sBEGIN\n", commentPrefix))

	// Declare variables
	content.WriteString(fmt.Sprintf("%s    DECLARE json_value JSON;\n", commentPrefix))
	content.WriteString(fmt.Sprintf("%s    SET json_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", commentPrefix, field.Number()))

	// Always handle null the same way - just different default values
	content.WriteString(fmt.Sprintf("%s    IF json_value IS NULL THEN\n", commentPrefix))
	content.WriteString(fmt.Sprintf("%s        RETURN %s;\n", commentPrefix, defaultValue))
	content.WriteString(fmt.Sprintf("%s    END IF;\n", commentPrefix))

	// Convert and return the value
	convertedExpr := converter("json_value")
	content.WriteString(fmt.Sprintf("%s    RETURN %s;\n", commentPrefix, convertedExpr))

	content.WriteString(fmt.Sprintf("%sEND $$\n\n", commentPrefix))
	return nil
}

// Now the original functions just call the unified one
func generateNullableGetter(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor, fieldFilterFunc FieldFilterFunc) error {
	returnType := GetProtobufType(field.Kind()).GetSqlTypeName()
	converter := func(valueVar string) string {
		return GetProtobufType(field.Kind()).GenerateNumberJsonToSqlExpression(valueVar)
	}
	return generateUnifiedFieldGetter(content, funcPrefix, field.ContainingMessage().FullName(), field, string(field.Name()), "__or", returnType, true, "default_value", converter, fieldFilterFunc)
}

// generateEnhancedSetter creates an enhanced setter with input validation
func generateEnhancedSetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string) error {
	// Use new naming: set_all_ for repeated/map fields, set_ for singular fields
	var setterFuncName string
	if field.IsList() || field.IsMap() {
		setterFuncName = fmt.Sprintf("%s_set_all_%s", funcPrefix, fieldName)
	} else {
		setterFuncName = fmt.Sprintf("%s_set_%s", funcPrefix, fieldName)
	}
	if err := validateFunctionName(setterFuncName, fullTypeName); err != nil {
		return err
	}

	var setterType string
	if field.IsList() || field.IsMap() {
		setterType = "JSON"
	} else {
		setterType = GetProtobufType(field.Kind()).GetSqlTypeName()
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", setterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, field_value %s) RETURNS JSON DETERMINISTIC\n", setterFuncName, setterType))
	content.WriteString("BEGIN\n")

	// Declare variables at the beginning (required by MySQL syntax)
	oneof := field.ContainingOneof()
	if oneof != nil {
		content.WriteString("    DECLARE temp_data JSON DEFAULT proto_data;\n")
	}

	// Add input validation
	if field.IsMap() {
		// For map fields, use the same conversion logic as put_all
		protobufType := GetProtobufType(field.MapValue().Kind())
		jsonConversionExpr := protobufType.GenerateJsonToNumberJsonExpression("current_value")
		needsConversion := jsonConversionExpr != "current_value"

		if needsConversion {
			content.WriteString("    DECLARE json_keys JSON;\n")
			content.WriteString("    DECLARE key_count INT;\n")
			content.WriteString("    DECLARE i INT;\n")
			content.WriteString("    DECLARE current_key VARCHAR(255);\n")
			content.WriteString("    DECLARE current_value JSON;\n")
			content.WriteString("    DECLARE converted_value JSON;\n")
			content.WriteString("    DECLARE result_map JSON DEFAULT JSON_OBJECT();\n")
		}

		content.WriteString("    IF field_value IS NULL THEN\n")
		content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
		content.WriteString("    END IF;\n")
		content.WriteString("    IF JSON_TYPE(field_value) != 'OBJECT' THEN\n")
		content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Map field value must be a JSON object';\n")
		content.WriteString("    END IF;\n")

		if needsConversion {
			// Need to convert each JSON value to number_json format (for float, double)
			content.WriteString("    SET json_keys = JSON_KEYS(field_value);\n")
			content.WriteString("    IF json_keys IS NOT NULL THEN\n")
			content.WriteString("        SET key_count = JSON_LENGTH(json_keys);\n")
			content.WriteString("        SET i = 0;\n")
			content.WriteString("        \n")
			content.WriteString("        WHILE i < key_count DO\n")
			content.WriteString("            SET current_key = JSON_UNQUOTE(JSON_EXTRACT(json_keys, CONCAT('$[', i, ']')));\n")
			content.WriteString("            SET current_value = JSON_EXTRACT(field_value, CONCAT('$.\"', current_key, '\"'));\n")
			content.WriteString(fmt.Sprintf("            SET converted_value = %s;\n", jsonConversionExpr))
			content.WriteString("            SET result_map = JSON_SET(result_map, CONCAT('$.\"', current_key, '\"'), converted_value);\n")
			content.WriteString("            SET i = i + 1;\n")
			content.WriteString("        END WHILE;\n")
			content.WriteString("    END IF;\n")
			content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', result_map);\n", field.Number()))
		} else {
			// For types that don't need conversion, use direct assignment
			content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
		}
		content.WriteString("END $$\n\n")
		return nil
	} else if field.IsList() {
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
		// Use unified JSON → NumberJSON conversion system
		convertedExpression := GetProtobufType(field.Kind()).GenerateJsonToNumberJsonExpression("element_value")
		content.WriteString(fmt.Sprintf("        SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', %s);\n", convertedExpression))

		content.WriteString("        SET i = i + 1;\n")
		content.WriteString("    END WHILE;\n")
		content.WriteString("\n")
		content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', converted_array);\n", field.Number()))
		content.WriteString("END $$\n\n")
		return nil
	} else if field.Message() != nil {
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

		// Set the new field value with proper JSON conversion - use unified type system
		jsonExpression := GetProtobufType(field.Kind()).GenerateSqlToNumberJsonExpression("field_value")
		content.WriteString(fmt.Sprintf("    RETURN JSON_SET(temp_data, '$.\"%.d\"', %s);\n", field.Number(), jsonExpression))
	} else {
		// Check if this is a proto3 field without presence and being set to default value
		// According to protonumberjson spec, such fields should be omitted
		if !field.HasPresence() && !field.IsList() && field.Message() == nil {
			// For proto3 fields without presence, omit default values - use unified type system
			content.WriteString("    -- Proto3 field without presence: omit default values per protonumberjson spec\n")
			setterLogic := GetProtobufType(field.Kind()).GenerateSetterWithZeroValueRemoval(int32(field.Number()), "field_value")
			content.WriteString("    " + setterLogic + "\n")
		} else {
			// Fields with presence: always set the value - use unified type system
			jsonExpression := GetProtobufType(field.Kind()).GenerateSqlToNumberJsonExpression("field_value")
			content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', %s);\n", field.Number(), jsonExpression))
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
func generateRepeatedFieldMethods(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string) error {
	// Generate add method for repeated fields
	addFuncName := fmt.Sprintf("%s_add_%s", funcPrefix, fieldName)
	if err := validateFunctionName(addFuncName, fullTypeName); err != nil {
		return err
	}

	elementSetterType := GetProtobufType(field.Kind()).GetSqlTypeName() // Get non-repeated type for element

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", addFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, element_value %s) RETURNS JSON DETERMINISTIC\n", addFuncName, elementSetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE current_array JSON;\n")
	content.WriteString(fmt.Sprintf("    SET current_array = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF current_array IS NULL THEN\n")
	content.WriteString("        SET current_array = JSON_ARRAY();\n")
	content.WriteString("    END IF;\n")

	// Handle type-specific conversions for proper protonumberjson format in arrays - use unified type system
	numberJsonExpression := GetProtobufType(field.Kind()).GenerateSqlToNumberJsonExpression("element_value")
	content.WriteString(fmt.Sprintf("    SET current_array = JSON_ARRAY_APPEND(current_array, '$', %s);\n", numberJsonExpression))

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

	elementGetterType := GetProtobufType(field.Kind()).GetSqlTypeName() // Get non-repeated type for element

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getIndexFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, index_value INT) RETURNS %s DETERMINISTIC\n", getIndexFuncName, elementGetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE array_value JSON;\n")
	content.WriteString("    DECLARE element_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF array_value IS NULL THEN\n")
	content.WriteString("        SET array_value = JSON_ARRAY(); -- Treat missing repeated field as empty array\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF index_value < 0 OR index_value >= JSON_LENGTH(array_value) THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Array index out of bounds';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    SET element_value = JSON_EXTRACT(array_value, CONCAT('$[', index_value, ']'));\n")

	// Use unified array element return logic
	expression := GetProtobufType(field.Kind()).GenerateNumberJsonToSqlExpression("element_value")
	content.WriteString(fmt.Sprintf("    RETURN %s;\n", expression))

	content.WriteString("END $$\n\n")

	// Generate index-based set method for repeated fields
	setIndexFuncName := fmt.Sprintf("%s_set_%s", funcPrefix, fieldName)
	if err := validateFunctionName(setIndexFuncName, fullTypeName); err != nil {
		return err
	}

	elementSetterType = GetProtobufType(field.Kind()).GetSqlTypeName() // Get non-repeated type for element

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", setIndexFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, index_value INT, element_value %s) RETURNS JSON DETERMINISTIC\n", setIndexFuncName, elementSetterType))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE array_value JSON;\n")
	content.WriteString("    DECLARE array_length INT;\n")
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF array_value IS NULL THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Array index out of bounds';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    SET array_length = JSON_LENGTH(array_value);\n")
	content.WriteString("    IF index_value < 0 OR index_value >= array_length THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Array index out of bounds';\n")
	content.WriteString("    END IF;\n")

	// Use unified array set logic
	numberJsonExpression = GetProtobufType(field.Kind()).GenerateSqlToNumberJsonExpression("element_value")
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_SET(array_value, CONCAT('$[', index_value, ']'), %s);\n", numberJsonExpression))

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
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF array_value IS NULL THEN\n")
	content.WriteString("        SET array_value = JSON_ARRAY();\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    SET array_length = JSON_LENGTH(array_value);\n")
	content.WriteString("    IF index_value < 0 OR index_value > array_length THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insert index out of bounds';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    \n")

	// Handle type-specific conversions for inserting
	// Use unified array insert logic with JSON_ARRAY_INSERT
	numberJsonExpression = GetProtobufType(field.Kind()).GenerateSqlToNumberJsonExpression("element_value")
	content.WriteString(fmt.Sprintf("    SET array_value = JSON_ARRAY_INSERT(array_value, CONCAT('$[', index_value, ']'), %s);\n", numberJsonExpression))

	content.WriteString("    \n")
	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', array_value);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate remove method for repeated fields
	removeFuncName := fmt.Sprintf("%s_remove_%s", funcPrefix, fieldName)
	if err := validateFunctionName(removeFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", removeFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, index_value INT) RETURNS JSON DETERMINISTIC\n", removeFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE array_length INT;\n")
	content.WriteString(fmt.Sprintf("    SET array_length = JSON_LENGTH(JSON_EXTRACT(proto_data, '$.\"%.d\"'));\n", field.Number()))
	content.WriteString("    IF array_length IS NULL THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Array index out of bounds';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF index_value < 0 OR index_value >= array_length THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Array index out of bounds';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    \n")
	content.WriteString("    IF array_length = 1 THEN\n")
	content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    ELSE\n")
	content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, CONCAT('$.\"%.d\"[', index_value, ']'));\n", field.Number()))
	content.WriteString("    END IF;\n")
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

	// Handle type-specific conversions for adding elements - use unified type system
	// This is JSON→NumberJSON conversion since we expect element_value to be JSON format
	numberJsonExpression = GetProtobufType(field.Kind()).GenerateJsonToNumberJsonExpression("element_value")
	content.WriteString(fmt.Sprintf("        SET current_array = JSON_ARRAY_APPEND(current_array, '$', %s);\n", numberJsonExpression))

	content.WriteString("        SET i = i + 1;\n")
	content.WriteString("    END WHILE;\n")
	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', current_array);\n", field.Number()))
	content.WriteString("END $$\n\n")

	return nil
}

// generateMapFieldMethods creates additional methods for map fields
func generateMapFieldMethods(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string) error {
	// Get MySQL types for key and value
	keySetterType := GetProtobufType(field.MapKey().Kind()).GetSqlTypeName()
	valueSetterType := GetProtobufType(field.MapValue().Kind()).GetSqlTypeName()
	valueGetterType := GetProtobufType(field.MapValue().Kind()).GetSqlTypeName()

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
	defaultValue := GetProtobufType(field.MapValue().Kind()).GetDefaultValueExpression()
	content.WriteString(fmt.Sprintf("        RETURN %s;\n", defaultValue))

	content.WriteString("    END IF;\n")

	// Convert key to string for JSON access (all map keys are stored as strings in JSON)
	if field.MapKey().Kind() == protoreflect.StringKind {
		content.WriteString("    SET element_value = JSON_EXTRACT(map_value, CONCAT('$.', JSON_QUOTE(key_value)));\n")
	} else {
		// For numeric keys, convert to string using unified type system
		keyPathExpression := GetProtobufType(field.MapKey().Kind()).GenerateSqlToMapKeyExpression("key_value")
		content.WriteString(fmt.Sprintf("    SET element_value = JSON_EXTRACT(map_value, CONCAT('$.\"', %s, '\"'));\n", keyPathExpression))
	}

	content.WriteString("    IF element_value IS NULL THEN\n")

	// Return appropriate default value for the value type when key not found
	content.WriteString(fmt.Sprintf("        RETURN %s;\n", defaultValue))

	content.WriteString("    END IF;\n")

	// Handle type-specific parsing for return value
	conversionExpression := GetProtobufType(field.MapValue().Kind()).GenerateNumberJsonToSqlExpression("element_value")
	content.WriteString(fmt.Sprintf("    RETURN %s;\n", conversionExpression))

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
	if field.MapKey().Kind() == protoreflect.StringKind {
		content.WriteString("    RETURN JSON_CONTAINS_PATH(map_value, 'one', CONCAT('$.', JSON_QUOTE(key_value)));\n")
	} else {
		// For numeric keys, convert to string using unified type system
		keyPathExpression := GetProtobufType(field.MapKey().Kind()).GenerateSqlToMapKeyExpression("key_value")
		content.WriteString(fmt.Sprintf("    RETURN JSON_CONTAINS_PATH(map_value, 'one', CONCAT('$.\"', %s, '\"'));\n", keyPathExpression))
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
	keyPathExpression := GetProtobufType(field.MapKey().Kind()).GenerateSqlToMapKeyExpression("key_value")
	valueConversionExpression := GetProtobufType(field.MapValue().Kind()).GenerateSqlToNumberJsonExpression("value_param")

	// Generate map key path based on key type
	var keyPath string
	if field.MapKey().Kind() == protoreflect.StringKind {
		keyPath = "CONCAT('$.', JSON_QUOTE(key_value))"
	} else {
		// For numeric keys, quote them as JSON object keys
		keyPath = fmt.Sprintf("CONCAT('$.\"', %s, '\"')", keyPathExpression)
	}

	content.WriteString(fmt.Sprintf("    SET map_value = JSON_SET(map_value, %s, %s);\n", keyPath, valueConversionExpression))

	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', map_value);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate put_all method for map fields
	putAllFuncName := fmt.Sprintf("%s_put_all_%s", funcPrefix, fieldName)
	if err := validateFunctionName(putAllFuncName, fullTypeName); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", putAllFuncName))
	// For primitive value types that need conversion (float, double), we need to convert each value
	protobufType := GetProtobufType(field.MapValue().Kind())
	jsonConversionExpr := protobufType.GenerateJsonToNumberJsonExpression("current_value")
	needsConversion := jsonConversionExpr != "current_value"

	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, new_entries JSON) RETURNS JSON DETERMINISTIC\n", putAllFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE map_value JSON;\n")
	if needsConversion {
		content.WriteString("    DECLARE json_keys JSON;\n")
		content.WriteString("    DECLARE key_count INT;\n")
		content.WriteString("    DECLARE i INT;\n")
		content.WriteString("    DECLARE current_key VARCHAR(255);\n")
		content.WriteString("    DECLARE current_value JSON;\n")
		content.WriteString("    DECLARE converted_value JSON;\n")
	}
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")
	content.WriteString("        SET map_value = JSON_OBJECT();\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF new_entries IS NOT NULL THEN\n")

	if needsConversion {
		// Need to convert each JSON value to number_json format (for float, double)
		content.WriteString("        SET json_keys = JSON_KEYS(new_entries);\n")
		content.WriteString("        SET key_count = JSON_LENGTH(json_keys);\n")
		content.WriteString("        SET i = 0;\n")
		content.WriteString("        \n")
		content.WriteString("        WHILE i < key_count DO\n")
		content.WriteString("            SET current_key = JSON_UNQUOTE(JSON_EXTRACT(json_keys, CONCAT('$[', i, ']')));\n")
		content.WriteString("            SET current_value = JSON_EXTRACT(new_entries, CONCAT('$.\"', current_key, '\"'));\n")
		content.WriteString(fmt.Sprintf("            SET converted_value = %s;\n", jsonConversionExpr))
		content.WriteString("            SET map_value = JSON_SET(map_value, CONCAT('$.\"', current_key, '\"'), converted_value);\n")
		content.WriteString("            SET i = i + 1;\n")
		content.WriteString("        END WHILE;\n")
	} else {
		// For types that don't need conversion (int32, string, etc.), use direct merge
		content.WriteString("        SET map_value = JSON_MERGE_PATCH(map_value, new_entries);\n")
	}

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
	if field.MapKey().Kind() == protoreflect.StringKind {
		content.WriteString("    SET map_value = JSON_REMOVE(map_value, CONCAT('$.', JSON_QUOTE(key_value)));\n")
	} else {
		// For numeric keys, convert to string using unified type system
		keyPathExpression := GetProtobufType(field.MapKey().Kind()).GenerateSqlToMapKeyExpression("key_value")
		content.WriteString(fmt.Sprintf("    SET map_value = JSON_REMOVE(map_value, CONCAT('$.\"', %s, '\"'));\n", keyPathExpression))
	}

	// If map becomes empty, remove the field entirely (proto3 default value omission)
	content.WriteString("    IF JSON_LENGTH(map_value) = 0 THEN\n")
	content.WriteString(fmt.Sprintf("        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    END IF;\n")
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

// generateMapKeyNullableGetter creates a nullable getter for individual map key access
func generateMapKeyNullableGetter(content *strings.Builder, funcPrefix string, fullTypeName protoreflect.FullName, field protoreflect.FieldDescriptor, fieldName string) error {
	// Get map key and value types
	isValueMessage := field.MapValue().Kind() == protoreflect.MessageKind

	keyMySQLType := GetProtobufType(field.MapKey().Kind()).GetSqlTypeName()
	valueMySQLType := GetProtobufType(field.MapValue().Kind()).GetSqlTypeName()

	// For float/double types, use LONGTEXT to accept binary format strings
	if field.MapValue().Kind() == protoreflect.FloatKind || field.MapValue().Kind() == protoreflect.DoubleKind {
		valueMySQLType = "LONGTEXT"
	}

	// Create function name with __or modifier for individual key access
	getterFuncName := fmt.Sprintf("%s_get_%s__or", funcPrefix, fieldName)

	if err := validateFunctionName(getterFuncName, fullTypeName); err != nil {
		return err
	}

	// Create the function signature
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, key_value %s, default_value %s) RETURNS %s DETERMINISTIC\n",
		getterFuncName, keyMySQLType, valueMySQLType, valueMySQLType))
	content.WriteString("BEGIN\n")

	// Declare variables
	content.WriteString("    DECLARE map_value JSON;\n")
	content.WriteString("    DECLARE element_value JSON;\n")

	// Get the map from the proto_data
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")
	content.WriteString("        RETURN default_value;\n")
	content.WriteString("    END IF;\n")

	// Get the element from the map using the key
	if field.MapKey().Kind() == protoreflect.StringKind {
		content.WriteString("    SET element_value = JSON_EXTRACT(map_value, CONCAT('$.', JSON_QUOTE(key_value)));\n")
	} else {
		// For non-string keys, use unified type system for key path generation
		keyPathExpression := GetProtobufType(field.MapKey().Kind()).GenerateSqlToMapKeyExpression("key_value")
		content.WriteString(fmt.Sprintf("    SET element_value = JSON_EXTRACT(map_value, CONCAT('$.\"', %s, '\"'));\n", keyPathExpression))
	}

	content.WriteString("    IF element_value IS NULL THEN\n")
	content.WriteString("        RETURN default_value;\n")
	content.WriteString("    END IF;\n")

	// Return the value, converting if necessary
	if isValueMessage {
		content.WriteString("    RETURN element_value;\n")
	} else {
		// For primitive types, use unified type conversion
		conversionExpression := GetProtobufType(field.MapValue().Kind()).GenerateNumberJsonToSqlExpression("element_value")
		content.WriteString(fmt.Sprintf("    RETURN %s;\n", conversionExpression))
	}

	content.WriteString("END $$\n\n")

	return nil
}

// generateOneOfMethods creates methods for oneOf groups
func generateOneOfMethods(content *strings.Builder, messageDesc protoreflect.MessageDescriptor, funcPrefix string) error {
	fullTypeName := messageDesc.FullName()
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

// generateEnumNameGetter creates getter functions that return enum names as strings
func generateEnumNameGetter(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor, isNullableVariant bool, fieldFilterFunc FieldFilterFunc, typePrefixFunc TypePrefixFunc) error {
	// Get enum conversion function name
	enumDesc := field.Enum()
	enumFullTypeName := enumDesc.FullName()
	enumPackageName := enumDesc.ParentFile().Package()
	enumFuncPrefix := typePrefixFunc(enumPackageName, enumFullTypeName)
	toStringFuncName := enumFuncPrefix + "_to_string"

	// Create converter that handles enum to string conversion
	converter := func(valueVar string) string {
		return fmt.Sprintf(`COALESCE(%s(CAST(%s AS SIGNED)), CAST(CAST(%s AS SIGNED) AS CHAR))`, toStringFuncName, valueVar, valueVar)
	}

	if isNullableVariant {
		return generateUnifiedFieldGetter(content, funcPrefix, field.ContainingMessage().FullName(), field, string(field.Name()), "__as_name_or", "LONGTEXT", true, "default_value", converter, fieldFilterFunc)
	} else {
		// For non-nullable enum, the default is the converted value of 0
		defaultValue := fmt.Sprintf(`COALESCE(%s(0), '0')`, toStringFuncName)
		return generateUnifiedFieldGetter(content, funcPrefix, field.ContainingMessage().FullName(), field, string(field.Name()), "__as_name", "LONGTEXT", false, defaultValue, converter, fieldFilterFunc)
	}
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
	content.WriteString(fmt.Sprintf("%s    END IF;\n", commentPrefix))
	content.WriteString(fmt.Sprintf("%s    RETURN JSON_SET(proto_data, '$.\"%.d\"', enum_value);\n", commentPrefix, field.Number()))

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
