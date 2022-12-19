include $(PWD)/.env

DB_URL = mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)
TEST_DB_URL = mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_TEST_NAME)

run:
	cd cmd; go run main.go

build:
	cd cmd; go build

build_run:
	cd cmd; go build; ./cmd

run_test:
	go test -v -cover ./repository/...; go test -v -cover ./service/...; go test -v -cover ./controller/...

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1

migrateup_test:
	migrate -path db/migrations -database "$(TEST_DB_URL)" -verbose up

migrateup1_test:
	migrate -path db/migrations -database "$(TEST_DB_URL)" -verbose up 1

migratedown_test:
	migrate -path db/migrations -database "$(TEST_DB_URL)" -verbose down

migratedown1_test:
	migrate -path db/migrations -database "$(TEST_DB_URL)" -verbose down 1
