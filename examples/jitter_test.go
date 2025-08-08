package re_test

import (
	"context"
	"fmt"
	"time"

	"github.com/da440dil/go-re"
)

func Example_jitter() {
	fn := func(ctx context.Context, x int) (bool, error) {
		// { 0 => true, 1 => false, 2 => true, ... => false }
		if x > 2 || x%2 != 0 {
			return false, fmt.Errorf("%w", re.ErrRetryable)
		}
		return true, nil
	}
	// Use linear algorithm with delay between retries 20 ms with maximum number of retries 3.
	// Set 5 ms maximum duration randomly added to or extracted from delay between retries to improve performance under high contention.
	fn = re.Tryable(fn, re.Jitter(re.Linear(time.Millisecond*20, 3), time.Millisecond*5))

	for i := range 3 {
		ok, err := fn(context.Background(), i)
		fmt.Printf("ok: %v, err: %v\n", ok, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*30)
	defer cancel()

	ok, err := fn(ctx, 3)
	fmt.Printf("ok: %v, err: %v\n", ok, err)
	// Output:
	// ok: true, err: <nil>
	// ok: false, err: retryable
	// ok: true, err: <nil>
	// ok: false, err: context deadline exceeded
}
