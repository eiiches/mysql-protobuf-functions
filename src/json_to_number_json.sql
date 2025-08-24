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

-- Helper function to parse JSON value as Protobuf string type
DROP FUNCTION IF EXISTS _pb_json_parse_string $$
CREATE FUNCTION _pb_json_parse_string(json_value JSON) RETURNS LONGTEXT DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;

	-- String fields must receive JSON string values
	IF JSON_TYPE(json_value) != 'STRING' THEN
		SET message_text = CONCAT('Invalid JSON type for string field: expected STRING, got ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	RETURN JSON_UNQUOTE(json_value);
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

-- Helper procedure to check if a type is a well-known type
-- TODO: Wrong and deprecated. Don't use this.
DROP PROCEDURE IF EXISTS _pb_is_well_known_type $$
CREATE PROCEDURE _pb_is_well_known_type(IN full_type_name TEXT, OUT is_wkt BOOLEAN)
BEGIN
	IF full_type_name LIKE '.google.protobuf.%' THEN
		SET is_wkt = TRUE;
	ELSE
		SET is_wkt = FALSE;
	END IF;
END $$

-- Helper function to convert JSON enum value to numeric value
-- Returns NULL for unknown values when ignore_unknown_enums is TRUE
DROP FUNCTION IF EXISTS _pb_convert_json_enum_to_number $$
CREATE FUNCTION _pb_convert_json_enum_to_number(
	descriptor_set_json JSON,
	full_enum_type_name TEXT,
	enum_value_json JSON,
	ignore_unknown_enums BOOLEAN
) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;
	DECLARE enum_type_index JSON;
	DECLARE type_paths JSON;
	DECLARE enum_name_index JSON;
	DECLARE enum_number_index JSON;
	DECLARE enum_descriptor JSON;
	DECLARE values_array JSON;
	DECLARE found_index INT;
	DECLARE value_descriptor JSON;
	DECLARE input_as_number INT;
	DECLARE enum_string_value TEXT;
	DECLARE is_numeric BOOLEAN DEFAULT FALSE;

	-- Handle number inputs directly
	IF JSON_TYPE(enum_value_json) = 'INTEGER' THEN
		RETURN CAST(enum_value_json AS SIGNED);
	END IF;

	-- Handle non-string inputs - this preserves the original TEXT input behavior
	IF JSON_TYPE(enum_value_json) != 'STRING' THEN
		SET message_text = CONCAT('_pb_convert_json_enum_to_number: invalid JSON type for enum field: ', JSON_TYPE(enum_value_json));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	SET enum_string_value = JSON_UNQUOTE(enum_value_json);

	-- Get enum type index (field 3 from DescriptorSet) - always signal error if missing
	SET enum_type_index = JSON_EXTRACT(descriptor_set_json, '$.\"3\"');
	IF enum_type_index IS NULL THEN
		SET message_text = '_pb_convert_json_enum_to_number: enum type index not found in descriptor set';
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get paths for the enum type - always signal error if missing
	SET type_paths = JSON_EXTRACT(enum_type_index, CONCAT('$.\"', full_enum_type_name, '\"'));
	IF type_paths IS NULL THEN
		SET message_text = CONCAT('_pb_convert_json_enum_to_number: enum type not found: ', full_enum_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Extract enum value indexes from EnumTypeIndex message
	SET enum_name_index = JSON_EXTRACT(type_paths, '$.\"3\"');
	SET enum_number_index = JSON_EXTRACT(type_paths, '$.\"4\"');

	-- Check if the input is a numeric value (protobuf JSON allows both string names and numeric values)
	SET is_numeric = (enum_string_value REGEXP '^-?[0-9]+$');
	IF is_numeric THEN
		SET input_as_number = CAST(enum_string_value AS SIGNED);
		-- Use number index for O(1) lookup
		SET found_index = JSON_EXTRACT(enum_number_index, CONCAT('$.\"', input_as_number, '\"'));
	ELSE
		-- Use name index for O(1) lookup
		SET found_index = JSON_EXTRACT(enum_name_index, CONCAT('$.\"', enum_string_value, '\"'));
	END IF;

	IF found_index IS NOT NULL THEN
		-- Get enum descriptor and values array to extract the number
		SET enum_descriptor = _pb_descriptor_set_get_enum_descriptor(descriptor_set_json, full_enum_type_name);
		SET values_array = JSON_EXTRACT(enum_descriptor, '$."2"');
		SET value_descriptor = JSON_EXTRACT(values_array, CONCAT('$[', found_index, ']'));
		RETURN JSON_EXTRACT(value_descriptor, '$."2"'); -- number field
	ELSE
		-- Not found, handle based on ignore_unknown_enums flag and numeric input
		-- ignore_unknown_enums only affects unknown VALUES, not missing type definitions
		IF is_numeric THEN
			-- For Proto3, unknown numeric enum values should be accepted as-is
			RETURN input_as_number;
		ELSEIF ignore_unknown_enums THEN
			RETURN NULL;  -- Return NULL to indicate unknown value should be ignored
		ELSE
			SET message_text = CONCAT('_pb_convert_json_enum_to_number: enum value not found: ', enum_string_value, ' in enum ', full_enum_type_name);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
	END IF;
END $$

-- Helper procedure to convert singular field value to ProtoNumberJSON
DROP PROCEDURE IF EXISTS _pb_convert_singular_field_to_number_json $$
CREATE PROCEDURE _pb_convert_singular_field_to_number_json(
	IN descriptor_set_json JSON,
	IN field_type INT,
	IN field_type_name TEXT,
	IN field_json_value JSON,
	IN ignore_unknown_fields BOOLEAN,
	IN ignore_unknown_enums BOOLEAN,
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
		SET converted_value = _pb_convert_json_wkt_to_number_json(field_type, field_type_name, field_json_value);
		IF converted_value IS NULL THEN -- Not handled by well-known type parser
			SET enum_numeric_value = _pb_convert_json_enum_to_number(descriptor_set_json, field_type_name, field_json_value, ignore_unknown_enums);
			SET converted_value = CAST(enum_numeric_value AS JSON);
			SET is_default = (enum_numeric_value = 0);
		ELSE
			SET enum_numeric_value = converted_value;
			SET is_default = (enum_numeric_value = 0);
		END IF;

	WHEN 11 THEN -- message
		IF field_json_value IS NULL THEN
			SET is_default = TRUE;
			SET converted_value = NULL;
		ELSE
			SET is_default = FALSE;
			SET converted_value = _pb_convert_json_wkt_to_number_json(field_type, field_type_name, field_json_value);
			IF converted_value IS NULL THEN -- Not handled by well-known type parser
				-- For regular (non-WKT) messages, validate that the JSON value is an object
				IF JSON_TYPE(field_json_value) != 'OBJECT' THEN
					SET message_text = CONCAT('Invalid JSON type for message field: expected OBJECT, got ', JSON_TYPE(field_json_value));
					SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
				END IF;
				CALL _pb_json_to_number_json_proc(descriptor_set_json, field_type_name, field_json_value, ignore_unknown_fields, ignore_unknown_enums, converted_value);
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
		SET str_value = _pb_json_parse_string(field_json_value);
		SET converted_value = JSON_QUOTE(str_value);
		SET is_default = (str_value = '');

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
	IN ignore_unknown_fields BOOLEAN,
	IN ignore_unknown_enums BOOLEAN,
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
	SET message_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, full_type_name);

	IF message_descriptor IS NULL THEN
		SET message_text = CONCAT('_pb_json_to_number_json_proc: message type not found: ', full_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get fields array (field 2 in DescriptorProto)
	SET fields = JSON_EXTRACT(message_descriptor, '$."2"');
	SET field_count = JSON_LENGTH(fields);
	SET field_index = 0;

	-- Get file descriptor to determine syntax
	SET file_descriptor = _pb_descriptor_set_get_file_descriptor(descriptor_set_json, full_type_name);
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
		SET proto3_optional = JSON_EXTRACT(field_descriptor, '$."17"'); -- proto3_optional
		SET oneof_index = JSON_EXTRACT(field_descriptor, '$."9"'); -- oneof_index

		SET is_repeated = (field_label = 3);
		-- Determine field presence
		SET has_presence = (syntax = 'proto2' AND field_label <> 3) -- proto2: all non-repeated fields
			OR (syntax = 'proto3'
				AND (
					(field_label = 1 AND proto3_optional) -- proto3 optional
					OR (field_label <> 3 AND field_type = 11) -- message fields
					OR (oneof_index IS NOT NULL) -- oneof fields
			));

		-- Check if this is a map field
		SET is_map = FALSE;
		IF is_repeated AND field_type = 11 AND field_type_name IS NOT NULL THEN -- TYPE_MESSAGE
			SET map_entry_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, field_type_name);
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

		-- If field is not found, skip to the next field processing.
		IF field_json_value IS NULL THEN
			SET field_index = field_index + 1;
			ITERATE field_loop;
		END IF;

		IF is_map THEN
			-- Explicit JSON null for a map field means the map is empty. This seems weird, but required by AllFieldAcceptNull.
			IF JSON_TYPE(field_json_value) = 'NULL' THEN
				SET field_index = field_index + 1;
				ITERATE field_loop;
			END IF;

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
				CALL _pb_convert_singular_field_to_number_json(descriptor_set_json, map_value_type, map_value_type_name, map_value_json, ignore_unknown_fields, ignore_unknown_enums, converted_value, is_default);

				-- Add converted value to map
				-- converted_value can be NULL if ignore_unknown_enums is set and enum name value is unknown.
				IF converted_value IS NOT NULL THEN
					SET converted_map = JSON_SET(converted_map, CONCAT('$."', map_key_name, '"'), converted_value);
				END IF;

				SET map_key_index = map_key_index + 1;
			END WHILE;

			-- In proto3, skip empty maps unless proto3_optional is true or it's a oneof field
			IF map_key_count > 0 THEN
				SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_map);
			END IF;
		ELSEIF is_repeated THEN
			-- Explicit JSON null for a repeated field means the list is empty. This seems weird, but required by AllFieldAcceptNull.
			IF JSON_TYPE(field_json_value) = 'NULL' THEN
				SET field_index = field_index + 1;
				ITERATE field_loop;
			END IF;

			-- Handle repeated fields (arrays)
			SET array_value = field_json_value;
			SET array_length = JSON_LENGTH(array_value);
			SET converted_array = JSON_ARRAY();
			SET array_index = 0;

			WHILE array_index < array_length DO
				SET array_element = JSON_EXTRACT(array_value, CONCAT('$[', array_index, ']'));

				-- Convert element using singular field conversion procedure
				CALL _pb_convert_singular_field_to_number_json(descriptor_set_json, field_type, field_type_name, array_element, ignore_unknown_fields, ignore_unknown_enums, converted_value, is_default);

				-- converted_value can be NULL if ignore_unknown_enums is set and enum name value is unknown.
				IF converted_value IS NOT NULL THEN
					SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', converted_value);
				END IF;

				SET array_index = array_index + 1;
			END WHILE;

			IF array_length > 0 THEN
				SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_array);
			END IF;
		ELSE
			-- Explicit JSON null for a field that has field presence tracking means the field is not set, except
			-- if the field is .google.protobuf.Value. Explicit JSON null for a Value field, should be recognized
			-- as a non-null Value message with kind.null_value set to NULL_VALUE.
			-- Explicit JSON null for a field that doesn't have a field presence tracking (>=proto3) means a zero
			-- value is set (and zero values are omitted in ProtoNumberJSON or on wire).
			IF JSON_TYPE(field_json_value) = 'NULL' AND (field_type_name IS NULL OR (field_type_name <> '.google.protobuf.Value' AND field_type_name <> '.google.protobuf.NullValue')) THEN
				SET field_index = field_index + 1;
				ITERATE field_loop;
			END IF;

			-- Handle singular fields
			CALL _pb_convert_singular_field_to_number_json(descriptor_set_json, field_type, field_type_name, field_json_value, ignore_unknown_fields, ignore_unknown_enums, converted_value, is_default);

			-- Include field unless it's a default value in proto3 without explicit presence
			-- converted_value can be NULL if ignore_unknown_enums is set and enum name value is unknown.
			IF converted_value IS NOT NULL AND has_presence OR NOT is_default THEN
				SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_value);
			END IF;
		END IF;

		SET field_index = field_index + 1;
	END WHILE field_loop;
