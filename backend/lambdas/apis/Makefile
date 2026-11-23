.PHONY : run tables docs connect seed

run:
	LOG_LEVEL=debug APP_ENV=local go run main.go

tables:
	bash ./local/createTables.sh

docs:
	swag init -g routes/routes.go

connect:
	psql -h localhost -p 5432 -U local -d ns-scoring

seed:
	cd seeding && tsx main.ts
