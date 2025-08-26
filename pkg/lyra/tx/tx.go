package tx

import "context"

type Manager interface {
    WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type InMemoryManager struct{}

func (InMemoryManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
    return fn(ctx)
}

