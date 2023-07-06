package spv

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/internal/testutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

func TestGetProofInfo(t *testing.T) {
	tests := map[string]struct {
		latestBlockHeight               uint
		transactionConfirmations        uint
		currentEpoch                    uint64
		currentEpochDifficulty          *big.Int
		previousEpochDifficulty         *big.Int
		expectedIsProofWithinRelayRange bool
		expectedRequiredConfirmations   uint
	}{
		"proof entirely within current epoch": {
			latestBlockHeight:               790277,
			transactionConfirmations:        3,
			currentEpoch:                    392,
			currentEpochDifficulty:          nil, // not needed
			previousEpochDifficulty:         nil, // not needed
			expectedIsProofWithinRelayRange: true,
			expectedRequiredConfirmations:   6,
		},
		"proof entirely within previous epoch": {
			latestBlockHeight:               790300,
			transactionConfirmations:        2041,
			currentEpoch:                    392,
			currentEpochDifficulty:          nil, // not needed
			previousEpochDifficulty:         nil, // not needed
			expectedIsProofWithinRelayRange: true,
			expectedRequiredConfirmations:   6,
		},
		"proof spans previous and current epochs and difficulty drops": {
			latestBlockHeight:               790300,
			transactionConfirmations:        31,
			currentEpoch:                    392,
			currentEpochDifficulty:          big.NewInt(50000000000000),
			previousEpochDifficulty:         big.NewInt(30000000000000),
			expectedIsProofWithinRelayRange: true,
			expectedRequiredConfirmations:   9,
		},
		"proof spans previous and current epochs and difficulty raises": {
			latestBlockHeight:               790300,
			transactionConfirmations:        31,
			currentEpoch:                    392,
			currentEpochDifficulty:          big.NewInt(30000000000000),
			previousEpochDifficulty:         big.NewInt(60000000000000),
			expectedIsProofWithinRelayRange: true,
			expectedRequiredConfirmations:   4,
		},
		"proof begins outside previous epoch": {
			latestBlockHeight:               790300,
			transactionConfirmations:        2048,
			currentEpoch:                    392,
			currentEpochDifficulty:          nil, // not needed
			previousEpochDifficulty:         nil, // not needed
			expectedIsProofWithinRelayRange: false,
			expectedRequiredConfirmations:   0,
		},
		"proof ends outside current epoch": {
			latestBlockHeight:               792285,
			transactionConfirmations:        3,
			currentEpoch:                    392,
			currentEpochDifficulty:          nil, // not needed
			previousEpochDifficulty:         nil, // not needed
			expectedIsProofWithinRelayRange: false,
			expectedRequiredConfirmations:   0,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			transactionHash, err := bitcoin.NewHashFromString(
				"44c568bc0eac07a2a9c2b46829be5b5d46e7d00e17bfb613f506a75ccf86a473",
				bitcoin.InternalByteOrder,
			)
			if err != nil {
				t.Fatal(err)
			}

			localChain := newLocalChain()

			btcChain := newLocalBitcoinChain()
			btcChain.addBlockHeader(
				test.latestBlockHeight,
				&bitcoin.BlockHeader{},
			)
			btcChain.addTransactionConfirmations(
				transactionHash,
				test.transactionConfirmations,
			)

			localChain.setTxProofDifficultyFactor(big.NewInt(6))
			localChain.setCurrentEpoch(test.currentEpoch)
			localChain.setCurrentAndPrevEpochDifficulty(
				test.currentEpochDifficulty,
				test.previousEpochDifficulty,
			)

			isProofWithinRelayRange, requiredConfirmations, err := getProofInfo(
				transactionHash,
				btcChain,
				localChain,
				localChain,
			)
			if err != nil {
				t.Fatal(err)
			}

			testutils.AssertBoolsEqual(
				t,
				"is proof within range",
				test.expectedIsProofWithinRelayRange,
				isProofWithinRelayRange,
			)

			testutils.AssertUintsEqual(
				t,
				"required confirmations",
				uint64(test.expectedRequiredConfirmations),
				uint64(requiredConfirmations),
			)
		})
	}
}
