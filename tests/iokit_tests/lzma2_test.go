package iokit_tests

import (
	"fmt"
	"gotest.tools/v3/assert"
	"math/rand"
	"sudachen.xyz/pkg/go-data/iokit"
	"testing"
)

func Test_Lzma2(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := iokit.Lzma2(iokit.Cache("test.lzma2").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		w.MustCommit()
	}()
	s := func() string {
		r := iokit.Compressed(iokit.Cache("test.lzma2").File()).MustOpen()
		defer r.Close()
		return string(r.MustReadAll())
	}()
	assert.Assert(t, s == S)
}
