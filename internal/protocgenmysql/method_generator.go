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
		if err := generateEnumMethods(&content, enumDesc, typePrefixFunc); err != nil {
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
		if err := generateEnumMethods(content, nestedEnumDesc, typePrefixFunc); err != nil {
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
	fields := messageDesc.Fields()

	for field := range protoreflectutils.Iterate(fields) {
		// Generate enhanced getter with better defaults and type safety
		if field.IsMap() {
			if err := generateMapGetAll(content, funcPrefix, field); err != nil {
				return err
			}
			if err := generateMapSetAll(content, funcPrefix, field); err != nil {
				return err
			}

			if err := generateMapFieldMethods(content, funcPrefix, field); err != nil {
				return err
			}

			// Generate individual key getter with default (__or) for map fields
			if err := generateUnifiedMapKeyGetter(content, funcPrefix, field, "__or", true, "default_value"); err != nil {
				return err
			}
		} else if field.IsList() {
			if err := generateListGetAll(content, funcPrefix, field); err != nil {
				return err
			}
			if err := generateListSetAll(content, funcPrefix, field); err != nil {
				return err
			}

			if err := generateRepeatedFieldMethods(content, funcPrefix, field); err != nil {
				return err
			}
		} else {
			// Generate basic singular getter using unified approach
			returnType := GetProtobufType(field.Kind()).GetSqlTypeName()
			var defaultValue string
			var converter ValueConverter

			if field.Message() != nil {
				// For message fields, return JSON object or empty object if not present
				defaultValue = "JSON_OBJECT()"
				converter = func(valueVar string) string {
					return valueVar // Just return the JSON value directly
				}
			} else {
				// For scalar fields, use default value and type conversion
				defaultValue = GetDefaultValue(field)
				converter = func(valueVar string) string {
					return GetProtobufType(field.Kind()).GenerateNumberJsonToSqlExpression(valueVar)
				}
			}

			if err := generateUnifiedFieldGetter(content, funcPrefix, field, "", returnType, false, defaultValue, converter, nil); err != nil {
				return err
			}

			if err := generateSingularSetter(content, funcPrefix, field, nil, nil); err != nil {
				return err
			}

			if field.HasPresence() {
				// Generate nullable getter with custom default for optional fields (__or)
				returnType := GetProtobufType(field.Kind()).GetSqlTypeName()
				converter := func(valueVar string) string {
					return GetProtobufType(field.Kind()).GenerateNumberJsonToSqlExpression(valueVar)
				}
				if err := generateUnifiedFieldGetter(content, funcPrefix, field, "__or", returnType, true, "default_value", converter, fieldFilterFunc); err != nil {
					return err
				}

				// Generate has method for fields with presence semantics
				if err := generateHasMethod(content, funcPrefix, field); err != nil {
					return err
				}
			}

			// Generate enum name getters for enum fields
			if field.Kind() == protoreflect.EnumKind {
				// Get enum conversion function names
				enumDesc := field.Enum()
				enumFullTypeName := enumDesc.FullName()
				enumPackageName := enumDesc.ParentFile().Package()
				enumFuncPrefix := typePrefixFunc(enumPackageName, enumFullTypeName)
				toStringFuncName := enumFuncPrefix + "_to_string"
				fromStringFuncName := enumFuncPrefix + "_from_string"

				// Create converter that handles enum to string conversion
				converter := func(valueVar string) string {
					return fmt.Sprintf(`COALESCE(%s(CAST(%s AS SIGNED)), CAST(CAST(%s AS SIGNED) AS CHAR))`, toStringFuncName, valueVar, valueVar)
				}

				// Generate regular enum name getter (__as_name)
				defaultValue := fmt.Sprintf(`COALESCE(%s(0), '0')`, toStringFuncName)
				if err := generateUnifiedFieldGetter(content, funcPrefix, field, "__as_name", "LONGTEXT", false, defaultValue, converter, fieldFilterFunc); err != nil {
					return err
				}

				// Generate nullable enum name getter (__as_name_or) for optional fields
				if field.HasPresence() {
					if err := generateUnifiedFieldGetter(content, funcPrefix, field, "__as_name_or", "LONGTEXT", true, "default_value", converter, fieldFilterFunc); err != nil {
						return err
					}
				}

				// Generate enum name setters for enum fields

				enumConverter := &InputConverter{
					Suffix:    "__from_name",
					InputType: "LONGTEXT",
					ConvertInput: func(inputVar string) (string, error) {
						return fmt.Sprintf("%s(%s)", fromStringFuncName, inputVar), nil
					},
					ValidationSQL: "IF converted_value IS NULL THEN SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid enum name'; END IF;",
				}

				if err := generateSingularSetter(content, funcPrefix, field, enumConverter, fieldFilterFunc); err != nil {
					return err
				}
			}
		}

		// Generate clear method for all fields
		if err := generateClearMethod(content, funcPrefix, field); err != nil {
			return err
		}
	}

	return nil
}

// generateListGetAll creates a getter for repeated fields
func generateListGetAll(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor) error {
	fieldName := string(field.Name())
	getterFuncName := fmt.Sprintf("%s_get_all_%s", funcPrefix, fieldName)
	if err := validateFunctionName(getterFuncName, field.ContainingMessage().FullName()); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS JSON DETERMINISTIC\n", getterFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE result_array JSON DEFAULT JSON_ARRAY();\n")
	content.WriteString("    DECLARE raw_array JSON;\n")
	content.WriteString("    DECLARE array_length INT;\n")
	content.WriteString("    DECLARE i INT DEFAULT 0;\n")
	content.WriteString("    DECLARE element_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET raw_array = COALESCE(JSON_EXTRACT(proto_data, '$.\"%.d\"'), JSON_ARRAY());\n", field.Number()))
	content.WriteString("    SET array_length = JSON_LENGTH(raw_array);\n")
	content.WriteString("    WHILE i < array_length DO\n")
	content.WriteString("        SET element_value = JSON_EXTRACT(raw_array, CONCAT('$[', i, ']'));\n")
	convertedExpression := GetProtobufType(field.Kind()).GenerateNumberJsonToJsonExpression("element_value")
	content.WriteString(fmt.Sprintf("        SET result_array = JSON_ARRAY_APPEND(result_array, '$', %s);\n", convertedExpression))
	content.WriteString("        SET i = i + 1;\n")
	content.WriteString("    END WHILE;\n")
	content.WriteString("    RETURN result_array;\n")
	content.WriteString("END $$\n\n")
	return nil
}

