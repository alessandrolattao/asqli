.PHONY: up down mysql postgres clean logs help
.DEFAULT_GOAL := help

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
	@echo "MySQL is ready! Connecting..."
	@echo ""
	@go run ./... --dbtype mysql --host localhost --port 3306 --user testuser --password testpass --db testdb

# Connect to PostgreSQL
postgres:
	@echo "Starting PostgreSQL..."
	@docker compose up -d postgres --wait
	@echo "PostgreSQL is ready! Connecting..."
	@echo ""
	@go run ./... --dbtype postgres --host localhost --port 5432 --user testuser --password testpass --db testdb

# Clean all data (remove volumes)
clean:
	@echo "Stopping services and removing volumes..."
	@docker compose down -v

# Show logs
logs:
	@docker compose logs -f

# Help
help:
	@echo "SQLAI Development Makefile"
	@echo ""
	@echo "Commands:"
	@echo "  make up        - Start MySQL and PostgreSQL databases"
	@echo "  make mysql     - Start MySQL and connect with SQLAI"
	@echo "  make postgres  - Start PostgreSQL and connect with SQLAI"
	@echo "  make down      - Stop all services"
	@echo "  make clean     - Stop services and remove all data"
	@echo "  make logs      - Show database logs"
	@echo ""
	@echo "Database credentials:"
	@echo "  User:     testuser"
	@echo "  Password: testpass"
	@echo "  Database: testdb"
	@echo ""
	@echo "Ports:"
	@echo "  MySQL:      localhost:3306"
	@echo "  PostgreSQL: localhost:5432"
