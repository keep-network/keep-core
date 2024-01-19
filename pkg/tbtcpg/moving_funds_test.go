package tbtcpg_test

import (
	"encoding/hex"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"github.com/keep-network/keep-core/pkg/tbtcpg"
)

func TestMovingFundsAction_FindTargetWallets_CommitmentNotSubmittedYet(t *testing.T) {
	var tests = map[string]struct {
		sourceWalletPublicKeyHash [20]byte
		walletBalance             uint64
		walletMaxBtcTransfer      uint64
		registeredWallets         []walletInfo
		liveWalletsCount          uint32
		expectedTargetWallets     [][20]byte
		expectedError             error
	}{
		"success scenario": {
			sourceWalletPublicKeyHash: hexToByte20(
				"ffb3f7538bfa98a511495dd96027cfbd57baf2fa",
			),
			walletBalance:        2000000,
			walletMaxBtcTransfer: 300000,
			registeredWallets: []walletInfo{
				{
					publicKeyHash: hexToByte20(
						"92a6ec889a8fa34f731e639edede4c75e184307c",
					),
					state: tbtc.StateLive,
				},
				{
					publicKeyHash: hexToByte20(
						"fdfa28e238734271f5e0d4f53d3843ae6cc09b24",
					),
					state: tbtc.StateLive,
				},
				{
					publicKeyHash: hexToByte20(
						"840dac51a6346e9372efbdc5d3503ed9fd32abdf",
					),
					state: tbtc.StateMovingFunds,
				},
				{
					publicKeyHash: hexToByte20(
						"3091d288521caec06ea912eacfd733edc5a36d6e",
					),
					state: tbtc.StateLive,
				},
				{
					publicKeyHash: hexToByte20(
						"c7302d75072d78be94eb8d36c4b77583c7abb06e",
					),
					state: tbtc.StateLive,
				},
			},
			liveWalletsCount: 4,
			expectedTargetWallets: [][20]byte{
				// The target wallets list should include all the Live wallets
				// and the wallets should be sorted according to their numerical
				// representation.
				hexToByte20(
					"3091d288521caec06ea912eacfd733edc5a36d6e",
				),
				hexToByte20(
					"92a6ec889a8fa34f731e639edede4c75e184307c",
				),
				hexToByte20(
					"c7302d75072d78be94eb8d36c4b77583c7abb06e",
				),
				hexToByte20(
					"fdfa28e238734271f5e0d4f53d3843ae6cc09b24",
				),
			},
			expectedError: nil,
		},
		"wallet max BTC transfer is zero": {
			sourceWalletPublicKeyHash: hexToByte20(
				"ffb3f7538bfa98a511495dd96027cfbd57baf2fa",
			),
			walletBalance:         10000,
			walletMaxBtcTransfer:  0, // Set to zero.
			registeredWallets:     []walletInfo{},
			liveWalletsCount:      4,
			expectedTargetWallets: nil,
			expectedError:         tbtcpg.ErrMaxBtcTransferZero,
		},
		"not enough live wallets": {
			sourceWalletPublicKeyHash: hexToByte20(
				"ffb3f7538bfa98a511495dd96027cfbd57baf2fa",
			),
			walletBalance:        2000000,
			walletMaxBtcTransfer: 300000,
			// Simulate there should be two Live wallets, but set only one Live
			// wallet. This could only happen if one of the target wallets
			// changed its state between getting the live wallets count and
			// getting wallets chain data.
			registeredWallets: []walletInfo{
				{
					publicKeyHash: hexToByte20(
						"92a6ec889a8fa34f731e639edede4c75e184307c",
					),
					state: tbtc.StateClosing,
				},
				{
					publicKeyHash: hexToByte20(
						"fdfa28e238734271f5e0d4f53d3843ae6cc09b24",
					),
					state: tbtc.StateLive,
				},
				{
					publicKeyHash: hexToByte20(
						"840dac51a6346e9372efbdc5d3503ed9fd32abdf",
					),
					state: tbtc.StateMovingFunds,
				},
			},
			liveWalletsCount:      2,
			expectedTargetWallets: nil,
			expectedError:         tbtcpg.ErrNotEnoughTargetWallets,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			tbtcChain := tbtcpg.NewLocalChain()
			tbtcChain.SetWalletParameters(
				0,
				0,
				0,
				0,
				0,
				test.walletMaxBtcTransfer,
				0,
			)

			for _, walletInfo := range test.registeredWallets {
				tbtcChain.AddPastNewWalletRegisteredEvent(
					nil,
					&tbtc.NewWalletRegisteredEvent{
						WalletPublicKeyHash: walletInfo.publicKeyHash,
					},
				)
				tbtcChain.SetWallet(
					walletInfo.publicKeyHash,
					&tbtc.WalletChainData{State: walletInfo.state},
				)
			}

			task := tbtcpg.NewMovingFundsTask(tbtcChain, nil)

			// Always simulate the moving funds commitment has not been
			// submitted yet.
			targetWalletsCommitmentHash := [32]byte{}

			targetWallets, alreadySubmitted, err := task.FindTargetWallets(
				&testutils.MockLogger{},
				test.sourceWalletPublicKeyHash,
				targetWalletsCommitmentHash,
				test.walletBalance,
				test.liveWalletsCount,
			)

			if !reflect.DeepEqual(test.expectedTargetWallets, targetWallets) {
				t.Errorf(
					"unexpected target wallets\nexpected: %v\nactual:   %v",
					test.expectedTargetWallets,
					targetWallets,
				)
			}

			// Returned value for the already submitted commitment should
			// always be false.
			expectedAlreadySubmitted := false
			testutils.AssertBoolsEqual(
				t,
				"already submitted flag",
				expectedAlreadySubmitted,
				alreadySubmitted,
			)

			testutils.AssertAnyErrorInChainMatchesTarget(
				t,
				test.expectedError,
				err,
			)
		})
	}
}

