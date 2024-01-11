lint:
	docker run --rm -v "$(shell pwd)":/app -w /app golangci/golangci-lint:v1.55.2 golangci-lint run -v

unit-tests:
	go test -race -v ./internal/...

integration-tests:
	AMQP_URL="amqp://guest:guest@localhost:5672/" AMQP_QUEUE="test" go test -race -v ./tests/...

deps:
	go mod tidy
	go mod vendor