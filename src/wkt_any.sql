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
DROP PROCEDURE IF EXISTS _pb_wkt_any_number_json_to_json $$
CREATE PROCEDURE _pb_wkt_any_number_json_to_json(IN number_json_value JSON, IN descriptor_set_jsons JSON, IN emit_default_values BOOLEAN, OUT result JSON)
proc_label: BEGIN
	DECLARE type_url TEXT;
	DECLARE value_base64 TEXT;
	DECLARE type_name TEXT;
	DECLARE message_descriptor JSON;
	DECLARE decoded_message LONGBLOB;
	DECLARE inner_json JSON;
	DECLARE inner_number_json JSON;
	DECLARE i INT DEFAULT 0;
	DECLARE descriptor_count INT;
	DECLARE current_descriptor_set JSON;
	DECLARE descriptor_set_json JSON;
	DECLARE wire_json JSON;
	DECLARE message LONGBLOB;
	DECLARE message_text TEXT;

	-- Extract type_url (field 1) and value (field 2) from ProtoNumberJSON
	SET type_url = COALESCE(JSON_UNQUOTE(JSON_EXTRACT(number_json_value, '$.\"1\"')), '');
	SET message = COALESCE(FROM_BASE64(JSON_UNQUOTE(JSON_EXTRACT(number_json_value, '$.\"2\"'))), '');

	-- If Any value is empty, use non-WKT form, as required by AnyWithNoType conformance test.
	IF type_url = '' AND message = '' THEN
		SET result = JSON_OBJECT();
		LEAVE proc_label;
	END IF;

	-- Extract type name from type_url
	SET type_name = _pb_wkt_any_extract_type_name(type_url);

	SET descriptor_set_json = _pb_wkt_get_descriptor_set(type_name);
	IF descriptor_set_json IS NOT NULL THEN -- wkt
		SET inner_number_json = _pb_message_to_number_json(descriptor_set_json, type_name, message, NULL);
		CALL _pb_convert_number_json_to_wkt(11, type_name, inner_number_json, descriptor_set_jsons, emit_default_values, inner_json);
		SET result = JSON_OBJECT('@type', type_url, 'value', inner_json);
	ELSE -- non-wkt
		-- Find the descriptor set that contains this message type
		SET descriptor_set_json = _pb_wkt_any_find_descriptor_set(type_name, descriptor_set_jsons);
		IF descriptor_set_json IS NULL THEN
			SET message_text = CONCAT('_pb_wkt_any_number_json_to_json: no descriptor found for type: ', type_name);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;
		SET inner_number_json = _pb_message_to_number_json(descriptor_set_json, type_name, message, NULL);
		CALL _pb_number_json_to_json_proc(descriptor_set_json, type_name, inner_number_json, emit_default_values, inner_json);
		SET result = JSON_SET(inner_json, '$."@type"', type_url);
	END IF;
END $$

-- Convert google.protobuf.Any from ProtoJSON to ProtoNumberJSON
DROP PROCEDURE IF EXISTS _pb_wkt_any_json_to_number_json $$
CREATE PROCEDURE _pb_wkt_any_json_to_number_json(IN proto_json_value JSON, IN descriptor_set_jsons JSON, OUT result JSON)
proc: BEGIN
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

	IF JSON_TYPE(proto_json_value) <> 'OBJECT' THEN
		SET message_text = CONCAT('_pb_wkt_any_json_to_number_json: expected OBJECT, but got: ', JSON_TYPE(proto_json_value));
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	IF JSON_LENGTH(proto_json_value) = 0 THEN
		SET result = JSON_OBJECT();
		LEAVE proc;
	END IF;

	-- Extract @type field
	SET type_url = JSON_UNQUOTE(JSON_EXTRACT(proto_json_value, '$.\"@type\"'));
	IF type_url IS NULL THEN
		SET message_text = CONCAT('_pb_wkt_any_json_to_number_json: @type is required for Any');
		SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
	END IF;

	-- Extract type name from type_url
	SET type_name = _pb_wkt_any_extract_type_name(type_url);

	-- If no content besides @type, return just type_url
	-- IF value_json IS NULL OR JSON_LENGTH(value_json) = 0 THEN
	-- 	SET result = JSON_OBJECT('1', COALESCE(type_url, ''));
	-- 	LEAVE _pb_wkt_any_json_to_number_json;
	-- END IF;

	-- Find the descriptor set that contains this message type
	SET descriptor_set_json = _pb_wkt_get_descriptor_set(type_name);
	IF descriptor_set_json IS NOT NULL THEN -- wkt
		-- TODO: support enum?
		CALL _pb_convert_json_wkt_to_number_json(11, type_name, JSON_EXTRACT(proto_json_value, '$.value'), descriptor_set_jsons, inner_number_json);
		SET message = _pb_number_json_to_message(descriptor_set_json, type_name, inner_number_json, NULL);
	ELSE -- non-wkt
		SET descriptor_set_json = _pb_wkt_any_find_descriptor_set(type_name, descriptor_set_jsons);
		IF descriptor_set_json IS NULL THEN
			SET message_text = CONCAT('_pb_wkt_any_json_to_number_json: no descriptor found for type: ', type_name);
			SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = message_text;
		END IF;

		CALL _pb_json_to_number_json_proc(descriptor_set_json, type_name, JSON_REMOVE(proto_json_value, '$.\"@type\"'), FALSE, FALSE, inner_number_json);
		SET message = _pb_number_json_to_message(descriptor_set_json, type_name, inner_number_json, NULL);
	END IF;

	-- Return ProtoNumberJSON format: field 1 = type_url, field 2 = base64 value (only if message is not empty)
	SET result = JSON_OBJECT('1', type_url);
	IF message IS NOT NULL AND LENGTH(message) > 0 THEN
		SET result = JSON_SET(result, '$.\"2\"', _pb_to_base64(message));
	END IF;
END $$
