package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v3"
)

type WireTypeAccessor struct {
	SqlType                       string
	SupportsPacked                bool
	GetFunction                   string
	GetElementFunction            string
	SetFunction                   string
	AddRepeatedElementFunction    string
	SetRepeatedElementFunction    string
	RemoveRepeatedElementFunction string
	InsertRepeatedElementFunction string
}

type Accessor struct {
	ProtoType      string
	SqlType        string
	ReturnExpr     string
	Procedure      *WireTypeAccessor
	Input          *Input
	ConvertExpr    string
	SupportsPacked bool
}

type Input struct {
	Kind    string // message or wire_json (used as part of the function or procedure name)
	Name    string // message or wire_json (used as a parameter name in the function)
	SqlType string // LONGBLOB or JSON (used as a parameter type in the function)
}

type RepeatedAsJsonAccessor struct {
	ProtoType           string
	SqlType             string
	Expr                string
	PackedUint64Decoder string
	WireType            int
	Suffix              string
}

func generateRepeatedNumbersAsJson() {
	accessors := []*RepeatedAsJsonAccessor{
		{
			ProtoType:           "int32",
			SqlType:             "INT",
			Expr:                "_pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_int64(uint_value))",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "uint32",
			SqlType:             "INT UNSIGNED",
			Expr:                "_pb_util_cast_uint64_as_uint32(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "int64",
			SqlType:             "BIGINT",
			Expr:                "_pb_util_reinterpret_uint64_as_int64(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "int64",
			SqlType:             "BIGINT",
			Expr:                "CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_string_array",
		},
		{
			ProtoType:           "uint64",
			SqlType:             "BIGINT UNSIGNED",
			Expr:                "uint_value",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "uint64",
			SqlType:             "BIGINT UNSIGNED",
			Expr:                "CAST(uint_value AS CHAR)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_string_array",
		},
		{
			ProtoType:           "sint32",
			SqlType:             "INT",
			Expr:                "_pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_sint64(uint_value))",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "sint64",
			SqlType:             "BIGINT",
			Expr:                "_pb_util_reinterpret_uint64_as_sint64(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "sint64",
			SqlType:             "BIGINT",
			Expr:                "CAST(_pb_util_reinterpret_uint64_as_sint64(uint_value) AS CHAR)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_string_array",
		},
		{
			ProtoType:           "enum",
			SqlType:             "INT",
			Expr:                "_pb_util_reinterpret_uint64_as_int64(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "bool",
			SqlType:             "BOOLEAN",
			Expr:                "uint_value <> 0",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "fixed32",
			SqlType:             "INT UNSIGNED",
			Expr:                "uint_value",
			PackedUint64Decoder: "_pb_wire_read_i32_as_uint32",
			WireType:            5,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "sfixed32",
			SqlType:             "INT",
			Expr:                "_pb_util_reinterpret_uint32_as_int32(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_i32_as_uint32",
			WireType:            5,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "float",
			SqlType:             "FLOAT",
			Expr:                "_pb_util_reinterpret_uint32_as_float(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_i32_as_uint32",
			WireType:            5,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "fixed64",
			SqlType:             "BIGINT UNSIGNED",
			Expr:                "uint_value",
			PackedUint64Decoder: "_pb_wire_read_i64_as_uint64",
			WireType:            1,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "fixed64",
			SqlType:             "BIGINT UNSIGNED",
			Expr:                "CAST(uint_value AS CHAR)",
			PackedUint64Decoder: "_pb_wire_read_i64_as_uint64",
			WireType:            1,
			Suffix:              "_as_json_string_array",
		},
		{
			ProtoType:           "sfixed64",
			SqlType:             "BIGINT",
			Expr:                "_pb_util_reinterpret_uint64_as_int64(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_i64_as_uint64",
			WireType:            1,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "sfixed64",
			SqlType:             "BIGINT",
			Expr:                "CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR)",
			PackedUint64Decoder: "_pb_wire_read_i64_as_uint64",
			WireType:            1,
			Suffix:              "_as_json_string_array",
		},
		{
			ProtoType:           "double",
			SqlType:             "DOUBLE",
			Expr:                "_pb_util_reinterpret_uint64_as_double(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_i64_as_uint64",
			WireType:            1,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "bytes",
			SqlType:             "LONGBLOB",
			Expr:                "TO_BASE64(bytes_value)",
			PackedUint64Decoder: "",
			WireType:            2,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "string",
			SqlType:             "LONGTEXT",
			Expr:                "CONVERT(bytes_value USING utf8mb4)",
			PackedUint64Decoder: "",
			WireType:            2,
			Suffix:              "_as_json_array",
		},
		{
			ProtoType:           "message",
			SqlType:             "LONGBLOB",
			Expr:                "TO_BASE64(bytes_value)",
			PackedUint64Decoder: "",
			WireType:            2,
			Suffix:              "_as_json_array",
		},
	}

	templateText := dedent.Pipe(`
		|
		|DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_{{.ProtoType}}_field{{.Suffix}} $$
		|CREATE PROCEDURE _pb_wire_json_get_repeated_{{.ProtoType}}_field{{.Suffix}}(IN wire_json JSON, IN field_number INT, OUT result JSON)
		|BEGIN
		|	DECLARE message_text TEXT;
		|	DECLARE uint_value BIGINT UNSIGNED;
		|	DECLARE bytes_value LONGBLOB;
		|	DECLARE wire_type INT;
		|	DECLARE wire_elements JSON;
		|	DECLARE wire_element JSON;
		|	DECLARE wire_element_index INT;
		|	DECLARE wire_element_count INT;
		|
		|	SET result = JSON_ARRAY();
		|
		|	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
		|	SET wire_element_index = 0;
		|	SET wire_element_count = JSON_LENGTH(wire_elements);
		|
		|	l1: WHILE wire_element_index < wire_element_count DO
		|		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		|		SET wire_type = JSON_EXTRACT(wire_element, '$.t');
		|
		|		CASE wire_type
		|		WHEN {{.WireType}} THEN
		|{{- if eq .WireType 2 }}
		|			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
		|{{- else }}
		|			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
		|{{- end }}
		|			SET result = JSON_ARRAY_APPEND(result, '$', {{.Expr}});
		|{{- if .PackedUint64Decoder }}
		|		WHEN 2 THEN -- LEN
		|			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
		|			WHILE LENGTH(bytes_value) <> 0 DO
		|				CALL {{.PackedUint64Decoder}}(bytes_value, uint_value, bytes_value);
		|				SET result = JSON_ARRAY_APPEND(result, '$', {{.Expr}});
		|			END WHILE;
		|{{- end }}
		|		ELSE
		|			SET message_text = CONCAT('_pb_wire_json_get_repeated_{{.ProtoType}}_field{{.Suffix}}: unexpected wire_type (', wire_type, ')');
		|			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		|		END CASE;
		|
		|		SET wire_element_index = wire_element_index + 1;
		|	END WHILE;
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_{{.ProtoType}}_field{{.Suffix}} $$
		|CREATE FUNCTION pb_wire_json_get_repeated_{{.ProtoType}}_field{{.Suffix}}(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
		|BEGIN
		|	DECLARE result JSON;
		|	CALL _pb_wire_json_get_repeated_{{.ProtoType}}_field{{.Suffix}}(wire_json, field_number, result);
		|	RETURN result;
		|END $$
		|
		|DROP PROCEDURE IF EXISTS _pb_message_get_repeated_{{.ProtoType}}_field{{.Suffix}} $$
		|CREATE PROCEDURE _pb_message_get_repeated_{{.ProtoType}}_field{{.Suffix}}(IN message LONGBLOB, IN field_number INT, OUT result JSON)
		|BEGIN
		|	DECLARE tag BIGINT;
		|	DECLARE tail LONGBLOB;
		|	DECLARE uint_value BIGINT UNSIGNED;
		|	DECLARE bytes_value LONGBLOB;
		|	DECLARE message_text TEXT;
		|	DECLARE current_field_number INT;
		|	DECLARE current_wire_type INT;
		|
		|	SET tail = message;
		|	SET result = JSON_ARRAY();
		|
		|	l1: WHILE LENGTH(tail) <> 0 DO
		|		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		|		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		|		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);
		|
		|		IF current_field_number != field_number THEN
		|			CALL _pb_wire_skip(tail, current_wire_type, tail);
		|			ITERATE l1;
		|		END IF;
		|
		|		CASE current_wire_type
		|		WHEN {{.WireType}} THEN
		|{{- if eq .WireType 0 }}
		|			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
		|{{- else if eq .WireType 2 }}
		|			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
		|{{- else if eq .WireType 1 }}
		|			CALL _pb_wire_read_i64_as_uint64(tail, uint_value, tail);
		|{{- else if eq .WireType 5 }}
		|			CALL _pb_wire_read_i32_as_uint32(tail, uint_value, tail);
		|{{- end }}
		|			SET result = JSON_ARRAY_APPEND(result, '$', {{.Expr}});
		|{{- if .PackedUint64Decoder }}
		|		WHEN 2 THEN
		|			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
		|			WHILE LENGTH(bytes_value) <> 0 DO
		|				CALL {{.PackedUint64Decoder}}(bytes_value, uint_value, bytes_value);
		|				SET result = JSON_ARRAY_APPEND(result, '$', {{.Expr}});
		|			END WHILE;
		|{{- end }}
		|		ELSE
		|			SET message_text = CONCAT('_pb_message_get_repeated_{{.ProtoType}}_field{{.Suffix}}: unexpected wire_type (', current_wire_type, ')');
		|			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		|		END CASE;
		|	END WHILE;
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_message_get_repeated_{{.ProtoType}}_field{{.Suffix}} $$
		|CREATE FUNCTION pb_message_get_repeated_{{.ProtoType}}_field{{.Suffix}}(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
		|BEGIN
		|	DECLARE result JSON;
		|	CALL _pb_message_get_repeated_{{.ProtoType}}_field{{.Suffix}}(message, field_number, result);
		|	RETURN result;
		|END $$
	`)

	tmpl, err := template.New("t").Parse(templateText)
	if err != nil {
		panic(err)
	}

	for _, accessor := range accessors {
		if err := tmpl.Execute(os.Stdout, accessor); err != nil {
			panic(err)
		}
	}
}

