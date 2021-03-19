package dataframe

import "sudachen.xyz/pkg/go-forge/lazy"

type Column struct { Sequence }

func Col(interface{}) Column {
	return Column{ nil }
}

func (c Column) Lazy() lazy.Source {
	return func()lazy.Stream {
		clen := lazy.Index(c.Len())
		return func(index lazy.Index)interface{}{
			if index < clen {
				return c.Index(int(index)).Val
			}
			return lazy.EndOfStream{}
		}
	}
}
