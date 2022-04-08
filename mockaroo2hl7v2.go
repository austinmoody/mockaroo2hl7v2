package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	inputFile    = flag.String("input", "", "Input File.  STDIN used if not specified")
	templateFile = flag.String("template", "", "Output Template File.  Default can be found in main.go")
	mllp         = flag.Bool("mllp", false, "Wrap each message in MLLP envelope?  Default is false.")
	print0       = flag.Bool("print0", false, "End each message with ASCII NUL?  Default is false.")
)

const JsonArray = "["
const JsonObject = "{"
const MllpStart = "\x0b"
const MllpEnd = "\x1c\x0d"
const NulChar = "\u0000"

func main() {

	flag.Parse()

	templateContent := TemplateContents()
	template := NewOutputTemplate(templateContent)
	encounter := NewEncounter(template)

	if *inputFile == "" {
		ProcessStdin(&encounter, *mllp, *print0)
	} else {
		ProcessInputFile(*inputFile, &encounter, *mllp, *print0)
	}
}

func ValidJson(rdr io.Reader) (bool, string) {
	// rename this and pass back err instead of bool...
	decoder := json.NewDecoder(rdr)

	firstToken, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	if delimiter, success := firstToken.(json.Delim); success {
		switch delimiter.String() {
		case JsonArray:
			return true, JsonArray
		case JsonObject:
			return true, JsonObject
		default:
			return false, ""
		}
	}

	return false, ""
}

func ProcessInputFile(fileName string, encounter *Encounter, mllp bool, print0 bool) {

	inputFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	jsonValid, jsonType := ValidJson(inputFile)

	if jsonValid != true {
		log.Fatal("Input JSON is invalid")
	}

	// reset so we can decode full input
	_, err = inputFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(inputFile)

	if jsonType == JsonArray {
		ProcessJsonInput(decoder, encounter, mllp, print0)
	} else {
		ProcessSingleJsonInput(decoder, encounter, mllp, print0)
	}
}

func ProcessStdin(encounter *Encounter, mllp bool, print0 bool) {

	jsonType := ""
	jsonValid := false

	buf := bytes.Buffer{}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if jsonType == "" {
			jsonValid, jsonType = ValidJson(bytes.NewReader(s.Bytes()))
			if jsonValid != true {
				log.Fatal("Input JSON is invalid")
			}
		}
		buf.Write(s.Bytes())
	}

	decoder := json.NewDecoder(bufio.NewReader(&buf))

	if jsonType == JsonArray {
		ProcessJsonInput(decoder, encounter, mllp, print0)
	} else {
		ProcessSingleJsonInput(decoder, encounter, mllp, print0)
	}
}

func ProcessSingleJsonInput(decoder *json.Decoder, encounter *Encounter, mllp bool, print0 bool) {
	//var encounter Encounter
	//encounter := NewEncounter()

	err := decoder.Decode(&encounter)
	if err != nil {
		log.Fatal(err)
	}

	mllpStart := ""
	mllpEnd := ""
	if mllp {
		mllpStart = MllpStart
		mllpEnd = MllpEnd
	}

	printZero := ""
	if print0 {
		printZero = NulChar
	}

	//encounter.Hl7Encoding = hl7template.encoding
	fmt.Printf("%s%s%s%s", mllpStart, encounter.AsHl7(), mllpEnd, printZero)
}

func ProcessJsonInput(decoder *json.Decoder, encounter *Encounter, mllp bool, print0 bool) {

	mllpStart := ""
	mllpEnd := ""
	if mllp {
		mllpStart = MllpStart
		mllpEnd = MllpEnd
	}

	printZero := ""
	if print0 {
		printZero = NulChar
	}

	// JSON array opening
	_, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	// Loop JSON Array Values
	for decoder.More() {
		//var encounter Encounter

		err := decoder.Decode(&encounter)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s%s%s%s", mllpStart, encounter.AsHl7(), mllpEnd, printZero)
	}

	// JSON array closing
	_, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
}

func TemplateContents() string {
	templateContent := ""
	if *templateFile == "" {
		templateContent = DefaultTemplate()
	} else {
		templateFileContents, err := os.Open(*templateFile)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(templateFileContents)
		for scanner.Scan() {
			templateContent += scanner.Text()
		}
	}

	return templateContent
}

func DefaultTemplate() string {
	msh := "MSH|^~\\&|Mirth|Hospital|HIE|HIE|{{ .DateTime.Format \"20060102150405\" }}||ADT^{{.Event}}|{{ .Utility.Uuid }}|P|2.5.1"
	evn := "EVN|{{.Event}}|{{ .DateTime.Format \"20060102150405\" }}"
	pid := "PID|1||{{ .Patient.IdentifiersAsHl7 .Output.Encoding }}||{{.Patient.Name.AsHl7 .Output.Encoding }}|{{ .Patient.MotherMaidenName }}|{{ .Patient.DOB.Format \"20060102\" }}|{{.Patient.Gender}}||{{.Patient.Race.AsHl7 .Output.Encoding}}|{{.Patient.AddressesAsHl7 .Output.Encoding}}||{{.Patient.HomePhoneNumbersAsHl7 .Output.Encoding}}|{{.Patient.WorkPhoneNumbersAsHl7 .Output.Encoding}}|{{.Patient.Language.Identifier}}|{{.Patient.MaritalStatus.AsHl7 .Output.Encoding}}|{{.Patient.Religion.Identifier}}|{{.Patient.IdentifierAsHl7 \"AN\" .Output.Encoding}}|{{.Patient.SSN}}|{{.Patient.DriversLicense.AsHl7 .Output.Encoding}}||{{.Patient.EthnicGroup.Identifier}}||{{.Patient.MultipleBirthIndicator}}|{{.Patient.BirthOrder}}||||{{.Patient.DeathDateTime.Format \"20060102\"}}|{{.Patient.DeathIndicator}}"
	pv1 := "PV1|1|{{ .Class.Identifier }}"

	segments := []string{msh, evn, pid, pv1}

	return strings.Join(segments, "\x0d")
}
