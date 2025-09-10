package protocgenmysql

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtobufType interface {
	GetKind() protoreflect.Kind
	GetSqlTypeName() string
	GetDefaultValueExpression() string
	GenerateNumberJsonToSqlExpression(numberJsonVar string) string
	GenerateSqlToNumberJsonExpression(sqlVar string) string
	GenerateNumberJsonToJsonExpression(numberJsonVar string) string
	GenerateJsonToNumberJsonExpression(sqlJsonVar string) string
	GenerateSqlToMapKeyExpression(keyVar string) string
	GenerateMapKeyToSqlExpression(mapKeyVar string) string
	GenerateSetterWithZeroValueRemoval(fieldNumber int32, sqlVar string) string
}

type protobufTypeImpl struct {
	Kind                               protoreflect.Kind
	SqlTypeName                        string
	DefaultValueExpression             string // SQL expression for default value
	NumberJsonToSqlExpressionTemplate  string // Template for converting NumberJSON to SQL type
	SqlToNumberJsonExpressionTemplate  string // Template for converting SQL type to NumberJSON
	NumberJsonToJsonExpressionTemplate string // Template for converting NumberJSON to SQL JSON
	JsonToNumberJsonExpressionTemplate string // Template for converting SQL JSON to NumberJSON
	MapKeyToSqlExpressionTemplate      string // Template for converting map key to SQL type
	SqlToMapKeyExpressionTemplate      string // Template for converting SQL type to map key
	ZeroValueConditionTemplate         string // Template for checking if value equals zero value
}

