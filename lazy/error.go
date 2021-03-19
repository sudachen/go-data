package lazy

func Error(err error, z ...Stream) Stream {
	return func(dx Index)interface{} {
		if dx == CloseSource && len(z) > 0 {
			for _,zz := range z{
				zz(CloseSource)
			}
		}
		return err
	}
}

func Wrap(e interface{}) Stream {
	if stream, ok := e.(Stream); ok {
		return stream
	} else {
		return Error(e.(error))
	}
}
