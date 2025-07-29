// Demo code for the Button primitive.
package main

import "github.com/malivvan/vv/pkg/cui"

func main() {
	app := cui.NewApplication()
	defer app.HandlePanic()

	app.EnableMouse(true)

	button := cui.NewButton("Hit Enter to close")
	button.SetRect(0, 0, 22, 3)
	button.SetSelectedFunc(func() {
		app.Stop()
	})

	app.SetRoot(button, false)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
