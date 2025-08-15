package protocgenmysql

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

// FileNameFunc defines how to transform a proto file path to a SQL file path
type FileNameFunc func(protoPath string) string

// TypePrefixFunc defines how to transform a proto package and type name to a SQL function prefix
type TypePrefixFunc func(packageName string, typeName string) string

// GenerateMethodFragments returns individual content fragments for each proto file
func GenerateMethodFragments(protoFiles []*descriptorpb.FileDescriptorProto, fileNameFunc FileNameFunc, typePrefixFunc TypePrefixFunc, schemaFunctionName string) map[string][]string {
	fileFragments := make(map[string][]string)

	// Generate fragments for each proto file
	for _, file := range protoFiles {
		if file.Name == nil {
			continue
		}
		filename := fileNameFunc(*file.Name)
		content := generateMethodsForFileContent(file, typePrefixFunc, schemaFunctionName)
		if content != "" {
			fileFragments[filename] = append(fileFragments[filename], content)
		}
	}

	return fileFragments
}

func generateMethodsForFileContent(file *descriptorpb.FileDescriptorProto, typePrefixFunc TypePrefixFunc, schemaFunctionName string) string {
	var content strings.Builder

	packageName := ""
	if file.Package != nil {
		packageName = *file.Package
	}

	// Generate methods for each message type
	for _, messageType := range file.MessageType {
		generateMessageMethods(&content, messageType, typePrefixFunc, packageName, schemaFunctionName)
	}

	return content.String()
}

func generateMessageMethods(content *strings.Builder, messageType *descriptorpb.DescriptorProto, typePrefixFunc TypePrefixFunc, packageName string, schemaFunctionName string) {
	generateMessageMethodsWithPath(content, messageType, typePrefixFunc, packageName, "", schemaFunctionName)
}

func generateMessageMethodsWithPath(content *strings.Builder, messageType *descriptorpb.DescriptorProto, typePrefixFunc TypePrefixFunc, packageName string, parentPath string, schemaFunctionName string) {
	if messageType.Name == nil {
		return
	}

	typeName := *messageType.Name

	// Build fully qualified type name
	var fullTypeName string
	if parentPath != "" {
		fullTypeName = parentPath + "." + typeName
	} else if packageName != "" {
		fullTypeName = packageName + "." + typeName
	} else {
		fullTypeName = typeName
	}

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
	for _, field := range messageType.Field {
		if field.Name == nil {
			continue
		}
		fieldName := *field.Name
		fieldType := field.GetType()
		isRepeated := field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED

		// Determine MySQL types for setter parameter and getter return
		setterType, getterType := getMySQLTypesForField(fieldType, isRepeated)

		// Generate setter
		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_set_%s $$\n", funcPrefix, fieldName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_set_%s(proto_data JSON, field_value %s) RETURNS JSON DETERMINISTIC\n", funcPrefix, fieldName, setterType))
		content.WriteString("BEGIN\n")
		content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.\"%.d\"', field_value);\n", field.GetNumber()))
		content.WriteString("END $$\n\n")

		// Generate getter
		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_get_%s $$\n", funcPrefix, fieldName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_get_%s(proto_data JSON) RETURNS %s DETERMINISTIC\n", funcPrefix, fieldName, getterType))
		content.WriteString("BEGIN\n")
		if isRepeated || fieldType == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			// For repeated fields and messages, return JSON directly
			content.WriteString(fmt.Sprintf("    RETURN JSON_EXTRACT(proto_data, '$.\"%.d\"');\n", field.GetNumber()))
		} else {
			// For scalar fields, unquote the JSON value with default value fallback
			defaultValue := getFieldDefaultValue(field)
			content.WriteString(fmt.Sprintf("    RETURN COALESCE(JSON_UNQUOTE(JSON_EXTRACT(proto_data, '$.\"%.d\"')), %s);\n", field.GetNumber(), defaultValue))
		}
		content.WriteString("END $$\n\n")
	}

	// Generate methods for nested message types
	for _, nestedType := range messageType.NestedType {
		generateMessageMethodsWithPath(content, nestedType, typePrefixFunc, packageName, fullTypeName, schemaFunctionName)
	}
}

// getMySQLTypesForField returns the appropriate MySQL types for setter parameter and getter return value
func getMySQLTypesForField(fieldType descriptorpb.FieldDescriptorProto_Type, isRepeated bool) (setterType, getterType string) {
	// For repeated fields and maps, always use JSON
	if isRepeated {
		return "JSON", "JSON"
	}

	// For scalar fields, map protobuf types to MySQL types
	switch fieldType {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return "DOUBLE", "DOUBLE"
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return "FLOAT", "FLOAT"
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		return "BIGINT", "BIGINT"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		return "BIGINT UNSIGNED", "BIGINT UNSIGNED"
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		return "INT", "INT"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return "BIGINT UNSIGNED", "BIGINT UNSIGNED"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		return "INT UNSIGNED", "INT UNSIGNED"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "BOOLEAN", "BOOLEAN"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "LONGTEXT", "LONGTEXT"
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "LONGBLOB", "LONGBLOB"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		return "INT UNSIGNED", "INT UNSIGNED"
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return "INT", "INT"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return "INT", "INT"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		return "BIGINT", "BIGINT"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		return "INT", "INT"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return "BIGINT", "BIGINT"
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		// Messages are represented as JSON (ProtoNumberJSON format)
		return "JSON", "JSON"
	default:
		// Fallback to JSON for unknown types
		return "JSON", "JSON"
	}
}

// getFieldDefaultValue returns the appropriate default value for a field based on proto version and custom defaults
func getFieldDefaultValue(field *descriptorpb.FieldDescriptorProto) string {
	fieldType := field.GetType()

	// Check if field has a custom default value (proto2)
	if field.DefaultValue != nil {
		defaultVal := *field.DefaultValue
		switch fieldType {
		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			// String defaults need to be quoted
			return fmt.Sprintf("'%s'", strings.ReplaceAll(defaultVal, "'", "''"))
		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			// Bytes defaults need to be hex-encoded
			return fmt.Sprintf("X'%s'", defaultVal)
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			// Boolean defaults
			if defaultVal == "true" {
				return "1"
			}
			return "0"
		default:
			// Numeric defaults can be used directly
			return defaultVal
		}
	}

	// Use proto3 zero values or proto2 implicit defaults
	switch fieldType {
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE, descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return "0.0"
	case descriptorpb.FieldDescriptorProto_TYPE_INT64, descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_INT32, descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64, descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32, descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32, descriptorpb.FieldDescriptorProto_TYPE_SINT64,
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return "0"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "0"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return "''"
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "X''"
	default:
		return "NULL"
	}
}
