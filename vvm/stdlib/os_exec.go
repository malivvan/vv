package stdlib

import (
	"context"
	"os/exec"

	"github.com/malivvan/vv/vvm"
)

func makeOSExecCommand(cmd *exec.Cmd) *vvm.ImmutableMap {
	return &vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			// combined_output() => bytes/error
			"combined_output": &vvm.BuiltinFunction{
				Name:  "combined_output",
				Value: FuncARYE(cmd.CombinedOutput),
			},
			// output() => bytes/error
			"output": &vvm.BuiltinFunction{
				Name:  "output",
				Value: FuncARYE(cmd.Output),
			}, //
			// run() => error
			"run": &vvm.BuiltinFunction{
				Name:  "run",
				Value: FuncARE(cmd.Run),
			}, //
			// start() => error
			"start": &vvm.BuiltinFunction{
				Name:  "start",
				Value: FuncARE(cmd.Start),
			}, //
			// wait() => error
			"wait": &vvm.BuiltinFunction{
				Name:  "wait",
				Value: FuncARE(cmd.Wait),
			}, //
			// set_path(path string)
			"set_path": &vvm.BuiltinFunction{
				Name: "set_path",
				Value: func(ctx context.Context, args ...vvm.Object) (vvm.Object, error) {
					if len(args) != 1 {
						return nil, vvm.ErrWrongNumArguments
					}
					s1, ok := vvm.ToString(args[0])
					if !ok {
						return nil, vvm.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
					}
					cmd.Path = s1
					return vvm.UndefinedValue, nil
				},
			},
			// set_dir(dir string)
			"set_dir": &vvm.BuiltinFunction{
				Name: "set_dir",
				Value: func(ctx context.Context, args ...vvm.Object) (vvm.Object, error) {
					if len(args) != 1 {
						return nil, vvm.ErrWrongNumArguments
					}
					s1, ok := vvm.ToString(args[0])
					if !ok {
						return nil, vvm.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "string(compatible)",
							Found:    args[0].TypeName(),
						}
					}
					cmd.Dir = s1
					return vvm.UndefinedValue, nil
				},
			},
			// set_env(env array(string))
			"set_env": &vvm.BuiltinFunction{
				Name: "set_env",
				Value: func(ctx context.Context, args ...vvm.Object) (vvm.Object, error) {
					if len(args) != 1 {
						return nil, vvm.ErrWrongNumArguments
					}

					var env []string
					var err error
					switch arg0 := args[0].(type) {
					case *vvm.Array:
						env, err = stringArray(arg0.Value, "first")
						if err != nil {
							return nil, err
						}
					case *vvm.ImmutableArray:
						env, err = stringArray(arg0.Value, "first")
						if err != nil {
							return nil, err
						}
					default:
						return nil, vvm.ErrInvalidArgumentType{
							Name:     "first",
							Expected: "array",
							Found:    arg0.TypeName(),
						}
					}
					cmd.Env = env
					return vvm.UndefinedValue, nil
				},
			},
			// process() => imap(process)
			"process": &vvm.BuiltinFunction{
				Name: "process",
				Value: func(ctx context.Context, args ...vvm.Object) (vvm.Object, error) {
					if len(args) != 0 {
						return nil, vvm.ErrWrongNumArguments
					}
					return makeOSProcess(cmd.Process), nil
				},
			},
		},
	}
}
