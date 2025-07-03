package testutils

import (
	"context"
	"maps"
	"slices"
	"testing"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"github.com/bufbuild/protocompile/wellknownimports"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type ProtoTestSupport struct {
	T     *testing.T
	Files linker.Files
	G     Gomega
}

func NewProtoTestSupport(t *testing.T, sources map[string]string) *ProtoTestSupport {
	g := NewWithT(t)

	compiler := protocompile.Compiler{
		Resolver: wellknownimports.WithStandardImports(&protocompile.SourceResolver{
			Accessor: protocompile.SourceAccessorFromMap(sources),
		}),
	}

	files, err := compiler.Compile(context.Background(), slices.Collect(maps.Keys(sources))...)
	g.Expect(err).NotTo(HaveOccurred())

	return &ProtoTestSupport{
		T:     t,
		G:     g,
		Files: files,
	}
}

func (this *ProtoTestSupport) GetFileDescriptorSet() *descriptorpb.FileDescriptorSet {
	fds := &descriptorpb.FileDescriptorSet{}
	for _, file := range this.Files {
		fileDescProto := protodesc.ToFileDescriptorProto(file)
		fds.File = append(fds.File, fileDescProto)
	}
	return fds
}

func (this *ProtoTestSupport) GetSerializedFileDescriptorSet() []byte {
	fds := this.GetFileDescriptorSet()
	fdsBytes, err := proto.Marshal(fds)
	if err != nil {
		this.G.Expect(err).NotTo(HaveOccurred())
	}
	return fdsBytes
}

func (this *ProtoTestSupport) JsonToDynamicMessage(name protoreflect.FullName, json string) protoreflect.Message {
	messageType, err := this.Files.AsResolver().FindMessageByName(name)
	this.G.Expect(err).NotTo(HaveOccurred())

	dynamicMessage := messageType.New()
	this.G.Expect(protojson.Unmarshal([]byte(json), dynamicMessage.Interface())).To(Succeed())

	return dynamicMessage
}

func (this *ProtoTestSupport) JsonToProtobuf(name protoreflect.FullName, json string) []byte {
	dynamicMessage := this.JsonToDynamicMessage(name, json)

	serialized, err := proto.Marshal(dynamicMessage.Interface())
	this.G.Expect(err).NotTo(HaveOccurred())

	return serialized
}
