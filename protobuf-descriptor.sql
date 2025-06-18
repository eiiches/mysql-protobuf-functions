DELIMITER $$

DROP PROCEDURE IF EXISTS _pb_insert_field_descriptor_proto $$
CREATE PROCEDURE _pb_insert_field_descriptor_proto(IN table_name VARCHAR(255), IN full_type_name TEXT, IN field_descriptor BLOB)
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
	SET json_name = IF(pb_message_has_string_field(field_descriptor, 8), pb_message_get_string_field(field_descriptor, 8, NULL), NULL);
	SET proto3_optional = pb_message_get_bool_field(field_descriptor, 17, NULL);
	SET field_options = pb_message_get_message_field(field_descriptor, 8, NULL);
	SET features = pb_message_get_message_field(field_options, 21, NULL);
	SET oneof_index = IF(pb_message_has_int32_field(field_descriptor, 9), pb_message_get_int32_field(field_descriptor, 9, NULL), NULL);

	SET @insert_sql = CONCAT('INSERT INTO `_Proto_FieldDescriptor_', table_name, '` (type_name, field_number, field_name, field_label, field_type, field_type_name, default_value, json_name, proto3_optional, oneof_index, field_options, features, field_descriptor) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)');
	PREPARE stmt FROM @insert_sql;
	SET @full_type_name = full_type_name;
	SET @field_number = field_number;
	SET @field_name = field_name;
	SET @field_label = field_label;
	SET @field_type = field_type;
	SET @field_type_name = field_type_name;
	SET @default_value = default_value;
	SET @json_name = json_name;
	SET @proto3_optional = proto3_optional;
	SET @oneof_index = oneof_index;
	SET @field_options = field_options;
	SET @features = features;
	SET @field_descriptor = field_descriptor;
	EXECUTE stmt USING @full_type_name, @field_number, @field_name, @field_label, @field_type, @field_type_name, @default_value, @json_name, @proto3_optional, @oneof_index, @field_options, @features, @field_descriptor;
	DEALLOCATE PREPARE stmt;
END $$

DROP PROCEDURE IF EXISTS _pb_insert_enum_value_descriptor$$
CREATE PROCEDURE _pb_insert_enum_value_descriptor(IN table_name VARCHAR(255), IN full_type_name TEXT, IN enum_value_descriptor BLOB)
BEGIN
	DECLARE enum_value_number INT;
	DECLARE enum_value_name TEXT;
	DECLARE enum_value_options BLOB;
	DECLARE features BLOB;

	SET enum_value_number = pb_message_get_int32_field(enum_value_descriptor, 2, NULL);
	SET enum_value_name = pb_message_get_string_field(enum_value_descriptor, 1, NULL);
	SET enum_value_options = pb_message_get_message_field(enum_value_descriptor, 3, NULL);
	SET features = pb_message_get_message_field(enum_value_options, 2, NULL);

	SET @insert_sql = CONCAT('INSERT INTO `_Proto_EnumValueDescriptor_', table_name, '` (type_name, enum_value_number, enum_value_name, enum_value_options, features, enum_value_descriptor) VALUES (?, ?, ?, ?, ?, ?)');
	PREPARE stmt FROM @insert_sql;
	SET @full_type_name = full_type_name;
	SET @enum_value_number = enum_value_number;
	SET @enum_value_name = enum_value_name;
	SET @enum_value_options = enum_value_options;
	SET @features = features;
	SET @enum_value_descriptor = enum_value_descriptor;
	EXECUTE stmt USING @full_type_name, @enum_value_number, @enum_value_name, @enum_value_options, @features, @enum_value_descriptor;
	DEALLOCATE PREPARE stmt;
END $$

DROP PROCEDURE IF EXISTS _pb_insert_enum_descriptor $$
CREATE PROCEDURE _pb_insert_enum_descriptor(IN table_name VARCHAR(255), IN file_name TEXT, IN parent_name TEXT, IN enum_descriptor BLOB)
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

	SET @insert_sql = CONCAT('INSERT INTO `_Proto_EnumDescriptor_', table_name, '` (file_name, type_name, enum_options, features, enum_descriptor) VALUES (?, ?, ?, ?, ?)');
	PREPARE stmt FROM @insert_sql;
	SET @file_name = file_name;
	SET @full_type_name = full_type_name;
	SET @enum_options = enum_options;
	SET @features = features;
	SET @enum_descriptor = enum_descriptor;
	EXECUTE stmt USING @file_name, @full_type_name, @enum_options, @features, @enum_descriptor;
	DEALLOCATE PREPARE stmt;

	SET enum_value_descriptor_count = pb_message_get_message_field_count(enum_descriptor, 2);
	SET enum_value_descriptor_index = 0;
	WHILE enum_value_descriptor_index < enum_value_descriptor_count DO
		SET enum_value_descriptor = pb_message_get_message_field(enum_descriptor, 2, enum_value_descriptor_index);
		CALL _pb_insert_enum_value_descriptor(table_name, full_type_name, enum_value_descriptor);
		SET enum_value_descriptor_index = enum_value_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS _pb_insert_message_descriptor $$
