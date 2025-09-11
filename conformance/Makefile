# MySQL Protobuf Conformance Testing Makefile
SHELL = /bin/bash

# Protobuf submodule paths
PROTOBUF_DIR = protobuf
CONFORMANCE_RUNNER = $(PROTOBUF_DIR)/conformance_test_runner
CONFORMANCE_PROTO = $(PROTOBUF_DIR)/conformance/conformance.proto

# MySQL conformance test program
MYSQL_CONFORMANCE_DIR = .
MYSQL_CONFORMANCE_BINARY = ./mysql-conformance
CONFORMANCE_PB_GO = ./conformance.pb.go

# Conformance test message descriptor set
CONFORMANCE_DESCRIPTOR_SET = ./conformance_test_messages.binpb
CONFORMANCE_SCHEMA_SQL = ./conformance_test_messages_schema.pb.sql

# Default target
.PHONY: all
all: $(CONFORMANCE_RUNNER) $(MYSQL_CONFORMANCE_BINARY) $(CONFORMANCE_SCHEMA_SQL)

# Build the official conformance test runner
.PHONY: build-runner
build-runner: $(CONFORMANCE_RUNNER)

$(CONFORMANCE_RUNNER):
	@echo "Building conformance test runner..."
	cd $(PROTOBUF_DIR) && cmake . -Dprotobuf_BUILD_CONFORMANCE=ON
	cd $(PROTOBUF_DIR) && cmake --build . -j4 --target conformance_test_runner

# Generate Go protobuf code from conformance.proto
$(CONFORMANCE_PB_GO): $(CONFORMANCE_PROTO)
	@echo "Generating Go protobuf code for conformance.proto..."
	protoc --proto_path=$(PROTOBUF_DIR)/conformance \
		--go_out=$(MYSQL_CONFORMANCE_DIR) \
		--go_opt=paths=source_relative \
		--go_opt=Mconformance.proto='github.com/eiiches/mysql-protobuf-functions/conformance;main' \
		$(CONFORMANCE_PROTO)

# Build the MySQL conformance test program
.PHONY: build-mysql-conformance
build-mysql-conformance: $(MYSQL_CONFORMANCE_BINARY)

# Build conformance test message schema
.PHONY: build-schema
build-schema: $(CONFORMANCE_SCHEMA_SQL)

.PHONY: $(MYSQL_CONFORMANCE_BINARY)
$(MYSQL_CONFORMANCE_BINARY): $(CONFORMANCE_PB_GO)
	@echo "Building MySQL conformance test program..."
	cd $(MYSQL_CONFORMANCE_DIR) && go build -o mysql-conformance .

# Generate conformance test message descriptor set
$(CONFORMANCE_DESCRIPTOR_SET):
	@echo "Generating conformance test message descriptor set..."
	protoc --descriptor_set_out=$(CONFORMANCE_DESCRIPTOR_SET) \
		--include_imports \
		--proto_path=$(PROTOBUF_DIR)/src \
		--proto_path=$(PROTOBUF_DIR)/conformance/test_protos \
		--proto_path=$(PROTOBUF_DIR)/editions/golden \
		$(PROTOBUF_DIR)/src/google/protobuf/test_messages_proto2.proto \
		$(PROTOBUF_DIR)/src/google/protobuf/test_messages_proto3.proto \
		$(PROTOBUF_DIR)/conformance/test_protos/test_messages_edition2023.proto \
		$(PROTOBUF_DIR)/editions/golden/test_messages_proto3_editions.proto

# Generate schema SQL function from descriptor set
$(CONFORMANCE_SCHEMA_SQL): $(CONFORMANCE_DESCRIPTOR_SET)
	@echo "Generating conformance test message schema SQL..."
	go run ../cmd/protoc-gen-mysql \
		--descriptor_set_in=$(CONFORMANCE_DESCRIPTOR_SET) \
		--name=conformance_test_messages_schema \
		--generate_methods=false \
		--validate=false \
		--mysql_out=.

