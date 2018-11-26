package gjkr

import (
	"fmt"
)

type keyGenerationState interface {
	execute([]Message) ([]Message, error)
	next() (keyGenerationState, error)
}

type commitmentState struct {
	member *CommittingMember
}

func (cs *commitmentState) execute() ([]Message, error) {
	var messages []Message

	peerSharesMessages, commitmentsMessage, err := cs.member.CalculateMembersSharesAndCommitments()
	if err != nil {
		return nil, fmt.Errorf("committing state execution failed [%v]", err)
	}

	for _, peerShareMessage := range peerSharesMessages {
		messages = append(messages, peerShareMessage)
	}
	messages = append(messages, commitmentsMessage)

	return messages, nil
}

func (cs *commitmentState) next() (keyGenerationState, error) {
	return &commitmentVerificationState{cs.member}, nil
}

type commitmentVerificationState struct {
	member *CommittingMember
}

func (cvs *commitmentVerificationState) execute(messages []Message) ([]Message, error) {
	var sharesMessages []*PeerSharesMessage
	var commitmentMessages []*MemberCommitmentsMessage
	for _, message := range messages {
		switch msg := message.(type) {
		case *PeerSharesMessage:
			sharesMessages = append(sharesMessages, msg)
		case *MemberCommitmentsMessage:
			commitmentMessages = append(commitmentMessages, msg)
		}
	}

	secretSharesAccusationMessage, err := cvs.member.VerifyReceivedSharesAndCommitmentsMessages(sharesMessages, commitmentMessages)
	if err != nil {
		return nil, fmt.Errorf("commitment verification state execution failed [%v]", err)
	}

	return []Message{secretSharesAccusationMessage}, nil
}

func (cvs *commitmentVerificationState) next() (keyGenerationState, error) {
	// return &joinState{is.member}, nil
	return nil, nil
}
