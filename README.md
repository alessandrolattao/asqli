# SQLAI

SQLAI is a command-line tool that generates SQL queries from natural language using AI. It connects to your database to understand the schema and create accurate queries based on your descriptions.

## Features

- Supports PostgreSQL, MySQL, and SQLite databases
- Generates SQL queries from natural language descriptions
- Automatically extracts schema information for context-aware queries
- Interactive CLI interface

## Installation

Download the pre-built binary for your platform from the [latest release](https://github.com/alessandrolattao/sqlai/releases/latest).

### Linux (x86_64)
```bash
tar -xzf sqlai-linux-amd64-X.Y.Z.tar.gz
./sqlai-linux-amd64
```

### Linux (ARM64)
```bash
tar -xzf sqlai-linux-arm64-X.Y.Z.tar.gz
./sqlai-linux-arm64
```

### macOS (Intel)
```bash
tar -xzf sqlai-darwin-amd64-X.Y.Z.tar.gz
./sqlai-darwin-amd64
```

### macOS (Apple Silicon)
```bash
tar -xzf sqlai-darwin-arm64-X.Y.Z.tar.gz
./sqlai-darwin-arm64
```

### Windows
Extract the zip file and run `sqlai-windows-amd64.exe`.

### From Source
If you prefer to build from source:

```bash
go install github.com/alessandrolattao/sqlai/cmd/sqlai@latest
```

## Usage

```bash
# Connect to a PostgreSQL database
sqlai --dbtype postgres --host localhost --port 5432 --user myuser --password mypassword --db mydb

# Connect to a MySQL database
sqlai --dbtype mysql --host localhost --port 3306 --user myuser --password mypassword --db mydb

# Connect to a SQLite database
sqlai --dbtype sqlite --file path/to/database.db

# Use a direct connection string
sqlai --dbtype postgres --connection "postgresql://user:password@localhost:5432/mydb?sslmode=disable"
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