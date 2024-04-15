## help: list available commands
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run/web: run the cmd/web application 
.PHONY: run/web
run/web:
	go run ./cmd/web

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