// generateMapGetAll creates a getter for map fields
func generateMapGetAll(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor) error {
	fieldName := string(field.Name())
	getterFuncName := fmt.Sprintf("%s_get_all_%s", funcPrefix, fieldName)
	if err := validateFunctionName(getterFuncName, field.ContainingMessage().FullName()); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS JSON DETERMINISTIC\n", getterFuncName))
	content.WriteString("BEGIN\n")

	// Check if value type needs conversion (like float/double)
	protobufType := GetProtobufType(field.MapValue().Kind())
	convertedExpression := protobufType.GenerateNumberJsonToJsonExpression("current_value")

	// For map fields with values that need conversion, iterate through and convert each value
	content.WriteString("    DECLARE result_map JSON DEFAULT JSON_OBJECT();\n")
	content.WriteString("    DECLARE raw_map JSON;\n")
	content.WriteString("    DECLARE json_keys JSON;\n")
	content.WriteString("    DECLARE key_count INT;\n")
	content.WriteString("    DECLARE i INT DEFAULT 0;\n")
	content.WriteString("    DECLARE current_key VARCHAR(255);\n")
	content.WriteString("    DECLARE current_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET raw_map = COALESCE(JSON_EXTRACT(proto_data, '$.\"%.d\"'), JSON_OBJECT());\n", field.Number()))
	content.WriteString("    SET json_keys = JSON_KEYS(raw_map);\n")
	content.WriteString("    IF json_keys IS NOT NULL THEN\n")
	content.WriteString("        SET key_count = JSON_LENGTH(json_keys);\n")
	content.WriteString("        WHILE i < key_count DO\n")
	content.WriteString("            SET current_key = JSON_UNQUOTE(JSON_EXTRACT(json_keys, CONCAT('$[', i, ']')));\n")
	content.WriteString("            SET current_value = JSON_EXTRACT(raw_map, CONCAT('$.\"', current_key, '\"'));\n")
	content.WriteString(fmt.Sprintf("            SET result_map = JSON_SET(result_map, CONCAT('$.\"', current_key, '\"'), %s);\n", convertedExpression))
	content.WriteString("            SET i = i + 1;\n")
	content.WriteString("        END WHILE;\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    RETURN result_map;\n")
	content.WriteString("END $$\n\n")
	return nil
}

