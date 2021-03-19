package lazy

import "reflect"

func (zf Source) Filter(f interface{}) Source {
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
					if fv.Call([]reflect.Value{x})[0].Bool() {
						return v
					}
				}
			}
			return NoValue
		}
	}
}
