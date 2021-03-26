package lazy

func (zf Source) First(n int) Source {
	return func(xs ...interface{}) Stream {
		z := zf(xs)
		count := 0
		return func(next bool) (interface{},int) {
			if count >= n {
				return EoS,0
			}
			v,i := z(next)
			switch v.(type) {
			case EndOfStream,struct{}:
			default:
				count++
			}
			return v,i
		}
	}
}
