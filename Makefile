mainPath = cmd/api/main.go

run:
	@go run $(mainPath)

build-escape-analysis:
	@go build -gcflags "-m" $(mainPath)