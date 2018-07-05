package relay

import (
	"fmt"
	"os"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/thresholdsignature"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
)

// NewNode returns an empty Node with no group, zero group count, and a nil last
// seen entry, tied to the given net.Provider.
func NewNode(
	stakeID string,
	netProvider net.Provider,
	blockCounter chain.BlockCounter,
	chainConfig config.Chain,
) Node {
	return Node{
		StakeID:      stakeID,
		netProvider:  netProvider,
		blockCounter: blockCounter,
		chainConfig:  chainConfig,
	}
}

func (n *Node) GenerateRelayEntryIfEligible(
	req event.Request,
	relayChain relaychain.Interface,
) {
	combinedEntryToSign := make([]byte, 0)
	combinedEntryToSign = append(combinedEntryToSign, req.PreviousEntry()...)
	combinedEntryToSign = append(combinedEntryToSign, req.Seed.Bytes()...)

	thresholdMember, groupChannel, err := n.memberAndGroupForRequest(req)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error joining group channel for request group [%s]: [%v]",
			req.RequestID.String(),
			err,
		)
	}
	if thresholdMember != nil {
		go func() {
			signature, err := thresholdsignature.Execute(
				combinedEntryToSign,
				n.blockCounter,
				groupChannel,
				thresholdMember,
			)
			if err != nil {
				fmt.Printf(
					"error creating threshold signature: [%v]",
					err,
				)
				return
			}
		}()
	}
}

func (n *Node) memberAndGroupForRequest(
	req event.Request,
) (*thresholdgroup.Member, net.BroadcastChannel, error) {
	// Use request to choose group.
	// See if we're in the group.
	// If we are, look up and return our member entry and our broadcast channel entry.
	return nil, nil, nil
}
