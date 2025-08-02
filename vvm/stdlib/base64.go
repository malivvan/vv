package stdlib

import (
	"encoding/base64"

	"github.com/malivvan/vv/vvm"
)

var base64Module = map[string]vvm.Object{
	"encode": &vvm.BuiltinFunction{
		Value: FuncAYRS(base64.StdEncoding.EncodeToString),
	},
	"decode": &vvm.BuiltinFunction{
		Value: FuncASRYE(base64.StdEncoding.DecodeString),
	},
	"raw_encode": &vvm.BuiltinFunction{
		Value: FuncAYRS(base64.RawStdEncoding.EncodeToString),
	},
	"raw_decode": &vvm.BuiltinFunction{
		Value: FuncASRYE(base64.RawStdEncoding.DecodeString),
	},
	"url_encode": &vvm.BuiltinFunction{
		Value: FuncAYRS(base64.URLEncoding.EncodeToString),
	},
	"url_decode": &vvm.BuiltinFunction{
		Value: FuncASRYE(base64.URLEncoding.DecodeString),
	},
	"raw_url_encode": &vvm.BuiltinFunction{
		Value: FuncAYRS(base64.RawURLEncoding.EncodeToString),
	},
	"raw_url_decode": &vvm.BuiltinFunction{
		Value: FuncASRYE(base64.RawURLEncoding.DecodeString),
	},
}
