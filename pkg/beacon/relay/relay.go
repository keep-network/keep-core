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
	combinedEntryToSign := combineEntryToSign(
		request.PreviousEntry(),
		request.Seed.Bytes(),
	)

	membership, err := n.memberAndGroupForRequest(request)
	if err != nil {
		return err
	}

	go func() {
		signature, err := thresholdsignature.Execute(
			combinedEntryToSign,
			n.blockCounter,
			membership.channel,
			membership.member,
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

func combineEntryToSign(previousEntry []byte, seed []byte) []byte {
	combinedEntryToSign := make([]byte, 0)
	combinedEntryToSign = append(combinedEntryToSign, previousEntry...)
	combinedEntryToSign = append(combinedEntryToSign, seed...)
	return combinedEntryToSign
}

func (n *Node) indexForNextGroup(request event.Request) *big.Int {
	var (
		entry     *big.Int
		nextGroup *big.Int
	)
	entry = entry.SetBytes(request.PreviousEntry())
	numberOfGroups := big.NewInt(int64(len(n.groupPublicKeys)))
	_ = big.NewInt(0).SetBytes([]byte{})

	if numberOfGroups.Int64() == 0 {
		return nextGroup
	}

	nextGroup.Mod(entry, numberOfGroups)
	return nextGroup
}

func (n *Node) memberAndGroupForRequest(
	request event.Request,
) (*membership, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	nextGroup := n.indexForNextGroup(request).Int64()
	// Search our list of memberships to see if we have a member entry.
	for _, membership := range n.myGroups {
		if membership.index == int(nextGroup) {
			return membership, nil
		}
	}
	// We don't have an entry - maybe in the pending groups?
	for _, membership := range n.pendingGroups {
		if membership.index == int(nextGroup) {
			return membership, nil
		}
	}

	return nil, fmt.Errorf(
		"nonexistant membership for requested group [%d]",
		nextGroup,
	)
}
