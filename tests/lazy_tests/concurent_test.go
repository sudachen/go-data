package lazy_tests

import (
	"gotest.tools/assert"
	"sort"
	"sudachen.xyz/pkg/go-forge/lazy"
	"testing"
)

func Test1_Concurrent(t *testing.T) {
	rs := lazy.List([]int{0, 1, 2, 3, 4}).
		Concurrent().
		MustCollectAny().([]int)
	assert.Assert(t, len(rs) == 5)
	for i, r := range rs {
		assert.Assert(t, r == i)
	}
}

func Test2_Concurrent(t *testing.T) {
	const count = 10000
	c := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			c <- i
		}
		close(c)
	}()
	rs := lazy.Chan(c).
		Concurrent().
		MustCollectAny().([]int)
	sort.Ints(rs)
	assert.Assert(t, len(rs) == count)
	for i, r := range rs {
		assert.Assert(t, r == i)
	}
}

