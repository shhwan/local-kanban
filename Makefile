.PHONY: up down build build-backend build-frontend test test-backend test-frontend vet fmt clean generate

up:
	docker compose up -d --build

down:
	docker compose down

build: build-backend build-frontend

build-backend:
	cd backend && go build -o bin/server .

build-frontend: generate
	cd frontend && go build -o bin/server .

generate:
	cd frontend && ~/go/bin/templ generate

test: test-backend test-frontend

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && go test ./...

vet:
	cd backend && go vet ./...
	cd frontend && go vet ./...

fmt:
	cd backend && gofmt -w .
	cd frontend && gofmt -w .

clean:
	cd backend && rm -rf bin/
	cd frontend && rm -rf bin/
