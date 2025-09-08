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

