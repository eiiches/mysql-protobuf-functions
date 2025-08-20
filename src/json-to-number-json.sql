DELIMITER $$

-- Helper function to parse JSON value as BIGINT with validation
DROP FUNCTION IF EXISTS _pb_json_parse_signed_int $$
CREATE FUNCTION _pb_json_parse_signed_int(json_value JSON) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;
	DECLARE double_value DOUBLE;
	DECLARE uint_value BIGINT UNSIGNED;

	IF JSON_TYPE(json_value) = 'STRING' THEN
		SET str_value = JSON_UNQUOTE(json_value);

		-- Early return for invalid inputs
		IF str_value = '' OR NOT (str_value REGEXP '^-?[0-9]+(\\.[0-9]+)?([eE][+-]?[0-9]+)?$') THEN
			SET message_text = CONCAT('Invalid number format for signed integer: ', str_value);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		SET json_value = CAST(str_value AS JSON);
	END IF;

	CASE JSON_TYPE(json_value)
	WHEN 'UNSIGNED INTEGER' THEN
		-- Cast to unsigned first, then check range for signed integer
		SET uint_value = CAST(json_value AS UNSIGNED);
		IF uint_value > 9223372036854775807 THEN
			SET message_text = CONCAT('Value too large for signed integer field: ', CAST(uint_value AS CHAR));
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		RETURN CAST(uint_value AS SIGNED);
	WHEN 'INTEGER' THEN
		-- JSON INTEGER type should be safe for signed integers
		RETURN CAST(json_value AS SIGNED);
	WHEN 'DOUBLE' THEN
		-- Handle JSON numbers with exponential notation (e.g., 1e5)
		-- But reject non-integer values (e.g., 0.5)
		SET double_value = CAST(json_value AS DOUBLE);
		IF double_value != FLOOR(double_value) THEN
			SET message_text = CONCAT('Non-integer value for signed integer field: ', CAST(double_value AS CHAR));
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		RETURN CAST(double_value AS SIGNED);
	ELSE
		RETURN CAST(json_value AS SIGNED);
	END CASE;
END $$

-- Helper function to parse JSON value as BIGINT UNSIGNED with validation
DROP FUNCTION IF EXISTS _pb_json_parse_unsigned_int $$
CREATE FUNCTION _pb_json_parse_unsigned_int(json_value JSON) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;
	DECLARE double_value DOUBLE;

	IF JSON_TYPE(json_value) = 'STRING' THEN
		SET str_value = JSON_UNQUOTE(json_value);

		-- Early return for invalid inputs
		IF str_value = '' OR NOT (str_value REGEXP '^-?[0-9]+(\\.[0-9]+)?([eE][+-]?[0-9]+)?$') THEN
			SET message_text = CONCAT('Invalid number format for signed integer: ', str_value);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		SET json_value = CAST(str_value AS JSON);
	END IF;

	CASE JSON_TYPE(json_value)
	WHEN 'INTEGER' THEN
		-- Check if the INTEGER value is negative (invalid for unsigned)
		IF CAST(json_value AS SIGNED) < 0 THEN
			SET message_text = CONCAT('Negative value for unsigned integer field: ', CAST(json_value AS CHAR));
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		RETURN CAST(json_value AS UNSIGNED);
	WHEN 'UNSIGNED INTEGER' THEN
		-- UNSIGNED INTEGER values are always valid for unsigned fields
		RETURN CAST(json_value AS UNSIGNED);
	WHEN 'DOUBLE' THEN
		-- Handle JSON numbers with exponential notation (e.g., 1e5)
		-- But reject non-integer values (e.g., 0.5)
		SET double_value = CAST(json_value AS DOUBLE);
		IF double_value != FLOOR(double_value) THEN
			SET message_text = CONCAT('Non-integer value for unsigned integer field: ', CAST(double_value AS CHAR));
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		RETURN CAST(double_value AS UNSIGNED);
	ELSE
		RETURN CAST(json_value AS UNSIGNED);
	END CASE;
END $$

