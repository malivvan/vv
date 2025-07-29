package stdlib

import (
	"os"
	"syscall"

	"github.com/malivvan/vv/vvm"
)

func makeOSProcessState(state *os.ProcessState) *vvm.ImmutableMap {
	return &vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			"exited": &vvm.UserFunction{
				Name:  "exited",
				Value: FuncARB(state.Exited),
			},
			"pid": &vvm.UserFunction{
				Name:  "pid",
				Value: FuncARI(state.Pid),
			},
			"string": &vvm.UserFunction{
				Name:  "string",
				Value: FuncARS(state.String),
			},
			"success": &vvm.UserFunction{
				Name:  "success",
				Value: FuncARB(state.Success),
			},
		},
	}
}

func makeOSProcess(proc *os.Process) *vvm.ImmutableMap {
	return &vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			"kill": &vvm.UserFunction{
				Name:  "kill",
				Value: FuncARE(proc.Kill),
			},
			"release": &vvm.UserFunction{
				Name:  "release",
				Value: FuncARE(proc.Release),
			},
			"signal": &vvm.UserFunction{
				Name: "signal",
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
					return wrapError(proc.Signal(syscall.Signal(i1))), nil
				},
			},
			"wait": &vvm.UserFunction{
				Name: "wait",
				Value: func(args ...vvm.Object) (vvm.Object, error) {
					if len(args) != 0 {
						return nil, vvm.ErrWrongNumArguments
					}
					state, err := proc.Wait()
					if err != nil {
						return wrapError(err), nil
					}
					return makeOSProcessState(state), nil
				},
			},
		},
	}
}
