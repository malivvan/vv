package sh

import (
	"flag"
	"github.com/malivvan/vv/pkg/sh/interp"
	"github.com/malivvan/vv/pkg/sh/readline/term"
	"io"
	"os"
	"strings"
)

func Exec(stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	stdio := interp.StdIO(stdin, stdout, stderr)

	flags := flag.NewFlagSet("vsh", flag.ExitOnError)
	var command string
	flags.StringVar(&command, "c", "", "Read and execute commands from the given string value.")
	if err := flags.Parse(args); err != nil {
		return err
	}
	args = flags.Args()

	if command != "" {
		runner, err := interp.New(stdio, interp.Params(args...))
		if err != nil {
			return err
		}
		return runReader(runner, strings.NewReader(command), "")
	}

	if len(args) == 0 {
		if r, ok := stdin.(*os.File); ok && term.IsTerminal(int(r.Fd())) {
			runner, err := interp.New(stdio, interp.Params(args...), interp.Interactive(true), interp.ExecHandlers(vvMiddleware))
			if err != nil {
				return err
			}
			if err := runInteractive(runner, stdout, stderr); err != nil {
				return err
			}
			return nil
		}
		runner, err := interp.New(stdio, interp.Params(args...))
		if err != nil {
			return err
		}
		return runReader(runner, stdin, "")
	}
	runner, err := interp.New(stdio, interp.Params(args...))
	if err != nil {
		return err
	}
	if err := runScript(runner, args[0]); err != nil {
		return err
	}
	return nil
}
