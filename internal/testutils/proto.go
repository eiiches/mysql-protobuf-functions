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
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

// CompositeResolver combines a local resolver with the global registry
type CompositeResolver struct {
	Local interface {
		FindMessageByName(name protoreflect.FullName) (protoreflect.MessageType, error)
		FindMessageByURL(url string) (protoreflect.MessageType, error)
		FindExtensionByName(name protoreflect.FullName) (protoreflect.ExtensionType, error)
		FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error)
	}
	Global interface {
		FindMessageByName(name protoreflect.FullName) (protoreflect.MessageType, error)
		FindMessageByURL(url string) (protoreflect.MessageType, error)
		FindExtensionByName(name protoreflect.FullName) (protoreflect.ExtensionType, error)
		FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error)
	}
}

func (r *CompositeResolver) FindMessageByName(name protoreflect.FullName) (protoreflect.MessageType, error) {
	// Try local resolver first
	if msgType, err := r.Local.FindMessageByName(name); err == nil {
		return msgType, nil
	}
	// Fall back to global registry
	return r.Global.FindMessageByName(name)
}

func (r *CompositeResolver) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	// Try local resolver first
	if msgType, err := r.Local.FindMessageByURL(url); err == nil {
		return msgType, nil
	}
	// Fall back to global registry
	return r.Global.FindMessageByURL(url)
}

func (r *CompositeResolver) FindExtensionByName(name protoreflect.FullName) (protoreflect.ExtensionType, error) {
	// Try local resolver first
	if extType, err := r.Local.FindExtensionByName(name); err == nil {
		return extType, nil
	}
	// Fall back to global registry
	return r.Global.FindExtensionByName(name)
}

func (r *CompositeResolver) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	// Try local resolver first
	if extType, err := r.Local.FindExtensionByNumber(message, field); err == nil {
		return extType, nil
	}
	// Fall back to global registry
	return r.Global.FindExtensionByNumber(message, field)
}

type ProtoTestSupport struct {
	T     *testing.T
	Files linker.Files
	G     Gomega
}

func NewProtoTestSupport(t *testing.T, sources map[string]string) *ProtoTestSupport {
	t.Helper()
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

	// Use UnmarshalOptions with a composite resolver that includes both local types and well-known types
	opts := protojson.UnmarshalOptions{
		Resolver: &CompositeResolver{
			Local:  this.Files.AsResolver(),
			Global: protoregistry.GlobalTypes,
		},
	}

	this.G.Expect(opts.Unmarshal([]byte(json), dynamicMessage.Interface())).To(Succeed())

	return dynamicMessage
}

func (this *ProtoTestSupport) GetMessageDescriptor(name protoreflect.FullName) protoreflect.MessageDescriptor {
	return this.GetMessageType(name).Descriptor()
}

func (this *ProtoTestSupport) GetMessageType(name protoreflect.FullName) protoreflect.MessageType {
	messageType, err := this.Files.AsResolver().FindMessageByName(name)
	this.G.Expect(err).NotTo(HaveOccurred())
	return messageType
}

func (this *ProtoTestSupport) JsonToProtobuf(name protoreflect.FullName, json string) []byte {
	dynamicMessage := this.JsonToDynamicMessage(name, json)

	serialized, err := proto.Marshal(dynamicMessage.Interface())
	this.G.Expect(err).NotTo(HaveOccurred())

	return serialized
}
