package iokit

import (
	"io"
	"os"
	"strings"
	"sudachen.xyz/pkg/go-data/errors"
)

type urlwriter struct {
	Whole
	iourl IoUrl
}

func (iourl IoUrl) createUrlWriter() (wh Whole, err error) {
	var f Whole
	if iourl.Cache.Exists() {
		if err = iourl.Cache.Remove(); err != nil {
			return nil, errors.Wrap(err, "can't delete existing cache file")
		}
	}
	if iourl.Cache.Defined() {
		f, err = File(iourl.Cache.Path() + "~").Create()
	} else {
		f, err = Tempfile("url-noncached-*")
	}
	wh = &urlwriter{f, iourl}
	return
}

func (uw *urlwriter) Commit() (err error) {
	var rd io.ReadCloser
	if err = uw.Whole.Commit(); err != nil {
		return
	}
	if uw.iourl.Cache.Defined() {
		// file was closed in Commit call
		if err = os.Rename(uw.iourl.Cache.Path()+"~", uw.iourl.Cache.Path()); err != nil {
			return
		}
		if rd, err = File(uw.iourl.Cache.Path()).Open(); err != nil {
			return
		}
	} else {
		rd = uw.Whole.(io.ReadCloser)
		if _, err = rd.(io.Seeker).Seek(0, 0); err != nil {
			return
		}
	}
	defer rd.Close()
	uw.Whole = nil
	err = uw.iourl.Upload(rd)
	return
}

func (uw *urlwriter) End() {
	if uw.Whole != nil {
		uw.Whole.End()
	}
}

func Upload(url string, reader io.Reader) error {
	j := strings.Index(url, "://")
	switch strings.ToLower(url[:j]) {
	case "s3":
		return S3Url(url).Upload(reader)
	case "gs":
		return GsUrl(url).Upload(reader)
	}
	return errors.New("can't read from url `" + url + "`")
}

func (iourl IoUrl) Upload(rd io.Reader) error {
	return Upload(iourl.Url, rd)
}
