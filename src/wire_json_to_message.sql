DELIMITER $$

-- Public function: Convert wire JSON to protobuf binary message
DROP FUNCTION IF EXISTS pb_wire_json_to_message $$
CREATE FUNCTION pb_wire_json_to_message(wire_json JSON) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE result LONGBLOB;
	CALL _pb_wire_json_to_message(wire_json, result);
	RETURN result;
END $$

-- Main function: Convert wire JSON to protobuf binary message
DROP PROCEDURE IF EXISTS _pb_wire_json_to_message $$
CREATE PROCEDURE _pb_wire_json_to_message(IN wire_json JSON, OUT message LONGBLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE done INT DEFAULT FALSE;
	DECLARE element_index INT;
	DECLARE field_number INT;
	DECLARE wire_type INT;
	DECLARE v_value JSON;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE tag_encoded LONGBLOB;
	DECLARE value_encoded LONGBLOB;
	DECLARE message_text TEXT;

	-- Cursor to read elements sorted by index 'i' using JSON_TABLE
	DECLARE element_cursor CURSOR FOR
		SELECT i, n, t, v
		FROM JSON_TABLE(
			JSON_EXTRACT(wire_json, '$.*[*]'),
			'$[*]' COLUMNS (
				i INT PATH '$.i',
				n INT PATH '$.n',
				t INT PATH '$.t',
				v JSON PATH '$.v'
			)
		) jt
		ORDER BY i;

	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET message = _binary '';

	OPEN element_cursor;

	read_loop: LOOP
		FETCH element_cursor INTO element_index, field_number, wire_type, v_value;
		IF done THEN
			LEAVE read_loop;
		END IF;

		-- Write tag (field number + wire type)
		CALL _pb_wire_write_tag(field_number, wire_type, tag_encoded);
		SET message = CONCAT(message, tag_encoded);

		-- Write value based on wire type
		CASE wire_type
		WHEN 0 THEN -- VARINT
			SET uint_value = CAST(v_value AS UNSIGNED);
			CALL _pb_wire_write_varint(uint_value, value_encoded);
			SET message = CONCAT(message, value_encoded);
		WHEN 1 THEN -- I64
			SET uint_value = CAST(v_value AS UNSIGNED);
			CALL _pb_wire_write_i64(uint_value, value_encoded);
			SET message = CONCAT(message, value_encoded);
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(v_value));
			CALL _pb_wire_write_len_type(bytes_value, value_encoded);
			SET message = CONCAT(message, value_encoded);
		WHEN 5 THEN -- I32
			SET uint_value = CAST(v_value AS UNSIGNED);
			CALL _pb_wire_write_i32(uint_value, value_encoded);
			SET message = CONCAT(message, value_encoded);
		ELSE
			SET message_text = CONCAT('_pb_wire_json_to_message: unsupported wire type (', wire_type, ')');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;

	CLOSE element_cursor;
END $$
