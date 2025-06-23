DELIMITER $$

CREATE TABLE IF NOT EXISTS _Proto_FileDescriptorSet (
	-- VARCHAR(64) is chosen so that users can use SHA-1 (20 bytes) or SHA-256 (32 bytes) as the name of the file descriptor set.
	set_name VARCHAR(64) NOT NULL,
	PRIMARY KEY (set_name)
) $$

CREATE TABLE IF NOT EXISTS _Proto_FileDescriptor (
	set_name VARCHAR(64) NOT NULL,
	file_name VARCHAR(512) NOT NULL,
	package_name VARCHAR(255) NOT NULL,
	syntax VARCHAR(63) NOT NULL,
	editions INT NOT NULL,
	file_options JSON NOT NULL,
	features JSON NOT NULL,
	file_descriptor JSON NOT NULL,
	PRIMARY KEY (`set_name`, `file_name`),
	FOREIGN KEY (`set_name`) REFERENCES _Proto_FileDescriptorSet (`set_name`) ON DELETE CASCADE
) $$

CREATE TABLE IF NOT EXISTS _Proto_MessageDescriptor (
	set_name VARCHAR(64) NOT NULL,
	file_name VARCHAR(512) NOT NULL,
	type_name VARCHAR(512) NOT NULL,
	message_options JSON NOT NULL,
	features JSON NOT NULL,
	message_descriptor JSON NOT NULL,
	PRIMARY KEY (`set_name`, `type_name`),
	FOREIGN KEY (`set_name`, `file_name`) REFERENCES _Proto_FileDescriptor (`set_name`, `file_name`) ON DELETE CASCADE
) $$

CREATE TABLE IF NOT EXISTS _Proto_FieldDescriptor (
	set_name VARCHAR(64) NOT NULL,
	type_name VARCHAR(512) NOT NULL,
	field_number INT NOT NULL,
	field_name TEXT NOT NULL,
	field_label INT NOT NULL,
	field_type INT NOT NULL,
	field_type_name TEXT NULL,
	default_value TEXT NULL,
	json_name TEXT NULL,
	proto3_optional BOOLEAN NOT NULL,
	oneof_index INT NULL,
	field_options JSON NOT NULL,
	features JSON NOT NULL,
	field_descriptor JSON NOT NULL,
	PRIMARY KEY (`set_name`, `type_name`, `field_number`),
	FOREIGN KEY (`set_name`, `type_name`) REFERENCES _Proto_MessageDescriptor (`set_name`, `type_name`) ON DELETE CASCADE
) $$

CREATE TABLE IF NOT EXISTS _Proto_EnumDescriptor (
	set_name VARCHAR(64) NOT NULL,
	file_name VARCHAR(512) NOT NULL,
	type_name VARCHAR(512) NOT NULL,
	enum_options JSON NOT NULL,
	features JSON NOT NULL,
	enum_descriptor JSON NOT NULL,
	PRIMARY KEY (`set_name`, `type_name`),
	FOREIGN KEY (`set_name`, `file_name`) REFERENCES `_Proto_FileDescriptor` (`set_name`, `file_name`) ON DELETE CASCADE
) $$

CREATE TABLE IF NOT EXISTS _Proto_EnumValueDescriptor (
	set_name VARCHAR(64) NOT NULL,
	type_name VARCHAR(512) NOT NULL,
	enum_value_number INT NOT NULL,
	enum_value_name TEXT NOT NULL,
	enum_value_options JSON NOT NULL,
	features JSON NOT NULL,
	enum_value_descriptor JSON NOT NULL,
	PRIMARY KEY (`set_name`, `type_name`, `enum_value_number`),
	FOREIGN KEY (`set_name`, `type_name`) REFERENCES _Proto_EnumDescriptor (`set_name`, `type_name`) ON DELETE CASCADE
) $$

