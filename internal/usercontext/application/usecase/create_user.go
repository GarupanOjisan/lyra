package usecase

import (
    "context"
    "errors"
    "time"

    "github.com/garupanojisan/lyra/internal/usercontext/domain/entity"
    "github.com/garupanojisan/lyra/internal/usercontext/domain/repo"
    "github.com/garupanojisan/lyra/pkg/lyra/tx"
)

var ErrEmailTaken = errors.New("email taken")

type CreateUser struct{ Email string }

type CreateUserHandler struct {
    Users repo.UserRepository
    Tx    tx.Manager
    Now   func() time.Time
    NewID func() entity.UserID
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUser) error {
    return h.Tx.WithinTx(ctx, func(ctx context.Context) error {
        if cmd.Email == "" { return ErrEmailTaken }
        email := entity.Email(cmd.Email)
        if existing, _ := h.Users.FindByEmail(ctx, email); existing != nil {
            return ErrEmailTaken
        }
        u := &entity.User{ID: h.NewID(), Email: email, CreatedAt: h.Now()}
        return h.Users.Save(ctx, u)
    })
}

