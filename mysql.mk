.PHONY: mysql-run
mysql-run: ensure-test-database
	$(MYSQL_COMMAND) -e "$(COMMAND)"