func generateAccessorsAction(ctx context.Context, command *cli.Command) error {
	inputs := []*Input{
		{Kind: "message", Name: "message", SqlType: "LONGBLOB"},
		{Kind: "wire_json", Name: "wire_json", SqlType: "JSON"},
	}

	templateText := dedent.Pipe(`
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_get_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_get_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT, default_value {{.SqlType}}) RETURNS {{.SqlType}} DETERMINISTIC
		|BEGIN
		|	DECLARE value {{.Procedure.SqlType}};
		|	DECLARE field_count INT;
		|	CALL {{.Procedure.GetFunction}}({{.Input.Name}}, field_number, NULL, value, field_count);
		|	IF field_count = 0 THEN
		|		RETURN default_value;
		|	END IF;
		|	RETURN {{.ReturnExpr}};
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_has_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_has_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT) RETURNS BOOLEAN DETERMINISTIC
		|BEGIN
		|	DECLARE value {{.Procedure.SqlType}};
		|	DECLARE field_count INT;
		|	CALL {{.Procedure.GetFunction}}({{.Input.Name}}, field_number, NULL, value, field_count);
		|	RETURN field_count > 0;
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_get_repeated_{{.ProtoType}}_field_element $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_get_repeated_{{.ProtoType}}_field_element({{.Input.Name}} {{.Input.SqlType}}, field_number INT, repeated_index INT) RETURNS {{.SqlType}} DETERMINISTIC
		|BEGIN
		|	DECLARE value {{.Procedure.SqlType}};
		|	DECLARE field_count INT;
		|	CALL {{.Procedure.GetElementFunction}}({{.Input.Name}}, field_number, repeated_index, value, field_count);
		|	RETURN {{.ReturnExpr}};
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_get_repeated_{{.ProtoType}}_field_count $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_get_repeated_{{.ProtoType}}_field_count({{.Input.Name}} {{.Input.SqlType}}, field_number INT) RETURNS INT DETERMINISTIC
		|BEGIN
		|	DECLARE value {{.Procedure.SqlType}};
		|	DECLARE field_count INT;
		|	CALL {{.Procedure.GetElementFunction}}({{.Input.Name}}, field_number, -1, value, field_count);
		|	RETURN field_count;
		|END $$
	`)

	os.Stdout.WriteString("DELIMITER $$\n")

	for _, input := range inputs {
		getVarintFieldAsUint64 := &WireTypeAccessor{
			SqlType:                       "BIGINT UNSIGNED",
			SupportsPacked:                true,
			GetFunction:                   fmt.Sprintf("_pb_%s_get_varint_field_as_uint64", input.Kind),
			GetElementFunction:            fmt.Sprintf("_pb_%s_get_varint_field_as_uint64", input.Kind),
			SetFunction:                   "_pb_wire_json_set_varint_field",
			AddRepeatedElementFunction:    "_pb_wire_json_add_repeated_varint_field_element",
			SetRepeatedElementFunction:    "_pb_wire_json_set_repeated_varint_field_element",
			RemoveRepeatedElementFunction: "_pb_wire_json_remove_repeated_varint_field_element",
			InsertRepeatedElementFunction: "_pb_wire_json_insert_repeated_varint_field_element",
		}

		getI64FieldAsUint64 := &WireTypeAccessor{
			SqlType:                       "BIGINT UNSIGNED",
			SupportsPacked:                true,
			GetFunction:                   fmt.Sprintf("_pb_%s_get_i64_field_as_uint64", input.Kind),
			GetElementFunction:            fmt.Sprintf("_pb_%s_get_i64_field_as_uint64", input.Kind),
			SetFunction:                   "_pb_wire_json_set_i64_field",
			AddRepeatedElementFunction:    "_pb_wire_json_add_repeated_i64_field_element",
			SetRepeatedElementFunction:    "_pb_wire_json_set_repeated_i64_field_element",
			RemoveRepeatedElementFunction: "_pb_wire_json_remove_repeated_i64_field_element",
			InsertRepeatedElementFunction: "_pb_wire_json_insert_repeated_i64_field_element",
		}

		getI32FieldAsUint64 := &WireTypeAccessor{
			SqlType:                       "INT UNSIGNED",
			SupportsPacked:                true,
			GetFunction:                   fmt.Sprintf("_pb_%s_get_i32_field_as_uint32", input.Kind),
			GetElementFunction:            fmt.Sprintf("_pb_%s_get_i32_field_as_uint32", input.Kind),
			SetFunction:                   "_pb_wire_json_set_i32_field",
			AddRepeatedElementFunction:    "_pb_wire_json_add_repeated_i32_field_element",
			SetRepeatedElementFunction:    "_pb_wire_json_set_repeated_i32_field_element",
			RemoveRepeatedElementFunction: "_pb_wire_json_remove_repeated_i32_field_element",
			InsertRepeatedElementFunction: "_pb_wire_json_insert_repeated_i32_field_element",
		}

		getLengthDelimitedField := &WireTypeAccessor{
			SqlType:                       "LONGBLOB",
			SupportsPacked:                false,
			GetFunction:                   fmt.Sprintf("_pb_%s_get_len_type_field", input.Kind),
			GetElementFunction:            fmt.Sprintf("_pb_%s_get_len_type_field", input.Kind),
			SetFunction:                   "_pb_wire_json_set_len_field",
			AddRepeatedElementFunction:    "_pb_wire_json_add_repeated_len_field_element",
			SetRepeatedElementFunction:    "_pb_wire_json_set_repeated_len_field_element",
			RemoveRepeatedElementFunction: "_pb_wire_json_remove_repeated_len_field_element",
			InsertRepeatedElementFunction: "_pb_wire_json_insert_repeated_len_field_element",
		}

		getLengthDelimitedConcatField := &WireTypeAccessor{ // for message field
			SqlType:                       "LONGBLOB",
			SupportsPacked:                false,
			GetFunction:                   fmt.Sprintf("_pb_%s_get_len_type_field_concatenated", input.Kind),
			GetElementFunction:            fmt.Sprintf("_pb_%s_get_len_type_field", input.Kind),
			SetFunction:                   "_pb_wire_json_set_len_field",
			AddRepeatedElementFunction:    "_pb_wire_json_add_repeated_len_field_element",
			SetRepeatedElementFunction:    "_pb_wire_json_set_repeated_len_field_element",
			RemoveRepeatedElementFunction: "_pb_wire_json_remove_repeated_len_field_element",
			InsertRepeatedElementFunction: "_pb_wire_json_insert_repeated_len_field_element",
		}

		accessors := []*Accessor{
			// VARINT
			{Input: input, ProtoType: "int32", SqlType: "INT", ReturnExpr: "_pb_util_reinterpret_uint64_as_int64(value)", Procedure: getVarintFieldAsUint64, ConvertExpr: "_pb_util_reinterpret_int64_as_uint64(value)", SupportsPacked: true},
			{Input: input, ProtoType: "int64", SqlType: "BIGINT", ReturnExpr: "_pb_util_reinterpret_uint64_as_int64(value)", Procedure: getVarintFieldAsUint64, ConvertExpr: "_pb_util_reinterpret_int64_as_uint64(value)", SupportsPacked: true},
			{Input: input, ProtoType: "uint32", SqlType: "INT UNSIGNED", ReturnExpr: "value", Procedure: getVarintFieldAsUint64, ConvertExpr: "value", SupportsPacked: true},
			{Input: input, ProtoType: "uint64", SqlType: "BIGINT UNSIGNED", ReturnExpr: "value", Procedure: getVarintFieldAsUint64, ConvertExpr: "value", SupportsPacked: true},
			{Input: input, ProtoType: "sint32", SqlType: "INT", ReturnExpr: "_pb_util_reinterpret_uint64_as_sint64(value)", Procedure: getVarintFieldAsUint64, ConvertExpr: "_pb_util_reinterpret_sint64_as_uint64(value)", SupportsPacked: true},
			{Input: input, ProtoType: "sint64", SqlType: "BIGINT", ReturnExpr: "_pb_util_reinterpret_uint64_as_sint64(value)", Procedure: getVarintFieldAsUint64, ConvertExpr: "_pb_util_reinterpret_sint64_as_uint64(value)", SupportsPacked: true},
			{Input: input, ProtoType: "enum", SqlType: "INT", ReturnExpr: "_pb_util_reinterpret_uint64_as_int64(value)", Procedure: getVarintFieldAsUint64, ConvertExpr: "_pb_util_reinterpret_int64_as_uint64(value)", SupportsPacked: true},
			{Input: input, ProtoType: "bool", SqlType: "BOOLEAN", ReturnExpr: "value <> 0", Procedure: getVarintFieldAsUint64, ConvertExpr: "IF(value, 1, 0)", SupportsPacked: true},

			// I32
			{Input: input, ProtoType: "fixed32", SqlType: "INT UNSIGNED", ReturnExpr: "value", Procedure: getI32FieldAsUint64, ConvertExpr: "value", SupportsPacked: true},
			{Input: input, ProtoType: "sfixed32", SqlType: "INT", ReturnExpr: "_pb_util_reinterpret_uint32_as_int32(value)", Procedure: getI32FieldAsUint64, ConvertExpr: "_pb_util_reinterpret_int32_as_uint32(value)", SupportsPacked: true},
			{Input: input, ProtoType: "float", SqlType: "FLOAT", ReturnExpr: "_pb_util_reinterpret_uint32_as_float(value)", Procedure: getI32FieldAsUint64, ConvertExpr: "_pb_util_reinterpret_float_as_uint32(value)", SupportsPacked: true},

			// I64
			{Input: input, ProtoType: "fixed64", SqlType: "BIGINT UNSIGNED", ReturnExpr: "value", Procedure: getI64FieldAsUint64, ConvertExpr: "value", SupportsPacked: true},
			{Input: input, ProtoType: "sfixed64", SqlType: "BIGINT", ReturnExpr: "_pb_util_reinterpret_uint64_as_int64(value)", Procedure: getI64FieldAsUint64, ConvertExpr: "_pb_util_reinterpret_int64_as_uint64(value)", SupportsPacked: true},
			{Input: input, ProtoType: "double", SqlType: "DOUBLE", ReturnExpr: "_pb_util_reinterpret_uint64_as_double(value)", Procedure: getI64FieldAsUint64, ConvertExpr: "_pb_util_reinterpret_double_as_uint64(value)", SupportsPacked: true},

			// LEN
			{Input: input, ProtoType: "bytes", SqlType: "LONGBLOB", ReturnExpr: "value", Procedure: getLengthDelimitedField, ConvertExpr: "value", SupportsPacked: false},
			{Input: input, ProtoType: "string", SqlType: "LONGTEXT", ReturnExpr: "CONVERT(value USING utf8mb4)", Procedure: getLengthDelimitedField, ConvertExpr: "CONVERT(value USING binary)", SupportsPacked: false},
			{Input: input, ProtoType: "message", SqlType: "LONGBLOB", ReturnExpr: "value", Procedure: getLengthDelimitedConcatField, ConvertExpr: "value", SupportsPacked: false},
		}

		tmpl, err := template.New("t").Parse(templateText)
		if err != nil {
			panic(err)
		}
		for _, accessor := range accessors {
			if err := tmpl.Execute(os.Stdout, accessor); err != nil {
				panic(err)
			}
		}

		// Generate setter functions only for wire_json input
		if input.Kind == "wire_json" {
			generateWireJsonSetters(accessors)
		}

		// Generate message setter functions only for message input
		if input.Kind == "message" {
			generateMessageSetters(accessors)
		}
	}

	generateRepeatedNumbersAsJson()

	generateBulkAddFunctions()

	return nil
}

