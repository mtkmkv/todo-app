-include .env
export

export PROJECT_ROOT := $(shell pwd)

.PHONY: env-up env-down env-cleanup \
        migrate-create migrate-up migrate-down migrate-action \
        todoapp-run

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Clean up all environment volume files? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ] || [ "$$ans" = "Y" ]; then \
		docker compose down todoapp-postgres && \
		rm -rf out/pgdata && \
		echo "Environment files cleaned up."; \
	else \
		echo "Cleanup canceled."; \
	fi

env-port-forward:
	@docker compose up -d todoapp-port-forwarder

env-port-close:
	@docker compose down todoapp-port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Error: Please provide a migration name via 'seq' variable."; \
		echo "Example: make migrate-create seq=init_database"; \
		exit 1; \
	fi; \
	MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make --no-print-directory migrate-action action=up

migrate-down:
	@make --no-print-directory migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Error: Please provide an action (up or down) via 'action' variable."; \
		echo "Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable \
		"$(action)"

todoapp-run:
	@go run cmd/todoapp/main.go