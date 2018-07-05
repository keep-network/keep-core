package relay

import (
	"fmt"
	"math/big"
	"os"
	"time"

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

// GenerateRelayEntryIfEligible takes a relay request and checks if this client
// is one of the nodes elected by that request to create a new relay entry.
// If it is, this client enters the threshold signature creation process and,
// upon successfully completing it, submits the signature as a new relay entry
// to the passed in relayChain. Note that this function returns immediately after
// determining whether the node is or is not is a member of the requested group, and
// signature creation and submission is performed in a background goroutine.
func (n *Node) GenerateRelayEntryIfEligible(
	request event.Request,
	relayChain relaychain.Interface,
) error {
	combinedEntryToSign := make([]byte, 0)
	combinedEntryToSign = append(combinedEntryToSign, request.PreviousEntry()...)
	combinedEntryToSign = append(combinedEntryToSign, request.Seed.Bytes()...)

	thresholdMember, groupChannel, err := n.memberAndGroupForRequest(request)
	if err != nil {
		return err
	}

	go func() {
		signature, err := thresholdsignature.Execute(
			combinedEntryToSign,
			n.blockCounter,
			groupChannel,
			thresholdMember,
		)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"error creating threshold signature: [%v]",
				err,
			)
			return
		}

		var (
			rightSizeSignature [32]byte
			previousEntry      *big.Int
		)
		previousEntry.SetBytes(request.PreviousEntry())
		for i := 0; i < 32; i++ {
			rightSizeSignature[i] = signature[i]
		}

		newEntry := &event.Entry{
			RequestID:     request.RequestID,
			Value:         rightSizeSignature,
			PreviousEntry: previousEntry,
			Timestamp:     time.Now().UTC(),
		}

		relayChain.SubmitRelayEntry(
			newEntry,
		).OnFailure(func(err error) {
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Failed submission of relay entry: [%v].\n",
					err,
				)
				return
			}
		})
	}()

	return nil
}

func (n *Node) memberAndGroupForRequest(
	request event.Request,
) (*thresholdgroup.Member, net.BroadcastChannel, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	var (
		found         bool
		memberOfGroup *thresholdgroup.Member
	)
	// Search our list of memberships to see if we have a member entry.
	if memberOfGroup, found = n.myGroups[request.RequestID.String()]; !found {
		// We don't have an entry - maybe in the pending groups?
		pendingGroup, found := n.pendingGroups[request.RequestID.String()]
		if !found {
			return nil, nil, fmt.Errorf(
				"nonexistant membership for requested group [%s]",
				request.RequestID.String(),
			)
		}
		memberOfGroup = pendingGroup
	}

	if memberOfGroup == nil {
		return nil, nil, fmt.Errorf(
			"nonexistant membership for requested group [%s]",
			request.RequestID.String(),
		)
	}

	// We have a valid member entry, find the broadcast channel entry.
	groupChannel, err := n.netProvider.ChannelFor(request.RequestID.String())
	if err != nil {
		return nil, nil, fmt.Errorf(
			"Error joining group channel for request group [%s]: [%v]",
			request.RequestID.String(),
			err,
		)
	}
	return memberOfGroup, groupChannel, nil
}
