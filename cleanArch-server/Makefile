createmigration:
	migrate create -ext=sql -dir=internal/infra/sql/migrations -seq init

migrate:
	migrate -path=internal/infra/sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up

migratedown:
	migrate -path=internal/infra/sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose down

.PHONY: migrate migratedown createmigration