END $$

-- Public function interface
DROP FUNCTION IF EXISTS _pb_json_to_number_json $$
CREATE FUNCTION _pb_json_to_number_json(descriptor_set_json JSON, type_name TEXT, proto_json JSON, json_unmarshal_options JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;
	DECLARE result JSON;
	DECLARE ignore_unknown_fields BOOLEAN DEFAULT FALSE;
	DECLARE ignore_unknown_enums BOOLEAN DEFAULT FALSE;

	-- Extract options using the generated accessor functions
	IF json_unmarshal_options IS NOT NULL THEN
		SET ignore_unknown_fields = pb_json_unmarshal_options_get_ignore_unknown_fields(json_unmarshal_options);
		SET ignore_unknown_enums = pb_json_unmarshal_options_get_ignore_unknown_enums(json_unmarshal_options);
	END IF;

	-- Validate type name starts with dot
	IF type_name NOT LIKE '.%' THEN
		SET message_text = CONCAT('_pb_json_to_number_json: type name `', type_name, '` must start with a dot');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF proto_json IS NULL THEN
		RETURN NULL;
	END IF;

	CALL _pb_json_to_number_json_proc(descriptor_set_json, type_name, proto_json, ignore_unknown_fields, ignore_unknown_enums, result);
	RETURN result;
END $$
