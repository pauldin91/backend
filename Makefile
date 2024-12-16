DB_URL=postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up $(times)

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down $(times)  


migrateversion:
	migrate -path db/migrations -database "$(DB_URL)" version


new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

docs:
	dbdocs build doc/db.dbml

schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/pauldin91/backend/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/pauldin91/backend/worker TaskDistributor

gen:
	sqlc generate

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=backend \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

test:
	go test -v -cover -short ./...

build: 
	mkdir -p bin
	go build -o bin/main 

clean:
	rm -rf bin

.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration db_docs db_schema sqlc test server mock proto evans redis