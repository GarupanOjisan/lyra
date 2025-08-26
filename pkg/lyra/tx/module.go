package tx

import (
    "context"

    "github.com/garupanojisan/lyra/pkg/lyra/di"
)

type Module struct{}

func (m *Module) Name() string { return "tx" }
func (m *Module) Boot(c *di.Container) error {
    di.Provide[Manager](c, InMemoryManager{})
    return nil
}
func (m *Module) Start(ctx context.Context) error { return nil }
func (m *Module) Stop(ctx context.Context) error  { return nil }

func ModuleDefault() *Module { return &Module{} }
