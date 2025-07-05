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

DROP FUNCTION IF EXISTS pb_message_get_repeated_int32_field_element $$
CREATE FUNCTION pb_message_get_repeated_int32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_int64_field_element $$
CREATE FUNCTION pb_message_get_repeated_int64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_uint32_field_element $$
CREATE FUNCTION pb_message_get_repeated_uint32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_uint64_field_element $$
CREATE FUNCTION pb_message_get_repeated_uint64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_sint32_field_element $$
CREATE FUNCTION pb_message_get_repeated_sint32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_sint64_field_element $$
CREATE FUNCTION pb_message_get_repeated_sint64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_enum_field_element $$
CREATE FUNCTION pb_message_get_repeated_enum_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_bool_field_element $$
CREATE FUNCTION pb_message_get_repeated_bool_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BOOLEAN DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_fixed32_field_element $$
CREATE FUNCTION pb_message_get_repeated_fixed32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_sfixed32_field_element $$
CREATE FUNCTION pb_message_get_repeated_sfixed32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_float_field_element $$
CREATE FUNCTION pb_message_get_repeated_float_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS FLOAT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_fixed64_field_element $$
CREATE FUNCTION pb_message_get_repeated_fixed64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_sfixed64_field_element $$
CREATE FUNCTION pb_message_get_repeated_sfixed64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_double_field_element $$
CREATE FUNCTION pb_message_get_repeated_double_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS DOUBLE DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_bytes_field_element $$
CREATE FUNCTION pb_message_get_repeated_bytes_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_string_field_element $$
CREATE FUNCTION pb_message_get_repeated_string_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGTEXT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_get_repeated_message_field_element $$
CREATE FUNCTION pb_message_get_repeated_message_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_message_set_int32_field $$
CREATE FUNCTION pb_message_set_int32_field(message LONGBLOB, field_number INT, value INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_int32_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_int32_field_element $$
CREATE FUNCTION pb_message_add_repeated_int32_field_element(message LONGBLOB, field_number INT, value INT, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_int32_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_int32_field_element $$
CREATE FUNCTION pb_message_set_repeated_int32_field_element(message LONGBLOB, field_number INT, repeated_index INT, value INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_int32_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_int32_field_element $$
CREATE FUNCTION pb_message_remove_repeated_int32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_int32_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_int32_field $$
CREATE FUNCTION pb_message_clear_int32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_int32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_int32_field $$
CREATE FUNCTION pb_message_clear_repeated_int32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_int32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_int64_field $$
CREATE FUNCTION pb_message_set_int64_field(message LONGBLOB, field_number INT, value BIGINT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_int64_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_int64_field_element $$
CREATE FUNCTION pb_message_add_repeated_int64_field_element(message LONGBLOB, field_number INT, value BIGINT, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_int64_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_int64_field_element $$
CREATE FUNCTION pb_message_set_repeated_int64_field_element(message LONGBLOB, field_number INT, repeated_index INT, value BIGINT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_int64_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_int64_field_element $$
CREATE FUNCTION pb_message_remove_repeated_int64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_int64_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_int64_field $$
CREATE FUNCTION pb_message_clear_int64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_int64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_int64_field $$
CREATE FUNCTION pb_message_clear_repeated_int64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_int64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_uint32_field $$
CREATE FUNCTION pb_message_set_uint32_field(message LONGBLOB, field_number INT, value INT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_uint32_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_uint32_field_element $$
CREATE FUNCTION pb_message_add_repeated_uint32_field_element(message LONGBLOB, field_number INT, value INT UNSIGNED, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_uint32_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_uint32_field_element $$
CREATE FUNCTION pb_message_set_repeated_uint32_field_element(message LONGBLOB, field_number INT, repeated_index INT, value INT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_uint32_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_uint32_field_element $$
CREATE FUNCTION pb_message_remove_repeated_uint32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_uint32_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_uint32_field $$
CREATE FUNCTION pb_message_clear_uint32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_uint32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_uint32_field $$
CREATE FUNCTION pb_message_clear_repeated_uint32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_uint32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_uint64_field $$
CREATE FUNCTION pb_message_set_uint64_field(message LONGBLOB, field_number INT, value BIGINT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_uint64_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_uint64_field_element $$
CREATE FUNCTION pb_message_add_repeated_uint64_field_element(message LONGBLOB, field_number INT, value BIGINT UNSIGNED, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_uint64_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_uint64_field_element $$
CREATE FUNCTION pb_message_set_repeated_uint64_field_element(message LONGBLOB, field_number INT, repeated_index INT, value BIGINT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_uint64_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_uint64_field_element $$
CREATE FUNCTION pb_message_remove_repeated_uint64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_uint64_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_uint64_field $$
CREATE FUNCTION pb_message_clear_uint64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_uint64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_uint64_field $$
CREATE FUNCTION pb_message_clear_repeated_uint64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_uint64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_sint32_field $$
CREATE FUNCTION pb_message_set_sint32_field(message LONGBLOB, field_number INT, value INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_sint32_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_sint32_field_element $$
CREATE FUNCTION pb_message_add_repeated_sint32_field_element(message LONGBLOB, field_number INT, value INT, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_sint32_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_sint32_field_element $$
CREATE FUNCTION pb_message_set_repeated_sint32_field_element(message LONGBLOB, field_number INT, repeated_index INT, value INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_sint32_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_sint32_field_element $$
CREATE FUNCTION pb_message_remove_repeated_sint32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_sint32_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_sint32_field $$
CREATE FUNCTION pb_message_clear_sint32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_sint32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_sint32_field $$
CREATE FUNCTION pb_message_clear_repeated_sint32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_sint32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_sint64_field $$
CREATE FUNCTION pb_message_set_sint64_field(message LONGBLOB, field_number INT, value BIGINT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_sint64_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_sint64_field_element $$
CREATE FUNCTION pb_message_add_repeated_sint64_field_element(message LONGBLOB, field_number INT, value BIGINT, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_sint64_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_sint64_field_element $$
CREATE FUNCTION pb_message_set_repeated_sint64_field_element(message LONGBLOB, field_number INT, repeated_index INT, value BIGINT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_sint64_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_sint64_field_element $$
CREATE FUNCTION pb_message_remove_repeated_sint64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_sint64_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_sint64_field $$
CREATE FUNCTION pb_message_clear_sint64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_sint64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_sint64_field $$
CREATE FUNCTION pb_message_clear_repeated_sint64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_sint64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_enum_field $$
CREATE FUNCTION pb_message_set_enum_field(message LONGBLOB, field_number INT, value INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_enum_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_enum_field_element $$
CREATE FUNCTION pb_message_add_repeated_enum_field_element(message LONGBLOB, field_number INT, value INT, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_enum_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_enum_field_element $$
CREATE FUNCTION pb_message_set_repeated_enum_field_element(message LONGBLOB, field_number INT, repeated_index INT, value INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_enum_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_enum_field_element $$
CREATE FUNCTION pb_message_remove_repeated_enum_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_enum_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_enum_field $$
CREATE FUNCTION pb_message_clear_enum_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_enum_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_enum_field $$
CREATE FUNCTION pb_message_clear_repeated_enum_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_enum_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_bool_field $$
CREATE FUNCTION pb_message_set_bool_field(message LONGBLOB, field_number INT, value BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_bool_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_bool_field_element $$
CREATE FUNCTION pb_message_add_repeated_bool_field_element(message LONGBLOB, field_number INT, value BOOLEAN, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_bool_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_bool_field_element $$
CREATE FUNCTION pb_message_set_repeated_bool_field_element(message LONGBLOB, field_number INT, repeated_index INT, value BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_bool_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_bool_field_element $$
CREATE FUNCTION pb_message_remove_repeated_bool_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_bool_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_bool_field $$
CREATE FUNCTION pb_message_clear_bool_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_bool_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_bool_field $$
CREATE FUNCTION pb_message_clear_repeated_bool_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_bool_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_fixed32_field $$
CREATE FUNCTION pb_message_set_fixed32_field(message LONGBLOB, field_number INT, value INT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_fixed32_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_fixed32_field_element $$
CREATE FUNCTION pb_message_add_repeated_fixed32_field_element(message LONGBLOB, field_number INT, value INT UNSIGNED, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_fixed32_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_fixed32_field_element $$
CREATE FUNCTION pb_message_set_repeated_fixed32_field_element(message LONGBLOB, field_number INT, repeated_index INT, value INT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_fixed32_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_fixed32_field_element $$
CREATE FUNCTION pb_message_remove_repeated_fixed32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_fixed32_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_fixed32_field $$
CREATE FUNCTION pb_message_clear_fixed32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_fixed32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_fixed32_field $$
CREATE FUNCTION pb_message_clear_repeated_fixed32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_fixed32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_sfixed32_field $$
CREATE FUNCTION pb_message_set_sfixed32_field(message LONGBLOB, field_number INT, value INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_sfixed32_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_sfixed32_field_element $$
CREATE FUNCTION pb_message_add_repeated_sfixed32_field_element(message LONGBLOB, field_number INT, value INT, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_sfixed32_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_sfixed32_field_element $$
CREATE FUNCTION pb_message_set_repeated_sfixed32_field_element(message LONGBLOB, field_number INT, repeated_index INT, value INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_sfixed32_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_sfixed32_field_element $$
CREATE FUNCTION pb_message_remove_repeated_sfixed32_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_sfixed32_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_sfixed32_field $$
CREATE FUNCTION pb_message_clear_sfixed32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_sfixed32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_sfixed32_field $$
CREATE FUNCTION pb_message_clear_repeated_sfixed32_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_sfixed32_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_float_field $$
CREATE FUNCTION pb_message_set_float_field(message LONGBLOB, field_number INT, value FLOAT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_float_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_float_field_element $$
CREATE FUNCTION pb_message_add_repeated_float_field_element(message LONGBLOB, field_number INT, value FLOAT, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_float_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_float_field_element $$
CREATE FUNCTION pb_message_set_repeated_float_field_element(message LONGBLOB, field_number INT, repeated_index INT, value FLOAT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_float_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_float_field_element $$
CREATE FUNCTION pb_message_remove_repeated_float_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_float_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_float_field $$
CREATE FUNCTION pb_message_clear_float_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_float_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_float_field $$
CREATE FUNCTION pb_message_clear_repeated_float_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_float_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_fixed64_field $$
CREATE FUNCTION pb_message_set_fixed64_field(message LONGBLOB, field_number INT, value BIGINT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_fixed64_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_fixed64_field_element $$
CREATE FUNCTION pb_message_add_repeated_fixed64_field_element(message LONGBLOB, field_number INT, value BIGINT UNSIGNED, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_fixed64_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_fixed64_field_element $$
CREATE FUNCTION pb_message_set_repeated_fixed64_field_element(message LONGBLOB, field_number INT, repeated_index INT, value BIGINT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_fixed64_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_fixed64_field_element $$
CREATE FUNCTION pb_message_remove_repeated_fixed64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_fixed64_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_fixed64_field $$
CREATE FUNCTION pb_message_clear_fixed64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_fixed64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_fixed64_field $$
CREATE FUNCTION pb_message_clear_repeated_fixed64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_fixed64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_sfixed64_field $$
CREATE FUNCTION pb_message_set_sfixed64_field(message LONGBLOB, field_number INT, value BIGINT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_sfixed64_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_sfixed64_field_element $$
CREATE FUNCTION pb_message_add_repeated_sfixed64_field_element(message LONGBLOB, field_number INT, value BIGINT, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_sfixed64_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_sfixed64_field_element $$
CREATE FUNCTION pb_message_set_repeated_sfixed64_field_element(message LONGBLOB, field_number INT, repeated_index INT, value BIGINT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_sfixed64_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_sfixed64_field_element $$
CREATE FUNCTION pb_message_remove_repeated_sfixed64_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_sfixed64_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_sfixed64_field $$
CREATE FUNCTION pb_message_clear_sfixed64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_sfixed64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_sfixed64_field $$
CREATE FUNCTION pb_message_clear_repeated_sfixed64_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_sfixed64_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_double_field $$
CREATE FUNCTION pb_message_set_double_field(message LONGBLOB, field_number INT, value DOUBLE) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_double_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_double_field_element $$
CREATE FUNCTION pb_message_add_repeated_double_field_element(message LONGBLOB, field_number INT, value DOUBLE, use_packed BOOLEAN) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_double_field_element(pb_message_to_wire_json(message), field_number, value, use_packed));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_double_field_element $$
CREATE FUNCTION pb_message_set_repeated_double_field_element(message LONGBLOB, field_number INT, repeated_index INT, value DOUBLE) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_double_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_double_field_element $$
CREATE FUNCTION pb_message_remove_repeated_double_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_double_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_double_field $$
CREATE FUNCTION pb_message_clear_double_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_double_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_double_field $$
CREATE FUNCTION pb_message_clear_repeated_double_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_double_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_bytes_field $$
CREATE FUNCTION pb_message_set_bytes_field(message LONGBLOB, field_number INT, value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_bytes_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_bytes_field_element $$
CREATE FUNCTION pb_message_add_repeated_bytes_field_element(message LONGBLOB, field_number INT, value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_bytes_field_element(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_bytes_field_element $$
CREATE FUNCTION pb_message_set_repeated_bytes_field_element(message LONGBLOB, field_number INT, repeated_index INT, value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_bytes_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_bytes_field_element $$
CREATE FUNCTION pb_message_remove_repeated_bytes_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_bytes_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_bytes_field $$
CREATE FUNCTION pb_message_clear_bytes_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_bytes_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_bytes_field $$
CREATE FUNCTION pb_message_clear_repeated_bytes_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_bytes_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_string_field $$
CREATE FUNCTION pb_message_set_string_field(message LONGBLOB, field_number INT, value LONGTEXT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_string_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_string_field_element $$
CREATE FUNCTION pb_message_add_repeated_string_field_element(message LONGBLOB, field_number INT, value LONGTEXT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_string_field_element(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_string_field_element $$
CREATE FUNCTION pb_message_set_repeated_string_field_element(message LONGBLOB, field_number INT, repeated_index INT, value LONGTEXT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_string_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_string_field_element $$
CREATE FUNCTION pb_message_remove_repeated_string_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_string_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_string_field $$
CREATE FUNCTION pb_message_clear_string_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_string_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_string_field $$
CREATE FUNCTION pb_message_clear_repeated_string_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_string_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_set_message_field $$
CREATE FUNCTION pb_message_set_message_field(message LONGBLOB, field_number INT, value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_message_field(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_add_repeated_message_field_element $$
CREATE FUNCTION pb_message_add_repeated_message_field_element(message LONGBLOB, field_number INT, value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_add_repeated_message_field_element(pb_message_to_wire_json(message), field_number, value));
END $$

DROP FUNCTION IF EXISTS pb_message_set_repeated_message_field_element $$
CREATE FUNCTION pb_message_set_repeated_message_field_element(message LONGBLOB, field_number INT, repeated_index INT, value LONGBLOB) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_set_repeated_message_field_element(pb_message_to_wire_json(message), field_number, repeated_index, value));
END $$

DROP FUNCTION IF EXISTS pb_message_remove_repeated_message_field_element $$
CREATE FUNCTION pb_message_remove_repeated_message_field_element(message LONGBLOB, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_remove_repeated_message_field_element(pb_message_to_wire_json(message), field_number, repeated_index));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_message_field $$
CREATE FUNCTION pb_message_clear_message_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_message_field(pb_message_to_wire_json(message), field_number));
END $$

DROP FUNCTION IF EXISTS pb_message_clear_repeated_message_field $$
CREATE FUNCTION pb_message_clear_repeated_message_field(message LONGBLOB, field_number INT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN pb_wire_json_to_message(pb_wire_json_clear_repeated_message_field(pb_message_to_wire_json(message), field_number));
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_int32_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_int32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_int64_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_int64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_uint32_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_uint32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_uint64_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_uint64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sint32_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_sint32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sint64_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_sint64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_enum_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_enum_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_bool_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_bool_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS BOOLEAN DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_fixed32_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_fixed32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT UNSIGNED DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sfixed32_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_sfixed32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS INT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_float_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_float_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS FLOAT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_fixed64_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_fixed64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT UNSIGNED DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sfixed64_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_sfixed64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS BIGINT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_double_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_double_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS DOUBLE DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_bytes_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_bytes_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_string_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_string_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS LONGTEXT DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_message_field_element $$
CREATE FUNCTION pb_wire_json_get_repeated_message_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS LONGBLOB DETERMINISTIC
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

DROP FUNCTION IF EXISTS pb_wire_json_set_int32_field $$
CREATE FUNCTION pb_wire_json_set_int32_field(wire_json JSON, field_number INT, value INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_varint_field(wire_json, field_number, _pb_util_reinterpret_int64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_int32_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_int32_field_element(wire_json JSON, field_number INT, value INT, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_varint_field_element(wire_json, field_number, _pb_util_reinterpret_int64_as_uint64(value), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_int32_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_int32_field_element(wire_json JSON, field_number INT, repeated_index INT, value INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_varint_field_element(wire_json, field_number, repeated_index, _pb_util_reinterpret_int64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_int32_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_int32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_varint_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_int32_field $$
CREATE FUNCTION pb_wire_json_clear_int32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_int32_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_int32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_int64_field $$
CREATE FUNCTION pb_wire_json_set_int64_field(wire_json JSON, field_number INT, value BIGINT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_varint_field(wire_json, field_number, _pb_util_reinterpret_int64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_int64_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_int64_field_element(wire_json JSON, field_number INT, value BIGINT, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_varint_field_element(wire_json, field_number, _pb_util_reinterpret_int64_as_uint64(value), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_int64_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_int64_field_element(wire_json JSON, field_number INT, repeated_index INT, value BIGINT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_varint_field_element(wire_json, field_number, repeated_index, _pb_util_reinterpret_int64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_int64_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_int64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_varint_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_int64_field $$
CREATE FUNCTION pb_wire_json_clear_int64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_int64_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_int64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_uint32_field $$
CREATE FUNCTION pb_wire_json_set_uint32_field(wire_json JSON, field_number INT, value INT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_varint_field(wire_json, field_number, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_uint32_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_uint32_field_element(wire_json JSON, field_number INT, value INT UNSIGNED, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_varint_field_element(wire_json, field_number, value, use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_uint32_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_uint32_field_element(wire_json JSON, field_number INT, repeated_index INT, value INT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_varint_field_element(wire_json, field_number, repeated_index, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_uint32_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_uint32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_varint_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_uint32_field $$
CREATE FUNCTION pb_wire_json_clear_uint32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_uint32_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_uint32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_uint64_field $$
CREATE FUNCTION pb_wire_json_set_uint64_field(wire_json JSON, field_number INT, value BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_varint_field(wire_json, field_number, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_uint64_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_uint64_field_element(wire_json JSON, field_number INT, value BIGINT UNSIGNED, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_varint_field_element(wire_json, field_number, value, use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_uint64_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_uint64_field_element(wire_json JSON, field_number INT, repeated_index INT, value BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_varint_field_element(wire_json, field_number, repeated_index, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_uint64_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_uint64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_varint_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_uint64_field $$
CREATE FUNCTION pb_wire_json_clear_uint64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_uint64_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_uint64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_sint32_field $$
CREATE FUNCTION pb_wire_json_set_sint32_field(wire_json JSON, field_number INT, value INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_varint_field(wire_json, field_number, _pb_util_reinterpret_sint64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_sint32_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_sint32_field_element(wire_json JSON, field_number INT, value INT, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_varint_field_element(wire_json, field_number, _pb_util_reinterpret_sint64_as_uint64(value), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_sint32_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_sint32_field_element(wire_json JSON, field_number INT, repeated_index INT, value INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_varint_field_element(wire_json, field_number, repeated_index, _pb_util_reinterpret_sint64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_sint32_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_sint32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_varint_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_sint32_field $$
CREATE FUNCTION pb_wire_json_clear_sint32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_sint32_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_sint32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_sint64_field $$
CREATE FUNCTION pb_wire_json_set_sint64_field(wire_json JSON, field_number INT, value BIGINT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_varint_field(wire_json, field_number, _pb_util_reinterpret_sint64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_sint64_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_sint64_field_element(wire_json JSON, field_number INT, value BIGINT, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_varint_field_element(wire_json, field_number, _pb_util_reinterpret_sint64_as_uint64(value), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_sint64_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_sint64_field_element(wire_json JSON, field_number INT, repeated_index INT, value BIGINT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_varint_field_element(wire_json, field_number, repeated_index, _pb_util_reinterpret_sint64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_sint64_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_sint64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_varint_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_sint64_field $$
CREATE FUNCTION pb_wire_json_clear_sint64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_sint64_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_sint64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_enum_field $$
CREATE FUNCTION pb_wire_json_set_enum_field(wire_json JSON, field_number INT, value INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_varint_field(wire_json, field_number, _pb_util_reinterpret_int64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_enum_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_enum_field_element(wire_json JSON, field_number INT, value INT, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_varint_field_element(wire_json, field_number, _pb_util_reinterpret_int64_as_uint64(value), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_enum_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_enum_field_element(wire_json JSON, field_number INT, repeated_index INT, value INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_varint_field_element(wire_json, field_number, repeated_index, _pb_util_reinterpret_int64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_enum_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_enum_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_varint_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_enum_field $$
CREATE FUNCTION pb_wire_json_clear_enum_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_enum_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_enum_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_bool_field $$
CREATE FUNCTION pb_wire_json_set_bool_field(wire_json JSON, field_number INT, value BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_varint_field(wire_json, field_number, IF(value, 1, 0));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_bool_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_bool_field_element(wire_json JSON, field_number INT, value BOOLEAN, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_varint_field_element(wire_json, field_number, IF(value, 1, 0), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_bool_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_bool_field_element(wire_json JSON, field_number INT, repeated_index INT, value BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_varint_field_element(wire_json, field_number, repeated_index, IF(value, 1, 0));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_bool_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_bool_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_varint_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_bool_field $$
CREATE FUNCTION pb_wire_json_clear_bool_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_bool_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_bool_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_fixed32_field $$
CREATE FUNCTION pb_wire_json_set_fixed32_field(wire_json JSON, field_number INT, value INT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_i32_field(wire_json, field_number, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_fixed32_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_fixed32_field_element(wire_json JSON, field_number INT, value INT UNSIGNED, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_i32_field_element(wire_json, field_number, value, use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_fixed32_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_fixed32_field_element(wire_json JSON, field_number INT, repeated_index INT, value INT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_i32_field_element(wire_json, field_number, repeated_index, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_fixed32_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_fixed32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_i32_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_fixed32_field $$
CREATE FUNCTION pb_wire_json_clear_fixed32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_fixed32_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_fixed32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_sfixed32_field $$
CREATE FUNCTION pb_wire_json_set_sfixed32_field(wire_json JSON, field_number INT, value INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_i32_field(wire_json, field_number, _pb_util_reinterpret_int32_as_uint32(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_sfixed32_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_sfixed32_field_element(wire_json JSON, field_number INT, value INT, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_i32_field_element(wire_json, field_number, _pb_util_reinterpret_int32_as_uint32(value), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_sfixed32_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_sfixed32_field_element(wire_json JSON, field_number INT, repeated_index INT, value INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_i32_field_element(wire_json, field_number, repeated_index, _pb_util_reinterpret_int32_as_uint32(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_sfixed32_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_sfixed32_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_i32_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_sfixed32_field $$
CREATE FUNCTION pb_wire_json_clear_sfixed32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_sfixed32_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_sfixed32_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_float_field $$
CREATE FUNCTION pb_wire_json_set_float_field(wire_json JSON, field_number INT, value FLOAT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_i32_field(wire_json, field_number, _pb_util_reinterpret_float_as_uint32(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_float_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_float_field_element(wire_json JSON, field_number INT, value FLOAT, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_i32_field_element(wire_json, field_number, _pb_util_reinterpret_float_as_uint32(value), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_float_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_float_field_element(wire_json JSON, field_number INT, repeated_index INT, value FLOAT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_i32_field_element(wire_json, field_number, repeated_index, _pb_util_reinterpret_float_as_uint32(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_float_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_float_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_i32_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_float_field $$
CREATE FUNCTION pb_wire_json_clear_float_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_float_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_float_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_fixed64_field $$
CREATE FUNCTION pb_wire_json_set_fixed64_field(wire_json JSON, field_number INT, value BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_i64_field(wire_json, field_number, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_fixed64_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_fixed64_field_element(wire_json JSON, field_number INT, value BIGINT UNSIGNED, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_i64_field_element(wire_json, field_number, value, use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_fixed64_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_fixed64_field_element(wire_json JSON, field_number INT, repeated_index INT, value BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_i64_field_element(wire_json, field_number, repeated_index, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_fixed64_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_fixed64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_i64_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_fixed64_field $$
CREATE FUNCTION pb_wire_json_clear_fixed64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_fixed64_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_fixed64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_sfixed64_field $$
CREATE FUNCTION pb_wire_json_set_sfixed64_field(wire_json JSON, field_number INT, value BIGINT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_i64_field(wire_json, field_number, _pb_util_reinterpret_int64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_sfixed64_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_sfixed64_field_element(wire_json JSON, field_number INT, value BIGINT, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_i64_field_element(wire_json, field_number, _pb_util_reinterpret_int64_as_uint64(value), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_sfixed64_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_sfixed64_field_element(wire_json JSON, field_number INT, repeated_index INT, value BIGINT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_i64_field_element(wire_json, field_number, repeated_index, _pb_util_reinterpret_int64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_sfixed64_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_sfixed64_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_i64_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_sfixed64_field $$
CREATE FUNCTION pb_wire_json_clear_sfixed64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_sfixed64_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_sfixed64_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_double_field $$
CREATE FUNCTION pb_wire_json_set_double_field(wire_json JSON, field_number INT, value DOUBLE) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_i64_field(wire_json, field_number, _pb_util_reinterpret_double_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_double_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_double_field_element(wire_json JSON, field_number INT, value DOUBLE, use_packed BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_i64_field_element(wire_json, field_number, _pb_util_reinterpret_double_as_uint64(value), use_packed);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_double_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_double_field_element(wire_json JSON, field_number INT, repeated_index INT, value DOUBLE) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_i64_field_element(wire_json, field_number, repeated_index, _pb_util_reinterpret_double_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_double_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_double_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_i64_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_double_field $$
CREATE FUNCTION pb_wire_json_clear_double_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_double_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_double_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_bytes_field $$
CREATE FUNCTION pb_wire_json_set_bytes_field(wire_json JSON, field_number INT, value LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_len_field(wire_json, field_number, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_bytes_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_bytes_field_element(wire_json JSON, field_number INT, value LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_len_field_element(wire_json, field_number, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_bytes_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_bytes_field_element(wire_json JSON, field_number INT, repeated_index INT, value LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_len_field_element(wire_json, field_number, repeated_index, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_bytes_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_bytes_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_len_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_bytes_field $$
CREATE FUNCTION pb_wire_json_clear_bytes_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_bytes_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_bytes_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_string_field $$
CREATE FUNCTION pb_wire_json_set_string_field(wire_json JSON, field_number INT, value LONGTEXT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_len_field(wire_json, field_number, CONVERT(value USING binary));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_string_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_string_field_element(wire_json JSON, field_number INT, value LONGTEXT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_len_field_element(wire_json, field_number, CONVERT(value USING binary));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_string_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_string_field_element(wire_json JSON, field_number INT, repeated_index INT, value LONGTEXT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_len_field_element(wire_json, field_number, repeated_index, CONVERT(value USING binary));
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_string_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_string_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_len_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_string_field $$
CREATE FUNCTION pb_wire_json_clear_string_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_string_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_string_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_message_field $$
CREATE FUNCTION pb_wire_json_set_message_field(wire_json JSON, field_number INT, value LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_len_field(wire_json, field_number, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_add_repeated_message_field_element $$
CREATE FUNCTION pb_wire_json_add_repeated_message_field_element(wire_json JSON, field_number INT, value LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_add_repeated_len_field_element(wire_json, field_number, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_set_repeated_message_field_element $$
CREATE FUNCTION pb_wire_json_set_repeated_message_field_element(wire_json JSON, field_number INT, repeated_index INT, value LONGBLOB) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_set_repeated_len_field_element(wire_json, field_number, repeated_index, value);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_remove_repeated_message_field_element $$
CREATE FUNCTION pb_wire_json_remove_repeated_message_field_element(wire_json JSON, field_number INT, repeated_index INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_remove_repeated_len_field_element(wire_json, field_number, repeated_index);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_message_field $$
CREATE FUNCTION pb_wire_json_clear_message_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP FUNCTION IF EXISTS pb_wire_json_clear_repeated_message_field $$
CREATE FUNCTION pb_wire_json_clear_repeated_message_field(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_wire_json_clear_field(wire_json, field_number);
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_int32_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_int32_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_int64(uint_value)));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_int64(uint_value)));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_int32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_uint64_as_uint32(uint_value));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_uint64_as_uint32(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_uint32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_uint64_as_uint32(uint_value));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_uint64_as_uint32(uint_value));
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_int64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_int64_field_as_json_string_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_int64_field_as_json_string_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_int64_field_as_json_string_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_int64_field_as_json_string_array $$
CREATE FUNCTION pb_wire_json_get_repeated_int64_field_as_json_string_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_int64_field_as_json_string_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_int64_field_as_json_string_array $$
CREATE PROCEDURE _pb_message_get_repeated_int64_field_as_json_string_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
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
			ITERATE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_int64_field_as_json_string_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_int64_field_as_json_string_array $$
CREATE FUNCTION pb_message_get_repeated_int64_field_as_json_string_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_int64_field_as_json_string_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_uint64_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_uint64_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_uint64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_uint64_field_as_json_string_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_uint64_field_as_json_string_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(uint_value AS CHAR));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(uint_value AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_uint64_field_as_json_string_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_uint64_field_as_json_string_array $$
CREATE FUNCTION pb_wire_json_get_repeated_uint64_field_as_json_string_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_uint64_field_as_json_string_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_uint64_field_as_json_string_array $$
CREATE PROCEDURE _pb_message_get_repeated_uint64_field_as_json_string_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
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
			ITERATE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(uint_value AS CHAR));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(uint_value AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_uint64_field_as_json_string_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_uint64_field_as_json_string_array $$
CREATE FUNCTION pb_message_get_repeated_uint64_field_as_json_string_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_uint64_field_as_json_string_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_sint32_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_sint32_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_sint64(uint_value)));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_cast_int64_as_int32(_pb_util_reinterpret_uint64_as_sint64(uint_value)));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sint32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_sint64(uint_value));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_sint64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sint64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_sint64_field_as_json_string_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_sint64_field_as_json_string_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_sint64(uint_value) AS CHAR));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_sint64(uint_value) AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sint64_field_as_json_string_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sint64_field_as_json_string_array $$
CREATE FUNCTION pb_wire_json_get_repeated_sint64_field_as_json_string_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_sint64_field_as_json_string_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_sint64_field_as_json_string_array $$
CREATE PROCEDURE _pb_message_get_repeated_sint64_field_as_json_string_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
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
			ITERATE l1;
		END IF;

		CASE current_wire_type
		WHEN 0 THEN
			CALL _pb_wire_read_varint_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_sint64(uint_value) AS CHAR));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_sint64(uint_value) AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_sint64_field_as_json_string_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sint64_field_as_json_string_array $$
CREATE FUNCTION pb_message_get_repeated_sint64_field_as_json_string_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_sint64_field_as_json_string_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_enum_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_enum_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_enum_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value <> 0);
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_varint_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value <> 0);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_bool_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 5 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_fixed32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 5 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_int32(uint_value));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_int32(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sfixed32_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 5 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_float(uint_value));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i32_as_uint32(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint32_as_float(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_float_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 1 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', uint_value);
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_fixed64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_fixed64_field_as_json_string_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_fixed64_field_as_json_string_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 1 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(uint_value AS CHAR));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(uint_value AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_fixed64_field_as_json_string_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_fixed64_field_as_json_string_array $$
CREATE FUNCTION pb_wire_json_get_repeated_fixed64_field_as_json_string_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_fixed64_field_as_json_string_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_fixed64_field_as_json_string_array $$
CREATE PROCEDURE _pb_message_get_repeated_fixed64_field_as_json_string_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
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
			ITERATE l1;
		END IF;

		CASE current_wire_type
		WHEN 1 THEN
			CALL _pb_wire_read_i64_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(uint_value AS CHAR));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(uint_value AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_fixed64_field_as_json_string_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_fixed64_field_as_json_string_array $$
CREATE FUNCTION pb_message_get_repeated_fixed64_field_as_json_string_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_fixed64_field_as_json_string_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_sfixed64_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_sfixed64_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 1 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_int64(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sfixed64_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_sfixed64_field_as_json_string_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_sfixed64_field_as_json_string_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 1 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_sfixed64_field_as_json_string_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_wire_json_get_repeated_sfixed64_field_as_json_string_array $$
CREATE FUNCTION pb_wire_json_get_repeated_sfixed64_field_as_json_string_array(wire_json JSON, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_wire_json_get_repeated_sfixed64_field_as_json_string_array(wire_json, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_message_get_repeated_sfixed64_field_as_json_string_array $$
CREATE PROCEDURE _pb_message_get_repeated_sfixed64_field_as_json_string_array(IN message LONGBLOB, IN field_number INT, OUT result JSON)
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
			ITERATE l1;
		END IF;

		CASE current_wire_type
		WHEN 1 THEN
			CALL _pb_wire_read_i64_as_uint64(tail, uint_value, tail);
			SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR));
		WHEN 2 THEN
			CALL _pb_wire_read_len_type(tail, bytes_value, tail);
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', CAST(_pb_util_reinterpret_uint64_as_int64(uint_value) AS CHAR));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_message_get_repeated_sfixed64_field_as_json_string_array: unexpected wire_type (', current_wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;
	END WHILE;
END $$

DROP FUNCTION IF EXISTS pb_message_get_repeated_sfixed64_field_as_json_string_array $$
CREATE FUNCTION pb_message_get_repeated_sfixed64_field_as_json_string_array(message LONGBLOB, field_number INT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	CALL _pb_message_get_repeated_sfixed64_field_as_json_string_array(message, field_number, result);
	RETURN result;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_repeated_double_field_as_json_array $$
CREATE PROCEDURE _pb_wire_json_get_repeated_double_field_as_json_array(IN wire_json JSON, IN field_number INT, OUT result JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 1 THEN
			SET uint_value = CAST(JSON_EXTRACT(wire_element, '$.v') AS UNSIGNED);
			SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_double(uint_value));
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			WHILE LENGTH(bytes_value) <> 0 DO
				CALL _pb_wire_read_i64_as_uint64(bytes_value, uint_value, bytes_value);
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_util_reinterpret_uint64_as_double(uint_value));
			END WHILE;
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_double_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 2 THEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			SET result = JSON_ARRAY_APPEND(result, '$', TO_BASE64(bytes_value));
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_bytes_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 2 THEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			SET result = JSON_ARRAY_APPEND(result, '$', CONVERT(bytes_value USING utf8mb4));
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_string_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
	DECLARE message_text TEXT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;
	DECLARE wire_type INT;
	DECLARE wire_elements JSON;
	DECLARE wire_element JSON;
	DECLARE wire_element_index INT;
	DECLARE wire_element_count INT;

	SET result = JSON_ARRAY();

	SET wire_elements = JSON_EXTRACT(wire_json, CONCAT('$."', field_number, '"'));
	SET wire_element_index = 0;
	SET wire_element_count = JSON_LENGTH(wire_elements);

	l1: WHILE wire_element_index < wire_element_count DO
		SET wire_element = JSON_EXTRACT(wire_elements, CONCAT('$[', wire_element_index, ']'));
		SET wire_type = JSON_EXTRACT(wire_element, '$.t');

		CASE wire_type
		WHEN 2 THEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(wire_element, '$.v')));
			SET result = JSON_ARRAY_APPEND(result, '$', TO_BASE64(bytes_value));
		ELSE
			SET message_text = CONCAT('_pb_wire_json_get_repeated_message_field_as_json_array: unexpected wire_type (', wire_type, ')');
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END CASE;

		SET wire_element_index = wire_element_index + 1;
	END WHILE;
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
			ITERATE l1;
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
