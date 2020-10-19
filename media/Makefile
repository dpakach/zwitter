SERVER_OUT := "bin/server"
PKG := "github.com/dpakach/zwitter/media"
SERVER_PKG_BUILD := "${PKG}/server"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
VERSION := 1.0

.PHONY: all api server client

all: server

swagger: api
	@echo "not implemented"

docs: swagger

dep: ## Get the dependencies
	@go get -v -d ./...

server: dep api ## Build the binary file for server
	@go build -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

clean: ## Remove previous builds
	@rm -rf $(SERVER_OUT) $(CLIENT_OUT)

docker-build:
	docker build -t zwiter/media:$(VERSION) --build-arg service_name=media -f ../service/Dockerfile  .. --no-cache

docker-run: docker-build
	docker run --rm -it --name=z-auth --network=host zwitter-auth

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

