package stdlib

import (
	"encoding/hex"
	"github.com/malivvan/vv/vvm"
)

var hexModule = map[string]vvm.Object{
	"encode": &vvm.BuiltinFunction{Value: FuncAYRS(hex.EncodeToString)},
	"decode": &vvm.BuiltinFunction{Value: FuncASRYE(hex.DecodeString)},
}
