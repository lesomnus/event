package event_test

import (
	"context"
	"testing"

	"github.com/lesomnus/event"
	"github.com/stretchr/testify/require"
)

func TestSkip(t *testing.T) {
	make_event := func() (context.Context, event.Event[string, int]) {
		return context.TODO(), event.SKip[string, int]()
	}

	t.Run("emitted value is discarded if the channel is full", func(t *testing.T) {
		ctx, e := make_event()

		l, close := e.Listen("", 1)
		defer close()

		e.Emit(ctx, "", 41)
		e.Emit(ctx, "", 42)

		v := <-l
		select {
		case v = <-l:
		default:
		}
		require.Equal(t, 41, v)
	})
	t.Run("emitted value is discarded if the receiver does not pending", func(t *testing.T) {
		ctx, e := make_event()

		l, close := e.Listen("", 0)
		defer close()

		e.Emit(ctx, "", 41)
		Delayed(func() { e.Emit(ctx, "", 42) })

		v := <-l
		require.Equal(t, 42, v)
	})
}
