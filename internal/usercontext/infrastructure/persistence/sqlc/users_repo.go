//go:build sqlc

package sqlc

import (
    "context"
    "database/sql"
    "time"

    "github.com/garupanojisan/lyra/internal/usercontext/domain/entity"
    "github.com/garupanojisan/lyra/internal/usercontext/domain/repo"
    "github.com/garupanojisan/lyra/pkg/lyra/sqldb"
    "github.com/garupanojisan/lyra/internal/usercontext/infrastructure/persistence/sqlcgen"
)

var _ repo.UserRepository = (*Users)(nil)

type Users struct { DB *sql.DB }

func NewUsers(db *sql.DB) *Users { return &Users{DB: db} }

func (r *Users) Save(ctx context.Context, u *entity.User) error {
    q := r.queries(ctx)
    return q.CreateUser(ctx, sqlcgen.CreateUserParams{
        ID:        string(u.ID),
        Email:     string(u.Email),
        CreatedAt: u.CreatedAt,
    })
}

func (r *Users) FindByEmail(ctx context.Context, e entity.Email) (*entity.User, error) {
    q := r.queries(ctx)
    row, err := q.GetUserByEmail(ctx, string(e))
    if err == sql.ErrNoRows { return nil, nil }
    if err != nil { return nil, err }
    return &entity.User{ID: entity.UserID(row.ID), Email: entity.Email(row.Email), CreatedAt: row.CreatedAt}, nil
}

func (r *Users) queries(ctx context.Context) *sqlcgen.Queries {
    if tx := sqldb.TxFrom(ctx); tx != nil {
        return sqlcgen.New(tx)
    }
    return sqlcgen.New(r.DB)
}

