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
	DECLARE wrapped_value JSON;
	-- Variables for Any handling
	DECLARE type_url TEXT;
	DECLARE any_data TEXT;
	-- Variables for FieldMask handling
	DECLARE paths_array JSON;
	DECLARE path_count INT;
	DECLARE path_index INT;
	DECLARE current_path TEXT;
	DECLARE result_str TEXT;
	-- Variables for Struct handling
	DECLARE struct_fields JSON;
	DECLARE struct_keys JSON;
	DECLARE struct_key_count INT;
	DECLARE struct_key_index INT;
	DECLARE struct_key_name TEXT;
	DECLARE struct_value_json JSON;
	DECLARE struct_converted_value JSON;
	DECLARE struct_result JSON;
	-- Variables for ListValue handling
	DECLARE list_values JSON;
	DECLARE list_length INT;
	DECLARE list_index INT;
	DECLARE list_element_json JSON;
	DECLARE list_converted_value JSON;
	DECLARE list_result JSON;

	CASE full_type_name
	WHEN '.google.protobuf.Timestamp' THEN
		-- Convert {seconds, nanos} to ISO 8601 timestamp
		CALL _pb_wkt_timestamp_number_json_to_json(number_json_value, proto_json_value);

	WHEN '.google.protobuf.Duration' THEN
		-- Convert {seconds, nanos} to duration string like "3.5s"
		CALL _pb_wkt_duration_number_json_to_json(number_json_value, proto_json_value);

	WHEN '.google.protobuf.StringValue' THEN
		-- {"1": "value"} becomes unwrapped "value", {} becomes ""
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = JSON_QUOTE('');
		ELSE
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."1"');
		END IF;

	WHEN '.google.protobuf.Int64Value' THEN
		-- {"1": value} becomes unwrapped "value" (as string for 64-bit), {} becomes "0"
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = JSON_QUOTE('0');
		ELSE
			SET wrapped_value = JSON_EXTRACT(number_json_value, '$."1"');
			SET proto_json_value = JSON_QUOTE(CAST(wrapped_value AS CHAR));
		END IF;

	WHEN '.google.protobuf.UInt64Value' THEN
		-- {"1": value} becomes unwrapped "value" (as string for 64-bit), {} becomes "0"
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = JSON_QUOTE('0');
		ELSE
			SET wrapped_value = JSON_EXTRACT(number_json_value, '$."1"');
			SET proto_json_value = JSON_QUOTE(CAST(wrapped_value AS CHAR));
		END IF;

	WHEN '.google.protobuf.Int32Value' THEN
		-- {"1": value} becomes unwrapped value (as number for 32-bit), {} becomes 0
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = CAST(0 AS JSON);
		ELSE
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."1"');
		END IF;

	WHEN '.google.protobuf.UInt32Value' THEN
		-- {"1": value} becomes unwrapped value (as number for 32-bit), {} becomes 0
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = CAST(0 AS JSON);
		ELSE
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."1"');
		END IF;

	WHEN '.google.protobuf.BoolValue' THEN
		-- {"1": value} becomes unwrapped value, {} becomes false
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = CAST(false AS JSON);
		ELSE
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."1"');
		END IF;

	WHEN '.google.protobuf.FloatValue' THEN
		-- {"1": value} becomes unwrapped value, {} becomes 0.0
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = CAST(0.0 AS JSON);
		ELSE
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."1"');
		END IF;

	WHEN '.google.protobuf.DoubleValue' THEN
		-- {"1": value} becomes unwrapped value, {} becomes 0.0
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = CAST(0.0 AS JSON);
		ELSE
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."1"');
		END IF;

	WHEN '.google.protobuf.BytesValue' THEN
		-- {"1": "value"} becomes unwrapped "value", {} becomes ""
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = JSON_QUOTE('');
		ELSE
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."1"');
		END IF;

	WHEN '.google.protobuf.Empty' THEN
		-- Empty object stays empty
		SET proto_json_value = JSON_OBJECT();

	WHEN '.google.protobuf.Value' THEN
		-- Convert Value to its unwrapped form
		-- Value has oneof fields: null_value(1), number_value(2), string_value(3), bool_value(4), struct_value(5), list_value(6)
		IF JSON_LENGTH(number_json_value) = 0 THEN
			SET proto_json_value = NULL;
		ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."1"') THEN
			-- null_value
			SET proto_json_value = NULL;
		ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."2"') THEN
			-- number_value
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."2"');
		ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."3"') THEN
			-- string_value
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."3"');
		ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."4"') THEN
			-- bool_value
			SET proto_json_value = JSON_EXTRACT(number_json_value, '$."4"');
		ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."5"') THEN
			-- struct_value - recursively convert
			CALL _pb_convert_number_json_to_wkt('.google.protobuf.Struct', JSON_EXTRACT(number_json_value, '$."5"'), proto_json_value);
		ELSEIF JSON_CONTAINS_PATH(number_json_value, 'one', '$."6"') THEN
			-- list_value - recursively convert
			CALL _pb_convert_number_json_to_wkt('.google.protobuf.ListValue', JSON_EXTRACT(number_json_value, '$."6"'), proto_json_value);
		ELSE
			SET proto_json_value = NULL;
		END IF;

	WHEN '.google.protobuf.Struct' THEN
		-- Convert Struct {"1": {field_map}} to {key: value, key: value}
		SET struct_fields = JSON_EXTRACT(number_json_value, '$."1"');
		IF struct_fields IS NULL OR JSON_LENGTH(struct_fields) = 0 THEN
			SET proto_json_value = JSON_OBJECT();
		ELSE
			SET struct_keys = JSON_KEYS(struct_fields);
			SET struct_key_count = JSON_LENGTH(struct_keys);
			SET struct_key_index = 0;
			SET struct_result = JSON_OBJECT();

			struct_loop: WHILE struct_key_index < struct_key_count DO
				SET struct_key_name = JSON_UNQUOTE(JSON_EXTRACT(struct_keys, CONCAT('$[', struct_key_index, ']')));
				SET struct_value_json = JSON_EXTRACT(struct_fields, CONCAT('$."', struct_key_name, '"'));
				-- Recursively convert the Value
				CALL _pb_convert_number_json_to_wkt('.google.protobuf.Value', struct_value_json, struct_converted_value);
				SET struct_result = JSON_SET(struct_result, CONCAT('$.', struct_key_name), struct_converted_value);
				SET struct_key_index = struct_key_index + 1;
			END WHILE struct_loop;

			SET proto_json_value = struct_result;
		END IF;

	WHEN '.google.protobuf.ListValue' THEN
		-- Convert ListValue {"1": [values]} to [value, value, value]
		SET list_values = JSON_EXTRACT(number_json_value, '$."1"');
		IF list_values IS NULL OR JSON_LENGTH(list_values) = 0 THEN
			SET proto_json_value = JSON_ARRAY();
		ELSE
			SET list_length = JSON_LENGTH(list_values);
			SET list_index = 0;
			SET list_result = JSON_ARRAY();

			list_loop: WHILE list_index < list_length DO
				SET list_element_json = JSON_EXTRACT(list_values, CONCAT('$[', list_index, ']'));
				-- Recursively convert the Value
				CALL _pb_convert_number_json_to_wkt('.google.protobuf.Value', list_element_json, list_converted_value);
				SET list_result = JSON_ARRAY_APPEND(list_result, '$', list_converted_value);
				SET list_index = list_index + 1;
			END WHILE list_loop;

			SET proto_json_value = list_result;
		END IF;

	WHEN '.google.protobuf.FieldMask' THEN
		-- Convert {"1": ["path1", "path2"]} to "path1,path2"
		SET paths_array = JSON_EXTRACT(number_json_value, '$."1"');
		IF paths_array IS NULL OR JSON_LENGTH(paths_array) = 0 THEN
			SET proto_json_value = JSON_QUOTE('');
		ELSE
			SET path_count = JSON_LENGTH(paths_array);
			SET path_index = 0;
			SET result_str = '';

			path_loop: WHILE path_index < path_count DO
				SET current_path = JSON_UNQUOTE(JSON_EXTRACT(paths_array, CONCAT('$[', path_index, ']')));
				IF path_index > 0 THEN
					SET result_str = CONCAT(result_str, ',');
				END IF;
				SET result_str = CONCAT(result_str, current_path);
				SET path_index = path_index + 1;
			END WHILE path_loop;

			SET proto_json_value = JSON_QUOTE(result_str);
		END IF;

	WHEN '.google.protobuf.Any' THEN
		-- {"1": "url", "2": "base64data"} -> {"@type": "url", "field": "value"}
		-- This is simplified - real Any handling is more complex
		SET type_url = JSON_UNQUOTE(JSON_EXTRACT(number_json_value, '$."1"'));
		SET any_data = JSON_UNQUOTE(JSON_EXTRACT(number_json_value, '$."2"'));
		-- Convert base64 data back to object (simplified)
		SET proto_json_value = JSON_OBJECT('@type', type_url);
		-- In reality, we would decode the base64 data and merge it

	ELSE
		-- Not a well-known type, return as-is
		SET proto_json_value = number_json_value;
	END CASE;
