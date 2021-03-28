package errors

import (
	"golang.org/x/xerrors"
	"strings"
)

type Panic struct{ Err error }

func (x Panic) stringify(indepth bool) string {
	s, e := stringifyError(x.Err)
	ns := []string{s}
	for e != nil && indepth {
		s, e = stringifyError(e)
		ns = append(ns, s)
	}
	return strings.Join(ns, "\n")
}

func (x Panic) Error() string {
	return x.stringify(false)
}

func (x Panic) String() string {
	return x.stringify(true)
}

func (x Panic) Unwrap() error {
	if w, ok := x.Err.(xerrors.Wrapper); ok {
		return w.Unwrap()
	}
	return x.Err
}
