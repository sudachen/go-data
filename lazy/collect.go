package lazy

import (
	"reflect"
	"sudachen.xyz/pkg/go-forge/errors"
)

//
// e := make([]int,0,0)
// err := Source().Collect(&e)
// Source().MustCollect(&e)
//


func (zf Source) Collect(to interface{}) error {
	values := reflect.ValueOf(to).Elem()
	err := zf.Drain(func(v interface{}) error {
		if v != DrainFailed && v != DrainSucceed {
			values = reflect.Append(values, reflect.ValueOf(v))
		}
		return nil
	})
	if err != nil { return err}
	reflect.ValueOf(to).Elem().Set(values)
	return nil
}

func (zf Source) MustCollect(to interface{}) {
	err := zf.Collect(to)
	if err != nil { panic(errors.PanicBtrace{err}) }
}

func (zf Source) CollectAny() (interface{},error) {
	var to reflect.Value
	err := zf.Drain(func(v interface{}) error {
		if v != DrainFailed && v != DrainSucceed {
			x := reflect.ValueOf(v)
			if !to.IsValid() {
				to = reflect.MakeSlice(reflect.SliceOf(x.Type()), 0, iniCollectLength)
			}
			to = reflect.Append(to, x)
		}
		return nil
	})
	if err != nil { return nil, err }
	if !to.IsValid() { return nil, errors.New("no values collected") }
	return to.Interface(), nil
}

func (zf Source) MustCollectAny() interface{} {
	ret, err := zf.CollectAny()
	if err != nil { panic(errors.PanicBtrace{err}) }
	return ret
}
