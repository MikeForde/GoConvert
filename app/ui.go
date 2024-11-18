package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fmt"
	"io/ioutil"
)

func Run() {
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})
	myWindow := myApp.NewWindow("HL7 to IPS MERN MongoDb Converter")

	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Paste HL7 v2.8 content here...")

	convertButton := widget.NewButton("Convert", func() {
		hl7Content := inputEntry.Text
		if hl7Content == "" {
			dialog.ShowError(fmt.Errorf("No HL7 content provided"), myWindow)
			return
		}

		convertedJSON, err := HL7toMongoDb(hl7Content)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		outputWindow := myApp.NewWindow("Converted MongoDB JSON")
		outputEntry := widget.NewMultiLineEntry()
		outputEntry.SetText(convertedJSON)
		outputEntry.Disable()

		saveButton := widget.NewButton("Save", func() {
			SaveToFile(convertedJSON, outputWindow)
		})

		outputWindow.SetContent(container.NewBorder(nil, saveButton, nil, nil, outputEntry))
		outputWindow.Resize(fyne.NewSize(600, 400))
		outputWindow.Show()
	})

	fileButton := widget.NewButton("Open File", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			content, err := ioutil.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			inputEntry.SetText(string(content))
		}, myWindow)
	})

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("HL7 v2.8 to IPS MERN MongoDB JSON Converter"),
		fileButton,
		inputEntry,
		convertButton,
	))
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.ShowAndRun()
}
