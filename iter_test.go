package re

import (
	"fmt"
	"iter"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConstant(t *testing.T) {
	it := Constant(time.Second, 4)
	next, stop := iter.Pull(it)
	defer stop()

	v, ok := next()
	require.Equal(t, time.Second, v)
	require.True(t, ok)

	var vs []time.Duration
	for v := range it {
		vs = append(vs, v)
	}
	require.Equal(t, []time.Duration{time.Second, time.Second, time.Second, time.Second}, vs)
}

func ExampleConstant() {
	next, stop := iter.Pull(Constant(time.Second, 4))
	defer stop()

	for i := range 5 {
		v, ok := next()
		fmt.Printf("#%v: { %v, %v }\n", i, v, ok)
	}
	// Output:
	// #0: { 1s, true }
	// #1: { 1s, true }
	// #2: { 1s, true }
	// #3: { 1s, true }
	// #4: { 0s, false }
}

func TestLinear(t *testing.T) {
	it := Linear(time.Second, 4)
	next, stop := iter.Pull(it)
	defer stop()

	v, ok := next()
	require.Equal(t, time.Second, v)
	require.True(t, ok)

	var vs []time.Duration
	for v := range it {
		vs = append(vs, v)
	}
	require.Equal(t, []time.Duration{time.Second, time.Second * 2, time.Second * 3, time.Second * 4}, vs)
}

func ExampleLinear() {
	next, stop := iter.Pull(Linear(time.Second, 4))
	defer stop()

	for i := range 5 {
		v, ok := next()
		fmt.Printf("#%v: { %v, %v }\n", i, v, ok)
	}
	// Output:
	// #0: { 1s, true }
	// #1: { 2s, true }
	// #2: { 3s, true }
	// #3: { 4s, true }
	// #4: { 0s, false }
}

func TestLinearRate(t *testing.T) {
	it := LinearRate(time.Second, time.Second*2, 4)
	next, stop := iter.Pull(it)
	defer stop()

	v, ok := next()
	require.Equal(t, time.Second, v)
	require.True(t, ok)

	var vs []time.Duration
	for v := range it {
		vs = append(vs, v)
	}
	require.Equal(t, []time.Duration{time.Second, time.Second * 3, time.Second * 5, time.Second * 7}, vs)
}

func ExampleLinearRate() {
	next, stop := iter.Pull(LinearRate(time.Second, time.Second*2, 4))
	defer stop()

	for i := range 5 {
		v, ok := next()
		fmt.Printf("#%v: { %v, %v }\n", i, v, ok)
	}
	// Output:
	// #0: { 1s, true }
	// #1: { 3s, true }
	// #2: { 5s, true }
	// #3: { 7s, true }
	// #4: { 0s, false }
}

func TestExponential(t *testing.T) {
	it := Exponential(time.Second, 4)
	next, stop := iter.Pull(it)
	defer stop()

	v, ok := next()
	require.Equal(t, time.Second, v)
	require.True(t, ok)

	var vs []time.Duration
	for v := range it {
		vs = append(vs, v)
	}
	require.Equal(t, []time.Duration{time.Second, time.Second * 2, time.Second * 4, time.Second * 8}, vs)
}

func ExampleExponential() {
	next, stop := iter.Pull(Exponential(time.Second, 4))
	defer stop()

	for i := range 5 {
		v, ok := next()
		fmt.Printf("#%v: { %v, %v }\n", i, v, ok)
	}
	// Output:
	// #0: { 1s, true }
	// #1: { 2s, true }
	// #2: { 4s, true }
	// #3: { 8s, true }
	// #4: { 0s, false }
}

func TestFibonacci(t *testing.T) {
	it := Fibonacci(time.Second, 5)
	next, stop := iter.Pull(it)
	defer stop()

	v, ok := next()
	require.Equal(t, time.Second, v)
	require.True(t, ok)

	var vs []time.Duration
	for v := range it {
		vs = append(vs, v)
	}
	require.Equal(t, []time.Duration{time.Second, time.Second * 2, time.Second * 3, time.Second * 5, time.Second * 8}, vs)
}

func ExampleFibonacci() {
	next, stop := iter.Pull(Fibonacci(time.Second, 5))
	defer stop()

	for i := range 6 {
		v, ok := next()
		fmt.Printf("#%v: { %v, %v }\n", i, v, ok)
	}
	// Output:
	// #0: { 1s, true }
	// #1: { 2s, true }
	// #2: { 3s, true }
	// #3: { 5s, true }
	// #4: { 8s, true }
	// #5: { 0s, false }
}

func TestJitter(t *testing.T) {
	it := Jitter(slices.Values([]time.Duration{time.Second, time.Second * 3, time.Second * 7}), time.Millisecond*100)
	next, stop := iter.Pull(it)
	defer stop()

	v, ok := next()
	require.True(t, time.Millisecond*900 <= v && v <= time.Millisecond*1100)
	require.True(t, ok)

	v, ok = next()
	require.True(t, time.Millisecond*2900 <= v && v <= time.Millisecond*3100)
	require.True(t, ok)

	v, ok = next()
	require.True(t, time.Millisecond*6900 <= v && v <= time.Millisecond*7100)
	require.True(t, ok)
}
