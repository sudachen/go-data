package lazy

import (
	"reflect"
	"sudachen.xyz/pkg/go-forge/fu"
)

func List(list interface{}) Source {
	return func() Stream {
		v := reflect.ValueOf(list)
		l := uint64(v.Len())
		flag := fu.AtomicFlag{Value: 1}
		return func(index Index) interface{} {
			if index < l && flag.State() {
				return v.Index(int(index)).Interface()
			}
			return EndOfStream{}
		}
	}
}
