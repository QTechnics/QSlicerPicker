package ui

import (
	_ "embed"
	"qslicerpicker/internal/config"
	"qslicerpicker/internal/i18n"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

//go:embed resources/q.svg
var qSVG []byte

//go:embed resources/q.png
var qPNG []byte

var mainApp fyne.App
var mainWindow fyne.Window

// RunApp starts the main application
func RunApp() {
	mainApp = app.NewWithID("com.qslicerpicker.app")

	// Set application icon from embedded SVG with custom color
	if len(qSVG) > 0 {
		// Replace fill color with #0193B1
		svgContent := strings.ReplaceAll(string(qSVG), `fill="#000000"`, `fill="#0193B1"`)
		svgContent = strings.ReplaceAll(svgContent, `fill='#000000'`, `fill="#0193B1"`)
		svgContent = strings.ReplaceAll(svgContent, `fill="black"`, `fill="#0193B1"`)

		// Create resource from modified SVG
		icon := fyne.NewStaticResource("q.svg", []byte(svgContent))
		mainApp.SetIcon(icon)
	} else if len(qPNG) > 0 {
		// Fallback to embedded PNG if SVG not available
		icon := fyne.NewStaticResource("q.png", qPNG)
		mainApp.SetIcon(icon)
	}

	// Initialize config
	config.GetConfig()

	// i18n is initialized automatically via init()

	// Create main window (hidden, for system tray)
	mainWindow = mainApp.NewWindow(i18n.T("app_title"))
	mainWindow.Resize(fyne.NewSize(800, 600))
	mainWindow.Hide() // Hide by default, show settings when needed

	// Create system tray
	createSystemTray()

	// Show settings window
	ShowSettings()

	mainApp.Run()
}

func createSystemTray() {
	if mainApp.Driver().Device().IsMobile() {
		return // System tray not available on mobile
	}

	// System tray is optional, can be added later if needed
	// For now, we'll skip it to avoid API issues
}

// GetApp returns the main Fyne app instance
func GetApp() fyne.App {
	return mainApp
}

// GetMainWindow returns the main window
func GetMainWindow() fyne.Window {
	return mainWindow
}
