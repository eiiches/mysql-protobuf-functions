.PHONY: test
test:
	docker exec -i test-mysql mysql -u root test < debug.sql
	docker exec -i test-mysql mysql -u root test < protobuf.sql
	docker exec -i test-mysql mysql -u root test < assert.sql
	docker exec -i test-mysql mysql -u root test < test.sql

.PHONY: show-logs
show-logs:
	docker exec -i test-mysql mysql -u root test -e 'SELECT * FROM DebugLog';

.PHONY: start-mysql
start-mysql:
	docker run --rm --name test-mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=true -e MYSQL_DATABASE=test mysql

.PHONY: stop-mysql
stop-mysql:
	docker stop test-mysql
