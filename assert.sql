DELIMITER $$

DROP PROCEDURE IF EXISTS assert_uint_eq $$
CREATE PROCEDURE assert_uint_eq(IN actual BIGINT UNSIGNED, IN expected BIGINT UNSIGNED, IN expression TEXT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE message_text TEXT;
	DECLARE call_expr TEXT;

	SET call_expr = CONCAT('assert_int_eq(actual = ', actual, ', expected = ', expected, ', expression = ', expression, ')');
	CALL debug_msg(call_expr);

	IF actual <> expected THEN
		SET message_text = CONCAT(call_expr, ' failed.');
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
	END IF;
END $$

DROP PROCEDURE IF EXISTS assert_int_eq $$
CREATE PROCEDURE assert_int_eq(IN actual BIGINT, IN expected BIGINT, IN expression TEXT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE message_text TEXT;
	DECLARE call_expr TEXT;

	SET call_expr = CONCAT('assert_int_eq(actual = ', actual, ', expected = ', expected, ', expression = ', expression, ')');
	CALL debug_msg(call_expr);

	IF actual <> expected THEN
		SET message_text = CONCAT(call_expr, ' failed.');
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
	END IF;
END $$

DROP PROCEDURE IF EXISTS assert_text_eq $$
CREATE PROCEDURE assert_text_eq(IN actual TEXT, IN expected TEXT, IN expression TEXT)
BEGIN
	DECLARE CUSTOM_EXCEPTION CONDITION FOR SQLSTATE '45000';

	DECLARE message_text TEXT;
	DECLARE call_expr TEXT;

	SET call_expr = CONCAT('assert_int_eq(actual = ', actual, ', expected = ', expected, ', expression = ', expression, ')');
	CALL debug_msg(call_expr);

	IF actual <> expected THEN
		SET message_text = CONCAT(call_expr, ' failed.');
		SIGNAL CUSTOM_EXCEPTION SET MESSAGE_TEXT = message_text;
	END IF;
END $$
