DELIMITER $$

-- Convert Struct JSON to number JSON format (stored procedure)
DROP PROCEDURE IF EXISTS _pb_wkt_struct_json_to_number_json $$
CREATE PROCEDURE _pb_wkt_struct_json_to_number_json(IN proto_json_value JSON, OUT result JSON)
BEGIN
	DECLARE struct_keys JSON;
	DECLARE struct_key_count INT;
	DECLARE struct_key_index INT;
	DECLARE struct_key_name TEXT;
	DECLARE struct_value_json JSON;
	DECLARE struct_converted_value JSON;
	DECLARE struct_result JSON;

	SET @@SESSION.max_sp_recursion_depth = 255;

	-- Convert Struct {key: value, key: value} to {"1": {field_map}}
	SET struct_keys = JSON_KEYS(proto_json_value);
	SET struct_key_count = JSON_LENGTH(struct_keys);
	SET struct_key_index = 0;
	SET struct_result = JSON_OBJECT();

	WHILE struct_key_index < struct_key_count DO
		SET struct_key_name = JSON_UNQUOTE(JSON_EXTRACT(struct_keys, CONCAT('$[', struct_key_index, ']')));
		SET struct_value_json = JSON_EXTRACT(proto_json_value, CONCAT('$.', JSON_QUOTE(struct_key_name)));
		-- Recursively convert the value as Value
		CALL _pb_wkt_value_json_to_number_json(struct_value_json, struct_converted_value);
		SET struct_result = JSON_SET(struct_result, CONCAT('$.', JSON_QUOTE(struct_key_name)), struct_converted_value);
		SET struct_key_index = struct_key_index + 1;
	END WHILE;

	-- Empty struct should result in empty object, not {"1": {}}, because default values (empty map field) are omitted in ProtoNumberJSON.
	IF struct_key_count = 0 THEN
		SET result = JSON_OBJECT();
	ELSE
		SET result = JSON_OBJECT('1', struct_result);
	END IF;
END $$

-- Convert Struct JSON to number JSON format (function wrapper)
DROP FUNCTION IF EXISTS _pb_wkt_struct_json_to_number_json $$
CREATE FUNCTION _pb_wkt_struct_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wkt_struct_json_to_number_json(proto_json_value, result);
	RETURN result;
END $$

-- Convert ListValue JSON to number JSON format (stored procedure)
DROP PROCEDURE IF EXISTS _pb_wkt_list_value_json_to_number_json $$
CREATE PROCEDURE _pb_wkt_list_value_json_to_number_json(IN proto_json_value JSON, OUT result JSON)
BEGIN
	DECLARE list_length INT;
	DECLARE list_index INT;
	DECLARE list_element_json JSON;
	DECLARE list_converted_value JSON;
	DECLARE list_result JSON;

	SET @@SESSION.max_sp_recursion_depth = 255;

	-- Convert ListValue [value, value, value] to {"1": [values]}
	SET list_length = JSON_LENGTH(proto_json_value);
	SET list_index = 0;
	SET list_result = JSON_ARRAY();

	WHILE list_index < list_length DO
		SET list_element_json = JSON_EXTRACT(proto_json_value, CONCAT('$[', list_index, ']'));
		-- Recursively convert the element as Value
		CALL _pb_wkt_value_json_to_number_json(list_element_json, list_converted_value);
		SET list_result = JSON_ARRAY_APPEND(list_result, '$', list_converted_value);
		SET list_index = list_index + 1;
	END WHILE;

	-- Empty array should result in empty object, not {"1": []}, because default values (empty repeated field) are omitted in ProtoNumberJSON.
	IF list_length = 0 THEN
		SET result = JSON_OBJECT();
	ELSE
		SET result = JSON_OBJECT('1', list_result);
	END IF;
END $$

-- Convert ListValue JSON to number JSON format (function wrapper)
DROP FUNCTION IF EXISTS _pb_wkt_list_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_list_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wkt_list_value_json_to_number_json(proto_json_value, result);
	RETURN result;
