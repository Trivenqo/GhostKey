.PHONY: setup run test lint

setup:
	docker-compose up -d postgres redis kafka

run:
	go run cmd/api/main.go

test:
	go test -v ./...

tidy:
	go mod tidy