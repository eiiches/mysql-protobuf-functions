DELIMITER $$

DROP FUNCTION IF EXISTS pb_message_to_wire_json $$
CREATE FUNCTION pb_message_to_wire_json(buf LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE wire_json JSON;
	IF buf IS NULL THEN
		RETURN NULL;
	END IF;
	CALL _pb_message_to_wire_json(buf, NULL, wire_json, buf);
	RETURN wire_json;
END $$

DROP PROCEDURE IF EXISTS _pb_message_to_wire_json $$
CREATE PROCEDURE _pb_message_to_wire_json(IN buf LONGBLOB, IN expected_egroup_field_number INT, OUT wire_json JSON, OUT tail LONGBLOB)
proc: BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag INT;
	DECLARE field_number INT;
	DECLARE wire_type INT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE i INT;
	DECLARE json_path TEXT;
	DECLARE wire_element JSON;
	DECLARE nested_wire_json JSON;

	SET wire_json = JSON_OBJECT();
	SET tail = buf;
	SET i = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);

		SET field_number = _pb_wire_get_field_number_from_tag(tag);
		SET wire_type = _pb_wire_get_wire_type_from_tag(tag);

		-- Validate field number is non-zero (protobuf spec requirement)
		IF field_number = 0 THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '_pb_message_to_wire_json: Invalid protobuf field number 0 is not allowed';
		END IF;

		CASE wire_type
		WHEN 0 THEN -- VARINT
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET wire_element = JSON_OBJECT('i', i, 'n', field_number, 't', wire_type, 'v', uint_value);
		WHEN 1 THEN -- I64
			CALL _pb_wire_read_i64_as_uint64(tail, uint_value, tail);
			SET wire_element = JSON_OBJECT('i', i, 'n', field_number, 't', wire_type, 'v', uint_value);
		WHEN 2 THEN -- LEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			SET wire_element = JSON_OBJECT('i', i, 'n', field_number, 't', wire_type, 'v', _pb_util_to_base64(bytes_value));
		WHEN 3 THEN -- SGROUP
			CALL _pb_message_to_wire_json(tail, field_number, nested_wire_json, tail);
			SET wire_element = JSON_OBJECT('i', i, 'n', field_number, 't', wire_type, 'v', nested_wire_json);
		WHEN 4 THEN -- EGROUP
			IF expected_egroup_field_number IS NULL THEN
				SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '_pb_message_to_wire_json: unexpected EGROUP found outside a group';
			END IF;
			IF expected_egroup_field_number <> field_number THEN
				SET message_text = CONCAT('_pb_message_to_wire_json: expected EGROUP field number (', expected_egroup_field_number, ') but got ', field_number);
				SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
			END IF;
			LEAVE proc;
		WHEN 5 THEN -- I32
			CALL _pb_wire_read_i32_as_uint32(tail, uint_value, tail);
			SET wire_element = JSON_OBJECT('i', i, 'n', field_number, 't', wire_type, 'v', uint_value);
		ELSE
			SET message_text = CONCAT('_pb_message_to_wire_json: unsupported wire type (', wire_type, ')');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END CASE;

		SET json_path = CONCAT('$.', JSON_QUOTE(CAST(field_number AS CHAR)));
		IF JSON_CONTAINS_PATH(wire_json, 'one', json_path) THEN
			SET wire_json = JSON_ARRAY_APPEND(wire_json, json_path, wire_element);
		ELSE
			SET wire_json = JSON_SET(wire_json, json_path, JSON_ARRAY(wire_element));
		END IF;

		SET i = i + 1;
	END WHILE;

	IF expected_egroup_field_number IS NOT NULL THEN
		SET message_text = CONCAT('_pb_message_to_wire_json: expected EGROUP field number (', expected_egroup_field_number, ') but reached EOF');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;
END $$
