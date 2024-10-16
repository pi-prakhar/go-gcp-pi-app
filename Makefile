COMPOSE_FILE=./deployments/docker-compose.yml
COMPOSE_FILE_PRODUCTION=./deployments/docker-compose.production.yml
DOCKER_FILE_PATH=./deployments/docker
PROJECT_NAME=go-gcp-pi-app
SERVICES := auth server user
TAGS := 1.0.0
BUILD_CONTEXT ?= .
DOCKER_ACCOUNT=16181181418

.PHONY: up-local up-debug up-production down restart build-images push-images list logs

# Start all services in the background
up-local:
	@echo "Starting services with Docker Compose..."
	docker compose -f $(COMPOSE_FILE) -p ${PROJECT_NAME} --env-file ./env/.env.local.docker up $(FLAGS)

up-debug:
	@echo "Starting services with Docker Compose..."
	docker compose -f $(COMPOSE_FILE) -p ${PROJECT_NAME} --env-file ./env/.env.debug up $(FLAGS)

up-production:
	@echo "Starting services with Docker Compose..."
	docker compose -f $(COMPOSE_FILE_PRODUCTION) -p ${PROJECT_NAME} up $(FLAGS)


# Stop all running services
down:
	@echo "Stopping services with Docker Compose..."
	docker compose -f $(COMPOSE_FILE) -p ${PROJECT_NAME} down

clean:
	@echo "Cleaning all data associated with Docker Compose..."
	docker compose -f $(COMPOSE_FILE) -p ${PROJECT_NAME} down --volumes --remove-orphans
	docker network prune


# Restart all services
restart: down up

# Build images
build-images:
	@for service in $(SERVICES); do \
		for tag in $(TAGS); do \
			docker build \
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

# View logs for current project
logs:
	docker compose -f $(COMPOSE_FILE) -p ${PROJECT_NAME} logs -f

# Push all binaries and secrets to production server
cp-to-prod:
	scp -r ./deployments/swarm/swarm.yml ./secrets ${REMOTE}:/home/prakhar/piapp/