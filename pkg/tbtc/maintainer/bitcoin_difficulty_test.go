package maintainer

import (
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

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
		btcChain: btcChain,
		chain:    nil,
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

	chain := &localBitcoinDifficultyChain{}

	chain.SetCurrentEpoch(299)
	chain.SetProofLength(3)

	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain: btcChain,
		chain:    chain,
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
