package lazy

import (
	"sudachen.xyz/pkg/go-forge/errors"
	"sudachen.xyz/pkg/go-forge/fu"
)

func (zf Source) Drain(sink func(interface{}) error) (err error) {
	z := zf()
	var i Index
loop:
	for {
		v := z(i)
		i++
		if v != NoValue {
			switch x := v.(type) {
			case Fail:
				err = x.Err
				break loop
			case EndOfStream:
				break loop
			default:
				if err = sink(v); err != nil {
					break loop
				}
			}
		}
	}
	z(CloseSource)
	if err != nil {
		return fu.Fnze(err,sink(DrainFailed))
	}
	return sink(DrainSucceed)
}

func (zf Source) MustDrain(sink func(interface{}) error) {
	if err := zf.Drain(sink); err != nil {
		panic(errors.PanicBtrace{ Err: err})
	}
}
