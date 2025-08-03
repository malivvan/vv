package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/malivvan/vv"
	"github.com/malivvan/vv/pkg/sh"
	"os"
	"path/filepath"
)

const (
	replPrompt = ">> "
)

var (
	compileOutput string
	showHelp      bool
	showVersion   bool
	version       = "dev"
)

func init() {
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.StringVar(&compileOutput, "o", "", "Compile output file")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	if showHelp {
		doHelp()
		os.Exit(2)
	} else if showVersion {
		fmt.Println(version)
		return
	}
	if len(os.Args) == 2 && os.Args[1] == "sh" {
		if err := sh.Exec(os.Stdin, os.Stdout, os.Stderr, []string{}); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error running shell: %s\n", err.Error())
			os.Exit(1)
		}
		return
	}

	inputFile := flag.Arg(0)
	if inputFile == "" {
		vv.RunREPL(ctx, os.Stdin, os.Stdout, replPrompt)
		return
	}

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"Error reading input file: %s\n", err.Error())
		os.Exit(1)
	}

	inputFile, err = filepath.Abs(inputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error file path: %s\n", err)
		os.Exit(1)
	}

	if len(inputData) > 1 && string(inputData[:2]) == "#!" {
		copy(inputData, "//")
	}

	if compileOutput != "" {
		err := vv.CompileOnly(inputData, inputFile, compileOutput)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else if string(inputData[:len(vv.Magic)]) != vv.Magic {
		err := vv.CompileAndRun(ctx, inputData, inputFile)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	} else {
		if err := vv.RunCompiled(ctx, inputData); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

func doHelp() {
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("	vv [flags] {file}")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println()
	fmt.Println("	-o        compile output file")
	fmt.Println("	-version  show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("	vv")
	fmt.Println()
	fmt.Println("	          Start vv REPL")
	fmt.Println()
	fmt.Println("	vv app.vv")
	fmt.Println()
	fmt.Println("	          Compile and run script file")
	fmt.Println()
	fmt.Println("	vv -o app app.vv")
	fmt.Println()
	fmt.Println("	          Compile script file into program file")
	fmt.Println()
	fmt.Println("	vv app")
	fmt.Println()
	fmt.Println("	          Run program or script file")
	fmt.Println()
	fmt.Println()
}
