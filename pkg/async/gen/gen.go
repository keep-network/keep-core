package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path"

	"golang.org/x/tools/imports"
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
	// Need to be provided only if package cannot be resolved by `imports.Process` function
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
			Type:       "string",
			Prefix:     "String",
			outputFile: "string_promise.go",
		},
		// Example for types from packages which need to be imported.
		// Provide `PackagedType` only if package cannot be resolved by the `imports.Process`
		// {
		// 	Type:         "ethereum.KeepRandomBeacon",
		// 	PackagedType: "github.com/keep-network/keep-core/pkg/chain/ethereum",
		// 	Prefix:       "KeepRandomBeacon",
		// 	outputFile:   "keep_random_beacon_promise.go",
		// },
	}

	for _, c := range configs {
		// Read a Promise's Template
		t, err := template.New(promiseTemplate).ParseFiles(promiseTemplate)
		if err != nil {
			log.Fatalf("template creation failed [%v]", err)
		}

		// Generate a Promise
		var buf bytes.Buffer
		err = t.Execute(&buf, c)
		if err != nil {
			log.Fatalf("generation failed [%v]", err)
		}

		organizeImports(&buf)

		// Store the Promise in a file
		f, err := os.Create(path.Join(outDir, c.outputFile))
		defer f.Close()
		if err != nil {
			log.Fatalf("output file creation failed [%v]", err)
		}
		buf.WriteTo(f)
	}
}

// Resolves imports in a code stored in a Buffer
func organizeImports(buf *bytes.Buffer) {
	code, err := imports.Process(outDir, buf.Bytes(), nil)
	if err != nil {
		log.Fatalf("failed to find/resove imports [%v]\nCode: %s", err, code)
	}

	// Write organized code to the buffer
	buf.Reset()
	length, err := buf.Write(code)
	if err != nil {
		log.Fatalf("cannot write code to buffer [%v]\nCode length: %d\nCode: %s", err, length, code)
	}
}
