.PHONY: build run test clean docker-build docker-up docker-down

# Build variables
BINARY_NAME=sense

build:
	go build -o bin/$(BINARY_NAME) cmd/main.go

run:
	go run cmd/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate:
	go run scripts/migrate.go
