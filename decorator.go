package re

import (
	"math/rand"
	"time"
)

// Decorator extends behavior of an iterable.
type Decorator func(Iterable) Iterable

type maxRetriesB struct {
	b Iterable
	n int
}

func (b maxRetriesB) Iterator() Iterator {
	return &maxRetriesI{b.b.Iterator(), b.n}
}

type maxRetriesI struct {
	i Iterator
	n int
}

func (i *maxRetriesI) Next() (time.Duration, bool) {
	if i.n > 0 {
		i.n--
		return i.i.Next()
	}
	return 0, true
}

// MaxRetries sets maximum number of retries.
func MaxRetries(n int) Decorator {
	return func(b Iterable) Iterable {
		return maxRetriesB{b, n}
	}
}

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type jitterB struct {
	b    Iterable
	n, j int64
}

func (b jitterB) Iterator() Iterator {
	return jitterI{b.b.Iterator(), b.n, b.j}
}

type jitterI struct {
	i    Iterator
	n, j int64
}

func (i jitterI) Next() (time.Duration, bool) {
	v, done := i.i.Next()
	if done {
		return 0, done
	}
	v = v + time.Duration(random.Int63n(i.n)-i.j)
	if v < 0 {
		v = 0
	}
	return v, done
}

// Jitter sets maximum duration randomly added to or extracted from delay between retries to improve performance under high contention.
func Jitter(d time.Duration) Decorator {
	return func(b Iterable) Iterable {
		j := int64(d)
		return jitterB{b, j*2 + 1, j}
	}
}
