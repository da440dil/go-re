package re_test

import (
	"context"
	"fmt"
	"time"

	"github.com/da440dil/go-re"
)

func Example_fibonacci() {
	fn := func(ctx context.Context, x int) (bool, error) {
		// { 0 => true, 1 => false, 2 => true, ... => false }
		if x > 2 || x%2 != 0 {
			return false, fmt.Errorf("%w", re.ErrRetryable)
		}
		return true, nil
	}
	// Use Fibonacci algorithm with delay between retries 10 ms & maximum number of retries 4.
	fn = re.Tryable(fn, re.Fibonacci(time.Millisecond*10, 4))

	for i := range 3 {
		ok, err := fn(context.Background(), i)
		fmt.Printf("ok: %v, err: %v\n", ok, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*40)
	defer cancel()

	ok, err := fn(ctx, 3)
	fmt.Printf("ok: %v, err: %v\n", ok, err)
	// Output:
	// ok: true, err: <nil>
	// ok: false, err: retryable
	// ok: true, err: <nil>
	// ok: false, err: context deadline exceeded
}
