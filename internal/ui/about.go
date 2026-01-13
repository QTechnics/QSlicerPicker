package ui

import (
	_ "embed"
	"net/url"
	"qslicerpicker/internal/i18n"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

//go:embed resources/q.png
var qPNGEmbed []byte

func createAboutTab() fyne.CanvasObject {
	// Load logo from embedded resource
	var logoObj fyne.CanvasObject
	if len(qPNGEmbed) > 0 {
		logo := fyne.NewStaticResource("q.png", qPNGEmbed)
		img := canvas.NewImageFromResource(logo)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(128, 128))
		logoObj = img
	} else {
		logoObj = widget.NewLabel("Logo not found")
	}

	title := widget.NewLabel("3D Slicer Picker")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	version := widget.NewLabel(i18n.T("version") + " 1.0.0")
	version.Alignment = fyne.TextAlignCenter

	// Author section
	authorLabel := widget.NewLabel(i18n.T("author") + ":")
	authorLabel.Alignment = fyne.TextAlignCenter
	authorUrl, _ := url.Parse("https://github.com/victorioustr")
	authorLink := widget.NewHyperlink("Muzaffer AKYIL (victorioustr)", authorUrl)
	authorLink.Alignment = fyne.TextAlignCenter
	authorContainer := container.NewVBox(
		authorLabel,
		container.NewCenter(authorLink),
	)

	// Source code section
	sourceLabel := widget.NewLabel(i18n.T("source_code") + ":")
	sourceLabel.Alignment = fyne.TextAlignCenter
	sourceCodeUrl, _ := url.Parse("https://github.com/QTechnics/QSlicerPicker")
	sourceLink := widget.NewHyperlink("GitHub Repository", sourceCodeUrl)
	sourceLink.Alignment = fyne.TextAlignCenter
	sourceContainer := container.NewVBox(
		sourceLabel,
		container.NewCenter(sourceLink),
	)

	copyright := widget.NewLabel("© 2026 QTeknoloji")
	copyright.Alignment = fyne.TextAlignCenter

	// Main content with spacing
	content := container.NewVBox(
		container.NewCenter(logoObj),
		widget.NewSeparator(),
		title,
		version,
		widget.NewSeparator(),
		authorContainer,
		widget.NewSeparator(),
		sourceContainer,
		widget.NewSeparator(),
		copyright,
	)

	// Center everything
	return container.NewCenter(
		container.NewVBox(
			content,
		),
	)
}
