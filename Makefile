.DEFAULT_GOAL := help

APPLICATION?=speedtest-go
COMMIT_SHA?=$(shell git rev-parse --short HEAD)
DOCKER?=docker
REGISTRY?=jonascheng
GOOS?=$(shell uname -s | awk '{print tolower($0)}')

.PHONY: setup
setup: ## setup go modules
	go mod tidy

.PHONY: clean
clean: ## cleans the binary
	go clean
	rm -rf ./bin

.PHONY: run
run: setup ## runs go run the application
	go run -race speedtest.go

.PHONY: test
test: ## runs go test the application
	go test -v ./...

.PHONY: build
build: clean ## build the application
	GOOS=${GOOS} GOARCH=amd64 go build -a -v -ldflags="-w -s" -o bin/${APPLICATION} speedtest.go

.PHONY: docker-login
docker-login: ## login docker registry
ifndef DOCKERHUB_USERNAME
	$(error DOCKERHUB_USERNAME not set on env)
endif
ifndef DOCKERHUB_PASSWORD
	$(error DOCKERHUB_PASSWORD not set on env)
endif
	${DOCKER} login --username ${DOCKERHUB_USERNAME} --password ${DOCKERHUB_PASSWORD}

.PHONY: docker-build
docker-build: clean ## build docker image with cache
	${DOCKER} build --pull -t ${REGISTRY}/${APPLICATION}:${COMMIT_SHA} .

.PHONY: docker-build-release
docker-build-release: clean ## build docker image without cache (slower than make docker-build)
	${DOCKER} build --pull --no-cache -t ${REGISTRY}/${APPLICATION}:${COMMIT_SHA} .

.PHONY: docker-push
docker-push: docker-build ## push the docker image to registry
	${DOCKER} push ${REGISTRY}/${APPLICATION}:${COMMIT_SHA}

.PHONY: help
help: ## prints this help message
	@echo "Usage: \n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