func generateWireJsonSetters(accessors []*Accessor) {
	setterTemplateText := dedent.Pipe(`
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_set_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_set_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT, value {{.SqlType}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN {{.Procedure.SetFunction}}({{.Input.Name}}, field_number, {{.ConvertExpr}});
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_add_repeated_{{.ProtoType}}_field_element $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_add_repeated_{{.ProtoType}}_field_element({{.Input.Name}} {{.Input.SqlType}}, field_number INT, value {{.SqlType}}{{if .SupportsPacked}}, use_packed BOOLEAN{{end}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|{{- if .Procedure.SupportsPacked}}
		|	RETURN {{.Procedure.AddRepeatedElementFunction}}({{.Input.Name}}, field_number, {{.ConvertExpr}}{{if .SupportsPacked}}, use_packed{{else}}, FALSE{{end}});
		|{{- else}}
		|	RETURN {{.Procedure.AddRepeatedElementFunction}}({{.Input.Name}}, field_number, {{.ConvertExpr}});
		|{{- end}}
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_add_all_repeated_{{.ProtoType}}_field_elements $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_add_all_repeated_{{.ProtoType}}_field_elements({{.Input.Name}} {{.Input.SqlType}}, field_number INT, value_array JSON{{if .SupportsPacked}}, use_packed BOOLEAN{{end}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN _pb_wire_json_add_all_repeated_{{.ProtoType}}_field_elements({{.Input.Name}}, field_number, value_array{{if .SupportsPacked}}, use_packed{{end}});
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_insert_repeated_{{.ProtoType}}_field_element $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_insert_repeated_{{.ProtoType}}_field_element({{.Input.Name}} {{.Input.SqlType}}, field_number INT, repeated_index INT, value {{.SqlType}}{{if .SupportsPacked}}, use_packed BOOLEAN{{end}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|{{- if .Procedure.SupportsPacked}}
		|	RETURN {{.Procedure.InsertRepeatedElementFunction}}({{.Input.Name}}, field_number, repeated_index, {{.ConvertExpr}}{{if .SupportsPacked}}, use_packed{{else}}, FALSE{{end}});
		|{{- else}}
		|	RETURN {{.Procedure.InsertRepeatedElementFunction}}({{.Input.Name}}, field_number, repeated_index, {{.ConvertExpr}});
		|{{- end}}
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_set_repeated_{{.ProtoType}}_field_element $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_set_repeated_{{.ProtoType}}_field_element({{.Input.Name}} {{.Input.SqlType}}, field_number INT, repeated_index INT, value {{.SqlType}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN {{.Procedure.SetRepeatedElementFunction}}({{.Input.Name}}, field_number, repeated_index, {{.ConvertExpr}});
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_remove_repeated_{{.ProtoType}}_field_element $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_remove_repeated_{{.ProtoType}}_field_element({{.Input.Name}} {{.Input.SqlType}}, field_number INT, repeated_index INT) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN {{.Procedure.RemoveRepeatedElementFunction}}({{.Input.Name}}, field_number, repeated_index);
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_clear_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_clear_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN _pb_wire_json_clear_field({{.Input.Name}}, field_number);
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_clear_repeated_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_clear_repeated_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN _pb_wire_json_clear_field({{.Input.Name}}, field_number);
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_set_repeated_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_set_repeated_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT, value_array JSON{{if .SupportsPacked}}, use_packed BOOLEAN{{end}}) RETURNS JSON DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_add_all_repeated_{{.ProtoType}}_field_elements(pb_wire_json_clear_repeated_{{.ProtoType}}_field({{.Input.Name}}, field_number), field_number, value_array{{if .SupportsPacked}}, use_packed{{end}});
		|END $$
	`)

	tmpl, err := template.New("setter").Parse(setterTemplateText)
	if err != nil {
		panic(err)
	}

	for _, accessor := range accessors {
		if err := tmpl.Execute(os.Stdout, accessor); err != nil {
			panic(err)
		}
	}
}

