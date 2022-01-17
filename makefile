mongo:
	sudo docker exec -it e-learning bash
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: test server