APP ?= vessels-api
GO_VERSION ?= 1.18
ENV ?= local
MIGRATIONS_DIR ?= db/migrations

print-% : ; @echo $($*) ## Print value of a variable (e.g. `make print-APP_VERSION`)
.PHONY: print-%

start_postgres:
	docker run --rm -p "5432:5432" --name postgresdb-$(APP) -e POSTGRES_PASSWORD=$${POSTGRES_PASSWORD} -d postgres
.PHONY: start_postgres

stop_postgres:
	docker stop postgresdb-$(APP)
.PHONY: stop_postgres

test:
	go test -v -race ./...
.PHONY: test

test_integration:
	go test -v -race -tags=integration ./...
.PHONY: test_integration

get_deps:
	go get -t -d ./...
.PHONY: get_deps

update_deps:
	go get -u -t -d ./...
	go mod tidy
.PHONY: update_deps

run:
	set -o allexport && . "config/$(ENV).env" && go run ./cmd/vesselsapi/main.go
.PHONY: run

# Migration commands
migrate_version:
	migrate -source file://$(shell pwd)/$(MIGRATIONS_DIR) -database "postgres://postgres:$${POSTGRES_PASSWORD}@0.0.0.0:5432/postgres?sslmode=disable" version
.PHONY: migrate_version

migrate_up:
	migrate -source file://$(shell pwd)/$(MIGRATIONS_DIR) -database "postgres://postgres:$${POSTGRES_PASSWORD}@0.0.0.0:5432/postgres?sslmode=disable" up
.PHONY: migrate_up

migrate_down:
	migrate -source file://$(shell pwd)/$(MIGRATIONS_DIR) -database "postgres://postgres:$${POSTGRES_PASSWORD}@0.0.0.0:5432/postgres?sslmode=disable" down
.PHONY: migrate_down