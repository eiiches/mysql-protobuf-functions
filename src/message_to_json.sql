DELIMITER $$

-- Public function interface
DROP FUNCTION IF EXISTS pb_message_to_json $$
CREATE FUNCTION pb_message_to_json(descriptor_set_json JSON, type_name TEXT, message LONGBLOB, unmarshal_options JSON, json_marshal_options JSON) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_number_json_to_json(descriptor_set_json, type_name, _pb_message_to_number_json(descriptor_set_json, type_name, message, unmarshal_options), json_marshal_options);
END $$

-- Public function interface for wire_json input
DROP FUNCTION IF EXISTS pb_wire_json_to_json $$
CREATE FUNCTION pb_wire_json_to_json(descriptor_set_json JSON, type_name TEXT, wire_json JSON, unmarshal_options JSON, json_marshal_options JSON) RETURNS JSON DETERMINISTIC
BEGIN
	RETURN _pb_number_json_to_json(descriptor_set_json, type_name, _pb_wire_json_to_number_json(descriptor_set_json, type_name, wire_json, unmarshal_options), json_marshal_options);
END $$
