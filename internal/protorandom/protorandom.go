package protorandom

import (
	"math"
	"math/rand"

	"github.com/eiiches/mysql-protobuf-functions/internal/protoreflectutils"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func Float(rng *rand.Rand, allowNan bool, allowInf bool) float32 {
	for {
		value := math.Float32frombits(rng.Uint32())
		if !allowNan && math.IsNaN(float64(value)) {
			continue
		}
		if !allowInf && math.IsInf(float64(value), 0) {
			continue
		}
		return value
	}
}

func Double(rng *rand.Rand, allowNan bool, allowInf bool) float64 {
	for {
		value := math.Float64frombits(rng.Uint64())
		if !allowNan && math.IsNaN(value) {
			continue
		}
		if !allowInf && math.IsInf(value, 0) {
			continue
		}
		return value
	}
}

func Int32(rng *rand.Rand) int32 {
	return int32(rng.Uint32()) //nolint:gosec // Intentional overflow for random data generation
}

func Int64(rng *rand.Rand) int64 {
	return int64(rng.Uint64()) //nolint:gosec // Intentional overflow for random data generation
}

func Uint32(rng *rand.Rand) uint32 {
	return rng.Uint32()
}

func Uint64(rng *rand.Rand) uint64 {
	return rng.Uint64()
}

func Bool(rng *rand.Rand) bool {
	return rng.Intn(2) == 1
}

func String(rng *rand.Rand, maxLength int) string {
	length := rng.Intn(maxLength) + 1 // Random length between 1 and 50
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 !@#$%^&*()_+-="
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rng.Intn(len(chars))]
	}
	return string(result)
}

func Bytes(rng *rand.Rand, maxLength int) []byte {
	length := rng.Intn(maxLength) + 1 // Random length between 1 and 50
	result := make([]byte, length)
	if _, err := rng.Read(result); err != nil {
		panic(err)
	}
	return result
}

func Enum(rng *rand.Rand, descriptor protoreflect.EnumDescriptor) protoreflect.EnumNumber {
	if descriptor.Values().Len() == 0 {
		return 0 // No values to choose from
	}
	index := rng.Intn(descriptor.Values().Len())
	return descriptor.Values().Get(index).Number()
}

func SingleValueForField(rng *rand.Rand, fieldDescriptor protoreflect.FieldDescriptor, config *Config) protoreflect.Value {
	switch fieldDescriptor.Kind() {
	case protoreflect.FloatKind:
		return protoreflect.ValueOfFloat32(Float(rng, config.GetAllowNan(), config.GetAllowInf()))
	case protoreflect.DoubleKind:
		return protoreflect.ValueOfFloat64(Double(rng, config.GetAllowNan(), config.GetAllowInf()))
	case protoreflect.Int32Kind:
		return protoreflect.ValueOfInt32(Int32(rng))
	case protoreflect.Int64Kind:
		return protoreflect.ValueOfInt64(Int64(rng))
	case protoreflect.Uint32Kind:
		return protoreflect.ValueOfUint32(Uint32(rng))
	case protoreflect.Uint64Kind:
		return protoreflect.ValueOfUint64(Uint64(rng))
	case protoreflect.BoolKind:
		return protoreflect.ValueOfBool(Bool(rng))
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(String(rng, config.GetMaxStringLength()))
	case protoreflect.BytesKind:
		return protoreflect.ValueOfBytes(Bytes(rng, config.GetMaxBytesLength()))
	case protoreflect.Sint32Kind:
		return protoreflect.ValueOfInt32(Int32(rng))
	case protoreflect.Sint64Kind:
		return protoreflect.ValueOfInt64(Int64(rng))
	case protoreflect.Fixed32Kind:
		return protoreflect.ValueOfUint32(Uint32(rng))
	case protoreflect.Sfixed32Kind:
		return protoreflect.ValueOfInt32(Int32(rng))
	case protoreflect.Fixed64Kind:
		return protoreflect.ValueOfUint64(Uint64(rng))
	case protoreflect.Sfixed64Kind:
		return protoreflect.ValueOfInt64(Int64(rng))
	case protoreflect.EnumKind:
		nestedEnum := Enum(rng, fieldDescriptor.Enum())
		return protoreflect.ValueOfEnum(nestedEnum)
	case protoreflect.MessageKind:
		nestedMessage := Message(rng, fieldDescriptor.Message(), config)
		return protoreflect.ValueOfMessage(nestedMessage)
	case protoreflect.GroupKind:
		panic("Groups are not supported")
	default:
		panic("Unsupported field kind: " + fieldDescriptor.Kind().String())
	}
}

func Message(rng *rand.Rand, descriptor protoreflect.MessageDescriptor, config *Config) protoreflect.Message {
	message := dynamicpb.NewMessage(descriptor)

	for fieldDescriptor := range protoreflectutils.Iterate(descriptor.Fields()) {
		switch {
		case fieldDescriptor.IsMap():
			length := rng.Intn(config.GetMaxMapSize() + 1) // Randomly choose a map size between 0 and 3
			for i := 0; i < length; i++ {
				key := SingleValueForField(rng, fieldDescriptor.MapKey(), config)
				value := SingleValueForField(rng, fieldDescriptor.MapValue(), config)
				message.Mutable(fieldDescriptor).Map().Set(key.MapKey(), value)
			}
		case fieldDescriptor.IsList():
			length := rng.Intn(config.GetMaxRepeatedSize() + 1) // Randomly choose a list length between 0 and 3
			for i := 0; i < length; i++ {
				value := SingleValueForField(rng, fieldDescriptor, config)
				message.Mutable(fieldDescriptor).List().Append(value)
			}
		default:
			if rng.Intn(2) == 1 { // Randomly decide whether to set the field
				value := SingleValueForField(rng, fieldDescriptor, config)
				message.Set(fieldDescriptor, value)
			}
		}
	}

	return message
}

type Config struct {
	AllowNan        *bool
	AllowInf        *bool
	MaxStringLength *int
	MaxBytesLength  *int
	MaxRepeatedSize *int
	MaxMapSize      *int
}

func (config *Config) GetMaxStringLength() int {
	if config != nil && config.MaxStringLength != nil {
		return *config.MaxStringLength
	}
	return 5
}

func (config *Config) GetMaxBytesLength() int {
	if config != nil && config.MaxBytesLength != nil {
		return *config.MaxBytesLength
	}
	return 5
}

func (config *Config) GetMaxRepeatedSize() int {
	if config != nil && config.MaxRepeatedSize != nil {
		return *config.MaxRepeatedSize
	}
	return 3
}

func (config *Config) GetMaxMapSize() int {
	if config != nil && config.MaxMapSize != nil {
		return *config.MaxMapSize
	}
	return 3
}

func (config *Config) GetAllowNan() bool {
	if config != nil && config.AllowNan != nil {
		return *config.AllowNan
	}
	return true
}

func (config *Config) GetAllowInf() bool {
	if config != nil && config.AllowInf != nil {
		return *config.AllowInf
	}
	return true
}
