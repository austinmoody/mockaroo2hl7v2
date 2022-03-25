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
)

var (
	componentDelimiter    = flag.String("component", "^", "HL7v2 Component Delimiter")
	subComponentDelimiter = flag.String("subcomponent", "&", "HL7v2 Subcomponent Delimiter")
	repetitionDelimiter   = flag.String("repetition", "~\r\n", "HL7v2 Repetition Delimiter")
	inputFile             = flag.String("input", "", "Input File")
	templateFile          = flag.String("template", "", "Output Template File")
)

const JsonArray = "["
const JsonObject = "{"

func main() {

	flag.Parse()

	// TODO read encoding characters from command line args
	hl7encoding := Hl7Encoding{} //Hl7Encoding{Component: "^", SubComponent: "&", Repetition: "~"}
	hl7encoding.Component = *componentDelimiter
	hl7encoding.SubComponent = *subComponentDelimiter
	hl7encoding.Repetition = *repetitionDelimiter

	if *inputFile == "" {
		ProcessStdin(hl7encoding)
	} else {
		ProcessInputFile(*inputFile, hl7encoding)
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

func ProcessInputFile(fileName string, hl7encoding Hl7Encoding) {

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
		ProcessJsonInput(decoder, hl7encoding)
	} else {
		ProcessSingleJsonInput(decoder, hl7encoding)
	}
}

func ProcessStdin(hl7encoding Hl7Encoding) {

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
		ProcessJsonInput(decoder, hl7encoding)
	} else {
		ProcessSingleJsonInput(decoder, hl7encoding)
	}

}

func ProcessSingleJsonInput(decoder *json.Decoder, hl7encoding Hl7Encoding) {
	var encounter Encounter

	err := decoder.Decode(&encounter)
	if err != nil {
		log.Fatal(err)
	}

	encounter.Hl7Encoding = hl7encoding
	fmt.Printf("%s%s", encounter.AsHl7(), "\u0000")
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
		fmt.Printf("%s%s", encounter.AsHl7(), "\u0000")
	}

	// JSON array closing
	_, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
}
