package ui

import (
	"fmt"
	"qslicerpicker/internal/config"
	"qslicerpicker/internal/i18n"
	"qslicerpicker/internal/slicer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var settingsWindow fyne.Window

// ShowSettings shows the settings window
func ShowSettings() {
	app := GetApp()
	if settingsWindow != nil {
		settingsWindow.Show()
		settingsWindow.RequestFocus()
		return
	}

	settingsWindow = app.NewWindow(i18n.T("settings"))
	settingsWindow.Resize(fyne.NewSize(900, 700))
	settingsWindow.CenterOnScreen()

	content := createSettingsContent()
	settingsWindow.SetContent(content)
	settingsWindow.Show()

	settingsWindow.SetOnClosed(func() {
		settingsWindow = nil
	})
}

func createSettingsContent() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		&container.TabItem{
			Text:    i18n.T("slicers"),
			Content: createSlicersTab(),
		},
		&container.TabItem{
			Text:    i18n.T("language"),
			Content: createLanguageTab(),
		},
		&container.TabItem{
			Text:    i18n.T("about"),
			Content: createAboutTab(),
		},
	)

	return tabs
}

func createSlicersTab() fyne.CanvasObject {
	allSlicers := slicer.LoadSlicers()

	// Create list of slicers with checkboxes
	var list *widget.List
	list = widget.NewList(
		func() int {
			return len(allSlicers)
		},
		func() fyne.CanvasObject {
			checkbox := widget.NewCheck("", nil)
			nameLabel := widget.NewLabel("")
			editBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)
			upBtn := widget.NewButtonWithIcon("", theme.MoveUpIcon(), nil)
			downBtn := widget.NewButtonWithIcon("", theme.MoveDownIcon(), nil)

			// Layout: Checkbox | Name | Spacer | Edit | Up | Down
			return container.NewBorder(
				nil, nil,
				checkbox,
				container.NewHBox(editBtn, upBtn, downBtn),
				nameLabel,
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(allSlicers) {
				return
			}

			s := allSlicers[id]
			borderContainer := obj.(*fyne.Container)

			// Find objects
			var checkbox *widget.Check
			var buttons *fyne.Container
			var nameLabel *widget.Label

			for _, obj := range borderContainer.Objects {
				if check, ok := obj.(*widget.Check); ok {
					checkbox = check
				} else if label, ok := obj.(*widget.Label); ok {
					nameLabel = label
				} else if cont, ok := obj.(*fyne.Container); ok {
					// Buttons container
					buttons = cont
				}
			}

			if checkbox == nil || buttons == nil || nameLabel == nil {
				return
			}

			editBtn := buttons.Objects[0].(*widget.Button)
			upBtn := buttons.Objects[1].(*widget.Button)
			downBtn := buttons.Objects[2].(*widget.Button)

			nameLabel.SetText(s.Name)
			checkbox.SetChecked(s.Enabled)

			// Update enabled state
			checkbox.OnChanged = func(checked bool) {
				updateSlicerEnabled(s.ID, checked)
			}

			// Edit button
			editBtn.OnTapped = func() {
				showSlicerDialog(&s, func(updatedSlicer slicer.Slicer) {
					// Update config based on slicer type
					cfg := config.GetConfig()

					if updatedSlicer.IsCustom {
						// Find and update custom slicer
						for i := range cfg.CustomSlicers {
							if fmt.Sprintf("custom_%d", i) == updatedSlicer.ID {
								cfg.CustomSlicers[i].Name = updatedSlicer.Name
								cfg.CustomSlicers[i].Path = updatedSlicer.Path
								cfg.CustomSlicers[i].Arguments = updatedSlicer.Arguments
								cfg.CustomSlicers[i].WorkingDir = updatedSlicer.WorkingDir
								cfg.CustomSlicers[i].Enabled = updatedSlicer.Enabled
								break
							}
						}
					} else {
						// Update default slicer config
						found := false
						for i := range cfg.Slicers {
							if cfg.Slicers[i].ID == updatedSlicer.ID {
								cfg.Slicers[i].CustomPath = updatedSlicer.Path
								cfg.Slicers[i].Arguments = updatedSlicer.Arguments
								cfg.Slicers[i].WorkingDir = updatedSlicer.WorkingDir
								cfg.Slicers[i].Enabled = updatedSlicer.Enabled
								found = true
								break
							}
						}
						if !found {
							// Add if not exists in config override
							cfg.Slicers = append(cfg.Slicers, config.SlicerConfig{
								ID:         updatedSlicer.ID,
								CustomPath: updatedSlicer.Path,
								Arguments:  updatedSlicer.Arguments,
								WorkingDir: updatedSlicer.WorkingDir,
								Enabled:    updatedSlicer.Enabled,
								Order:      updatedSlicer.Order,
							})
						}
					}

					config.SaveConfig()
					if settingsWindow != nil {
						settingsWindow.SetContent(createSettingsContent())
					}
				})
			}

			// Move up
			upBtn.OnTapped = func() {
				if id > 0 {
					moveSlicer(id, id-1)
					if settingsWindow != nil {
						settingsWindow.SetContent(createSettingsContent())
					}
				}
			}

			// Move down
			downBtn.OnTapped = func() {
				if id < len(allSlicers)-1 {
					moveSlicer(id, id+1)
					if settingsWindow != nil {
						settingsWindow.SetContent(createSettingsContent())
					}
				}
			}

			// Disable up button for first item
			upBtn.Disable()
			if id > 0 {
				upBtn.Enable()
			}

			// Disable down button for last item
			downBtn.Disable()
			if id < len(allSlicers)-1 {
				downBtn.Enable()
			}
		},
	)

	// Add custom slicer button
	addBtn := widget.NewButton(i18n.T("add_custom_slicer"), func() {
		showAddCustomSlicerDialog()
	})

	return container.NewBorder(
		nil,
		addBtn,
		nil, nil,
		list,
	)
}

