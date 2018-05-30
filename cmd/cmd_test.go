package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const (
	defaultCLIBinaryName  = "../keep-client" // Can use CLI_BINARY env variable
	defaultVerboseLogging = false
	testFileMode          = os.FileMode(0640)
)

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
	var test = validateIntType{"xX%9999999", "smoke-test", defaultGroupSize}
	got := test.validateInt()
	if got != test.want {
		t.Errorf("%s %s => \nGOT\n%d\nWANTED\n%d", cliBinaryName, test.cmdAndFlags, got, test.want)
	}
}

func TestFileExists(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "foobar")

	if FileExists(file) {
		t.Errorf("File %q should not exist", file)
	}

	err = createEmptyFile(file)
	if err != nil {
		t.Fatal(err)
	}

	if !FileExists(file) {
		t.Errorf("File %q should not exist", file)
	}
}

func TestGetConfigFilePath(t *testing.T) {
	path, err := GetConfigFilePath("")
	if err != nil {
		t.Errorf("GetConfigFilePath(\"\") returned and error: %v", err)
	} else if path != DefaultConfigPath {
		t.Errorf("Calling GetConfigFilePath(\"\"), got %s, want: %s", path, DefaultConfigPath)
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

func createEmptyFile(path string) (err error) {
	return ioutil.WriteFile(path, []byte(""), testFileMode)
}
