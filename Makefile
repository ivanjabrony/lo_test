
PROJECT_NAME := test_lo
DOCKER_COMPOSE := docker-compose

GO_PACKAGES := $(shell go list ./...)
GO_TEST_FLAGS := -v -cover

.PHONY: build run test unit-test clean 

build:
	@echo "Building containers..."
	$(DOCKER_COMPOSE) build

run: build
	@echo "Starting application..."
	$(DOCKER_COMPOSE) up

test: unit-test

unit-test:
	@echo "Running unit tests..."
	@go test $(GO_TEST_FLAGS) $(GO_PACKAGES)

clean:
	@echo "Cleaning up Docker resources..."
	$(DOCKER_COMPOSE) down -v --remove-orphans