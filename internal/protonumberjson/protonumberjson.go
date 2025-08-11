package protonumberjson

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/eiiches/mysql-protobuf-functions/internal/protoreflectutils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ToJsonTree converts a protobuf message to a JSON tree structure using field numbers as keys
func ToJsonTree(m proto.Message) (interface{}, error) {
	if m == nil {
		return nil, nil
	}

	return marshalMessage(m.ProtoReflect())
}

// Marshal serializes a protobuf message to JSON using field numbers as keys
func Marshal(m proto.Message) ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	result, err := ToJsonTree(m)
	if err != nil {
		return nil, err
	}

	return json.Marshal(result)
}

func marshalMessage(msg protoreflect.Message) (interface{}, error) {
	if msg == nil {
		return nil, nil
	}

	result := make(map[string]interface{})
	fields := msg.Descriptor().Fields()

	for field := range protoreflectutils.Iterate(fields) {
		if field.HasPresence() && !msg.Has(field) {
			continue
		}

		jsonKey := strconv.Itoa(int(field.Number()))
		value := msg.Get(field)

		jsonValue, err := marshalFieldValue(value, field)
		if err != nil {
			return nil, fmt.Errorf("marshaling field %d: %w", field.Number(), err)
		}

		result[jsonKey] = jsonValue
	}

	return result, nil
}

func marshalFieldValue(value protoreflect.Value, field protoreflect.FieldDescriptor) (interface{}, error) {
	switch {
	case field.IsMap():
		return marshalMapField(value.Map(), field)
	case field.IsList():
		return marshalListField(value.List(), field)
	default:
		return marshalSingularField(value, field)
	}
}

func marshalListField(list protoreflect.List, field protoreflect.FieldDescriptor) ([]interface{}, error) {
	result := make([]interface{}, 0, list.Len())

	for value := range protoreflectutils.Iterate(list) {
		jsonValue, err := marshalSingularField(value, field)
		if err != nil {
			return nil, fmt.Errorf("marshaling list element: %w", err)
		}
		result = append(result, jsonValue)
	}

	return result, nil
}

func marshalMapField(mapValue protoreflect.Map, field protoreflect.FieldDescriptor) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var err error
	mapValue.Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
		keyStr := marshalMapKey(key, field.MapKey().Kind())

		var jsonValue interface{}
		jsonValue, err = marshalSingularField(value, field.MapValue())
		if err != nil {
			err = fmt.Errorf("marshaling map value for key %s: %w", keyStr, err)
			return false
		}

		result[keyStr] = jsonValue
		return true
	})

	return result, err
}

func marshalMapKey(key protoreflect.MapKey, kind protoreflect.Kind) string {
	switch kind {
	case protoreflect.BoolKind:
		return strconv.FormatBool(key.Bool())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return strconv.FormatInt(key.Int(), 10)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return strconv.FormatUint(key.Uint(), 10)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return strconv.FormatInt(key.Int(), 10)
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return strconv.FormatUint(key.Uint(), 10)
	case protoreflect.StringKind:
		return key.String()
	case protoreflect.EnumKind, protoreflect.FloatKind, protoreflect.DoubleKind,
		protoreflect.BytesKind, protoreflect.MessageKind, protoreflect.GroupKind:
		// These types cannot be map keys in protobuf
		return key.String()
	default:
		return key.String()
	}
}

func marshalSingularField(value protoreflect.Value, field protoreflect.FieldDescriptor) (interface{}, error) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return value.Bool(), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		val := value.Int()
		if val > int64(^uint32(0)>>1) || val < int64(-1<<31) {
			return nil, fmt.Errorf("int32 overflow: %d", val)
		}
		return int32(val), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		val := value.Uint()
		if val > uint64(^uint32(0)) {
			return nil, fmt.Errorf("uint32 overflow: %d", val)
		}
		return uint32(val), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return value.Int(), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return value.Uint(), nil
	case protoreflect.FloatKind:
		return float32(value.Float()), nil
	case protoreflect.DoubleKind:
		return value.Float(), nil
	case protoreflect.StringKind:
		return value.String(), nil
	case protoreflect.BytesKind:
		return base64.StdEncoding.EncodeToString(value.Bytes()), nil
	case protoreflect.EnumKind:
		// For enum values, we just return the number for simplicity
		return int32(value.Enum()), nil
	case protoreflect.MessageKind:
		return marshalMessage(value.Message())
	case protoreflect.GroupKind:
		// Groups are deprecated, but handle them like messages
		return marshalMessage(value.Message())
	default:
		return nil, fmt.Errorf("unsupported kind: %v", field.Kind())
	}
}
