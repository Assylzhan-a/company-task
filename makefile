# Makefile

DB_URL=postgres://user:password@localhost:5432/company_db?sslmode=disable

.PHONY: migrate-up migrate-down

migrate-up:
	goose -dir ./migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir ./migrations postgres "$(DB_URL)" down

run:
	go run cmd/api/main.go

build:
	go build -o company-service cmd/api/main.go

test:
	go test ./...

lint:
	golangci-lint run
