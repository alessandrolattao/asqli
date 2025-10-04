.PHONY: up down mysql postgres clean logs help
.DEFAULT_GOAL := help

# Variables
PROVIDER ?= openai
MODEL ?=

# Start all databases
up:
	@echo "Starting MySQL and PostgreSQL..."
	@docker compose up -d

# Stop all services
down:
	@docker compose down

# Connect to MySQL
mysql:
	@echo "Starting MySQL..."
	@docker compose up -d mysql --wait
	@echo "MySQL is ready! Connecting with provider: $(PROVIDER)$(if $(MODEL), model: $(MODEL),)..."
	@echo ""
	@go run ./... --provider $(PROVIDER) $(if $(MODEL),--model $(MODEL),) --dbtype mysql --host localhost --port 3306 --user testuser --password testpass --db testdb

# Connect to PostgreSQL
postgres:
	@echo "Starting PostgreSQL..."
	@docker compose up -d postgres --wait
	@echo "PostgreSQL is ready! Connecting with provider: $(PROVIDER)$(if $(MODEL), model: $(MODEL),)..."
	@echo ""
	@go run ./... --provider $(PROVIDER) $(if $(MODEL),--model $(MODEL),) --dbtype postgres --host localhost --port 5432 --user testuser --password testpass --db testdb

# Clean all data (remove volumes)
clean:
	@echo "Stopping services and removing volumes..."
	@docker compose down -v

# Show logs
logs:
	@docker compose logs -f

# Help
help:
	@echo "ASQLI Development Makefile"
	@echo ""
	@echo "Commands:"
	@echo "  make up        - Start MySQL and PostgreSQL databases"
	@echo "  make mysql     - Start MySQL and connect with ASQLI"
	@echo "  make postgres  - Start PostgreSQL and connect with ASQLI"
	@echo "  make down      - Stop all services"
	@echo "  make clean     - Stop services and remove all data"
	@echo "  make logs      - Show database logs"
	@echo ""
	@echo "AI Provider Options:"
	@echo "  PROVIDER=<name>  - AI provider (default: openai)"
	@echo "                     Options: openai, gemini"
	@echo "  MODEL=<name>     - AI model to use (optional)"
	@echo ""
	@echo "Examples:"
	@echo "  make postgres                           # Use default OpenAI provider"
	@echo "  make postgres PROVIDER=gemini           # Use Gemini with default model"
	@echo "  make mysql PROVIDER=openai MODEL=gpt-4  # Use OpenAI with GPT-4"
	@echo ""
	@echo "Database credentials:"
	@echo "  User:     testuser"
	@echo "  Password: testpass"
	@echo "  Database: testdb"
	@echo ""
	@echo "Ports:"
	@echo "  MySQL:      localhost:3306"
	@echo "  PostgreSQL: localhost:5432"
