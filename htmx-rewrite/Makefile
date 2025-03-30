.PHONY: dev build templ generate server up down restart logs help

dev: ## Run the dev loop
	find . -name '*.templ' | entr -r make build serve

build: ## Build the static files
	@make templ
	@make generate

serve: ## Start the server
	@go run ./cmd/server

templ: ## Process the .templ files
	templ generate

generate: ## Generate the static files
	go run ./cmd/generate

up: ## Start the Dynamo db
	docker compose up -d

down: ## Stop the Dynamo db
	docker compose down

reload: ## Start livereload server
	livereload public

restart: ## Restart the Dynamo db
	$(MAKE) down
	$(MAKE) up

logs: ## View the logs of the Dynamo db
	docker compose logs -f

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
