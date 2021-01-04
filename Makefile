DOCKER_COMPOSE_ARGS ?= -f deployments/docker-compose.yml

.PHONY: dev-docker-compose-down
dev-docker-compose-down: ## stop container network
	@docker-compose ${DOCKER_COMPOSE_ARGS} down -v

.PHONY: dev-docker-compose-up
dev-docker-compose-up: ## start container network
	@docker-compose ${DOCKER_COMPOSE_ARGS} up -d --build

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
