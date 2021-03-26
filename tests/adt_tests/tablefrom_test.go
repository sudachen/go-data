package adt_tests

import (
	"gotest.tools/v3/assert"
	"sudachen.xyz/pkg/go-data/adt"
	"sudachen.xyz/pkg/go-data/lazy"
	"testing"
)

func Test1_TableFromChannel(t *testing.T) {
	c := make(chan Color)
	go func() {
		for _, x := range colors {
			c <- x
		}
		close(c)
	}()
	var x adt.Table
	lazy.Chan(c).Map1(adt.StructToRow).MustDrain(x.Sink())
	assert.Assert(t, x.Len() == len(colors))
	for i, v := range colors {
		assert.Assert(t, x.Col("Color").At(i).String() == v.Color)
		assert.Assert(t, x.Col("Index").At(i).Int() == v.Index)
		assert.Assert(t, x.Row(i).String() == structToRowString(v))
	}
}

func Test1_TableFromList(t *testing.T) {
	var x adt.Table
	lazy.List(colors).Map1(adt.StructToRow).MustDrain(x.Sink())
	assert.Assert(t, x.Len() == len(colors))
	for i, v := range colors {
		assert.Assert(t, x.Col("Color").At(i).String() == v.Color)
		assert.Assert(t, x.Col("Index").At(i).Int() == v.Index)
		assert.Assert(t, x.Row(i).String() == structToRowString(v))
	}
}

func Test2_TableFromList(t *testing.T) {
	var x adt.Table
	lazy.List(colors).Map1(adt.StructToRow).MustDrain(x.Sink())
	assert.Assert(t, x.Len() == len(colors))
	for i, v := range colors {
		assert.Assert(t, x.Col("Color").At(i).String() == v.Color)
		assert.Assert(t, x.Col("Index").At(i).Int() == v.Index)
		assert.Assert(t, x.Row(i).String() == structToRowString(v))
	}
}

func Test1_ConcurrentTableSink(t *testing.T) {
	var x adt.Table
	lazy.List(colors_N).
		Map1(adt.StructToRow).
		MustDrain(x.Sink(),8)
	assert.Assert(t, x.Len() == len(colors_N))
	for i, v := range colors_N {
		//fmt.Println(x.Row(i), x.Col("Color").At(i).String(), v.Color)
		assert.Assert(t, x.Col("Color").At(i).String() == v.Color)
		assert.Assert(t, x.Col("Index").At(i).Int() == v.Index)
	}
}
