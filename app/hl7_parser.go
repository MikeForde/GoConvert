package app

import (
	"encoding/json"
	"strings"
	"time"
	"myapp/models"
)

func HL7toFHIR(hl7Message string) (string, error) {
	lines := strings.Split(hl7Message, "\n")

	// Define the data - note we only need to add the array elements as the single elements are already members
	data := models.HL7FHIRData{
        Medication:    []models.Medication{},
        Allergies:     []models.Allergy{},
        Conditions:    []models.Condition{},
        Observations:  []models.Observation{},
        Immunizations: []models.Immunization{},
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
			data.Patient = models.Patient{
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
                data.Medication = append(data.Medication, models.Medication{
                    Name:   medicationName,
                    Date:   date,
                    Dosage: segments[6],
                })
            } else {
                parts := strings.Split(medicationName, "^")
                immunization := models.Immunization{
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
                data.Allergies = append(data.Allergies, models.Allergy{
                    Name:        parseSubfield(segments[3], 1),
                    Criticality: map[string]string{"U": "low", "SV": "high", "MO": "moderate", "MI": "mild"}[segments[4]],
                    Date:        segments[6],
                })
            }
		case "DG1":
			if len(segments) > 5 {
                data.Conditions = append(data.Conditions, models.Condition{
                    Name: parseSubfield(segments[3], 1),
                    Date: segments[5],
                })
            }
		case "OBX":
			if len(segments) > 12 {
                data.Observations = append(data.Observations, models.Observation{
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

// Helper functions for HL7 parsing
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
	return t.Format("2006-01-02T15:04:05.000Z"), nil
}

func parseHL7Date(input string) (string, error) {
	t, err := time.Parse("20060102", input)
	if err != nil {
		return "", err
	}
	return t.Format(time.RFC3339), nil
}
