package dataframe

import (
	"sudachen.xyz/pkg/go-forge/lazy"
)

type Table struct { Frame }

func (t Table) Column(name string) Column {
	return Column{t.Frame.Column(name)}
}

func (t Table) Lazy() lazy.Source {
	return func()lazy.Stream {
		clen := lazy.Index(t.Len())
		return func(index lazy.Index)interface{}{
			if index < clen {
				return t.Row(int(index))
			}
			return lazy.EndOfStream{}
		}
	}
}

