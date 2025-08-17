DELIMITER $$

-- Helper procedure to normalize timestamp seconds and nanoseconds
-- Ensures nanos is non-negative and within [0, 999999999] range
-- Even for negative seconds, nanos must be non-negative and count forward in time
DROP PROCEDURE IF EXISTS _pb_wkt_timestamp_normalize_fields $$
CREATE PROCEDURE _pb_wkt_timestamp_normalize_fields(INOUT seconds BIGINT, INOUT nanos INT)
BEGIN
	-- Handle case where nanos is outside [0, 999999999] range
	-- For negative seconds with fractional part, nanos should still be positive
	-- Example: -1.5 seconds = seconds=-2, nanos=500000000 (not seconds=-1, nanos=-500000000)

	DECLARE extra_seconds BIGINT;

	-- Handle nanos overflow/underflow
	-- Calculate how many whole seconds are represented by nanos
	SET extra_seconds = nanos DIV 1000000000;
	SET nanos = nanos % 1000000000;

	-- Handle negative modulo result (MySQL modulo can return negative values)
	IF nanos < 0 THEN
		SET extra_seconds = extra_seconds - 1;
		SET nanos = nanos + 1000000000;
	END IF;

	-- Add the extra seconds to the original seconds
	SET seconds = seconds + extra_seconds;
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

	-- Normalize seconds and nanos using helper procedure
	CALL _pb_wkt_timestamp_normalize_fields(seconds, nanos);

	-- Validate timestamp range: [0001-01-01T00:00:00Z, 9999-12-31T23:59:59.999999999Z]
	-- This corresponds to seconds range: [-62135596800, 253402300799]
	-- Allow for 1 second tolerance in case of nanosecond normalization
	IF seconds < -62135596800 OR seconds > 253402300800 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Timestamp out of range';
	END IF;

	-- Convert seconds since Unix epoch to datetime string using TIMESTAMPADD
	SET datetime_part = TIMESTAMPADD(SECOND, seconds, '1970-01-01 00:00:00');

	RETURN JSON_QUOTE(CONCAT(REPLACE(datetime_part, " ", "T"), _pb_json_wkt_time_common_format_fractional_seconds(nanos), "Z"));
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

	-- Normalize seconds and nanos using helper procedure
	CALL _pb_wkt_timestamp_normalize_fields(seconds, nanos);

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

-- Helper procedure to convert Timestamp from ProtoJSON to ProtoNumberJSON
DROP PROCEDURE IF EXISTS _pb_wkt_timestamp_json_to_number_json $$
CREATE PROCEDURE _pb_wkt_timestamp_json_to_number_json(
	IN proto_json_value JSON,
	OUT number_json_value JSON
)
BEGIN
	DECLARE timestamp_str TEXT;
	DECLARE seconds_part BIGINT;
	DECLARE nanos_part INT;
	DECLARE dot_pos INT;
	DECLARE nanos_str TEXT;

	-- Convert ISO 8601 timestamp to {seconds, nanos}
	SET timestamp_str = JSON_UNQUOTE(proto_json_value);
	-- Parse RFC3339 timestamp format: "1972-01-01T10:00:20.021Z"
	-- Since RFC3339 timestamps are UTC, parse and convert to Unix epoch
	-- Use CONVERT_TZ to ensure we're treating input as UTC
	SET seconds_part = UNIX_TIMESTAMP(CONVERT_TZ(STR_TO_DATE(LEFT(timestamp_str, 19), '%Y-%m-%dT%H:%i:%s'), '+00:00', @@session.time_zone));

	-- Extract nanoseconds part if present
	SET dot_pos = LOCATE('.', timestamp_str);
	IF dot_pos > 0 THEN
		SET nanos_str = SUBSTRING(timestamp_str, dot_pos + 1);
		-- Remove trailing 'Z'
		SET nanos_str = LEFT(nanos_str, LENGTH(nanos_str) - 1);
		-- Pad to 9 digits (nanoseconds)
		SET nanos_str = RPAD(nanos_str, 9, '0');
		SET nanos_part = CAST(nanos_str AS UNSIGNED);
	ELSE
		SET nanos_part = 0;
	END IF;

	-- Build result with proto3 zero-value omission
	SET number_json_value = JSON_OBJECT();
	IF seconds_part != 0 THEN
		SET number_json_value = JSON_SET(number_json_value, '$."1"', seconds_part);
	END IF;
	IF nanos_part != 0 THEN
		SET number_json_value = JSON_SET(number_json_value, '$."2"', nanos_part);
	END IF;
END $$

-- Helper procedure to convert Timestamp from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_wkt_timestamp_number_json_to_json $$
CREATE PROCEDURE _pb_wkt_timestamp_number_json_to_json(
	IN number_json_value JSON,
	OUT proto_json_value JSON
)
BEGIN
	DECLARE seconds_part BIGINT;
	DECLARE nanos_part INT;
	DECLARE timestamp_str TEXT;

	-- Convert {seconds, nanos} to ISO 8601 timestamp
	-- Handle empty object (zero timestamp)
	IF JSON_LENGTH(number_json_value) = 0 THEN
		SET proto_json_value = JSON_QUOTE('1970-01-01T00:00:00Z');
	ELSE
		SET seconds_part = COALESCE(JSON_EXTRACT(number_json_value, '$."1"'), 0);
		SET nanos_part = COALESCE(JSON_EXTRACT(number_json_value, '$."2"'), 0);

		-- Normalize seconds and nanos using helper procedure
		CALL _pb_wkt_timestamp_normalize_fields(seconds_part, nanos_part);

		-- Convert seconds since Unix epoch to datetime string using TIMESTAMPADD
		SET timestamp_str = TIMESTAMPADD(SECOND, seconds_part, '1970-01-01 00:00:00');
		-- Format as ISO8601 with T separator
		SET timestamp_str = REPLACE(timestamp_str, ' ', 'T');
		-- Add fractional seconds using common function
		SET timestamp_str = CONCAT(timestamp_str, _pb_json_wkt_time_common_format_fractional_seconds(nanos_part), 'Z');
		SET proto_json_value = JSON_QUOTE(timestamp_str);
	END IF;
END $$
