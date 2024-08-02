bin_dir = bin/ecomm

db_entry_dir = cmd/db
migrations_dir = $(db_entry_dir)/migrations

up_arg = -up
down_arg = -down

build: 
	@go build -race -gcflags "-m" -o $(bin_dir) $(entryPointPath)

run: build
	@./$(bin_dir)

migration:
	@$(GOPATH)/bin/migrate create -ext sql -dir $(migrations_dir) $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run $(db_entry_dir)/main.go $(up_arg)

migrate-down:
	@go run $(db_entry_dir)/main.go $(down_arg)

test: 
	@go test -v ./tests/...