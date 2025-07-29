package stdlib

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/malivvan/vv/vvm"
)

var osModule = map[string]vvm.Object{
	"o_rdonly":            &vvm.Int{Value: int64(os.O_RDONLY)},
	"o_wronly":            &vvm.Int{Value: int64(os.O_WRONLY)},
	"o_rdwr":              &vvm.Int{Value: int64(os.O_RDWR)},
	"o_append":            &vvm.Int{Value: int64(os.O_APPEND)},
	"o_create":            &vvm.Int{Value: int64(os.O_CREATE)},
	"o_excl":              &vvm.Int{Value: int64(os.O_EXCL)},
	"o_sync":              &vvm.Int{Value: int64(os.O_SYNC)},
	"o_trunc":             &vvm.Int{Value: int64(os.O_TRUNC)},
	"mode_dir":            &vvm.Int{Value: int64(os.ModeDir)},
	"mode_append":         &vvm.Int{Value: int64(os.ModeAppend)},
	"mode_exclusive":      &vvm.Int{Value: int64(os.ModeExclusive)},
	"mode_temporary":      &vvm.Int{Value: int64(os.ModeTemporary)},
	"mode_symlink":        &vvm.Int{Value: int64(os.ModeSymlink)},
	"mode_device":         &vvm.Int{Value: int64(os.ModeDevice)},
	"mode_named_pipe":     &vvm.Int{Value: int64(os.ModeNamedPipe)},
	"mode_socket":         &vvm.Int{Value: int64(os.ModeSocket)},
	"mode_setuid":         &vvm.Int{Value: int64(os.ModeSetuid)},
	"mode_setgui":         &vvm.Int{Value: int64(os.ModeSetgid)},
	"mode_char_device":    &vvm.Int{Value: int64(os.ModeCharDevice)},
	"mode_sticky":         &vvm.Int{Value: int64(os.ModeSticky)},
	"mode_type":           &vvm.Int{Value: int64(os.ModeType)},
	"mode_perm":           &vvm.Int{Value: int64(os.ModePerm)},
	"path_separator":      &vvm.Char{Value: os.PathSeparator},
	"path_list_separator": &vvm.Char{Value: os.PathListSeparator},
	"dev_null":            &vvm.String{Value: os.DevNull},
	"seek_set":            &vvm.Int{Value: int64(io.SeekStart)},
	"seek_cur":            &vvm.Int{Value: int64(io.SeekCurrent)},
	"seek_end":            &vvm.Int{Value: int64(io.SeekEnd)},
	"args": &vvm.BuiltinFunction{
		Name:      "args",
		Value:     osArgs,
		NeedVMObj: true,
	}, // args() => array(string)
	"chdir": &vvm.UserFunction{
		Name:  "chdir",
		Value: FuncASRE(os.Chdir),
	}, // chdir(dir string) => error
	"chmod": osFuncASFmRE("chmod", os.Chmod), // chmod(name string, mode int) => error
	"chown": &vvm.UserFunction{
		Name:  "chown",
		Value: FuncASIIRE(os.Chown),
	}, // chown(name string, uid int, gid int) => error
	"clearenv": &vvm.UserFunction{
		Name:  "clearenv",
		Value: FuncAR(os.Clearenv),
	}, // clearenv()
	"environ": &vvm.UserFunction{
		Name:  "environ",
		Value: FuncARSs(os.Environ),
	}, // environ() => array(string)
	"exit": &vvm.UserFunction{
		Name:  "exit",
		Value: FuncAIR(os.Exit),
	}, // exit(code int)
	"expand_env": &vvm.UserFunction{
		Name:  "expand_env",
		Value: osExpandEnv,
	}, // expand_env(s string) => string
	"getegid": &vvm.UserFunction{
		Name:  "getegid",
		Value: FuncARI(os.Getegid),
	}, // getegid() => int
	"getenv": &vvm.UserFunction{
		Name:  "getenv",
		Value: FuncASRS(os.Getenv),
	}, // getenv(s string) => string
	"geteuid": &vvm.UserFunction{
		Name:  "geteuid",
		Value: FuncARI(os.Geteuid),
	}, // geteuid() => int
	"getgid": &vvm.UserFunction{
		Name:  "getgid",
		Value: FuncARI(os.Getgid),
	}, // getgid() => int
	"getgroups": &vvm.UserFunction{
		Name:  "getgroups",
		Value: FuncARIsE(os.Getgroups),
	}, // getgroups() => array(string)/error
	"getpagesize": &vvm.UserFunction{
		Name:  "getpagesize",
		Value: FuncARI(os.Getpagesize),
	}, // getpagesize() => int
	"getpid": &vvm.UserFunction{
		Name:  "getpid",
		Value: FuncARI(os.Getpid),
	}, // getpid() => int
	"getppid": &vvm.UserFunction{
		Name:  "getppid",
		Value: FuncARI(os.Getppid),
	}, // getppid() => int
	"getuid": &vvm.UserFunction{
		Name:  "getuid",
		Value: FuncARI(os.Getuid),
	}, // getuid() => int
	"getwd": &vvm.UserFunction{
		Name:  "getwd",
		Value: FuncARSE(os.Getwd),
	}, // getwd() => string/error
	"hostname": &vvm.UserFunction{
		Name:  "hostname",
		Value: FuncARSE(os.Hostname),
	}, // hostname() => string/error
	"lchown": &vvm.UserFunction{
		Name:  "lchown",
		Value: FuncASIIRE(os.Lchown),
	}, // lchown(name string, uid int, gid int) => error
	"link": &vvm.UserFunction{
		Name:  "link",
		Value: FuncASSRE(os.Link),
	}, // link(oldname string, newname string) => error
	"lookup_env": &vvm.UserFunction{
		Name:  "lookup_env",
		Value: osLookupEnv,
	}, // lookup_env(key string) => string/false
	"mkdir":     osFuncASFmRE("mkdir", os.Mkdir),        // mkdir(name string, perm int) => error
	"mkdir_all": osFuncASFmRE("mkdir_all", os.MkdirAll), // mkdir_all(name string, perm int) => error
	"readlink": &vvm.UserFunction{
		Name:  "readlink",
		Value: FuncASRSE(os.Readlink),
	}, // readlink(name string) => string/error
	"remove": &vvm.UserFunction{
		Name:  "remove",
		Value: FuncASRE(os.Remove),
	}, // remove(name string) => error
	"remove_all": &vvm.UserFunction{
		Name:  "remove_all",
		Value: FuncASRE(os.RemoveAll),
	}, // remove_all(name string) => error
	"rename": &vvm.UserFunction{
		Name:  "rename",
		Value: FuncASSRE(os.Rename),
	}, // rename(oldpath string, newpath string) => error
	"setenv": &vvm.UserFunction{
		Name:  "setenv",
		Value: FuncASSRE(os.Setenv),
	}, // setenv(key string, value string) => error
	"symlink": &vvm.UserFunction{
		Name:  "symlink",
		Value: FuncASSRE(os.Symlink),
	}, // symlink(oldname string newname string) => error
	"temp_dir": &vvm.UserFunction{
		Name:  "temp_dir",
		Value: FuncARS(os.TempDir),
	}, // temp_dir() => string
	"truncate": &vvm.UserFunction{
		Name:  "truncate",
		Value: FuncASI64RE(os.Truncate),
	}, // truncate(name string, size int) => error
	"unsetenv": &vvm.UserFunction{
		Name:  "unsetenv",
		Value: FuncASRE(os.Unsetenv),
	}, // unsetenv(key string) => error
	"create": &vvm.UserFunction{
		Name:  "create",
		Value: osCreate,
	}, // create(name string) => imap(file)/error
	"open": &vvm.UserFunction{
		Name:  "open",
		Value: osOpen,
	}, // open(name string) => imap(file)/error
	"open_file": &vvm.UserFunction{
		Name:  "open_file",
		Value: osOpenFile,
	}, // open_file(name string, flag int, perm int) => imap(file)/error
	"find_process": &vvm.UserFunction{
		Name:  "find_process",
		Value: osFindProcess,
	}, // find_process(pid int) => imap(process)/error
	"start_process": &vvm.UserFunction{
		Name:  "start_process",
		Value: osStartProcess,
	}, // start_process(name string, argv array(string), dir string, env array(string)) => imap(process)/error
	"exec_look_path": &vvm.UserFunction{
		Name:  "exec_look_path",
		Value: FuncASRSE(exec.LookPath),
	}, // exec_look_path(file) => string/error
	"exec": &vvm.UserFunction{
		Name:  "exec",
		Value: osExec,
	}, // exec(name, args...) => command
	"stat": &vvm.UserFunction{
		Name:  "stat",
		Value: osStat,
	}, // stat(name) => imap(fileinfo)/error
	"read_file": &vvm.UserFunction{
		Name:  "read_file",
		Value: osReadFile,
	}, // readfile(name) => array(byte)/error
}

