package beacon

import (
	"encoding/hex"
	"fmt"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/beacon/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/entry"
	"github.com/keep-network/keep-core/pkg/beacon/event"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"

	"github.com/keep-network/keep-core/pkg/beacon/registry"
	"github.com/keep-network/keep-core/pkg/net"
)

// node represents the current state of a beacon node.
type node struct {
	operatorPublicKey *operator.PublicKey

	// External interactors.
	netProvider net.Provider
	beaconChain beaconchain.Interface

	groupRegistry *registry.Groups
}

// newNode returns an empty node with no group, zero group count, and a nil last
// seen entry, tied to the given net.Provider.
func newNode(
	operatorPublicKey *operator.PublicKey,
	netProvider net.Provider,
	beaconChain beaconchain.Interface,
	groupRegistry *registry.Groups,
) *node {
	return &node{
		operatorPublicKey: operatorPublicKey,
		netProvider:       netProvider,
		beaconChain:       beaconChain,
		groupRegistry:     groupRegistry,
	}
}

// IsInGroup checks if this node is a member of the group which was selected to
// join a group which undergoes the process of generating a threshold relay entry.
func (n *node) IsInGroup(groupPublicKey []byte) bool {
	return len(n.groupRegistry.GetGroup(groupPublicKey)) > 0
}

// JoinDKGIfEligible takes a seed value and undergoes the process of the
// distributed key generation if this node's operator proves to be eligible for
// the group generated by that seed. This is an interactive on-chain process,
// and JoinDKGIfEligible can block for an extended period of time while it
// completes the on-chain operation.
func (n *node) JoinDKGIfEligible(
	dkgSeed *big.Int,
	dkgStartBlockNumber uint64,
) {
	logger.Infof(
		"checking eligibility for DKG with seed [0x%x]",
		dkgSeed,
	)

	selectedOperators, err := n.beaconChain.SelectGroup(dkgSeed)
	if err != nil {
		logger.Errorf(
			"failed to select group with seed [0x%x]: [%v]",
			dkgSeed,
			err,
		)
		return
	}

	if len(selectedOperators) > n.beaconChain.GetConfig().GroupSize {
		logger.Errorf(
			"group size larger than supported: [%v]",
			len(selectedOperators),
		)
		return
	}

	signing := n.beaconChain.Signing()

	operatorAddress, err := signing.PublicKeyToAddress(n.operatorPublicKey)
	if err != nil {
		logger.Errorf("failed to get operator address: [%v]", err)
		return
	}

	indexes := make([]uint8, 0)
	for index, selectedOperator := range selectedOperators {
		// See if we are amongst those chosen
		if selectedOperator == operatorAddress {
			indexes = append(indexes, uint8(index))
		}
	}

	// create temporary broadcast channel name for DKG using the
	// group selection seed
	channelName := dkgSeed.Text(16)

	if len(indexes) > 0 {
		logger.Infof(
			"joining DKG with seed [0x%x] and controlling [%v] group members",
			dkgSeed,
			len(indexes),
		)

		broadcastChannel, err := n.netProvider.BroadcastChannelFor(channelName)
		if err != nil {
			logger.Errorf("failed to get broadcast channel: [%v]", err)
			return
		}

		membershipValidator := group.NewOperatorsMembershipValidator(
			selectedOperators,
			signing,
		)

		err = broadcastChannel.SetFilter(membershipValidator.IsInGroup)
		if err != nil {
			logger.Errorf(
				"could not set filter for channel [%v]: [%v]",
				broadcastChannel.Name(),
				err,
			)
		}

		for _, index := range indexes {
			// Capture the player index for the goroutine.
			playerIndex := index

			go func() {
				signer, err := dkg.ExecuteDKG(
					dkgSeed,
					playerIndex,
					dkgStartBlockNumber,
					n.beaconChain,
					broadcastChannel,
					membershipValidator,
					selectedOperators,
				)
				if err != nil {
					logger.Errorf("failed to execute dkg: [%v]", err)
					return
				}

				groupPublicKey := hex.EncodeToString(
					signer.GroupPublicKeyBytesCompressed(),
				)

				// TODO: Consider snapshotting the key material just in case.
				err = n.groupRegistry.RegisterGroup(signer, groupPublicKey)
				if err != nil {
					logger.Errorf(
						"[member:%v] failed to register a group [%v]: [%v]",
						signer.MemberID(),
						groupPublicKey,
						err,
					)
					return
				}

				logger.Infof(
					"[member:%v] group [%v] registered successfully",
					signer.MemberID(),
					groupPublicKey,
				)
			}()
		}
	} else {
		logger.Infof("not eligible for DKG with seed [0x%x]", dkgSeed)
	}

	return
}

// ForwardSignatureShares enables the ability to forward signature shares
// messages to other nodes even if this node is not a part of the group which
// signs the relay entry.
func (n *node) ForwardSignatureShares(groupPublicKeyBytes []byte) {
	name, err := channelNameForPublicKeyBytes(groupPublicKeyBytes)
	if err != nil {
		logger.Warningf("could not forward signature shares: [%v]", err)
		return
	}

	n.netProvider.BroadcastChannelForwarderFor(name)
}

