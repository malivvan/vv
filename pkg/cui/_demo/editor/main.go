// Demo code for the Box primitive.
package main

import (
	_ "embed"
	"github.com/gdamore/tcell/v2"
	"github.com/malivvan/vv/pkg/cui"
	"github.com/malivvan/vv/pkg/cui/editor"
)

//go:embed main.go
var mainGo string

func main() {
	app := cui.NewApplication()
	defer app.HandlePanic()

	buf := editor.NewBufferFromString(mainGo, "main.go")
	view := editor.NewView(buf)
	view.SetTheme("darcula")
	view.SetBorder(true)
	view.SetBorderAttributes(tcell.AttrBold)

	app.SetRoot(view, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
