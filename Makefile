APP_NAME := mya

format: 
	go fmt
lint:
	golangci-lint run
