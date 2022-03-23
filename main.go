package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	elementDelimiter    = flag.String("component", "^", "HL7v2 Component Delimiter")
	subElementDelimiter = flag.String("subcomponent", "&", "HL7v2 Subcomponent Delimiter")
	segmentDelimiter    = flag.String("repetition", "~\r\n", "HL7v2 Repetition Delimiter")
	inputFile           = flag.String("input", "", "Input File")
	templateFile        = flag.String("template", "", "Output Template File")
)

func main() {

	flag.Parse()

	// TODO read encoding characters from command line args
	hl7Encoding := Hl7Encoding{Component: "^", SubComponent: "&", Repetition: "~"}

	if *inputFile == "" {
		ProcessStdin(hl7Encoding)
	} else {
		ProcessInputFile(*inputFile, hl7Encoding)
	}

}

func ProcessInputFile(fileName string, hl7encoding Hl7Encoding) {
	inputFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(inputFile)

	ProcessJsonInput(decoder, hl7encoding)
}

func ProcessStdin(hl7encoding Hl7Encoding) {
	decoder := json.NewDecoder(os.Stdin)
	ProcessJsonInput(decoder, hl7encoding)
}

func ProcessJsonInput(decoder *json.Decoder, hl7encoding Hl7Encoding) {

	// JSON array opening
	_, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	// Loop JSON Array Values
	for decoder.More() {
		var encounter Encounter

		err := decoder.Decode(&encounter)
		if err != nil {
			log.Fatal(err)
		}

		encounter.Hl7Encoding = hl7encoding
		fmt.Print(encounter.AsHl7())
	}

	// JSON array closing
	_, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
}
