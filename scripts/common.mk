.PHONY: test
test: purge reload ensure-test-database
	go test ./tests -database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" $${GO_TEST_FLAGS:-}

.PHONY: build
build:
	go run cmd/protobuf-accessors/main.go > protobuf-accessors.sql

.PHONY: reload
reload: build ensure-test-database
	$(MYSQL_COMMAND) < debug.sql
	$(MYSQL_COMMAND) < protobuf.sql
	$(MYSQL_COMMAND) < protobuf-accessors.sql
	$(MYSQL_COMMAND) < protobuf-descriptor.sql
	$(MYSQL_COMMAND) < protobuf-json.sql

.PHONY: purge
purge: ensure-test-database
	$(MYSQL_COMMAND) < purge.sql

.PHONY: show-logs
show-logs: ensure-test-database
	$(MYSQL_COMMAND) -e 'SELECT * FROM DebugLog';

.PHONY: start-profiling
start-profiling: ensure-test-database
	$(MYSQL_COMMAND) -e "UPDATE performance_schema.setup_consumers SET ENABLED = 'YES' WHERE NAME = 'events_statements_history_long';"
	$(MYSQL_COMMAND) -e "UPDATE performance_schema.setup_instruments SET ENABLED = 'YES', TIMED = 'YES' WHERE NAME LIKE 'statement/%';"
	$(MYSQL_COMMAND) -e "TRUNCATE TABLE performance_schema.events_statements_history_long;"

.PHONY: stop-profiling
stop-profiling: ensure-test-database
	$(MYSQL_COMMAND) -e 'UPDATE performance_schema.setup_consumers SET ENABLED="NO" WHERE NAME = "events_statements_history_long"; select count(*) from performance_schema.events_statements_history_long;'
	$(MYSQL_COMMAND) < scripts/perf-report.sql

.PHONY: flamegraph
flamegraph:
	set -exuo pipefail; \
	output=flamegraph-$$(date +%s).svg; \
	go run cmd/mysql-profiler/main.go -database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" | flamegraph.pl > $$output \
		&& xdg-open $$output

.PHONY: ensure-test-database
ensure-test-database: download-mysql
	$(MYSQL_COMMAND_NO_DB) -e 'CREATE DATABASE IF NOT EXISTS test';

.PHONY: coverage
coverage: instrument-files load-instrumented-files run-coverage-tests generate-coverage-report
	xdg-open coverage-html/index.html

.PHONY: instrument-files
instrument-files:
	go run cmd/mysql-coverage/main.go instrument protobuf.sql protobuf-accessors.sql protobuf-descriptor.sql protobuf-json.sql

.PHONY: load-instrumented-files
load-instrumented-files: instrument-files ensure-test-database
	$(MYSQL_COMMAND) < purge.sql
	$(MYSQL_COMMAND) < debug.sql
	go run cmd/mysql-coverage/main.go init --database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)"
	$(MYSQL_COMMAND) < protobuf.sql.instrumented
	$(MYSQL_COMMAND) < protobuf-accessors.sql.instrumented
	$(MYSQL_COMMAND) < protobuf-descriptor.sql.instrumented
	$(MYSQL_COMMAND) < protobuf-json.sql.instrumented

.PHONY: run-coverage-tests
run-coverage-tests: load-instrumented-files
	go test ./tests -database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" $${GO_TEST_FLAGS:-}

.PHONY: generate-coverage-report
generate-coverage-report: run-coverage-tests
	go run cmd/mysql-coverage/main.go lcov --database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" --output coverage.lcov
	genhtml coverage.lcov --output-directory coverage-html --title "MySQL Protobuf Functions Coverage Report"
	@echo ""
	@echo "=== COVERAGE REPORT GENERATED ==="
	@echo "HTML Report: coverage-html/index.html"
	@echo "LCOV Data: coverage.lcov"
	@echo ""

.PHONY: format
format:
	go tool gofumpt -l -w .
