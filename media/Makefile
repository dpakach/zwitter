SERVER_OUT := "bin/server"
PKG := "github.com/dpakach/zwitter/media"
SERVER_PKG_BUILD := "${PKG}/server"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
VERSION := 1.0

.PHONY: all api server

all: server

swagger: api
	@echo "not implemented"

docs: swagger ## Generate swagger documentation files

dep: ## download go dependencies
	@go get -v -d ./...

server: dep api ## Build the server binary
	@go build -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

clean: ## clean all build products
	@rm -rf $(SERVER_OUT)

docker-build: ## Build the docker image for the service
	docker build -t zwiter/media:$(VERSION) --build-arg service_name=media -f ../service/Dockerfile  .. --no-cache

docker-run: docker-build ## Build and run the service in a docker container
	docker run --rm -it --name=z-auth --network=host zwitter-auth

help: ## Display this help screen
	@echo
	@echo Zwitter:Auth $(VERSION)
	@echo ------------------------------------------
	@echo
	@echo Available commands
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

