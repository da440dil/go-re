# go-re

Re-execution for functions with configurable limits.

[Example](./examples/jitter/main.go) usage:

```go
import (
	"context"
	"fmt"
	"time"

	"github.com/da440dil/go-re"
)

func main() {
	fn := func(ctx context.Context, i int) (bool, error) {
		if i == 0 {
			return true, nil
		}
		return false, fmt.Errorf("%w", re.ErrRetryable)
	}
	// Use linear algorithm with delay between retries 10ms with maximum number of retries 5.
	// Set 5ms maximum duration randomly added to or extracted from delay between retries
	// to improve performance under high contention.
	fn = re.Tryable(
		fn,
		re.Linear(time.Millisecond*10),
		re.WithMaxRetries(5),
		re.WithJitter(time.Millisecond*5),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	for i := 0; i < 3; i++ {
		ok, err := fn(ctx, i)
		fmt.Printf("{ i: %v, ok: %v, err: %v }\n", i, ok, err)
	}
	// Output:
	// { i: 0, ok: true, err: <nil> }
	// { i: 1, ok: false, err: retryable }
	// { i: 2, ok: false, err: context deadline exceeded }
}
```
