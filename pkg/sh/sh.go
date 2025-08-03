package sh

import (
	"errors"
	"github.com/malivvan/vv/pkg/sh/readline/term"
	"io"
	"mvdan.cc/sh/v3/interp"
	"os"
	"strings"
)

type Config struct {
	Args     []string
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Command  string
	Executor func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc
}

func (conf *Config) runner(interactive bool) (*interp.Runner, error) {
	if conf.Stdin == nil {
		return nil, errors.New("stdin is required")
	}
	if conf.Stdout == nil {
		return nil, errors.New("stdout is required")
	}
	if conf.Stderr == nil {
		return nil, errors.New("stderr is required")
	}
	options := []interp.RunnerOption{
		interp.StdIO(conf.Stdin, conf.Stdout, conf.Stderr),
		interp.Params(conf.Args...),
	}
	if interactive {
		options = append(options, interp.Interactive(true))
	}
	if conf.Executor != nil {
		options = append(options, interp.ExecHandlers(conf.Executor))
	}
	return interp.New(options...)
}

func Exec(conf *Config) error {

	if conf.Command != "" {
		runner, err := conf.runner(false)
		if err != nil {
			return err
		}
		return runReader(runner, strings.NewReader(conf.Command), "")
	}

	if len(conf.Args) == 0 {
		if r, ok := conf.Stdin.(*os.File); ok && term.IsTerminal(int(r.Fd())) {
			runner, err := conf.runner(true)
			if err != nil {
				return err
			}
			if err := runInteractive(runner, conf.Stdout, conf.Stderr); err != nil {
				return err
			}
			return nil
		}
		runner, err := conf.runner(false)
		if err != nil {
			return err
		}
		return runReader(runner, conf.Stdin, "")
	}

	runner, err := conf.runner(false)
	if err != nil {
		return err
	}
	if err := runScript(runner, conf.Args[0]); err != nil {
		return err
	}
	return nil
}
