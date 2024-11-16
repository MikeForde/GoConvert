package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"fyne.io/fyne/v2/storage"

)

type myTheme struct{}

func (m myTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameDisabled:
		return color.RGBA{R: 120, G: 120, B: 120, A: 255} // Light gray for disabled text
	case theme.ColorNameForeground:
		return color.RGBA{R: 0, G: 0, B: 0, A: 255} // Black for foreground text
	case theme.ColorNameBackground:
		return color.RGBA{R: 250, G: 250, B: 250, A: 255} // Very light gray for background
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 240, G: 240, B: 240, A: 255} // Slightly darker gray for input fields
	case theme.ColorNameButton:
		return color.RGBA{R: 220, G: 220, B: 220, A: 255} // Light gray for button background
	case theme.ColorNameHover:
		return color.RGBA{R: 180, G: 180, B: 180, A: 255} // Slightly darker gray for hover effects
	case theme.ColorNameOverlayBackground:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255} // White for overlay background
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

type HL7FHIRData struct {
    PackageUUID   string       `json:"packageUUID"`
    TimeStamp     string       `json:"timeStamp"`
    Patient       Patient      `json:"patient"`
    Medication    []Medication `json:"medication"`
    Allergies     []Allergy    `json:"allergies"`
    Conditions    []Condition  `json:"conditions"`
    Observations  []Observation `json:"observations"`
    Immunizations []Immunization `json:"immunizations"`
}

type Patient struct {
    Name         string `json:"name"`
    Given        string `json:"given"`
    DOB          string `json:"dob"`
    Gender       string `json:"gender"`
    Practitioner string `json:"practitioner"`
    Nation       string `json:"nation"`
    Organization string `json:"organization"`
}

type Medication struct {
    Name   string `json:"name"`
    Date   string `json:"date"`
    Dosage string `json:"dosage"`
}

type Allergy struct {
    Name        string `json:"name"`
    Criticality string `json:"criticality"`
    Date        string `json:"date"`
}

type Condition struct {
    Name string `json:"name"`
    Date string `json:"date"`
}

type Observation struct {
    Name  string `json:"name"`
    Date  string `json:"date"`
    Value string `json:"value"`
}

type Immunization struct {
    Name   string `json:"name"`
    System string `json:"system"`
    Date   string `json:"date"`
}


// HL7toFHIR to convert HL7 v2.8 to MongoDB JSON (as used on IPS MERN/SERN)
func HL7toFHIR(hl7Message string) (string, error) {
	lines := strings.Split(hl7Message, "\n")

	// Define the data - note we only need to add the array elements as the single elements are already members
	data := HL7FHIRData{
        Medication:    []Medication{},
        Allergies:     []Allergy{},
        Conditions:    []Condition{},
        Observations:  []Observation{},
        Immunizations: []Immunization{},
    }

	for _, line := range lines {
		segments := strings.Split(line, "|")
		if len(segments) < 2 {
			continue
		}

		switch segments[0] {
		case "MSH":
			if timestamp, err := parseHL7DateTime(segments[6]); err == nil {
                data.TimeStamp = timestamp
            }
            data.PackageUUID = segments[9]
		case "PID":
			data.Patient = Patient{
                Name:         parseSubfield(segments[5], 0),
                Given:        parseSubfield(segments[5], 1),
                Nation:       parseSubfield(segments[11], 3),
                Organization: parseSubfield(segments[3], 3),
            }
            if dob, err := parseHL7Date(segments[7]); err == nil {
                data.Patient.DOB = dob
            }
            switch strings.ToLower(segments[8]) {
            case "m":
                data.Patient.Gender = "male"
            case "f":
                data.Patient.Gender = "female"
            default:
                data.Patient.Gender = "other"
            }
		case "IVC":
            data.Patient.Practitioner = segments[2]
		case "RXA":
            medicationName := segments[5]
            date, _ := parseHL7DateTime(segments[3])
            if len(segments) > 6 {
                data.Medication = append(data.Medication, Medication{
                    Name:   medicationName,
                    Date:   date,
                    Dosage: segments[6],
                })
            } else {
                parts := strings.Split(medicationName, "^")
                immunization := Immunization{
                    Name:   parts[0],
                    Date:   date,
                }
                if len(parts) > 1 {
                    immunization.System = parts[1]
                } else {
                    immunization.System = "unknown"
                }
                data.Immunizations = append(data.Immunizations, immunization)
            }
		case "AL1":
			if len(segments) > 6 {
                data.Allergies = append(data.Allergies, Allergy{
                    Name:        parseSubfield(segments[3], 1),
                    Criticality: map[string]string{"U": "low", "SV": "high", "MO": "moderate", "MI": "mild"}[segments[4]],
                    Date:        segments[6],
                })
            }
		case "DG1":
			if len(segments) > 5 {
                data.Conditions = append(data.Conditions, Condition{
                    Name: parseSubfield(segments[3], 1),
                    Date: segments[5],
                })
            }
		case "OBX":
			if len(segments) > 12 {
                data.Observations = append(data.Observations, Observation{
                    Name:  parseSubfield(segments[3], 1),
                    Value: strings.TrimSpace(segments[5] + " " + segments[6]),
                    Date:  segments[12],
                })
            }
		
		default:
			// Probably not needed but keep for now
		}
	}

	fhirJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(fhirJSON), nil
}