END $$

-- Convert Value JSON to number JSON format (stored procedure)
DROP PROCEDURE IF EXISTS _pb_wkt_value_json_to_number_json $$
CREATE PROCEDURE _pb_wkt_value_json_to_number_json(IN proto_json_value JSON, OUT result JSON)
BEGIN
	DECLARE converted_value JSON;
	DECLARE message_text TEXT;

	SET @@SESSION.max_sp_recursion_depth = 255;

	CASE JSON_TYPE(proto_json_value)
	WHEN 'NULL' THEN
		-- null_value (field 1, enum value 0)
		SET result = JSON_OBJECT('1', 0);
	WHEN 'INTEGER' THEN
		-- number_value (field 2)
		SET result = JSON_OBJECT('2', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(CAST(proto_json_value AS DOUBLE))));
	WHEN 'UNSIGNED INTEGER' THEN
		-- number_value (field 2)
		SET result = JSON_OBJECT('2', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(CAST(proto_json_value AS DOUBLE))));
	WHEN 'DOUBLE' THEN
		-- number_value (field 2)
		SET result = JSON_OBJECT('2', _pb_convert_double_uint64_to_number_json(_pb_util_reinterpret_double_as_uint64(CAST(proto_json_value AS DOUBLE))));
	WHEN 'STRING' THEN
		-- string_value (field 3)
		SET result = JSON_OBJECT('3', proto_json_value);
	WHEN 'BOOLEAN' THEN
		-- bool_value (field 4)
		SET result = JSON_OBJECT('4', proto_json_value);
	WHEN 'OBJECT' THEN
		-- struct_value (field 5) - recursively convert as Struct
		CALL _pb_wkt_struct_json_to_number_json(proto_json_value, converted_value);
		SET result = JSON_OBJECT('5', converted_value);
	WHEN 'ARRAY' THEN
		-- list_value (field 6) - recursively convert as ListValue
		CALL _pb_wkt_list_value_json_to_number_json(proto_json_value, converted_value);
		SET result = JSON_OBJECT('6', converted_value);
	ELSE
		-- Unknown JSON type, signal error
		SET message_text = CONCAT('Unsupported JSON type for Value: ', JSON_TYPE(proto_json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$

-- Convert Value JSON to number JSON format (function wrapper)
DROP FUNCTION IF EXISTS _pb_wkt_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_value_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wkt_value_json_to_number_json(proto_json_value, result);
	RETURN result;
END $$

-- Helper function to convert google.protobuf.NullValue from ProtoJSON to ProtoNumberJSON
DROP FUNCTION IF EXISTS _pb_wkt_null_value_json_to_number_json $$
CREATE FUNCTION _pb_wkt_null_value_json_to_number_json(proto_json_value JSON) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;
	-- google.protobuf.NullValue in ProtoJSON can be:
	-- 1. JSON null -> should convert to enum value 0 (NULL_VALUE)
	-- 2. String "NULL_VALUE" -> should convert to enum value 0
	-- 3. Number 0 -> should convert to enum value 0

	IF JSON_TYPE(proto_json_value) = 'NULL' THEN
		-- JSON null represents NULL_VALUE (enum value 0)
		RETURN 0;
	ELSEIF JSON_TYPE(proto_json_value) = 'STRING' THEN
		-- String name "NULL_VALUE"
		IF JSON_UNQUOTE(proto_json_value) = 'NULL_VALUE' THEN
			RETURN 0;
		ELSE
			-- TODO: STRING '0'?
			-- Invalid string value for NullValue enum
			-- TODO: What if ignore_unknown_enums is set?
			SET message_text = CONCAT('Invalid NullValue enum string: ', JSON_UNQUOTE(proto_json_value));
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
	ELSEIF JSON_TYPE(proto_json_value) IN ('INTEGER', 'UNSIGNED INTEGER') THEN
		-- Numeric value (should be 0 for NULL_VALUE)
		RETURN proto_json_value;
	ELSE
		-- Invalid JSON type for NullValue
		SET message_text = CONCAT('Invalid JSON type for NullValue: ', JSON_TYPE(proto_json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;
END $$

-- Helper function to convert google.protobuf.NullValue from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_null_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_null_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;
	DECLARE enum_value INT;

	-- google.protobuf.NullValue in ProtoNumberJSON is just the enum numeric value
	-- It should always be 0 (NULL_VALUE), and converts back to JSON null

	IF JSON_TYPE(number_json_value) IN ('INTEGER', 'UNSIGNED INTEGER') THEN
		SET enum_value = CAST(number_json_value AS SIGNED);
		IF enum_value = 0 THEN
			-- NULL_VALUE (enum value 0) converts to JSON null
			RETURN CAST(NULL AS JSON);
		ELSE
			-- Invalid numeric value for NullValue enum
			SET message_text = CONCAT('Invalid NullValue enum number in ProtoNumberJSON: ', enum_value);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
	ELSE
		-- Invalid JSON type for NullValue in ProtoNumberJSON
		SET message_text = CONCAT('Invalid JSON type for NullValue in ProtoNumberJSON: ', JSON_TYPE(number_json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;
END $$

-- Helper procedure to convert google.protobuf.Value from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_wkt_value_number_json_to_json $$
CREATE PROCEDURE _pb_wkt_value_number_json_to_json(IN number_json_value JSON, OUT result JSON)
BEGIN
	DECLARE struct_converted_value JSON;
	DECLARE list_converted_value JSON;
	DECLARE uint64_value BIGINT UNSIGNED;

	SET @@SESSION.max_sp_recursion_depth = 255;

	-- Convert Value to its unwrapped form
	-- Value has oneof fields: null_value(1), number_value(2), string_value(3), bool_value(4), struct_value(5), list_value(6)
	IF JSON_LENGTH(number_json_value) = 0 THEN
		SET result = CAST(NULL AS JSON);
	ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."1"') THEN
		-- null_value
		SET result = CAST(NULL AS JSON);
	ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."2"') THEN
		-- number_value
		SET uint64_value = _pb_json_parse_double_as_uint64(JSON_EXTRACT(number_json_value, '$."2"'), TRUE);

		-- Check for special values and reject them as per google.protobuf.Value specification
		IF (uint64_value & 0x7FF0000000000000) = 0x7FF0000000000000 THEN
			-- Exponent is all 1s - this is NaN or Infinity, which google.protobuf.Value cannot represent
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'google.protobuf.Value cannot represent NaN or Infinity values';
		END IF;

		SET result = _pb_convert_double_uint64_to_json(uint64_value);
	ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."3"') THEN
		-- string_value
		SET result = JSON_EXTRACT(number_json_value, '$."3"');
	ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."4"') THEN
		-- bool_value
		SET result = JSON_EXTRACT(number_json_value, '$."4"');
	ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."5"') THEN
		-- struct_value - recursively convert
		CALL _pb_wkt_struct_number_json_to_json(JSON_EXTRACT(number_json_value, '$."5"'), struct_converted_value);
		SET result = struct_converted_value;
	ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."6"') THEN
		-- list_value - recursively convert
		CALL _pb_wkt_list_value_number_json_to_json(JSON_EXTRACT(number_json_value, '$."6"'), list_converted_value);
		SET result = list_converted_value;
	ELSE
		SET result = CAST(NULL AS JSON);
	END IF;
END $$

-- Helper function to convert google.protobuf.Value from ProtoNumberJSON to ProtoJSON (function wrapper)
DROP FUNCTION IF EXISTS _pb_wkt_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wkt_value_number_json_to_json(number_json_value, result);
	RETURN result;
END $$

-- Helper procedure to convert google.protobuf.Struct from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_wkt_struct_number_json_to_json $$
CREATE PROCEDURE _pb_wkt_struct_number_json_to_json(IN number_json_value JSON, OUT result JSON)
BEGIN
	DECLARE struct_fields JSON;
	DECLARE struct_keys JSON;
	DECLARE struct_key_count INT;
	DECLARE struct_key_index INT;
	DECLARE struct_key_name TEXT;
	DECLARE struct_value_json JSON;
	DECLARE struct_converted_value JSON;
	DECLARE struct_result JSON;

	SET @@SESSION.max_sp_recursion_depth = 255;

	-- Convert Struct {"1": {field_map}} to {key: value, key: value}
	SET struct_fields = JSON_EXTRACT(number_json_value, '$."1"');
	IF struct_fields IS NULL OR JSON_LENGTH(struct_fields) = 0 THEN
		SET result = JSON_OBJECT();
	ELSE
		SET struct_keys = JSON_KEYS(struct_fields);
		SET struct_key_count = JSON_LENGTH(struct_keys);
		SET struct_key_index = 0;
		SET struct_result = JSON_OBJECT();

		struct_loop: WHILE struct_key_index < struct_key_count DO
			SET struct_key_name = JSON_UNQUOTE(JSON_EXTRACT(struct_keys, CONCAT('$[', struct_key_index, ']')));
			SET struct_value_json = JSON_EXTRACT(struct_fields, CONCAT('$.', JSON_QUOTE(struct_key_name)));
			-- Recursively convert the Value
			CALL _pb_wkt_value_number_json_to_json(struct_value_json, struct_converted_value);
			SET struct_result = JSON_SET(struct_result, CONCAT('$.', JSON_QUOTE(struct_key_name)), struct_converted_value);
			SET struct_key_index = struct_key_index + 1;
		END WHILE struct_loop;

		SET result = struct_result;
	END IF;
END $$

-- Helper function to convert google.protobuf.Struct from ProtoNumberJSON to ProtoJSON (function wrapper)
DROP FUNCTION IF EXISTS _pb_wkt_struct_number_json_to_json $$
CREATE FUNCTION _pb_wkt_struct_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wkt_struct_number_json_to_json(number_json_value, result);
	RETURN result;
END $$

-- Helper procedure to convert google.protobuf.ListValue from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_wkt_list_value_number_json_to_json $$
CREATE PROCEDURE _pb_wkt_list_value_number_json_to_json(IN number_json_value JSON, OUT result JSON)
BEGIN
	DECLARE list_values JSON;
	DECLARE list_length INT;
	DECLARE list_index INT;
	DECLARE list_element_json JSON;
	DECLARE list_converted_value JSON;
	DECLARE list_result JSON;

	SET @@SESSION.max_sp_recursion_depth = 255;

	-- Convert ListValue {"1": [values]} to [value, value, value]
	SET list_values = JSON_EXTRACT(number_json_value, '$."1"');
	IF list_values IS NULL OR JSON_LENGTH(list_values) = 0 THEN
		SET result = JSON_ARRAY();
	ELSE
		SET list_length = JSON_LENGTH(list_values);
		SET list_index = 0;
		SET list_result = JSON_ARRAY();

		list_loop: WHILE list_index < list_length DO
			SET list_element_json = JSON_EXTRACT(list_values, CONCAT('$[', list_index, ']'));
			-- Recursively convert the Value
			CALL _pb_wkt_value_number_json_to_json(list_element_json, list_converted_value);
			SET list_result = JSON_ARRAY_APPEND(list_result, '$', list_converted_value);
			SET list_index = list_index + 1;
		END WHILE list_loop;

		SET result = list_result;
	END IF;
END $$

-- Helper function to convert google.protobuf.ListValue from ProtoNumberJSON to ProtoJSON (function wrapper)
DROP FUNCTION IF EXISTS _pb_wkt_list_value_number_json_to_json $$
CREATE FUNCTION _pb_wkt_list_value_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wkt_list_value_number_json_to_json(number_json_value, result);
	RETURN result;
END $$