// ResumeSigningIfEligible enables a client to rejoin the ongoing signing process
// after it was crashed or restarted and if it belongs to the signing group.
func (n *node) ResumeSigningIfEligible() {
	isEntryInProgress, err := n.beaconChain.IsEntryInProgress()
	if err != nil {
		logger.Errorf(
			"failed checking if an entry is in progress: [%v]",
			err,
		)
		return
	}

	if isEntryInProgress {
		previousEntry, err := n.beaconChain.CurrentRequestPreviousEntry()
		if err != nil {
			logger.Errorf(
				"failed to get a previous entry for the current request: [%v]",
				err,
			)
			return
		}
		entryStartBlock, err := n.beaconChain.CurrentRequestStartBlock()
		if err != nil {
			logger.Errorf(
				"failed to get a start block for the current request: [%v]",
				err,
			)
			return
		}
		groupPublicKey, err := n.beaconChain.CurrentRequestGroupPublicKey()
		if err != nil {
			logger.Errorf(
				"failed to get a group public key for the current request: [%v]",
				err,
			)
			return
		}

		logger.Infof(
			"attempting to rejoin the current signing process [0x%x]",
			groupPublicKey,
		)
		n.GenerateRelayEntry(
			previousEntry,
			groupPublicKey,
			entryStartBlock.Uint64(),
		)
	}
}

// MonitorRelayEntry is listetning to the chain for a new relay entry.
// When a processing group which is supposed to deliver a relay entry does not
// fulfill its work, then this node notifies the chain about it. In the case of
// delivering a relay entry by a processing group, this node does nothing.
func (n *node) MonitorRelayEntry(
	relayRequestBlockNumber uint64,
) {
	logger.Infof("monitoring chain for a new relay entry")

	blockCounter, err := n.beaconChain.BlockCounter()
	if err != nil {
		logger.Errorf("failed to get block counter: [%v]", err)
		return
	}

	chainConfig := n.beaconChain.GetConfig()

	timeoutWaiterChannel, err := blockCounter.BlockHeightWaiter(
		relayRequestBlockNumber + chainConfig.RelayEntryTimeout,
	)
	if err != nil {
		logger.Errorf("waiter for a relay entry timeout block failed: [%v]", err)
		return
	}

	onEntrySubmittedChannel := make(chan *event.RelayEntrySubmitted)

	subscription := n.beaconChain.OnRelayEntrySubmitted(
		func(event *event.RelayEntrySubmitted) {
			onEntrySubmittedChannel <- event
		},
	)

	for {
		select {
		case blockNumber := <-timeoutWaiterChannel:
			subscription.Unsubscribe()
			close(onEntrySubmittedChannel)
			logger.Warningf(
				"relay entry was not submitted on time, reporting timeout at block [%v]",
				blockNumber,
			)
			err = n.beaconChain.ReportRelayEntryTimeout()
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
func (n *node) GenerateRelayEntry(
	previousEntry []byte,
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

	// Each signer of the given group should have the same picture of other
	// group operators. Otherwise, they would not be in the group registry.
	// That said, take the group operators from the first signer.
	groupMembers := n.groupRegistry.GetGroup(groupPublicKey)[0].
		Signer.
		GroupOperators()

	membershipValidator := group.NewOperatorsMembershipValidator(
		groupMembers,
		n.beaconChain.Signing(),
	)

	err = channel.SetFilter(membershipValidator.IsInGroup)
	if err != nil {
		logger.Errorf(
			"could not set filter for channel [%v]: [%v]",
			channel.Name(),
			err,
		)
	}

	blockCounter, err := n.beaconChain.BlockCounter()
	if err != nil {
		logger.Errorf("failed to get block counter: [%v]", err)
		return
	}

	chainConfig := n.beaconChain.GetConfig()

	for _, member := range memberships {
		go func(member *registry.Membership) {
			err = entry.SignAndSubmit(
				blockCounter,
				channel,
				n.beaconChain,
				previousEntry,
				chainConfig.HonestThreshold,
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

// channelNameForPublicKey takes group public key represented by marshalled
// G2 point and transforms it into a broadcast channel name.
// Broadcast channel name for group is the hexadecimal representation of
// compressed public key of the group.
func channelNameForPublicKeyBytes(groupPublicKey []byte) (string, error) {
	g2 := new(bn256.G2)

	if _, err := g2.Unmarshal(groupPublicKey); err != nil {
		return "", fmt.Errorf("could not create channel name: [%v]", err)
	}

	return channelNameForPublicKey(g2), nil
}

// channelNameForPublicKey takes group public key represented by G2 point
// and transforms it into a broadcast channel name.
// Broadcast channel name for group is the hexadecimal representation of
// compressed public key of the group.
func channelNameForPublicKey(groupPublicKey *bn256.G2) string {
	altbn128GroupPublicKey := altbn128.G2Point{G2: groupPublicKey}
	return hex.EncodeToString(altbn128GroupPublicKey.Compress())
}
