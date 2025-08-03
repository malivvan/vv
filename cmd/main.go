package main

import (
	"context"
	"fmt"
	"github.com/malivvan/vv"
	"github.com/malivvan/vv/pkg/cli"
	"github.com/malivvan/vv/pkg/cui"
	"github.com/malivvan/vv/pkg/cui/vte"
	"os"
	"os/exec"
)

var (
	serial   string
	commit   string
	version  string
	compiled string
)

func main() {
	ctx := context.Background()
	app, err := vv.NewCli(func(c *cli.Context) error {
		if c.Args().Len() == 0 {
			xapp := cui.NewApplication()
			xvte := vte.NewTerminal(xapp, exec.Command(os.Args[0], "sh"))
			xvte.SetBorder(true)
			xapp.SetRoot(xvte, true)
			if err := xapp.Run(); err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating app: %s\n", err.Error())
		os.Exit(1)
	}

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
