package app

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	. "myapp/models"
)

// Converts from MongoDB JSON to FHiR JSON - we can chain this for HL7 to FHiR
func GenerateIPSBundle(ipsRecord HL7FHIRData) (string, error) {
	// Generate UUIDs - may drop this inline with same suggestion for web version - replace with simple ids
	compositionUUID := uuid.New().String()
	patientUUID := uuid.New().String()
	practitionerUUID := uuid.New().String()
	organizationUUID := uuid.New().String()

	// Current date/time for textual reference
	currentDateTime := time.Now().Format("2006-01-02T15:04:05.000Z")

	// Medications
	medications := []map[string]interface{}{}
	for _, med := range ipsRecord.Medication {
		medicationUUID := uuid.New().String()
		medications = append(medications, map[string]interface{}{
			"fullUrl": "urn:uuid:" + medicationUUID,
			"resource": map[string]interface{}{
				"resourceType": "Medication",
				"id":           medicationUUID,
				"code": map[string]interface{}{
					"coding": []map[string]interface{}{
						{"display": med.Name},
					},
				},
			},
		})
	}

	// MedicationStatements
	medicationStatements := []map[string]interface{}{}
	for i, med := range medications {
		medicationStatementUUID := uuid.New().String()
		medicationStatements = append(medicationStatements, map[string]interface{}{
			"fullUrl": "urn:uuid:" + medicationStatementUUID,
			"resource": map[string]interface{}{
				"resourceType": "MedicationStatement",
				"id":           medicationStatementUUID,
				"medicationReference": map[string]interface{}{
					"reference": "Medication/" + med["resource"].(map[string]interface{})["id"].(string),
					"display":   med["resource"].(map[string]interface{})["code"].(map[string]interface{})["coding"].([]map[string]interface{})[0]["display"].(string),
				},
				"subject": map[string]interface{}{
					"reference": "Patient/" + patientUUID,
				},
				"effectivePeriod": map[string]interface{}{
					"start": ipsRecord.Medication[i].Date,
				},
				"dosage": []map[string]interface{}{
					{"text": ipsRecord.Medication[i].Dosage},
				},
			},
		})
	}

	// AllergyIntolerances
	allergyIntolerances := []map[string]interface{}{}
	for _, allergy := range ipsRecord.Allergies {
		allergyIntoleranceUUID := uuid.New().String()
		allergyIntolerances = append(allergyIntolerances, map[string]interface{}{
			"fullUrl": "urn:uuid:" + allergyIntoleranceUUID,
			"resource": map[string]interface{}{
				"resourceType": "AllergyIntolerance",
				"id":           allergyIntoleranceUUID,
				"type":         "allergy",
				"category":     []string{"medication"},
				"criticality":  allergy.Criticality,
				"code": map[string]interface{}{
					"coding": []map[string]interface{}{
						{"display": allergy.Name},
					},
				},
				"patient": map[string]interface{}{
					"reference": "Patient/" + patientUUID,
				},
				"onsetDateTime": allergy.Date,
			},
		})
	}

	// Conditions
	conditions := []map[string]interface{}{}
	for _, condition := range ipsRecord.Conditions {
		conditionUUID := uuid.New().String()
		conditions = append(conditions, map[string]interface{}{
			"fullUrl": "urn:uuid:" + conditionUUID,
			"resource": map[string]interface{}{
				"resourceType": "Condition",
				"id":           conditionUUID,
				"code": map[string]interface{}{
					"coding": []map[string]interface{}{
						{"display": condition.Name},
					},
				},
				"subject": map[string]interface{}{
					"reference": "Patient/" + patientUUID,
				},
				"onsetDateTime": condition.Date,
			},
		})
	}

	// Observations
	observations := []map[string]interface{}{}
	for _, observation := range ipsRecord.Observations {
		observationUUID := uuid.New().String()
		observations = append(observations, map[string]interface{}{
			"fullUrl": "urn:uuid:" + observationUUID,
			"resource": map[string]interface{}{
				"resourceType": "Observation",
				"id":           observationUUID,
				"code": map[string]interface{}{
					"coding": []map[string]interface{}{
						{"display": observation.Name},
					},
				},
				"subject": map[string]interface{}{
					"reference": "Patient/" + patientUUID,
				},
				"effectiveDateTime": observation.Date,
				"valueString":       observation.Value,
			},
		})
	}

	// Immunizations
	immunizations := []map[string]interface{}{}
	for _, immunization := range ipsRecord.Immunizations {
		immunizationUUID := uuid.New().String()
		immunizations = append(immunizations, map[string]interface{}{
			"fullUrl": "urn:uuid:" + immunizationUUID,
			"resource": map[string]interface{}{
				"resourceType": "Immunization",
				"id":           immunizationUUID,
				"status":       "completed",
				"vaccineCode": map[string]interface{}{
					"coding": []map[string]interface{}{
						{"system": immunization.System, "code": immunization.Name},
					},
				},
				"patient": map[string]interface{}{
					"reference": "Patient/" + patientUUID,
				},
				"occurrenceDateTime": immunization.Date,
			},
		})
	}

	// Composition
composition := map[string]interface{}{
	"fullUrl": "urn:uuid:" + compositionUUID,
	"resource": map[string]interface{}{
		"resourceType": "Composition",
		"id":           compositionUUID,
		"type": map[string]interface{}{
			"coding": []map[string]interface{}{
				{
					"system": "http://loinc.org",
					"code":   "60591-5",
					"display": "Patient summary Document",
				},
			},
		},
		"subject": map[string]interface{}{
			"reference": "Patient/" + patientUUID,
		},
		"date":  currentDateTime,
		"title": "Patient Summary as of " + currentDateTime,
		"author": []map[string]interface{}{
			{
				"reference": "Practitioner/" + practitionerUUID,
			},
		},
		"custodian": map[string]interface{}{
			"reference": "Organization/" + organizationUUID,
		},
		"section": []map[string]interface{}{
			{
				"title": "Medication",
				"code": map[string]interface{}{
					"coding": []map[string]interface{}{
						{
							"system": "http://loinc.org",
							"code":   "10160-0",
							"display": "History of Medication use Narrative",
						},
					},
				},
				"entry": func() []map[string]interface{} {
					entries := []map[string]interface{}{}
					for _, medStatement := range medicationStatements {
						entries = append(entries, map[string]interface{}{
							"reference": "MedicationStatement/" + medStatement["resource"].(map[string]interface{})["id"].(string),
						})
					}
					return entries
				}(),
			},
			{
				"title": "Allergies and Intolerances",
				"code": map[string]interface{}{
					"coding": []map[string]interface{}{
						{
							"system": "http://loinc.org",
							"code":   "48765-2",
							"display": "Allergies and adverse reactions Document",
						},
					},
				},
				"entry": func() []map[string]interface{} {
					entries := []map[string]interface{}{}
					for _, allergy := range allergyIntolerances {
						entries = append(entries, map[string]interface{}{
							"reference": "AllergyIntolerance/" + allergy["resource"].(map[string]interface{})["id"].(string),
						})
					}
					return entries
				}(),
			},
			{
				"title": "Conditions",
				"code": map[string]interface{}{
					"coding": []map[string]interface{}{
						{
							"system": "http://loinc.org",
							"code":   "11450-4",
							"display": "Problem List",
						},
					},
				},
				"entry": func() []map[string]interface{} {
					entries := []map[string]interface{}{}
					for _, condition := range conditions {
						entries = append(entries, map[string]interface{}{
							"reference": "Condition/" + condition["resource"].(map[string]interface{})["id"].(string),
						})
					}
					return entries
				}(),
			},
			{
				"title": "Observations",
				"code": map[string]interface{}{
					"coding": []map[string]interface{}{
						{
							"system": "http://loinc.org",
							"code":   "61150-9",
							"display": "Vital signs, weight, length, head circumference, oxygen saturation and BMI Panel",
						},
					},
				},
				"entry": func() []map[string]interface{} {
					entries := []map[string]interface{}{}
					for _, observation := range observations {
						entries = append(entries, map[string]interface{}{
							"reference": "Observation/" + observation["resource"].(map[string]interface{})["id"].(string),
						})
					}
					return entries
				}(),
			},
			{
				"title": "Immunizations",
				"code": map[string]interface{}{
					"coding": []map[string]interface{}{
						{
							"system": "http://loinc.org",
							"code":   "11369-6",
							"display": "Immunization Activity",
						},
					},
				},
				"entry": func() []map[string]interface{} {
					entries := []map[string]interface{}{}
					for _, immunization := range immunizations {
						entries = append(entries, map[string]interface{}{
							"reference": "Immunization/" + immunization["resource"].(map[string]interface{})["id"].(string),
						})
					}
					return entries
				}(),
			},
		},
	},
}


	// Construct FHIR Bundle
fhirBundle := map[string]interface{}{
	"resourceType": "Bundle",
	"id":           ipsRecord.PackageUUID,
	"type":         "document",
	"timestamp":    ipsRecord.TimeStamp,
	"entry": append(
		[]map[string]interface{}{
			// Composition Resource
			composition,
			// Patient Resource
			{
				"fullUrl": "urn:uuid:" + patientUUID,
				"resource": map[string]interface{}{
					"resourceType": "Patient",
					"id":           patientUUID,
					"name": []map[string]interface{}{
						{"family": ipsRecord.Patient.Name, "given": []string{ipsRecord.Patient.Given}},
					},
					"gender":    ipsRecord.Patient.Gender,
					"birthDate": ipsRecord.Patient.DOB,
					"address":   []map[string]interface{}{{"country": ipsRecord.Patient.Nation}},
				},
			},
		},
		mergeResources(
			medicationStatements,
			medications,
			allergyIntolerances,
			conditions,
			observations,
			immunizations,
		)...,
	),
}

	// Convert to JSON
	fhirJSON, err := json.MarshalIndent(fhirBundle, "", "  ")
	if err != nil {
		return "", err
	}

	return string(fhirJSON), nil
}

// Merge resources
func mergeResources(resources ...[]map[string]interface{}) []map[string]interface{} {
	merged := []map[string]interface{}{}
	for _, resourceGroup := range resources {
		merged = append(merged, resourceGroup...)
	}
	return merged
}
