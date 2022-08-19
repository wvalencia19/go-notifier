.PHONY: help check test fmt vet lint shellcheck fix-fmt compile tidy build integration

help: ## Show this help
	@echo "Help"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'

##
### Code validation
check: ## Run all checks: test lint
	@bash scripts/check.sh

test: ## Run tests for all go packages
	@bash scripts/test.sh

compile: ## Compile the binary
	@bash scripts/compile.sh
