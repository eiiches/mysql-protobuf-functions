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
	file_options BLOB NOT NULL,
	features BLOB NOT NULL,
	file_descriptor BLOB NOT NULL,
	PRIMARY KEY (`set_name`, `file_name`),
	FOREIGN KEY (`set_name`) REFERENCES _Proto_FileDescriptorSet (`set_name`) ON DELETE CASCADE
) $$

CREATE TABLE IF NOT EXISTS _Proto_MessageDescriptor (
	set_name VARCHAR(64) NOT NULL,
	file_name VARCHAR(512) NOT NULL,
	type_name VARCHAR(512) NOT NULL,
	message_options BLOB NOT NULL,
	features BLOB NOT NULL,
	message_descriptor BLOB NOT NULL,
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
	field_options BLOB NOT NULL,
	features BLOB NOT NULL,
	field_descriptor BLOB NOT NULL,
	PRIMARY KEY (`set_name`, `type_name`, `field_number`),
	FOREIGN KEY (`set_name`, `type_name`) REFERENCES _Proto_MessageDescriptor (`set_name`, `type_name`) ON DELETE CASCADE
) $$

CREATE TABLE IF NOT EXISTS _Proto_EnumDescriptor (
	set_name VARCHAR(64) NOT NULL,
	file_name VARCHAR(512) NOT NULL,
	type_name VARCHAR(512) NOT NULL,
	enum_options BLOB NOT NULL,
	features BLOB NOT NULL,
	enum_descriptor BLOB NOT NULL,
	PRIMARY KEY (`set_name`, `type_name`),
	FOREIGN KEY (`set_name`, `file_name`) REFERENCES `_Proto_FileDescriptor` (`set_name`, `file_name`) ON DELETE CASCADE
) $$

CREATE TABLE IF NOT EXISTS _Proto_EnumValueDescriptor (
	set_name VARCHAR(64) NOT NULL,
	type_name VARCHAR(512) NOT NULL,
	enum_value_number INT NOT NULL,
	enum_value_name TEXT NOT NULL,
	enum_value_options BLOB NOT NULL,
	features BLOB NOT NULL,
	enum_value_descriptor BLOB NOT NULL,
	PRIMARY KEY (`set_name`, `type_name`, `enum_value_number`),
	FOREIGN KEY (`set_name`, `type_name`) REFERENCES _Proto_EnumDescriptor (`set_name`, `type_name`) ON DELETE CASCADE
) $$

DROP PROCEDURE IF EXISTS _pb_insert_field_descriptor_proto $$
CREATE PROCEDURE _pb_insert_field_descriptor_proto(IN set_name VARCHAR(64), IN full_type_name TEXT, IN field_descriptor BLOB)
BEGIN
	DECLARE field_number INT;
	DECLARE field_options BLOB;
	DECLARE features BLOB;
	DECLARE field_name TEXT;
	DECLARE field_label INT;
	DECLARE field_type INT;
	DECLARE field_type_name TEXT;
	DECLARE default_value TEXT;
	DECLARE json_name TEXT;
	DECLARE proto3_optional BOOLEAN;
	DECLARE oneof_index INT;

	SET field_number = pb_message_get_int32_field(field_descriptor, 3, NULL);
	SET field_name = pb_message_get_string_field(field_descriptor, 1, NULL);
	SET field_label = pb_message_get_enum_field(field_descriptor, 4, NULL);
	SET field_type = pb_message_get_enum_field(field_descriptor, 5, NULL);
	SET field_type_name = IF(pb_message_has_string_field(field_descriptor, 6), pb_message_get_string_field(field_descriptor, 6, NULL), NULL);
	SET default_value = IF(pb_message_has_string_field(field_descriptor, 7), pb_message_get_string_field(field_descriptor, 7, NULL), NULL);
	SET json_name = IF(pb_message_has_string_field(field_descriptor, 10), pb_message_get_string_field(field_descriptor, 10, NULL), NULL);
	SET proto3_optional = pb_message_get_bool_field(field_descriptor, 17, NULL);
	SET field_options = pb_message_get_message_field(field_descriptor, 8, NULL);
	SET features = pb_message_get_message_field(field_options, 21, NULL);
	SET oneof_index = IF(pb_message_has_int32_field(field_descriptor, 9), pb_message_get_int32_field(field_descriptor, 9, NULL), NULL);

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
		field_descriptor
	);
