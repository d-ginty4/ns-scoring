.PHONY : run tables docs

run:
	LOG_LEVEL=debug APP_ENV=local go run main.go

tables:
	bash ./local/createTables.sh

docs:
	swag init -g routes/routes.go
