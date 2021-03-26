package iokit_tests

import (
	"fmt"
	"gotest.tools/v3/assert"
	"math/rand"
	"sudachen.xyz/pkg/go-data/iokit"
	"testing"
)

func Test_Gzip(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := iokit.Gzip(iokit.Cache("test.gz").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		w.MustCommit()
	}()
	s := func() string {
		r := iokit.Compressed(iokit.Cache("test.gz").File()).MustOpen()
		defer r.Close()
		return string(r.MustReadAll())
	}()
	assert.Assert(t, s == S)
}

func Test_Gzip_Fail(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := iokit.Gzip(iokit.Cache("test.gz").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		// no commit here
	}()
	_, err := iokit.Compressed(iokit.Cache("test.gz").File()).Open()
	assert.Assert(t, err != nil)
}
