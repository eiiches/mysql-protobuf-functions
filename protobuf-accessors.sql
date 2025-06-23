DELIMITER $$

DROP FUNCTION IF EXISTS pb_message_get_int32_field $$
CREATE FUNCTION pb_message_get_int32_field(message LONGBLOB, field_number INT, default_value INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_has_int32_field $$
CREATE FUNCTION pb_message_has_int32_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_int32_field $$
CREATE FUNCTION pb_message_get_repeated_int32_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_int32_field_count $$
CREATE FUNCTION pb_message_get_repeated_int32_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_int64_field $$
CREATE FUNCTION pb_message_get_int64_field(message LONGBLOB, field_number INT, default_value BIGINT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_has_int64_field $$
CREATE FUNCTION pb_message_has_int64_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_int64_field $$
CREATE FUNCTION pb_message_get_repeated_int64_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_int64_field_count $$
CREATE FUNCTION pb_message_get_repeated_int64_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_uint32_field $$
CREATE FUNCTION pb_message_get_uint32_field(message LONGBLOB, field_number INT, default_value INT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_has_uint32_field $$
CREATE FUNCTION pb_message_has_uint32_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_uint32_field $$
CREATE FUNCTION pb_message_get_repeated_uint32_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_uint32_field_count $$
CREATE FUNCTION pb_message_get_repeated_uint32_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_uint64_field $$
CREATE FUNCTION pb_message_get_uint64_field(message LONGBLOB, field_number INT, default_value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_has_uint64_field $$
CREATE FUNCTION pb_message_has_uint64_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_uint64_field $$
CREATE FUNCTION pb_message_get_repeated_uint64_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_uint64_field_count $$
CREATE FUNCTION pb_message_get_repeated_uint64_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_sint32_field $$
CREATE FUNCTION pb_message_get_sint32_field(message LONGBLOB, field_number INT, default_value INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_has_sint32_field $$
CREATE FUNCTION pb_message_has_sint32_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sint32_field $$
CREATE FUNCTION pb_message_get_repeated_sint32_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sint32_field_count $$
CREATE FUNCTION pb_message_get_repeated_sint32_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_sint64_field $$
CREATE FUNCTION pb_message_get_sint64_field(message LONGBLOB, field_number INT, default_value BIGINT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_has_sint64_field $$
CREATE FUNCTION pb_message_has_sint64_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sint64_field $$
CREATE FUNCTION pb_message_get_repeated_sint64_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sint64_field_count $$
CREATE FUNCTION pb_message_get_repeated_sint64_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_enum_field $$
CREATE FUNCTION pb_message_get_enum_field(message LONGBLOB, field_number INT, default_value INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_has_enum_field $$
CREATE FUNCTION pb_message_has_enum_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_enum_field $$
CREATE FUNCTION pb_message_get_repeated_enum_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_enum_field_count $$
CREATE FUNCTION pb_message_get_repeated_enum_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_bool_field $$
CREATE FUNCTION pb_message_get_bool_field(message LONGBLOB, field_number INT, default_value BOOLEAN) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value <> 0;
END $$

DROP FUNCTION IF EXISTS pb_message_has_bool_field $$
CREATE FUNCTION pb_message_has_bool_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_bool_field $$
CREATE FUNCTION pb_message_get_repeated_bool_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN value <> 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_bool_field_count $$
CREATE FUNCTION pb_message_get_repeated_bool_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_varint_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_fixed32_field $$
CREATE FUNCTION pb_message_get_fixed32_field(message LONGBLOB, field_number INT, default_value INT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_has_fixed32_field $$
CREATE FUNCTION pb_message_has_fixed32_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_fixed32_field $$
CREATE FUNCTION pb_message_get_repeated_fixed32_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_fixed32_field_count $$
CREATE FUNCTION pb_message_get_repeated_fixed32_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_sfixed32_field $$
CREATE FUNCTION pb_message_get_sfixed32_field(message LONGBLOB, field_number INT, default_value INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint32_as_int32(value);
END $$

DROP FUNCTION IF EXISTS pb_message_has_sfixed32_field $$
CREATE FUNCTION pb_message_has_sfixed32_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sfixed32_field $$
CREATE FUNCTION pb_message_get_repeated_sfixed32_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint32_as_int32(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sfixed32_field_count $$
CREATE FUNCTION pb_message_get_repeated_sfixed32_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_float_field $$
CREATE FUNCTION pb_message_get_float_field(message LONGBLOB, field_number INT, default_value FLOAT) RETURNS FLOAT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint32_as_float(value);
END $$

DROP FUNCTION IF EXISTS pb_message_has_float_field $$
CREATE FUNCTION pb_message_has_float_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_float_field $$
CREATE FUNCTION pb_message_get_repeated_float_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS FLOAT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint32_as_float(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_float_field_count $$
CREATE FUNCTION pb_message_get_repeated_float_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i32_field_as_uint32(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_fixed64_field $$
CREATE FUNCTION pb_message_get_fixed64_field(message LONGBLOB, field_number INT, default_value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_has_fixed64_field $$
CREATE FUNCTION pb_message_has_fixed64_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_fixed64_field $$
CREATE FUNCTION pb_message_get_repeated_fixed64_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_fixed64_field_count $$
CREATE FUNCTION pb_message_get_repeated_fixed64_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_sfixed64_field $$
CREATE FUNCTION pb_message_get_sfixed64_field(message LONGBLOB, field_number INT, default_value BIGINT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_has_sfixed64_field $$
CREATE FUNCTION pb_message_has_sfixed64_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sfixed64_field $$
CREATE FUNCTION pb_message_get_repeated_sfixed64_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sfixed64_field_count $$
CREATE FUNCTION pb_message_get_repeated_sfixed64_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_double_field $$
CREATE FUNCTION pb_message_get_double_field(message LONGBLOB, field_number INT, default_value DOUBLE) RETURNS DOUBLE DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_double(value);
END $$

DROP FUNCTION IF EXISTS pb_message_has_double_field $$
CREATE FUNCTION pb_message_has_double_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_double_field $$
CREATE FUNCTION pb_message_get_repeated_double_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS DOUBLE DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_double(value);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_double_field_count $$
CREATE FUNCTION pb_message_get_repeated_double_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_message_get_i64_field_as_uint64(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_bytes_field $$
CREATE FUNCTION pb_message_get_bytes_field(message LONGBLOB, field_number INT, default_value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_has_bytes_field $$
CREATE FUNCTION pb_message_has_bytes_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_bytes_field $$
CREATE FUNCTION pb_message_get_repeated_bytes_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_bytes_field_count $$
CREATE FUNCTION pb_message_get_repeated_bytes_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_string_field $$
CREATE FUNCTION pb_message_get_string_field(message LONGBLOB, field_number INT, default_value LONGTEXT) RETURNS LONGTEXT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN CONVERT(value USING utf8mb4);
END $$

DROP FUNCTION IF EXISTS pb_message_has_string_field $$
CREATE FUNCTION pb_message_has_string_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_string_field $$
CREATE FUNCTION pb_message_get_repeated_string_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGTEXT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, repeated_index, value, field_count);
	RETURN CONVERT(value USING utf8mb4);
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_string_field_count $$
CREATE FUNCTION pb_message_get_repeated_string_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_message_get_message_field $$
CREATE FUNCTION pb_message_get_message_field(message LONGBLOB, field_number INT, default_value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_has_message_field $$
CREATE FUNCTION pb_message_has_message_field(message LONGBLOB, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_message_field $$
CREATE FUNCTION pb_message_get_repeated_message_field(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_message_field_count $$
CREATE FUNCTION pb_message_get_repeated_message_field_count(message LONGBLOB, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_message_get_len_type_field(message, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_int32_field $$
CREATE FUNCTION pb_wire_json_get_int32_field(wire_json JSON, field_number INT, default_value INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_int32_field $$
CREATE FUNCTION pb_wire_json_has_int32_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_int32_field $$
CREATE FUNCTION pb_wire_json_get_repeated_int32_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_int32_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_int32_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_int64_field $$
CREATE FUNCTION pb_wire_json_get_int64_field(wire_json JSON, field_number INT, default_value BIGINT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_int64_field $$
CREATE FUNCTION pb_wire_json_has_int64_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_int64_field $$
CREATE FUNCTION pb_wire_json_get_repeated_int64_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_int64_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_int64_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_uint32_field $$
CREATE FUNCTION pb_wire_json_get_uint32_field(wire_json JSON, field_number INT, default_value INT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_uint32_field $$
CREATE FUNCTION pb_wire_json_has_uint32_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_uint32_field $$
CREATE FUNCTION pb_wire_json_get_repeated_uint32_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_uint32_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_uint32_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_uint64_field $$
CREATE FUNCTION pb_wire_json_get_uint64_field(wire_json JSON, field_number INT, default_value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_uint64_field $$
CREATE FUNCTION pb_wire_json_has_uint64_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_uint64_field $$
CREATE FUNCTION pb_wire_json_get_repeated_uint64_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_uint64_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_uint64_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_sint32_field $$
CREATE FUNCTION pb_wire_json_get_sint32_field(wire_json JSON, field_number INT, default_value INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_sint32_field $$
CREATE FUNCTION pb_wire_json_has_sint32_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sint32_field $$
CREATE FUNCTION pb_wire_json_get_repeated_sint32_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sint32_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_sint32_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_sint64_field $$
CREATE FUNCTION pb_wire_json_get_sint64_field(wire_json JSON, field_number INT, default_value BIGINT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_sint64_field $$
CREATE FUNCTION pb_wire_json_has_sint64_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sint64_field $$
CREATE FUNCTION pb_wire_json_get_repeated_sint64_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_sint64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sint64_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_sint64_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_enum_field $$
CREATE FUNCTION pb_wire_json_get_enum_field(wire_json JSON, field_number INT, default_value INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_enum_field $$
CREATE FUNCTION pb_wire_json_has_enum_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_enum_field $$
CREATE FUNCTION pb_wire_json_get_repeated_enum_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_enum_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_enum_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_bool_field $$
CREATE FUNCTION pb_wire_json_get_bool_field(wire_json JSON, field_number INT, default_value BOOLEAN) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value <> 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_bool_field $$
CREATE FUNCTION pb_wire_json_has_bool_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_bool_field $$
CREATE FUNCTION pb_wire_json_get_repeated_bool_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN value <> 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_bool_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_bool_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_varint_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_fixed32_field $$
CREATE FUNCTION pb_wire_json_get_fixed32_field(wire_json JSON, field_number INT, default_value INT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_fixed32_field $$
CREATE FUNCTION pb_wire_json_has_fixed32_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_fixed32_field $$
CREATE FUNCTION pb_wire_json_get_repeated_fixed32_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_fixed32_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_fixed32_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_sfixed32_field $$
CREATE FUNCTION pb_wire_json_get_sfixed32_field(wire_json JSON, field_number INT, default_value INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint32_as_int32(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_sfixed32_field $$
CREATE FUNCTION pb_wire_json_has_sfixed32_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sfixed32_field $$
CREATE FUNCTION pb_wire_json_get_repeated_sfixed32_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint32_as_int32(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sfixed32_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_sfixed32_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_float_field $$
CREATE FUNCTION pb_wire_json_get_float_field(wire_json JSON, field_number INT, default_value FLOAT) RETURNS FLOAT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint32_as_float(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_float_field $$
CREATE FUNCTION pb_wire_json_has_float_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_float_field $$
CREATE FUNCTION pb_wire_json_get_repeated_float_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS FLOAT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint32_as_float(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_float_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_float_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value INT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i32_field_as_uint32(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_fixed64_field $$
CREATE FUNCTION pb_wire_json_get_fixed64_field(wire_json JSON, field_number INT, default_value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_fixed64_field $$
CREATE FUNCTION pb_wire_json_has_fixed64_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_fixed64_field $$
CREATE FUNCTION pb_wire_json_get_repeated_fixed64_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_fixed64_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_fixed64_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_sfixed64_field $$
CREATE FUNCTION pb_wire_json_get_sfixed64_field(wire_json JSON, field_number INT, default_value BIGINT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_sfixed64_field $$
CREATE FUNCTION pb_wire_json_has_sfixed64_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sfixed64_field $$
CREATE FUNCTION pb_wire_json_get_repeated_sfixed64_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_int64(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sfixed64_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_sfixed64_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_double_field $$
CREATE FUNCTION pb_wire_json_get_double_field(wire_json JSON, field_number INT, default_value DOUBLE) RETURNS DOUBLE DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN _pb_util_reinterpret_uint64_as_double(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_double_field $$
CREATE FUNCTION pb_wire_json_has_double_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_double_field $$
CREATE FUNCTION pb_wire_json_get_repeated_double_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS DOUBLE DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, repeated_index, value, field_count);
	RETURN _pb_util_reinterpret_uint64_as_double(value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_double_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_double_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value BIGINT UNSIGNED;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_i64_field_as_uint64(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_bytes_field $$
CREATE FUNCTION pb_wire_json_get_bytes_field(wire_json JSON, field_number INT, default_value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_bytes_field $$
CREATE FUNCTION pb_wire_json_has_bytes_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_bytes_field $$
CREATE FUNCTION pb_wire_json_get_repeated_bytes_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_bytes_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_bytes_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_string_field $$
CREATE FUNCTION pb_wire_json_get_string_field(wire_json JSON, field_number INT, default_value LONGTEXT) RETURNS LONGTEXT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN CONVERT(value USING utf8mb4);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_string_field $$
CREATE FUNCTION pb_wire_json_has_string_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_string_field $$
CREATE FUNCTION pb_wire_json_get_repeated_string_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS LONGTEXT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, repeated_index, value, field_count);
	RETURN CONVERT(value USING utf8mb4);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_string_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_string_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_message_field $$
CREATE FUNCTION pb_wire_json_get_message_field(wire_json JSON, field_number INT, default_value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, NULL, value, field_count);
	IF field_count = 0 THEN
		RETURN default_value;
	END IF;
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_has_message_field $$
CREATE FUNCTION pb_wire_json_has_message_field(wire_json JSON, field_number INT) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, NULL, value, field_count);
	RETURN field_count > 0;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_message_field $$
CREATE FUNCTION pb_wire_json_get_repeated_message_field(wire_json JSON, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, repeated_index, value, field_count);
	RETURN value;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_message_field_count $$
CREATE FUNCTION pb_wire_json_get_repeated_message_field_count(wire_json JSON, field_number INT) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE value LONGBLOB;
	DECLARE field_count INT;
	CALL _pb_wire_json_get_len_type_field(wire_json, field_number, -1, value, field_count);
	RETURN field_count;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_int32_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_int32_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 0 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_int64(uint_value)));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_int64(uint_value)));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_int32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_int32_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_int32_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_int32_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_int32_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_int32_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_int64(uint_value)));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_int64(uint_value)));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_int32_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_int32_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_int32_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_int32_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_uint32_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_uint32_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 0 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(uint_value));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_uint32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_uint32_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_uint32_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_uint32_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_uint32_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_uint32_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(uint_value));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_uint32_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_uint32_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_uint32_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_uint32_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_int64_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_int64_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 0 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_int64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_int64_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_int64_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_int64_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_int64_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_int64_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_int64_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_int64_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_int64_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_int64_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_uint64_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_uint64_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 0 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_uint64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_uint64_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_uint64_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_uint64_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_uint64_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_uint64_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_uint64_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_uint64_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_uint64_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_uint64_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_sint32_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_sint32_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 0 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_sint64(uint_value)));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_sint64(uint_value)));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sint32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sint32_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_sint32_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_sint32_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_sint32_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_sint32_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_sint64(uint_value)));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_sint64(uint_value)));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_sint32_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sint32_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_sint32_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_sint32_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_sint64_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_sint64_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 0 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_sint64(uint_value));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_sint64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sint64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sint64_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_sint64_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_sint64_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_sint64_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_sint64_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_sint64(uint_value));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_sint64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_sint64_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sint64_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_sint64_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_sint64_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_enum_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_enum_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 0 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_enum_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_enum_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_enum_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_enum_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_enum_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_enum_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_enum_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_enum_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_enum_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_enum_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_bool_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_bool_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 0 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value <> 0);
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value <> 0);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_bool_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_bool_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_bool_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_bool_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_bool_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_bool_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value <> 0);
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value <> 0);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_bool_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_bool_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_bool_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_bool_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_fixed32_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_fixed32_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 5 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_fixed32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_fixed32_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_fixed32_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_fixed32_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_fixed32_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_fixed32_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 5 THEN
			CALL _pb_wire_read_i32_as_uint32(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_fixed32_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_fixed32_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_fixed32_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_fixed32_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_sfixed32_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_sfixed32_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 5 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_int32(uint_value));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_int32(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sfixed32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sfixed32_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_sfixed32_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_sfixed32_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_sfixed32_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_sfixed32_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 5 THEN
			CALL _pb_wire_read_i32_as_uint32(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_int32(uint_value));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_int32(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_sfixed32_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sfixed32_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_sfixed32_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_sfixed32_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_float_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_float_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 5 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_float(uint_value));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_float(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_float_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_float_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_float_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_float_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_float_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_float_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 5 THEN
			CALL _pb_wire_read_i32_as_uint32(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_float(uint_value));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_float(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_float_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_float_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_float_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_float_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_fixed64_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_fixed64_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 1 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_fixed64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_fixed64_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_fixed64_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_fixed64_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_fixed64_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_fixed64_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 1 THEN
			CALL _pb_wire_read_i64_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_fixed64_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_fixed64_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_fixed64_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_fixed64_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_sfixed64_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_sfixed64_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 1 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sfixed64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sfixed64_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_sfixed64_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_sfixed64_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_sfixed64_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_sfixed64_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 1 THEN
			CALL _pb_wire_read_i64_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_sfixed64_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sfixed64_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_sfixed64_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_sfixed64_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_double_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_double_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 1 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_double(uint_value));
		WHEN 2 THEN -- LEN
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_double(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_double_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_double_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_double_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_double_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_double_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_double_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 1 THEN
			CALL _pb_wire_read_i64_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_double(uint_value));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_double(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_double_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_double_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_double_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_double_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_bytes_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_bytes_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 2 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', TO_BASE64(bytes_value));
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_bytes_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_bytes_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_bytes_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_bytes_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_bytes_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_bytes_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', TO_BASE64(bytes_value));
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_bytes_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_bytes_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_bytes_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_bytes_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_string_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_string_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 2 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', CONVERT(bytes_value USING utf8mb4));
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_string_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_string_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_string_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_string_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_string_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_string_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', CONVERT(bytes_value USING utf8mb4));
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_string_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_string_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_string_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_string_field_as_json_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_message_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_message_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE done TINYINT DEFAULT FALSE;
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;

	DECLARE cur CURSOR FOR
		SELECT
			jt.wire_type,
			jt.uint_value,
			FROM_BASE64(jt.bytes_value)
		FROM JSON_TABLE(wire_json, '$[*]' COLUMNS (
			field_number INT PATH '$.field_number',
			wire_type INT PATH '$.wire_type',
			uint_value BIGINT UNSIGNED PATH '$.value.uint',
			bytes_value TEXT PATH '$.value.bytes'
		)) AS jt
		WHERE jt.field_number = field_number;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	SET result = JSON_ARRAY();

	OPEN cur;
	l1: LOOP
		FETCH cur INTO wire_type, uint_value, bytes_value;
		IF done THEN
			LEAVE l1;
		END IF;

		CASE wire_type
		WHEN 2 THEN
			SET result = JSON_ARRAY_APPEND(result, '$', TO_BASE64(bytes_value));
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_message_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END LOOP;
	CLOSE cur;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_message_field_as_json_array $$
CREATE FUNCTION pb_wire_json_get_repeated_message_field_as_json_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_message_field_as_json_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_message_field_as_json_array $$
CREATE PROCEDURE _pb_message_get_repeated_message_field_as_json_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE tag BIGINT;
	DECLARE tail LONGBLOB;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE message_text TEXT;
	DECLARE current_field_number INT;
	DECLARE current_wire_type INT;

	SET tail = message;
	SET result = JSON_ARRAY();

	l1: WHILE LENGTH(tail) <> 0 DO
		CALL _pb_wire_read_varint_as_uint64(tail, tag, tail);
		SET current_field_number = _pb_wire_get_field_number_from_tag(tag);
		SET current_wire_type = _pb_wire_get_wire_type_from_tag(tag);

		IF current_field_number != field_number THEN
			CALL _pb_wire_skip(tail, current_wire_type, tail);
			LEAVE l1;
		END IF;

		CASE current_wire_type
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', TO_BASE64(bytes_value));
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_message_field_as_json_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_message_field_as_json_array $$
CREATE FUNCTION pb_message_get_repeated_message_field_as_json_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_message_field_as_json_array(message, field_number, result);
	RETURN result;
END $$
