package filehandler

import (
	"fmt"
	"os"
	"qslicerpicker/internal/config"
	"qslicerpicker/internal/slicer"
	"qslicerpicker/internal/ui"

	"fyne.io/fyne/v2/app"
)

// HandleFile handles a file that should be opened with a slicer
func HandleFile(filePath string) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File not found: %s\n", filePath)
		os.Exit(1)
	}

	// Initialize config and i18n
	config.GetConfig()
	// i18n is initialized automatically via init()

	// Create a minimal app for the dialog
	fyneApp := app.NewWithID("com.qslicerpicker.selector")

	// Get enabled slicers
	enabledSlicers := slicer.GetEnabledSlicers()

	if len(enabledSlicers) == 0 {
		fmt.Fprintf(os.Stderr, "No slicers available\n")
		os.Exit(1)
	}

	// Show selector dialog
	selectedSlicer := ui.ShowSlicerSelectorWithApp(fyneApp, filePath, enabledSlicers)

	if selectedSlicer == nil {
		// User cancelled
		os.Exit(0)
	}

	// Launch slicer with file
	if err := slicer.LaunchSlicer(*selectedSlicer, filePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error launching slicer: %v\n", err)
		os.Exit(1)
	}
}
