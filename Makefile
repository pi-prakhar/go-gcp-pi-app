COMPOSE_FILE=./deployments/docker-compose.yml
COMPOSE_FILE_PRODUCTION=./deployments/docker-compose.production.yml
DOCKER_FILE_PATH=./deployments/docker
PROJECT_NAME=go-gcp-pi-app
SERVICES := auth server user
TAGS := latest
BUILD_CONTEXT ?= .
DOCKER_ACCOUNT=16181181418

# Colors for terminal output
BLUE := \033[1;34m
GREEN := \033[1;32m
RED := \033[1;31m
YELLOW := \033[1;33m
NC := \033[0m # No Color

.PHONY: up-local up-debug up-production down restart build-images push-images list logs

# Default target when just running 'make'
.DEFAULT_GOAL := help

# Help target to display available commands
help: ## Display this help
	@echo "$(BLUE)Available commands:$(NC)"
	@echo
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# Start all services in the background
up-local: ## Start all services locally
	@echo "$(GREEN)Starting services...$(NC)"
	docker compose -f $(COMPOSE_FILE) -p ${PROJECT_NAME} --env-file ./env/.env.local.docker up $(FLAGS)
	@echo "$(GREEN)Services started successfully!$(NC)"

up-debug: ## Start all services debug mode
	@echo "$(GREEN)Starting services in debug mode...$(NC)"
	docker compose -f $(COMPOSE_FILE) -p ${PROJECT_NAME} --env-file ./env/.env.debug up $(FLAGS)
	@echo "$(GREEN)Services started successfully!$(NC)"

up-production: ## Start all services production mode
	@echo "$(GREEN)Starting services in production mode...$(NC)"
	docker compose -f $(COMPOSE_FILE_PRODUCTION) -p ${PROJECT_NAME} up $(FLAGS)
	@echo "$(GREEN)Services started successfully!$(NC)"

# Stop all running services
down: ## Stop all services
	@echo "$(YELLOW)Stopping services...$(NC)"
	docker compose -f $(COMPOSE_FILE) -p ${PROJECT_NAME} down
	@echo "$(YELLOW)Services stopped successfully!$(NC)"

clean-volumes: ## Remove all project volumes
	@echo "$(RED)Removing volumes...$(NC)"
	@docker volume rm -f $(PROJECT_NAME)_caddy_data $(PROJECT_NAME)_grafana-data $(PROJECT_NAME)_caddy_config 2>/dev/null || true
	@echo "$(RED)Volumes removed!$(NC)"

clean-networks: ## Remove project network
	@echo "$(RED)Removing network...$(NC)"
	@docker network rm $(PROJECT_NAME)_gg-network 2>/dev/null || true
	@echo "$(RED)Network removed!$(NC)"

clean-containers: ## Remove all project containers
	@echo "$(RED)Removing containers...$(NC)"
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) rm -f -s -v
	@echo "$(RED)Containers removed!$(NC)"

clean-all: down clean-containers clean-volumes clean-networks ## Complete cleanup (containers, volumes, networks)
	@echo "$(GREEN)Cleanup completed!$(NC)"

restart: ## Restart a specific service (usage: make restart SERVICE=service-name)
	@if [ "$(SERVICE)" = "" ]; then \
		echo "$(RED)Please specify a service name: make restart SERVICE=service-name$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Restarting $(SERVICE)...$(NC)"
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) restart $(SERVICE)
	@echo "$(GREEN)Service restarted!$(NC)"

# Build images
build-images:
	@for service in $(SERVICES); do \
		for tag in $(TAGS); do \
			docker build \
				--no-cache \
				-t $(DOCKER_ACCOUNT)/${PROJECT_NAME}_$$service:$$tag \
				-f $(DOCKER_FILE_PATH)/$$service.Dockerfile \
				$(BUILD_CONTEXT); \
		done; \
	done

# Push Images to Docker Repository
push-images:
	@for service in $(SERVICES); do \
		for tag in $(TAGS); do \
			docker push $(DOCKER_ACCOUNT)/${PROJECT_NAME}_$$service:$$tag; \
		done; \
	done

# List all services and tags available for build
list:
	@echo "Available services:"
	@for service in $(SERVICES); do echo "  - $$service"; done
	@echo "\nAvailable tags:"
	@for tag in $(TAGS); do echo "  - $$tag"; done

# Development helper targets
dev-logs: ## View logs of a specific service (usage: make dev-logs SERVICE=service-name)
	@if [ "$(SERVICE)" = "" ]; then \
		echo "$(RED)Please specify a service name: make dev-logs SERVICE=service-name$(NC)"; \
		exit 1; \
	fi
	@docker compose -p $(PROJECT_NAME) -f $(COMPOSE_FILE) logs -f $(SERVICE)

# Push all binaries and secrets to production server
cp-to-prod:
	scp -r ./deployments/swarm/swarm.yml ./secrets ${REMOTE}:/home/prakhar/piapp/