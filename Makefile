SHELL = /bin/bash

.PHONY: test
test: reload
	set -exuo pipefail; \
	host=$$(docker inspect test-mysql | jq -r '.[].NetworkSettings.IPAddress'); \
	cd tests && go test -database "root@tcp($$host:3306)/test" $${GO_TEST_FLAGS:-}

.PHONY: await-mysql
await-mysql:
	set -euxo pipefail; \
	until docker exec test-mysql mysql -u root test -e 'select 1'; do \
		sleep 1; \
	done; \
	echo "MySQL is ready.";

.PHONY: test-descriptor
test-descriptor: purge reload
	docker exec -i test-mysql mysql -u root test < test-descriptor.sql

.PHONY: reload
reload:
	docker exec -i test-mysql mysql -u root test < debug.sql
	docker exec -i test-mysql mysql -u root test < protobuf.sql
	docker exec -i test-mysql mysql -u root test < protobuf-descriptor.sql
	docker exec -i test-mysql mysql -u root test < protobuf-json.sql

.PHONY: purge
purge:
	docker exec -i test-mysql mysql -u root test < purge.sql

.PHONY: show-logs
show-logs:
	docker exec -i test-mysql mysql -u root test -e 'SELECT * FROM DebugLog';

.PHONY: start-mysql
start-mysql:
	docker run -d --rm --name test-mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=true -e MYSQL_DATABASE=test mysql:8.3.0

.PHONY: stop-mysql
stop-mysql:
	docker stop test-mysql

.PHONY: mysql-shell
mysql-shell:
	docker exec -it test-mysql mysql -u root test