DROP PROCEDURE IF EXISTS _pb_insert_field_descriptor_proto $$
CREATE PROCEDURE _pb_insert_field_descriptor_proto(IN set_name VARCHAR(64), IN full_type_name TEXT, IN field_descriptor LONGBLOB)
BEGIN
	DECLARE field_number INT;
	DECLARE field_options JSON;
	DECLARE features JSON;
	DECLARE field_name TEXT;
	DECLARE field_label INT;
	DECLARE field_type INT;
	DECLARE field_type_name TEXT;
	DECLARE default_value TEXT;
	DECLARE json_name TEXT;
	DECLARE proto3_optional BOOLEAN;
	DECLARE oneof_index INT;
	DECLARE field_descriptor_wire_json JSON;
	DECLARE message_text TEXT;

	SET field_descriptor_wire_json = pb_message_to_wire_json(field_descriptor);

	SET field_number = pb_wire_json_get_int32_field(field_descriptor_wire_json, 3, 0);
	SET field_name = pb_wire_json_get_string_field(field_descriptor_wire_json, 1, NULL);
	SET field_label = pb_wire_json_get_enum_field(field_descriptor_wire_json, 4, 0);
	SET field_type = pb_wire_json_get_enum_field(field_descriptor_wire_json, 5, 0);
	SET field_type_name = pb_wire_json_get_string_field(field_descriptor_wire_json, 6, NULL);
	SET default_value = pb_wire_json_get_string_field(field_descriptor_wire_json, 7, NULL);
	SET json_name = pb_wire_json_get_string_field(field_descriptor_wire_json, 10, NULL);
	SET proto3_optional = pb_wire_json_get_bool_field(field_descriptor_wire_json, 17, FALSE);
	SET oneof_index = pb_wire_json_get_int32_field(field_descriptor_wire_json, 9, NULL);
	SET field_options = pb_message_to_wire_json(pb_wire_json_get_message_field(field_descriptor_wire_json, 8, _binary X''));

	SET features = pb_message_to_wire_json(pb_wire_json_get_message_field(field_options, 21, _binary X''));

	INSERT INTO _Proto_FieldDescriptor (
		set_name,
		type_name,
		field_number,
		field_name,
		field_label,
		field_type,
		field_type_name,
		default_value,
		json_name,
		proto3_optional,
		oneof_index,
		field_options,
		features,
		field_descriptor
	) VALUES (
		set_name,
		full_type_name,
		field_number,
		field_name,
		field_label,
		field_type,
		field_type_name,
		default_value,
		json_name,
		proto3_optional,
		oneof_index,
		field_options,
		features,
		field_descriptor_wire_json
	);
END $$

DROP PROCEDURE IF EXISTS _pb_insert_enum_value_descriptor$$
CREATE PROCEDURE _pb_insert_enum_value_descriptor(IN set_name VARCHAR(64), IN full_type_name TEXT, IN enum_value_descriptor LONGBLOB)
BEGIN
	DECLARE enum_value_number INT;
	DECLARE enum_value_name TEXT;
	DECLARE enum_value_options JSON;
	DECLARE features JSON;
	DECLARE enum_value_descriptor_wire_json JSON;

	SET enum_value_descriptor_wire_json = pb_message_to_wire_json(enum_value_descriptor);

	SET enum_value_number = pb_wire_json_get_int32_field(enum_value_descriptor_wire_json, 2, 0);
	SET enum_value_name = pb_wire_json_get_string_field(enum_value_descriptor_wire_json, 1, '');
	SET enum_value_options = pb_message_to_wire_json(pb_wire_json_get_message_field(enum_value_descriptor_wire_json, 3, _binary X''));

	SET features = pb_message_to_wire_json(pb_wire_json_get_message_field(enum_value_options, 2, _binary X''));

	INSERT INTO _Proto_EnumValueDescriptor (
		set_name,
		type_name,
		enum_value_number,
		enum_value_name,
		enum_value_options,
		features,
		enum_value_descriptor
	) VALUES (
		set_name,
		full_type_name,
		enum_value_number,
		enum_value_name,
		enum_value_options,
		features,
		enum_value_descriptor_wire_json
	);
END $$

