DELIMITER $$

-- Convert snake_case to camelCase
-- Implements logic from protobuf-go JSONCamelCase:
-- https://github.com/protocolbuffers/protobuf-go/blob/v1.28.1/internal/strs/strings.go
DROP FUNCTION IF EXISTS _pb_util_snake_to_camel $$
CREATE FUNCTION _pb_util_snake_to_camel(snake_name TEXT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE result TEXT DEFAULT '';
	DECLARE i INT DEFAULT 1;
	DECLARE char_val BINARY(1);
	DECLARE was_underscore BOOLEAN DEFAULT FALSE;
	DECLARE binary_name LONGBLOB;

	-- Handle empty or null input
	IF snake_name IS NULL OR LENGTH(snake_name) = 0 THEN
		RETURN snake_name;
	END IF;

	-- Cast to binary for proper character comparisons
	SET binary_name = CAST(snake_name AS BINARY);

	-- Process character by character (proto identifiers are always ASCII)
	WHILE i <= LENGTH(binary_name) DO
		SET char_val = SUBSTRING(binary_name, i, 1);

		IF char_val != _binary '_' THEN
			IF was_underscore AND (char_val >= _binary 'a' AND char_val <= _binary 'z') THEN
				-- Convert lowercase to uppercase after underscore
				SET result = CONCAT(result, UPPER(CONVERT(char_val USING utf8mb4)));
			ELSE
				-- Keep character as-is
				SET result = CONCAT(result, CONVERT(char_val USING utf8mb4));
			END IF;
		END IF;

		SET was_underscore = (char_val = _binary '_');
		SET i = i + 1;
	END WHILE;

	RETURN result;
END $$

-- Convert camelCase to snake_case
-- Implements logic from protobuf-go JSONSnakeCase:
-- https://github.com/protocolbuffers/protobuf-go/blob/v1.28.1/internal/strs/strings.go
DROP FUNCTION IF EXISTS _pb_util_camel_to_snake $$
CREATE FUNCTION _pb_util_camel_to_snake(camel_name TEXT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE result TEXT DEFAULT '';
	DECLARE i INT DEFAULT 1;
	DECLARE char_val BINARY(1);
	DECLARE binary_name LONGBLOB;

	-- Handle empty or null input
	IF camel_name IS NULL OR LENGTH(camel_name) = 0 THEN
		RETURN camel_name;
	END IF;

	-- Cast to binary for proper character comparisons
	SET binary_name = CAST(camel_name AS BINARY);

	-- Process character by character (proto identifiers are always ASCII)
	WHILE i <= LENGTH(binary_name) DO
		SET char_val = SUBSTRING(binary_name, i, 1);

		IF char_val >= _binary 'A' AND char_val <= _binary 'Z' THEN
			-- Add underscore before uppercase letter, then convert to lowercase
			SET result = CONCAT(result, '_', LOWER(CONVERT(char_val USING utf8mb4)));
		ELSE
			-- Keep character as-is
			SET result = CONCAT(result, CONVERT(char_val USING utf8mb4));
		END IF;

		SET i = i + 1;
	END WHILE;

	RETURN result;
END $$

-- Safe snake_case to camelCase conversion with round-trip validation
-- Matches protobuf-go's marshalFieldMask validation logic
DROP FUNCTION IF EXISTS _pb_util_snake_to_camel_safe $$
CREATE FUNCTION _pb_util_snake_to_camel_safe(snake_name TEXT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE camel_name TEXT;
	DECLARE roundtrip_name TEXT;
	DECLARE message_text TEXT;

	-- Convert snake_case to camelCase
	SET camel_name = _pb_util_snake_to_camel(snake_name);

	-- Round-trip validation: snake→camel→snake should equal original
	SET roundtrip_name = _pb_util_camel_to_snake(camel_name);

	IF snake_name != roundtrip_name THEN
		-- Contains irreversible value, fail like protobuf-go
		SET message_text = CONCAT('Contains irreversible value "', snake_name, '"');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	RETURN camel_name;
END $$

DELIMITER ;