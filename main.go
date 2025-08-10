package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	window := app.NewWindow("OpenCalcc")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Result will appear here")
	output.Disable()

	window.SetContent(container.NewVBox(
		widget.NewLabel("OpenCalcc - A simple calculator"),
		input,
		output,
		widget.NewButton("Press Enter", PressedEnter),
	))

	input.OnSubmitted = func(text string) {
		if text == "" {
			return
		}
		result, err := evalMath(text)
		if err != nil {
			output.SetText("Error: " + err.Error())
		} else {
			output.SetText(fmt.Sprintf("%v", round(result, 4))) // TODO: ability to change decimal place settings while removing floating point errors
		}
	}

	window.Resize(fyne.NewSize(1000, 750))
	window.ShowAndRun()
	onLeave()
}

func onLeave() {
	fmt.Println("Exited")
}

func PressedEnter() {
	fmt.Println("Enter key pressed")
}

func evalMath(expr string) (float64, error) {
	expr = strings.ReplaceAll(expr, " ", "")
	var opperator rune
	for _, char := range "+-*/" {
		if strings.ContainsRune(expr, char) {
			opperator = char
			parts := strings.Split(expr, string(opperator))
			if len(parts) != 2 {
				return 0, fmt.Errorf("invalid expression: %s", expr)
			}
			part1, error1 := strconv.ParseFloat(parts[0], 64)
			part2, error2 := strconv.ParseFloat(parts[1], 64)
			if error1 != nil || error2 != nil {
				return 0, fmt.Errorf("invalid number in expression: %s", expr)
			}

			switch opperator {
			case '+':
				return part1 + part2, nil
			case '-':
				return part1 - part2, nil
			case '*':
				return part1 * part2, nil
			case '/':
				if part2 == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				return part1 / part2, nil
			}
		}
	}
	return 0, fmt.Errorf("no valid operator found in expression: %s", expr)
}

func round(num float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(num*pow) / pow
}
