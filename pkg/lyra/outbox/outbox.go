package outbox

import (
    "context"
)

type Event struct {
    Name string
    Data []byte
}

type Outbox interface {
    Enqueue(ctx context.Context, e Event) error
}

type InMemory struct{ ch chan Event }

func NewInMemory(buffer int) *InMemory { return &InMemory{ch: make(chan Event, buffer)} }

func (o *InMemory) Enqueue(_ context.Context, e Event) error { o.ch <- e; return nil }

type Worker struct {
    In  *InMemory
    Run func(context.Context, Event) error
}

func (w *Worker) Start(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case e := <-w.In.ch:
            if w.Run != nil {
                _ = w.Run(ctx, e)
            }
        }
    }
}

