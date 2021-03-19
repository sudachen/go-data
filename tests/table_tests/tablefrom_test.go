package table_tests

import (
	"gotest.tools/assert"
	"sudachen.xyz/pkg/go-forge/dataframe"
	"sudachen.xyz/pkg/go-forge/lazy"
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
	var x dataframe.Table
	lazy.Chan(c).Map1(dataframe.StructToRow).MustDrain(dataframe.Sink(&x))
	assert.Assert(t, x.Len() == len(colors))
	for i, v := range colors {
		assert.Assert(t, x.Column("Color").Index(i).String() == v.Color)
		assert.Assert(t, x.Column("Index").Index(i).Int() == v.Index)
	}
}

func Test1_TableFromList(t *testing.T) {
	var x dataframe.Table
	lazy.List(colors).Map1(dataframe.StructToRow).MustDrain(dataframe.Sink(&x))
	assert.Assert(t, x.Len() == len(colors))
	for i, v := range colors {
		assert.Assert(t, x.Column("Color").Index(i).String() == v.Color)
		assert.Assert(t, x.Column("Index").Index(i).Int() == v.Index)
	}
}
