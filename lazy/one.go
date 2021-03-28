package lazy

import (
	"sudachen.xyz/pkg/go-data/errors"
)

func (zf Source) GetOne() (ret interface{}, err error) {
	err = zf.First(1).Drain(func(int)[]Worker{
		return []Worker{func(_ int, v interface{}, _ error) (_ error) {
			if v != nil {
				ret = v
			}
			return
		}}
	})
	return
}

func (zf Source) MustGetOne() interface{} {
	x, err := zf.GetOne()
	if err != nil {
		panic(errors.Panic{Err: err})
	}
	return x
}
