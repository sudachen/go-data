package lazy

import (
	"reflect"
)

func Chan(c interface{}, stop ...chan struct{}) Source {
	return func(...interface{}) Stream {
		ni := 0
		scase := []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)}}
		return func(next bool) (interface{}, int) {
			if !next {
				for _, s := range stop {
					close(s)
				}
			}
			if _, r, ok := reflect.Select(scase); ok {
				i := ni
				ni++
				return r.Interface(), i
			}
			return EoS, ni
		}
	}
}
