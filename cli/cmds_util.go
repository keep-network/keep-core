package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

type captureStdout struct {
	origStdout *os.File
	r          *os.File
	w          *os.File
	err        error
	outChan    chan string
	out        string
	printOut   bool
}

func (c *captureStdout) capture() {
	c.origStdout = os.Stdout
	c.r, c.w, c.err = os.Pipe()
	if c.err != nil {
		fmt.Printf("error while capturing stdout (pipe error): %v\n", c.err)
		return
	}
	os.Stdout = c.w
	c.outChan = make(chan string)
	// Copy the output in a separate goroutine s.t. printing won't block
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, c.r)
		if err != nil {
			fmt.Printf("error while capturing stdout (copy error): %v\n", c.err)
			return
		}
		c.outChan <- buf.String()
	}()
}

func (c *captureStdout) reset() {
	// back to normal state
	err := c.w.Close()
	if err != nil {
		fmt.Printf("Unable to close stdout: %v\n", err)
	}
	os.Stdout = c.origStdout // restoring the real stdout
	c.out = <-c.outChan
	// reading our temp stdout
	if c.printOut {
		fmt.Print(c.out)
	}
}

// GetConfigFilePath returns the full path to the project confiuration file
func GetConfigFilePath(configPath string) (string, error) {
	if configPath == "" {
		configPath = DefaultConfigPath
		displayConfigWarning()
	}
	if exists := FileExists(configPath); !exists {
		return "", fmt.Errorf("config file (%s) not found", configPath)
	}
	return configPath, nil
}

// FileExists returns true if a file at the given path exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func runningTest() bool {
	found := false
	i := 0
	_, callerPath, _, _ := runtime.Caller(0)
	for len(callerPath) > 0 && !found {
		i++
		_, callerPath, _, _ = runtime.Caller(i)
		if strings.HasSuffix(callerPath, "testing.go") {
			found = true
		}
	}
	return found
}

func initializeBls() {

	if !blsInitialized {
		// Initialize BLS library
		err := bls.Init(bls.CurveSNARK1)
		if err != nil {
			log.Fatal("Failed to initialize BLS.", err)
		} else {
			blsInitialized = true
		}
	}
}

// Only display this warning message once per test suite
func displayConfigWarning() {
	if !configWarningDisplayed {
		fmt.Printf("WARNING using default config path (%s)\n", DefaultConfigPath)
		configWarningDisplayed = true
	}
}
