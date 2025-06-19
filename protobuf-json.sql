DELIMITER $$

DROP PROCEDURE IF EXISTS _pb_message_to_json $$
CREATE PROCEDURE _pb_message_to_json(IN set_name VARCHAR(64), IN full_type_name VARCHAR(512), IN buf LONGBLOB, OUT result JSON)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE field_number INT;
	DECLARE field_name TEXT;
	DECLARE field_type INT;
	DECLARE field_type_name TEXT;
	DECLARE json_name TEXT;
	DECLARE proto3_optional BOOLEAN;

	DECLARE bytes_value LONGBLOB;
	DECLARE nested_json_value JSON;

	-- NOTE: always use alias in select columns, to avoid confusion with variables with the same name.
	DECLARE cur CURSOR FOR
		SELECT
			t.field_number,
			t.field_name,
			t.field_type,
			t.field_type_name,
			t.json_name,
			t.proto3_optional
		FROM _Proto_FieldDescriptor t
		WHERE
			t.set_name = set_name
			AND t.type_name = full_type_name;

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

	OPEN cur;

	l1: LOOP
		FETCH cur INTO field_number, field_name, field_type, field_type_name, json_name, proto3_optional;

		IF done THEN
			LEAVE l1;
		END IF;

		CASE field_type
			WHEN 1 THEN -- double
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_double_field(buf, field_number, NULL));
			WHEN 2 THEN -- float
				SET result = JSON_SET(result, CONCAT('$.', field_name), 
					pb_message_get_double_field(buf, field_number, NULL));
			WHEN 3 THEN -- int64
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_int64_field(buf, field_number, NULL));
			WHEN 4 THEN -- uint64
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_uint64_field(buf, field_number, NULL));
			WHEN 5 THEN -- int32
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_int32_field(buf, field_number, NULL));
			WHEN 6 THEN -- fixed64
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_fixed64_field(buf, field_number, NULL));
			WHEN 7 THEN -- fixed32
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_fixed32_field(buf, field_number, NULL));
			WHEN 8 THEN -- bool
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_bool_field(buf, field_number, NULL));
			WHEN 9 THEN -- string
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_string_field(buf, field_number, NULL));
			WHEN 10 THEN -- group
				SET message_text = CONCAT('_pb_message_to_json: unsupported field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
				SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
			WHEN 11 THEN -- message
				SET bytes_value = pb_message_get_message_field(buf, field_number, NULL);
				CALL _pb_message_to_json(set_name, field_type_name, bytes_value, nested_json_value);
				SET result = JSON_SET(result, CONCAT('$.', field_name), nested_json_value);
			WHEN 12 THEN -- bytes
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					TO_BASE64(pb_message_get_bytes_field(buf, field_number, NULL)));
			WHEN 13 THEN -- uint32
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_uint32_field(buf, field_number, NULL));
			WHEN 14 THEN -- enum
				-- TODO: convert to enum name
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_enum_field(buf, field_number, NULL));
			WHEN 15 THEN -- sfixed32
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_sfixed32_field(buf, field_number, NULL));
			WHEN 16 THEN -- sfixed64
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_sfixed64_field(buf, field_number, NULL));
			WHEN 17 THEN -- sint32
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_sint32_field(buf, field_number, NULL));
			WHEN 18 THEN -- sint64
				SET result = JSON_SET(result, CONCAT('$.', field_name),
					pb_message_get_sint64_field(buf, field_number, NULL));
			ELSE
				SET message_text = CONCAT('_pb_message_to_json: unknown field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
				SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END CASE;

	END LOOP;

	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_message_to_json $$
CREATE FUNCTION pb_message_to_json(set_name VARCHAR(64), full_type_name VARCHAR(512), buf LONGBLOB) RETURNS JSON READS SQL DATA
BEGIN
	DECLARE result JSON;
	CALL _pb_message_to_json(set_name, full_type_name, buf, result);
	RETURN result;
END $$
