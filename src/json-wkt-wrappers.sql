DELIMITER $$

DROP FUNCTION IF EXISTS _pb_wkt_int64_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_int64_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE int64_val BIGINT;

	IF proto_json_value IS NULL THEN
		SET int64_val = 0;
	ELSE
		SET int64_val = _pb_json_parse_signed_int(proto_json_value);
	END IF;

	IF int64_val != 0 THEN
		RETURN JSON_OBJECT('1', int64_val);
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wkt_uint64_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_uint64_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE uint64_val BIGINT UNSIGNED;

	IF proto_json_value IS NULL THEN
		SET uint64_val = 0;
	ELSE
		SET uint64_val = _pb_json_parse_unsigned_int(proto_json_value);
	END IF;

	IF uint64_val != 0 THEN
		RETURN JSON_OBJECT('1', uint64_val);
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wkt_int32_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_int32_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE int32_val INT;

	IF proto_json_value IS NULL THEN
		SET int32_val = 0;
	ELSE
		SET int32_val = _pb_json_parse_signed_int(proto_json_value);
	END IF;

	IF int32_val != 0 THEN
		RETURN JSON_OBJECT('1', int32_val);
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wkt_uint32_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_uint32_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE uint32_val INT UNSIGNED;

	IF proto_json_value IS NULL THEN
		SET uint32_val = 0;
	ELSE
		SET uint32_val = _pb_json_parse_unsigned_int(proto_json_value);
	END IF;

	IF uint32_val != 0 THEN
		RETURN JSON_OBJECT('1', uint32_val);
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wkt_float_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_float_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE float_val JSON;

	IF proto_json_value IS NULL THEN
		SET float_val = CAST(0.0 AS JSON);
	ELSE
		SET float_val = _pb_json_parse_float(proto_json_value);
	END IF;

	IF float_val != CAST(0.0 AS JSON) THEN
		RETURN JSON_OBJECT('1', float_val);
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wkt_double_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_double_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE double_val JSON;

	IF proto_json_value IS NULL THEN
		SET double_val = CAST(0.0 AS JSON);
	ELSE
		SET double_val = _pb_json_parse_double(proto_json_value);
	END IF;

	IF double_val != CAST(0.0 AS JSON) THEN
		RETURN JSON_OBJECT('1', double_val);
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wkt_bytes_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_bytes_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE bytes_val LONGBLOB;

	IF proto_json_value IS NULL THEN
		SET str_value = '';
	ELSE
		SET str_value = JSON_UNQUOTE(proto_json_value);
	END IF;

	IF str_value != '' THEN
		-- Use the parsing function to decode and get standard Base64
		SET bytes_val = _pb_json_parse_bytes(proto_json_value);
		RETURN JSON_OBJECT('1', TO_BASE64(bytes_val));
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wkt_bool_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_bool_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE bool_val BOOLEAN;

	IF proto_json_value IS NULL THEN
		SET bool_val = FALSE;
	ELSE
		SET bool_val = CAST(proto_json_value AS UNSIGNED);
	END IF;

	IF bool_val != FALSE THEN
		RETURN JSON_OBJECT('1', bool_val);
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wkt_string_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_string_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;

	IF proto_json_value IS NULL THEN
		SET str_value = '';
	ELSE
		SET str_value = JSON_UNQUOTE(proto_json_value);
	END IF;

	IF str_value != '' THEN
		RETURN JSON_OBJECT('1', proto_json_value);
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$
