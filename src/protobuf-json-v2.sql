DELIMITER $$

-- Helper function to get message descriptor from descriptor set JSON
DROP FUNCTION IF EXISTS _pb_get_message_descriptor $$
CREATE FUNCTION _pb_get_message_descriptor(descriptor_set_json JSON, type_name TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE type_index JSON;
	DECLARE type_paths JSON;
	DECLARE kind INT;
	DECLARE file_path TEXT;
	DECLARE type_path TEXT;

	-- Get type index (element 2)
	SET type_index = JSON_EXTRACT(descriptor_set_json, '$[2]');

	-- Get paths for the type
	SET type_paths = JSON_EXTRACT(type_index, CONCAT('$."', type_name, '"'));

	IF type_paths IS NULL THEN
		RETURN NULL;
	END IF;

	-- Extract kind, file path and type path
	SET kind = JSON_EXTRACT(type_paths, '$[0]');
	SET file_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$[1]'));
	SET type_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$[2]'));

	-- Verify this is a message type (kind = 11)
	IF kind <> 11 THEN
		RETURN NULL;
	END IF;

	-- Return the message descriptor
	RETURN JSON_EXTRACT(descriptor_set_json, type_path);
END $$

-- Helper function to get enum descriptor from descriptor set JSON
DROP FUNCTION IF EXISTS _pb_get_enum_descriptor $$
CREATE FUNCTION _pb_get_enum_descriptor(descriptor_set_json JSON, type_name TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE type_index JSON;
	DECLARE type_paths JSON;
	DECLARE kind INT;
	DECLARE file_path TEXT;
	DECLARE type_path TEXT;

	-- Get type index (element 2)
	SET type_index = JSON_EXTRACT(descriptor_set_json, '$[2]');

	-- Get paths for the type
	SET type_paths = JSON_EXTRACT(type_index, CONCAT('$."', type_name, '"'));

	IF type_paths IS NULL THEN
		RETURN NULL;
	END IF;

	-- Extract kind, file path and type path
	SET kind = JSON_EXTRACT(type_paths, '$[0]');
	SET file_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$[1]'));
	SET type_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$[2]'));

	-- Verify this is an enum type (kind = 14)
	IF kind <> 14 THEN
		RETURN NULL;
	END IF;

	-- Return the enum descriptor
	RETURN JSON_EXTRACT(descriptor_set_json, type_path);
END $$

-- Helper procedure to convert enum value to JSON using descriptor set
DROP PROCEDURE IF EXISTS _pb_enum_to_json $$
CREATE PROCEDURE _pb_enum_to_json(IN descriptor_set_json JSON, IN full_type_name TEXT, IN enum_value_number INT, OUT result JSON)
proc: BEGIN
	DECLARE enum_descriptor JSON;
	DECLARE enum_values JSON;
	DECLARE enum_value JSON;
	DECLARE enum_count INT;
	DECLARE enum_index INT;
	DECLARE current_number INT;
	DECLARE current_name TEXT;

	SET enum_descriptor = _pb_get_enum_descriptor(descriptor_set_json, full_type_name);

	IF enum_descriptor IS NULL THEN
		SET result = NULL;
		LEAVE proc;
	END IF;

	-- Get enum values array (field 2 in EnumDescriptorProto)
	SET enum_values = JSON_EXTRACT(enum_descriptor, '$."2"');

	IF enum_values IS NULL THEN
		SET result = NULL;
		LEAVE proc;
	END IF;

	SET enum_count = JSON_LENGTH(enum_values);
	SET enum_index = 0;

	-- Find enum value by number
	WHILE enum_index < enum_count DO
		SET enum_value = JSON_EXTRACT(enum_values, CONCAT('$[', enum_index, ']'));
		SET current_number = JSON_EXTRACT(enum_value, '$."2"'); -- number field

		IF current_number = enum_value_number THEN
			SET current_name = JSON_UNQUOTE(JSON_EXTRACT(enum_value, '$."1"')); -- name field
			SET result = JSON_QUOTE(current_name);
			LEAVE proc;
		END IF;

		SET enum_index = enum_index + 1;
	END WHILE;

	-- If not found, return the numeric value (Proto3 behavior for unknown enum values)
	IF result IS NULL THEN
		SET result = CAST(enum_value_number AS JSON);
	END IF;
END $$

-- Helper function to get file descriptor for a type
DROP FUNCTION IF EXISTS _pb_get_file_descriptor $$
CREATE FUNCTION _pb_get_file_descriptor(descriptor_set_json JSON, type_name TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE type_index JSON;
	DECLARE type_paths JSON;
	DECLARE file_path TEXT;

	-- Get type index (element 2)
	SET type_index = JSON_EXTRACT(descriptor_set_json, '$[2]');

	-- Get paths for the type
	SET type_paths = JSON_EXTRACT(type_index, CONCAT('$."', type_name, '"'));

	IF type_paths IS NULL THEN
		RETURN NULL;
	END IF;

	-- Extract file path (now at index 1)
	SET file_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$[1]'));

	-- Return the file descriptor
	RETURN JSON_EXTRACT(descriptor_set_json, file_path);
END $$

-- Helper function to get the appropriate descriptor set for Google well-known types
DROP FUNCTION IF EXISTS _pb_get_wkt_descriptor_set $$
CREATE FUNCTION _pb_get_wkt_descriptor_set(full_type_name TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	-- For Google well-known types, use built-in descriptor functions
	CASE
	WHEN full_type_name IN ('.google.protobuf.Struct', '.google.protobuf.Value', '.google.protobuf.ListValue', '.google.protobuf.NullValue') THEN
		RETURN _pb_wkt_struct_proto();
	WHEN full_type_name = '.google.protobuf.FieldMask' THEN
		RETURN _pb_wkt_field_mask_proto();
	WHEN full_type_name IN ('.google.protobuf.DoubleValue', '.google.protobuf.FloatValue', '.google.protobuf.Int64Value', '.google.protobuf.UInt64Value', '.google.protobuf.Int32Value', '.google.protobuf.UInt32Value', '.google.protobuf.BoolValue', '.google.protobuf.StringValue', '.google.protobuf.BytesValue') THEN
		RETURN _pb_wkt_wrappers_proto();
	WHEN full_type_name = '.google.protobuf.Empty' THEN
		RETURN _pb_wkt_empty_proto();
	WHEN full_type_name = '.google.protobuf.Any' THEN
		RETURN _pb_wkt_any_proto();
	WHEN full_type_name = '.google.protobuf.Timestamp' THEN
		RETURN _pb_wkt_timestamp_proto();
	WHEN full_type_name = '.google.protobuf.Duration' THEN
		RETURN _pb_wkt_duration_proto();
	ELSE
		-- Return NULL for types that don't match or should use regular WKT handling
		RETURN NULL;
	END CASE;
END $$

-- Main procedure for converting protobuf message to JSON using descriptor set
DROP PROCEDURE IF EXISTS _pb_wire_json_to_json_proc $$
CREATE PROCEDURE _pb_wire_json_to_json_proc(IN descriptor_set_json JSON, IN full_type_name TEXT, IN wire_json JSON, IN as_number_json BOOLEAN, IN emit_default_values BOOLEAN, OUT result JSON)
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

	-- Handle well-known types first (only for regular JSON, not number JSON)
	IF full_type_name LIKE '.google.protobuf.%' THEN
		IF as_number_json THEN -- For ProtoNumberJSON, no special WKT handling is performed.
			SET wkt_descriptor_set = _pb_get_wkt_descriptor_set(full_type_name);
			IF wkt_descriptor_set IS NOT NULL THEN
				SET descriptor_set_json = wkt_descriptor_set;
			END IF;
		ELSE -- For ProtoJSON, we use special WKT decoders.
			SET result = _pb_wire_json_decode_wkt_as_json(wire_json, full_type_name);
			IF result IS NOT NULL THEN
				LEAVE proc;
			END IF;
		END IF;
	END IF;

	-- Get message descriptor
	SET message_descriptor = _pb_get_message_descriptor(descriptor_set_json, full_type_name);

	IF message_descriptor IS NULL THEN
		SET message_text = CONCAT('_pb_wire_json_to_json: message type `', full_type_name, '` not found in descriptor set');
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get file descriptor to determine syntax
	SET file_descriptor = _pb_get_file_descriptor(descriptor_set_json, full_type_name);
	SET syntax = JSON_UNQUOTE(JSON_EXTRACT(file_descriptor, '$."12"')); -- syntax field
	IF syntax IS NULL THEN
		SET syntax = 'proto2'; -- default
	END IF;

	SET result = JSON_OBJECT();
	SET oneofs = JSON_OBJECT();

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
				SET map_entry_descriptor = _pb_get_message_descriptor(descriptor_set_json, field_type_name);
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
					SET map_key_field = JSON_EXTRACT(map_entry_descriptor, '$."2"[0]'); -- first field (key)
					SET map_value_field = JSON_EXTRACT(map_entry_descriptor, '$."2"[1]'); -- second field (value)
					SET map_key_type = JSON_EXTRACT(map_key_field, '$."5"');
					SET map_value_type = JSON_EXTRACT(map_value_field, '$."5"');
					SET map_value_type_name = JSON_UNQUOTE(JSON_EXTRACT(map_value_field, '$."6"'));

					WHILE element_index < element_count DO
						SET element = pb_message_to_wire_json(FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(elements, CONCAT('$[', element_index, ']')))));
						CALL _pb_wire_json_get_primitive_field_as_json(element, 1, map_key_type, FALSE, FALSE, as_number_json, map_key);

						IF map_value_type = 11 THEN -- message
							CALL _pb_message_to_json(descriptor_set_json, map_value_type_name, pb_wire_json_get_message_field(element, 2, NULL), as_number_json, emit_default_values, map_value);
						ELSEIF map_value_type = 14 THEN -- enum
							IF as_number_json THEN
								SET map_value = CAST(pb_wire_json_get_enum_field(element, 2, 0) AS JSON);
							ELSE
								CALL _pb_enum_to_json(descriptor_set_json, map_value_type_name, pb_wire_json_get_enum_field(element, 2, 0), map_value);
							END IF;
						ELSE
							CALL _pb_wire_json_get_primitive_field_as_json(element, 2, map_value_type, FALSE, FALSE, as_number_json, map_value);
						END IF;

						IF JSON_TYPE(map_key) = 'STRING' THEN
							SET field_json_value = JSON_SET(field_json_value, CONCAT('$.', map_key), map_value);
						ELSE
							SET field_json_value = JSON_SET(field_json_value, CONCAT('$."', map_key, '"'), map_value);
						END IF;

						SET element_index = element_index + 1;
					END WHILE;

					IF NOT emit_default_values AND element_count = 0 THEN
						SET field_json_value = NULL;
					END IF;

				ELSEIF is_repeated THEN
					-- Handle repeated message fields
					SET element_count = COALESCE(pb_wire_json_get_repeated_message_field_count(wire_json, field_number), 0);
					SET element_index = 0;
					SET field_json_value = JSON_ARRAY();

					WHILE element_index < element_count DO
						SET bytes_value = pb_wire_json_get_repeated_message_field_element(wire_json, field_number, element_index);
						CALL _pb_message_to_json(descriptor_set_json, field_type_name, bytes_value, as_number_json, emit_default_values, nested_json_value);
						SET field_json_value = JSON_ARRAY_APPEND(field_json_value, '$', nested_json_value);
						SET element_index = element_index + 1;
					END WHILE;

					IF NOT emit_default_values AND element_count = 0 THEN
						SET field_json_value = NULL;
					END IF;
				ELSE
					-- Handle singular message fields
					SET bytes_value = pb_wire_json_get_message_field(wire_json, field_number, NULL);
					CALL _pb_message_to_json(descriptor_set_json, field_type_name, bytes_value, as_number_json, emit_default_values, field_json_value);
				END IF;

			WHEN 14 THEN -- TYPE_ENUM
				IF is_repeated THEN
					SET elements = pb_wire_json_get_repeated_enum_field_as_json_array(wire_json, field_number);
					SET element_count = COALESCE(JSON_LENGTH(elements), 0);
					SET element_index = 0;
					SET field_json_value = JSON_ARRAY();

					WHILE element_index < element_count DO
						SET element = JSON_EXTRACT(elements, CONCAT('$[', element_index, ']'));
						IF as_number_json THEN
							SET field_json_value = JSON_ARRAY_APPEND(field_json_value, '$', CAST(element AS JSON));
						ELSE
							CALL _pb_enum_to_json(descriptor_set_json, field_type_name, element, nested_json_value);
							SET field_json_value = JSON_ARRAY_APPEND(field_json_value, '$', nested_json_value);
						END IF;
						SET element_index = element_index + 1;
					END WHILE;

					IF NOT emit_default_values AND element_count = 0 THEN
						SET field_json_value = NULL;
					END IF;
				ELSE
					-- Handle singular enum fields
					SET field_enum_value = pb_wire_json_get_enum_field(wire_json, field_number, NULL);
					IF syntax = 'proto3' AND NOT has_field_presence AND field_enum_value = 0 THEN
						SET field_enum_value = NULL;
					END IF;

					SET field_json_value = NULL;
					IF as_number_json THEN -- ProtoNumberJSON: always omit default values
						IF field_enum_value IS NOT NULL THEN
							SET field_json_value = CAST(field_enum_value AS JSON);
						END IF;
					ELSE -- ProtoJSON:
						IF field_enum_value IS NOT NULL OR (emit_default_values AND NOT has_field_presence) THEN
							IF field_enum_value IS NULL THEN
								SET field_enum_value = 0;
							END IF;
							CALL _pb_enum_to_json(descriptor_set_json, field_type_name, field_enum_value, field_json_value);
						END IF;
					END IF;
				END IF;
			ELSE
				-- Handle primitive types using existing function
				CALL _pb_wire_json_get_primitive_field_as_json(wire_json, field_number, field_type, is_repeated, has_field_presence, as_number_json, field_json_value);
				IF is_repeated THEN
					IF NOT emit_default_values AND JSON_LENGTH(field_json_value) = 0 THEN
						SET field_json_value = NULL;
					END IF;
				ELSE
					IF NOT has_field_presence THEN
						IF syntax = 'proto3' AND _pb_is_proto3_default_value(field_type, field_json_value) THEN
							SET field_json_value = NULL;
						END IF;
						IF emit_default_values AND field_json_value IS NULL THEN
							SET field_json_value = _pb_get_proto3_default_value(field_type, as_number_json);
						END IF;
					END IF;
				END IF;
			END CASE;

			-- Add field to result if it has a value
			IF field_json_value IS NOT NULL THEN
				IF as_number_json THEN
					SET json_field_name = CAST(field_number AS CHAR);
				ELSE
					SET json_field_name = IF(json_name IS NOT NULL, json_name, _pb_util_snake_to_camel(field_name));
				END IF;

				IF oneof_index IS NOT NULL AND NOT proto3_optional THEN
					-- Handle oneof fields
					SET elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
					SET oneof_priority = JSON_EXTRACT(elements, CONCAT('$[', JSON_LENGTH(elements)-1, '].i'));
					SET oneof_priority_prev = JSON_EXTRACT(oneofs, CONCAT('$."', oneof_index, '".i'));

					IF oneof_priority_prev IS NULL OR oneof_priority_prev < oneof_priority THEN
						SET oneofs = JSON_SET(oneofs, CONCAT('$."', oneof_index, '"'), JSON_OBJECT('i', oneof_priority, 'v', JSON_OBJECT(json_field_name, field_json_value)));
					END IF;
				ELSE
					-- Regular field
					IF as_number_json THEN
						-- For number JSON format, field names are numeric and need to be quoted in JSON paths
						SET result = JSON_SET(result, CONCAT('$."', json_field_name, '"'), field_json_value);
					ELSE
						SET result = JSON_SET(result, CONCAT('$.', json_field_name), field_json_value);
					END IF;
				END IF;
			END IF;

			SET field_index = field_index + 1;
		END WHILE;
	END IF;

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

-- Wrapper procedure that converts LONGBLOB to wire_json and delegates
DROP PROCEDURE IF EXISTS _pb_message_to_json $$
CREATE PROCEDURE _pb_message_to_json(IN descriptor_set_json JSON, IN full_type_name TEXT, IN message LONGBLOB, IN as_number_json BOOLEAN, IN emit_default_values BOOLEAN, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;

	-- Validate type name starts with dot
	IF full_type_name NOT LIKE '.%' THEN
		SET message_text = CONCAT('_pb_message_to_json: type name `', full_type_name, '` must start with a dot');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF message IS NULL THEN
		SET result = NULL;
	ELSE
		CALL _pb_wire_json_to_json_proc(descriptor_set_json, full_type_name, pb_message_to_wire_json(message), as_number_json, emit_default_values, result);
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_to_json $$
CREATE PROCEDURE _pb_wire_json_to_json(IN descriptor_set_json JSON, IN full_type_name TEXT, IN wire_json JSON, IN as_number_json BOOLEAN, IN emit_default_values BOOLEAN, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;

	-- Validate type name starts with dot
	IF full_type_name NOT LIKE '.%' THEN
		SET message_text = CONCAT('_pb_wire_json_to_json: type name `', full_type_name, '` must start with a dot');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF wire_json IS NULL THEN
		SET result = NULL;
	ELSE
		CALL _pb_wire_json_to_json_proc(descriptor_set_json, full_type_name, wire_json, as_number_json, emit_default_values, result);
	END IF;
END $$

-- Public function interface
DROP FUNCTION IF EXISTS pb_message_to_json $$
CREATE FUNCTION pb_message_to_json(descriptor_set_json JSON, type_name TEXT, message LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_to_json(descriptor_set_json, type_name, message, FALSE, TRUE, result);
	RETURN result;
END $$

-- Private function interface for protonumberjson format
DROP FUNCTION IF EXISTS _pb_message_to_number_json $$
CREATE FUNCTION _pb_message_to_number_json(descriptor_set_json JSON, type_name TEXT, message LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_to_json(descriptor_set_json, type_name, message, TRUE, FALSE, result);
	RETURN result;
END $$

-- Public function interface for wire_json input
DROP FUNCTION IF EXISTS pb_wire_json_to_json $$
CREATE FUNCTION pb_wire_json_to_json(descriptor_set_json JSON, type_name TEXT, wire_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_to_json(descriptor_set_json, type_name, wire_json, FALSE, TRUE, result);
	RETURN result;
END $$

-- Private function interface for wire_json input with number JSON format
DROP FUNCTION IF EXISTS _pb_wire_json_to_number_json $$
CREATE FUNCTION _pb_wire_json_to_number_json(descriptor_set_json JSON, type_name TEXT, wire_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_to_json(descriptor_set_json, type_name, wire_json, TRUE, FALSE, result);
	RETURN result;
END $$
