SHELL = /bin/bash

.PHONY: test
test: purge reload
	set -exuo pipefail; \
	host=$$(docker inspect test-mysql | jq -r '.[].NetworkSettings.IPAddress'); \
	go test ./tests -database "root@tcp($$host:3306)/test" $${GO_TEST_FLAGS:-}

.PHONY: await-mysql
await-mysql:
	set -euxo pipefail; \
	until docker exec test-mysql mysql -u root test -e 'select 1'; do \
		sleep 1; \
	done; \
	echo "MySQL is ready.";

.PHONY: build
build:
	go run cmd/protobuf-accessors/main.go > protobuf-accessors.sql

.PHONY: reload
reload: build
	docker exec -i test-mysql mysql -u root test < debug.sql
	docker exec -i test-mysql mysql -u root test < protobuf.sql
	docker exec -i test-mysql mysql -u root test < protobuf-accessors.sql
	docker exec -i test-mysql mysql -u root test < protobuf-descriptor.sql
	docker exec -i test-mysql mysql -u root test < protobuf-json.sql

.PHONY: purge
purge:
	docker exec -i test-mysql mysql -u root test < purge.sql

.PHONY: show-logs
show-logs:
	docker exec test-mysql mysql -u root test -e 'SELECT * FROM DebugLog';

.PHONY: start-mysql
start-mysql:
	docker run -d --rm --name test-mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=true -e MYSQL_DATABASE=test mysql:8.0.17 --performance-schema-events-statements-history-long-size=1000000

.PHONY: start-profiling
start-profiling:
	docker exec test-mysql mysql -u root test -e "UPDATE performance_schema.setup_consumers SET ENABLED = 'YES' WHERE NAME = 'events_statements_history_long';"
	docker exec test-mysql mysql -u root test -e "UPDATE performance_schema.setup_instruments SET ENABLED = 'YES', TIMED = 'YES' WHERE NAME LIKE 'statement/%';"
	docker exec test-mysql mysql -u root test -e "TRUNCATE TABLE performance_schema.events_statements_history_long;"

.PHONY: stop-profiling
stop-profiling:
	docker exec test-mysql mysql -u root test -e 'UPDATE performance_schema.setup_consumers SET ENABLED="NO" WHERE NAME = "events_statements_history_long"; select count(*) from performance_schema.events_statements_history_long;'
	docker exec -i test-mysql mysql -u root test < scripts/perf-report.sql

.PHONY: flamegraph
flamegraph:
	set -exuo pipefail; \
	host=$$(docker inspect test-mysql | jq -r '.[].NetworkSettings.IPAddress'); \
	output=flamegraph-$$(date +%s).svg; \
	go run cmd/mysql-profiler/main.go -database "root@tcp($$host:3306)/test" | flamegraph.pl > $$output \
		&& xdg-open $$output

.PHONY: stop-mysql
stop-mysql:
	docker stop test-mysql

.PHONY: mysql-shell
mysql-shell:
	docker exec -it test-mysql mysql -u root test
