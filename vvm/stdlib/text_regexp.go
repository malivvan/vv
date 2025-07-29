package stdlib

import (
	"regexp"

	"github.com/malivvan/vv/vvm"
)

func makeTextRegexp(re *regexp.Regexp) *vvm.ImmutableMap {
	return &vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			// match(text) => bool
			"match": &vvm.UserFunction{
				Value: func(args ...vvm.Object) (
					ret vvm.Object,
					err error,
				) {
					if len(args) != 1 {
						err = vvm.ErrWrongNumArguments
						return
					}

					s1, ok := vvm.ToString(args[0])
					if !ok {
						err = vvm.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
						return
					}

					if re.MatchString(s1) {
						ret = vvm.TrueValue
					} else {
						ret = vvm.FalseValue
					}

					return
				},
			},

			// find(text) 			=> array(array({text:,begin:,end:}))/undefined
			// find(text, maxCount) => array(array({text:,begin:,end:}))/undefined
			"find": &vvm.UserFunction{
				Value: func(args ...vvm.Object) (
					ret vvm.Object,
					err error,
				) {
					numArgs := len(args)
					if numArgs != 1 && numArgs != 2 {
						err = vvm.ErrWrongNumArguments
						return
					}

					s1, ok := vvm.ToString(args[0])
					if !ok {
						err = vvm.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
						return
					}

					if numArgs == 1 {
						m := re.FindStringSubmatchIndex(s1)
						if m == nil {
							ret = vvm.UndefinedValue
							return
						}

						arr := &vvm.Array{}
						for i := 0; i < len(m); i += 2 {
							arr.Value = append(arr.Value,
								&vvm.ImmutableMap{
									Value: map[string]vvm.Object{
										"text": &vvm.String{
											Value: s1[m[i]:m[i+1]],
										},
										"begin": &vvm.Int{
											Value: int64(m[i]),
										},
										"end": &vvm.Int{
											Value: int64(m[i+1]),
										},
									}})
						}

						ret = &vvm.Array{Value: []vvm.Object{arr}}

						return
					}

					i2, ok := vvm.ToInt(args[1])
					if !ok {
						err = vvm.ErrInvalidArgumentType{
							Name:     "second",
							Expected: "int(compatible)",
							Found:    args[1].TypeName(),
						}
						return
					}
					m := re.FindAllStringSubmatchIndex(s1, i2)
					if m == nil {
						ret = vvm.UndefinedValue
						return
					}

					arr := &vvm.Array{}
					for _, m := range m {
						subMatch := &vvm.Array{}
						for i := 0; i < len(m); i += 2 {
							subMatch.Value = append(subMatch.Value,
								&vvm.ImmutableMap{
									Value: map[string]vvm.Object{
										"text": &vvm.String{
											Value: s1[m[i]:m[i+1]],
										},
										"begin": &vvm.Int{
											Value: int64(m[i]),
										},
										"end": &vvm.Int{
											Value: int64(m[i+1]),
										},
									}})
						}

						arr.Value = append(arr.Value, subMatch)
					}

					ret = arr

					return
				},
			},

			// replace(src, repl) => string
			"replace": &vvm.UserFunction{
				Value: func(args ...vvm.Object) (
					ret vvm.Object,
					err error,
				) {
					if len(args) != 2 {
						err = vvm.ErrWrongNumArguments
						return
					}

					s1, ok := vvm.ToString(args[0])
					if !ok {
						err = vvm.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
						return
					}

					s2, ok := vvm.ToString(args[1])
					if !ok {
						err = vvm.ErrInvalidArgumentType{
							Name:     "second",
							Expected: "string(compatible)",
							Found:    args[1].TypeName(),
						}
						return
					}

					s, ok := doTextRegexpReplace(re, s1, s2)
					if !ok {
						return nil, vvm.ErrStringLimit
					}

					ret = &vvm.String{Value: s}

					return
				},
			},

			// split(text) 			 => array(string)
			// split(text, maxCount) => array(string)
			"split": &vvm.UserFunction{
				Value: func(args ...vvm.Object) (
					ret vvm.Object,
					err error,
				) {
					numArgs := len(args)
					if numArgs != 1 && numArgs != 2 {
						err = vvm.ErrWrongNumArguments
						return
					}

					s1, ok := vvm.ToString(args[0])
					if !ok {
						err = vvm.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
						return
					}

					var i2 = -1
					if numArgs > 1 {
						i2, ok = vvm.ToInt(args[1])
						if !ok {
							err = vvm.ErrInvalidArgumentType{
								Name:     "second",
								Expected: "int(compatible)",
								Found:    args[1].TypeName(),
							}
							return
						}
					}

					arr := &vvm.Array{}
					for _, s := range re.Split(s1, i2) {
						arr.Value = append(arr.Value,
							&vvm.String{Value: s})
					}

					ret = arr

					return
				},
			},
		},
	}
}

// Size-limit checking implementation of regexp.ReplaceAllString.
func doTextRegexpReplace(re *regexp.Regexp, src, repl string) (string, bool) {
	idx := 0
	out := ""
	for _, m := range re.FindAllStringSubmatchIndex(src, -1) {
		var exp []byte
		exp = re.ExpandString(exp, repl, src, m)
		if len(out)+m[0]-idx+len(exp) > vvm.MaxStringLen {
			return "", false
		}
		out += src[idx:m[0]] + string(exp)
		idx = m[1]
	}
	if idx < len(src) {
		if len(out)+len(src)-idx > vvm.MaxStringLen {
			return "", false
		}
		out += src[idx:]
	}
	return out, true
}
