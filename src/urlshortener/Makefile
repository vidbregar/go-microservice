SHELL = $$SHELL

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show help
	@awk -F ':|##' '/^[^\t].+:.*##/ { printf "\033[36mmake %-28s\033[0m -%s\n", $$1, $$NF }' $(MAKEFILE_LIST) | sort

generate: ## Generate Go server and types from OpenAPI spec
	oapi-codegen -generate server,spec -package api api/openapi.yaml > internal/api/server.gen.go
	oapi-codegen -generate types -package api api/openapi.yaml > internal/api/types.gen.go

run: ## Start dev urlshortener microservice
	go run cmd/urlshortener/main.go