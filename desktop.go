package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type entryWriter struct {
	entry *widget.Entry
}

func (ew entryWriter) Write(p []byte) (int, error) {
	ew.entry.Append(string(p))
	return len(p), nil
}

// View
func desktopView(controller *Controller) {

	a := app.New()
	w := a.NewWindow("Okkular Scraper")

	logOutputTextArea := widget.NewMultiLineEntry()

	writer := entryWriter{logOutputTextArea}
	log.SetOutput(writer)

	// Set the minimum size for the window
	// w.SetFixedSize(true)
	w.Resize(fyne.NewSize(1200, 600))

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter or select Excel file path here")

	// File chooser button

	fileChooserButton := widget.NewButton("Select File", func() {
		fileChooser := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				input.SetText(reader.URI().Path())
				controller.UpdateExcelFilePath(input.Text)
			}
		}, w)
		fileChooser.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx", "xls"}))
		fileChooser.Show()
	})

	exitButton := widget.NewButton("  Quit   ", func() {
		w.Close()
	})

	// Scrape button
	scrapeButton := widget.NewButton(" Scrape ", func() {
		controller.model.UpdateExcelFilePath(input.Text)
		logOutputTextArea.SetText("")
		if err := controller.HandleScrape(); err != nil {
			dialog.NewError(err, w).Show()
		}
	})

	// Create a layout container
	content := container.NewBorder(
		container.NewBorder(
			nil,
			nil,
			widget.NewLabel("Excel File Path: "),
			fileChooserButton,
			input,
		),
		container.NewHBox(layout.NewSpacer(), scrapeButton, exitButton),
		nil,
		nil,
		container.NewScroll(
			logOutputTextArea,
		),
	)

	w.SetContent(content)
	w.ShowAndRun()
}
