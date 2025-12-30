postgres:
	docker run --name postgres-db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres

createdb: 
	docker exec -it postgres-db createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres-db dropdb simplebank

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simplebank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test