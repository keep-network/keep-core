//go:generate sh -c "SOLIDITY_DIR=../../../contracts/solidity make"
// Code generation execution command requires the package to be set to `main`.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/tools/imports"
)

// Main function. Expects to be invoked as:
//
//   <executable> [input.abi] contract/[contract_output.go] cmd/[cmd_output.go]
//
// The first file will receive a contract binding that is slightly higher-level
// than abigen's output, including an event-based interface for contract event
// interaction, support for revert error reporting, serialized transaction
// submission, and simplified transactor handling.
//
// The second file will receive an urfave/cli-compatible command initialization
// that can be used to add command-line interaction with the specified contract
// by adding the relevant commands to a top-level urfave/cli.App object.
func main() {
	if len(os.Args) != 4 {
		panic(fmt.Sprintf(
			"Expected `%v [input.abi] [contract_output.go] [cmd_output.go]`, but got [%v].",
			os.Args[0],
			os.Args,
		))
	}

	abiPath := os.Args[1]
	contractOutputPath := os.Args[2]
	commandOutputPath := os.Args[3]

	abiFile, err := ioutil.ReadFile(abiPath)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to read ABI file at [%v]: [%v].",
			abiPath,
			err,
		))
	}

	templates, err := template.ParseGlob("*.go.tmpl")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse templates: [%v].", err))
	}

	abi, err := abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to parse ABI at [%v]: [%v].",
			abiPath,
			err,
		))
	}

	var payableInfo []methodPayableInfo
	err = json.Unmarshal(abiFile, &payableInfo)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to parse additional ABI metadata at [%v]: [%v].",
			abiPath,
			err,
		))
	}

	// The name of the ABI binding Go class is the same as the filename of the
	// ABI file, minus the extension.
	abiClassName := path.Base(abiPath)
	abiClassName = abiClassName[0 : len(abiClassName)-4] // strip .abi
	contractInfo := buildContractInfo(abiClassName, &abi, payableInfo)

	contractBuf, err := generateCode(
		contractOutputPath,
		templates,
		"contract.go.tmpl",
		&contractInfo,
	)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to generate Go file for contract [%v] at [%v]: [%v].",
			contractInfo.AbiClass,
			contractOutputPath,
			err,
		))
	}

	commandBuf, err := generateCode(
		commandOutputPath,
		templates,
		"command.go.tmpl",
		&contractInfo,
	)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to generate Go file at [%v]: [%v].",
			commandOutputPath,
			err,
		))
	}

	// Save the contract code to a file.
	if err := saveBufferToFile(contractBuf, contractOutputPath); err != nil {
		panic(fmt.Sprintf(
			"Failed to save Go file at [%v]: [%v].",
			contractOutputPath,
			err,
		))
	}

	// Save the command code to a file.
	if err := saveBufferToFile(commandBuf, commandOutputPath); err != nil {
		panic(fmt.Sprintf(
			"Failed to save Go file at [%v]: [%v].",
			commandOutputPath,
			err,
		))
	}
}

// Generates code by applying the named template in the passed template bundle
// to the specified data object. Writes the output to a buffer and then
// formats and organizes the imports on that buffer, returning the final result
// ready for emission onto the filesystem.
//
// Note that this means the generated file must compile, or import organization
// will fail. The error message in case of compilation failure will be bubbled
// up, but the file contents currently will not be written.
func generateCode(
	outFile string,
	templat *template.Template,
	templateName string,
	data interface{},
) (*bytes.Buffer, error) {
	var buffer bytes.Buffer

	if err := templat.ExecuteTemplate(&buffer, templateName, data); err != nil {
		return nil, fmt.Errorf(
			"generating code failed: [%v]",
			err,
		)
	}

	if err := organizeImports(outFile, &buffer); err != nil {
		return nil, err
	}

	return &buffer, nil
}

// Resolves imports in a code stored in a Buffer.
func organizeImports(outFile string, buf *bytes.Buffer) error {
	// Resolve imports
	code, err := imports.Process(outFile, buf.Bytes(), nil)
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
