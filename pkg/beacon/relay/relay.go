package relay

import (
	"fmt"
	"math/big"
	"os"
	"time"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry"
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
	groupRegistry *registry.Groups,
) Node {
	return Node{
		Staker:        staker,
		netProvider:   netProvider,
		blockCounter:  blockCounter,
		chainConfig:   chainConfig,
		stakeIDs:      make([]string, 100),
		groupRegistry: groupRegistry,
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
	previousEntry *big.Int,
	seed *big.Int,
	relayChain relaychain.RelayEntryInterface,
	groupPublicKey []byte,
	startBlockHeight uint64,
) {
	combinedEntryToSign := combineEntryToSign(
		previousEntry.Bytes(),
		seed.Bytes(),
	)

	memberships := n.groupRegistry.GetGroup(groupPublicKey)

	if len(memberships) < 1 {
		return
	}

	for _, signer := range memberships {
		go func(signer *registry.Membership) {
			channel, err := n.netProvider.ChannelFor(signer.ChannelName)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"could not create broadcast channel with name [%v]: [%v]\n",
					signer.ChannelName,
					err,
				)
				return
			}

			signature, err := thresholdsignature.Execute(
				requestID,
				combinedEntryToSign,
				n.chainConfig.HonestThreshold(),
				n.blockCounter,
				channel,
				signer.Signer,
				startBlockHeight,
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
				PreviousEntry: previousEntry,
				Timestamp:     time.Now().UTC(),
				GroupPubKey:   signer.Signer.GroupPublicKeyBytes(),
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
