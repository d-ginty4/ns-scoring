.PHONY : run tables docs

run:
	APP_ENV=local go run main.go

tables:
	bash ./local/createTables.sh

docs:
	swag init -g routes/routes.go
