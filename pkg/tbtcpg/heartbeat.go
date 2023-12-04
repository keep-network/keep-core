package tbtcpg

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/keep-network/keep-core/pkg/tbtc"
)

// HeartbeatTask is a task that may produce a heartbeat proposal.
type HeartbeatTask struct {
	chain Chain
}

func NewHeartbeatTask(chain Chain) *HeartbeatTask {
	return &HeartbeatTask{
		chain: chain,
	}
}

func (ht *HeartbeatTask) Run(request *tbtc.CoordinationProposalRequest) (
	tbtc.CoordinationProposal,
	bool,
	error,
) {
	walletPublicKeyHash := request.WalletPublicKeyHash

	blockCounter, err := ht.chain.BlockCounter()
	if err != nil {
		return nil, false, fmt.Errorf("failed to get block counter: [%v]", err)
	}

	block, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, false, fmt.Errorf("failed to get current block: [%v]", err)
	}
	blockBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(blockBytes, block)

	hash := sha256.Sum256(append(walletPublicKeyHash[:], blockBytes...))

	message := [16]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		hash[0], hash[1], hash[2], hash[3], hash[4], hash[5], hash[6], hash[7],
	}

	proposal := &tbtc.HeartbeatProposal{
		Message: message,
	}

	if err := ht.chain.ValidateHeartbeatProposal(proposal); err != nil {
		return nil, false, fmt.Errorf(
			"failed to verify heartbeat proposal: [%v]",
			err,
		)
	}

	return proposal, true, nil
}

func (ht *HeartbeatTask) ActionType() tbtc.WalletActionType {
	return tbtc.ActionHeartbeat
}
