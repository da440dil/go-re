package re

import (
	"iter"
	"math/rand"
	"time"
)

// Constant creates a sequence of durations whose values are constant.
func Constant(d time.Duration, n int) iter.Seq[time.Duration] {
	return func(yield func(time.Duration) bool) {
		for i := n; i > 0; i-- {
			if !yield(d) {
				return
			}
		}
	}
}

// Linear creates a sequence of durations whose values grow linearly.
func Linear(d time.Duration, n int) iter.Seq[time.Duration] {
	return func(yield func(time.Duration) bool) {
		v := d
		for i := n; i > 0; i-- {
			if !yield(v) {
				return
			}
			v += d
		}
	}
}

// LinearRate creates a sequence of durations whose values grow linearly with specified rate.
func LinearRate(d, rate time.Duration, n int) iter.Seq[time.Duration] {
	return func(yield func(time.Duration) bool) {
		v := d
		for i := n; i > 0; i-- {
			if !yield(v) {
				return
			}
			v += rate
		}
	}
}

// Exponential creates a sequence of durations whose values grow exponentially.
func Exponential(d time.Duration, n int) iter.Seq[time.Duration] {
	return func(yield func(time.Duration) bool) {
		v := d
		for i := n; i > 0; i-- {
			if !yield(v) {
				return
			}
			v += v
		}
	}
}

// Fibonacci creates a sequence of durations whose values grow according to the Fibonacci algorithm.
func Fibonacci(d time.Duration, n int) iter.Seq[time.Duration] {
	return func(yield func(time.Duration) bool) {
		x, v := time.Duration(0), d
		for i := n; i > 0; i-- {
			x, v = v, x+v
			if !yield(v) {
				return
			}
		}
	}
}

// Jitter sets maximum duration randomly added to or extracted from delay between retries to improve performance under high contention.
func Jitter(it iter.Seq[time.Duration], d time.Duration) iter.Seq[time.Duration] {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := int64(d)
	n := x*2 + 1
	return func(yield func(time.Duration) bool) {
		for v := range it {
			v = max(v+time.Duration(random.Int63n(n)-x), 0)
			if !yield(v) {
				return
			}
		}
	}
}
