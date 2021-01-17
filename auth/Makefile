SERVER_OUT := "bin/server"
CLIENT_OUT := "bin/client"
API_OUT := "api/authpb/api.pb.go"
API_REST_OUT := "api/authpb/api.pb.gw.go"
PKG := "github.com/dpakach/zwitter/auth"
SERVER_PKG_BUILD := "${PKG}/server"
CLIENT_PKG_BUILD := "${PKG}/client"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
VERSION := 1.0
GOPATH ?= '${HOME}/go'

.PHONY: all api server client

all: server

api/authpb/api.pb.go: api/auth_api.proto
	@protoc -I api/ \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=api \
		--go-grpc_out=api \
		api/auth_api.proto

api/authpb/api.pb.gw.go: api/auth_api.proto
	@protoc -I api/ \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:api \
		api/auth_api.proto

swagger: api
	@protoc -I api/ \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--openapiv2_out=logtostderr=true:./swagger \
		api/auth_api.proto

docs: swagger ## Generate swagger documentation files

api: api/authpb/api.pb.go api/authpb/api.pb.gw.go ## Auto-generate grpc go sources

dep: ## Download go dependencies
	@go get -v -d ./...

server: dep ## Build the server binary
	@go build -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

client: dep ## Build the client binary
	@go build -v -o $(CLIENT_OUT) $(CLIENT_PKG_BUILD)

clean: ## Clean all build products
	@rm -rf $(SERVER_OUT) $(CLIENT_OUT)

docker-build: ## Build the docker image of the service
	docker build -t zwiter/auth:$(VERSION) --build-arg service_name=auth -f ../service/Dockerfile  .. --no-cache

docker-run: docker-build ## Build and run the service in a docker container
	docker run --rm -it --name=z-auth --network=host zwitter-auth

help: ## List available make commands
	@echo
	@echo Zwitter:Auth $(VERSION)
	@echo ------------------------------------------
	@echo
	@echo Available commands
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

