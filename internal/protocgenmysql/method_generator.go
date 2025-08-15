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
func GenerateMethodFragments(protoFiles []*descriptorpb.FileDescriptorProto, fileNameFunc FileNameFunc, typePrefixFunc TypePrefixFunc) map[string][]string {
	fileFragments := make(map[string][]string)

	// Generate fragments for each proto file
	for _, file := range protoFiles {
		if file.Name == nil {
			continue
		}
		filename := fileNameFunc(*file.Name)
		content := generateMethodsForFileContent(file, typePrefixFunc)
		if content != "" {
			fileFragments[filename] = append(fileFragments[filename], content)
		}
	}

	return fileFragments
}

func generateMethodsForFileContent(file *descriptorpb.FileDescriptorProto, typePrefixFunc TypePrefixFunc) string {
	var content strings.Builder

	packageName := ""
	if file.Package != nil {
		packageName = *file.Package
	}

	// Generate methods for each message type
	for _, messageType := range file.MessageType {
		generateMessageMethods(&content, messageType, typePrefixFunc, packageName)
	}

	return content.String()
}

func generateMessageMethods(content *strings.Builder, messageType *descriptorpb.DescriptorProto, typePrefixFunc TypePrefixFunc, packageName string) {
	generateMessageMethodsWithPath(content, messageType, typePrefixFunc, packageName, "")
}

func generateMessageMethodsWithPath(content *strings.Builder, messageType *descriptorpb.DescriptorProto, typePrefixFunc TypePrefixFunc, packageName string, parentPath string) {
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
	content.WriteString("    -- not implemented\n")
	content.WriteString("    RETURN JSON_OBJECT();\n")
	content.WriteString("END $$\n\n")

	// Generate from_json
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_from_json $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_from_json(json_data JSON) RETURNS JSON DETERMINISTIC\n", funcPrefix))
	content.WriteString("BEGIN\n")
	content.WriteString("    -- not implemented\n")
	content.WriteString("    RETURN json_data;\n")
	content.WriteString("END $$\n\n")

	// Generate from_message
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_from_message $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_from_message(message_data LONGBLOB) RETURNS JSON DETERMINISTIC\n", funcPrefix))
	content.WriteString("BEGIN\n")
	content.WriteString("    -- not implemented\n")
	content.WriteString("    RETURN JSON_OBJECT();\n")
	content.WriteString("END $$\n\n")

	// Generate to_json
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_to_json $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_to_json(proto_data JSON) RETURNS JSON DETERMINISTIC\n", funcPrefix))
	content.WriteString("BEGIN\n")
	content.WriteString("    -- not implemented\n")
	content.WriteString("    RETURN proto_data;\n")
	content.WriteString("END $$\n\n")

	// Generate to_message
	content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_to_message $$\n", funcPrefix))
	content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_to_message(proto_data JSON) RETURNS LONGBLOB DETERMINISTIC\n", funcPrefix))
	content.WriteString("BEGIN\n")
	content.WriteString("    -- not implemented\n")
	content.WriteString("    RETURN NULL;\n")
	content.WriteString("END $$\n\n")

	// Generate setter and getter methods for each field
	for _, field := range messageType.Field {
		if field.Name == nil {
			continue
		}
		fieldName := *field.Name

		// Generate setter
		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_set_%s $$\n", funcPrefix, fieldName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_set_%s(proto_data JSON, field_value JSON) RETURNS JSON DETERMINISTIC\n", funcPrefix, fieldName))
		content.WriteString("BEGIN\n")
		content.WriteString("    -- not implemented\n")
		content.WriteString(fmt.Sprintf("    RETURN JSON_SET(proto_data, '$.%s', field_value);\n", fieldName))
		content.WriteString("END $$\n\n")

		// Generate getter
		content.WriteString(fmt.Sprintf("DROP FUNCTION IF EXISTS %s_get_%s $$\n", funcPrefix, fieldName))
		content.WriteString(fmt.Sprintf("CREATE FUNCTION %s_get_%s(proto_data JSON) RETURNS JSON DETERMINISTIC\n", funcPrefix, fieldName))
		content.WriteString("BEGIN\n")
		content.WriteString("    -- not implemented\n")
		content.WriteString(fmt.Sprintf("    RETURN JSON_EXTRACT(proto_data, '$.%s');\n", fieldName))
		content.WriteString("END $$\n\n")
	}

	// Generate methods for nested message types
	for _, nestedType := range messageType.NestedType {
		generateMessageMethodsWithPath(content, nestedType, typePrefixFunc, packageName, fullTypeName)
	}
}
