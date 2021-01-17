SERVER_OUT := "bin/server"
CLIENT_OUT := "bin/client"
API_OUT := "api/postspb/api.pb.go"
API_REST_OUT := "api/postspb/api.pb.gw.go"
PKG := "github.com/dpakach/zwitter/posts"
SERVER_PKG_BUILD := "${PKG}/server"
CLIENT_PKG_BUILD := "${PKG}/client"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
VERSION := 1.0
GOPATH ?= '${HOME}/go'

.PHONY: all api server client

all: server

api/postspb/api.pb.go: api/post_api.proto
	@protoc -I api/ \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=api \
		--go-grpc_out=api \
		api/post_api.proto

api/postspb/api.pb.gw.go: api/post_api.proto
	@protoc -I api/ \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:api \
		api/post_api.proto

swagger: api
	@protoc -I api/ \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--openapiv2_out=logtostderr=true:./swagger \
		api/post_api.proto

docs: swagger ## Generate the swagger documentation files

api: api/postspb/api.pb.go api/postspb/api.pb.gw.go ## Auto-generate grpc go sources

dep: ## Download the go dependencies
	@go get -v -d ./...

server: dep ## Build the server binary
	@go build -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

client: dep ## Build the client binary
	@go build -v -o $(CLIENT_OUT) $(CLIENT_PKG_BUILD)

clean: ## Clean all build products
	@rm -rf $(SERVER_OUT) $(CLIENT_OUT)

docker-build: ## Build the docker images for the service
	docker build -t zwiter/posts:$(VERSION) --build-arg service_name=posts -f ../service/Dockerfile  .. --no-cache

docker-run: docker-build ## Build and run the service in a docker container
	docker run --rm -it --name=z-posts --network=host zwiter/posts

help: ## Display this help screen
	@echo
	@echo Zwitter:Posts $(VERSION)
	@echo ------------------------------------------
	@echo
	@echo Available commands
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

