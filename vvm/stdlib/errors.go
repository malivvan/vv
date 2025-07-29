package stdlib

import (
	"github.com/malivvan/vv/vvm"
)

func wrapError(err error) vvm.Object {
	if err == nil {
		return vvm.TrueValue
	}
	return &vvm.Error{Value: &vvm.String{Value: err.Error()}}
}
