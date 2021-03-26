package log_tests

import (
	"bytes"
	"gotest.tools/v3/assert"
	"strings"
	"sudachen.xyz/pkg/go-data/log"
	"testing"
)

func Test_Init(t *testing.T) {
	defer log.Config{Name: "test log", Verbose: true}.Init().Close()
	log.Info("hello log!")
}

func Test_LogWriter(t *testing.T) {
	bf := bytes.Buffer{}
	func() {
		defer log.Config{Name: "test log", LogWriter: &bf}.Init().Close()
		log.Info("hello log!")
	}()
	assert.Assert(t, strings.Contains(bf.String(), "hello log!"))
}
