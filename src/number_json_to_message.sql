DELIMITER $$

-- Helper procedure to set primitive field values in wire_json format
DROP PROCEDURE IF EXISTS _pb_json_set_primitive_field_as_wire_json $$
CREATE PROCEDURE _pb_json_set_primitive_field_as_wire_json(IN wire_json JSON, IN field_number INT, IN field_type INT, IN is_repeated BOOLEAN, IN json_value JSON, IN use_packed BOOLEAN, IN syntax TEXT, IN has_field_presence BOOLEAN, OUT result JSON)
proc: BEGIN
	DECLARE element_count INT;
	DECLARE element_index INT;
	DECLARE element JSON;
	DECLARE temp_wire_json JSON;
	DECLARE str_value TEXT;
	DECLARE hex_value TEXT;
	DECLARE uint32_bits INT UNSIGNED;
	DECLARE uint64_bits BIGINT UNSIGNED;

	SET result = wire_json;

	-- Skip encoding proto3 default values for fields without explicit presence
	IF NOT is_repeated AND syntax = 'proto3' AND NOT has_field_presence AND _pb_number_json_is_proto3_default_value(field_type, json_value) THEN
		-- Do not encode default values in proto3 for fields without explicit presence
		LEAVE proc;
	END IF;

	IF is_repeated THEN
		-- Handle repeated primitive fields
		SET element_count = JSON_LENGTH(json_value);
		SET element_index = 0;

		WHILE element_index < element_count DO
			SET element = JSON_EXTRACT(json_value, CONCAT('$[', element_index, ']'));

			CASE field_type
			WHEN 1 THEN -- TYPE_DOUBLE
				SET uint64_bits = _pb_json_parse_double_as_uint64(element, TRUE);
				-- TODO: This is a workaround and should be replaced with generated code by @cmd/protobuf-accessors/
				SET result = pb_wire_json_add_repeated_fixed64_field_element(result, field_number, uint64_bits, use_packed);
			WHEN 2 THEN -- TYPE_FLOAT
				SET uint32_bits = _pb_json_parse_float_as_uint32(element, TRUE);
				-- TODO: This is a workaround and should be replaced with generated code by @cmd/protobuf-accessors/
				SET result = pb_wire_json_add_repeated_fixed32_field_element(result, field_number, uint32_bits, use_packed);
			WHEN 3 THEN -- TYPE_INT64
				SET result = pb_wire_json_add_repeated_int64_field_element(result, field_number, _pb_json_parse_signed_int(element), use_packed);
			WHEN 4 THEN -- TYPE_UINT64
				SET result = pb_wire_json_add_repeated_uint64_field_element(result, field_number, _pb_json_parse_unsigned_int(element), use_packed);
			WHEN 5 THEN -- TYPE_INT32
				SET result = pb_wire_json_add_repeated_int32_field_element(result, field_number, _pb_json_parse_signed_int(element), use_packed);
			WHEN 6 THEN -- TYPE_FIXED64
				SET result = pb_wire_json_add_repeated_fixed64_field_element(result, field_number, _pb_json_parse_unsigned_int(element), use_packed);
			WHEN 7 THEN -- TYPE_FIXED32
				SET result = pb_wire_json_add_repeated_fixed32_field_element(result, field_number, _pb_json_parse_unsigned_int(element), use_packed);
			WHEN 8 THEN -- TYPE_BOOL
				SET result = pb_wire_json_add_repeated_bool_field_element(result, field_number, IF(element, TRUE, FALSE), use_packed);
			WHEN 9 THEN -- TYPE_STRING
				SET result = pb_wire_json_add_repeated_string_field_element(result, field_number, JSON_UNQUOTE(element));
			WHEN 10 THEN -- TYPE_SFIXED64
				SET result = pb_wire_json_add_repeated_sfixed64_field_element(result, field_number, _pb_json_parse_signed_int(element), use_packed);
			WHEN 11 THEN -- TYPE_SFIXED32
				SET result = pb_wire_json_add_repeated_sfixed32_field_element(result, field_number, _pb_json_parse_signed_int(element), use_packed);
			WHEN 12 THEN -- TYPE_BYTES
				SET result = pb_wire_json_add_repeated_bytes_field_element(result, field_number, _pb_util_from_base64_url(JSON_UNQUOTE(element)));
			WHEN 13 THEN -- TYPE_UINT32
				SET result = pb_wire_json_add_repeated_uint32_field_element(result, field_number, _pb_json_parse_unsigned_int(element), use_packed);
			WHEN 15 THEN -- TYPE_SFIXED32 (duplicate, but keeping for completeness)
				SET result = pb_wire_json_add_repeated_sfixed32_field_element(result, field_number, _pb_json_parse_signed_int(element), use_packed);
			WHEN 16 THEN -- TYPE_SFIXED64 (duplicate, but keeping for completeness)
				SET result = pb_wire_json_add_repeated_sfixed64_field_element(result, field_number, _pb_json_parse_signed_int(element), use_packed);
			WHEN 17 THEN -- TYPE_SINT32
				SET result = pb_wire_json_add_repeated_sint32_field_element(result, field_number, _pb_json_parse_signed_int(element), use_packed);
			WHEN 18 THEN -- TYPE_SINT64
				SET result = pb_wire_json_add_repeated_sint64_field_element(result, field_number, _pb_json_parse_signed_int(element), use_packed);
			END CASE;

			SET element_index = element_index + 1;
		END WHILE;
	ELSE
		-- Handle singular primitive fields
		CASE field_type
		WHEN 1 THEN -- TYPE_DOUBLE
			SET uint64_bits = _pb_json_parse_double_as_uint64(json_value, TRUE);
			-- TODO: This is a workaround and should be replaced with generated code by @cmd/protobuf-accessors/
			SET result = pb_wire_json_set_fixed64_field(result, field_number, uint64_bits);
		WHEN 2 THEN -- TYPE_FLOAT
			SET uint32_bits = _pb_json_parse_float_as_uint32(json_value, TRUE);
			-- TODO: This is a workaround and should be replaced with generated code by @cmd/protobuf-accessors/
			SET result = pb_wire_json_set_fixed32_field(result, field_number, uint32_bits);
		WHEN 3 THEN -- TYPE_INT64
			SET result = pb_wire_json_set_int64_field(result, field_number, _pb_json_parse_signed_int(json_value));
		WHEN 4 THEN -- TYPE_UINT64
			SET result = pb_wire_json_set_uint64_field(result, field_number, _pb_json_parse_unsigned_int(json_value));
		WHEN 5 THEN -- TYPE_INT32
			SET result = pb_wire_json_set_int32_field(result, field_number, _pb_json_parse_signed_int(json_value));
		WHEN 6 THEN -- TYPE_FIXED64
			SET result = pb_wire_json_set_fixed64_field(result, field_number, _pb_json_parse_unsigned_int(json_value));
		WHEN 7 THEN -- TYPE_FIXED32
			SET result = pb_wire_json_set_fixed32_field(result, field_number, _pb_json_parse_unsigned_int(json_value));
		WHEN 8 THEN -- TYPE_BOOL
			SET result = pb_wire_json_set_bool_field(result, field_number, IF(json_value, TRUE, FALSE));
		WHEN 9 THEN -- TYPE_STRING
			SET result = pb_wire_json_set_string_field(result, field_number, JSON_UNQUOTE(json_value));
		WHEN 10 THEN -- TYPE_SFIXED64
			SET result = pb_wire_json_set_sfixed64_field(result, field_number, _pb_json_parse_signed_int(json_value));
		WHEN 11 THEN -- TYPE_SFIXED32
			SET result = pb_wire_json_set_sfixed32_field(result, field_number, _pb_json_parse_signed_int(json_value));
		WHEN 12 THEN -- TYPE_BYTES
			SET result = pb_wire_json_set_bytes_field(result, field_number, _pb_util_from_base64_url(JSON_UNQUOTE(json_value)));
		WHEN 13 THEN -- TYPE_UINT32
			SET result = pb_wire_json_set_uint32_field(result, field_number, _pb_json_parse_unsigned_int(json_value));
		WHEN 15 THEN -- TYPE_SFIXED32 (duplicate, but keeping for completeness)
			SET result = pb_wire_json_set_sfixed32_field(result, field_number, _pb_json_parse_signed_int(json_value));
		WHEN 16 THEN -- TYPE_SFIXED64 (duplicate, but keeping for completeness)
			SET result = pb_wire_json_set_sfixed64_field(result, field_number, _pb_json_parse_signed_int(json_value));
		WHEN 17 THEN -- TYPE_SINT32
			SET result = pb_wire_json_set_sint32_field(result, field_number, _pb_json_parse_signed_int(json_value));
		WHEN 18 THEN -- TYPE_SINT64
			SET result = pb_wire_json_set_sint64_field(result, field_number, _pb_json_parse_signed_int(json_value));
		END CASE;
	END IF;
