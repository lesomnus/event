package event

import (
	"context"
	"sync"
	"sync/atomic"
)

type skip[K comparable, V any] struct {
	event[K, V]
}

func SKip[K comparable, V any]() Event[K, V] {
	return &skip[K, V]{}
}

func (e *skip[K, V]) Listen(k K, n int) (<-chan V, func()) {
	ls := e.load(k)

	m := sync.Mutex{}
	d := atomic.Bool{}
	t := e.tickets.Add(1)
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

	ls.Store(t, a)
	return c, func() {
		if d.Swap(true) {
			return
		}
		ls.Delete(t)

		m.Lock()
		defer m.Unlock()
		close(c)
	}
}
