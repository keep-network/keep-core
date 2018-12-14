package publish

import (
	"math/big"
	"testing"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func TestPhase14_pt1(t *testing.T) {
	threshold := 2
	groupSize := 5
	blockStep := 2 // T_step

	chainHandle, initialBlock, err := initChainHandle2(threshold, groupSize)
	if err != nil {
		t.Fatal(err)
	}

	resultToPublish := &relayChain.DKGResult{
		GroupPublicKey: big.NewInt(12345),
	}
	_ = resultToPublish

	var tests = map[string]struct {
		correctResult   *relayChain.DKGResult
		publishingIndex int
		xyzzy           int
	}{
		"base call": {
			correctResult:   resultToPublish,
			publishingIndex: 0,
			xyzzy:           0,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			publisher := &Publisher{
				ID:              gjkr.MemberID(test.publishingIndex + 1),
				RequestID:       big.NewInt(101),
				publishingIndex: test.publishingIndex,
				chainHandle:     chainHandle,
				blockStep:       blockStep,
			}

			// Reinitialize chain to reset block counter
			publisher.chainHandle, initialBlock, err = initChainHandle2(threshold, groupSize)
			if err != nil {
				t.Fatalf("chain initialization failed [%v]", err)
			}

			chainRelay := publisher.chainHandle.ThresholdRelay()
			_ = chainRelay

			blockCounter, err := publisher.chainHandle.BlockCounter()
			if err != nil {
				t.Fatalf("unexpected error [%v]", err)
			}
			_ = blockCounter

			if true {
				// func (pm *Publisher) Phase14(correctResult *relayChain.DKGResult) error {
				publisher.Phase14(test.correctResult)
			}
		})
	}
}

func initChainHandle2(threshold, groupSize int) (chainHandle chain.Handle, initialBlock int, err error) {
	chainHandle = local.Connect(groupSize, threshold)
	blockCounter, err := chainHandle.BlockCounter() // PJS - save blockCounter?
	if err != nil {
		return nil, -1, err
	}
	err = blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, -1, err
	}

	initialBlock, err = blockCounter.CurrentBlock() // PJS - need CurrentBlock to make this work
	if err != nil {
		return nil, -1, err
	}
	return
}
