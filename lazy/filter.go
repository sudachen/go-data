package lazy

import "reflect"

func (zf Source) Filter(f interface{}) Source {
	return func(xs ...interface{}) Stream {
		z := zf(xs)
		fv := reflect.ValueOf(f)
		return func(next bool) (interface{},int) {
			switch v,i := z(next); v.(type) {
			case EndOfStream, struct{}:
				return v,i
			default:
				x := reflect.ValueOf(v)
				/*if x.Kind() == reflect.Interface {
					x = x.Elem()
				}*/
				if fv.Call([]reflect.Value{x})[0].Bool() {
					return v,i
				}
				return NoValue, i
			}
		}
	}
}
