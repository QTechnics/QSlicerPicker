package main

import (
	"os"

	"qslicerpicker/internal/filehandler"
	"qslicerpicker/internal/ui"
)

func main() {
	// Check if a file path is provided as argument
	if len(os.Args) > 1 {
		filePath := os.Args[1]
		// Show selector dialog and handle file
		filehandler.HandleFile(filePath)
		return
	}

	// Otherwise, show main application window
	ui.RunApp()
}
