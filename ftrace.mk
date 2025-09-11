ALL_SQL_FILES_FTRACED := $(patsubst build/%.sql,build/%.sql.ftraced,$(ALL_SQL_FILES))

build/%.sql.ftraced: build/%.sql mysql-ftrace
	./mysql-ftrace instrument --trace-statements $<

.PHONY: ftrace-instrument
ftrace-instrument: $(ALL_SQL_FILES_FTRACED)

.PHONY: ftrace-init
ftrace-init: ensure-test-database
	./mysql-ftrace init --database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)"

.PHONY: ftrace-load
ftrace-load: $(ALL_SQL_FILES_FTRACED) ftrace-init
	$(foreach file,$(ALL_SQL_FILES_FTRACED),$(MYSQL_COMMAND) < $(file);)

.PHONY: ftrace-report
ftrace-report:
	./mysql-ftrace report --database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)"

.PHONY: clean
clean::
	$(RM) $(ALL_SQL_FILES_FTRACED)
