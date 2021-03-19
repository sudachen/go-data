package lazy

import (
	"math"
	"sudachen.xyz/pkg/go-forge/errors"
)

type Index = uint64
type Stream func(Index) interface{}
type Source func() Stream
type Sink func(interface{}) error

const iniCollectLength = 13
const CloseSource = Index(math.MaxUint64)

var NoValue interface{} = nil
type EndOfStream struct{}
type Fail struct { Err error }

type LazyValue interface {
	Recycle()
}

var DrainSucceed interface{} = EndOfStream{}
var DrainFailed interface{} = Fail{errors.New("drain failed")}
