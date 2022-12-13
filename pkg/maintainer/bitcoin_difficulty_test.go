package maintainer

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestVerifySubmissionEligibility(t *testing.T) {
	tests := map[string]struct {
		ready                 bool
		authorizationRequired bool
		operatorAuthorized    bool
		expectedError         error
	}{
		"chain not ready": {
			ready:                 false,
			authorizationRequired: false,
			operatorAuthorized:    false,
			expectedError:         errNoGenesis,
		},
		"authorization not required": {
			ready:                 true,
			authorizationRequired: false,
			operatorAuthorized:    false,
			expectedError:         nil,
		},
		"operator not authorized": {
			ready:                 true,
			authorizationRequired: true,
			operatorAuthorized:    false,
			expectedError:         errNotAuthorized,
		},
		"operator authorized": {
			ready:                 true,
			authorizationRequired: true,
			operatorAuthorized:    true,
			expectedError:         nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			difficultyChain := connectLocalBitcoinDifficultyChain()
			operatorAddress := difficultyChain.Signing().Address()

			difficultyChain.SetReady(test.ready)
			difficultyChain.SetAuthorizationRequired(test.authorizationRequired)
			difficultyChain.SetAuthorizedOperator(
				operatorAddress,
				test.operatorAuthorized,
			)

			bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
				btcChain:           nil,
				chain:              difficultyChain,
				idleBackOffTime:    defaultIdleBackOffTime,
				restartBackOffTime: defaultRestartBackoffTime,
			}

			err := bitcoinDifficultyMaintainer.verifySubmissionEligibility()
			testutils.AssertAnyErrorInChainMatchesTarget(
				t,
				test.expectedError,
				err,
			)
		})
	}
}

func TestProveNextEpoch(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	btcChain := connectLocalBitcoinChain()

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

	difficultyChain := connectLocalBitcoinDifficultyChain()

	difficultyChain.SetCurrentEpoch(299)
	difficultyChain.SetProofLength(3)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:           btcChain,
		chain:              difficultyChain,
		idleBackOffTime:    defaultIdleBackOffTime,
		restartBackOffTime: defaultRestartBackoffTime,
	}

	result, err := bitcoinDifficultyMaintainer.proveNextEpoch(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult := true
	if result != expectedResult {
		t.Fatalf(
			"unexpected result returned\nexpected: %v\nactual:   %v\n",
			expectedResult,
			result,
		)
	}

	expectedNumberOfRetargetEvents := 1
	retargetEvents := difficultyChain.RetargetEvents()
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

	// Call once more, this time without any new epoch to prove
	result, err = bitcoinDifficultyMaintainer.proveNextEpoch(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expectedResult = false
	if result != expectedResult {
		t.Fatalf(
			"unexpected result returned\nexpected: %v\nactual:   %v\n",
			expectedResult,
			result,
		)
	}
}

func TestGetBlockHeaders(t *testing.T) {
	btcChain := connectLocalBitcoinChain()

	blockHeaders := map[uint]*bitcoin.BlockHeader{
		700000: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000000,
			Bits:                    1111111,
			Nonce:                   30,
		},
		700001: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000100,
			Bits:                    1111111,
			Nonce:                   40,
		},
		700002: {
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000200,
			Bits:                    2222222,
			Nonce:                   50,
		},
	}
	btcChain.SetBlockHeaders(blockHeaders)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:           btcChain,
		chain:              nil,
		idleBackOffTime:    defaultIdleBackOffTime,
		restartBackOffTime: defaultRestartBackoffTime,
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

func TestWaitForCurrentEpochUpdate_Successful(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	currentEpoch := uint64(299)
	targetEpoch := uint64(300)

	difficultyChain := connectLocalBitcoinDifficultyChain()
	difficultyChain.SetCurrentEpoch(currentEpoch)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:           nil,
		chain:              difficultyChain,
		idleBackOffTime:    2 * time.Second,
		restartBackOffTime: 2 * time.Second,
	}

	// Run function on a goroutine. The function should wait until the current
	// epoch is set to the target epoch.
	errChan := make(chan error, 1)
	go func() {
		err := bitcoinDifficultyMaintainer.waitForCurrentEpochUpdate(
			ctx,
			targetEpoch,
		)

		errChan <- err
	}()

	// Make sure the function keeps waiting for the current epoch to be updated.
	select {
	case <-time.After(1500 * time.Millisecond):
	case <-errChan:
		t.Fatal("Unexpected return from function")
	}

	// Update the current epoch to allow the waiting function to return.
	difficultyChain.SetCurrentEpoch(targetEpoch)
	select {
	case <-time.After(1500 * time.Millisecond):
		t.Fatal("Function did not return on time")
	case err := <-errChan:
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestWaitForCurrentEpochUpdate_Cancelled(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	currentEpoch := uint64(299)
	targetEpoch := uint64(300)

	difficultyChain := connectLocalBitcoinDifficultyChain()
	difficultyChain.SetCurrentEpoch(currentEpoch)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:           nil,
		chain:              difficultyChain,
		idleBackOffTime:    2 * time.Second,
		restartBackOffTime: 2 * time.Second,
	}

	// Run function on a goroutine. The function should wait until the current
	// epoch is set to the target epoch.
	errChan := make(chan error, 1)
	go func() {
		err := bitcoinDifficultyMaintainer.waitForCurrentEpochUpdate(
			ctx,
			targetEpoch,
		)

		errChan <- err
	}()

	// Cancel context while the function is waiting for current epoch to be
	// updated.
	cancelCtx()

	var err error
	select {
	case <-time.After(1500 * time.Millisecond):
		t.Fatal("Function did not return on time")
	case err = <-errChan:
	}

	testutils.AssertAnyErrorInChainMatchesTarget(t, context.Canceled, err)
}

