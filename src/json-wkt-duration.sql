DELIMITER $$

-- Helper procedure to normalize duration seconds and nanoseconds
-- Follows protobuf Duration specification:
-- - nanos range: -999,999,999 to +999,999,999
-- - For durations >= 1 second: nanos must have same sign as seconds
-- - For durations < 1 second: seconds = 0, nanos can be positive or negative
DROP PROCEDURE IF EXISTS _pb_wkt_duration_normalize_fields $$
CREATE PROCEDURE _pb_wkt_duration_normalize_fields(INOUT seconds BIGINT, INOUT nanos INT)
BEGIN
	DECLARE extra_seconds BIGINT;
	DECLARE abs_nanos INT;

	-- Handle nanos overflow/underflow (outside [-999999999, 999999999])
	IF ABS(nanos) > 999999999 THEN
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
	END IF;

	-- Apply Duration-specific sign rules:
	-- For durations >= 1 second: nanos must have same sign as seconds
	-- For durations < 1 second: seconds = 0
	IF seconds > 0 AND nanos < 0 THEN
		-- Positive duration with negative nanos: adjust to maintain same sign
		SET seconds = seconds - 1;
		SET nanos = 1000000000 + nanos;
	ELSEIF seconds < 0 AND nanos > 0 THEN
		-- Negative duration with positive nanos: adjust to maintain same sign
		SET seconds = seconds + 1;
		SET nanos = nanos - 1000000000;
	END IF;

	-- Validate final ranges
	-- seconds: -315,576,000,000 to +315,576,000,000
	-- nanos: -999,999,999 to +999,999,999
	IF seconds < -315576000000 OR seconds > 315576000000 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Duration seconds out of range';
	END IF;

	IF ABS(nanos) > 999999999 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Duration nanos out of range';
	END IF;
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
		RETURN JSON_QUOTE(CONCAT('-0', _pb_json_wkt_time_common_format_fractional_seconds(ABS(nanos)), 's'));
	ELSE
		RETURN JSON_QUOTE(CONCAT(CAST(seconds AS CHAR), _pb_json_wkt_time_common_format_fractional_seconds(ABS(nanos)), 's'));
	END IF;
END $$

