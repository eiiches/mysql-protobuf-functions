package protonumberjson

import (
	"encoding/json"
	"testing"

	"github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestMarshalBasicTypes(t *testing.T) {
	g := gomega.NewWithT(t)

	// Test Timestamp
	ts := &timestamppb.Timestamp{Seconds: 1234567890, Nanos: 123456789}
	result, err := Marshal(ts)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var timestampStr string
	err = json.Unmarshal(result, &timestampStr)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(timestampStr).To(gomega.Equal("2009-02-13T23:31:30.123456789Z"))

	// Test Duration
	dur := &durationpb.Duration{Seconds: 3600, Nanos: 500000000}
	result, err = Marshal(dur)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var durationStr string
	err = json.Unmarshal(result, &durationStr)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(durationStr).To(gomega.Equal("3600.500s"))

	// Test Empty
	empty := &emptypb.Empty{}
	result, err = Marshal(empty)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(string(result)).To(gomega.Equal("{}"))

	// Test StringValue
	str := &wrapperspb.StringValue{Value: "hello world"}
	result, err = Marshal(str)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var strVal string
	err = json.Unmarshal(result, &strVal)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(strVal).To(gomega.Equal("hello world"))

	// Test Int64Value
	int64Val := &wrapperspb.Int64Value{Value: 9223372036854775807}
	result, err = Marshal(int64Val)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var int64Result float64
	err = json.Unmarshal(result, &int64Result)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(int64Result).To(gomega.Equal(float64(9223372036854775807)))

	// Test BoolValue
	boolVal := &wrapperspb.BoolValue{Value: true}
	result, err = Marshal(boolVal)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var boolResult bool
	err = json.Unmarshal(result, &boolResult)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(boolResult).To(gomega.Equal(true))

	// Test FieldMask
	fm := &fieldmaskpb.FieldMask{Paths: []string{"user.name", "user.email"}}
	result, err = Marshal(fm)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var fmStr string
	err = json.Unmarshal(result, &fmStr)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(fmStr).To(gomega.Equal("user.name,user.email"))
}

func TestMarshalStruct(t *testing.T) {
	g := gomega.NewWithT(t)

	// Create a simple Struct
	structValue, err := structpb.NewStruct(map[string]interface{}{
		"name":   "test",
		"count":  42,
		"active": true,
	})
	g.Expect(err).ToNot(gomega.HaveOccurred())

	result, err := Marshal(structValue)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var jsonObj map[string]interface{}
	err = json.Unmarshal(result, &jsonObj)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	// The struct should be marshaled directly as an object in ProtoJSON format
	g.Expect(jsonObj["name"]).To(gomega.Equal("test"))
	g.Expect(jsonObj["count"]).To(gomega.Equal(float64(42)))
	g.Expect(jsonObj["active"]).To(gomega.Equal(true))
}

func TestMarshalListValue(t *testing.T) {
	g := gomega.NewWithT(t)

	// Create a ListValue
	listValue, err := structpb.NewList([]interface{}{
		"string",
		42,
		true,
		nil,
	})
	g.Expect(err).ToNot(gomega.HaveOccurred())

	result, err := Marshal(listValue)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var jsonArray []interface{}
	err = json.Unmarshal(result, &jsonArray)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	// The list should be marshaled directly as an array in ProtoJSON format
	g.Expect(jsonArray).To(gomega.Equal([]interface{}{
		"string",
		float64(42),
		true,
		nil,
	}))
}

func TestMarshalValue(t *testing.T) {
	g := gomega.NewWithT(t)

	// Test different Value types
	testCases := []struct {
		name     string
		value    *structpb.Value
		expected interface{}
	}{
		{
			name:     "number value",
			value:    structpb.NewNumberValue(3.14),
			expected: 3.14,
		},
		{
			name:     "string value",
			value:    structpb.NewStringValue("hello"),
			expected: "hello",
		},
		{
			name:     "bool value",
			value:    structpb.NewBoolValue(true),
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Marshal(tc.value)
			g.Expect(err).ToNot(gomega.HaveOccurred())

			var jsonValue interface{}
			err = json.Unmarshal(result, &jsonValue)
			g.Expect(err).ToNot(gomega.HaveOccurred())
			g.Expect(jsonValue).To(gomega.Equal(tc.expected))
		})
	}

	// Test null value separately
	t.Run("null value", func(t *testing.T) {
		nullValue := structpb.NewNullValue()
		result, err := Marshal(nullValue)
		g.Expect(err).ToNot(gomega.HaveOccurred())

		var jsonValue interface{}
		err = json.Unmarshal(result, &jsonValue)
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(jsonValue).To(gomega.BeNil())
	})
}

func TestMarshalAny(t *testing.T) {
	g := gomega.NewWithT(t)

	// Create an Any containing a StringValue
	stringVal := &wrapperspb.StringValue{Value: "test"}
	anyVal, err := anypb.New(stringVal)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	result, err := Marshal(anyVal)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var jsonObj map[string]interface{}
	err = json.Unmarshal(result, &jsonObj)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	// Any should have field "1" for type_url and field "2" for value (base64 encoded)
	g.Expect(jsonObj["1"]).To(gomega.Equal("type.googleapis.com/google.protobuf.StringValue"))
	g.Expect(jsonObj["2"]).To(gomega.BeAssignableToTypeOf(""))
}

func TestMarshalNilMessage(t *testing.T) {
	g := gomega.NewWithT(t)

	var nilMessage proto.Message
	result, err := Marshal(nilMessage)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(result).To(gomega.BeNil())
}
