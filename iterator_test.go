package re

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConstant(t *testing.T) {
	b := Constant(time.Second)
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
		require.Equal(t, time.Second, d)
		require.False(t, done)
	}
}

func ExampleConstant() {
	it := Constant(time.Second).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 1s, false }
	// #2: { 1s, false }
	// #3: { 1s, false }
}

func TestLinear(t *testing.T) {
	b := Linear(time.Second)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*2, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*3, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*4, d)
		require.False(t, done)
	}
}

func ExampleLinear() {
	it := Linear(time.Second).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 2s, false }
	// #2: { 3s, false }
	// #3: { 4s, false }
}

func TestLinearRate(t *testing.T) {
	b := LinearRate(time.Second, time.Second*2)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*3, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*5, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*7, d)
		require.False(t, done)
	}
}

func ExampleLinearRate() {
	it := LinearRate(time.Second, time.Second*2).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 3s, false }
	// #2: { 5s, false }
	// #3: { 7s, false }
}

func TestExponential(t *testing.T) {
	b := Exponential(time.Second)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*2, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*4, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*8, d)
		require.False(t, done)
	}
}

func ExampleExponential() {
	it := Exponential(time.Second).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 2s, false }
	// #2: { 4s, false }
	// #3: { 8s, false }
}

func TestExponentialRate(t *testing.T) {
	b := ExponentialRate(time.Second, 0.2)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Millisecond*1000, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*1200, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*1440, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*1728, d)
		require.False(t, done)
	}
}

func ExampleExponentialRate() {
	it := ExponentialRate(time.Second, 0.2).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 1.2s, false }
	// #2: { 1.44s, false }
	// #3: { 1.728s, false }
}

func TestFibonacci(t *testing.T) {
	b := Fibonacci(time.Millisecond * 10)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Millisecond*10, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*20, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*30, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*50, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*80, d)
		require.False(t, done)
	}
}

func ExampleFibonacci() {
	it := Fibonacci(time.Millisecond * 10).Iterator()
	for i := 0; i < 5; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 10ms, false }
	// #1: { 20ms, false }
	// #2: { 30ms, false }
	// #3: { 50ms, false }
	// #4: { 80ms, false }
}
