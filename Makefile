include .env
export


GOPATH=$(shell go env GOPATH)
MOCKERY=$(GOPATH)/bin/mockery

run:
	@go run cmd/main.go

migrate:
	@echo "Applying migrations..."
	@go run migrate.go
	@echo "Migrations applied successfully!"

clean-db:
	@echo "Cleaning database..."
	@go run migrate.go clean
	@echo "Database cleaned!"


.PHONY: .install-mockery
.install-mockery:
	@[ -f $(MOCKERY) ] || go install github.com/vektra/mockery/v2@latest

.PHONY: mocks
mocks: .install-mockery
	mockery