package main

import (
	"context"
	"fmt"
	"github.com/eiiches/mysql-protobuf-functions/internal/dedent"
	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"text/template"
)

type AccessProcedure struct {
	Name    string
	SqlType string
}

type Accessor struct {
	ProtoType  string
	SqlType    string
	ReturnExpr string
	Procedure  *AccessProcedure
	Input      *Input
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
}

func generateRepeatedNumbersAsJson() {

	accessors := []*RepeatedAsJsonAccessor{
		{
			ProtoType:           "int32",
			SqlType:             "INT",
			Expr:                "_pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_int64(uint_value))",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
		},
		{
			ProtoType:           "uint32",
			SqlType:             "INT UNSIGNED",
			Expr:                "_pb_util_cast_int64_as_int32(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
		},
		{
			ProtoType:           "int64",
			SqlType:             "BIGINT",
			Expr:                "_pb_util_reinterpret_uint64_as_int64(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
		},
		{
			ProtoType:           "uint64",
			SqlType:             "BIGINT UNSIGNED",
			Expr:                "uint_value",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
		},
		{
			ProtoType:           "sint32",
			SqlType:             "INT",
			Expr:                "_pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_sint64(uint_value))",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
		},
		{
			ProtoType:           "sint64",
			SqlType:             "BIGINT",
			Expr:                "_pb_util_reinterpret_uint64_as_sint64(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
		},
		{
			ProtoType:           "enum",
			SqlType:             "INT",
			Expr:                "_pb_util_reinterpret_uint64_as_int64(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
		},
		{
			ProtoType:           "bool",
			SqlType:             "BOOLEAN",
			Expr:                "uint_value <> 0",
			PackedUint64Decoder: "_pb_wire_read_varint_as_uint64",
			WireType:            0,
		},
		{
			ProtoType:           "fixed32",
			SqlType:             "INT UNSIGNED",
			Expr:                "uint_value",
			PackedUint64Decoder: "_pb_wire_read_i32_as_uint32",
			WireType:            5,
		},
		{
			ProtoType:           "sfixed32",
			SqlType:             "INT",
			Expr:                "_pb_util_reinterpret_uint32_as_int32(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_i32_as_uint32",
			WireType:            5,
		},
		{
			ProtoType:           "float",
			SqlType:             "FLOAT",
			Expr:                "_pb_util_reinterpret_uint32_as_float(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_i32_as_uint32",
			WireType:            5,
		},
		{
			ProtoType:           "fixed64",
			SqlType:             "BIGINT UNSIGNED",
			Expr:                "uint_value",
			PackedUint64Decoder: "_pb_wire_read_i64_as_uint64",
			WireType:            1,
		},
		{
			ProtoType:           "sfixed64",
			SqlType:             "BIGINT",
			Expr:                "_pb_util_reinterpret_uint64_as_int64(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_i64_as_uint64",
			WireType:            1,
		},
		{
			ProtoType:           "double",
			SqlType:             "DOUBLE",
			Expr:                "_pb_util_reinterpret_uint64_as_double(uint_value)",
			PackedUint64Decoder: "_pb_wire_read_i64_as_uint64",
			WireType:            1,
		},
		{
			ProtoType:           "bytes",
			SqlType:             "LONGBLOB",
			Expr:                "TO_BASE64(bytes_value)",
			PackedUint64Decoder: "",
			WireType:            2,
		},
		{
			ProtoType:           "string",
			SqlType:             "LONGTEXT",
			Expr:                "CONVERT(bytes_value USING utf8mb4)",
			PackedUint64Decoder: "",
			WireType:            2,
		},
		{
			ProtoType:           "message",
			SqlType:             "LONGBLOB",
			Expr:                "TO_BASE64(bytes_value)",
			PackedUint64Decoder: "",
			WireType:            2,
		},
	}

	templateText := dedent.Pipe(`
		|
		|DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_{{.ProtoType}}_field_as_json_array $$
		|CREATE PROCEDURE _pb_wire_json_get_repeated_{{.ProtoType}}_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
		|BEGIN
		|	DECLARE done TINYINT DEFAULT FALSE;
		|	DECLARE message_text TEXT;
		|	DECLARE uint_value BIGINT UNSIGNED;
		|	DECLARE bytes_value LONGBLOB;
		|	DECLARE wire_type INT;
		|
		|	DECLARE cur CURSOR FOR
		|		SELECT
		|			jt.wire_type,
		|			jt.uint_value,
		|			FROM_BASE64(jt.bytes_value)
		|		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
		|			field_number INT PATH '$.field_number',
		|			wire_type INT PATH '$.wire_type',
		|			uint_value BIGINT UNSIGNED PATH '$.value.uint',
		|			bytes_value TEXT PATH '$.value.bytes'
		|		)) AS jt
		|		WHERE jt.field_number = field_number;
		|	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
		|
		|	SET result = JSON_ARRAY();
		|
		|	OPEN cur;
		|	l1: LOOP
		|		FETCH cur INTO wire_type, uint_value, bytes_value;
		|		IF done THEN
		|			LEAVE l1;
		|		END IF;
		|
		|		CASE wire_type
		|		WHEN {{.WireType}} THEN
		|			SET result = JSON_ARRAY_APPEND(result, '$', {{.Expr}});
		|{{- if .PackedUint64Decoder }}
		|		WHEN 2 THEN -- LEN
		|			WHILE LENGTH(bytes_value) <> 0 DO
		|				CALL {{.PackedUint64Decoder}}(bytes_value, uint_value, bytes_value);
		|				SET result = JSON_ARRAY_APPEND(result, '$', {{.Expr}});
		|			END WHILE;
		|{{- end }}
		|		ELSE
		|			SET message_text = CONCAT('_pb_wire_json_get_repeated_{{.ProtoType}}_field_as_json_array: unexpected wire_type (', wire_type, ')');
		|			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		|		END CASE;
		|	END LOOP;
		|	CLOSE cur;
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_{{.ProtoType}}_field_as_json_array $$
		|CREATE FUNCTION pb_wire_json_get_repeated_{{.ProtoType}}_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
		|BEGIN
		|	DECLARE result JSON;
		|	CALL _pb_wire_json_get_repeated_{{.ProtoType}}_field_as_json_array(wire_json, field_number, result);
		|	RETURN result;
		|END $$
		|
		|DROP PROCEDURE IF EXISTS _pb_message_get_repeated_{{.ProtoType}}_field_as_json_array $$
		|CREATE PROCEDURE _pb_message_get_repeated_{{.ProtoType}}_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
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
		|			LEAVE l1;
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
		|			SET message_text = CONCAT('_pb_message_get_repeated_{{.ProtoType}}_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
		|			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		|		END CASE;
		|	END WHILE;
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_message_get_repeated_{{.ProtoType}}_field_as_json_array $$
		|CREATE FUNCTION pb_message_get_repeated_{{.ProtoType}}_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
		|BEGIN
		|	DECLARE result JSON;
		|	CALL _pb_message_get_repeated_{{.ProtoType}}_field_as_json_array(message, field_number, result);
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
		|	CALL {{.Procedure.Name}}({{.Input.Name}}, field_number, NULL, value, field_count);
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
		|	CALL {{.Procedure.Name}}({{.Input.Name}}, field_number, NULL, value, field_count);
		|	RETURN field_count > 0;
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_get_repeated_{{.ProtoType}}_field $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_get_repeated_{{.ProtoType}}_field({{.Input.Name}} {{.Input.SqlType}}, field_number INT, repeated_index INT) RETURNS {{.SqlType}} DETERMINISTIC
		|BEGIN
		|	DECLARE value {{.Procedure.SqlType}};
		|	DECLARE field_count INT;
		|	CALL {{.Procedure.Name}}({{.Input.Name}}, field_number, repeated_index, value, field_count);
		|	RETURN {{.ReturnExpr}};
		|END $$
		|
		|DROP FUNCTION IF EXISTS pb_{{.Input.Kind}}_get_repeated_{{.ProtoType}}_field_count $$
		|CREATE FUNCTION pb_{{.Input.Kind}}_get_repeated_{{.ProtoType}}_field_count({{.Input.Name}} {{.Input.SqlType}}, field_number INT) RETURNS INT DETERMINISTIC
		|BEGIN
		|	DECLARE value {{.Procedure.SqlType}};
		|	DECLARE field_count INT;
		|	CALL {{.Procedure.Name}}({{.Input.Name}}, field_number, -1, value, field_count);
		|	RETURN field_count;
		|END $$
	`)

	os.Stdout.WriteString("DELIMITER $$\n")

	for _, input := range inputs {

		getVarintFieldAsUint64 := &AccessProcedure{
			Name:    fmt.Sprintf("_pb_%s_get_varint_field_as_uint64", input.Kind),
			SqlType: "BIGINT UNSIGNED",
		}

		getI64FieldAsUint64 := &AccessProcedure{
			Name:    fmt.Sprintf("_pb_%s_get_i64_field_as_uint64", input.Kind),
			SqlType: "BIGINT UNSIGNED",
		}

		getI32FieldAsUint64 := &AccessProcedure{
			Name:    fmt.Sprintf("_pb_%s_get_i32_field_as_uint32", input.Kind),
			SqlType: "INT UNSIGNED",
		}

		getLengthDelimitedField := &AccessProcedure{
			Name:    fmt.Sprintf("_pb_%s_get_len_type_field", input.Kind),
			SqlType: "LONGBLOB",
		}

		accessors := []*Accessor{
			// VARINT
			{Input: input, ProtoType: "int32", SqlType: "INT", ReturnExpr: "_pb_util_reinterpret_uint64_as_int64(value)", Procedure: getVarintFieldAsUint64},
			{Input: input, ProtoType: "int64", SqlType: "BIGINT", ReturnExpr: "_pb_util_reinterpret_uint64_as_int64(value)", Procedure: getVarintFieldAsUint64},
			{Input: input, ProtoType: "uint32", SqlType: "INT UNSIGNED", ReturnExpr: "value", Procedure: getVarintFieldAsUint64},
			{Input: input, ProtoType: "uint64", SqlType: "BIGINT UNSIGNED", ReturnExpr: "value", Procedure: getVarintFieldAsUint64},
			{Input: input, ProtoType: "sint32", SqlType: "INT", ReturnExpr: "_pb_util_reinterpret_uint64_as_sint64(value)", Procedure: getVarintFieldAsUint64},
			{Input: input, ProtoType: "sint64", SqlType: "BIGINT", ReturnExpr: "_pb_util_reinterpret_uint64_as_sint64(value)", Procedure: getVarintFieldAsUint64},
			{Input: input, ProtoType: "enum", SqlType: "INT", ReturnExpr: "_pb_util_reinterpret_uint64_as_int64(value)", Procedure: getVarintFieldAsUint64},
			{Input: input, ProtoType: "bool", SqlType: "BOOLEAN", ReturnExpr: "value <> 0", Procedure: getVarintFieldAsUint64},

			// I32
			{Input: input, ProtoType: "fixed32", SqlType: "INT UNSIGNED", ReturnExpr: "value", Procedure: getI32FieldAsUint64},
			{Input: input, ProtoType: "sfixed32", SqlType: "INT", ReturnExpr: "_pb_util_reinterpret_uint32_as_int32(value)", Procedure: getI32FieldAsUint64},
			{Input: input, ProtoType: "float", SqlType: "FLOAT", ReturnExpr: "_pb_util_reinterpret_uint32_as_float(value)", Procedure: getI32FieldAsUint64},

			// I64
			{Input: input, ProtoType: "fixed64", SqlType: "BIGINT UNSIGNED", ReturnExpr: "value", Procedure: getI64FieldAsUint64},
			{Input: input, ProtoType: "sfixed64", SqlType: "BIGINT", ReturnExpr: "_pb_util_reinterpret_uint64_as_int64(value)", Procedure: getI64FieldAsUint64},
			{Input: input, ProtoType: "double", SqlType: "DOUBLE", ReturnExpr: "_pb_util_reinterpret_uint64_as_double(value)", Procedure: getI64FieldAsUint64},

			// LEN
			{Input: input, ProtoType: "bytes", SqlType: "LONGBLOB", ReturnExpr: "value", Procedure: getLengthDelimitedField},
			{Input: input, ProtoType: "string", SqlType: "LONGTEXT", ReturnExpr: "CONVERT(value USING utf8mb4)", Procedure: getLengthDelimitedField},
			{Input: input, ProtoType: "message", SqlType: "LONGBLOB", ReturnExpr: "value", Procedure: getLengthDelimitedField},
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
	}

	generateRepeatedNumbersAsJson()

	return nil
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
