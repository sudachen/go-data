package dataframe

import (
	"reflect"
	"sudachen.xyz/pkg/go-forge/fu"
	"sync"
)

const maxVarpartLength = 32*1024/8-16

type varpart struct {
	columns []reflect.Value
	na      []fu.Bits
	offset,length  int
}

type varframe struct {
	factory rowFactory
	parts []*varpart
	length int
	pool sync.Pool
}

type varseq struct {
	fr *varframe
	colndx int
}

func (fr *varframe) locate(row int) (*varpart,int) {
	return nil,0
}

func (fr *varframe) Len() int {
	return fr.length
}

func (fr *varframe) Names() []string {
	return fr.factory.names
}

func (fr *varframe) Column(column string) Sequence {
	colndx := 0
	return &varseq{fr, colndx}
}

func (fr *varframe) Row(index int) *Row {
	row := fr.factory.New()
	part, ndx := fr.locate(index)
	for i,c := range part.columns {
		if !part.na[i].Bit(ndx) {
			row.data[i].Val = c.Index(i).Interface()
		}
	}
	return row
}

func (sq *varseq) Len() int {
	return sq.fr.length
}

func (sq *varseq) Index(index int) Cell {
	part,ndx := sq.fr.locate(index)
	if !part.na[sq.colndx].Bit(ndx) {
		return Cell{part.columns[sq.colndx].Index(ndx).Interface()}
	}
	return Cell{}
}

func (sq *varseq) Na(index int) bool {
	part,ndx := sq.fr.locate(index)
	return part.na[sq.colndx].Bit(ndx)
}

func (sq *varseq) Copy(to interface{}, offset,length int) {
	v := reflect.ValueOf(to)
	for i := 0; i < length; i++ {
		part,ndx := sq.fr.locate(offset+i)
		if !part.na[sq.colndx].Bit(ndx) {
			v.Index(i).Set(part.columns[sq.colndx].Index(ndx))
		}
	}
}

