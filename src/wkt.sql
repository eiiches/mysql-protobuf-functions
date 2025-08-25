DELIMITER $$

-- Helper function to get the appropriate descriptor set for Google well-known types
DROP FUNCTION IF EXISTS _pb_wkt_get_descriptor_set $$
CREATE FUNCTION _pb_wkt_get_descriptor_set(full_type_name TEXT) RETURNS JSON DETERMINISTIC
BEGIN
	-- For Google well-known types, use built-in descriptor functions
	CASE
	WHEN full_type_name IN ('.google.protobuf.Struct', '.google.protobuf.Value', '.google.protobuf.ListValue', '.google.protobuf.NullValue') THEN
		RETURN _pb_wkt_struct_proto();
	WHEN full_type_name = '.google.protobuf.FieldMask' THEN
		RETURN _pb_wkt_field_mask_proto();
	WHEN full_type_name IN ('.google.protobuf.DoubleValue', '.google.protobuf.FloatValue', '.google.protobuf.Int64Value', '.google.protobuf.UInt64Value', '.google.protobuf.Int32Value', '.google.protobuf.UInt32Value', '.google.protobuf.BoolValue', '.google.protobuf.StringValue', '.google.protobuf.BytesValue') THEN
		RETURN _pb_wkt_wrappers_proto();
	WHEN full_type_name = '.google.protobuf.Empty' THEN
		RETURN _pb_wkt_empty_proto();
	WHEN full_type_name = '.google.protobuf.Any' THEN
		RETURN _pb_wkt_any_proto();
	WHEN full_type_name = '.google.protobuf.Timestamp' THEN
		RETURN _pb_wkt_timestamp_proto();
	WHEN full_type_name = '.google.protobuf.Duration' THEN
		RETURN _pb_wkt_duration_proto();
	ELSE
		-- Return NULL for types that don't match or should use regular WKT handling
		RETURN NULL;
	END CASE;
END $$

-- Helper function to convert well-known type from ProtoJSON to ProtoNumberJSON
DROP FUNCTION IF EXISTS _pb_convert_json_wkt_to_number_json $$
CREATE FUNCTION _pb_convert_json_wkt_to_number_json(field_type INT, full_type_name TEXT, proto_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	CASE field_type
	WHEN 14 THEN -- enum
		CASE full_type_name
		WHEN '.google.protobuf.NullValue' THEN
			RETURN CAST(_pb_wkt_null_value_json_to_number_json(proto_json_value) AS JSON);
		ELSE
			RETURN NULL;
		END CASE;
	WHEN 11 THEN -- message
		CASE full_type_name
		WHEN '.google.protobuf.Timestamp' THEN
			RETURN _pb_wkt_timestamp_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.Duration' THEN
			RETURN _pb_wkt_duration_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.FieldMask' THEN
			RETURN _pb_wkt_field_mask_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.Value' THEN
			RETURN _pb_wkt_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.Struct' THEN
			RETURN _pb_wkt_struct_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.ListValue' THEN
			RETURN _pb_wkt_list_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.StringValue' THEN
			RETURN _pb_wkt_string_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.Int64Value' THEN
			RETURN _pb_wkt_int64_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.UInt64Value' THEN
			RETURN _pb_wkt_uint64_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.Int32Value' THEN
			RETURN _pb_wkt_int32_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.UInt32Value' THEN
			RETURN _pb_wkt_uint32_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.BoolValue' THEN
			RETURN _pb_wkt_bool_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.FloatValue' THEN
			RETURN _pb_wkt_float_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.DoubleValue' THEN
			RETURN _pb_wkt_double_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.BytesValue' THEN
			RETURN _pb_wkt_bytes_value_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.Empty' THEN
			RETURN _pb_wkt_empty_json_to_number_json(proto_json_value);
		WHEN '.google.protobuf.Any' THEN
			RETURN _pb_wkt_any_json_to_number_json(proto_json_value);
		ELSE
			RETURN NULL;
		END CASE;
	END CASE;
END $$

-- Helper function to convert well-known type from ProtoNumberJSON to ProtoJSON
DROP FUNCTION IF EXISTS _pb_convert_number_json_to_wkt $$
CREATE FUNCTION _pb_convert_number_json_to_wkt(field_type INT, full_type_name TEXT, number_json_value JSON) RETURNS JSON DETERMINISTIC
BEGIN
	CASE field_type
	WHEN 14 THEN -- enum
		CASE full_type_name
		WHEN '.google.protobuf.NullValue' THEN
			RETURN _pb_wkt_null_value_number_json_to_json(number_json_value);
		ELSE
			RETURN NULL;
		END CASE;
	WHEN 11 THEN -- message
		CASE full_type_name
		WHEN '.google.protobuf.Timestamp' THEN
			RETURN _pb_wkt_timestamp_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.Duration' THEN
			RETURN _pb_wkt_duration_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.StringValue' THEN
			RETURN _pb_wkt_string_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.Int64Value' THEN
			RETURN _pb_wkt_int64_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.UInt64Value' THEN
			RETURN _pb_wkt_uint64_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.Int32Value' THEN
			RETURN _pb_wkt_int32_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.UInt32Value' THEN
			RETURN _pb_wkt_uint32_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.BoolValue' THEN
			RETURN _pb_wkt_bool_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.FloatValue' THEN
			RETURN _pb_wkt_float_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.DoubleValue' THEN
			RETURN _pb_wkt_double_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.BytesValue' THEN
			RETURN _pb_wkt_bytes_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.Empty' THEN
			RETURN _pb_wkt_empty_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.Value' THEN
			RETURN _pb_wkt_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.Struct' THEN
			RETURN _pb_wkt_struct_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.ListValue' THEN
			RETURN _pb_wkt_list_value_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.FieldMask' THEN
			RETURN _pb_wkt_field_mask_number_json_to_json(number_json_value);
		WHEN '.google.protobuf.Any' THEN
			RETURN _pb_wkt_any_number_json_to_json(number_json_value);
		ELSE
			RETURN NULL;
		END CASE;
	END CASE;
END $$
