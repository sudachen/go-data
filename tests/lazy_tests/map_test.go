package lazy_tests

import (
	"fmt"
	"gotest.tools/assert"
	"sudachen.xyz/pkg/go-forge/lazy"
	"testing"
)

func Test1_Map1(t *testing.T) {
	rs := lazy.List([]int{0, 1, 2, 3, 4}).
		Map(func(r int) string { return fmt.Sprint(r) }).
		Parallel().
		MustCollectAny().([]string)
	assert.Assert(t, len(rs) == 5)
	for i, r := range rs {
		assert.Assert(t, r == fmt.Sprint(i))
	}
}

func Test2_Map1(t *testing.T) {
	rs := lazy.List([]int{0, 1, 2, 3, 4}).
		Map(func(r int) string { return fmt.Sprint(r) }).
		Parallel().
		MustCollectAny().([]string)
	assert.Assert(t, len(rs) == 5)
	for i, r := range rs {
		assert.Assert(t, r == fmt.Sprint(i))
	}
}

func Test2_Map2(t *testing.T) {
	dta := []int{7, 2, 3, 4, 5}
	rs := lazy.List(dta).
		Map1(func(ctx int)interface{}{
				return func(r int)string{
					return fmt.Sprint(r+ctx)
				}
			}).
		Parallel().
		MustCollectAny().([]string)
	assert.Assert(t, len(rs) == 5)
	for i, r := range rs {
		assert.Assert(t, r == fmt.Sprint(dta[i]+dta[0]))
	}
}
