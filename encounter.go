package main

import (
	"bytes"
	"strings"
	"text/template"
	"time"
)

type Hl7Encoding struct {
	Component    string `default:"^"`
	SubComponent string `default:"&"`
	Repetition   string `default:"~"`
}

type Encounter struct {
	Hl7Encoding  Hl7Encoding
	Patient      Patient      `json:"Patient"`
	Providers    []Provider   `json:"Providers"`
	PatientClass CodedElement `json:"PatientClass"`
}

func (e Encounter) AsHl7() string {

	templateString := GetTemplate()

	template, err := template.New("").Parse(templateString)
	//template, err := template.ParseFiles("default.gohl7")
	if err != nil {
		panic(err)
	}

	hl7 := new(bytes.Buffer)

	template.Execute(hl7, e)

	return hl7.String()
}

type Patient struct {
	// PID
	Ids                    []Identifier                `json:"Ids"`
	Name                   Name                        `json:"Name"`
	MotherMaidenName       string                      `json:"MotherMaidenName"`
	DOB                    time.Time                   `json:"DOB"`
	Gender                 string                      `json:"Gender"`
	Race                   CodedElement                `json:"Race"`
	Addresses              []ExtendedAddress           `json:"Address"`
	Phones                 []ExtendedTelecommunication `json:"Comm"`
	Language               CodedElement                `json:"Language"`
	MaritalStatus          CodedElement                `json:"MaritalStatus"`
	Religion               CodedElement                `json:"Religion"`
	SSN                    string                      `json:"SSN"`
	DriversLicense         DriversLicense              `json:"DriversLicense"`
	EthnicGroup            CodedElement                `json:"EthnicGroup"`
	MultipleBirthIndicator string                      `json:"MultipleBirthIndicator"`
	BirthOrder             int                         `json:"BirthOrder"`
	DeathIndicator         string                      `json:"DeathIndicator"`
	DeathDateTime          time.Time                   `json:"DeathDateTime"`
}

func (p Patient) IdentifierAsHl7(idTypeCode string, encoding Hl7Encoding) string {
	var id Identifier
	for i := range p.Ids {
		if p.Ids[i].IdType == idTypeCode {
			id = p.Ids[i]
			break
		}
	}
	return id.AsHl7(encoding)
}

func (p Patient) IdentifiersAsHl7(encoding Hl7Encoding) string {
	var identifiers []string
	for _, id := range p.Ids {
		identifiers = append(identifiers, id.AsHl7(encoding))
	}

	return strings.Join(identifiers, encoding.Repetition)
}

type Name struct {
	Family string `json:"Family"`
	Given  string `json:"Given"`
	Second string `json:"Second"`
	Suffix string `json:"Suffix"`
	Prefix string `json:"Prefix"`
}

func (n Name) AsHl7(encoding Hl7Encoding) string {
	names := []string{n.Family, n.Given, n.Second, n.Suffix, n.Prefix}
	return strings.Join(names, encoding.Component)
}

type Identifier struct {
	IdNumber         string `json:"IdNumber"`
	CheckDigit       string `json:"CheckDigit"`
	CheckDigitScheme string `json:"CheckDigitScheme"`
	IdType           string `json:"IdType"`
}

func (id Identifier) AsHl7(encoding Hl7Encoding) string {
	idSlice := []string{id.IdNumber, id.CheckDigit, id.CheckDigitScheme, "", id.IdType}
	return strings.Join(idSlice, encoding.Component)
}

type CodedElement struct {
	// CE HL7v2 Data Type
	Identifier   string `json:"Id"`
	Text         string `json:"Text"`
	CodingSystem string `json:"CodingSystem"`
}

func (ce CodedElement) AsHl7(encoding Hl7Encoding) string {
	ces := []string{ce.Identifier, ce.Text, ce.CodingSystem}
	return strings.Join(ces, encoding.Component)
}

type ExtendedAddress struct {
	// XAD HL7v2 Data Type - yes I know the address lines aren't named as in spec :)
	// also haven't filled out for all possible fields for XAD because I haven't
	// needed them yet.
	Line1   string `json:"Line1"`
	Line2   string `json:"Line2"`
	City    string `json:"City"`
	State   string `json:"State"`
	Zip     string `json:"Zip"`
	Country string `json:"Country"`
	Type    string `json:"Type"`
}

type ExtendedTelecommunication struct {
	// XTN HL7v2 Data Type
	TelephoneNumber            string `json:"PhoneNumber"`
	UseCode                    string `json:"UseCode"`
	EquipmentType              string `json:"Type"`
	EmailAddress               string `json:"Email"`
	CountryCode                string `json:"CountryCode"`
	AreaCode                   string `json:"AreaCode"`
	LocalNumber                string `json:"LocalNumber"`
	Extension                  string `json:"Extension"`
	AnyText                    string `json:"AnyText"`
	ExtensionPrefix            string `json:"ExtensionPrefix"`
	SpeedDialCode              string `json:"SpeedDialCode"`
	UnformattedTelephoneNumber string `json:"UnformattedTelephoneNumber"`
}

type DriversLicense struct {
	LicenseNumber  string    `json:"Number"`
	State          string    `json:"State"`
	ExpirationDate time.Time `json:"ExpirationDate"`
}

type Provider struct {
	Identifier Identifier   `json:"Id"`
	Name       Name         `json:"Name"`
	Degree     string       `json:"Degree"`
	Role       CodedElement `json:"Role"`
}

func GetTemplate() string {
	// TODO - obviously allow user to specify external template
	return `PID|||{{ .Patient.IdentifiersAsHl7 .Hl7Encoding }}||{{.Patient.Name.AsHl7 .Hl7Encoding }}|{{ .Patient.MotherMaidenName }}|{{ .Patient.DOB.Format "20060102" }}|{{.Patient.Gender}}||{{.Patient.Race.AsHl7 .Hl7Encoding}}|pickup w/ address`
}
