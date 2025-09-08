ALL_SQL_FILES_FTRACED := $(patsubst build/%.sql,build/%.sql.ftraced,$(ALL_SQL_FILES))

# Function tracing targets
.PHONY: instrument-ftrace-files
instrument-ftrace-files: build/protobuf.sql.ftraced build/protobuf-json.sql.ftraced

build/protobuf.sql.ftraced: build/protobuf.sql cmd/mysql-ftrace/main.go
	go tool pigeon -o internal/mysql/sqlflowparser/mysql_ast_parser.go internal/mysql/sqlflowparser/mysql_ast.peg
	go run cmd/mysql-ftrace/main.go instrument --trace-statements build/protobuf.sql

build/protobuf-json.sql.ftraced: build/protobuf-json.sql cmd/mysql-ftrace/main.go
	go tool pigeon -o internal/mysql/sqlflowparser/mysql_ast_parser.go internal/mysql/sqlflowparser/mysql_ast.peg
	go run cmd/mysql-ftrace/main.go instrument --trace-statements build/protobuf-json.sql

.PHONY: load-ftrace-instrumented-files
load-ftrace-instrumented-files: instrument-ftrace-files ensure-test-database
	go run cmd/mysql-ftrace/main.go init --database "root@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)"
	$(MYSQL_COMMAND) < build/protobuf.sql.ftraced
	$(MYSQL_COMMAND) < build/protobuf-json.sql.ftraced
	@echo ""

.PHONY: clean
clean::
	$(RM) $(ALL_SQL_FILES_FTRACED)
