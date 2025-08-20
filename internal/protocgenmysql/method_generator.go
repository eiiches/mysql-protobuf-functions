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
func GenerateMethodFragments(files *protoregistry.Files, fileNameFunc FileNameFunc, typePrefixFunc TypePrefixFunc, schemaFunctionName string) map[string][]string {
	fileFragments := make(map[string][]string)

	// Generate fragments for each proto file
	files.RangeFiles(func(fileDesc protoreflect.FileDescriptor) bool {
		filename := fileNameFunc(fileDesc.Path())
		content := generateMethodsForFile(fileDesc, typePrefixFunc, schemaFunctionName)
		if content != "" {
			fileFragments[filename] = append(fileFragments[filename], content)
		}
		return true // continue iteration
	})

	return fileFragments
}

func generateMethodsForFile(fileDesc protoreflect.FileDescriptor, typePrefixFunc TypePrefixFunc, schemaFunctionName string) string {
	var content strings.Builder

	// Generate methods for each message type
	messages := fileDesc.Messages()
	for messageDesc := range protoreflectutils.Iterate(messages) {
		generateMessageMethods(&content, messageDesc, typePrefixFunc, schemaFunctionName)
	}

	// Generate methods for each enum type
	enums := fileDesc.Enums()
	for enumDesc := range protoreflectutils.Iterate(enums) {
		generateEnumMethods(&content, enumDesc, typePrefixFunc, schemaFunctionName)
	}

	return content.String()
}

func generateMessageMethods(content *strings.Builder, messageDesc protoreflect.MessageDescriptor, typePrefixFunc TypePrefixFunc, schemaFunctionName string) {
	// Use FullName from descriptor - no manual string construction needed
	fullTypeName := messageDesc.FullName()
	packageName := messageDesc.ParentFile().Package()

	// Get prefix for this specific type
	funcPrefix := typePrefixFunc(packageName, fullTypeName)

	// Generate constructor
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_new $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_new() RETURNS JSON DETERMINISTIC\n", funcPrefix))
	content.WriteString("BEGIN\n")
	content.WriteString("    RETURN JSON_OBJECT();\n")
	content.WriteString("END $$\n\n")

	// Generate from_json
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_from_json $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_from_json(json_data JSON) RETURNS JSON DETERMINISTIC\n", funcPrefix))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_json_to_number_json(%s(), '.%s', json_data);\n", schemaFunctionName, fullTypeName))
	content.WriteString("END $$\n\n")

	// Generate from_message
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_from_message $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_from_message(message_data LONGBLOB) RETURNS JSON DETERMINISTIC\n", funcPrefix))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_message_to_number_json(%s(), '.%s', message_data);\n", schemaFunctionName, fullTypeName))
	content.WriteString("END $$\n\n")

	// Generate to_json
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_to_json $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_to_json(proto_data JSON) RETURNS JSON DETERMINISTIC\n", funcPrefix))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_number_json_to_json(%s(), '.%s', proto_data, TRUE);\n", schemaFunctionName, fullTypeName))
	content.WriteString("END $$\n\n")

	// Generate to_message
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_to_message $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_to_message(proto_data JSON) RETURNS LONGBLOB DETERMINISTIC\n", funcPrefix))
	content.WriteString("BEGIN\n")
	content.WriteString(fmt.Sprintf("    RETURN _pb_number_json_to_message(%s(), '.%s', proto_data);\n", schemaFunctionName, fullTypeName))
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
		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_set_%s $$\n", funcPrefix, fieldName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_set_%s(proto_data JSON, field_value %s) RETURNS JSON DETERMINISTIC\n", funcPrefix, fieldName, setterType))
		content.WriteString("BEGIN\n")
		content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.Number()))
		content.WriteString("END $$\n\n")

		// Generate getter
		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_get_%s $$\n", funcPrefix, fieldName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_get_%s(proto_data JSON) RETURNS %s DETERMINISTIC\n", funcPrefix, fieldName, getterType))
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
		generateMessageMethods(content, nestedMessageDesc, typePrefixFunc, schemaFunctionName)
	}

	// Generate methods for nested enum types
	nestedEnums := messageDesc.Enums()
	for nestedEnumDesc := range protoreflectutils.Iterate(nestedEnums) {
		generateEnumMethods(content, nestedEnumDesc, typePrefixFunc, schemaFunctionName)
	}
}

func generateEnumMethods(content *strings.Builder, enumDesc protoreflect.EnumDescriptor, typePrefixFunc TypePrefixFunc, schemaFunctionName string) {
	// Use FullName from descriptor
	fullTypeName := enumDesc.FullName()
	packageName := enumDesc.ParentFile().Package()

	// Get prefix for this specific type
	funcPrefix := typePrefixFunc(packageName, fullTypeName)

	// Generate from_string method
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_from_string $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_from_string(enum_name LONGTEXT) RETURNS INT DETERMINISTIC\n", funcPrefix))
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
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_to_string $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_to_string(enum_value INT) RETURNS LONGTEXT DETERMINISTIC\n", funcPrefix))
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
