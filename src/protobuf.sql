DELIMITER $$

DROP PROCEDURE IF EXISTS _pb_wire_read_i32_as_uint32 $$
CREATE PROCEDURE _pb_wire_read_i32_as_uint32(IN buf LONGBLOB, OUT value INT UNSIGNED, OUT tail LONGBLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(buf) < 4 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_read_i32_as_uint32: Unexpected end of BLOB.';
	END IF;

	SET value = _pb_util_swap_endian_32(_pb_util_bin_as_uint32(LEFT(buf, 4)));
	SET tail = SUBSTRING(buf, 5);
END $$

DROP PROCEDURE IF EXISTS _pb_wire_read_i64_as_uint64 $$
CREATE PROCEDURE _pb_wire_read_i64_as_uint64(IN buf LONGBLOB, OUT value BIGINT UNSIGNED, OUT tail LONGBLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(buf) < 8 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_read_i64_as_uint64: Unexpected end of BLOB.';
	END IF;

	SET value = _pb_util_swap_endian_64(_pb_util_bin_as_uint64(LEFT(buf, 8)));
	SET tail = SUBSTRING(buf, 9);
END $$

DROP PROCEDURE IF EXISTS _pb_wire_skip_i32 $$
CREATE PROCEDURE _pb_wire_skip_i32(IN buf LONGBLOB, OUT tail LONGBLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(buf) < 4 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_skip_i32: Unexpected end of BLOB.';
	END IF;

	SET tail = SUBSTRING(buf, 5);
END $$

DROP PROCEDURE IF EXISTS _pb_wire_skip_i64 $$
CREATE PROCEDURE _pb_wire_skip_i64(IN buf LONGBLOB, OUT tail LONGBLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(buf) < 8 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_skip_i64: Unexpected end of BLOB.';
	END IF;

	SET tail = SUBSTRING(buf, 9);
END $$

DROP PROCEDURE IF EXISTS _pb_wire_skip_len_type $$
CREATE PROCEDURE _pb_wire_skip_len_type(IN buf LONGBLOB, OUT tail LONGBLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE len BIGINT;

	SET tail = buf;
	CALL _pb_wire_read_varint_as_uint64(tail, len, tail);

	IF LENGTH(tail) < len THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_skip_len_type: Unexpected end of BLOB.';
	END IF;

	SET tail = SUBSTRING(tail, len + 1);
END $$

DROP PROCEDURE IF EXISTS _pb_wire_read_len_type $$
CREATE PROCEDURE _pb_wire_read_len_type(IN buf LONGBLOB, OUT value LONGBLOB, OUT tail LONGBLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE len BIGINT;

	SET tail = buf;
	CALL _pb_wire_read_varint_as_uint64(tail, len, tail);

	IF LENGTH(tail) < len THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_read_len_type: Unexpected end of BLOB.';
	END IF;

	SET value = LEFT(tail, len);
	SET tail = SUBSTRING(tail, len + 1);
END $$

DROP PROCEDURE IF EXISTS _pb_wire_skip_varint $$
CREATE PROCEDURE _pb_wire_skip_varint(IN buf LONGBLOB, OUT tail LONGBLOB)
BEGIN
	DECLARE head INT;
	DECLARE byte_index INT;
	DECLARE buf_len INT;

	SET byte_index = 0;
	SET buf_len = LENGTH(buf);

	l1: LOOP
		IF byte_index >= buf_len THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '_pb_wire_skip_varint: Unexpected end of BLOB.';
		END IF;

		SET head = ORD(SUBSTRING(buf, byte_index + 1, 1));
		SET byte_index = byte_index + 1;

		IF (head & 0x80) = 0 THEN
			LEAVE l1;
		END IF;

		IF byte_index > 10 THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '_pb_wire_skip_varint: Varint cannot exceed 10 bytes.';
		END IF;
	END LOOP;

	SET tail = SUBSTRING(buf, byte_index + 1);
END $$

DROP PROCEDURE IF EXISTS _pb_wire_skip $$
CREATE PROCEDURE _pb_wire_skip(IN buf LONGBLOB, IN wire_type INT, OUT tail LONGBLOB)
BEGIN
	DECLARE dummy BIGINT UNSIGNED;
	DECLARE message_text TEXT;

	CASE wire_type
	WHEN 0 THEN -- VARINT
		CALL _pb_wire_skip_varint(buf, tail);
	WHEN 1 THEN -- I64
		CALL _pb_wire_skip_i64(buf, tail);
	WHEN 2 THEN -- LEN
		CALL _pb_wire_skip_len_type(buf, tail);
	WHEN 5 THEN -- I32
		CALL _pb_wire_skip_i32(buf, tail);
	ELSE
		SET message_text = CONCAT('_pb_wire_skip: unknown wire_type (', wire_type, ')');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_read_varint_as_uint64 $$
CREATE PROCEDURE _pb_wire_read_varint_as_uint64(IN buf LONGBLOB, OUT value BIGINT UNSIGNED, OUT tail LONGBLOB)
BEGIN
	DECLARE head INT;
	DECLARE byte_index INT;
	DECLARE buf_len INT;

	SET value = 0;
	SET byte_index = 0;
	SET buf_len = LENGTH(buf);

	l1: LOOP
		IF byte_index >= buf_len THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '_pb_wire_read_varint_as_uint64: Unexpected end of BLOB.';
		END IF;

		SET head = ORD(SUBSTRING(buf, byte_index + 1, 1));
		SET value = value + ((head & 0x7f) << (7 * byte_index));
		SET byte_index = byte_index + 1;

		IF (head & 0x80) = 0 THEN
			LEAVE l1;
		END IF;

		IF byte_index > 10 THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '_pb_wire_read_varint_as_uint64: Varint cannot exceed 10 bytes.';
		END IF;
	END LOOP;

	SET tail = SUBSTRING(buf, byte_index + 1);
END $$

DROP FUNCTION IF EXISTS _pb_wire_read_varint_as_uint64 $$
CREATE FUNCTION _pb_wire_read_varint_as_uint64(buf LONGBLOB) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE tail LONGBLOB;
	DECLARE value BIGINT;
	CALL _pb_wire_read_varint_as_uint64(buf, value, tail);
	RETURN value;
END $$


DROP FUNCTION IF EXISTS _pb_wire_get_field_number_from_tag $$
CREATE FUNCTION _pb_wire_get_field_number_from_tag(tag INT) RETURNS INT DETERMINISTIC
BEGIN
	RETURN tag >> 3;
END $$

DROP FUNCTION IF EXISTS _pb_wire_get_wire_type_from_tag $$
CREATE FUNCTION _pb_wire_get_wire_type_from_tag(tag INT) RETURNS INT DETERMINISTIC
BEGIN
	RETURN tag & 0b111;
END $$

DROP FUNCTION IF EXISTS _pb_wire_type_name $$
CREATE FUNCTION _pb_wire_type_name(num INT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';
	CASE num
	WHEN 0 THEN RETURN 'VARINT';
	WHEN 1 THEN RETURN 'I64';
	WHEN 2 THEN RETURN 'LEN';
	WHEN 3 THEN RETURN 'SGROUP';
	WHEN 4 THEN RETURN 'EGROUP';
	WHEN 5 THEN RETURN 'I32';
	ELSE
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_type_name: unsupported wire_type';
	END CASE;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_len_type_field $$
CREATE PROCEDURE _pb_message_get_len_type_field(IN buf LONGBLOB, IN field_number INT, IN repeated_index INT, OUT value LONGBLOB, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;

	SET tail = buf;
	SET field_count = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 2 /* VARINT */ THEN
			SET message_text = CONCAT('_pb_message_get_len_type_field: string or bytes value cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
		WHEN 1 THEN -- I64
			CALL _pb_wire_skip_i64(tail, tail);
		WHEN 2 THEN -- LEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				IF repeated_index IS NULL OR repeated_index = field_count THEN
					SET value = bytes_value;
				END IF;
				SET field_count = field_count + 1;
			END IF;
		WHEN 5 THEN -- I32
			CALL _pb_wire_skip_i32(tail, tail);
		ELSE
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_len_type_field: unsupported wire_type';
		END CASE;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_len_type_field: repeated index out of range';
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_i32_field_as_uint32 $$
CREATE PROCEDURE _pb_message_get_i32_field_as_uint32(IN buf LONGBLOB, IN field_number INT, IN repeated_index INT, OUT value INT UNSIGNED, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE packed_value LONGBLOB;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;

	SET tail = buf;
	SET field_count = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 5 /* I32 */ AND (repeated_index IS NULL OR _pb_wire_get_wire_type_from_tag(tag) <> 2 /* LEN */) THEN
			SET message_text = CONCAT('_pb_message_get_i32_field_as_uint32: I32 value cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
		WHEN 1 THEN -- I64
			CALL _pb_wire_skip_i64(tail, tail);
		WHEN 2 THEN -- LEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number AND repeated_index IS NOT NULL THEN
				WHILE LENGTH(bytes_value) <> 0 DO
					CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
					IF repeated_index = field_count THEN
						SET value = uint_value;
					END IF;
					SET field_count = field_count + 1;
				END WHILE;
			END IF;
		WHEN 5 THEN -- I32
			CALL _pb_wire_read_i32_as_uint32(tail, uint_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				IF repeated_index IS NULL OR repeated_index = field_count THEN
					SET value = uint_value;
				END IF;
				SET field_count = field_count + 1;
			END IF;
		ELSE
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_i32_field_as_uint32: unsupported wire_type';
		END CASE;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_i32_field_as_uint32: repeated index out of range';
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_i64_field_as_uint64 $$
CREATE PROCEDURE _pb_message_get_i64_field_as_uint64(IN buf LONGBLOB, IN field_number INT, IN repeated_index INT, OUT value BIGINT UNSIGNED, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE packed_value LONGBLOB;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;

	SET tail = buf;
	SET field_count = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 1 /* I64 */ AND (repeated_index IS NULL OR _pb_wire_get_wire_type_from_tag(tag) <> 2 /* LEN */) THEN
			SET message_text = CONCAT('_pb_message_get_i64_field_as_uint64: I64 value cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
		WHEN 1 THEN -- I64
			CALL _pb_wire_read_i64_as_uint64(tail, uint_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				IF repeated_index IS NULL OR repeated_index = field_count THEN
					SET value = uint_value;
				END IF;
				SET field_count = field_count + 1;
			END IF;
		WHEN 2 THEN -- LEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number AND repeated_index IS NOT NULL THEN
				WHILE LENGTH(bytes_value) <> 0 DO
					CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
					IF repeated_index = field_count THEN
						SET value = uint_value;
					END IF;
					SET field_count = field_count + 1;
				END WHILE;
			END IF;
		WHEN 5 THEN -- I32
			CALL _pb_wire_skip_i32(tail, tail);
		ELSE
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_i64_field_as_uint64: unsupported wire_type';
		END CASE;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_i64_field_as_uint64: repeated index out of range';
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_varint_field_as_uint64 $$
CREATE PROCEDURE _pb_message_get_varint_field_as_uint64(IN buf LONGBLOB, IN field_number INT, IN repeated_index INT, OUT value BIGINT UNSIGNED, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = buf;
	SET field_count = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number = field_number AND current_wire_type <> 0 /* VARINT */ AND (repeated_index IS NULL OR current_wire_type <> 2 /* LEN */) THEN
			SET message_text = CONCAT('_pb_message_get_varint_field_as_uint64: uint32 or uint64 value cannot be parsed from ', _pb_wire_type_name(current_wire_type), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN -- VARINT
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			IF current_field_number = field_number THEN
				IF repeated_index IS NULL OR repeated_index = field_count THEN
					SET value = uint_value;
				END IF;
				SET field_count = field_count + 1;
			END IF;
		WHEN 1 THEN -- I64
			CALL _pb_wire_skip_i64(tail, tail);
		WHEN 2 THEN -- LEN
			IF current_field_number = field_number AND repeated_index IS NOT NULL THEN
				CALL _pb_wire_read_len_type(tail, bytes_value, tail);
				WHILE LENGTH(bytes_value) <> 0 DO
					CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
					IF repeated_index = field_count THEN
						SET value = uint_value;
					END IF;
					SET field_count = field_count + 1;
				END WHILE;
			ELSE
				CALL _pb_wire_skip_len_type(tail, tail);
			END IF;
		WHEN 5 THEN -- I32
			CALL _pb_wire_skip_i32(tail, tail);
		ELSE
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_varint_field_as_uint64: unsupported wire_type';
		END CASE;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_varint_field_as_uint64: repeated index out of range';
	END IF;
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

		SET json_path = CONCAT('$."', field_number, '"');
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

DROP FUNCTION IF EXISTS pb_message_new $$
CREATE FUNCTION pb_message_new() RETURNS LONGBLOB DETERMINISTIC
BEGIN
	-- Return an empty protobuf message
	RETURN _binary X'';
END $$

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

DROP PROCEDURE IF EXISTS pb_wire_json_as_table $$
CREATE PROCEDURE pb_wire_json_as_table(IN wire_json JSON)
BEGIN
	SELECT * FROM JSON_TABLE(JSON_EXTRACT(wire_json, '$.*[*]'), '$[*]' COLUMNS (i INT PATH '$.i', n INT PATH '$.n', t INT PATH '$.t', v JSON PATH '$.v')) jt;
END $$

DROP PROCEDURE IF EXISTS pb_message_show_wire_format $$
CREATE PROCEDURE pb_message_show_wire_format(IN buf LONGBLOB)
BEGIN
	CALL pb_wire_json_as_table(pb_message_to_wire_json(buf));
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_varint_field_as_uint64 $$
CREATE PROCEDURE _pb_wire_json_get_varint_field_as_uint64(IN wire_json JSON, IN field_number INT, IN repeated_index INT, OUT value BIGINT UNSIGNED, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET field_count = 0;

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN -- VARINT
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			IF repeated_index IS NULL OR repeated_index = field_count THEN
				SET value = uint_value;
			END IF;
			SET field_count = field_count + 1;
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			IF repeated_index IS NULL THEN
				SET message_text = CONCAT('_pb_wire_json_get_varint_field_as_uint64: unexpected wire_type (', wire_type, ')');
				SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
			END IF;
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				IF repeated_index = field_count THEN
					SET value = uint_value;
				END IF;
				SET field_count = field_count + 1;
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_varint_field_as_uint64: unexpected wire_type (', wire_type, ')');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_json_get_varint_field_as_uint64: repeated index out of range';
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_i64_field_as_uint64 $$
CREATE PROCEDURE _pb_wire_json_get_i64_field_as_uint64(IN wire_json JSON, IN field_number INT, IN repeated_index INT, OUT value BIGINT UNSIGNED, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET field_count = 0;

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 1 THEN -- I64
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			IF repeated_index IS NULL OR repeated_index = field_count THEN
				SET value = uint_value;
			END IF;
			SET field_count = field_count + 1;
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			IF repeated_index IS NULL THEN
				SET message_text = CONCAT('_pb_wire_json_get_i64_field_as_uint64: unexpected wire_type (', wire_type, ')');
				SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
			END IF;
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				IF repeated_index = field_count THEN
					SET value = uint_value;
				END IF;
				SET field_count = field_count + 1;
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_i64_field_as_uint64: unexpected wire_type (', wire_type, ')');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_json_get_i64_field_as_uint64: repeated index out of range';
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_i32_field_as_uint32 $$
CREATE PROCEDURE _pb_wire_json_get_i32_field_as_uint32(IN wire_json JSON, IN field_number INT, IN repeated_index INT, OUT value INT UNSIGNED, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET field_count = 0;

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 5 THEN -- I32
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			IF repeated_index IS NULL OR repeated_index = field_count THEN
				SET value = uint_value;
			END IF;
			SET field_count = field_count + 1;
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			IF repeated_index IS NULL THEN
				SET message_text = CONCAT('_pb_wire_json_get_i32_field_as_uint32: unexpected wire_type (', wire_type, ')');
				SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
			END IF;
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				IF repeated_index = field_count THEN
					SET value = uint_value;
				END IF;
				SET field_count = field_count + 1;
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_i32_field_as_uint32: unexpected wire_type (', wire_type, ')');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_json_get_i32_field_as_uint32: repeated index out of range';
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_len_type_field $$
CREATE PROCEDURE _pb_wire_json_get_len_type_field(IN wire_json JSON, IN field_number INT, IN repeated_index INT, OUT value LONGBLOB, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET field_count = 0;

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			IF repeated_index IS NULL OR repeated_index = field_count THEN
				SET value = bytes_value;
			END IF;
			SET field_count = field_count + 1;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_len_type_field: unexpected wire_type (', wire_type, ')');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_wire_json_get_len_type_field: repeated index out of range';
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_len_type_field_concatenated $$
CREATE PROCEDURE _pb_wire_json_get_len_type_field_concatenated(IN wire_json JSON, IN field_number INT, IN repeated_index INT, OUT concatenated_value LONGBLOB, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET field_count = 0;
	SET concatenated_value = _binary X'';

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			SET concatenated_value = CONCAT(concatenated_value, bytes_value);
			SET field_count = field_count + 1;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_len_type_field_concatenated: unexpected wire_type (', wire_type, ')');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_len_type_field_concatenated $$
CREATE PROCEDURE _pb_message_get_len_type_field_concatenated(IN buf LONGBLOB, IN field_number INT, IN repeated_index INT, OUT concatenated_value LONGBLOB, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;

	SET tail = buf;
	SET field_count = 0;
	SET concatenated_value = _binary X'';

	WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 2 /* LEN */ THEN
			SET message_text = CONCAT('_pb_message_get_len_type_field_concatenated: message field cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
		WHEN 1 THEN -- I64
			CALL _pb_wire_skip_i64(tail, tail);
		WHEN 2 THEN -- LEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				SET concatenated_value = CONCAT(concatenated_value, bytes_value);
				SET field_count = field_count + 1;
			END IF;
		WHEN 5 THEN -- I32
			CALL _pb_wire_skip_i32(tail, tail);
		ELSE
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_len_type_field_concatenated: unsupported wire_type';
		END CASE;
	END WHILE;
END $$

-- =============================================================================
-- Wire Encoding Functions (for converting back to protobuf binary format)
-- =============================================================================

-- Encode varint (unsigned integer with variable length encoding)
DROP PROCEDURE IF EXISTS _pb_wire_write_varint $$
CREATE PROCEDURE _pb_wire_write_varint(IN value BIGINT UNSIGNED, OUT encoded LONGBLOB)
BEGIN
	DECLARE result LONGBLOB DEFAULT _binary '';
	WHILE value >= 0x80 DO
		SET result = CONCAT(result, CHAR((value & 0x7F) | 0x80));
		SET value = value >> 7;
	END WHILE;
	SET result = CONCAT(result, CHAR(value & 0x7F));
	SET encoded = result;
END $$

-- Write tag (field number + wire type)
DROP PROCEDURE IF EXISTS _pb_wire_write_tag $$
CREATE PROCEDURE _pb_wire_write_tag(IN field_number INT, IN wire_type INT, OUT encoded LONGBLOB)
BEGIN
	DECLARE tag BIGINT UNSIGNED DEFAULT (field_number << 3) | wire_type;

	-- Validate field number is non-zero (protobuf spec requirement)
	IF field_number = 0 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '_pb_wire_write_tag: Invalid protobuf field number 0 is not allowed';
	END IF;

	CALL _pb_wire_write_varint(tag, encoded);
END $$

-- Write 32-bit integer (little-endian)
DROP PROCEDURE IF EXISTS _pb_wire_write_i32 $$
CREATE PROCEDURE _pb_wire_write_i32(IN value INT UNSIGNED, OUT encoded LONGBLOB)
BEGIN
	SET encoded = _pb_util_uint32_to_bin(_pb_util_swap_endian_32(value));
END $$

-- Write 64-bit integer (little-endian)
DROP PROCEDURE IF EXISTS _pb_wire_write_i64 $$
CREATE PROCEDURE _pb_wire_write_i64(IN value BIGINT UNSIGNED, OUT encoded LONGBLOB)
BEGIN
	SET encoded = _pb_util_uint64_to_bin(_pb_util_swap_endian_64(value));
END $$

-- Write length-delimited data (length prefix + data)
DROP PROCEDURE IF EXISTS _pb_wire_write_len_type $$
CREATE PROCEDURE _pb_wire_write_len_type(IN data LONGBLOB, OUT encoded LONGBLOB)
BEGIN
	DECLARE length_encoded LONGBLOB;
	CALL _pb_wire_write_varint(LENGTH(data), length_encoded);
	SET encoded = CONCAT(length_encoded, data);
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

DROP FUNCTION IF EXISTS pb_wire_json_new $$
CREATE FUNCTION pb_wire_json_new() RETURNS JSON DETERMINISTIC
BEGIN
	-- Return an empty wire JSON object
	RETURN JSON_OBJECT();
END $$

-- Public function: Convert wire JSON to protobuf binary message
DROP FUNCTION IF EXISTS pb_wire_json_to_message $$
CREATE FUNCTION pb_wire_json_to_message(wire_json JSON) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE result LONGBLOB;
	CALL _pb_wire_json_to_message(wire_json, result);
	RETURN result;
END $$

-- =============================================================================
-- Wire JSON Setter Functions (for modifying wire JSON fields)
-- =============================================================================

-- Helper function to get the next available index for wire JSON elements
DROP FUNCTION IF EXISTS _pb_wire_json_get_next_index $$
CREATE FUNCTION _pb_wire_json_get_next_index(wire_json JSON) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE max_index INT;

	-- Use the same JSON_TABLE pattern as pb_wire_json_as_table to get max index
	SELECT COALESCE(MAX(i), -1) INTO max_index
	FROM JSON_TABLE(
		JSON_EXTRACT(wire_json, '$.*[*]'),
		'$[*]' COLUMNS (
			i INT PATH '$.i'
		)
	) jt;

	RETURN max_index + 1;
END $$

-- Helper function to set a field in wire JSON with proper wire type
DROP FUNCTION IF EXISTS _pb_wire_json_set_field $$
CREATE FUNCTION _pb_wire_json_set_field(
	wire_json JSON,
	field_number INT,
	wire_type INT,
	value JSON
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE next_index INT;
	DECLARE new_element JSON;
	DECLARE message_text TEXT;

	IF value IS NULL OR field_number IS NULL OR wire_type IS NULL THEN
		SET message_text = CONCAT('_pb_wire_json_set_field: invalid null parameter');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Get the next available index
	SET next_index = _pb_wire_json_get_next_index(wire_json);

	-- Create the new wire element
	SET new_element = JSON_OBJECT('i', next_index, 'n', field_number, 't', wire_type, 'v', value);

	-- Replace the entire field
	RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
END $$

-- Private: Add to repeated field (appends new element)
DROP FUNCTION IF EXISTS _pb_wire_json_add_repeated_field_element $$
CREATE FUNCTION _pb_wire_json_add_repeated_field_element(
	wire_json JSON,
	field_number INT,
	wire_type INT,
	value JSON
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE next_index INT;
	DECLARE new_element JSON;

	-- Get the next available index
	SET next_index = _pb_wire_json_get_next_index(wire_json);

	-- Create the new wire element
	SET new_element = JSON_OBJECT('i', next_index, 'n', field_number, 't', wire_type, 'v', value);

	-- If field doesn't exist, create it; otherwise append
	IF NOT JSON_CONTAINS_PATH(wire_json, 'one', field_path) THEN
		RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
	ELSE
		RETURN JSON_ARRAY_APPEND(wire_json, field_path, new_element);
	END IF;
END $$

-- Private: Clear field (removes the entire field from wire JSON)
DROP FUNCTION IF EXISTS _pb_wire_json_clear_field $$
CREATE FUNCTION _pb_wire_json_clear_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	RETURN JSON_REMOVE(wire_json, field_path);
END $$

-- Private: Set repeated VARINT field element at specific index
DROP FUNCTION IF EXISTS _pb_wire_json_set_repeated_varint_field_element $$
CREATE FUNCTION _pb_wire_json_set_repeated_varint_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT,
	value BIGINT UNSIGNED
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE current_index INT DEFAULT 0;
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE original_index INT;
	DECLARE new_element JSON;
	DECLARE element_path TEXT;
	-- Variables for packed field handling
	DECLARE packed_data LONGBLOB;
	DECLARE new_packed_data LONGBLOB DEFAULT '';
	DECLARE temp_value BIGINT UNSIGNED;
	DECLARE temp_encoded LONGBLOB;
	DECLARE found_target BOOLEAN DEFAULT FALSE;
	DECLARE message_text TEXT;

	-- Get the field array
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Find the target index across all array elements
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		CASE element_wire_type
		WHEN 2 THEN
			-- Packed field - decode, modify, and re-encode
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			SET new_packed_data = '';
			SET found_target = FALSE;

			-- Decode packed varints and rebuild with modification
			WHILE LENGTH(packed_data) > 0 DO
				CALL _pb_wire_read_varint_as_uint64(packed_data, temp_value, packed_data);

				IF current_index = repeated_index THEN
					-- Replace this element with new value
					CALL _pb_wire_write_varint(value, temp_encoded);
					SET found_target = TRUE;
				ELSE
					-- Keep original element
					CALL _pb_wire_write_varint(temp_value, temp_encoded);
				END IF;

				SET new_packed_data = CONCAT(new_packed_data, temp_encoded);
				SET current_index = current_index + 1;
			END WHILE;

			-- If we found the target, update and return
			IF found_target THEN
				SET original_index = JSON_EXTRACT(element, '$.i');
				SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(new_packed_data));
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				RETURN JSON_SET(wire_json, element_path, new_element);
			END IF;
		WHEN 0 THEN
			-- Non-packed field (wire type 0 for VARINT)
			IF current_index = repeated_index THEN
				-- Found the target index - preserve original index and replace value
				SET original_index = JSON_EXTRACT(element, '$.i');
				SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 0, 'v', CAST(value AS JSON));
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				RETURN JSON_SET(wire_json, element_path, new_element);
			END IF;
			SET current_index = current_index + 1;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_set_repeated_varint_field_element: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET array_index = array_index + 1;
	END WHILE;

	-- Index not found
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
END $$

-- Private: Set repeated I64 field element at specific index
DROP FUNCTION IF EXISTS _pb_wire_json_set_repeated_i64_field_element $$
CREATE FUNCTION _pb_wire_json_set_repeated_i64_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT,
	value BIGINT UNSIGNED
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE current_index INT DEFAULT 0;
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE original_index INT;
	DECLARE new_element JSON;
	DECLARE element_path TEXT;
	-- Variables for packed field handling
	DECLARE packed_data LONGBLOB;
	DECLARE new_packed_data LONGBLOB DEFAULT '';
	DECLARE temp_value BIGINT UNSIGNED;
	DECLARE temp_encoded LONGBLOB;
	DECLARE found_target BOOLEAN DEFAULT FALSE;
	DECLARE message_text TEXT;

	-- Get the field array
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Find the target index across all array elements
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		CASE element_wire_type
		WHEN 2 THEN
			-- Packed field - decode, modify, and re-encode
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			SET new_packed_data = '';
			SET found_target = FALSE;

			-- Decode packed I64 values and rebuild with modification
			WHILE LENGTH(packed_data) > 0 DO
				CALL _pb_wire_read_i64_as_uint64(packed_data, temp_value, packed_data);

				IF current_index = repeated_index THEN
					-- Replace this element with new value
					CALL _pb_wire_write_i64(value, temp_encoded);
					SET found_target = TRUE;
				ELSE
					-- Keep original element
					CALL _pb_wire_write_i64(temp_value, temp_encoded);
				END IF;

				SET new_packed_data = CONCAT(new_packed_data, temp_encoded);
				SET current_index = current_index + 1;
			END WHILE;

			-- If we found the target, update and return
			IF found_target THEN
				SET original_index = JSON_EXTRACT(element, '$.i');
				SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(new_packed_data));
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				RETURN JSON_SET(wire_json, element_path, new_element);
			END IF;
		WHEN 1 THEN
			-- Non-packed field (wire type 1 for I64)
			IF current_index = repeated_index THEN
				-- Found the target index - preserve original index and replace value
				SET original_index = JSON_EXTRACT(element, '$.i');
				SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 1, 'v', CAST(value AS JSON));
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				RETURN JSON_SET(wire_json, element_path, new_element);
			END IF;
			SET current_index = current_index + 1;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_set_repeated_i64_field_element: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET array_index = array_index + 1;
	END WHILE;

	-- Index not found
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
END $$

-- Private: Set repeated I32 field element at specific index
DROP FUNCTION IF EXISTS _pb_wire_json_set_repeated_i32_field_element $$
CREATE FUNCTION _pb_wire_json_set_repeated_i32_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT,
	value INT UNSIGNED
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE current_index INT DEFAULT 0;
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE original_index INT;
	DECLARE new_element JSON;
	DECLARE element_path TEXT;
	-- Variables for packed field handling
	DECLARE packed_data LONGBLOB;
	DECLARE new_packed_data LONGBLOB DEFAULT '';
	DECLARE temp_value INT UNSIGNED;
	DECLARE temp_encoded LONGBLOB;
	DECLARE found_target BOOLEAN DEFAULT FALSE;
	DECLARE message_text TEXT;

	-- Get the field array
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Find the target index across all array elements
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		CASE element_wire_type
		WHEN 2 THEN
			-- Packed field - decode, modify, and re-encode
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			SET new_packed_data = '';
			SET found_target = FALSE;

			-- Decode packed I32 values and rebuild with modification
			WHILE LENGTH(packed_data) > 0 DO
				CALL _pb_wire_read_i32_as_uint32(packed_data, temp_value, packed_data);

				IF current_index = repeated_index THEN
					-- Replace this element with new value
					CALL _pb_wire_write_i32(value, temp_encoded);
					SET found_target = TRUE;
				ELSE
					-- Keep original element
					CALL _pb_wire_write_i32(temp_value, temp_encoded);
				END IF;

				SET new_packed_data = CONCAT(new_packed_data, temp_encoded);
				SET current_index = current_index + 1;
			END WHILE;

			-- If we found the target, update and return
			IF found_target THEN
				SET original_index = JSON_EXTRACT(element, '$.i');
				SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(new_packed_data));
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				RETURN JSON_SET(wire_json, element_path, new_element);
			END IF;
		WHEN 5 THEN
			-- Non-packed field (wire type 5 for I32)
			IF current_index = repeated_index THEN
				-- Found the target index - preserve original index and replace value
				SET original_index = JSON_EXTRACT(element, '$.i');
				SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 5, 'v', CAST(value AS JSON));
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				RETURN JSON_SET(wire_json, element_path, new_element);
			END IF;
			SET current_index = current_index + 1;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_set_repeated_i32_field_element: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET array_index = array_index + 1;
	END WHILE;

	-- Index not found
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
END $$

-- Private: Set repeated LEN field element at specific index
DROP FUNCTION IF EXISTS _pb_wire_json_set_repeated_len_field_element $$
CREATE FUNCTION _pb_wire_json_set_repeated_len_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT,
	value LONGBLOB
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE current_index INT DEFAULT 0;
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE original_index INT;
	DECLARE new_element JSON;
	DECLARE element_path TEXT;
	DECLARE message_text TEXT;

	-- Get the field array
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Find the target index across all array elements
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		-- LEN fields are always wire type 2, so all elements should match
		CASE element_wire_type
		WHEN 2 THEN
			IF current_index = repeated_index THEN
				-- Found the target index - preserve original index and replace value
				SET original_index = JSON_EXTRACT(element, '$.i');
				SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(value));
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				RETURN JSON_SET(wire_json, element_path, new_element);
			END IF;
			SET current_index = current_index + 1;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_set_repeated_len_field_element: unexpected wire_type (', element_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET array_index = array_index + 1;
	END WHILE;

	-- Index not found
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
END $$

-- Private: Remove repeated VARINT field element at specific index
DROP FUNCTION IF EXISTS _pb_wire_json_remove_repeated_varint_field_element $$
CREATE FUNCTION _pb_wire_json_remove_repeated_varint_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE current_index INT DEFAULT 0;
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE element_path TEXT;
	-- Variables for packed field handling
	DECLARE packed_data LONGBLOB;
	DECLARE new_packed_data LONGBLOB DEFAULT '';
	DECLARE temp_value BIGINT UNSIGNED;
	DECLARE temp_encoded LONGBLOB;
	DECLARE found_target BOOLEAN DEFAULT FALSE;
	DECLARE original_index INT;
	DECLARE new_element JSON;

	-- Get the field array
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Check if field exists
	IF field_array IS NULL THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
	END IF;

	-- Find the target index across all array elements
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		IF element_wire_type = 2 THEN
			-- Packed field - decode, remove element, and re-encode
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			SET new_packed_data = '';
			SET found_target = FALSE;

			-- Decode packed varints and rebuild without target element
			WHILE LENGTH(packed_data) > 0 DO
				CALL _pb_wire_read_varint_as_uint64(packed_data, temp_value, packed_data);

				IF current_index = repeated_index THEN
					-- Skip this element (remove it)
					SET found_target = TRUE;
				ELSE
					-- Keep this element
					CALL _pb_wire_write_varint(temp_value, temp_encoded);
					SET new_packed_data = CONCAT(new_packed_data, temp_encoded);
				END IF;

				SET current_index = current_index + 1;
			END WHILE;

			-- If we found the target, update and return
			IF found_target THEN
				SET original_index = JSON_EXTRACT(element, '$.i');
				-- If no elements left, remove the entire field entry
				IF LENGTH(new_packed_data) = 0 THEN
					-- Check if this is the only element in the field array
					IF JSON_LENGTH(field_array) = 1 THEN
						-- Remove the entire field
						RETURN JSON_REMOVE(wire_json, field_path);
					ELSE
						-- Remove just this packed element
						SET element_path = CONCAT(field_path, '[', array_index, ']');
						RETURN JSON_REMOVE(wire_json, element_path);
					END IF;
				ELSE
					SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(new_packed_data));
					SET element_path = CONCAT(field_path, '[', array_index, ']');
					RETURN JSON_SET(wire_json, element_path, new_element);
				END IF;
			END IF;
		ELSE
			-- Non-packed field (wire type 0 for VARINT)
			IF current_index = repeated_index THEN
				-- Found the target index - remove this element
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				-- Check if this is the last element in the field
				IF JSON_LENGTH(field_array) = 1 THEN
					-- Remove the entire field
					RETURN JSON_REMOVE(wire_json, field_path);
				ELSE
					-- Remove just this element
					RETURN JSON_REMOVE(wire_json, element_path);
				END IF;
			END IF;
			SET current_index = current_index + 1;
		END IF;

		SET array_index = array_index + 1;
	END WHILE;

	-- Index not found
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
END $$

-- Private: Remove repeated I64 field element at specific index
DROP FUNCTION IF EXISTS _pb_wire_json_remove_repeated_i64_field_element $$
CREATE FUNCTION _pb_wire_json_remove_repeated_i64_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE current_index INT DEFAULT 0;
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE element_path TEXT;
	-- Variables for packed field handling
	DECLARE packed_data LONGBLOB;
	DECLARE new_packed_data LONGBLOB DEFAULT '';
	DECLARE temp_value BIGINT UNSIGNED;
	DECLARE temp_encoded LONGBLOB;
	DECLARE found_target BOOLEAN DEFAULT FALSE;
	DECLARE original_index INT;
	DECLARE new_element JSON;

	-- Get the field array
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Check if field exists
	IF field_array IS NULL THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
	END IF;

	-- Find the target index across all array elements
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		IF element_wire_type = 2 THEN
			-- Packed field - decode, remove element, and re-encode
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			SET new_packed_data = '';
			SET found_target = FALSE;

			-- Decode packed I64 values and rebuild without target element
			WHILE LENGTH(packed_data) > 0 DO
				CALL _pb_wire_read_i64_as_uint64(packed_data, temp_value, packed_data);

				IF current_index = repeated_index THEN
					-- Skip this element (remove it)
					SET found_target = TRUE;
				ELSE
					-- Keep this element
					CALL _pb_wire_write_i64(temp_value, temp_encoded);
					SET new_packed_data = CONCAT(new_packed_data, temp_encoded);
				END IF;

				SET current_index = current_index + 1;
			END WHILE;

			-- If we found the target, update and return
			IF found_target THEN
				SET original_index = JSON_EXTRACT(element, '$.i');
				-- If no elements left, remove the entire field entry
				IF LENGTH(new_packed_data) = 0 THEN
					-- Check if this is the only element in the field array
					IF JSON_LENGTH(field_array) = 1 THEN
						-- Remove the entire field
						RETURN JSON_REMOVE(wire_json, field_path);
					ELSE
						-- Remove just this packed element
						SET element_path = CONCAT(field_path, '[', array_index, ']');
						RETURN JSON_REMOVE(wire_json, element_path);
					END IF;
				ELSE
					SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(new_packed_data));
					SET element_path = CONCAT(field_path, '[', array_index, ']');
					RETURN JSON_SET(wire_json, element_path, new_element);
				END IF;
			END IF;
		ELSE
			-- Non-packed field (wire type 1 for I64)
			IF current_index = repeated_index THEN
				-- Found the target index - remove this element
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				-- Check if this is the last element in the field
				IF JSON_LENGTH(field_array) = 1 THEN
					-- Remove the entire field
					RETURN JSON_REMOVE(wire_json, field_path);
				ELSE
					-- Remove just this element
					RETURN JSON_REMOVE(wire_json, element_path);
				END IF;
			END IF;
			SET current_index = current_index + 1;
		END IF;

		SET array_index = array_index + 1;
	END WHILE;

	-- Index not found
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
END $$

-- Private: Remove repeated I32 field element at specific index
DROP FUNCTION IF EXISTS _pb_wire_json_remove_repeated_i32_field_element $$
CREATE FUNCTION _pb_wire_json_remove_repeated_i32_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE current_index INT DEFAULT 0;
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE element_path TEXT;
	-- Variables for packed field handling
	DECLARE packed_data LONGBLOB;
	DECLARE new_packed_data LONGBLOB DEFAULT '';
	DECLARE temp_value INT UNSIGNED;
	DECLARE temp_encoded LONGBLOB;
	DECLARE found_target BOOLEAN DEFAULT FALSE;
	DECLARE original_index INT;
	DECLARE new_element JSON;

	-- Get the field array
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Check if field exists
	IF field_array IS NULL THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
	END IF;

	-- Find the target index across all array elements
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		IF element_wire_type = 2 THEN
			-- Packed field - decode, remove element, and re-encode
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			SET new_packed_data = '';
			SET found_target = FALSE;

			-- Decode packed I32 values and rebuild without target element
			WHILE LENGTH(packed_data) > 0 DO
				CALL _pb_wire_read_i32_as_uint32(packed_data, temp_value, packed_data);

				IF current_index = repeated_index THEN
					-- Skip this element (remove it)
					SET found_target = TRUE;
				ELSE
					-- Keep this element
					CALL _pb_wire_write_i32(temp_value, temp_encoded);
					SET new_packed_data = CONCAT(new_packed_data, temp_encoded);
				END IF;

				SET current_index = current_index + 1;
			END WHILE;

			-- If we found the target, update and return
			IF found_target THEN
				SET original_index = JSON_EXTRACT(element, '$.i');
				-- If no elements left, remove the entire field entry
				IF LENGTH(new_packed_data) = 0 THEN
					-- Check if this is the only element in the field array
					IF JSON_LENGTH(field_array) = 1 THEN
						-- Remove the entire field
						RETURN JSON_REMOVE(wire_json, field_path);
					ELSE
						-- Remove just this packed element
						SET element_path = CONCAT(field_path, '[', array_index, ']');
						RETURN JSON_REMOVE(wire_json, element_path);
					END IF;
				ELSE
					SET new_element = JSON_OBJECT('i', original_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(new_packed_data));
					SET element_path = CONCAT(field_path, '[', array_index, ']');
					RETURN JSON_SET(wire_json, element_path, new_element);
				END IF;
			END IF;
		ELSE
			-- Non-packed field (wire type 5 for I32)
			IF current_index = repeated_index THEN
				-- Found the target index - remove this element
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				-- Check if this is the last element in the field
				IF JSON_LENGTH(field_array) = 1 THEN
					-- Remove the entire field
					RETURN JSON_REMOVE(wire_json, field_path);
				ELSE
					-- Remove just this element
					RETURN JSON_REMOVE(wire_json, element_path);
				END IF;
			END IF;
			SET current_index = current_index + 1;
		END IF;

		SET array_index = array_index + 1;
	END WHILE;

	-- Index not found
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
END $$

-- Private: Remove repeated LEN field element at specific index
DROP FUNCTION IF EXISTS _pb_wire_json_remove_repeated_len_field_element $$
CREATE FUNCTION _pb_wire_json_remove_repeated_len_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE current_index INT DEFAULT 0;
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE element_path TEXT;

	-- Get the field array
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Check if field exists
	IF field_array IS NULL THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
	END IF;

	-- Find the target index across all array elements
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		-- LEN fields are always wire type 2 and don't support packed encoding
		IF element_wire_type = 2 THEN
			IF current_index = repeated_index THEN
				-- Found the target index - remove this element
				SET element_path = CONCAT(field_path, '[', array_index, ']');
				-- Check if this is the last element in the field
				IF JSON_LENGTH(field_array) = 1 THEN
					-- Remove the entire field
					RETURN JSON_REMOVE(wire_json, field_path);
				ELSE
					-- Remove just this element
					RETURN JSON_REMOVE(wire_json, element_path);
				END IF;
			END IF;
			SET current_index = current_index + 1;
		END IF;

		SET array_index = array_index + 1;
	END WHILE;

	-- Index not found
	SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
END $$

-- Private: Add to packed varint field
DROP FUNCTION IF EXISTS _pb_wire_json_add_packed_varint_field $$
CREATE FUNCTION _pb_wire_json_add_packed_varint_field(wire_json JSON, field_number INT, value BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE packed_data LONGBLOB;
	DECLARE new_varint LONGBLOB;
	DECLARE last_element JSON;
	DECLARE next_index INT;
	DECLARE new_element JSON;
	DECLARE last_index INT;

	-- Get the field array (null if doesn't exist)
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Encode the new varint value
	CALL _pb_wire_write_varint(value, new_varint);

	-- Check if field exists and last element is LEN (wire type 2)
	IF field_array IS NOT NULL THEN
		SET last_index = JSON_LENGTH(field_array) - 1;
		SET last_element = JSON_EXTRACT(field_array, CONCAT('$[', last_index, ']'));
		IF JSON_EXTRACT(last_element, '$.t') = 2 THEN
			-- Append to existing packed data
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(last_element, '$.v')));
			SET packed_data = CONCAT(packed_data, new_varint);
			RETURN JSON_SET(wire_json, CONCAT(field_path, '[', last_index, '].v'), _pb_util_to_base64(packed_data));
		END IF;
	END IF;

	-- Create new packed element
	SET next_index = _pb_wire_json_get_next_index(wire_json);
	SET new_element = JSON_OBJECT('i', next_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(new_varint));

	IF field_array IS NULL THEN
		RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
	ELSE
		RETURN JSON_ARRAY_APPEND(wire_json, field_path, new_element);
	END IF;
END $$

-- Private: Add to packed I64 field
DROP FUNCTION IF EXISTS _pb_wire_json_add_packed_i64_field $$
CREATE FUNCTION _pb_wire_json_add_packed_i64_field(wire_json JSON, field_number INT, value BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE packed_data LONGBLOB;
	DECLARE new_i64 LONGBLOB;
	DECLARE last_element JSON;
	DECLARE next_index INT;
	DECLARE new_element JSON;
	DECLARE last_index INT;

	-- Get the field array (null if doesn't exist)
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Encode the new I64 value (8 bytes little-endian)
	CALL _pb_wire_write_i64(value, new_i64);

	-- Check if field exists and last element is LEN (wire type 2)
	IF field_array IS NOT NULL THEN
		SET last_index = JSON_LENGTH(field_array) - 1;
		SET last_element = JSON_EXTRACT(field_array, CONCAT('$[', last_index, ']'));
		IF JSON_EXTRACT(last_element, '$.t') = 2 THEN
			-- Append to existing packed data
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(last_element, '$.v')));
			SET packed_data = CONCAT(packed_data, new_i64);
			RETURN JSON_SET(wire_json, CONCAT(field_path, '[', last_index, '].v'), _pb_util_to_base64(packed_data));
		END IF;
	END IF;

	-- Create new packed element
	SET next_index = _pb_wire_json_get_next_index(wire_json);
	SET new_element = JSON_OBJECT('i', next_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(new_i64));

	IF field_array IS NULL THEN
		RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
	ELSE
		RETURN JSON_ARRAY_APPEND(wire_json, field_path, new_element);
	END IF;
END $$

-- Private: Add to packed I32 field
DROP FUNCTION IF EXISTS _pb_wire_json_add_packed_i32_field $$
CREATE FUNCTION _pb_wire_json_add_packed_i32_field(wire_json JSON, field_number INT, value INT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE packed_data LONGBLOB;
	DECLARE new_i32 LONGBLOB;
	DECLARE last_element JSON;
	DECLARE next_index INT;
	DECLARE new_element JSON;
	DECLARE last_index INT;

	-- Get the field array (null if doesn't exist)
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Encode the new I32 value (4 bytes little-endian)
	CALL _pb_wire_write_i32(value, new_i32);

	-- Check if field exists and last element is LEN (wire type 2)
	IF field_array IS NOT NULL THEN
		SET last_index = JSON_LENGTH(field_array) - 1;
		SET last_element = JSON_EXTRACT(field_array, CONCAT('$[', last_index, ']'));
		IF JSON_EXTRACT(last_element, '$.t') = 2 THEN
			-- Append to existing packed data
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(last_element, '$.v')));
			SET packed_data = CONCAT(packed_data, new_i32);
			RETURN JSON_SET(wire_json, CONCAT(field_path, '[', last_index, '].v'), _pb_util_to_base64(packed_data));
		END IF;
	END IF;

	-- Create new packed element
	SET next_index = _pb_wire_json_get_next_index(wire_json);
	SET new_element = JSON_OBJECT('i', next_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(new_i32));

	IF field_array IS NULL THEN
		RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
	ELSE
		RETURN JSON_ARRAY_APPEND(wire_json, field_path, new_element);
	END IF;
END $$

-- Private: Set VARINT field (int32, int64, uint32, uint64, sint32, sint64, enum, bool)
DROP FUNCTION IF EXISTS _pb_wire_json_set_varint_field $$
CREATE FUNCTION _pb_wire_json_set_varint_field(wire_json JSON, field_number INT, value BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_field(wire_json, field_number, 0, CAST(value AS JSON));
END $$

-- Private: Set I64 field (fixed64, sfixed64, double)
DROP FUNCTION IF EXISTS _pb_wire_json_set_i64_field $$
CREATE FUNCTION _pb_wire_json_set_i64_field(wire_json JSON, field_number INT, value BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_field(wire_json, field_number, 1, CAST(value AS JSON));
END $$

-- Private: Set I32 field (fixed32, sfixed32, float)
DROP FUNCTION IF EXISTS _pb_wire_json_set_i32_field $$
CREATE FUNCTION _pb_wire_json_set_i32_field(wire_json JSON, field_number INT, value INT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_field(wire_json, field_number, 5, CAST(value AS JSON));
END $$

-- Private: Set length-delimited field (string, bytes, message)
DROP FUNCTION IF EXISTS _pb_wire_json_set_len_field $$
CREATE FUNCTION _pb_wire_json_set_len_field(wire_json JSON, field_number INT, value LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_field(wire_json, field_number, 2, JSON_QUOTE(_pb_util_to_base64(value)));
END $$

-- Private: Add to repeated VARINT field
DROP FUNCTION IF EXISTS _pb_wire_json_add_repeated_varint_field_element $$
CREATE FUNCTION _pb_wire_json_add_repeated_varint_field_element(wire_json JSON, field_number INT, value BIGINT UNSIGNED, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	IF use_packed THEN
		RETURN _pb_wire_json_add_packed_varint_field(wire_json, field_number, value);
	ELSE
		RETURN _pb_wire_json_add_repeated_field_element(wire_json, field_number, 0, CAST(value AS JSON));
	END IF;
END $$

-- Private: Add to repeated I64 field
DROP FUNCTION IF EXISTS _pb_wire_json_add_repeated_i64_field_element $$
CREATE FUNCTION _pb_wire_json_add_repeated_i64_field_element(wire_json JSON, field_number INT, value BIGINT UNSIGNED, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	IF use_packed THEN
		RETURN _pb_wire_json_add_packed_i64_field(wire_json, field_number, value);
	ELSE
		RETURN _pb_wire_json_add_repeated_field_element(wire_json, field_number, 1, CAST(value AS JSON));
	END IF;
END $$

-- Private: Add to repeated I32 field
DROP FUNCTION IF EXISTS _pb_wire_json_add_repeated_i32_field_element $$
CREATE FUNCTION _pb_wire_json_add_repeated_i32_field_element(wire_json JSON, field_number INT, value INT UNSIGNED, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	IF use_packed THEN
		RETURN _pb_wire_json_add_packed_i32_field(wire_json, field_number, value);
	ELSE
		RETURN _pb_wire_json_add_repeated_field_element(wire_json, field_number, 5, CAST(value AS JSON));
	END IF;
END $$

-- Private: Add to repeated length-delimited field
DROP FUNCTION IF EXISTS _pb_wire_json_add_repeated_len_field_element $$
CREATE FUNCTION _pb_wire_json_add_repeated_len_field_element(wire_json JSON, field_number INT, value LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_field_element(wire_json, field_number, 2, JSON_QUOTE(_pb_util_to_base64(value)));
END $$

-- Private: Insert into repeated VARINT field
DROP FUNCTION IF EXISTS _pb_wire_json_insert_repeated_varint_field_element $$
CREATE FUNCTION _pb_wire_json_insert_repeated_varint_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT,
	value BIGINT UNSIGNED,
	use_packed BOOLEAN
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE new_field_array JSON DEFAULT JSON_ARRAY();
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE next_wire_index INT;
	DECLARE new_element JSON;
	DECLARE logical_position INT DEFAULT 0;
	DECLARE inserted BOOLEAN DEFAULT FALSE;
	-- Variables for processing elements
	DECLARE packed_data LONGBLOB;
	DECLARE temp_value BIGINT UNSIGNED;
	DECLARE temp_encoded LONGBLOB;
	DECLARE result_packed_data LONGBLOB DEFAULT '';

	-- Get the field array (null if doesn't exist)
	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- Calculate wire index once
	SET next_wire_index = _pb_wire_json_get_next_index(wire_json);

	-- If field doesn't exist, create new field with single element
	IF field_array IS NULL THEN
		IF repeated_index = 0 THEN
			IF use_packed THEN
				CALL _pb_wire_write_varint(value, temp_encoded);
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(temp_encoded));
				SET next_wire_index = next_wire_index + 1;
			ELSE
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 0, 'v', CAST(value AS JSON));
				SET next_wire_index = next_wire_index + 1;
			END IF;
			RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
		ELSE
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
		END IF;
	END IF;

	-- Process existing elements while building new array directly
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		IF element_wire_type = 2 THEN
			-- Packed field - process each packed value
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			WHILE LENGTH(packed_data) > 0 DO
				CALL _pb_wire_read_varint_as_uint64(packed_data, temp_value, packed_data);

				-- Check if we need to insert here
				IF NOT inserted AND logical_position = repeated_index THEN
					IF use_packed THEN
						CALL _pb_wire_write_varint(value, temp_encoded);
						SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
					ELSE
						SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 0, 'v', CAST(value AS JSON));
						SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
						SET next_wire_index = next_wire_index + 1;
					END IF;
					SET inserted = TRUE;
				END IF;

				-- Add current value
				IF use_packed THEN
					CALL _pb_wire_write_varint(temp_value, temp_encoded);
					SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
				ELSE
					SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 0, 'v', CAST(temp_value AS JSON));
					SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
					SET next_wire_index = next_wire_index + 1;
				END IF;

				SET logical_position = logical_position + 1;
			END WHILE;
		ELSE
			-- Unpacked field - process single value
			SET temp_value = CAST(JSON_EXTRACT(element, '$.v') AS UNSIGNED);

			-- Check if we need to insert here
			IF NOT inserted AND logical_position = repeated_index THEN
				IF use_packed THEN
					CALL _pb_wire_write_varint(value, temp_encoded);
					SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
				ELSE
					SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 0, 'v', CAST(value AS JSON));
					SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
					SET next_wire_index = next_wire_index + 1;
				END IF;
				SET inserted = TRUE;
			END IF;

			-- Add current value
			IF use_packed THEN
				CALL _pb_wire_write_varint(temp_value, temp_encoded);
				SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
			ELSE
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 0, 'v', CAST(temp_value AS JSON));
				SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
				SET next_wire_index = next_wire_index + 1;
			END IF;

			SET logical_position = logical_position + 1;
		END IF;

		SET array_index = array_index + 1;
	END WHILE;

	-- Check if we need to append at the end
	IF NOT inserted THEN
		IF repeated_index = logical_position THEN
			IF use_packed THEN
				CALL _pb_wire_write_varint(value, temp_encoded);
				SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
			ELSE
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 0, 'v', CAST(value AS JSON));
				SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
			END IF;
		ELSE
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
		END IF;
	END IF;

	-- Return result based on format
	IF use_packed THEN
		SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(result_packed_data));
		RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
	ELSE
		RETURN JSON_SET(wire_json, field_path, new_field_array);
	END IF;
END $$

-- Private: Insert into repeated I64 field
DROP FUNCTION IF EXISTS _pb_wire_json_insert_repeated_i64_field_element $$
CREATE FUNCTION _pb_wire_json_insert_repeated_i64_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT,
	value BIGINT UNSIGNED,
	use_packed BOOLEAN
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE new_field_array JSON DEFAULT JSON_ARRAY();
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE next_wire_index INT;
	DECLARE new_element JSON;
	DECLARE logical_position INT DEFAULT 0;
	DECLARE inserted BOOLEAN DEFAULT FALSE;
	-- Variables for processing elements
	DECLARE packed_data LONGBLOB;
	DECLARE temp_value BIGINT UNSIGNED;
	DECLARE temp_encoded LONGBLOB;
	DECLARE result_packed_data LONGBLOB DEFAULT '';

	-- Get the field array (null if doesn't exist)
	-- Calculate wire index once
	SET next_wire_index = _pb_wire_json_get_next_index(wire_json);

	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- If field doesn't exist, create new field with single element
	IF field_array IS NULL THEN
		IF repeated_index = 0 THEN
			IF use_packed THEN
				CALL _pb_wire_write_i64(value, temp_encoded);
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(temp_encoded));
				SET next_wire_index = next_wire_index + 1;
			ELSE
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 1, 'v', CAST(value AS JSON));
				SET next_wire_index = next_wire_index + 1;
			END IF;
			RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
		ELSE
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
		END IF;
	END IF;

	-- Process existing elements while building new array directly
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		IF element_wire_type = 2 THEN
			-- Packed field - process each packed value
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			WHILE LENGTH(packed_data) > 0 DO
				CALL _pb_wire_read_i64_as_uint64(packed_data, temp_value, packed_data);

				-- Check if we need to insert here
				IF NOT inserted AND logical_position = repeated_index THEN
					IF use_packed THEN
						CALL _pb_wire_write_i64(value, temp_encoded);
						SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
					ELSE
						SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 1, 'v', CAST(value AS JSON));
						SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
						SET next_wire_index = next_wire_index + 1;
					END IF;
					SET inserted = TRUE;
				END IF;

				-- Add current value
				IF use_packed THEN
					CALL _pb_wire_write_i64(temp_value, temp_encoded);
					SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
				ELSE
					SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 1, 'v', CAST(temp_value AS JSON));
					SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
					SET next_wire_index = next_wire_index + 1;
				END IF;

				SET logical_position = logical_position + 1;
			END WHILE;
		ELSE
			-- Unpacked field - process single value
			SET temp_value = CAST(JSON_EXTRACT(element, '$.v') AS UNSIGNED);

			-- Check if we need to insert here
			IF NOT inserted AND logical_position = repeated_index THEN
				IF use_packed THEN
					CALL _pb_wire_write_i64(value, temp_encoded);
					SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
				ELSE
					SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 1, 'v', CAST(value AS JSON));
					SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
					SET next_wire_index = next_wire_index + 1;
				END IF;
				SET inserted = TRUE;
			END IF;

			-- Add current value
			IF use_packed THEN
				CALL _pb_wire_write_i64(temp_value, temp_encoded);
				SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
			ELSE
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 1, 'v', CAST(temp_value AS JSON));
				SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
				SET next_wire_index = next_wire_index + 1;
			END IF;

			SET logical_position = logical_position + 1;
		END IF;

		SET array_index = array_index + 1;
	END WHILE;

	-- Check if we need to append at the end
	IF NOT inserted THEN
		IF repeated_index = logical_position THEN
			IF use_packed THEN
				CALL _pb_wire_write_i64(value, temp_encoded);
				SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
			ELSE
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 1, 'v', CAST(value AS JSON));
				SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
			END IF;
		ELSE
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
		END IF;
	END IF;

	-- Return result based on format
	IF use_packed THEN
		SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(result_packed_data));
		RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
	ELSE
		RETURN JSON_SET(wire_json, field_path, new_field_array);
	END IF;
END $$

-- Private: Insert into repeated I32 field
DROP FUNCTION IF EXISTS _pb_wire_json_insert_repeated_i32_field_element $$
CREATE FUNCTION _pb_wire_json_insert_repeated_i32_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT,
	value INT UNSIGNED,
	use_packed BOOLEAN
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE new_field_array JSON DEFAULT JSON_ARRAY();
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE element_wire_type INT;
	DECLARE next_wire_index INT;
	DECLARE new_element JSON;
	DECLARE logical_position INT DEFAULT 0;
	DECLARE inserted BOOLEAN DEFAULT FALSE;
	-- Variables for processing elements
	DECLARE packed_data LONGBLOB;
	DECLARE temp_value INT UNSIGNED;
	DECLARE temp_encoded LONGBLOB;
	DECLARE result_packed_data LONGBLOB DEFAULT '';

	-- Get the field array (null if doesn't exist)
	-- Calculate wire index once
	SET next_wire_index = _pb_wire_json_get_next_index(wire_json);

	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- If field doesn't exist, create new field with single element
	IF field_array IS NULL THEN
		IF repeated_index = 0 THEN
			IF use_packed THEN
				CALL _pb_wire_write_i32(value, temp_encoded);
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(temp_encoded));
				SET next_wire_index = next_wire_index + 1;
			ELSE
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 5, 'v', CAST(value AS JSON));
				SET next_wire_index = next_wire_index + 1;
			END IF;
			RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
		ELSE
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
		END IF;
	END IF;

	-- Process existing elements while building new array directly
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET element_wire_type = JSON_EXTRACT(element, '$.t');

		IF element_wire_type = 2 THEN
			-- Packed field - process each packed value
			SET packed_data = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			WHILE LENGTH(packed_data) > 0 DO
				CALL _pb_wire_read_i32_as_uint32(packed_data, temp_value, packed_data);

				-- Check if we need to insert here
				IF NOT inserted AND logical_position = repeated_index THEN
					IF use_packed THEN
						CALL _pb_wire_write_i32(value, temp_encoded);
						SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
					ELSE
						SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 5, 'v', CAST(value AS JSON));
						SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
						SET next_wire_index = next_wire_index + 1;
					END IF;
					SET inserted = TRUE;
				END IF;

				-- Add current value
				IF use_packed THEN
					CALL _pb_wire_write_i32(temp_value, temp_encoded);
					SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
				ELSE
					SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 5, 'v', CAST(temp_value AS JSON));
					SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
					SET next_wire_index = next_wire_index + 1;
				END IF;

				SET logical_position = logical_position + 1;
			END WHILE;
		ELSE
			-- Unpacked field - process single value
			SET temp_value = CAST(JSON_EXTRACT(element, '$.v') AS UNSIGNED);

			-- Check if we need to insert here
			IF NOT inserted AND logical_position = repeated_index THEN
				IF use_packed THEN
					CALL _pb_wire_write_i32(value, temp_encoded);
					SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
				ELSE
					SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 5, 'v', CAST(value AS JSON));
					SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
					SET next_wire_index = next_wire_index + 1;
				END IF;
				SET inserted = TRUE;
			END IF;

			-- Add current value
			IF use_packed THEN
				CALL _pb_wire_write_i32(temp_value, temp_encoded);
				SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
			ELSE
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 5, 'v', CAST(temp_value AS JSON));
				SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
				SET next_wire_index = next_wire_index + 1;
			END IF;

			SET logical_position = logical_position + 1;
		END IF;

		SET array_index = array_index + 1;
	END WHILE;

	-- Check if we need to append at the end
	IF NOT inserted THEN
		IF repeated_index = logical_position THEN
			IF use_packed THEN
				CALL _pb_wire_write_i32(value, temp_encoded);
				SET result_packed_data = CONCAT(result_packed_data, temp_encoded);
			ELSE
				SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 5, 'v', CAST(value AS JSON));
				SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
			END IF;
		ELSE
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
		END IF;
	END IF;

	-- Return result based on format
	IF use_packed THEN
		SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 2, 'v', _pb_util_to_base64(result_packed_data));
		RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
	ELSE
		RETURN JSON_SET(wire_json, field_path, new_field_array);
	END IF;
END $$

-- Private: Insert into repeated length-delimited field
DROP FUNCTION IF EXISTS _pb_wire_json_insert_repeated_len_field_element $$
CREATE FUNCTION _pb_wire_json_insert_repeated_len_field_element(
	wire_json JSON,
	field_number INT,
	repeated_index INT,
	value LONGBLOB
) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_path TEXT DEFAULT CONCAT('$."', field_number, '"');
	DECLARE field_array JSON;
	DECLARE new_field_array JSON DEFAULT JSON_ARRAY();
	DECLARE logical_values JSON DEFAULT JSON_ARRAY();
	DECLARE array_index INT DEFAULT 0;
	DECLARE element JSON;
	DECLARE next_wire_index INT;
	DECLARE new_element JSON;
	DECLARE i INT DEFAULT 0;
	DECLARE temp_value JSON;
	DECLARE encoded_value JSON DEFAULT JSON_QUOTE(_pb_util_to_base64(value));

	-- Get the field array (null if doesn't exist)
	-- Calculate wire index once
	SET next_wire_index = _pb_wire_json_get_next_index(wire_json);

	SET field_array = JSON_EXTRACT(wire_json, field_path);

	-- If field doesn't exist, create new field with single element
	IF field_array IS NULL THEN
		IF repeated_index = 0 THEN
			SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 2, 'v', encoded_value);
			SET next_wire_index = next_wire_index + 1;
			RETURN JSON_SET(wire_json, field_path, JSON_ARRAY(new_element));
		ELSE
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
		END IF;
	END IF;

	-- Extract all logical values from existing wire elements (length-delimited fields are never packed)
	WHILE array_index < JSON_LENGTH(field_array) DO
		SET element = JSON_EXTRACT(field_array, CONCAT('$[', array_index, ']'));
		SET temp_value = JSON_EXTRACT(element, '$.v');
		SET logical_values = JSON_ARRAY_APPEND(logical_values, '$', temp_value);
		SET array_index = array_index + 1;
	END WHILE;

	-- Check bounds
	IF repeated_index < 0 OR repeated_index > JSON_LENGTH(logical_values) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Index out of bounds';
	END IF;

	-- Insert new value at the specified position
	SET logical_values = JSON_ARRAY_INSERT(logical_values, CONCAT('$[', repeated_index, ']'), encoded_value);

	-- Rebuild as unpacked wire elements
	SET i = 0;
	WHILE i < JSON_LENGTH(logical_values) DO
		SET temp_value = JSON_EXTRACT(logical_values, CONCAT('$[', i, ']'));
		SET new_element = JSON_OBJECT('i', next_wire_index, 'n', field_number, 't', 2, 'v', temp_value);
		SET next_wire_index = next_wire_index + 1;
		SET new_field_array = JSON_ARRAY_APPEND(new_field_array, '$', new_element);
		SET i = i + 1;
	END WHILE;

	-- Replace the field array with the new one
	RETURN JSON_SET(wire_json, field_path, new_field_array);
END $$
