DELIMITER $$

DROP FUNCTION IF EXISTS pb_message_new $$
CREATE FUNCTION pb_message_new() RETURNS LONGBLOB DETERMINISTIC
BEGIN
	-- Return an empty protobuf message
	RETURN _binary X'';
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
