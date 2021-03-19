package iokit

import (
	"io"
	"net/http"
)

type HttpUrl string

func (url HttpUrl) Download(wr io.Writer) error {
	resp, err := http.Get(string(url))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(wr, resp.Body)
	return err
}
