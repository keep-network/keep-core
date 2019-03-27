package result

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/net"
)

// represents a given state in the state machine for signing dkg results
type signingState = state.State

// resultSigningState is the state during which group members sign their preferred
// dkg result (by hashing their dkg result, and then signing the result), and
// share this over the broadcast channel.
//
// State covers phase 13 of the protocol.
type resultSigningState struct {
	channel net.BroadcastChannel
	member  *SigningMember

	signedHashResults []*DKGResultHashSignatureMessage
}

func (rs *resultSigningState) ActiveBlocks() int { return 3 }

func (rs *resultSigningState) Initiate() error {
	message, err := rs.member.SignDKGResult(nil, nil)
	if err != nil {
		return err
	}
	if err := rs.channel.Send(message); err != nil {
		return err
	}
	return nil
}

func (rs *resultSigningState) Receive(msg net.Message) error {
	// is message from self?
	// is message sender accepted?
	// then add it to our list
	return nil
}

func (rs *resultSigningState) Next() signingState {
	// set up the verification state, phase 14
	return nil
}

func (rs *resultSigningState) MemberIndex() member.Index {
	return rs.member.index
}
