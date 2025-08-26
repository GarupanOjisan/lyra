package usercontext

import (
    "context"
    "time"

    "github.com/garupanojisan/lyra/internal/usercontext/application/usecase"
    "github.com/garupanojisan/lyra/internal/usercontext/domain/entity"
    "github.com/garupanojisan/lyra/internal/usercontext/domain/repo"
    httpi "github.com/garupanojisan/lyra/internal/usercontext/interfaces/http"
    "github.com/garupanojisan/lyra/pkg/lyra/di"
    "github.com/garupanojisan/lyra/pkg/lyra/httpx"
    "github.com/garupanojisan/lyra/pkg/lyra/tx"
)

type Module struct{}

func (m *Module) Name() string { return "usercontext" }
func (m *Module) Boot(c *di.Container) error {
    var users repo.UserRepository = getUserRepo(c)
    di.Provide[repo.UserRepository](c, users)

    handler := &usecase.CreateUserHandler{
        Users: users,
        Tx:    di.MustGet[tx.Manager](c),
        Now:   time.Now,
        NewID: func() entity.UserID { return entity.UserID(time.Now().UTC().Format("20060102T150405.000000000")) },
    }

    api := &httpi.UsersAPI{Create: handler}
    r := di.MustGet[*httpx.Router](c)
    api.Register(r)
    return nil
}
func (m *Module) Start(ctx context.Context) error { return nil }
func (m *Module) Stop(ctx context.Context) error  { return nil }

func ModuleDefault() *Module { return &Module{} }