// ValueConverter is a function that converts an extracted JSON value to the final return type
type ValueConverter func(valueVar string) string

// InputConverter handles custom input conversion for setter functions
type InputConverter struct {
	Suffix        string                                // Function name suffix (e.g., "__from_name")
	InputType     string                                // SQL input parameter type (e.g., "LONGTEXT")
	ConvertInput  func(inputVar string) (string, error) // Convert input to field value, returns SQL expression and error message
	ValidationSQL string                                // Optional SQL validation for converted value
}

// generateUnifiedFieldGetter creates regular field getter variants (__or, __as_name, __as_name_or) - NOT for maps
func generateUnifiedFieldGetter(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor, suffix string, returnType string, hasDefaultParam bool, defaultValue string, converter ValueConverter, fieldFilterFunc FieldFilterFunc) error {
	fieldName := string(field.Name())
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
		if err := validateFunctionName(getterFuncName, field.ContainingMessage().FullName()); err != nil {
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

// generateSingularSetter creates a setter for singular scalar/message fields
func generateSingularSetter(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor, inputConverter *InputConverter, fieldFilterFunc FieldFilterFunc) error {
	fieldName := string(field.Name())

	// Determine function name and input type based on converter
	var setterFuncName string
	var setterType string
	if inputConverter != nil {
		setterFuncName = fmt.Sprintf("%s_set_%s%s", funcPrefix, fieldName, inputConverter.Suffix)
		setterType = inputConverter.InputType
	} else {
		setterFuncName = fmt.Sprintf("%s_set_%s", funcPrefix, fieldName)
		setterType = GetProtobufType(field.Kind()).GetSqlTypeName()
	}

	// Apply field filtering
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
	if err := validateFunctionName(setterFuncName, field.ContainingMessage().FullName()); err != nil {
		if commentPrefix == "" {
			return err
		}
	}

	content.WriteString(fmt.Sprintf("%sDROP FUNCTION IF EXISTS %s $$\n", commentPrefix, setterFuncName))
	content.WriteString(fmt.Sprintf("%sCREATE FUNCTION %s(proto_data JSON, field_value %s) RETURNS JSON DETERMINISTIC\n", commentPrefix, setterFuncName, setterType))
	content.WriteString(fmt.Sprintf("%sBEGIN\n", commentPrefix))

	// Declare variables at the beginning (required by MySQL syntax)
	oneof := field.ContainingOneof()
	if oneof != nil {
		content.WriteString(fmt.Sprintf("%s    DECLARE temp_data JSON DEFAULT proto_data;\n", commentPrefix))
	}

	// Handle input conversion if needed
	var actualFieldValue string
	if inputConverter != nil {
		content.WriteString(fmt.Sprintf("%s    DECLARE converted_value %s;\n", commentPrefix, GetProtobufType(field.Kind()).GetSqlTypeName()))
		convertExpr, err := inputConverter.ConvertInput("field_value")
		if err != nil {
			return err
		}
		content.WriteString(fmt.Sprintf("%s    SET converted_value = %s;\n", commentPrefix, convertExpr))

		// Add validation if provided
		if inputConverter.ValidationSQL != "" {
			content.WriteString(fmt.Sprintf("%s    %s\n", commentPrefix, inputConverter.ValidationSQL))
		}
		actualFieldValue = "converted_value"
	} else {
		actualFieldValue = "field_value"
	}

	// Handle message field null case
	if field.Message() != nil {
		content.WriteString(fmt.Sprintf("%s    IF %s IS NULL THEN\n", commentPrefix, actualFieldValue))
		content.WriteString(fmt.Sprintf("%s        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n", commentPrefix, field.Number()))
		content.WriteString(fmt.Sprintf("%s    END IF;\n", commentPrefix))
	}

	// Handle oneOf field mutual exclusion
	if oneof != nil {
		content.WriteString(fmt.Sprintf("%s    -- OneOf field mutual exclusion: clear other fields in the same oneOf group\n", commentPrefix))

		// Clear all other fields in the same oneOf group
		fields := oneof.Fields()
		for otherField := range protoreflectutils.Iterate(fields) {
			if otherField.Number() != field.Number() {
				content.WriteString(fmt.Sprintf("%s    SET temp_data = JSON_REMOVE(temp_data, '$.\"%.d\"');\n", commentPrefix, otherField.Number()))
			}
		}

		// Set the new field value with proper JSON conversion - use unified type system
		jsonExpression := GetProtobufType(field.Kind()).GenerateSqlToNumberJsonExpression(actualFieldValue)
		content.WriteString(fmt.Sprintf("%s    RETURN JSON_SET(temp_data, '$.\"%.d\"', %s);\n", commentPrefix, field.Number(), jsonExpression))
	} else {
		// Check if this is a proto3 field without presence and being set to default value
		// According to protonumberjson spec, such fields should be omitted
		if !field.HasPresence() && field.Message() == nil {
			// For proto3 fields without presence, omit default values - use unified type system
			content.WriteString(fmt.Sprintf("%s    -- Proto3 field without presence: omit default values per protonumberjson spec\n", commentPrefix))
			setterLogic := GetProtobufType(field.Kind()).GenerateSetterWithZeroValueRemoval(int32(field.Number()), actualFieldValue)
			content.WriteString(fmt.Sprintf("%s    %s\n", commentPrefix, setterLogic))
		} else {
			// Fields with presence: always set the value - use unified type system
			jsonExpression := GetProtobufType(field.Kind()).GenerateSqlToNumberJsonExpression(actualFieldValue)
			content.WriteString(fmt.Sprintf("%s    RETURN JSON_SET(proto_data, '$.\"%.d\"', %s);\n", commentPrefix, field.Number(), jsonExpression))
		}
	}

	content.WriteString(fmt.Sprintf("%sEND $$\n\n", commentPrefix))
	return nil
}

// generateListSetAll creates a setter for repeated fields
func generateListSetAll(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor) error {
	fieldName := string(field.Name())
	setterFuncName := fmt.Sprintf("%s_set_all_%s", funcPrefix, fieldName)
	if err := validateFunctionName(setterFuncName, field.ContainingMessage().FullName()); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", setterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, field_value JSON) RETURNS JSON DETERMINISTIC\n", setterFuncName))
	content.WriteString("BEGIN\n")

	// For repeated fields, handle JSON array input
	content.WriteString("    DECLARE array_length INT;\n")
	content.WriteString("    DECLARE i INT DEFAULT 0;\n")
	content.WriteString("    DECLARE element_value JSON;\n")
	content.WriteString("    DECLARE converted_array JSON DEFAULT JSON_ARRAY();\n")
	content.WriteString("\n")
	content.WriteString("    IF field_value IS NULL THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'field_value cannot be NULL';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF JSON_TYPE(field_value) != 'ARRAY' THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'field_value must be a JSON array';\n")
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
}

