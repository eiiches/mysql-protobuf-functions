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
