.PHONY: all auth posts users

## set SERVICE env variable for service specific commands

SERVICE?=

all: build ## Build the binaries and docker-images for all the services

service: ## build the binary and docker image of $SERVICE service
	cd $(SERVICE) && \
	make && make docker-build

restart-service: ## Restart a service $SERVICE
	make -C $(SERVICE) api
	make -C $(SERVICE)
	docker-compose -f ./docker-compose.envoy.yml restart $(SERVICE)

build: frontend-build ## Build the binaries and docker-images for all the services
	cd auth && \
	make && make docker-build && \
	cd ../posts && \
	make && make docker-build && \
	cd ../users && \
	make && make docker-build && \
	cd ../media && \
	make && make docker-build && \
	cd ../web && \
	make && make docker-build

clean: ## Clean the build products for all the services
	cd auth && \
	make clean && \
	cd ../posts && \
	make clean && \
	cd ../users && \
	make clean && \
	cd ../media && \
	make clean && \
	cd ../web && \
	make clean

config: ## Create default configuraion files
	cd auth && \
	cp config/config.example.yaml config/config.yaml && \
	cd ../posts && \
	cp config/config.example.yaml config/config.yaml && \
	cd ../users && \
	cp config/config.example.yaml config/config.yaml && \
	cd ../media && \
	cp config/config.example.yaml config/config.yaml && \
	cd ../web && \
	cp config/config.example.yaml config/config.yaml

docker-run: ## Run zwitter with vanilla docker-compose setup
	docker-compose up

envoy-run: ## Run zwitter using docker-compose with envoy proxy
	docker-compose -f ./docker-compose.envoy.yml up -d

js-deps: ## Install javascript dependencies
	cd frontend && \
	yarn

frontend-build: js-deps ## Start frontend dev environment
	cd frontend && \
	yarn build

frontend-watch: js-deps ## Start frontend dev environment
	cd frontend && \
	yarn start

logs: ## Show docker compose logs for service $SERVICE, shows all logs if $SERVICE is not set
	docker-compose -f ./docker-compose.envoy.yml logs --tail=250 -ft $(SERVICE)

stop: ## Stop the services
	docker-compose -f ./docker-compose.envoy.yml down

dev: envoy-run frontend-watch ## Start the development environment

run: envoy-run

fmt: ## Format the go code with gofmt
	go fmt ./...

help: ## Display this help screen
	@echo
	@echo Zwitter
	@echo ------------------------------------------
	@echo
	@echo Available commands
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