func parseSubfield(field string, index int) string {
	parts := strings.Split(field, "^")
	if len(parts) > index {
		return parts[index]
	}
	return ""
}

func parseHL7DateTime(input string) (string, error) {
    t, err := time.Parse("20060102150405", input)
    if err != nil {
        return "", err
    }
    // We use this instead of RFC3339 so we can have the .000 bit at the end
    return t.Format("2006-01-02T15:04:05.000Z"), nil
}

// parseHL7Date parses an HL7 date (YYYYMMDD) as occurs in some fields.
func parseHL7Date(input string) (string, error) {
	t, err := time.Parse("20060102", input)
	if err != nil {
		return "", err
	}
	return t.Format(time.RFC3339), nil
}

func main() {
	myApp := app.New()
	// Apply some custom theme - otherwise some elements barely readable
	myApp.Settings().SetTheme(&myTheme{})
	myWindow := myApp.NewWindow("HL7 to IPS MERN MongoDb Converter")

	// The UI elements
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetPlaceHolder("Paste HL7 v2.8 content here...")

	convertButton := widget.NewButton("Convert", func() {
		hl7Content := inputEntry.Text
		if hl7Content == "" {
			dialog.ShowError(fmt.Errorf("No HL7 content provided"), myWindow)
			return
		}

		convertedJSON, err := HL7toFHIR(hl7Content)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

		outputWindow := myApp.NewWindow("Converted MongoDB JSON")
		outputEntry := widget.NewMultiLineEntry()
		outputEntry.SetText(convertedJSON)
		// Next one is the main benefit from custom theme
		outputEntry.Disable()

		saveButton := widget.NewButton("Save", func() {
			patientName := "Unknown"
			packageUUID := "Unknown"
	
			var data map[string]interface{}
		
			// We will use the patient name and packageUUID from the converted JSON in the save file name
			if err := json.Unmarshal([]byte(convertedJSON), &data); err == nil {
				if patient, ok := data["patient"].(map[string]interface{}); ok {
					if name, exists := patient["name"].(string); exists {
						patientName = name
					}
				}
				if uuid, exists := data["packageUUID"].(string); exists {
					packageUUID = uuid
				}
			}
		
			defaultFilename := fmt.Sprintf("%s_%s.json", patientName, packageUUID)
		
			dialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, _ error) {
				if writer == nil {
					return
				}
				defer writer.Close()
		
				_, err := writer.Write([]byte(convertedJSON))
				if err != nil {
					dialog.ShowError(err, myWindow)
				}
			}, outputWindow)
		
			dialog.SetFileName(defaultFilename)
			dialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		
			dialog.Show()
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

	// Basic Layout
	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("HL7 v2.8 to IPS MERN MongoDB JSON Converter"),
		fileButton,
		inputEntry,
		convertButton,
	))
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.ShowAndRun()
}
