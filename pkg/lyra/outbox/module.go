package outbox

import (
    "context"
    "log"

    "github.com/garupanojisan/lyra/pkg/lyra/di"
)

type Module struct {
    box    *InMemory
    worker *Worker
    cancel context.CancelFunc
}

func (m *Module) Name() string { return "outbox" }

func (m *Module) Boot(c *di.Container) error {
    m.box = NewInMemory(128)
    di.Provide[Outbox](c, m.box)
    m.worker = &Worker{In: m.box, Run: func(ctx context.Context, e Event) error {
        log.Printf("outbox delivered: %s len=%d", e.Name, len(e.Data))
        return nil
    }}
    return nil
}

func (m *Module) Start(ctx context.Context) error {
    ctx, m.cancel = context.WithCancel(ctx)
    go m.worker.Start(ctx)
    return nil
}

func (m *Module) Stop(ctx context.Context) error {
    if m.cancel != nil { m.cancel() }
    return nil
}

func ModuleDefault() *Module { return &Module{} }
