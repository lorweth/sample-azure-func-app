# Variables
include .env.dev
export

DOCKER_COMPOSE = docker-compose --file docker/docker-compose.yml --project-directory . --project-name ${PROJECT}

# Run azure functions core tool
.PHONY: setup start new
setup: build-tool-img build-binaries

start:
	@echo "Start azure function tool..."
	@echo "\033[1;31m Make sure you have enabled rosetta feature on Docker Desktop \033[0m"
	@${DOCKER_COMPOSE} run --rm --service-ports func sh -c "func start --custom"

new:
	@${DOCKER_COMPOSE} run --rm func sh -c "func new"

# Helper target
.PHONY: build-tool-img build-api-img build-binaries
build-tool-img:
	@docker build -f docker/tool.Dockerfile -t virsavik/az-function-tool:latest .

publish:
	@${FUNC} sh -c "az login --service-principal -u $$ARM_CLIENT_ID --password $$ARM_CLIENT_SECRET --tenant $$ARM_TENANT_ID \
	 && cd build \
	 && func azure functionapp publish $$PROJECT-$$ENVIRONMENT-function-app"

build-binaries:
	@${DOCKER_COMPOSE} run --rm builder sh -c "GOOS=linux GOARCH=amd64 go build -o ./build/handler/ ./cmd/serverd"

open-terminal:
	@${DOCKER_COMPOSE} run --rm func /bin/sh
