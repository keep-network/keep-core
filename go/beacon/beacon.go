package beacon

import (
	"fmt"

	"./membership"
)

type participantState int

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
		fmt.Printf("Could not resolve current relay state, aborting: [%s]", err)
		panic("Unknown participant state, aborting.")
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

func libp2pConnected() {
	// get
	participantState := staked
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
		startRelay()
	default:
		panic(fmt.Sprintf("Unexpected participant state [%d].", participantState))
	}
}

func startRelay() {
	// Start listening for beacon requests on chain.
}
