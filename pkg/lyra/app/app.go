package app

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "sync"
    "syscall"

    "github.com/garupanojisan/lyra/pkg/lyra/di"
)

type Module interface {
    Name() string
    Boot(c *di.Container) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}

type App struct {
    c       *di.Container
    modules []Module
    mu      sync.Mutex
}

func New() *App { return &App{c: di.New()} }

func (a *App) Container() *di.Container { return a.c }

func (a *App) Use(m Module) { a.modules = append(a.modules, m) }

func (a *App) Boot() error {
    for _, m := range a.modules {
        if err := m.Boot(a.c); err != nil {
            return fmt.Errorf("boot %s: %w", m.Name(), err)
        }
    }
    return nil
}

func (a *App) Start(ctx context.Context) error {
    for _, m := range a.modules {
        if err := m.Start(ctx); err != nil {
            return fmt.Errorf("start %s: %w", m.Name(), err)
        }
    }
    return nil
}

func (a *App) Stop(ctx context.Context) error {
    for i := len(a.modules) - 1; i >= 0; i-- {
        m := a.modules[i]
        if err := m.Stop(ctx); err != nil {
            return fmt.Errorf("stop %s: %w", m.Name(), err)
        }
    }
    return nil
}

func (a *App) Run(ctx context.Context) error {
    if err := a.Boot(); err != nil { return err }
    if err := a.Start(ctx); err != nil { return err }

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
    return a.Stop(context.Background())
}

