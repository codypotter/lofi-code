.PHONY: dev build generate generate-prod serve up down reload restart logs help login-ecr tag-lambda push cf-deploy deploy nuke check-clean upload-static invalidate-cache

## Local Development

dev: ## Run the dev loop
	@find . -name '*.templ' | entr -r $(MAKE) generate serve

generate: ## Generate the static files using "go run"
	@templ generate
	@go run ./cmd/generate

generate-prod: ## Generate production static files; update/replace this for CI usage
	@templ generate
	@ENVIRONMENT=production AWS_REGION=us-east-1 go run ./cmd/generate

serve: ## Start the server locally
	@go run ./cmd/server

up: ## Start the local DynamoDB
	@docker compose up -d

down: ## Stop the local DynamoDB
	@docker compose down

reload: ## Run livereload on the public folder
	@livereload public

logs: ## Tail logs from the local DynamoDB
	@docker compose logs -f

## Lambda/Docker/ECR/CloudFormation Deployment

build: ## Build the Lambda Docker image
	@docker build -t loficode-api -f Dockerfile .

login-ecr: ## Log in to Amazon ECR
	@aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 812100404712.dkr.ecr.us-east-1.amazonaws.com/blog-api

tag-lambda: ## Tag the Lambda image for ECR push
	@docker tag loficode-api 812100404712.dkr.ecr.us-east-1.amazonaws.com/blog-api:$(shell git rev-parse --short HEAD)
	@docker tag loficode-api 812100404712.dkr.ecr.us-east-1.amazonaws.com/blog-api:latest

push: login-ecr tag-lambda ## Push the Lambda image to ECR
	@docker push 812100404712.dkr.ecr.us-east-1.amazonaws.com/blog-api:$(shell git rev-parse --short HEAD)
	@docker push 812100404712.dkr.ecr.us-east-1.amazonaws.com/blog-api:latest

cf-deploy: ## Deploy the CloudFormation stack
	@echo "Deploying Lambda image: 812100404712.dkr.ecr.us-east-1.amazonaws.com/blog-api:$(shell git rev-parse --short HEAD)"
	@aws cloudformation deploy \
		--stack-name loficode-blog \
		--template-file infra.yaml \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameter-overrides \
            LambdaImageUri=812100404712.dkr.ecr.us-east-1.amazonaws.com/blog-api:$(shell git rev-parse --short HEAD) \
		--region us-east-1

deploy: check-clean generate-prod build push cf-deploy upload-static invalidate-cache ## Full deployment: generate prod assets, build & push image, deploy stack, update static assets, invalidate cache

nuke: ## Delete the CloudFormation stack
	@aws cloudformation delete-stack --stack-name loficode-blog --region us-east-1

check-clean: ## Ensure there are no uncommitted changes
	@git diff-index --quiet HEAD -- || { echo "❌ Uncommitted changes. Commit them before deploying."; exit 1; }

upload-static: ## Sync local public/ to the S3 static site bucket
	@aws s3 sync public/ s3://loficode-site --delete

invalidate-cache: ## Invalidate the CloudFront cache
	@aws cloudfront create-invalidation --distribution-id E2PTHK2S97TEHJ --paths "/*"

help: ## List available targets
	@echo "Available commands:"; \
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