-- Helper function to convert Duration string to wire_json
DROP FUNCTION IF EXISTS _pb_json_encode_wkt_duration_as_wire_json $$
CREATE FUNCTION _pb_json_encode_wkt_duration_as_wire_json(duration_str TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	DECLARE seconds BIGINT;
	DECLARE nanos INT;
	DECLARE dot_pos INT;
	DECLARE s_pos INT;
	DECLARE nanos_str TEXT;
	DECLARE is_negative BOOLEAN DEFAULT FALSE;
	DECLARE duration_without_s TEXT;

	SET result = JSON_OBJECT();

	-- Validate duration format - must end with 's'
	IF duration_str IS NULL OR duration_str = '' THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid duration format - empty duration';
	END IF;

	-- Find 's' suffix
	SET s_pos = LOCATE('s', duration_str);
	IF s_pos = 0 OR s_pos <> LENGTH(duration_str) THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid duration format - must end with s';
	END IF;

	-- Remove 's' suffix and check for negative sign
	SET duration_without_s = LEFT(duration_str, s_pos - 1);
	SET is_negative = LEFT(TRIM(duration_without_s), 1) = '-';

	SET dot_pos = LOCATE('.', duration_without_s);

	IF dot_pos > 0 THEN
		-- Has fractional seconds
		SET seconds = CAST(LEFT(duration_without_s, dot_pos - 1) AS SIGNED);
		SET nanos_str = SUBSTRING(duration_without_s, dot_pos + 1);
		-- Pad to 9 digits for nanoseconds
		WHILE LENGTH(nanos_str) < 9 DO
			SET nanos_str = CONCAT(nanos_str, '0');
		END WHILE;
		SET nanos_str = LEFT(nanos_str, 9);
		SET nanos = CAST(nanos_str AS SIGNED);

		-- Handle negative durations: if seconds is negative or zero but original had minus, nanos should be negative
		IF seconds < 0 OR (seconds = 0 AND is_negative) THEN
			SET nanos = -nanos;
		END IF;
	ELSE
		-- Whole seconds only
		SET seconds = CAST(duration_without_s AS SIGNED);
		SET nanos = 0;
	END IF;

	-- Validate duration range: [-315576000000, +315576000000] seconds
	IF seconds < -315576000000 OR seconds > 315576000000 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Duration out of range';
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

-- Helper procedure to convert Duration from ProtoJSON to ProtoNumberJSON
DROP PROCEDURE IF EXISTS _pb_wkt_duration_json_to_number_json $$
CREATE PROCEDURE _pb_wkt_duration_json_to_number_json(
	IN proto_json_value JSON,
	OUT number_json_value JSON
)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';
	DECLARE message_text TEXT;
	DECLARE duration_str TEXT;
	DECLARE seconds_part BIGINT;
	DECLARE nanos_part INT;
	DECLARE dot_pos INT;
	DECLARE seconds_str TEXT;
	DECLARE nanos_str TEXT;
	DECLARE suffix_char CHAR(1);

	-- Convert duration string like "3.5s" to {seconds, nanos}
	SET duration_str = JSON_UNQUOTE(proto_json_value);
	SET suffix_char = RIGHT(duration_str, 1);

	IF suffix_char != 's' THEN
		SET message_text = CONCAT('_pb_wkt_duration_json_to_number_json: invalid Duration format: ', duration_str);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Remove 's' suffix
	SET duration_str = LEFT(duration_str, LENGTH(duration_str) - 1);

	-- Split on decimal point
	SET dot_pos = LOCATE('.', duration_str);
	IF dot_pos > 0 THEN
		SET seconds_str = LEFT(duration_str, dot_pos - 1);
		SET nanos_str = SUBSTRING(duration_str, dot_pos + 1);
		-- Pad to 9 digits
		SET nanos_str = RPAD(nanos_str, 9, '0');
		SET seconds_part = CAST(seconds_str AS SIGNED);
		SET nanos_part = CAST(nanos_str AS UNSIGNED);
	ELSE
		SET seconds_part = CAST(duration_str AS SIGNED);
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

-- Helper procedure to convert Duration from ProtoNumberJSON to ProtoJSON
DROP PROCEDURE IF EXISTS _pb_wkt_duration_number_json_to_json $$
CREATE PROCEDURE _pb_wkt_duration_number_json_to_json(
	IN number_json_value JSON,
	OUT proto_json_value JSON
)
BEGIN
	DECLARE seconds_part BIGINT;
	DECLARE nanos_part INT;
	DECLARE duration_str TEXT;

	-- Convert {seconds, nanos} to duration string like "3.5s"
	-- Handle empty object (zero duration)
	IF JSON_LENGTH(number_json_value) = 0 THEN
		SET proto_json_value = JSON_QUOTE('0s');
	ELSE
		SET seconds_part = COALESCE(JSON_EXTRACT(number_json_value, '$."1"'), 0);
		SET nanos_part = COALESCE(JSON_EXTRACT(number_json_value, '$."2"'), 0);

		-- Normalize seconds and nanos using duration-specific helper procedure
		CALL _pb_wkt_duration_normalize_fields(seconds_part, nanos_part);

		IF nanos_part = 0 THEN
			SET duration_str = CONCAT(seconds_part, 's');
		ELSE
			-- Use common function to format fractional seconds properly
			-- Duration nanos can be negative, so use absolute value for formatting
			SET duration_str = CONCAT(seconds_part, _pb_json_wkt_time_common_format_fractional_seconds(ABS(nanos_part)), 's');
		END IF;

		SET proto_json_value = JSON_QUOTE(duration_str);
	END IF;
END $$
