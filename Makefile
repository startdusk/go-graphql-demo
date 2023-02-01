include .env

.PHONY: mock
mock:
	@mockery --all --keeptree

.PHONY: test
test: clean
	@go test ./...

.PHONY: clean
clean:
	@go mod tidy
	@go vet ./...
	@go fmt ./...

MIGRATE_SOURCE=file://data/postgres/migrations

.PHONY: migrate
migrate:
	@migrate -source ${MIGRATE_SOURCE} \
			-database ${DATABASE_URL} up

.PHONY: rollback
rollback:
	@migrate -source ${MIGRATE_SOURCE} \
			-database ${DATABASE_URL} down 1


.PHONY: drop
drop:
	@migrate -source ${MIGRATE_SOURCE} \
			-database ${DATABASE_URL} drop

.PHONY: migration
migration:
	@read -p "Enter migration name: " name; \
		migrate create -ext sql -dir data/postgres/migrations $$name

DATABASE_NAME=go-graphql-demo

.PHONY: initdb
initdb:
	@docker run \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=password \
    -e POSTGRES_DB=${DATABASE_NAME} \
    -p 5432:5432 \
    -d \
    --name=${DATABASE_NAME} \
    postgres -N 1000

run:
	@go run ./cmd/twitter/...

generate:
	@go get github.com/vektah/dataloaden && go get github.com/99designs/gqlgen && go generate ./...
	

