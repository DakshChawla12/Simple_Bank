postgres:
	docker run --name postgres-db --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres

createdb: 
	docker exec -it postgres-db createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres-db dropdb simplebank

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simplebank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simplebank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simplebank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/simplebank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

mock:
	 mockgen -package mockdb -destination db/mock/store.go  github.com/DakshChawla/simplebank/db/sqlc Store

test:
	go test -v ./...
server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1