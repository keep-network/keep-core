package beacon

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/beacon/entry"
	"github.com/keep-network/keep-core/pkg/beacon/membership"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
)

type participantState int

// FIXME To become something more real...
type libp2pHandle int

// Config contains the config data needed for the beacon to operate.
type Config struct {
	GroupSize int
	Threshold int
}

// ChainInterface represents the interface that the beacon expects to interact
// with the anchoring blockchain on.
type ChainInterface interface {
	GetConfig() Config
}

const (
	unstaked participantState = iota
	staked
	waitingForGroup
	inIncompleteGroup
	inCompleteGroup
	inInitializingGroup
	inInitializedGroup
	inActivatingGroup
	inActiveGroup
)

func initialize() {
	if curParticipantState, err := checkParticipantState(); err != nil {
		panic(fmt.Sprintf("Could not resolve current relay state, aborting: [%s]", err))
	} else {
		switch curParticipantState {
		case unstaked:
			// check for stake command-line parameter to initialize staking?
		default:
			// connect to libp2p
		}
	}
}

func checkParticipantState() (participantState, error) {
	return unstaked, nil
}

func checkNetworkParticipantState(handle libp2pHandle) (participantState, error) {
	// FIXME This will be a real handle, and we will do something real with it.
	fmt.Println(handle)

	// FIXME This will return a real libp2p-based state: are we waiting for a
	// group, in an incomplete one, in a complete one, etc.
	return unstaked, nil
}

func libp2pConnected(handle libp2pHandle) {
	if participantState, err := checkNetworkParticipantState(handle); err != nil {
		panic(fmt.Sprintf("Could not resolve current relay state from libp2p, aborting: [%s]", err))
	} else {
		switch participantState {
		case staked:
			membership.WaitForGroup()
		case waitingForGroup:
			membership.WaitForGroup()
		case inIncompleteGroup:
			membership.WaitForGroupCompletion()
		case inCompleteGroup:
			membership.InitializeMembership()
		case inInitializingGroup:
			membership.InitializeMembership()
		case inInitializedGroup:
			membership.ActivateMembership()
		case inActivatingGroup:
			membership.ActivateMembership()
		case inActiveGroup:
			// FIXME We should have a non-empty state at this point ;)
			entry.ServeRequests(relay.EmptyState())
		default:
			panic(fmt.Sprintf("Unexpected participant state [%d].", participantState))
		}
	}
}
