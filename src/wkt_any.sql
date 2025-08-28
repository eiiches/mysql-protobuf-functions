DELIMITER $$

-- Helper function to find descriptor set for a given type name
DROP FUNCTION IF EXISTS _pb_wkt_any_find_descriptor_set $$
CREATE FUNCTION _pb_wkt_any_find_descriptor_set(type_name TEXT, descriptor_set_jsons JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE i INT DEFAULT 0;
	DECLARE descriptor_count INT;
	DECLARE current_descriptor_set JSON;
	DECLARE message_descriptor JSON;

	IF descriptor_set_jsons IS NULL OR type_name IS NULL THEN
		RETURN NULL;
	END IF;

	SET descriptor_count = JSON_LENGTH(descriptor_set_jsons);

	WHILE i < descriptor_count DO
		SET current_descriptor_set = JSON_EXTRACT(descriptor_set_jsons, CONCAT('$[', i, ']'));
		SET message_descriptor = _pb_descriptor_set_get_message_descriptor(current_descriptor_set, type_name);

		IF message_descriptor IS NOT NULL THEN
			RETURN current_descriptor_set;
		END IF;

		SET i = i + 1;
	END WHILE;

	RETURN NULL;
END $$

-- Helper function to extract type name from type_url
DROP FUNCTION IF EXISTS _pb_wkt_any_extract_type_name $$
CREATE FUNCTION _pb_wkt_any_extract_type_name(type_url TEXT) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE last_slash_pos INT;

	IF type_url IS NULL THEN
		RETURN NULL;
	END IF;

	-- Extract type name from URLs like "type.googleapis.com/google.protobuf.StringValue"
	SET last_slash_pos = CHAR_LENGTH(type_url) - CHAR_LENGTH(SUBSTRING_INDEX(type_url, '/', -1));
	IF last_slash_pos > 0 THEN
		RETURN CONCAT('.', SUBSTRING(type_url, last_slash_pos + 1));
	END IF;

	-- If no slash, assume the entire string is the type name
	RETURN CONCAT('.', type_url);
END $$

-- Convert google.protobuf.Any from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_wkt_any_number_json_to_json $$
CREATE FUNCTION _pb_wkt_any_number_json_to_json(number_json_value JSON, descriptor_set_jsons JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE type_url TEXT;
	DECLARE value_base64 TEXT;
	DECLARE type_name TEXT;
	DECLARE message_descriptor JSON;
	DECLARE decoded_message LONGBLOB;
	DECLARE inner_json JSON;
	DECLARE inner_number_json JSON;
	DECLARE result JSON;
	DECLARE i INT DEFAULT 0;
	DECLARE descriptor_count INT;
	DECLARE current_descriptor_set JSON;
	DECLARE descriptor_set_json JSON;
	DECLARE wire_json JSON;
	DECLARE message LONGBLOB;
	DECLARE message_text TEXT;

	-- Extract type_url (field 1) and value (field 2) from ProtoNumberJSON
	SET type_url = JSON_UNQUOTE(JSON_EXTRACT(number_json_value, '$.\"1\"'));
	SET message = FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(number_json_value, '$.\"2\"')));

	-- If no value, return just the type
	IF message IS NULL THEN
		RETURN JSON_OBJECT('@type', type_url);
	END IF;

	-- Extract type name from type_url
	SET type_name = _pb_wkt_any_extract_type_name(type_url);

	SET descriptor_set_json = _pb_wkt_get_descriptor_set(type_name);
	IF descriptor_set_json IS NOT NULL THEN -- wkt
		-- TODO: set json_marshal_option
		SET inner_number_json = _pb_message_to_number_json(descriptor_set_json, type_name, message, NULL);
		CALL _pb_convert_number_json_to_wkt(11, type_name, inner_number_json, descriptor_set_jsons, inner_json);
		RETURN JSON_OBJECT('@type', type_url, 'value', inner_json);
	ELSE -- non-wkt
		-- Find the descriptor set that contains this message type
		SET descriptor_set_json = _pb_wkt_any_find_descriptor_set(type_name, descriptor_set_jsons);
		IF descriptor_set_json IS NULL THEN
			SET message_text = CONCAT('_pb_wkt_any_number_json_to_json: no descriptor found for type: ', type_name);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		SET inner_number_json = _pb_message_to_number_json(descriptor_set_json, type_name, message, NULL);
		-- TODO: set json_marshal_option
		CALL _pb_number_json_to_json_proc(descriptor_set_json, type_name, inner_number_json, FALSE, inner_json);
		RETURN JSON_SET(inner_json, '$."@type"', type_url);
	END IF;
END $$

