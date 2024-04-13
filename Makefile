## help: list available commands
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run/web: run the cmd/web application 
run/web:
	go run ./cmd/web
