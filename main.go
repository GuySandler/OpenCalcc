package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"

	"opencalcc/mathcat"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	app := app.New()
	window := app.NewWindow("OpenCalcc")

	// graph
	function := widget.NewEntry()
	function.SetPlaceHolder("Enter function...")

	domainMin := widget.NewEntry()
	domainMin.SetPlaceHolder("Domain Min")
	domainMax := widget.NewEntry()
	domainMax.SetPlaceHolder("Domain Max")

	rangeMin := widget.NewEntry()
	rangeMin.SetPlaceHolder("Range Min")
	rangeMax := widget.NewEntry()
	rangeMax.SetPlaceHolder("Range Max")

	file := "opencalccgraph.png"
	makeGraph(file, function.Text, domainMin.Text, domainMax.Text, rangeMin.Text, rangeMax.Text)
	graph := canvas.NewImageFromFile(file)
	graph.FillMode = canvas.ImageFillOriginal
	graph.Resize(fyne.NewSize(800, 600))

	regenGraph := widget.NewButton("Regenerate Graph", func() {
		makeGraph(file, function.Text, domainMin.Text, domainMax.Text, rangeMin.Text, rangeMax.Text)
		graph.File = file
		graph.Refresh()
	})

	// trace funcs
	traceXnum := widget.NewEntry()
	traceXnum.SetPlaceHolder("Trace X")
	traceXresult := widget.NewLabel("")
	traceX := widget.NewButton("Trace X", func() {
		x, err := strconv.ParseFloat(traceXnum.Text, 64)
		if err != nil {
			return
		}

		fn, err := parseFunction(function.Text)
		if err != nil {
			return
		}

		y := fn.F(x)
		output := fmt.Sprintf("f(%s) = %f", traceXnum.Text, y)
		traceXresult.SetText(output)
	})

	traceYnum := widget.NewEntry()
	traceYnum.SetPlaceHolder("Trace Y")
	traceYresult := widget.NewLabel("")
	traceY := widget.NewButton("Trace Y", func() {
		y, err := strconv.ParseFloat(traceYnum.Text, 64)
		if err != nil {
			return
		}

		fn, err := parseFunction(function.Text)
		if err != nil {
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
		output := fmt.Sprintf("f⁻¹(%s) = %f", traceYnum.Text, x)
		traceYresult.SetText(output)
	})

	graphContent := container.NewVBox(
		graph,
		function,
		regenGraph,
		domainMin,
		domainMax,
		rangeMin,
		rangeMax,
		widget.NewLabel("Trace Functions"),
		traceXnum,
		traceX,
		traceXresult,
		traceYnum,
		traceY,
		traceYresult,
	)

	// Calc
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Result will appear here")
	output.Disable()

	history := container.NewVBox(
		widget.NewLabel("History"),
	)

	calcmode := false
	calcmodeSwitch := widget.NewCheck("Exact Mode", func(checked bool) {
		calcmode = checked
	})

	calcContent := container.NewVBox(
		widget.NewLabel("OpenCalcc"),
		input,
		output,
		widget.NewButton("Press Enter", func() {
			PressedEnter(input, output, history, calcmode)
		}),
		calcmodeSwitch,
		history,
	)

	// tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Calc", calcContent),
		container.NewTabItem("Graph", graphContent),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	window.SetContent(tabs)

	input.OnSubmitted = func(text string) {
		PressedEnter(input, output, history, calcmode)
	}

	window.Resize(fyne.NewSize(1000, 750))
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
		if !calcmode { // float else exact
			floatResult, _ := result.Float64()
			outputText = fmt.Sprintf("%.6g", floatResult)
		} else {
			outputText = fmt.Sprintf("%s", result)
		}
		output.(*widget.Entry).SetText(outputText)
	}

	history.(*fyne.Container).Add(widget.NewLabel(fmt.Sprintf("%s = %s", expression.(*widget.Entry).Text, output.(*widget.Entry).Text)))
}

func makeGraph(filename string, function string, domainMin string, domainMax string, rangeMin string, rangeMax string) {
	p := plot.New()

	p.Title.Text = "Functions"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// plottedfunc := plotter.NewFunction(func(x float64) float64 { return mathcat.Parse(function) })
	// plottedfunc.Color = color.RGBA{B: 255, A: 255}

	fn, _ := parseFunction(function)

	// legend
	p.Add(fn)
	p.Legend.Add(function, fn)
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	// Convert domainMin, domainMax, rangeMin, rangeMax to float64
	importedDomainMin, err := strconv.ParseFloat(domainMin, 64)
	if err != nil {
		importedDomainMin = 0
	}
	importedDomainMax, err := strconv.ParseFloat(domainMax, 64)
	if err != nil {
		importedDomainMax = 10
	}
	importedRangeMin, err := strconv.ParseFloat(rangeMin, 64)
	if err != nil {
		importedRangeMin = 0
	}
	importedRangeMax, err := strconv.ParseFloat(rangeMax, 64)
	if err != nil {
		importedRangeMax = 10
	}

	p.X.Min = importedDomainMin
	p.X.Max = importedDomainMax
	p.Y.Min = importedRangeMin
	p.Y.Max = importedRangeMax

	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		panic(err)
	}
}

func parseFunction(exprStr string) (*plotter.Function, error) {
	fn := plotter.NewFunction(func(x float64) float64 {
		variables := map[string]*big.Rat{
			"x": new(big.Rat).SetFloat64(x), // Use SetFloat64 instead of converting to int64
			"pi": func() *big.Rat {
				rat, _ := new(big.Float).SetFloat64(math.Pi).Rat(nil)
				return rat
			}(),
		}

		res, err := mathcat.Exec(exprStr, variables)
		if err != nil {
			log.Printf("eval error: %v", err)
			return math.NaN()
		}

		floatVal, _ := res.Float64()
		return floatVal
	})
	fn.Samples = 1000
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
