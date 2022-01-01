test:
	go test -cover ./...
server:
	go run main.go
.PHONY: test server