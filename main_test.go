package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/urfave/cli"
)

const (
	configFilePath = "../../test/config.toml"
)

var (
	wd, _  = os.Getwd()
	stdOut = captureStdout{}
)

var flagvar string

func init() {
	flag.StringVar(&flagvar, "config", "/tmp/config.toml", "path to config file")
}

func TestFlags(t *testing.T) {
	cases := []struct {
		testArgs    []string
		expectedErr error
	}{
		{[]string{"get-info", "argOne", "argTwo", "--break"}, errors.New("flag provided but not defined: -break")},
		{[]string{"--config", configFilePath, "get-info"}, nil},
	}

	for _, c := range cases {

		app := NewApp(Version, Revision)
		app.Writer = ioutil.Discard
		set := flag.NewFlagSet("test", 0)
		set.StringVar(&configPath, "config", "", "")
		set.Parse(c.testArgs)

		context := cli.NewContext(app, set, nil)

		command := findCommandByName(app.Commands, "get-info")
		command.Action = func(_ *cli.Context) error { return nil }
		command.SkipFlagParsing = false
		command.SkipArgReorder = false
		command.UseShortOptionHandling = false

		expectEqual(t, command.Run(context), c.expectedErr)
	}
}

func TestGetInfo(t *testing.T) {
	cases := []struct {
		testArgs    []string
		expectedErr error
	}{
		{[]string{"--config", configFilePath, "get-info"}, nil},
		{[]string{"get-info"}, nil},
	}

	for _, c := range cases {

		app := NewApp(Version, Revision)
		app.Writer = ioutil.Discard
		set := flag.NewFlagSet("test", 0)
		set.StringVar(&configPath, "config", "", "")
		set.Parse(c.testArgs)

		context := cli.NewContext(app, set, nil)

		command := findCommandByName(app.Commands, "get-info")

		output := captureOutput(func() {
			command.Run(context)
		})

		text := `Keep client: %s

Description: %s
Version:     %s
Revision:    %s
Config Path: %s
`
		expectedOutput := fmt.Sprintf(text, app.Name, app.Description, Version, Revision, context.GlobalString("config"))
		expectEqual(t, expectedOutput, output)
	}
}

//-------------------------------------------------------------------------------
// Helpers
//-------------------------------------------------------------------------------

func findCommandByName(commands []cli.Command, name string) cli.Command {
	for _, cmd := range commands {
		if cmd.Name == name {
			return cmd
		}
	}
	return cli.Command{}
}

func expectEqual(t *testing.T, a interface{}, b interface{}) {
	_, fn, line, _ := runtime.Caller(1)
	fn = strings.Replace(fn, wd+"/", "", -1)

	if !reflect.DeepEqual(a, b) {
		t.Errorf("(%s:%d) WANT %v (type %v) - GOT %v (type %v)", fn, line, b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func captureOutput(f func()) string {
	stdOut.capture()
	f()
	stdOut.reset()
	return stdOut.out
}

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
	err := c.w.Close()
	if err != nil {
		fmt.Printf("Unable to close stdout: %v\n", err)
	}
	os.Stdout = c.origStdout // Restore the original stdout
	c.out = <-c.outChan
	// Read our temp stdout
	if c.printOut {
		fmt.Print(c.out)
	}
}