var protobufTypeImpls = map[protoreflect.Kind]protobufTypeImpl{
	protoreflect.BoolKind: {
		Kind:                               protoreflect.BoolKind,
		SqlTypeName:                        "BOOLEAN",
		DefaultValueExpression:             "FALSE",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_bool({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST(({{sqlVar}} IS TRUE) AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST((_pb_json_parse_bool({{numberJsonVar}}) IS TRUE) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "IF({{keyVar}}, 'true', 'false')",
		MapKeyToSqlExpressionTemplate:      "({{mapKeyVar}} = 'true')",
		ZeroValueConditionTemplate:         "{{sqlVar}} IS FALSE",
	},
	protoreflect.StringKind: {
		Kind:                               protoreflect.StringKind,
		SqlTypeName:                        "LONGTEXT",
		DefaultValueExpression:             "''",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_string({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST(JSON_QUOTE({{sqlVar}}) AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(JSON_QUOTE(_pb_json_parse_string({{numberJsonVar}})) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "{{keyVar}}",
		MapKeyToSqlExpressionTemplate:      "{{mapKeyVar}}",
		ZeroValueConditionTemplate:         "{{sqlVar}} = ''",
	},
	protoreflect.BytesKind: {
		Kind:                               protoreflect.BytesKind,
		SqlTypeName:                        "LONGBLOB",
		DefaultValueExpression:             "X''",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_bytes({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST(JSON_QUOTE(_pb_util_to_base64({{sqlVar}})) AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(JSON_QUOTE(_pb_util_to_base64(_pb_json_parse_bytes({{numberJsonVar}}))) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "CAST(JSON_QUOTE(_pb_util_to_base64(_pb_json_parse_bytes({{sqlJsonVar}}))) AS JSON)",
		SqlToMapKeyExpressionTemplate:      "", // bytes cannot be map keys
		MapKeyToSqlExpressionTemplate:      "", // bytes cannot be map keys
		ZeroValueConditionTemplate:         "LENGTH({{sqlVar}}) = 0",
	},
	protoreflect.FloatKind: {
		Kind:                               protoreflect.FloatKind,
		SqlTypeName:                        "FLOAT",
		DefaultValueExpression:             "0.0",
		NumberJsonToSqlExpressionTemplate:  "_pb_util_reinterpret_uint32_as_float(_pb_json_parse_float_as_uint32({{numberJsonVar}}, TRUE))",
		SqlToNumberJsonExpressionTemplate:  "_pb_convert_float_uint32_to_number_json(_pb_util_reinterpret_float_as_uint32({{sqlVar}}))",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_util_reinterpret_uint32_as_float(_pb_json_parse_float_as_uint32({{numberJsonVar}}, TRUE)) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "_pb_convert_float_uint32_to_number_json(_pb_json_parse_float_as_uint32({{sqlJsonVar}}, FALSE))",
		SqlToMapKeyExpressionTemplate:      "", // float cannot be map key
		MapKeyToSqlExpressionTemplate:      "", // float cannot be map key
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0.0",
	},
	protoreflect.DoubleKind: {
		Kind:                               protoreflect.DoubleKind,
		SqlTypeName:                        "DOUBLE",
		DefaultValueExpression:             "0.0",
		NumberJsonToSqlExpressionTemplate:  "_pb_util_reinterpret_uint64_as_double(_pb_json_parse_double_as_uint64({{numberJsonVar}}, TRUE))",
		SqlToNumberJsonExpressionTemplate:  "_pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64({{sqlVar}}))",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_util_reinterpret_uint64_as_double(_pb_json_parse_double_as_uint64({{numberJsonVar}}, TRUE)) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "_pb_convert_double_uint64_to_number_json(_pb_json_parse_double_as_uint64({{sqlJsonVar}}, FALSE))",
		SqlToMapKeyExpressionTemplate:      "", // double cannot be map key
		MapKeyToSqlExpressionTemplate:      "", // double cannot be map key
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0.0",
	},
	protoreflect.Int32Kind: {
		Kind:                               protoreflect.Int32Kind,
		SqlTypeName:                        "INT",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_signed_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_signed_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS SIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.Sint32Kind: {
		Kind:                               protoreflect.Sint32Kind,
		SqlTypeName:                        "INT",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_signed_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_signed_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS SIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.Sfixed32Kind: {
		Kind:                               protoreflect.Sfixed32Kind,
		SqlTypeName:                        "INT",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_signed_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_signed_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS SIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.Uint32Kind: {
		Kind:                               protoreflect.Uint32Kind,
		SqlTypeName:                        "INT UNSIGNED",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_unsigned_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_unsigned_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS UNSIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.Fixed32Kind: {
		Kind:                               protoreflect.Fixed32Kind,
		SqlTypeName:                        "INT UNSIGNED",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_unsigned_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_unsigned_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS UNSIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.Int64Kind: {
		Kind:                               protoreflect.Int64Kind,
		SqlTypeName:                        "BIGINT",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_signed_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_signed_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS SIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.Sint64Kind: {
		Kind:                               protoreflect.Sint64Kind,
		SqlTypeName:                        "BIGINT",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_signed_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_signed_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS SIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.Sfixed64Kind: {
		Kind:                               protoreflect.Sfixed64Kind,
		SqlTypeName:                        "BIGINT",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_signed_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_signed_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS SIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.Uint64Kind: {
		Kind:                               protoreflect.Uint64Kind,
		SqlTypeName:                        "BIGINT UNSIGNED",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_unsigned_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_unsigned_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS UNSIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.Fixed64Kind: {
		Kind:                               protoreflect.Fixed64Kind,
		SqlTypeName:                        "BIGINT UNSIGNED",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_unsigned_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_unsigned_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS UNSIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.EnumKind: {
		Kind:                               protoreflect.EnumKind,
		SqlTypeName:                        "INT",
		DefaultValueExpression:             "0",
		NumberJsonToSqlExpressionTemplate:  "_pb_json_parse_signed_int({{numberJsonVar}})",
		SqlToNumberJsonExpressionTemplate:  "CAST({{sqlVar}} AS JSON)",
		NumberJsonToJsonExpressionTemplate: "CAST(_pb_json_parse_signed_int({{numberJsonVar}}) AS JSON)",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "CAST({{keyVar}} AS CHAR)",
		MapKeyToSqlExpressionTemplate:      "CAST({{mapKeyVar}} AS SIGNED)",
		ZeroValueConditionTemplate:         "{{sqlVar}} = 0",
	},
	protoreflect.MessageKind: {
		Kind:                               protoreflect.MessageKind,
		SqlTypeName:                        "JSON",
		DefaultValueExpression:             "JSON_OBJECT()",
		NumberJsonToSqlExpressionTemplate:  "{{numberJsonVar}}",
		SqlToNumberJsonExpressionTemplate:  "{{sqlVar}}",
		NumberJsonToJsonExpressionTemplate: "{{numberJsonVar}}",
		JsonToNumberJsonExpressionTemplate: "{{sqlJsonVar}}",
		SqlToMapKeyExpressionTemplate:      "", // message cannot be map key
		MapKeyToSqlExpressionTemplate:      "", // message cannot be map key
		ZeroValueConditionTemplate:         "", // messages don't have zero value removal
	},
}

// Helper function to replace template variables
func replaceTemplateVars(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}
	return result
}

// Implement ProtobufType interface for protobufTypeImpl
func (p protobufTypeImpl) GetKind() protoreflect.Kind {
	return p.Kind
}

func (p protobufTypeImpl) GetSqlTypeName() string {
	return p.SqlTypeName
}

func (p protobufTypeImpl) GetDefaultValueExpression() string {
	return p.DefaultValueExpression
}

func (p protobufTypeImpl) GenerateNumberJsonToSqlExpression(numberJsonVar string) string {
	return replaceTemplateVars(p.NumberJsonToSqlExpressionTemplate, map[string]string{
		"numberJsonVar": numberJsonVar,
	})
}

func (p protobufTypeImpl) GenerateSqlToNumberJsonExpression(sqlVar string) string {
	return replaceTemplateVars(p.SqlToNumberJsonExpressionTemplate, map[string]string{
		"sqlVar": sqlVar,
	})
}

func (p protobufTypeImpl) GenerateNumberJsonToJsonExpression(numberJsonVar string) string {
	return replaceTemplateVars(p.NumberJsonToJsonExpressionTemplate, map[string]string{
		"numberJsonVar": numberJsonVar,
	})
}

func (p protobufTypeImpl) GenerateJsonToNumberJsonExpression(sqlJsonVar string) string {
	return replaceTemplateVars(p.JsonToNumberJsonExpressionTemplate, map[string]string{
		"sqlJsonVar": sqlJsonVar,
	})
}

func (p protobufTypeImpl) GenerateSqlToMapKeyExpression(keyVar string) string {
	if p.SqlToMapKeyExpressionTemplate == "" {
		panic(fmt.Sprintf("%s type cannot be used as map key", p.Kind))
	}
	return replaceTemplateVars(p.SqlToMapKeyExpressionTemplate, map[string]string{
		"keyVar": keyVar,
	})
}

func (p protobufTypeImpl) GenerateMapKeyToSqlExpression(mapKeyVar string) string {
	if p.MapKeyToSqlExpressionTemplate == "" {
		panic(fmt.Sprintf("%s type cannot be used as map key", p.Kind))
	}
	return replaceTemplateVars(p.MapKeyToSqlExpressionTemplate, map[string]string{
		"mapKeyVar": mapKeyVar,
	})
}

func (p protobufTypeImpl) GenerateSetterWithZeroValueRemoval(fieldNumber int32, sqlVar string) string {
	if p.ZeroValueConditionTemplate == "" {
		// Messages and other types that don't support zero value removal
		jsonExpression := p.GenerateSqlToNumberJsonExpression(sqlVar)
		return fmt.Sprintf("RETURN JSON_SET(proto_data, '$.\"%.d\"', %s);", fieldNumber, jsonExpression)
	}

	condition := replaceTemplateVars(p.ZeroValueConditionTemplate, map[string]string{
		"sqlVar": sqlVar,
	})
	jsonExpression := p.GenerateSqlToNumberJsonExpression(sqlVar)

	return fmt.Sprintf("IF %s THEN\n        RETURN JSON_REMOVE(proto_data, '$.\"%.d\"');\n    END IF;\n    RETURN JSON_SET(proto_data, '$.\"%.d\"', %s);",
		condition, fieldNumber, fieldNumber, jsonExpression)
}

// GetProtobufType returns the ProtobufType for the given kind, fallback to string type
func GetProtobufType(kind protoreflect.Kind) ProtobufType {
	if ptype, ok := protobufTypeImpls[kind]; ok {
		return ptype
	}
	panic(fmt.Sprintf("unsupported protobuf kind: %s", kind))
}

// GetDefaultValue returns the appropriate default value for a field
func GetDefaultValue(field protoreflect.FieldDescriptor) string {
	// Handle explicit default values (proto2 only)
	if field.HasDefault() {
		defaultVal := field.Default()
		switch field.Kind() {
		case protoreflect.BoolKind:
			if defaultVal.Bool() {
				return "TRUE"
			}
			return "FALSE"
		case protoreflect.StringKind:
			return "'" + strings.ReplaceAll(defaultVal.String(), "'", "''") + "'"
		case protoreflect.BytesKind:
			return fmt.Sprintf("X'%X'", defaultVal.Bytes())
		case protoreflect.FloatKind:
			return "_pb_util_reinterpret_uint32_as_float(" + defaultVal.String() + ")"
		case protoreflect.DoubleKind:
			return "_pb_util_reinterpret_uint64_as_double(" + defaultVal.String() + ")"
		case protoreflect.EnumKind:
			return defaultVal.String()
		default:
			// Numeric types
			return defaultVal.String()
		}
	}

	// Proto3 and proto2 fields without explicit defaults use zero values from ProtobufType
	ptype := GetProtobufType(field.Kind())
	return ptype.GetDefaultValueExpression()
}
