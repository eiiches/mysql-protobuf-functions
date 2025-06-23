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
		CALL _pb_wire_skip_varint(tail, tail);
	WHEN 1 THEN -- I64
		CALL _pb_wire_skip_i64(tail, tail);
	WHEN 2 THEN -- LEN
		CALL _pb_wire_skip_len_type(tail, tail);
	WHEN 5 THEN -- I32
		CALL _pb_wire_skip_i32(tail, tail);
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

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_int64;
CREATE FUNCTION _pb_util_reinterpret_uint64_as_int64(value BIGINT UNSIGNED) RETURNS BIGINT DETERMINISTIC
BEGIN
	IF value <= 0x7fffffffffffffff THEN
		RETURN CAST(value AS SIGNED);
	ELSE
		RETURN value - 18446744073709551616; -- 2^64
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint32_as_int32;
CREATE FUNCTION _pb_util_reinterpret_uint32_as_int32(value INT UNSIGNED) RETURNS INT DETERMINISTIC
BEGIN
	IF value <= 0x7fffffff THEN
		RETURN CAST(value AS SIGNED);
	ELSE
		RETURN CAST(value AS SIGNED) - 4294967296; -- 2^32
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_int32;
CREATE FUNCTION _pb_util_reinterpret_uint64_as_int32(value BIGINT UNSIGNED) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF value > 0xffffffff THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_util_reinterpret_uint64_as_int32: value is larger than 0xffffffff';
	END IF;

	RETURN IF(
		value <= 0x7fffffff,
		value,
		CAST(value AS SIGNED) - 4294967296 -- 2^32
	);
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_uint32;
CREATE FUNCTION _pb_util_reinterpret_uint64_as_uint32(value BIGINT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	RETURN value;
END $$

DROP FUNCTION IF EXISTS _pb_util_zigzag_decode;
CREATE FUNCTION _pb_util_zigzag_decode(value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	RETURN (value >> 1) ^ - (value & 1);
END $$

DROP FUNCTION IF EXISTS _pb_util_zigzag_encode_64;
CREATE FUNCTION _pb_util_zigzag_encode_64(value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	RETURN (value << 1) ^ (value >> 63);
END $$

DROP FUNCTION IF EXISTS _pb_util_swap_endian_32;
CREATE FUNCTION _pb_util_swap_endian_32(value INT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	RETURN ((value & 0xff) << 24)
		| ((value >> 8) & 0xff) << 16
		| ((value >> 16) & 0xff) << 8
		| ((value >> 24) & 0xff);
END $$

DROP FUNCTION IF EXISTS _pb_util_swap_endian_64;
CREATE FUNCTION _pb_util_swap_endian_64(value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	RETURN ((value & 0xff) << 56)
		| ((value >> 8) & 0xff) << 48
		| ((value >> 16) & 0xff) << 40
		| ((value >> 24) & 0xff) << 32
		| ((value >> 32) & 0xff) << 24
		| ((value >> 40) & 0xff) << 16
		| ((value >> 48) & 0xff) << 8
		| ((value >> 56) & 0xff);
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_sint64;
CREATE FUNCTION _pb_util_reinterpret_uint64_as_sint64(value BIGINT UNSIGNED) RETURNS BIGINT DETERMINISTIC
BEGIN
	RETURN _pb_util_reinterpret_uint64_as_int64(_pb_util_zigzag_decode(value));
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_double $$
CREATE FUNCTION _pb_util_reinterpret_uint64_as_double(bits BIGINT UNSIGNED) RETURNS DOUBLE DETERMINISTIC
BEGIN
	DECLARE sign INT;
	DECLARE exponent INT;
	DECLARE fraction DOUBLE;

	SET sign = IF(bits >> 63 = 0, 1, -1); -- sign: +1 or -1
	SET exponent = (bits >> 52) & 0x7FF; -- exponent (11 bits)
	SET fraction = bits & 0xFFFFFFFFFFFFF; -- fraction (52 bits)

	IF exponent = 2047 THEN -- special case
		IF fraction = 0 THEN
			RETURN sign * NULL;  -- +Inf or -Inf
		ELSE
			RETURN NULL; -- NaN
		END IF;
	ELSEIF exponent = 0 THEN -- subnormal number
		RETURN sign * POW(2, -1022) * (fraction / POW(2, 52));
	ELSE -- normal number
		RETURN sign * POW(2, exponent - 1023) * (1 + (fraction / POW(2, 52)));
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint32_as_float $$
CREATE FUNCTION _pb_util_reinterpret_uint32_as_float(bits INT UNSIGNED) RETURNS FLOAT DETERMINISTIC
BEGIN
    DECLARE sign INT;
    DECLARE exponent INT;
    DECLARE fraction DOUBLE;

    SET sign = IF(bits >> 31 = 0, 1, -1); -- sign: +1 or -1
    SET exponent = (bits >> 23) & 0xFF; -- exponent (8 bits)
    SET fraction = bits & 0x7FFFFF; -- fraction (23 bits)

    IF exponent = 255 THEN -- special case
        IF fraction = 0 THEN
            RETURN sign * NULL; -- +Inf or -Inf
        ELSE
            RETURN NULL; -- NaN
        END IF;
    ELSEIF exponent = 0 THEN -- subnormal number
        RETURN sign * POW(2, -126) * (fraction / POW(2, 23));
    ELSE -- normal number
        RETURN sign * POW(2, exponent - 127) * (1 + (fraction / POW(2, 23)));
    END IF;
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
CREATE PROCEDURE _pb_message_to_wire_json(IN buf LONGBLOB, OUT wire_json JSON)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag INT;
	DECLARE field_number INT;
	DECLARE wire_type INT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE i INT;
	DECLARE json_path TEXT;
	DECLARE wire_element JSON;

	SET wire_json = JSON_OBJECT();
	SET tail = buf;
	SET i = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);

		SET field_number = _pb_wire_get_field_number_from_tag(tag);
		SET wire_type = _pb_wire_get_wire_type_from_tag(tag);

		CASE wire_type
		WHEN 0 THEN -- VARINT
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET wire_element = JSON_OBJECT('i', i, 'n', field_number, 't', wire_type, 'v', uint_value);
		WHEN 1 THEN -- I64
			CALL _pb_wire_read_i64_as_uint64(tail, uint_value, tail);
			SET wire_element = JSON_OBJECT('i', i, 'n', field_number, 't', wire_type, 'v', uint_value);
		WHEN 2 THEN -- LEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			SET wire_element = JSON_OBJECT('i', i, 'n', field_number, 't', wire_type, 'v', TO_BASE64(bytes_value));
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
END $$

DROP FUNCTION IF EXISTS pb_message_to_wire_json $$
CREATE FUNCTION pb_message_to_wire_json(buf LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE wire_json JSON;
	CALL _pb_message_to_wire_json(buf, wire_json);
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

DROP FUNCTION IF EXISTS _pb_util_cast_uint64_as_uint32 $$
CREATE FUNCTION _pb_util_cast_uint64_as_uint32(value BIGINT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	RETURN value;
END $$

DROP FUNCTION IF EXISTS _pb_util_cast_int64_as_int32 $$
CREATE FUNCTION _pb_util_cast_int64_as_int32(value BIGINT) RETURNS INT DETERMINISTIC
BEGIN
	RETURN value;
END $$
