package iokit_tests

import (
	"fmt"
	"gotest.tools/v3/assert"
	"math/rand"
	"os"
	"sudachen.xyz/pkg/go-data/iokit"
	"testing"
)

/*
	GS tests use GS_URL environment variables
	GS_ENCTEST_URL = gs://bucket/prefix:password:/abspath/credential.json.enc
*/

func Test_GsPath1(t *testing.T) {
	if len(os.Getenv("enctest")) == 0 {
		t.Skip("envirnomet variable 'enctest' does not set")
	}
	url := "gs://$enctest/test_gspath1.txt"
	S := fmt.Sprintf(`Hello world! %d`, rand.Int())
	wh := iokit.Url(url).MustCreate()
	defer wh.End()
	wh.MustWrite([]byte(S))
	wh.MustCommit()
	rd := iokit.Url(url).MustOpen()
	defer rd.Close()
	q := rd.MustReadAll()
	assert.Assert(t, string(q) == S)
}
