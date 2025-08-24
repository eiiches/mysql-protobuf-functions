package protocgenmysql

import (
	"fmt"
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

	// Generate fragments for each proto file
	var generationErr error
	files.RangeFiles(func(fileDesc protoreflect.FileDescriptor) bool {
		filename := fileNameFunc(fileDesc.Path())
		content, err := generateMethodsForFile(fileDesc, typePrefixFunc, schemaFunctionName)
		if err != nil {
			generationErr = err
			return false // stop iteration
		}
		if content != "" {
			fileFragments[filename] = append(fileFragments[filename], content)
		}
		return true // continue iteration
	})

	if generationErr != nil {
		return nil, generationErr
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

	// Generate constructor
	newFuncName := funcPrefix + "_new"
	if err := validateFunctionName(newFuncName, fullTypeName); err != nil {
		return err
	}
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", newFuncName))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s() RETURNS JSON DETERMINISTIC\n", newFuncName))
	content.WriteString("BEGIN\n")
	content.WriteString("    RETURN JSON_OBJECT();\n")
	content.WriteString("END $$\n\n")

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

	// Generate setter and getter methods for each field
	fields := messageDesc.Fields()
	for field := range protoreflectutils.Iterate(fields) {
		fieldName := string(field.Name())
		fieldType := field.Kind()
		isRepeated := field.Cardinality() == protoreflect.Repeated

		// Determine MySQL types for setter parameter and getter return
		setterType, getterType := getMySQLTypesForFieldFromReflection(fieldType, isRepeated)

		// Generate setter
		setterFuncName := fmt.Sprintf("%s_set_%s", funcPrefix, fieldName)
		if err := validateFunctionName(setterFuncName, fullTypeName); err != nil {
			return err
		}
		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", setterFuncName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON, field_value %s) RETURNS JSON DETERMINISTIC\n", setterFuncName, setterType))
		content.WriteString("BEGIN\n")
		content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
		content.WriteString("END $$\n\n")

		// Generate getter
		getterFuncName := fmt.Sprintf("%s_get_%s", funcPrefix, fieldName)
		if err := validateFunctionName(getterFuncName, fullTypeName); err != nil {
			return err
		}
		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s $$\n", getterFuncName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s(proto_data JSON) RETURNS %s DETERMINISTIC\n", getterFuncName, getterType))
		content.WriteString("BEGIN\n")
		if isRepeated || fieldType == protoreflect.MessageKind || fieldType == protoreflect.BoolKind {
			// For repeated fields and messages, return JSON directly
			content.WriteString(fmt.Sprintf("    RETURN JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.Number()))
		} else {
			// For scalar fields, unquote the JSON value with default value fallback
			defaultValue := getFieldDefaultValueFromReflection(field)
			content.WriteString(fmt.Sprintf("    RETURN COALESCE(JSON_UNQUOTE(JSON_EXTRACT(proto_data, '$.\"%.d\"')), %s);\n", field.Number(), defaultValue))
		}
		content.WriteString("END $$\n\n")
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

func generateEnumMethods(content *strings.Builder, enumDesc protoreflect.EnumDescriptor, typePrefixFunc TypePrefixFunc, schemaFunctionName string) error {
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
