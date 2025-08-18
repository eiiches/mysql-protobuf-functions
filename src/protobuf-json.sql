DELIMITER $$


DROP PROCEDURE IF EXISTS _pb_wire_json_get_primitive_field_as_json $$
CREATE PROCEDURE _pb_wire_json_get_primitive_field_as_json(IN wire_json JSON, IN field_number INT, IN field_type INT, IN is_repeated BOOLEAN, IN has_field_presence BOOLEAN, IN emit_64bit_integers_as_numbers BOOLEAN, OUT field_json_value JSON)
BEGIN
	DECLARE message_text TEXT;
	DECLARE boolean_value BOOLEAN;
	DECLARE uint_value BIGINT UNSIGNED;
	DECLARE int_value BIGINT;

	CASE field_type
	WHEN 1 THEN -- double
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_double_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_double_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 2 THEN -- float
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_float_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_float_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 3 THEN -- int64
		IF is_repeated THEN
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = pb_wire_json_get_repeated_int64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_int64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			SET int_value = pb_wire_json_get_int64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = CAST(int_value AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(int_value AS CHAR));
			END IF;
		END IF;
	WHEN 4 THEN -- uint64
		IF is_repeated THEN
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = pb_wire_json_get_repeated_uint64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_uint64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			SET uint_value = pb_wire_json_get_uint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = CAST(uint_value AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(uint_value AS CHAR));
			END IF;
		END IF;
	WHEN 5 THEN -- int32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_int32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_int32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 6 THEN -- fixed64
		IF is_repeated THEN
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = pb_wire_json_get_repeated_fixed64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_fixed64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			SET uint_value = pb_wire_json_get_fixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = CAST(uint_value AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(uint_value AS CHAR));
			END IF;
		END IF;
	WHEN 7 THEN -- fixed32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_fixed32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_fixed32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 8 THEN -- bool
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_bool_field_as_json_array(wire_json, field_number);
		ELSE
			SET boolean_value = pb_wire_json_get_bool_field(wire_json, field_number, IF(has_field_presence, NULL, FALSE));
			IF boolean_value IS NULL THEN
				SET field_json_value = NULL;
			ELSE
				-- See https://bugs.mysql.com/bug.php?id=79813
				SET field_json_value = CAST((boolean_value IS TRUE) AS JSON);
			END IF;
		END IF;
	WHEN 9 THEN -- string
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_string_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(pb_wire_json_get_string_field(wire_json, field_number, IF(has_field_presence, NULL, '')));
		END IF;
	WHEN 12 THEN -- bytes
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_bytes_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = JSON_QUOTE(TO_BASE64(pb_wire_json_get_bytes_field(wire_json, field_number, IF(has_field_presence, NULL, _binary X''))));
		END IF;
	WHEN 13 THEN -- uint32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_uint32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_uint32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 15 THEN -- sfixed32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_sfixed32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_sfixed32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 16 THEN -- sfixed64
		IF is_repeated THEN
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = pb_wire_json_get_repeated_sfixed64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_sfixed64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			SET int_value = pb_wire_json_get_sfixed64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = CAST(int_value AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(int_value AS CHAR));
			END IF;
		END IF;
	WHEN 17 THEN -- sint32
		IF is_repeated THEN
			SET field_json_value = pb_wire_json_get_repeated_sint32_field_as_json_array(wire_json, field_number);
		ELSE
			SET field_json_value = CAST(pb_wire_json_get_sint32_field(wire_json, field_number, IF(has_field_presence, NULL, 0)) AS JSON);
		END IF;
	WHEN 18 THEN -- sint64
		IF is_repeated THEN
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = pb_wire_json_get_repeated_sint64_field_as_json_array(wire_json, field_number);
			ELSE
				SET field_json_value = pb_wire_json_get_repeated_sint64_field_as_json_string_array(wire_json, field_number);
			END IF;
		ELSE
			SET int_value = pb_wire_json_get_sint64_field(wire_json, field_number, IF(has_field_presence, NULL, 0));
			IF emit_64bit_integers_as_numbers THEN
				SET field_json_value = CAST(int_value AS JSON);
			ELSE
				SET field_json_value = JSON_QUOTE(CAST(int_value AS CHAR));
			END IF;
		END IF;
	ELSE
		SET message_text = CONCAT('_pb_message_to_json: unknown field_type `', field_type, '` for field `', field_name, '` (', field_number, ').');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$
