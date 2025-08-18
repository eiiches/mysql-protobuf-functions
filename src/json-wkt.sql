DELIMITER $$

DROP FUNCTION IF EXISTS _pb_wire_json_decode_wkt_as_json $$
CREATE FUNCTION _pb_wire_json_decode_wkt_as_json(wire_json JSON, full_type_name TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	CASE full_type_name
	WHEN '.google.protobuf.Timestamp' THEN
		RETURN _pb_wire_json_decode_wkt_timestamp_as_json(wire_json);
	WHEN '.google.protobuf.Duration' THEN
		RETURN _pb_wire_json_decode_wkt_duration_as_json(wire_json);
	WHEN '.google.protobuf.Struct' THEN
		RETURN _pb_wire_json_decode_wkt_struct_as_json(wire_json);
	WHEN '.google.protobuf.Value' THEN
		RETURN _pb_wire_json_decode_wkt_value_as_json(wire_json);
	WHEN '.google.protobuf.ListValue' THEN
		RETURN _pb_wire_json_decode_wkt_list_value_as_json(wire_json);
	WHEN '.google.protobuf.Empty' THEN
		RETURN JSON_OBJECT();
	WHEN '.google.protobuf.DoubleValue' THEN
		RETURN CAST(pb_wire_json_get_double_field(wire_json, 1, 0.0) AS JSON);
	WHEN '.google.protobuf.FloatValue' THEN
		RETURN CAST(pb_wire_json_get_float_field(wire_json, 1, 0.0) AS JSON);
	WHEN '.google.protobuf.Int64Value' THEN
		RETURN JSON_QUOTE(CAST(pb_wire_json_get_int64_field(wire_json, 1, 0) AS CHAR));
	WHEN '.google.protobuf.UInt64Value' THEN
		RETURN JSON_QUOTE(CAST(pb_wire_json_get_uint64_field(wire_json, 1, 0) AS CHAR));
	WHEN '.google.protobuf.Int32Value' THEN
		RETURN CAST(pb_wire_json_get_int32_field(wire_json, 1, 0) AS JSON);
	WHEN '.google.protobuf.UInt32Value' THEN
		RETURN CAST(pb_wire_json_get_uint32_field(wire_json, 1, 0) AS JSON);
	WHEN '.google.protobuf.BoolValue' THEN
		RETURN CAST((pb_wire_json_get_bool_field(wire_json, 1, FALSE) IS TRUE) AS JSON);
	WHEN '.google.protobuf.StringValue' THEN
		RETURN JSON_QUOTE(pb_wire_json_get_string_field(wire_json, 1, ''));
	WHEN '.google.protobuf.BytesValue' THEN
		RETURN JSON_QUOTE(TO_BASE64(pb_wire_json_get_bytes_field(wire_json, 1, _binary X'')));
	WHEN '.google.protobuf.FieldMask' THEN
		RETURN _pb_wire_json_decode_wkt_field_mask_as_json(wire_json);
	ELSE
		RETURN NULL;
	END CASE;
END $$

-- Helper function to encode well-known types from JSON to wire_json
DROP FUNCTION IF EXISTS _pb_json_encode_wkt_as_wire_json $$
CREATE FUNCTION _pb_json_encode_wkt_as_wire_json(json_value JSON, full_type_name TEXT, from_number_json BOOLEAN) RETURNS JSON DETERMINISTIC
BEGIN
	DECLARE result JSON;
	DECLARE float_value FLOAT;
	DECLARE double_value DOUBLE;

	CASE full_type_name
	WHEN '.google.protobuf.Timestamp' THEN
		-- Parse RFC 3339 timestamp string like "1996-12-19T16:39:57.000340012Z"
		IF JSON_TYPE(json_value) = 'NULL' THEN
			-- Null timestamp should not produce any field in protobuf
			RETURN NULL;
		ELSEIF JSON_TYPE(json_value) = 'STRING' THEN
			RETURN _pb_json_encode_wkt_timestamp_as_wire_json(JSON_UNQUOTE(json_value));
		END IF;

	WHEN '.google.protobuf.Duration' THEN
		-- Parse duration string like "1.000340012s" or "3600s"
		IF JSON_TYPE(json_value) = 'NULL' THEN
			-- Null duration should not produce any field in protobuf
			RETURN NULL;
		ELSEIF JSON_TYPE(json_value) = 'STRING' THEN
			RETURN _pb_json_encode_wkt_duration_as_wire_json(JSON_UNQUOTE(json_value));
		END IF;

	WHEN '.google.protobuf.FieldMask' THEN
		-- Parse comma-separated field names like "path1,path2"
		IF JSON_TYPE(json_value) = 'STRING' THEN
			RETURN _pb_json_encode_wkt_field_mask_as_wire_json(JSON_UNQUOTE(json_value));
		END IF;

	WHEN '.google.protobuf.Empty' THEN
		-- Always return empty wire_json
		RETURN JSON_OBJECT();

	WHEN '.google.protobuf.Struct' THEN
		-- For number JSON format, use regular descriptor-based processing
		IF from_number_json THEN
			RETURN NULL;
		END IF;
		-- Convert JSON object to Struct with repeated fields map
		CALL _pb_json_encode_wkt_struct_as_wire_json(json_value, from_number_json, result);
		IF result IS NOT NULL THEN
			RETURN result;
		END IF;

	WHEN '.google.protobuf.Value' THEN
		-- For number JSON format, use regular descriptor-based processing
		IF from_number_json THEN
			RETURN NULL;
		END IF;
		-- Handle different JSON value types
		CALL _pb_json_encode_wkt_value_as_wire_json(json_value, from_number_json, result);
		IF result IS NOT NULL THEN
			RETURN result;
		END IF;

	WHEN '.google.protobuf.ListValue' THEN
		-- For number JSON format, use regular descriptor-based processing
		IF from_number_json THEN
			RETURN NULL;
		END IF;
		-- Convert JSON array to ListValue with repeated Value fields
		CALL _pb_json_encode_wkt_list_value_as_wire_json(json_value, from_number_json, result);
		IF result IS NOT NULL THEN
			RETURN result;
		END IF;

	-- Wrapper types
	WHEN '.google.protobuf.Int32Value' THEN
		IF JSON_TYPE(json_value) IN ('INTEGER', 'DECIMAL', 'STRING') THEN
			SET result = JSON_OBJECT();
			-- Only encode non-default values (proto3 behavior)
			IF _pb_json_to_signed_int(json_value) <> 0 THEN
				SET result = pb_wire_json_set_int32_field(result, 1, _pb_json_to_signed_int(json_value));
			END IF;
			RETURN result;
		END IF;

	WHEN '.google.protobuf.UInt32Value' THEN
		IF JSON_TYPE(json_value) IN ('INTEGER', 'DECIMAL', 'STRING') THEN
			SET result = JSON_OBJECT();
			-- Only encode non-default values (proto3 behavior)
			IF _pb_json_to_unsigned_int(json_value) <> 0 THEN
				SET result = pb_wire_json_set_uint32_field(result, 1, _pb_json_to_unsigned_int(json_value));
			END IF;
			RETURN result;
		END IF;

	WHEN '.google.protobuf.Int64Value' THEN
		IF JSON_TYPE(json_value) IN ('INTEGER', 'DECIMAL', 'STRING') THEN
			SET result = JSON_OBJECT();
			-- Only encode non-default values (proto3 behavior)
			IF _pb_json_to_signed_int(json_value) <> 0 THEN
				SET result = pb_wire_json_set_int64_field(result, 1, _pb_json_to_signed_int(json_value));
			END IF;
			RETURN result;
		END IF;

	WHEN '.google.protobuf.UInt64Value' THEN
		IF JSON_TYPE(json_value) IN ('INTEGER', 'DECIMAL', 'STRING') THEN
			SET result = JSON_OBJECT();
			-- Only encode non-default values (proto3 behavior)
			IF _pb_json_to_unsigned_int(json_value) <> 0 THEN
				SET result = pb_wire_json_set_uint64_field(result, 1, _pb_json_to_unsigned_int(json_value));
			END IF;
			RETURN result;
		END IF;

	WHEN '.google.protobuf.FloatValue' THEN
		IF JSON_TYPE(json_value) IN ('INTEGER', 'DECIMAL', 'DOUBLE', 'STRING') THEN
			SET result = JSON_OBJECT();
			IF JSON_TYPE(json_value) = 'STRING' THEN
				SET float_value = CAST(JSON_UNQUOTE(json_value) AS FLOAT);
			ELSE
				SET float_value = CAST(json_value AS FLOAT);
			END IF;
			-- Only encode non-default values (proto3 behavior)
			IF float_value <> 0.0 THEN
				SET result = pb_wire_json_set_float_field(result, 1, float_value);
			END IF;
			RETURN result;
		END IF;

	WHEN '.google.protobuf.DoubleValue' THEN
		IF JSON_TYPE(json_value) IN ('INTEGER', 'DECIMAL', 'DOUBLE', 'STRING') THEN
			SET result = JSON_OBJECT();
			IF JSON_TYPE(json_value) = 'STRING' THEN
				SET double_value = CAST(JSON_UNQUOTE(json_value) AS DOUBLE);
			ELSE
				SET double_value = CAST(json_value AS DOUBLE);
			END IF;
			-- Only encode non-default values (proto3 behavior)
			IF double_value <> 0.0 THEN
				SET result = pb_wire_json_set_double_field(result, 1, double_value);
			END IF;
			RETURN result;
		END IF;

	WHEN '.google.protobuf.BoolValue' THEN
		IF JSON_TYPE(json_value) = 'BOOLEAN' THEN
			SET result = JSON_OBJECT();
			-- Only encode non-default values (proto3 behavior)
			IF json_value THEN
				SET result = pb_wire_json_set_bool_field(result, 1, TRUE);
			END IF;
			RETURN result;
		END IF;

	WHEN '.google.protobuf.StringValue' THEN
		IF JSON_TYPE(json_value) = 'STRING' THEN
			SET result = JSON_OBJECT();
			-- Only encode non-default values (proto3 behavior)
			IF JSON_UNQUOTE(json_value) <> '' THEN
				SET result = pb_wire_json_set_string_field(result, 1, JSON_UNQUOTE(json_value));
			END IF;
			RETURN result;
		END IF;

	WHEN '.google.protobuf.BytesValue' THEN
		IF JSON_TYPE(json_value) = 'STRING' THEN
			SET result = JSON_OBJECT();
			-- Only encode non-default values (proto3 behavior)
			IF JSON_UNQUOTE(json_value) <> '' THEN
				SET result = pb_wire_json_set_bytes_field(result, 1, _pb_util_from_base64_url(JSON_UNQUOTE(json_value)));
			END IF;
			RETURN result;
		END IF;
	END CASE;

	-- Return NULL to fall back to normal message handling
	RETURN NULL;
END $$
