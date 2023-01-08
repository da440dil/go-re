package re_test

import (
	"context"
	"fmt"
	"time"

	"github.com/da440dil/go-re"
)

func Example_fibonacci() {
	fn := func(ctx context.Context, i int) (bool, error) {
		if i == 0 {
			return true, nil
		}
		return false, fmt.Errorf("%w", re.ErrRetryable)
	}
	// Use exponential algorithm with delay between retries 10 ms with maximum number of retries 5.
	fn = re.Tryable(fn, re.Exponential(time.Millisecond*10), re.MaxRetries(5))

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
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
