rabbitmq:
	docker run -d --name rabbitmq -p5672:5672 -p15672:15672 rabbitmq:3-management

rm-rabbitmq:
	docker rm -f rabbitmq

tidy:
	go mod tidy

.PHONY: rabbitmq rm-rabbitmq build