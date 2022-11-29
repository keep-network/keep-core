package maintainer

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

func TestVerifySubmissionEligibility(t *testing.T) {
	tests := map[string]struct {
		ready                 bool
		authorizationRequired bool
		operatorAuthorized    bool
		expectedErr           error
	}{
		"chain not ready": {
			ready:                 false,
			authorizationRequired: false,
			operatorAuthorized:    false,
			expectedErr:           fmt.Errorf("relay genesis has not been performed"),
		},
		"authorization not required": {
			ready:                 true,
			authorizationRequired: false,
			operatorAuthorized:    false,
			expectedErr:           nil,
		},
		"operator not authorized": {
			ready:                 true,
			authorizationRequired: true,
			operatorAuthorized:    false,
			expectedErr: fmt.Errorf(
				"relay maintainer has not been authorized to submit block headers",
			),
		},
		"operator authorized": {
			ready:                 true,
			authorizationRequired: true,
			operatorAuthorized:    true,
			expectedErr:           nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			localChain := ConnectLocal()
			operatorAddress := localChain.Signing().Address()

			localChain.SetReady(test.ready)
			localChain.SetAuthorizationRequired(test.authorizationRequired)
			localChain.SetAuthorizedOperator(
				operatorAddress,
				test.operatorAuthorized,
			)

			bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
				btcChain:               nil,
				chain:                  localChain,
				epochProvenBackOffTime: defaultEpochProvenBackOffTime,
				restartBackOffTime:     defaultRestartBackoffTime,
			}

			err := bitcoinDifficultyMaintainer.verifySubmissionEligibility()
			if !reflect.DeepEqual(test.expectedErr, err) {
				t.Errorf(
					"unexpected error\nexpected: %v\nactual:   %v\n",
					test.expectedErr,
					err,
				)
			}
		})
	}
}

func TestProveSingleEpoch(t *testing.T) {
	btcChain := bitcoin.ConnectLocal()

	// Set three block headers on each side of the retarget. The old epoch
	// number is 299, the new epoch number is 300.
	blockHeaders := map[uint]*bitcoin.BlockHeader{
		604797: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000000,
			Bits:                    1111111,
			Nonce:                   10,
		},
		604798: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000100,
			Bits:                    1111111,
			Nonce:                   20,
		},
		604799: { // Last block of the old epoch (epoch 299)
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000200,
			Bits:                    1111111,
			Nonce:                   30,
		},
		604800: { // First block of the new epoch (epoch 300)
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000300,
			Bits:                    2222222,
			Nonce:                   40,
		},
		604801: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000400,
			Bits:                    2222222,
			Nonce:                   50,
		},
		604802: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000500,
			Bits:                    2222222,
			Nonce:                   60,
		},
	}
	btcChain.SetBlockHeaders(blockHeaders)

	chain := ConnectLocal()

	chain.SetCurrentEpoch(299)
	chain.SetProofLength(3)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:               btcChain,
		chain:                  chain,
		epochProvenBackOffTime: defaultEpochProvenBackOffTime,
		restartBackOffTime:     defaultRestartBackoffTime,
	}

	err := bitcoinDifficultyMaintainer.proveSingleEpoch()
	if err != nil {
		t.Fatal(err)
	}

	expectedNumberOfRetargetEvents := 1
	retargetEvents := chain.RetargetEvents()
	if len(retargetEvents) != expectedNumberOfRetargetEvents {
		t.Fatalf(
			"unexpected number of retarget events\nexpected: %v\nactual:   %v\n",
			expectedNumberOfRetargetEvents,
			len(retargetEvents),
		)
	}

	eventsOldDifficulty := retargetEvents[0].oldDifficulty
	expectedOldDifficulty := blockHeaders[604799].Bits
	if eventsOldDifficulty != expectedOldDifficulty {
		t.Fatalf(
			"unexpected old difficulty of the retarget event \n"+
				"expected: %v\nactual:   %v\n",
			expectedOldDifficulty,
			eventsOldDifficulty,
		)
	}

	eventsNewDifficulty := retargetEvents[0].newDifficulty
	expectedNewDifficulty := blockHeaders[604800].Bits
	if eventsNewDifficulty != expectedNewDifficulty {
		t.Fatalf(
			"unexpected new difficulty of the retarget event \n"+
				"expected: %v\nactual:   %v\n",
			expectedNewDifficulty,
			eventsNewDifficulty,
		)
	}
}

