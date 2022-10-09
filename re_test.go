package re

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type imock struct {
	d    time.Duration
	done bool
}

func (m imock) Next() (time.Duration, bool) {
	return m.d, m.done
}

func (m imock) Iterator() Iterator {
	return m
}

func TestTryable(t *testing.T) {
	b := imock{d: 0, done: true}
	ctx := context.Background()

	ok, err := Tryable(func(ctx context.Context, v int) (bool, error) {
		return true, nil
	}, b)(ctx, 0)
	require.NoError(t, err)
	require.True(t, ok)

	e := errors.New("some error")
	ok, err = Tryable(func(ctx context.Context, v int) (bool, error) {
		return false, e
	}, b)(ctx, 0)
	require.Equal(t, e, err)
	require.False(t, ok)

	e = fmt.Errorf("some error: %w", ErrRetryable)
	ok, err = Tryable(func(ctx context.Context, v int) (bool, error) {
		return false, e
	}, b)(ctx, 0)
	require.Equal(t, e, err)
	require.False(t, ok)

	b2 := imock{d: time.Millisecond * 100, done: false}
	i := 0
	fn := func(ctx context.Context, v *int) (bool, error) {
		*v++
		if *v%2 == 0 {
			return true, nil
		}
		return false, e
	}
	ok, err = Tryable(fn, b2)(ctx, &i)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = Tryable(fn, b, func(b Iterable) Iterable {
		return b2
	})(ctx, &i)
	require.NoError(t, err)
	require.True(t, ok)

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*150)
	defer cancel()
	ok, err = Tryable(func(ctx context.Context, v int) (bool, error) {
		return false, e
	}, b2)(ctx, 0)
	require.Equal(t, context.DeadlineExceeded, err)
	require.False(t, ok)
}
