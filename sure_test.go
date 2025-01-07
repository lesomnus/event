package event_test

import (
	"context"
	"testing"
	"time"

	"github.com/lesomnus/event"
	"github.com/stretchr/testify/require"
)

func TestSure(t *testing.T) {
	make_slot := func() (context.Context, event.Slot[int]) {
		return context.TODO(), event.Sure[int]()
	}

	t.Run("signal is blocked until the value is received", func(t *testing.T) {
		ctx, e := make_slot()

		l, close := e.Connect(0)
		defer close()

		t0 := time.Now()
		v := 0
		Delayed(func() { v = <-l })

		e.Signal(ctx, 42)
		dt := time.Since(t0)
		require.GreaterOrEqual(t, dt, Delay)
		require.Equal(t, 42, v)
	})
	t.Run("signal is unblocked by context", func(t *testing.T) {
		ctx, e := make_slot()
		ctx, cancel := context.WithCancel(ctx)

		_, close := e.Connect(0)
		defer close()

		t0 := time.Now()
		v := 36
		Delayed(cancel)

		e.Signal(ctx, 42)
		dt := time.Since(t0)
		require.GreaterOrEqual(t, dt, Delay)
		require.Equal(t, 36, v)
	})
	t.Run("signal is unblocked by closing the connection", func(t *testing.T) {
		ctx, e := make_slot()

		_, close := e.Connect(0)
		defer close()

		t0 := time.Now()
		v := 36
		Delayed(close)

		e.Signal(ctx, 42)
		dt := time.Since(t0)
		require.GreaterOrEqual(t, dt, Delay)
		require.Equal(t, 36, v)
	})
	t.Run("no connections no block", func(t *testing.T) {
		_, s := make_slot()
		s.Signal(context.TODO(), 42)
	})
	t.Run("emit on closed connection is not blocked", func(t *testing.T) {
		ctx, e := make_slot()

		_, close := e.Connect(0)
		close()

		t0 := time.Now()
		Delayed(func() {})

		e.Signal(ctx, 42)
		dt := time.Since(t0)
		require.Less(t, dt, Delay)
	})
}
