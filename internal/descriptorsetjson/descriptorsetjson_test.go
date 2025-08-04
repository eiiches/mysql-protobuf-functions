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

		// Verify it's a 3-element array
		resultArray, ok := result.([]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(resultArray).To(HaveLen(3))

		// Verify first element is the version number
		version, ok := resultArray[0].(float64)
		g.Expect(ok).To(BeTrue())
		g.Expect(version).To(Equal(float64(1)))

		// Verify second element is the FileDescriptorSet
		fileDescriptorSetData := resultArray[1]
		g.Expect(fileDescriptorSetData).ToNot(BeNil())

		// Verify third element is the type index
		typeIndexData, ok := resultArray[2].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(typeIndexData).ToNot(BeEmpty())

		// Verify some expected types exist in the index
		expectedTypes := []string{
			".google.protobuf.FileDescriptorSet",
			".google.protobuf.FileDescriptorProto",
			".google.protobuf.DescriptorProto",
			".google.protobuf.FieldDescriptorProto",
			".google.protobuf.FieldDescriptorProto.Type",
			".google.protobuf.FieldDescriptorProto.Label",
		}

		for _, expectedType := range expectedTypes {
			g.Expect(typeIndexData).To(HaveKey(expectedType), "missing type: %s", expectedType)

			// Verify TypeIndex structure
			typeIndex, ok := typeIndexData[expectedType].([]interface{})
			g.Expect(ok).To(BeTrue(), "type index should be array for %s", expectedType)
			g.Expect(typeIndex).To(HaveLen(3), "type index should have 3 elements for %s", expectedType)

			// Verify kind (first element)
			kind, ok := typeIndex[0].(float64)
			g.Expect(ok).To(BeTrue(), "kind should be number for %s", expectedType)
			g.Expect(kind).To(SatisfyAny(Equal(float64(11)), Equal(float64(14))), "kind should be 11 (message) or 14 (enum) for %s", expectedType)

			// Verify file path (second element)
			filePath, ok := typeIndex[1].(string)
			g.Expect(ok).To(BeTrue(), "file path should be string for %s", expectedType)
			g.Expect(filePath).To(MatchRegexp(`^\$\[1\]\."1"\[\d+\]$`), "file path format for %s", expectedType)

			// Verify type path (third element)
			typePath, ok := typeIndex[2].(string)
			g.Expect(ok).To(BeTrue(), "type path should be string for %s", expectedType)
			g.Expect(typePath).To(ContainSubstring(filePath), "type path should contain file path for %s", expectedType)
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

		// Verify it's a 3-element array
		resultArray, ok := result.([]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(resultArray).To(HaveLen(3))

		// Verify first element is the version number
		version, ok := resultArray[0].(float64)
		g.Expect(ok).To(BeTrue())
		g.Expect(version).To(Equal(float64(1)))

		// Verify type index is empty
		typeIndexData, ok := resultArray[2].(map[string]interface{})
		g.Expect(ok).To(BeTrue())
		g.Expect(typeIndexData).To(BeEmpty())
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
		g.Expect(result).To(HaveLen(3))

		// Verify version
		version := result[0]
		g.Expect(version).To(Equal(1))

		// Verify FileDescriptorSet part
		fileDescriptorSetData := result[1]
		g.Expect(fileDescriptorSetData).ToNot(BeNil())

		// Verify type index part
		typeIndexData, ok := result[2].(map[string]TypeIndex)
		g.Expect(ok).To(BeTrue())
		g.Expect(typeIndexData).To(HaveKey(".test.TestMessage"))

		// Verify the TestMessage type index
		testMessageIndex := typeIndexData[".test.TestMessage"]
		g.Expect(testMessageIndex[0]).To(Equal(11)) // TYPE_MESSAGE
		g.Expect(testMessageIndex[1]).To(Equal("$[1].\"1\"[0]"))
		g.Expect(testMessageIndex[2]).To(Equal("$[1].\"1\"[0].\"4\"[0]"))
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

		typeIndexData, ok := result[2].(map[string]TypeIndex)
		g.Expect(ok).To(BeTrue())

		// Verify outer message
		g.Expect(typeIndexData).To(HaveKey(".nested.OuterMessage"))
		outerIndex := typeIndexData[".nested.OuterMessage"]
		g.Expect(outerIndex[0]).To(Equal(11)) // TYPE_MESSAGE

		// Verify nested message
		g.Expect(typeIndexData).To(HaveKey(".nested.OuterMessage.InnerMessage"))
		innerMsgIndex := typeIndexData[".nested.OuterMessage.InnerMessage"]
		g.Expect(innerMsgIndex[0]).To(Equal(11)) // TYPE_MESSAGE
		g.Expect(innerMsgIndex[2]).To(Equal("$[1].\"1\"[0].\"4\"[0].\"3\"[0]"))

		// Verify nested enum
		g.Expect(typeIndexData).To(HaveKey(".nested.OuterMessage.InnerEnum"))
		innerEnumIndex := typeIndexData[".nested.OuterMessage.InnerEnum"]
		g.Expect(innerEnumIndex[0]).To(Equal(14)) // TYPE_ENUM
		g.Expect(innerEnumIndex[2]).To(Equal("$[1].\"1\"[0].\"4\"[0].\"4\"[0]"))
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

		typeIndexData, ok := result[2].(map[string]TypeIndex)
		g.Expect(ok).To(BeTrue())

		// Verify enum type
		g.Expect(typeIndexData).To(HaveKey(".enums.Status"))
		statusIndex := typeIndexData[".enums.Status"]
		g.Expect(statusIndex[0]).To(Equal(14)) // TYPE_ENUM
		g.Expect(statusIndex[1]).To(Equal("$[1].\"1\"[0]"))
		g.Expect(statusIndex[2]).To(Equal("$[1].\"1\"[0].\"5\"[0]"))
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

		typeIndexData, ok := result[2].(map[string]TypeIndex)
		g.Expect(ok).To(BeTrue())

		// Verify both messages exist with different file paths
		g.Expect(typeIndexData).To(HaveKey(".pkg1.Message1"))
		g.Expect(typeIndexData).To(HaveKey(".pkg2.Message2"))

		message1Index := typeIndexData[".pkg1.Message1"]
		message2Index := typeIndexData[".pkg2.Message2"]

		g.Expect(message1Index[1]).To(Equal("$[1].\"1\"[0]")) // First file
		g.Expect(message2Index[1]).To(Equal("$[1].\"1\"[1]")) // Second file
	})

	t.Run("with nil FileDescriptorSet", func(t *testing.T) {
		result, err := ToJsonTree(nil)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("cannot be nil"))
		g.Expect(result).To(Equal([3]interface{}{}))
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
