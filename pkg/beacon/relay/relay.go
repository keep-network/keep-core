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
	staker chain.Staker,
	netProvider net.Provider,
	blockCounter chain.BlockCounter,
	chainConfig *config.Chain,
) Node {
	return Node{
		Staker:          staker,
		netProvider:     netProvider,
		blockCounter:    blockCounter,
		chainConfig:     chainConfig,
		stakeIDs:        make([]string, 100),
		groupPublicKeys: make([][]byte, 0),
		seenPublicKeys:  make(map[string]bool),
		myGroups:        make(map[string][]*membership),
		pendingGroups:   make(map[string][]*membership),
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
	requestID *big.Int,
	previousValue *big.Int,
	seed *big.Int,
	relayChain relaychain.RelayEntryInterface,
) {
	combinedEntryToSign := combineEntryToSign(
		previousValue.Bytes(),
		seed.Bytes(),
	)

	memberships := n.membershipsForRequest(previousValue)
	if len(memberships) < 1 {
		return
	}

	for _, signer := range memberships {
		go func(signer *membership) {
			signature, err := thresholdsignature.Execute(
				combinedEntryToSign,
				n.blockCounter,
				signer.channel,
				signer.member,
			)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error creating threshold signature: [%v]\n",
					err,
				)
				return
			}

			rightSizeSignature := big.NewInt(0).SetBytes(signature[:32])

			newEntry := &event.Entry{
				RequestID:     requestID,
				Value:         rightSizeSignature,
				PreviousEntry: previousValue,
				Timestamp:     time.Now().UTC(),
				GroupID:       &big.Int{},
				Seed:          seed,
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
		}(signer)
	}
}

func combineEntryToSign(previousEntry []byte, seed []byte) []byte {
	combinedEntryToSign := make([]byte, 0)
	combinedEntryToSign = append(combinedEntryToSign, previousEntry...)
	combinedEntryToSign = append(combinedEntryToSign, seed...)
	return combinedEntryToSign
}

func (n *Node) indexForNextGroup(previousValue *big.Int) *big.Int {
	numberOfGroups := big.NewInt(int64(len(n.groupPublicKeys)))

	return nextGroupIndex(previousValue, numberOfGroups)
}

func nextGroupIndex(entry *big.Int, numberOfGroups *big.Int) *big.Int {
	if numberOfGroups.Cmp(&big.Int{}) == 0 {
		return &big.Int{}
	}

	return (&big.Int{}).Mod(entry, numberOfGroups)
}

func (n *Node) membershipsForRequest(previousValue *big.Int) []*membership {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	nextGroup := n.indexForNextGroup(previousValue).Int64()
	// Search our list of memberships to see if we have a member entries.
	membershipsForRequest := make([]*membership, 0)
	for _, memberships := range n.myGroups {
		for _, membership := range memberships {
			if membership.index == int(nextGroup) {
				membershipsForRequest = append(
					membershipsForRequest, membership,
				)
			}
		}
	}

	return membershipsForRequest
}
