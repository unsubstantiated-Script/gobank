postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb --username=root simple_bank

migrateup:
	 migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlcversion:
	docker run --rm -v $(CURDIR):/src -w /src kjconroy/sqlc version

sqlcupgrade:
	docker run --rm -v $(CURDIR):/src -w /src kjconroy/sqlc upgrade

sqlcinit:
	docker run --rm -v $(CURDIR):/src -w /src kjconroy/sqlc init

sqlcgenerate:
	docker run --rm -v $(CURDIR):/src -w //src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlcversion sqlcupgrade sqlcinit sqlcgenerate test server