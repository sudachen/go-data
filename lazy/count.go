package lazy

import "sudachen.xyz/pkg/go-forge/errors"

func (zf Source) Count() (count int, err error) {
	err = zf.Drain(func(v interface{}) error {
		if v != DrainFailed && v != DrainSucceed { count++ }
		return nil
	})
	return
}

func (zf Source) MustCount() int{
	count, err := zf.Count()
	if err != nil { panic(errors.PanicBtrace{err}) }
	return count
}
