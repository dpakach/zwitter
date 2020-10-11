.PHONY: all auth posts users

SERVICE?=

service:
	cd $(SERVICE) && \
	make && make docker-build

restart-service:
	docker-compose -f ./docker-compose.envoy.yml restart $(SERVICE)

build-all:
	cd auth && \
	make && make docker-build && \
	cd ../posts && \
	make && make docker-build && \
	cd ../users && \
	make && make docker-build

docker-run:
	docker-compose up

envoy-run:
	docker-compose -f ./docker-compose.envoy.yml up -d

frontend-watch:
	cd frontend && \
	yarn start

logs:
	docker-compose -f ./docker-compose.envoy.yml logs --tail=250 -ft $(SERVICE)

stop:
	docker-compose -f ./docker-compose.envoy.yml down

dev: envoy-run frontend-watch

fmt:
	go fmt ./...

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

