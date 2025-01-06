package event

import (
	"context"
	"sync/atomic"
)

type Emitter[K comparable, V any] interface {
	Emit(ctx context.Context, k K, v V)
	WithKey(k K) KeyedEmitter[V]
}

type KeyedEmitter[V any] interface {
	Emit(ctx context.Context, v V)
}

type Listener[K comparable, V any] interface {
	Listen(k K, n int) (<-chan V, func())
}

type Event[K comparable, V any] interface {
	Emitter[K, V]
	Listener[K, V]
}

type handler[V any] func(ctx context.Context, v V)

type event[K comparable, V any] struct {
	slots   gm[K, *gm[int64, handler[V]]]
	tickets atomic.Int64
}

func (e *event[K, V]) Emit(ctx context.Context, k K, v V) {
	ls, ok := e.slots.Load(k)
	if !ok {
		return
	}
	ls.Range(func(_ int64, l handler[V]) bool {
		l(ctx, v)
		return true
	})
}

func (e *event[K, V]) WithKey(k K) KeyedEmitter[V] {
	return WithKey(e, k)
}

func (e *event[K, V]) load(k K) *gm[int64, handler[V]] {
	ls, ok := e.slots.Load(k)
	if ok {
		return ls
	}

	ls_ := &gm[int64, handler[V]]{}
	ls, ok = e.slots.LoadOrStore(k, ls_)
	if ok {
		return ls
	}
	return ls_
}

type keyed[K comparable, V any] struct {
	e Emitter[K, V]
	k K
}

func WithKey[K comparable, V any](e Emitter[K, V], k K) KeyedEmitter[V] {
	return &keyed[K, V]{e, k}
}

func (e *keyed[K, V]) Emit(ctx context.Context, v V) {
	e.e.Emit(ctx, e.k, v)
}
