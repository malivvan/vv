// Demo code for the CheckBox primitive.
package main

import (
	"github.com/malivvan/vv/pkg/cui"
)

func main() {
	app := cui.NewApplication()
	defer app.HandlePanic()

	app.EnableMouse(true)

	checkbox := cui.NewCheckBox()
	checkbox.SetLabel("Hit Enter to check box: ")

	app.SetRoot(checkbox, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
