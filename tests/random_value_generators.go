package main

import (
	"math/rand"

	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Generator function type for protobuf values
type ValueGenerator func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value)

// Shared generator functions for different protobuf types
var (
	RandomFloatGenerator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Float(rng, false, false)
		return newValue, protoreflect.ValueOfFloat32(newValue)
	}

	RandomDoubleGenerator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Double(rng, false, false)
		return newValue, protoreflect.ValueOfFloat64(newValue)
	}

	RandomInt32Generator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int32(rng)
		return newValue, protoreflect.ValueOfInt32(newValue)
	}

	RandomInt64Generator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Int64(rng)
		return newValue, protoreflect.ValueOfInt64(newValue)
	}

	RandomUint32Generator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint32(rng)
		return newValue, protoreflect.ValueOfUint32(newValue)
	}

	RandomUint64Generator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Uint64(rng)
		return newValue, protoreflect.ValueOfUint64(newValue)
	}

	RandomBoolGenerator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Bool(rng)
		return newValue, protoreflect.ValueOfBool(newValue)
	}

	RandomStringGenerator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.String(rng, 5)
		return newValue, protoreflect.ValueOfString(newValue)
	}

	RandomBytesGenerator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Bytes(rng, 5)
		return newValue, protoreflect.ValueOfBytes(newValue)
	}

	RandomEnumGenerator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Enum(rng, fieldDescriptor.Enum())
		return newValue, protoreflect.ValueOfEnum(newValue)
	}

	RandomMessageGenerator = func(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor) (interface{}, protoreflect.Value) {
		newValue := protorandom.Message(rng, fieldDescriptor.Message(), nil)
		return newValue.Interface(), protoreflect.ValueOfMessage(newValue)
	}
)
