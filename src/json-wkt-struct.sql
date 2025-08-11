DELIMITER $$

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

-- Helper procedure to convert JSON object to Struct wire_json (allows recursion)
DROP PROCEDURE IF EXISTS _pb_json_encode_wkt_struct_as_wire_json $$
CREATE PROCEDURE _pb_json_encode_wkt_struct_as_wire_json(IN json_value JSON, IN from_number_json BOOLEAN, OUT result JSON)
BEGIN
	DECLARE struct_keys JSON;
	DECLARE struct_key_count INT;
	DECLARE struct_key_index INT;
	DECLARE struct_key_name TEXT;
	DECLARE struct_value_json JSON;
	DECLARE struct_value_wire_json JSON;
	DECLARE struct_entry_wire_json JSON;

	SET @@SESSION.max_sp_recursion_depth = 255;
	SET result = JSON_OBJECT();

	IF JSON_TYPE(json_value) = 'OBJECT' THEN
		SET struct_keys = JSON_KEYS(json_value);
		SET struct_key_count = JSON_LENGTH(struct_keys);
		SET struct_key_index = 0;

		WHILE struct_key_index < struct_key_count DO
			SET struct_key_name = JSON_UNQUOTE(JSON_EXTRACT(struct_keys, CONCAT('$[', struct_key_index, ']')));
			SET struct_value_json = JSON_EXTRACT(json_value, CONCAT('$."', struct_key_name, '"'));

			-- Create map entry with key=1, value=2
			SET struct_entry_wire_json = JSON_OBJECT();
			SET struct_entry_wire_json = pb_wire_json_set_string_field(struct_entry_wire_json, 1, struct_key_name);

			-- Convert value to Value type (recursive call)
			CALL _pb_json_encode_wkt_value_as_wire_json(struct_value_json, from_number_json, struct_value_wire_json);
			IF struct_value_wire_json IS NOT NULL THEN
				SET struct_entry_wire_json = pb_wire_json_set_message_field(struct_entry_wire_json, 2, pb_wire_json_to_message(struct_value_wire_json));
				SET result = pb_wire_json_add_repeated_message_field_element(result, 1, pb_wire_json_to_message(struct_entry_wire_json));
			END IF;

			SET struct_key_index = struct_key_index + 1;
		END WHILE;
	END IF;
END $$

-- Helper procedure to convert JSON array to ListValue wire_json (allows recursion)
DROP PROCEDURE IF EXISTS _pb_json_encode_wkt_listvalue_as_wire_json $$
CREATE PROCEDURE _pb_json_encode_wkt_listvalue_as_wire_json(IN json_value JSON, IN from_number_json BOOLEAN, OUT result JSON)
BEGIN
	DECLARE list_element_count INT;
	DECLARE list_element_index INT;
	DECLARE list_element JSON;
	DECLARE list_value_wire_json JSON;

	SET @@SESSION.max_sp_recursion_depth = 255;
	SET result = JSON_OBJECT();

	IF JSON_TYPE(json_value) = 'ARRAY' THEN
		SET list_element_count = JSON_LENGTH(json_value);
		SET list_element_index = 0;

		WHILE list_element_index < list_element_count DO
			SET list_element = JSON_EXTRACT(json_value, CONCAT('$[', list_element_index, ']'));

			-- Convert element to Value type (recursive call)
			CALL _pb_json_encode_wkt_value_as_wire_json(list_element, from_number_json, list_value_wire_json);
			IF list_value_wire_json IS NOT NULL THEN
				SET result = pb_wire_json_add_repeated_message_field_element(result, 1, pb_wire_json_to_message(list_value_wire_json));
			END IF;

			SET list_element_index = list_element_index + 1;
		END WHILE;
	END IF;
END $$

-- Helper procedure to convert JSON to google.protobuf.Value wire_json (allows recursion)
DROP PROCEDURE IF EXISTS _pb_json_encode_wkt_value_as_wire_json $$
CREATE PROCEDURE _pb_json_encode_wkt_value_as_wire_json(IN json_value JSON, IN from_number_json BOOLEAN, OUT result JSON)
BEGIN
	DECLARE struct_wire_json JSON;
	DECLARE list_wire_json JSON;

	SET @@SESSION.max_sp_recursion_depth = 255;
	SET result = JSON_OBJECT();

	CASE JSON_TYPE(json_value)
	WHEN 'NULL' THEN
		-- null_value = 0 (field 1, enum)
		SET result = pb_wire_json_set_enum_field(result, 1, 0);
	WHEN 'BOOLEAN' THEN
		-- bool_value (field 4)
		SET result = pb_wire_json_set_bool_field(result, 4, IF(json_value, TRUE, FALSE));
	WHEN 'INTEGER' THEN
		-- number_value (field 2)
		SET result = pb_wire_json_set_double_field(result, 2, CAST(json_value AS DOUBLE));
	WHEN 'DECIMAL' THEN
		-- number_value (field 2)
		SET result = pb_wire_json_set_double_field(result, 2, CAST(json_value AS DOUBLE));
	WHEN 'DOUBLE' THEN
		-- number_value (field 2)
		SET result = pb_wire_json_set_double_field(result, 2, CAST(json_value AS DOUBLE));
	WHEN 'STRING' THEN
		-- string_value (field 3)
		SET result = pb_wire_json_set_string_field(result, 3, JSON_UNQUOTE(json_value));
	WHEN 'DATETIME' THEN
		-- string_value (field 3) - convert datetime to string
		SET result = pb_wire_json_set_string_field(result, 3, JSON_UNQUOTE(json_value));
	WHEN 'DATE' THEN
		-- string_value (field 3) - convert date to string
		SET result = pb_wire_json_set_string_field(result, 3, JSON_UNQUOTE(json_value));
	WHEN 'TIME' THEN
		-- string_value (field 3) - convert time to string
		SET result = pb_wire_json_set_string_field(result, 3, JSON_UNQUOTE(json_value));
	WHEN 'OBJECT' THEN
		-- struct_value (field 5) - convert to Struct
		CALL _pb_json_encode_wkt_struct_as_wire_json(json_value, from_number_json, struct_wire_json);
		IF struct_wire_json IS NOT NULL THEN
			SET result = pb_wire_json_set_message_field(result, 5, pb_wire_json_to_message(struct_wire_json));
		END IF;
	WHEN 'ARRAY' THEN
		-- list_value (field 6) - convert to ListValue
		CALL _pb_json_encode_wkt_listvalue_as_wire_json(json_value, from_number_json, list_wire_json);
		IF list_wire_json IS NOT NULL THEN
			SET result = pb_wire_json_set_message_field(result, 6, pb_wire_json_to_message(list_wire_json));
		END IF;
	WHEN 'BLOB' THEN
		-- string_value (field 3) - treat binary as string
		SET result = pb_wire_json_set_string_field(result, 3, JSON_UNQUOTE(json_value));
	WHEN 'OPAQUE' THEN
		-- string_value (field 3) - treat opaque as string
		SET result = pb_wire_json_set_string_field(result, 3, JSON_UNQUOTE(json_value));
	END CASE;
END $$
