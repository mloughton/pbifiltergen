all: build

build:
	@echo "Building..."

	@go build -o main cmd/api/main.go && ./main

run:
	@go run cmd/api/main.go && ./main