func TestMovingFundsAction_FindTargetWallets_CommitmentAlreadySubmitted(t *testing.T) {
	var tests = map[string]struct {
		sourceWalletPublicKeyHash   [20]byte
		targetWalletsCommitmentHash [32]byte
		targetWallets               []walletInfo
		currentBlock                uint64
		averageBlockTime            time.Duration
		expectedTargetWallets       [][20]byte
		expectedError               error
	}{
		"success scenario": {
			sourceWalletPublicKeyHash: hexToByte20(
				"ffb3f7538bfa98a511495dd96027cfbd57baf2fa",
			),
			targetWalletsCommitmentHash: hexToByte32(
				"9d9368117956680760fa27bb9542ceba2d4fcc398d640a5a0769f5a9593afb0e",
			),
			targetWallets: []walletInfo{
				{
					publicKeyHash: hexToByte20(
						"92a6ec889a8fa34f731e639edede4c75e184307c",
					),
					state: tbtc.StateLive,
				},
				{
					publicKeyHash: hexToByte20(
						"c7302d75072d78be94eb8d36c4b77583c7abb06e",
					),
					state: tbtc.StateLive,
				},
				{
					publicKeyHash: hexToByte20(
						"fdfa28e238734271f5e0d4f53d3843ae6cc09b24",
					),
					state: tbtc.StateLive,
				},
			},
			currentBlock:     1000000,
			averageBlockTime: 10 * time.Second,
			expectedTargetWallets: [][20]byte{
				hexToByte20("92a6ec889a8fa34f731e639edede4c75e184307c"),
				hexToByte20("c7302d75072d78be94eb8d36c4b77583c7abb06e"),
				hexToByte20("fdfa28e238734271f5e0d4f53d3843ae6cc09b24"),
			},
			expectedError: nil,
		},
		"target wallet is not live": {
			sourceWalletPublicKeyHash: hexToByte20(
				"ffb3f7538bfa98a511495dd96027cfbd57baf2fa",
			),
			targetWalletsCommitmentHash: hexToByte32(
				"9d9368117956680760fa27bb9542ceba2d4fcc398d640a5a0769f5a9593afb0e",
			),
			targetWallets: []walletInfo{
				{
					publicKeyHash: hexToByte20(
						"92a6ec889a8fa34f731e639edede4c75e184307c",
					),
					state: tbtc.StateLive,
				},
				{
					publicKeyHash: hexToByte20(
						"c7302d75072d78be94eb8d36c4b77583c7abb06e",
					),
					state: tbtc.StateTerminated, // wrong state
				},
				{
					publicKeyHash: hexToByte20(
						"fdfa28e238734271f5e0d4f53d3843ae6cc09b24",
					),
					state: tbtc.StateLive,
				},
			},
			currentBlock:          1000000,
			averageBlockTime:      10 * time.Second,
			expectedTargetWallets: nil,
			expectedError:         tbtcpg.ErrTargetWalletNotLive,
		},
		"target wallet commitment hash mismatch": {
			sourceWalletPublicKeyHash: hexToByte20(
				"ffb3f7538bfa98a511495dd96027cfbd57baf2fa",
			),
			targetWalletsCommitmentHash: hexToByte32(
				"9d9368117956680760fa27bb9542ceba2d4fcc398d640a5a0769f5a9593afb0e",
			),
			targetWallets: []walletInfo{
				{ // Use only one target wallet.
					publicKeyHash: hexToByte20(
						"92a6ec889a8fa34f731e639edede4c75e184307c",
					),
					state: tbtc.StateLive,
				},
			},
			currentBlock:          1000000,
			averageBlockTime:      10 * time.Second,
			expectedTargetWallets: nil,
			expectedError:         tbtcpg.ErrWrongCommitmentHash,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			tbtcChain := tbtcpg.NewLocalChain()

			tbtcChain.SetAverageBlockTime(test.averageBlockTime)

			blockCounter := tbtcpg.NewMockBlockCounter()
			blockCounter.SetCurrentBlock(test.currentBlock)
			tbtcChain.SetBlockCounter(blockCounter)

			startBlock := test.currentBlock - uint64(
				tbtcpg.MovingFundsCommitmentLookBackPeriod.Seconds(),
			)/uint64(test.averageBlockTime.Seconds())

			targetWallets := [][20]byte{}
			for _, walletInfo := range test.targetWallets {
				targetWallets = append(targetWallets, walletInfo.publicKeyHash)
				tbtcChain.SetWallet(
					walletInfo.publicKeyHash,
					&tbtc.WalletChainData{State: walletInfo.state},
				)
			}

			err := tbtcChain.AddPastMovingFundsCommitmentSubmittedEvent(
				&tbtc.MovingFundsCommitmentSubmittedEventFilter{
					StartBlock: startBlock,
					WalletPublicKeyHash: [][20]byte{
						test.sourceWalletPublicKeyHash,
					},
				},
				&tbtc.MovingFundsCommitmentSubmittedEvent{
					TargetWallets: targetWallets,
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			task := tbtcpg.NewMovingFundsTask(tbtcChain, nil)

			// Live wallets count and wallet's balance don't matter, as we are
			// retrieving target wallets from an already submitted commitment.
			walletBalance := uint64(2000000)
			liveWalletsCount := uint32(5)

			targetWallets, alreadySubmitted, err := task.FindTargetWallets(
				&testutils.MockLogger{},
				test.sourceWalletPublicKeyHash,
				test.targetWalletsCommitmentHash,
				walletBalance,
				liveWalletsCount,
			)

			if !reflect.DeepEqual(test.expectedTargetWallets, targetWallets) {
				t.Errorf(
					"unexpected target wallets\nexpected: %v\nactual:   %v",
					test.expectedTargetWallets,
					targetWallets,
				)
			}

			// Returned value for the already submitted commitment should
			// always be true.
			expectedAlreadySubmitted := true
			testutils.AssertBoolsEqual(
				t,
				"already submitted flag",
				expectedAlreadySubmitted,
				alreadySubmitted,
			)

			testutils.AssertAnyErrorInChainMatchesTarget(
				t,
				test.expectedError,
				err,
			)
		})
	}
}

func TestMovingFundsAction_GetWalletMembersInfo(t *testing.T) {
	var tests = map[string]struct {
		walletOperators          []operatorInfo
		executingOperator        chain.Address
		expectedMemberIDs        []uint32
		expectedOperatorPosition uint32
		expectedError            error
	}{
		"success case": {
			walletOperators: []operatorInfo{
				{"5df232b0348928793658dd05dfc6b05a59d11ae8", 3},
				{"dcc895d32b74b34cef2baa6546884fcda65da1e9", 1},
				{"28759deda2ea33bd72f68ea2e8f60cd670c2549f", 2},
				{"f7891d42f3c61a49e0aed1e31b151877c0905cf7", 4},
			},
			executingOperator:        "28759deda2ea33bd72f68ea2e8f60cd670c2549f",
			expectedMemberIDs:        []uint32{3, 1, 2, 4},
			expectedOperatorPosition: 3,
			expectedError:            nil,
		},
		"executing operator not among operators": {
			walletOperators: []operatorInfo{
				{"5df232b0348928793658dd05dfc6b05a59d11ae8", 2},
				{"dcc895d32b74b34cef2baa6546884fcda65da1e9", 1},
			},
			executingOperator:        "28759deda2ea33bd72f68ea2e8f60cd670c2549f",
			expectedMemberIDs:        nil,
			expectedOperatorPosition: 0,
			expectedError:            tbtcpg.ErrNoExecutingOperator,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			tbtcChain := tbtcpg.NewLocalChain()

			task := tbtcpg.NewMovingFundsTask(tbtcChain, nil)

			walletOperators := []chain.Address{}
			for _, operatorInfo := range test.walletOperators {
				err := tbtcChain.SetOperatorID(
					operatorInfo.Address,
					operatorInfo.OperatorID,
				)
				if err != nil {
					t.Fatal(err)
				}
				walletOperators = append(walletOperators, operatorInfo.Address)
			}

			memberIDs, operatorPosition, err := task.GetWalletMembersInfo(
				walletOperators,
				test.executingOperator,
			)

			if !reflect.DeepEqual(test.expectedMemberIDs, memberIDs) {
				t.Errorf(
					"unexpected memberIDs\nexpected: %v\nactual:   %v",
					test.expectedMemberIDs,
					memberIDs,
				)
			}

			testutils.AssertUintsEqual(
				t,
				"operator position",
				uint64(test.expectedOperatorPosition),
				uint64(operatorPosition),
			)

			testutils.AssertAnyErrorInChainMatchesTarget(
				t,
				test.expectedError,
				err,
			)
		})
	}
}

type walletInfo struct {
	publicKeyHash [20]byte
	state         tbtc.WalletState
}

type operatorInfo struct {
	chain.Address
	chain.OperatorID
}

func hexToByte20(hexStr string) [20]byte {
	if len(hexStr) != 40 {
		panic("hex string length incorrect")
	}
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	var result [20]byte
	copy(result[:], decoded)
	return result
}

func hexToByte32(hexStr string) [32]byte {
	if len(hexStr) != 64 {
		panic("hex string length incorrect")
	}
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	var result [32]byte
	copy(result[:], decoded)
	return result
}
