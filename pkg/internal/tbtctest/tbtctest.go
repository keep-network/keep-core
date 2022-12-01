// Package tbtctest contains scenarios meant to be used for Bitcoin-related
// tests in the pkg/tbtc package. Here are the details of specific scenarios:
//
// - deposit_sweep_scenario_1.json: Bitcoin deposit sweep transaction in which
//   three inputs (a P2WPKH main UTXO, a P2SH deposit, and a P2WSH deposit)
//   were swept into a P2WPKH main UTXO. For reference see:
//   https://live.blockcypher.com/btc-testnet/tx/435d4aff6d4bc34134877bd3213c17970142fdd04d4113d534120033b9eecb2e
package tbtctest

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"io/fs"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	testDataDirFormat              = "%s/testdata"
	depositSweepTestDataFilePrefix = "deposit_sweep_scenario"
)

// Deposit holds the deposit data in the given test scenario.
type Deposit struct {
	Utxo                *bitcoin.UnspentTransactionOutput
	Depositor           [20]byte
	BlindingFactor      [8]byte
	WalletPublicKeyHash [20]byte
	RefundPublicKeyHash [20]byte
	RefundLocktime      [4]byte
	Vault               [20]byte
}

// DepositSweepTestScenario represents a deposit sweep test scenario.
type DepositSweepTestScenario struct {
	Title             string
	WalletPublicKey   *ecdsa.PublicKey
	WalletPrivateKey  *big.Int
	WalletMainUtxo    *bitcoin.UnspentTransactionOutput
	Deposits          []*Deposit
	InputTransactions []*bitcoin.Transaction
	Fee               int64
	Signatures        []*bitcoin.SignatureContainer

	ExpectedSigHashes                   []*big.Int
	ExpectedSweepTransaction            *bitcoin.Transaction
	ExpectedSweepTransactionHash        bitcoin.Hash
	ExpectedSweepTransactionWitnessHash bitcoin.Hash
}

// LoadDepositSweepTestScenarios loads all scenarios related with deposit sweep.
func LoadDepositSweepTestScenarios() ([]*DepositSweepTestScenario, error) {
	filePaths, err := detectTestDataFiles(depositSweepTestDataFilePrefix)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot detect test data files: [%v]",
			err,
		)
	}

	scenarios := make([]*DepositSweepTestScenario, 0)

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

		var scenario DepositSweepTestScenario
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