func generateMessageSetters(accessors []*Accessor) {
	messageSetterTemplateText := dedent.Pipe(`
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_set_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_set_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT, value {{.SqlType}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_to_message(pb_wire_json_set_{{.ProtoType}}_field(pb_message_to_wire_json({{.Input.Name}}), field_number, value));
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_add_repeated_{{.ProtoType}}_field_element $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_add_repeated_{{.ProtoType}}_field_element({{.Input.Name}} {{.Input.SqlType}}, field_number INT, value {{.SqlType}}{{if .SupportsPacked}}, use_packed BOOLEAN{{end}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_{{.ProtoType}}_field_element(pb_message_to_wire_json({{.Input.Name}}), field_number, value{{if .SupportsPacked}}, use_packed{{end}}));
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_add_all_repeated_{{.ProtoType}}_field_elements $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_add_all_repeated_{{.ProtoType}}_field_elements({{.Input.Name}} {{.Input.SqlType}}, field_number INT, value_array JSON{{if .SupportsPacked}}, use_packed BOOLEAN{{end}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_to_message(pb_wire_json_add_all_repeated_{{.ProtoType}}_field_elements(pb_message_to_wire_json({{.Input.Name}}), field_number, value_array{{if .SupportsPacked}}, use_packed{{end}}));
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_insert_repeated_{{.ProtoType}}_field_element $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_insert_repeated_{{.ProtoType}}_field_element({{.Input.Name}} {{.Input.SqlType}}, field_number INT, repeated_index INT, value {{.SqlType}}{{if .SupportsPacked}}, use_packed BOOLEAN{{end}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_to_message(pb_wire_json_insert_repeated_{{.ProtoType}}_field_element(pb_message_to_wire_json({{.Input.Name}}), field_number, repeated_index, value{{if .SupportsPacked}}, use_packed{{end}}));
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_set_repeated_{{.ProtoType}}_field_element $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_set_repeated_{{.ProtoType}}_field_element({{.Input.Name}} {{.Input.SqlType}}, field_number INT, repeated_index INT, value {{.SqlType}}) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_{{.ProtoType}}_field_element(pb_message_to_wire_json({{.Input.Name}}), field_number, repeated_index, value));
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_remove_repeated_{{.ProtoType}}_field_element $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_remove_repeated_{{.ProtoType}}_field_element({{.Input.Name}} {{.Input.SqlType}}, field_number INT, repeated_index INT) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_{{.ProtoType}}_field_element(pb_message_to_wire_json({{.Input.Name}}), field_number, repeated_index));
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_clear_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_clear_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_to_message(pb_wire_json_clear_{{.ProtoType}}_field(pb_message_to_wire_json({{.Input.Name}}), field_number));
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_clear_repeated_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_clear_repeated_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT) RETURNS {{.Input.SqlType}} DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_{{.ProtoType}}_field(pb_message_to_wire_json({{.Input.Name}}), field_number));
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_set_repeated_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_set_repeated_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT, value_array JSON{{if .SupportsPacked}}, use_packed BOOLEAN{{end}}) RETURNS LONGBLOB DETERMINISTIC
		|BEGIN
		|	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_{{.ProtoType}}_field(pb_message_to_wire_json({{.Input.Name}}), field_number, value_array{{if .SupportsPacked}}, use_packed{{end}}));
		|END $$
	`)

	tmpl, err := template.New("messageSetter").Parse(messageSetterTemplateText)
	if err != nil {
		panic(err)
	}

	for _, accessor := range accessors {
		if err := tmpl.Execute(os.Stdout, accessor); err != nil {
			panic(err)
		}
	}
}

