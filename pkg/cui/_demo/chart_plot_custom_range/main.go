package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/malivvan/vv/pkg/cui"
	"github.com/malivvan/vv/pkg/cui/chart"
	"math"
)

func main() {

	app := cui.NewApplication()

	sinData := func() [][]float64 {
		n := 220
		data := make([][]float64, 2)
		data[0] = make([]float64, n)
		data[1] = make([]float64, n)
		for i := 0; i < n; i++ {
			data[0][i] = math.Sin(float64(i+1) / 5)
			// Avoid taking Cos(0) because it creates a high point of 2 that
			// will never be hit again and makes the graph look a little funny
			data[1][i] = math.Cos(float64(i+1) / 5)
		}
		return data
	}()

	bmLineChart := chart.NewPlot()
	bmLineChart.SetBorder(true)
	bmLineChart.SetTitle("line chart (braille mode)")
	bmLineChart.SetLineColor([]tcell.Color{
		tcell.ColorSteelBlue,
		tcell.ColorGreen,
	})
	bmLineChart.SetMarker(chart.PlotMarkerBraille)
	bmLineChart.SetYAxisAutoScaleMin(false)
	bmLineChart.SetYAxisAutoScaleMax(false)
	bmLineChart.SetYRange(-1.5, 1.5)
	bmLineChart.SetData(sinData)

	bmLineChart.SetDrawXAxisLabel(false)

	dmLineChart := chart.NewPlot()
	dmLineChart.SetBorder(true)
	dmLineChart.SetTitle("line chart (dot mode)")
	dmLineChart.SetLineColor([]tcell.Color{
		tcell.ColorDarkOrange,
	})
	dmLineChart.SetAxesLabelColor(tcell.ColorGold)
	dmLineChart.SetAxesColor(tcell.ColorGold)
	dmLineChart.SetMarker(chart.PlotMarkerDot)
	dmLineChart.SetDotMarkerRune('\u25c9')

	sampleData1 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sampleData2 := []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	dotModeChartData := [][]float64{sampleData1}
	dotModeChartData[0] = append(dotModeChartData[0], sampleData2...)
	dotModeChartData[0] = append(dotModeChartData[0], sampleData1[:5]...)
	dotModeChartData[0] = append(dotModeChartData[0], sampleData2[5:]...)
	dotModeChartData[0] = append(dotModeChartData[0], sampleData1[:7]...)
	dotModeChartData[0] = append(dotModeChartData[0], sampleData2[3:]...)
	dmLineChart.SetYAxisAutoScaleMin(false)
	dmLineChart.SetYAxisAutoScaleMax(false)
	dmLineChart.SetYRange(0, 3)
	dmLineChart.SetData(dotModeChartData)

	scatterPlotData := make([][]float64, 2)
	scatterPlotData[0] = []float64{1, 2, 3, 4, 5}
	scatterPlotData[1] = sinData[1][4:]
	dmScatterPlot := chart.NewPlot()

	dmScatterPlot.SetBorder(true)
	dmScatterPlot.SetTitle("scatter plot (dot mode)")
	dmScatterPlot.SetLineColor([]tcell.Color{
		tcell.ColorMediumSlateBlue,
		tcell.ColorLightSkyBlue,
	})
	dmScatterPlot.SetPlotType(chart.PlotTypeScatter)
	dmScatterPlot.SetMarker(chart.PlotMarkerDot)
	dmScatterPlot.SetYAxisAutoScaleMin(false)
	dmScatterPlot.SetYAxisAutoScaleMax(false)
	dmScatterPlot.SetYRange(-1, 3)
	dmScatterPlot.SetData(scatterPlotData)
	dmScatterPlot.SetDrawYAxisLabel(false)

	bmScatterPlot := chart.NewPlot()
	bmScatterPlot.SetBorder(true)
	bmScatterPlot.SetTitle("scatter plot (braille mode)")
	bmScatterPlot.SetLineColor([]tcell.Color{
		tcell.ColorGold,
		tcell.ColorLightSkyBlue,
	})
	bmScatterPlot.SetPlotType(chart.PlotTypeScatter)
	bmScatterPlot.SetMarker(chart.PlotMarkerBraille)
	bmScatterPlot.SetYAxisAutoScaleMin(false)
	bmScatterPlot.SetYAxisAutoScaleMax(false)
	bmScatterPlot.SetYRange(-1, 5)
	bmScatterPlot.SetData(scatterPlotData)

	firstRow := cui.NewFlex()
	firstRow.SetDirection(cui.FlexColumn)
	firstRow.AddItem(dmLineChart, 0, 1, false)
	firstRow.AddItem(bmLineChart, 0, 1, false)
	firstRow.SetRect(0, 0, 100, 15)

	secondRow := cui.NewFlex()
	secondRow.SetDirection(cui.FlexColumn)
	secondRow.AddItem(dmScatterPlot, 0, 1, false)
	secondRow.AddItem(bmScatterPlot, 0, 1, false)
	secondRow.SetRect(0, 0, 100, 15)

	layout := cui.NewFlex()
	layout.SetDirection(cui.FlexRow)
	layout.AddItem(firstRow, 0, 1, false)
	layout.AddItem(secondRow, 0, 1, false)
	layout.SetRect(0, 0, 100, 30)

	app.SetRoot(layout, false)
	app.EnableMouse(true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
