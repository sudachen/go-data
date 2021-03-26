package lazy

func Error(err error, z ...Stream) Stream {
	return func(next bool) (interface{},int) {
		if !next {
			for _, zz := range z {
				zz.Close()
			}
		}
		return EndOfStream{err}, 0
	}
}

func Wrap(e interface{}) Stream {
	if stream, ok := e.(Stream); ok {
		return stream
	} else {
		return Error(e.(error))
	}
}

func ErrorSource(e error) Source {
	return func(...interface{}) Stream {
		return Error(e)
	}
}

func ErrorSink(e error) WorkerFactory {
	return func(int)[]Worker{
		return []Worker{func(int, interface{}, error)error{
			return e
		}}
	}
}
