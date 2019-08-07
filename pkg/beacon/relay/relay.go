package relay

import (
	"math/big"

	"github.com/ipfs/go-log"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/entry"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-relay")

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
		groupRegistry: groupRegistry,
	}
}

// MonitorRelayEntryOnChain is listetning to the chain for a new relay entry.
// When a processing group which is supposed to deliver a relay entry does not
// fulfill its work, then this Node notifies the chain about it. In the case of
// delivering a relay entry by a processing group, this Node does nothing.
func (n *Node) MonitorRelayEntryOnChain(
	relayChain relayChain.Interface,
	relayRequestBlockNumber uint64,
	chainConfig *config.Chain,
) {
	logger.Infof("observing chain for a new relay entry")

	timeoutWaiterChannel, err := n.blockCounter.BlockHeightWaiter(relayRequestBlockNumber + chainConfig.RelayEntryTimeout)
	if err != nil {
		logger.Errorf("block height waiter failure [%v]", err)
	}

	onSubmittedResultChannel := make(chan *event.Entry)

	subscription, err := relayChain.OnSignatureSubmitted(
		func(event *event.Entry) {
			onSubmittedResultChannel <- event
		},
	)
	if err != nil {
		close(onSubmittedResultChannel)
		logger.Errorf("could not watch for a signature submission: [%v]", err)
		return
	}

	for {
		select {
		case <-timeoutWaiterChannel:
			subscription.Unsubscribe()
			close(onSubmittedResultChannel)
			err = relayChain.ReportRelayEntryTimeout()
			if err != nil {
				logger.Errorf("could not report a relay entry timeout: [%v]", err)
			}
			return
		case <-onSubmittedResultChannel:
			return
		}
	}
}

// GenerateRelayEntry takes a relay request and checks if this client
// is one of the nodes elected by that request to create a new relay entry.
// If it is, this client enters the threshold signature creation process and,
// upon successfully completing it, submits the signature as a new relay entry
// to the passed in relayChain. Note that this function returns immediately after
// determining whether the node is or is not is a member of the requested group, and
// signature creation and submission is performed in a background goroutine.
func (n *Node) GenerateRelayEntry(
	previousEntry *big.Int,
	seed *big.Int,
	relayChain relayChain.Interface,
	groupPublicKey []byte,
	startBlockHeight uint64,
) {
	memberships := n.groupRegistry.GetGroup(groupPublicKey)

	if len(memberships) < 1 {
		return
	}

	for _, signer := range memberships {
		go func(signer *registry.Membership) {
			channel, err := n.netProvider.ChannelFor(signer.ChannelName)
			if err != nil {
				logger.Errorf(
					"could not create broadcast channel with name [%v]: [%v]",
					signer.ChannelName,
					err,
				)
				return
			}

			err = entry.SignAndSubmit(
				n.blockCounter,
				channel,
				relayChain,
				previousEntry,
				seed,
				n.chainConfig.HonestThreshold(),
				signer.Signer,
				startBlockHeight,
			)
			if err != nil {
				logger.Errorf(
					"error creating threshold signature: [%v]",
					err,
				)
				return
			}
		}(signer)
	}
}
