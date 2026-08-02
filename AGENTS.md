# Bus Booking API — Agent Guidelines

## Testing

Run all tests:

```bash
make test        # go test ./...
```

Repository tests require a live PostgreSQL. Start it via Docker and run migrations first:

```bash
docker compose up -d db
make migrate-up
make seed
```

### Conventions

- Use `testify/require` (not `assert`) for fatal assertions — fail fast.
- Table-driven subtests with `t.Run("case name", ...)` when there are multiple cases.
- Service and handler tests use testify mocks (`internal/mocks/`); never hit the DB in those layers.
- Repository tests hit the real Postgres (following the `bun` dsn pattern); use a separate schema or the seeded data.
- Handler tests use `httptest.NewRecorder` + `gin.TestMode`; initialize i18n and validator once in `TestMain`.
- Auth handler tests send `multipart/form-data` (the backend uses form binding, not JSON).
- Test file naming: `<file>_test.go` next to the file under test.
- Keep tests independent — no shared mutable state between subtests.
- Use `mock.MatchedBy` when you need to assert on struct fields, `mock.Anything` otherwise.

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

## Code style

- `gofmt` clean (enforced).
- `go vet` clean.
- Compact code — collapse duplicate else branches, avoid unnecessary nesting.
- Don't add/remove comments unless asked.
- Follow existing error handling style: wrap with `apperror.*` constructors, don't over-try/catch.

## Database

- Migrations: `migrations/` via goose (`make migrate-up`).
- Seed: `migrations/*_seed_data.sql` goose migration, applied via `make migrate-up` or re-applied via `make seed` (goose redo).
- 18 tables; only `users` and `roles` have service/handler code so far.
