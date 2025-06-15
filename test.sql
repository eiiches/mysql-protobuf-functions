DELIMITER $$

DROP PROCEDURE IF EXISTS _pb_run_tests $$
CREATE PROCEDURE _pb_run_tests()
BEGIN
	CALL assert_int_eq(pb_wire_read_varint(_binary X'9601'), 150, 'pb_wire_read_varint(_binary X''9601'')');
	CALL assert_int_eq(pb_wire_read_varint(_binary X'960133'), 150, 'pb_wire_read_varint(_binary X''960133'')');

	CALL assert_int_eq(_pb_util_bin_as_int32(_binary X'00'), 0, '_pb_util_bin_as_int32(_binary X''00'')');
	CALL assert_int_eq(_pb_util_bin_as_int32(_binary X'7fffffff'), 2147483647, '_pb_util_bin_as_int32(_binary X''7fffffff'')');
	CALL assert_int_eq(_pb_util_bin_as_int32(_binary X'80000000'), -2147483648, '_pb_util_bin_as_int32(_binary X''80000000'')');
	CALL assert_int_eq(_pb_util_bin_as_int32(_binary X'ffffffff'), -1, '_pb_util_bin_as_int32(_binary X''ffffffff'')');

	CALL assert_uint_eq(_pb_util_bin_as_uint32(_binary X'00'), 0, '_pb_util_bin_as_uint32(_binary X''00'')');
	CALL assert_uint_eq(_pb_util_bin_as_uint32(_binary X'7fffffff'), 2147483647, '_pb_util_bin_as_uint32(_binary X''7fffffff'')');
	CALL assert_uint_eq(_pb_util_bin_as_uint32(_binary X'80000000'), 2147483648, '_pb_util_bin_as_uint32(_binary X''80000000'')');
	CALL assert_uint_eq(_pb_util_bin_as_uint32(_binary X'ffffffff'), 4294967295, '_pb_util_bin_as_uint32(_binary X''ffffffff'')');

	CALL assert_int_eq(_pb_util_bin_as_int64(_binary X'00'), 0, '_pb_util_bin_as_int64(_binary X''00'')');
	CALL assert_int_eq(_pb_util_bin_as_int64(_binary X'7fffffffffffffff'), 9223372036854775807, '_pb_util_bin_as_int64(_binary X''7fffffffffffffff'')');
	CALL assert_int_eq(_pb_util_bin_as_int64(_binary X'8000000000000000'), -9223372036854775808, '_pb_util_bin_as_int64(_binary X''8000000000000000'')');
	CALL assert_int_eq(_pb_util_bin_as_int64(_binary X'ffffffffffffffff'), -1, '_pb_util_bin_as_int64(_binary X''ffffffffffffffff'')');

	CALL assert_uint_eq(_pb_util_bin_as_uint64(_binary X'00'), 0, '_pb_util_bin_as_uint64(_binary X''00'')');
	CALL assert_uint_eq(_pb_util_bin_as_uint64(_binary X'7fffffffffffffff'), 9223372036854775807, '_pb_util_bin_as_uint64(_binary X''7fffffffffffffff'')');
	CALL assert_uint_eq(_pb_util_bin_as_uint64(_binary X'8000000000000000'), 9223372036854775808, '_pb_util_bin_as_uint64(_binary X''8000000000000000'')');
	CALL assert_uint_eq(_pb_util_bin_as_uint64(_binary X'ffffffffffffffff'), 18446744073709551615, '_pb_util_bin_as_uint64(_binary X''ffffffffffffffff'')');

	CALL assert_int_eq(_pb_wire_get_field_number_from_tag(0x08), 1, '_pb_wire_get_field_number_from_tag(0x08)');
	CALL assert_int_eq(_pb_wire_get_wire_type_from_tag(0x08), 0, '_pb_wire_get_wire_type_from_tag(0x08)');

	CALL assert_int_eq(pb_message_get_int32_field(_binary X'100a080a', 1, 0), 10, 'pb_message_get_int32_field(_binary X''100a080a'', 1, 0)');
	CALL assert_int_eq(pb_message_get_int64_field(_binary X'100a080a', 1, 0), 10, 'pb_message_get_int64_field(_binary X''100a080a'', 1, 0)');
	CALL assert_int_eq(pb_message_get_int32_field(_binary X'100a080a', 1, NULL), 10, 'pb_message_get_int32_field(_binary X''100a080a'', 1, NULL)');
	CALL assert_int_eq(pb_message_get_int64_field(_binary X'100a080a', 1, NULL), 10, 'pb_message_get_int64_field(_binary X''100a080a'', 1, NULL)');
	CALL assert_text_eq(pb_message_get_string_field(_binary X'100a2a03616263', 5, NULL), 'abc', 'pb_message_get_string_field(_binary X''100a2a03616263'', 5, NULL)');
	CALL assert_text_eq(pb_message_get_bytes_field(_binary X'100a2a03616263', 5, NULL), _binary 'abc', 'pb_message_get_bytes_field(_binary X''100a2a03616263'', 5, NULL)');

	-- packed repeated int32
	CALL assert_int_eq(pb_message_get_int32_field(_binary X'3a03010203', 7, 0), 1, 'pb_message_get_int32_field(_binary X''3a03010203'', 7, 0)');
	CALL assert_int_eq(pb_message_get_int32_field(_binary X'3a03010203', 7, 1), 2, 'pb_message_get_int32_field(_binary X''3a03010203'', 7, 1)');
	CALL assert_int_eq(pb_message_get_int32_field(_binary X'3a03010203', 7, 2), 3, 'pb_message_get_int32_field(_binary X''3a03010203'', 7, 2)');

	CALL assert_int_eq(pb_message_get_int32_field(_binary X'', 1, NULL), 0, 'pb_message_get_int32_field(_binary X'''', 1, NULL)');
	CALL assert_int_eq(pb_message_get_int32_field(_binary X'08ffffffff07', 1, 0), 2147483647, 'pb_message_get_int32_field(_binary X''08ffffffff07'', 1, 0)');
	CALL assert_int_eq(pb_message_get_int32_field(_binary X'08ffffffffffffffffff01', 1, 0), -1, 'pb_message_get_int32_field(_binary X''08ffffffffffffffffff01'', 1, 0)');
	CALL assert_int_eq(pb_message_get_int32_field(_binary X'0880808080f8ffffffff01', 1, 0), -2147483648, 'pb_message_get_int32_field(_binary X''0880808080f8ffffffff01'', 1, 0)');
END $$
CALL _pb_run_tests;
