package main

import (
	"fmt"
	"image/color"
	"math"
	"math/big"
	"strconv"

	"opencalcc/mathcat"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

	historyTitle := widget.NewRichTextFromMarkdown("#d# History")
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

func makeGraph(filename string, function1 string, function2 string, function3 string, function4 string, domainMin string, domainMax string, rangeMin string, rangeMax string) {
	p := plot.New()

	p.Title.Text = "Functions"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// plottedfunc := plotter.NewFunction(func(x float64) float64 { return mathcat.Parse(function) })
	// plottedfunc.Color = color.RGBA{B: 255, A: 255}

	fn1, _ := parseFunction(function1)
	fn2, _ := parseFunction(function2)
	fn3, _ := parseFunction(function3)
	fn4, _ := parseFunction(function4)

	fn1.Color = color.RGBA{R: 255, A: 255}
	fn2.Color = color.RGBA{G: 255, A: 255}
	fn3.Color = color.RGBA{B: 255, A: 255}

	// legend
	p.Add(fn1)
	p.Add(fn2)
	p.Add(fn3)
	p.Add(fn4)
	p.Legend.Add(function1, fn1)
	p.Legend.Add(function2, fn2)
	p.Legend.Add(function3, fn3)
	p.Legend.Add(function4, fn4)
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

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
	if exprStr == "" {
		fn := plotter.NewFunction(func(x float64) float64 {
			return math.NaN()
		})
		fn.Samples = 1000
		return fn, nil
	}
	fn := plotter.NewFunction(func(x float64) float64 {
		variables := map[string]*big.Rat{
			"x": new(big.Rat).SetFloat64(x),
			"pi": func() *big.Rat {
				rat, _ := new(big.Float).SetFloat64(math.Pi).Rat(nil)
				return rat
			}(),
			"e": func() *big.Rat {
				rat, _ := new(big.Float).SetFloat64(math.E).Rat(nil)
				return rat
			}(),
		}

		res, err := mathcat.Exec(exprStr, variables)
		if err != nil || res == nil {
			return math.NaN()
		}
		floatVal, ok := res.Float64()
		if !ok {
			return math.NaN()
		}
		if math.IsInf(floatVal, 0) || math.IsNaN(floatVal) {
			return math.NaN()
		}
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
