package event

import (
	"context"
	"sync/atomic"

	"github.com/lesomnus/event/internal/gm"
)

type Signaler[V any] interface {
	Signal(ctx context.Context, v V)
}

type Connector[V any] interface {
	Connect(n int) (<-chan V, func())
}

type Slot[V any] interface {
	Signaler[V]
	Connector[V]
}

type slot[V any] struct {
	ls gm.Map[int64, func(ctx context.Context, v V)]
	tc atomic.Int64
}

func (s *slot[V]) Signal(ctx context.Context, v V) {
	s.ls.Range(func(_ int64, l func(ctx context.Context, v V)) bool {
		l(ctx, v)
		return true
	})
}
