package lazy

import (
	"sudachen.xyz/pkg/go-data/errors"
	"sudachen.xyz/pkg/go-data/fu"
	"sync"
)

type Worker func(int, interface{}, error)error
type WorkerFactory func(concurrency int) []Worker

func (zf Source) Drain(wf WorkerFactory, concurrency ...int) error {
	switch ws := wf(fu.Fnzi(fu.Fnzi(concurrency...),1)); len(ws) {
	case 0:
		return errors.New("no workers provided for drain")
	case 1:
		z := zf.Open()
		defer z.Close()
		for {
			switch v, i := z(true); x := v.(type) {
			case struct{}:
			case EndOfStream:
				return fu.Fnze(ws[0](i, nil, x.Err), x.Err)
			default:
				if err := ws[0](i, v, nil); err != nil {
					return err
				}
			}
		}
	default:
		return concurrentDrain(zf, ws)
	}
}

func concurrentDrain(zf Source, ws []Worker) error {
	width := len(ws)
	wg := sync.WaitGroup{}
	wg.Add(width)
	errc := make(chan error,width)
	pf := &prefetch{width: width}
	for i:=0; i<width; i++ {
		go func(index int){
			defer wg.Done()
			z := zf(func()(int,int,Prefetch){
				return index, width, pf.Get
			})
		loop:
			for {
				switch v, i := z(true); x := v.(type) {
				case EndOfStream:
					if x.Err != nil {
						errc <- x.Err
					}
					_ = ws[index](i,NoValue,nil)
					break loop
				default:
					if err := ws[index](i,v,nil); err != nil {
						errc <- err
						break loop
					}
				}
			}
		}(i)
	}
	wg.Wait()
	close(errc)
	err, _ := <- errc
	return fu.Fnze(ws[0](0, nil, err),err)
}

func (zf Source) MustDrain(wf WorkerFactory, concurrency ...int) {
	if err := zf.Drain(wf, concurrency...); err != nil {
		panic(errors.Panic{Err: err})
	}
}

