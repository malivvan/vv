// Demo code for the InputField primitive.
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/malivvan/vv/pkg/cui"
)

func main() {
	app := cui.NewApplication()
	defer app.HandlePanic()

	app.EnableMouse(true)

	inputField := cui.NewInputField()
	inputField.SetLabel("Enter a number: ")
	inputField.SetPlaceholder("E.g. 1234")
	inputField.SetFieldWidth(10)
	inputField.SetAcceptanceFunc(cui.InputFieldInteger)
	inputField.SetDoneFunc(func(key tcell.Key) {
		app.Stop()
	})

	app.SetRoot(inputField, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
