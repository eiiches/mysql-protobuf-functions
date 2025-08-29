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
	DECLARE uint32_value INT UNSIGNED;

	IF proto_json_value IS NULL THEN
		SET uint32_value = 0;
	ELSE
		SET uint32_value = _pb_json_parse_float_as_uint32(proto_json_value, FALSE);
	END IF;

	IF uint32_value <> 0 THEN
		RETURN JSON_OBJECT('1', _pb_convert_float_uint32_to_number_json(uint32_value));
	ELSE
		RETURN JSON_OBJECT();
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wkt_double_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_double_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE uint64_value BIGINT UNSIGNED;

	IF proto_json_value IS NULL THEN
		SET uint64_value = 0;
	ELSE
		SET uint64_value = _pb_json_parse_double_as_uint64(proto_json_value, FALSE);
	END IF;

	IF uint64_value <> 0 THEN
		RETURN JSON_OBJECT('1', _pb_convert_double_uint64_to_number_json(uint64_value));
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
		RETURN JSON_OBJECT('1', _pb_to_base64(bytes_val));
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
		SET bool_val = _pb_json_parse_bool(proto_json_value);
	END IF;

	IF bool_val THEN
		RETURN JSON_OBJECT('1', CAST('true' AS JSON));
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

-- Helper function to convert StringValue from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_string_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_string_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	-- {"1": "value"} becomes unwrapped "value", {} becomes ""
	IF JSON_LENGTH(number_json_value) = 0 THEN
		RETURN JSON_QUOTE('');
	ELSE
		RETURN JSON_EXTRACT(number_json_value, '$."1"');
	END IF;
END $$

-- Helper function to convert Int64Value from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_int64_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_int64_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE wrapped_value JSON;
	-- {"1": value} becomes unwrapped "value" (as string for 64-bit), {} becomes "0"
	IF JSON_LENGTH(number_json_value) = 0 THEN
		RETURN JSON_QUOTE('0');
	ELSE
		SET wrapped_value = JSON_EXTRACT(number_json_value, '$."1"');
		RETURN JSON_QUOTE(CAST(wrapped_value AS CHAR));
	END IF;
END $$

-- Helper function to convert UInt64Value from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_uint64_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_uint64_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE wrapped_value JSON;
	-- {"1": value} becomes unwrapped "value" (as string for 64-bit), {} becomes "0"
	IF JSON_LENGTH(number_json_value) = 0 THEN
		RETURN JSON_QUOTE('0');
	ELSE
		SET wrapped_value = JSON_EXTRACT(number_json_value, '$."1"');
		RETURN JSON_QUOTE(CAST(wrapped_value AS CHAR));
	END IF;
END $$

-- Helper function to convert Int32Value from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_int32_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_int32_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	-- {"1": value} becomes unwrapped value (as number for 32-bit), {} becomes 0
	IF JSON_LENGTH(number_json_value) = 0 THEN
		RETURN CAST(0 AS JSON);
	ELSE
		RETURN JSON_EXTRACT(number_json_value, '$."1"');
	END IF;
END $$

-- Helper function to convert UInt32Value from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_uint32_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_uint32_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	-- {"1": value} becomes unwrapped value (as number for 32-bit), {} becomes 0
	IF JSON_LENGTH(number_json_value) = 0 THEN
		RETURN CAST(0 AS JSON);
	ELSE
		RETURN JSON_EXTRACT(number_json_value, '$."1"');
	END IF;
END $$

-- Helper function to convert BoolValue from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_bool_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_bool_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	-- {"1": value} becomes unwrapped value, {} becomes false
	IF JSON_LENGTH(number_json_value) = 0 THEN
		RETURN CAST(false AS JSON);
	ELSE
		RETURN JSON_EXTRACT(number_json_value, '$."1"');
	END IF;
END $$

-- Helper function to convert FloatValue from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_float_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_float_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_value JSON;
	DECLARE uint32_value INT UNSIGNED;

	-- {"1": value} becomes unwrapped value, {} becomes 0.0
	IF JSON_LENGTH(number_json_value) = 0 THEN
		RETURN CAST(0.0 AS JSON);
	ELSE
		SET field_value = JSON_EXTRACT(number_json_value, '$."1"');
		-- Parse the field value (which may be in binary32 format) and convert to uint32
		SET uint32_value = _pb_json_parse_float_as_uint32(field_value, TRUE);
		-- Convert uint32 to proper JSON representation
		RETURN _pb_convert_float_uint32_to_json(uint32_value);
	END IF;
END $$

-- Helper function to convert DoubleValue from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_double_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_double_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_value JSON;
	DECLARE uint64_value BIGINT UNSIGNED;

	-- {"1": value} becomes unwrapped value, {} becomes 0.0
	IF JSON_LENGTH(number_json_value) = 0 THEN
		RETURN CAST(0.0 AS JSON);
	ELSE
		SET field_value = JSON_EXTRACT(number_json_value, '$."1"');
		-- Parse the field value (which may be in binary64 format) and convert to uint64
		SET uint64_value = _pb_json_parse_double_as_uint64(field_value, TRUE);
		-- Convert uint64 to proper JSON representation
		RETURN _pb_convert_double_uint64_to_json(uint64_value);
	END IF;
END $$

-- Helper function to convert BytesValue from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_bytes_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_bytes_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	-- {"1": "value"} becomes unwrapped "value", {} becomes ""
	IF JSON_LENGTH(number_json_value) = 0 THEN
		RETURN JSON_QUOTE('');
	ELSE
		RETURN JSON_EXTRACT(number_json_value, '$."1"');
	END IF;
END $$
