package stdlib

import (
	"fmt"

	"github.com/malivvan/vv/vvm"
)

var fmtModule = map[string]vvm.Object{
	"print":   &vvm.BuiltinFunction{Value: fmtPrint, NeedVMObj: true},
	"printf":  &vvm.BuiltinFunction{Value: fmtPrintf, NeedVMObj: true},
	"println": &vvm.BuiltinFunction{Value: fmtPrintln, NeedVMObj: true},
	"sprintf": &vvm.UserFunction{Name: "sprintf", Value: fmtSprintf},
}

func fmtPrint(args ...vvm.Object) (ret vvm.Object, err error) {
	vm := args[0].(*vvm.VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(vm.Out, printArgs...)
	return nil, nil
}

func fmtPrintf(args ...vvm.Object) (ret vvm.Object, err error) {
	vm := args[0].(*vvm.VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	numArgs := len(args)
	if numArgs == 0 {
		return nil, vvm.ErrWrongNumArguments
	}

	format, ok := args[0].(*vvm.String)
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		fmt.Fprint(vm.Out, format)
		return nil, nil
	}

	s, err := vvm.Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}
	fmt.Fprint(vm.Out, s)
	return nil, nil
}

func fmtPrintln(args ...vvm.Object) (ret vvm.Object, err error) {
	vm := args[0].(*vvm.VMObj).Value
	args = args[1:] // the first arg is VMObj inserted by VM
	printArgs, err := getPrintArgs(args...)
	if err != nil {
		return nil, err
	}
	printArgs = append(printArgs, "\n")
	fmt.Fprint(vm.Out, printArgs...)
	return nil, nil
}

func fmtSprintf(args ...vvm.Object) (ret vvm.Object, err error) {
	numArgs := len(args)
	if numArgs == 0 {
		return nil, vvm.ErrWrongNumArguments
	}

	format, ok := args[0].(*vvm.String)
	if !ok {
		return nil, vvm.ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		// okay to return 'format' directly as String is immutable
		return format, nil
	}
	s, err := vvm.Format(format.Value, args[1:]...)
	if err != nil {
		return nil, err
	}
	return &vvm.String{Value: s}, nil
}

func getPrintArgs(args ...vvm.Object) ([]interface{}, error) {
	var printArgs []interface{}
	l := 0
	for _, arg := range args {
		s, _ := vvm.ToString(arg)
		slen := len(s)
		// make sure length does not exceed the limit
		if l+slen > vvm.MaxStringLen {
			return nil, vvm.ErrStringLimit
		}
		l += slen
		printArgs = append(printArgs, s)
	}
	return printArgs, nil
}
