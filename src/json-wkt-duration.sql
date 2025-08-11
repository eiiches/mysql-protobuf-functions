DELIMITER $$

DROP FUNCTION IF EXISTS _pb_json_wkt_duration_format_fractional_seconds $$
CREATE FUNCTION _pb_json_wkt_duration_format_fractional_seconds(nanos INT) RETURNS TEXT DETERMINISTIC
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
		RETURN JSON_QUOTE(CONCAT('-0', _pb_json_wkt_duration_format_fractional_seconds(nanos), 's'));
	ELSE
		RETURN JSON_QUOTE(CONCAT(CAST(seconds AS CHAR), _pb_json_wkt_duration_format_fractional_seconds(nanos), 's'));
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
