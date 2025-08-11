DELIMITER $$

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
			-- Use add_repeated_string_field_element for repeated field
			SET result = pb_wire_json_add_repeated_string_field_element(result, 1, path);
		END IF;
	END WHILE;

	RETURN result;
END $$