CREATE PROCEDURE _pb_insert_message_descriptor(IN table_name VARCHAR(255), IN file_name TEXT, IN parent_name TEXT, IN message_descriptor BLOB)
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

	SET @insert_sql = CONCAT('INSERT INTO `_Proto_MessageDescriptor_', table_name, '` (file_name, type_name, message_options, features, message_descriptor) VALUES (?, ?, ?, ?, ?)');
	PREPARE stmt FROM @insert_sql;
	SET @file_name = file_name;
	SET @full_type_name = full_type_name;
	SET @message_options = message_options;
	SET @message_descriptor = message_descriptor;
	EXECUTE stmt USING @file_name, @full_type_name, @message_options, @features, @message_descriptor;
	DEALLOCATE PREPARE stmt;

	SET field_descriptor_count = pb_message_get_message_field_count(message_descriptor, 2);
	SET field_descriptor_index = 0;
	WHILE field_descriptor_index < field_descriptor_count DO
		SET field_descriptor = pb_message_get_message_field(message_descriptor, 2, field_descriptor_index);
		CALL _pb_insert_field_descriptor_proto(table_name, full_type_name, field_descriptor);
		SET field_descriptor_index = field_descriptor_index + 1;
	END WHILE;

	SET nested_descriptor_count = pb_message_get_message_field_count(message_descriptor, 3);
	SET nested_descriptor_index = 0;
	WHILE nested_descriptor_index < nested_descriptor_count DO
		SET nested_descriptor = pb_message_get_message_field(message_descriptor, 3, nested_descriptor_index);
		CALL _pb_insert_message_descriptor(table_name, file_name, full_type_name, nested_descriptor);
		SET nested_descriptor_index = nested_descriptor_index + 1;
	END WHILE;

	SET enum_descriptor_count = pb_message_get_message_field_count(message_descriptor, 4);
	SET enum_descriptor_index = 0;
	WHILE enum_descriptor_index < enum_descriptor_count DO
		SET enum_descriptor = pb_message_get_message_field(message_descriptor, 4, enum_descriptor_index);
		CALL _pb_insert_enum_descriptor(table_name, file_name, full_type_name, enum_descriptor);
		SET enum_descriptor_index = enum_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS _pb_insert_file_descriptor $$
CREATE PROCEDURE _pb_insert_file_descriptor(IN table_name VARCHAR(255), IN file_descriptor BLOB)
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

	SET @insert_sql = CONCAT('INSERT INTO `_Proto_FileDescriptor_', table_name, '` (file_name, package_name, syntax, editions, file_options, features, file_descriptor) VALUES (?, ?, ?, ?, ?, ?, ?)');
	PREPARE stmt FROM @insert_sql;
	SET @file_name = file_name;
	SET @package_name = package_name;
	SET @syntax = syntax;
	SET @editions = editions;
	SET @file_options = file_options;
	SET @features = features;
	SET @file_descriptor = file_descriptor;
	EXECUTE stmt USING @file_name, @package_name, @syntax, @editions, @file_options, @features, @file_descriptor;
	DEALLOCATE PREPARE stmt;

	SET message_descriptor_count = pb_message_get_message_field_count(file_descriptor, 4);
	SET message_descriptor_index = 0;
	WHILE message_descriptor_index < message_descriptor_count DO
		SET message_descriptor = pb_message_get_message_field(file_descriptor, 4, message_descriptor_index);
		CALL _pb_insert_message_descriptor(table_name, file_name, package_name, message_descriptor);
		SET message_descriptor_index = message_descriptor_index + 1;
	END WHILE;

	SET enum_descriptor_count = pb_message_get_message_field_count(file_descriptor, 5);
	SET enum_descriptor_index = 0;
	WHILE enum_descriptor_index < enum_descriptor_count DO
		SET enum_descriptor = pb_message_get_message_field(file_descriptor, 5, enum_descriptor_index);
		CALL _pb_insert_enum_descriptor(table_name, file_name, package_name, enum_descriptor);
		SET enum_descriptor_index = enum_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS pb_load_file_descriptor_set $$
