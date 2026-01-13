package ui

import (
	"qslicerpicker/internal/i18n"
	"qslicerpicker/internal/slicer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ShowSlicerSelector shows a dialog to select a slicer (uses main app)
func ShowSlicerSelector(filePath string, slicers []slicer.Slicer) *slicer.Slicer {
	app := GetApp()
	return ShowSlicerSelectorWithApp(app, filePath, slicers)
}

// ShowSlicerSelectorWithApp shows a dialog to select a slicer with a specific app instance
func ShowSlicerSelectorWithApp(fyneApp fyne.App, filePath string, slicers []slicer.Slicer) *slicer.Slicer {
	if len(slicers) == 0 {
		return nil
	}

	resultChan := make(chan *slicer.Slicer, 1)
	var selectedSlicer *slicer.Slicer

	win := fyneApp.NewWindow(i18n.T("open_in"))
	win.Resize(fyne.NewSize(400, 300))
	win.CenterOnScreen()
	win.SetFixedSize(true)

	// Create list widget
	list := widget.NewList(
		func() int {
			return len(slicers)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.SetText(slicers[id].Name)
		},
	)

	// Handle selection
	list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(slicers) {
			selectedSlicer = &slicers[id]
		}
	}

	// Select first item by default
	if len(slicers) > 0 {
		list.Select(0)
		selectedSlicer = &slicers[0]
	}

	// Create buttons
	cancelBtn := widget.NewButton(i18n.T("cancel"), func() {
		resultChan <- nil
		win.Close()
		fyneApp.Quit()
	})

	openBtn := widget.NewButton(i18n.T("open"), func() {
		resultChan <- selectedSlicer
		win.Close()
		fyneApp.Quit()
	})
	openBtn.Importance = widget.HighImportance

	// Create content with proper layout
	titleLabel := widget.NewLabel(i18n.T("choose_slicer"))
	titleLabel.Alignment = fyne.TextAlignCenter
	titleLabel.TextStyle = fyne.TextStyle{}

	// Buttons container - left aligned
	buttonsContainer := container.NewHBox(
		cancelBtn,
		openBtn,
	)

	// Main content: Title at top, list in center, buttons at bottom
	content := container.NewBorder(
		titleLabel,       // Top
		buttonsContainer, // Bottom
		nil, nil,         // Left, Right
		list, // Center
	)

	// Add padding
	paddedContent := container.NewBorder(
		nil, nil, nil, nil,
		container.NewPadded(content),
	)

	win.SetContent(paddedContent)
	win.Show()

	// Run app event loop (blocking) - this must be called on main thread
	fyneApp.Run()

	// After app.Run() returns, get the result
	select {
	case result := <-resultChan:
		return result
	default:
		return nil
	}
}
