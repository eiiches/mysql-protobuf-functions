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

DROP FUNCTION IF EXISTS pb_message_new $$
CREATE FUNCTION pb_message_new() RETURNS LONGBLOB DETERMINISTIC
BEGIN
	-- Return an empty protobuf message
	RETURN _binary X'';
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
