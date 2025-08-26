# Persistence (sqlc) — Lyra

This repository includes scaffolding for a sqlc-based persistence option.

## Layout
- `migrations/` — SQL schema (users table as example)
- `sqlc.yaml` — sqlc config (PostgreSQL engine)
- `internal/usercontext/infrastructure/persistence/sqlc/queries/*.sql` — queries
- Generated code target: `internal/usercontext/infrastructure/persistence/sqlcgen`
- Repo impl (build-tagged): `internal/.../persistence/sqlc/users_repo.go` (`//go:build sqlc`)

## Generate
1) Install sqlc locally and a DB driver (e.g., pgx) in your environment.
2) Run: `sqlc generate`
3) Build with the tag to enable sqlc repo: `go build -tags sqlc ./...`

## Runtime (optional)
- Provide a DSN/driver via env to enable the SQL Tx module:
  - `LYRA_DB_DRIVER=pgx`
  - `LYRA_DB_DSN=postgres://user:pass@localhost:5432/lyra?sslmode=disable`
- Register in `cmd/api` by replacing the Tx module with `sqldb.ModuleFromEnv()` if you want DB at runtime.

Note: By default, the app uses in-memory repositories for easy bootstrapping.

