package usercontext

import (
    "github.com/garupanojisan/lyra/internal/usercontext/domain/repo"
    "github.com/garupanojisan/lyra/internal/usercontext/infrastructure/persistence/memory"
    "github.com/garupanojisan/lyra/pkg/lyra/di"
)

// default (no build tag): use in-memory repository
func getUserRepo(c *di.Container) repo.UserRepository {
    return memory.NewUsers()
}

