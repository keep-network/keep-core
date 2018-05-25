package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const defaultCLIBinaryName = "./keep-core"
const defaultVerboseLogging = false

var (
	cliBinaryName  string
	verboseLogging bool
)

//-------------------------------------------------------------------------------
// Types and Methods
//-------------------------------------------------------------------------------
type validateStringType struct {
	keepEthPwd string
	//cmdPath     string  // <- use cliBinaryName
	cmdAndFlags string
	want        string
}

func (data *validateStringType) validateString() (got string) {
	cmd := exec.Command(cliBinaryName, data.cmdAndFlags)
	cmd.Env = append(os.Environ(), "KEEP_ETHEREUM_PASSWORD="+data.keepEthPwd)
	out, _ := cmd.Output()
	got = fmt.Sprintf("%s", out)
	got = strings.TrimSuffix(got, "\n")
	return got
}

type validateIntType struct {
	keepEthPwd  string
	cmdAndFlags string
	want        int
}

func (data *validateIntType) validateInt() (got int) {
	cmd := exec.Command(cliBinaryName, data.cmdAndFlags)
	cmd.Env = append(os.Environ(), "KEEP_ETHEREUM_PASSWORD="+data.keepEthPwd)
	out, _ := cmd.Output()
	if verboseLogging {
		fmt.Printf("out:\n%s\n", out)
	} else {
		fmt.Printf("")
	}
	got = countOccurences(fmt.Sprintf("%s", out))
	return
}

//-------------------------------------------------------------------------------
// Main
//-------------------------------------------------------------------------------
func TestMain(m *testing.M) {
	setup()

	code := m.Run()
	os.Exit(code)

	teardown()
}

func setup() {
	cliBinaryName = os.Getenv("CLI_BINARY")
	if len(cliBinaryName) == 0 {
		cliBinaryName = defaultCLIBinaryName
	}
	verboseLogging = defaultVerboseLogging
	if os.Getenv("CLI_VERBOSE_LOGGING") == "true" {
		verboseLogging = true
	}
}

func teardown() {
	// global test teardown instructions
}

//-------------------------------------------------------------------------------
// Tests
//-------------------------------------------------------------------------------
func TestCliValidateConfig(t *testing.T) {
	const (
		allFailedMsg = `Config validation failed:
* invalid password, failed (atLeast8) rule
* invalid password, failed (number) rule
* invalid password, failed (upper) rule
* invalid password, failed (special) rule`
		allButUpperFailedMsg = `Config validation failed:
* invalid password, failed (atLeast8) rule
* invalid password, failed (number) rule
* invalid password, failed (special) rule`
		atLeast8AndNumberFailedMsg = `Config validation failed:
* invalid password, failed (atLeast8) rule
* invalid password, failed (number) rule`
		atLeast8FailedMsg = `Config validation failed:
* invalid password, failed (atLeast8) rule`
	)
	var tests = []validateStringType{
		{"x", "validate-config", allFailedMsg},
		{"xX", "validate-config", allButUpperFailedMsg},
		{"xX%", "validate-config", atLeast8AndNumberFailedMsg},
		{"xX%9", "validate-config", atLeast8FailedMsg},
		{"xX%9999999", "validate-config", "validate-config success!"},
	}
	for _, test := range tests {
		got := test.validateString()
		if got != test.want {
			t.Errorf("%s %s => \nGOT\n%s\nWANTED\n%v", cliBinaryName, test.cmdAndFlags, test.want, got)
		}
	}
}

func TestCliVerston(t *testing.T) {
	const version = "%s version 0.0.1 (revision deadbeef)"
	var test = validateStringType{"xX%9999999", "--version", fmt.Sprintf(version, filepath.Base(cliBinaryName))}
	got := test.validateString()
	if got != test.want {
		t.Errorf("%s %s => \nGOT\n%s\nWANTED\n%v", cliBinaryName, test.cmdAndFlags, got, test.want)
	}
}

func TestCliSmokeTest(t *testing.T) {
	var test = validateIntType{"xX%9999999", "smoke-test", 10}
	got := test.validateInt()
	if got != test.want {
		t.Errorf("%s %s => \nGOT\n%d\nWANTED\n%d", cliBinaryName, test.cmdAndFlags, got, test.want)
	}
}

//-------------------------------------------------------------------------------
// Helpers
//-------------------------------------------------------------------------------
func countOccurences(got string) (occurrencesFound int) {
	findRegex := regexp.MustCompile("Did we get it\\? true")
	matches := findRegex.FindAllStringIndex(got, -1)
	return len(matches)
}
