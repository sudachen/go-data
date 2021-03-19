package errors

import (
	"golang.org/x/xerrors"
	"strings"
)

type PanicBtrace struct{ Err error }

func (x PanicBtrace) stringify(indepth bool) string {
	s, e := stringifyError(x.Err)
	ns := []string{s}
	for e != nil && indepth {
		s, e = stringifyError(e)
		ns = append(ns, s)
	}
	return strings.Join(ns, "\n")
}

func (x PanicBtrace) Error() string {
	return x.stringify(false)
}

func (x PanicBtrace) String() string {
	return x.stringify(true)
}

func (x PanicBtrace) Unwrap() error {
	if w, ok := x.Err.(xerrors.Wrapper); ok {
		return w.Unwrap()
	}
	return x.Err
}
