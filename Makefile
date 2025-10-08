DB_NAME=backand_GO
DB_PASSWORD=1
DB_PORT=5436
DB_HOST=localhost
APP_PORT=8000
MIGRATIONS_PATH=./migrations
DATABASE_URL=postgres://postgres:$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: help db-start db-stop db-clean migrate-up migrate-down run build clean

help: ## Show this help
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

db-start: ## Start PostgreSQL database container
	@echo "Starting PostgreSQL database..."
	docker run --name=todo-db \
	 -e POSTGRES_PASSWORD='$(DB_PASSWORD)' \
	 -e POSTGRES_DB='$(DB_NAME)' \
	 -p $(DB_PORT):5432 -d --rm postgres
	@sleep 3
	@echo "Database started on port $(DB_PORT)"

db-stop: ## Stop PostgreSQL database container
	@echo "Stopping database..."
	docker rm -f todo-db || true

db-clean: db-stop ## Stop and remove database container
	@echo "Database container cleaned"

db-restart: db-stop db-start ## Restart database container

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down

migrate-force: ## Force migration to specific version (usage: make migrate-force version=1)
	@echo "Forcing migration to version $(version)..."
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(version)

run: ## Run the application
	@echo "Starting application..."
	go run cmd/main.go

build: ## Build the application
	@echo "Building application..."
	go build -o bin/main cmd/main.go

dev: db-start migrate-up run ## Start database, run migrations and start application

clean: db-stop ## Clean up (stop database)
	@echo "Cleaning up..."
	rm -rf bin/

test: ## Run tests
	@echo "Running tests..."
	go test ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy

# Database connection check
db-check: ## Check database connection
	@echo "Checking database connection..."
	@pg_isready -h $(DB_HOST) -p $(DB_PORT) -d $(DB_NAME) -U postgres || echo "Database is not ready"

# Show database status
db-status: ## Show database status
	@docker ps -f name=todo-db

# Reset everything: clean database and run migrations
reset: db-clean db-start migrate-up ## Complete reset: clean DB, start fresh and migrate
	@echo "Reset complete!"

# Default target
default: help