package cli

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/cli/runtime"
	"github.com/urfave/cli"
)

const (
	testFileMode        = os.FileMode(0640)
	dkgSignatureSuccess = "Did we get it\\? true"
	commandKey          = "COMMAND"
)

var (
	Version       = "1.0.0"
	Revision      = "591fbe1"
	cliBinaryName = "keep-core"
	osArgs        = []string{cliBinaryName, "test-commands"}
)

//-------------------------------------------------------------------------------
// Main
//-------------------------------------------------------------------------------
func TestMain(m *testing.M) {
	setup()

	code := m.Run()
	os.Exit(code)

	teardown()
}

// Global test setup instructions
func setup() {

	// Initialize BLS library
	err := bls.Init(bls.CurveSNARK1)
	if err != nil {
		log.Fatal("Failed to initialize BLS.", err)
	}

	//TODO: Add logging framework with verbose logging capability
	successMsg = ""
}

// Global test teardown instructions
func teardown() {
}

//-------------------------------------------------------------------------------
// Tests
//-------------------------------------------------------------------------------

// TestCliSmokeTest tests the smoketest command
func TestCliSmokeTest(t *testing.T) {

	smokeTestCmdName := "smoke-test"

	type testData struct {
		args map[string]string
		want int
	}

	data := []testData{
		{
			map[string]string{
				"--config": "/tmp/xxx",
				commandKey: "smoke-test",
				"-g":       strconv.Itoa(minGroupSize),
				"-t":       strconv.Itoa(minThreshold),
			},
			minGroupSize,
		},
	}

	for _, test := range data {

		err := RunCLI(osArgs, Version, Revision, runtime.Command(smokeTestCmdName))

		if err != nil {
			t.Errorf("(%s -g %s -t %s) command returned and error: %v", test.args[commandKey], test.args["-g"], test.args["-t"], err)
		} else {
			got := countOccurences(os.Getenv("SMOKETEST_OUT"))
			if test.want != got {
				t.Errorf("WANT %d, GOT: %d\n", test.want, got)
			}
		}

	}
}

// TestValidateConfig tests the validate-config command
func TestValidateConfig(t *testing.T) {

	app := cli.NewApp()

	type testData struct {
		arguments []string
		want      string
	}

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
	var tests = []testData{
		{[]string{"x"}, allFailedMsg},
		{[]string{"xX"}, allButUpperFailedMsg},
		{[]string{"xX%"}, atLeast8AndNumberFailedMsg},
		{[]string{"xX%9"}, atLeast8FailedMsg},
		{[]string{"xX%9999999"}, successMsg},
	}

	for _, test := range tests {
		set := flag.NewFlagSet("", 0)
		set.Parse(test.arguments)

		ctx := cli.NewContext(app, set, nil)
		os.Setenv("KEEP_ETHEREUM_PASSWORD", "xX%999999999")
		err := ValidateConfig(ctx)
		if err != nil {
			t.Errorf("WANT\n%s\nGOT\n%v", test.want, err.Error())
		}
	}
}

// TestFileExists test the FileExists function
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

// TestGetConfigFilePath test the GetConfigFilePath function
func TestGetConfigFilePath(t *testing.T) {
	path, err := GetConfigFilePath("")
	if err != nil {
		t.Errorf("GetConfigFilePath(\"\") returned and error: %v", err)
	} else if path != DefaultConfigPath {
		t.Errorf("Calling GetConfigFilePath(\"\"), got %s, want: %s\n", path, DefaultConfigPath)
	}
}

//-------------------------------------------------------------------------------
// Helpers
//-------------------------------------------------------------------------------
func countOccurences(got string) (occurrencesFound int) {
	findRegex := regexp.MustCompile(dkgSignatureSuccess)
	matches := findRegex.FindAllStringIndex(got, -1)
	return len(matches)
}

func createEmptyFile(path string) (err error) {
	return ioutil.WriteFile(path, []byte(""), testFileMode)
}
