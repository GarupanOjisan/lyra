SHELL := /bin/bash

.PHONY: build test fmt vet run-api run-worker sqlc-generate db-run

build:
	go build ./...

test:
	go test ./... -race -cover

fmt:
	gofmt -s -w .

vet:
	go vet ./...

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

sqlc-generate:
	sqlc generate

# Run API with sqlc-backed repository (requires generated code and DB env)
db-run:
	# Example: make db-run LYRA_DB_DRIVER=pgx LYRA_DB_DSN='postgres://user:pass@localhost:5432/lyra?sslmode=disable'
	LYRA_DB_DRIVER?=$(LYRA_DB_DRIVER)
	LYRA_DB_DSN?=$(LYRA_DB_DSN)
	LYRA_DB_DRIVER=$${LYRA_DB_DRIVER} LYRA_DB_DSN=$${LYRA_DB_DSN} go run -tags sqlc ./cmd/api
