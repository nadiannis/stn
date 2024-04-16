include .env

DB_DSN = ${DB_USER}:${DB_PASS}@${DB_PROTO}\(${DB_HOST}:${DB_PORT}\)/${DB_NAME}?parseTime=true

## help: list available commands
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/web: run the cmd/web application 
.PHONY: run/web
run/web:
	go run ./cmd/web -db-dsn=${DB_DSN}

## db/mysql: connect to the database using mysql
.PHONY: db/mysql
db/mysql:
	mysql -D ${DB_NAME} -u ${DB_USER} -p

## db/migrations/new name=$1: create new database migration files 
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Apply all migrations...'
	migrate -path ./migrations -database mysql://${DB_DSN} up

## db/migrations/down: revert all migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Rollback migrations...'
	migrate -path ./migrations -database mysql://${DB_DSN} down

## audit: tidy dependencies, format code, & vet code
.PHONY: audit
audit:
	@echo 'Tidying & verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...