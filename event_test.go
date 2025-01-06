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

func TestCommon(t *testing.T) {
	tcs := []struct {
		name string
		init func() event.Event[string, int]
	}{
		{
			name: "sure",
			init: func() event.Event[string, int] { return event.Sure[string, int]() },
		},
		{
			name: "skip",
			init: func() event.Event[string, int] { return event.SKip[string, int]() },
		},
	}
	for _, tc := range tcs {
		desc := func(n string) string {
			return fmt.Sprintf("[%s] %s", tc.name, n)
		}
		make_event := func() (context.Context, event.Event[string, int]) {
			return context.TODO(), tc.init()
		}

		t.Run(desc("Listen with size"), func(t *testing.T) {
			ctx, e := make_event()

			l, close := e.Listen("", 3)
			defer close()

			e.Emit(ctx, "", 41)
			e.Emit(ctx, "", 42)
			e.Emit(ctx, "", 43)

			vs := []int{<-l, <-l, <-l}
			require.Equal(t, []int{41, 42, 43}, vs)
		})
		t.Run(desc("multiple Listenions"), func(t *testing.T) {
			ctx, e := make_event()

			l1, close1 := e.Listen("", 3)
			defer close1()
			l2, close2 := e.Listen("", 3)
			defer close2()
			l3, close3 := e.Listen("", 3)
			defer close3()

			e.Emit(ctx, "", 41)
			e.Emit(ctx, "", 42)
			e.Emit(ctx, "", 43)

			vs := []int{<-l1, <-l1, <-l1}
			require.Equal(t, []int{41, 42, 43}, vs)
			vs = []int{<-l2, <-l2, <-l2}
			require.Equal(t, []int{41, 42, 43}, vs)
			vs = []int{<-l3, <-l3, <-l3}
			require.Equal(t, []int{41, 42, 43}, vs)
		})
		t.Run(desc("close the Listenion closes the channel"), func(t *testing.T) {
			_, e := make_event()

			l, close := e.Listen("", 0)

			t0 := time.Now()
			Delayed(close)

			<-l
			dt := time.Since(t0)
			require.GreaterOrEqual(t, dt, Delay)
		})
	}
}