# Run conformance tests
.PHONY: test
test: $(CONFORMANCE_RUNNER) $(MYSQL_CONFORMANCE_BINARY)
	@echo "Running MySQL protobuf conformance tests..."
	@OUTPUT_FILE="test-`date +%s`.out"; \
	$(CONFORMANCE_RUNNER) --failure_list expected_failures.txt --output_dir . --enforce_recommended $(MYSQL_CONFORMANCE_BINARY) \
		--database "root@tcp(127.0.0.100:13306)/test" \
		--debug \
		> $$OUTPUT_FILE 2>&1; \
	echo "Test results written to $$OUTPUT_FILE"

# Run conformance tests with specific test filter
.PHONY: test-filter
test-filter: $(CONFORMANCE_RUNNER) $(MYSQL_CONFORMANCE_BINARY)
ifndef TEST_NAME
	$(error TEST_NAME is required. Usage: make test-filter TEST_NAME=SomeTestName)
endif
ifeq ($(shell echo "$(TEST_NAME)" | grep -E "^(Recommended\.|Required\.)" | wc -l),0)
	$(error TEST_NAME must start with "Recommended." or "Required." (with dot). Regex or pattern matching is not supported. Got: $(TEST_NAME))
endif
	@echo "Running filtered MySQL protobuf conformance tests..."
	$(CONFORMANCE_RUNNER) --test "$(TEST_NAME)" --debug --enforce_recommended $(MYSQL_CONFORMANCE_BINARY) \
		--database "root@tcp(127.0.0.100:13306)/test"

# Load conformance test message schema into MySQL
.PHONY: load-schema
load-schema: $(CONFORMANCE_SCHEMA_SQL)
	@echo "Loading conformance test message schema into MySQL..."
	mysql -h 127.0.0.100 -P 13306 -u root test < $(CONFORMANCE_SCHEMA_SQL)

# Clean generated files
.PHONY: clean
clean:
	@echo "Cleaning conformance test files..."
	$(RM) $(CONFORMANCE_PB_GO)
	$(RM) $(MYSQL_CONFORMANCE_BINARY)
	$(RM) $(CONFORMANCE_DESCRIPTOR_SET)
	$(RM) $(CONFORMANCE_SCHEMA_SQL)
	cd $(PROTOBUF_DIR) && $(MAKE) clean || true

# Clean everything including the built protobuf submodule
.PHONY: clean-all
clean-all: clean
	@echo "Cleaning protobuf submodule build..."
	cd $(PROTOBUF_DIR) && git clean -fdx || true

# Show conformance test runner help
.PHONY: help-runner
help-runner: $(CONFORMANCE_RUNNER)
	$(CONFORMANCE_RUNNER) --help || true

# Validate the setup
.PHONY: validate
validate: $(CONFORMANCE_RUNNER) $(MYSQL_CONFORMANCE_BINARY)
	@echo "Validating conformance test setup..."
	@echo "✓ Conformance test runner built: $(CONFORMANCE_RUNNER)"
	@echo "✓ MySQL conformance program built: $(MYSQL_CONFORMANCE_BINARY)"
	@echo "✓ Protobuf Go code generated: $(CONFORMANCE_PB_GO)"
	@echo "Setup validation complete!"

.PHONY: help
help:
	@echo "MySQL Protobuf Conformance Testing"
	@echo ""
	@echo "Available targets:"
	@echo "  all              - Build everything (default)"
	@echo "  build-runner     - Build the official conformance test runner"
	@echo "  build-mysql-conformance - Build the MySQL conformance test program"
	@echo "  build-schema     - Build conformance test message schema"
	@echo "  load-schema      - Load conformance test message schema into MySQL"
	@echo "  test             - Run conformance tests"
	@echo "  test-filter      - Run filtered conformance tests (requires TEST_NAME=...)"
	@echo "  test-debug       - Run conformance tests in debug mode"
	@echo "  validate         - Validate the setup"
	@echo "  help-runner      - Show conformance test runner help"
	@echo "  clean            - Clean generated files"
	@echo "  clean-all        - Clean everything including protobuf build"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make test"
	@echo "  make test-filter TEST_NAME=BinaryProto3.*"
	@echo "  make test-debug"
