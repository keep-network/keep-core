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
	"github.com/keep-network/keep-core/pkg/tbtc"
)

const (
	testDataDirFormat                      = "%s/testdata"
	findDepositsToSweepTestDataFilePrefix  = "find_deposits"
	proposeDepositsSweepTestDataFilePrefix = "propose_sweep"
)

// Wallet holds the wallet data in the given test scenario.
type Wallet struct {
	WalletPublicKeyHash     [20]byte
	RegistrationBlockNumber uint64
}

// Deposit holds the deposit data in the given test scenario.
type Deposit struct {
	FundingTxHash          bitcoin.Hash
	FundingOutputIndex     uint32
	FundingTxConfirmations uint

	WalletPublicKeyHash [20]byte

	RevealBlockNumber uint64
	SweptAt           time.Time
}

// FindDepositsToSweepTestScenario represents a deposit sweep test scenario.
type FindDepositsToSweepTestScenario struct {
	Title string

	MaxNumberOfDeposits uint16
	WalletPublicKeyHash [20]byte

	Wallets  []*Wallet
	Deposits []*Deposit

	ExpectedWalletPublicKeyHash [20]byte
	ExpectedUnsweptDeposits     []*coordinator.DepositSweepDetails

	SweepTxFee             int64
	EstimateSatPerVByteFee int64
}

// LoadFindDepositsToSweepTestScenario loads all scenarios related with deposit sweep.
func LoadFindDepositsToSweepTestScenario() ([]*FindDepositsToSweepTestScenario, error) {
	filePaths, err := detectTestDataFiles(findDepositsToSweepTestDataFilePrefix)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot detect test data files: [%v]",
			err,
		)
	}

	scenarios := make([]*FindDepositsToSweepTestScenario, 0)

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

		var scenario FindDepositsToSweepTestScenario
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

type ProposeSweepDepositsData struct {
	coordinator.DepositSweepDetails
	Transaction            *bitcoin.Transaction
	FundingTxConfirmations uint
}

// ProposeSweepTestScenario represents a deposit sweep test scenario.
type ProposeSweepTestScenario struct {
	Title string

	WalletPublicKeyHash          [20]byte
	DepositTxMaxFee              uint64
	Deposits                     []*ProposeSweepDepositsData
	SweepTxFee                   int64
	EstimateSatPerVByteFee       int64
	ExpectedDepositSweepProposal *tbtc.DepositSweepProposal
}

func (p *ProposeSweepTestScenario) DepositsSweepDetails() []*coordinator.DepositSweepDetails {
	result := make([]*coordinator.DepositSweepDetails, len(p.Deposits))
	for i, d := range p.Deposits {
		result[i] = &coordinator.DepositSweepDetails{
			FundingTxHash:      d.FundingTxHash,
			FundingOutputIndex: d.FundingOutputIndex,
			RevealBlock:        d.RevealBlock,
		}
	}

	return result
}

// LoadProposeSweepTestScenario loads all scenarios related with deposit sweep.
func LoadProposeSweepTestScenario() ([]*ProposeSweepTestScenario, error) {
	filePaths, err := detectTestDataFiles(proposeDepositsSweepTestDataFilePrefix)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot detect test data files: [%v]",
			err,
		)
	}

	scenarios := make([]*ProposeSweepTestScenario, 0)

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

		var scenario ProposeSweepTestScenario
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
