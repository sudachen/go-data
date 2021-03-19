package dataframe

type StruCtx struct {}

func NewCtx(stru interface{}) *StruCtx {
	return &StruCtx{}
}

func (*StruCtx) Wrap(stru interface{}) *Row {
	return nil
}

func StructToRow(stru interface{}) interface{} {
	ctx := NewCtx(stru)
	return func(x interface{})*Row {
		return ctx.Wrap(stru)
	}
}


