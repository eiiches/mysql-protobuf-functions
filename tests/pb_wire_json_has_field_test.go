package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/eiiches/mysql-protobuf-functions/internal/protorandom"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestRandomizedWireJsonHasField(t *testing.T) {
	test := func(t *testing.T, protoFieldType string, hasFunction string) {
		t.Run(protoFieldType, func(t *testing.T) {
			GivenFieldDefinitions(t, fmt.Sprintf("int32 a = 1; %s value = 2; int32 b = 3;", protoFieldType), func(messageType protoreflect.MessageType) {
				valueField := messageType.Descriptor().Fields().ByName("value")

				seed := time.Now().UnixNano()
				t.Logf("Using seed = %d.", seed)
				rng := rand.New(rand.NewSource(seed))
				for i := 0; i < iterations; i++ {
					input := protorandom.Message(rng, messageType.Descriptor(), nil)
					expectedHasValue := input.Has(valueField)

					RunTestThatExpression(t, fmt.Sprintf("%s(pb_message_to_wire_json(?), 2)", strings.ReplaceAll(hasFunction, "{kind}", "wire_json")), input.Interface()).
						IsEqualTo(expectedHasValue)

					RunTestThatExpression(t, fmt.Sprintf("%s(?, 2)", strings.ReplaceAll(hasFunction, "{kind}", "message")), input.Interface()).
						IsEqualTo(expectedHasValue)
				}
			})
		})
	}

	test(t, "optional float", "pb_{kind}_has_float_field")
	test(t, "optional double", "pb_{kind}_has_double_field")
	test(t, "optional int32", "pb_{kind}_has_int32_field")
	test(t, "optional int64", "pb_{kind}_has_int64_field")
	test(t, "optional uint32", "pb_{kind}_has_uint32_field")
	test(t, "optional uint64", "pb_{kind}_has_uint64_field")
	test(t, "optional bool", "pb_{kind}_has_bool_field")
	test(t, "optional string", "pb_{kind}_has_string_field")
	test(t, "optional bytes", "pb_{kind}_has_bytes_field")
	test(t, "optional sint32", "pb_{kind}_has_sint32_field")
	test(t, "optional sint64", "pb_{kind}_has_sint64_field")
	test(t, "optional fixed32", "pb_{kind}_has_fixed32_field")
	test(t, "optional fixed64", "pb_{kind}_has_fixed64_field")
	test(t, "optional sfixed32", "pb_{kind}_has_sfixed32_field")
	test(t, "optional sfixed64", "pb_{kind}_has_sfixed64_field")
	test(t, "optional EnumType", "pb_{kind}_has_enum_field")
	test(t, "optional MessageType", "pb_{kind}_has_message_field")

	// In proto3, field presence for message fields is always tracked.
	test(t, "MessageType", "pb_{kind}_has_message_field")
}
