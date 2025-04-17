# Variables
PROJECT_ID ?= stackpulse-development
REGION ?= europe-west1
SERVICE_NAME ?= slack-mcp-sse
IMAGE_NAME ?= gcr.io/$(PROJECT_ID)/$(SERVICE_NAME):latest

# Load environment variables from .env file if it exists
-include .env
export

# Default environment variables (used if not set in .env)
SLACK_BOT_TOKEN ?= unset
SLACK_TEAM_ID ?= unset
MCP_SERVER_BEARER_TOKEN ?= unset
BASE_HOST ?= unset

# Targets
.PHONY: build deploy

build:
	@echo "Building Docker image $(IMAGE_NAME)..."
	docker build --platform linux/amd64 -t $(IMAGE_NAME) .
	# Optional: Push the image if using a registry other than the project's default GCR
	docker push $(IMAGE_NAME)

deploy: build
	@echo "Deploying $(SERVICE_NAME) to Cloud Run in $(PROJECT_ID) [$(REGION)]..."
	gcloud run deploy $(SERVICE_NAME) \
		--image=$(IMAGE_NAME) \
		--project=$(PROJECT_ID) \
		--region=$(REGION) \
		--platform=managed \
		--allow-unauthenticated \
		--set-env-vars "SLACK_BOT_TOKEN=${SLACK_BOT_TOKEN},SLACK_TEAM_ID=${SLACK_TEAM_ID},MCP_SERVER_BEARER_TOKEN=${MCP_SERVER_BEARER_TOKEN},BASE_HOST=${BASE_HOST}"

# You can override variables like this:
# make deploy PROJECT_ID=my-other-project REGION=us-central1
