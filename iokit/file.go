package iokit

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sudachen.xyz/pkg/go-forge/errors"
)

func File(path string) StrictInputOutput {
	return StrictInputOutput{file(path)}
}

type file string

func expand(path string) (string, error) {
	if len(path) > 0 && path[0] == '$' {
		j := strings.IndexRune(path, '/')
		e := strings.ToLower(path[1:j])
		found := false
		for _, ev := range os.Environ() {
			k := strings.IndexRune(ev, '=')
			if k > 0 {
				ex := strings.ToLower(ev[:k])
				if ex == e {
					path = ev[k+1:] + path[j:]
					found = true
					break
				}
			}
		}
		if !found {
			return "", errors.New("can't expand path `" + path + "`")
		}
	}
	return path, nil
}

func (f file) Open() (io.ReadCloser, error) {
	path, err := expand(string(f))
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func (f file) Create() (Whole, error) {
	path, err := expand(string(f))
	if err != nil {
		return nil, err
	}
	dir, _ := filepath.Split(path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	x, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &whole{regular{x}}, nil
}

type regular struct {
	*os.File
}

func (f regular) Reset() error {
	_, err := f.File.Seek(0, 0)
	return err
}

func (f regular) Size() int64 {
	st, _ := f.File.Stat()
	return st.Size()
}

func (f regular) Fail() {
	fname := f.File.Name()
	_ = f.File.Truncate(0)
	_ = f.File.Close()
	_ = os.Remove(fname)
}
