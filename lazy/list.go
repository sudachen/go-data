package lazy

import (
	"reflect"
	"sudachen.xyz/pkg/go-data/errors"
)

func List(list interface{}) Source {
	return func(xs ...interface{}) Stream {
		index, stride := 0, 1
		for _, x := range xs {
			if f, ok := x.(func() (int, int, Prefetch)); ok {
				index, stride, _ = f()
			} else {
				return Error(errors.Errorf("unsupported source option: %v", x))
			}
		}
		v := reflect.ValueOf(list)
		return func(next bool) (r interface{}, i int) {
			if next && index < v.Len() {
				r, i = v.Index(index).Interface(), index
				index += stride
				return
			}
			return EoS, index
		}
	}
}

func Sequence(gen func(int) interface{}) Source {
	return func(xs ...interface{}) Stream {
		pf := NoPrefetch
		worker := 0
		for _, x := range xs {
			if f, ok := x.(func() (int, int, Prefetch)); ok {
				worker, _, pf = f()
			} else {
				return Error(errors.Errorf("unsupported source option: %v", x))
			}
		}
		return pf(worker, func() Stream {
			n := 0
			return func(next bool) (v interface{}, i int) {
				if next {
					v, i = gen(n), n
					n++
					return
				}
				return EoS, n
			}
		})
	}
}

func Generator(gen func(int) interface{}) Source {
	return func(xs ...interface{}) Stream {
		index, stride := 0, 1
		for _, x := range xs {
			if f, ok := x.(func() (int, int, Prefetch)); ok {
				index, stride, _ = f()
			} else {
				return Error(errors.Errorf("unsupported source option: %v", x))
			}
		}
		return func(next bool) (v interface{}, i int) {
			if next {
				v, i = gen(index), index
				index += stride
				return
			}
			return EoS, index
		}
	}
}
