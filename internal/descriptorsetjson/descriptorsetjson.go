package descriptorsetjson

import (
	"encoding/json"
	"fmt"

	"github.com/eiiches/mysql-protobuf-functions/internal/protonumberjson"
	"google.golang.org/protobuf/types/descriptorpb"
)

// TypeIndex represents a type reference with kind, file path, and type path
type TypeIndex [3]interface{}

// Result represents the complete descriptor set JSON structure
type Result struct {
	FileDescriptorSet interface{}          `json:"fileDescriptorSet"`
	TypeIndex         map[string]TypeIndex `json:"typeIndex"`
}

// ToJson converts a FileDescriptorSet to the MySQL-compatible JSON format
// Returns a 3-element array: [version, fileDescriptorSet, typeIndex]
func ToJson(fileDescriptorSet *descriptorpb.FileDescriptorSet) (string, error) {
	jsonTree, err := ToJsonTree(fileDescriptorSet)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(jsonTree)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON tree: %w", err)
	}

	return string(jsonBytes), nil
}

// ToJsonTree converts a FileDescriptorSet to the MySQL-compatible JSON tree structure
// Returns a 3-element array: [version, fileDescriptorSet, typeIndex]
func ToJsonTree(fileDescriptorSet *descriptorpb.FileDescriptorSet) ([3]interface{}, error) {
	if fileDescriptorSet == nil {
		return [3]interface{}{}, fmt.Errorf("fileDescriptorSet cannot be nil")
	}

	// Convert fileDescriptorSet to JSON tree using protonumberjson
	fileDescriptorSetTree, err := protonumberjson.ToJsonTree(fileDescriptorSet)
	if err != nil {
		return [3]interface{}{}, fmt.Errorf("failed to convert FileDescriptorSet to JSON tree: %w", err)
	}

	// Build the type index
	typeIndex := buildTypeIndex(fileDescriptorSet)

	return [3]interface{}{1, fileDescriptorSetTree, typeIndex}, nil
}

// buildTypeIndex creates a mapping of fully-qualified type names to their JSON path locations
func buildTypeIndex(fileDescriptorSet *descriptorpb.FileDescriptorSet) map[string]TypeIndex {
	index := make(map[string]TypeIndex)

	for fileIndex, fileDesc := range fileDescriptorSet.File {
		filePackage := ""
		if fileDesc.Package != nil {
			filePackage = *fileDesc.Package
		}

		filePath := fmt.Sprintf("$[1].\"1\"[%d]", fileIndex)

		// Process message types
		for msgIndex, msgDesc := range fileDesc.MessageType {
			msgPath := fmt.Sprintf("%s.\"4\"[%d]", filePath, msgIndex)
			msgName := buildTypeName(filePackage, *msgDesc.Name)
			index[msgName] = TypeIndex{11, filePath, msgPath} // TYPE_MESSAGE

			// Process nested types recursively
			buildNestedTypes(index, msgDesc, msgName, msgPath, filePath)
		}

		// Process enum types
		for enumIndex, enumDesc := range fileDesc.EnumType {
			enumPath := fmt.Sprintf("%s.\"5\"[%d]", filePath, enumIndex)
			enumName := buildTypeName(filePackage, *enumDesc.Name)
			index[enumName] = TypeIndex{14, filePath, enumPath} // TYPE_ENUM
		}
	}

	return index
}

// buildNestedTypes recursively processes nested message and enum types
func buildNestedTypes(index map[string]TypeIndex, msgDesc *descriptorpb.DescriptorProto, parentName, parentPath, filePath string) {
	// Process nested message types
	for nestedMsgIndex, nestedMsgDesc := range msgDesc.NestedType {
		nestedMsgPath := fmt.Sprintf("%s.\"3\"[%d]", parentPath, nestedMsgIndex)
		nestedMsgName := parentName + "." + *nestedMsgDesc.Name
		index[nestedMsgName] = TypeIndex{11, filePath, nestedMsgPath} // TYPE_MESSAGE

		// Recursively process further nested types
		buildNestedTypes(index, nestedMsgDesc, nestedMsgName, nestedMsgPath, filePath)
	}

	// Process nested enum types
	for nestedEnumIndex, nestedEnumDesc := range msgDesc.EnumType {
		nestedEnumPath := fmt.Sprintf("%s.\"4\"[%d]", parentPath, nestedEnumIndex)
		nestedEnumName := parentName + "." + *nestedEnumDesc.Name
		index[nestedEnumName] = TypeIndex{14, filePath, nestedEnumPath} // TYPE_ENUM
	}
}

// buildTypeName constructs a fully-qualified type name
func buildTypeName(packageName, typeName string) string {
	if packageName == "" {
		return "." + typeName
	}
	return "." + packageName + "." + typeName
}
