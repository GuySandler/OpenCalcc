package main

import (
	"fmt"
	"image/color"
	"math"
	"opencalcc/mathcat"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func main() {
	app := app.New()
	window := app.NewWindow("OpenCalcc")

	// graph
	functionLabel := widget.NewRichTextFromMarkdown("## Function Input")
	function1 := widget.NewEntry()
	function1.SetPlaceHolder("Enter function...")
	function1.Resize(fyne.NewSize(400, 40))

	function2 := widget.NewEntry()
	function2.SetPlaceHolder("Enter function...")
	function2.Resize(fyne.NewSize(400, 40))

	function3 := widget.NewEntry()
	function3.SetPlaceHolder("Enter function...")
	function3.Resize(fyne.NewSize(400, 40))

	function4 := widget.NewEntry()
	function4.SetPlaceHolder("Enter function...")
	function4.Resize(fyne.NewSize(400, 40))

	domainMin := widget.NewEntry()
	domainMin.SetPlaceHolder("Min")
	domainMin.Resize(fyne.NewSize(80, 50))
	domainMax := widget.NewEntry()
	domainMax.SetPlaceHolder("Max")
	domainMax.Resize(fyne.NewSize(80, 50))

	domainContainer := container.NewHBox(
		widget.NewLabel("Domain:"),
		domainMin,
		widget.NewLabel("to"),
		domainMax,
	)

	rangeMin := widget.NewEntry()
	rangeMin.SetPlaceHolder("Min")
	rangeMin.Resize(fyne.NewSize(80, 50))
	rangeMax := widget.NewEntry()
	rangeMax.SetPlaceHolder("Max")
	rangeMax.Resize(fyne.NewSize(80, 50))

	rangeContainer := container.NewHBox(
		widget.NewLabel("Range:"),
		rangeMin,
		widget.NewLabel("to"),
		rangeMax,
	)

	file := "opencalccgraph.png"
	makeGraph(file, function1.Text, function2.Text, function3.Text, function4.Text, domainMin.Text, domainMax.Text, rangeMin.Text, rangeMax.Text)
	graph := canvas.NewImageFromFile(file)
	graph.FillMode = canvas.ImageFillOriginal
	graph.Resize(fyne.NewSize(600, 450))

	regenGraph := widget.NewButton("Regenerate Graph", func() {
		makeGraph(file, function1.Text, function2.Text, function3.Text, function4.Text, domainMin.Text, domainMax.Text, rangeMin.Text, rangeMax.Text)
		graph.File = file
		graph.Refresh()
	})
	function1.OnSubmitted = func(text string) {
		makeGraph(file, function1.Text, function2.Text, function3.Text, function4.Text, domainMin.Text, domainMax.Text, rangeMin.Text, rangeMax.Text)
		graph.File = file
		graph.Refresh()
	}
	function2.OnSubmitted = function1.OnSubmitted
	function3.OnSubmitted = function1.OnSubmitted
	function4.OnSubmitted = function1.OnSubmitted
	domainMin.OnSubmitted = function1.OnSubmitted
	domainMax.OnSubmitted = function1.OnSubmitted
	rangeMin.OnSubmitted = function1.OnSubmitted
	rangeMax.OnSubmitted = function1.OnSubmitted

	// trace funcs
	traceLabel := widget.NewRichTextFromMarkdown("## Trace Graph")

	traceSelector := widget.NewLabel("Select a function to trace:")
	selectedfunc := "Func 1"
	traceFunctionSelector := widget.NewSelect([]string{"Func 1", "Func 2", "Func 3", "Func 4"}, func(selected string) {
		selectedfunc = selected
	})

	traceXnum := widget.NewEntry()
	traceXnum.SetPlaceHolder("X value")
	traceXnum.Resize(fyne.NewSize(100, 35))

	traceXresult := widget.NewLabel("Result will appear here")
	traceXresult.Wrapping = fyne.TextWrapWord

	traceX := widget.NewButton("Trace X", func() {
		x, err := strconv.ParseFloat(traceXnum.Text, 64)
		if err != nil {
			traceXresult.SetText("Invalid X value")
			return
		}
		var fn *plotter.Function
		switch selectedfunc {
		case "Func 1":
			fn, err = parseFunction(function1.Text)
		case "Func 2":
			fn, err = parseFunction(function2.Text)
		case "Func 3":
			fn, err = parseFunction(function3.Text)
		case "Func 4":
			fn, err = parseFunction(function4.Text)
		}
		if err != nil {
			traceXresult.SetText("Invalid function")
			return
		}

		y := fn.F(x)
		output := fmt.Sprintf("f(%s) = %.6g", traceXnum.Text, y)
		traceXresult.SetText(output)
	})

	traceYnum := widget.NewEntry()
	traceYnum.SetPlaceHolder("Y value")
	traceYnum.Resize(fyne.NewSize(100, 35))

	traceYresult := widget.NewLabel("Result will appear here")
	traceYresult.Wrapping = fyne.TextWrapWord

	traceY := widget.NewButton("Trace Y", func() {
		y, err := strconv.ParseFloat(traceYnum.Text, 64)
		if err != nil {
			traceYresult.SetText("Invalid Y value")
			return
		}
		var fn *plotter.Function
		switch selectedfunc {
		case "Func 1":
			fn, err = parseFunction(function1.Text)
		case "Func 2":
			fn, err = parseFunction(function2.Text)
		case "Func 3":
			fn, err = parseFunction(function3.Text)
		case "Func 4":
			fn, err = parseFunction(function4.Text)
		}
		if err != nil {
			traceYresult.SetText("Invalid function")
			return
		}

		dMin, err1 := strconv.ParseFloat(domainMin.Text, 64)
		if err1 != nil {
			dMin = 0
		}
		dMax, err2 := strconv.ParseFloat(domainMax.Text, 64)
		if err2 != nil {
			dMax = 10
		}

		x, found := findInverse(fn, y, dMin, dMax, 1e-6)
		if !found {
			traceYresult.SetText("No solution found")
			return
		}
		output := fmt.Sprintf("f⁻¹(%s) = %.6g", traceYnum.Text, x)
		traceYresult.SetText(output)
	})

	traceXcontainer := container.NewHBox(
		traceXnum,
		traceX,
	)
	traceYcontainer := container.NewHBox(
		traceYnum,
		traceY,
	)

	intersectionSelection := []string{}
	intersectionSelect1 := widget.NewSelect([]string{"Func 1", "Func 2", "Func 3", "Func 4"}, func(selected string) {
		intersectionSelection = append(intersectionSelection, selected)
	})
	intersectionSelect2 := widget.NewSelect([]string{"Func 1", "Func 2", "Func 3", "Func 4"}, func(selected string) {
		intersectionSelection = append(intersectionSelection, selected)
	})

	findIntersectionResult := widget.NewLabel("Intersection result will appear here")

	findIntersection := widget.NewButton("Find Intersection between 2 functions", func() {
		var fn1, fn2 *plotter.Function
		var err error

		if len(intersectionSelection) < 2 {
			findIntersectionResult.SetText("Please select 2 functions")
			return
		}

		// Get first function
		switch intersectionSelection[0] {
		case "Func 1":
			fn1, err = parseFunction(function1.Text)
		case "Func 2":
			fn1, err = parseFunction(function2.Text)
		case "Func 3":
			fn1, err = parseFunction(function3.Text)
		case "Func 4":
			fn1, err = parseFunction(function4.Text)
		}
		if err != nil || fn1 == nil {
			findIntersectionResult.SetText("Invalid first function")
			return
		}

		// Get second function
		switch intersectionSelection[1] {
		case "Func 1":
			fn2, err = parseFunction(function1.Text)
		case "Func 2":
			fn2, err = parseFunction(function2.Text)
		case "Func 3":
			fn2, err = parseFunction(function3.Text)
		case "Func 4":
			fn2, err = parseFunction(function4.Text)
		}
		if err != nil || fn2 == nil {
			findIntersectionResult.SetText("Invalid second function")
			return
		}

		dMin, err1 := strconv.ParseFloat(domainMin.Text, 64)
		if err1 != nil {
			dMin = -10
		}
		dMax, err2 := strconv.ParseFloat(domainMax.Text, 64)
		if err2 != nil {
			dMax = 10
		}

		x, y, found := findIntersection(fn1, fn2, dMin, dMax, 1e-6)
		if !found {
			findIntersectionResult.SetText("No intersection found")
			return
		}
		output := fmt.Sprintf("Intersection at (%.6g, %.6g)", x, y)
		findIntersectionResult.SetText(output)
	})

	// graph control panel
	controlPanel := container.NewVBox(
		functionLabel,
		function1,
		function2,
		function3,
		function4,
		container.NewHBox(layout.NewSpacer(), regenGraph, layout.NewSpacer()),
		widget.NewSeparator(),
		domainContainer,
		rangeContainer,
		widget.NewSeparator(),
		traceLabel,
		traceSelector,
		traceFunctionSelector,
		widget.NewLabel("Trace X"),
		traceXcontainer,
		traceXresult,
		widget.NewLabel("Trace Y"),
		traceYcontainer,
		traceYresult,
		widget.NewSeparator(),
		intersectionSelect1,
		intersectionSelect2,
		findIntersection,
		findIntersectionResult,
	)

	graphContent := container.NewHSplit(
		graph,
		container.NewScroll(controlPanel),
	)
	graphContent.Offset = 0.7

	// Calc
	calcTitle := widget.NewRichTextFromMarkdown("## Calculator")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	input.Resize(fyne.NewSize(400, 40))

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Result will appear here")
	output.Disable()
	output.Resize(fyne.NewSize(400, 40))

	historyTitle := widget.NewRichTextFromMarkdown("## History")
	historyScroll := container.NewScroll(container.NewVBox())
	historyScroll.SetMinSize(fyne.NewSize(400, 200))

	calcmode := false
	var calcmodeSwitch *widget.Check
	calcmodeSwitch = widget.NewCheck("Exact Mode", func(checked bool) {
		calcmode = checked
		if checked {
			calcmodeSwitch.SetText("Float Mode")
		} else {
			calcmodeSwitch.SetText("Exact Mode")
		}
	})

	calcButton := widget.NewButton("Calculate", func() {
		PressedEnter(input, output, historyScroll.Content, calcmode)
	})
	calcButton.Importance = widget.HighImportance

	clearButton := widget.NewButton("Clear Input", func() {
		input.SetText("")
		output.SetText("")
	})

	clearHistory := widget.NewButton("Clear History", func() {
		historyScroll.Content.(*fyne.Container).Objects = []fyne.CanvasObject{}
		historyScroll.Refresh()
	})

	buttonContainer := container.NewHBox(
		calcButton,
		clearButton,
		clearHistory,
	)

	calcContent := container.NewVBox(
		calcTitle,
		widget.NewSeparator(),
		input,
		output,
		container.NewHBox(layout.NewSpacer(), buttonContainer, layout.NewSpacer()),
		calcmodeSwitch,
		widget.NewSeparator(),
		historyTitle,
		historyScroll,
	)

	// tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Calculator", calcContent),
		container.NewTabItem("Graph", graphContent),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	window.SetContent(tabs)

	input.OnSubmitted = func(text string) {
		PressedEnter(input, output, historyScroll.Content, calcmode)
	}

	window.Resize(fyne.NewSize(1200, 800))
	window.ShowAndRun()
	onLeave()
}

func onLeave() {
	fmt.Println("Exited")
}

func PressedEnter(expression fyne.CanvasObject, output fyne.CanvasObject, history fyne.CanvasObject, calcmode bool) {
	if expression.(*widget.Entry).Text == "" {
		output.(*widget.Entry).SetText("")
		return
	}
	result, err := mathcat.Eval(expression.(*widget.Entry).Text)
	if err != nil {
		output.(*widget.Entry).SetText("Error: " + err.Error())
	} else {
		var outputText string
		if !calcmode {
			floatResult, _ := result.Float64()
			outputText = fmt.Sprintf("%.6g", floatResult)
		} else {
			outputText = result.String()
		}
		output.(*widget.Entry).SetText(outputText)
	}

	history.(*fyne.Container).Add(widget.NewLabel(fmt.Sprintf("%s = %s", expression.(*widget.Entry).Text, output.(*widget.Entry).Text)))
}

func makeGraph(filename string, function1 string, function2 string, function3 string, function4 string, domainMin string, domainMax string, rangeMin string, rangeMax string) {
	p := plot.New()

	p.Title.Text = "Functions"
	p.Title.TextStyle.Font.Size = vg.Points(24)
	p.X.Label.Text = "Y"
	p.Y.Label.Text = "X"
	p.X.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Label.TextStyle.Font.Size = vg.Points(20)
	p.X.Tick.Label.Font.Size = vg.Points(16)
	p.Y.Tick.Label.Font.Size = vg.Points(16)
	p.Legend.TextStyle.Font.Size = vg.Points(16)

	p.X.Width = vg.Points(2)
	p.Y.Width = vg.Points(2)

	importedDomainMin, err := strconv.ParseFloat(domainMin, 64)
	if err != nil {
		importedDomainMin = -10
	}
	importedDomainMax, err := strconv.ParseFloat(domainMax, 64)
	if err != nil {
		importedDomainMax = 10
	}
	importedRangeMin, err := strconv.ParseFloat(rangeMin, 64)
	if err != nil {
		importedRangeMin = -10
	}
	importedRangeMax, err := strconv.ParseFloat(rangeMax, 64)
	if err != nil {
		importedRangeMax = 10
	}

	majorInterval := 1.0
	if math.Abs(importedDomainMax-importedDomainMin) >= 15 {
		majorInterval = 5.0
	} else if math.Abs(importedDomainMax-importedDomainMin) >= 50 {
		majorInterval = 10.0
	} else if math.Abs(importedDomainMax-importedDomainMin) >= 100 {
		majorInterval = 20.0
	}

	p.X.Tick.Marker = SubTicker{Major: majorInterval, Minor: majorInterval / 4}
	p.Y.Tick.Marker = SubTicker{Major: majorInterval, Minor: majorInterval / 4}

	// grids
	mainGrid := plotter.NewGrid()
	mainGrid.Vertical.Width = vg.Points(2)
	mainGrid.Horizontal.Width = vg.Points(2)
	mainGrid.Vertical.Color = color.RGBA{A: 180}
	mainGrid.Horizontal.Color = color.RGBA{A: 180}

	// small grid
	fineGrid := plotter.NewGrid()
	fineGrid.Vertical.Width = vg.Points(0.5)
	fineGrid.Horizontal.Width = vg.Points(0.5)
	fineGrid.Vertical.Color = color.RGBA{A: 40}
	fineGrid.Horizontal.Color = color.RGBA{A: 40}

	p.Add(fineGrid)
	p.Add(mainGrid)

	// Draw axes
	if importedDomainMin <= 0 && importedDomainMax >= 0 {
		yAxis := plotter.NewFunction(func(x float64) float64 { return 0 })
		yAxis.Width = vg.Points(4)
		yAxis.Color = color.RGBA{A: 255}
		p.Add(yAxis)
	}
	if importedRangeMin <= 0 && importedRangeMax >= 0 {
		xAxis := plotter.NewFunction(func(x float64) float64 { return 0 })
		xAxis.Width = vg.Points(4)
		xAxis.Color = color.RGBA{A: 255}
		p.Add(xAxis)
	}

	funcs := []string{function1, function2, function3, function4}
	colors := []color.Color{
		color.RGBA{R: 255, A: 255},
		color.RGBA{G: 255, A: 255},
		color.RGBA{B: 255, A: 255},
		color.RGBA{B: 255, G: 255, A: 255},
	}

	for i, f := range funcs {
		if f == "" {
			continue
		}
		pts := generatePoints(f, importedDomainMin, importedDomainMax)
		line, err := plotter.NewLine(pts)
		if err != nil {
			continue
		}
		line.Color = colors[i]
		p.Add(line)
		p.Legend.Add(f, line)
	}

	p.X.Min = importedDomainMin
	p.X.Max = importedDomainMax
	p.Y.Min = importedRangeMin
	p.Y.Max = importedRangeMax
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	w := 8 * vg.Inch
	h := 8 * vg.Inch

	img := vgimg.New(w, h)
	c := draw.New(img)
	p.Draw(c)

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(f); err != nil {
		panic(err)
	}
}

type Points []struct{ X, Y float64 }

func (p Points) Len() int                    { return len(p) }
func (p Points) XY(i int) (float64, float64) { return p[i].X, p[i].Y }

func generatePoints(expr string, xmin, xmax float64) Points {
	if expr == "" {
		return nil
	}

	initialCapacity := 1000
	points := make(Points, 0, initialCapacity)
	steps := 7000
	dx := (xmax - xmin) / float64(steps)
	lastY := math.NaN()

	growthFactor := 1.5

	for i := 0; i <= steps; i++ {
		x := xmin + float64(i)*dx
		res, err := mathcat.Eval(strings.ReplaceAll(expr, "x", fmt.Sprintf("(%g)", x)))
		if err != nil {
			continue
		}
		y, ok := res.Float64()
		if !ok || math.IsInf(y, 0) || math.IsNaN(y) {
			lastY = math.NaN()
			continue
		}

		if !math.IsNaN(lastY) {
			delta := math.Abs(y - lastY)
			threshold := math.Max(math.Abs(y), math.Abs(lastY)) * 0.5
			if delta > threshold {
				lastY = math.NaN()
				continue
			}
		}

		// dynamically change amount of points
		if len(points) == cap(points) {
			newCap := int(float64(cap(points)) * growthFactor)
			newPoints := make(Points, len(points), newCap)
			copy(newPoints, points)
			points = newPoints
		}

		points = append(points, struct{ X, Y float64 }{x, y})
		lastY = y
	}
	return points
}

func parseFunction(exprStr string) (*plotter.Function, error) {
	if exprStr == "" {
		return plotter.NewFunction(func(x float64) float64 { return math.NaN() }), nil
	}
	pts := generatePoints(exprStr, -10, 10)
	if len(pts) == 0 {
		return nil, fmt.Errorf("invalid function")
	}

	fn := plotter.NewFunction(func(x float64) float64 {
		for i := 0; i < len(pts)-1; i++ {
			if x >= pts[i].X && x <= pts[i+1].X {
				dx := pts[i+1].X - pts[i].X
				dy := pts[i+1].Y - pts[i].Y
				return pts[i].Y + (x-pts[i].X)*dy/dx
			}
		}
		return math.NaN()
	})
	return fn, nil
}

func findInverse(fn *plotter.Function, y float64, Dmin, Dmax, tol float64) (float64, bool) {
	for i := 0; i < 100; i++ {
		mid := (Dmin + Dmax) / 2
		fmid := fn.F(mid)
		if math.IsNaN(fmid) {
			return math.NaN(), false
		}
		if math.Abs(fmid-y) < tol {
			return mid, true
		}
		if fmid < y {
			Dmin = mid
		} else {
			Dmax = mid
		}
	}
	return math.NaN(), false
}

func findIntersection(fn1, fn2 *plotter.Function, Dmin, Dmax, tol float64) (float64, float64, bool) {
	if fn1 == nil || fn2 == nil {
		return math.NaN(), math.NaN(), false
	}

	diffFn := plotter.NewFunction(func(x float64) float64 {
		return fn1.F(x) - fn2.F(x)
	})

	for i := 0; i < 100; i++ {
		mid := (Dmin + Dmax) / 2
		fmid := diffFn.F(mid)
		if math.IsNaN(fmid) {
			return math.NaN(), math.NaN(), false
		}
		if math.Abs(fmid) < tol {
			return mid, fn1.F(mid), true
		}
		if fmid < 0 {
			Dmin = mid
		} else {
			Dmax = mid
		}
	}
	return math.NaN(), math.NaN(), false

}

type SubTicker struct {
	Major, Minor float64
}

func (t SubTicker) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick

	majorStart := math.Ceil(min/t.Major) * t.Major
	for x := majorStart; x <= max; x += t.Major {
		ticks = append(ticks, plot.Tick{Value: x, Label: fmt.Sprintf("%.0f", x)})
	}

	minorStart := math.Ceil(min/t.Minor) * t.Minor
	for x := minorStart; x <= max; x += t.Minor {
		if math.Mod(x, t.Major) != 0 {
			ticks = append(ticks, plot.Tick{Value: x})
		}
	}

	return ticks
}
