package main

import (
	"html/template"
	"log"
	"os"
	"path"
)

// Promises code config.
// Execute `go generate` command in current directory to generate Promises code.

// Do not remove next comment!
//go:generate go run gen.go

// Name of the promises template file
const promiseTemplate string = "promise.go.tmpl"

// Directory to which generated code will be exported
const outDir string = "../"

// Configuration for the generator
type config struct {
	// Type which Promise will handle
	Type string
	// Package of the `Type`
	PackagedType string
	// Prefix for naming the promises
	Prefix string
	// Name of the generated file
	outputFile string
}

func main() {
	configs := []config{
		// Promise for `string` type. There is a test for this promise named `string_promise_test.go`.
		// We need a test to validate correctness of generated promises.
		{
			Type:         "string",
			PackagedType: "",
			Prefix:       "String",
			outputFile:   "string_promise.go",
		},
		// Example for types from packages which need to be imported:
		// {
		// 	Type:         "ethereum.KeepRandomBeacon",
		// 	PackagedType: "github.com/keep-network/keep-core/pkg/chain/ethereum",
		// 	Prefix:       "KeepRandomBeacon",
		// 	outputFile:   "keep_random_beacon_promise.go",
		// },
	}

	for _, c := range configs {
		// Read template
		t, err := template.New(promiseTemplate).ParseFiles(promiseTemplate)
		if err != nil {
			log.Fatalf("template creation failed [%v]", err)
		}

		// Create output file
		outputPath := path.Join(outDir, c.outputFile)
		f, err := os.Create(outputPath)
		defer f.Close()
		if err != nil {
			log.Fatalf("output file creation failed [%v]", err)
		}

		// Generate files
		err = t.Execute(f, c)
		if err != nil {
			log.Fatalf("generation failed [%v]", err)
		}
	}
}
