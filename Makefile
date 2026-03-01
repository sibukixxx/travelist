.PHONY: dev dev-api dev-frontend build test lint clean

# Development
dev:
	docker compose up --build

dev-api:
	cd api && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev

# Build
build:
	docker compose build

build-api:
	cd api && go build -o ../tmp/server ./cmd/server

build-frontend:
	cd frontend && npm run build

# Test
test: test-api test-frontend

test-api:
	cd api && GOCACHE=/tmp/go-build-cache go test ./...

test-frontend:
	cd frontend && npm test

# Lint
lint: lint-api lint-frontend

lint-api:
	cd api && GOCACHE=/tmp/go-build-cache go vet ./...

lint-frontend:
	cd frontend && npm run lint

# Clean
clean:
	rm -rf tmp/ frontend/dist/
