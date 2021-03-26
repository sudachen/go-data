package fu_tests

import (
	"gotest.tools/v3/assert"
	"sudachen.xyz/pkg/go-data/fu"
	"testing"
)

func Test_Atomic1(t *testing.T) {
	f := fu.AtomicFlag{1}
	assert.Assert(t, f.State() == true)
	f.Clear()
	assert.Assert(t, f.State() == false)
	f.Set()
	assert.Assert(t, f.State() == true)
	f.Clear()
	assert.Assert(t, f.State() == false)

	f = fu.AtomicFlag{0}
	assert.Assert(t, f.State() == false)
	f.Clear()
	assert.Assert(t, f.State() == false)
	f.Set()
	assert.Assert(t, f.State() == true)
}
