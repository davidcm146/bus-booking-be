# Bus Booking API

Backend API for the Bus Booking system, written in Go (Gin + GORM + PostgreSQL).

## Tech stack

- **Language:** Go 1.26
- **Web framework:** Gin
- **ORM:** GORM (over `database/sql` + `lib/pq`)
- **Database:** PostgreSQL 16
- **Migrations:** Goose (SQL-based, in `migrations/`)
- **Auth:** JWT (`golang-jwt/v5`), bcrypt password hashing
- **i18n:** `go-i18n/v2` with embedded `vi` / `en` locales
- **Validation:** `go-playground/validator/v10` (custom `phone`, `strong_password` rules)
- **Hot reload:** Air (dev only)
- **API docs:** Swagger via `swaggo/swag`

## Prerequisites

| Tool | Why | Install |
|------|-----|---------|
| Go 1.26+ | Build & run | https://go.dev/dl/ |
| PostgreSQL 16+ | Database | or use the Docker compose service |
| [Air](https://github.com/air-verse/air) | Hot reload in dev | `go install github.com/air-verse/air@latest` |
| [Goose](https://github.com/pressly/goose) | Run migrations | `go install github.com/pressly/goose/v3/cmd/goose@latest` |
| [golangci-lint](https://golangci-lint.run/) | Lint (optional) | see site |
| [Swag](https://github.com/swaggo/swag) | Generate Swagger docs (optional) | `go install github.com/swaggo/swag/cmd/swag@latest` |
| Docker + Docker Compose | Containerized run | https://docs.docker.com/ |

## Configuration

All config is read from environment variables (with sensible dev defaults). Copy the example and edit:

```bash
cp .env.example .env
```

| Variable | Default | Notes |
|----------|---------|-------|
| `APP_ENV` | `development` | `production`/`staging` switches logger to JSON and **requires** `JWT_SECRET` |
| `APP_PORT` | `8080` | HTTP listen port |
| `DB_HOST` | `localhost` | |
| `DB_PORT` | `5432` | |
| `DB_USER` | `postgres` | |
| `DB_PASSWORD` | `10042003` | **change in any non-local setup** |
| `DB_NAME` | `bus_booking` | database must exist before running migrations |
| `DB_SSLMODE` | `disable` | |
| `DB_TIMEZONE` | `Asia/Ho_Chi_Minh` | |
| `DB_MAX_OPEN_CONNS` | `25` | |
| `DB_MAX_IDLE_CONNS` | `10` | |
| `DB_CONN_MAX_LIFETIME` | `1h` | |
| `DB_CONN_MAX_IDLE_TIME` | `15m` | |
| `DATABASE_URL` | — | lib/pq URL used **only by Goose** (Makefile `migrate-*` targets) |
| `JWT_SECRET` | — | **required in production** |
| `JWT_EXPIRES_IN` | `24h` | Go duration string |

> The `Makefile` does `include .env` and `export`, so every `make` target picks up your `.env` automatically.

## Getting started

### Option A — Docker (recommended)

Spins up PostgreSQL + the API together. The app waits for Postgres to be healthy before starting.

```bash
cp .env.example .env          # tweak DB_PASSWORD / JWT_SECRET
docker compose up -d --build
```

Then run migrations against the containerized DB (from the host, using the `DATABASE_URL` in your `.env`):

```bash
make migrate-up
```

API is at `http://localhost:8080`. Health check: `GET http://localhost:8080/health`.

To stop and remove containers (data volume `bus_booking_db_data` is preserved):

```bash
docker compose down
```

### Option B — Local Go toolchain

1. **Start PostgreSQL** (or point `DB_HOST` at a remote instance) and create the database:
   ```bash
   createdb bus_booking
   ```

2. **Install dev tools** (Air, Goose):
   ```bash
   go install github.com/air-verse/air@latest
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

3. **Run migrations:**
   ```bash
   make migrate-up
   ```

4. **Run the app** (hot reload):
   ```bash
   make run      # uses Air, rebuilds on file changes
   ```
   …or build & run a binary:
   ```bash
   make build && ./bin/app
   ```

## Database migrations

Migrations live in `migrations/` as Goose SQL files. The initial schema (`20260712094943_init_schema.sql`) creates all tables: `roles`, `users`, `operators`, `administrative_areas`, `locations`, `routes`, `route_stops`, `vehicles`, `trips`, `bookings`, `seat_inventories`, `payments`, `passenger_checkins`, `promotions`, `booking_promotions`, `notifications`, `reviews`, `audit_logs` — plus `updated_at` triggers and trigram/GIN indexes for i18n (`vi`/`en`) name search.

```bash
make migrate-up                 # apply all pending
make migrate-down               # roll back one
make migrate-up-to version=20260712094943
make migrate-down-to version=20260712094943
make migrate-status
make migrate-create name=add_foo   # scaffold a new migration
```

## Seeding

The seed is a goose migration (`migrations/*_seed_data.sql`) that populates roles, 5 operators, 15 cities, 30 routes, 10 vehicles, ~1680 trips across 7 days, ~60k seat inventories, and two demo users. It is idempotent (ON CONFLICT, DELETE+INSERT).

```bash
make migrate-up    # applies schema + seed in one go
make seed          # re-applies just the seed migration (goose redo)
```

Demo users (password `password123` for both):

| Role | Phone | Email |
|------|-------|-------|
| Operator | `+84900000010` | `operator@busgo.app` |
| Passenger | `+84900000020` | `passenger@busgo.app` |

## Testing

```bash
make test      # go test ./...
```

### Conventions

- **`testify/require`** (not `assert`) for fatal assertions — fail fast.
- **Table-driven subtests** with `t.Run("case name", ...)` when there are multiple cases.
- **Service and handler tests** use testify mocks (`internal/mocks/`); never hit the DB.
- **Repository tests** hit the real Postgres — start Docker, run migrations, and seed first:
  ```bash
  docker compose up -d db
  make migrate-up
  make seed
  ```
- **Handler tests** use `httptest.NewRecorder` + `gin.TestMode`; i18n and validator are initialized once in `TestMain`.
- **Auth handler tests** send `multipart/form-data` (the backend uses form binding, not JSON).
- Test files sit next to the file under test as `<file>_test.go`.
- Tests are independent — no shared mutable state between subtests.
- Use `mock.MatchedBy` to assert on struct fields, `mock.Anything` otherwise.

### Layer responsibilities

| Layer | Tests against | Mocks |
|-------|---------------|-------|
| `utils` | pure functions | none |
| `service` | mock repository | `mocks.MockUserRepository` |
| `handler` | mock service | `mocks.MockAuthService`, `mocks.MockUserService` |
| `middleware` | stub auth context | none |
| `repository` | live Postgres | none |
| `shared/i18n` | embedded locale files | none |
| `shared/auth` | context helpers | none |
| `shared/apperror` | error constructors | none |

## API endpoints

Base path: `/api/v1`. All auth endpoints use `multipart/form-data`.

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/health` | — | Liveness probe |
| `POST` | `/api/v1/auth/signup` | — | Register (`phone`, `password`) → returns user + JWT |
| `POST` | `/api/v1/auth/login` | — | Login (`phone`, `password`) → returns user + JWT |
| `GET` | `/api/v1/auth/me` | JWT | Current user profile |
| `POST` | `/api/v1/users` | — | Create user (`name`, `email`, `password`, `phone?`, `role?`) |
| `GET` | `/api/v1/users` | — | List users |
| `GET` | `/api/v1/users/:id` | — | Get user by ID |
| `PUT` | `/api/v1/users/:id` | — | Update user |
| `DELETE` | `/api/v1/users/:id` | — | Delete user |

Protected routes expect `Authorization: Bearer <jwt>`.

### Quick signup / login

```bash
# Signup (multipart form)
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -F phone=+84901234567 \
  -F password='Str0ng!Pass'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -F phone=+84901234567 \
  -F password='Str0ng!Pass'

# Me (authenticated)
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <token>"
```

## Swagger / API docs

Swagger annotations live in the handlers. Regenerate the docs bundle:

```bash
make swag     # swag init -g cmd/app/main.go -o docs
```

Generated `docs/` is gitignored by default (see `.gitignore`).

## Project structure

```
cmd/app/main.go          # Entrypoint: loads .env, boots Application, graceful shutdown
configs/                 # Config structs + env loaders (app, database, jwt)
internal/
  app/                   # Application bootstrap + DI container
  router/                # Gin engine + route registration
  handler/               # HTTP handlers (auth, user)
  service/               # Business logic
  repository/            # GORM data access
  module/                # Per-domain wiring (user, auth)
  middleware/            # Locale + JWT auth middleware
  dto/                   # Request / response DTOs
  model/                 # GORM domain models
  infra/                 # database (postgres, gorm), logger
  shared/                # i18n, validator, response, apperror, auth, constant
  utils/                 # password hashing, JWT generation
migrations/              # Goose SQL migrations
```

## Makefile targets

| Target | Action |
|--------|--------|
| `make run` | Hot-reload dev server via Air |
| `make build` | Compile `./bin/app` |
| `make test` | `go test ./...` |
| `make fmt` | `go fmt ./...` |
| `make lint` | `golangci-lint run` |
| `make swag` | Regenerate Swagger docs |
| `make migrate-*` | Goose migration commands (see above) |
| `make seed` | Re-apply the seed migration via `goose redo` |

## Notes

- The app **does not auto-run migrations** on startup — run `make migrate-up` explicitly after the DB is up.
- `DB_PASSWORD` default (`10042003`) is the original author's local value; override it everywhere except a throwaway local DB.
- i18n locales (`vi`, `en`) are embedded into the binary via `go:embed`; no external locale files needed at runtime.
