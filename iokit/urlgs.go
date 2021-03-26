package iokit

import (
	"io"
	"sudachen.xyz/pkg/go-data/iokit/gcp"
)

type GsUrl string

func (gsurl GsUrl) Download(wr io.Writer) error {
	return gcp.Download(string(gsurl), wr)
}

func (gsurl GsUrl) Upload(rd io.Reader, metadata ...map[string]string) error {
	mdp := map[string]string(nil)
	if len(metadata) > 0 {
		mdp = metadata[0]
	}
	return gcp.Upload(string(gsurl), rd, mdp)
}
