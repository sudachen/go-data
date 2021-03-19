package lazy

import (
	"reflect"
	"sudachen.xyz/pkg/go-forge/fu"
	"sync"
)

func (zf Source) Map(f interface{}) Source {
	return func() Stream {
		z := zf()
		fv := reflect.ValueOf(f)
		return func(index Index) interface{} {
			if v := z(index); v != NoValue {
				switch v.(type) {
				case Fail, EndOfStream:
					return v
				default:
					x := reflect.ValueOf(v)
					if x.Kind() == reflect.Interface {
						x = x.Elem()
					}
					return fv.Call([]reflect.Value{x})[0].Interface()
				}
			}
			return NoValue
		}
	}
}

//
// source.Map1(f) => source.Map(source.Map(f).MustGetOne())
// 1. calculates fx = f(first stream value)
// 2. maps all as fx(stream value)

func (zf Source) Map1(fx interface{}) Source {
	return func() Stream {
		z := zf()
		var fv reflect.Value
		ctxf := fu.AtomicFlag{}
		mux := sync.Mutex{}
		return func(index Index) interface{} {
			if v := z(index); v != NoValue {
				switch v.(type) {
				case Fail, EndOfStream:
					return v
				default:
					x := reflect.ValueOf(v)
					if x.Kind() == reflect.Interface {
						x = x.Elem()
					}
					if !ctxf.State() {
						mux.Lock()
						if !ctxf.State() {
							fv = reflect.ValueOf(fx).Call([]reflect.Value{x})[0]
							if fv.Kind() == reflect.Interface {
								fv = fv.Elem()
							}
						}
						ctxf.Set()
						mux.Unlock()
					}
					return fv.Call([]reflect.Value{x})[0].Interface()
				}
			}
			return NoValue
		}
	}
}
