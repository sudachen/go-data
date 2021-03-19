package iokit

import (
	"errors"
	"io"
	"os"
	"strings"
)

func (iourl IoUrl) openUrlReader() (rd io.ReadCloser, err error) {
	if iourl.Cache.Exists() {
		return iourl.Cache.Open()
	}
	var f Whole
	if iourl.Cache.Defined() {
		f, err = File(iourl.Cache.Path() + "~").Create()
	} else {
		f, err = Tempfile("url-noncached-*")
	}
	defer func() {
		if f != nil {
			f.End()
		}
	}()
	if err = iourl.Download(f); err != nil {
		return
	}
	if err = f.Commit(); err != nil {
		return
	}
	if iourl.Cache.Defined() {
		// file was closed in Commit call
		if err = os.Rename(iourl.Cache.Path()+"~", iourl.Cache.Path()); err != nil {
			return
		}
		if rd, err = File(iourl.Cache.Path()).Open(); err != nil {
			return
		}
	} else {
		rd = f.(io.ReadWriteCloser)
		if _, err = rd.(io.Seeker).Seek(0, 0); err != nil {
			return
		}
		f = nil // do not close tempfile, it will be removed on close later
	}
	return
}

func Download(url string, writer io.Writer) error {
	j := strings.Index(url, "://")
	switch strings.ToLower(url[:j]) {
	case "http", "https":
		return HttpUrl(url).Download(writer)
	case "s3":
		return S3Url(url).Download(writer)
	case "gs":
		return GsUrl(url).Download(writer)
	}
	return errors.New("can't read from url `" + url + "`")
}

func (iourl IoUrl) Download(wr io.Writer) error {
	return Download(iourl.Url, wr)
}
