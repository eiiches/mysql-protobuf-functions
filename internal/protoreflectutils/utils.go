package protoreflectutils

import (
	"github.com/eiiches/mysql-protobuf-functions/internal/moremaps"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func BuildFileDescriptorSetWithDependencies(fileDescriptors ...protoreflect.FileDescriptor) *descriptorpb.FileDescriptorSet {
	visited := map[string]protoreflect.FileDescriptor{}

	var recurse func(fileDescriptor protoreflect.FileDescriptor)
	recurse = func(fileDescriptor protoreflect.FileDescriptor) {
		if _, ok := visited[fileDescriptor.Path()]; ok {
			return
		}
		visited[fileDescriptor.Path()] = fileDescriptor
		for importedFile := range Iterate(fileDescriptor.Imports()) {
			recurse(importedFile)
		}
	}

	for _, fileDescriptor := range fileDescriptors {
		recurse(fileDescriptor)
	}

	result := &descriptorpb.FileDescriptorSet{}

	for _, fileDescriptor := range moremaps.SortedEntries(visited) {
		fileDescriptorProto := protodesc.ToFileDescriptorProto(fileDescriptor)
		result.File = append(result.File, fileDescriptorProto)
	}

	return result
}
