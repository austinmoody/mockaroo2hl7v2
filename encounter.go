package main

import (
	"bytes"
	"github.com/google/uuid"
	"log"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type Hl7Encoding struct {
	Field        string `default:"|"`
	Component    string `default:"^"`
	SubComponent string `default:"&"`
	Repetition   string `default:"~"`
	Escape       string `default:"\\"`
}

type Utility string

func (u Utility) Uuid() string {
	return uuid.New().String()
}

type Encounter struct {
	Output    OutputTemplate
	Patient   Patient      `json:"Patient"`
	Providers []Provider   `json:"Providers"`
	Class     CodedElement `json:"PatientClass"`
	Event     string       `json:"Event"`
	DateTime  time.Time
	Utility   Utility
}

func NewEncounter(hl7template OutputTemplate) Encounter {
	return Encounter{DateTime: time.Now(), Output: hl7template}
}

func (e Encounter) AsHl7() string {

	hl7 := new(bytes.Buffer)

	err := e.Output.Template.Execute(hl7, e)
	if err != nil {
		log.Fatal(err)
	}

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

func (p Patient) PhoneNumbersAsHl7(phones []ExtendedTelecommunication, encoding Hl7Encoding) string {
	var phonesHl7 []string
	for _, phone := range phones {
		phonesHl7 = append(phonesHl7, phone.AsHl7(encoding))
	}

	return strings.Join(phonesHl7, encoding.Repetition)
}

func (p Patient) HomePhoneNumbersAsHl7(encoding Hl7Encoding) string {
	phones := p.PhonesByUseCode("PRN")
	phones = append(phones, p.PhonesByUseCode("ORN")...)
	phones = append(phones, p.PhonesByUseCode("NET")...)

	return p.PhoneNumbersAsHl7(phones, encoding)
}

func (p Patient) WorkPhoneNumbersAsHl7(encoding Hl7Encoding) string {
	phones := p.PhonesByUseCode("WPN")
	return p.PhoneNumbersAsHl7(phones, encoding)
}

func (p Patient) PhonesByUseCode(useCode string) []ExtendedTelecommunication {
	var phones []ExtendedTelecommunication

	for i := range p.Phones {
		if p.Phones[i].UseCode == useCode {
			phones = append(phones, p.Phones[i])
		}
	}

	return phones
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

func (p Patient) AddressesAsHl7(encoding Hl7Encoding) string {
	var addresses []string
	for _, address := range p.Addresses {
		addresses = append(addresses, address.AsHl7(encoding))
	}

	return strings.Join(addresses, encoding.Repetition)
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

func (ea ExtendedAddress) AsHl7(encoding Hl7Encoding) string {
	eas := []string{ea.Line1, ea.Line2, ea.City, ea.State, ea.Zip, ea.Country, ea.Type}
	return strings.Join(eas, encoding.Component)
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

func (et ExtendedTelecommunication) AsHl7(encoding Hl7Encoding) string {
	ets := []string{
		et.TelephoneNumber,
		et.UseCode,
		et.EquipmentType,
		et.EmailAddress,
		et.CountryCode,
		et.AreaCode,
		et.LocalNumber,
		et.Extension,
		et.AnyText,
		et.ExtensionPrefix,
		et.SpeedDialCode,
		et.UnformattedTelephoneNumber,
	}

	return strings.Join(ets, encoding.Component)
}

type DriversLicense struct {
	LicenseNumber  string    `json:"Number"`
	State          string    `json:"State"`
	ExpirationDate time.Time `json:"ExpirationDate"`
}

func (dl DriversLicense) AsHl7(encoding Hl7Encoding) string {
	dls := []string{dl.LicenseNumber, dl.State, dl.ExpirationDate.Format("20060102")}
	return strings.Join(dls, encoding.Component)
}

type Provider struct {
	Identifier Identifier   `json:"Id"`
	Name       Name         `json:"Name"`
	Degree     string       `json:"Degree"`
	Role       CodedElement `json:"Role"`
}

type OutputTemplate struct {
	Encoding Hl7Encoding
	Contents string
	Template *template.Template
}

func NewOutputTemplate(templateContents string) OutputTemplate {
	t := OutputTemplate{}
	t.Encoding = GetHl7Encoding(templateContents)
	t.Contents = templateContents
	var err error
	t.Template, err = template.New("").Parse(templateContents)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func GetHl7Encoding(templateContents string) Hl7Encoding {
	var encoding = Hl7Encoding{}
	rgx := regexp.MustCompile(`^(?P<segment>MSH)(?P<field>.)(?P<component>.)(?P<repetition>.)(?P<escape>.)(?P<subcomponent>.)(?P<field2>.)`)

	matches := rgx.FindStringSubmatch(templateContents)

	segment := rgx.SubexpIndex("segment")
	field := rgx.SubexpIndex("field")
	component := rgx.SubexpIndex("component")
	repetition := rgx.SubexpIndex("repetition")
	escape := rgx.SubexpIndex("escape")
	subcomponent := rgx.SubexpIndex("subcomponent")
	field2 := rgx.SubexpIndex("field2")

	// We expect the segment to be MSH and the first match after the segment and the last matched character to be the same
	// So given MSH|^~\&|... we expect 1st after the segment and last matched to be |
	if matches != nil && matches[segment] == "MSH" && matches[field] == matches[field2] {
		encoding.Component = matches[component]
		encoding.SubComponent = matches[subcomponent]
		encoding.Repetition = matches[repetition]
		encoding.Escape = matches[escape]
		encoding.Field = matches[field]
	}

	return encoding
}
