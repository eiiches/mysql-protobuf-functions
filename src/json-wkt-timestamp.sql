DELIMITER $$

DROP FUNCTION IF EXISTS _pb_json_wkt_timestamp_format_fractional_seconds $$
CREATE FUNCTION _pb_json_wkt_timestamp_format_fractional_seconds(nanos INT) RETURNS TEXT DETERMINISTIC
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

	RETURN JSON_QUOTE(CONCAT(REPLACE(datetime_part, " ", "T"), _pb_json_wkt_timestamp_format_fractional_seconds(nanos), "Z"));
END $$

-- Helper function to convert Timestamp string to wire_json
DROP FUNCTION IF EXISTS _pb_json_encode_wkt_timestamp_as_wire_json $$
CREATE FUNCTION _pb_json_encode_wkt_timestamp_as_wire_json(timestamp_str TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	DECLARE seconds BIGINT;
	DECLARE nanos INT;
	DECLARE dot_pos INT;
	DECLARE nanos_str TEXT;
	DECLARE target_datetime DATETIME;
	DECLARE timezone_offset TEXT;

	-- Validate RFC 3339 format - must end with uppercase 'Z'
	IF timestamp_str IS NULL OR timestamp_str = '' THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid timestamp format - empty timestamp';
	END IF;

	-- Validate RFC 3339 format (supports uppercase Z suffix or timezone offsets)
	IF timestamp_str NOT REGEXP '^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\\.[0-9]{1,9})?(Z|[+-][0-9]{2}:[0-9]{2})$' COLLATE utf8mb4_bin THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid timestamp format - must follow RFC 3339 format';
	END IF;

	SET result = JSON_OBJECT();
	-- Convert timestamp string to seconds since Unix epoch, handling timezone offsets
	-- Extract timezone from the input string and convert to UTC
	IF timestamp_str LIKE '%Z' THEN
		-- UTC timezone, parse directly
		SET target_datetime = STR_TO_DATE(LEFT(timestamp_str, 19), '%Y-%m-%dT%H:%i:%s');
	ELSEIF timestamp_str REGEXP '[+-][0-9]{2}:[0-9]{2}$' THEN
		-- Handle timezone offset (+08:00, -08:00)
		SET timezone_offset = RIGHT(timestamp_str, 6);
		SET target_datetime = STR_TO_DATE(LEFT(timestamp_str, 19), '%Y-%m-%dT%H:%i:%s');
		-- Convert from local timezone to UTC using CONVERT_TZ
		SET target_datetime = CONVERT_TZ(target_datetime, timezone_offset, '+00:00');
		IF target_datetime IS NULL THEN
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid timezone offset in timestamp';
		END IF;
	ELSE
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid timestamp format - must end with Z or timezone offset';
	END IF;
	SET seconds = TIMESTAMPDIFF(SECOND, '1970-01-01 00:00:00', target_datetime);

	-- Extract nanoseconds if present
	SET nanos = 0;
	SET dot_pos = LOCATE('.', timestamp_str);
	IF dot_pos > 0 THEN
		SET nanos_str = SUBSTRING(timestamp_str, dot_pos + 1);
		-- Remove timezone suffix (Z or +/-HH:MM)
		IF nanos_str LIKE '%Z' THEN
			SET nanos_str = LEFT(nanos_str, LENGTH(nanos_str) - 1);
		ELSEIF nanos_str REGEXP '[+-][0-9]{2}:[0-9]{2}$' THEN
			SET nanos_str = LEFT(nanos_str, LENGTH(nanos_str) - 6);
		END IF;
		-- Pad or truncate to 9 digits for nanoseconds
		WHILE LENGTH(nanos_str) < 9 DO
			SET nanos_str = CONCAT(nanos_str, '0');
		END WHILE;
		SET nanos_str = LEFT(nanos_str, 9);
		SET nanos = CAST(nanos_str AS UNSIGNED);
	END IF;

	-- Validate timestamp range: [0001-01-01T00:00:00Z, 9999-12-31T23:59:59.999999999Z]
	-- This corresponds to seconds range: [-62135596800, 253402300799]
	IF seconds < -62135596800 OR seconds > 253402300799 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Timestamp out of range';
	END IF;

	-- For proto3 semantics, omit default values (seconds=0 and nanos=0)
	IF seconds = 0 AND nanos = 0 THEN
		RETURN result; -- Return empty wire_json
	END IF;

	-- Add non-default values
	IF seconds <> 0 THEN
		SET result = pb_wire_json_set_int64_field(result, 1, seconds);
	END IF;

	IF nanos <> 0 THEN
		SET result = pb_wire_json_set_int32_field(result, 2, nanos);
	END IF;

	RETURN result;
END $$
