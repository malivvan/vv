package stdlib

import (
	"context"
	"github.com/malivvan/vv/pkg/cui"
	"github.com/malivvan/vv/pkg/cui/chart"
	"github.com/malivvan/vv/pkg/cui/editor"
	"github.com/malivvan/vv/pkg/cui/menu"
	"github.com/malivvan/vv/pkg/cui/vte"
	"github.com/malivvan/vv/vvm"
)

var _ cui.Primitive = (*chart.BarChart)(nil)
var _ cui.Primitive = (*menu.MenuBar)(nil)
var _ cui.Primitive = (*vte.Terminal)(nil)
var _ cui.Primitive = (*editor.View)(nil)

var cuiModule = map[string]vvm.Object{
	"new": &vvm.BuiltinFunction{
		Value: func(ctx context.Context, args ...vvm.Object) (vvm.Object, error) {
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
