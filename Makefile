SHELL := /bin/bash
GO ?= go

.PHONY: help setup gen gen-check dev dev-web dev-api build lint fmt test test-web test-api e2e storybook clean

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

setup: ## Install all dependencies (node + go tools)
	pnpm install --frozen-lockfile || pnpm install
	cd apps/api && $(GO) mod download
	cd apps/api && $(GO) install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
	cd apps/api && $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint

gen: ## Generate Go server + TS types from openapi/openapi.yaml
	cd apps/api && $(GO) generate ./...
	pnpm run gen:api

gen-check: gen ## Fail if generated code is out of date
	git diff --exit-code -- apps/api/internal/openapi apps/web/src/lib/api/schema.gen.ts

dev: ## Run web + api together
	$(MAKE) -j2 dev-api dev-web

dev-web: ## Run Next.js dev server
	pnpm --filter web dev

dev-api: ## Run Go API server
	cd apps/api && $(GO) run ./cmd/api

build: ## Build web + api
	pnpm --filter web build
	cd apps/api && $(GO) build -o bin/api ./cmd/api

lint: ## Lint everything
	pnpm run lint
	cd apps/api && golangci-lint run ./...

fmt: ## Format everything
	pnpm run lint:fix
	cd apps/api && $(GO) fmt ./...

test: test-web test-api ## Run unit tests

test-web: ## Typecheck the web app
	pnpm run typecheck

test-api: ## Run Go tests
	cd apps/api && $(GO) test ./...

e2e: ## Run Playwright e2e tests
	pnpm --filter e2e test

storybook: ## Run Storybook
	pnpm --filter web storybook

clean:
	rm -rf apps/web/.next apps/web/storybook-static apps/api/bin e2e/playwright-report e2e/test-results