func createLanguageTab() fyne.CanvasObject {
	currentLang := i18n.GetLanguage()

	langRadio := widget.NewRadioGroup(
		[]string{
			i18n.T("turkish"),
			i18n.T("english"),
			i18n.T("german"),
			i18n.T("french"),
		},
		func(selected string) {
			langMap := map[string]string{
				i18n.T("turkish"): "tr",
				i18n.T("english"): "en",
				i18n.T("german"):  "de",
				i18n.T("french"):  "fr",
			}

			if langCode, ok := langMap[selected]; ok {
				// Only update if language actually changed
				if langCode == i18n.GetLanguage() {
					return
				}

				i18n.SetLanguage(langCode)
				// Refresh settings window
				if settingsWindow != nil {
					settingsWindow.SetContent(createSettingsContent())
				}
			}
		},
	)

	// Set current selection
	langNames := map[string]string{
		"tr": i18n.T("turkish"),
		"en": i18n.T("english"),
		"de": i18n.T("german"),
		"fr": i18n.T("french"),
	}

	if name, ok := langNames[currentLang]; ok {
		langRadio.SetSelected(name)
	}

	return container.NewVBox(
		widget.NewLabel(i18n.T("language")),
		langRadio,
	)
}

func updateSlicerEnabled(id string, enabled bool) {
	cfg := config.GetConfig()
	for i := range cfg.Slicers {
		if cfg.Slicers[i].ID == id {
			cfg.Slicers[i].Enabled = enabled
			config.SaveConfig()
			return
		}
	}
	// If not found, add it
	cfg.Slicers = append(cfg.Slicers, config.SlicerConfig{
		ID:      id,
		Enabled: enabled,
	})
	config.SaveConfig()
}

func moveSlicer(from, to int) {
	cfg := config.GetConfig()
	allSlicers := slicer.LoadSlicers()

	if from < 0 || from >= len(allSlicers) || to < 0 || to >= len(allSlicers) {
		return
	}

	// Update order values
	fromOrder := allSlicers[from].Order
	toOrder := allSlicers[to].Order

	fromID := allSlicers[from].ID
	toID := allSlicers[to].ID

	// Find and update in config (for default slicers)
	for i := range cfg.Slicers {
		if cfg.Slicers[i].ID == fromID {
			cfg.Slicers[i].Order = toOrder
		} else if cfg.Slicers[i].ID == toID {
			cfg.Slicers[i].Order = fromOrder
		}
	}

	// Handle custom slicers
	if allSlicers[from].IsCustom {
		// Find custom slicer index
		for i := range cfg.CustomSlicers {
			if fmt.Sprintf("custom_%d", i) == fromID {
				cfg.CustomSlicers[i].Order = toOrder
				break
			}
		}
	}
	if allSlicers[to].IsCustom {
		// Find custom slicer index
		for i := range cfg.CustomSlicers {
			if fmt.Sprintf("custom_%d", i) == toID {
				cfg.CustomSlicers[i].Order = fromOrder
				break
			}
		}
	}

	config.SaveConfig()
}

