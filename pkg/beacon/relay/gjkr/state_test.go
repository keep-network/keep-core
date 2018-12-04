package gjkr

import (
	"github.com/keep-network/keep-core/pkg/chain/local"
)

func initChain(
	threshold, groupSize, expectedProtocolDuration, blockStep int,
) (*Chain, error) {
	chainHandle := local.Connect(groupSize, threshold)
	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		return nil, err
	}

	// This is really wrong **** - should not have to wait for 1 block to get current block height.

	err = blockCounter.WaitForBlocks(1)
	if err != nil {
		return nil, err
	}

	initialBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, err
	}

	return &Chain{
		handle:                   chainHandle,
		expectedProtocolDuration: expectedProtocolDuration, // T_dkg
		blockStep:                blockStep,                // T_step
		initialBlockHeight:       initialBlockHeight,       // T_init
	}, nil
}
