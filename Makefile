BINARY      := ./bin/api
BUILD_FLAGS := -trimpath -ldflags="-s -w"
PKG         := golabs-api/...

.PHONY: build run dev test lint tidy docker-up docker-down migrate

## build: Compile the API binary
build:
	go build $(BUILD_FLAGS) -o $(BINARY) ./cmd/api

## run: Build then run locally
run: build
	$(BINARY)

## dev: Hot-reload with 'air' (go install github.com/air-verse/air@latest)
dev:
	air

## test: Run all tests
test:
	go test ./...

## lint: Run static analysis tools
lint:
	go vet ./...
	@command -v staticcheck >/dev/null 2>&1 && staticcheck ./... || echo "Install staticcheck: go install honnef.co/go/tools/cmd/staticcheck@latest"

## tidy: Tidy and verify modules
tidy:
	go mod tidy
	go mod verify

## docker-up: Build image and start all containers
docker-up:
	docker compose up --build -d

## docker-down: Stop all containers
docker-down:
	docker compose down

## migrate: Print migration files to apply (manual target — use a tool for production)
migrate:
	@echo "Applying migrations from deployments/database/init/ ..."
	@for f in deployments/database/init/V*.sql; do \
		echo "  -> $$f"; \
	done
	@echo "Run these against your DB with: mysql -u<user> -p<pass> <db> < <file>"
