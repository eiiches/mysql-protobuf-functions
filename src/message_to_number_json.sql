DELIMITER $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_primitive_field_as_number_json $$
CREATE PROCEDURE _pb_wire_json_get_primitive_field_as_number_json(IN wire_json JSON, IN field_number INT, IN field_type INT, IN is_repeated BOOLEAN, IN has_field_presence BOOLEAN, OUT field_json_value JSON)
BEGIN
	-- Note: This procedure is optimized for number JSON format (64-bit integers as numbers, floats as hex strings)

	DECLARE message_text TEXT;
	DECLARE boolean_value BOOLEAN;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE int_value BIGINT;
	DECLARE element_count INT;
	DECLARE element_index INT;

	CASE field_type
	WHEN 1 THEN -- double
		IF is_repeated THEN
			-- Handle repeated double fields with binary64 format (emit_floats_as_hex_strings is always TRUE)
			-- TODO: This should be replaced with generated code by @cmd/protobuf-accessors/ that directly returns binary64 format
			SET field_json_value = pb_wire_json_get_repeated_fixed64_field_as_json_array(wire_json, field_number);
			SET element_count = JSON_LENGTH(field_json_value);
			SET element_index = 0;

			WHILE element_index < element_count DO
				SET uint_value = CAST(JSON_EXTRACT(field_json_value, CONCAT('$[', element_index, ']')) AS UNSIGNED);
				SET field_json_value = JSON_SET(field_json_value, CONCAT('$[', element_index, ']'), _pb_convert_double_uint64_to_number_json(uint_value));
				SET element_index = element_index + 1;
			END WHILE;
		ELSE
			-- IEEE 754 binary format (emit_floats_as_hex_strings is always TRUE)
			-- TODO: This is a workaround and should be replaced with generated code by @cmd/protobuf-accessors/
			SET uint_value = pb_wire_json_get_fixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			IF uint_value IS NULL THEN
				SET field_json_value = NULL;
			ELSE
				SET field_json_value = _pb_convert_double_uint64_to_number_json(uint_value);
			END IF;
		END IF;
	WHEN 2 THEN -- float
		IF is_repeated THEN
			-- Handle repeated float fields with binary32 format (emit_floats_as_hex_strings is always TRUE)
			-- TODO: This should be replaced with generated code by @cmd/protobuf-accessors/ that directly returns binary32 format
			SET field_json_value = pb_wire_json_get_repeated_fixed32_field_as_json_array(wire_json, field_number);
			SET element_count = JSON_LENGTH(field_json_value);
			SET element_index = 0;

			WHILE element_index < element_count DO
				SET uint_value = CAST(JSON_EXTRACT(field_json_value, CONCAT('$[', element_index, ']')) AS UNSIGNED);
				SET field_json_value = JSON_SET(field_json_value, CONCAT('$[', element_index, ']'), _pb_convert_float_uint32_to_number_json(uint_value));
				SET element_index = element_index + 1;
			END WHILE;
		ELSE
			-- IEEE 754 binary format (emit_floats_as_hex_strings is always TRUE)
			-- TODO: This is a workaround and should be replaced with generated code by @cmd/protobuf-accessors/
			SET uint_value = pb_wire_json_get_fixed32_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			IF uint_value IS NULL THEN
				SET field_json_value = NULL;
			ELSE
				SET field_json_value = _pb_convert_float_uint32_to_number_json(uint_value);
			END IF;
		END IF;
	WHEN 3 THEN -- int64
		IF is_repeated THEN
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = pb_wire_json_get_repeated_int64_field_as_json_array(wire_json, field_number);
		ELSE
			SET int_value = pb_wire_json_get_int64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = CAST(int_value AS JSON);
		END IF;
	WHEN 4 THEN -- uint64
		IF is_repeated THEN
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = pb_wire_json_get_repeated_uint64_field_as_json_array(wire_json, field_number);
		ELSE
			SET uint_value = pb_wire_json_get_uint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = CAST(uint_value AS JSON);
		END IF;
	WHEN 5 THEN -- int32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_int32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_int32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 6 THEN -- fixed64
		IF is_repeated THEN
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = pb_wire_json_get_repeated_fixed64_field_as_json_array(wire_json, field_number);
		ELSE
			SET uint_value = pb_wire_json_get_fixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = CAST(uint_value AS JSON);
		END IF;
	WHEN 7 THEN -- fixed32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_fixed32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_fixed32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 8 THEN -- bool
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_bool_field_as_json_array(wire_json, field_number);
		ELSE
			SET boolean_value = pb_wire_json_get_bool_field(wire_json, field_number, IF(has_field_presence, NULL, FALSE));
			IF boolean_value IS NULL THEN
				SET field_json_value = NULL;
			ELSE
				-- See https://bugs.mysql.com/bug.php?id=79813
				SET field_json_value = CAST((boolean_value IS TRUE) AS JSON);
			END IF;
		END IF;
	WHEN 9 THEN -- string
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_string_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(pb_wire_json_get_string_field(wire_json, field_number, IF(has_field_presence, NULL, '')));
		END IF;
	WHEN 12 THEN -- bytes
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_bytes_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(_pb_to_base64(pb_wire_json_get_bytes_field(wire_json, field_number, IF(has_field_presence, NULL, _binary X''))));
		END IF;
	WHEN 13 THEN -- uint32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_uint32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_uint32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 15 THEN -- sfixed32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_sfixed32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_sfixed32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 16 THEN -- sfixed64
		IF is_repeated THEN
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = pb_wire_json_get_repeated_sfixed64_field_as_json_array(wire_json, field_number);
		ELSE
			SET int_value = pb_wire_json_get_sfixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = CAST(int_value AS JSON);
		END IF;
	WHEN 17 THEN -- sint32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_sint32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_sint32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 18 THEN -- sint64
		IF is_repeated THEN
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = pb_wire_json_get_repeated_sint64_field_as_json_array(wire_json, field_number);
		ELSE
			SET int_value = pb_wire_json_get_sint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			-- emit_64bit_integers_as_numbers is always TRUE
			SET field_json_value = CAST(int_value AS JSON);
		END IF;
	ELSE
		SET message_text = CONCAT('_pb_message_to_json: unknown field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$

-- Main procedure for converting protobuf message to number JSON using descriptor set
DROP PROCEDURE IF EXISTS _pb_wire_json_to_number_json_proc $$
CREATE PROCEDURE _pb_wire_json_to_number_json_proc(IN descriptor_set_json JSON, IN full_type_name TEXT, IN wire_json JSON, OUT result JSON)
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
	DECLARE has_field_presence BOOLEAN;
	DECLARE field_json_value JSON;
	DECLARE json_field_name TEXT;
	DECLARE bytes_value LONGBLOB;
	DECLARE nested_json_value JSON;
	DECLARE elements JSON;
	DECLARE element JSON;
	DECLARE element_count INT;
	DECLARE element_index INT;
	DECLARE field_enum_value INT;

	-- Map handling
	DECLARE is_map BOOLEAN;
	DECLARE map_entry_descriptor JSON;
	DECLARE map_key_field JSON;
	DECLARE map_value_field JSON;
	DECLARE map_key_type INT;
	DECLARE map_value_type INT;
	DECLARE map_value_type_name TEXT;
	DECLARE map_key JSON;
	DECLARE map_value JSON;

	-- Oneof handling
	DECLARE oneofs JSON;
	DECLARE oneof_priority INT;
	DECLARE oneof_priority_prev INT;

	DECLARE wkt_descriptor_set JSON;

	SET @@SESSION.max_sp_recursion_depth = 255;

	IF wire_json IS NULL THEN
		SET result = NULL;
		LEAVE proc;
	END IF;

	-- Handle well-known types first (only for regular JSON, not number JSON)
	SET wkt_descriptor_set = _pb_wkt_get_descriptor_set(full_type_name);
	IF wkt_descriptor_set IS NOT NULL THEN
		SET descriptor_set_json = wkt_descriptor_set;
	END IF;

	-- Get message descriptor
	SET message_descriptor = _pb_descriptor_set_get_message_descriptor(descriptor_set_json, full_type_name);
	IF message_descriptor IS NULL THEN
		SET message_text = CONCAT('_pb_wire_json_to_json: message type `', full_type_name, '` not found in descriptor set');
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get file descriptor to determine syntax
	SET file_descriptor = _pb_descriptor_set_get_file_descriptor(descriptor_set_json, full_type_name);
	IF file_descriptor IS NULL THEN
		SET message_text = CONCAT('_pb_wire_json_to_json: file descriptor not found for type `', full_type_name, '`');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	SET syntax = _pb_file_descriptor_proto_get_syntax__or(file_descriptor, 'proto2');

	SET result = JSON_OBJECT();
	SET oneofs = JSON_OBJECT();

	-- Process each field in the message descriptor
	SET field_index = 0;
	SET field_count = _pb_descriptor_proto_count_field(message_descriptor);
	WHILE field_index < field_count DO
		SET field_descriptor = _pb_descriptor_proto_get_field(message_descriptor, field_index);

		-- Extract field properties from FieldDescriptorProto
		SET field_number = _pb_field_descriptor_proto_get_number(field_descriptor);
		SET field_name = _pb_field_descriptor_proto_get_name(field_descriptor);
		SET field_label = _pb_field_descriptor_proto_get_label(field_descriptor);
		SET field_type = _pb_field_descriptor_proto_get_type(field_descriptor);
		SET field_type_name = _pb_field_descriptor_proto_get_type_name__or(field_descriptor, NULL);
		SET json_name = _pb_field_descriptor_proto_get_json_name__or(field_descriptor, NULL);
		SET proto3_optional = _pb_field_descriptor_proto_get_proto3_optional(field_descriptor);
		SET oneof_index = _pb_field_descriptor_proto_get_oneof_index__or(field_descriptor, NULL);

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

		CASE field_type
		WHEN 10 THEN -- TYPE_GROUP (unsupported)
			SET message_text = CONCAT('_pb_wire_json_to_json: unsupported field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;

		WHEN 11 THEN -- TYPE_MESSAGE
			IF is_map THEN
				-- Handle map fields
				SET elements = pb_wire_json_get_repeated_message_field_as_json_array(wire_json, field_number);
				SET element_count = COALESCE(JSON_LENGTH(elements), 0);
				SET element_index = 0;
				SET field_json_value = JSON_OBJECT();

				-- Get map key/value field descriptors
				SET map_key_field = JSON_EXTRACT(map_entry_descriptor, '$."2"[0]'); -- first field (key) -- FIXME: don't assume specific index
				SET map_value_field = JSON_EXTRACT(map_entry_descriptor, '$."2"[1]'); -- second field (value) -- FIXME: don't assume specific index
				SET map_key_type = _pb_field_descriptor_proto_get_type(map_key_field);
				SET map_value_type = _pb_field_descriptor_proto_get_type(map_value_field);
				SET map_value_type_name = _pb_field_descriptor_proto_get_type_name__or(map_value_field, NULL);

				WHILE element_index < element_count DO
					SET element = pb_message_to_wire_json(FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(elements, CONCAT('$[', element_index, ']')))));
					CALL _pb_wire_json_get_primitive_field_as_number_json(element, 1, map_key_type, FALSE, FALSE, map_key);

					IF map_value_type = 11 THEN -- message
						CALL _pb_wire_json_to_number_json_proc(descriptor_set_json, map_value_type_name, pb_message_to_wire_json(pb_wire_json_get_message_field(element, 2, NULL)), map_value);
					ELSEIF map_value_type = 14 THEN -- enum
						SET map_value = CAST(pb_wire_json_get_enum_field(element, 2, 0) AS JSON);
					ELSE
						CALL _pb_wire_json_get_primitive_field_as_number_json(element, 2, map_value_type, FALSE, FALSE, map_value);
					END IF;

					IF JSON_TYPE(map_key) = 'STRING' THEN
						SET field_json_value = JSON_SET(field_json_value, CONCAT('$.', map_key), map_value);
					ELSE
						SET field_json_value = JSON_SET(field_json_value, CONCAT('$."', map_key, '"'), map_value);
					END IF;

					SET element_index = element_index + 1;
				END WHILE;

				IF element_count = 0 THEN
					SET field_json_value = NULL;
				END IF;

			ELSEIF is_repeated THEN
				-- Handle repeated message fields
				SET element_count = COALESCE(pb_wire_json_get_repeated_message_field_count(wire_json, field_number), 0);
				SET element_index = 0;
				SET field_json_value = JSON_ARRAY();

				WHILE element_index < element_count DO
					SET bytes_value = pb_wire_json_get_repeated_message_field_element(wire_json, field_number, element_index);
					CALL _pb_wire_json_to_number_json_proc(descriptor_set_json, field_type_name, pb_message_to_wire_json(bytes_value), nested_json_value);
					SET field_json_value = JSON_ARRAY_APPEND(field_json_value, '$', nested_json_value);
					SET element_index = element_index + 1;
				END WHILE;

				IF element_count = 0 THEN
					SET field_json_value = NULL;
				END IF;
			ELSE
				-- Handle singular message fields
				SET bytes_value = pb_wire_json_get_message_field(wire_json, field_number, NULL);
				CALL _pb_wire_json_to_number_json_proc(descriptor_set_json, field_type_name, pb_message_to_wire_json(bytes_value), field_json_value);
			END IF;

		WHEN 14 THEN -- TYPE_ENUM
			IF is_repeated THEN
				SET elements = pb_wire_json_get_repeated_enum_field_as_json_array(wire_json, field_number);
				SET element_count = COALESCE(JSON_LENGTH(elements), 0);
				SET element_index = 0;
				SET field_json_value = JSON_ARRAY();

				WHILE element_index < element_count DO
					SET element = JSON_EXTRACT(elements, CONCAT('$[', element_index, ']'));
					SET field_json_value = JSON_ARRAY_APPEND(field_json_value, '$', CAST(element AS JSON));
					SET element_index = element_index + 1;
				END WHILE;

				IF element_count = 0 THEN
					SET field_json_value = NULL;
				END IF;
			ELSE
				-- Handle singular enum fields
				SET field_enum_value = pb_wire_json_get_enum_field(wire_json, field_number, NULL);
				IF syntax = 'proto3' AND NOT has_field_presence AND field_enum_value = 0 THEN
					SET field_enum_value = NULL;
				END IF;

				SET field_json_value = NULL;
				IF field_enum_value IS NOT NULL THEN
					SET field_json_value = CAST(field_enum_value AS JSON);
				END IF;
			END IF;
		ELSE
			-- Handle primitive types using existing function
			CALL _pb_wire_json_get_primitive_field_as_number_json(wire_json, field_number, field_type, is_repeated, has_field_presence, field_json_value);
			IF is_repeated THEN
				IF JSON_LENGTH(field_json_value) = 0 THEN
					SET field_json_value = NULL;
				END IF;
			ELSE
				IF NOT has_field_presence THEN
					IF syntax = 'proto3' AND _pb_number_json_is_proto3_default_value(field_type, field_json_value) THEN
						SET field_json_value = NULL;
					END IF;
					-- emit_default_values is FALSE, so we never emit default values
				END IF;
			END IF;
		END CASE;

		-- Add field to result if it has a value
		IF field_json_value IS NOT NULL THEN
			SET json_field_name = CAST(field_number AS CHAR);

			IF oneof_index IS NOT NULL AND NOT proto3_optional THEN
				-- Handle oneof fields
				SET elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
				SET oneof_priority = JSON_EXTRACT(elements, CONCAT('$[', JSON_LENGTH(elements)-1, '].i'));
				SET oneof_priority_prev = JSON_EXTRACT(oneofs, CONCAT('$."', oneof_index, '".i'));

				IF oneof_priority_prev IS NULL OR oneof_priority_prev < oneof_priority THEN
					SET oneofs = JSON_SET(oneofs, CONCAT('$."', oneof_index, '"'), JSON_OBJECT('i', oneof_priority, 'v', JSON_OBJECT(json_field_name, field_json_value)));
				END IF;
			ELSE
				-- For number JSON format, field names are numeric and need to be quoted in JSON paths
				SET result = JSON_SET(result, CONCAT('$."', json_field_name, '"'), field_json_value);
			END IF;
		END IF;

		SET field_index = field_index + 1;
	END WHILE;

	-- Add oneof fields to result
	SET elements = JSON_EXTRACT(oneofs, '$.*.v');
	SET element_count = JSON_LENGTH(elements);
	SET element_index = 0;

	WHILE element_index < element_count DO
		SET field_json_value = JSON_EXTRACT(elements, CONCAT('$[', element_index, ']'));
		SET result = JSON_MERGE(result, field_json_value);
		SET element_index = element_index + 1;
	END WHILE;
