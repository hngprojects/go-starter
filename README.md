# go-starter

A standard Go web application scaffold — stdlib HTTP, structured logging, graceful shutdown, and a SQL migration workflow.

---

## Stack

- **Go 1.23+**
- **`net/http`** with the Go 1.22+ method-aware `ServeMux` (no third-party router)
- **`log/slog`** — structured JSON logging
- **`golang-migrate`** — SQL migrations
- **Postgres** — assumed default (any `golang-migrate`-supported DB works)

---

## Project layout

```
.
├── cmd/server/main.go        # Entry point — server bootstrap, graceful shutdown
├── internal/
│   ├── config/               # Env loading
│   ├── handlers/             # HTTP handlers (health, etc.)
│   ├── middleware/           # Logger, recover, chain helper
│   ├── models/               # Domain models
│   ├── repository/           # Data access layer
│   ├── service/              # Business logic
│   └── server/               # Router + middleware wiring
├── pkg/logger/               # slog wrapper
├── api/                      # OpenAPI / proto specs
├── migrations/               # SQL migrations (golang-migrate format)
├── scripts/                  # Build / deploy helpers
├── test/                     # Integration tests
├── .env.example
├── go.mod
├── Makefile
└── README.md
```

---

## Prerequisites

| Tool                  | Why                          | Install                                                                     |
| --------------------- | ---------------------------- | --------------------------------------------------------------------------- |
| Go 1.23+ (linux-amd64)| Build & run                  | https://go.dev/dl/                                                          |
| Postgres 14+          | Database (you can swap)      | `apt install postgresql` / Docker / managed                                 |
| `migrate` CLI         | Run migrations               | `make install-tools` (installs `golang-migrate v4.17.1` with postgres tag)  |

> **Note:** if `go test -race` errors with *"-race is not supported on linux/386"*, your Go install is the 32-bit build. Reinstall from `https://go.dev/dl/go<version>.linux-amd64.tar.gz`.

---

## Quick start

```bash
# 1. Configure environment
cp .env.example .env
# edit .env — set DATABASE_URL

# 2. Install dev tools (migrate CLI)
make install-tools

# 3. Run migrations against your DB
export $(grep -v '^#' .env | xargs)   # load .env into the shell
make migrate-up

# 4. Run the server
make run
```

Verify:

```bash
curl http://localhost:8080/health
# {"status":"ok","time":"..."}
```

---

## All commands

Run `make help` to see this list with descriptions:

```bash
make help
```

### Build & run

| Command       | What it does                                |
| ------------- | ------------------------------------------- |
| `make run`    | Run the server (`go run ./cmd/server`)      |
| `make build`  | Build a binary at `bin/server`              |
| `make clean`  | Remove `bin/`, `tmp/`, `dist/`, coverage    |

### Test & quality

| Command            | What it does                                              |
| ------------------ | --------------------------------------------------------- |
| `make test`        | `go test -race -cover ./...`                              |
| `make test-cover`  | Generate `coverage.html` and print the open hint          |
| `make lint`        | `go vet ./...`                                            |
| `make fmt`         | `gofmt -s -w .`                                           |
| `make tidy`        | `go mod tidy`                                             |

### Tooling

| Command               | What it does                                                       |
| --------------------- | ------------------------------------------------------------------ |
| `make install-tools`  | Install `golang-migrate` CLI to `$GOPATH/bin` (or `$GOBIN`)        |

Make sure `$(go env GOPATH)/bin` (default `~/go/bin`) is on your `PATH`.

### Migrations

All migration targets read `DATABASE_URL` from the environment and fail fast with a clear error if it isn't set.

| Command                                       | What it does                                                |
| --------------------------------------------- | ----------------------------------------------------------- |
| `make migrate-up`                             | Apply every pending migration                               |
| `make migrate-down`                           | Roll back the most recent migration (one step)              |
| `make migrate-create NAME=add_users`          | Create a new pair of `<seq>_add_users.{up,down}.sql` files  |
| `make migrate-version`                        | Print the current schema version                            |
| `make migrate-force VERSION=1`                | Force a version (fixes a "dirty" migration state)           |
| `make migrate-drop`                           | **DESTRUCTIVE** — drop everything in the database           |

#### Migration workflow example

```bash
# 1. Create a migration
make migrate-create NAME=add_orders_table
# -> migrations/000002_add_orders_table.up.sql
# -> migrations/000002_add_orders_table.down.sql

# 2. Edit both files (write the CREATE in .up.sql, the DROP in .down.sql)

# 3. Apply it
make migrate-up

# 4. Check the version
make migrate-version
# -> 2

# 5. Roll back if needed
make migrate-down
```

#### Recovering from a dirty migration

If a migration fails halfway, golang-migrate marks the version as *dirty* and refuses to continue. Fix the underlying SQL (or manually clean up the partial state in the DB), then:

```bash
make migrate-force VERSION=<the-version-that-failed>
make migrate-up
```

---

## Configuration

All config is loaded from environment variables in `internal/config/config.go`. See `.env.example`.

| Variable        | Default       | Description                                       |
| --------------- | ------------- | ------------------------------------------------- |
| `APP_ENV`       | `development` | Environment label (`development`, `production`)   |
| `PORT`          | `8080`        | HTTP listen port                                  |
| `LOG_LEVEL`     | `info`        | One of `debug`, `info`, `warn`, `error`           |
| `LOG_FORMAT`    | _(auto)_      | `pretty` or `json`. Defaults to `pretty` when `APP_ENV=development`, else `json` |
| `DATABASE_URL`  | _(unset)_     | Postgres connection string — required for migrations |

---

## Renaming the module

The scaffold ships as `github.com/user/go-starter`. To rename:

```bash
NEW=github.com/your-org/your-app
go mod edit -module $NEW
grep -rl 'github.com/user/go-starter' . | xargs sed -i "s|github.com/user/go-starter|$NEW|g"
go mod tidy
```

---

## Endpoints (out of the box)

| Method | Path      | Response                                                  |
| ------ | --------- | --------------------------------------------------------- |
| GET    | `/`       | `{"message":"go-starter is running"}`                     |
| GET    | `/health` | `{"status":"ok","time":"<RFC3339>"}`                      |

The router uses Go 1.22+ method-aware patterns (`GET /health`), so `POST /health` returns `405`, and unknown paths return `404`.