DROP PROCEDURE IF EXISTS _pb_insert_enum_descriptor $$
CREATE PROCEDURE _pb_insert_enum_descriptor(IN set_name VARCHAR(64), IN file_name TEXT, IN parent_name TEXT, IN enum_descriptor LONGBLOB)
BEGIN
	DECLARE full_type_name TEXT;
	DECLARE simple_type_name TEXT;
	DECLARE enum_options JSON;
	DECLARE features JSON;
	DECLARE enum_value_descriptor_count INT;
	DECLARE enum_value_descriptor_index INT;
	DECLARE enum_value_descriptor LONGBLOB;
	DECLARE enum_descriptor_wire_json JSON;

	SET enum_descriptor_wire_json = pb_message_to_wire_json(enum_descriptor);

	SET simple_type_name = pb_wire_json_get_string_field(enum_descriptor_wire_json, 1, '');
	SET full_type_name = CONCAT(parent_name, '.', simple_type_name);
	SET enum_options = pb_message_to_wire_json(pb_wire_json_get_message_field(enum_descriptor_wire_json, 3, _binary X''));

	SET features = pb_message_to_wire_json(pb_wire_json_get_message_field(enum_options, 7, _binary X''));

	INSERT INTO _Proto_EnumDescriptor (
		set_name,
		file_name,
		type_name,
		enum_options,
		features,
		enum_descriptor
	) VALUES (
		set_name,
		file_name,
		full_type_name,
		enum_options,
		features,
		enum_descriptor_wire_json
	);

	SET enum_value_descriptor_count = pb_wire_json_get_repeated_message_field_count(enum_descriptor_wire_json, 2);
	SET enum_value_descriptor_index = 0;
	WHILE enum_value_descriptor_index < enum_value_descriptor_count DO
		SET enum_value_descriptor = pb_wire_json_get_repeated_message_field(enum_descriptor_wire_json, 2, enum_value_descriptor_index);
		CALL _pb_insert_enum_value_descriptor(set_name, full_type_name, enum_value_descriptor);
		SET enum_value_descriptor_index = enum_value_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS _pb_insert_message_descriptor $$
CREATE PROCEDURE _pb_insert_message_descriptor(IN set_name VARCHAR(64), IN file_name TEXT, IN parent_name TEXT, IN message_descriptor LONGBLOB)
BEGIN
	DECLARE full_type_name TEXT;
	DECLARE simple_type_name TEXT;
	DECLARE field_descriptor_count INT;
	DECLARE field_descriptor_index INT;
	DECLARE field_descriptor LONGBLOB;
	DECLARE nested_descriptor_count INT;
	DECLARE nested_descriptor_index INT;
	DECLARE nested_descriptor LONGBLOB;
	DECLARE enum_descriptor_count INT;
	DECLARE enum_descriptor_index INT;
	DECLARE enum_descriptor LONGBLOB;
	DECLARE message_options JSON;
	DECLARE features JSON;
	DECLARE message_descriptor_wire_json JSON;

	SET message_descriptor_wire_json = pb_message_to_wire_json(message_descriptor);

	SET simple_type_name = pb_wire_json_get_string_field(message_descriptor_wire_json, 1, '');
	SET full_type_name = CONCAT(parent_name, '.', simple_type_name);
	SET message_options = pb_message_to_wire_json(pb_wire_json_get_message_field(message_descriptor_wire_json, 7, _binary X''));

	SET features = pb_message_to_wire_json(pb_wire_json_get_message_field(message_options, 12, _binary X''));

	INSERT INTO _Proto_MessageDescriptor (
		set_name,
		file_name,
		type_name,
		message_options,
		features,
		message_descriptor
	) VALUES (
		set_name,
		file_name,
		full_type_name,
		message_options,
		features,
		message_descriptor_wire_json
	);

	SET field_descriptor_count = pb_wire_json_get_repeated_message_field_count(message_descriptor_wire_json, 2);
	SET field_descriptor_index = 0;
	WHILE field_descriptor_index < field_descriptor_count DO
		SET field_descriptor = pb_wire_json_get_repeated_message_field(message_descriptor_wire_json, 2, field_descriptor_index);
		CALL _pb_insert_field_descriptor_proto(set_name, full_type_name, field_descriptor);
		SET field_descriptor_index = field_descriptor_index + 1;
	END WHILE;

	SET nested_descriptor_count = pb_wire_json_get_repeated_message_field_count(message_descriptor_wire_json, 3);
	SET nested_descriptor_index = 0;
	WHILE nested_descriptor_index < nested_descriptor_count DO
		SET nested_descriptor = pb_wire_json_get_repeated_message_field(message_descriptor_wire_json, 3, nested_descriptor_index);
		CALL _pb_insert_message_descriptor(set_name, file_name, full_type_name, nested_descriptor);
		SET nested_descriptor_index = nested_descriptor_index + 1;
	END WHILE;

	SET enum_descriptor_count = pb_wire_json_get_repeated_message_field_count(message_descriptor_wire_json, 4);
	SET enum_descriptor_index = 0;
	WHILE enum_descriptor_index < enum_descriptor_count DO
		SET enum_descriptor = pb_wire_json_get_repeated_message_field(message_descriptor_wire_json, 4, enum_descriptor_index);
		CALL _pb_insert_enum_descriptor(set_name, file_name, full_type_name, enum_descriptor);
		SET enum_descriptor_index = enum_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS _pb_insert_file_descriptor $$
