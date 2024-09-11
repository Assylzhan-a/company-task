DB_URL=postgres://user:password@localhost:5432/company_db?sslmode=disable

.PHONY: migrate-up migrate-down run build test lint

migrate-up:
	goose -dir ./migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir ./migrations postgres "$(DB_URL)" down

run:
	docker-compose up

build:
	go build -o company-service cmd/api/main.go

test:
	docker-compose -f tests/docker-compose-test.yml run --rm app_test

lint:
	golangci-lint run
