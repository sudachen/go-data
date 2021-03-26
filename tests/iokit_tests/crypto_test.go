package iokit_tests

import (
	"crypto/rand"
	"gotest.tools/v3/assert"
	"sudachen.xyz/pkg/go-data/fu"
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