END $$

-- Helper procedure to convert map fields from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_convert_map_number_json_to_proto_json $$
CREATE PROCEDURE _pb_convert_map_number_json_to_proto_json(
	IN descriptor_set_json JSON,
	IN map_entry_type_name TEXT,
	IN map_number_json JSON,
	OUT map_proto_json JSON
)
BEGIN
	DECLARE map_entry_descriptor JSON;
	DECLARE key_field_descriptor JSON;
	DECLARE value_field_descriptor JSON;
	DECLARE value_field_type INT;
	DECLARE value_field_type_name TEXT;
	DECLARE map_keys JSON;
	DECLARE key_count INT;
	DECLARE key_index INT;
	DECLARE current_key TEXT;
	DECLARE current_value JSON;
	DECLARE converted_value JSON;
	DECLARE enum_string_value TEXT;
	DECLARE result JSON;

	-- Get the map entry descriptor
	SET map_entry_descriptor = _pb_get_message_descriptor(descriptor_set_json, map_entry_type_name);

	-- Get key and value field descriptors (map entries always have field 1 = key, field 2 = value)
	SET key_field_descriptor = JSON_EXTRACT(map_entry_descriptor, '$."2"[0]'); -- field 1 (key)
	SET value_field_descriptor = JSON_EXTRACT(map_entry_descriptor, '$."2"[1]'); -- field 2 (value)

	-- Get value field type information
	SET value_field_type = JSON_EXTRACT(value_field_descriptor, '$."5"'); -- type
	SET value_field_type_name = JSON_UNQUOTE(JSON_EXTRACT(value_field_descriptor, '$."6"')); -- type_name

	-- Initialize result object
	SET result = JSON_OBJECT();

	-- Get all keys from the input map
	SET map_keys = JSON_KEYS(map_number_json);
	SET key_count = JSON_LENGTH(map_keys);
	SET key_index = 0;

	-- Process each key-value pair
	WHILE key_index < key_count DO
		SET current_key = JSON_UNQUOTE(JSON_EXTRACT(map_keys, CONCAT('$[', key_index, ']')));
		SET current_value = JSON_EXTRACT(map_number_json, CONCAT('$."', current_key, '"'));

		-- Convert the value based on its type
		CASE value_field_type
		WHEN 14 THEN -- enum
			-- Convert enum number to string
			CALL _pb_convert_number_enum_to_json(descriptor_set_json, value_field_type_name, current_value, enum_string_value);
			SET converted_value = JSON_QUOTE(enum_string_value);
		WHEN 11 THEN -- message
			-- Recursively convert nested message
			CALL _pb_number_json_to_json_proc(descriptor_set_json, value_field_type_name, current_value, TRUE, converted_value);
		WHEN 3 THEN -- int64 (convert number to string)
			SET converted_value = JSON_QUOTE(CAST(current_value AS CHAR));
		WHEN 4 THEN -- uint64 (convert number to string)
			SET converted_value = JSON_QUOTE(CAST(current_value AS CHAR));
		WHEN 6 THEN -- fixed64 (convert number to string)
			SET converted_value = JSON_QUOTE(CAST(current_value AS CHAR));
		WHEN 16 THEN -- sfixed64 (convert number to string)
			SET converted_value = JSON_QUOTE(CAST(current_value AS CHAR));
		WHEN 18 THEN -- sint64 (convert number to string)
			SET converted_value = JSON_QUOTE(CAST(current_value AS CHAR));
		ELSE
			-- Other types (primitives) stay the same
			SET converted_value = current_value;
		END CASE;

		-- Add to result object
		SET result = JSON_SET(result, CONCAT('$."', current_key, '"'), converted_value);

		SET key_index = key_index + 1;
	END WHILE;

	SET map_proto_json = result;
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
	-- Map handling
	DECLARE is_map BOOLEAN;
	DECLARE map_entry_descriptor JSON;

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

		-- Check if this is a map field
		SET is_map = FALSE;
		IF field_type = 11 AND field_type_name IS NOT NULL THEN -- TYPE_MESSAGE
			SET map_entry_descriptor = _pb_get_message_descriptor(descriptor_set_json, field_type_name);
			SET is_map = COALESCE(CAST(JSON_EXTRACT(map_entry_descriptor, '$."7"."7"') AS UNSIGNED), FALSE); -- map_entry
		END IF;

		-- Check if field exists in source JSON (by field number)
		IF JSON_CONTAINS_PATH(number_json, 'one', CONCAT('$."', CAST(field_number AS CHAR), '"')) THEN
			SET field_json_value = JSON_EXTRACT(number_json, CONCAT('$."', CAST(field_number AS CHAR), '"'));

			IF is_map THEN
				-- Handle map fields: convert object keys/values properly
				-- For maps, field_json_value is an object like {"key1": value1, "key2": value2}
				-- We need to convert the values based on the map value type
				CALL _pb_convert_map_number_json_to_proto_json(descriptor_set_json, field_type_name, field_json_value, nested_json);
				SET result = JSON_SET(result, CONCAT('$.', target_field_name), nested_json);
			ELSEIF is_repeated THEN
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
					WHEN 16 THEN -- sfixed64 (convert number to string)
						SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', CAST(array_element AS CHAR));
					WHEN 18 THEN -- sint64 (convert number to string)
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
				WHEN 16 THEN -- sfixed64 (convert number to string)
					SET result = JSON_SET(result, CONCAT('$.', target_field_name), CAST(field_json_value AS CHAR));
				WHEN 18 THEN -- sint64 (convert number to string)
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
				-- Exception: map fields and repeated fields should always emit default values regardless of presence
				-- Only non-optional singular primitive fields lack presence
				IF is_map OR is_repeated THEN
					SET has_presence = FALSE; -- Maps and repeated fields always emit defaults
				ELSE
					SET has_presence = proto3_optional OR (JSON_EXTRACT(field_descriptor, '$.\"9\"') IS NOT NULL) OR (field_type = 11); -- oneof_index or message type
				END IF;

				-- Only emit defaults for non-presence-sensing fields
				IF NOT has_presence THEN
					IF is_map THEN
						-- Empty object for map fields
						SET result = JSON_SET(result, CONCAT('$.', target_field_name), JSON_OBJECT());
					ELSEIF is_repeated THEN
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
							IF is_map THEN
								-- For map fields, default is empty object
								SET result = JSON_SET(result, CONCAT('$.', target_field_name), JSON_OBJECT());
							ELSE
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