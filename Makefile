APP ?= vessels-api
GO_VERSION ?= 1.18
ENV ?= local

print-% : ; @echo $($*) ## Print value of a variable (e.g. `make print-APP_VERSION`)
.PHONY: print-%

start_postgres:
	docker run --rm -p "5432:5432" --name postgresdb-$(APP) -e POSTGRES_PASSWORD=$${POSTGRES_PASSWORD} -d postgres
.PHONY: start_postgres

stop_postgres:
	docker stop postgresdb-$(APP)
.PHONY: stop_postgres

test:
	go test -v ./...
.PHONY: test

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