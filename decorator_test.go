package re

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMaxRetries(t *testing.T) {
	b := Constant(time.Second)
	b = MaxRetries(3)(b)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Duration(0), d)
		require.True(t, done)
	}
}

func ExampleMaxRetries_constant() {
	it := MaxRetries(3)(Constant(time.Second)).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 1s, false }
	// #2: { 1s, false }
	// #3: { 0s, true }
}

func ExampleMaxRetries_linear() {
	it := MaxRetries(3)(Linear(time.Second)).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 2s, false }
	// #2: { 3s, false }
	// #3: { 0s, true }
}

func ExampleMaxRetries_exponential() {
	it := MaxRetries(3)(Exponential(time.Second)).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 2s, false }
	// #2: { 4s, false }
	// #3: { 0s, true }
}

func TestJitter(t *testing.T) {
	b := Linear(time.Second)
	b = MaxRetries(3)(b)
	b = Jitter(time.Millisecond * 100)(b)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.True(t, time.Millisecond*900 <= d && d <= time.Millisecond*1100)
		require.False(t, done)
		d, done = it.Next()
		require.True(t, time.Millisecond*1900 <= d && d <= time.Millisecond*2100)
		require.False(t, done)
		d, done = it.Next()
		require.True(t, time.Millisecond*2900 <= d && d <= time.Millisecond*3100)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Duration(0), d)
		require.True(t, done)
	}

	// for test coverage
	b = Linear(time.Millisecond * -100)
	b = Jitter(time.Millisecond * 100)(b)
	it := b.Iterator()
	d, done := it.Next()
	require.Equal(t, time.Duration(0), d)
	require.False(t, done)
	d, done = it.Next()
	require.Equal(t, time.Duration(0), d)
	require.False(t, done)
}
