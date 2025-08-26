package repo

import (
    "context"
    "github.com/garupanojisan/lyra/internal/usercontext/domain/entity"
)

type UserRepository interface {
    Save(ctx context.Context, u *entity.User) error
    FindByEmail(ctx context.Context, e entity.Email) (*entity.User, error)
}

