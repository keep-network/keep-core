package relay

import (
	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/entry"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/registry"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-relay")

const maxGroupSize = 255

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

// MonitorRelayEntry is listetning to the chain for a new relay entry.
// When a processing group which is supposed to deliver a relay entry does not
// fulfill its work, then this Node notifies the chain about it. In the case of
// delivering a relay entry by a processing group, this Node does nothing.
func (n *Node) MonitorRelayEntry(
	relayChain relayChain.Interface,
	relayRequestBlockNumber uint64,
	chainConfig *config.Chain,
) {
	logger.Infof("monitoring chain for a new relay entry")

	timeoutWaiterChannel, err := n.blockCounter.BlockHeightWaiter(relayRequestBlockNumber + chainConfig.RelayEntryTimeout)
	if err != nil {
		logger.Errorf("waiter for a relay entry timeout block failed: [%v]", err)
	}

	onEntrySubmittedChannel := make(chan *event.EntrySubmitted)

	subscription, err := relayChain.OnRelayEntrySubmitted(
		func(event *event.EntrySubmitted) {
			onEntrySubmittedChannel <- event
		},
	)
	if err != nil {
		close(onEntrySubmittedChannel)
		logger.Errorf("could not watch for a signature submission: [%v]", err)
		return
	}

	for {
		select {
		case blockNumber := <-timeoutWaiterChannel:
			subscription.Unsubscribe()
			close(onEntrySubmittedChannel)
			logger.Warnf(
				"relay entry was not submitted on time, reporting timeout at block [%v]",
				blockNumber,
			)
			err = relayChain.ReportRelayEntryTimeout()
			if err != nil {
				logger.Errorf("could not report a relay entry timeout: [%v]", err)
			}
			return
		case entry := <-onEntrySubmittedChannel:
			logger.Infof(
				"relay entry was submitted by the selected group on time at block [%v]",
				entry.BlockNumber,
			)
			return
		}
	}
}

// GenerateRelayEntry is triggered for a new relay request and checks if this
// client is one of the group members selected to create a new relay entry.
// If it is, this client enters the threshold signature creation process and,
// upon successfully completing it, submits the signature as a new relay entry.
// Note that this function returns immediately after determining whether the
// node is or is not a member of the requested group, and signature creation
// and submission is performed in a background goroutine.
func (n *Node) GenerateRelayEntry(
	previousEntry []byte,
	relayChain relayChain.Interface,
	signing chain.Signing,
	groupPublicKey []byte,
	startBlockHeight uint64,
) {
	memberships := n.groupRegistry.GetGroup(groupPublicKey)

	if len(memberships) < 1 {
		return
	}

	channel, err := n.netProvider.BroadcastChannelFor(memberships[0].ChannelName)
	if err != nil {
		logger.Errorf("could not create broadcast channel: [%v]", err)
		return
	}

	entry.RegisterUnmarshallers(channel)

	groupMembers, err := relayChain.GetGroupMembers(groupPublicKey)
	if err != nil {
		logger.Errorf("could not get group members: [%v]", err)
		return
	}

	membershipValidator := group.NewStakersMembershipValidator(
		groupMembers,
		signing,
	)

	err = channel.SetFilter(membershipValidator.IsInGroup)
	if err != nil {
		logger.Errorf(
			"could not set filter for channel [%v]: [%v]",
			channel.Name(),
			err,
		)
	}

	for _, member := range memberships {
		go func(member *registry.Membership) {
			err = entry.SignAndSubmit(
				n.blockCounter,
				channel,
				relayChain,
				previousEntry,
				n.chainConfig.HonestThreshold,
				member.Signer,
				startBlockHeight,
			)
			if err != nil {
				logger.Errorf(
					"error creating threshold signature: [%v]",
					err,
				)
				return
			}
		}(member)
	}
}
