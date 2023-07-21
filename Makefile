.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: test
test: ## Run tests
	go test -cover -race ./...

.PHONY: cover
cover: ## Generate Go coverage report
	@echo "mode: count" > coverage.out
	@go test -coverprofile coverage.tmp ./...
	@tail -n +2 coverage.tmp >> coverage.out
	@go tool cover -html=coverage.out

.PHONY: mocks
mocks: ## Generate mocks for unit tests
	mockery
