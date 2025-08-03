package sh

import (
	"mvdan.cc/sh/v3/interp"
	"os"
	"strings"
)

type Prompt struct {
	runner *interp.Runner
}

func NewPrompt(runner *interp.Runner) *Prompt {
	return &Prompt{runner: runner}
}
func (p *Prompt) String() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	username := os.Getenv("USER")
	if username == "" {
		username = "unknown"
	}
	workdir := strings.Replace(p.runner.Dir, p.runner.Env.Get("HOME").String(), "~", 1)

	return "\033[01;32m" + hostname + "@" + username + "\033[0m" + ":" + "\033[01;34m" + workdir + "\033[0m" + "> "
}
