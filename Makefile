# Default target
.DEFAULT_GOAL := help

.PHONY: build run fmt lint staticcheck tidy clean test test-cover open-cover help

BINARY := wingit-mcp

build: ## Build the stdio server
	@go build -o bin/$(BINARY) ./cmd/wingit-mcp

run: build ## Run (normally your MCP host launches this)
	@./bin/$(BINARY)

fmt: ## go fmt
	@go fmt ./...

lint: ## go vet + staticcheck
	@GOCACHE=/tmp/go-build go vet ./...
	@$(MAKE) staticcheck

staticcheck: ## Run staticcheck
	@if command -v staticcheck >/dev/null 2>&1; then \
	  XDG_CACHE_HOME=/tmp GOCACHE=/tmp/go-build staticcheck ./...; \
	else \
	  echo "staticcheck not found. Install with: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	  exit 1; \
	fi

tidy: ## go mod tidy
	@go mod tidy

test: ## Run tests with race detector and coverage
	@echo "Running tests..."
	@if command -v gotestsum >/dev/null 2>&1; then \
	  gotestsum --format=short-verbose -- -count=1 -race -covermode=atomic -coverpkg=./... -coverprofile=coverage.out ./...; \
	else \
	  go test -count=1 -race -covermode=atomic -coverpkg=./... -coverprofile=coverage.out ./...; \
	fi
	@echo "Coverage summary:"; \
	go tool cover -func=coverage.out

test-cover: test ## Generate HTML coverage report
	@go tool cover -html=coverage.out -o coverage.html
	@echo "🗂  Wrote coverage.html"

open-cover: test-cover ## Open HTML coverage report (macOS-friendly)
	@command -v open >/dev/null 2>&1 && open coverage.html || true

clean: ## Remove build and test artifacts
	@rm -rf bin
	@rm -f coverage.out coverage.html junit.xml
	@rm -rf test-results/

help: ## Show this help
	@echo "Make targets:"; \
	awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_.-]+:.*##/ {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
