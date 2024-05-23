PWD=c:/Wieee/Work/simple-bank
postgres:
	docker run --name postgres --network bank-networkfi -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
	#migrate -path db/migration -database "postgresql://root:oxDqG9zltIZr9MGr5dS6@simple-bank.c5us88q8g0ov.ap-southeast-1.rds.amazonaws.com:5432/simple_bank" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq <migration_name>

sqlc:
	docker run --rm -v $(PWD):/src -w /src sqlc/sqlc generate

test:
	go test -cover -v ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simple-bank/db/sqlc Store

proto:
#	Generate only gRPC
#	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

protoall:
#	rm -f pb/*.go
#	rm -f doc/swagger/*.swagger.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto

evans:
	evans --host localhost --port 8081 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2.5-alpine

.PHONY: postgres created migrateup proto
