package dataframe


type Row struct {
	Factory RowFactory
	data []Cell
}

func (r *Row) Recycle() {
	r.Factory.Recycle(r)
}

func (r *Row) Len() int {
	return len(r.data)
}

func (r *Row) Index(name string) (int,bool) {
	return r.Factory.Index(name)
}

func (r *Row) ByName(name string) (*Cell,bool) {
	if i, ok := r.Index(name); ok {
		return &r.data[i], true
	}
	return nil, false
}

func (r *Row) ByIndex(index int) (*Cell,bool) {
	if index >= 0 && index < len(r.data) {
		return &r.data[index], true
	}
	return nil, false
}

