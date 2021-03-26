package lazy

import (
	"sudachen.xyz/pkg/go-data/errors"
)

type Source func(...interface{}) Stream

func (zf Source) Open() Stream { return zf() }

type Stream func(bool) (interface{},int)

func (z Stream) Close() {
	z(false)
}

func (z Stream) Next() interface{} {
	v,_ := z(true)
	return v
}

type EndOfStream struct{ Err error }

func Fail(err error) interface{} { return EndOfStream{err} }

var EoS interface{} = EndOfStream{}
var NoValue interface{} = struct{}{}

var S Source = func(xs ...interface{}) Stream {
	for _, s := range xs {
		if f, ok := s.(func()Stream); ok {
			return f()
		}
	}
	return Error(errors.New("there is no stream provider function"))
}

func (zf Source) Link(xf Source) Source {
	return func(xs ...interface{}) Stream {
		return xf(func()Stream{ return zf(xs...) })
	}
}

