// Package coordinatortest contains scenarios meant to be used for pkg/coordinator
// package tests.
package coordinatortest

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/coordinator"
)

const (
	testDataDirFormat               = "%s/testdata"
	depositsSweepTestDataFilePrefix = "deposits_sweep_scenario"
)

// Wallet holds the wallet data in the given test scenario.
type Wallet struct {
	WalletPublicKeyHash     coordinator.WalletPublicKeyHash
	RegistrationBlockNumber uint64
}

// Deposit holds the deposit data in the given test scenario.
type Deposit struct {
	FundingTxHash          bitcoin.Hash
	FundingOutputIndex     uint32
	FundingTxConfirmations uint

	WalletPublicKeyHash coordinator.WalletPublicKeyHash

	RevealBlockNumber uint64
	SweptAt           time.Time
}

// DepositsSweepTestScenario represents a deposit sweep test scenario.
type DepositsSweepTestScenario struct {
	Title string

	MaxNumberOfDeposits uint16
	WalletPublicKeyHash *coordinator.WalletPublicKeyHash

	Wallets  []*Wallet
	Deposits []*Deposit

	ExpectedWalletPublicKeyHash coordinator.WalletPublicKeyHash
	ExpectedUnsweptDeposits     []*coordinator.DepositSweepDetails
}

// LoadDepositsSweepTestScenarios loads all scenarios related with deposit sweep.
func LoadDepositsSweepTestScenarios() ([]*DepositsSweepTestScenario, error) {
	filePaths, err := detectTestDataFiles(depositsSweepTestDataFilePrefix)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot detect test data files: [%v]",
			err,
		)
	}

	scenarios := make([]*DepositsSweepTestScenario, 0)

	for _, filePath := range filePaths {
		// #nosec G304 (file path provided as taint input)
		// This line is used to read a test fixture file.
		// There is no user input.
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot read file [%v]: [%v]",
				filePath,
				err,
			)
		}

		var scenario DepositsSweepTestScenario
		if err = json.Unmarshal(fileBytes, &scenario); err != nil {
			return nil, fmt.Errorf(
				"cannot unmarshal scenario for file [%v]: [%v]",
				filePath,
				err,
			)
		}

		scenarios = append(scenarios, &scenario)
	}

	return scenarios, nil
}

func detectTestDataFiles(prefix string) ([]string, error) {
	_, callerFileName, _, _ := runtime.Caller(0)
	sourceDirName := filepath.Dir(callerFileName)
	testDataDirName := fmt.Sprintf(testDataDirFormat, sourceDirName)
	filePaths := make([]string, 0)

	err := filepath.Walk(
		testDataDirName,
		func(path string, info fs.FileInfo, err error) error {
			if strings.HasPrefix(info.Name(), prefix) {
				filePaths = append(filePaths, path)
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return filePaths, nil
}
