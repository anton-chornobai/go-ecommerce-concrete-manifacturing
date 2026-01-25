GO_FILES := $(wildcard *.go)

formatall:
	go fmt ./...
lint:
	golangci-lint run
