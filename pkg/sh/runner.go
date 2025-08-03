package sh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/malivvan/vv/pkg/sh/readline"
	"io"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"os"
	"strings"
)

func runScript(runner *interp.Runner, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return runReader(runner, f, file)

}

func runReader(runner *interp.Runner, reader io.Reader, name string) error {

	prog, err := syntax.NewParser().Parse(reader, name)
	if err != nil {
		return err
	}
	runner.Reset()
	return runner.Run(context.Background(), prog)
}

func runInteractive(runner *interp.Runner, stdout, stderr io.Writer) error {
	prompt := NewPrompt(runner)
	l, err := readline.NewFromConfig(&readline.Config{
		HistoryFile:       "/tmp/readline.tmp",
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			switch r {
			// block CtrlZ feature
			case readline.CharCtrlZ:
				return r, false
			}
			return r, true
		},
		Undo: true,
	})
	if err != nil {
		return err
	}
	defer l.Close()

	l.CaptureExitSignal()

	setPasswordCfg := l.GeneratePasswordConfig()
	setPasswordCfg.Listener = func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	}

	parser := syntax.NewParser()
	var runnerErr error
	parserFn := func(stmts []*syntax.Stmt) bool {
		if parser.Incomplete() {
			fmt.Fprintf(stdout, "> ")
			return true
		}
		ctx := context.Background()
		for _, stmt := range stmts {
			runnerErr = runner.Run(ctx, stmt)
			if runner.Exited() {
				return false
			}
		}
		//	fmt.Fprintf(stdout, "$ ")
		return true
	}

	for {
		l.SetPrompt(prompt.String())

		line, err := l.ReadLine()
		if errors.Is(err, readline.ErrInterrupt) {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		switch line := strings.TrimSpace(line); {
		case line == "":
		case line == "exit":
			return nil
		default:
			if err := parser.Interactive(bytes.NewBufferString(line+"\n"), parserFn); err != nil {
				return err
			}
			if runnerErr != nil {
				return runnerErr
			}
		}
	}
	return nil
}
