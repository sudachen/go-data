package dataframe

import (
	"sort"
	"sync"
)

type RowFactory interface {
	New() *Row
	Recycle(*Row)
	Index(n string) (int, bool)
	With(with []string) RowFactory
	Except(except []string) RowFactory
}

type rowFactory struct {
	pool sync.Pool
	names []string
}

func NewRowFactory(names []string) RowFactory{
	ns := make([]string,len(names))
	copy(ns,names)
	sort.Strings(ns)
	return &rowFactory{pool: sync.Pool{}, names: ns}
}

func (f *rowFactory) New() *Row {
	if x := f.pool.Get(); x != nil {
		return x.(*Row)
	}
	return &Row{f, make([]Cell,len(f.names))}
}

func (f *rowFactory) With(with []string) RowFactory {
	ns := make([]string,len(f.names),len(with))
	copy(ns,f.names)
	for _,n := range with {
		j := sort.SearchStrings(ns,n)
		if ns[j] != n {
			k := len(ns)
			ns = append(ns,n)
			copy(ns[j+1:],ns[j:k])
			ns[j] = n
		}
	}
	return &rowFactory{pool: sync.Pool{}, names: ns}
}

func (f *rowFactory) Except(except []string) RowFactory {
	ns := make([]string,0,len(f.names))
loop:
	for _,n := range f.names {
		for _,e := range except {
			if e == n {
				continue loop
			}
		}
		ns = append(ns,n)
	}
	return &rowFactory{pool: sync.Pool{}, names: ns}
}

func (f *rowFactory) Recycle(r *Row) {
	for i := range r.data {
		r.data[i] = Cell{nil}
	}
	f.pool.Put(r)
}

func (f *rowFactory) Index(n string) (int,bool) {
	for i,s := range f.names {
		if s == n {
			return i, true
		}
	}
	return -1,false
}