END $$

-- Main procedure for converting JSON to protobuf wire_json using descriptor set
DROP PROCEDURE IF EXISTS _pb_number_json_to_wire_json_proc $$
CREATE PROCEDURE _pb_number_json_to_wire_json_proc(IN descriptor_set_json JSON, IN full_type_name TEXT, IN json_value JSON, OUT result JSON)
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
	DECLARE oneof_index INT;
	DECLARE default_value TEXT;

	-- Processing variables
	DECLARE is_repeated BOOLEAN;
	DECLARE has_field_presence BOOLEAN;
	DECLARE use_packed BOOLEAN;
	DECLARE field_json_value JSON;
	DECLARE json_field_name TEXT;
	DECLARE element JSON;
	DECLARE element_count INT;
	DECLARE element_index INT;
	DECLARE nested_wire_json JSON;
	DECLARE enum_number INT;

	-- Map handling
	DECLARE is_map BOOLEAN;
	DECLARE map_entry_descriptor JSON;
	DECLARE map_key_field JSON;
	DECLARE map_value_field JSON;
	DECLARE map_key_type INT;
	DECLARE map_value_type INT;
	DECLARE map_value_type_name TEXT;
	DECLARE map_keys JSON;
	DECLARE map_key_count INT;
	DECLARE map_key_index INT;
	DECLARE map_key_name TEXT;
	DECLARE map_entry_wire_json JSON;
	DECLARE map_value_json JSON;
	DECLARE map_value_wire_json JSON;
	-- Well-known type handling
	DECLARE wkt_descriptor_set JSON;

	SET @@SESSION.max_sp_recursion_depth = 255;

	IF JSON_TYPE(json_value) = 'NULL' THEN
		-- Null value should not produce any field in protobuf
		SET result = NULL;
		LEAVE proc;
	END IF;

	-- Get message descriptor
	SET message_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, full_type_name);

	IF message_descriptor IS NULL AND full_type_name LIKE '.google.protobuf.%' THEN
		-- Try to get well-known type descriptor set
		SET wkt_descriptor_set = _pb_wkt_get_descriptor_set(full_type_name);
		IF wkt_descriptor_set IS NOT NULL THEN
			SET descriptor_set_json = wkt_descriptor_set;
			SET message_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, full_type_name);
		END IF;
	END IF;

	IF message_descriptor IS NULL THEN
		SET message_text = CONCAT('_pb_json_to_wire_json: message type `', full_type_name, '` not found in descriptor set');
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get file descriptor to determine syntax
	SET file_descriptor = _pb_descriptor_set_get_file_descriptor(descriptor_set_json, full_type_name);
	SET syntax = JSON_UNQUOTE(JSON_EXTRACT(file_descriptor, '$."12"')); -- syntax field
	IF syntax IS NULL THEN
		SET syntax = 'proto2'; -- default
	END IF;

	SET result = JSON_OBJECT();

	-- Get fields array (field 2 in DescriptorProto)
	SET fields = JSON_EXTRACT(message_descriptor, '$."2"');

	IF fields IS NOT NULL THEN
		SET field_count = JSON_LENGTH(fields);
		SET field_index = 0;

		WHILE field_index < field_count DO
			SET field_descriptor = JSON_EXTRACT(fields, CONCAT('$[', field_index, ']'));

			-- Extract field properties from FieldDescriptorProto
			SET field_number = JSON_EXTRACT(field_descriptor, '$."3"'); -- number
			SET field_name = JSON_UNQUOTE(JSON_EXTRACT(field_descriptor, '$."1"')); -- name
			SET field_label = JSON_EXTRACT(field_descriptor, '$."4"'); -- label
			SET field_type = JSON_EXTRACT(field_descriptor, '$."5"'); -- type
			SET field_type_name = JSON_UNQUOTE(JSON_EXTRACT(field_descriptor, '$."6"')); -- type_name
			SET json_name = JSON_UNQUOTE(JSON_EXTRACT(field_descriptor, '$."10"')); -- json_name
			SET proto3_optional = COALESCE(CAST(JSON_EXTRACT(field_descriptor, '$."17"') AS UNSIGNED), FALSE); -- proto3_optional
			SET oneof_index = JSON_EXTRACT(field_descriptor, '$."9"'); -- oneof_index
			SET default_value = JSON_UNQUOTE(JSON_EXTRACT(field_descriptor, '$."7"')); -- default_value

			SET is_repeated = (field_label = 3); -- LABEL_REPEATED

			-- Check if this is a map field
			SET is_map = FALSE;
			IF field_type = 11 AND field_type_name IS NOT NULL THEN -- TYPE_MESSAGE
				SET map_entry_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, field_type_name);
				SET is_map = COALESCE(CAST(JSON_EXTRACT(map_entry_descriptor, '$."7"."7"') AS UNSIGNED), FALSE); -- map_entry
			END IF;

			-- Determine field presence
			SET has_field_presence =
				(syntax = 'proto2' AND field_label <> 3) -- proto2: all non-repeated fields
				OR (syntax = 'proto3'
					AND (
						(field_label = 1 AND proto3_optional) -- proto3 optional
						OR (field_label <> 3 AND field_type = 11) -- message fields
						OR (oneof_index IS NOT NULL) -- oneof fields
					));

			-- Determine packed encoding for repeated fields
			-- Check field options for explicit packed setting (field 8.2 in FieldDescriptorProto)
			SET use_packed = CAST(JSON_EXTRACT(field_descriptor, '$."8"."2"') AS UNSIGNED);

			-- Use syntax default if field option not set
			IF use_packed IS NULL THEN
				-- Proto3: packed by default, Proto2: unpacked by default
				SET use_packed = (syntax = 'proto3');
			END IF;

			SET field_json_value = JSON_EXTRACT(json_value, CONCAT('$."', field_number, '"'));

			-- Process field if it exists in JSON
			IF field_json_value IS NOT NULL THEN
				CASE field_type
				WHEN 10 THEN -- TYPE_GROUP (unsupported)
					SET message_text = CONCAT('_pb_json_to_wire_json: unsupported field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
					SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;

				WHEN 11 THEN -- TYPE_MESSAGE
					IF is_map THEN
						-- Handle map fields - convert JSON object to repeated message entries
						SET map_keys = JSON_KEYS(field_json_value);
						SET map_key_count = JSON_LENGTH(map_keys);
						SET map_key_index = 0;

						-- Get map key/value field descriptors
						SET map_key_field = JSON_EXTRACT(map_entry_descriptor, '$."2"[0]'); -- first field (key)
						SET map_value_field = JSON_EXTRACT(map_entry_descriptor, '$."2"[1]'); -- second field (value)
						SET map_key_type = JSON_EXTRACT(map_key_field, '$."5"');
						SET map_value_type = JSON_EXTRACT(map_value_field, '$."5"');
						SET map_value_type_name = JSON_UNQUOTE(JSON_EXTRACT(map_value_field, '$."6"'));

						WHILE map_key_index < map_key_count DO
							SET map_key_name = JSON_UNQUOTE(JSON_EXTRACT(map_keys, CONCAT('$[', map_key_index, ']')));
							SET map_value_json = JSON_EXTRACT(field_json_value, CONCAT('$."', map_key_name, '"'));

							-- Create map entry with key=1, value=2
							SET map_entry_wire_json = JSON_OBJECT();

							-- Add key field (always field 1)
							-- Convert map key to proper JSON type based on key type
							-- Map keys always have presence and should always be encoded
							CASE map_key_type
							WHEN 8 THEN -- bool
								CALL _pb_json_set_primitive_field_as_wire_json(map_entry_wire_json, 1, map_key_type, FALSE, CAST((map_key_name = 'true') AS JSON), FALSE, syntax, TRUE, map_entry_wire_json);
							ELSE
								CALL _pb_json_set_primitive_field_as_wire_json(map_entry_wire_json, 1, map_key_type, FALSE, JSON_QUOTE(map_key_name), FALSE, syntax, TRUE, map_entry_wire_json);
							END CASE;

							-- Add value field (always field 2)
							IF map_value_type = 11 THEN -- message
								CALL _pb_number_json_to_wire_json_proc(descriptor_set_json, map_value_type_name, map_value_json, map_value_wire_json);
								SET map_entry_wire_json = pb_wire_json_set_message_field(map_entry_wire_json, 2, pb_wire_json_to_message(map_value_wire_json));
							ELSEIF map_value_type = 14 THEN -- enum
								SET enum_number = _pb_convert_json_enum_to_number(descriptor_set_json, map_value_type_name, map_value_json, FALSE);
								SET map_entry_wire_json = pb_wire_json_set_enum_field(map_entry_wire_json, 2, enum_number);
							ELSE
								-- Map values also always have presence in map entries
								CALL _pb_json_set_primitive_field_as_wire_json(map_entry_wire_json, 2, map_value_type, FALSE, map_value_json, FALSE, syntax, TRUE, map_entry_wire_json);
							END IF;

							-- Add map entry to result
							SET result = pb_wire_json_add_repeated_message_field_element(result, field_number, pb_wire_json_to_message(map_entry_wire_json));
							SET map_key_index = map_key_index + 1;
						END WHILE;

					ELSEIF is_repeated THEN
						-- Handle repeated message fields
						SET element_count = JSON_LENGTH(field_json_value);
						SET element_index = 0;

						WHILE element_index < element_count DO
							SET element = JSON_EXTRACT(field_json_value, CONCAT('$[', element_index, ']'));
							CALL _pb_number_json_to_wire_json_proc(descriptor_set_json, field_type_name, element, nested_wire_json);
							SET result = pb_wire_json_add_repeated_message_field_element(result, field_number, pb_wire_json_to_message(nested_wire_json));
							SET element_index = element_index + 1;
						END WHILE;
					ELSE
						-- Handle singular message fields
						CALL _pb_number_json_to_wire_json_proc(descriptor_set_json, field_type_name, field_json_value, nested_wire_json);
						IF nested_wire_json IS NOT NULL THEN
							SET result = pb_wire_json_set_message_field(result, field_number, pb_wire_json_to_message(nested_wire_json));
						END IF;
					END IF;

				WHEN 14 THEN -- TYPE_ENUM
					IF is_repeated THEN
						SET element_count = JSON_LENGTH(field_json_value);
						SET element_index = 0;

						WHILE element_index < element_count DO
							SET element = JSON_EXTRACT(field_json_value, CONCAT('$[', element_index, ']'));
							SET enum_number = _pb_convert_json_enum_to_number(descriptor_set_json, field_type_name, element, FALSE);
							SET result = pb_wire_json_add_repeated_enum_field_element(result, field_number, enum_number, use_packed);
							SET element_index = element_index + 1;
						END WHILE;
					ELSE
						SET enum_number = _pb_convert_json_enum_to_number(descriptor_set_json, field_type_name, field_json_value, FALSE);
						-- Skip encoding proto3 default values for fields without explicit presence
						IF NOT (syntax = 'proto3' AND NOT has_field_presence AND enum_number = 0) THEN
							SET result = pb_wire_json_set_enum_field(result, field_number, enum_number);
						END IF;
					END IF;

				ELSE
					-- Handle primitive types
					CALL _pb_json_set_primitive_field_as_wire_json(result, field_number, field_type, is_repeated, field_json_value, use_packed, syntax, has_field_presence, result);
				END CASE;
			END IF;

			SET field_index = field_index + 1;
		END WHILE;
	END IF;
END $$

-- Wrapper procedure for number JSON to wire_json conversion
DROP PROCEDURE IF EXISTS _pb_number_json_to_wire_json $$
CREATE PROCEDURE _pb_number_json_to_wire_json(IN descriptor_set_json JSON, IN full_type_name TEXT, IN json_value JSON, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;

	-- Validate type name starts with dot
	IF full_type_name NOT LIKE '.%' THEN
		SET message_text = CONCAT('_pb_number_json_to_wire_json: type name `', full_type_name, '` must start with a dot');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF json_value IS NULL THEN
		SET result = NULL;
	ELSE
		CALL _pb_number_json_to_wire_json_proc(descriptor_set_json, full_type_name, json_value, result);
	END IF;
END $$

-- Wrapper procedure for JSON to message conversion
DROP PROCEDURE IF EXISTS _pb_number_json_to_message $$
CREATE PROCEDURE _pb_number_json_to_message(IN descriptor_set_json JSON, IN full_type_name TEXT, IN json_value JSON, OUT result LONGBLOB)
BEGIN
	DECLARE message_text TEXT;
	DECLARE wire_json JSON;

	-- Validate type name starts with dot
	IF full_type_name NOT LIKE '.%' THEN
		SET message_text = CONCAT('_pb_number_json_to_message: type name `', full_type_name, '` must start with a dot');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF json_value IS NULL THEN
		SET result = NULL;
	ELSE
		CALL _pb_number_json_to_wire_json_proc(descriptor_set_json, full_type_name, json_value, wire_json);
		SET result = pb_wire_json_to_message(wire_json);
	END IF;
END $$

-- Private function interface for number JSON to wire_json conversion
DROP FUNCTION IF EXISTS _pb_number_json_to_wire_json $$
CREATE FUNCTION _pb_number_json_to_wire_json(descriptor_set_json JSON, type_name TEXT, json_value JSON, marshal_options JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	-- For now, marshal_options is accepted but not yet used - keeping current behavior
	CALL _pb_number_json_to_wire_json(descriptor_set_json, type_name, json_value, result);
	RETURN result;
END $$

-- Private function interface for number JSON to message conversion
DROP FUNCTION IF EXISTS _pb_number_json_to_message $$
CREATE FUNCTION _pb_number_json_to_message(descriptor_set_json JSON, type_name TEXT, json_value JSON, marshal_options JSON) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE result LONGBLOB;
	-- For now, marshal_options is accepted but not yet used - keeping current behavior
	CALL _pb_number_json_to_message(descriptor_set_json, type_name, json_value, result);
	RETURN result;
END $$
