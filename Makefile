LDFLAGS := -ldflags "-X main.Version=${VERSION}"

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## run the api server
	go run ${LDFLAGS} cmd/api/main.go

.PHONY: build
build:  ## build the API server binary
		CGO_ENABLED=0 go build ${LDFLAGS} -a -o server $(MODULE)/cmd/api

.PHONY: docker-up
docker-up: ## start production environment
	docker compose --env-file .env -f deployments/docker-compose.yml up -d
.PHONY: docker-dev-up
docker-dev-up: ## start development environment
	docker compose --env-file .env -f deployments/docker-compose.yml -f deployments/docker-compose.dev.yml up -d
.PHONY: docker-down
docker-down: ## stop all environments
	docker compose --env-file .env -f deployments/docker-compose.yml down
.PHONY: docker-dev-down
docker-dev-down: ## stop development environment
	docker compose --env-file .env -f deployments/docker-compose.yml -f deployments/docker-compose.dev.yml down

.PHONY: docker-logs
docker-logs: ## show logs: make docker-logs SERVICE=api
	docker compose --env-file .env -f deployments/docker-compose.yml logs -f $(SERVICE)
.PHONY: docker-dev-logs
docker-dev-logs: ## show dev logs: make docker-dev-logs SERVICE=api
	docker compose --env-file .env -f deployments/docker-compose.yml -f deployments/docker-compose.dev.yml logs -f $(SERVICE)

.PHONY: docs-watch
	docs-watch: ## watch docs changes and rebuild

docs-watch:
	$(MAKE) -j2 typespec-watch scalar-preview

typespec-watch:
	cd docs/typespec && tsp compile . --watch

scalar-preview:
	cd docs && npx @scalar/cli project preview