END $$

DROP PROCEDURE IF EXISTS _pb_insert_enum_value_descriptor$$
CREATE PROCEDURE _pb_insert_enum_value_descriptor(IN set_name VARCHAR(64), IN full_type_name TEXT, IN enum_value_descriptor BLOB)
BEGIN
	DECLARE enum_value_number INT;
	DECLARE enum_value_name TEXT;
	DECLARE enum_value_options BLOB;
	DECLARE features BLOB;

	SET enum_value_number = pb_message_get_int32_field(enum_value_descriptor, 2, NULL);
	SET enum_value_name = pb_message_get_string_field(enum_value_descriptor, 1, NULL);
	SET enum_value_options = pb_message_get_message_field(enum_value_descriptor, 3, NULL);
	SET features = pb_message_get_message_field(enum_value_options, 2, NULL);

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
		enum_value_descriptor
	);
END $$

DROP PROCEDURE IF EXISTS _pb_insert_enum_descriptor $$
CREATE PROCEDURE _pb_insert_enum_descriptor(IN set_name VARCHAR(64), IN file_name TEXT, IN parent_name TEXT, IN enum_descriptor BLOB)
BEGIN
	DECLARE full_type_name TEXT;
	DECLARE simple_type_name TEXT;
	DECLARE enum_options BLOB;
	DECLARE features BLOB;
	DECLARE enum_value_descriptor_count INT;
	DECLARE enum_value_descriptor_index INT;
	DECLARE enum_value_descriptor BLOB;

	SET simple_type_name = pb_message_get_string_field(enum_descriptor, 1, NULL);
	SET full_type_name = CONCAT(IF(parent_name = '', '', CONCAT(parent_name, '.')), simple_type_name);
	SET enum_options = pb_message_get_message_field(enum_descriptor, 3, NULL);
	SET features = pb_message_get_message_field(enum_options, 7, NULL);

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
		enum_descriptor
	);

	SET enum_value_descriptor_count = pb_message_get_message_field_count(enum_descriptor, 2);
	SET enum_value_descriptor_index = 0;
	WHILE enum_value_descriptor_index < enum_value_descriptor_count DO
		SET enum_value_descriptor = pb_message_get_message_field(enum_descriptor, 2, enum_value_descriptor_index);
		CALL _pb_insert_enum_value_descriptor(set_name, full_type_name, enum_value_descriptor);
		SET enum_value_descriptor_index = enum_value_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS _pb_insert_message_descriptor $$
CREATE PROCEDURE _pb_insert_message_descriptor(IN set_name VARCHAR(64), IN file_name TEXT, IN parent_name TEXT, IN message_descriptor BLOB)
BEGIN
	DECLARE full_type_name TEXT;
	DECLARE simple_type_name TEXT;
	DECLARE field_descriptor_count INT;
	DECLARE field_descriptor_index INT;
	DECLARE field_descriptor BLOB;
	DECLARE nested_descriptor_count INT;
	DECLARE nested_descriptor_index INT;
	DECLARE nested_descriptor BLOB;
	DECLARE enum_descriptor_count INT;
	DECLARE enum_descriptor_index INT;
	DECLARE enum_descriptor BLOB;
	DECLARE message_options BLOB;
	DECLARE features BLOB;

	SET simple_type_name = pb_message_get_string_field(message_descriptor, 1, NULL);
	SET full_type_name = CONCAT(IF(parent_name = '', '', CONCAT(parent_name, '.')), simple_type_name);
	SET message_options = pb_message_get_message_field(message_descriptor, 7, NULL);
	SET features = pb_message_get_message_field(message_options, 12, NULL);

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
		message_descriptor
	);

	SET field_descriptor_count = pb_message_get_message_field_count(message_descriptor, 2);
	SET field_descriptor_index = 0;
	WHILE field_descriptor_index < field_descriptor_count DO
		SET field_descriptor = pb_message_get_message_field(message_descriptor, 2, field_descriptor_index);
		CALL _pb_insert_field_descriptor_proto(set_name, full_type_name, field_descriptor);
		SET field_descriptor_index = field_descriptor_index + 1;
	END WHILE;

	SET nested_descriptor_count = pb_message_get_message_field_count(message_descriptor, 3);
	SET nested_descriptor_index = 0;
	WHILE nested_descriptor_index < nested_descriptor_count DO
		SET nested_descriptor = pb_message_get_message_field(message_descriptor, 3, nested_descriptor_index);
		CALL _pb_insert_message_descriptor(set_name, file_name, full_type_name, nested_descriptor);
		SET nested_descriptor_index = nested_descriptor_index + 1;
	END WHILE;

	SET enum_descriptor_count = pb_message_get_message_field_count(message_descriptor, 4);
	SET enum_descriptor_index = 0;
	WHILE enum_descriptor_index < enum_descriptor_count DO
		SET enum_descriptor = pb_message_get_message_field(message_descriptor, 4, enum_descriptor_index);
		CALL _pb_insert_enum_descriptor(set_name, file_name, full_type_name, enum_descriptor);
		SET enum_descriptor_index = enum_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS _pb_insert_file_descriptor $$
