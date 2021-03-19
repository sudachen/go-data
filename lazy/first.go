package lazy

import (
	"sudachen.xyz/pkg/go-forge/fu"
)

func (zf Source) First(n int) Source {
	return func() Stream {
		z := zf()
		count := 0
		wc := fu.WaitCounter{Value: 0}
		return func(index Index) interface{} {
			v := z(index)
			if index != CloseSource && wc.Wait(index) {
				if count < n {
					if v != NoValue {
						switch v.(type) {
						case EndOfStream, Fail:
							wc.Stop()
						default:
							if count < n {
								count++
								wc.Inc()
							}
						}
					}
					return v
				}
			}
			return EndOfStream{}
		}
	}
}

