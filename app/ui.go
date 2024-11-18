package app

import (
	"fmt"
	"io/ioutil"
	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	. "myapp/models"
)

func Run() {
	myApp := app.New()
	myApp.Settings().SetTheme(&myTheme{})
	myWindow := myApp.NewWindow("HL7 to FHIR Converter")

	// UI Elements
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Paste HL7 2.x or IPS MERN MongoDB JSON content here...")

	// Dropdown for selecting conversion type
	conversionTypes := []string{
		"HL7 2.x to IPS MERN MongoDb JSON",
		"IPS MERN MongoDb JSON to IPS FHiR",
		"HL7 2.x to IPS FHiR",
	}
	conversionSelect := widget.NewSelect(conversionTypes, nil)
	conversionSelect.SetSelected(conversionTypes[0]) // Default to "HL7 to MongoDB"

	// Convert Button
	convertButton := widget.NewButton("Convert", func() {
		content := inputEntry.Text
		if content == "" {
			dialog.ShowError(fmt.Errorf("No content provided"), myWindow)
			return
		}

		var convertedJSON string
		var err error

		switch conversionSelect.Selected {
		case "HL7 2.x to IPS MERN MongoDb JSON":
			convertedJSON, err = HL7toMongoDb(content)
		case "IPS MERN MongoDb JSON to IPS FHiR":
			convertedJSON, err = GenerateIPSBundleFromMongo(content)
		case "HL7 2.x to IPS FHiR":
			mongoJSON, err := HL7toMongoDb(content)
			if err == nil {
				convertedJSON, err = GenerateIPSBundleFromMongo(mongoJSON)
			}
		default:
			err = fmt.Errorf("invalid conversion type selected")
		}

		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		// Display converted JSON in a new window
		outputWindow := myApp.NewWindow(fmt.Sprintf("Converted JSON (%s)", conversionSelect.Selected))
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

	// File Button
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

	// Layout
	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Select Conversion Type:"),
		conversionSelect,
		fileButton,
		inputEntry,
		convertButton,
	))
	myWindow.Resize(fyne.NewSize(500, 400))
	myWindow.ShowAndRun()
}

// GenerateIPSBundleFromMongo wraps the GenerateIPSBundle to take string input
func GenerateIPSBundleFromMongo(mongoJSON string) (string, error) {
	var ipsRecord HL7FHIRData
	err := json.Unmarshal([]byte(mongoJSON), &ipsRecord)
	if err != nil {
		return "", fmt.Errorf("failed to parse MongoDB JSON: %v", err)
	}

	return GenerateIPSBundle(ipsRecord)
}
