DELIMITER $$

DROP TABLE IF EXISTS DebugLog $$
CREATE TABLE IF NOT EXISTS DebugLog (message TEXT) $$

DROP PROCEDURE IF EXISTS `debug_msg` $$
CREATE PROCEDURE debug_msg(IN msg VARCHAR(255))
BEGIN
	INSERT INTO DebugLog (message) VALUES (msg);
END $$
