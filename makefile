mongo:
	docker exec -it e-learning bash
test:
	go test -cover ./...
server:
	go run main.go
.PHONY: test server