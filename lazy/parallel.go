package lazy

import (
	"runtime"
	"sudachen.xyz/pkg/go-forge/fu"
)

func (zf Source) Parallel(concurrency ...int) Source {
	return func() Stream {
		z := zf()
		ccrn := fu.Fnzi(fu.Fnzi(concurrency...), runtime.NumCPU()+1)
		index := fu.AtomicCounter{0}
		wc := fu.WaitCounter{Value: 0}
		c := make(chan interface{}, ccrn)
		stop := make(chan struct{})
		alive := fu.AtomicCounter{uint64(ccrn)}
		for i := 0; i < ccrn; i++ {
			go func() {
			loop:
				for !wc.Stopped() {
					n := index.PostIncIfLess(CloseSource)
					if n == CloseSource {
						wc.Stop()
						break loop
					}
					v := z(n)
					if v != NoValue && wc.Wait(n) {
						switch v.(type) {
						case EndOfStream:
							wc.Stop()
							break loop
						case Fail:
							wc.Stop()
						}
						select {
						case c <- v:
							wc.Inc()
						case <-stop:
							wc.Stop()
							break loop
						}
					}
				}
				if alive.Dec() == 0 { // returns new value
					close(c)
				}
			}()
		}
		return func(index Index) interface{} {
			if index == CloseSource {
				close(stop)
				return z(CloseSource)
			}
			if v, ok := <-c; ok {
				return v
			}
			return EndOfStream{}
		}
	}
}

func (zf Source) Concurrent(concurrency ...int) Source {
	return func() Stream {
		z := zf()
		ccrn := fu.Fnzi(fu.Fnzi(concurrency...), runtime.NumCPU())
		index := fu.AtomicCounter{0}
		c := make(chan interface{})
		stop := make(chan struct{})
		alive := fu.AtomicCounter{uint64(ccrn)}
		eof := fu.AtomicFlag{}
		for i := 0; i < ccrn; i++ {
			go func() {
			loop:
				for !eof.State() {
					n := index.PostIncIfLess(CloseSource)
					if n == CloseSource  {
						break loop
					}
					v := z(n)
					if v != NoValue {
						switch v.(type) {
						case EndOfStream:
							break loop
						case Fail:
							eof.Set()
						}
						select {
						case c <- v:
							// nothing
						case <-stop:
							eof.Set()
							break loop
						}
					}
				}
				if alive.Dec() == 0 { // returns new value
					close(c)
				}
			}()
		}
		return func(index Index) interface{} {
			if index == CloseSource {
				close(stop)
				return z(CloseSource)
			}
			if v, ok := <-c; ok {
				return v
			}
			return EndOfStream{}
		}
	}
}
