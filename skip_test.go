package event_test

import (
	"context"
	"testing"

	"github.com/lesomnus/event"
	"github.com/stretchr/testify/require"
)

func TestSkip(t *testing.T) {
	make_slot := func() (context.Context, event.Slot[int]) {
		return context.TODO(), event.Skip[int]()
	}

	t.Run("emitted value is discarded if the channel is full", func(t *testing.T) {
		ctx, e := make_slot()

		l, close := e.Connect(1)
		defer close()

		e.Signal(ctx, 41)
		e.Signal(ctx, 42)

		v := <-l
		select {
		case v = <-l:
		default:
		}
		require.Equal(t, 41, v)
	})
	t.Run("emitted value is discarded if the receiver is not in waiting", func(t *testing.T) {
		ctx, e := make_slot()

		l, close := e.Connect(0)
		defer close()

		e.Signal(ctx, 41)
		Delayed(func() { e.Signal(ctx, 42) })

		v := <-l
		require.Equal(t, 42, v)
	})
}
