include .env

PGDB_URI = ${POSTGRES_URI}

.PHONY: serve
serve:
	go run main.go serve

.PHONY: seed
seed:
	go run db/seeders/main.go

.PHONY: migration
migration:
	migrate create -ext sql -dir db/migrations $(name)

.PHONY: migrate
migrate:
	migrate -path db/migrations -database "$(PGDB_URI)?sslmode=disable -verbose up

.PHONY: rollback
rollback:
	migrate -path db/migrations -database "$(PGDB_URI)?sslmode=disable -verbose down 1

.PHONY: rollback-all
rollback-all:
	migrate -path db/migrations -database "$(PGDB_URI)?sslmode=disable -verbose down -all

.PHONY: force-migrate
force-migrate:
	migrate -path db/migrations -database "$(PGDB_URI)?sslmode=disable -verbose force $(version)