-- Helper function to parse JSON value as DOUBLE with validation
DROP FUNCTION IF EXISTS _pb_json_parse_double $$
CREATE FUNCTION _pb_json_parse_double(json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;

	IF JSON_TYPE(json_value) = 'STRING' THEN
		SET str_value = JSON_UNQUOTE(json_value);

		-- Reject empty strings
		IF str_value = '' THEN
			SET message_text = 'Empty string is not a valid number for double field';
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		-- Reject non-numeric strings
		IF NOT (str_value REGEXP '^[+-]?([0-9]*\\.?[0-9]+([eE][+-]?[0-9]+)?|Infinity|-Infinity|NaN)$') THEN
			SET message_text = CONCAT('Invalid number format for double field: ', str_value);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		-- Handle special values that can't be cast to JSON
		IF str_value IN ('Infinity', '-Infinity', 'NaN') THEN
			-- Return the special value as JSON string
			RETURN JSON_QUOTE(str_value);
		END IF;

		-- Convert string to JSON for further processing
		SET json_value = CAST(str_value AS JSON);
	END IF;

	CASE JSON_TYPE(json_value)
	WHEN 'INTEGER' THEN
		RETURN CAST(CAST(json_value AS DOUBLE) AS JSON);
	WHEN 'UNSIGNED INTEGER' THEN
		RETURN CAST(CAST(json_value AS DOUBLE) AS JSON);
	WHEN 'DECIMAL' THEN
		RETURN CAST(CAST(json_value AS DOUBLE) AS JSON);
	WHEN 'DOUBLE' THEN
		RETURN CAST(CAST(json_value AS DOUBLE) AS JSON);
	ELSE
		SET message_text = CONCAT('Invalid JSON type for double field: ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$

-- Helper function to parse JSON value as FLOAT with validation
DROP FUNCTION IF EXISTS _pb_json_parse_float $$
CREATE FUNCTION _pb_json_parse_float(json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;
	DECLARE double_value DOUBLE;

	IF JSON_TYPE(json_value) = 'STRING' THEN
		SET str_value = JSON_UNQUOTE(json_value);

		-- Reject empty strings
		IF str_value = '' THEN
			SET message_text = 'Empty string is not a valid number for float field';
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		-- Reject non-numeric strings
		IF NOT (str_value REGEXP '^[+-]?([0-9]*\\.?[0-9]+([eE][+-]?[0-9]+)?|Infinity|-Infinity|NaN)$') THEN
			SET message_text = CONCAT('Invalid number format for float field: ', str_value);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		-- Handle special values that can't be cast to JSON
		IF str_value IN ('Infinity', '-Infinity', 'NaN') THEN
			-- Return the special value as JSON string
			RETURN JSON_QUOTE(str_value);
		END IF;

		-- Convert string to JSON for further processing
		SET json_value = CAST(str_value AS JSON);
	END IF;

	CASE JSON_TYPE(json_value)
	WHEN 'INTEGER' THEN
		SET double_value = CAST(json_value AS DOUBLE);
	WHEN 'UNSIGNED INTEGER' THEN
		SET double_value = CAST(json_value AS DOUBLE);
	WHEN 'DECIMAL' THEN
		SET double_value = CAST(json_value AS DOUBLE);
	WHEN 'DOUBLE' THEN
		SET double_value = CAST(json_value AS DOUBLE);
	ELSE
		SET message_text = CONCAT('Invalid JSON type for float field: ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;

	-- Convert through float to get proper precision, then back to JSON
	RETURN CAST(CAST(double_value AS FLOAT) AS JSON);
END $$

-- Helper function to parse JSON value as bytes with Base64 decoding
DROP FUNCTION IF EXISTS _pb_json_parse_bytes $$
CREATE FUNCTION _pb_json_parse_bytes(json_value JSON) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;

	-- Bytes must be a JSON string
	IF JSON_TYPE(json_value) != 'STRING' THEN
		SET message_text = CONCAT('Invalid JSON type for bytes field: expected STRING, got ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Unquote the JSON string
	SET str_value = JSON_UNQUOTE(json_value);

	-- Decode from Base64/Base64URL and return as binary data
	RETURN _pb_util_from_base64_url(str_value);
END $$

-- Helper function to parse JSON value as Protobuf bool type
DROP FUNCTION IF EXISTS _pb_json_parse_bool $$
CREATE FUNCTION _pb_json_parse_bool(json_value JSON) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE bool_value BOOLEAN;
	DECLARE message_text TEXT;

	IF JSON_TYPE(json_value) <> 'BOOLEAN' THEN
		SET message_text = CONCAT('Invalid JSON type for bool field: expected BOOLEAN, got ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	SET bool_value = json_value;

	RETURN bool_value;
END $$

-- Helper function to check if a type is a protobuf wrapper type
DROP FUNCTION IF EXISTS _pb_is_wrapper_type $$
CREATE FUNCTION _pb_is_wrapper_type(full_type_name TEXT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	RETURN full_type_name IN (
		'.google.protobuf.DoubleValue',
		'.google.protobuf.FloatValue',
		'.google.protobuf.Int64Value',
		'.google.protobuf.UInt64Value',
		'.google.protobuf.Int32Value',
		'.google.protobuf.UInt32Value',
		'.google.protobuf.BoolValue',
		'.google.protobuf.StringValue',
		'.google.protobuf.BytesValue'
	);
END $$

-- Helper procedure to check if a type is a well-known type
DROP PROCEDURE IF EXISTS _pb_is_well_known_type $$
CREATE PROCEDURE _pb_is_well_known_type(IN full_type_name TEXT, OUT is_wkt BOOLEAN)
BEGIN
	IF full_type_name LIKE '.google.protobuf.%' THEN
		SET is_wkt = TRUE;
	ELSE
		SET is_wkt = FALSE;
	END IF;
END $$

-- Helper function to check if a well-known type has a zero/default value
DROP FUNCTION IF EXISTS _pb_is_wkt_zero_value $$
CREATE FUNCTION _pb_is_wkt_zero_value(full_type_name TEXT, converted_value JSON) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	CASE full_type_name
	WHEN '.google.protobuf.Timestamp' THEN
		-- Zero timestamp: empty object or seconds=0 and nanos=0
		RETURN JSON_LENGTH(converted_value) = 0 OR (JSON_EXTRACT(converted_value, '$."1"') = 0 AND JSON_EXTRACT(converted_value, '$."2"') = 0);
	WHEN '.google.protobuf.Duration' THEN
		-- Zero duration: empty object or seconds=0 and nanos=0
		RETURN JSON_LENGTH(converted_value) = 0 OR (JSON_EXTRACT(converted_value, '$."1"') = 0 AND JSON_EXTRACT(converted_value, '$."2"') = 0);
	WHEN '.google.protobuf.StringValue' THEN
		-- Empty string: empty object or empty string value
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_UNQUOTE(JSON_EXTRACT(converted_value, '$."1"')) = '';
	WHEN '.google.protobuf.Int64Value' THEN
		-- Zero value: empty object or zero value
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_EXTRACT(converted_value, '$."1"') = 0;
	WHEN '.google.protobuf.UInt64Value' THEN
		-- Zero value: empty object or zero value
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_EXTRACT(converted_value, '$."1"') = 0;
	WHEN '.google.protobuf.Int32Value' THEN
		-- Zero value: empty object or zero value
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_EXTRACT(converted_value, '$."1"') = 0;
	WHEN '.google.protobuf.UInt32Value' THEN
		-- Zero value: empty object or zero value
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_EXTRACT(converted_value, '$."1"') = 0;
	WHEN '.google.protobuf.BoolValue' THEN
		-- False value: empty object or false value
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_EXTRACT(converted_value, '$."1"') = false;
	WHEN '.google.protobuf.FloatValue' THEN
		-- Zero value: empty object or zero value
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_EXTRACT(converted_value, '$."1"') = 0.0;
	WHEN '.google.protobuf.DoubleValue' THEN
		-- Zero value: empty object or zero value
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_EXTRACT(converted_value, '$."1"') = 0.0;
	WHEN '.google.protobuf.BytesValue' THEN
		-- Empty bytes: empty object or empty string value
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_UNQUOTE(JSON_EXTRACT(converted_value, '$."1"')) = '';
	WHEN '.google.protobuf.Empty' THEN
		-- Empty is always considered zero
		RETURN TRUE;
	WHEN '.google.protobuf.FieldMask' THEN
		-- Empty paths array or no paths field
		RETURN JSON_LENGTH(converted_value) = 0 OR JSON_LENGTH(COALESCE(JSON_EXTRACT(converted_value, '$."1"'), JSON_ARRAY())) = 0;
	WHEN '.google.protobuf.Value' THEN
		-- Empty object means zero value
		RETURN JSON_LENGTH(converted_value) = 0;
	WHEN '.google.protobuf.Struct' THEN
		-- Empty Struct: either {} or {"1": {}}
		RETURN JSON_LENGTH(converted_value) = 0 OR
		       (JSON_CONTAINS_PATH(converted_value, 'one', '$."1"') AND JSON_LENGTH(JSON_EXTRACT(converted_value, '$."1"')) = 0);
	WHEN '.google.protobuf.ListValue' THEN
		-- Empty ListValue: either {} or {"1": []}
		RETURN JSON_LENGTH(converted_value) = 0 OR
		       (JSON_CONTAINS_PATH(converted_value, 'one', '$."1"') AND JSON_LENGTH(JSON_EXTRACT(converted_value, '$."1"')) = 0);
	ELSE
		-- Unknown WKT, don't consider it zero
		RETURN FALSE;
	END CASE;
END $$

-- Helper procedure to convert enum string name or numeric value to numeric value
DROP PROCEDURE IF EXISTS _pb_convert_json_enum_to_number $$
CREATE PROCEDURE _pb_convert_json_enum_to_number(
	IN descriptor_set_json JSON,
	IN full_enum_type_name TEXT,
	IN enum_string_value TEXT,
	OUT enum_numeric_value INT
)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';
	DECLARE message_text TEXT;
	DECLARE enum_descriptor JSON;
	DECLARE values_array JSON;
	DECLARE value_count INT;
	DECLARE value_index INT;
	DECLARE value_descriptor JSON;
	DECLARE value_name TEXT;
	DECLARE value_number INT;
	DECLARE input_as_number INT;
	DECLARE is_numeric BOOLEAN DEFAULT FALSE;

	-- Get enum descriptor
	SET enum_descriptor = _pb_get_enum_descriptor(descriptor_set_json, full_enum_type_name);

	IF enum_descriptor IS NULL THEN
		SET message_text = CONCAT('_pb_convert_json_enum_to_number: enum type not found: ', full_enum_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Check if the input is a numeric value (protobuf JSON allows both string names and numeric values)
	SET is_numeric = (enum_string_value REGEXP '^-?[0-9]+$');
	IF is_numeric THEN
		SET input_as_number = CAST(enum_string_value AS SIGNED);
	END IF;

	-- Get the values array (field 2 in EnumDescriptor)
	SET values_array = JSON_EXTRACT(enum_descriptor, '$."2"');
	SET value_count = JSON_LENGTH(values_array);
	SET value_index = 0;

	-- Search for either the string name or numeric value
	search_loop: WHILE value_index < value_count DO
		SET value_descriptor = JSON_EXTRACT(values_array, CONCAT('$[', value_index, ']'));
		SET value_name = JSON_UNQUOTE(JSON_EXTRACT(value_descriptor, '$."1"')); -- name field
		SET value_number = JSON_EXTRACT(value_descriptor, '$."2"'); -- number field

		-- Check for match by name or by numeric value
		IF value_name = enum_string_value OR (is_numeric AND value_number = input_as_number) THEN
			SET enum_numeric_value = value_number;
			LEAVE search_loop;
		END IF;

		SET value_index = value_index + 1;
	END WHILE search_loop;

	-- If not found, signal error
	IF value_index >= value_count THEN
		SET message_text = CONCAT('_pb_convert_json_enum_to_number: enum value not found: ', enum_string_value, ' in enum ', full_enum_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;
END $$

-- Helper function to convert well-known type from ProtoJSON to ProtoNumberJSON
DROP FUNCTION IF EXISTS _pb_convert_json_wkt_to_number_json $$
CREATE FUNCTION _pb_convert_json_wkt_to_number_json(
	full_type_name TEXT,
	proto_json_value JSON
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';
	DECLARE message_text TEXT;
	-- Variables for Any handling
	DECLARE type_url TEXT;
	DECLARE remaining_object JSON;

	CASE full_type_name
	WHEN '.google.protobuf.Timestamp' THEN
		RETURN _pb_wkt_timestamp_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Duration' THEN
		RETURN _pb_wkt_duration_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.FieldMask' THEN
		RETURN _pb_wkt_field_mask_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Value' THEN
		RETURN _pb_wkt_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Struct' THEN
		RETURN _pb_wkt_struct_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.ListValue' THEN
		RETURN _pb_wkt_list_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.StringValue' THEN
		RETURN _pb_wkt_string_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Int64Value' THEN
		RETURN _pb_wkt_int64_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.UInt64Value' THEN
		RETURN _pb_wkt_uint64_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Int32Value' THEN
		RETURN _pb_wkt_int32_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.UInt32Value' THEN
		RETURN _pb_wkt_uint32_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.BoolValue' THEN
		RETURN _pb_wkt_bool_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.FloatValue' THEN
		RETURN _pb_wkt_float_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.DoubleValue' THEN
		RETURN _pb_wkt_double_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.BytesValue' THEN
		RETURN _pb_wkt_bytes_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Empty' THEN
		RETURN _pb_wkt_empty_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Any' THEN
		RETURN _pb_wkt_any_json_to_number_json(proto_json_value);
	ELSE
		RETURN NULL;
	END CASE;
END $$

-- Helper procedure to convert singular field value to ProtoNumberJSON
DROP PROCEDURE IF EXISTS _pb_convert_singular_field_to_number_json $$
CREATE PROCEDURE _pb_convert_singular_field_to_number_json(
	IN descriptor_set_json JSON,
	IN field_type INT,
	IN field_type_name TEXT,
	IN field_json_value JSON,
	OUT converted_value JSON,
	OUT is_default BOOLEAN
)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';
	DECLARE message_text TEXT;
	-- Type-specific variables
	DECLARE enum_string_value TEXT;
	DECLARE enum_numeric_value INT;
	DECLARE int64_value BIGINT;
	DECLARE uint64_value BIGINT UNSIGNED;
	DECLARE int32_value INT;
	DECLARE uint32_value INT UNSIGNED;
	DECLARE double_json_value JSON;
	DECLARE float_json_value JSON;
	DECLARE str_value TEXT;
	DECLARE bool_value BOOLEAN;
	DECLARE nested_json JSON;

	CASE field_type
	WHEN 14 THEN -- enum
		SET enum_string_value = JSON_UNQUOTE(field_json_value);
		CALL _pb_convert_json_enum_to_number(descriptor_set_json, field_type_name, enum_string_value, enum_numeric_value);
		SET converted_value = CAST(enum_numeric_value AS JSON);
		SET is_default = (enum_numeric_value = 0);

	WHEN 11 THEN -- message
		IF field_json_value IS NULL THEN
			SET is_default = TRUE;
			SET converted_value = NULL;
		ELSE
			SET is_default = FALSE;
			SET converted_value = _pb_convert_json_wkt_to_number_json(field_type_name, field_json_value);
			IF converted_value IS NULL THEN
				CALL _pb_json_to_number_json_proc(descriptor_set_json, field_type_name, field_json_value, converted_value);
			END IF;
		END IF;

	WHEN 3 THEN -- int64 (convert string to number)
		SET int64_value = _pb_json_parse_signed_int(field_json_value);
		SET converted_value = CAST(int64_value AS JSON);
		SET is_default = (int64_value = 0);

	WHEN 4 THEN -- uint64 (convert string to number)
		SET uint64_value = _pb_json_parse_unsigned_int(field_json_value);
		SET converted_value = CAST(uint64_value AS JSON);
		SET is_default = (uint64_value = 0);

	WHEN 6 THEN -- fixed64 (convert string to number)
		SET uint64_value = _pb_json_parse_unsigned_int(field_json_value);
		SET converted_value = CAST(uint64_value AS JSON);
		SET is_default = (uint64_value = 0);

	WHEN 16 THEN -- sfixed64 (convert string to number)
		SET int64_value = _pb_json_parse_signed_int(field_json_value);
		SET converted_value = CAST(int64_value AS JSON);
		SET is_default = (int64_value = 0);

	WHEN 18 THEN -- sint64 (convert string to number)
		SET int64_value = _pb_json_parse_signed_int(field_json_value);
		SET converted_value = CAST(int64_value AS JSON);
		SET is_default = (int64_value = 0);

	WHEN 5 THEN -- int32 (handle string numbers including exponential notation)
		SET int32_value = _pb_json_parse_signed_int(field_json_value);
		SET converted_value = CAST(int32_value AS JSON);
		SET is_default = (int32_value = 0);

	WHEN 13 THEN -- uint32 (handle string numbers including exponential notation)
		SET uint32_value = _pb_json_parse_unsigned_int(field_json_value);
		SET converted_value = CAST(uint32_value AS JSON);
		SET is_default = (uint32_value = 0);

	WHEN 7 THEN -- fixed32 (handle with range validation)
		SET uint32_value = _pb_json_parse_unsigned_int(field_json_value);
		SET converted_value = CAST(uint32_value AS JSON);
		SET is_default = (uint32_value = 0);

	WHEN 15 THEN -- sfixed32 (handle with range validation)
		SET int32_value = _pb_json_parse_signed_int(field_json_value);
		SET converted_value = CAST(int32_value AS JSON);
		SET is_default = (int32_value = 0);

	WHEN 17 THEN -- sint32 (handle with range validation)
		SET int32_value = _pb_json_parse_signed_int(field_json_value);
		SET converted_value = CAST(int32_value AS JSON);
		SET is_default = (int32_value = 0);

	WHEN 1 THEN -- double (handle with validation)
		SET double_json_value = _pb_json_parse_double(field_json_value);
		SET converted_value = double_json_value;
		SET is_default = (double_json_value = CAST(0.0 AS JSON));

	WHEN 2 THEN -- float (handle with validation)
		SET float_json_value = _pb_json_parse_float(field_json_value);
		SET converted_value = float_json_value;
		SET is_default = (float_json_value = CAST(0.0 AS JSON));

	WHEN 12 THEN -- bytes
		-- Decode from JSON Base64/Base64URL and re-encode as standard Base64
		SET str_value = JSON_UNQUOTE(field_json_value);
		SET converted_value = JSON_QUOTE(TO_BASE64(_pb_json_parse_bytes(field_json_value)));
		SET is_default = (str_value = '');

	WHEN 8 THEN -- bool
		SET bool_value = CAST(_pb_json_parse_bool(field_json_value) AS JSON);
		SET converted_value = IF(bool_value, CAST('true' AS JSON), CAST('false' AS JSON));
		SET is_default = NOT bool_value;

	WHEN 9 THEN -- string
		SET converted_value = field_json_value;
		SET is_default = (JSON_UNQUOTE(field_json_value) = '');

	ELSE
		-- Unknown field type
		SET message_text = CONCAT('_pb_convert_singular_field_to_number_json: unknown field type: ', field_type);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$

-- Main conversion procedure from ProtoJSON to ProtoNumberJSON
DROP PROCEDURE IF EXISTS _pb_json_to_number_json_proc $$
CREATE PROCEDURE _pb_json_to_number_json_proc(
	IN descriptor_set_json JSON,
	IN full_type_name TEXT,
	IN proto_json JSON,
	OUT result JSON
)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';
	DECLARE message_text TEXT;
	DECLARE message_descriptor JSON;
	DECLARE file_descriptor JSON;
	DECLARE syntax TEXT;
	DECLARE fields JSON;
	DECLARE field_count INT;
	DECLARE field_index INT;
	DECLARE field_descriptor JSON;
	-- Field properties
	DECLARE field_number INT;
	DECLARE field_name TEXT;
	DECLARE field_label INT;
	DECLARE field_type INT;
	DECLARE field_type_name TEXT;
	DECLARE json_name TEXT;
	DECLARE proto3_optional BOOLEAN;
	DECLARE oneof_index INT;
	-- Processing variables
	DECLARE is_repeated BOOLEAN;
	DECLARE has_presence BOOLEAN;
	DECLARE is_default BOOLEAN;
	DECLARE field_json_value JSON;
	DECLARE source_field_name TEXT;
	DECLARE converted_value JSON;
	-- Array processing
	DECLARE array_value JSON;
	DECLARE array_length INT;
	DECLARE array_index INT;
	DECLARE array_element JSON;
	DECLARE converted_array JSON;
	-- Nested message processing
	DECLARE nested_json JSON;
	-- Map processing
	DECLARE is_map BOOLEAN DEFAULT FALSE;
	DECLARE map_entry_descriptor JSON;
	DECLARE map_value_field JSON;
	DECLARE map_value_type INT;
	DECLARE map_value_type_name TEXT;
	DECLARE map_keys JSON;
	DECLARE map_key_count INT;
	DECLARE map_key_index INT;
	DECLARE map_key_name TEXT;
	DECLARE map_value_json JSON;
	DECLARE converted_map JSON;

	-- Set recursion limit for nested message processing
	SET @@SESSION.max_sp_recursion_depth = 255;

	-- Initialize result as empty object
	SET result = JSON_OBJECT();

	-- Get message descriptor
	SET message_descriptor = _pb_get_message_descriptor(descriptor_set_json, full_type_name);

	IF message_descriptor IS NULL THEN
		SET message_text = CONCAT('_pb_json_to_number_json_proc: message type not found: ', full_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get fields array (field 2 in DescriptorProto)
	SET fields = JSON_EXTRACT(message_descriptor, '$."2"');
	SET field_count = JSON_LENGTH(fields);
	SET field_index = 0;

	-- Get file descriptor to determine syntax
	SET file_descriptor = _pb_get_file_descriptor(descriptor_set_json, full_type_name);
	SET syntax = COALESCE(JSON_UNQUOTE(JSON_EXTRACT(file_descriptor, '$."12"')), 'proto2');

	-- Process each field in the message descriptor
	field_loop: WHILE field_index < field_count DO
		SET field_descriptor = JSON_EXTRACT(fields, CONCAT('$[', field_index, ']'));

		-- Extract field metadata using protobuf field numbers
		SET field_number = JSON_EXTRACT(field_descriptor, '$."3"'); -- number
		SET field_name = JSON_UNQUOTE(JSON_EXTRACT(field_descriptor, '$."1"')); -- name
		SET field_label = COALESCE(JSON_EXTRACT(field_descriptor, '$."4"'), 1); -- label
		SET field_type = JSON_EXTRACT(field_descriptor, '$."5"'); -- type
		SET field_type_name = JSON_UNQUOTE(JSON_EXTRACT(field_descriptor, '$."6"')); -- type_name
		SET json_name = JSON_UNQUOTE(JSON_EXTRACT(field_descriptor, '$."10"')); -- json_name
		SET proto3_optional = CAST(JSON_EXTRACT(field_descriptor, '$."17"') AS UNSIGNED) = 1; -- proto3_optional
		SET oneof_index = JSON_EXTRACT(field_descriptor, '$."9"'); -- oneof_index

		SET is_repeated = (field_label = 3);
		SET has_presence = (syntax != 'proto3' OR proto3_optional OR oneof_index IS NOT NULL);

		-- Check if this is a map field
		SET is_map = FALSE;
		IF is_repeated AND field_type = 11 AND field_type_name IS NOT NULL THEN -- TYPE_MESSAGE
			SET map_entry_descriptor = _pb_get_message_descriptor(descriptor_set_json, field_type_name);
			SET is_map = COALESCE(CAST(JSON_EXTRACT(map_entry_descriptor, '$."7"."7"') AS UNSIGNED), FALSE); -- map_entry
		END IF;

		-- Try multiple field name variations:
		-- 1. json_name if specified in proto
		-- 2. camelCase version of field name
		-- 3. original proto field name
		SET field_json_value = NULL;

		-- First try json_name if specified
		IF json_name IS NOT NULL THEN
			IF JSON_CONTAINS_PATH(proto_json, 'one', CONCAT('$.', json_name)) THEN
				SET field_json_value = JSON_EXTRACT(proto_json, CONCAT('$.', json_name));
			END IF;
		END IF;

		-- If not found and json_name is different from camelCase version, try camelCase
		IF field_json_value IS NULL THEN
			SET source_field_name = _pb_util_snake_to_camel(field_name);
			IF json_name IS NULL OR json_name != source_field_name THEN
				IF JSON_CONTAINS_PATH(proto_json, 'one', CONCAT('$.', source_field_name)) THEN
					SET field_json_value = JSON_EXTRACT(proto_json, CONCAT('$.', source_field_name));
				END IF;
			END IF;
		END IF;

		-- If still not found, try original proto field name
		IF field_json_value IS NULL THEN
			IF JSON_CONTAINS_PATH(proto_json, 'one', CONCAT('$.', field_name)) THEN
				SET field_json_value = JSON_EXTRACT(proto_json, CONCAT('$.', field_name));
			END IF;
		END IF;

		-- Process field if it exists in JSON
		IF field_json_value IS NOT NULL THEN
			-- Handle null values: null means "field is absent" for both singular and repeated fields
			IF JSON_TYPE(field_json_value) = 'NULL' THEN
				-- For both singular and repeated fields, null means the field is absent - skip processing
				SET field_index = field_index + 1;
				ITERATE field_loop;
			END IF;

			IF is_map THEN
				-- Handle map fields - convert JSON object to ProtoNumberJSON format
				-- Maps in ProtoJSON are objects like {"key1": "value1", "key2": "value2"}
				-- In ProtoNumberJSON, values may need conversion based on their types
				-- Get map value field type for conversion
				SET map_value_field = JSON_EXTRACT(map_entry_descriptor, '$."2"[1]'); -- second field (value)
				SET map_value_type = JSON_EXTRACT(map_value_field, '$."5"');
				SET map_value_type_name = JSON_UNQUOTE(JSON_EXTRACT(map_value_field, '$."6"'));

				-- Convert map values if necessary
				SET map_keys = JSON_KEYS(field_json_value);
				SET map_key_count = JSON_LENGTH(map_keys);
				SET converted_map = JSON_OBJECT();
				SET map_key_index = 0;

				-- TODO: validate map_key
				WHILE map_key_index < map_key_count DO
					SET map_key_name = JSON_UNQUOTE(JSON_EXTRACT(map_keys, CONCAT('$[', map_key_index, ']')));
					SET map_value_json = JSON_EXTRACT(field_json_value, CONCAT('$."', map_key_name, '"'));

					-- Use singular field conversion procedure
					CALL _pb_convert_singular_field_to_number_json(descriptor_set_json, map_value_type, map_value_type_name, map_value_json, converted_value, is_default);

					-- Add converted value to map
					SET converted_map = JSON_SET(converted_map, CONCAT('$."', map_key_name, '"'), converted_value);
					SET map_key_index = map_key_index + 1;
				END WHILE;

				-- In proto3, skip empty maps unless proto3_optional is true or it's a oneof field
				IF map_key_count > 0 THEN
					SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_map);
				END IF;
			ELSEIF is_repeated THEN
				-- Handle repeated fields (arrays)
				SET array_value = field_json_value;
				SET array_length = JSON_LENGTH(array_value);
				SET converted_array = JSON_ARRAY();
				SET array_index = 0;

				WHILE array_index < array_length DO
					SET array_element = JSON_EXTRACT(array_value, CONCAT('$[', array_index, ']'));

					-- Convert element using singular field conversion procedure
					CALL _pb_convert_singular_field_to_number_json(descriptor_set_json, field_type, field_type_name, array_element, converted_value, is_default);
					SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', converted_value);

					SET array_index = array_index + 1;
				END WHILE;

				IF array_length > 0 THEN
					SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_array);
				END IF;
			ELSE
				-- Handle singular fields
				CALL _pb_convert_singular_field_to_number_json(descriptor_set_json, field_type, field_type_name, field_json_value, converted_value, is_default);
				-- Include field unless it's a default value in proto3 without explicit presence
				IF has_presence OR NOT is_default THEN
					SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_value);
				END IF;
			END IF;
		END IF;

		SET field_index = field_index + 1;
	END WHILE field_loop;
END $$

-- Public function interface
DROP FUNCTION IF EXISTS _pb_json_to_number_json $$
CREATE FUNCTION _pb_json_to_number_json(descriptor_set_json JSON, type_name TEXT, proto_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;
	DECLARE result JSON;

	-- Validate type name starts with dot
	IF type_name NOT LIKE '.%' THEN
		SET message_text = CONCAT('_pb_json_to_number_json: type name `', type_name, '` must start with a dot');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF proto_json IS NULL THEN
		RETURN NULL;
	END IF;

	CALL _pb_json_to_number_json_proc(descriptor_set_json, type_name, proto_json, result);
	RETURN result;
END $$
