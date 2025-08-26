package di

import (
    "reflect"
    "sync"
)

type Container struct {
    mu       sync.RWMutex
    services map[reflect.Type]any
}

func New() *Container { return &Container{services: make(map[reflect.Type]any)} }

func Provide[T any](c *Container, svc T) {
    c.mu.Lock()
    defer c.mu.Unlock()
    t := reflect.TypeOf((*T)(nil)).Elem()
    c.services[t] = svc
}

func Get[T any](c *Container) (T, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    var zero T
    t := reflect.TypeOf((*T)(nil)).Elem()
    v, ok := c.services[t]
    if !ok {
        return zero, false
    }
    return v.(T), true
}

func MustGet[T any](c *Container) T {
    v, ok := Get[T](c)
    if !ok {
        panic("di: missing dependency")
    }
    return v
}
