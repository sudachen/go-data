package iokit_tests

import (
	"fmt"
	"gotest.tools/assert"
	"math/rand"
	"sudachen.xyz/pkg/go-forge/iokit"
	"testing"
)

func Test_Zip(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := iokit.Zip("test.txt", iokit.Cache("test.zip").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		w.MustCommit()
	}()
	s := func() string {
		r := iokit.ZipFile("test.txt", iokit.Cache("test.zip").File()).MustOpen()
		defer r.Close()
		return string(r.MustReadAll())
	}()
	assert.Assert(t, s == S)
}

func Test_Zip_Fail(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := iokit.Zip("test.txt", iokit.Cache("test.zip").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		// no commit here
	}()
	_, err := iokit.ZipFile("test.txt", iokit.Cache("test.zip").File()).Open()
	assert.Assert(t, err != nil)
}
