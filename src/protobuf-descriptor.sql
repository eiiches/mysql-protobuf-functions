DELIMITER $$

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
	INOUT type_index JSON
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
			
			-- Add to type index: [kind=11, file_path, type_path]
			SET type_entry = JSON_ARRAY(11, file_path, nested_msg_path);
			SET type_index = JSON_SET(type_index, CONCAT('$."', nested_type_name, '"'), type_entry);
			
			-- Recursively process further nested types
			CALL _pb_build_nested_types(nested_msg_descriptor, nested_type_name, nested_msg_path, file_path, type_index);
			
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
			
			-- Add to type index: [kind=14, file_path, type_path]
			SET type_entry = JSON_ARRAY(14, file_path, nested_enum_path);
			SET type_index = JSON_SET(type_index, CONCAT('$."', nested_type_name, '"'), type_entry);
			
			SET nested_enum_index = nested_enum_index + 1;
		END WHILE;
	END IF;
END $$

-- Public function to generate type index from FileDescriptorSet in protonumberjson format
DROP FUNCTION IF EXISTS _pb_build_type_index_from_descriptor_set $$
CREATE FUNCTION _pb_build_type_index_from_descriptor_set(file_descriptor_set_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE type_index JSON DEFAULT JSON_OBJECT();
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

	-- Extract files array (field 1 in FileDescriptorSet)
	SET files = JSON_EXTRACT(file_descriptor_set_json, '$."1"');
	
	IF files IS NULL THEN
		RETURN type_index;
	END IF;
	
	SET file_count = JSON_LENGTH(files);
	SET file_index = 0;
	
	-- Iterate through each file
	WHILE file_index < file_count DO
		SET file_descriptor = JSON_EXTRACT(files, CONCAT('$[', file_index, ']'));
		SET file_package = COALESCE(JSON_UNQUOTE(JSON_EXTRACT(file_descriptor, '$."2"')), ''); -- package field
		SET file_path = CONCAT('$[1]."1"[', file_index, ']');
		
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
				
				-- Add to type index: [kind=11, file_path, type_path]
				SET type_entry = JSON_ARRAY(11, file_path, message_path);
				SET type_index = JSON_SET(type_index, CONCAT('$."', full_type_name, '"'), type_entry);
				
				-- Process nested types recursively
				CALL _pb_build_nested_types(message_descriptor, full_type_name, message_path, file_path, type_index);
				
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
				
				-- Add to type index: [kind=14, file_path, type_path]
				SET type_entry = JSON_ARRAY(14, file_path, enum_path);
				SET type_index = JSON_SET(type_index, CONCAT('$."', full_type_name, '"'), type_entry);
				
				SET enum_index = enum_index + 1;
			END WHILE;
		END IF;
		
		SET file_index = file_index + 1;
	END WHILE;
	
	RETURN type_index;
END $$

-- Public function to convert FileDescriptorSet LONGBLOB to descriptor set JSON
-- Returns a 2-element JSON array: [fileDescriptorSet, typeIndex]
DROP FUNCTION IF EXISTS pb_build_descriptor_set_json $$
CREATE FUNCTION pb_build_descriptor_set_json(file_descriptor_set_blob LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE file_descriptor_set_number_json JSON;
	DECLARE type_index JSON;
	DECLARE result JSON;
	
	-- Convert FileDescriptorSet LONGBLOB to protonumberjson format
	SET file_descriptor_set_number_json = _pb_message_to_number_json(
		_pb_get_descriptor_proto_set(),
		'.google.protobuf.FileDescriptorSet',
		file_descriptor_set_blob
	);

	-- Build type index from the FileDescriptorSet
	SET type_index = _pb_build_type_index_from_descriptor_set(file_descriptor_set_number_json);
	
	-- Return 3-element array: [version, fileDescriptorSet, typeIndex]
	SET result = JSON_ARRAY(1, file_descriptor_set_number_json, type_index);
	
	RETURN result;
END $$
