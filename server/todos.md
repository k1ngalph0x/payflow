REF:

https://github.com/tsawler/go-microservices/

GIN MIDDLE WARE -
https://gin-gonic.com/en/docs/examples/using-middleware/
https://gin-gonic.com/en/docs/examples/grouping-routes/

AUTH SERVICE:

1.

PACKAGES
go get github.com/joho/godotenv

go get github.com/lib/pq

go get -u github.com/gin-gonic/gin

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate create -ext sql -dir db/migrations -seq create_payflow_auth_table
migrate create -ext sql -dir db/migrations -seq create_payflow_wallets_table

go get -u github.com/golang-jwt/jwt/v5

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

go get google.golang.org/grpc

protoc \ --go_out=. \ --go-grpc_out=. \ proto/wallet.proto

protoc --go_out=. --go-grpc_out=. proto/wallet.proto - run this in wallet-service

make migrate-up # Run migrations
make migrate-down # Rollback last migration
make migrate-create # Create new migration

SELECT \* from schema_migrations;
UPDATE schema_migrations SET dirty = FALSE WHERE version = 2;
