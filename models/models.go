package models

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
