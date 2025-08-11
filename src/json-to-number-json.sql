DELIMITER $$

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
	DECLARE timestamp_str TEXT;
	DECLARE duration_str TEXT;
	DECLARE seconds_part BIGINT;
	DECLARE nanos_part INT;
	DECLARE dot_pos INT;
	DECLARE seconds_str TEXT;
	DECLARE nanos_str TEXT;
	DECLARE suffix_char CHAR(1);
	-- Variables for Any handling
	DECLARE type_url TEXT;
	DECLARE remaining_object JSON;

	CASE full_type_name
	WHEN '.google.protobuf.Timestamp' THEN
		-- Convert ISO 8601 timestamp to {seconds, nanos}
		SET timestamp_str = JSON_UNQUOTE(proto_json_value);
		-- Parse RFC3339 timestamp format: "1972-01-01T10:00:20.021Z"
		-- For simplicity, we'll use MySQL's built-in timestamp parsing
		SET seconds_part = UNIX_TIMESTAMP(STR_TO_DATE(LEFT(timestamp_str, 19), '%Y-%m-%dT%H:%i:%s'));

		-- Extract nanoseconds part if present
		SET dot_pos = LOCATE('.', timestamp_str);
		IF dot_pos > 0 THEN
			SET nanos_str = SUBSTRING(timestamp_str, dot_pos + 1);
			-- Remove trailing 'Z'
			SET nanos_str = LEFT(nanos_str, LENGTH(nanos_str) - 1);
			-- Pad to 9 digits (nanoseconds)
			SET nanos_str = RPAD(nanos_str, 9, '0');
			SET nanos_part = CAST(nanos_str AS UNSIGNED);
		ELSE
			SET nanos_part = 0;
		END IF;

		SET number_json_value = JSON_OBJECT('1', seconds_part, '2', nanos_part);

	WHEN '.google.protobuf.Duration' THEN
		-- Convert duration string like "3.5s" to {seconds, nanos}
		SET duration_str = JSON_UNQUOTE(proto_json_value);
		SET suffix_char = RIGHT(duration_str, 1);

		IF suffix_char != 's' THEN
			SET message_text = CONCAT('_pb_convert_json_wkt_to_number_json: invalid Duration format: ', duration_str);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		-- Remove 's' suffix
		SET duration_str = LEFT(duration_str, LENGTH(duration_str) - 1);

		-- Split on decimal point
		SET dot_pos = LOCATE('.', duration_str);
		IF dot_pos > 0 THEN
			SET seconds_str = LEFT(duration_str, dot_pos - 1);
			SET nanos_str = SUBSTRING(duration_str, dot_pos + 1);
			-- Pad to 9 digits
			SET nanos_str = RPAD(nanos_str, 9, '0');
			SET seconds_part = CAST(seconds_str AS SIGNED);
			SET nanos_part = CAST(nanos_str AS UNSIGNED);
		ELSE
			SET seconds_part = CAST(duration_str AS SIGNED);
			SET nanos_part = 0;
		END IF;

		SET number_json_value = JSON_OBJECT('1', seconds_part, '2', nanos_part);

	WHEN '.google.protobuf.StringValue' THEN
		-- Unwrapped string becomes {"1": "value"}
		SET number_json_value = JSON_OBJECT('1', proto_json_value);

	WHEN '.google.protobuf.Int64Value' THEN
		-- Unwrapped number becomes {"1": value}
		SET number_json_value = JSON_OBJECT('1', CAST(JSON_UNQUOTE(proto_json_value) AS SIGNED));

	WHEN '.google.protobuf.UInt64Value' THEN
		SET number_json_value = JSON_OBJECT('1', CAST(JSON_UNQUOTE(proto_json_value) AS UNSIGNED));

	WHEN '.google.protobuf.Int32Value' THEN
		SET number_json_value = JSON_OBJECT('1', CAST(JSON_UNQUOTE(proto_json_value) AS SIGNED));

	WHEN '.google.protobuf.UInt32Value' THEN
		SET number_json_value = JSON_OBJECT('1', CAST(JSON_UNQUOTE(proto_json_value) AS UNSIGNED));

	WHEN '.google.protobuf.BoolValue' THEN
		SET number_json_value = JSON_OBJECT('1', proto_json_value);

	WHEN '.google.protobuf.FloatValue' THEN
		SET number_json_value = JSON_OBJECT('1', proto_json_value);

	WHEN '.google.protobuf.DoubleValue' THEN
		SET number_json_value = JSON_OBJECT('1', proto_json_value);

	WHEN '.google.protobuf.BytesValue' THEN
		SET number_json_value = JSON_OBJECT('1', proto_json_value);

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
	DECLARE numeric_value DECIMAL(20,0);
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

		-- Check if field exists in source JSON
		IF JSON_CONTAINS_PATH(proto_json, 'one', CONCAT('$.', source_field_name)) THEN
			SET field_json_value = JSON_EXTRACT(proto_json, CONCAT('$.', source_field_name));

			IF is_repeated THEN
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
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), converted_value);
					ELSE
						-- Recursively convert nested message
						CALL _pb_json_to_number_json_proc(descriptor_set_json, field_type_name, field_json_value, nested_json);
						-- Always include nested messages in proto3 (they represent explicit field presence)
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), nested_json);
					END IF;
				WHEN 3 THEN -- int64 (convert string to number)
					SET numeric_value = CAST(JSON_UNQUOTE(field_json_value) AS DECIMAL(20,0));
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (numeric_value = 0 AND JSON_UNQUOTE(field_json_value) = '0') THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), numeric_value);
					END IF;
				WHEN 4 THEN -- uint64 (convert string to number)
					SET numeric_value = CAST(JSON_UNQUOTE(field_json_value) AS DECIMAL(20,0));
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (numeric_value = 0 AND JSON_UNQUOTE(field_json_value) = '0') THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), numeric_value);
					END IF;
				WHEN 6 THEN -- fixed64 (convert string to number)
					SET numeric_value = CAST(JSON_UNQUOTE(field_json_value) AS DECIMAL(20,0));
					-- In proto3, skip zero values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (numeric_value = 0 AND JSON_UNQUOTE(field_json_value) = '0') THEN
						SET result = JSON_SET(result, CONCAT('$."', field_number, '"'), numeric_value);
					END IF;
				ELSE
					-- Other primitive types: int32, uint32, float, double, bool, string, bytes
					-- In proto3, skip zero/default values unless proto3_optional is true
					IF syntax != 'proto3' OR proto3_optional OR NOT (
						(field_type IN (5, 13, 15, 17, 18) AND field_json_value = 0) OR  -- int32, uint32, sint32, sint64 = 0
						(field_type IN (1, 2) AND field_json_value = 0.0) OR            -- double, float = 0.0
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

DELIMITER ;