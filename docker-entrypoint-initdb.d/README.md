# Docker Setup for ASQLI Testing

This directory contains Docker configuration for running a MySQL test database with sample data.

## Quick Start

### Option 1: Run ASQLI in Docker with Watch Mode (Recommended for Development)

This will start both MySQL and ASQLI in containers with automatic rebuild on code changes:

```bash
# From the project root directory
docker compose watch
```

This will:
- Start MySQL with sample data
- Build and run ASQLI automatically
- Watch for file changes and rebuild instantly
- Connect ASQLI to MySQL automatically

You'll see the ASQLI CLI connected to the MySQL database. Any code changes will trigger an automatic rebuild.

To stop:
```bash
# Press Ctrl+C, then
docker compose down
```

### Option 2: MySQL Only (Run ASQLI Locally)

Start just the MySQL database and connect with local ASQLI:

```bash
# Start MySQL only
docker compose up -d mysql

# Wait for MySQL to be ready (about 10-30 seconds)
docker compose logs -f mysql

# Connect with local ASQLI binary
./asqli --dbtype mysql --host localhost --port 3306 --user testuser --password testpass --db testdb
```

Or using connection string:

```bash
./asqli --dbtype mysql --connection "testuser:testpass@tcp(localhost:3306)/testdb"
```

### Sample Queries to Try

Once connected, try these natural language queries:

- `show me all users`
- `how many products do we have?`
- `list all orders with status shipped`
- `show me the top 5 most expensive products`
- `which users are from USA?`
- `show me orders from user john_doe`
- `what's the average product rating?`
- `list products with low stock (less than 50)`

Or use raw SQL with `#` prefix:

- `# SELECT * FROM users LIMIT 10`
- `# SELECT name, price FROM products WHERE category = 'Electronics'`
- `# SELECT u.username, COUNT(o.id) as order_count FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id`

## Database Schema

### Tables

- **users**: User accounts (10 sample users)
- **products**: Product catalog (15 sample products)
- **orders**: Customer orders (10 sample orders)
- **order_items**: Order line items
- **reviews**: Product reviews (12 sample reviews)

### Credentials

- **Host**: localhost
- **Port**: 3306
- **Database**: testdb
- **Username**: testuser
- **Password**: testpass
- **Root Password**: rootpassword

## Management Commands

```bash
# Run with watch mode (auto rebuild on changes)
docker compose watch

# Start services in background
docker compose up -d

# Stop services
docker compose down

# Stop and remove all data
docker compose down -v

# Restart services
docker compose restart

# View logs
docker compose logs -f

# View MySQL logs only
docker compose logs -f mysql

# View ASQLI logs only
docker compose logs -f asqli

# Access MySQL CLI
docker compose exec mysql mysql -u testuser -ptestpass testdb

# Rebuild ASQLI container
docker compose build asqli

# Run ASQLI interactively
docker compose run --rm asqli
```

## Data Persistence

Data is persisted in a Docker volume named `mysql_data`. To start fresh:

```bash
docker-compose down -v
docker-compose up -d
```

This will recreate the database with the initial sample data.
