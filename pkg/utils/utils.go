package utils

import (
	"fmt"
	"strings"
)

func ColorText(text string, color string) string {
	var colorCode string

	switch strings.ToLower(color) {
	case "verde", "green":
		colorCode = "\033[32m"
	case "rojo", "red":
		colorCode = "\033[31m"
	case "amarillo", "yellow":
		colorCode = "\033[33m"
	case "azul", "blue":
		colorCode = "\033[34m"
	default:
		colorCode = "\033[0m"
	}

	reset := "\033[0m"
	return colorCode + text + reset
}

func FormatTitle(title string) string {
	return ColorText("["+title+"]", "verde")
}

func PrintError(message string) {
	fmt.Println(ColorText("\n[ERROR]", "verde") + " " + message)
}

func PrintSuccess(message string) {
	fmt.Println(ColorText("\n[ÉXITO]", "verde") + " " + message)
}

func PrintWarning(message string) {
	fmt.Println(ColorText("\n[AVISO]", "verde") + " " + message)
}
