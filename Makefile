#!/usr/bin/make
# Makefile readme (ru): <http://linux.yaroslavl.ru/docs/prog/gnu_make_3-79_russian_manual.html>
# Makefile readme (en): <https://www.gnu.org/software/make/manual/html_node/index.html#SEC_Contents>

SHELL = /bin/sh
LDFLAGS = "-s -w"

DOCKER_BIN = $(shell command -v docker 2> /dev/null)
DC_BIN = $(shell command -v docker-compose 2> /dev/null)
DC_RUN_ARGS = --rm --user "$(shell id -u):$(shell id -g)" addreality
DC_RUN_TEST_ARGS = --rm --user "$(shell id -u):$(shell id -g)" addreality_test
APP_NAME = $(notdir $(CURDIR))
GO_RUN_ARGS ?=

.PHONY : help build fmt lint gotest test cover run shell image clean
.DEFAULT_GOAL : help
.SILENT : test shell

# This will output the help for each task. thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Show this help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[32m%-11s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build app binary file
	$(DC_BIN) run $(DC_RUN_ARGS) go build -ldflags=$(LDFLAGS) -o './$(APP_NAME)' cmd/main.go

fmt: ## Run source code formatter tools
	$(DC_BIN) run $(DC_RUN_ARGS) sh -c 'GO111MODULE=off go get golang.org/x/tools/cmd/goimports && $$GOPATH/bin/goimports -d -w .'
	$(DC_BIN) run $(DC_RUN_ARGS) gofmt -s -w -d .

lint: ## Run app linters
	$(DOCKER_BIN) run --rm -t -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run -v

gotest: ## Run app tests
	$(DC_BIN) -f docker-compose-test.yml run $(DC_RUN_TEST_ARGS) go test -v -race ./...

test: lint gotest ## Run app tests and linters
	@printf "\n   \e[30;42m %s \033[0m\n\n" 'All tests passed!';

cover: ## Run app tests with coverage report
	$(DC_BIN) run $(DC_RUN_ARGS) sh -c 'go test -race -covermode=atomic -coverprofile /tmp/cp.out ./... && go tool cover -html=/tmp/cp.out -o ./coverage.html'
	-sensible-browser ./coverage.html && sleep 2 && rm -f ./coverage.html

run: ## Run app without building binary file
	$(DC_BIN) run $(DC_RUN_ARGS) go run . $(GO_RUN_ARGS)

shell: ## Start shell into container with golang
	$(DC_BIN) run $(DC_RUN_ARGS) bash


image: ## Build docker image with app
	$(DOCKER_BIN) build -f ./Dockerfile -t $(APP_NAME) .
	@printf "\n   \e[1;45m %s \033[0m\n\n" "Now you can run \`docker run $(APP_NAME) <your_args>\`";

clean: ## Make clean
	$(DC_BIN) down -v -t 1
	$(DOCKER_BIN) rmi $(APP_NAME) -f
