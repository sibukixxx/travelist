.PHONY: dev dev-api dev-frontend build test lint clean setup

# Dependencies
frontend/node_modules: frontend/package-lock.json frontend/package.json
	cd frontend && npm ci
	@touch $@

# Setup
setup:
	cd frontend && npm ci

# Development
dev:
	docker compose up --build

dev-api:
	cd api && go run ./cmd/server

dev-frontend: frontend/node_modules
	cd frontend && npm run dev

# Build
build:
	docker compose build

build-api:
	cd api && go build -o ../tmp/server ./cmd/server

build-frontend: frontend/node_modules
	cd frontend && npm run build

# Test
test: test-api test-frontend

test-api:
	cd api && GOCACHE=/tmp/go-build-cache go test ./...

test-frontend: frontend/node_modules
	cd frontend && npm test

# Lint
lint: lint-api lint-frontend

lint-api:
	cd api && GOCACHE=/tmp/go-build-cache go vet ./...

lint-frontend: frontend/node_modules
	cd frontend && npm run lint

# Clean
clean:
	rm -rf tmp/ frontend/dist/
