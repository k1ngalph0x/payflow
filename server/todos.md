REF:

https://github.com/tsawler/go-microservices/

GIN MIDDLE WARE -
https://gin-gonic.com/en/docs/examples/using-middleware/
https://gin-gonic.com/en/docs/examples/grouping-routes/

AUTH SERVICE:

1.

PACKAGES
github.com/joho/godotenv

github.com/lib/pq

go get -u github.com/gin-gonic/gin

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

grate create -ext sql -dir db/migrations -seq create_payflow_auth_table

go get -u github.com/golang-jwt/jwt/v5

make migrate-up # Run migrations
make migrate-down # Rollback last migration
make migrate-create # Create new migration
