// Demo code for the bar chart primitive.
package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/malivvan/vv/pkg/cui"
	"github.com/malivvan/vv/pkg/cui/chart"
)

func main() {
	app := cui.NewApplication()
	barGraph := chart.NewBarChart()
	barGraph.SetRect(4, 2, 50, 20)
	barGraph.SetBorder(true)
	barGraph.SetTitle("System Resource Usage")
	// display system metric usage
	barGraph.AddBar("cpu", 80, tcell.ColorBlue)
	barGraph.AddBar("mem", 20, tcell.ColorRed)
	barGraph.AddBar("swap", 40, tcell.ColorGreen)
	barGraph.AddBar("disk", 40, tcell.ColorOrange)
	barGraph.SetMaxValue(100)
	barGraph.SetAxesColor(tcell.ColorAntiqueWhite)
	barGraph.SetAxesLabelColor(tcell.ColorAntiqueWhite)

	app.SetRoot(barGraph, false)
	app.EnableMouse(true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
