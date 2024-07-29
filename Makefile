entryPointPath = cmd/api/main.go
binPath = bin/ecomm

build: 
	@go build -race -gcflags "-m" -o $(binPath) $(entryPointPath)

run: build
	@./$(binPath)

test: 
	@go test -v ./tests/...