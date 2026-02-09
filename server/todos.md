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

go get github.com/google/uuid

github.com/rabbitmq/amqp091-go

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

## wallet -> identity -> merchant -> payment

{
"email": "tester@gmail.com",
"password": "12345678",
"role": "user"
}

{
"message": "User created successfully",
"role": "user",
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNWI2ZWVjNTUtN2Y4ZC00YjY3LWFmN2QtYTQyYzVkODFkNzE3IiwiZW1haWwiOiJ0ZXN0ZXJAZ21haWwuY29tIiwicm9sZSI6InVzZXIiLCJpc3MiOiJwYXlmbG93LWF1dGgiLCJzdWIiOiI1YjZlZWM1NS03ZjhkLTRiNjctYWY3ZC1hNDJjNWQ4MWQ3MTciLCJleHAiOjE3NzA3MTE3MTQsImlhdCI6MTc3MDYyNTMxNH0.Sw2QeTT7RqcAoqj618ede5-nELMoYrrV7Sh3olzxxt8"
}

MERCHANT
{
"email": "merchant@test.com",
"password": "password123",
"role": "merchant"
}

{
"message": "User created successfully",
"role": "merchant",
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMGVjYzg5YzMtZWIxZC00OTk3LWI3NjEtNjE1MjIwNzgzMTFjIiwiZW1haWwiOiJtZXJjaGFudEB0ZXN0LmNvbSIsInJvbGUiOiJtZXJjaGFudCIsImlzcyI6InBheWZsb3ctYXV0aCIsInN1YiI6IjBlY2M4OWMzLWViMWQtNDk5Ny1iNzYxLTYxNTIyMDc4MzExYyIsImV4cCI6MTc3MDcxMTg5NSwiaWF0IjoxNzcwNjI1NDk1fQ.NArdCeUeSCubt_fvApQl4bwSRmHnIK-TgQjl0x-aMbo"
}

merchant usr id auth - 0ecc89c3-eb1d-4997-b761-61522078311c

{
"merchant_id": "7ccacc97-b9e5-4251-a131-0414f7c1294a",
"status": "ACTIVE"
}

DROP TABLE IF EXISTS payflow_merchants CASCADE;
DROP TABLE IF EXISTS schema_migrations;

INTERNAL:
email - platform@payflow.internal
role - platform
password - internal
id : "646b4308-a1d0-4407-85a5-8208d0bb67e0"

RABBIT MQ
docker run -d --hostname payflow-rabbit \
 --name rabbitmq \
 -p 5672:5672 \
 -p 15672:15672 \
 rabbitmq:3-management

http://localhost:15672
guest / guest

docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
