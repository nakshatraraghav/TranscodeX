start:
	@go run cmd/main.go

build:
	@go build -o bin/api cmd/main.go

# Run PostgreSQL container
run_pg:
	@docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres

# Stop PostgreSQL container
stop_pg:
	@docker container stop postgres

# Remove PostgreSQL container
remove_pg_container:
	@docker container rm postgres

# Create database
createdb:
	@docker exec -it postgres createdb --username=postgres --owner=postgres transcodex

# Drop database
dropdb:
	@docker exec -it postgres dropdb --username=postgres transcodex

# Migrate Up
migrateup:
	@docker run --rm -v $(PWD)/db/migration:/migrations --network host migrate/migrate -path=/migrations -database "postgresql://postgres:password@127.0.0.1:5432/transcodex?sslmode=disable" -verbose up

# Migrate Down
migratedown:
	@docker run --rm -v $(PWD)/db/migration:/migrations --network host migrate/migrate -path=/migrations -database "postgresql://postgres:password@127.0.0.1:5432/transcodex?sslmode=disable" -verbose down -all

# Clean up all (stop and remove container)
clean:
	@make stop_postgres
	@make remove_postgres_container
