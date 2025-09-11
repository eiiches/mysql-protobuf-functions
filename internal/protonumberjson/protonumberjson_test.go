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

	// Test Timestamp - should now use field numbers
	ts := &timestamppb.Timestamp{Seconds: 1234567890, Nanos: 123456789}
	result, err := Marshal(ts)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var timestampObj map[string]interface{}
	err = json.Unmarshal(result, &timestampObj)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(timestampObj["1"]).To(gomega.Equal(float64(1234567890))) // seconds field
	g.Expect(timestampObj["2"]).To(gomega.Equal(float64(123456789)))  // nanos field

	// Test Duration - should now use field numbers
	dur := &durationpb.Duration{Seconds: 3600, Nanos: 500000000}
	result, err = Marshal(dur)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var durationObj map[string]interface{}
	err = json.Unmarshal(result, &durationObj)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(durationObj["1"]).To(gomega.Equal(float64(3600)))      // seconds field
	g.Expect(durationObj["2"]).To(gomega.Equal(float64(500000000))) // nanos field

	// Test Empty - empty message should still be empty
	empty := &emptypb.Empty{}
	result, err = Marshal(empty)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(string(result)).To(gomega.Equal("{}"))

	// Test StringValue - should now use field numbers
	str := &wrapperspb.StringValue{Value: "hello world"}
	result, err = Marshal(str)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var strObj map[string]interface{}
	err = json.Unmarshal(result, &strObj)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(strObj["1"]).To(gomega.Equal("hello world")) // value field

	// Test Int64Value - should now use field numbers
	int64Val := &wrapperspb.Int64Value{Value: 9223372036854775807}
	result, err = Marshal(int64Val)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var int64Obj map[string]interface{}
	err = json.Unmarshal(result, &int64Obj)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(int64Obj["1"]).To(gomega.Equal(float64(9223372036854775807))) // value field

	// Test BoolValue - should now use field numbers
	boolVal := &wrapperspb.BoolValue{Value: true}
	result, err = Marshal(boolVal)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var boolObj map[string]interface{}
	err = json.Unmarshal(result, &boolObj)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(boolObj["1"]).To(gomega.Equal(true)) // value field

	// Test FieldMask - should now use field numbers
	fm := &fieldmaskpb.FieldMask{Paths: []string{"user.name", "user.email"}}
	result, err = Marshal(fm)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	var fmObj map[string]interface{}
	err = json.Unmarshal(result, &fmObj)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(fmObj["1"]).To(gomega.Equal([]interface{}{"user.name", "user.email"})) // paths field
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

	// The struct should now be marshaled using field numbers (field 1 is the fields map)
	fieldsMap := jsonObj["1"].(map[string]interface{})
	g.Expect(fieldsMap["name"]).To(gomega.BeAssignableToTypeOf(map[string]interface{}{}))
	g.Expect(fieldsMap["count"]).To(gomega.BeAssignableToTypeOf(map[string]interface{}{}))
	g.Expect(fieldsMap["active"]).To(gomega.BeAssignableToTypeOf(map[string]interface{}{}))
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

	var jsonObj map[string]interface{}
	err = json.Unmarshal(result, &jsonObj)
	g.Expect(err).ToNot(gomega.HaveOccurred())

	// The list should now be marshaled using field numbers (field 1 is the values array)
	valuesArray := jsonObj["1"].([]interface{})
	g.Expect(len(valuesArray)).To(gomega.Equal(4))
	// Each element should be a Value message with its own field structure
	for _, elem := range valuesArray {
		g.Expect(elem).To(gomega.BeAssignableToTypeOf(map[string]interface{}{}))
	}
}

func TestMarshalValue(t *testing.T) {
	g := gomega.NewWithT(t)

	// Test different Value types - now using field numbers
	testCases := []struct {
		name          string
		value         *structpb.Value
		expectedField string
		expectedValue interface{}
	}{
		{
			name:          "number value",
			value:         structpb.NewNumberValue(3.14),
			expectedField: "2", // number_value is field 2
			expectedValue: "binary64:0x40091eb851eb851f",
		},
		{
			name:          "string value",
			value:         structpb.NewStringValue("hello"),
			expectedField: "3", // string_value is field 3
			expectedValue: "hello",
		},
		{
			name:          "bool value",
			value:         structpb.NewBoolValue(true),
			expectedField: "4", // bool_value is field 4
			expectedValue: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Marshal(tc.value)
			g.Expect(err).ToNot(gomega.HaveOccurred())

			var jsonObj map[string]interface{}
			err = json.Unmarshal(result, &jsonObj)
			g.Expect(err).ToNot(gomega.HaveOccurred())
			g.Expect(jsonObj[tc.expectedField]).To(gomega.Equal(tc.expectedValue))
		})
	}

	// Test null value separately
	t.Run("null value", func(t *testing.T) {
		nullValue := structpb.NewNullValue()
		result, err := Marshal(nullValue)
		g.Expect(err).ToNot(gomega.HaveOccurred())

		var jsonObj map[string]interface{}
		err = json.Unmarshal(result, &jsonObj)
		g.Expect(err).ToNot(gomega.HaveOccurred())
		// null_value is field 1 with enum value 0
		g.Expect(jsonObj["1"]).To(gomega.Equal(float64(0)))
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
