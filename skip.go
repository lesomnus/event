package event

import (
	"context"
	"sync"
	"sync/atomic"
)

type skip[V any] struct {
	slot[V]
}

func Skip[V any]() Slot[V] {
	return &skip[V]{}
}

func (e *skip[V]) Connect(n int) (<-chan V, func()) {
	m := sync.Mutex{}
	d := atomic.Bool{}
	t := e.tc.Add(1)
	c := make(chan V, n)
	a := func(ctx context.Context, v V) {
		m.Lock()
		defer m.Unlock()
		if d.Load() {
			return
		}

		select {
		case c <- v:
		default:
		}
	}

	e.ls.Store(t, a)
	return c, func() {
		if d.Swap(true) {
			return
		}
		e.ls.Delete(t)

		m.Lock()
		defer m.Unlock()
		close(c)
	}
}
