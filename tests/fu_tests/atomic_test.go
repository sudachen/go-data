package fu_tests

import (
	"gotest.tools/assert"
	"sudachen.xyz/pkg/go-forge/fu"
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
