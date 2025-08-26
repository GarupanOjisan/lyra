package httpx

import (
    "context"

    "github.com/garupanojisan/lyra/pkg/lyra/app"
    "github.com/garupanojisan/lyra/pkg/lyra/di"
)

type Module struct {
    addr   string
    router *Router
    server *Server
}

func (m *Module) Name() string { return "http" }

func (m *Module) Boot(c *di.Container) error {
    if m.addr == "" { m.addr = ":8080" }
    m.router = NewRouter()
    m.router.Use(Recover(), Logger())
    di.Provide(c, m.router)
    m.server = NewServer(m.addr, m.router)
    di.Provide(c, m.server)
    return nil
}

func (m *Module) Start(ctx context.Context) error { return m.server.Start(ctx) }
func (m *Module) Stop(ctx context.Context) error  { return m.server.Stop(ctx) }

func ModuleWithAddr(addr string) app.Module { return &Module{addr: addr} }
func ModuleDefault() app.Module           { return &Module{addr: ":8080"} }
