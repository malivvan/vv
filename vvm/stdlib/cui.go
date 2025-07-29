package stdlib

import (
	"github.com/malivvan/vv/pkg/cui"
	"github.com/malivvan/vv/vvm"
)

var cuiModule = map[string]vvm.Object{
	"new": &vvm.UserFunction{
		Value: func(args ...vvm.Object) (vvm.Object, error) {
			app := cui.NewApplication()
			txt := cui.NewTextView()
			app.SetRoot(txt, true)
			return NewObject[*cui.Application](app, "cui.App", func(app *cui.Application) string {
				return app.GetScreen().CharacterSet()
			}, map[string]any{
				"run": app.Run,
				"write": func(b []byte) (int, error) {
					defer app.QueueUpdateDraw(func() {}, txt)
					return txt.Write(b)
				},
			}), nil
		},
	},
}
