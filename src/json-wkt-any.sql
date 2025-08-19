DELIMITER $$

DROP FUNCTION IF EXISTS _pb_wkt_any_json_to_number_json $$
CREATE FUNCTION _pb_wkt_any_json_to_number_json(proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE type_url TEXT;
	DECLARE remaining_object JSON;

	-- {"@type": "url", "field": "value"} -> {"1": "url", "2": "base64data"}
	-- This is simplified - real Any handling is more complex
	SET type_url = JSON_UNQUOTE(JSON_EXTRACT(proto_json_value, '$."@type"'));
	SET remaining_object = JSON_REMOVE(proto_json_value, '$."@type"');

	-- Convert remaining object to base64-encoded bytes (simplified)
	RETURN JSON_OBJECT('1', type_url, '2', TO_BASE64(remaining_object));
END $$