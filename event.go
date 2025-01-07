package event

import (
	"context"

	"github.com/lesomnus/event/internal/gm"
)

type Emitter[K comparable, V any] interface {
	Emit(ctx context.Context, k K, v V)
}

type Listener[K comparable, V any] interface {
	Listen(k K, n int) (<-chan V, func())
}

type Event[K comparable, V any] interface {
	Emitter[K, V]
	Listener[K, V]
}

type event[K comparable, V any] struct {
	ss   gm.Map[K, Slot[V]]
	init func() Slot[V]
}

func New[K comparable, V any](init func() Slot[V]) Event[K, V] {
	return &event[K, V]{init: Sure[V]}
}

func (e *event[K, V]) Emit(ctx context.Context, k K, v V) {
	s, ok := e.ss.Load(k)
	if !ok {
		return
	}

	s.Signal(ctx, v)
}

func (e *event[K, V]) Listen(k K, n int) (<-chan V, func()) {
	s := e.WithKey(k)
	return s.Connect(n)
}

func (e *event[K, V]) WithKey(k K) Slot[V] {
	s, ok := e.ss.Load(k)
	if ok {
		return s
	}

	s_ := e.init()
	s, ok = e.ss.LoadOrStore(k, s_)
	if ok {
		return s
	}
	return s_
}

type keyed[K comparable, V any] struct {
	e Emitter[K, V]
	k K
}

func WithKey[K comparable, V any](e Emitter[K, V], k K) Signaler[V] {
	if v, ok := e.(interface{ WithKey(k K) Slot[V] }); ok {
		return v.WithKey(k)
	}

	return &keyed[K, V]{e, k}
}

func (e *keyed[K, V]) Signal(ctx context.Context, v V) {
	e.e.Emit(ctx, e.k, v)
}
