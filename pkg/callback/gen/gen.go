package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
)

// Promises code config.
// Execute `go generate` command in current directory to generate Promises code.

// Do not remove next comment!
//go:generate go run gen.go

const promiseTemplate string = "promise.go.tmpl"

const outDir string = "../"

// Data to be replaced in the `promises template`
type Data struct {
	// Type
	Type string
	// Package of the `Type`
	PackagedType string
	// Prefix for naming the premises
	Prefix string
	// Empty value for given `Type`
	NotDefinedValue string
}

// Configuration for the generator
type config struct {
	// Name of the promises template file
	templateFile string
	// Values to replace in the template file
	data Data
	// Directory to which generated code will be exported
	outputDir string
	// Name of the generated file
	outputFile string
}

//TODO move this executor to `main.go`
func main() {
	err := GeneratePromises()

	if err != nil {
		log.Fatalf("Error when generating promises\n%v", err)
	}
}

func GeneratePromises() error {
	configs := []config{
		// Promise for `string` type. There is a test for this promise `string_promise_test.go`. We need test to validate correctness of promises.
		{
			data: Data{
				PackagedType:    "",
				Type:            "string",
				Prefix:          "String",
				NotDefinedValue: ""},
			outputFile: "string_promise.go",
		},
		// This is just an example which should be replaced with types of needed promises
		{
			data: Data{
				PackagedType:    "pkg/beacon/chain",
				Type:            "chain.BeaconConfig",
				Prefix:          "BeaconConfig",
				NotDefinedValue: "chain.BeaconConfig{}"},
			outputFile: "beacon_config_promise.go",
		},
	}

	for _, c := range configs {
		// Read template
		t, err := template.New(promiseTemplate).ParseFiles(promiseTemplate)
		if err != nil {
			return fmt.Errorf("template creation failed: %v", err)
		}

		// Create output file
		outputPath := path.Join(outDir, c.outputFile)
		f, err := os.Create(outputPath)
		defer f.Close()

		// Generate files
		err = t.Execute(f, c.data)
		if err != nil {
			return fmt.Errorf("generation failed: %v", err)
		}
	}
	return nil
}
