db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

server-air:
	air

server-air-debug: # https://github.com/cosmtrek/air?tab=readme-ov-file#debug
	air -d

tidy:
	go mod tidy

build:
	go mod tidy
	go build main.go

.PHONY: db_docs db_schema sqlc test server