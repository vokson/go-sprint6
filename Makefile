GO            ?= go
GOLANGCI_LINT ?= golangci-lint

.PHONY: run test lint fmt

run:
	$(GO) run ./cmd

test:
	$(GO) test ./...

lint:
	$(GO) vet ./internal/... ./cmd/...
	$(GOLANGCI_LINT) run ./internal/... ./cmd/...

fmt:
	$(GOLANGCI_LINT) fmt ./...
