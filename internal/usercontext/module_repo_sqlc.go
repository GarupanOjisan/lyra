//go:build sqlc

package usercontext

import (
    "database/sql"

    "github.com/garupanojisan/lyra/internal/usercontext/domain/repo"
    "github.com/garupanojisan/lyra/internal/usercontext/infrastructure/persistence/sqlc"
    "github.com/garupanojisan/lyra/internal/usercontext/infrastructure/persistence/memory"
    "github.com/garupanojisan/lyra/pkg/lyra/di"
)

// sqlc build: prefer SQL-backed repo if DB is provided; fallback to memory.
func getUserRepo(c *di.Container) repo.UserRepository {
    if db, ok := di.Get[*sql.DB](c); ok && db != nil {
        return sqlc.NewUsers(db)
    }
    return memory.NewUsers()
}

