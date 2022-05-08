mongo:
	docker run -d --name e-learning -p 27017:27017 -v ~/mongodb_data:/data/db mongo
startdb:
	docker start e-learning
dbshell:
	docker exec -it e-learning bash
test:
	go test -v -cover ./...
server:
	go run main.go
.PHONY: mongo startdb test server dbshell