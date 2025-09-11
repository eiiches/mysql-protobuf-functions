DELIMITER $$

-- Helper function to check if a type is a well-known type
-- TODO: Wrong and deprecated. Don't use this.
DROP FUNCTION IF EXISTS _pb_is_well_known_type $$
CREATE FUNCTION _pb_is_well_known_type(full_type_name TEXT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	RETURN full_type_name LIKE '.google.protobuf.%';
END $$

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
	SET type_paths = JSON_EXTRACT(enum_type_index, CONCAT('$.', JSON_QUOTE(full_enum_type_name)));
	IF type_paths IS NULL THEN
		SET message_text = CONCAT('_pb_convert_number_enum_to_json: enum type not found: ', full_enum_type_name);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Extract enum number index from EnumTypeIndex message
	SET enum_number_index = JSON_EXTRACT(type_paths, '$.\"4\"');

	-- Use number index for O(1) lookup
	SET found_index = JSON_EXTRACT(enum_number_index, CONCAT('$.', JSON_QUOTE(CAST(enum_numeric_value AS CHAR))));

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
-- TODO: this needs to be rewritten; number_json format uses numbers for 64bit integer types,
--   but we should still be able to handle strings, etc.
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
		IF _pb_is_well_known_type(field_type_name) THEN
			CALL _pb_convert_number_json_to_wkt(field_type, field_type_name, field_number_json_value, JSON_ARRAY(descriptor_set_json), emit_default_values, converted_value);
		ELSE
			SET enum_numeric_value = JSON_EXTRACT(field_number_json_value, '$');
			SET converted_value = _pb_convert_number_enum_to_json(descriptor_set_json, field_type_name, enum_numeric_value);
		END IF;
	WHEN 11 THEN -- message
		-- Check if it's a well-known type
		IF _pb_is_well_known_type(field_type_name) THEN
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
	SET key_field_descriptor = JSON_EXTRACT(map_entry_descriptor, '$."2"[0]'); -- field 1 (key) -- FIXME: don't assume specific index
	SET value_field_descriptor = JSON_EXTRACT(map_entry_descriptor, '$."2"[1]'); -- field 2 (value) -- FIXME: don't assume specific index

	-- Get value field type information
	SET value_field_type = _pb_field_descriptor_proto_get_type(value_field_descriptor);
	SET value_field_type_name = _pb_field_descriptor_proto_get_type_name__or(value_field_descriptor, NULL);

	-- Initialize result object
	SET result = JSON_OBJECT();

	-- Get all keys from the input map
	SET map_keys = JSON_KEYS(map_number_json);
	SET key_count = JSON_LENGTH(map_keys);
	SET key_index = 0;

	-- Process each key-value pair
	WHILE key_index < key_count DO
		SET current_key = JSON_UNQUOTE(JSON_EXTRACT(map_keys, CONCAT('$[', key_index, ']')));
		SET current_value = JSON_EXTRACT(map_number_json, CONCAT('$.', JSON_QUOTE(current_key)));

		-- Convert the value based on its type, handling null values with appropriate defaults
		IF current_value IS NULL THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'BUG: Values in JSON object cannot be SQL NULL';
		ELSE
			-- Value is not null, convert using unified singular field conversion
			CALL _pb_convert_singular_field_from_number_json(descriptor_set_json, value_field_type, value_field_type_name, current_value, emit_default_values, converted_value);
		END IF;

		-- Add to result object
		SET result = JSON_SET(result, CONCAT('$.', JSON_QUOTE(current_key)), converted_value);

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

	-- Get file descriptor to determine syntax
	SET file_descriptor = _pb_descriptor_set_get_file_descriptor(descriptor_set_json, full_type_name);
	SET syntax = _pb_file_descriptor_proto_get_syntax__or(file_descriptor, 'proto2');

	-- Process each field in the message descriptor
	SET field_index = 0;
	SET field_count = _pb_descriptor_proto_count_field(message_descriptor);
	field_loop: WHILE field_index < field_count DO
		SET field_descriptor = _pb_descriptor_proto_get_field(message_descriptor, field_index);

		-- Extract field metadata using protobuf field numbers
		SET field_number = _pb_field_descriptor_proto_get_number(field_descriptor);
		SET field_name = _pb_field_descriptor_proto_get_name(field_descriptor);
		SET field_label = _pb_field_descriptor_proto_get_label(field_descriptor);
		SET field_type = _pb_field_descriptor_proto_get_type(field_descriptor);
		SET field_type_name = _pb_field_descriptor_proto_get_type_name__or(field_descriptor, NULL);
		SET json_name = _pb_field_descriptor_proto_get_json_name__or(field_descriptor, NULL);
		SET proto3_optional = _pb_field_descriptor_proto_get_proto3_optional(field_descriptor);
		SET oneof_index = _pb_field_descriptor_proto_get_oneof_index__or(field_descriptor, NULL);

		-- Determine target field name (json_name takes precedence over field_name)
		SET target_field_name = COALESCE(json_name, field_name);
		SET is_repeated = (field_label = 3);

		-- Check if this is a map field
		SET is_map = FALSE;
		IF is_repeated AND field_type = 11 AND field_type_name IS NOT NULL THEN -- TYPE_MESSAGE
			SET map_entry_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, field_type_name);
			SET is_map = COALESCE(CAST(JSON_EXTRACT(map_entry_descriptor, '$."7"."7"') AS UNSIGNED), FALSE); -- map_entry
		END IF;

		-- Check if field exists in source JSON (by field number) and extract it
		SET field_json_value = JSON_EXTRACT(number_json, CONCAT('$.', JSON_QUOTE(CAST(field_number AS CHAR))));

		-- Check for unsupported field types first
		IF field_type = 10 THEN
			IF field_json_value IS NOT NULL THEN -- TYPE_GROUP (unsupported)
				SET message_text = CONCAT('_pb_number_json_to_json: unsupported field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
				SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
			END IF;
			SET field_index = field_index + 1;
			ITERATE field_loop;
		END IF;

		IF is_map THEN
			IF field_json_value IS NULL THEN
				-- Field is missing from number JSON - emit default value if requested
				-- Maps don't have field presence, so always emit defaults when requested
				IF emit_default_values THEN
					-- Empty object for map fields
					SET result = JSON_SET(result, CONCAT('$.', JSON_QUOTE(target_field_name)), JSON_OBJECT());
				END IF;
				SET field_index = field_index + 1;
				ITERATE field_loop;
			END IF;

			-- Handle map fields: convert object keys/values properly
			-- For maps, field_json_value is an object like {"key1": value1, "key2": value2}
			-- We need to convert the values based on the map value type
			CALL _pb_convert_map_number_json_to_proto_json(descriptor_set_json, field_type_name, field_json_value, emit_default_values, nested_json);
			SET result = JSON_SET(result, CONCAT('$.', JSON_QUOTE(target_field_name)), nested_json);
		ELSEIF is_repeated THEN
			IF field_json_value IS NULL THEN
				-- Field is missing from number JSON - emit default value if requested
				-- Repeated fields don't have field presence, so always emit defaults when requested
				IF emit_default_values THEN
					-- Empty array for repeated fields
					SET result = JSON_SET(result, CONCAT('$.', JSON_QUOTE(target_field_name)), JSON_ARRAY());
				END IF;
				SET field_index = field_index + 1;
				ITERATE field_loop;
			END IF;

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

			SET result = JSON_SET(result, CONCAT('$.', JSON_QUOTE(target_field_name)), converted_array);
		ELSE
			SET has_presence = (syntax = 'proto2' AND field_label <> 3) -- proto2: all non-repeated fields
				OR (syntax = 'proto3'
					AND (
						(field_label = 1 AND proto3_optional) -- proto3 optional
						OR (field_label <> 3 AND field_type = 11) -- message fields
						OR (oneof_index IS NOT NULL) -- oneof fields
				));

			IF field_json_value IS NULL THEN
				-- Field is missing from number JSON - emit default value if requested for non-optional fields
				-- Only emit defaults for non-presence-sensing fields
				IF emit_default_values AND NOT has_presence THEN
					-- Default values for singular fields
					CASE field_type
					WHEN 14 THEN -- enum
						-- Get the first (zero) enum value
						SET converted_value = _pb_convert_number_enum_to_json(descriptor_set_json, field_type_name, 0);
						SET result = JSON_SET(result, CONCAT('$.', JSON_QUOTE(target_field_name)), converted_value);
					WHEN 11 THEN -- message
						SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'BUG: this never happens; message fields always have field presence';
					ELSE
						-- Use the existing function for primitive types (false = don't emit 64bit as numbers, use strings)
						SET converted_value = _pb_json_get_proto3_default_value(field_type, false);
						SET result = JSON_SET(result, CONCAT('$.', JSON_QUOTE(target_field_name)), converted_value);
					END CASE;
				END IF;
				SET field_index = field_index + 1;
				ITERATE field_loop;
			END IF;

			-- Handle singular fields using unified conversion
			CALL _pb_convert_singular_field_from_number_json(descriptor_set_json, field_type, field_type_name, field_json_value, emit_default_values, converted_value);
			SET result = JSON_SET(result, CONCAT('$.', JSON_QUOTE(target_field_name)), converted_value);
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
