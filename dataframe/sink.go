package dataframe

import "sudachen.xyz/pkg/go-forge/lazy"

//
// var t Table
// source.MustDrain(dataframe.Sink(&t))
//


func Sink(t *Table, preserve ...int) func(interface{})error{
	var fr *varframe
	return func(v interface{})error{
		if v == lazy.DrainSucceed {
			*t = Table{fr}
		} else if v != lazy.DrainFailed {
			// fill frame
			r := v.(*Row)
			//
			r.Recycle()
		}
		return nil
	}
}
