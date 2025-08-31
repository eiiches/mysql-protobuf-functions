DELIMITER $$

-- Helper function to convert enum numeric value to JSON (string name for known values, number for unknown values)
DROP FUNCTION IF EXISTS _pb_convert_number_enum_to_json $$
CREATE FUNCTION _pb_convert_number_enum_to_json(
	descriptor_set_json JSON,
	full_enum_type_name TEXT,
	enum_numeric_value INT
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;
	DECLARE enum_type_index JSON;
	DECLARE type_paths JSON;
	DECLARE enum_number_index JSON;
	DECLARE enum_descriptor JSON;
	DECLARE values_array JSON;
	DECLARE found_index INT;
	DECLARE value_descriptor JSON;
	DECLARE value_name TEXT;

	-- Get enum type index (field 3 from DescriptorSet)
	SET enum_type_index = JSON_EXTRACT(descriptor_set_json, '$.\"3\"');
	IF enum_type_index IS NULL THEN
		SET message_text = CONCAT('_pb_convert_number_enum_to_json: enum type index not found in descriptor set');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get paths for the enum type
	SET type_paths = JSON_EXTRACT(enum_type_index, CONCAT('$.\"', full_enum_type_name, '\"'));
	IF type_paths IS NULL THEN
		SET message_text = CONCAT('_pb_convert_number_enum_to_json: enum type not found: ', full_enum_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Extract enum number index from EnumTypeIndex message
	SET enum_number_index = JSON_EXTRACT(type_paths, '$.\"4\"');

	-- Use number index for O(1) lookup
	SET found_index = JSON_EXTRACT(enum_number_index, CONCAT('$.\"', enum_numeric_value, '\"'));

	IF found_index IS NOT NULL THEN
		-- Get enum descriptor and values array to extract the name
		SET enum_descriptor = _pb_descriptor_set_get_enum_descriptor(descriptor_set_json, full_enum_type_name);
		SET values_array = JSON_EXTRACT(enum_descriptor, '$."2"');
		SET value_descriptor = JSON_EXTRACT(values_array, CONCAT('$[', found_index, ']'));
		SET value_name = JSON_UNQUOTE(JSON_EXTRACT(value_descriptor, '$."1"')); -- name field
		RETURN JSON_QUOTE(value_name);
	ELSE
		-- If not found, return the numeric value as JSON number for Proto3 unknown enum values
		-- For Proto3, unknown enum values should be serialized as their numeric values
		RETURN CAST(enum_numeric_value AS JSON);
	END IF;
END $$

-- Helper procedure to convert singular field value from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_convert_singular_field_from_number_json $$
CREATE PROCEDURE _pb_convert_singular_field_from_number_json(
	IN descriptor_set_json JSON,
	IN field_type INT,
	IN field_type_name TEXT,
	IN field_number_json_value JSON,
	IN emit_default_values BOOLEAN,
	OUT converted_value JSON
)
BEGIN
	DECLARE enum_numeric_value INT;
	DECLARE nested_json JSON;
	DECLARE str_value TEXT;

	CASE field_type
	WHEN 14 THEN -- enum
		-- Check if it's a well-known type
		CALL _pb_is_well_known_type(field_type_name, @is_wkt);
		IF @is_wkt THEN
			CALL _pb_convert_number_json_to_wkt(field_type, field_type_name, field_number_json_value, JSON_ARRAY(descriptor_set_json), emit_default_values, converted_value);
		ELSE
			SET enum_numeric_value = JSON_EXTRACT(field_number_json_value, '$');
			SET converted_value = _pb_convert_number_enum_to_json(descriptor_set_json, field_type_name, enum_numeric_value);
		END IF;
	WHEN 11 THEN -- message
		-- Check if it's a well-known type
		CALL _pb_is_well_known_type(field_type_name, @is_wkt);
		IF @is_wkt THEN
			CALL _pb_convert_number_json_to_wkt(field_type, field_type_name, field_number_json_value, JSON_ARRAY(descriptor_set_json), emit_default_values, converted_value);
		ELSE
			-- Recursively convert nested message
			CALL _pb_number_json_to_json_proc(descriptor_set_json, field_type_name, field_number_json_value, emit_default_values, nested_json);
			SET converted_value = nested_json;
		END IF;
	WHEN 3 THEN -- int64 (convert number to string)
		SET converted_value = JSON_QUOTE(CAST(field_number_json_value AS CHAR));
	WHEN 4 THEN -- uint64 (convert number to string)
		SET converted_value = JSON_QUOTE(CAST(field_number_json_value AS CHAR));
	WHEN 6 THEN -- fixed64 (convert number to string)
		SET converted_value = JSON_QUOTE(CAST(field_number_json_value AS CHAR));
	WHEN 16 THEN -- sfixed64 (convert number to string)
		SET converted_value = JSON_QUOTE(CAST(field_number_json_value AS CHAR));
	WHEN 18 THEN -- sint64 (convert number to string)
		SET converted_value = JSON_QUOTE(CAST(field_number_json_value AS CHAR));
	WHEN 1 THEN -- double (check for IEEE 754 binary format)
		SET converted_value = _pb_convert_double_uint64_to_json(_pb_json_parse_double_as_uint64(field_number_json_value, TRUE));
	WHEN 2 THEN -- float (check for IEEE 754 binary format)
		SET converted_value = _pb_convert_float_uint32_to_json(_pb_json_parse_float_as_uint32(field_number_json_value, TRUE));
	ELSE
		-- Other primitive types stay the same
		SET converted_value = field_number_json_value;
	END CASE;
END $$

-- Helper procedure to convert map fields from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_convert_map_number_json_to_proto_json $$
CREATE PROCEDURE _pb_convert_map_number_json_to_proto_json(
	IN descriptor_set_json JSON,
	IN map_entry_type_name TEXT,
	IN map_number_json JSON,
	IN emit_default_values BOOLEAN,
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
	DECLARE result JSON;

	-- Get the map entry descriptor
	SET map_entry_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, map_entry_type_name);

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

		-- Convert the value based on its type, handling null values with appropriate defaults
		IF current_value IS NULL THEN
			-- Handle null values with appropriate defaults based on type
			CASE value_field_type
			WHEN 14 THEN -- enum
				-- Default enum value is first enum value (typically 0)
				SET converted_value = JSON_QUOTE('');
			WHEN 11 THEN -- message
				-- Default message value is empty object
				SET converted_value = JSON_OBJECT();
			WHEN 3 THEN -- int64
				SET converted_value = JSON_QUOTE('0');
			WHEN 4 THEN -- uint64
				SET converted_value = JSON_QUOTE('0');
			WHEN 6 THEN -- fixed64
				SET converted_value = JSON_QUOTE('0');
			WHEN 16 THEN -- sfixed64
				SET converted_value = JSON_QUOTE('0');
			WHEN 18 THEN -- sint64
				SET converted_value = JSON_QUOTE('0');
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
				-- Unknown type, use appropriate default
				SET converted_value = JSON_QUOTE('');
			END CASE;
		ELSE
			-- Value is not null, convert using unified singular field conversion
			CALL _pb_convert_singular_field_from_number_json(descriptor_set_json, value_field_type, value_field_type_name, current_value, emit_default_values, converted_value);
		END IF;

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
proc: BEGIN
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
	-- WKT handling
	DECLARE wkt_result JSON;

	-- Set recursion limit for nested message processing
	SET @@SESSION.max_sp_recursion_depth = 255;

	-- Check if this is a well-known type and handle it specially
	CALL _pb_convert_number_json_to_wkt(11, full_type_name, number_json, JSON_ARRAY(descriptor_set_json), emit_default_values, wkt_result);
	IF wkt_result IS NOT NULL THEN
		SET result = wkt_result;
		LEAVE proc;
	END IF;

	-- Initialize result as empty object
	SET result = JSON_OBJECT();

	-- Get message descriptor
	SET message_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, full_type_name);

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
		SET proto3_optional = COALESCE(CAST(JSON_EXTRACT(field_descriptor, '$."17"') AS UNSIGNED), FALSE); -- proto3_optional

		-- Determine target field name (json_name takes precedence over field_name)
		SET target_field_name = COALESCE(json_name, field_name);
		SET is_repeated = (field_label = 3);

		-- Check if this is a map field
		SET is_map = FALSE;
		IF field_type = 11 AND field_type_name IS NOT NULL THEN -- TYPE_MESSAGE
			SET map_entry_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, field_type_name);
			SET is_map = COALESCE(CAST(JSON_EXTRACT(map_entry_descriptor, '$."7"."7"') AS UNSIGNED), FALSE); -- map_entry
		END IF;

		-- Check if field exists in source JSON (by field number)
		IF JSON_CONTAINS_PATH(number_json, 'one', CONCAT('$."', CAST(field_number AS CHAR), '"')) THEN
			SET field_json_value = JSON_EXTRACT(number_json, CONCAT('$."', CAST(field_number AS CHAR), '"'));

			IF is_map THEN
				-- Handle map fields: convert object keys/values properly
				-- For maps, field_json_value is an object like {"key1": value1, "key2": value2}
				-- We need to convert the values based on the map value type
				CALL _pb_convert_map_number_json_to_proto_json(descriptor_set_json, field_type_name, field_json_value, emit_default_values, nested_json);
				SET result = JSON_SET(result, CONCAT('$.', target_field_name), nested_json);
			ELSEIF is_repeated THEN
				-- Handle repeated fields (arrays)
				SET array_value = field_json_value;
				SET array_length = JSON_LENGTH(array_value);
				SET converted_array = JSON_ARRAY();
				SET array_index = 0;

				array_loop: WHILE array_index < array_length DO
					SET array_element = JSON_EXTRACT(array_value, CONCAT('$[', array_index, ']'));

					-- Convert element using unified singular field conversion
					CALL _pb_convert_singular_field_from_number_json(descriptor_set_json, field_type, field_type_name, array_element, emit_default_values, converted_value);
					SET converted_array = JSON_ARRAY_APPEND(converted_array, '$', converted_value);

					SET array_index = array_index + 1;
				END WHILE array_loop;

				SET result = JSON_SET(result, CONCAT('$.', target_field_name), converted_array);
			ELSE
				-- Handle singular fields using unified conversion
				CALL _pb_convert_singular_field_from_number_json(descriptor_set_json, field_type, field_type_name, field_json_value, emit_default_values, converted_value);
				SET result = JSON_SET(result, CONCAT('$.', target_field_name), converted_value);
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
							SET converted_value = _pb_convert_number_enum_to_json(descriptor_set_json, field_type_name, 0);
							SET result = JSON_SET(result, CONCAT('$.', target_field_name), converted_value);
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
							SET converted_value = _pb_json_get_proto3_default_value(field_type, false);
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
CREATE FUNCTION _pb_number_json_to_json(descriptor_set_json JSON, type_name TEXT, number_json JSON, json_marshal_options JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;
	DECLARE result JSON;
	DECLARE emit_default_values BOOLEAN DEFAULT FALSE;

	-- Extract emit_default_values from json_marshal_options
	IF json_marshal_options IS NOT NULL THEN
		SET emit_default_values = pb_json_marshal_options_get_emit_default_values(json_marshal_options);
	END IF;

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
