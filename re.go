// Package re provides re-execution for functions with configurable limits.
package re

import (
	"context"
	"errors"
	"iter"
	"time"
)

// ErrRetryable is error to signal function execution could be retried.
var ErrRetryable = errors.New("retryable")

// Retryable is a function which execution could be retried if error is ErrRetryable.
type Retryable[T any, U any] func(ctx context.Context, v T) (U, error)

// Tryable returns passed function with retry logic.
func Tryable[T any, U any](fn Retryable[T, U], it iter.Seq[time.Duration]) Retryable[T, U] {
	return func(ctx context.Context, v T) (U, error) {
		var next func() (time.Duration, bool)
		var timer *time.Timer
		for {
			r, err := fn(ctx, v)
			if err == nil {
				return r, err
			}
			if !errors.Is(err, ErrRetryable) {
				return r, err
			}

			if next == nil {
				var stop func()
				next, stop = iter.Pull(it)
				defer stop()
			}
			d, ok := next()
			if !ok {
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
