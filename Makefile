.PHONY: run build up down migrate-up migrate-down lint test test-integration

run:
	go run ./cmd/api

build:
	go build -o bin/kinoza ./cmd/api

up:
	docker compose up --build

down:
	docker compose down

migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down

lint:
	golangci-lint run

test:
	go test ./...

test-integration:
	go test -tags=integration ./...