END $$

-- Private function interface for protonumberjson format
DROP FUNCTION IF EXISTS _pb_message_to_number_json $$
CREATE FUNCTION _pb_message_to_number_json(descriptor_set_json JSON, type_name TEXT, message LONGBLOB, unmarshal_options JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	DECLARE message_text TEXT;

	-- Validate type name starts with dot
	IF type_name NOT LIKE '.%' THEN
		SET message_text = CONCAT('_pb_message_to_number_json: type name `', type_name, '` must start with a dot');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF message IS NULL THEN
		RETURN NULL;
	END IF;

	-- Embedded logic from _pb_message_to_number_json procedure
	CALL _pb_wire_json_to_number_json_proc(descriptor_set_json, type_name, pb_message_to_wire_json(message), result);
	RETURN result;
END $$

-- Private function interface for wire_json input with number JSON format
DROP FUNCTION IF EXISTS _pb_wire_json_to_number_json $$
CREATE FUNCTION _pb_wire_json_to_number_json(descriptor_set_json JSON, type_name TEXT, wire_json JSON, unmarshal_options JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	DECLARE message_text TEXT;

	-- Validate type name starts with dot
	IF type_name NOT LIKE '.%' THEN
		SET message_text = CONCAT('_pb_wire_json_to_number_json: type name `', type_name, '` must start with a dot');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF wire_json IS NULL THEN
		RETURN NULL;
	END IF;

	-- Embedded logic from _pb_wire_json_to_json
	CALL _pb_wire_json_to_number_json_proc(descriptor_set_json, type_name, wire_json, result);
	RETURN result;
END $$
