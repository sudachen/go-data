package iokit_tests

import (
	"crypto/rand"
	"gotest.tools/assert"
	"sudachen.xyz/pkg/go-forge/fu"
	"testing"
)

func Test_Crypto(t *testing.T) {
	buf := make([]byte, 4096)
	rand.Read(buf)
	dat, err := fu.Encrypt("helloworld", buf)
	assert.NilError(t, err)
	dat, err = fu.Decrypt("helloworld", dat)
	assert.NilError(t, err)
	assert.DeepEqual(t, dat, buf)
}
