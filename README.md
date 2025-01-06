# event

## Usage

```go
import (
	"github.com/lesomnus/event"
)

func main() {
	ctx := context.TODO()
	e := event.Sure[string, int]()
	
	foo, close_foo := e.Listen("value", 0)
	defer close_foo()
	time.AfterFunc(100 * time.Milliseconds, func() {
		// v == 42
		v := <-foo
	})

	bar, close_bar := e.Listen("value", 0)
	defer close_bar()
	time.AfterFunc(100 * time.Milliseconds, func() {
		// v == 42
		v := <-bar
	})

	// Unblocked after 100ms
	e.Emit("value", 42)
}
```

## Implementations

- `Sure.Emit` blocks until all receivers have received the value.
- `Skip.Emit` skips any receiver that cannot receive a value at the time. Therefore, some receivers may not receive all emitted messages. However, if the channel size is sufficient, all messages will be delivered.
