DB_URL=postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable

push: 
	git add .
	git commit -m "$(message)"
	git push -u origin develop

network:
	docker network create bank-network

postgres:
	docker run --name postgres --network bank-network -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres

createdb:
	docker exec -it src-postgres-1 createdb --username=backend --owner=backend backend

dropdb:
	docker exec -it src-postgres-1 dropdb --username=backend backend

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up $(times)

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down $(times)  


migrateversion:
	migrate -path db/migrations -database "$(DB_URL)" version


new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

gen:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

build: 
	go build -o main 

mock:
	mockgen -package mockdb -destination db/mock/store.go backend/db/sqlc Store
	#mockgen -package mockwk -destination worker/mock/distributor.go backend/worker TaskDistributor

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

fproto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto


evans:
	docker run --name evans -d -p 9090:9090 evans

redis:
	docker run --name src-redis-1 -p 6379:6379 -d redis:7-alpine

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redis