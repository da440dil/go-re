package re

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTryable(t *testing.T) {
	fn := func(x *int, e error) func(ctx context.Context, v int) (bool, error) {
		return func(ctx context.Context, v int) (bool, error) {
			*x += v
			// { 0 => true, 1 => false, 2 => true, ... => false }
			if *x > 2 || *x%2 != 0 {
				return false, e
			}
			return true, nil
		}
	}
	it := slices.Values([]time.Duration{time.Millisecond * 100})
	ctx := context.Background()
	e := errors.New("some error")

	x := 1
	ok, err := Tryable(fn(&x, e), it)(ctx, 1)
	require.True(t, ok)
	require.NoError(t, err)
	require.Equal(t, 2, x)

	x = 0
	ok, err = Tryable(fn(&x, e), it)(ctx, 1)
	require.False(t, ok)
	require.Equal(t, e, err)
	require.Equal(t, 1, x)

	x = 0
	e = fmt.Errorf("some error: %w", ErrRetryable)
	ok, err = Tryable(fn(&x, e), it)(ctx, 1)
	require.True(t, ok)
	require.NoError(t, err)
	require.Equal(t, 2, x)

	x = 2
	ok, err = Tryable(fn(&x, e), it)(ctx, 1)
	require.False(t, ok)
	require.Equal(t, e, err)
	require.Equal(t, 4, x)

	it = slices.Values([]time.Duration{time.Millisecond * 100, time.Millisecond * 100})
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*150)
	defer cancel()

	x = 2
	ok, err = Tryable(fn(&x, e), it)(ctx, 1)
	require.False(t, ok)
	require.Equal(t, context.DeadlineExceeded, err)
	require.Equal(t, 4, x)
}
