ALL_SQL_FILES_INSTRUMENTED := $(patsubst build/%.sql,build/%.sql.instrumented,$(ALL_SQL_FILES))

build/%.sql.instrumented: build/%.sql mysql-coverage
	./mysql-coverage instrument $<

.PHONY: coverage
coverage: purge coverage-load coverage-run coverage-report-html
	xdg-open coverage-html/index.html

.PHONY: coverage-ci
coverage-ci: purge coverage-load coverage-run coverage-report-lcov
	go run cmd/mysql-coverage/main.go github-comment --lcov-file coverage.lcov --exclude-file _pb_options_proto.pb.sql --exclude-file well_known_proto.pb.sql --github-pr-number $(GITHUB_PR_NUMBER)

.PHONY: coverage-instrument
coverage-instrument: $(ALL_SQL_FILES_INSTRUMENTED)

.PHONY: coverage-init
coverage-init: ensure-test-database
	go run cmd/mysql-coverage/main.go init --database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)"

.PHONY: coverage-load
coverage-load: $(ALL_SQL_FILES_INSTRUMENTED) coverage-init
	$(foreach file,$(ALL_SQL_FILES_INSTRUMENTED),$(MYSQL_COMMAND) < $(file);)

.PHONY: coverage-run
coverage-run: coverage-init
	go test ./tests -database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" -fuzz-iterations 20 $${GO_TEST_FLAGS:-}

.PHONY: coverage-report-html
coverage-report-html: coverage-report-lcov
	genhtml coverage.lcov --output-directory coverage-html --title "MySQL Protobuf Functions Coverage Report" --flat --exclude "build/well_known_proto.pb.sql" --exclude "build/_pb_options_proto.pb.sql" --missed --show-navigation
	@echo ""
	@echo "=== COVERAGE REPORT GENERATED ==="
	@echo "HTML Report: coverage-html/index.html"
	@echo "LCOV Data: coverage.lcov"

.PHONY: coverage-report-lcov
coverage-report-lcov:
	go run cmd/mysql-coverage/main.go lcov --database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)" --output coverage.lcov
	@echo ""
	@echo "=== LCOV COVERAGE DATA GENERATED ==="
	@echo "LCOV Data: coverage.lcov"
	@echo ""

.PHONY: clean
clean::
	$(RM) $(ALL_SQL_FILES_INSTRUMENTED)