CREATE PROCEDURE _pb_insert_file_descriptor(IN set_name VARCHAR(64), IN file_descriptor LONGBLOB)
BEGIN
	DECLARE message_descriptor LONGBLOB;
	DECLARE message_descriptor_count INT;
	DECLARE message_descriptor_index INT;
	DECLARE enum_descriptor_count INT;
	DECLARE enum_descriptor_index INT;
	DECLARE enum_descriptor LONGBLOB;
	DECLARE file_name TEXT;
	DECLARE package_name TEXT;
	DECLARE syntax TEXT;
	DECLARE editions INT;
	DECLARE file_options JSON;
	DECLARE features JSON;
	DECLARE file_descriptor_wire_json JSON;

	SET file_descriptor_wire_json = pb_message_to_wire_json(file_descriptor);

	SET file_name = pb_wire_json_get_string_field(file_descriptor_wire_json, 1, '');
	SET package_name = pb_wire_json_get_string_field(file_descriptor_wire_json, 2, '');
	SET syntax = pb_wire_json_get_string_field(file_descriptor_wire_json, 12, '');
	SET editions = pb_wire_json_get_enum_field(file_descriptor_wire_json, 14, 0);
	SET file_options = pb_message_to_wire_json(pb_wire_json_get_message_field(file_descriptor_wire_json, 8, _binary X''));

	SET features = pb_message_to_wire_json(pb_wire_json_get_message_field(file_options, 50, _binary X''));

	INSERT INTO _Proto_FileDescriptor (
		set_name,
		file_name,
		package_name,
		syntax,
		editions,
		file_options,
		features,
		file_descriptor
	) VALUES (
		set_name,
		file_name,
		package_name,
		syntax,
		editions,
		file_options,
		features,
		file_descriptor_wire_json
	);

	SET message_descriptor_count = pb_wire_json_get_repeated_message_field_count(file_descriptor_wire_json, 4);
	SET message_descriptor_index = 0;
	WHILE message_descriptor_index < message_descriptor_count DO
		SET message_descriptor = pb_wire_json_get_repeated_message_field(file_descriptor_wire_json, 4, message_descriptor_index);
		CALL _pb_insert_message_descriptor(set_name, file_name, IF(package_name = '', '', CONCAT('.', package_name)), message_descriptor);
		SET message_descriptor_index = message_descriptor_index + 1;
	END WHILE;

	SET enum_descriptor_count = pb_wire_json_get_repeated_message_field_count(file_descriptor_wire_json, 5);
	SET enum_descriptor_index = 0;
	WHILE enum_descriptor_index < enum_descriptor_count DO
		SET enum_descriptor = pb_wire_json_get_repeated_message_field(file_descriptor_wire_json, 5, enum_descriptor_index);
		CALL _pb_insert_enum_descriptor(set_name, file_name, IF(package_name = '', '', CONCAT('.', package_name)), enum_descriptor);
		SET enum_descriptor_index = enum_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS pb_descriptor_set_load $$
