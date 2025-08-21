DELIMITER $$

-- Helper function to get message descriptor from descriptor set JSON
DROP FUNCTION IF EXISTS _pb_get_message_descriptor $$
CREATE FUNCTION _pb_get_message_descriptor(descriptor_set_json JSON, type_name TEXT) RETURNS JSON DETERMINISTIC
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
DROP FUNCTION IF EXISTS _pb_get_enum_descriptor $$
CREATE FUNCTION _pb_get_enum_descriptor(descriptor_set_json JSON, type_name TEXT) RETURNS JSON DETERMINISTIC
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
DROP FUNCTION IF EXISTS _pb_get_file_descriptor $$
CREATE FUNCTION _pb_get_file_descriptor(descriptor_set_json JSON, type_name TEXT) RETURNS JSON DETERMINISTIC
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