func showSlicerDialog(s *slicer.Slicer, onSave func(s slicer.Slicer)) {
	if settingsWindow == nil {
		return
	}

	isEdit := s != nil
	title := i18n.T("add_custom_slicer")
	if isEdit {
		// Edit title
		title = s.Name
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder(i18n.T("name"))

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder(i18n.T("path"))

	argsEntry := widget.NewEntry()
	argsEntry.SetPlaceHolder(i18n.T("arguments"))

	workingDirEntry := widget.NewEntry()
	workingDirEntry.SetPlaceHolder(i18n.T("working_directory"))

	enabledCheck := widget.NewCheck(i18n.T("enabled"), nil)
	enabledCheck.SetChecked(true)

	if isEdit {
		nameEntry.SetText(s.Name)
		pathEntry.SetText(s.Path)
		if len(s.Arguments) > 0 {
			// Join arguments
			argsStr := ""
			for i, arg := range s.Arguments {
				if i > 0 {
					argsStr += " "
				}
				argsStr += arg
			}
			argsEntry.SetText(argsStr)
		}
		workingDirEntry.SetText(s.WorkingDir)
		enabledCheck.SetChecked(s.Enabled)

		// If not custom, disable name editing
		if !s.IsCustom {
			nameEntry.Disable()
		}
	}

	// Browse buttons
	browsePathBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			pathEntry.SetText(reader.URI().Path())
		}, settingsWindow)
	})

	browseDirBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			workingDirEntry.SetText(uri.Path())
		}, settingsWindow)
	})

	// Dialog instance
	var d dialog.Dialog

	saveBtn := widget.NewButton(i18n.T("save"), func() {
		if nameEntry.Text == "" || pathEntry.Text == "" {
			return
		}

		newSlicer := slicer.Slicer{
			Name:       nameEntry.Text,
			Path:       pathEntry.Text,
			WorkingDir: workingDirEntry.Text,
			Enabled:    enabledCheck.Checked,
			IsCustom:   true, // Default to true, logic will handle override
		}

		if isEdit {
			newSlicer.ID = s.ID
			newSlicer.Order = s.Order
			newSlicer.IsCustom = s.IsCustom
		}

		if argsEntry.Text != "" {
			// Simple split (TODO: better argument parsing)
			newSlicer.Arguments = []string{argsEntry.Text}
		} else {
			newSlicer.Arguments = []string{}
		}

		onSave(newSlicer)
		d.Hide()
	})
	saveBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem(i18n.T("name"), nameEntry),
			widget.NewFormItem(i18n.T("path"), container.NewBorder(nil, nil, nil, browsePathBtn, pathEntry)),
			widget.NewFormItem(i18n.T("arguments"), argsEntry),
			widget.NewFormItem(i18n.T("working_directory"), container.NewBorder(nil, nil, nil, browseDirBtn, workingDirEntry)),
		),
		enabledCheck,
		saveBtn, // Only save button, dismiss button is handled by dialog
	)

	d = dialog.NewCustom(title, i18n.T("cancel"), content, settingsWindow)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}

func showAddCustomSlicerDialog() {
	cfg := config.GetConfig()
	showSlicerDialog(nil, func(s slicer.Slicer) {
		customSlicer := config.CustomSlicer{
			Name:       s.Name,
			Path:       s.Path,
			Arguments:  s.Arguments,
			WorkingDir: s.WorkingDir,
			Enabled:    s.Enabled,
			Order:      len(cfg.CustomSlicers) * 10,
		}
		cfg.CustomSlicers = append(cfg.CustomSlicers, customSlicer)
		config.SaveConfig()

		if settingsWindow != nil {
			settingsWindow.SetContent(createSettingsContent())
		}
	})
}
