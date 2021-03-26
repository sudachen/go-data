package lazy_tests

import (
	"fmt"
	"gotest.tools/v3/assert"
	"sudachen.xyz/pkg/go-data/fu"
	"sudachen.xyz/pkg/go-data/lazy"
	"testing"
)

func Test_Drain_1(t *testing.T) {
	i := 0
	lazy.List(colors).MustDrain(lazy.Sink(func(v interface{}, _ error)(_ error){
		if v != nil {
			x := v.(Color)
			assert.Assert(t, colors[i].Color == x.Color)
			i++
		} else {
			assert.Assert(t, i == len(colors))
		}
		return
	}))
}

func Test_Drain_2(t *testing.T) {
	i := 0
	lazy.List(colors).MustDrain(lazy.Sink(func(v interface{}, _ error)(_ error){
		if v != nil {
			x := v.(Color)
			assert.Assert(t, colors[i].Color == x.Color)
			i++
		} else {
			assert.Assert(t, i == len(colors))
		}
		return
	}),8)
}

func Test_CcrDrain_1(t *testing.T) {
	c := fu.AtomicCounter{}
	lazy.List(colors).MustDrain(func(_ int)[]lazy.Worker {
		wrk := make([]lazy.Worker,8) // concurrency = 8
		for i := range wrk {
			wrk[i] = func(i int, v interface{}, _ error) error {
				if v != nil {
					switch x := v.(type) {
					case Color:
						assert.Assert(t, colors[i].Color == x.Color)
						i++
						c.Inc()

					}
				}
				return nil
			}
		}
		return wrk
	})
	fmt.Println(c.Value)
	assert.Assert(t, int(c.Value) == len(colors))
}
