DELIMITER $$

-- Helper function to convert uint32 IEEE 754 bits to JSON float value
DROP FUNCTION IF EXISTS _pb_convert_float_uint32_to_json $$
CREATE FUNCTION _pb_convert_float_uint32_to_json(uint32_bits INT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	-- Check for special values first
	IF (uint32_bits & 0x7F800000) = 0x7F800000 THEN
		-- Exponent is all 1s (255)
		IF (uint32_bits & 0x007FFFFF) = 0 THEN
			-- Mantissa is zero - this is infinity
			IF (uint32_bits & 0x80000000) = 0 THEN
				RETURN JSON_QUOTE('Infinity');
			ELSE
				RETURN JSON_QUOTE('-Infinity');
			END IF;
		ELSE
			-- Mantissa is non-zero - this is NaN
			RETURN JSON_QUOTE('NaN');
		END IF;
	END IF;

	-- For regular numbers, convert back to float and cast to JSON
	RETURN CAST(_pb_util_reinterpret_uint32_as_float(uint32_bits) AS JSON);
END $$

-- Helper function to convert uint64 IEEE 754 bits to JSON double value
DROP FUNCTION IF EXISTS _pb_convert_double_uint64_to_json $$
CREATE FUNCTION _pb_convert_double_uint64_to_json(uint64_bits BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	-- Check for special values first
	IF (uint64_bits & 0x7FF0000000000000) = 0x7FF0000000000000 THEN
		-- Exponent is all 1s (2047)
		IF (uint64_bits & 0x000FFFFFFFFFFFFF) = 0 THEN
			-- Mantissa is zero - this is infinity
			IF (uint64_bits & 0x8000000000000000) = 0 THEN
				RETURN JSON_QUOTE('Infinity');
			ELSE
				RETURN JSON_QUOTE('-Infinity');
			END IF;
		ELSE
			-- Mantissa is non-zero - this is NaN
			RETURN JSON_QUOTE('NaN');
		END IF;
	END IF;

	-- For regular numbers, convert back to double and cast to JSON
	RETURN CAST(_pb_util_reinterpret_uint64_as_double(uint64_bits) AS JSON);
END $$

-- Helper function to parse JSON value as BIGINT with validation
DROP FUNCTION IF EXISTS _pb_json_parse_signed_int $$
CREATE FUNCTION _pb_json_parse_signed_int(json_value JSON) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;
	DECLARE double_value DOUBLE;
	DECLARE uint_value BIGINT UNSIGNED;

	IF JSON_TYPE(json_value) = 'STRING' THEN
		SET str_value = JSON_UNQUOTE(json_value);

		-- Early return for invalid inputs
		IF str_value = '' OR NOT (str_value REGEXP '^-?[0-9]+(\\.[0-9]+)?([eE][+-]?[0-9]+)?$') THEN
			SET message_text = CONCAT('Invalid number format for signed integer: ', str_value);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		SET json_value = CAST(str_value AS JSON);
	END IF;

	CASE JSON_TYPE(json_value)
	WHEN 'UNSIGNED INTEGER' THEN
		-- Cast to unsigned first, then check range for signed integer
		SET uint_value = CAST(json_value AS UNSIGNED);
		IF uint_value > 9223372036854775807 THEN
			SET message_text = CONCAT('Value too large for signed integer field: ', CAST(uint_value AS CHAR));
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		RETURN CAST(uint_value AS SIGNED);
	WHEN 'INTEGER' THEN
		-- JSON INTEGER type should be safe for signed integers
		RETURN CAST(json_value AS SIGNED);
	WHEN 'DOUBLE' THEN
		-- Handle JSON numbers with exponential notation (e.g., 1e5)
		-- But reject non-integer values (e.g., 0.5)
		SET double_value = CAST(json_value AS DOUBLE);
		IF double_value != FLOOR(double_value) THEN
			SET message_text = CONCAT('Non-integer value for signed integer field: ', CAST(double_value AS CHAR));
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		RETURN CAST(double_value AS SIGNED);
	ELSE
		RETURN CAST(json_value AS SIGNED);
	END CASE;
END $$

-- Helper function to parse JSON value as BIGINT UNSIGNED with validation
DROP FUNCTION IF EXISTS _pb_json_parse_unsigned_int $$
CREATE FUNCTION _pb_json_parse_unsigned_int(json_value JSON) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;
	DECLARE double_value DOUBLE;

	IF JSON_TYPE(json_value) = 'STRING' THEN
		SET str_value = JSON_UNQUOTE(json_value);

		-- Early return for invalid inputs
		IF str_value = '' OR NOT (str_value REGEXP '^-?[0-9]+(\\.[0-9]+)?([eE][+-]?[0-9]+)?$') THEN
			SET message_text = CONCAT('Invalid number format for signed integer: ', str_value);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		SET json_value = CAST(str_value AS JSON);
	END IF;

	CASE JSON_TYPE(json_value)
	WHEN 'INTEGER' THEN
		-- Check if the INTEGER value is negative (invalid for unsigned)
		IF CAST(json_value AS SIGNED) < 0 THEN
			SET message_text = CONCAT('Negative value for unsigned integer field: ', CAST(json_value AS CHAR));
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		RETURN CAST(json_value AS UNSIGNED);
	WHEN 'UNSIGNED INTEGER' THEN
		-- UNSIGNED INTEGER values are always valid for unsigned fields
		RETURN CAST(json_value AS UNSIGNED);
	WHEN 'DOUBLE' THEN
		-- Handle JSON numbers with exponential notation (e.g., 1e5)
		-- But reject non-integer values (e.g., 0.5)
		SET double_value = CAST(json_value AS DOUBLE);
		IF double_value != FLOOR(double_value) THEN
			SET message_text = CONCAT('Non-integer value for unsigned integer field: ', CAST(double_value AS CHAR));
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		RETURN CAST(double_value AS UNSIGNED);
	ELSE
		RETURN CAST(json_value AS UNSIGNED);
	END CASE;
END $$

-- Helper function to parse JSON value as DOUBLE with validation
DROP FUNCTION IF EXISTS _pb_json_parse_double_as_uint64 $$
CREATE FUNCTION _pb_json_parse_double_as_uint64(json_value JSON, allow_hex_strings BOOLEAN) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;
	DECLARE double_value DOUBLE;
	DECLARE hex_value TEXT;

	IF JSON_TYPE(json_value) = 'STRING' THEN
		SET str_value = JSON_UNQUOTE(json_value);

		-- Reject empty strings
		IF str_value = '' THEN
			SET message_text = 'Empty string is not a valid number for double field';
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		-- Handle binary64 hex format if allowed
		IF allow_hex_strings AND str_value REGEXP '^binary64:0x[0-9a-fA-F]{16}$' THEN
			SET hex_value = SUBSTRING(str_value, 12); -- Skip "binary64:0x"
			RETURN CONV(hex_value, 16, 10);
		END IF;

		-- Handle special values
		IF str_value IN ('Infinity', '-Infinity', 'NaN') THEN
			CASE str_value
			WHEN 'Infinity' THEN
				RETURN 0x7FF0000000000000;
			WHEN '-Infinity' THEN
				RETURN 0xFFF0000000000000;
			WHEN 'NaN' THEN
				RETURN 0x7FF8000000000000;
			END CASE;
		END IF;

		-- Reject non-numeric strings (but allow binary64 format)
		IF NOT (str_value REGEXP '^[+-]?([0-9]*\\.?[0-9]+([eE][+-]?[0-9]+)?|Infinity|-Infinity|NaN|binary64:0x[0-9a-fA-F]{16})$') THEN
			SET message_text = CONCAT('Invalid number format for double field: ', str_value);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		-- Convert string to JSON for further processing
		SET json_value = CAST(str_value AS JSON);
	END IF;

	-- Convert JSON value to double
	CASE JSON_TYPE(json_value)
	WHEN 'INTEGER' THEN
		SET double_value = CAST(json_value AS DOUBLE);
	WHEN 'UNSIGNED INTEGER' THEN
		SET double_value = CAST(json_value AS DOUBLE);
	WHEN 'DECIMAL' THEN
		SET double_value = CAST(json_value AS DOUBLE);
	WHEN 'DOUBLE' THEN
		SET double_value = CAST(json_value AS DOUBLE);
	ELSE
		SET message_text = CONCAT('Invalid JSON type for double field: ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;

	-- Convert to IEEE 754 binary representation
	RETURN _pb_util_reinterpret_double_as_uint64(double_value);
END $$

-- Helper function to parse JSON value as FLOAT with validation
DROP FUNCTION IF EXISTS _pb_json_parse_float_as_uint32 $$
CREATE FUNCTION _pb_json_parse_float_as_uint32(json_value JSON, allow_hex_strings BOOLEAN) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;
	DECLARE float_value FLOAT;
	DECLARE hex_value TEXT;

	IF JSON_TYPE(json_value) = 'STRING' THEN
		SET str_value = JSON_UNQUOTE(json_value);

		-- Reject empty strings
		IF str_value = '' THEN
			SET message_text = 'Empty string is not a valid number for float field';
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		-- Handle binary32 hex format if allowed
		IF allow_hex_strings AND str_value REGEXP '^binary32:0x[0-9a-fA-F]{8}$' THEN
			SET hex_value = SUBSTRING(str_value, 12); -- Skip "binary32:0x"
			RETURN CONV(hex_value, 16, 10);
		END IF;

		-- Handle special values
		IF str_value IN ('Infinity', '-Infinity', 'NaN') THEN
			CASE str_value
			WHEN 'Infinity' THEN
				RETURN 0x7F800000;
			WHEN '-Infinity' THEN
				RETURN 0xFF800000;
			WHEN 'NaN' THEN
				RETURN 0x7FC00000;
			END CASE;
		END IF;

		-- Reject non-numeric strings (but allow binary32 format)
		IF NOT (str_value REGEXP '^[+-]?([0-9]*\\.?[0-9]+([eE][+-]?[0-9]+)?|Infinity|-Infinity|NaN|binary32:0x[0-9a-fA-F]{8})$') THEN
			SET message_text = CONCAT('Invalid number format for float field: ', str_value);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		-- Convert string to JSON for further processing
		SET json_value = CAST(str_value AS JSON);
	END IF;

	-- Convert JSON value to float
	CASE JSON_TYPE(json_value)
	WHEN 'INTEGER' THEN
		SET float_value = CAST(json_value AS FLOAT);
	WHEN 'UNSIGNED INTEGER' THEN
		SET float_value = CAST(json_value AS FLOAT);
	WHEN 'DECIMAL' THEN
		SET float_value = CAST(json_value AS FLOAT);
	WHEN 'DOUBLE' THEN
		SET float_value = CAST(json_value AS FLOAT);
	ELSE
		SET message_text = CONCAT('Invalid JSON type for float field: ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;

	-- Convert to IEEE 754 binary representation
	RETURN _pb_util_reinterpret_float_as_uint32(float_value);
END $$

-- Helper function to parse JSON value as bytes with Base64 decoding
DROP FUNCTION IF EXISTS _pb_json_parse_bytes $$
CREATE FUNCTION _pb_json_parse_bytes(json_value JSON) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE str_value TEXT;
	DECLARE message_text TEXT;

	-- Bytes must be a JSON string
	IF JSON_TYPE(json_value) != 'STRING' THEN
		SET message_text = CONCAT('Invalid JSON type for bytes field: expected STRING, got ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Unquote the JSON string
	SET str_value = JSON_UNQUOTE(json_value);

	-- Decode from Base64/Base64URL and return as binary data
	RETURN _pb_util_from_base64_url(str_value);
END $$

-- Helper function to parse JSON value as Protobuf string type
DROP FUNCTION IF EXISTS _pb_json_parse_string $$
CREATE FUNCTION _pb_json_parse_string(json_value JSON) RETURNS LONGTEXT DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;

	-- String fields must receive JSON string values
	IF JSON_TYPE(json_value) != 'STRING' THEN
		SET message_text = CONCAT('Invalid JSON type for string field: expected STRING, got ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	RETURN JSON_UNQUOTE(json_value);
END $$

-- Helper function to parse JSON value as Protobuf bool type
DROP FUNCTION IF EXISTS _pb_json_parse_bool $$
CREATE FUNCTION _pb_json_parse_bool(json_value JSON) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE bool_value BOOLEAN;
	DECLARE message_text TEXT;

	IF JSON_TYPE(json_value) <> 'BOOLEAN' THEN
		SET message_text = CONCAT('Invalid JSON type for bool field: expected BOOLEAN, got ', JSON_TYPE(json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	SET bool_value = json_value;

	RETURN bool_value;
END $$

-- Helper function to get proto3 default value for a field type
DROP FUNCTION IF EXISTS _pb_json_get_proto3_default_value $$
CREATE FUNCTION _pb_json_get_proto3_default_value(field_type INT, emit_64bit_integers_as_numbers BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;
	CASE field_type
	WHEN 1 THEN -- TYPE_DOUBLE
		RETURN CAST(0.0 AS JSON);
	WHEN 2 THEN -- TYPE_FLOAT
		RETURN CAST(0.0 AS JSON);
	WHEN 3 THEN -- TYPE_INT64
		IF emit_64bit_integers_as_numbers THEN
			RETURN CAST(0 AS JSON);
		ELSE
			RETURN JSON_QUOTE('0');
		END IF;
	WHEN 4 THEN -- TYPE_UINT64
		IF emit_64bit_integers_as_numbers THEN
			RETURN CAST(0 AS JSON);
		ELSE
			RETURN JSON_QUOTE('0');
		END IF;
	WHEN 5 THEN -- TYPE_INT32
		RETURN CAST(0 AS JSON);
	WHEN 6 THEN -- TYPE_FIXED64
		IF emit_64bit_integers_as_numbers THEN
			RETURN CAST(0 AS JSON);
		ELSE
			RETURN JSON_QUOTE('0');
		END IF;
	WHEN 7 THEN -- TYPE_FIXED32
		RETURN CAST(0 AS JSON);
	WHEN 8 THEN -- TYPE_BOOL
		RETURN CAST(false AS JSON);
	WHEN 9 THEN -- TYPE_STRING
		RETURN JSON_QUOTE('');
	WHEN 11 THEN -- TYPE_MESSAGE
		RETURN JSON_OBJECT();
	WHEN 12 THEN -- TYPE_BYTES
		RETURN JSON_QUOTE('');
	WHEN 13 THEN -- TYPE_UINT32
		RETURN CAST(0 AS JSON);
	WHEN 15 THEN -- TYPE_SFIXED32
		RETURN CAST(0 AS JSON);
	WHEN 16 THEN -- TYPE_SFIXED64
		IF emit_64bit_integers_as_numbers THEN
			RETURN CAST(0 AS JSON);
		ELSE
			RETURN JSON_QUOTE('0');
		END IF;
	WHEN 17 THEN -- TYPE_SINT32
		RETURN CAST(0 AS JSON);
	WHEN 18 THEN -- TYPE_SINT64
		IF emit_64bit_integers_as_numbers THEN
			RETURN CAST(0 AS JSON);
		ELSE
			RETURN JSON_QUOTE('0');
		END IF;
	WHEN 14 THEN -- TYPE_ENUM
		RETURN CAST(0 AS JSON);
	ELSE
		-- For unknown types, raise error
		SET message_text = CONCAT('_pb_json_get_proto3_default_value: unsupported field_type ', field_type);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$
