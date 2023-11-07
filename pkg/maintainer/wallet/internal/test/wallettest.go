package test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	walletmtr "github.com/keep-network/keep-core/pkg/maintainer/wallet"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

const (
	testDataDirFormat                        = "%s/testdata"
	findDepositsToSweepTestDataFilePrefix    = "find_deposits"
	proposeDepositsSweepTestDataFilePrefix   = "propose_sweep"
	findPendingRedemptionsTestDataFilePrefix = "find_pending_redemptions"
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
	FundingTx              *bitcoin.Transaction

	WalletPublicKeyHash [20]byte

	RevealBlockNumber uint64
	SweptAt           time.Time
}

// FindDepositsToSweepTestScenario represents a test scenario of finding deposits to sweep.
type FindDepositsToSweepTestScenario struct {
	Title string

	MaxNumberOfDeposits uint16
	WalletPublicKeyHash [20]byte

	Wallets  []*Wallet
	Deposits []*Deposit

	ExpectedWalletPublicKeyHash [20]byte
	ExpectedUnsweptDeposits     []*walletmtr.DepositReference

	SweepTxFee             int64
	EstimateSatPerVByteFee int64
}

// LoadFindDepositsToSweepTestScenario loads all scenarios related with deposit sweep.
func LoadFindDepositsToSweepTestScenario() ([]*FindDepositsToSweepTestScenario, error) {
	return loadTestScenarios[*FindDepositsToSweepTestScenario](findDepositsToSweepTestDataFilePrefix)
}

type ProposeSweepDepositsData struct {
	walletmtr.DepositReference

	Transaction            *bitcoin.Transaction
	FundingTxConfirmations uint
}

// ProposeSweepTestScenario represents a test scenario of proposing deposits sweep.
type ProposeSweepTestScenario struct {
	Title string

	WalletPublicKeyHash          [20]byte
	DepositTxMaxFee              uint64
	Deposits                     []*ProposeSweepDepositsData
	SweepTxFee                   int64
	EstimateSatPerVByteFee       int64
	ExpectedDepositSweepProposal *tbtc.DepositSweepProposal
	ExpectedErr                  error
}

func (psts *ProposeSweepTestScenario) DepositsReferences() []*walletmtr.DepositReference {
	result := make([]*walletmtr.DepositReference, len(psts.Deposits))
	for i, d := range psts.Deposits {
		result[i] = &walletmtr.DepositReference{
			FundingTxHash:      d.FundingTxHash,
			FundingOutputIndex: d.FundingOutputIndex,
			RevealBlock:        d.RevealBlock,
		}
	}

	return result
}

// LoadProposeSweepTestScenario loads all scenarios related with deposit sweep.
func LoadProposeSweepTestScenario() ([]*ProposeSweepTestScenario, error) {
	return loadTestScenarios[*ProposeSweepTestScenario](proposeDepositsSweepTestDataFilePrefix)
}

// RedemptionRequest holds the redemption request data in the given test scenario.
type RedemptionRequest struct {
	WalletPublicKeyHash  [20]byte
	RedeemerOutputScript bitcoin.Script
	RequestedAmount      uint64
	RequestedAt          time.Time
	RequestBlock         uint64
}

// FindPendingRedemptionsTestScenario represents a test scenario of finding
// pending redemptions.
type FindPendingRedemptionsTestScenario struct {
	Title           string
	ChainParameters struct {
		AverageBlockTime time.Duration
		CurrentBlock     uint64
		RequestTimeout   uint32
		RequestMinAge    uint32
	}
	Filter                            walletmtr.PendingRedemptionsFilter
	Wallets                           []*Wallet
	PendingRedemptions                []*RedemptionRequest
	ExpectedWalletsPendingRedemptions map[[20]byte][]bitcoin.Script
}

// LoadFindPendingRedemptionsTestScenario loads all scenarios related with
// finding pending redemptions.
func LoadFindPendingRedemptionsTestScenario() (
	[]*FindPendingRedemptionsTestScenario,
	error,
) {
	return loadTestScenarios[*FindPendingRedemptionsTestScenario](
		findPendingRedemptionsTestDataFilePrefix,
	)
}

func loadTestScenarios[T json.Unmarshaler](testDataFilePrefix string) ([]T, error) {
	filePaths, err := detectTestDataFiles(testDataFilePrefix)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot detect test data files: [%v]",
			err,
		)
	}

	scenarios := make([]T, 0)

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

		var scenario T
		if err = json.Unmarshal(fileBytes, &scenario); err != nil {
			return nil, fmt.Errorf(
				"cannot unmarshal scenario for file [%v]: [%v]",
				filePath,
				err,
			)
		}

		scenarios = append(scenarios, scenario)
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
