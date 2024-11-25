package app

import (
    "fmt"
	"encoding/json"
	"strings"
	"time"
	. "myapp/models"
)

func HL7toMongoDb(hl7Message string) (string, error) {
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
			if timestamp, err := parseHL7DateOrDateTime(segments[6]); err == nil {
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
            if dob, err := parseHL7DateOrDateTime(segments[7]); err == nil {
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
            date, _ := parseHL7DateOrDateTime(segments[3])
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
                date, _ := parseHL7DateOrDateTime(segments[6])
                data.Allergies = append(data.Allergies, Allergy{
                    Name:        parseSubfield(segments[3], 1),
                    Criticality: map[string]string{"U": "low", "SV": "high", "MO": "moderate", "MI": "mild"}[segments[4]],
                    Date:        date,
                })
            }
		case "DG1":
			if len(segments) > 5 {
                date, _ := parseHL7DateOrDateTime(segments[5])
                data.Conditions = append(data.Conditions, Condition{
                    Name: parseSubfield(segments[3], 1),
                    Date: date,
                })
            }
		case "OBX":
			if len(segments) > 12 {
                date, _ := parseHL7DateOrDateTime(segments[12])
                data.Observations = append(data.Observations, Observation{
                    Name:  parseSubfield(segments[3], 1),
                    Value: strings.TrimSpace(segments[5] + " " + segments[6]),
                    Date:  date,
                })
            }
		
		default:
			// Probably not needed but keep for now
		}
	}

	mongodbJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(mongodbJSON), nil
}

// Helper functions for HL7 parsing
func parseSubfield(field string, index int) string {
	parts := strings.Split(field, "^")
	if len(parts) > index {
		return parts[index]
	}
	return ""
}

func parseHL7DateOrDateTime(input string) (string, error) {
	// Try to parse as long form (date and time)
	t, err := time.Parse("20060102150405", input)
	if err == nil {
		return t.Format("2006-01-02T15:04:05.000Z"), nil
	}

	// If parsing as long form fails, try short form (date only)
	t, err = time.Parse("20060102", input)
	if err == nil {
		return t.Format("2006-01-02T15:04:05.000Z"), nil
	}

	// If both parsing attempts fail, return an error
	return "", fmt.Errorf("invalid HL7 date format: %s", input)
}

