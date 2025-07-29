package stdlib

import (
	"os"

	"github.com/malivvan/vv/vvm"
)

func makeOSFile(file *os.File) *vvm.ImmutableMap {
	return &vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			// chdir() => true/error
			"chdir": &vvm.UserFunction{
				Name:  "chdir",
				Value: FuncARE(file.Chdir),
			}, //
			// chown(uid int, gid int) => true/error
			"chown": &vvm.UserFunction{
				Name:  "chown",
				Value: FuncAIIRE(file.Chown),
			}, //
			// close() => error
			"close": &vvm.UserFunction{
				Name:  "close",
				Value: FuncARE(file.Close),
			}, //
			// name() => string
			"name": &vvm.UserFunction{
				Name:  "name",
				Value: FuncARS(file.Name),
			}, //
			// readdirnames(n int) => array(string)/error
			"readdirnames": &vvm.UserFunction{
				Name:  "readdirnames",
				Value: FuncAIRSsE(file.Readdirnames),
			}, //
			// sync() => error
			"sync": &vvm.UserFunction{
				Name:  "sync",
				Value: FuncARE(file.Sync),
			}, //
			// write(bytes) => int/error
			"write": &vvm.UserFunction{
				Name:  "write",
				Value: FuncAYRIE(file.Write),
			}, //
			// write(string) => int/error
			"write_string": &vvm.UserFunction{
				Name:  "write_string",
				Value: FuncASRIE(file.WriteString),
			}, //
			// read(bytes) => int/error
			"read": &vvm.UserFunction{
				Name:  "read",
				Value: FuncAYRIE(file.Read),
			}, //
			// chmod(mode int) => error
			"chmod": &vvm.UserFunction{
				Name: "chmod",
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
					return wrapError(file.Chmod(os.FileMode(i1))), nil
				},
			},
			// seek(offset int, whence int) => int/error
			"seek": &vvm.UserFunction{
				Name: "seek",
				Value: func(args ...vvm.Object) (vvm.Object, error) {
					if len(args) != 2 {
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
					i2, ok := vvm.ToInt(args[1])
					if !ok {
						return nil, vvm.ErrInvalidArgumentType{
							Name:     "second",
							Expected: "int(compatible)",
							Found:    args[1].TypeName(),
						}
					}
					res, err := file.Seek(i1, i2)
					if err != nil {
						return wrapError(err), nil
					}
					return &vvm.Int{Value: res}, nil
				},
			},
			// stat() => imap(fileinfo)/error
			"stat": &vvm.UserFunction{
				Name: "stat",
				Value: func(args ...vvm.Object) (vvm.Object, error) {
					if len(args) != 0 {
						return nil, vvm.ErrWrongNumArguments
					}
					return osStat(&vvm.String{Value: file.Name()})
				},
			},
		},
	}
}
