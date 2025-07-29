package stdlib

import (
	"bytes"
	gojson "encoding/json"

	"github.com/malivvan/vv/vvm"
	"github.com/malivvan/vv/vvm/stdlib/json"
)

var jsonModule = map[string]vvm.Object{
	"decode": &vvm.UserFunction{
		Name:  "decode",
		Value: jsonDecode,
	},
	"encode": &vvm.UserFunction{
		Name:  "encode",
		Value: jsonEncode,
	},
	"indent": &vvm.UserFunction{
		Name:  "encode",
		Value: jsonIndent,
	},
	"html_escape": &vvm.UserFunction{
		Name:  "html_escape",
		Value: jsonHTMLEscape,
	},
}

func jsonDecode(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		return nil, vvm.ErrWrongNumArguments
	}

	switch o := args[0].(type) {
	case *vvm.Bytes:
		v, err := json.Decode(o.Value)
		if err != nil {
			return &vvm.Error{
				Value: &vvm.String{Value: err.Error()},
			}, nil
		}
		return v, nil
	case *vvm.String:
		v, err := json.Decode([]byte(o.Value))
		if err != nil {
			return &vvm.Error{
				Value: &vvm.String{Value: err.Error()},
			}, nil
		}
		return v, nil
	default:
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}

func jsonEncode(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		return nil, vvm.ErrWrongNumArguments
	}

	b, err := json.Encode(args[0])
	if err != nil {
		return &vvm.Error{Value: &vvm.String{Value: err.Error()}}, nil
	}

	return &vvm.Bytes{Value: b}, nil
}

func jsonIndent(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 3 {
		return nil, vvm.ErrWrongNumArguments
	}

	prefix, ok := vvm.ToString(args[1])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "prefix",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	indent, ok := vvm.ToString(args[2])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "indent",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	switch o := args[0].(type) {
	case *vvm.Bytes:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, o.Value, prefix, indent)
		if err != nil {
			return &vvm.Error{
				Value: &vvm.String{Value: err.Error()},
			}, nil
		}
		return &vvm.Bytes{Value: dst.Bytes()}, nil
	case *vvm.String:
		var dst bytes.Buffer
		err := gojson.Indent(&dst, []byte(o.Value), prefix, indent)
		if err != nil {
			return &vvm.Error{
				Value: &vvm.String{Value: err.Error()},
			}, nil
		}
		return &vvm.Bytes{Value: dst.Bytes()}, nil
	default:
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}

func jsonHTMLEscape(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		return nil, vvm.ErrWrongNumArguments
	}

	switch o := args[0].(type) {
	case *vvm.Bytes:
		var dst bytes.Buffer
		gojson.HTMLEscape(&dst, o.Value)
		return &vvm.Bytes{Value: dst.Bytes()}, nil
	case *vvm.String:
		var dst bytes.Buffer
		gojson.HTMLEscape(&dst, []byte(o.Value))
		return &vvm.Bytes{Value: dst.Bytes()}, nil
	default:
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "bytes/string",
			Found:    args[0].TypeName(),
		}
	}
}
