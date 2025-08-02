package main

import (
	"testing"

	"github.com/eiiches/mysql-protobuf-functions/internal/descriptorsetjson"
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
	test := func(fileDescriptor protoreflect.FileDescriptor) {
		name := fileDescriptor.Path()
		t.Run(name, func(t *testing.T) {
			g := NewWithT(t)

			// Build FileDescriptorSet with dependencies
			fileDescriptorSet := protoreflectutils.BuildFileDescriptorSetWithDependencies(fileDescriptor)
			fileDescriptorSetBytes, err := proto.Marshal(fileDescriptorSet)
			g.Expect(err).NotTo(HaveOccurred())

			// Get expected JSON from Go package
			expectedJSON, err := descriptorsetjson.ToJson(fileDescriptorSet)
			g.Expect(err).NotTo(HaveOccurred())

			// Test that MySQL function output matches expected JSON
			RunTestThatExpression(t, "pb_build_descriptor_set_json(?)", fileDescriptorSetBytes).
				IsEqualToJsonString(expectedJSON)
		})
	}

	// Test all Well-Known Types
	test(anypb.File_google_protobuf_any_proto)
	test(apipb.File_google_protobuf_api_proto)
	test(durationpb.File_google_protobuf_duration_proto)
	test(descriptorpb.File_google_protobuf_descriptor_proto)
	test(emptypb.File_google_protobuf_empty_proto)
	test(fieldmaskpb.File_google_protobuf_field_mask_proto)
	test(sourcecontextpb.File_google_protobuf_source_context_proto)
	test(structpb.File_google_protobuf_struct_proto)
	test(timestamppb.File_google_protobuf_timestamp_proto)
	test(typepb.File_google_protobuf_type_proto)
	test(wrapperspb.File_google_protobuf_wrappers_proto)
}
