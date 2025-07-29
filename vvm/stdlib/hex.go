package stdlib

import (
	"encoding/hex"
	"github.com/malivvan/vv/vvm"
)

var hexModule = map[string]vvm.Object{
	"encode": &vvm.UserFunction{Value: FuncAYRS(hex.EncodeToString)},
	"decode": &vvm.UserFunction{Value: FuncASRYE(hex.DecodeString)},
}