CREATE PROCEDURE pb_load_file_descriptor_set(IN table_name VARCHAR(255), IN persist BOOLEAN, IN file_descriptor_set BLOB)
BEGIN
	DECLARE file_descriptor BLOB;
	DECLARE file_descriptor_count INT;
	DECLARE file_descriptor_index INT;

	SET @@SESSION.max_sp_recursion_depth = 255;

	SET @create_sql = CONCAT('CREATE', IF(persist, '', ' TEMPORARY'), ' TABLE `_Proto_FileDescriptor_', table_name, '` (file_name VARCHAR(512) NOT NULL, package_name VARCHAR(255) NOT NULL, syntax VARCHAR(63) NOT NULL, editions INT NOT NULL, file_options BLOB NOT NULL, features BLOB NOT NULL, file_descriptor BLOB NOT NULL, PRIMARY KEY (`file_name`))');
	PREPARE stmt FROM @create_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;

	SET @create_sql = CONCAT('CREATE', IF(persist, '', ' TEMPORARY'), ' TABLE `_Proto_MessageDescriptor_', table_name, '` (file_name VARCHAR(512) NOT NULL, type_name VARCHAR(512) NOT NULL, message_options BLOB NOT NULL, features BLOB NOT NULL, message_descriptor BLOB NOT NULL, PRIMARY KEY (`type_name`), FOREIGN KEY (`file_name`) REFERENCES `_Proto_FileDescriptor_', table_name, '` (`file_name`))');
	PREPARE stmt FROM @create_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;

	SET @create_sql = CONCAT('CREATE', IF(persist, '', ' TEMPORARY'), ' TABLE `_Proto_FieldDescriptor_', table_name, '` (type_name VARCHAR(512) NOT NULL, field_number INT NOT NULL, field_name TEXT NOT NULL, field_label INT NOT NULL, field_type INT NOT NULL, field_type_name TEXT NULL, default_value TEXT NULL, json_name TEXT NULL, proto3_optional BOOLEAN NOT NULL, oneof_index INT NULL, field_options BLOB NOT NULL, features BLOB NOT NULL, field_descriptor BLOB NOT NULL, PRIMARY KEY (`type_name`, `field_number`), FOREIGN KEY (`type_name`) REFERENCES `_Proto_MessageDescriptor_', table_name, '` (`type_name`))');
	PREPARE stmt FROM @create_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;

	SET @create_sql = CONCAT('CREATE', IF(persist, '', ' TEMPORARY'), ' TABLE `_Proto_EnumDescriptor_', table_name, '` (file_name VARCHAR(512) NOT NULL, type_name VARCHAR(512) NOT NULL, enum_options BLOB NOT NULL, features BLOB NOT NULL, enum_descriptor BLOB NOT NULL, PRIMARY KEY (`type_name`), FOREIGN KEY (`file_name`) REFERENCES `_Proto_FileDescriptor_', table_name, '` (`file_name`))');
	PREPARE stmt FROM @create_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;

	SET @create_sql = CONCAT('CREATE', IF(persist, '', ' TEMPORARY'), ' TABLE `_Proto_EnumValueDescriptor_', table_name, '` (type_name VARCHAR(512) NOT NULL, enum_value_number INT NOT NULL, enum_value_name TEXT NOT NULL, enum_value_options BLOB NOT NULL, features BLOB NOT NULL, enum_value_descriptor BLOB NOT NULL, PRIMARY KEY (`type_name`, `enum_value_number`), FOREIGN KEY (`type_name`) REFERENCES `_Proto_EnumDescriptor_', table_name, '` (`type_name`))');
	PREPARE stmt FROM @create_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;

	SET file_descriptor_count = pb_message_get_message_field_count(file_descriptor_set, 1);
	SET file_descriptor_index = 0;
	WHILE file_descriptor_index < file_descriptor_count DO
		SET file_descriptor = pb_message_get_message_field(file_descriptor_set, 1, file_descriptor_index);
		CALL _pb_insert_file_descriptor(table_name, file_descriptor);
		SET file_descriptor_index = file_descriptor_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS pb_delete_file_descriptor_set $$
CREATE PROCEDURE pb_delete_file_descriptor_set(IN table_name VARCHAR(255))
BEGIN
	SET @drop_sql = CONCAT('DROP TABLE IF EXISTS `_Proto_EnumValueDescriptor_', table_name, '`');
	PREPARE stmt FROM @drop_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;

	SET @drop_sql = CONCAT('DROP TABLE IF EXISTS `_Proto_EnumDescriptor_', table_name, '`');
	PREPARE stmt FROM @drop_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;

	SET @drop_sql = CONCAT('DROP TABLE IF EXISTS `_Proto_FieldDescriptor_', table_name, '`');
	PREPARE stmt FROM @drop_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;

	SET @drop_sql = CONCAT('DROP TABLE IF EXISTS `_Proto_MessageDescriptor_', table_name, '`');
	PREPARE stmt FROM @drop_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;

	SET @drop_sql = CONCAT('DROP TABLE IF EXISTS `_Proto_FileDescriptor_', table_name, '`');
	PREPARE stmt FROM @drop_sql;
	EXECUTE stmt;
	DEALLOCATE PREPARE stmt;
END $$
