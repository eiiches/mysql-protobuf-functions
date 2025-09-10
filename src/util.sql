DELIMITER $$

DROP FUNCTION IF EXISTS _pb_util_bin_as_int32 $$
CREATE FUNCTION _pb_util_bin_as_int32(b BLOB) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(b) > 4 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_util_bin_as_int32: value must not be longer than 4 bytes.';
	END IF;

	IF LPAD(b, 4, _binary X'00') & _binary X'80000000' = _binary X'00000000' THEN
		RETURN CONV(HEX(b), 16, 10);
	ELSE
		RETURN -(CONV(HEX(~b), 16, 10) + 1);
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_bin_as_uint32 $$
CREATE FUNCTION _pb_util_bin_as_uint32(b BLOB) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	RETURN CONV(HEX(b), 16, 10);
END $$

DROP FUNCTION IF EXISTS _pb_util_bin_as_int64 $$
CREATE FUNCTION _pb_util_bin_as_int64(b BLOB) RETURNS BIGINT DETERMINISTIC
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	IF LENGTH(b) > 8 THEN
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = '_pb_util_bin_as_int64: value must not be longer than 4 bytes.';
	END IF;

	IF LPAD(b, 8, _binary X'00') & _binary X'8000000000000000' = _binary X'0000000000000000' THEN
		RETURN CONV(HEX(b), 16, 10);
	ELSE
		RETURN -(CONV(HEX(~b), 16, 10) + 1);
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_bin_as_uint64 $$
CREATE FUNCTION _pb_util_bin_as_uint64(b BLOB) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	RETURN CONV(HEX(b), 16, 10);
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_int64 $$
CREATE FUNCTION _pb_util_reinterpret_uint64_as_int64(value BIGINT UNSIGNED) RETURNS BIGINT DETERMINISTIC
BEGIN
	IF value <= 0x7fffffffffffffff THEN
		RETURN CAST(value AS SIGNED);
	ELSE
		RETURN value - 18446744073709551616; -- 2^64
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint32_as_int32 $$
CREATE FUNCTION _pb_util_reinterpret_uint32_as_int32(value INT UNSIGNED) RETURNS INT DETERMINISTIC
BEGIN
	IF value <= 0x7fffffff THEN
		RETURN CAST(value AS SIGNED);
	ELSE
		RETURN CAST(value AS SIGNED) - 4294967296; -- 2^32
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_int32 $$
CREATE FUNCTION _pb_util_reinterpret_uint64_as_int32(value BIGINT UNSIGNED) RETURNS INT DETERMINISTIC
BEGIN
	DECLARE truncated_value INT UNSIGNED;

	-- Step 1: Truncate to 32-bit by masking with 0xFFFFFFFF
	SET truncated_value = CAST(value & 0xFFFFFFFF AS UNSIGNED);

	-- Step 2: Convert to signed 32-bit using existing function
	RETURN _pb_util_reinterpret_uint32_as_int32(truncated_value);
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_uint32 $$
CREATE FUNCTION _pb_util_reinterpret_uint64_as_uint32(value BIGINT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	-- Truncate to 32-bit by masking with 0xFFFFFFFF
	RETURN CAST(value & 0xFFFFFFFF AS UNSIGNED);
END $$

DROP FUNCTION IF EXISTS _pb_util_zigzag_decode_uint64 $$
CREATE FUNCTION _pb_util_zigzag_decode_uint64(value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	RETURN (value >> 1) ^ - (value & 1);
END $$

DROP FUNCTION IF EXISTS _pb_util_zigzag_decode_uint32 $$
CREATE FUNCTION _pb_util_zigzag_decode_uint32(value INT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	IF value & 1 = 0 THEN
		RETURN value >> 1; -- Positive number
	ELSE
		RETURN (value >> 1) ^ 0xFFFFFFFF; -- Negative number
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_zigzag_encode_uint64 $$
CREATE FUNCTION _pb_util_zigzag_encode_uint64(value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	-- ZigZag encoding formula: (n << 1) ^ (n >> 63)
	-- where >> is arithmetic right shift (sign extension)
	--
	-- For signed integers interpreted as unsigned:
	-- - Positive: (n << 1) ^ 0 = 2n
	-- - Negative: (n << 1) ^ -1 = ~(2n) = -2n - 1
	--
	-- Since MySQL's >> is logical shift on unsigned values,
	-- we simulate arithmetic shift: negative numbers have
	-- high bit set, so (value >> 63) = 1, and we need -1
	-- which is 0xFFFFFFFFFFFFFFFF in two's complement
	RETURN (value << 1) ^ -(value >> 63);
END $$

DROP FUNCTION IF EXISTS _pb_util_swap_endian_32 $$
CREATE FUNCTION _pb_util_swap_endian_32(value INT UNSIGNED) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	RETURN ((value & 0xff) << 24)
		| ((value >> 8) & 0xff) << 16
		| ((value >> 16) & 0xff) << 8
		| ((value >> 24) & 0xff);
END $$

DROP FUNCTION IF EXISTS _pb_util_swap_endian_64 $$
CREATE FUNCTION _pb_util_swap_endian_64(value BIGINT UNSIGNED) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	RETURN ((value & 0xff) << 56)
		| ((value >> 8) & 0xff) << 48
		| ((value >> 16) & 0xff) << 40
		| ((value >> 24) & 0xff) << 32
		| ((value >> 32) & 0xff) << 24
		| ((value >> 40) & 0xff) << 16
		| ((value >> 48) & 0xff) << 8
		| ((value >> 56) & 0xff);
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_sint64 $$
CREATE FUNCTION _pb_util_reinterpret_uint64_as_sint64(value BIGINT UNSIGNED) RETURNS BIGINT DETERMINISTIC
BEGIN
	RETURN _pb_util_reinterpret_uint64_as_int64(_pb_util_zigzag_decode_uint64(value));
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_sint32 $$
CREATE FUNCTION _pb_util_reinterpret_uint64_as_sint32(value BIGINT UNSIGNED) RETURNS INT DETERMINISTIC
BEGIN
	-- For sint32: first truncate to 32 bits, then apply zigzag decoding
	-- This follows protobuf specification for sint64 values parsed as sint32
	DECLARE truncated_value INT UNSIGNED;

	-- Step 1: Truncate the varint to 32 bits
	SET truncated_value = CAST(value & 0xFFFFFFFF AS UNSIGNED);

	-- Step 2: Apply zigzag decoding and convert to signed 32-bit
	RETURN _pb_util_reinterpret_uint32_as_int32(_pb_util_zigzag_decode_uint32(truncated_value));
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint64_as_double $$
CREATE FUNCTION _pb_util_reinterpret_uint64_as_double(bits BIGINT UNSIGNED) RETURNS DOUBLE DETERMINISTIC
BEGIN
	DECLARE sign INT;
	DECLARE exponent INT;
	DECLARE fraction DOUBLE;

	SET sign = IF(bits >> 63 = 0, 1, -1); -- sign: +1 or -1
	SET exponent = (bits >> 52) & 0x7FF; -- exponent (11 bits)
	SET fraction = bits & 0xFFFFFFFFFFFFF; -- fraction (52 bits)

	IF exponent = 2047 THEN -- special case
		IF fraction = 0 THEN
			RETURN sign * NULL;  -- +Inf or -Inf
		ELSE
			RETURN NULL; -- NaN
		END IF;
	ELSEIF exponent = 0 THEN -- subnormal number
		RETURN sign * POW(2, -1022) * (fraction / POW(2, 52));
	ELSE -- normal number
		RETURN sign * POW(2, exponent - 1023) * (1 + (fraction / POW(2, 52)));
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_uint32_as_float $$
CREATE FUNCTION _pb_util_reinterpret_uint32_as_float(bits INT UNSIGNED) RETURNS FLOAT DETERMINISTIC
BEGIN
	DECLARE sign INT;
	DECLARE exponent INT;
	DECLARE fraction DOUBLE;

	SET sign = IF(bits >> 31 = 0, 1, -1); -- sign: +1 or -1
	SET exponent = (bits >> 23) & 0xFF; -- exponent (8 bits)
	SET fraction = bits & 0x7FFFFF; -- fraction (23 bits)

	IF exponent = 255 THEN -- special case
		IF fraction = 0 THEN
			RETURN sign * NULL; -- +Inf or -Inf
		ELSE
			RETURN NULL; -- NaN
		END IF;
	ELSEIF exponent = 0 THEN -- subnormal number
		RETURN sign * POW(2, -126) * (fraction / POW(2, 23));
	ELSE -- normal number
		RETURN sign * POW(2, exponent - 127) * (1 + (fraction / POW(2, 23)));
	END IF;
END $$

-- Missing reverse conversion functions needed for setters
DROP FUNCTION IF EXISTS _pb_util_reinterpret_int64_as_uint64 $$
CREATE FUNCTION _pb_util_reinterpret_int64_as_uint64(value BIGINT) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	RETURN CAST(value AS UNSIGNED);
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_int32_as_uint32 $$
CREATE FUNCTION _pb_util_reinterpret_int32_as_uint32(value INT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	-- Handle negative values using 2's complement representation
	IF value < 0 THEN
		RETURN CAST(4294967296 + value AS UNSIGNED);
	ELSE
		RETURN CAST(value AS UNSIGNED);
	END IF;
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_float_as_uint32 $$
CREATE FUNCTION _pb_util_reinterpret_float_as_uint32(value FLOAT) RETURNS INT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE bits BIGINT UNSIGNED;
	DECLARE sign_bit BIGINT UNSIGNED;
	DECLARE exponent BIGINT;
	DECLARE biased_exponent BIGINT;
	DECLARE fraction BIGINT UNSIGNED;

	IF value IS NULL THEN
		RETURN 0x7FC00000; -- NaN
	ELSEIF value != value THEN -- NaN check
		RETURN 0x7FC00000;
	END IF;

	-- Handle zero values (including negative zero)
	IF value = 0 THEN
		-- Use string conversion to detect negative zero
		-- Negative zero shows as "-0" in string representation
		SET sign_bit = IF(CAST(value AS CHAR) LIKE '-%', 1, 0);
		-- Return signed zero: +0.0 = 0x00000000, -0.0 = 0x80000000
		RETURN sign_bit << 31;
	END IF;

	-- Capture sign for non-zero values
	SET sign_bit = IF(value < 0, 1, 0);
	SET value = ABS(value);

	-- Check for infinity
	IF value >= 3.4028235e+38 THEN
		RETURN (sign_bit << 31) | 0x7F800000;
	END IF;

	IF value < 1.17549435082228751e-38 THEN -- subnormal threshold
		SET exponent = -127;
		SET biased_exponent = 0;
		SET fraction = ROUND(value / 1.4012985e-45); -- 2^-149
		IF fraction > 0x7FFFFF THEN
			SET fraction = 0x7FFFFF;
		END IF;
	ELSE -- normal number
		SET exponent = FLOOR(LOG(2, value));
		SET biased_exponent = exponent + 127;

		-- Unfortunately, LOG(2, value) is not always correct. E.g. LOG(2, 1.1754943508222875e-38) may return incorrect result.
		-- We already know the unbiased exponent should be >= -126 for normal numbers.
		-- Exponent is adjusted so that 1 <= ABS(value)/POW(2, exponent) < 2.
		IF exponent >= 128 OR POW(2, exponent) > value THEN
			SET exponent = exponent - 1;
			SET biased_exponent = biased_exponent - 1;
		END IF;
		IF value / POW(2, exponent) >= 2 THEN
			SET exponent = exponent + 1;
			SET biased_exponent = biased_exponent + 1;
		END IF;

		IF biased_exponent < 0 THEN
			SET biased_exponent = 0;
			SET fraction = 0;
		ELSEIF biased_exponent >= 255 THEN
			RETURN (sign_bit << 31) | 0x7F800000; -- infinity
		ELSE
			SET fraction = ROUND((value / POW(2, exponent) - 1) * POW(2, 23));
			IF fraction > 0x7FFFFF THEN
				SET fraction = 0x7FFFFF;
			END IF;
		END IF;
	END IF;

	RETURN (sign_bit << 31) | (CAST(biased_exponent AS UNSIGNED) << 23) | (fraction & 0x7FFFFF);
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_sint64_as_uint64 $$
CREATE FUNCTION _pb_util_reinterpret_sint64_as_uint64(value BIGINT) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	RETURN _pb_util_zigzag_encode_uint64(_pb_util_reinterpret_int64_as_uint64(value));
END $$

DROP FUNCTION IF EXISTS _pb_util_reinterpret_double_as_uint64 $$
CREATE FUNCTION _pb_util_reinterpret_double_as_uint64(value DOUBLE) RETURNS BIGINT UNSIGNED DETERMINISTIC
BEGIN
	DECLARE bits BIGINT UNSIGNED;
	DECLARE sign_bit BIGINT UNSIGNED;
	DECLARE exponent BIGINT;
	DECLARE biased_exponent BIGINT;
	DECLARE fraction BIGINT UNSIGNED;

	IF value IS NULL THEN
		RETURN 0x7FF8000000000000; -- NaN
	ELSEIF value != value THEN -- NaN check
		RETURN 0x7FF8000000000000;
	END IF;

	-- Handle zero values (including negative zero)
	IF value = 0 THEN
		-- Use string conversion to detect negative zero
		-- Negative zero shows as "-0" in string representation
		SET sign_bit = IF(CAST(value AS CHAR) LIKE '-%', 1, 0);
		-- Return signed zero: +0.0 = 0x0000000000000000, -0.0 = 0x8000000000000000
		RETURN sign_bit << 63;
	END IF;

	-- Capture sign for non-zero values
	SET sign_bit = IF(value < 0, 1, 0);
	SET value = ABS(value);

	-- Check for infinity. This never happens because MySQL doesn't support Inf or -Inf.
	-- From this condition, we know the unbiased exponent is less than 1024 (and not equal to or greater than 1024).
	IF value > 1.7976931348623157e+308 THEN
		RETURN (sign_bit << 63) | 0x7FF0000000000000;
	END IF;

	IF value < 2.2250738585072014e-308 THEN -- subnormal threshold
		SET exponent = -1023;
		SET biased_exponent = 0;
		SET fraction = ROUND(value / 4.9406564584124654e-324); -- 2^-1074
		IF fraction > 0xFFFFFFFFFFFFF THEN
			SET fraction = 0xFFFFFFFFFFFFF;
		END IF;
	ELSE -- normal number
		SET exponent = FLOOR(LOG(2, value));
		SET biased_exponent = exponent + 1023;

		-- Unfortunately, LOG(2, value) is not always correct. E.g. LOG(2, 1.7976931348623157e+308) may return incorrect result.
		-- We already know the unbiased exponent should be >= -1022 for normal numbers.
		-- Exponent is adjusted so that 1 <= ABS(value)/POW(2, exponent) < 2.
		IF exponent >= 1024 OR POW(2, exponent) > value THEN
			SET exponent = exponent - 1;
			SET biased_exponent = biased_exponent - 1;
		END IF;
		IF value / POW(2, exponent) >= 2 THEN
			SET exponent = exponent + 1;
			SET biased_exponent = biased_exponent + 1;
		END IF;

		IF biased_exponent < 0 THEN
			SET biased_exponent = 0;
			SET fraction = 0;
		ELSEIF biased_exponent >= 2047 THEN
			RETURN (sign_bit << 63) | 0x7FF0000000000000; -- infinity
		ELSE
			SET fraction = ROUND((value / POW(2, exponent) - 1) * POW(2, 52));
			IF fraction > 0xFFFFFFFFFFFFF THEN
				SET fraction = 0xFFFFFFFFFFFFF;
			END IF;
		END IF;
	END IF;

	RETURN (sign_bit << 63) | (CAST(biased_exponent AS UNSIGNED) << 52) | (fraction & 0xFFFFFFFFFFFFF);
END $$

-- Convert unsigned integer to binary (little-endian)
DROP FUNCTION IF EXISTS _pb_util_uint32_to_bin $$
CREATE FUNCTION _pb_util_uint32_to_bin(value INT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	RETURN UNHEX(CONCAT(
		LPAD(HEX((value >> 24) & 0xFF), 2, '0'),
		LPAD(HEX((value >> 16) & 0xFF), 2, '0'),
		LPAD(HEX((value >> 8) & 0xFF), 2, '0'),
		LPAD(HEX(value & 0xFF), 2, '0')
	));
END $$

-- Convert unsigned 64-bit integer to binary (little-endian)
DROP FUNCTION IF EXISTS _pb_util_uint64_to_bin $$
CREATE FUNCTION _pb_util_uint64_to_bin(value BIGINT UNSIGNED) RETURNS LONGBLOB DETERMINISTIC
BEGIN
	DECLARE low32 INT UNSIGNED DEFAULT value & 0xFFFFFFFF;
	DECLARE high32 INT UNSIGNED DEFAULT (value >> 32) & 0xFFFFFFFF;
	RETURN CONCAT(_pb_util_uint32_to_bin(high32), _pb_util_uint32_to_bin(low32));
END $$
