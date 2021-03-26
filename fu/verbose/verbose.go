package verbose

import (
	"fmt"
	"sudachen.xyz/pkg/go-data/log"
)

type VerboseKind int

const (
	Silent VerboseKind = iota
	Print
	Log
)

var Verbose = Silent

func Markup() string {
	return "### "
}

func Println(a ...interface{}) {
	switch Verbose {
	case Print:
		fmt.Println(append([]interface{}{Markup()}, a...))
	case Log:
		log.Info(a...)
	}
}

func Printf(f string, a ...interface{}) {
	switch Verbose {
	case Print:
		fmt.Printf(f+"\n", a...)
	case Log:
		log.Infof(f, a...)
	}
}

func BeVerbose(kind VerboseKind) (old VerboseKind) {
	old = Verbose
	Verbose = kind
	return
}

func (old VerboseKind) Revert() {
	Verbose = old
}
