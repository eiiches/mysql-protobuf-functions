DELIMITER $$

-- Public function interface for JSON to wire_json conversion
DROP FUNCTION IF EXISTS pb_json_to_wire_json $$
CREATE FUNCTION pb_json_to_wire_json(descriptor_set_json JSON, type_name TEXT, json_value JSON, json_unmarshal_options JSON, marshal_options JSON) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_number_json_to_wire_json(descriptor_set_json, type_name, _pb_json_to_number_json(descriptor_set_json, type_name, json_value, json_unmarshal_options), marshal_options);
END $$

-- Public function interface for JSON to message conversion
DROP FUNCTION IF EXISTS pb_json_to_message $$
CREATE FUNCTION pb_json_to_message(descriptor_set_json JSON, type_name TEXT, json_value JSON, json_unmarshal_options JSON, marshal_options JSON) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN _pb_number_json_to_message(descriptor_set_json, type_name, _pb_json_to_number_json(descriptor_set_json, type_name, json_value, json_unmarshal_options), marshal_options);
END $$