func generateBulkAddFunctions() {
	// Group accessors by type category
	type BulkAccessor struct {
		ProtoType       string
		SqlType         string
		ConvertExpr     string
		JsonExtractExpr string
		WireType        int
		WriteFunction   string
		SupportsPacked  bool
	}

	bulkAccessors := []*BulkAccessor{
		// VARINT types
		{ProtoType: "int32", SqlType: "INT", ConvertExpr: "_pb_util_reinterpret_int64_as_uint64(current_value)", JsonExtractExpr: "JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))", WireType: 0, WriteFunction: "_pb_wire_write_varint", SupportsPacked: true},
		{ProtoType: "int64", SqlType: "BIGINT", ConvertExpr: "_pb_util_reinterpret_int64_as_uint64(current_value)", JsonExtractExpr: "CAST(JSON_UNQUOTE(JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))) AS SIGNED)", WireType: 0, WriteFunction: "_pb_wire_write_varint", SupportsPacked: true},
		{ProtoType: "uint32", SqlType: "INT UNSIGNED", ConvertExpr: "current_value", JsonExtractExpr: "JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))", WireType: 0, WriteFunction: "_pb_wire_write_varint", SupportsPacked: true},
		{ProtoType: "uint64", SqlType: "BIGINT UNSIGNED", ConvertExpr: "current_value", JsonExtractExpr: "CAST(JSON_UNQUOTE(JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))) AS UNSIGNED)", WireType: 0, WriteFunction: "_pb_wire_write_varint", SupportsPacked: true},
		{ProtoType: "sint32", SqlType: "INT", ConvertExpr: "_pb_util_reinterpret_sint64_as_uint64(current_value)", JsonExtractExpr: "JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))", WireType: 0, WriteFunction: "_pb_wire_write_varint", SupportsPacked: true},
		{ProtoType: "sint64", SqlType: "BIGINT", ConvertExpr: "_pb_util_reinterpret_sint64_as_uint64(current_value)", JsonExtractExpr: "CAST(JSON_UNQUOTE(JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))) AS SIGNED)", WireType: 0, WriteFunction: "_pb_wire_write_varint", SupportsPacked: true},
		{ProtoType: "enum", SqlType: "INT", ConvertExpr: "_pb_util_reinterpret_int64_as_uint64(current_value)", JsonExtractExpr: "JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))", WireType: 0, WriteFunction: "_pb_wire_write_varint", SupportsPacked: true},
		{ProtoType: "bool", SqlType: "BOOLEAN", ConvertExpr: "IF(current_value, 1, 0)", JsonExtractExpr: "JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))", WireType: 0, WriteFunction: "_pb_wire_write_varint", SupportsPacked: true},

		// I32 types
		{ProtoType: "fixed32", SqlType: "INT UNSIGNED", ConvertExpr: "current_value", JsonExtractExpr: "JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))", WireType: 5, WriteFunction: "_pb_wire_write_i32", SupportsPacked: true},
		{ProtoType: "sfixed32", SqlType: "INT", ConvertExpr: "_pb_util_reinterpret_int32_as_uint32(current_value)", JsonExtractExpr: "JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))", WireType: 5, WriteFunction: "_pb_wire_write_i32", SupportsPacked: true},
		{ProtoType: "float", SqlType: "FLOAT", ConvertExpr: "_pb_util_reinterpret_float_as_uint32(current_value)", JsonExtractExpr: "JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))", WireType: 5, WriteFunction: "_pb_wire_write_i32", SupportsPacked: true},

		// I64 types
		{ProtoType: "fixed64", SqlType: "BIGINT UNSIGNED", ConvertExpr: "current_value", JsonExtractExpr: "CAST(JSON_UNQUOTE(JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))) AS UNSIGNED)", WireType: 1, WriteFunction: "_pb_wire_write_i64", SupportsPacked: true},
		{ProtoType: "sfixed64", SqlType: "BIGINT", ConvertExpr: "_pb_util_reinterpret_int64_as_uint64(current_value)", JsonExtractExpr: "CAST(JSON_UNQUOTE(JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))) AS SIGNED)", WireType: 1, WriteFunction: "_pb_wire_write_i64", SupportsPacked: true},
		{ProtoType: "double", SqlType: "DOUBLE", ConvertExpr: "_pb_util_reinterpret_double_as_uint64(current_value)", JsonExtractExpr: "JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))", WireType: 1, WriteFunction: "_pb_wire_write_i64", SupportsPacked: true},

		// Length-delimited types (no packed encoding)
		{ProtoType: "bytes", SqlType: "LONGBLOB", ConvertExpr: "current_value", JsonExtractExpr: "FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))))", WireType: 2, WriteFunction: "", SupportsPacked: false},
		{ProtoType: "string", SqlType: "LONGTEXT", ConvertExpr: "CONVERT(current_value USING binary)", JsonExtractExpr: "JSON_UNQUOTE(JSON_EXTRACT(value_array, CONCAT('$[', i, ']')))", WireType: 2, WriteFunction: "", SupportsPacked: false},
		{ProtoType: "message", SqlType: "LONGBLOB", ConvertExpr: "current_value", JsonExtractExpr: "FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(value_array, CONCAT('$[', i, ']'))))", WireType: 2, WriteFunction: "", SupportsPacked: false},
	}

	bulkTemplateText := dedent.Pipe(`
		|
		|DROP FUNCTION IF EXISTS _pb_wire_json_add_all_repeated_{{.ProtoType}}_field_elements $$
		|CREATE FUNCTION _pb_wire_json_add_all_repeated_{{.ProtoType}}_field_elements(wire_json JSON, field_number INT, value_array JSON{{if .SupportsPacked}}, use_packed BOOLEAN{{end}}) RETURNS JSON DETERMINISTIC
		|BEGIN
		|	DECLARE i INT DEFAULT 0;
		|	DECLARE array_length INT;
		|	DECLARE current_value {{.SqlType}};
		|	DECLARE result JSON;
		|{{- if .SupportsPacked}}
		|	DECLARE packed_data LONGBLOB DEFAULT '';
		|	DECLARE temp_encoded LONGBLOB;
		|	DECLARE new_element JSON;
		|	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
		|	DECLARE field_array JSON;
		|	DECLARE next_index INT;
		|{{- end}}
		|
		|	SET result = wire_json;
		|	SET array_length = JSON_LENGTH(value_array);
		|
		|	IF array_length = 0 THEN
		|		RETURN result;
		|	END IF;
		|
		|{{- if .SupportsPacked}}
		|	IF use_packed THEN
		|		-- Build packed data
		|		WHILE i < array_length DO
		|			SET current_value = {{.JsonExtractExpr}};
		|			CALL {{.WriteFunction}}({{.ConvertExpr}}, temp_encoded);
		|			SET packed_data = CONCAT(packed_data, temp_encoded);
		|			SET i = i + 1;
		|		END WHILE;
		|
		|		-- Get existing field array
		|		SET field_array = JSON_EXTRACT(result, field_path);
		|
		|		-- Check if we can append to existing packed field
		|		IF field_array IS NOT NULL AND JSON_LENGTH(field_array) > 0 THEN
		|			-- Check if last element is packed (wire type 2)
		|			IF JSON_EXTRACT(field_array, CONCAT('$[', JSON_LENGTH(field_array) - 1, '].t')) = 2 THEN
		|				-- Append to existing packed data
		|				SET packed_data = CONCAT(
		|					FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(field_array, CONCAT('$[', JSON_LENGTH(field_array) - 1, '].v')))),
		|					packed_data
		|				);
		|				RETURN JSON_SET(result, CONCAT(field_path, '[', JSON_LENGTH(field_array) - 1, '].v'), TO_BASE64(packed_data));
		|			END IF;
		|		END IF;
		|
		|		-- Create new packed element
		|		SET next_index = _pb_wire_json_get_next_index(result);
		|		SET new_element = JSON_OBJECT('i', next_index, 'n', field_number, 't', 2, 'v', TO_BASE64(packed_data));
		|
		|		IF field_array IS NULL THEN
		|			RETURN JSON_SET(result, field_path, JSON_ARRAY(new_element));
		|		ELSE
		|			RETURN JSON_ARRAY_APPEND(result, field_path, new_element);
		|		END IF;
		|	ELSE
		|		-- Add unpacked elements
		|		WHILE i < array_length DO
		|			SET current_value = {{.JsonExtractExpr}};
		|			SET result = _pb_wire_json_add_repeated_{{if eq .WireType 0}}varint{{else if eq .WireType 1}}i64{{else if eq .WireType 5}}i32{{else if eq .WireType 2}}len{{end}}_field_element(result, field_number, {{.ConvertExpr}}, FALSE);
		|			SET i = i + 1;
		|		END WHILE;
		|		RETURN result;
		|	END IF;
		|{{- else}}
		|	-- Non-packable types - always add as separate elements
		|	WHILE i < array_length DO
		|		SET current_value = {{.JsonExtractExpr}};
		|		SET result = _pb_wire_json_add_repeated_len_field_element(result, field_number, {{.ConvertExpr}});
		|		SET i = i + 1;
		|	END WHILE;
		|	RETURN result;
		|{{- end}}
		|END $$
	`)

	tmpl, err := template.New("bulk").Parse(bulkTemplateText)
	if err != nil {
		panic(err)
	}

	for _, accessor := range bulkAccessors {
		if err := tmpl.Execute(os.Stdout, accessor); err != nil {
			panic(err)
		}
	}
}

func main() {
	cmd := &cli.Command{
		Name:   "protobuf-accessors",
		Usage:  "Generates protobuf-accessors.sql",
		Action: generateAccessorsAction,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
