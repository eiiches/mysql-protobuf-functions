DELIMITER $$

-- Encode bytes to Base64 string without newlines
-- MySQL's TO_BASE64 adds newlines every 76 characters, this function removes them
DROP FUNCTION IF EXISTS _pb_to_base64 $$
CREATE FUNCTION _pb_to_base64(data LONGBLOB) RETURNS LONGTEXT DETERMINISTIC
BEGIN
	-- Handle null/empty input
	IF data IS NULL THEN
		RETURN NULL;
	END IF;
	
	-- Use TO_BASE64 and remove all newline characters
	RETURN REPLACE(TO_BASE64(data), '\n', '');
END $$

-- Decode Base64 encoded string to bytes
-- Supports both standard Base64 (+/) and Base64 URL (-_) encoding
-- Handles input with or without padding as per protobuf JSON spec
DROP FUNCTION IF EXISTS _pb_util_from_base64_url $$
CREATE FUNCTION _pb_util_from_base64_url(encoded_value LONGTEXT) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE standard_base64 LONGTEXT;

	-- Handle null/empty input
	IF encoded_value IS NULL OR LENGTH(encoded_value) = 0 THEN
		RETURN encoded_value;
	END IF;

	-- Convert Base64 URL encoding to standard Base64 if needed:
	-- Replace - with + and _ with / (no-op for standard Base64)
	SET standard_base64 = REPLACE(REPLACE(encoded_value, '-', '+'), '_', '/');

	-- Add padding if needed (handles both Base64 URL without padding and standard Base64 without padding)
	-- Base64 strings must be multiples of 4 characters
	WHILE LENGTH(standard_base64) % 4 != 0 DO
		SET standard_base64 = CONCAT(standard_base64, '=');
	END WHILE;

	-- Decode using MySQL's standard FROM_BASE64 function
	RETURN FROM_BASE64(standard_base64);
END $$

DELIMITER ;