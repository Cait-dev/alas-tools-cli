package main

import (
	"fmt"
	"os"

	"github.com/Cait-dev/alas-tools-cli/internal/config"
	"github.com/Cait-dev/alas-tools-cli/internal/ui"
)

var (
	version = "dev"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("Alas-Tools-Cli versión %s\n", version)
		return
	}

	config.LoadEnv()

	ui.ShowStartScreen()

	ui.StartMainMenu()

	fmt.Println("\n¡Hasta pronto! Gracias por usar Alas-Tools-Cli.")
}
