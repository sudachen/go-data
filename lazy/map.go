package lazy

import (
	"reflect"
)

func (zf Source) Map(f interface{}) Source {
	return func(xs ...interface{}) Stream {
		z := zf(xs...)
		fv := reflect.ValueOf(f)
		return func(next bool) (interface{}, int) {
			switch v,i := z(next); v.(type) {
			case EndOfStream, struct{}:
				return v, i
			default:
				x := reflect.ValueOf(v)
				/*for x.Kind() == reflect.Interface {
					x = x.Elem()
				}*/
				//fmt.Println(x.Interface(), x.Type())
				return fv.Call([]reflect.Value{x})[0].Interface(), i
			}
		}
	}
}

//
// source.Map1(f) => source.Map(source.Map(f).MustGetOne())
// 1. calculates fx = f(first stream value)
// 2. maps all as fx(stream value)

func (zf Source) Map1(fx interface{}) Source {
	return func(xs ...interface{}) Stream {
		z := zf(xs...)
		var fv func(interface{}) interface{}
		return func(next bool) (interface{}, int) {
			switch v, i := z(next); v.(type) {
			case EndOfStream, struct{}:
				return v,i
			default:
				if fv == nil {
					f := reflect.ValueOf(fx).Call([]reflect.Value{reflect.ValueOf(v)})[0]
					var ok bool
					for f.Kind() == reflect.Interface {
						f = f.Elem()
					}
					if fv, ok = f.Interface().(func(interface{}) interface{}); !ok {
						fv = func(v interface{}) interface{} {
							x := reflect.ValueOf(v)
							/*if x.Kind() == reflect.Interface {
								x = x.Elem()
							}*/
							return f.Call([]reflect.Value{x})[0].Interface()
						}
					}
				}
				return fv(v),i
			}
		}
	}
}

