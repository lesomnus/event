package event_test

import (
	"context"
	"testing"
	"time"

	"github.com/lesomnus/event"
	"github.com/stretchr/testify/require"
)

func TestSure(t *testing.T) {
	make_event := func() (context.Context, event.Event[string, int]) {
		return context.TODO(), event.Sure[string, int]()
	}

	t.Run("emit is blocked until received", func(t *testing.T) {
		ctx, e := make_event()

		l, close := e.Listen("", 0)
		defer close()

		t0 := time.Now()
		v := 0
		Delayed(func() { v = <-l })

		e.Emit(ctx, "", 42)
		dt := time.Since(t0)
		require.GreaterOrEqual(t, dt, Delay)
		require.Equal(t, 42, v)
	})
	t.Run("emit is unblocked by context", func(t *testing.T) {
		ctx, e := make_event()
		ctx, cancel := context.WithCancel(ctx)

		_, close := e.Listen("", 0)
		defer close()

		t0 := time.Now()
		v := 36
		Delayed(cancel)

		e.Emit(ctx, "", 42)
		dt := time.Since(t0)
		require.GreaterOrEqual(t, dt, Delay)
		require.Equal(t, 36, v)
	})
	t.Run("emit is unblocked by Listenion close", func(t *testing.T) {
		ctx, e := make_event()

		_, close := e.Listen("", 0)
		defer close()

		t0 := time.Now()
		v := 36
		Delayed(close)

		e.Emit(ctx, "", 42)
		dt := time.Since(t0)
		require.GreaterOrEqual(t, dt, Delay)
		require.Equal(t, 36, v)
	})
	t.Run("no Listens no block", func(t *testing.T) {
		_, e := make_event()
		e.Emit(context.TODO(), "", 42)
	})
	t.Run("closed Listenion does not block emit", func(t *testing.T) {
		ctx, e := make_event()

		_, close := e.Listen("", 0)
		close()

		t0 := time.Now()
		Delayed(func() {})

		e.Emit(ctx, "", 42)
		dt := time.Since(t0)
		require.Less(t, dt, Delay)
	})
}