// generateMapSetAll creates a setter for map fields
func generateMapSetAll(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor) error {
	fieldName := string(field.Name())
	setterFuncName := fmt.Sprintf("%s_set_all_%s", funcPrefix, fieldName)
	if err := validateFunctionName(setterFuncName, field.ContainingMessage().FullName()); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", setterFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, field_value JSON) RETURNS JSON DETERMINISTIC\n", setterFuncName))
	content.WriteString("BEGIN\n")

	// For map fields, use the same conversion logic as put_all
	protobufType := GetProtobufType(field.MapValue().Kind())
	jsonConversionExpr := protobufType.GenerateJsonToNumberJsonExpression("current_value")

	content.WriteString("    DECLARE json_keys JSON;\n")
	content.WriteString("    DECLARE key_count INT;\n")
	content.WriteString("    DECLARE i INT;\n")
	content.WriteString("    DECLARE current_key VARCHAR(255);\n")
	content.WriteString("    DECLARE current_value JSON;\n")
	content.WriteString("    DECLARE converted_value JSON;\n")
	content.WriteString("    DECLARE result_map JSON DEFAULT JSON_OBJECT();\n")

	content.WriteString("    IF field_value IS NULL THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'field_value cannot be NULL';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF JSON_TYPE(field_value) != 'OBJECT' THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'field_value must be a JSON object';\n")
	content.WriteString("    END IF;\n")

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
	content.WriteString("END $$\n\n")
	return nil
}

