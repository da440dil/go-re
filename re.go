// Package re provides re-execution for functions with configurable limits.
package re

import (
	"context"
	"errors"
	"time"
)

// ErrRetryable is error to signal function execution could be retried.
var ErrRetryable = errors.New("retryable")

// Retryable is a function which execution could be retried if error is ErrRetryable.
type Retryable[T any, U any] func(ctx context.Context, v T) (U, error)

// Tryable returns passed function with retry logic.
func Tryable[T any, U any](fn Retryable[T, U], b Iterable, fns ...Decorator) Retryable[T, U] {
	for _, fn := range fns {
		b = fn(b)
	}
	return func(ctx context.Context, v T) (U, error) {
		var it Iterator
		var timer *time.Timer
		for {
			r, err := fn(ctx, v)
			if err == nil {
				return r, err
			}
			if !errors.Is(err, ErrRetryable) {
				return r, err
			}
			if it == nil {
				it = b.Iterator()
			}
			d, done := it.Next()
			if done {
				return r, err
			}
			if timer == nil {
				timer = time.NewTimer(d)
				defer timer.Stop()
			} else {
				timer.Reset(d)
			}
			select {
			case <-ctx.Done():
				return r, ctx.Err()
			case <-timer.C:
			}
		}
	}
}