func TestGetBlockHeaders(t *testing.T) {
	btcChain := bitcoin.ConnectLocal()

	blockHeaders := map[uint]*bitcoin.BlockHeader{
		700000: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    11111,
			Bits:                    2222,
			Nonce:                   3333,
		},
		700001: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    222,
			Bits:                    333,
			Nonce:                   444,
		},
		700002: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    555,
			Bits:                    555,
			Nonce:                   666,
		},
	}
	btcChain.SetBlockHeaders(blockHeaders)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:               btcChain,
		chain:                  nil,
		epochProvenBackOffTime: defaultEpochProvenBackOffTime,
		restartBackOffTime:     defaultRestartBackoffTime,
	}

	headers, err := bitcoinDifficultyMaintainer.getBlockHeaders(700000, 700002)
	if err != nil {
		t.Fatal(err)
	}

	expectedHeaders := []*bitcoin.BlockHeader{
		blockHeaders[700000], blockHeaders[700001], blockHeaders[700002],
	}

	if !reflect.DeepEqual(expectedHeaders, headers) {
		t.Errorf("\nexpected: %v\nactual:   %v", expectedHeaders, headers)
	}
}

func TestProveEpochs_ErrorVerifyingSubmissionEligibility(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// Do not authorize the maintainer to trigger an error.
	chain := ConnectLocal()
	chain.SetReady(true)
	chain.SetAuthorizationRequired(true)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:               nil,
		chain:                  chain,
		epochProvenBackOffTime: defaultEpochProvenBackOffTime,
		restartBackOffTime:     defaultRestartBackoffTime,
	}

	err := bitcoinDifficultyMaintainer.proveEpochs(ctx)
	expectedError := fmt.Errorf(
		"cannot proceed with proving Bitcoin blockchain epochs in the relay " +
			"chain [relay maintainer has not been authorized to submit block " +
			"headers]",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestProveEpochs_ErrorProvingSingleEpoch(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	localChain := ConnectLocal()
	maintainerAddress := localChain.Signing().Address()

	localChain.SetReady(true)
	localChain.SetAuthorizationRequired(true)
	localChain.SetAuthorizedOperator(
		maintainerAddress,
		true,
	)

	// Do not set block headers in the Bitcoin chain to trigger an error.
	btcChain := bitcoin.ConnectLocal()

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:               btcChain,
		chain:                  localChain,
		epochProvenBackOffTime: defaultEpochProvenBackOffTime,
		restartBackOffTime:     defaultRestartBackoffTime,
	}

	err := bitcoinDifficultyMaintainer.proveEpochs(ctx)
	expectedError := fmt.Errorf(
		"cannot prove Bitcoin blockchain epoch to the relay chain [failed to " +
			"get current block number [blockchain does not contain any blocks]]",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestProveEpochs_Successful(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	localChain := ConnectLocal()
	maintainerAddress := localChain.Signing().Address()

	localChain.SetReady(true)
	localChain.SetAuthorizationRequired(true)
	localChain.SetAuthorizedOperator(
		maintainerAddress,
		true,
	)
	localChain.SetProofLength(1)
	localChain.SetCurrentEpoch(299)

	btcChain := bitcoin.ConnectLocal()

	// Set one block header on each side of the retarget. The old epoch number
	// is 299, the new epoch number is 300.
	blockHeaders := map[uint]*bitcoin.BlockHeader{
		604799: { // Last block of the old epoch (epoch 299)
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000200,
			Bits:                    1111111,
			Nonce:                   30,
		},
		604800: { // First block of the new epoch (epoch 300)
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000300,
			Bits:                    2222222,
			Nonce:                   40,
		},
	}
	btcChain.SetBlockHeaders(blockHeaders)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:               btcChain,
		chain:                  localChain,
		epochProvenBackOffTime: 2 * time.Second,
		restartBackOffTime:     2 * time.Second,
	}

	// Run a goroutine that will cancel the context while the maintainer is
	// proving epochs. Maintainer should prove a single epoch and quit.
	go func() {
		time.Sleep(time.Second)
		cancelCtx()
	}()
	err := bitcoinDifficultyMaintainer.proveEpochs(ctx)

	expectedError := fmt.Errorf("context canceled")
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestBitcoinDifficultyMaintainer_Integration(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	localChain := ConnectLocal()
	maintainerAddress := localChain.Signing().Address()

	localChain.SetReady(true)
	localChain.SetAuthorizationRequired(true)
	localChain.SetAuthorizedOperator(
		maintainerAddress,
		true,
	)
	localChain.SetProofLength(1)
	localChain.SetCurrentEpoch(299)

	btcChain := bitcoin.ConnectLocal()

	epochProvenBackOffTime := 500 * time.Millisecond
	restartBackOffTime := 1 * time.Second

	initializeBitcoinDifficultyMaintainer(
		ctx,
		btcChain,
		localChain,
		epochProvenBackOffTime,
		restartBackOffTime,
	)

	//************ Loop restart on error ************
	// Do not set any headers in the Bitcoin chain, so that an error is
	// triggered. Wait for a moment to make sure the Bitcoin difficulty
	// maintainer started processing headers.
	time.Sleep(100 * time.Millisecond)

	//************ Prove epoch for the first time ************
	// Set headers in the Bitcoin chain. The headers will be used to prove
	// epochs 300 and 301 in this and subsequent test.
	blockHeaders := map[uint]*bitcoin.BlockHeader{
		604799: { // Last block of the epoch 299
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000200,
			Bits:                    1111111,
			Nonce:                   30,
		},
		604800: { // First block of the epoch 300
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000300,
			Bits:                    2222222,
			Nonce:                   40,
		},
		606815: { // Last block of the epoch 300
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000400,
			Bits:                    3333333,
			Nonce:                   50,
		},
		606816: { // First block of the epoch 301
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000500,
			Bits:                    4444444,
			Nonce:                   60,
		},
		608831: { // Last block of the epoch 301
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000600,
			Bits:                    5555555,
			Nonce:                   70,
		},
		608832: { // First block of the epoch 302
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000700,
			Bits:                    6666666,
			Nonce:                   80,
		},
	}
	btcChain.SetBlockHeaders(blockHeaders)

	// Wait for the Bitcoin difficulty maintainer to try processing headers
	// again after the previous attempt that ended in an error.
	time.Sleep(restartBackOffTime)

	// Make sure the first new epoch has been proven.
	expectedNumberOfRetargetEvents := 1
	retargetEvents := localChain.RetargetEvents()
	if len(retargetEvents) != expectedNumberOfRetargetEvents {
		t.Fatalf(
			"unexpected number of retarget events\nexpected: %v\nactual:   %v\n",
			expectedNumberOfRetargetEvents,
			len(retargetEvents),
		)
	}

	eventsOldDifficulty := retargetEvents[0].oldDifficulty
	expectedOldDifficulty := blockHeaders[604799].Bits
	if eventsOldDifficulty != expectedOldDifficulty {
		t.Fatalf(
			"unexpected old difficulty of the retarget event \n"+
				"expected: %v\nactual:   %v\n",
			expectedOldDifficulty,
			eventsOldDifficulty,
		)
	}

	eventsNewDifficulty := retargetEvents[0].newDifficulty
	expectedNewDifficulty := blockHeaders[604800].Bits
	if eventsNewDifficulty != expectedNewDifficulty {
		t.Fatalf(
			"unexpected new difficulty of the retarget event \n"+
				"expected: %v\nactual:   %v\n",
			expectedNewDifficulty,
			eventsNewDifficulty,
		)
	}

	//************ Prove epoch for the second time ************
	// Wait for the Bitcoin difficulty maintainer to try processing headers
	// again after the previous successful attempt.
	time.Sleep(epochProvenBackOffTime)

	// Make sure the second new epoch has been proven.
	expectedNumberOfRetargetEvents = 2
	retargetEvents = localChain.RetargetEvents()
	if len(retargetEvents) != expectedNumberOfRetargetEvents {
		t.Fatalf(
			"unexpected number of retarget events\nexpected: %v\nactual:   %v\n",
			expectedNumberOfRetargetEvents,
			len(retargetEvents),
		)
	}

	eventsOldDifficulty = retargetEvents[1].oldDifficulty
	expectedOldDifficulty = blockHeaders[606815].Bits
	if eventsOldDifficulty != expectedOldDifficulty {
		t.Fatalf(
			"unexpected old difficulty of the retarget event \n"+
				"expected: %v\nactual:   %v\n",
			expectedOldDifficulty,
			eventsOldDifficulty,
		)
	}

	eventsNewDifficulty = retargetEvents[1].newDifficulty
	expectedNewDifficulty = blockHeaders[606816].Bits
	if eventsNewDifficulty != expectedNewDifficulty {
		t.Fatalf(
			"unexpected new difficulty of the retarget event \n"+
				"expected: %v\nactual:   %v\n",
			expectedNewDifficulty,
			eventsNewDifficulty,
		)
	}

	//************ Cancel context ************
	// Cancel the context to force the Bitcoin difficulty maintainer to stop.
	cancelCtx()

	// Wait before proceeding with testing. If the Bitcoin difficulty maintainer
	// has not stopped, it will prove another epoch.
	time.Sleep(epochProvenBackOffTime)

	// Make sure the Bitcoin difficulty maintainer has stopped and the number
	// of proven epochs has not changed.
	expectedNumberOfRetargetEvents = 2
	retargetEvents = localChain.RetargetEvents()
	if len(retargetEvents) != expectedNumberOfRetargetEvents {
		t.Fatalf(
			"unexpected number of retarget events\nexpected: %v\nactual:   %v\n",
			expectedNumberOfRetargetEvents,
			len(retargetEvents),
		)
	}
}
