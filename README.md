# SQLAI

SQLAI is a command-line tool that generates SQL queries from natural language using AI. It connects to your database to understand the schema and create accurate queries based on your descriptions.

## Features

- Supports PostgreSQL, MySQL, and SQLite databases
- Generates SQL queries from natural language descriptions
- Automatically extracts schema information for context-aware queries
- Interactive CLI interface

## Installation

```bash
# Clone the repository
git clone https://github.com/alessandrolattao/sqlai.git

# Build the project
cd sqlai
go build ./cmd/sqlai
```

## Usage

```bash
# Connect to a PostgreSQL database
./sqlai --dbtype postgres --host localhost --port 5432 --user myuser --password mypassword --db mydb

# Connect to a MySQL database
./sqlai --dbtype mysql --host localhost --port 3306 --user myuser --password mypassword --db mydb

# Connect to a SQLite database
./sqlai --dbtype sqlite --file path/to/database.db

# Use a direct connection string
./sqlai --dbtype postgres --connection "postgresql://user:password@localhost:5432/mydb?sslmode=disable"
```

### Connection Parameters

| Parameter    | Description                                   | Default   |
|--------------|-----------------------------------------------|-----------|
| `--dbtype`   | Database type (postgres, mysql, sqlite)       | postgres  |
| `--connection` | Full connection string (overrides other params) | |
| `--host`     | Database host                                 | |
| `--port`     | Database port                                 | 5432      |
| `--user`     | Database username                             | |
| `--password` | Database password                             | |
| `--db`       | Database name                                 | |
| `--sslmode`  | PostgreSQL SSL mode                           | disable   |
| `--parsetime` | MySQL: parse time values to Go time.Time     | true      |
| `--file`     | SQLite database file path                     | |
| `--version`  | Print the version and exit                    | |

## Example

Once connected to your database, you can ask for queries in natural language:

```
> Show me all users who registered in the last month
```

The tool will generate and display the corresponding SQL query based on your database schema.

## License

MIT

## Warning

This software is in early development and not ready for production use.