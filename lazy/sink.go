package lazy

import "sudachen.xyz/pkg/go-data/fu"

func Sink(f func(interface{}, error)error) WorkerFactory {
	return func(w int) []Worker {
		if w == 1 {
			return []Worker{ func(_ int, v interface{}, err error) error{
				return f(v,err)
			}}
		}
		return synchronizedSink(w,f)
	}
}

func synchronizedSink(w int,f func(interface{},error)error) (r []Worker) {
	r = make([]Worker,w)
	c := make([]chan interface{},w)
	e := make(chan struct{})
	errc := make(chan error,1)
	for i := range r {
		n := i
		c[n] = make(chan interface{})
		r[n] = func(_ int, v interface{}, err error) error{
			if v == nil {
				if n == 0 {
					close(e)
					return fu.Fnze(f(nil,err),<-errc)
				}
			} else {
				select {
				case c[n] <- v:
				case <-e:
				case err := <-errc:
					return err
				}
			}
			return nil
		}
	}
	go func(){
		defer close(errc)
		for {
			for i := range c {
				select {
				case <-e:
					return
				case v := <- c[i]:
					switch v.(type) {
					case struct{}:
					default:
						if err := f(v, nil); err != nil {
							errc <- err
							return
						}
					}
				}
			}
		}
	}()
	return
}
