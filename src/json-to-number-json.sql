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

-- Helper procedure to convert enum string name to numeric value
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

	-- Get enum descriptor
	SET enum_descriptor = _pb_get_enum_descriptor(descriptor_set_json, full_enum_type_name);

	IF enum_descriptor IS NULL THEN
		SET message_text = CONCAT('_pb_convert_json_enum_to_number: enum type not found: ', full_enum_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get the values array (field 2 in EnumDescriptor)
	SET values_array = JSON_EXTRACT(enum_descriptor, '$."2"');
	SET value_count = JSON_LENGTH(values_array);
	SET value_index = 0;

	-- Search for the string value
	search_loop: WHILE value_index < value_count DO
		SET value_descriptor = JSON_EXTRACT(values_array, CONCAT('$[', value_index, ']'));
		SET value_name = JSON_UNQUOTE(JSON_EXTRACT(value_descriptor, '$."1"')); -- name field

		IF value_name = enum_string_value THEN
			SET value_number = JSON_EXTRACT(value_descriptor, '$."2"'); -- number field
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

-- Helper procedure to convert well-known type from ProtoJSON to ProtoNumberJSON
DROP PROCEDURE IF EXISTS _pb_convert_json_wkt_to_number_json $$
CREATE PROCEDURE _pb_convert_json_wkt_to_number_json(
	IN full_type_name TEXT,
	IN proto_json_value JSON,
	OUT number_json_value JSON
)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';
	DECLARE message_text TEXT;
	-- Variables for Any handling
	DECLARE type_url TEXT;
	DECLARE remaining_object JSON;
	-- Variables for wrapper types
	DECLARE int64_val BIGINT;
	DECLARE uint64_val BIGINT UNSIGNED;
	DECLARE int32_val INT;
	DECLARE uint32_val INT UNSIGNED;

	CASE full_type_name
	WHEN '.google.protobuf.Timestamp' THEN
		SET number_json_value = _pb_wkt_timestamp_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Duration' THEN
		SET number_json_value = _pb_wkt_duration_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.FieldMask' THEN
		SET number_json_value = _pb_wkt_field_mask_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Value' THEN
		SET number_json_value = _pb_wkt_value_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.Struct' THEN
		SET number_json_value = _pb_wkt_struct_json_to_number_json(proto_json_value);
	WHEN '.google.protobuf.ListValue' THEN
		SET number_json_value = _pb_wkt_list_value_json_to_number_json(proto_json_value);

	WHEN '.google.protobuf.StringValue' THEN
		-- Unwrapped string becomes {"1": "value"} - but omit if empty string
		IF JSON_UNQUOTE(proto_json_value) != '' THEN
			SET number_json_value = JSON_OBJECT('1', proto_json_value);
		ELSE
			SET number_json_value = JSON_OBJECT();
		END IF;

	WHEN '.google.protobuf.Int64Value' THEN
		-- Unwrapped number becomes {"1": value} - but omit if zero
		IF proto_json_value IS NULL THEN
			SET int64_val = 0;
		ELSE
			SET int64_val = CAST(JSON_UNQUOTE(proto_json_value) AS SIGNED);
		END IF;
		IF int64_val != 0 THEN
			SET number_json_value = JSON_OBJECT('1', int64_val);
		ELSE
			SET number_json_value = JSON_OBJECT();
		END IF;

	WHEN '.google.protobuf.UInt64Value' THEN
		IF proto_json_value IS NULL THEN
			SET uint64_val = 0;
		ELSE
			SET uint64_val = CAST(JSON_UNQUOTE(proto_json_value) AS UNSIGNED);
		END IF;
		IF uint64_val != 0 THEN
			SET number_json_value = JSON_OBJECT('1', uint64_val);
		ELSE
			SET number_json_value = JSON_OBJECT();
		END IF;

	WHEN '.google.protobuf.Int32Value' THEN
		IF proto_json_value IS NULL THEN
			SET int32_val = 0;
		ELSE
			SET int32_val = CAST(JSON_UNQUOTE(proto_json_value) AS SIGNED);
		END IF;
		IF int32_val != 0 THEN
			SET number_json_value = JSON_OBJECT('1', int32_val);
		ELSE
			SET number_json_value = JSON_OBJECT();
		END IF;

	WHEN '.google.protobuf.UInt32Value' THEN
		IF proto_json_value IS NULL THEN
			SET uint32_val = 0;
		ELSE
			SET uint32_val = CAST(JSON_UNQUOTE(proto_json_value) AS UNSIGNED);
		END IF;
		IF uint32_val != 0 THEN
			SET number_json_value = JSON_OBJECT('1', uint32_val);
		ELSE
			SET number_json_value = JSON_OBJECT();
		END IF;

	WHEN '.google.protobuf.BoolValue' THEN
		IF proto_json_value != false THEN
			SET number_json_value = JSON_OBJECT('1', proto_json_value);
		ELSE
			SET number_json_value = JSON_OBJECT();
		END IF;

	WHEN '.google.protobuf.FloatValue' THEN
		IF proto_json_value != 0.0 THEN
			SET number_json_value = JSON_OBJECT('1', proto_json_value);
		ELSE
			SET number_json_value = JSON_OBJECT();
		END IF;

	WHEN '.google.protobuf.DoubleValue' THEN
		IF proto_json_value != 0.0 THEN
			SET number_json_value = JSON_OBJECT('1', proto_json_value);
		ELSE
			SET number_json_value = JSON_OBJECT();
		END IF;

	WHEN '.google.protobuf.BytesValue' THEN
		IF JSON_UNQUOTE(proto_json_value) != '' THEN
			SET number_json_value = JSON_OBJECT('1', proto_json_value);
		ELSE
			SET number_json_value = JSON_OBJECT();
		END IF;

	WHEN '.google.protobuf.Empty' THEN
		-- Empty object stays empty
		SET number_json_value = JSON_OBJECT();

	WHEN '.google.protobuf.Any' THEN
		-- {"@type": "url", "field": "value"} -> {"1": "url", "2": "base64data"}
		-- This is simplified - real Any handling is more complex
		SET type_url = JSON_UNQUOTE(JSON_EXTRACT(proto_json_value, '$."@type"'));
		SET remaining_object = JSON_REMOVE(proto_json_value, '$."@type"');
		-- Convert remaining object to base64-encoded bytes (simplified)
		SET number_json_value = JSON_OBJECT('1', type_url, '2', TO_BASE64(remaining_object));

	ELSE
		-- Not a well-known type, return as-is
		SET number_json_value = proto_json_value;
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
	-- Processing variables
	DECLARE is_repeated BOOLEAN;
	DECLARE field_json_value JSON;
	DECLARE source_field_name TEXT;
	DECLARE converted_value JSON;
	DECLARE int64_value BIGINT;
	DECLARE uint64_value BIGINT UNSIGNED;
	DECLARE int32_value INT;
	DECLARE uint32_value INT UNSIGNED;
	DECLARE double_json_value JSON;
	DECLARE float_json_value JSON;
	DECLARE str_value TEXT;
	DECLARE enum_string_value TEXT;
	DECLARE enum_numeric_value INT;
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

		-- Determine source field name (json_name takes precedence over field_name)
		SET source_field_name = COALESCE(json_name, field_name);
		SET is_repeated = (field_label = 3);

		-- Check if this is a map field
		SET is_map = FALSE;
		IF is_repeated AND field_type = 11 AND field_type_name IS NOT NULL THEN -- TYPE_MESSAGE
			SET map_entry_descriptor = _pb_get_message_descriptor(descriptor_set_json, field_type_name);
			SET is_map = COALESCE(CAST(JSON_EXTRACT(map_entry_descriptor, '$."7"."7"') AS UNSIGNED), FALSE); -- map_entry
		END IF;

		-- Check if field exists in source JSON
		IF JSON_CONTAINS_PATH(proto_json, 'one', CONCAT('$.', source_field_name)) THEN
			SET field_json_value = JSON_EXTRACT(proto_json, CONCAT('$.', source_field_name));

			-- Handle null values: null is allowed for singular fields (treated as absent),
			-- but not for repeated fields
			IF JSON_TYPE(field_json_value) = 'NULL' THEN
				IF is_repeated THEN
					SET message_text = CONCAT('Invalid null value for repeated field `', field_name, '` (', field_number, ')');
					SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
				ELSE
					-- For singular fields, null means the field is absent - skip processing
					SET field_index = field_index + 1;
					ITERATE field_loop;
				END IF;
			END IF;

			IF is_repeated THEN
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

					WHILE map_key_index < map_key_count DO
						SET map_key_name = JSON_UNQUOTE(JSON_EXTRACT(map_keys, CONCAT('$[', map_key_index, ']')));
						SET map_value_json = JSON_EXTRACT(field_json_value, CONCAT('$."', map_key_name, '"'));

						-- Convert value based on type, handling null values with appropriate defaults
						CASE map_value_type
						WHEN 14 THEN -- enum
							IF map_value_json IS NULL THEN
								-- Default enum value is 0
								SET converted_value = CAST(0 AS JSON);
							ELSE
								CALL _pb_convert_json_enum_to_number(descriptor_set_json, map_value_type_name, JSON_UNQUOTE(map_value_json), enum_numeric_value);
								SET converted_value = CAST(enum_numeric_value AS JSON);
							END IF;
						WHEN 11 THEN -- message
							IF map_value_json IS NULL THEN
								-- Default message value is empty object
								SET converted_value = JSON_OBJECT();
							ELSE
								-- Check if it's a well-known type
								CALL _pb_is_well_known_type(map_value_type_name, @is_wkt);
								IF @is_wkt THEN
									CALL _pb_convert_json_wkt_to_number_json(map_value_type_name, map_value_json, converted_value);
								ELSE
									-- Recursively convert nested message
									CALL _pb_json_to_number_json_proc(descriptor_set_json, map_value_type_name, map_value_json, converted_value);
								END IF;
							END IF;
						WHEN 3 THEN -- int64 (convert string to number)
							IF map_value_json IS NULL THEN
								SET converted_value = CAST(0 AS JSON);
							ELSE
								SET converted_value = CAST(JSON_UNQUOTE(map_value_json) AS JSON);
							END IF;
						WHEN 4 THEN -- uint64 (convert string to number)
							IF map_value_json IS NULL THEN
								SET converted_value = CAST(0 AS JSON);
							ELSE
								SET converted_value = CAST(JSON_UNQUOTE(map_value_json) AS JSON);
							END IF;
						WHEN 6 THEN -- fixed64 (convert string to number)
							IF map_value_json IS NULL THEN
								SET converted_value = CAST(0 AS JSON);
							ELSE
								SET converted_value = CAST(JSON_UNQUOTE(map_value_json) AS JSON);
							END IF;
						WHEN 16 THEN -- sfixed64 (convert string to number)
							IF map_value_json IS NULL THEN
								SET converted_value = CAST(0 AS JSON);
							ELSE
								SET converted_value = CAST(JSON_UNQUOTE(map_value_json) AS JSON);
							END IF;
						WHEN 18 THEN -- sint64 (convert string to number)
							IF map_value_json IS NULL THEN
								SET converted_value = CAST(0 AS JSON);
							ELSE
								SET converted_value = CAST(JSON_UNQUOTE(map_value_json) AS JSON);
							END IF;
						ELSE
							-- Other primitive types: handle null values with appropriate defaults
							IF map_value_json IS NULL THEN
								CASE map_value_type
								WHEN 1 THEN -- double
									SET converted_value = CAST(0.0 AS JSON);
								WHEN 2 THEN -- float
									SET converted_value = CAST(0.0 AS JSON);
								WHEN 5 THEN -- int32
									SET converted_value = CAST(0 AS JSON);
								WHEN 7 THEN -- fixed32
									SET converted_value = CAST(0 AS JSON);
								WHEN 8 THEN -- bool
									SET converted_value = CAST(FALSE AS JSON);
								WHEN 9 THEN -- string
									SET converted_value = JSON_QUOTE('');
								WHEN 12 THEN -- bytes
									SET converted_value = JSON_QUOTE('');
								WHEN 13 THEN -- uint32
									SET converted_value = CAST(0 AS JSON);
								WHEN 15 THEN -- sfixed32
									SET converted_value = CAST(0 AS JSON);
								WHEN 17 THEN -- sint32
									SET converted_value = CAST(0 AS JSON);
								ELSE
									-- Unknown type, use null
									SET converted_value = NULL;
								END CASE;
							ELSE
								SET converted_value = map_value_json;
							END IF;
						END CASE;

						-- Add converted value to map
						SET converted_map = JSON_SET(converted_map, CONCAT('$."', map_key_name, '"'), converted_value);
						SET map_key_index = map_key_index + 1;
					END WHILE;

					-- In proto3, skip empty maps unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR JSON_LENGTH(converted_map) > 0 THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_map);
					END IF;
				ELSE
					-- Handle repeated fields (arrays)
					SET array_value = field_json_value;
					SET array_length = JSON_LENGTH(array_value);
					SET converted_array = JSON_ARRAY();
					SET array_index = 0;

				array_loop: WHILE array_index < array_length DO
					SET array_element = JSON_EXTRACT(array_value, CONCAT('$[', array_index, ']'));

					-- Convert element based on field type
					CASE field_type
					WHEN 14 THEN -- enum
						SET enum_string_value = JSON_UNQUOTE(array_element);
						CALL _pb_convert_json_enum_to_number(descriptor_set_json, field_type_name, enum_string_value, enum_numeric_value);
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', enum_numeric_value);
					WHEN 11 THEN -- message
						-- Check if it's a well-known type
						CALL _pb_is_well_known_type(field_type_name, @is_wkt);
						IF @is_wkt THEN
							CALL _pb_convert_json_wkt_to_number_json(field_type_name, array_element, converted_value);
							SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', converted_value);
						ELSE
							-- Recursively convert nested message
							CALL _pb_json_to_number_json_proc(descriptor_set_json, field_type_name, array_element, nested_json);
							SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', nested_json);
						END IF;
					WHEN 3 THEN -- int64 (convert string to number)
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', CAST(JSON_UNQUOTE(array_element) AS DECIMAL(20,0)));
					WHEN 4 THEN -- uint64 (convert string to number)
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', CAST(JSON_UNQUOTE(array_element) AS DECIMAL(20,0)));
					WHEN 6 THEN -- fixed64 (convert string to number)
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', CAST(JSON_UNQUOTE(array_element) AS DECIMAL(20,0)));
					WHEN 16 THEN -- sfixed64 (convert string to number)
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', CAST(JSON_UNQUOTE(array_element) AS DECIMAL(20,0)));
					WHEN 18 THEN -- sint64 (convert string to number)
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', CAST(JSON_UNQUOTE(array_element) AS DECIMAL(20,0)));
					WHEN 1 THEN -- double (validate and convert)
						SET double_json_value = _pb_json_parse_double(array_element);
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', double_json_value);
					WHEN 2 THEN -- float (validate and convert)
						SET float_json_value = _pb_json_parse_float(array_element);
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', float_json_value);
					ELSE
						-- Other primitive types stay the same
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', array_element);
					END CASE;

					SET array_index = array_index + 1;
				END WHILE array_loop;

					-- In proto3, skip empty arrays unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR array_length > 0 THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_array);
					END IF;
				END IF;
			ELSE
				-- Handle singular fields
				CASE field_type
				WHEN 14 THEN -- enum
					SET enum_string_value = JSON_UNQUOTE(field_json_value);
					CALL _pb_convert_json_enum_to_number(descriptor_set_json, field_type_name, enum_string_value, enum_numeric_value);
					-- In proto3, skip enum fields with zero value unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR enum_numeric_value != 0 THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), enum_numeric_value);
					END IF;
				WHEN 11 THEN -- message
					-- Check if it's a well-known type
					CALL _pb_is_well_known_type(field_type_name, @is_wkt);
					IF @is_wkt THEN
						CALL _pb_convert_json_wkt_to_number_json(field_type_name, field_json_value, converted_value);
						-- For WKTs, always include the field but use empty object for zero values
						IF _pb_is_wkt_zero_value(field_type_name, converted_value) THEN
							SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), JSON_OBJECT());
						ELSE
							SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_value);
						END IF;
					ELSE
						-- Recursively convert nested message
						CALL _pb_json_to_number_json_proc(descriptor_set_json, field_type_name, field_json_value, nested_json);
						-- Always include nested messages in proto3 (they represent explicit field presence)
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), nested_json);
					END IF;
				WHEN 3 THEN -- int64 (convert string to number)
					SET int64_value = _pb_json_parse_signed_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (int64_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), int64_value);
					END IF;
				WHEN 4 THEN -- uint64 (convert string to number)
					SET uint64_value = _pb_json_parse_unsigned_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (uint64_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), uint64_value);
					END IF;
				WHEN 6 THEN -- fixed64 (convert string to number)
					SET uint64_value = _pb_json_parse_unsigned_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (uint64_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), uint64_value);
					END IF;
				WHEN 16 THEN -- sfixed64 (convert string to number)
					SET int64_value = _pb_json_parse_signed_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (int64_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), int64_value);
					END IF;
				WHEN 18 THEN -- sint64 (convert string to number)
					SET int64_value = _pb_json_parse_signed_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (int64_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), int64_value);
					END IF;
				WHEN 5 THEN -- int32 (handle string numbers including exponential notation)
					SET int32_value = _pb_json_parse_signed_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (int32_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), int32_value);
					END IF;
				WHEN 13 THEN -- uint32 (handle string numbers including exponential notation)
					SET uint32_value = _pb_json_parse_unsigned_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (uint32_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), uint32_value);
					END IF;
				WHEN 7 THEN -- fixed32 (handle with range validation)
					SET uint32_value = _pb_json_parse_unsigned_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (uint32_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), uint32_value);
					END IF;
				WHEN 15 THEN -- sfixed32 (handle with range validation)
					SET int32_value = _pb_json_parse_signed_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (int32_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), int32_value);
					END IF;
				WHEN 17 THEN -- sint32 (handle with range validation)
					SET int32_value = _pb_json_parse_signed_int(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (int32_value = 0) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), int32_value);
					END IF;
				WHEN 1 THEN -- double (handle with validation)
					SET double_json_value = _pb_json_parse_double(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (double_json_value = CAST(0.0 AS JSON)) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), double_json_value);
					END IF;
				WHEN 2 THEN -- float (handle with validation)
					SET float_json_value = _pb_json_parse_float(field_json_value);
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (float_json_value = CAST(0.0 AS JSON)) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), float_json_value);
					END IF;
				ELSE
					-- Other primitive types: bool, string, bytes
					-- In proto3, skip zero/default values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (
						(field_type = 8 AND field_json_value = false) OR                -- bool = false
						(field_type = 9 AND JSON_UNQUOTE(field_json_value) = '') OR     -- string = ""
						(field_type = 12 AND JSON_UNQUOTE(field_json_value) = '')       -- bytes = ""
					) THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), field_json_value);
					END IF;
				END CASE;
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