func osReadFile(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		return nil, vvm.ErrWrongNumArguments
	}
	fname, ok := vvm.ToString(args[0])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return wrapError(err), nil
	}
	if len(bytes) > vvm.MaxBytesLen {
		return nil, vvm.ErrBytesLimit
	}
	return &vvm.Bytes{Value: bytes}, nil
}

func osStat(args ...vvm.Object) (ret vvm.Object, err error) {
	if len(args) != 1 {
		return nil, vvm.ErrWrongNumArguments
	}
	fname, ok := vvm.ToString(args[0])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	stat, err := os.Stat(fname)
	if err != nil {
		return wrapError(err), nil
	}
	fstat := &vvm.ImmutableMap{
		Value: map[string]vvm.Object{
			"name":  &vvm.String{Value: stat.Name()},
			"mtime": &vvm.Time{Value: stat.ModTime()},
			"size":  &vvm.Int{Value: stat.Size()},
			"mode":  &vvm.Int{Value: int64(stat.Mode())},
		},
	}
	if stat.IsDir() {
		fstat.Value["directory"] = vvm.TrueValue
	} else {
		fstat.Value["directory"] = vvm.FalseValue
	}
	return fstat, nil
}

func osCreate(args ...vvm.Object) (vvm.Object, error) {
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
	res, err := os.Create(s1)
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSFile(res), nil
}

