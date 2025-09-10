DELIMITER $$

-- Helper function to get message descriptor from descriptor set JSON
DROP FUNCTION IF EXISTS _pb_descriptor_set_get_message_descriptor $$
CREATE FUNCTION _pb_descriptor_set_get_message_descriptor(descriptor_set_json JSON, type_name TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE message_type_index JSON;
	DECLARE type_paths JSON;
	DECLARE file_path TEXT;
	DECLARE type_path TEXT;

	-- Get message type index (field 2 from DescriptorSet)
	SET message_type_index = JSON_EXTRACT(descriptor_set_json, '$.\"2\"');

	-- Get paths for the type
	SET type_paths = JSON_EXTRACT(message_type_index, CONCAT('$."', type_name, '"'));

	IF type_paths IS NULL THEN
		RETURN NULL;
	END IF;

	-- Extract file path and type path from MessageTypeIndex message
	SET file_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$.\"1\"'));
	SET type_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$.\"2\"'));

	-- Get FileDescriptorSet (field "1" of DescriptorSet), then apply type_path directly
	RETURN JSON_EXTRACT(JSON_EXTRACT(descriptor_set_json, '$.\"1\"'), type_path);
END $$

-- Helper function to get enum descriptor from descriptor set JSON
DROP FUNCTION IF EXISTS _pb_descriptor_set_get_enum_descriptor $$
CREATE FUNCTION _pb_descriptor_set_get_enum_descriptor(descriptor_set_json JSON, type_name TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE enum_type_index JSON;
	DECLARE type_paths JSON;
	DECLARE file_path TEXT;
	DECLARE type_path TEXT;

	-- Get enum type index (field 3 from DescriptorSet)
	SET enum_type_index = JSON_EXTRACT(descriptor_set_json, '$.\"3\"');

	-- Get paths for the type
	SET type_paths = JSON_EXTRACT(enum_type_index, CONCAT('$."', type_name, '"'));

	IF type_paths IS NULL THEN
		RETURN NULL;
	END IF;

	-- Extract file path and type path from EnumTypeIndex message
	SET file_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$.\"1\"'));
	SET type_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$.\"2\"'));

	-- Get FileDescriptorSet (field "1" of DescriptorSet), then apply type_path directly
	RETURN JSON_EXTRACT(JSON_EXTRACT(descriptor_set_json, '$.\"1\"'), type_path);
END $$


-- Helper function to get file descriptor for a type
DROP FUNCTION IF EXISTS _pb_descriptor_set_get_file_descriptor $$
CREATE FUNCTION _pb_descriptor_set_get_file_descriptor(descriptor_set_json JSON, type_name TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE type_index JSON;
	DECLARE type_paths JSON;
	DECLARE file_path TEXT;

	-- Get type index (field 2 from DescriptorSet)
	SET type_index = JSON_EXTRACT(descriptor_set_json, '$.\"2\"');

	-- Get paths for the type
	SET type_paths = JSON_EXTRACT(type_index, CONCAT('$."', type_name, '"'));

	IF type_paths IS NULL THEN
		RETURN NULL;
	END IF;

	-- Extract file path from MessageTypeIndex message (field "1")
	SET file_path = JSON_UNQUOTE(JSON_EXTRACT(type_paths, '$.\"1\"'));

	-- Get FileDescriptorSet (field "1" of DescriptorSet), then apply file_path
	RETURN JSON_EXTRACT(JSON_EXTRACT(descriptor_set_json, '$.\"1\"'), file_path);
END $$

-- Helper function to build fully-qualified type name
DROP FUNCTION IF EXISTS _pb_build_type_name $$
CREATE FUNCTION _pb_build_type_name(package_name TEXT, type_name TEXT) RETURNS TEXT DETERMINISTIC
BEGIN
	IF package_name = '' OR package_name IS NULL THEN
		RETURN CONCAT('.', type_name);
	ELSE
		RETURN CONCAT('.', package_name, '.', type_name);
	END IF;
END $$

-- Helper procedure to process nested types recursively
DROP PROCEDURE IF EXISTS _pb_build_nested_types $$
CREATE PROCEDURE _pb_build_nested_types(
	IN message_descriptor JSON,
	IN parent_name TEXT,
	IN parent_path TEXT,
	IN file_path TEXT,
	INOUT message_type_index JSON,
	INOUT enum_type_index JSON
)
proc: BEGIN
	DECLARE nested_messages JSON;
	DECLARE nested_enums JSON;
	DECLARE nested_msg_count INT DEFAULT 0;
	DECLARE nested_enum_count INT DEFAULT 0;
	DECLARE nested_msg_index INT DEFAULT 0;
	DECLARE nested_enum_index INT DEFAULT 0;
	DECLARE nested_msg_descriptor JSON;
	DECLARE nested_enum_descriptor JSON;
	DECLARE nested_msg_name TEXT;
	DECLARE nested_enum_name TEXT;
	DECLARE nested_msg_path TEXT;
	DECLARE nested_enum_path TEXT;
	DECLARE nested_type_name TEXT;
	DECLARE type_entry JSON;

	-- Enum indexing variables
	DECLARE enum_name_index JSON;
	DECLARE enum_number_index JSON;
	DECLARE enum_values_array JSON;
	DECLARE enum_value_count INT DEFAULT 0;
	DECLARE enum_value_idx INT DEFAULT 0;
	DECLARE enum_value_desc JSON;
	DECLARE enum_value_name TEXT;
	DECLARE enum_value_number INT;

	-- Nested field indexing variables
	DECLARE nested_field_name_index JSON;
	DECLARE nested_field_number_index JSON;
	DECLARE nested_field_array JSON;
	DECLARE nested_field_count INT DEFAULT 0;
	DECLARE nested_field_idx INT DEFAULT 0;
	DECLARE nested_field_desc JSON;
	DECLARE nested_field_name_val TEXT;
	DECLARE nested_field_number_val INT;

	-- Process nested message types (field 3 in DescriptorProto)
	SET nested_messages = JSON_EXTRACT(message_descriptor, '$."3"');

	IF nested_messages IS NOT NULL THEN
		SET nested_msg_count = JSON_LENGTH(nested_messages);
		SET nested_msg_index = 0;

		WHILE nested_msg_index < nested_msg_count DO
			SET nested_msg_descriptor = JSON_EXTRACT(nested_messages, CONCAT('$[', nested_msg_index, ']'));
			SET nested_msg_name = JSON_UNQUOTE(JSON_EXTRACT(nested_msg_descriptor, '$."1"')); -- name field
			SET nested_msg_path = CONCAT(parent_path, '."3"[', nested_msg_index, ']');
			SET nested_type_name = CONCAT(parent_name, '.', nested_msg_name);

			-- Build field indexes for nested message
			SET nested_field_name_index = JSON_OBJECT();
			SET nested_field_number_index = JSON_OBJECT();
			SET nested_field_array = JSON_EXTRACT(nested_msg_descriptor, '$."2"'); -- field array
			IF nested_field_array IS NOT NULL THEN
				SET nested_field_count = JSON_LENGTH(nested_field_array);
				SET nested_field_idx = 0;
				WHILE nested_field_idx < nested_field_count DO
					SET nested_field_desc = JSON_EXTRACT(nested_field_array, CONCAT('$[', nested_field_idx, ']'));
					SET nested_field_name_val = JSON_UNQUOTE(JSON_EXTRACT(nested_field_desc, '$."1"')); -- name
					SET nested_field_number_val = JSON_EXTRACT(nested_field_desc, '$."3"'); -- number
					IF nested_field_name_val IS NOT NULL THEN
						SET nested_field_name_index = JSON_SET(nested_field_name_index, CONCAT('$."', nested_field_name_val, '"'), nested_field_idx);
					END IF;
					IF nested_field_number_val IS NOT NULL THEN
						SET nested_field_number_index = JSON_SET(nested_field_number_index, CONCAT('$."', nested_field_number_val, '"'), nested_field_idx);
					END IF;
					SET nested_field_idx = nested_field_idx + 1;
				END WHILE;
			END IF;

			-- Add to message type index: MessageTypeIndex format
			SET type_entry = JSON_OBJECT(
				'1', file_path,
				'2', nested_msg_path
			);
			-- Only include field indexes if they're non-empty
			IF JSON_LENGTH(JSON_KEYS(nested_field_name_index)) > 0 THEN
				SET type_entry = JSON_SET(type_entry, '$.\"3\"', nested_field_name_index);
			END IF;
			IF JSON_LENGTH(JSON_KEYS(nested_field_number_index)) > 0 THEN
				SET type_entry = JSON_SET(type_entry, '$.\"4\"', nested_field_number_index);
			END IF;
			SET message_type_index = JSON_SET(message_type_index, CONCAT('$."', nested_type_name, '"'), type_entry);

			-- Recursively process further nested types
			CALL _pb_build_nested_types(nested_msg_descriptor, nested_type_name, nested_msg_path, file_path, message_type_index, enum_type_index);

			SET nested_msg_index = nested_msg_index + 1;
		END WHILE;
	END IF;

	-- Process nested enum types (field 4 in DescriptorProto)
	SET nested_enums = JSON_EXTRACT(message_descriptor, '$."4"');

	IF nested_enums IS NOT NULL THEN
		SET nested_enum_count = JSON_LENGTH(nested_enums);
		SET nested_enum_index = 0;

		WHILE nested_enum_index < nested_enum_count DO
			SET nested_enum_descriptor = JSON_EXTRACT(nested_enums, CONCAT('$[', nested_enum_index, ']'));
			SET nested_enum_name = JSON_UNQUOTE(JSON_EXTRACT(nested_enum_descriptor, '$."1"')); -- name field
			SET nested_enum_path = CONCAT(parent_path, '."4"[', nested_enum_index, ']');
			SET nested_type_name = CONCAT(parent_name, '.', nested_enum_name);

			-- Build enum value indexes
			SET enum_name_index = JSON_OBJECT();
			SET enum_number_index = JSON_OBJECT();
			SET enum_values_array = JSON_EXTRACT(nested_enum_descriptor, '$."2"'); -- value array
			IF enum_values_array IS NOT NULL THEN
				SET enum_value_count = JSON_LENGTH(enum_values_array);
				SET enum_value_idx = 0;
				WHILE enum_value_idx < enum_value_count DO
					SET enum_value_desc = JSON_EXTRACT(enum_values_array, CONCAT('$[', enum_value_idx, ']'));
					SET enum_value_name = JSON_UNQUOTE(JSON_EXTRACT(enum_value_desc, '$."1"')); -- name
					SET enum_value_number = JSON_EXTRACT(enum_value_desc, '$."2"'); -- number
					IF enum_value_name IS NOT NULL THEN
						SET enum_name_index = JSON_SET(enum_name_index, CONCAT('$."', enum_value_name, '"'), enum_value_idx);
					END IF;
					IF enum_value_number IS NOT NULL THEN
						SET enum_number_index = JSON_SET(enum_number_index, CONCAT('$."', enum_value_number, '"'), enum_value_idx);
					END IF;
					SET enum_value_idx = enum_value_idx + 1;
				END WHILE;
			END IF;

			-- Add to enum type index: EnumTypeIndex format
			SET type_entry = JSON_OBJECT(
				'1', file_path,
				'2', nested_enum_path,
				'3', enum_name_index,
				'4', enum_number_index
			);
			SET enum_type_index = JSON_SET(enum_type_index, CONCAT('$."', nested_type_name, '"'), type_entry);

			SET nested_enum_index = nested_enum_index + 1;
		END WHILE;
	END IF;
END $$

-- Public procedure to generate type indexes from FileDescriptorSet in protonumberjson format
DROP PROCEDURE IF EXISTS _pb_build_type_indexes_from_descriptor_set $$
CREATE PROCEDURE _pb_build_type_indexes_from_descriptor_set(
	IN file_descriptor_set_json JSON,
	OUT message_type_index JSON,
	OUT enum_type_index JSON
)
proc: BEGIN
	DECLARE files JSON;
	DECLARE file_count INT DEFAULT 0;
	DECLARE file_index INT DEFAULT 0;
	DECLARE file_descriptor JSON;
	DECLARE file_package TEXT;
	DECLARE file_path TEXT;
	DECLARE message_types JSON;
	DECLARE enum_types JSON;
	DECLARE msg_count INT DEFAULT 0;
	DECLARE enum_count INT DEFAULT 0;
	DECLARE msg_index INT DEFAULT 0;
	DECLARE enum_index INT DEFAULT 0;
	DECLARE message_descriptor JSON;
	DECLARE enum_descriptor JSON;
	DECLARE message_name TEXT;
	DECLARE enum_name TEXT;
	DECLARE message_path TEXT;
	DECLARE enum_path TEXT;
	DECLARE full_type_name TEXT;
	DECLARE type_entry JSON;

	-- Field indexing variables
	DECLARE field_name_index JSON;
	DECLARE field_number_index JSON;
	DECLARE field_array JSON;
	DECLARE field_count INT DEFAULT 0;
	DECLARE field_idx INT DEFAULT 0;
	DECLARE field_desc JSON;
	DECLARE field_name_val TEXT;
	DECLARE field_number_val INT;

	-- Enum indexing variables
	DECLARE enum_name_index JSON;
	DECLARE enum_number_index JSON;
	DECLARE enum_values_array JSON;
	DECLARE enum_value_count INT DEFAULT 0;
	DECLARE enum_value_idx INT DEFAULT 0;
	DECLARE enum_value_desc JSON;
	DECLARE enum_value_name TEXT;
	DECLARE enum_value_number INT;

	-- Initialize indexes
	SET message_type_index = JSON_OBJECT();
	SET enum_type_index = JSON_OBJECT();

	-- Extract files array (field 1 in FileDescriptorSet)
	SET files = JSON_EXTRACT(file_descriptor_set_json, '$."1"');

	IF files IS NULL THEN
		-- Return empty indexes (OUT parameters already initialized)
		LEAVE proc;
	END IF;

	SET file_count = JSON_LENGTH(files);
	SET file_index = 0;

	-- Iterate through each file
	WHILE file_index < file_count DO
		SET file_descriptor = JSON_EXTRACT(files, CONCAT('$[', file_index, ']'));
		SET file_package = COALESCE(JSON_UNQUOTE(JSON_EXTRACT(file_descriptor, '$."2"')), ''); -- package field
		SET file_path = CONCAT('$."1"[', file_index, ']');

		-- Process message types (field 4 in FileDescriptorProto)
		SET message_types = JSON_EXTRACT(file_descriptor, '$."4"');

		IF message_types IS NOT NULL THEN
			SET msg_count = JSON_LENGTH(message_types);
			SET msg_index = 0;

			WHILE msg_index < msg_count DO
				SET message_descriptor = JSON_EXTRACT(message_types, CONCAT('$[', msg_index, ']'));
				SET message_name = JSON_UNQUOTE(JSON_EXTRACT(message_descriptor, '$."1"')); -- name field
				SET message_path = CONCAT(file_path, '."4"[', msg_index, ']');
				SET full_type_name = _pb_build_type_name(file_package, message_name);

				-- Build field indexes for message
				SET field_name_index = JSON_OBJECT();
				SET field_number_index = JSON_OBJECT();
				SET field_array = JSON_EXTRACT(message_descriptor, '$."2"'); -- field array
				IF field_array IS NOT NULL THEN
					SET field_count = JSON_LENGTH(field_array);
					SET field_idx = 0;
					WHILE field_idx < field_count DO
						SET field_desc = JSON_EXTRACT(field_array, CONCAT('$[', field_idx, ']'));
						SET field_name_val = JSON_UNQUOTE(JSON_EXTRACT(field_desc, '$."1"')); -- name
						SET field_number_val = JSON_EXTRACT(field_desc, '$."3"'); -- number
						IF field_name_val IS NOT NULL THEN
							SET field_name_index = JSON_SET(field_name_index, CONCAT('$."', field_name_val, '"'), field_idx);
						END IF;
						IF field_number_val IS NOT NULL THEN
							SET field_number_index = JSON_SET(field_number_index, CONCAT('$."', field_number_val, '"'), field_idx);
						END IF;
						SET field_idx = field_idx + 1;
					END WHILE;
				END IF;

				-- Add to message type index: MessageTypeIndex format
				SET type_entry = JSON_OBJECT(
					'1', file_path,
					'2', message_path
				);
				-- Only include field indexes if they're non-empty
				IF JSON_LENGTH(JSON_KEYS(field_name_index)) > 0 THEN
					SET type_entry = JSON_SET(type_entry, '$.\"3\"', field_name_index);
				END IF;
				IF JSON_LENGTH(JSON_KEYS(field_number_index)) > 0 THEN
					SET type_entry = JSON_SET(type_entry, '$.\"4\"', field_number_index);
				END IF;
				SET message_type_index = JSON_SET(message_type_index, CONCAT('$."', full_type_name, '"'), type_entry);

				-- Process nested types recursively
				CALL _pb_build_nested_types(message_descriptor, full_type_name, message_path, file_path, message_type_index, enum_type_index);

				SET msg_index = msg_index + 1;
			END WHILE;
		END IF;

		-- Process enum types (field 5 in FileDescriptorProto)
		SET enum_types = JSON_EXTRACT(file_descriptor, '$."5"');

		IF enum_types IS NOT NULL THEN
			SET enum_count = JSON_LENGTH(enum_types);
			SET enum_index = 0;

			WHILE enum_index < enum_count DO
				SET enum_descriptor = JSON_EXTRACT(enum_types, CONCAT('$[', enum_index, ']'));
				SET enum_name = JSON_UNQUOTE(JSON_EXTRACT(enum_descriptor, '$."1"')); -- name field
				SET enum_path = CONCAT(file_path, '."5"[', enum_index, ']');
				SET full_type_name = _pb_build_type_name(file_package, enum_name);

				-- Build enum value indexes
				SET enum_name_index = JSON_OBJECT();
				SET enum_number_index = JSON_OBJECT();
				SET enum_values_array = JSON_EXTRACT(enum_descriptor, '$."2"'); -- value array
				IF enum_values_array IS NOT NULL THEN
					SET enum_value_count = JSON_LENGTH(enum_values_array);
					SET enum_value_idx = 0;
					WHILE enum_value_idx < enum_value_count DO
						SET enum_value_desc = JSON_EXTRACT(enum_values_array, CONCAT('$[', enum_value_idx, ']'));
						SET enum_value_name = JSON_UNQUOTE(JSON_EXTRACT(enum_value_desc, '$."1"')); -- name
						SET enum_value_number = JSON_EXTRACT(enum_value_desc, '$."2"'); -- number
						IF enum_value_name IS NOT NULL THEN
							SET enum_name_index = JSON_SET(enum_name_index, CONCAT('$."', enum_value_name, '"'), enum_value_idx);
						END IF;
						IF enum_value_number IS NOT NULL THEN
							SET enum_number_index = JSON_SET(enum_number_index, CONCAT('$."', enum_value_number, '"'), enum_value_idx);
						END IF;
						SET enum_value_idx = enum_value_idx + 1;
					END WHILE;
				END IF;

				-- Add to enum type index: EnumTypeIndex format
				SET type_entry = JSON_OBJECT(
					'1', file_path,
					'2', enum_path,
					'3', enum_name_index,
					'4', enum_number_index
				);
				SET enum_type_index = JSON_SET(enum_type_index, CONCAT('$."', full_type_name, '"'), type_entry);

				SET enum_index = enum_index + 1;
			END WHILE;
		END IF;

		SET file_index = file_index + 1;
	END WHILE;
END $$

-- Public function to convert FileDescriptorSet LONGBLOB to descriptor set JSON
-- Returns a DescriptorSet message in protonumberjson format
DROP FUNCTION IF EXISTS pb_descriptor_set_build $$
CREATE FUNCTION pb_descriptor_set_build(file_descriptor_set_blob LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE file_descriptor_set_number_json JSON;
	DECLARE message_type_index JSON;
	DECLARE enum_type_index JSON;
	DECLARE result JSON;

	-- Convert FileDescriptorSet LONGBLOB to protonumberjson format
	SET file_descriptor_set_number_json = _pb_message_to_number_json(
		_pb_descriptor_proto(),
		'.google.protobuf.FileDescriptorSet',
		file_descriptor_set_blob,
		NULL  -- unmarshal_options: use default behavior
	);

	-- Build type indexes from the FileDescriptorSet
	CALL _pb_build_type_indexes_from_descriptor_set(
		file_descriptor_set_number_json,
		message_type_index,
		enum_type_index
	);

	-- Return DescriptorSet message: {"1": fileDescriptorSet, "2": messageTypeIndex, "3": enumTypeIndex}
	SET result = JSON_OBJECT(
		'1', file_descriptor_set_number_json,
		'2', message_type_index
	);
	-- Only include enumTypeIndex if it's non-empty
	IF JSON_LENGTH(JSON_KEYS(enum_type_index)) > 0 THEN
		SET result = JSON_SET(result, '$.\"3\"', enum_type_index);
	END IF;

	RETURN result;
END $$
