package log_tests

import (
	"bytes"
	"gotest.tools/assert"
	"strings"
	"sudachen.xyz/pkg/go-forge/log"
	"testing"
)

func Test_Init(t *testing.T) {
	defer log.Config{Name: "test log", Verbose: true}.Init().Close()
	log.Info("hello logger!")
}

func Test_LogWriter(t *testing.T) {
	bf := bytes.Buffer{}
	func() {
		defer log.Config{Name: "test log", LogWriter: &bf}.Init().Close()
		log.Info("hello logger!")
	}()
	assert.Assert(t, strings.Contains(bf.String(), "hello logger!"))
}
