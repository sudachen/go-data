package iokit_tests

import (
	"fmt"
	"gotest.tools/v3/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"sudachen.xyz/pkg/go-data/iokit"
	"testing"
)

/*
	S3 tests use S3_AWS_TEST_URL and S3_DO_TEST_URL environment variables
	S3_name_URL = s3://key:secret@region.endpoint/bucket/prefix
*/

func Test_S3Path1(t *testing.T) {
	if len(os.Getenv("do_test")) == 0 {
		t.Skip("envirnomet variable 'do_test' does not set")
	}
	if len(os.Getenv("aws_test")) == 0 {
		t.Skip("envirnomet variable 'aws_test' does not set")
	}
	for i, url := range []string{
		"s3://$do_test/go-iokit/test_s3path.txt",
		"s3://$aws_test/go-iokit/test_s3path.txt"} {

		S := fmt.Sprintf(`Hello world! %d`, i)
		wh := iokit.Url(url).MustCreate()
		defer wh.End()
		wh.MustWrite([]byte(S))
		wh.MustCommit()
		rd := iokit.Url(url).MustOpen()
		defer rd.Close()
		q := rd.MustReadAll()
		assert.Assert(t, string(q) == S)
	}
}

func Test_S3Path2(t *testing.T) {
	if len(os.Getenv("do_test")) == 0 {
		t.Skip("envirnomet variable 'do_test' does not set")
	}
	if len(os.Getenv("aws_test")) == 0 {
		t.Skip("envirnomet variable 'aws_test' does not set")
	}
	for i, url := range []string{
		"s3://$do_test/go-iokit/test_s3path.txt",
		"s3://$aws_test/go-iokit/test_s3path.txt"} {

		S := fmt.Sprintf(`Hello world! %d`, i)
		file := iokit.Url(url, iokit.Cache("go-iokit/test_s3path.txt"))
		wh, err := file.Create()
		assert.NilError(t, err)
		defer wh.End()
		_, err = wh.Write([]byte(S))
		assert.NilError(t, err)
		err = wh.Commit()
		assert.NilError(t, err)
		wh.End()
		rd, err := file.Open()
		assert.NilError(t, err)
		defer rd.Close()
		q, err := ioutil.ReadAll(rd)
		assert.NilError(t, err)
		assert.Assert(t, string(q) == S)
		rd.Close()
	}
	for i, url := range []string{
		"s3://$do_test/go-iokit/test_s3path.txt",
		"s3://$aws_test/go-iokit/test_s3path.txt"} {
		S := fmt.Sprintf(`Hello world! %d`, i)
		rd, err := iokit.Url(url).Open()
		assert.NilError(t, err)
		defer rd.Close()
		q, err := ioutil.ReadAll(rd)
		assert.NilError(t, err)
		assert.Assert(t, string(q) == S)
	}
}

func Test_Example1(t *testing.T) {
	if len(os.Getenv("do_test")) == 0 {
		t.Skip("envirnomet variable 'do_test' does not set")
	}
	url := "s3://$do_test/go-iokit/test_s3path.txt"
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

func Test_Example2(t *testing.T) {
	if len(os.Getenv("do_test")) == 0 {
		t.Skip("envirnomet variable 'do_test' does not set")
	}
	url := "s3://$do_test/go-iokit/test_s3path.txt"
	S := fmt.Sprintf(`Hello world! %d`, rand.Int())
	wh, err := iokit.Url(url).Create()
	assert.NilError(t, err)
	defer wh.End()
	_, err = wh.Write([]byte(S))
	assert.NilError(t, err)
	err = wh.Commit()
	assert.NilError(t, err)
	rd, err := iokit.Url(url).Open()
	assert.NilError(t, err)
	defer rd.Close()
	q, err := ioutil.ReadAll(rd)
	assert.NilError(t, err)
	assert.Assert(t, string(q) == S)
}
