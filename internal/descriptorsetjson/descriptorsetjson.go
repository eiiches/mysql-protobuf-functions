package descriptorsetjson

import (
	"encoding/json"
	"fmt"

	"github.com/eiiches/mysql-protobuf-functions/internal/descriptorsetpb"
	"github.com/eiiches/mysql-protobuf-functions/internal/protonumberjson"
	"google.golang.org/protobuf/types/descriptorpb"
)

// ToJson converts a FileDescriptorSet to the MySQL-compatible JSON format
// Returns a DescriptorSet message serialized using protonumberjson format
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
// Returns a DescriptorSet message serialized using protonumberjson format
func ToJsonTree(fileDescriptorSet *descriptorpb.FileDescriptorSet) (interface{}, error) {
	if fileDescriptorSet == nil {
		return nil, fmt.Errorf("fileDescriptorSet cannot be nil")
	}

	// Build the type indexes
	messageTypeIndex, enumTypeIndex := buildTypeIndexes(fileDescriptorSet)

	// Create the DescriptorSet message
	descriptorSet := &descriptorsetpb.DescriptorSet{
		FileDescriptorSet: fileDescriptorSet,
		MessageTypeIndex:  messageTypeIndex,
		EnumTypeIndex:     enumTypeIndex,
	}

	// Convert to JSON tree using protonumberjson
	return protonumberjson.ToJsonTree(descriptorSet)
}

// buildTypeIndexes creates mappings of fully-qualified type names to their respective index messages
func buildTypeIndexes(fileDescriptorSet *descriptorpb.FileDescriptorSet) (map[string]*descriptorsetpb.MessageTypeIndex, map[string]*descriptorsetpb.EnumTypeIndex) {
	messageIndex := make(map[string]*descriptorsetpb.MessageTypeIndex)
	enumIndex := make(map[string]*descriptorsetpb.EnumTypeIndex)

	for fileIndex, fileDesc := range fileDescriptorSet.File {
		filePackage := ""
		if fileDesc.Package != nil {
			filePackage = *fileDesc.Package
		}

		// JSON paths now reference the FileDescriptorSet root directly
		filePath := fmt.Sprintf("$.\"1\"[%d]", fileIndex)

		// Process message types
		for msgIndex, msgDesc := range fileDesc.MessageType {
			msgPath := fmt.Sprintf("%s.\"4\"[%d]", filePath, msgIndex)
			msgName := buildTypeName(filePackage, *msgDesc.Name)

			messageTypeIndex := &descriptorsetpb.MessageTypeIndex{
				FileJsonPath:     filePath,
				TypeJsonPath:     msgPath,
				FieldNameIndex:   buildFieldNameIndex(msgDesc),
				FieldNumberIndex: buildFieldNumberIndex(msgDesc),
			}
			messageIndex[msgName] = messageTypeIndex

			// Process nested types recursively
			buildNestedTypes(messageIndex, enumIndex, msgDesc, msgName, msgPath, filePath)
		}

		// Process enum types
		for enumIdx, enumDesc := range fileDesc.EnumType {
			enumPath := fmt.Sprintf("%s.\"5\"[%d]", filePath, enumIdx)
			enumName := buildTypeName(filePackage, *enumDesc.Name)

			enumTypeIndex := &descriptorsetpb.EnumTypeIndex{
				FileJsonPath:    filePath,
				TypeJsonPath:    enumPath,
				EnumNameIndex:   buildEnumNameIndex(enumDesc),
				EnumNumberIndex: buildEnumNumberIndex(enumDesc),
			}
			enumIndex[enumName] = enumTypeIndex
		}
	}

	return messageIndex, enumIndex
}

// buildNestedTypes recursively processes nested message and enum types
func buildNestedTypes(messageIndex map[string]*descriptorsetpb.MessageTypeIndex, enumIndex map[string]*descriptorsetpb.EnumTypeIndex, msgDesc *descriptorpb.DescriptorProto, parentName, parentPath, filePath string) {
	// Process nested message types
	for nestedMsgIndex, nestedMsgDesc := range msgDesc.NestedType {
		nestedMsgPath := fmt.Sprintf("%s.\"3\"[%d]", parentPath, nestedMsgIndex)
		nestedMsgName := parentName + "." + *nestedMsgDesc.Name

		messageTypeIndex := &descriptorsetpb.MessageTypeIndex{
			FileJsonPath:     filePath,
			TypeJsonPath:     nestedMsgPath,
			FieldNameIndex:   buildFieldNameIndex(nestedMsgDesc),
			FieldNumberIndex: buildFieldNumberIndex(nestedMsgDesc),
		}
		messageIndex[nestedMsgName] = messageTypeIndex

		// Recursively process further nested types
		buildNestedTypes(messageIndex, enumIndex, nestedMsgDesc, nestedMsgName, nestedMsgPath, filePath)
	}

	// Process nested enum types
	for nestedEnumIndex, nestedEnumDesc := range msgDesc.EnumType {
		nestedEnumPath := fmt.Sprintf("%s.\"4\"[%d]", parentPath, nestedEnumIndex)
		nestedEnumName := parentName + "." + *nestedEnumDesc.Name

		enumTypeIndex := &descriptorsetpb.EnumTypeIndex{
			FileJsonPath:    filePath,
			TypeJsonPath:    nestedEnumPath,
			EnumNameIndex:   buildEnumNameIndex(nestedEnumDesc),
			EnumNumberIndex: buildEnumNumberIndex(nestedEnumDesc),
		}
		enumIndex[nestedEnumName] = enumTypeIndex
	}
}

// buildTypeName constructs a fully-qualified type name
func buildTypeName(packageName, typeName string) string {
	if packageName == "" {
		return "." + typeName
	}
	return "." + packageName + "." + typeName
}

// buildFieldNameIndex creates a map from field names to their array indices in the message descriptor
func buildFieldNameIndex(msgDesc *descriptorpb.DescriptorProto) map[string]int32 {
	index := make(map[string]int32)
	for i, field := range msgDesc.Field {
		if field.Name != nil {
			index[*field.Name] = int32(i)
		}
	}
	return index
}

// buildFieldNumberIndex creates a map from field numbers to their array indices in the message descriptor
func buildFieldNumberIndex(msgDesc *descriptorpb.DescriptorProto) map[int32]int32 {
	index := make(map[int32]int32)
	for i, field := range msgDesc.Field {
		if field.Number != nil {
			index[*field.Number] = int32(i)
		}
	}
	return index
}

// buildEnumNameIndex creates a map from enum value names to their array indices in the enum descriptor
func buildEnumNameIndex(enumDesc *descriptorpb.EnumDescriptorProto) map[string]int32 {
	index := make(map[string]int32)
	for i, value := range enumDesc.Value {
		if value.Name != nil {
			index[*value.Name] = int32(i)
		}
	}
	return index
}

// buildEnumNumberIndex creates a map from enum value numbers to their array indices in the enum descriptor
func buildEnumNumberIndex(enumDesc *descriptorpb.EnumDescriptorProto) map[int32]int32 {
	index := make(map[int32]int32)
	for i, value := range enumDesc.Value {
		if value.Number != nil {
			index[*value.Number] = int32(i)
		}
	}
	return index
}
