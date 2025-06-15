DELIMITER $$

DROP FUNCTION IF EXISTS _pb_util_bin_as_int32 $$
CREATE FUNCTION _pb_util_bin_as_int32(b BLOB) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(b) > 4 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_util_bin_as_int32: value must not be longer than 4 bytes.';
	END IF;

	IF LPAD(b, 4, _binary X'00') & _binary X'80000000' = _binary X'00000000' THEN
		RETURN CONV(HEX(b), 16, 10);
	ELSE
		RETURN -(CONV(HEX(~b), 16, 10) + 1);
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_bin_as_uint32 $$
CREATE FUNCTION _pb_util_bin_as_uint32(b BLOB) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	RETURN CONV(HEX(b), 16, 10);
END $$

DROP FUNCTION IF EXISTS _pb_util_bin_as_int64 $$
CREATE FUNCTION _pb_util_bin_as_int64(b BLOB) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(b) > 8 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_util_bin_as_int64: value must not be longer than 4 bytes.';
	END IF;

	IF LPAD(b, 8, _binary X'00') & _binary X'8000000000000000' = _binary X'0000000000000000' THEN
		RETURN CONV(HEX(b), 16, 10);
	ELSE
		RETURN -(CONV(HEX(~b), 16, 10) + 1);
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_bin_as_uint64 $$
CREATE FUNCTION _pb_util_bin_as_uint64(b BLOB) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	RETURN CONV(HEX(b), 16, 10);
END $$