-- Convert google.protobuf.Any from ProtoJSON to ProtoNumberJSON
DROP FUNCTION IF EXISTS _pb_wkt_any_json_to_number_json $$
CREATE FUNCTION _pb_wkt_any_json_to_number_json(proto_json_value JSON, descriptor_set_jsons JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE type_url TEXT;
	DECLARE type_name TEXT;
	DECLARE message_descriptor JSON;
	DECLARE value_json JSON;
	DECLARE encoded_message LONGBLOB;
	DECLARE value_base64 TEXT;
	DECLARE i INT DEFAULT 0;
	DECLARE descriptor_count INT;
	DECLARE current_descriptor_set JSON;
	DECLARE descriptor_set_json JSON;
	DECLARE wire_json JSON;
	DECLARE inner_number_json JSON;
	DECLARE message LONGBLOB;
	DECLARE message_text TEXT;

	-- Extract @type field
	SET type_url = JSON_UNQUOTE(JSON_EXTRACT(proto_json_value, '$.\"@type\"'));

	-- Extract type name from type_url
	SET type_name = _pb_wkt_any_extract_type_name(type_url);

	-- If no content besides @type, return just type_url
	-- IF value_json IS NULL OR JSON_LENGTH(value_json) = 0 THEN
	-- 	RETURN JSON_OBJECT('1', COALESCE(type_url, ''));
	-- END IF;

	-- Find the descriptor set that contains this message type
	SET descriptor_set_json = _pb_wkt_get_descriptor_set(type_name);
	IF descriptor_set_json IS NOT NULL THEN -- wkt
		-- TODO: support enum?
		CALL _pb_convert_json_wkt_to_number_json(11, type_name, JSON_EXTRACT(proto_json_value, '$.value'), descriptor_set_json, inner_number_json);
		SET message = _pb_number_json_to_message(descriptor_set_json, type_name, inner_number_json, NULL);
	ELSE -- non-wkt
		SET descriptor_set_json = _pb_wkt_any_find_descriptor_set(type_name, descriptor_set_jsons);
		IF descriptor_set_json IS NULL THEN
			SET message_text = CONCAT('_pb_wkt_any_json_to_number_json: no descriptor found for type: ', type_name);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		SET inner_number_json = _pb_json_to_number_json(descriptor_set_json, type_name, JSON_REMOVE(proto_json_value, '$.\"@type\"'), NULL);
		SET message = _pb_number_json_to_message(descriptor_set_json, type_name, inner_number_json, NULL);
	END IF;

	-- Return ProtoNumberJSON format: field 1 = type_url, field 2 = base64 value
	RETURN JSON_OBJECT(
		'1', type_url,
		'2', TO_BASE64(message)
	);
END $$

-- Helper function to decode Any value field
DROP FUNCTION IF EXISTS _pb_wkt_any_decode_value $$
CREATE FUNCTION _pb_wkt_any_decode_value(type_url TEXT, value_base64 TEXT, descriptor_set_jsons JSON) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE type_name TEXT;
	DECLARE message_descriptor JSON;
	DECLARE decoded_message LONGBLOB;
	DECLARE decoded_json JSON;
	DECLARE i INT DEFAULT 0;
	DECLARE descriptor_count INT;
	DECLARE current_descriptor_set JSON;
	DECLARE first_descriptor_set JSON;
	DECLARE wire_json JSON;

	-- Extract type name from type_url
	SET type_name = _pb_wkt_any_extract_type_name(type_url);

	-- Decode base64 value to binary message
	IF value_base64 IS NOT NULL AND value_base64 != '' THEN
		SET decoded_message = FROM_BASE64(value_base64);
	ELSE
		RETURN JSON_OBJECT();
	END IF;

	-- Find the descriptor set that contains this message type
	SET first_descriptor_set = _pb_wkt_any_find_descriptor_set(type_name, descriptor_set_jsons);
	IF first_descriptor_set IS NOT NULL THEN
		SET message_descriptor = _pb_descriptor_set_get_message_descriptor(first_descriptor_set, type_name);
	END IF;

	-- Convert binary message to JSON
	IF message_descriptor IS NOT NULL THEN
		SET decoded_json = pb_message_to_json(first_descriptor_set, type_name, decoded_message, NULL, NULL);
	ELSE
		-- Fallback: convert to wire JSON without schema
		SET wire_json = pb_message_to_wire_json(decoded_message);
		SET decoded_json = wire_json;
	END IF;

	RETURN COALESCE(decoded_json, JSON_OBJECT());
END $$

-- Helper function to encode Any value field
DROP FUNCTION IF EXISTS _pb_wkt_any_encode_value $$
CREATE FUNCTION _pb_wkt_any_encode_value(type_url TEXT, value_json JSON, descriptor_set_jsons JSON) RETURNS TEXT DETERMINISTIC
BEGIN
	DECLARE type_name TEXT;
	DECLARE message_descriptor JSON;
	DECLARE encoded_message LONGBLOB;
	DECLARE i INT DEFAULT 0;
	DECLARE descriptor_count INT;
	DECLARE current_descriptor_set JSON;
	DECLARE first_descriptor_set JSON;
	DECLARE wire_json JSON;

	-- Extract type name from type_url
	SET type_name = _pb_wkt_any_extract_type_name(type_url);

	-- Find the descriptor set that contains this message type
	SET first_descriptor_set = _pb_wkt_any_find_descriptor_set(type_name, descriptor_set_jsons);
	IF first_descriptor_set IS NOT NULL THEN
		SET message_descriptor = _pb_descriptor_set_get_message_descriptor(first_descriptor_set, type_name);
	END IF;

	-- Convert JSON to binary message
	IF message_descriptor IS NOT NULL AND value_json IS NOT NULL THEN
		SET encoded_message = _pb_json_to_message_with_descriptor(value_json, message_descriptor, descriptor_set_json);
	ELSE
		-- Fallback: try to encode as raw message without schema
		IF value_json IS NOT NULL THEN
			SET encoded_message = pb_json_to_message(value_json, descriptor_set_json);
		END IF;
	END IF;

	-- Return base64 encoded value
	IF encoded_message IS NOT NULL THEN
		RETURN TO_BASE64(encoded_message);
	ELSE
		RETURN '';
	END IF;
END $$

DELIMITER ;
