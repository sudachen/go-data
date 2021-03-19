package errors

import (
	"bytes"
	"fmt"
)

type errorPrinter struct {
	bytes.Buffer
	details bool
}

func (ep *errorPrinter) Print(args ...interface{}) {
	ep.Buffer.WriteString(fmt.Sprint(args...))
}

func (ep *errorPrinter) Printf(format string, args ...interface{}) {
	ep.Buffer.WriteString(fmt.Sprintf(format, args...))
}

func (ep errorPrinter) Detail() bool {
	return ep.details
}
