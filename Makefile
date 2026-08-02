include .env
export
APP=cmd/app/main.go

run:
	air

build:
	go build -o bin/app $(APP)

test:
	go test ./...

fmt:
	go fmt ./...

# LINTER
lint:
	golangci-lint run

# API DOCUMENTATION
swag:
	swag init -g cmd/app/main.go -o docs

# DATABASE MIGRATIONS
migrate-create:
	goose -dir migrations create $(name) sql

migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-up-to:
	goose -dir migrations postgres "$(DATABASE_URL)" up-to $(version)

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-down-to:
	goose -dir migrations postgres "$(DATABASE_URL)" down-to $(version)

migrate-status:
	goose -dir migrations postgres "$(DATABASE_URL)" status