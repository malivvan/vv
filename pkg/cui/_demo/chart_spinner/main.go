package main

import (
	"github.com/malivvan/vv/pkg/cui"
	"github.com/malivvan/vv/pkg/cui/chart"
	"time"
)

func main() {
	app := cui.NewApplication()
	grid := cui.NewGrid()
	grid.SetBorder(true)
	grid.SetTitle("Spinners")

	spinners := [][]*chart.Spinner{
		{
			chart.NewSpinner().SetStyle(chart.SpinnerDotsCircling),
			chart.NewSpinner().SetStyle(chart.SpinnerDotsUpDown),
			chart.NewSpinner().SetStyle(chart.SpinnerBounce),
			chart.NewSpinner().SetStyle(chart.SpinnerLine),
		},
		{
			chart.NewSpinner().SetStyle(chart.SpinnerCircleQuarters),
			chart.NewSpinner().SetStyle(chart.SpinnerSquareCorners),
			chart.NewSpinner().SetStyle(chart.SpinnerCircleHalves),
			chart.NewSpinner().SetStyle(chart.SpinnerCorners),
		},
		{
			chart.NewSpinner().SetStyle(chart.SpinnerArrows),
			chart.NewSpinner().SetStyle(chart.SpinnerHamburger),
			chart.NewSpinner().SetStyle(chart.SpinnerStack),
			chart.NewSpinner().SetStyle(chart.SpinnerStar),
		},
		{
			chart.NewSpinner().SetStyle(chart.SpinnerGrowHorizontal),
			chart.NewSpinner().SetStyle(chart.SpinnerGrowVertical),
			chart.NewSpinner().SetStyle(chart.SpinnerBoxBounce),
			chart.NewSpinner().SetCustomStyle([]rune{'🕛', '🕐', '🕑', '🕒', '🕓', '🕔', '🕕', '🕖', '🕗', '🕘', '🕙', '🕚'}),
		},
	}

	for rowIdx, row := range spinners {
		for colIdx, spinner := range row {
			grid.AddItem(spinner, rowIdx, colIdx, 1, 1, 1, 1, false)
		}
	}

	update := func() {
		tick := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-tick.C:
				for _, row := range spinners {
					for _, spinner := range row {
						spinner.Pulse()
					}
				}
				app.Draw()
			}
		}
	}
	go update()

	app.SetRoot(grid, false)
	app.EnableMouse(true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
