// Demo code for the bar chart primitive.
package main

import (
	"github.com/malivvan/vv/pkg/cui"
	"github.com/malivvan/vv/pkg/cui/chart"
)

func main() {
	app := cui.NewApplication()
	dialog := chart.NewMessageDialog(chart.ErrorDailog)
	dialog.SetTitle("error dialog")
	dialog.SetMessage("This is first line of error\nThis is second line of the error message")
	dialog.SetDoneFunc(func() {
		app.Stop()
	})

	app.SetRoot(dialog, true)
	app.EnableMouse(true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
