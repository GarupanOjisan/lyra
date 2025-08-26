package memory

import (
    "context"
    "errors"
    "sync"

    "github.com/garupanojisan/lyra/internal/usercontext/domain/entity"
    "github.com/garupanojisan/lyra/internal/usercontext/domain/repo"
)

type Users struct {
    mu     sync.RWMutex
    byID   map[entity.UserID]*entity.User
    byMail map[entity.Email]*entity.User
}

func NewUsers() *Users {
    return &Users{byID: make(map[entity.UserID]*entity.User), byMail: make(map[entity.Email]*entity.User)}
}

var _ repo.UserRepository = (*Users)(nil)

func (r *Users) Save(_ context.Context, u *entity.User) error {
    r.mu.Lock(); defer r.mu.Unlock()
    if u == nil { return errors.New("nil user") }
    r.byID[u.ID] = u
    r.byMail[u.Email] = u
    return nil
}

func (r *Users) FindByEmail(_ context.Context, e entity.Email) (*entity.User, error) {
    r.mu.RLock(); defer r.mu.RUnlock()
    if u, ok := r.byMail[e]; ok { return u, nil }
    return nil, nil
}

