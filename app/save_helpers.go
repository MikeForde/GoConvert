package app

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func SaveToFile(convertedJSON string, parentWindow fyne.Window) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(convertedJSON), &data); err != nil {
		dialog.ShowError(err, parentWindow)
		return
	}

	patientName := "Unknown"
	packageUUID := "Unknown"

	if patient, ok := data["patient"].(map[string]interface{}); ok {
		if name, exists := patient["name"].(string); exists {
			patientName = name
		}
	}
	if uuid, exists := data["packageUUID"].(string); exists {
		packageUUID = uuid
	}

	defaultFilename := fmt.Sprintf("%s_%s.json", patientName, packageUUID)

	dialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, _ error) {
		if writer == nil {
			return
		}
		defer writer.Close()

		_, err := writer.Write([]byte(convertedJSON))
		if err != nil {
			dialog.ShowError(err, parentWindow)
		}
	}, parentWindow)

	dialog.SetFileName(defaultFilename)
	dialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	dialog.Show()
}
