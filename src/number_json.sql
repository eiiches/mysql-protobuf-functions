DELIMITER $$

-- Helper function to convert uint32 IEEE 754 bits to binary32 number JSON format
DROP FUNCTION IF EXISTS _pb_convert_float_uint32_to_number_json $$
CREATE FUNCTION _pb_convert_float_uint32_to_number_json(uint32_bits INT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	-- Always produce binary32 format for number JSON
	RETURN JSON_QUOTE(CONCAT('binary32:0x', LPAD(LOWER(HEX(uint32_bits)), 8, '0')));
END $$

-- Helper function to convert uint64 IEEE 754 bits to binary64 number JSON format
DROP FUNCTION IF EXISTS _pb_convert_double_uint64_to_number_json $$
CREATE FUNCTION _pb_convert_double_uint64_to_number_json(uint64_bits BIGINT UNSIGNED) RETURNS JSON DETERMINISTIC
BEGIN
	-- Always produce binary64 format for number JSON
	RETURN JSON_QUOTE(CONCAT('binary64:0x', LPAD(LOWER(HEX(uint64_bits)), 16, '0')));
END $$

-- Helper function to check if a value is a proto3 default value (number JSON format)
-- TODO: Deprecated, remove this function.
DROP FUNCTION IF EXISTS _pb_number_json_is_proto3_default_value $$
CREATE FUNCTION _pb_number_json_is_proto3_default_value(field_type INT, json_value JSON) RETURNS BOOLEAN DETERMINISTIC
BEGIN
	DECLARE message_text TEXT;

	CASE field_type
	WHEN 1 THEN -- TYPE_DOUBLE
		RETURN _pb_json_parse_double_as_uint64(json_value, TRUE) = 0;
	WHEN 2 THEN -- TYPE_FLOAT
		RETURN _pb_json_parse_float_as_uint32(json_value, TRUE) = 0;
	WHEN 3 THEN -- TYPE_INT64
		RETURN _pb_json_parse_signed_int(json_value) = 0;
	WHEN 4 THEN -- TYPE_UINT64
		RETURN _pb_json_parse_unsigned_int(json_value) = 0;
	WHEN 5 THEN -- TYPE_INT32
		RETURN _pb_json_parse_signed_int(json_value) = 0;
	WHEN 6 THEN -- TYPE_FIXED64
		RETURN _pb_json_parse_unsigned_int(json_value) = 0;
	WHEN 7 THEN -- TYPE_FIXED32
		RETURN _pb_json_parse_unsigned_int(json_value) = 0;
	WHEN 8 THEN -- TYPE_BOOL
		RETURN json_value = CAST(false AS JSON);
	WHEN 9 THEN -- TYPE_STRING
		RETURN JSON_UNQUOTE(json_value) = '';
	WHEN 11 THEN -- TYPE_MESSAGE
		RETURN JSON_LENGTH(json_value) = 0;
	WHEN 12 THEN -- TYPE_BYTES
		RETURN JSON_UNQUOTE(json_value) = '';
	WHEN 13 THEN -- TYPE_UINT32
		RETURN _pb_json_parse_unsigned_int(json_value) = 0;
	WHEN 15 THEN -- TYPE_SFIXED32
		RETURN _pb_json_parse_signed_int(json_value) = 0;
	WHEN 16 THEN -- TYPE_SFIXED64
		RETURN _pb_json_parse_signed_int(json_value) = 0;
	WHEN 17 THEN -- TYPE_SINT32
		RETURN _pb_json_parse_signed_int(json_value) = 0;
	WHEN 18 THEN -- TYPE_SINT64
		RETURN _pb_json_parse_signed_int(json_value) = 0;
	WHEN 14 THEN -- TYPE_ENUM
		-- Expects numeric enum value (conversion should be done elsewhere)
		RETURN CAST(json_value AS SIGNED) = 0;
	ELSE
		-- For unknown types, raise error
		SET message_text = CONCAT('_pb_number_json_is_proto3_default_value: unsupported field_type ', field_type);
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END CASE;
END $$
