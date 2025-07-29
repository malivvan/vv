// Demo code for the Frame primitive.
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/malivvan/vv/pkg/cui"
)

func main() {
	app := cui.NewApplication()
	defer app.HandlePanic()

	app.EnableMouse(true)

	box := cui.NewBox()
	box.SetBackgroundColor(tcell.ColorBlue.TrueColor())

	frame := cui.NewFrame(box)
	frame.SetBorders(2, 2, 2, 2, 4, 4)
	frame.AddText("Header left", true, cui.AlignLeft, tcell.ColorWhite.TrueColor())
	frame.AddText("Header middle", true, cui.AlignCenter, tcell.ColorWhite.TrueColor())
	frame.AddText("Header right", true, cui.AlignRight, tcell.ColorWhite.TrueColor())
	frame.AddText("Header second middle", true, cui.AlignCenter, tcell.ColorRed.TrueColor())
	frame.AddText("Footer middle", false, cui.AlignCenter, tcell.ColorGreen.TrueColor())
	frame.AddText("Footer second middle", false, cui.AlignCenter, tcell.ColorGreen.TrueColor())

	app.SetRoot(frame, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
