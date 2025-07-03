package main

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/protoreflectutils"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/apipb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/sourcecontextpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/typepb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestDescriptorSetLoadWkt(t *testing.T) {
	test := func(fileDescriptor protoreflect.FileDescriptor, fn func(setName string, t *testing.T)) {
		name := fileDescriptor.Path()
		t.Run(name, func(t *testing.T) {
			g := NewWithT(t)
			fileDescriptorSet := protoreflectutils.BuildFileDescriptorSetWithDependencies(fileDescriptor)
			fileDescriptorSetBytes, err := proto.Marshal(fileDescriptorSet)
			g.Expect(err).NotTo(HaveOccurred())
			AssertThatCall(t, "pb_descriptor_set_delete(?)", name).ShouldSucceed()
			AssertThatCall(t, "pb_descriptor_set_load(?, ?)", name, fileDescriptorSetBytes).ShouldSucceed()
			AssertThatExpression(t, "pb_descriptor_set_exists(?)", name).IsTrue()
			fn(name, t)
		})
	}

	test(anypb.File_google_protobuf_any_proto, func(setName string, t *testing.T) {
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Any").IsTrue()
	})

	test(apipb.File_google_protobuf_api_proto, func(setName string, t *testing.T) {})

	test(durationpb.File_google_protobuf_duration_proto, func(setName string, t *testing.T) {
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Duration").IsTrue()
	})

	test(descriptorpb.File_google_protobuf_descriptor_proto, func(setName string, t *testing.T) {
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.FileDescriptorSet").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.DescriptorProto").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.FieldDescriptorProto").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.EnumDescriptorProto").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.EnumValueDescriptorProto").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.ServiceDescriptorProto").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_enum_type(?, ?)", setName, ".google.protobuf.FieldDescriptorProto.Type").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_enum_type(?, ?)", setName, ".google.protobuf.FieldDescriptorProto.Label").IsTrue()
	})

	test(emptypb.File_google_protobuf_empty_proto, func(setName string, t *testing.T) {
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Empty").IsTrue()
	})

	test(fieldmaskpb.File_google_protobuf_field_mask_proto, func(setName string, t *testing.T) {
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.FieldMask").IsTrue()
	})

	test(sourcecontextpb.File_google_protobuf_source_context_proto, func(setName string, t *testing.T) {})

	test(structpb.File_google_protobuf_struct_proto, func(setName string, t *testing.T) {
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Struct").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Value").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.ListValue").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_enum_type(?, ?)", setName, ".google.protobuf.NullValue").IsTrue()
	})

	test(timestamppb.File_google_protobuf_timestamp_proto, func(setName string, t *testing.T) {
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Timestamp").IsTrue()
	})

	test(typepb.File_google_protobuf_type_proto, func(setName string, t *testing.T) {
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Type").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Field").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Enum").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.EnumValue").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Option").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_enum_type(?, ?)", setName, ".google.protobuf.Field.Kind").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_enum_type(?, ?)", setName, ".google.protobuf.Field.Cardinality").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_enum_type(?, ?)", setName, ".google.protobuf.Syntax").IsTrue()
	})

	test(wrapperspb.File_google_protobuf_wrappers_proto, func(setName string, t *testing.T) {
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.DoubleValue").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.FloatValue").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Int64Value").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.UInt64Value").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.Int32Value").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.UInt32Value").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.BoolValue").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.StringValue").IsTrue()
		RunTestThatExpression(t, "pb_descriptor_set_contains_message_type(?, ?)", setName, ".google.protobuf.BytesValue").IsTrue()

		type FieldDescriptor struct {
			FieldName      string
			FieldNumber    int32
			Proto3Optional bool
			OneofIndex     *int32
		}

		RunTestThatStatement[FieldDescriptor](t, "SELECT field_name, field_number, proto3_optional, oneof_index FROM _Proto_FieldDescriptor WHERE set_name = ? AND type_name = ? AND field_number = ?", setName, ".google.protobuf.DoubleValue", 1).
			ShouldReturnSingleRow(Equal(&FieldDescriptor{
				FieldName:      "value",
				FieldNumber:    1,
				Proto3Optional: false,
				OneofIndex:     nil,
			}))
	})
}
