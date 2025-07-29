package main

import (
	"fmt"
	"os"

	"github.com/malivvan/vv/pkg/cui/vte"
)

func main() {
	fmt.Println("----- tcell-term parser example -----")
	fmt.Println("reading from stdin")
	parser := vte.NewParser(os.Stdin)
	for {
		seq := parser.Next()
		fmt.Printf("%s\n", seq)
		switch seq.(type) {
		case vte.EOF:
			return
		}

	}
}
