DELIMITER $$

DROP FUNCTION IF EXISTS _pb_wkt_empty_json_to_number_json $$
CREATE FUNCTION _pb_wkt_empty_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	-- Empty always returns an empty JSON object regardless of input
	RETURN JSON_OBJECT();
END $$