// Demo code for the bar chart primitive.
package main

import (
	"github.com/malivvan/vv/pkg/cui"
	"github.com/malivvan/vv/pkg/cui/chart"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	app := cui.NewApplication()
	gauge := chart.NewActivityModeGauge()
	gauge.SetTitle("activity mode gauge")
	gauge.SetPgBgColor(tcell.ColorOrange)
	gauge.SetRect(10, 4, 50, 3)
	gauge.SetBorder(true)

	update := func() {
		tick := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-tick.C:
				gauge.Pulse()
				app.Draw()
			}
		}
	}
	go update()

	app.SetRoot(gauge, false)
	app.EnableMouse(true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
