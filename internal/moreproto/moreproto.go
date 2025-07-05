package moreproto

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// EqualOrClose compares two proto messages with a custom floating-point comparison function.
// It works like proto.Equal but uses the provided equalOrCloseFn for float32/float64 comparisons.
// The floatSize parameter indicates whether the values are 32-bit or 64-bit floats.
func EqualOrClose(actual, expected proto.Message, equalOrCloseFn func(a, b float64, floatSize int) bool) bool {
	return equalWithTolerance(actual.ProtoReflect(), expected.ProtoReflect(), equalOrCloseFn)
}

func equalWithTolerance(actual, expected protoreflect.Message, equalOrCloseFn func(a, b float64, floatSize int) bool) bool {
	actualDesc := actual.Descriptor()
	expectedDesc := expected.Descriptor()

	if actualDesc.FullName() != expectedDesc.FullName() {
		return false
	}

	// Check all fields in expected message
	expectedFields := expected.Descriptor().Fields()
	for i := 0; i < expectedFields.Len(); i++ {
		field := expectedFields.Get(i)

		expectedValue := expected.Get(field)
		actualValue := actual.Get(field)

		if !equalValueWithTolerance(actualValue, expectedValue, field, equalOrCloseFn) {
			return false
		}
	}

	// Check all fields in actual message to ensure no extra fields
	actualFields := actual.Descriptor().Fields()
	for i := 0; i < actualFields.Len(); i++ {
		field := actualFields.Get(i)

		// If actual has a field that expected doesn't have, they're different
		if actual.Has(field) && !expected.Has(field) {
			return false
		}
	}

	// Check for unknown fields - these should also be equal
	actualUnknown := actual.GetUnknown()
	expectedUnknown := expected.GetUnknown()

	// Compare unknown fields by their byte representation
	if len(actualUnknown) != len(expectedUnknown) {
		return false
	}
	for i := 0; i < len(actualUnknown); i++ {
		if actualUnknown[i] != expectedUnknown[i] {
			return false
		}
	}

	return true
}

func equalValueWithTolerance(actual, expected protoreflect.Value, field protoreflect.FieldDescriptor, equalOrCloseFn func(a, b float64, floatSize int) bool) bool {
	if field.IsList() {
		return equalListWithTolerance(actual.List(), expected.List(), field, equalOrCloseFn)
	}

	if field.IsMap() {
		return equalMapWithTolerance(actual.Map(), expected.Map(), field, equalOrCloseFn)
	}

	return equalScalarWithTolerance(actual, expected, field.Kind(), equalOrCloseFn)
}

func equalListWithTolerance(actual, expected protoreflect.List, field protoreflect.FieldDescriptor, equalOrCloseFn func(a, b float64, floatSize int) bool) bool {
	if actual.Len() != expected.Len() {
		return false
	}

	for i := 0; i < actual.Len(); i++ {
		actualValue := actual.Get(i)
		expectedValue := expected.Get(i)

		if field.Kind() == protoreflect.MessageKind {
			if !equalWithTolerance(actualValue.Message(), expectedValue.Message(), equalOrCloseFn) {
				return false
			}
		} else {
			if !equalScalarWithTolerance(actualValue, expectedValue, field.Kind(), equalOrCloseFn) {
				return false
			}
		}
	}

	return true
}

func equalMapWithTolerance(actual, expected protoreflect.Map, field protoreflect.FieldDescriptor, equalOrCloseFn func(a, b float64, floatSize int) bool) bool {
	if actual.Len() != expected.Len() {
		return false
	}

	equal := true
	expected.Range(func(key protoreflect.MapKey, expectedValue protoreflect.Value) bool {
		actualValue := actual.Get(key)
		if !actualValue.IsValid() {
			equal = false
			return false
		}

		if field.MapValue().Kind() == protoreflect.MessageKind {
			if !equalWithTolerance(actualValue.Message(), expectedValue.Message(), equalOrCloseFn) {
				equal = false
				return false
			}
		} else {
			if !equalScalarWithTolerance(actualValue, expectedValue, field.MapValue().Kind(), equalOrCloseFn) {
				equal = false
				return false
			}
		}
		return true
	})

	return equal
}

func equalScalarWithTolerance(actual, expected protoreflect.Value, kind protoreflect.Kind, equalOrCloseFn func(a, b float64, floatSize int) bool) bool {
	switch kind {
	case protoreflect.FloatKind:
		actualFloat := float64(actual.Float())
		expectedFloat := float64(expected.Float())
		return equalOrCloseFn(actualFloat, expectedFloat, 32)
	case protoreflect.DoubleKind:
		actualDouble := actual.Float()
		expectedDouble := expected.Float()
		return equalOrCloseFn(actualDouble, expectedDouble, 64)
	case protoreflect.MessageKind:
		return equalWithTolerance(actual.Message(), expected.Message(), equalOrCloseFn)
	default:
		// For non-floating point types, use exact equality
		return actual.Equal(expected)
	}
}
