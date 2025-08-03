---
search:
  boost: 2
---

Let's build a simple "greeter" application to demonstrate the basic structure of a `cli` app. This example will create a command that prints a friendly greeting.

Start by creating a directory named `greet`, and within it, add a file,
`greet.go` with the following code in it:

<!-- {
  "output": "Hello friend!"
} -->
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/malivvan/vv/pkg/cli"
)

func main() {
	app := &cli.App{
		Name:  "greet",
		Usage: "fight the loneliness!",
		Action: func(*cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
```

Install our command to the `$GOPATH/bin` directory:

```sh-session
$ go install
```

Finally run our new command:

```sh-session
$ greet
Hello friend!
```

cli also generates neat help text:

```sh-session
$ greet help
NAME:
    greet - fight the loneliness!

USAGE:
    greet [global options] command [command options] [arguments...]

COMMANDS:
    help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS
    --help, -h  show help (default: false)
```
