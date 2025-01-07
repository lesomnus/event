package event

import (
	"context"
	"sync"
	"sync/atomic"
)

type sure[V any] struct {
	slot[V]
}

func Sure[V any]() Slot[V] {
	return &sure[V]{}
}

func (e *sure[V]) Connect(n int) (<-chan V, func()) {
	m := sync.Mutex{}
	d := atomic.Bool{}
	t := e.tc.Add(1)
	c := make(chan V, n)
	s := make(chan interface{})
	a := func(ctx context.Context, v V) {
		m.Lock()
		defer m.Unlock()
		if d.Load() {
			return
		}

		select {
		case <-ctx.Done():
			return
		case c <- v:
			return
		case <-s:
			close(c)
			c = nil
		}
	}

	e.ls.Store(t, a)
	return c, func() {
		if d.Swap(true) {
			return
		}
		close(s)
		e.ls.Delete(t)

		m.Lock()
		defer m.Unlock()
		if c != nil {
			close(c)
		}
	}
}
