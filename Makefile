.DEFAULT_GOAL := help

.PHONY: clean
clean: ## clean temporary files
	go clean

.PHONY: clean
fmt: ## proper code formatting
	go fmt

.PHONY: build
build: clean fmt  ## build this project
	GOOS=linux GOARCH=amd64 go build lambda_client.go

.PHONY: run
run: build ## start this project with a sample call
	echo -e "{\"foo\": \"bar\"}" | ./lambda_client

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
