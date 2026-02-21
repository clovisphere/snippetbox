# 📇 Snippet-Box

Share your text snippets effortlessly! 😊

> A hands-on project built while working through the excellent [Let's Go](https://lets-go.alexedwards.net/) book by Alex Edwards.

## 🛠️ Prerequisites

Before running `Snippet-Box`, make sure you have the following installed:

- **Go** (version 1.25 or higher) - [installation guide](https://go.dev/dl/)
- **MySQL** (version 8.0.45 or higher) - [releases](https://dev.mysql.com/downloads/mysql/)

Optionally, having [Docker](https://www.docker.com/) installed is useful for running development services like MySQL.

## ⚙️ Installation & Running

> Clone the repository

  ```sh
  git clone https://github.com/clovisphere/snippetbox.git
  cd snippetbox
  ```

> Install dependencies

```sh
go mod download
```

> Setup the database

```sh
# Start MySQL (via Docker)
make start

# Ensure the database exists and apply migrations
make migration-up
```

> Build and run the application

```sh
go run main.go
```

Or just use the [Makefile](./Makefile)

```sh
make help  # Lists all available commands and usage
```

By default, the application will listen on port `4000`. To use a custom port:

```sh
go run main.go addr=":6969"
# or using Makefile
make local PORT=6969
```

### 🗃️ Database Migrations

Snippet-Box uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations. With the [Makefile](./Makefile), you can manage migrations easily:

```sh
# Create a new migration
make migration-create NAME=add_column_phone

# Apply all pending migrations
make migration-up

# Rollback last migration (interactive)
make migration-down

# Show current migration version
make migration-status

# Force database to a specific migration version (interactive)
make migration-force
```

For development, you can run MySQL with Docker using:

```sh
make start   # Starts services
make stop    # Stops services
make restart # Restarts services
make logs    # Tail logs
```

### 💡 Tips

- Use `make clean` to remove build artifacts and coverage reports.
- Keep your migration files in `./migrations` for compatibility with the Makefile targets.
- The project structure follows [Go best practices](https://go.dev/doc/modules/layout#server-project) and is easy to extend with new features.
