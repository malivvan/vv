package stdlib

import (
	"math/rand"

	"github.com/malivvan/vv/vvm"
)

var randModule = map[string]vvm.Object{
	"int": &vvm.BuiltinFunction{
		Name:  "int",
		Value: FuncARI64(rand.Int63),
	},
	"float": &vvm.BuiltinFunction{
		Name:  "float",
		Value: FuncARF(rand.Float64),
	},
	"intn": &vvm.BuiltinFunction{
		Name:  "intn",
		Value: FuncAI64RI64(rand.Int63n),
	},
	"exp_float": &vvm.BuiltinFunction{
		Name:  "exp_float",
		Value: FuncARF(rand.ExpFloat64),
	},
	"norm_float": &vvm.BuiltinFunction{
		Name:  "norm_float",
		Value: FuncARF(rand.NormFloat64),
	},
	"perm": &vvm.BuiltinFunction{
		Name:  "perm",
		Value: FuncAIRIs(rand.Perm),
	},
	"seed": &vvm.BuiltinFunction{
		Name:  "seed",
		Value: FuncAI64R(rand.Seed),
	},
	"read": &vvm.BuiltinFunction{
		Name: "read",
		Value: func(args ...vvm.Object) (ret vvm.Object, err error) {
			if len(args) != 1 {
				return nil, vvm.ErrWrongNumArguments
			}
			y1, ok := args[0].(*vvm.Bytes)
			if !ok {
				return nil, vvm.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "bytes",
					Found:    args[0].TypeName(),
				}
			}
			res, err := rand.Read(y1.Value)
			if err != nil {
				ret = wrapError(err)
				return
			}
			return &vvm.Int{Value: int64(res)}, nil
		},
	},
	"rand": &vvm.BuiltinFunction{
		Name: "rand",
		Value: func(args ...vvm.Object) (vvm.Object, error) {
			if len(args) != 1 {
				return nil, vvm.ErrWrongNumArguments
			}
			i1, ok := vvm.ToInt64(args[0])
			if !ok {
				return nil, vvm.ErrInvalidArgumentType{
					Name:     "first",
					Expected: "int(compatible)",
					Found:    args[0].TypeName(),
				}
			}
			src := rand.NewSource(i1)
			return randRand(rand.New(src)), nil
		},
	},
}

func randRand(r *rand.Rand) *vvm.ImmutableMap {
	return &vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			"int": &vvm.BuiltinFunction{
				Name:  "int",
				Value: FuncARI64(r.Int63),
			},
			"float": &vvm.BuiltinFunction{
				Name:  "float",
				Value: FuncARF(r.Float64),
			},
			"intn": &vvm.BuiltinFunction{
				Name:  "intn",
				Value: FuncAI64RI64(r.Int63n),
			},
			"exp_float": &vvm.BuiltinFunction{
				Name:  "exp_float",
				Value: FuncARF(r.ExpFloat64),
			},
			"norm_float": &vvm.BuiltinFunction{
				Name:  "norm_float",
				Value: FuncARF(r.NormFloat64),
			},
			"perm": &vvm.BuiltinFunction{
				Name:  "perm",
				Value: FuncAIRIs(r.Perm),
			},
			"seed": &vvm.BuiltinFunction{
				Name:  "seed",
				Value: FuncAI64R(r.Seed),
			},
			"read": &vvm.BuiltinFunction{
				Name: "read",
				Value: func(args ...vvm.Object) (
					ret vvm.Object,
					err error,
				) {
					if len(args) != 1 {
						return nil, vvm.ErrWrongNumArguments
					}
					y1, ok := args[0].(*vvm.Bytes)
					if !ok {
						return nil, vvm.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "bytes",
							Found:    args[0].TypeName(),
						}
					}
					res, err := r.Read(y1.Value)
					if err != nil {
						ret = wrapError(err)
						return
					}
					return &vvm.Int{Value: int64(res)}, nil
				},
			},
		},
	}
}
