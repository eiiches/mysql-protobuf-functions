DELIMITER $$

DROP FUNCTION IF EXISTS _pb_wkt_empty_json_to_number_json $$
CREATE FUNCTION _pb_wkt_empty_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	-- Empty always returns an empty JSON object regardless of input
	RETURN JSON_OBJECT();
END $$

-- Helper function to convert Empty from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_empty_number_json_to_json $$
CREATE FUNCTION _pb_wkt_empty_number_json_to_json(number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	-- Empty object stays empty
	RETURN JSON_OBJECT();
END $$