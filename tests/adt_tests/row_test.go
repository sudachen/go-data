package adt_tests

import (
	"fmt"
	"gotest.tools/v3/assert"
	"reflect"
	"strings"
	"sudachen.xyz/pkg/go-data/adt"
	"sudachen.xyz/pkg/go-data/lazy"
	"testing"
)

func structToString(st interface{}) string {
	tv := reflect.ValueOf(st)
	tp := tv.Type()
	var ss []string
	for i := 0; i < tp.NumField(); i++ {
		ss = append(ss, fmt.Sprintf("%s: %v(%v)", tp.Field(i).Name, tp.Field(i).Type, tv.Field(i).Interface()))
	}
	return strings.Join(ss, ", ")
}

func structToRowString(st interface{}) string {
	return "Row{" + structToString(st) + "}"
}

func Test1_RowString(t *testing.T) {
	w, err := adt.NewWrapper(Color{})
	assert.NilError(t, err)
	for _, v := range colors {
		r, err := w.Wrap(v)
		assert.NilError(t, err)
		assert.Assert(t, r.String() == structToRowString(v))
	}
}

func Test1_CollectRowString(t *testing.T) {
	var r []*adt.Row
	lazy.List(colors).Map1(adt.StructToRow).MustCollect(&r)
	for i, v := range colors {
		assert.Assert(t, r[i].String() == structToRowString(v))
	}
}

func Test1_CollectAnyRowString(t *testing.T) {
	r := lazy.List(colors).Map1(adt.StructToRow).MustCollectAny()
	assert.Assert(t, reflect.TypeOf(r) == reflect.TypeOf([]*adt.Row{}))
	for i, v := range colors {
		assert.Assert(t, r.([]*adt.Row)[i].String() == structToRowString(v))
	}
}

func Test1_CollectValueToRow(t *testing.T) {
	var r []*adt.Row
	lazy.List(colors).Map(func(c Color)int{return c.Index}).Map1(adt.ValueToRow("Index")).MustCollect(&r)
	for i, v := range colors {
		assert.Assert(t, r[i].Col("Index").Int() == v.Index)
	}
}
