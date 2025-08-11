DELIMITER $$

DROP FUNCTION IF EXISTS _pb_util_snake_to_lower_camel $$
CREATE FUNCTION _pb_util_snake_to_lower_camel(s TEXT) RETURNS TEXT DETERMINISTIC
BEGIN
	-- TODO: implement
	RETURN s;
END $$

DROP PROCEDURE IF EXISTS _pb_wire_json_get_primitive_field_as_json $$
CREATE PROCEDURE _pb_wire_json_get_primitive_field_as_json(IN wire_json JSON, IN field_number INT, IN field_type INT, IN is_repeated BOOLEAN, IN has_field_presence BOOLEAN, IN as_number_json BOOLEAN, OUT field_json_value JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE boolean_value BOOLEAN;

	CASE field_type
	WHEN 1 THEN -- double
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_double_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_double_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 2 THEN -- float
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_float_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_float_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 3 THEN -- int64
		IF is_repeated THEN
			IF as_number_json THEN
				SET field_json_value = pb_wire_json_get_repeated_int64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_int64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			IF as_number_json THEN
				SET field_json_value = CAST(pb_wire_json_get_int64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_int64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
			END IF;
		END IF;
	WHEN 4 THEN -- uint64
		IF is_repeated THEN
			IF as_number_json THEN
				SET field_json_value = pb_wire_json_get_repeated_uint64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_uint64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			IF as_number_json THEN
				SET field_json_value = CAST(pb_wire_json_get_uint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_uint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
			END IF;
		END IF;
	WHEN 5 THEN -- int32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_int32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_int32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 6 THEN -- fixed64
		IF is_repeated THEN
			IF as_number_json THEN
				SET field_json_value = pb_wire_json_get_repeated_fixed64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_fixed64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			IF as_number_json THEN
				SET field_json_value = CAST(pb_wire_json_get_fixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_fixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
			END IF;
		END IF;
	WHEN 7 THEN -- fixed32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_fixed32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_fixed32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 8 THEN -- bool
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_bool_field_as_json_array(wire_json, field_number);
		ELSE
			SET boolean_value = pb_wire_json_get_bool_field(wire_json, field_number, IF(has_field_presence, NULL, FALSE));
			IF boolean_value IS NULL THEN
				SET field_json_value = NULL;
			ELSE
				-- See https://bugs.mysql.com/bug.php?id=79813
				SET field_json_value = CAST((boolean_value IS TRUE) AS JSON);
			END IF;
		END IF;
	WHEN 9 THEN -- string
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_string_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(pb_wire_json_get_string_field(wire_json, field_number, IF(has_field_presence, NULL, '')));
		END IF;
	WHEN 12 THEN -- bytes
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_bytes_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(TO_BASE64(pb_wire_json_get_bytes_field(wire_json, field_number, IF(has_field_presence, NULL, _binary X''))));
		END IF;
	WHEN 13 THEN -- uint32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_uint32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_uint32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 15 THEN -- sfixed32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_sfixed32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_sfixed32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 16 THEN -- sfixed64
		IF is_repeated THEN
			IF as_number_json THEN
				SET field_json_value = pb_wire_json_get_repeated_sfixed64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_sfixed64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			IF as_number_json THEN
				SET field_json_value = CAST(pb_wire_json_get_sfixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_sfixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
			END IF;
		END IF;
	WHEN 17 THEN -- sint32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_sint32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_sint32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 18 THEN -- sint64
		IF is_repeated THEN
			IF as_number_json THEN
				SET field_json_value = pb_wire_json_get_repeated_sint64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_sint64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			IF as_number_json THEN
				SET field_json_value = CAST(pb_wire_json_get_sint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(pb_wire_json_get_sint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS CHAR));
			END IF;
		END IF;
	ELSE
		SET message_text = CONCAT('_pb_message_to_json: unknown field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$

DROP FUNCTION IF EXISTS _pb_wire_json_decode_wkt_timestamp_as_json $$
CREATE FUNCTION _pb_wire_json_decode_wkt_timestamp_as_json(wire_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE seconds BIGINT;
	DECLARE nanos INT;

	DECLARE elements JSON;
	DECLARE element JSON;
	DECLARE element_count INT;
	DECLARE element_index INT;
	DECLARE wire_type INT;
	DECLARE field_number INT;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE datetime_part TEXT;

	SET seconds = 0;
	SET nanos = 0;

	SET elements = JSON_EXTRACT(wire_json, '$.*[*]');
	SET element_index = 0;
	SET element_count = JSON_LENGTH(elements);
	WHILE element_index < element_count DO
		SET element = JSON_EXTRACT(elements, CONCAT('$[', element_index, ']'));
		SET wire_type = JSON_EXTRACT(element, '$.t');
		SET field_number = JSON_EXTRACT(element, '$.n');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(element, '$.v') AS UNSIGNED);
			CASE field_number
			WHEN 1 THEN
				SET seconds = _pb_util_reinterpret_uint64_as_int64(uint_value);
			WHEN 2 THEN
				SET nanos = _pb_util_reinterpret_uint64_as_int64(uint_value);
			END CASE;
		END CASE;

		SET element_index = element_index + 1;
	END WHILE;

	SET seconds = seconds + (nanos DIV 1000000000);
	SET nanos = nanos % 1000000000;

	-- Validate timestamp range: [0001-01-01T00:00:00Z, 9999-12-31T23:59:59.999999999Z]
	-- This corresponds to seconds range: [-62135596800, 253402300799]
	-- Allow for 1 second tolerance in case of nanosecond normalization
	IF seconds < -62135596800 OR seconds > 253402300800 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Timestamp out of range';
	END IF;

	-- Convert seconds since Unix epoch to datetime string using TIMESTAMPADD
	SET datetime_part = TIMESTAMPADD(SECOND, seconds, '1970-01-01 00:00:00');

	RETURN JSON_QUOTE(CONCAT(REPLACE(datetime_part, " ", "T"), _pb_util_format_fractional_seconds(nanos), "Z"));
END $$

DROP FUNCTION IF EXISTS _pb_wire_json_decode_wkt_duration_as_json $$
CREATE FUNCTION _pb_wire_json_decode_wkt_duration_as_json(wire_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE seconds BIGINT;
	DECLARE nanos INT;

	DECLARE elements JSON;
	DECLARE element JSON;
	DECLARE element_count INT;
	DECLARE element_index INT;
	DECLARE wire_type INT;
	DECLARE field_number INT;
	DECLARE uint_value BIGINT UNSIGNED;

	SET seconds = 0;
	SET nanos = 0;

	SET elements = JSON_EXTRACT(wire_json, '$.*[*]');
	SET element_index = 0;
	SET element_count = JSON_LENGTH(elements);
	WHILE element_index < element_count DO
		SET element = JSON_EXTRACT(elements, CONCAT('$[', element_index, ']'));
		SET wire_type = JSON_EXTRACT(element, '$.t');
		SET field_number = JSON_EXTRACT(element, '$.n');

		CASE wire_type
		WHEN 0 THEN
			SET uint_value = CAST(JSON_EXTRACT(element, '$.v') AS UNSIGNED);
			CASE field_number
			WHEN 1 THEN
				SET seconds = _pb_util_reinterpret_uint64_as_int64(uint_value);
			WHEN 2 THEN
				SET nanos = _pb_util_reinterpret_uint64_as_int64(uint_value);
			END CASE;
		END CASE;

		SET element_index = element_index + 1;
	END WHILE;

	SET seconds = seconds + (nanos DIV 1000000000);
	SET nanos = nanos % 1000000000;

	-- Validate duration range: [-315576000000, +315576000000] seconds
	IF seconds < -315576000000 OR seconds > 315576000000 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Duration out of range';
	END IF;

	-- Handle case where seconds=0 but nanos<0 (e.g., -0.5s)
	IF seconds = 0 AND nanos < 0 THEN
		RETURN JSON_QUOTE(CONCAT('-0', _pb_util_format_fractional_seconds(nanos), 's'));
	ELSE
		RETURN JSON_QUOTE(CONCAT(CAST(seconds AS CHAR), _pb_util_format_fractional_seconds(nanos), 's'));
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wire_json_decode_wkt_struct_as_json $$
CREATE FUNCTION _pb_wire_json_decode_wkt_struct_as_json(wire_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE elements JSON;
	DECLARE element JSON;
	DECLARE element_count INT;
	DECLARE element_index INT;
	DECLARE wire_type INT;
	DECLARE field_number INT;
	DECLARE object_key TEXT;
	DECLARE object_value JSON;
	DECLARE result JSON;

	SET result = JSON_OBJECT();

	SET elements = JSON_EXTRACT(wire_json, '$.*[*]');
	SET element_index = 0;
	SET element_count = JSON_LENGTH(elements);
	WHILE element_index < element_count DO
		SET element = JSON_EXTRACT(elements, CONCAT('$[', element_index, ']'));
		SET wire_type = JSON_EXTRACT(element, '$.t');
		SET field_number = JSON_EXTRACT(element, '$.n');

		CASE wire_type
		WHEN 2 THEN
			SET element = pb_message_to_wire_json(FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v'))));
			CASE field_number
			WHEN 1 THEN
				SET object_key = pb_wire_json_get_string_field(element, 1, '');
				SET object_value = _pb_wire_json_decode_wkt_value_as_json(pb_message_to_wire_json(pb_wire_json_get_message_field(element, 2, _binary X'')));
				SET result = JSON_MERGE(result, JSON_OBJECT(object_key, object_value));
			END CASE;
		END CASE;

		SET element_index = element_index + 1;
	END WHILE;

	RETURN result;
END $$

DROP FUNCTION IF EXISTS _pb_wire_json_decode_wkt_value_as_json $$
CREATE FUNCTION _pb_wire_json_decode_wkt_value_as_json(wire_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE elements JSON;
	DECLARE element JSON;
	DECLARE element_count INT;
	DECLARE element_index INT;
	DECLARE wire_type INT;
	DECLARE field_number INT;
	DECLARE result JSON;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE bytes_value LONGBLOB;

	SET result = JSON_OBJECT();

	SET elements = JSON_EXTRACT(wire_json, '$.*[*]');
	SET element_index = 0;
	SET element_count = JSON_LENGTH(elements);
	WHILE element_index < element_count DO
		SET element = JSON_EXTRACT(elements, CONCAT('$[', element_index, ']'));
		SET wire_type = JSON_EXTRACT(element, '$.t');
		SET field_number = JSON_EXTRACT(element, '$.n');

		CASE wire_type
		WHEN 0 THEN -- VARINT
			SET uint_value = CAST(JSON_EXTRACT(element, '$.v') AS UNSIGNED);
			CASE field_number
			WHEN 1 THEN -- null_value
				SET result = NULL;
			WHEN 4 THEN -- bool_value
				SET result = CAST(((uint_value <> 0) IS TRUE) AS JSON);
			END CASE;
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			CASE field_number
			WHEN 3 THEN -- string_value
				SET result = JSON_QUOTE(CONVERT(bytes_value USING utf8mb4));
			WHEN 5 THEN -- struct_value
				SET result = _pb_wire_json_decode_wkt_struct_as_json(pb_message_to_wire_json(bytes_value));
			WHEN 6 THEN -- list_value
				SET result = _pb_wire_json_decode_wkt_list_value_as_json(pb_message_to_wire_json(bytes_value));
			END CASE;
		WHEN 1 THEN -- I64
			SET uint_value = CAST(JSON_EXTRACT(element, '$.v') AS UNSIGNED);
			CASE field_number
			WHEN 2 THEN -- double_value
				SET result = CAST(_pb_util_reinterpret_uint64_as_double(uint_value) AS JSON);
			END CASE;
		END CASE;

		SET element_index = element_index + 1;
	END WHILE;

	RETURN result;
END $$

DROP FUNCTION IF EXISTS _pb_wire_json_decode_wkt_list_value_as_json $$
CREATE FUNCTION _pb_wire_json_decode_wkt_list_value_as_json(wire_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE elements JSON;
	DECLARE element JSON;
	DECLARE element_count INT;
	DECLARE element_index INT;
	DECLARE wire_type INT;
	DECLARE field_number INT;
	DECLARE result JSON;
	DECLARE bytes_value LONGBLOB;

	SET result = JSON_ARRAY();

	SET elements = JSON_EXTRACT(wire_json, '$.*[*]');
	SET element_index = 0;
	SET element_count = JSON_LENGTH(elements);
	WHILE element_index < element_count DO
		SET element = JSON_EXTRACT(elements, CONCAT('$[', element_index, ']'));
		SET wire_type = JSON_EXTRACT(element, '$.t');
		SET field_number = JSON_EXTRACT(element, '$.n');

		CASE wire_type
		WHEN 2 THEN -- LEN
			SET bytes_value = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v')));
			CASE field_number
			WHEN 1 THEN -- values
				SET result = JSON_ARRAY_APPEND(result, '$', _pb_wire_json_decode_wkt_value_as_json(pb_message_to_wire_json(bytes_value)));
			END CASE;
		END CASE;

		SET element_index = element_index + 1;
	END WHILE;

	RETURN result;
END $$

DROP FUNCTION IF EXISTS _pb_util_format_fractional_seconds $$
CREATE FUNCTION _pb_util_format_fractional_seconds(nanos INT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE abs_nanos INT;

	SET nanos = nanos % 1000000000;
	IF nanos = 0 THEN
		RETURN '';
	END IF;

	-- Handle negative nanoseconds
	SET abs_nanos = ABS(nanos);

	IF abs_nanos % 1000000 = 0 THEN
		RETURN CONCAT('.', LPAD(CAST(abs_nanos DIV 1000000 AS CHAR), 3, '0')); -- 3 digits
	ELSEIF abs_nanos % 1000 = 0 THEN
		RETURN CONCAT('.', LPAD(CAST(abs_nanos DIV 1000 AS CHAR), 6, '0')); -- 6 digits
	ELSE
		RETURN CONCAT('.', LPAD(CAST(abs_nanos AS CHAR), 9, '0')); -- 9 digits
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_wire_json_decode_wkt_field_mask_as_json $$
CREATE FUNCTION _pb_wire_json_decode_wkt_field_mask_as_json(wire_json JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE elements JSON;
	DECLARE element JSON;
	DECLARE element_count INT;
	DECLARE element_index INT;
	DECLARE wire_type INT;
	DECLARE field_number INT;
	DECLARE result TEXT;
	DECLARE string_value TEXT;
	DECLARE sep TEXT;

	SET result = '';
	SET sep = '';

	SET elements = JSON_EXTRACT(wire_json, '$.*[*]');
	SET element_index = 0;
	SET element_count = JSON_LENGTH(elements);
	WHILE element_index < element_count DO
		SET element = JSON_EXTRACT(elements, CONCAT('$[', element_index, ']'));
		SET wire_type = JSON_EXTRACT(element, '$.t');
		SET field_number = JSON_EXTRACT(element, '$.n');

		CASE wire_type
		WHEN 2 THEN -- LEN
			SET string_value = CONVERT(FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(element, '$.v'))) USING utf8mb4);
			CASE field_number
			WHEN 1 THEN -- values
				SET result = CONCAT(result, sep, string_value);
				SET sep = ',';
			END CASE;
		END CASE;

		SET element_index = element_index + 1;
	END WHILE;

	RETURN JSON_QUOTE(result);
END $$

DROP FUNCTION IF EXISTS _pb_wire_json_decode_wkt_as_json $$
CREATE FUNCTION _pb_wire_json_decode_wkt_as_json(wire_json JSON, full_type_name TEXT, as_number_json BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	CASE full_type_name
	WHEN '.google.protobuf.Timestamp' THEN
		RETURN _pb_wire_json_decode_wkt_timestamp_as_json(wire_json);
	WHEN '.google.protobuf.Duration' THEN
		RETURN _pb_wire_json_decode_wkt_duration_as_json(wire_json);
	WHEN '.google.protobuf.Struct' THEN
		RETURN _pb_wire_json_decode_wkt_struct_as_json(wire_json);
	WHEN '.google.protobuf.Value' THEN
		RETURN _pb_wire_json_decode_wkt_value_as_json(wire_json);
	WHEN '.google.protobuf.ListValue' THEN
		RETURN _pb_wire_json_decode_wkt_list_value_as_json(wire_json);
	WHEN '.google.protobuf.Empty' THEN
		RETURN JSON_OBJECT();
	WHEN '.google.protobuf.DoubleValue' THEN
		RETURN CAST(pb_wire_json_get_double_field(wire_json, 1, 0.0) AS JSON);
	WHEN '.google.protobuf.FloatValue' THEN
		RETURN CAST(pb_wire_json_get_float_field(wire_json, 1, 0.0) AS JSON);
	WHEN '.google.protobuf.Int64Value' THEN
		IF as_number_json THEN
			RETURN CAST(pb_wire_json_get_int64_field(wire_json, 1, 0) AS JSON);
		ELSE
			RETURN JSON_QUOTE(CAST(pb_wire_json_get_int64_field(wire_json, 1, 0) AS CHAR));
		END IF;
	WHEN '.google.protobuf.UInt64Value' THEN
		IF as_number_json THEN
			RETURN CAST(pb_wire_json_get_uint64_field(wire_json, 1, 0) AS JSON);
		ELSE
			RETURN JSON_QUOTE(CAST(pb_wire_json_get_uint64_field(wire_json, 1, 0) AS CHAR));
		END IF;
	WHEN '.google.protobuf.Int32Value' THEN
		RETURN CAST(pb_wire_json_get_int32_field(wire_json, 1, 0) AS JSON);
	WHEN '.google.protobuf.UInt32Value' THEN
		RETURN CAST(pb_wire_json_get_uint32_field(wire_json, 1, 0) AS JSON);
	WHEN '.google.protobuf.BoolValue' THEN
		RETURN CAST((pb_wire_json_get_bool_field(wire_json, 1, FALSE) IS TRUE) AS JSON);
	WHEN '.google.protobuf.StringValue' THEN
		RETURN JSON_QUOTE(pb_wire_json_get_string_field(wire_json, 1, ''));
	WHEN '.google.protobuf.BytesValue' THEN
		RETURN JSON_QUOTE(TO_BASE64(pb_wire_json_get_bytes_field(wire_json, 1, _binary X'')));
	WHEN '.google.protobuf.FieldMask' THEN
		RETURN _pb_wire_json_decode_wkt_field_mask_as_json(wire_json);
	ELSE
		RETURN NULL;
	END CASE;
END $$
