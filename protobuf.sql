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
	CALL pb_wire_read_varint_as_uint64(tail, len, tail);

	IF LENGTH(tail) < len THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'pb_wire_read_len_type: Unexpected end of BLOB.';
	END IF;

	SET value = LEFT(tail, len);
	SET tail = SUBSTRING(tail, len + 1);
END $$

DROP PROCEDURE IF EXISTS pb_wire_read_varint_as_uint64 $$
CREATE PROCEDURE pb_wire_read_varint_as_uint64(IN buf BLOB, OUT value BIGINT UNSIGNED, OUT tail BLOB)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE head INT;
	DECLARE byte_index INT;

	SET value = 0;
	SET tail = buf;
	SET byte_index = 0;

	l1: LOOP
		IF LENGTH(tail) = 0 THEN
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'pb_wire_read_varint_as_uint64: Unexpected end of BLOB.';
		END IF;

		SET head = _pb_util_bin_as_int32(LEFT(tail, 1));
		SET tail = SUBSTRING(tail, 2);

		SET value = value + ((head & 0x7f) << (7 * byte_index));

		IF (head & 0x80) = 0 THEN
			LEAVE l1;
		END IF;

		SET byte_index = byte_index + 1;
		IF byte_index > 10 THEN
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = 'pb_wire_read_varint_as_uint64: Varint cannot exceed 10 bytes.';
		END IF;
	END LOOP;
END $$

