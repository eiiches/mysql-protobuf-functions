DELIMITER $$

-- Convert snake_case field path to camelCase for FieldMask JSON output
-- Implements logic from:
-- https://github.com/protocolbuffers/protobuf/blob/89c585602af7d28646ce92cd3abba07cfdad7fa6/src/google/protobuf/json/internal/unparser.cc#L696-L733
DROP FUNCTION IF EXISTS _pb_field_mask_snake_to_camel $$
CREATE FUNCTION _pb_field_mask_snake_to_camel(snake_path TEXT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE result TEXT DEFAULT '';
	DECLARE i INT DEFAULT 1;
	DECLARE char_val BINARY(1);
	DECLARE saw_underscore BOOLEAN DEFAULT FALSE;
	DECLARE binary_path LONGBLOB;

	-- Handle empty or null input
	IF snake_path IS NULL OR LENGTH(snake_path) = 0 THEN
		RETURN snake_path;
	END IF;

	-- Cast to binary for proper character comparisons
	SET binary_path = CAST(snake_path AS BINARY);

	-- Process character by character
	char_loop: WHILE i <= LENGTH(binary_path) DO
		SET char_val = SUBSTRING(binary_path, i, 1);

		IF (char_val >= _binary 'a' AND char_val <= _binary 'z') AND saw_underscore THEN
			-- Convert lowercase to uppercase after underscore
			SET result = CONCAT(result, UPPER(CONVERT(char_val USING utf8mb4)));
		ELSEIF (char_val >= _binary '0' AND char_val <= _binary '9') OR (char_val >= _binary 'a' AND char_val <= _binary 'z') OR char_val = _binary '.' THEN
			-- Keep lowercase letters, digits, and dots as-is
			SET result = CONCAT(result, CONVERT(char_val USING utf8mb4));
		ELSEIF char_val = _binary '_' AND NOT saw_underscore THEN
			-- Allow single underscore (will set saw_underscore = true and continue)
			SET saw_underscore = TRUE;
			SET i = i + 1;
			ITERATE char_loop;
		ELSE
			-- Invalid character in field path (allow_legacy_nonconformant_behavior is disabled)
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid character in FieldMask path';
		END IF;

		SET saw_underscore = FALSE;

		SET i = i + 1;
	END WHILE;

	RETURN result;
END $$

-- Convert camelCase field path to snake_case for FieldMask proto format
-- Implements logic from:
-- https://github.com/protocolbuffers/protobuf/blob/89c585602af7d28646ce92cd3abba07cfdad7fa6/src/google/protobuf/json/internal/parser.cc#L1001-L1036
DROP FUNCTION IF EXISTS _pb_field_mask_camel_to_snake $$
CREATE FUNCTION _pb_field_mask_camel_to_snake(camel_path TEXT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE result TEXT DEFAULT '';
	DECLARE i INT DEFAULT 1;
	DECLARE char_val BINARY(1);
	DECLARE binary_path LONGBLOB;

	-- Handle empty or null input
	IF camel_path IS NULL OR LENGTH(camel_path) = 0 THEN
		RETURN camel_path;
	END IF;

	-- Cast to binary for proper character comparisons
	SET binary_path = CAST(camel_path AS BINARY);

	-- Process character by character
	WHILE i <= LENGTH(binary_path) DO
		SET char_val = SUBSTRING(binary_path, i, 1);

		IF (char_val >= _binary '0' AND char_val <= _binary '9') OR (char_val >= _binary 'a' AND char_val <= _binary 'z') OR char_val = _binary '.' THEN
			-- Keep lowercase letters, digits, and dots as-is
			SET result = CONCAT(result, CONVERT(char_val USING utf8mb4));
		ELSEIF char_val >= _binary 'A' AND char_val <= _binary 'Z' THEN
			-- Convert uppercase to lowercase with preceding underscore
			SET result = CONCAT(result, '_', LOWER(CONVERT(char_val USING utf8mb4)));
		ELSE
			-- Invalid character in field path (allow_legacy_nonconformant_behavior is disabled)
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Invalid character in FieldMask path';
		END IF;

		SET i = i + 1;
	END WHILE;

	RETURN result;
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
				-- Convert snake_case proto field path to camelCase JSON path
				SET result = CONCAT(result, sep, _pb_field_mask_snake_to_camel(string_value));
				SET sep = ',';
			END CASE;
		END CASE;

		SET element_index = element_index + 1;
	END WHILE;

	RETURN JSON_QUOTE(result);
END $$

-- Helper function to convert FieldMask string to wire_json
DROP FUNCTION IF EXISTS _pb_json_encode_wkt_field_mask_as_wire_json $$
CREATE FUNCTION _pb_json_encode_wkt_field_mask_as_wire_json(field_mask_str TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	DECLARE comma_pos INT;
	DECLARE path TEXT;
	DECLARE remaining TEXT;

	SET result = JSON_OBJECT();
	SET remaining = field_mask_str;

	WHILE remaining IS NOT NULL AND LENGTH(remaining) > 0 DO
		SET comma_pos = LOCATE(',', remaining);
		IF comma_pos > 0 THEN
			SET path = TRIM(LEFT(remaining, comma_pos - 1));
			SET remaining = SUBSTRING(remaining, comma_pos + 1);
		ELSE
			SET path = TRIM(remaining);
			SET remaining = NULL;
		END IF;

		IF LENGTH(path) > 0 THEN
			-- Convert camelCase JSON field path to snake_case proto path
			-- Use add_repeated_string_field_element for repeated field
			SET result = pb_wire_json_add_repeated_string_field_element(result, 1, _pb_field_mask_camel_to_snake(path));
		END IF;
	END WHILE;

	RETURN result;
END $$

-- Convert FieldMask JSON to number JSON format
DROP FUNCTION IF EXISTS _pb_wkt_field_mask_json_to_number_json $$
CREATE FUNCTION _pb_wkt_field_mask_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE field_mask_str TEXT;
	DECLARE paths_array JSON;
	DECLARE comma_pos INT;
	DECLARE current_path TEXT;
	DECLARE remaining_str TEXT;

	-- Convert comma-separated string to paths array
	-- "path1,path2" -> {"1": ["path1", "path2"]}

	SET field_mask_str = JSON_UNQUOTE(proto_json_value);
	IF field_mask_str = '' THEN
		RETURN JSON_OBJECT();
	END IF;

	-- Split comma-separated string into array
	SET paths_array = JSON_ARRAY();
	-- Simple implementation: split by comma and add each path
	-- Note: This is a simplified implementation
	SET remaining_str = field_mask_str;

	split_loop: WHILE LENGTH(remaining_str) > 0 DO
		SET comma_pos = LOCATE(',', remaining_str);
		IF comma_pos > 0 THEN
			SET current_path = TRIM(LEFT(remaining_str, comma_pos - 1));
			SET remaining_str = TRIM(SUBSTRING(remaining_str, comma_pos + 1));
		ELSE
			SET current_path = TRIM(remaining_str);
			SET remaining_str = '';
		END IF;

		IF LENGTH(current_path) > 0 THEN
			-- Convert camelCase JSON field path to snake_case proto path
			SET paths_array = JSON_ARRAY_APPEND(paths_array, '$', _pb_field_mask_camel_to_snake(current_path));
		END IF;
	END WHILE split_loop;

	RETURN JSON_OBJECT('1', paths_array);
END $$

-- Convert FieldMask from number JSON format to JSON format
-- Extracts from number-json-to-json.sql and applies proper snake_to_camel conversion
DROP FUNCTION IF EXISTS _pb_wkt_field_mask_number_json_to_json $$
CREATE FUNCTION _pb_wkt_field_mask_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE paths_array JSON;
	DECLARE path_count INT;
	DECLARE path_index INT;
	DECLARE current_path TEXT;
	DECLARE result_str TEXT;

	-- Convert {"1": ["path1", "path2"]} to "camelPath1,camelPath2"
	SET paths_array = JSON_EXTRACT(number_json_value, '$.\"1\"');
	IF paths_array IS NULL OR JSON_LENGTH(paths_array) = 0 THEN
		RETURN JSON_QUOTE('');
	ELSE
		SET path_count = JSON_LENGTH(paths_array);
		SET path_index = 0;
		SET result_str = '';

		path_loop: WHILE path_index < path_count DO
			SET current_path = JSON_UNQUOTE(JSON_EXTRACT(paths_array, CONCAT('$[', path_index, ']')));
			IF path_index > 0 THEN
				SET result_str = CONCAT(result_str, ',');
			END IF;
			-- Convert snake_case proto field path to camelCase JSON path
			SET result_str = CONCAT(result_str, _pb_field_mask_snake_to_camel(current_path));
			SET path_index = path_index + 1;
		END WHILE path_loop;

		RETURN JSON_QUOTE(result_str);
	END IF;
END $$
