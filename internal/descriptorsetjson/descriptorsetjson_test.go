package descriptorsetjson

import (
	"encoding/json"
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/protoreflectutils"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestToJson(t *testing.T) {
	g := NewWithT(t)

	t.Run("with descriptor.proto", func(t *testing.T) {
		// Get the descriptor.proto file descriptor
		fileDescriptor := descriptorpb.File_google_protobuf_descriptor_proto

		// Build a FileDescriptorSet with dependencies
		fileDescriptorSet := protoreflectutils.BuildFileDescriptorSetWithDependencies(fileDescriptor)

		// Convert to JSON
		jsonStr, err := ToJson(fileDescriptorSet)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(jsonStr).ToNot(BeEmpty())

		// Verify it's valid JSON
		var result interface{}
		err = json.Unmarshal([]byte(jsonStr), &result)
		g.Expect(err).ToNot(HaveOccurred())

		// Verify it's a DescriptorSet message object with protonumberjson format
		resultMap, ok := result.(map[string]interface{})
		g.Expect(ok).To(BeTrue())

		// Verify field "1" is the FileDescriptorSet
		fileDescriptorSetData := resultMap["1"]
		g.Expect(fileDescriptorSetData).ToNot(BeNil())

		// Verify field "2" is the message type index
		messageTypeIndexData, ok := resultMap["2"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(messageTypeIndexData).ToNot(BeEmpty())

		// Verify field "3" is the enum type index
		enumTypeIndexData, ok := resultMap["3"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())

		// Verify some expected message types exist in the message index
		expectedMessageTypes := []string{
			".google.protobuf.FileDescriptorSet",
			".google.protobuf.FileDescriptorProto",
			".google.protobuf.DescriptorProto",
			".google.protobuf.FieldDescriptorProto",
		}

		for _, expectedType := range expectedMessageTypes {
			g.Expect(messageTypeIndexData).To(HaveKey(expectedType), "missing message type: %s", expectedType)

			// Verify MessageTypeIndex structure (protonumberjson format)
			typeIndex, ok := messageTypeIndexData[expectedType].(map[string]interface{})
			g.Expect(ok).To(BeTrue(), "message type index should be map for %s", expectedType)

			// Verify file path (field "1")
			filePath, ok := typeIndex["1"].(string)
			g.Expect(ok).To(BeTrue(), "file path should be string for %s", expectedType)
			g.Expect(filePath).To(MatchRegexp(`^\$\."1"\[\d+\]$`), "file path format for %s", expectedType)

			// Verify type path (field "2")
			typePath, ok := typeIndex["2"].(string)
			g.Expect(ok).To(BeTrue(), "type path should be string for %s", expectedType)
			g.Expect(typePath).To(ContainSubstring(filePath), "type path should contain file path for %s", expectedType)

			// Verify field name index (field "3")
			_, ok = typeIndex["3"].(map[string]interface{})
			g.Expect(ok).To(BeTrue(), "field name index should be map for %s", expectedType)

			// Verify field number index (field "4")
			_, ok = typeIndex["4"].(map[string]interface{})
			g.Expect(ok).To(BeTrue(), "field number index should be map for %s", expectedType)
		}

		// Verify some expected enum types exist in the enum index
		expectedEnumTypes := []string{
			".google.protobuf.FieldDescriptorProto.Type",
			".google.protobuf.FieldDescriptorProto.Label",
		}

		for _, expectedType := range expectedEnumTypes {
			g.Expect(enumTypeIndexData).To(HaveKey(expectedType), "missing enum type: %s", expectedType)

			// Verify EnumTypeIndex structure (protonumberjson format)
			typeIndex, ok := enumTypeIndexData[expectedType].(map[string]interface{})
			g.Expect(ok).To(BeTrue(), "enum type index should be map for %s", expectedType)

			// Verify file path (field "1")
			filePath, ok := typeIndex["1"].(string)
			g.Expect(ok).To(BeTrue(), "file path should be string for %s", expectedType)
			g.Expect(filePath).To(MatchRegexp(`^\$\."1"\[\d+\]$`), "file path format for %s", expectedType)

			// Verify type path (field "2")
			typePath, ok := typeIndex["2"].(string)
			g.Expect(ok).To(BeTrue(), "type path should be string for %s", expectedType)
			g.Expect(typePath).To(ContainSubstring(filePath), "type path should contain file path for %s", expectedType)

			// Verify enum name index (field "3")
			_, ok = typeIndex["3"].(map[string]interface{})
			g.Expect(ok).To(BeTrue(), "enum name index should be map for %s", expectedType)

			// Verify enum number index (field "4")
			_, ok = typeIndex["4"].(map[string]interface{})
			g.Expect(ok).To(BeTrue(), "enum number index should be map for %s", expectedType)
		}
	})

	t.Run("with empty FileDescriptorSet", func(t *testing.T) {
		fileDescriptorSet := &descriptorpb.FileDescriptorSet{
			File: []*descriptorpb.FileDescriptorProto{},
		}

		jsonStr, err := ToJson(fileDescriptorSet)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(jsonStr).ToNot(BeEmpty())

		// Verify it's valid JSON
		var result interface{}
		err = json.Unmarshal([]byte(jsonStr), &result)
		g.Expect(err).ToNot(HaveOccurred())

		// Verify it's a DescriptorSet message object with protonumberjson format
		resultMap, ok := result.(map[string]interface{})
		g.Expect(ok).To(BeTrue())

		// Verify field "1" is the FileDescriptorSet
		fileDescriptorSetData := resultMap["1"]
		g.Expect(fileDescriptorSetData).ToNot(BeNil())

		// Verify field "2" (message type index) is empty or doesn't exist
		messageTypeIndexData, exists := resultMap["2"]
		if exists {
			messageTypeIndexMap, ok := messageTypeIndexData.(map[string]interface{})
			g.Expect(ok).To(BeTrue())
			g.Expect(messageTypeIndexMap).To(BeEmpty())
		}

		// Verify field "3" (enum type index) is empty or doesn't exist
		enumTypeIndexData, exists := resultMap["3"]
		if exists {
			enumTypeIndexMap, ok := enumTypeIndexData.(map[string]interface{})
			g.Expect(ok).To(BeTrue())
			g.Expect(enumTypeIndexMap).To(BeEmpty())
		}
	})

	t.Run("with nil FileDescriptorSet", func(t *testing.T) {
		jsonStr, err := ToJson(nil)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("cannot be nil"))
		g.Expect(jsonStr).To(BeEmpty())
	})
}

