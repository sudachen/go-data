package lazy_tests

import (
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
	"sudachen.xyz/pkg/go-data/lazy"
	"testing"
)

func Test_CollectFromChan_1(t *testing.T) {
	c := make(chan Color)
	go func() {
		for _, x := range colors {
			c <- x
		}
		close(c)
	}()
	var e []Color
	lazy.Chan(c).MustCollect(&e)
	assert.DeepEqual(t, e, colors)
}

func Test_CollectFromChan_2(t *testing.T) {
	c := make(chan Color)
	go func() {
		for _, x := range colors {
			c <- x
		}
		close(c)
	}()
	var e []Color
	lazy.Chan(c).MustCollect(&e, 8)
	assert.DeepEqual(t, e, colors)
}

func Test_CollectFromList_1(t *testing.T) {
	var e []Color
	lazy.List(colors).MustCollect(&e)
	assert.DeepEqual(t, e, colors)
}

func Test_CollectFromList_2(t *testing.T) {
	var e []Color
	lazy.List(colors).MustCollect(&e, 8)
	assert.DeepEqual(t, e, colors)
}

func Test_CollectFromEmpty_1(t *testing.T) {
	var e []Color
	lazy.List([]Color{}).MustCollect(&e)
	assert.Assert(t, len(e) == 0)
}

func Test_CollectFromEmpty_2(t *testing.T) {
	var e []Color
	lazy.List([]Color{}).MustCollect(&e, 8)
	assert.Assert(t, len(e) == 0)
}

func Test_CollectAnyFromList_1(t *testing.T) {
	e := lazy.List(colors).MustCollectAny()
	assert.DeepEqual(t, e, colors)
}

func Test_CollectAnyFromList_2(t *testing.T) {
	e := lazy.List(colors).MustCollectAny(8)
	assert.DeepEqual(t, e, colors)
}

func Test_CollectAnyFromEmpty(t *testing.T) {
	assert.Assert(t, cmp.Panics(func() {
		lazy.List([]Color{}).MustCollectAny(8)
	}))
}
