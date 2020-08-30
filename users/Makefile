SERVER_OUT := "bin/server"
CLIENT_OUT := "bin/client"
API_OUT := "api/userspb/api.pb.go"
API_REST_OUT := "api/userspb/api.pb.gw.go"
PKG := "github.com/dpakach/zwitter/users"
SERVER_PKG_BUILD := "${PKG}/server"
CLIENT_PKG_BUILD := "${PKG}/client"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all api server client

all: server client

api/userspb/api.pb.go: api/users_api.proto
	@protoc -I api/ \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:api \
		api/users_api.proto

api/userspb/api.pb.gw.go: api/users_api.proto
	@protoc -I api/ \
		-I/usr/local/include \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:api \
		api/users_api.proto

api: api/userspb/api.pb.go api/userspb/api.pb.gw.go ## Auto-generate grpc go sources

dep: ## Get the dependencies
	@go get -v -d ./...

server: dep api ## Build the binary file for server
	@go build -i -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

client: dep api ## Build the binary file for client
	@go build -i -v -o $(CLIENT_OUT) $(CLIENT_PKG_BUILD)

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(CLIENT_OUT) $(API_OUT) $(API_REST_OUT)

docker-build: server
	docker build -t zwitter-users . --no-cache

docker-run: docker-build
	docker run --rm -it --name=z-users --network=host zwitter-users

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

