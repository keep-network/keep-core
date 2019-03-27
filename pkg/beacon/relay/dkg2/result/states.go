package result

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
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

	result                gjkr.Result
	disqualifiedMemberIDs []gjkr.MemberID
	inactiveMemberIDs     []gjkr.MemberID

	signedHashResults []*DKGResultHashSignatureMessage
}

func (rs *resultSigningState) ActiveBlocks() int { return 3 }

func (rs *resultSigningState) Initiate() error {
	rs.disqualifiedMemberIDs = rs.result.Disqualified
	rs.inactiveMemberIDs = rs.result.Inactive

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
	switch signedMessage := msg.Payload().(type) {
	case *DKGResultHashSignatureMessage:
		// ignore messages from ourselves
		if signedMessage.senderIndex == rs.member.index {
			return nil
		}

		// ignore messages from DQ
		for _, disqualifiedMember := range rs.disqualifiedMemberIDs {
			if signedMessage.senderIndex == disqualifiedMember {
				return nil
			}
		}

		// ignore messages from IA
		for _, inactiveMemeber := range rs.inactiveMemberIDs {
			if signedMessage.senderIndex == inactiveMemeber {
				return nil
			}
		}

		// then add it to our list
		rs.signedHashResults = append(rs.signedHashResults)
	}
	return nil
}

func (rs *resultSigningState) Next() signingState {
	// set up the verification state, phase 14
	return nil
}

func (rs *resultSigningState) MemberIndex() member.Index {
	return rs.member.index
}
