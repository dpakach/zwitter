SERVER_OUT := "bin/server"
PKG := "github.com/dpakach/zwitter/web"
SERVER_PKG_BUILD := "${PKG}/server"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
VERSION := 1.0

.PHONY: all server

all: server

dep: ## download the go dependencies
	@go get -v -d ./...

server: dep ## Build the server binary
	@go build -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

clean: ## Clean all the build products
	@rm -rf $(SERVER_OUT)

docker-build: ## build the docker images for the service
	docker build -t zwiter/web:$(VERSION) --build-arg service_name=web -f ../service/Dockerfile  .. --no-cache

docker-run: docker-build ## Build and run the service in a docker container
	docker run --rm -it --name=z-web --network=host zwitter-web

help: ## Display this help screen
	@echo
	@echo Zwitter:Users $(VERSION)
	@echo ------------------------------------------
	@echo
	@echo Available make commands
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

