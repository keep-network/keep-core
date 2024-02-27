package tbtcpg_test

import (
	"testing"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"github.com/keep-network/keep-core/pkg/tbtcpg"
)

func TestFindMovingFundsTxData(t *testing.T) {
	walletPublicKeyHash := hexToByte20(
		"92a6ec889a8fa34f731e639edede4c75e184307c",
	)

	currentBlock := uint64(1000000)

	tbtcChain := tbtcpg.NewLocalChain()

	blockCounter := tbtcpg.NewMockBlockCounter()
	blockCounter.SetCurrentBlock(currentBlock)
	tbtcChain.SetBlockCounter(blockCounter)

	startBlock := currentBlock - tbtcpg.MovedFundsSweepLookBackBlocks

	// Test data on all the moving funds commitments.
	movingFundsCommitments := []struct {
		sourceWallet  [20]byte
		targetWallets [][20]byte
	}{
		// Moving funds 1: transfers funds to the test wallet. The proof has
		//                 not been submitted yet.
		{
			sourceWallet: hexToByte20(
				"1e37ea706ec1300d518ac6e4f58d1f4390443566",
			),
			targetWallets: [][20]byte{
				// unrelated wallet at index 0
				hexToByte20("c7302d75072d78be94eb8d36c4b77583c7abb06e"),
				// test wallet at index 1
				hexToByte20("92a6ec889a8fa34f731e639edede4c75e184307c"),
			},
		},

		// Moving funds 2: transfers funds to the test wallet. The proof has
		//                 been submitted but the moved funds sweep request is
		//                 already processed.
		{
			sourceWallet: hexToByte20(
				"d4162e0c6522eafac43650c78928910ceda4aef7",
			),
			targetWallets: [][20]byte{
				// test wallet at index 0
				hexToByte20("92a6ec889a8fa34f731e639edede4c75e184307c"),
			},
		},

		// Moving funds 3: transfers funds to the test wallet. The proof has
		//                 been submitted and the moved funds sweep request is
		//                 pending.
		{
			sourceWallet: hexToByte20(
				"ffb3f7538bfa98a511495dd96027cfbd57baf2fa",
			),
			targetWallets: [][20]byte{
				// unrelated wallet at index 0
				hexToByte20("3091d288521caec06ea912eacfd733edc5a36d6e"),
				// test wallet at index 1
				hexToByte20("92a6ec889a8fa34f731e639edede4c75e184307c"),
				// unrelated wallet at index 2
				hexToByte20("c7302d75072d78be94eb8d36c4b77583c7abb06e"),
			},
		},

		// Moving funds 4: transfers funds to unrelated wallet.
		{
			sourceWallet: hexToByte20(
				"a37c0a3e56decce984e4e6ccda7394c09b507fb5",
			),
			targetWallets: [][20]byte{
				// unrelated wallet at index 0
				hexToByte20("3091d288521caec06ea912eacfd733edc5a36d6e"),
			},
		},
	}

	for _, movingFundsTx := range movingFundsCommitments {
		err := tbtcChain.AddPastMovingFundsCommitmentSubmittedEvent(
			&tbtc.MovingFundsCommitmentSubmittedEventFilter{
				StartBlock: startBlock,
			},
			&tbtc.MovingFundsCommitmentSubmittedEvent{
				WalletPublicKeyHash: movingFundsTx.sourceWallet,
				TargetWallets:       movingFundsTx.targetWallets,
			},
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	// The filter for moving funds completed events contains the starting block
	// and the public key hashes of all the source wallets that committed to
	// move funds to the test wallet.
	movingFundsCompletedFilter := &tbtc.MovingFundsCompletedEventFilter{
		StartBlock: startBlock,
		WalletPublicKeyHash: [][20]byte{
			hexToByte20("1e37ea706ec1300d518ac6e4f58d1f4390443566"),
			hexToByte20("d4162e0c6522eafac43650c78928910ceda4aef7"),
			hexToByte20("ffb3f7538bfa98a511495dd96027cfbd57baf2fa"),
		},
	}

	// Tests data on all the completed moving funds transactions. Note it does
	// not contain data from moving funds 1 as the transaction is not completed.
	movingFundsTransactions := []struct {
		sourceWallet [20]byte
		txHash       bitcoin.Hash
	}{
		// Data from Moving funds 2.
		{
			sourceWallet: hexToByte20(
				"d4162e0c6522eafac43650c78928910ceda4aef7",
			),
			txHash: hashFromString(
				"722b27db8c84795a0eff3103394d0d34469dc79b822981eeb2ae0c023d7ce1e0",
			),
		},

		// Data from Moving funds 3.
		{
			sourceWallet: hexToByte20(
				"ffb3f7538bfa98a511495dd96027cfbd57baf2fa",
			),
			txHash: hashFromString(
				"9a745afae8ad4e3a50f64da9ad23e6d14f9bd5e13ae6266d2099e6d16388161c",
			),
		},
	}

	for _, movingFundsTx := range movingFundsTransactions {
		err := tbtcChain.AddPastMovingFundsCompletedEvent(
			movingFundsCompletedFilter,
			&tbtc.MovingFundsCompletedEvent{
				WalletPublicKeyHash: movingFundsTx.sourceWallet,
				MovingFundsTxHash:   movingFundsTx.txHash,
			},
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Data to build moved funds sweep requests. Only the
	movedFundsSweepRequests := []struct {
		txHash      bitcoin.Hash
		outputIndex uint32
		state       tbtc.MovedFundsSweepRequestState
	}{
		// Data from Moving funds 2.
		{
			txHash: hashFromString(
				"722b27db8c84795a0eff3103394d0d34469dc79b822981eeb2ae0c023d7ce1e0",
			),
			outputIndex: 0,
			state:       tbtc.MovedFundsStateProcessed,
		},

		// Data from Moving funds 3.
		{
			txHash: hashFromString(
				"9a745afae8ad4e3a50f64da9ad23e6d14f9bd5e13ae6266d2099e6d16388161c",
			),
			outputIndex: 1,
			state:       tbtc.MovedFundsStatePending,
		},
	}

	for _, request := range movedFundsSweepRequests {
		tbtcChain.SetMovedFundsSweepRequest(
			request.txHash,
			request.outputIndex,
			&tbtc.MovedFundsSweepRequest{
				State: request.state,
			},
		)
	}

	task := tbtcpg.NewMovedFundsSweepTask(
		tbtcChain,
		nil,
	)

	movingFundsTxHash, movingFundsOutputIdx, err := task.FindMovingFundsTxData(
		walletPublicKeyHash,
	)
	if err != nil {
		t.Fatal(err)
	}

	// The expected result comes from the Moving funds 3.
	expectedOutputIdx := uint32(1)
	expectedTxHash := hashFromString(
		"9a745afae8ad4e3a50f64da9ad23e6d14f9bd5e13ae6266d2099e6d16388161c",
	)

	testutils.AssertBytesEqual(
		t,
		expectedTxHash[:],
		movingFundsTxHash[:],
	)

	testutils.AssertUintsEqual(
		t,
		"fee",
		uint64(expectedOutputIdx),
		uint64(movingFundsOutputIdx),
	)
}

func TestEstimateMovedFundsSweepFee(t *testing.T) {
	var tests = map[string]struct {
		sweepTxMaxTotalFee uint64
		hasMainUtxo        bool
		expectedFee        uint64
		expectedError      error
	}{
		"estimated fee correct, one input": {
			sweepTxMaxTotalFee: 3000,
			hasMainUtxo:        false,
			expectedFee:        1760,
			expectedError:      nil,
		},
		"estimated fee correct, two inputs": {
			sweepTxMaxTotalFee: 3000,
			hasMainUtxo:        true,
			expectedFee:        2848,
			expectedError:      nil,
		},
		"estimated fee too high": {
			sweepTxMaxTotalFee: 2500,
			hasMainUtxo:        true,
			expectedFee:        0,
			expectedError:      tbtcpg.ErrSweepTxFeeTooHigh,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			btcChain := tbtcpg.NewLocalBitcoinChain()
			btcChain.SetEstimateSatPerVByteFee(1, 16)

			actualFee, err := tbtcpg.EstimateMovedFundsSweepFee(
				btcChain,
				test.hasMainUtxo,
				test.sweepTxMaxTotalFee,
			)

			testutils.AssertUintsEqual(
				t,
				"fee",
				test.expectedFee,
				uint64(actualFee),
			)

			testutils.AssertAnyErrorInChainMatchesTarget(
				t,
				test.expectedError,
				err,
			)
		})
	}
}

func hashFromString(str string) bitcoin.Hash {
	hash, err := bitcoin.NewHashFromString(
		str,
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		panic(err)
	}

	return hash
}
