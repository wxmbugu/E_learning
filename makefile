mongo:
	docker run -d --name e-learning -p 27017:27017 -v ~/mongodb_data:/data/db mongo
startdb:
	docker start e-learning
dbshell:
	docker exec -it e-learning bash
test:
	go test -v -cover ./...
rabbitmq:
	docker run -d --name rabbitmq -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password -p 8080:15672 -p 5672:5672 rabbitmq:3-management
startrabbitmq:
	docker start rabbitmq
server:
	go run main.go
.PHONY: mongo startdb test server dbshell rabbitmq