func TestToJsonTree(t *testing.T) {
	g := NewWithT(t)

	t.Run("with simple FileDescriptorSet", func(t *testing.T) {
		// Create a simple FileDescriptorSet with one message
		fileDescriptorSet := &descriptorpb.FileDescriptorSet{
			File: []*descriptorpb.FileDescriptorProto{
				{
					Name:    proto.String("test.proto"),
					Package: proto.String("test"),
					MessageType: []*descriptorpb.DescriptorProto{
						{
							Name: proto.String("TestMessage"),
							Field: []*descriptorpb.FieldDescriptorProto{
								{
									Name:   proto.String("field1"),
									Number: proto.Int32(1),
									Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
								},
							},
						},
					},
				},
			},
		}

		result, err := ToJsonTree(fileDescriptorSet)
		g.Expect(err).ToNot(HaveOccurred())

		// Verify it's a map (DescriptorSet message)
		resultMap, ok := result.(map[string]interface{})
		g.Expect(ok).To(BeTrue())

		// Verify field "1" is the FileDescriptorSet
		fileDescriptorSetData := resultMap["1"]
		g.Expect(fileDescriptorSetData).ToNot(BeNil())

		// Verify field "2" is the message type index
		messageTypeIndexData, ok := resultMap["2"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(messageTypeIndexData).To(HaveKey(".test.TestMessage"))

		// Verify the TestMessage type index
		testMessageIndex, ok := messageTypeIndexData[".test.TestMessage"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(testMessageIndex["1"]).To(Equal("$.\"1\"[0]"))          // file_json_path
		g.Expect(testMessageIndex["2"]).To(Equal("$.\"1\"[0].\"4\"[0]")) // type_json_path

		// Verify field indexes exist
		g.Expect(testMessageIndex).To(HaveKey("3")) // field_name_index
		g.Expect(testMessageIndex).To(HaveKey("4")) // field_number_index

		// Verify field "3" is the enum type index (may be empty)
		if enumTypeData, exists := resultMap["3"]; exists {
			_, ok = enumTypeData.(map[string]interface{})
			g.Expect(ok).To(BeTrue())
		}
	})

	t.Run("with nested types", func(t *testing.T) {
		// Create a FileDescriptorSet with nested types
		fileDescriptorSet := &descriptorpb.FileDescriptorSet{
			File: []*descriptorpb.FileDescriptorProto{
				{
					Name:    proto.String("nested.proto"),
					Package: proto.String("nested"),
					MessageType: []*descriptorpb.DescriptorProto{
						{
							Name: proto.String("OuterMessage"),
							NestedType: []*descriptorpb.DescriptorProto{
								{
									Name: proto.String("InnerMessage"),
								},
							},
							EnumType: []*descriptorpb.EnumDescriptorProto{
								{
									Name: proto.String("InnerEnum"),
									Value: []*descriptorpb.EnumValueDescriptorProto{
										{
											Name:   proto.String("VALUE1"),
											Number: proto.Int32(0),
										},
									},
								},
							},
						},
					},
				},
			},
		}

		result, err := ToJsonTree(fileDescriptorSet)
		g.Expect(err).ToNot(HaveOccurred())

		resultMap, ok := result.(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		typeIndexData, ok := resultMap["2"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())

		// Verify outer message
		g.Expect(typeIndexData).To(HaveKey(".nested.OuterMessage"))
		outerIndex, ok := typeIndexData[".nested.OuterMessage"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(outerIndex["1"]).To(Equal("$.\"1\"[0]"))          // file_json_path
		g.Expect(outerIndex["2"]).To(Equal("$.\"1\"[0].\"4\"[0]")) // type_json_path

		// Verify nested message
		g.Expect(typeIndexData).To(HaveKey(".nested.OuterMessage.InnerMessage"))
		innerMsgIndex, ok := typeIndexData[".nested.OuterMessage.InnerMessage"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(innerMsgIndex["1"]).To(Equal("$.\"1\"[0]"))                   // file_json_path
		g.Expect(innerMsgIndex["2"]).To(Equal("$.\"1\"[0].\"4\"[0].\"3\"[0]")) // type_json_path

		// Verify field "3" is the enum type index and contains nested enum
		enumTypeIndexData, ok := resultMap["3"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(enumTypeIndexData).To(HaveKey(".nested.OuterMessage.InnerEnum"))
		innerEnumIndex, ok := enumTypeIndexData[".nested.OuterMessage.InnerEnum"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(innerEnumIndex["1"]).To(Equal("$.\"1\"[0]"))                   // file_json_path
		g.Expect(innerEnumIndex["2"]).To(Equal("$.\"1\"[0].\"4\"[0].\"4\"[0]")) // type_json_path
	})

	t.Run("with enum types", func(t *testing.T) {
		fileDescriptorSet := &descriptorpb.FileDescriptorSet{
			File: []*descriptorpb.FileDescriptorProto{
				{
					Name:    proto.String("enum.proto"),
					Package: proto.String("enums"),
					EnumType: []*descriptorpb.EnumDescriptorProto{
						{
							Name: proto.String("Status"),
							Value: []*descriptorpb.EnumValueDescriptorProto{
								{
									Name:   proto.String("ACTIVE"),
									Number: proto.Int32(1),
								},
								{
									Name:   proto.String("INACTIVE"),
									Number: proto.Int32(2),
								},
							},
						},
					},
				},
			},
		}

		result, err := ToJsonTree(fileDescriptorSet)
		g.Expect(err).ToNot(HaveOccurred())

		resultMap, ok := result.(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		enumTypeIndexData, ok := resultMap["3"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())

		// Verify enum type
		g.Expect(enumTypeIndexData).To(HaveKey(".enums.Status"))
		statusIndex, ok := enumTypeIndexData[".enums.Status"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(statusIndex["1"]).To(Equal("$.\"1\"[0]"))          // file_json_path
		g.Expect(statusIndex["2"]).To(Equal("$.\"1\"[0].\"5\"[0]")) // type_json_path

		// Verify enum indexes exist
		g.Expect(statusIndex).To(HaveKey("3")) // enum_name_index
		g.Expect(statusIndex).To(HaveKey("4")) // enum_number_index
	})

	t.Run("with multiple files", func(t *testing.T) {
		fileDescriptorSet := &descriptorpb.FileDescriptorSet{
			File: []*descriptorpb.FileDescriptorProto{
				{
					Name:    proto.String("file1.proto"),
					Package: proto.String("pkg1"),
					MessageType: []*descriptorpb.DescriptorProto{
						{
							Name: proto.String("Message1"),
						},
					},
				},
				{
					Name:    proto.String("file2.proto"),
					Package: proto.String("pkg2"),
					MessageType: []*descriptorpb.DescriptorProto{
						{
							Name: proto.String("Message2"),
						},
					},
				},
			},
		}

		result, err := ToJsonTree(fileDescriptorSet)
		g.Expect(err).ToNot(HaveOccurred())

		resultMap, ok := result.(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		typeIndexData, ok := resultMap["2"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())

		// Verify both messages exist with different file paths
		g.Expect(typeIndexData).To(HaveKey(".pkg1.Message1"))
		g.Expect(typeIndexData).To(HaveKey(".pkg2.Message2"))

		message1Index, ok := typeIndexData[".pkg1.Message1"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		message2Index, ok := typeIndexData[".pkg2.Message2"].(map[string]interface{})
		g.Expect(ok).To(BeTrue())

		g.Expect(message1Index["1"]).To(Equal("$.\"1\"[0]"))          // file_json_path - First file
		g.Expect(message1Index["2"]).To(Equal("$.\"1\"[0].\"4\"[0]")) // type_json_path - First file, first message
		g.Expect(message2Index["1"]).To(Equal("$.\"1\"[1]"))          // file_json_path - Second file
		g.Expect(message2Index["2"]).To(Equal("$.\"1\"[1].\"4\"[0]")) // type_json_path - Second file, first message
	})

	t.Run("with nil FileDescriptorSet", func(t *testing.T) {
		result, err := ToJsonTree(nil)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("cannot be nil"))
		g.Expect(result).To(BeNil())
	})
}

func TestBuildTypeName(t *testing.T) {
	g := NewWithT(t)

	tests := []struct {
		packageName string
		typeName    string
		expected    string
	}{
		{"", "Message", ".Message"},
		{"pkg", "Message", ".pkg.Message"},
		{"com.example", "Message", ".com.example.Message"},
	}

	for _, test := range tests {
		result := buildTypeName(test.packageName, test.typeName)
		g.Expect(result).To(Equal(test.expected))
	}
}
