DELIMITER $$

DROP FUNCTION IF EXISTS _pb_json_wkt_time_common_format_fractional_seconds $$
CREATE FUNCTION _pb_json_wkt_time_common_format_fractional_seconds(nanos INT) RETURNS TEXT DETERMINISTIC
BEGIN
	-- Validate that nanos is within [0, 999999999] range
	-- Caller must ensure proper normalization before calling this function
	IF nanos < 0 OR nanos > 999999999 THEN
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '_pb_json_wkt_time_common_format_fractional_seconds: nanos must be in range [0, 999999999]';
	END IF;

	IF nanos = 0 THEN
		RETURN '';
	END IF;

	IF nanos % 1000000 = 0 THEN
		RETURN CONCAT('.', LPAD(CAST(nanos DIV 1000000 AS CHAR), 3, '0')); -- 3 digits
	ELSEIF nanos % 1000 = 0 THEN
		RETURN CONCAT('.', LPAD(CAST(nanos DIV 1000 AS CHAR), 6, '0')); -- 6 digits
	ELSE
		RETURN CONCAT('.', LPAD(CAST(nanos AS CHAR), 9, '0')); -- 9 digits
	END IF;
END $$
