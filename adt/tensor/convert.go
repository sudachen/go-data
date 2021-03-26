package tensor

import (
	"golang.org/x/xerrors"
	"reflect"
	"sudachen.xyz/pkg/go-data/adt"
	"sudachen.xyz/pkg/go-data/errors"
	"sudachen.xyz/pkg/go-data/fu"
)

type Xtensor struct{ T reflect.Type }

func (t Xtensor) Type() reflect.Type {
	return fu.TensorType
}

func (t Xtensor) Convert(value string, data *interface{}, index, _ int) (err error) {
	*data, err = DecodeTensor(value)
	return
}

func tensorOf(x interface{}, tp reflect.Type, width int) (adt.Tensor, error) {
	if x != nil {
		return x.(adt.Tensor), nil
	}
	switch tp {
	case fu.Float32:
		return MakeFloat32Tensor(1, 1, width, nil), nil
	//case fu.Fixed8Type:
	//	return MakeFixed8Tensor(1, 1, width, nil), nil
	//case fu.Int:
	//	return MakeIntTensor(1, 1, width, nil), nil
	case fu.Byte:
		return MakeByteTensor(1, 1, width, nil), nil
	default:
		return adt.Tensor{}, xerrors.Errorf("unknown tensor value type " + tp.String())
	}
}

func (t Xtensor) ConvertElm(value string, data *interface{}, index, width int) (err error) {
	z, err := tensorOf(*data, t.T, width)
	if err != nil {
		return
	}
	if *data == nil {
		*data = z
	}
	return z.ConvertElem(value, index)
}

func (Xtensor) Format(x interface{}) (string,error) {
	if x == nil {
		return "", nil
	}
	if tz, ok := x.(adt.Tensor); ok {
		return tz.String(), nil
	}
	return "", errors.Errorf("`%v` is not a tensor value", x)
}