CREATE PROCEDURE _pb_insert_file_descriptor(IN set_name VARCHAR(64), IN file_descriptor BLOB)
BEGIN
	DECLARE message_descriptor BLOB;
	DECLARE message_descriptor_count INT;
	DECLARE message_descriptor_index INT;
	DECLARE enum_descriptor_count INT;
	DECLARE enum_descriptor_index INT;
	DECLARE enum_descriptor BLOB;
	DECLARE file_name TEXT;
	DECLARE package_name TEXT;
	DECLARE syntax TEXT;
	DECLARE editions INT;
	DECLARE file_options BLOB;
	DECLARE features BLOB;

	SET file_name = pb_message_get_string_field(file_descriptor, 1, NULL);
	SET package_name = pb_message_get_string_field(file_descriptor, 2, NULL);
	SET syntax = pb_message_get_string_field(file_descriptor, 12, NULL);
	SET editions = pb_message_get_enum_field(file_descriptor, 14, NULL);
	SET file_options = pb_message_get_message_field(file_descriptor, 8, NULL);
	SET features = pb_message_get_message_field(file_options, 50, NULL);

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
		file_descriptor
	);

	SET message_descriptor_count = pb_message_get_message_field_count(file_descriptor, 4);
	SET message_descriptor_index = 0;
	WHILE message_descriptor_index < message_descriptor_count DO
		SET message_descriptor = pb_message_get_message_field(file_descriptor, 4, message_descriptor_index);
		CALL _pb_insert_message_descriptor(set_name, file_name, package_name, message_descriptor);
		SET message_descriptor_index = message_descriptor_index + 1;
	END WHILE;

	SET enum_descriptor_count = pb_message_get_message_field_count(file_descriptor, 5);
	SET enum_descriptor_index = 0;
	WHILE enum_descriptor_index < enum_descriptor_count DO
		SET enum_descriptor = pb_message_get_message_field(file_descriptor, 5, enum_descriptor_index);
		CALL _pb_insert_enum_descriptor(set_name, file_name, package_name, enum_descriptor);
		SET enum_descriptor_index = enum_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS pb_descriptor_set_load $$
CREATE PROCEDURE pb_descriptor_set_load(IN set_name VARCHAR(64), IN file_descriptor_set BLOB)
BEGIN
	DECLARE file_descriptor BLOB;
	DECLARE file_descriptor_count INT;
	DECLARE file_descriptor_index INT;

	SET @@SESSION.max_sp_recursion_depth = 255;

	INSERT INTO _Proto_FileDescriptorSet (set_name) VALUES (set_name);

	SET file_descriptor_count = pb_message_get_message_field_count(file_descriptor_set, 1);
	SET file_descriptor_index = 0;
	WHILE file_descriptor_index < file_descriptor_count DO
		SET file_descriptor = pb_message_get_message_field(file_descriptor_set, 1, file_descriptor_index);
		CALL _pb_insert_file_descriptor(set_name, file_descriptor);
		SET file_descriptor_index = file_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS pb_descriptor_set_delete $$
CREATE PROCEDURE pb_descriptor_set_delete(IN set_name VARCHAR(64))
BEGIN
	DELETE FROM _Proto_FileDescriptorSet WHERE set_name = set_name;
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
