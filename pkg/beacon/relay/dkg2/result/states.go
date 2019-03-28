package result

import (
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

// represents a given state in the state machine for signing dkg results
type signingState = state.State

// resultSigningState is the state during which group members sign their preferred
// dkg result (by hashing their dkg result, and then signing the result), and
// share this over the broadcast channel.
//
// State is part of phase 13 of the protocol.
type resultSigningState struct {
	channel net.BroadcastChannel
	handle  chain.Handle

	member *SigningMember

	requestID             *big.Int
	result                *relayChain.DKGResult
	disqualifiedMemberIDs []gjkr.MemberID
	inactiveMemberIDs     []gjkr.MemberID

	signedHashResults []*DKGResultHashSignatureMessage
}

func (rss *resultSigningState) ActiveBlocks() int { return 3 }

func (rss *resultSigningState) Initiate() error {
	message, err := rss.member.SignDKGResult(rss.result, rss.handle)
	if err != nil {
		return err
	}
	if err := rss.channel.Send(message); err != nil {
		return err
	}
	return nil
}

func (rss *resultSigningState) Receive(msg net.Message) error {
	switch signedMessage := msg.Payload().(type) {
	case *DKGResultHashSignatureMessage:
		// ignore messages from ourselves
		if signedMessage.senderIndex == rss.member.index {
			return nil
		}

		// ignore messages from DQ
		for _, disqualifiedMember := range rss.disqualifiedMemberIDs {
			if signedMessage.senderIndex == disqualifiedMember {
				return nil
			}
		}

		// ignore messages from IA
		for _, inactiveMemeber := range rss.inactiveMemberIDs {
			if signedMessage.senderIndex == inactiveMemeber {
				return nil
			}
		}

		// then add it to our list
		rss.signedHashResults = append(rss.signedHashResults)
	}
	return nil
}

func (rss *resultSigningState) Next() signingState {
	// set up the verification state, phase 13 part 2
	return &verificationState{rss, nil}

}

func (rss *resultSigningState) MemberIndex() member.Index {
	return rss.member.index
}

// verificationState is the state during which group members verify all signatures
// that valid submitters sent over the broadcast channel in the previous state.
// Valid signatures are added to the state.
//
// State is part of phase 13 of the protocol.
type verificationState struct {
	*resultSigningState

	signatures map[member.Index]operator.Signature
}

func (vs *verificationState) ActiveBlocks() int { return 0 }

func (vs *verificationState) Initiate() error {
	signatures, err := vs.member.VerifyDKGResultSignatures(vs.signedHashResults)
	if err != nil {
		return err
	}

	vs.signatures = signatures
	return nil
}

func (vs *verificationState) Receive(msg net.Message) error {
	return nil
}

func (vs *verificationState) Next() signingState {
	return &resultSubmissionState{
		channel:    vs.channel,
		handle:     vs.handle,
		member:     NewSubmittingMember(vs.member.index),
		requestID:  vs.requestID,
		result:     vs.result,
		signatures: vs.signatures,
	}

}

func (vs *verificationState) MemberIndex() member.Index {
	return vs.member.index
}

// resultSubmissionState is the state during which group members submit the dkg
// result to the chain.
//
// State covers phase 14 of the protocol.
type resultSubmissionState struct {
	channel net.BroadcastChannel
	handle  chain.Handle

	member *SubmittingMember

	requestID  *big.Int
	result     *relayChain.DKGResult
	signatures map[member.Index]operator.Signature
}

func (rss *resultSubmissionState) ActiveBlocks() int { return 3 }

func (rss *resultSubmissionState) Initiate() error {
	return rss.member.SubmitDKGResult(
		rss.requestID,
		rss.result,
		rss.signatures,
		rss.handle,
	)
}

func (rss *resultSubmissionState) Receive(msg net.Message) error {
	return nil
}

func (rss *resultSubmissionState) Next() signingState {
	return nil
}

func (rss *resultSubmissionState) MemberIndex() member.Index {
	return rss.member.index
}
