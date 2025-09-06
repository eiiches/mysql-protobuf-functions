SHELL = /bin/bash

MYSQL_HOST = 127.0.0.100
MYSQL_PORT = 23306
MYSQL_BIND_ADDRESS = 127.0.0.100
MYSQL_DATABASE = test
MYSQL_COMMAND = docker exec -i test-mysql mysql -u root $(MYSQL_DATABASE)
MYSQL_COMMAND_NO_DB = docker exec -i test-mysql mysql -u root
MYSQL_COMMAND_WITH_TERMINAL = docker exec -it test-mysql mysql -u root test

include scripts/common.mk

.PHONY: start-mysql
start-mysql: download-mysql
	docker run -d --rm --name test-mysql --tmpfs /var/lib/mysql -p $(MYSQL_BIND_ADDRESS):$(MYSQL_PORT):3306 -e MYSQL_ALLOW_EMPTY_PASSWORD=true -e MYSQL_DATABASE=test mysql:8.0.17 --performance-schema-events-statements-history-long-size=1000000 --stored_program_cache=2048
	set -euxo pipefail; \
	until $(MYSQL_COMMAND) -e 'select 1'; do \
		sleep 1; \
	done; \
	echo "MySQL is ready.";

.PHONY: stop-mysql
stop-mysql:
	docker stop test-mysql

.PHONY: mysql-shell
mysql-shell: ensure-test-database
	$(MYSQL_COMMAND_WITH_TERMINAL)

.PHONY: download-mysql
download-mysql:
