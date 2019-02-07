//go:generate sh -c "rm -f ./async/*_promise.go; go run async.go"
// Code generation execution command requires the package to be set to `main`.
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"

	"golang.org/x/tools/imports"
)

// Promises code generator.
// Execute `go generate` command in current directory to generate Promises code.

// Name of the promise template file.
const promiseTemplateFile string = "async_promise.go.tmpl"

// Directory to which generated code will be exported.
const outDir string = "./async"

// Configuration for the generator
type promiseConfig struct {
	// Type which promise will handle.
	Type string
	// Prefix for naming the promises.
	Prefix string
	// Name of the generated file.
	outputFile string
}

func main() {
	configs := []promiseConfig{
		// Promise for `*big.Int` type.
		// There is a test for this promise named `big_int_promise_test.go`.
		// We need a test to validate correctness of generated promises.
		{
			Type:       "*big.Int",
			Prefix:     "BigInt",
			outputFile: "big_int_promise.go",
		},
		{
			Type:       "*event.Entry",
			Prefix:     "RelayEntry",
			outputFile: "relay_entry_promise.go",
		},
		{
			Type:       "*event.GroupRegistration",
			Prefix:     "GroupRegistration",
			outputFile: "group_registration_promise.go",
		},
		{
			Type:       "*event.Request",
			Prefix:     "RelayRequest",
			outputFile: "relay_entry_requested_promise.go",
		},
		{
			Type:       "*event.DKGResultPublication",
			Prefix:     "DKGResultPublication",
			outputFile: "dkg_result_publication_promise.go",
		},
		{
			Type:       "*event.GroupTicketSubmission",
			Prefix:     "GroupTicket",
			outputFile: "group_ticket_submission_promise.go",
		},
	}

	if err := generatePromisesCode(configs); err != nil {
		log.Fatalf("promises generation failed [%v]", err)
	}
}

// Generates promises based on a given `promiseConfig`
func generatePromisesCode(promisesConfig []promiseConfig) error {
	// Read a promise's template.
	promiseTemplate, err := template.New(promiseTemplateFile).ParseFiles(promiseTemplateFile)
	if err != nil {
		return fmt.Errorf("template creation failed [%v]", err)
	}

	for _, promiseConfig := range promisesConfig {
		// Generate a promise code.
		buf, err := generateCode(promiseTemplate, &promiseConfig)
		if err != nil {
			return fmt.Errorf("generation failed [%v]", err)
		}

		// Save the promise code to a file.
		if err := saveBufferToFile(buf, path.Join(outDir, promiseConfig.outputFile)); err != nil {
			return fmt.Errorf("saving promise code to file failed [%v]", err)
		}
	}
	return nil
}

// Generates a code from template and configuration.
// Returns a buffered code.
func generateCode(tmpl *template.Template, config *promiseConfig) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, fmt.Errorf("generating code for type %s failed [%v]", config.Type, err)
	}

	if err := organizeImports(&buf); err != nil {
		return nil, err
	}

	return &buf, nil
}

// Resolves imports in a code stored in a Buffer.
func organizeImports(buf *bytes.Buffer) error {
	// Resolve imports
	code, err := imports.Process(outDir, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("failed to find/resove imports [%v]", err)
	}

	// Write organized code to the buffer.
	buf.Reset()
	if _, err := buf.Write(code); err != nil {
		return fmt.Errorf("cannot write code to buffer [%v]", err)
	}

	return nil
}

// Stores the Buffer `buf` content to a file in `filePath`
func saveBufferToFile(buf *bytes.Buffer, filePath string) error {
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("output file %s creation failed [%v]", filePath, err)
	}

	if _, err := buf.WriteTo(file); err != nil {
		return fmt.Errorf("writing to output file %s failed [%v]", filePath, err)
	}

	return nil
}
