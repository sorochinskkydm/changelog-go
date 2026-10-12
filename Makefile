.PHONY: run build test lint up down

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

lint:
	golangci-lint run

up:
	docker compose up -d

down:
	docker compose down