func osOpen(args ...vvm.Object) (vvm.Object, error) {
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
	res, err := os.Open(s1)
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSFile(res), nil
}

func osOpenFile(args ...vvm.Object) (vvm.Object, error) {
	if len(args) != 3 {
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
	i2, ok := vvm.ToInt(args[1])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
	}
	i3, ok := vvm.ToInt(args[2])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "third",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
	}
	res, err := os.OpenFile(s1, i2, os.FileMode(i3))
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSFile(res), nil
}

func osArgs(args ...vvm.Object) (vvm.Object, error) {
	vm := args[0].(*vvm.VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	if len(args) != 0 {
		return nil, vvm.ErrWrongNumArguments
	}
	arr := &vvm.Array{}
	for _, osArg := range vm.Args {
		if len(osArg) > vvm.MaxStringLen {
			return nil, vvm.ErrStringLimit
		}
		arr.Value = append(arr.Value, &vvm.String{Value: osArg})
	}
	return arr, nil
}

func osFuncASFmRE(
	name string,
	fn func(string, os.FileMode) error,
) *vvm.UserFunction {
	return &vvm.UserFunction{
		Name: name,
		Value: func(args ...vvm.Object) (vvm.Object, error) {
			if len(args) != 2 {
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
			i2, ok := vvm.ToInt64(args[1])
			if !ok {
				return nil, vvm.ErrInvalidArgumentType{
					Name:     "second",
					Expected: "int(compatible)",
					Found:    args[1].TypeName(),
				}
			}
			return wrapError(fn(s1, os.FileMode(i2))), nil
		},
	}
}

func osLookupEnv(args ...vvm.Object) (vvm.Object, error) {
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
	res, ok := os.LookupEnv(s1)
	if !ok {
		return vvm.FalseValue, nil
	}
	if len(res) > vvm.MaxStringLen {
		return nil, vvm.ErrStringLimit
	}
	return &vvm.String{Value: res}, nil
}

func osExpandEnv(args ...vvm.Object) (vvm.Object, error) {
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
	var vlen int
	var failed bool
	s := os.Expand(s1, func(k string) string {
		if failed {
			return ""
		}
		v := os.Getenv(k)

		// this does not count the other texts that are not being replaced
		// but the code checks the final length at the end
		vlen += len(v)
		if vlen > vvm.MaxStringLen {
			failed = true
			return ""
		}
		return v
	})
	if failed || len(s) > vvm.MaxStringLen {
		return nil, vvm.ErrStringLimit
	}
	return &vvm.String{Value: s}, nil
}

func osExec(args ...vvm.Object) (vvm.Object, error) {
	if len(args) == 0 {
		return nil, vvm.ErrWrongNumArguments
	}
	name, ok := vvm.ToString(args[0])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	var execArgs []string
	for idx, arg := range args[1:] {
		execArg, ok := vvm.ToString(arg)
		if !ok {
			return nil, vvm.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("args[%d]", idx),
				Expected: "string(compatible)",
				Found:    args[1+idx].TypeName(),
			}
		}
		execArgs = append(execArgs, execArg)
	}
	return makeOSExecCommand(exec.Command(name, execArgs...)), nil
}

func osFindProcess(args ...vvm.Object) (vvm.Object, error) {
	if len(args) != 1 {
		return nil, vvm.ErrWrongNumArguments
	}
	i1, ok := vvm.ToInt(args[0])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	proc, err := os.FindProcess(i1)
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSProcess(proc), nil
}

func osStartProcess(args ...vvm.Object) (vvm.Object, error) {
	if len(args) != 4 {
		return nil, vvm.ErrWrongNumArguments
	}
	name, ok := vvm.ToString(args[0])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "first",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}
	var argv []string
	var err error
	switch arg1 := args[1].(type) {
	case *vvm.Array:
		argv, err = stringArray(arg1.Value, "second")
		if err != nil {
			return nil, err
		}
	case *vvm.ImmutableArray:
		argv, err = stringArray(arg1.Value, "second")
		if err != nil {
			return nil, err
		}
	default:
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "second",
			Expected: "array",
			Found:    arg1.TypeName(),
		}
	}

	dir, ok := vvm.ToString(args[2])
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "third",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	var env []string
	switch arg3 := args[3].(type) {
	case *vvm.Array:
		env, err = stringArray(arg3.Value, "fourth")
		if err != nil {
			return nil, err
		}
	case *vvm.ImmutableArray:
		env, err = stringArray(arg3.Value, "fourth")
		if err != nil {
			return nil, err
		}
	default:
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "fourth",
			Expected: "array",
			Found:    arg3.TypeName(),
		}
	}

	proc, err := os.StartProcess(name, argv, &os.ProcAttr{
		Dir: dir,
		Env: env,
	})
	if err != nil {
		return wrapError(err), nil
	}
	return makeOSProcess(proc), nil
}

func stringArray(arr []vvm.Object, argName string) ([]string, error) {
	var sarr []string
	for idx, elem := range arr {
		str, ok := elem.(*vvm.String)
		if !ok {
			return nil, vvm.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("%s[%d]", argName, idx),
				Expected: "string",
				Found:    elem.TypeName(),
			}
		}
		sarr = append(sarr, str.Value)
	}
	return sarr, nil
}
