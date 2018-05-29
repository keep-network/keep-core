package ethereum

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/ssh/terminal"
)

var client *provider

func TestEthereumBlockTest(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		wait         int
		want         time.Duration
		errorMessage string
	}{
		"does wait for a block": {
			wait:         1,
			want:         time.Duration(100000000),
			errorMessage: "Failed to wait for a single block",
		},
		"waited for a longer time": {
			wait:         2,
			want:         time.Duration(200000000),
			errorMessage: "Failed to wait for 2 blocks",
		},
		"doesn't wait if 0 blocks": {
			wait:         0,
			want:         time.Duration(1000),
			errorMessage: "Failed for a 0 block wait",
		},
		"invalid value": {
			wait:         -1,
			want:         time.Duration(0),
			errorMessage: "Waiting for a time when it should have errored",
		},
	}

	var e chan interface{}

	tim := 240 // Force test to fail if not completed in 4 minutes
	tick := time.NewTimer(time.Duration(tim) * time.Second)

	go func() {
		select {
		case e <- tick:
			t.Fatal("Test ran too long - it failed")
		}
	}()

	waitForBlock := BlockCounter(client)

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			start := time.Now().UTC()
			waitForBlock.WaitForBlocks(test.wait)
			end := time.Now().UTC()

			elapsed := end.Sub(start)
			if elapsed < test.want {
				t.Error(test.errorMessage)
			}
		})
	}
}

func TestMain(m *testing.M) {

	fn := "test/config.toml"

	Config, err := readConfig(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: Error reading configuration: %s\n", err)
		os.Exit(1)
	}

	client, err = Connect(Config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: Failed to connect to Ethereum: %s\n", err)
		os.Exit(1)
	}

	// os.Exit(m.Run())
}

// Environment variable with 'prompt' for prompting for password or the password.
const passwordEnvVariable = "KEEP_ETHEREUM_PASSWORD"

// Top level config structure from the config file specified on the command line.
type ethereum_config struct {
	Ethereum Config
}

// ReadConfig reads in the configuration file in .toml format.
func readConfig(filePath string) (cfg Config, err error) {

	var ec ethereum_config

	if _, err = toml.DecodeFile(filePath, &ec); err != nil {
		return ec.Ethereum, fmt.Errorf("unable to decode .toml file [%s] error [%s]", filePath, err)
	}

	var password string
	envPassword := os.Getenv(passwordEnvVariable)
	if envPassword == "prompt" {
		if password, err = readPassword("Enter Account Password: "); err != nil {
			return ec.Ethereum, err
		}
		ec.Ethereum.Account.KeyFilePassword = password
	} else {
		ec.Ethereum.Account.KeyFilePassword = envPassword
	}

	if ec.Ethereum.Account.KeyFilePassword == "" {
		return cfg, fmt.Errorf("Password is required.  Set " + passwordEnvVariable + " environment variable to password or 'prompt'")
	}

	return ec.Ethereum, nil
}

// ReadPassword prompts a user to enter a password.   The read password uses
// the system password reading call that helps to prevent key loggers from
// capturing the password.
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("Unable to read password, error [%s]", err)
	}
	return strings.TrimSpace(string(bytePassword)), nil
}
