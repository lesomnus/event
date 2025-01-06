package event

import (
	"context"
	"sync"
	"sync/atomic"
)

type sure[K comparable, V any] struct {
	event[K, V]
}

func Sure[K comparable, V any]() Event[K, V] {
	return &sure[K, V]{}
}

func (e *sure[K, V]) Listen(k K, n int) (<-chan V, func()) {
	ls := e.load(k)

	m := sync.Mutex{}
	d := atomic.Bool{}
	t := e.tickets.Add(1)
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

	ls.Store(t, a)
	return c, func() {
		if d.Swap(true) {
			return
		}
		close(s)
		ls.Delete(t)

		m.Lock()
		defer m.Unlock()
		if c != nil {
			close(c)
		}
	}
}
