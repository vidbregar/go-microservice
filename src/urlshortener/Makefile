REVISION ?= $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_TAG ?= $(shell git tag --points-at HEAD 2>/dev/null)
IMAGE_TAG ?= latest

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show help
	@awk -F ':|##' '/^[^\t].+:.*##/ { printf "\033[36mmake %-28s\033[0m -%s\n", $$1, $$NF }' $(MAKEFILE_LIST) | sort

.PHONY: generate
generate: ## Generate Go server and types from OpenAPI spec
	rm -rf internal/api/oapi/*
	oapi-codegen -generate server,spec -package oapi api/openapi.yaml > internal/api/oapi/server.gen.go
	oapi-codegen -generate types -package oapi api/openapi.yaml > internal/api/oapi/types.gen.go

.PHONY: run
run: ## Start urlshortener in Docker
	docker stop urlshortener || true && docker rm urlshortener || true
	docker run -p 8080:8080 --name urlshortener urlshortener:latest

.PHONY: build
build: ## Builds urlshortener image
	docker image rm k3d-registry.localhost:50000/urlshortener:$(IMAGE_TAG) || true
	docker build --build-arg REVISION=$(REVISION) \
				 --build-arg GIT_TAG=$(GIT_TAG) \
				 -t k3d-registry.localhost:50000/urlshortener:$(IMAGE_TAG) .
.PHONY: push
push: ## Pushes urlshortener image to k3d registry
	docker push k3d-registry.localhost:50000/urlshortener:$(IMAGE_TAG)
