package lazy

import (
	"sudachen.xyz/pkg/go-forge/errors"
)

func (zf Source) GetOne() (ret interface{}, err error) {
	err = zf.First(1).Drain(func(v interface{})error {
		if v != DrainFailed && v != DrainSucceed {
			ret = v
		}
		return nil
	})
	return
}

func (zf Source) MustGetOne() interface{} {
	x,err := zf.GetOne();
	if err != nil {
		panic(errors.PanicBtrace{Err: err})
	}
	return x
}

