DELIMITER $$

DROP FUNCTION IF EXISTS _pb_util_snake_to_lower_camel $$
CREATE FUNCTION _pb_util_snake_to_lower_camel(s TEXT) RETURNS TEXT DETERMINISTIC
BEGIN
	-- TODO: implement
	RETURN s;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_primitive_field_as_json $$
CREATE PROCEDURE _pb_wire_json_get_primitive_field_as_json(IN wire_json JSON, IN field_number INT, IN field_type INT, IN is_repeated BOOLEAN, IN has_field_presence BOOLEAN, OUT field_json_value JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE boolean_value BOOLEAN;

	CASE field_type
	WHEN 1 THEN -- double
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_double_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_double_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 2 THEN -- float
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_float_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_float_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 3 THEN -- int64
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_int64_field_as_json_string_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_int64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
		END IF;
	WHEN 4 THEN -- uint64
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_uint64_field_as_json_string_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_uint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
		END IF;
	WHEN 5 THEN -- int32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_int32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_int32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 6 THEN -- fixed64
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_fixed64_field_as_json_string_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_fixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
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
			SET field_json_value = JSON_QUOTE(TO_BASE64(pb_wire_json_get_bytes_field(wire_json, field_number, IF(has_field_presence, NULL, _binary X''))));
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
			SET field_json_value = pb_wire_json_get_repeated_sfixed64_field_as_json_string_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_sfixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
		END IF;
	WHEN 17 THEN -- sint32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_sint32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_sint32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 18 THEN -- sint64
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_sint64_field_as_json_string_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_sint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
		END IF;
	ELSE
		SET message_text = CONCAT('_pb_message_to_json: unknown field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$

DROP PROCEDURE IF EXISTS _pb_enum_to_json $$
CREATE PROCEDURE _pb_enum_to_json(IN set_name VARCHAR(64), IN full_type_name VARCHAR(512), IN enum_value_number INT, OUT result JSON)
BEGIN
	DECLARE enum_value_name TEXT;
	DECLARE message_text TEXT;

	IF NOT pb_descriptor_set_exists(set_name) THEN
		SET message_text = CONCAT('_pb_enum_to_json: descriptor set `', set_name, '` does not exist');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;
	IF NOT pb_descriptor_set_contains_enum_type(set_name, full_type_name) THEN
		SET message_text = CONCAT('_pb_enum_to_json: message type `', full_type_name, '` does not exist in descriptor set `', set_name, '`');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	SELECT
		enum_value.enum_value_name
	FROM _Proto_EnumValueDescriptor enum_value
	WHERE
		enum_value.set_name = set_name
		AND enum_value.type_name = full_type_name
		AND enum_value.enum_value_number = enum_value_number
	INTO
		enum_value_name;

	IF enum_value_name IS NULL THEN
		SET result = NULL;
	ELSE
		SET result = JSON_QUOTE(enum_value_name);
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_message_to_json $$
CREATE PROCEDURE _pb_message_to_json(IN set_name VARCHAR(64), IN full_type_name VARCHAR(512), IN buf LONGBLOB, OUT result JSON)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE field_number INT;
	DECLARE field_name TEXT;
	DECLARE field_label INT;
	DECLARE field_type INT;
	DECLARE field_type_name TEXT;
	DECLARE json_name TEXT;
	DECLARE proto3_optional BOOLEAN;
	DECLARE wire_json JSON;
	DECLARE syntax TEXT;
	DECLARE is_map BOOLEAN;
	DECLARE is_repeated BOOLEAN;
	DECLARE has_field_presence BOOLEAN;
	DECLARE oneof_index INT;
	DECLARE default_value TEXT;
	DECLARE field_index INT;
	DECLARE field_count INT;
	DECLARE field_descriptor JSON;
	DECLARE map_key_type INT;
	DECLARE map_value_type INT;
	DECLARE map_value_type_name TEXT;

	DECLARE bytes_value LONGBLOB;
	DECLARE nested_json_value JSON;
	DECLARE field_json_value JSON;
	DECLARE elements JSON;
	DECLARE element JSON;
	DECLARE map_key JSON;
	DECLARE map_value JSON;

	DECLARE oneofs JSON;
	DECLARE oneof_priority INT;
	DECLARE oneof_priority_prev INT;

	DECLARE json_field_name TEXT;

	-- NOTE: always use alias in select columns, to avoid confusion with variables with the same name.
	DECLARE cur CURSOR FOR
		SELECT
			file.syntax,
			field_message_type.map_entry,
			field.field_number,
			field.field_name,
			field.field_label,
			field.field_type,
			field.field_type_name,
			field.json_name,
			field.proto3_optional,
			field.oneof_index,
			field.default_value,
			map_key.field_type AS map_key_type,
			map_value.field_type AS map_value_type,
			map_value.field_type_name AS map_value_type_name,
			field.field_descriptor
		FROM _Proto_FieldDescriptor field
			INNER JOIN _Proto_MessageDescriptor message USING (set_name, type_name)
			INNER JOIN _Proto_FileDescriptor file USING (set_name, file_name)
			LEFT JOIN _Proto_MessageDescriptor field_message_type ON
				field.field_type = 11
				AND field.set_name = field_message_type.set_name
				AND field.field_type_name = field_message_type.type_name
			LEFT JOIN _Proto_EnumDescriptor field_enum_type ON
				field.field_type = 14
				AND field.set_name = field_enum_type.set_name
				AND field.field_type_name = field_enum_type.type_name
			LEFT JOIN _Proto_FieldDescriptor map_key ON
				field_message_type.map_entry IS TRUE
				AND map_key.set_name = field.set_name
				AND map_key.type_name = field_message_type.type_name
				AND map_key.field_number = 1
			LEFT JOIN _Proto_FieldDescriptor map_value ON
				field_message_type.map_entry IS TRUE
				AND map_value.set_name = field.set_name
				AND map_value.type_name = field_message_type.type_name
				AND map_value.field_number = 2
		WHERE
			field.set_name = set_name
			AND field.type_name = full_type_name;

	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET @@SESSION.max_sp_recursion_depth = 255;

	IF NOT pb_descriptor_set_exists(set_name) THEN
		SET message_text = CONCAT('_pb_message_to_json: descriptor set `', set_name, '` does not exist');
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
	END IF;
	IF NOT pb_descriptor_set_contains_message_type(set_name, full_type_name) THEN
		SET message_text = CONCAT('_pb_message_to_json: message type `', full_type_name, '` does not exist in descriptor set `', set_name, '`');
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
	END IF;

	SET result = JSON_OBJECT();
	SET oneofs = JSON_OBJECT();
	SET wire_json = pb_message_to_wire_json(buf);

	OPEN cur;

	-- TODO: Support oneof.

	l1: LOOP
		-- FETCH cur INTO syntax, field_descriptor;
		FETCH cur INTO syntax, is_map, field_number, field_name, field_label, field_type, field_type_name, json_name, proto3_optional, oneof_index, default_value, map_key_type, map_value_type, map_value_type_name, field_descriptor;
		IF done THEN
			LEAVE l1;
		END IF;

		-- SET field_number = pb_wire_json_get_int32_field(field_descriptor, 3, 0);
		-- SET field_name = pb_wire_json_get_string_field(field_descriptor, 1, NULL);
		-- SET field_label = pb_wire_json_get_enum_field(field_descriptor, 4, 0);
		-- SET field_type = pb_wire_json_get_enum_field(field_descriptor, 5, 0);
		-- SET field_type_name = pb_wire_json_get_string_field(field_descriptor, 6, NULL);
		-- SET default_value = pb_wire_json_get_string_field(field_descriptor, 7, NULL);
		-- SET json_name = pb_wire_json_get_string_field(field_descriptor, 10, NULL);
		-- SET proto3_optional = pb_wire_json_get_bool_field(field_descriptor, 17, FALSE);
		-- SET oneof_index = pb_wire_json_get_int32_field(field_descriptor, 9, NULL);

		-- SET field_number = 0;
		-- SET field_name = NULL;
		-- SET field_label = 0;
		-- SET field_type = 0;
		-- SET field_type_name = NULL;
		-- SET default_value = NULL;
		-- SET json_name = NULL;
		-- SET proto3_optional = FALSE;
		-- SET oneof_index = NULL;
		-- CALL _pb_util_decode_field_descriptor2(field_descriptor, field_number, field_name, field_label, field_type, field_type_name, default_value, json_name, proto3_optional, oneof_index);

		SET is_repeated = field_label = 3; /* repeated */

		-- For proto2, all fields except repeated ones have field presence.
		SET has_field_presence =
			(syntax = 'proto2' AND field_label <> 3 /* repeated */)
			OR (syntax = 'proto3'
				AND (
					(field_label = 1 /* optional */ AND proto3_optional) -- This line is redundant because proto3 optional is a oneof.
					OR (field_label <> 3 /* repeated */ AND field_type = 11 /* TYPE_MESSAGE */)
					OR (oneof_index IS NOT NULL)
				));

		CASE field_type
		WHEN 10 THEN -- group
			SET message_text = CONCAT('_pb_message_to_json: unsupported field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		WHEN 11 THEN -- message
			-- TODO: use pb_wire_json_get_repeated_message_field_as_json_array
			IF is_map THEN
				SET elements = pb_wire_json_get_repeated_message_field_as_json_array(wire_json, field_number);
				SET field_count = JSON_LENGTH(elements);
				SET field_index = 0;
				SET field_json_value = JSON_OBJECT();
				WHILE field_index < field_count DO
					SET element = pb_message_to_wire_json(FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(elements, CONCAT('$[', field_index, ']')))));
					CALL _pb_wire_json_get_primitive_field_as_json(element, 1, map_key_type, FALSE, FALSE, map_key);
					IF map_value_type = 11 THEN -- message
						CALL _pb_message_to_json(set_name, map_value_type_name, pb_wire_json_get_message_field(element, 2, NULL), map_value);
					ELSEIF map_value_type = 14 THEN -- enum
						CALL _pb_enum_to_json(set_name, map_value_type_name, pb_wire_json_get_enum_field(element, 2, NULL), map_value);
					ELSE
						CALL _pb_wire_json_get_primitive_field_as_json(element, 2, map_value_type, FALSE, TRUE, map_value);
					END IF;
					IF JSON_TYPE(map_key) = 'STRING' THEN
						SET field_json_value = JSON_SET(field_json_value, CONCAT('$.', map_key), map_value);
					ELSE
						SET field_json_value = JSON_SET(field_json_value, CONCAT('$."', map_key, '"'), map_value);
					END IF;
					SET field_index = field_index + 1;
				END WHILE;
			ELSEIF is_repeated THEN
				SET field_count = pb_wire_json_get_repeated_message_field_count(wire_json, field_number);
				SET field_index = 0;
				SET field_json_value = JSON_ARRAY();
				WHILE field_index < field_count DO
					SET bytes_value = pb_wire_json_get_repeated_message_field(wire_json, field_number, field_index);
					CALL _pb_message_to_json(set_name, field_type_name, bytes_value, nested_json_value);
					SET field_json_value = JSON_ARRAY_APPEND(field_json_value, '$', nested_json_value);
					SET field_index = field_index + 1;
				END WHILE;
			ELSE
				SET bytes_value = pb_wire_json_get_message_field(wire_json, field_number, NULL); -- message fields always have field presence
				IF bytes_value IS NULL THEN
					SET field_json_value = NULL;
				ELSE
					CALL _pb_message_to_json(set_name, field_type_name, bytes_value, nested_json_value);
					SET field_json_value = nested_json_value;
				END IF;
			END IF;
		WHEN 14 THEN -- enum
			IF is_repeated THEN
				SET elements = pb_wire_json_get_repeated_enum_field_as_json_array(wire_json, field_number);
				SET field_count = JSON_LENGTH(elements);
				SET field_index = 0;
				SET field_json_value = JSON_ARRAY();
				WHILE field_index < field_count DO
					SET element = JSON_EXTRACT(elements, CONCAT('$[', field_index, ']'));
					CALL _pb_enum_to_json(set_name, field_type_name, element, nested_json_value);
					SET field_json_value = JSON_ARRAY_APPEND(field_json_value, '$', nested_json_value);
					SET field_index = field_index + 1;
				END WHILE;
			ELSE
				CALL _pb_enum_to_json(set_name, field_type_name, pb_wire_json_get_enum_field(wire_json, field_number, IF(has_field_presence, NULL, 0)), field_json_value);
			END IF;
		ELSE
			CALL _pb_wire_json_get_primitive_field_as_json(wire_json, field_number, field_type, is_repeated, has_field_presence, field_json_value);
		END CASE;

		IF field_json_value IS NOT NULL THEN
			SET json_field_name = IF(json_name IS NOT NULL, json_name, _pb_util_snake_to_lower_camel(field_name));
			IF oneof_index IS NOT NULL AND NOT proto3_optional THEN
				SET elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
				SET oneof_priority = JSON_EXTRACT(elements, CONCAT('$[', JSON_LENGTH(elements)-1, '].i'));
				SET oneof_priority_prev = JSON_EXTRACT(oneofs, CONCAT('$."', oneof_index, '".i'));
				IF oneof_priority_prev IS NULL OR oneof_priority_prev < oneof_priority THEN
					SET oneofs = JSON_SET(oneofs, CONCAT('$."', oneof_index, '"'), JSON_OBJECT('i', oneof_priority, 'v', JSON_OBJECT(json_field_name, field_json_value)));
				END IF;
			ELSE
				SET result = JSON_SET(result, CONCAT('$.', json_field_name), field_json_value);
			END IF;
		END IF;
	END LOOP;

	SET elements = JSON_EXTRACT(oneofs, '$.*.v');
	SET field_index = 0;
	SET field_count = JSON_LENGTH(elements);

	WHILE field_index < field_count DO
		SET field_json_value = JSON_EXTRACT(elements, CONCAT('$[', field_index, ']'));
		SET result = JSON_MERGE(result, field_json_value);
		SET field_index = field_index + 1;
	END WHILE;

	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_message_to_json $$
CREATE FUNCTION pb_message_to_json(set_name VARCHAR(64), full_type_name VARCHAR(512), buf LONGBLOB) RETURNS JSON READS SQL DATA
BEGIN
	DECLARE result JSON;
	CALL _pb_message_to_json(set_name, full_type_name, buf, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_util_decode_field_descriptor $$
CREATE PROCEDURE _pb_util_decode_field_descriptor(
	IN wire_json JSON,
	OUT out_field_number INT,
	OUT field_name TEXT,
	OUT field_label INT,
	OUT field_type INT,
	OUT field_type_name TEXT,
	OUT default_value TEXT,
	OUT json_name TEXT,
	OUT proto3_optional BOOLEAN,
	OUT oneof_index INT)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE field_number INT;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.field_number,
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	OPEN cur;
	l1: LOOP
		FETCH cur INTO field_number, wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 0 THEN
			CASE field_number
			WHEN 3 THEN
				SET out_field_number = _pb_util_reinterpret_uint64_as_int64(uint_value);
			WHEN 4 THEN
				SET field_label = _pb_util_reinterpret_uint64_as_int64(uint_value);
			WHEN 5 THEN
				SET field_type = _pb_util_reinterpret_uint64_as_int64(uint_value);
			WHEN 17 THEN
				SET proto3_optional = uint_value <> 0;
			WHEN 9 THEN
				SET oneof_index = _pb_util_reinterpret_uint64_as_int64(uint_value);
			END CASE;
		WHEN 2 THEN
			CASE field_number
			WHEN 1 THEN
				SET field_name = CONVERT(bytes_value USING utf8mb4);
			WHEN 6 THEN
				SET field_type_name = CONVERT(bytes_value USING utf8mb4);
			WHEN 7 THEN
				SET default_value = CONVERT(bytes_value USING utf8mb4);
			WHEN 10 THEN
				SET json_name = CONVERT(bytes_value USING utf8mb4);
			END CASE;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP PROCEDURE IF EXISTS _pb_util_decode_field_descriptor2 $$
CREATE PROCEDURE _pb_util_decode_field_descriptor2(
	IN wire_json JSON,
	OUT out_field_number INT,
	OUT field_name TEXT,
	OUT field_label INT,
	OUT field_type INT,
	OUT field_type_name TEXT,
	OUT default_value TEXT,
	OUT json_name TEXT,
	OUT proto3_optional BOOLEAN,
	OUT oneof_index INT)
BEGIN
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE field_number INT;
	DECLARE wire_type INT;
	DECLARE wire_entry JSON;
	DECLARE wire_json_index INT;
	DECLARE wire_json_length INT;
	DECLARE wire_elements JSON;

	SET wire_elements = JSON_EXTRACT(wire_json, '$.*[*]');

	SET wire_json_index = 0;
	SET wire_json_length = JSON_LENGTH(wire_elements);
	l1: WHILE wire_json_index < wire_json_length DO
		SET wire_entry = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_json_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_entry, '$.t');
		SET field_number = JSON_EXTRACT(wire_entry, '$.n');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_entry, '$.v') AS UNSIGNED);
			CASE field_number
			WHEN 3 THEN
				SET out_field_number = _pb_util_reinterpret_uint64_as_int64(uint_value);
			WHEN 4 THEN
				SET field_label = _pb_util_reinterpret_uint64_as_int64(uint_value);
			WHEN 5 THEN
				SET field_type = _pb_util_reinterpret_uint64_as_int64(uint_value);
			WHEN 17 THEN
				SET proto3_optional = uint_value <> 0;
			WHEN 9 THEN
				SET oneof_index = _pb_util_reinterpret_uint64_as_int64(uint_value);
			END CASE;
		WHEN 2 THEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_entry, '$.v')));
			CASE field_number
			WHEN 1 THEN
				SET field_name = CONVERT(bytes_value USING utf8mb4);
			WHEN 6 THEN
				SET field_type_name = CONVERT(bytes_value USING utf8mb4);
			WHEN 7 THEN
				SET default_value = CONVERT(bytes_value USING utf8mb4);
			WHEN 10 THEN
				SET json_name = CONVERT(bytes_value USING utf8mb4);
			END CASE;
		END CASE;
		SET wire_json_index = wire_json_index + 1;
	END WHILE;
END $$
