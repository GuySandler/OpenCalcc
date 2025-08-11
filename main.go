package main

import (
	"fmt"
	"opencalcc/mathcat"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	window := app.NewWindow("OpenCalcc")

	// graph
	graphContent := widget.NewLabel("Graph goes here")

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
