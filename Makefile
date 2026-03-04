.PHONY: up down build test vet fmt clean

up:
	docker compose up -d --build

down:
	docker compose down

build:
	cd backend && go build -o bin/server .

test:
	cd backend && go test ./...

vet:
	cd backend && go vet ./...

fmt:
	cd backend && gofmt -w .

clean:
	cd backend && rm -rf bin/
