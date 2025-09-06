SHELL = /bin/bash

MYSQL = tools/mysql-9.3.0-linux-glibc2.28-x86_64-minimal
MYSQLD_BIN = ./$(MYSQL)/bin/mysqld
MYSQL_BIN = ./$(MYSQL)/bin/mysql
MYSQL_DATADIR = /tmp/mysql-57915e48-2946-404a-a03b-d7245f8e0cde
MYSQL_HOST = 127.0.0.100
MYSQL_PORT = 13306
MYSQL_BIND_ADDRESS = 127.0.0.100
MYSQL_DATABASE = test
MYSQL_COMMAND = $(MYSQL_BIN) --host $(MYSQL_HOST) --port $(MYSQL_PORT) -u root $(MYSQL_DATABASE)
MYSQL_COMMAND_NO_DB = $(MYSQL_BIN) --host $(MYSQL_HOST) --port $(MYSQL_PORT) -u root

include scripts/common.mk

.PHONY: start-mysql
start-mysql: download-mysql
	if test -f mysql.pid; then echo "MySQL is already running. Delete mysql.pid if you are sure MySQL is stopped."; exit 1; fi
	$(RM) -r $(MYSQL_DATADIR)
	$(MYSQLD_BIN) --console --initialize-insecure --datadir=$(MYSQL_DATADIR)
	nohup $(MYSQLD_BIN) --datadir=$(MYSQL_DATADIR) --bind-address=$(MYSQL_BIND_ADDRESS) --port=$(MYSQL_PORT) --performance-schema-events-statements-history-long-size=1000000 --stored_program_cache=2048 > mysql.out & echo $$! > mysql.pid
	set -euxo pipefail; \
	until $(MYSQL_COMMAND_NO_DB) -e 'select 1'; do \
		sleep 1; \
	done; \
	echo "MySQL is ready.";

.PHONY: stop-mysql
stop-mysql:
	kill -KILL $$(cat mysql.pid)
	$(RM) mysql.pid

.PHONY: mysql-shell
mysql-shell: ensure-test-database
	$(MYSQL_COMMAND)

# ----- linux perf -----

.PHONY: perf-record
perf-record:
	perf record -g -p $$(cat mysql.pid)

.PHONY: perf-flamegraph
perf-flamegraph: tools/FlameGraph
	set -exuo pipefail; \
	output=flamegraph-$$(date +%s).svg; \
	perf script | ./tools/FlameGraph/stackcollapse-perf.pl | ./tools/FlameGraph/flamegraph.pl > $$output; \
		xdg-open $$output

# ----- tools -----

.PHONY: download-flamegraph
download-flamegraph: FlameGraph

tools/FlameGraph:
	mkdir -p tools/
	git clone https://github.com/brendangregg/FlameGraph tools/FlameGraph

.PHONY: download-mysql
download-mysql: $(MYSQL)

tools/mysql-9.3.0-linux-glibc2.28-x86_64-minimal.tar.xz:
	mkdir -p tools/
	curl -f -L -o tools/mysql-9.3.0-linux-glibc2.28-x86_64-minimal.tar.xz.tmp https://dev.mysql.com/get/Downloads/MySQL-9.3/mysql-9.3.0-linux-glibc2.28-x86_64-minimal.tar.xz
	mv tools/mysql-9.3.0-linux-glibc2.28-x86_64-minimal.tar.xz.tmp tools/mysql-9.3.0-linux-glibc2.28-x86_64-minimal.tar.xz

tools/mysql-9.3.0-linux-glibc2.28-x86_64-minimal: tools/mysql-9.3.0-linux-glibc2.28-x86_64-minimal.tar.xz
	cd tools && tar xvf mysql-9.3.0-linux-glibc2.28-x86_64-minimal.tar.xz

tools/mysql-8.0.42-linux-glibc2.17-x86_64-minimal.tar.xz:
	mkdir -p tools/
	curl -f -L -o tools/mysql-8.0.42-linux-glibc2.17-x86_64-minimal.tar.xz.tmp https://dev.mysql.com/get/Downloads/MySQL-8.0/mysql-8.0.42-linux-glibc2.17-x86_64-minimal.tar.xz
	mv tools/mysql-8.0.42-linux-glibc2.17-x86_64-minimal.tar.xz.tmp tools/mysql-8.0.42-linux-glibc2.17-x86_64-minimal.tar.xz

tools/mysql-8.0.42-linux-glibc2.17-x86_64-minimal: tools/mysql-8.0.42-linux-glibc2.17-x86_64-minimal.tar.xz
	cd tools && tar xvf mysql-8.0.42-linux-glibc2.17-x86_64-minimal.tar.xz
