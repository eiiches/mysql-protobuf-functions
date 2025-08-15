DELIMITER $$

-- Helper procedure to convert enum numeric value to string name
DROP PROCEDURE IF EXISTS _pb_convert_number_enum_to_json $$
CREATE PROCEDURE _pb_convert_number_enum_to_json(
	IN descriptor_set_json JSON,
	IN full_enum_type_name TEXT,
	IN enum_numeric_value INT,
	OUT enum_string_value TEXT
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
		SET message_text = CONCAT('_pb_convert_number_enum_to_json: enum type not found: ', full_enum_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get the values array (field 2 in EnumDescriptor)
	SET values_array = JSON_EXTRACT(enum_descriptor, '$."2"');
	SET value_count = JSON_LENGTH(values_array);
	SET value_index = 0;

	-- Search for the numeric value
	search_loop: WHILE value_index < value_count DO
		SET value_descriptor = JSON_EXTRACT(values_array, CONCAT('$[', value_index, ']'));
		SET value_number = JSON_EXTRACT(value_descriptor, '$."2"'); -- number field

		IF value_number = enum_numeric_value THEN
			SET value_name = JSON_UNQUOTE(JSON_EXTRACT(value_descriptor, '$."1"')); -- name field
			SET enum_string_value = value_name;
			LEAVE search_loop;
		END IF;

		SET value_index = value_index + 1;
	END WHILE search_loop;

	-- If not found, signal error
	IF value_index >= value_count THEN
		SET message_text = CONCAT('_pb_convert_number_enum_to_json: enum value not found: ', enum_numeric_value, ' in enum ', full_enum_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;
END $$

-- Helper procedure to convert well-known type from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_convert_number_json_to_wkt $$
CREATE PROCEDURE _pb_convert_number_json_to_wkt(
	IN full_type_name TEXT,
	IN number_json_value JSON,
	OUT proto_json_value JSON
)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';
	DECLARE message_text TEXT;
	DECLARE seconds_part BIGINT;
	DECLARE nanos_part INT;
	DECLARE timestamp_str TEXT;
	DECLARE duration_str TEXT;
	DECLARE nanos_str TEXT;
	DECLARE wrapped_value JSON;
	-- Variables for Any handling
	DECLARE type_url TEXT;
	DECLARE any_data TEXT;

	CASE full_type_name
	WHEN '.google.protobuf.Timestamp' THEN
		-- Convert {seconds, nanos} to ISO 8601 timestamp
		SET seconds_part = JSON_EXTRACT(number_json_value, '$.1');
		SET nanos_part = COALESCE(JSON_EXTRACT(number_json_value, '$.2'), 0);

		-- Convert to RFC3339 format
		SET timestamp_str = FROM_UNIXTIME(seconds_part, '%Y-%m-%dT%H:%i:%s');

		-- Add nanoseconds if present
		IF nanos_part > 0 THEN
			SET nanos_str = LPAD(nanos_part, 9, '0');
			-- Remove trailing zeros by converting back to number and to string
			SET nanos_str = CAST(CAST(nanos_str AS UNSIGNED) AS CHAR);
			SET timestamp_str = CONCAT(timestamp_str, '.', nanos_str);
		END IF;

		SET timestamp_str = CONCAT(timestamp_str, 'Z');
		SET proto_json_value = JSON_QUOTE(timestamp_str);

	WHEN '.google.protobuf.Duration' THEN
		-- Convert {seconds, nanos} to duration string like "3.5s"
		SET seconds_part = JSON_EXTRACT(number_json_value, '$.1');
		SET nanos_part = COALESCE(JSON_EXTRACT(number_json_value, '$.2'), 0);

		IF nanos_part = 0 THEN
			SET duration_str = CONCAT(seconds_part, 's');
		ELSE
			SET nanos_str = LPAD(nanos_part, 9, '0');
			-- Remove trailing zeros by converting back to number and to string
			SET nanos_str = CAST(CAST(nanos_str AS UNSIGNED) AS CHAR);
			SET duration_str = CONCAT(seconds_part, '.', nanos_str, 's');
		END IF;

		SET proto_json_value = JSON_QUOTE(duration_str);

	WHEN '.google.protobuf.StringValue' THEN
		-- {"1": "value"} becomes unwrapped "value"
		SET proto_json_value = JSON_EXTRACT(number_json_value, '$.1');

	WHEN '.google.protobuf.Int64Value' THEN
		-- {"1": value} becomes unwrapped "value" (as string for 64-bit)
		SET wrapped_value = JSON_EXTRACT(number_json_value, '$.1');
		SET proto_json_value = JSON_QUOTE(CAST(wrapped_value AS CHAR));

	WHEN '.google.protobuf.UInt64Value' THEN
		SET wrapped_value = JSON_EXTRACT(number_json_value, '$.1');
		SET proto_json_value = JSON_QUOTE(CAST(wrapped_value AS CHAR));

	WHEN '.google.protobuf.Int32Value' THEN
		-- {"1": value} becomes unwrapped value (as number for 32-bit)
		SET proto_json_value = JSON_EXTRACT(number_json_value, '$.1');

	WHEN '.google.protobuf.UInt32Value' THEN
		SET proto_json_value = JSON_EXTRACT(number_json_value, '$.1');

	WHEN '.google.protobuf.BoolValue' THEN
		SET proto_json_value = JSON_EXTRACT(number_json_value, '$.1');

	WHEN '.google.protobuf.FloatValue' THEN
		SET proto_json_value = JSON_EXTRACT(number_json_value, '$.1');

	WHEN '.google.protobuf.DoubleValue' THEN
		SET proto_json_value = JSON_EXTRACT(number_json_value, '$.1');

	WHEN '.google.protobuf.BytesValue' THEN
		SET proto_json_value = JSON_EXTRACT(number_json_value, '$.1');

	WHEN '.google.protobuf.Empty' THEN
		-- Empty object stays empty
		SET proto_json_value = JSON_OBJECT();

	WHEN '.google.protobuf.Any' THEN
		-- {"1": "url", "2": "base64data"} -> {"@type": "url", "field": "value"}
		-- This is simplified - real Any handling is more complex
		SET type_url = JSON_UNQUOTE(JSON_EXTRACT(number_json_value, '$.1'));
		SET any_data = JSON_UNQUOTE(JSON_EXTRACT(number_json_value, '$.2'));
		-- Convert base64 data back to object (simplified)
		SET proto_json_value = JSON_OBJECT('@type', type_url);
		-- In reality, we would decode the base64 data and merge it

	ELSE
		-- Not a well-known type, return as-is
		SET proto_json_value = number_json_value;
	END CASE;
END $$

-- Main conversion procedure from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_number_json_to_json_proc $$
CREATE PROCEDURE _pb_number_json_to_json_proc(
	IN descriptor_set_json JSON,
	IN full_type_name TEXT,
	IN number_json JSON,
	IN emit_default_values BOOLEAN,
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
	DECLARE target_field_name TEXT;
	DECLARE converted_value JSON;
	DECLARE enum_numeric_value INT;
	DECLARE enum_string_value TEXT;
	-- Array processing
	DECLARE array_value JSON;
	DECLARE array_length INT;
	DECLARE array_index INT;
	DECLARE array_element JSON;
	DECLARE converted_array JSON;
	-- Nested message processing
	DECLARE nested_json JSON;
	-- Field presence detection
	DECLARE has_presence BOOLEAN;

	-- Set recursion limit for nested message processing
	SET @@SESSION.max_sp_recursion_depth = 255;

	-- Initialize result as empty object
	SET result = JSON_OBJECT();

	-- Get message descriptor
	SET message_descriptor = _pb_get_message_descriptor(descriptor_set_json, full_type_name);

	IF message_descriptor IS NULL THEN
		SET message_text = CONCAT('_pb_number_json_to_json_proc: message type not found: ', full_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get fields array (field 2 in DescriptorProto)
	SET fields = JSON_EXTRACT(message_descriptor, '$."2"');
	SET field_count = JSON_LENGTH(fields);
	SET field_index = 0;

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
		SET proto3_optional = COALESCE(JSON_EXTRACT(field_descriptor, '$."17"'), 0) = 1; -- proto3_optional

		-- Determine target field name (json_name takes precedence over field_name)
		SET target_field_name = COALESCE(json_name, field_name);
		SET is_repeated = (field_label = 3);

		-- Check if field exists in source JSON (by field number)
		IF JSON_CONTAINS_PATH(number_json, 'one', CONCAT('$."', field_number, '"')) THEN
			SET field_json_value = JSON_EXTRACT(number_json, CONCAT('$."', field_number, '"'));

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
						SET enum_numeric_value = JSON_EXTRACT(array_element, '$');
						CALL _pb_convert_number_enum_to_json(descriptor_set_json, field_type_name, enum_numeric_value, enum_string_value);
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', enum_string_value);
					WHEN 11 THEN -- message
						-- Check if it's a well-known type
						CALL _pb_is_well_known_type(field_type_name, @is_wkt);
						IF @is_wkt THEN
							CALL _pb_convert_number_json_to_wkt(field_type_name, array_element, converted_value);
							SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', converted_value);
						ELSE
							-- Recursively convert nested message
							CALL _pb_number_json_to_json_proc(descriptor_set_json, field_type_name, array_element, emit_default_values, nested_json);
							SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', nested_json);
						END IF;
					WHEN 3 THEN -- int64 (convert number to string)
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', CAST(array_element AS CHAR));
					WHEN 4 THEN -- uint64 (convert number to string)
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', CAST(array_element AS CHAR));
					WHEN 6 THEN -- fixed64 (convert number to string)
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', CAST(array_element AS CHAR));
					ELSE
						-- Other primitive types stay the same
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', array_element);
					END CASE;

					SET array_index = array_index + 1;
				END WHILE array_loop;

				SET result = JSON_SET(result, CONCAT('$.', target_field_name), converted_array);
			ELSE
				-- Handle singular fields
				CASE field_type
				WHEN 14 THEN -- enum
					SET enum_numeric_value = JSON_EXTRACT(field_json_value, '$');
					CALL _pb_convert_number_enum_to_json(descriptor_set_json, field_type_name, enum_numeric_value, enum_string_value);
					SET result = JSON_SET(result, CONCAT('$.', target_field_name), enum_string_value);
				WHEN 11 THEN -- message
					-- Check if it's a well-known type
					CALL _pb_is_well_known_type(field_type_name, @is_wkt);
					IF @is_wkt THEN
						CALL _pb_convert_number_json_to_wkt(field_type_name, field_json_value, converted_value);
						SET result = JSON_SET(result, CONCAT('$.', target_field_name), converted_value);
					ELSE
						-- Recursively convert nested message
						CALL _pb_number_json_to_json_proc(descriptor_set_json, field_type_name, field_json_value, emit_default_values, nested_json);
						SET result = JSON_SET(result, CONCAT('$.', target_field_name), nested_json);
					END IF;
				WHEN 3 THEN -- int64 (convert number to string)
					SET result = JSON_SET(result, CONCAT('$.', target_field_name), CAST(field_json_value AS CHAR));
				WHEN 4 THEN -- uint64 (convert number to string)
					SET result = JSON_SET(result, CONCAT('$.', target_field_name), CAST(field_json_value AS CHAR));
				WHEN 6 THEN -- fixed64 (convert number to string)
					SET result = JSON_SET(result, CONCAT('$.', target_field_name), CAST(field_json_value AS CHAR));
				ELSE
					-- Other primitive types stay the same
					SET result = JSON_SET(result, CONCAT('$.', target_field_name), field_json_value);
				END CASE;
			END IF;
		ELSE
			-- Field is missing from number JSON - emit default value if requested for non-optional fields
			IF emit_default_values THEN
				-- Determine if field has presence-sensing
				-- In proto3: message fields always have presence, optional fields have presence, oneof fields have presence
				-- Only non-optional primitive fields lack presence
				SET has_presence = proto3_optional OR (JSON_EXTRACT(field_descriptor, '$.\"9\"') IS NOT NULL) OR (field_type = 11); -- oneof_index or message type

				-- Only emit defaults for non-presence-sensing fields
				IF NOT has_presence THEN
					IF is_repeated THEN
						-- Empty array for repeated fields
						SET result = JSON_SET(result, CONCAT('$.', target_field_name), JSON_ARRAY());
					ELSE
						-- Default values for singular fields
						CASE field_type
						WHEN 14 THEN -- enum
							-- Get the first (zero) enum value
							CALL _pb_convert_number_enum_to_json(descriptor_set_json, field_type_name, 0, enum_string_value);
							SET result = JSON_SET(result, CONCAT('$.', target_field_name), enum_string_value);
						WHEN 11 THEN -- message
							-- Check if it's a well-known type
							CALL _pb_is_well_known_type(field_type_name, @is_wkt);
							IF @is_wkt THEN
								-- For WKTs, use empty object as default (could be improved)
								SET result = JSON_SET(result, CONCAT('$.', target_field_name), JSON_OBJECT());
							ELSE
								-- Recursively convert empty nested message
								CALL _pb_number_json_to_json_proc(descriptor_set_json, field_type_name, JSON_OBJECT(), emit_default_values, nested_json);
								SET result = JSON_SET(result, CONCAT('$.', target_field_name), nested_json);
							END IF;
						ELSE
							-- Use the existing function for primitive types (false = don't emit 64bit as numbers, use strings)
							SET converted_value = _pb_get_proto3_default_value(field_type, false);
							SET result = JSON_SET(result, CONCAT('$.', target_field_name), converted_value);
						END CASE;
					END IF;
				END IF;
			END IF;
		END IF;

		SET field_index = field_index + 1;
	END WHILE field_loop;
END $$

-- Public function interface
DROP FUNCTION IF EXISTS _pb_number_json_to_json $$
CREATE FUNCTION _pb_number_json_to_json(descriptor_set_json JSON, type_name TEXT, number_json JSON, emit_default_values BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;
	DECLARE result JSON;

	-- Validate type name starts with dot
	IF type_name NOT LIKE '.%' THEN
		SET message_text = CONCAT('_pb_number_json_to_json: type name `', type_name, '` must start with a dot');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF number_json IS NULL THEN
		RETURN NULL;
	END IF;

	CALL _pb_number_json_to_json_proc(descriptor_set_json, type_name, number_json, emit_default_values, result);
	RETURN result;
END $$

DELIMITER ;