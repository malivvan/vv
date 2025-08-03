package vv

import (
	"bufio"
	"context"
	"fmt"
	"github.com/malivvan/vv/vvm"
	"github.com/malivvan/vv/vvm/parser"
	"github.com/malivvan/vv/vvm/stdlib"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var Modules = stdlib.GetModuleMap(stdlib.AllModuleNames()...)

// CompileOnly compiles the source code and writes the compiled binary into
// outputFile.
func CompileOnly(data []byte, inputFile, outputFile string) (err error) {
	program, err := compileSrc(data, inputFile)
	if err != nil {
		return
	}

	if outputFile == "" {
		outputFile = basename(inputFile) + ".out"
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = out.Close()
		} else {
			err = out.Close()
		}
	}()

	b, err := program.Marshal()
	if err != nil {
		return
	}
	_, err = out.Write(b)
	if err != nil {
		return fmt.Errorf("error writing to output file %s: %w", outputFile, err)
	}

	return
}

// CompileAndRun compiles the source code and executes it.
func CompileAndRun(ctx context.Context, data []byte, inputFile string) (err error) {
	p, err := compileSrc(data, inputFile)
	if err != nil {
		return
	}
	err = p.RunContext(ctx)
	return
}

// RunCompiled reads the compiled binary from file and executes it.
func RunCompiled(ctx context.Context, data []byte) (err error) {
	p := &Program{}
	err = p.Unmarshal(data)
	if err != nil {
		return
	}
	err = p.RunContext(ctx)
	return
}

// RunREPL starts REPL.
func RunREPL(ctx context.Context, in io.Reader, out io.Writer, prompt string) {
	stdin := bufio.NewScanner(in)
	fileSet := parser.NewFileSet()
	globals := make([]vvm.Object, vvm.GlobalsSize)
	symbolTable := vvm.NewSymbolTable()
	for idx, fn := range vvm.GetAllBuiltinFunctions() {
		symbolTable.DefineBuiltin(idx, fn.Name)
	}

	// embed println function
	symbol := symbolTable.Define("__repl_println__")
	globals[symbol.Index] = &vvm.BuiltinFunction{
		Name: "println",
		Value: func(ctx context.Context, args ...vvm.Object) (ret vvm.Object, err error) {
			var printArgs []interface{}
			for _, arg := range args {
				if _, isUndefined := arg.(*vvm.Undefined); isUndefined {
					printArgs = append(printArgs, "<undefined>")
				} else {
					s, _ := vvm.ToString(arg)
					printArgs = append(printArgs, s)
				}
			}
			printArgs = append(printArgs, "\n")
			_, _ = fmt.Print(printArgs...)
			return
		},
	}

	var constants []vvm.Object
	for {
		_, _ = fmt.Fprint(out, prompt)
		scanned := stdin.Scan()
		if !scanned {
			return
		}

		line := stdin.Text()
		srcFile := fileSet.AddFile("repl", -1, len(line))
		p := parser.NewParser(srcFile, []byte(line), nil)
		file, err := p.ParseFile()
		if err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}

		file = addPrints(file)
		c := vvm.NewCompiler(srcFile, symbolTable, constants, Modules, nil)
		if err := c.Compile(file); err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}

		bytecode := c.Bytecode()
		machine := vvm.NewVM(ctx, bytecode, globals, -1)
		if err := machine.Run(); err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}
		constants = bytecode.Constants
	}
}

func compileSrc(src []byte, inputFile string) (*Program, error) {
	s := NewScript(src)
	s.SetName(inputFile)
	s.SetImports(Modules)
	s.EnableFileImport(true)
	if err := s.SetImportDir(filepath.Dir(inputFile)); err != nil {
		return nil, fmt.Errorf("error setting import dir: %w", err)
	}
	return s.Compile()
}

func addPrints(file *parser.File) *parser.File {
	var stmts []parser.Stmt
	for _, s := range file.Stmts {
		switch s := s.(type) {
		case *parser.ExprStmt:
			stmts = append(stmts, &parser.ExprStmt{
				Expr: &parser.CallExpr{
					Func: &parser.Ident{Name: "__repl_println__"},
					Args: []parser.Expr{s.Expr},
				},
			})
		case *parser.AssignStmt:
			stmts = append(stmts, s)

			stmts = append(stmts, &parser.ExprStmt{
				Expr: &parser.CallExpr{
					Func: &parser.Ident{
						Name: "__repl_println__",
					},
					Args: s.LHS,
				},
			})
		default:
			stmts = append(stmts, s)
		}
	}
	return &parser.File{
		InputFile: file.InputFile,
		Stmts:     stmts,
	}
}

func basename(s string) string {
	s = filepath.Base(s)
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}
