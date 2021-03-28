package lazy

import (
	"reflect"
	"sudachen.xyz/pkg/go-data/errors"
)

const iniCollectLength = 13

//
// e := make([]int,0,0)
// err := Source().Collect(&e)
// Source().MustCollect(&e)
//

func SinkTo(to interface{}, reserve ...int) WorkerFactory {
	values := reflect.ValueOf(to).Elem()
	return Sink(func(v interface{}, err error)(_ error) {
		if v != nil {
			values = reflect.Append(values, reflect.ValueOf(v))
		} else if err == nil {
			reflect.ValueOf(to).Elem().Set(values)
		}
		return
	})
}

func (zf Source) Collect(to interface{}, reserve ...int) error {
	return zf.Drain(SinkTo(to, reserve...))
}

func (zf Source) MustCollect(to interface{}, reserve ...int) {
	err := zf.Collect(to, reserve...)
	if err != nil {
		panic(errors.Panic{err})
	}
}

func (zf Source) CollectAny(concurrency ...int) (interface{}, error) {
	var to reflect.Value
	err := zf.Drain(Sink(func(v interface{},err error)(_ error){
		if v != nil {
			x := reflect.ValueOf(v)
			if !to.IsValid() {
				to = reflect.MakeSlice(reflect.SliceOf(x.Type()), 0, iniCollectLength)
			}
			to = reflect.Append(to, x)
		}
		return
	}), concurrency...)
	if err != nil {
		return nil, err
	}
	if !to.IsValid() {
		return nil, errors.New("nothing to collect")
	}
	return to.Interface(), nil
}

func (zf Source) MustCollectAny(concurrency ...int) interface{} {
	ret, err := zf.CollectAny(concurrency...)
	if err != nil {
		panic(errors.Panic{err})
	}
	return ret
}