// generateHasMethod creates a method to check field presence
func generateHasMethod(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor) error {
	fieldName := string(field.Name())
	hasFuncName := fmt.Sprintf("%s_has_%s", funcPrefix, fieldName)
	if err := validateFunctionName(hasFuncName, field.ContainingMessage().FullName()); err != nil {
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
func generateClearMethod(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor) error {
	fieldName := string(field.Name())
	clearFuncName := fmt.Sprintf("%s_clear_%s", funcPrefix, fieldName)
	if err := validateFunctionName(clearFuncName, field.ContainingMessage().FullName()); err != nil {
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
func generateRepeatedFieldMethods(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor) error {
	// Generate add method for repeated fields
	fieldName := string(field.Name())
	addFuncName := fmt.Sprintf("%s_add_%s", funcPrefix, fieldName)
	if err := validateFunctionName(addFuncName, field.ContainingMessage().FullName()); err != nil {
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
	if err := validateFunctionName(countFuncName, field.ContainingMessage().FullName()); err != nil {
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
	if err := validateFunctionName(getIndexFuncName, field.ContainingMessage().FullName()); err != nil {
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
	if err := validateFunctionName(setIndexFuncName, field.ContainingMessage().FullName()); err != nil {
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
	if err := validateFunctionName(insertFuncName, field.ContainingMessage().FullName()); err != nil {
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
	if err := validateFunctionName(removeFuncName, field.ContainingMessage().FullName()); err != nil {
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
	if err := validateFunctionName(addAllFuncName, field.ContainingMessage().FullName()); err != nil {
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
	content.WriteString("    IF elements_array IS NULL THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'elements_array cannot be NULL';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF JSON_TYPE(elements_array) != 'ARRAY' THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'elements_array must be a JSON array';\n")
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

// generateUnifiedMapKeyGetter creates key-based getters for map fields (both standard and nullable variants)
func generateUnifiedMapKeyGetter(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor, suffix string, hasDefaultParam bool, defaultValue string) error {
	fieldName := string(field.Name())
	keyMySQLType := GetProtobufType(field.MapKey().Kind()).GetSqlTypeName()
	valueMySQLType := GetProtobufType(field.MapValue().Kind()).GetSqlTypeName()
	isValueMessage := field.MapValue().Kind() == protoreflect.MessageKind

	// For float/double types in nullable variant, use LONGTEXT to accept binary format strings
	if hasDefaultParam && (field.MapValue().Kind() == protoreflect.FloatKind || field.MapValue().Kind() == protoreflect.DoubleKind) {
		valueMySQLType = "LONGTEXT"
	}

	// Generate function name
	getKeyFuncName := fmt.Sprintf("%s_get_%s%s", funcPrefix, fieldName, suffix)
	if err := validateFunctionName(getKeyFuncName, field.ContainingMessage().FullName()); err != nil {
		return err
	}

	// Build parameters
	var parameters string
	if hasDefaultParam {
		parameters = fmt.Sprintf("proto_data JSON, key_value %s, default_value %s", keyMySQLType, valueMySQLType)
	} else {
		parameters = fmt.Sprintf("proto_data JSON, key_value %s", keyMySQLType)
	}

	// Generate function
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getKeyFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(%s) RETURNS %s DETERMINISTIC\n", getKeyFuncName, parameters, valueMySQLType))
	content.WriteString("BEGIN\n")

	// Declare variables
	content.WriteString("    DECLARE map_value JSON;\n")
	content.WriteString("    DECLARE element_value JSON;\n")

	// Get the map from the proto_data
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")
	content.WriteString(fmt.Sprintf("        RETURN %s;\n", defaultValue))
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
	content.WriteString(fmt.Sprintf("        RETURN %s;\n", defaultValue))
	content.WriteString("    END IF;\n")

	// Return the value, converting if necessary
	if hasDefaultParam && isValueMessage {
		// For nullable variant with message types, return without conversion
		content.WriteString("    RETURN element_value;\n")
	} else {
		// For primitive types or standard variant, use unified type conversion
		conversionExpression := GetProtobufType(field.MapValue().Kind()).GenerateNumberJsonToSqlExpression("element_value")
		content.WriteString(fmt.Sprintf("    RETURN %s;\n", conversionExpression))
	}

	content.WriteString("END $$\n\n")
	return nil
}

// generateMapFieldMethods creates additional methods for map fields
func generateMapFieldMethods(content *strings.Builder, funcPrefix string, field protoreflect.FieldDescriptor) error {
	// Generate key-based get method for map fields
	if err := generateUnifiedMapKeyGetter(content, funcPrefix, field, "", false, "NULL"); err != nil {
		return err
	}

	// Get MySQL types for key and value
	fieldName := string(field.Name())
	keySetterType := GetProtobufType(field.MapKey().Kind()).GetSqlTypeName()
	valueSetterType := GetProtobufType(field.MapValue().Kind()).GetSqlTypeName()

	// Generate contains method for map fields
	containsFuncName := fmt.Sprintf("%s_contains_%s", funcPrefix, fieldName)
	if err := validateFunctionName(containsFuncName, field.ContainingMessage().FullName()); err != nil {
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
	if err := validateFunctionName(putFuncName, field.ContainingMessage().FullName()); err != nil {
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
	content.WriteString("    IF key_value IS NULL THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'key_value cannot be NULL';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF value_param IS NULL THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'value_param cannot be NULL';\n")
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
	if err := validateFunctionName(putAllFuncName, field.ContainingMessage().FullName()); err != nil {
		return err
	}

	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", putAllFuncName))
	// For primitive value types that need conversion (float, double), we need to convert each value
	protobufType := GetProtobufType(field.MapValue().Kind())
	jsonConversionExpr := protobufType.GenerateJsonToNumberJsonExpression("current_value")

	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, new_entries JSON) RETURNS JSON DETERMINISTIC\n", putAllFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    DECLARE map_value JSON;\n")
	content.WriteString("    DECLARE json_keys JSON;\n")
	content.WriteString("    DECLARE key_count INT;\n")
	content.WriteString("    DECLARE i INT;\n")
	content.WriteString("    DECLARE current_key VARCHAR(255);\n")
	content.WriteString("    DECLARE current_value JSON;\n")
	content.WriteString("    DECLARE converted_value JSON;\n")
	content.WriteString(fmt.Sprintf("    SET map_value = JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
	content.WriteString("    IF map_value IS NULL THEN\n")
	content.WriteString("        SET map_value = JSON_OBJECT();\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF new_entries IS NULL THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'new_entries cannot be NULL';\n")
	content.WriteString("    END IF;\n")
	content.WriteString("    IF JSON_TYPE(new_entries) != 'OBJECT' THEN\n")
	content.WriteString("        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'new_entries must be a JSON object';\n")
	content.WriteString("    END IF;\n")

	// Need to convert each JSON value to number_json format (for float, double)
	content.WriteString("    SET json_keys = JSON_KEYS(new_entries);\n")
	content.WriteString("    SET key_count = JSON_LENGTH(json_keys);\n")
	content.WriteString("    SET i = 0;\n")
	content.WriteString("    \n")
	content.WriteString("    WHILE i < key_count DO\n")
	content.WriteString("        SET current_key = JSON_UNQUOTE(JSON_EXTRACT(json_keys, CONCAT('$[', i, ']')));\n")
	content.WriteString("        SET current_value = JSON_EXTRACT(new_entries, CONCAT('$.\"', current_key, '\"'));\n")
	content.WriteString(fmt.Sprintf("        SET converted_value = %s;\n", jsonConversionExpr))
	content.WriteString("        SET map_value = JSON_SET(map_value, CONCAT('$.\"', current_key, '\"'), converted_value);\n")
	content.WriteString("        SET i = i + 1;\n")
	content.WriteString("    END WHILE;\n")
	content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', map_value);\n", field.Number()))
	content.WriteString("END $$\n\n")

	// Generate remove method for map fields
	removeFuncName := fmt.Sprintf("%s_remove_%s", funcPrefix, fieldName)
	if err := validateFunctionName(removeFuncName, field.ContainingMessage().FullName()); err != nil {
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
	if err := validateFunctionName(countFuncName, field.ContainingMessage().FullName()); err != nil {
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
func generateOneOfMethods(content *strings.Builder, messageDesc protoreflect.MessageDescriptor, funcPrefix string) error {
	fullTypeName := messageDesc.FullName()
	oneofs := messageDesc.Oneofs()
	for oneof := range protoreflectutils.Iterate(oneofs) {
		// Skip synthetic oneofs created by proto3 optional fields
		// These are implementation details and shouldn't have user-facing accessor methods
		if oneof.IsSynthetic() {
			continue
		}

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

func generateEnumMethods(content *strings.Builder, enumDesc protoreflect.EnumDescriptor, typePrefixFunc TypePrefixFunc) error {
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
