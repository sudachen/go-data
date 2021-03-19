package dataframe

type Frame interface {
	Len() int
	Names() []string
	Column(column string) Sequence
	Row(index int) *Row
}

type Sequence interface {
	Len() int
	Index(index int) Cell
	Na(index int) bool
	Copy(to interface{}, offset,length int)
}