DROP PROCEDURE IF EXISTS pb_wire_read_i32 $$
CREATE PROCEDURE pb_wire_read_i32(IN buf BLOB, OUT value BINARY(4), OUT tail BLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(buf) < 4 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'pb_wire_read_i32: Unexpected end of BLOB.';
	END IF;

	SET value = LEFT(buf, 4);
	SET tail = SUBSTRING(buf, 5);
END $$

DROP FUNCTION IF EXISTS pb_wire_read_i32 $$
CREATE FUNCTION pb_wire_read_i32(buf BLOB) RETURNS BINARY(4) DETERMINISTIC
BEGIN
	DECLARE tail BLOB;
	DECLARE value BINARY(4);
	CALL pb_wire_read_i32(buf, value, tail);
	RETURN value;
END $$

DROP PROCEDURE IF EXISTS pb_wire_read_i64 $$
CREATE PROCEDURE pb_wire_read_i64(IN buf BLOB, OUT value BINARY(8), OUT tail BLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(buf) < 8 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'pb_wire_read_i64: Unexpected end of BLOB.';
	END IF;

	SET value = LEFT(buf, 8);
	SET tail = SUBSTRING(buf, 9);
END $$

DROP FUNCTION IF EXISTS pb_wire_read_i64 $$
CREATE FUNCTION pb_wire_read_i64(buf BLOB) RETURNS BINARY(8) DETERMINISTIC
BEGIN
	DECLARE tail BLOB;
	DECLARE value BINARY(8);
	CALL pb_wire_read_i64(buf, value, tail);
	RETURN value;
END $$

DROP PROCEDURE IF EXISTS pb_wire_read_len_type $$
CREATE PROCEDURE pb_wire_read_len_type(IN buf BLOB, OUT value BLOB, OUT tail BLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE len BIGINT;

	SET tail = buf;
	CALL pb_wire_read_varint(tail, len, tail);

	IF LENGTH(tail) < len THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'pb_wire_read_len_type: Unexpected end of BLOB.';
	END IF;

	SET value = LEFT(tail, len);
	SET tail = SUBSTRING(tail, len + 1);
END $$

DROP PROCEDURE IF EXISTS pb_wire_read_varint $$
CREATE PROCEDURE pb_wire_read_varint(IN buf BLOB, OUT value BIGINT, OUT tail BLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE head INT;
	DECLARE byte_index INT;

	SET value = 0;
	SET tail = buf;
	SET byte_index = 0;

	l1: LOOP
		IF LENGTH(tail) = 0 THEN
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'pb_wire_read_varint: Unexpected end of BLOB.';
		END IF;

		SET head = _pb_util_bin_as_int32(LEFT(tail, 1));
		SET tail = SUBSTRING(tail, 2);

		SET value = value + ((head & 0x7f) << (7 * byte_index));

		IF (head & 0x80) = 0 THEN
			LEAVE l1;
		END IF;

		SET byte_index = byte_index + 1;
		IF byte_index > 10 THEN
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'pb_wire_read_varint: Varint cannot exceed 10 bytes.';
		END IF;
	END LOOP;
END $$

DROP FUNCTION IF EXISTS pb_wire_read_varint $$
CREATE FUNCTION pb_wire_read_varint(buf BLOB) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE tail BLOB;
	DECLARE value BIGINT;
	CALL pb_wire_read_varint(buf, value, tail);
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
CREATE PROCEDURE _pb_message_get_len_type_field(IN buf BLOB, IN field_number INT, IN repeated_index INT, OUT value BLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail BLOB;
	DECLARE int_value BIGINT;
	DECLARE bytes_value BLOB;
	DECLARE occurrence INT;
	DECLARE message_text TEXT;

	SET value = 0; -- proto3 default value for an integer field
	SET tail = buf;
	SET occurrence = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL pb_wire_read_varint(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 2 /* VARINT */ THEN
			SET message_text = CONCAT('_pb_message_get_len_type_field: string or bytes value cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL pb_wire_read_varint(tail, int_value, tail);
		WHEN 1 THEN -- I64
			CALL pb_wire_read_i64(tail, bytes_value, tail);
		WHEN 2 THEN -- LEN
			CALL pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				IF repeated_index IS NULL OR repeated_index = occurrence THEN
					SET value = bytes_value;
					SET occurrence = occurrence + 1;
				END IF;
			END IF;
		WHEN 5 THEN -- I32
			CALL pb_wire_read_i32(tail, bytes_value, tail);
		ELSE
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_len_type_field: unsupported wire_type';
		END CASE;
	END WHILE;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_int32_or_int64_field $$
CREATE PROCEDURE _pb_message_get_int32_or_int64_field(IN buf BLOB, IN field_number INT, IN repeated_index INT, OUT value BIGINT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail BLOB;
	DECLARE int_value BIGINT;
	DECLARE bytes_value BLOB;
	DECLARE occurrence INT;
	DECLARE message_text TEXT;

	SET value = 0; -- proto3 default value for an integer field
	SET tail = buf;
	SET occurrence = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL pb_wire_read_varint(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 0 /* VARINT */ AND (repeated_index IS NULL OR _pb_wire_get_wire_type_from_tag(tag) <> 2 /* LEN */) THEN
			SET message_text = CONCAT('_pb_message_get_int32_or_int64_field: int32 or int64 value cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL pb_wire_read_varint(tail, int_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				IF repeated_index IS NULL OR repeated_index = occurrence THEN
					SET value = int_value;
				END IF;
				SET occurrence = occurrence + 1;
			END IF;
		WHEN 1 THEN -- I64
			CALL pb_wire_read_i64(tail, bytes_value, tail);
		WHEN 2 THEN -- LEN
			CALL pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number AND repeated_index IS NOT NULL THEN
				WHILE LENGTH(bytes_value) <> 0 DO
					CALL pb_wire_read_varint(bytes_value, int_value, bytes_value);
					IF repeated_index = occurrence THEN
						SET value = int_value;
					END IF;
					SET occurrence = occurrence + 1;
				END WHILE;
			END IF;
		WHEN 5 THEN -- I32
			CALL pb_wire_read_i32(tail, bytes_value, tail);
		ELSE
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_int32_or_int64_field: unsupported wire_type';
		END CASE;
	END WHILE;

	IF repeated_index IS NOT NULL AND occurrence <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_int32_or_int64_field: repeated index out of range';
	END IF;
END $$

DROP FUNCTION IF EXISTS pb_message_get_int32_field $$
CREATE FUNCTION pb_message_get_int32_field(buf BLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT;
	CALL _pb_message_get_int32_or_int64_field(buf, field_number, repeated_index, value);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_uint32_field $$
CREATE FUNCTION pb_message_get_uint32_field(buf BLOB, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_int64_field $$
CREATE FUNCTION pb_message_get_int64_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT;
	CALL _pb_message_get_int32_or_int64_field(buf, field_number, repeated_index, value);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_uint64_field $$
CREATE FUNCTION pb_message_get_uint64_field(buf BLOB, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_string_field $$
CREATE FUNCTION pb_message_get_string_field(buf BLOB, field_number INT, repeated_index INT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	CALL _pb_message_get_len_type_field(buf, field_number, repeated_index, value);
	RETURN CONVERT(value USING utf8mb4);
END $$

DROP FUNCTION IF EXISTS pb_message_get_bytes_field $$
CREATE FUNCTION pb_message_get_bytes_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BLOB DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	CALL _pb_message_get_len_type_field(buf, field_number, repeated_index, value);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_message_field $$
CREATE FUNCTION pb_message_get_message_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BLOB DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	CALL _pb_message_get_len_type_field(buf, field_number, repeated_index, value);
	RETURN value;
END $$
