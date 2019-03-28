package result

import (
	"math/big"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
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
	channel    net.BroadcastChannel
	relayChain relayChain.Interface

	member *SigningMember

	requestID             *big.Int
	result                *relayChain.DKGResult
	disqualifiedMemberIDs []gjkr.MemberID
	inactiveMemberIDs     []gjkr.MemberID

	signatureMessages []*DKGResultHashSignatureMessage
}

func (rss *resultSigningState) ActiveBlocks() int { return 3 }

func (rss *resultSigningState) Initiate() error {
	message, err := rss.member.SignDKGResult(rss.result, rss.relayChain)
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
		rss.signatureMessages = append(rss.signatureMessages, signedMessage)
	}
	return nil
}

func (rss *resultSigningState) Next() signingState {
	// set up the verification state, phase 13 part 2
	return &signaturesVerificationState{
		channel:           rss.channel,
		relayChain:        rss.relayChain,
		member:            rss.member,
		requestID:         rss.requestID,
		result:            rss.result,
		signatureMessages: rss.signatureMessages,
		validSignatures:   make(map[member.Index]operator.Signature),
	}

}

func (rss *resultSigningState) MemberIndex() member.Index {
	return rss.member.index
}

// signaturesVerificationState is the state during which group members verify all validSignatures
// that valid submitters sent over the broadcast channel in the previous state.
// Valid validSignatures are added to the state.
//
// State is part of phase 13 of the protocol.
type signaturesVerificationState struct {
	channel    net.BroadcastChannel
	relayChain relayChain.Interface

	member *SigningMember

	requestID *big.Int
	result    *relayChain.DKGResult

	signatureMessages []*DKGResultHashSignatureMessage
	validSignatures   map[member.Index]operator.Signature
}

func (svs *signaturesVerificationState) ActiveBlocks() int { return 0 }

func (svs *signaturesVerificationState) Initiate() error {
	signatures, err := svs.member.VerifyDKGResultSignatures(svs.signatureMessages)
	if err != nil {
		return err
	}

	svs.validSignatures = signatures
	return nil
}

func (svs *signaturesVerificationState) Receive(msg net.Message) error {
	return nil
}

func (svs *signaturesVerificationState) Next() signingState {
	return &resultSubmissionState{
		channel:    svs.channel,
		relayChain: svs.relayChain,
		member:     NewSubmittingMember(svs.member.index),
		requestID:  svs.requestID,
		result:     svs.result,
		signatures: svs.validSignatures,
	}

}

func (svs *signaturesVerificationState) MemberIndex() member.Index {
	return svs.member.index
}

// resultSubmissionState is the state during which group members submit the dkg
// result to the chain. This state concludes the DKG protocol.
//
// State covers, the final phase, phase 14 of the protocol.
type resultSubmissionState struct {
	channel    net.BroadcastChannel
	relayChain relayChain.Interface

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
		rss.relayChain,
	)
}

func (rss *resultSubmissionState) Receive(msg net.Message) error {
	return nil
}

func (rss *resultSubmissionState) Next() signingState {
	// returning nil represents this is the final state
	return nil
}

func (rss *resultSubmissionState) MemberIndex() member.Index {
	return rss.member.index
}
