.PHONY: dev build templ generate server up down restart logs help

dev: ## Run the dev loop
	find . -name '*.templ' | entr -r $(MAKE) build serve

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

lambda_image     = loficode-api
lambda_repo      = 812100404712.dkr.ecr.us-east-1.amazonaws.com/blog-api
lambda_region    = us-east-1
git_sha          = $(shell git rev-parse --short HEAD)

build-lambda:
	docker build -t $(lambda_image) -f Dockerfile .

login-ecr: ## Log in to ECR
	aws ecr get-login-password --region $(lambda_region) | docker login --username AWS --password-stdin $(lambda_repo)

tag-lambda: ## Tag the Lambda image
	docker tag ${lambda_image} $(lambda_repo):$(git_sha)
	docker tag $(lambda_image) $(lambda_repo):latest

push-lambda: login-ecr tag-lambda ## Push the Lambda image to ECR
	docker push $(lambda_repo):$(git_sha)
	docker push $(lambda_repo):latest

cf-deploy: ## Deploy CloudFormation stack
	@echo "Deploying Lambda image: $(lambda_repo):$(git_sha)"
	aws cloudformation deploy \
		--stack-name loficode-blog \
		--template-file infra.yaml \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameter-overrides LambdaImageUri=$(lambda_repo):$(git_sha) \
		--region $(lambda_region)

deploy: check-clean build-lambda push-lambda cf-deploy ## Build, push, and deploy the Lambda

nuke: ## Delete the CloudFormation stack
	aws cloudformation delete-stack --stack-name loficode-blog --region us-east-1

check-clean:
	@git diff-index --quiet HEAD -- || { echo "❌ Uncommitted changes. Commit them before deploying."; exit 1; }

upload-static: ## Sync public/ to the S3 static site bucket
	aws s3 sync public/ s3://loficode-site --delete

invalidate-cache: ## Invalidate CloudFront cache
	aws cloudfront create-invalidation \
		--distribution-id E2PTHK2S97TEHJ \
		--paths "/*"

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
