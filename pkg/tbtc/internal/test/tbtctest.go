// Package tbtctest contains scenarios meant to be used for Bitcoin-related
// tests in the pkg/tbtc package. Here are the details of specific scenarios:
//
//   - deposit_sweep_scenario_0.json: Bitcoin deposit sweep transaction in which
//     three inputs (a P2WPKH main UTXO, a P2SH deposit, and a P2WSH deposit)
//     were swept into a P2WPKH main UTXO. For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/435d4aff6d4bc34134877bd3213c17970142fdd04d4113d534120033b9eecb2e
//
//   - deposit_sweep_scenario_1.json: Bitcoin deposit sweep transaction in which
//     two inputs (a P2SH deposit, and a P2WSH deposit) were swept into a P2WPKH
//     main UTXO. For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668
//
//   - deposit_sweep_scenario_2.json: Bitcoin deposit sweep transaction in which
//     one input (a P2SH deposit) was swept into a P2WPKH main UTXO.
//     For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/f5b9ad4e8cd5317925319ebc64dc923092bef3b56429c6b1bc2261bbdc73f351
//
//   - deposit_sweep_scenario_3.json: Bitcoin deposit sweep transaction in which
//     one input (a P2WSH deposit) was swept into a P2WPKH main UTXO.
//     For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/9efc9d555233e12e06378a35a7b988d54f7043b5c3156adc79c7af0a0fd6f1a0
//
//   - redemption_scenario_0.json: Bitcoin redemption transaction that uses a
//     single P2WPKH input to pay a single P2PKH redeemer script and a P2WPKH change.
//     For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/c437f1117db977682334b53a71fbe63a42aab42f6e0976c35b69977f86308c20
//
//   - redemption_scenario_1.json: Bitcoin redemption transaction that uses a
//     single P2WPKH input to pay a single P2WPKH redeemer script and a P2WPKH change.
//     For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/925e61dc31396e7f2cbcc8bc9b4009b4f24ba679257762df078b7e9b875ea110
//
//   - redemption_scenario_2.json: Bitcoin redemption transaction that uses a
//     single P2WPKH input to pay a single P2SH redeemer script and a P2WPKH change.
//     For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/ef25c9c8f4df673def035c0c1880278c90030b3c94a56668109001a591c2c521
//
//   - redemption_scenario_3.json: Bitcoin redemption transaction that uses a
//     single P2WPKH input to pay a single P2WSH redeemer script and a P2WPKH change.
//     For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/3d28bb5bf73379da51bc683f4d0ed31d7b024466c619d80ebd9378077d900be3
//
//   - redemption_scenario_4.json: Bitcoin redemption transaction that uses a
//     single P2WPKH input to pay redeemer scripts (P2PKH, P2WPKH, P2SH and P2WSH)
//     and a P2WPKH change.
//     For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/f70ff89fd2b6226183e4b8143cc5f0f457f05dd1dca0c6151ab66f4523d972b7
//
//   - redemption_scenario_5.json: Bitcoin redemption transaction that uses a
//     single P2WPKH input to pay redeemer scripts (P2PKH, P2WPKH) without a change.
//     For reference see:
//     https://live.blockcypher.com/btc-testnet/tx/afcdf8f91273b73abc40018873978c22bbb7c3d8d669ef2faffa0c4b0898c8eb
package test

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

const (
	testDataDirFormat              = "%s/testdata"
	depositSweepTestDataFilePrefix = "deposit_sweep_scenario"
	redemptionTestDataFilePrefix   = "redemption_scenario"
)

// Deposit holds the deposit data in the given test scenario.
type Deposit struct {
	Utxo                *bitcoin.UnspentTransactionOutput
	Depositor           chain.Address
	BlindingFactor      [8]byte
	WalletPublicKeyHash [20]byte
	RefundPublicKeyHash [20]byte
	RefundLocktime      [4]byte
	Vault               *chain.Address
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
	return loadTestScenarios[*DepositSweepTestScenario](depositSweepTestDataFilePrefix)
}

// RedemptionRequest holds the redemption request data in the given test scenario.
type RedemptionRequest struct {
	Redeemer             chain.Address
	RedeemerOutputScript []byte
	RequestedAmount      uint64
	TreasuryFee          uint64
	TxMaxFee             uint64
	RequestedAt          time.Time
}

// RedemptionTestScenario represents a redemption test scenario.
type RedemptionTestScenario struct {
	Title              string
	WalletPublicKey    *ecdsa.PublicKey
	WalletPrivateKey   *big.Int
	WalletMainUtxo     *bitcoin.UnspentTransactionOutput
	RedemptionRequests []*RedemptionRequest
	InputTransaction   *bitcoin.Transaction
	FeeShares          []int64
	Signature          *bitcoin.SignatureContainer

	ExpectedSigHash                          *big.Int
	ExpectedRedemptionTransaction            *bitcoin.Transaction
	ExpectedRedemptionTransactionHash        bitcoin.Hash
	ExpectedRedemptionTransactionWitnessHash bitcoin.Hash
}

// LoadRedemptionTestScenarios loads all scenarios related with redemption.
func LoadRedemptionTestScenarios() ([]*RedemptionTestScenario, error) {
	return loadTestScenarios[*RedemptionTestScenario](redemptionTestDataFilePrefix)
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
