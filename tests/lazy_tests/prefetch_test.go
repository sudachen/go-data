package lazy_tests

import (
	"gotest.tools/v3/assert"
	"reflect"
	"sudachen.xyz/pkg/go-data/errors"
	"sudachen.xyz/pkg/go-data/lazy"
	"testing"
)

func linlist(list interface{}) lazy.Source {
	return func(xs ...interface{}) lazy.Stream {
		worker := 0
		open := lazy.NoPrefetch
		for _, x := range xs {
			if f, ok := x.(func()(int,int,lazy.Prefetch)); ok {
				worker, _, open = f()
			} else {
				return lazy.Error(errors.Errorf("unsupported source option: %v", x))
			}
		}
		v := reflect.ValueOf(list)
		return open(worker,func()lazy.Stream{
			index := 0
			return func(next bool) (r interface{}, i int) {
				if next && index < v.Len() {
					i = index
					r = v.Index(i).Interface()
					index++
					return r, i
				}
				return lazy.EoS, index
			}
		})
	}
}

func Test_Prefetch_1(t *testing.T) {
	i := 0
	linlist(colors).MustDrain(lazy.Sink(func(v interface{}, _ error)(_ error){
		if v != nil {
			x := v.(Color)
			y := colors[i]
			assert.Assert(t, y.Color == x.Color)
			i++
		} else {
			assert.Assert(t, i == len(colors))
		}
		return
	}),2)
}

func Test_Prefetch_2(t *testing.T) {
	a := lazy.Sequence(func(i int)interface{}{
		if i < 100 { return i }
		return lazy.EoS
	}).MustCollectAny(2).([]int)
	assert.Assert(t, len(a) == 100)
}
