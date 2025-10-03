# SQLAI

SQLAI is a command-line tool that generates SQL queries from natural language using AI. It connects to your database to understand the schema and create accurate queries based on your descriptions.

https://github.com/user-attachments/assets/d5ac054f-57b7-421d-b29f-0178f327ca95

## Features

- 🗄️ **Database Support**: PostgreSQL, MySQL, and SQLite
- 🤖 **Multiple AI Providers**: OpenAI, Claude (Anthropic), Google Gemini, and Ollama (local models)
- 💬 **Natural Language to SQL**: Generate queries from plain English descriptions
- 🔍 **Schema-Aware**: Automatically extracts database schema for accurate queries
- 🎨 **Interactive TUI**: Beautiful terminal interface with table navigation
- ⚡ **Fast & Efficient**: Token usage tracking and caching support
- 🔧 **Raw SQL Mode**: Execute direct SQL with `#` prefix
- 📊 **Query History**: Navigate and reuse previous queries
- 🔐 **Password Management**: Support for PostgreSQL `.pgpass` file

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

## Prerequisites

SQLAI requires an API key from one of the supported AI providers (except Ollama which runs locally):

### OpenAI (default)
```bash
export OPENAI_API_KEY="sk-..."
```
Get your API key from: https://platform.openai.com/api-keys

### Claude (Anthropic)
```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```
Get your API key from: https://console.anthropic.com/

### Google Gemini
```bash
export GEMINI_API_KEY="AIza..."
```
Get your API key from: https://aistudio.google.com/app/apikey

### Ollama (Local - No API Key Required)
Install Ollama and pull a model:
```bash
# Install Ollama from https://ollama.ai
ollama pull llama3.2
# Or any other model: llama3.3, deepseek-r1, qwen2.5, etc.
```

Ollama will automatically detect and use:
- Running models (priority)
- Locally available models (fallback)

## Usage

```bash
# Connect to a PostgreSQL database (using OpenAI by default)
sqlai --dbtype postgres --host localhost --port 5432 --user myuser --password mypassword --db mydb

# Use Claude (Anthropic)
sqlai --provider claude --dbtype postgres --host localhost --port 5432 --user myuser --password mypassword --db mydb

# Use a specific Claude model
sqlai --provider claude --model claude-opus-4-1 --dbtype postgres --connection "postgresql://..."

# Use Google Gemini
sqlai --provider gemini --dbtype postgres --host localhost --port 5432 --user myuser --password mypassword --db mydb

# Use Ollama (local models - no API key needed)
sqlai --provider ollama --dbtype postgres --host localhost --port 5432 --user myuser --password mypassword --db mydb

# Use Ollama with a specific model
sqlai --provider ollama --model llama3.3 --dbtype sqlite --file database.db

# Use Ollama on a remote server
export OLLAMA_HOST=http://192.168.1.100:11434
sqlai --provider ollama --dbtype postgres --connection "postgresql://..."

# Connect to a MySQL database
sqlai --dbtype mysql --host localhost --port 3306 --user myuser --password mypassword --db mydb

# Connect to a SQLite database
sqlai --dbtype sqlite --file path/to/database.db
```

### Parameters

#### AI Provider
| Parameter    | Description                                          | Default   |
|--------------|------------------------------------------------------|-----------|
| `--provider` | AI provider (openai, claude, gemini, ollama)         | openai    |
| `--model`    | AI model to use (provider-specific, optional)        |           |

#### Database Connection
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

#### Other
| Parameter    | Description                                   | Default   |
|--------------|-----------------------------------------------|-----------|
| `--version`  | Print the version and exit                    | |

## Password Management with `.pgpass`

SQLAI supports the PostgreSQL `.pgpass` file for secure password storage. This allows you to omit the `--password` flag from the command line.

### Setup

Create a `.pgpass` file in your home directory:

**Linux/macOS:**
```bash
# Create the file
cat > ~/.pgpass << EOF
hostname:port:database:username:password
localhost:5432:mydb:myuser:mypassword
localhost:5432:*:postgres:adminpass
*:*:testdb:testuser:testpass
EOF

# Set correct permissions (required)
chmod 0600 ~/.pgpass
```

**Windows:**
```
# Location: %APPDATA%\postgresql\pgpass.conf
hostname:port:database:username:password
```

### Format

```
hostname:port:database:username:password
```

- Use `*` as wildcard in any of the first four fields
- Escape `:` or `\` characters with `\`
- First matching line is used
- Comments start with `#`

### Usage

```bash
# Password automatically loaded from .pgpass
sqlai --dbtype postgres --host localhost --user myuser --db mydb

# Override with PGPASSFILE environment variable
export PGPASSFILE=/path/to/custom/pgpass
sqlai --dbtype postgres --host localhost --user myuser --db mydb
```

**Note:** `.pgpass` works for both PostgreSQL and MySQL connections.

## Interactive Usage

Once connected, SQLAI provides a beautiful terminal interface:

### Natural Language Queries
```
sqlai > show me all users from Italy
sqlai > count active subscriptions by plan
sqlai > list top 10 customers by revenue
```

### Raw SQL Mode (prefix with `#`)
```
sqlai > # SELECT * FROM users WHERE created_at > NOW() - INTERVAL '7 days'
```

### Keyboard Shortcuts
- `↑`/`↓`/`←`/`→` - Navigate table results
- `Ctrl+↑`/`Ctrl+↓` - Navigate query history
- `Ctrl+r` - Open history list
- `Ctrl+p` - View last query details (prompt, SQL, tokens)
- `Ctrl+c` - Copy table as TSV
- `Esc` - Clear input
- `Ctrl+q` - Quit

### AI Provider Models

**OpenAI**: `gpt-4o`, `gpt-4-turbo`, `gpt-3.5-turbo`, `o1-mini`, `o3-mini`

**Claude**: `claude-sonnet-4-5` (default), `claude-opus-4-1`, `claude-3-7-sonnet-latest`, `claude-3-5-haiku-latest`

**Gemini**: `gemini-2.0-flash-exp` (default), `gemini-1.5-pro`, `gemini-1.5-flash`

**Ollama**: Auto-detects running/available models. Popular: `llama3.2`, `llama3.3`, `deepseek-r1`, `qwen2.5`, `mistral`, `codellama`

## License

MIT

## Warning

This software is in early development and not ready for production use.
