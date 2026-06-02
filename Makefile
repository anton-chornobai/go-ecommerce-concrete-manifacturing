.SILENT:

GO_FILES := $(wildcard *.go)

formatall:
	go fmt ./...
lint:
	golangci-lint run
migration_up:	
	goose -dir ./internal/migrations postgres ${DB_CONN_STR} up
migration_down:	
	goose -dir ./internal/migrations postgres ${DB_CONN_STR} down
migration_test:
	goose -dir ./internal/migrations postgres ${DB_TEST_CONN_STR} up
test: 
	go test -v ./...