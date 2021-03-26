package iokit

import (
	"io"
	"sudachen.xyz/pkg/go-data/iokit/s3p"
)

type S3Url string

func (s3url S3Url) Download(wr io.Writer) error {
	return s3p.Download(string(s3url), wr.(io.WriterAt))
}

func (s3url S3Url) Upload(rd io.Reader, metadata ...map[string]string) error {
	mdp := map[string]string(nil)
	if len(metadata) > 0 {
		mdp = metadata[0]
	}
	return s3p.Upload(string(s3url), rd, mdp)
}
