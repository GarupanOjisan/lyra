# Repository Guidelines

## Project Structure & Module Organization
- `cmd/api`, `cmd/worker`: Executable entry points (HTTP API, outbox worker).
- `pkg/lyra/{app,di,httpx,tx,outbox}`: Framework core (modules, DI, HTTP, Tx, Outbox).
- `internal/<bounded-context>/{domain,application,interfaces,infrastructure}`: DDD layers; inward-only deps.
- Example: `internal/usercontext/interfaces/http` registers routes on `httpx.Router`.

## Build, Test, and Development Commands
- `make build`: Compile all packages.
- `make test`: Run `go test ./... -race -cover`.
- `make fmt` / `make vet`: Format and static analysis.
- `make run-api` / `make run-worker`: Launch API server (on :8080) / worker.
- Direct Go: `go build ./...`, `go run ./cmd/api`.

## Coding Style & Naming Conventions
- Go 1.21+, `gofmt -s` required; imports grouped and sorted.
- Packages: short, lowercase, no underscores; tests end with `_test.go`.
- DDD naming: `CreateUserHandler.Handle(ctx, cmd)`, repos per aggregate (e.g., `UserRepository`).
- Errors: wrap with `%w`; HTTP errors via `httpx.Problem` (RFC7807).
- DI: constructor injection; register via module `Boot()` using `di.Provide(c, dep)`.

## Testing Guidelines
- Use Go `testing`; table-driven tests preferred.
- Domain: pure unit tests; Application: in-memory repos + Tx; HTTP: `httptest` + JSON assertions.
- Name tests `TestXxx`; fixtures under `testdata/` if needed.
- Run: `make test` (coverage on) or `go test ./... -race -cover`.

## Commit & Pull Request Guidelines
- Commits: Conventional Commits (e.g., `feat(httpx): add problem+json`).
- PRs: clear description, scope-limited, link issues (`Closes #123`).
- Include tests/docs for behavior changes; describe trade-offs and follow-ups.

## Security & Configuration Tips
- Do not commit secrets; prefer env vars (e.g., `LYRA_HTTP_ADDR`, tokens).
- `.env` and `go.work` are ignored—use locally only. Artifacts/coverage are gitignored.

## Architecture Overview
- Layers: Domain <- Application <- Interfaces/Infrastructure; one use case = one Tx.
- Modules: implement `Boot/Start/Stop` and add with `app.Use(...)` in `cmd/*`.