DROP FUNCTION IF EXISTS pb_wire_read_varint_as_uint64 $$
CREATE FUNCTION pb_wire_read_varint_as_uint64(buf BLOB) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE tail BLOB;
	DECLARE value BIGINT;
	CALL pb_wire_read_varint_as_uint64(buf, value, tail);
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
CREATE PROCEDURE _pb_message_get_len_type_field(IN buf BLOB, IN field_number INT, IN repeated_index INT, OUT value BLOB, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail BLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value BLOB;
	DECLARE message_text TEXT;

	SET value = _binary X''; -- proto3 default value for an string, bytes, and message field
	SET tail = buf;
	SET field_count = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL pb_wire_read_varint_as_uint64(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 2 /* VARINT */ THEN
			SET message_text = CONCAT('_pb_message_get_len_type_field: string or bytes value cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL pb_wire_read_varint_as_uint64(tail, uint_value, tail);
		WHEN 1 THEN -- I64
			CALL pb_wire_read_i64(tail, bytes_value, tail);
		WHEN 2 THEN -- LEN
			CALL pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				IF repeated_index IS NULL OR repeated_index = field_count THEN
					SET value = bytes_value;
				END IF;
				SET field_count = field_count + 1;
			END IF;
		WHEN 5 THEN -- I32
			CALL pb_wire_read_i32(tail, bytes_value, tail);
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
CREATE PROCEDURE _pb_message_get_i32_field_as_uint32(IN buf BLOB, IN field_number INT, IN repeated_index INT, OUT value INT UNSIGNED, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail BLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE packed_value BLOB;
	DECLARE bytes_value BLOB;
	DECLARE message_text TEXT;

	SET value = 0; -- proto3 default value for an integer field
	SET tail = buf;
	SET field_count = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL pb_wire_read_varint_as_uint64(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 5 /* I32 */ AND (repeated_index IS NULL OR _pb_wire_get_wire_type_from_tag(tag) <> 2 /* LEN */) THEN
			SET message_text = CONCAT('_pb_message_get_i32_field_as_uint32: I32 value cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL pb_wire_read_varint_as_uint64(tail, uint_value, tail);
		WHEN 1 THEN -- I64
			CALL pb_wire_read_i64(tail, bytes_value, tail);
		WHEN 2 THEN -- LEN
			CALL pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number AND repeated_index IS NOT NULL THEN
				WHILE LENGTH(bytes_value) <> 0 DO
					CALL pb_wire_read_i32(bytes_value, packed_value, bytes_value);
					IF repeated_index = field_count THEN
						SET value = _pb_util_swap_endian_32(_pb_util_bin_as_uint32(packed_value));
					END IF;
					SET field_count = field_count + 1;
				END WHILE;
			END IF;
		WHEN 5 THEN -- I32
			CALL pb_wire_read_i32(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				IF repeated_index IS NULL OR repeated_index = field_count THEN
					SET value = _pb_util_swap_endian_32(_pb_util_bin_as_uint32(bytes_value));
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
CREATE PROCEDURE _pb_message_get_i64_field_as_uint64(IN buf BLOB, IN field_number INT, IN repeated_index INT, OUT value BIGINT UNSIGNED, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail BLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE packed_value BLOB;
	DECLARE bytes_value BLOB;
	DECLARE message_text TEXT;

	SET value = 0; -- proto3 default value for an integer field
	SET tail = buf;
	SET field_count = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL pb_wire_read_varint_as_uint64(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 1 /* I64 */ AND (repeated_index IS NULL OR _pb_wire_get_wire_type_from_tag(tag) <> 2 /* LEN */) THEN
			SET message_text = CONCAT('_pb_message_get_i64_field_as_uint64: I64 value cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL pb_wire_read_varint_as_uint64(tail, uint_value, tail);
		WHEN 1 THEN -- I64
			CALL pb_wire_read_i64(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				IF repeated_index IS NULL OR repeated_index = field_count THEN
					SET value = _pb_util_swap_endian_64(_pb_util_bin_as_uint64(bytes_value));
				END IF;
				SET field_count = field_count + 1;
			END IF;
		WHEN 2 THEN -- LEN
			CALL pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number AND repeated_index IS NOT NULL THEN
				WHILE LENGTH(bytes_value) <> 0 DO
					CALL pb_wire_read_i64(bytes_value, packed_value, bytes_value);
					IF repeated_index = field_count THEN
						SET value = _pb_util_swap_endian_64(_pb_util_bin_as_uint64(packed_value));
					END IF;
					SET field_count = field_count + 1;
				END WHILE;
			END IF;
		WHEN 5 THEN -- I32
			CALL pb_wire_read_i32(tail, bytes_value, tail);
		ELSE
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_i64_field_as_uint64: unsupported wire_type';
		END CASE;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_i64_field_as_uint64: repeated index out of range';
	END IF;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_uint32_or_uint64_field $$
CREATE PROCEDURE _pb_message_get_uint32_or_uint64_field(IN buf BLOB, IN field_number INT, IN repeated_index INT, OUT value BIGINT UNSIGNED, OUT field_count INT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE tag BIGINT;
	DECLARE tail BLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value BLOB;
	DECLARE message_text TEXT;

	SET value = 0; -- proto3 default value for an integer field
	SET tail = buf;
	SET field_count = 0;

	WHILE LENGTH(tail) <> 0 DO
		CALL pb_wire_read_varint_as_uint64(tail, tag, tail);

		IF _pb_wire_get_field_number_from_tag(tag) = field_number AND _pb_wire_get_wire_type_from_tag(tag) <> 0 /* VARINT */ AND (repeated_index IS NULL OR _pb_wire_get_wire_type_from_tag(tag) <> 2 /* LEN */) THEN
			SET message_text = CONCAT('_pb_message_get_uint32_or_uint64_field: uint32 or uint64 value cannot be parsed from ', _pb_wire_type_name(_pb_wire_get_wire_type_from_tag(tag)), ' wire type.');
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
		END IF;

		CASE _pb_wire_get_wire_type_from_tag(tag)
		WHEN 0 THEN -- VARINT
			CALL pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number THEN
				IF repeated_index IS NULL OR repeated_index = field_count THEN
					SET value = uint_value;
				END IF;
				SET field_count = field_count + 1;
			END IF;
		WHEN 1 THEN -- I64
			CALL pb_wire_read_i64(tail, bytes_value, tail);
		WHEN 2 THEN -- LEN
			CALL pb_wire_read_len_type(tail, bytes_value, tail);
			IF _pb_wire_get_field_number_from_tag(tag) = field_number AND repeated_index IS NOT NULL THEN
				WHILE LENGTH(bytes_value) <> 0 DO
					CALL pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
					IF repeated_index = field_count THEN
						SET value = uint_value;
					END IF;
					SET field_count = field_count + 1;
				END WHILE;
			END IF;
		WHEN 5 THEN -- I32
			CALL pb_wire_read_i32(tail, bytes_value, tail);
		ELSE
			SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_uint32_or_uint64_field: unsupported wire_type';
		END CASE;
	END WHILE;

	-- Negative repeated_index is used when just counting the number of repeated elements.
	IF repeated_index IS NOT NULL AND repeated_index >= 0 AND field_count <= repeated_index THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_message_get_uint32_or_uint64_field: repeated index out of range';
	END IF;
END $$

DROP FUNCTION IF EXISTS pb_message_get_bool_field $$
CREATE FUNCTION pb_message_get_bool_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value, field_count);
	RETURN value <> 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_bool_field_count $$
CREATE FUNCTION pb_message_get_bool_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_bool_field $$
CREATE FUNCTION pb_message_has_bool_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_enum_field $$
CREATE FUNCTION pb_message_get_enum_field(buf BLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value, field_count);
	RETURN value <> 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_enum_field_count $$
CREATE FUNCTION pb_message_get_enum_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_enum_field $$
CREATE FUNCTION pb_message_has_enum_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_int32_field $$
CREATE FUNCTION pb_message_get_int32_field(buf BLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_int32_field_count $$
CREATE FUNCTION pb_message_get_int32_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_int32_field $$
CREATE FUNCTION pb_message_has_int32_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_uint32_field $$
CREATE FUNCTION pb_message_get_uint32_field(buf BLOB, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_uint32_field_count $$
CREATE FUNCTION pb_message_get_uint32_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_uint32_field $$
CREATE FUNCTION pb_message_has_uint32_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_sint32_field $$
CREATE FUNCTION pb_message_get_sint32_field(buf BLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_sint32_field_count $$
CREATE FUNCTION pb_message_get_sint32_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_sint32_field $$
CREATE FUNCTION pb_message_has_sint32_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_fixed32_field $$
CREATE FUNCTION pb_message_get_fixed32_field(buf BLOB, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(buf, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_fixed32_field_count $$
CREATE FUNCTION pb_message_get_fixed32_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_fixed32_field $$
CREATE FUNCTION pb_message_has_fixed32_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_sfixed32_field $$
CREATE FUNCTION pb_message_get_sfixed32_field(buf BLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(buf, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint32_as_int32(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_sfixed32_field_count $$
CREATE FUNCTION pb_message_get_sfixed32_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_sfixed32_field $$
CREATE FUNCTION pb_message_has_sfixed32_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_float_field $$
CREATE FUNCTION pb_message_get_float_field(buf BLOB, field_number INT, repeated_index INT) RETURNS float DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(buf, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint32_as_float(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_float_field_count $$
CREATE FUNCTION pb_message_get_float_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_float_field $$
CREATE FUNCTION pb_message_has_float_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_int64_field $$
CREATE FUNCTION pb_message_get_int64_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_int64_field_count $$
CREATE FUNCTION pb_message_get_int64_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_int64_field $$
CREATE FUNCTION pb_message_has_int64_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_uint64_field $$
CREATE FUNCTION pb_message_get_uint64_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_uint64_field_count $$
CREATE FUNCTION pb_message_get_uint64_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_uint64_field $$
CREATE FUNCTION pb_message_has_uint64_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_sint64_field $$
CREATE FUNCTION pb_message_get_sint64_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_sint64_field_count $$
CREATE FUNCTION pb_message_get_sint64_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_sint64_field $$
CREATE FUNCTION pb_message_has_sint64_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_uint32_or_uint64_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_fixed64_field $$
CREATE FUNCTION pb_message_get_fixed64_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(buf, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_fixed64_field_count $$
CREATE FUNCTION pb_message_get_fixed64_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_fixed64_field $$
CREATE FUNCTION pb_message_has_fixed64_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_sfixed64_field $$
CREATE FUNCTION pb_message_get_sfixed64_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(buf, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_sfixed64_field_count $$
CREATE FUNCTION pb_message_get_sfixed64_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_sfixed64_field $$
CREATE FUNCTION pb_message_has_sfixed64_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_double_field $$
CREATE FUNCTION pb_message_get_double_field(buf BLOB, field_number INT, repeated_index INT) RETURNS DOUBLE DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(buf, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_double(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_double_field_count $$
CREATE FUNCTION pb_message_get_double_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_double_field $$
CREATE FUNCTION pb_message_has_double_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_string_field $$
CREATE FUNCTION pb_message_get_string_field(buf BLOB, field_number INT, repeated_index INT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(buf, field_number, repeated_index, value, field_count);
	RETURN CONVERT(value USING utf8mb4);
END $$

DROP FUNCTION IF EXISTS pb_message_get_string_field_count $$
CREATE FUNCTION pb_message_get_string_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_string_field $$
CREATE FUNCTION pb_message_has_string_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_bytes_field $$
CREATE FUNCTION pb_message_get_bytes_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BLOB DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(buf, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_bytes_field_count $$
CREATE FUNCTION pb_message_get_bytes_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_bytes_field $$
CREATE FUNCTION pb_message_has_bytes_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_message_field $$
CREATE FUNCTION pb_message_get_message_field(buf BLOB, field_number INT, repeated_index INT) RETURNS BLOB DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(buf, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_message_field_count $$
CREATE FUNCTION pb_message_get_message_field_count(buf BLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(buf, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_has_message_field $$
CREATE FUNCTION pb_message_has_message_field(buf BLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(buf, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$