CREATE PROCEDURE pb_descriptor_set_load(IN set_name VARCHAR(64), IN file_descriptor_set LONGBLOB)
BEGIN
	DECLARE file_descriptor LONGBLOB;
	DECLARE file_descriptor_count INT;
	DECLARE file_descriptor_index INT;
	DECLARE file_descriptor_set_wire_json JSON;

	-- DECLARE done TINYINT DEFAULT FALSE;
	-- DECLARE cur CURSOR FOR
	-- 	SELECT
	-- 		FROM_BASE64(jt.file_descriptor)
	-- 	FROM JSON_TABLE(pb_wire_json_get_repeated_message_field_as_json(file_descriptor_set_wire_json, 1), '$[*]' COLUMNS (
	-- 		file_descriptor LONGBLOB PATH '$'
	-- 	)) AS jt;
	-- DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET file_descriptor_set_wire_json = pb_message_to_wire_json(file_descriptor_set);

	SET @@SESSION.max_sp_recursion_depth = 255;

	INSERT INTO _Proto_FileDescriptorSet (set_name) VALUES (set_name);

	-- OPEN cur;
	-- l1: LOOP
	-- 	FETCH cur INTO file_descriptor;
	-- 	IF done THEN
	-- 		LEAVE l1;
	-- 	END IF;
	-- 	CALL _pb_insert_file_descriptor(set_name, file_descriptor);
	-- END LOOP;
	-- CLOSE cur;

	SET file_descriptor_count = pb_wire_json_get_repeated_message_field_count(file_descriptor_set_wire_json, 1);
	SET file_descriptor_index = 0;
	WHILE file_descriptor_index < file_descriptor_count DO
		SET file_descriptor = pb_wire_json_get_repeated_message_field(file_descriptor_set_wire_json, 1, file_descriptor_index);
		CALL _pb_insert_file_descriptor(set_name, file_descriptor);
		SET file_descriptor_index = file_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS pb_descriptor_set_delete $$
CREATE PROCEDURE pb_descriptor_set_delete(IN set_name VARCHAR(64))
BEGIN
	DELETE FROM _Proto_FileDescriptorSet t WHERE t.set_name = set_name;
END $$

DROP FUNCTION IF EXISTS pb_descriptor_set_exists $$
CREATE FUNCTION pb_descriptor_set_exists(set_name VARCHAR(64)) RETURNS BOOLEAN READS SQL DATA
BEGIN
	DECLARE exist INT;
	SET exist = (SELECT count(*) FROM _Proto_FileDescriptorSet t WHERE t.set_name = set_name) > 0;
	RETURN exist > 0;
END $$

DROP FUNCTION IF EXISTS pb_descriptor_set_contains_message_type $$
CREATE FUNCTION pb_descriptor_set_contains_message_type(set_name VARCHAR(64), full_type_name VARCHAR(512)) RETURNS BOOLEAN READS SQL DATA
BEGIN
	DECLARE exist INT;
	SET exist = (SELECT count(*) FROM _Proto_MessageDescriptor t WHERE t.set_name = set_name AND t.type_name = full_type_name) > 0;
	RETURN exist > 0;
END $$

DROP FUNCTION IF EXISTS pb_descriptor_set_contains_enum_type $$
CREATE FUNCTION pb_descriptor_set_contains_enum_type(set_name VARCHAR(64), full_type_name VARCHAR(512)) RETURNS BOOLEAN READS SQL DATA
BEGIN
	DECLARE exist INT;
	SET exist = (SELECT count(*) FROM _Proto_EnumDescriptor t WHERE t.set_name = set_name AND t.type_name = full_type_name) > 0;
	RETURN exist > 0;
END $$