func TestProveEpochs_ErrorVerifyingSubmissionEligibility(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// Do not authorize the maintainer to trigger an error.
	difficultyChain := connectLocalBitcoinDifficultyChain()
	difficultyChain.SetReady(true)
	difficultyChain.SetAuthorizationRequired(true)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:           nil,
		chain:              difficultyChain,
		idleBackOffTime:    defaultIdleBackOffTime,
		restartBackOffTime: defaultRestartBackoffTime,
	}

	err := bitcoinDifficultyMaintainer.proveEpochs(ctx)
	testutils.AssertAnyErrorInChainMatchesTarget(t, errNotAuthorized, err)
}

func TestProveEpochs_ErrorProvingSingleEpoch(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	difficultyChain := connectLocalBitcoinDifficultyChain()
	maintainerAddress := difficultyChain.Signing().Address()

	difficultyChain.SetReady(true)
	difficultyChain.SetAuthorizationRequired(true)
	difficultyChain.SetAuthorizedOperator(
		maintainerAddress,
		true,
	)

	// Do not set block headers in the Bitcoin chain to trigger an error.
	btcChain := connectLocalBitcoinChain()

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:           btcChain,
		chain:              difficultyChain,
		idleBackOffTime:    defaultIdleBackOffTime,
		restartBackOffTime: defaultRestartBackoffTime,
	}

	err := bitcoinDifficultyMaintainer.proveEpochs(ctx)
	testutils.AssertAnyErrorInChainMatchesTarget(t, errNoBlocksSet, err)
}

func TestProveEpochs_Successful(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	difficultyChain := connectLocalBitcoinDifficultyChain()
	maintainerAddress := difficultyChain.Signing().Address()

	difficultyChain.SetReady(true)
	difficultyChain.SetAuthorizationRequired(true)
	difficultyChain.SetAuthorizedOperator(
		maintainerAddress,
		true,
	)
	difficultyChain.SetProofLength(1)
	difficultyChain.SetCurrentEpoch(299)

	btcChain := connectLocalBitcoinChain()

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
		btcChain:           btcChain,
		chain:              difficultyChain,
		idleBackOffTime:    2 * time.Second,
		restartBackOffTime: 2 * time.Second,
	}

	// Run a goroutine that will cancel the context while the maintainer is
	// proving epochs. Maintainer should prove a single epoch and quit.
	go func() {
		time.Sleep(time.Second)
		cancelCtx()
	}()

	err := bitcoinDifficultyMaintainer.proveEpochs(ctx)
	testutils.AssertAnyErrorInChainMatchesTarget(t, context.Canceled, err)
}

func TestBitcoinDifficultyMaintainer_Integration(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	difficultyChain := connectLocalBitcoinDifficultyChain()
	maintainerAddress := difficultyChain.Signing().Address()

	difficultyChain.SetReady(true)
	difficultyChain.SetAuthorizationRequired(true)
	difficultyChain.SetAuthorizedOperator(
		maintainerAddress,
		true,
	)
	difficultyChain.SetProofLength(1)
	difficultyChain.SetCurrentEpoch(299)

	btcChain := connectLocalBitcoinChain()

	idleBackOffTime := 500 * time.Millisecond
	restartBackOffTime := 1 * time.Second

	initializeBitcoinDifficultyMaintainer(
		ctx,
		btcChain,
		difficultyChain,
		idleBackOffTime,
		restartBackOffTime,
	)

	//************ Loop restart on error ************
	// Do not set any headers in the Bitcoin chain, so that an error is
	// triggered. Wait for a moment to make sure the Bitcoin difficulty
	// maintainer started processing headers.
	time.Sleep(100 * time.Millisecond)

	//************ Prove two epochs ************
	// Set block headers for epochs 300 and 301 in the Bitcoin chain.
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
			Bits:                    2222222,
			Nonce:                   50,
		},
		606816: { // First block of the epoch 301
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000500,
			Bits:                    3333333,
			Nonce:                   60,
		},
	}
	btcChain.SetBlockHeaders(blockHeaders)

	// Wait for the Bitcoin difficulty maintainer to try processing headers
	// again after the previous attempt that ended in an error.
	time.Sleep(restartBackOffTime)

	// Make sure the first new epoch has been proven.
	expectedNumberOfRetargetEvents := 2
	retargetEvents := difficultyChain.RetargetEvents()
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

	// Set block headers for epoch 302 in the Bitcoin chain.
	blockHeaders = map[uint]*bitcoin.BlockHeader{
		608831: { // Last block of the epoch 301
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000600,
			Bits:                    3333333,
			Nonce:                   70,
		},
		608832: { // First block of the epoch 302
			Version:                 0,
			PreviousBlockHeaderHash: bitcoin.Hash{},
			MerkleRootHash:          bitcoin.Hash{},
			Time:                    1000700,
			Bits:                    4444444,
			Nonce:                   80,
		},
	}
	btcChain.SetBlockHeaders(blockHeaders)

	// Wait before proceeding with testing. If the Bitcoin difficulty maintainer
	// has not stopped, it will prove another epoch.
	time.Sleep(restartBackOffTime)

	// Make sure the Bitcoin difficulty maintainer has stopped and the number
	// of proven epochs has not changed.
	expectedNumberOfRetargetEvents = 2
	retargetEvents = difficultyChain.RetargetEvents()
	if len(retargetEvents) != expectedNumberOfRetargetEvents {
		t.Fatalf(
			"unexpected number of retarget events\nexpected: %v\nactual:   %v\n",
			expectedNumberOfRetargetEvents,
			len(retargetEvents),
		)
	}
}
