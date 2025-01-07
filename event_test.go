package event_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lesomnus/event"
	"github.com/stretchr/testify/require"
)

const Delay = 100 * time.Millisecond

func Delayed(f func()) {
	time.AfterFunc(Delay, f)
}

func TestSlotCommon(t *testing.T) {
	tcs := []struct {
		name string
		init func() event.Slot[int]
	}{
		{
			name: "sure",
			init: event.Sure[int],
		},
		{
			name: "skip",
			init: event.Skip[int],
		},
	}
	for _, tc := range tcs {
		desc := func(n string) string {
			return fmt.Sprintf("[%s] %s", tc.name, n)
		}
		make_slot := func() (context.Context, event.Slot[int]) {
			return context.TODO(), tc.init()
		}

		t.Run(desc("connect with size"), func(t *testing.T) {
			ctx, e := make_slot()

			l, close := e.Connect(3)
			defer close()

			e.Signal(ctx, 41)
			e.Signal(ctx, 42)
			e.Signal(ctx, 43)

			vs := []int{<-l, <-l, <-l}
			require.Equal(t, []int{41, 42, 43}, vs)
		})
		t.Run(desc("multiple connections"), func(t *testing.T) {
			ctx, e := make_slot()

			l1, close1 := e.Connect(3)
			defer close1()
			l2, close2 := e.Connect(3)
			defer close2()
			l3, close3 := e.Connect(3)
			defer close3()

			e.Signal(ctx, 41)
			e.Signal(ctx, 42)
			e.Signal(ctx, 43)

			vs := []int{<-l1, <-l1, <-l1}
			require.Equal(t, []int{41, 42, 43}, vs)
			vs = []int{<-l2, <-l2, <-l2}
			require.Equal(t, []int{41, 42, 43}, vs)
			vs = []int{<-l3, <-l3, <-l3}
			require.Equal(t, []int{41, 42, 43}, vs)
		})
		t.Run(desc("channel is closed when the connection is closed"), func(t *testing.T) {
			_, e := make_slot()

			l, close := e.Connect(0)

			t0 := time.Now()
			Delayed(close)

			<-l
			dt := time.Since(t0)
			require.GreaterOrEqual(t, dt, Delay)
		})
	}
}
