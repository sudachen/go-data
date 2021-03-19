package lazy

import (
	"reflect"
	"sudachen.xyz/pkg/go-forge/fu"
)

func Chan(c interface{}, stop ...chan struct{}) Source {
	return func() Stream {
		scase := []reflect.SelectCase{{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(c)}}
		wc := fu.WaitCounter{Value: 0}
		return func(index Index) interface{} {
			if index == CloseSource {
				wc.Stop()
				for _, s := range stop {
					close(s)
				}
			}
			if wc.Wait(index) {
				_, r, ok := reflect.Select(scase)
				if wc.Inc() && ok {
					return r.Interface()
				}
			}
			return EndOfStream{}
		}
	}
}

