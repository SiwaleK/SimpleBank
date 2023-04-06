postgres: 
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.2-alpine
createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres15 dropdb simple_bank
connectdb:
	docker exec -it postgres15 psql -U root simple_bank
migrateup : 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown :
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc: 
	docker run --rm -v "$(pwd):/src" -w /src kjconroy/sqlc:1.17.0 /workspace/sqlc generate
sqlc_init:
	docker run --rm -v "$(pwd):/src" -w /src kjconroy/sqlc:1.17.0 /workspace/sqlc init
test:
	go test -v -cover ./...
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
