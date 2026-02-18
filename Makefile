migrate: ## Run migrations against non test DB
	docker-compose up -d migrate

migrate-create: ## Create a DB migration. You need to pass the file name, e.g. `make migrate-create name=migration-name`
	docker-compose run --rm migrate create -ext sql -dir /migrations $(name)