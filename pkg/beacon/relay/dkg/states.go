package dkg

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
)

type keyGenerationState interface {
	initiate() error
	groupMember() thresholdgroup.BaseMember
	// activePeriod is the period during which this state is active, in blocks.
	activePeriod() int
	receive(msg net.Message) error
	nextState() keyGenerationState
}

// initializationState is the starting state of key generation; it waits for
// activePeriod and then enters joinState. No messages are valid in this state.
type initializationState struct {
	channel net.BroadcastChannel
	member  *thresholdgroup.LocalMember
}

func (is *initializationState) groupMember() thresholdgroup.BaseMember {
	return is.member
}

func (is *initializationState) activePeriod() int { return 15 }

func (is *initializationState) initiate() error {
	return nil
}

func (is *initializationState) receive(msg net.Message) error {
	return fmt.Errorf("unexpected message for initialization state: [%#v]", msg)
}

func (is *initializationState) nextState() keyGenerationState {
	return nil
}
