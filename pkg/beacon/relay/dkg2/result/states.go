package result

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/beacon/relay/states"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

type resultState states.DKGState

type resultSigningState struct {
	channel     net.BroadcastChannel
	member      *SigningMember
	chainHandle chain.Handle

	gjkrResult *gjkr.Result

	phaseMessages []*DKGResultHashSignatureMessage
}

func (rss *resultSigningState) activeBlocks() int { return 3 }

func (rss *resultSigningState) initiate() error {
	chainResult := convertResult(rss.gjkrResult, rss.groupSize) // TODO: Move to chain specific package?

	message, err := rss.member.SignDKGResult(chainResult, rss.chainHandle)
	if err != nil {
		return err
	}

	if err := rss.channel.Send(message); err != nil {
		return err
	}

	return nil
}

func (rss *resultSigningState) receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *DKGResultHashSignatureMessage:
		if !isMessageFromSelf(rss, phaseMessage) &&
			isSenderAccepted(rss.member, phaseMessage) {
			rss.phaseMessages = append(rss.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (rss *resultSigningState) nextState() resultState {
	return &resultSignaturesVerificationState{
		channel: rss.channel,
		member:  rss.member.InitializeSignaturesVerification(),

		previousPhaseMessages: rss.phaseMessages,
	}
}

func (rss *resultSigningState) memberID() gjkr.MemberID {
	return rss.member.index
}

// commitmentsVerificationState is the state during which members validate
// shares and commitments computed and published by other members in the
// previous phase. `gjkr.SecretShareAccusationMessage`s are valid in this state.
//
// State covers phase 4 of the protocol.
type resultSignaturesVerificationState struct {
	channel net.BroadcastChannel
	member  *SignaturesVerifyingMember

	previousPhaseMessages         []*DKGResultHashSignatureMessage
	receivedValidResultSignatures map[gjkr.MemberID]operator.Signature
}

func (rsvs *resultSignaturesVerificationState) activeBlocks() int { return 3 }

func (rsvs *resultSignaturesVerificationState) initiate() error {
	cvs.member.MarkInactiveMembers(
		cvs.previousPhaseMessages,
	)

	cvs.receivedValidResultSignatures, err = cvs.member.VerifyReceivedSharesAndCommitmentsMessages(
		cvs.previousPhaseMessages,
	)
	if err != nil {
		return err
	}

	return nil
}

func (rsvs *resultSignaturesVerificationState) receive(msg net.Message) error {
	return nil
}

func (rsvs *resultSignaturesVerificationState) nextState() resultState {
	return nil
}

func (rsvs *resultSignaturesVerificationState) memberID() gjkr.MemberID {
	return cvs.member.Index
}
