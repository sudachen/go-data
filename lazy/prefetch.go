package lazy

import (
	"runtime"
	"sync"
	"sync/atomic"
)

const PrefetchFactor = 2
type Prefetch func(int, func()Stream)Stream

type prefetch struct{
	width int
	mu sync.Mutex
	ring   []struct{val interface{}; i int64}
	stop, eos int32
}

func (p *prefetch) Get(worker int, open func()Stream) Stream {
	// this function is called once for every drain worker
	if p.width < 2 {
		return open()
	}
	// concurrent stream processing with p.width workers
	p.mu.Lock()
	if p.ring == nil {
		p.ring = make([]struct{val interface{}; i int64},p.width*p.width*PrefetchFactor)
		go func(stream Stream) {
			j := 0
			for atomic.LoadInt32(&p.stop) == 0 {
				v,_ := stream(true)
				t := 0
				for atomic.LoadInt32(&p.stop) == 0 {
					k := j % len(p.ring)
					if atomic.LoadInt64(&p.ring[k].i) == 0 {
						p.ring[k].val = v
						atomic.StoreInt64(&p.ring[k].i, int64(j+1))
						j++
						break
					} else {
						t++
						if t == p.width {
							runtime.Gosched()
							t = 0
						}
					}
				}
				if _, ok := v.(EndOfStream); ok {
					atomic.StoreInt32(&p.eos,1)
					break
				}
			}
			stream(false)
		}(open())
	}
	p.mu.Unlock()
	return p.substream(worker)
}

func (p *prefetch) substream(worker int) Stream {
	i := worker
	return func(next bool)(v interface{}, j int){
		if !next {
			atomic.StoreInt32(&p.stop,1)
			return EoS, 0
		}
		for {
			k := i % len(p.ring)
			if int64(i+1) == atomic.LoadInt64(&p.ring[k].i) {
				v, j = p.ring[k].val, i
				i += p.width
				atomic.StoreInt64(&p.ring[k].i,0)
				return
			}
			if atomic.LoadInt32(&p.eos) != 0 {
				return EoS, i
			}
			// switch to another goroutine
			runtime.Gosched()
		}
	}
}

func NoPrefetch(_ int,open func()Stream)Stream {
	return open()
}

