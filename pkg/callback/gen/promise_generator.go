package main

import (
	"html/template"
	"log"
	"os"
)

//go:generate go run promise_generator.go

const promiseTemplate string = "promise.go.tmpl"

type Data struct {
	PackagedType string
	Type         string
	Prefix       string
}

type generator struct {
	templateFile string
	data         Data
	outputDir    string
	outputFile   string
}

func main() {
	generators := []generator{
		{
			templateFile: promiseTemplate,
			data:         Data{"pkg/beacon/chain/BeaconConfig", "BeaconConfig", "BeaconConfig"},
			outputDir:    "out",
			outputFile:   "beacon_config_promise.go",
		},
		{
			templateFile: promiseTemplate,
			data:         Data{"", "string", "String"},
			outputDir:    "out",
			outputFile:   "string_promise.go",
		},
	}

	for _, generator := range generators {
		// Read template.
		t, err := template.New(generator.templateFile).ParseFiles(generator.templateFile)
		if err != nil {
			log.Fatalf("template creation failed: %v", err)
		}

		// Create output file.
		os.Mkdir(generator.outputDir, 0700)
		outputPath := generator.outputDir + generator.outputFile
		os.Remove(outputPath)
		f, err := os.Create(outputPath)
		defer f.Close()

		// Generate files
		err = t.Execute(f, generator.data)
		if err != nil {
			log.Fatalf("files generation failed: %v", err)
		}
	}
}
