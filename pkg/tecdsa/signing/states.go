package signing

import (
	"context"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

const (
	silentStateDelayBlocks  = 0
	silentStateActiveBlocks = 0

	ephemeralKeyPairStateDelayBlocks  = 1
	ephemeralKeyPairStateActiveBlocks = 5

	tssRoundOneStateDelayBlocks  = 1
	tssRoundOneStateActiveBlocks = 5

	tssRoundTwoStateDelayBlocks  = 1
	tssRoundTwoStateActiveBlocks = 5

	tssRoundThreeStateDelayBlocks  = 1
	tssRoundThreeStateActiveBlocks = 5

	tssRoundFourStateDelayBlocks  = 1
	tssRoundFourStateActiveBlocks = 5
)

// ProtocolBlocks returns the total number of blocks it takes to execute
// all the required work defined by the signing protocol.
func ProtocolBlocks() uint64 {
	return ephemeralKeyPairStateDelayBlocks +
		ephemeralKeyPairStateActiveBlocks +
		tssRoundOneStateDelayBlocks +
		tssRoundOneStateActiveBlocks +
		tssRoundTwoStateDelayBlocks +
		tssRoundTwoStateActiveBlocks +
		tssRoundThreeStateDelayBlocks +
		tssRoundThreeStateActiveBlocks +
		tssRoundFourStateDelayBlocks +
		tssRoundFourStateActiveBlocks
}

// ephemeralKeyPairGenerationState is the state during which members broadcast
// public ephemeral keys generated for other members of the group.
// `ephemeralPublicKeyMessage`s are valid in this state.
type ephemeralKeyPairGenerationState struct {
	channel net.BroadcastChannel
	member  *ephemeralKeyPairGeneratingMember

	phaseMessages []*ephemeralPublicKeyMessage
}

func (ekpgs *ephemeralKeyPairGenerationState) DelayBlocks() uint64 {
	return ephemeralKeyPairStateDelayBlocks
}

func (ekpgs *ephemeralKeyPairGenerationState) ActiveBlocks() uint64 {
	return ephemeralKeyPairStateActiveBlocks
}

func (ekpgs *ephemeralKeyPairGenerationState) Initiate(ctx context.Context) error {
	message, err := ekpgs.member.generateEphemeralKeyPair()
	if err != nil {
		return err
	}

	if err := ekpgs.channel.Send(ctx, message); err != nil {
		return err
	}
	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *ephemeralPublicKeyMessage:
		if ekpgs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && ekpgs.member.sessionID == phaseMessage.sessionID {
			ekpgs.phaseMessages = append(ekpgs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) Next() (state.State, error) {
	return &symmetricKeyGenerationState{
		channel:               ekpgs.channel,
		member:                ekpgs.member.initializeSymmetricKeyGeneration(),
		previousPhaseMessages: ekpgs.phaseMessages,
	}, nil
}

func (ekpgs *ephemeralKeyPairGenerationState) MemberIndex() group.MemberIndex {
	return ekpgs.member.id
}

// symmetricKeyGenerationState is the state during which members compute
// symmetric keys from the previously exchanged ephemeral public keys.
// No messages are valid in this state.
type symmetricKeyGenerationState struct {
	channel net.BroadcastChannel
	member  *symmetricKeyGeneratingMember

	previousPhaseMessages []*ephemeralPublicKeyMessage
}

func (skgs *symmetricKeyGenerationState) DelayBlocks() uint64 {
	return silentStateDelayBlocks
}

func (skgs *symmetricKeyGenerationState) ActiveBlocks() uint64 {
	return silentStateActiveBlocks
}

func (skgs *symmetricKeyGenerationState) Initiate(ctx context.Context) error {
	skgs.member.markInactiveMembers(skgs.previousPhaseMessages)

	if len(skgs.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(skgs.member.group.InactiveMemberIDs())
	}

	return skgs.member.generateSymmetricKeys(skgs.previousPhaseMessages)
}

func (skgs *symmetricKeyGenerationState) Receive(msg net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) Next() (state.State, error) {
	return &tssRoundOneState{
		channel:     skgs.channel,
		member:      skgs.member.initializeTssRoundOne(),
		outcomeChan: make(chan error),
	}, nil
}

func (skgs *symmetricKeyGenerationState) MemberIndex() group.MemberIndex {
	return skgs.member.id
}

// tssRoundOneState is the state during which members broadcast TSS
// round one messages.
// `tssRoundOneMessage`s are valid in this state.
type tssRoundOneState struct {
	channel net.BroadcastChannel
	member  *tssRoundOneMember

	outcomeChan chan error

	phaseMessages []*tssRoundOneMessage
}

func (tros *tssRoundOneState) DelayBlocks() uint64 {
	return tssRoundOneStateDelayBlocks
}

func (tros *tssRoundOneState) ActiveBlocks() uint64 {
	return tssRoundOneStateActiveBlocks
}

func (tros *tssRoundOneState) Initiate(ctx context.Context) error {
	// TSS computations can be time-consuming and can exceed the current
	// state's time window. The ctx parameter is scoped to the lifetime of
	// the current state so, it can be used as a timeout signal. However,
	// that ctx is cancelled upon state's end only after Initiate returns.
	// In order to make that working, Initiate must trigger the computations
	// in a separate goroutine and return before the end of the state.
	go func() {
		message, err := tros.member.tssRoundOne(ctx)
		if err != nil {
			tros.outcomeChan <- err
			return
		}

		if err := tros.channel.Send(ctx, message); err != nil {
			tros.outcomeChan <- err
			return
		}

		close(tros.outcomeChan)
	}()

	return nil
}

func (tros *tssRoundOneState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundOneMessage:
		if tros.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && tros.member.sessionID == phaseMessage.sessionID {
			tros.phaseMessages = append(tros.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (tros *tssRoundOneState) Next() (state.State, error) {
	err := <-tros.outcomeChan
	if err != nil {
		return nil, err
	}

	return &tssRoundTwoState{
		channel:               tros.channel,
		member:                tros.member.initializeTssRoundTwo(),
		outcomeChan:           make(chan error),
		previousPhaseMessages: tros.phaseMessages,
	}, nil
}

func (tros *tssRoundOneState) MemberIndex() group.MemberIndex {
	return tros.member.id
}

// tssRoundTwoState is the state during which members broadcast TSS
// round two messages.
// `tssRoundTwoMessage`s are valid in this state.
type tssRoundTwoState struct {
	channel net.BroadcastChannel
	member  *tssRoundTwoMember

	outcomeChan chan error

	previousPhaseMessages []*tssRoundOneMessage

	phaseMessages []*tssRoundTwoMessage
}

func (trts *tssRoundTwoState) DelayBlocks() uint64 {
	return tssRoundTwoStateDelayBlocks
}

func (trts *tssRoundTwoState) ActiveBlocks() uint64 {
	return tssRoundTwoStateActiveBlocks
}

func (trts *tssRoundTwoState) Initiate(ctx context.Context) error {
	trts.member.markInactiveMembers(trts.previousPhaseMessages)

	if len(trts.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(trts.member.group.InactiveMemberIDs())
	}

	// TSS computations can be time-consuming and can exceed the current
	// state's time window. The ctx parameter is scoped to the lifetime of
	// the current state so, it can be used as a timeout signal. However,
	// that ctx is cancelled upon state's end only after Initiate returns.
	// In order to make that working, Initiate must trigger the computations
	// in a separate goroutine and return before the end of the state.
	go func() {
		message, err := trts.member.tssRoundTwo(ctx, trts.previousPhaseMessages)
		if err != nil {
			trts.outcomeChan <- err
			return
		}

		if err := trts.channel.Send(ctx, message); err != nil {
			trts.outcomeChan <- err
			return
		}

		close(trts.outcomeChan)
	}()

	return nil
}

func (trts *tssRoundTwoState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundTwoMessage:
		if trts.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && trts.member.sessionID == phaseMessage.sessionID {
			trts.phaseMessages = append(trts.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (trts *tssRoundTwoState) Next() (state.State, error) {
	err := <-trts.outcomeChan
	if err != nil {
		return nil, err
	}

	return &tssRoundThreeState{
		channel:               trts.channel,
		member:                trts.member.initializeTssRoundThree(),
		outcomeChan:           make(chan error),
		previousPhaseMessages: trts.phaseMessages,
	}, nil
}

func (trts *tssRoundTwoState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// tssRoundThreeState is the state during which members broadcast TSS
// round three messages.
// `tssRoundThreeMessage`s are valid in this state.
type tssRoundThreeState struct {
	channel net.BroadcastChannel
	member  *tssRoundThreeMember

	outcomeChan chan error

	previousPhaseMessages []*tssRoundTwoMessage

	phaseMessages []*tssRoundThreeMessage
}

func (trts *tssRoundThreeState) DelayBlocks() uint64 {
	return tssRoundThreeStateDelayBlocks
}

func (trts *tssRoundThreeState) ActiveBlocks() uint64 {
	return tssRoundThreeStateActiveBlocks
}

func (trts *tssRoundThreeState) Initiate(ctx context.Context) error {
	trts.member.markInactiveMembers(trts.previousPhaseMessages)

	if len(trts.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(trts.member.group.InactiveMemberIDs())
	}

	// TSS computations can be time-consuming and can exceed the current
	// state's time window. The ctx parameter is scoped to the lifetime of
	// the current state so, it can be used as a timeout signal. However,
	// that ctx is cancelled upon state's end only after Initiate returns.
	// In order to make that working, Initiate must trigger the computations
	// in a separate goroutine and return before the end of the state.
	go func() {
		message, err := trts.member.tssRoundThree(ctx, trts.previousPhaseMessages)
		if err != nil {
			trts.outcomeChan <- err
			return
		}

		if err := trts.channel.Send(ctx, message); err != nil {
			trts.outcomeChan <- err
			return
		}

		close(trts.outcomeChan)
	}()

	return nil
}

func (trts *tssRoundThreeState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundThreeMessage:
		if trts.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && trts.member.sessionID == phaseMessage.sessionID {
			trts.phaseMessages = append(trts.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (trts *tssRoundThreeState) Next() (state.State, error) {
	err := <-trts.outcomeChan
	if err != nil {
		return nil, err
	}

	return &tssRoundFourState{
		channel:               trts.channel,
		member:                trts.member.initializeTssRoundFour(),
		outcomeChan:           make(chan error),
		previousPhaseMessages: trts.phaseMessages,
	}, nil
}

func (trts *tssRoundThreeState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// tssRoundFourState is the state during which members broadcast TSS
// round four messages.
// `tssRoundMessage`s are valid in this state.
type tssRoundFourState struct {
	channel net.BroadcastChannel
	member  *tssRoundFourMember

	outcomeChan chan error

	previousPhaseMessages []*tssRoundThreeMessage

	phaseMessages []*tssRoundFourMessage
}

func (trfs *tssRoundFourState) DelayBlocks() uint64 {
	return tssRoundFourStateDelayBlocks
}

func (trfs *tssRoundFourState) ActiveBlocks() uint64 {
	return tssRoundFourStateActiveBlocks
}

func (trfs *tssRoundFourState) Initiate(ctx context.Context) error {
	trfs.member.markInactiveMembers(trfs.previousPhaseMessages)

	if len(trfs.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(trfs.member.group.InactiveMemberIDs())
	}

	// TSS computations can be time-consuming and can exceed the current
	// state's time window. The ctx parameter is scoped to the lifetime of
	// the current state so, it can be used as a timeout signal. However,
	// that ctx is cancelled upon state's end only after Initiate returns.
	// In order to make that working, Initiate must trigger the computations
	// in a separate goroutine and return before the end of the state.
	go func() {
		message, err := trfs.member.tssRoundFour(ctx, trfs.previousPhaseMessages)
		if err != nil {
			trfs.outcomeChan <- err
			return
		}

		if err := trfs.channel.Send(ctx, message); err != nil {
			trfs.outcomeChan <- err
			return
		}

		close(trfs.outcomeChan)
	}()

	return nil
}

func (trfs *tssRoundFourState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundFourMessage:
		if trfs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && trfs.member.sessionID == phaseMessage.sessionID {
			trfs.phaseMessages = append(trfs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (trfs *tssRoundFourState) Next() (state.State, error) {
	err := <-trfs.outcomeChan
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (trfs *tssRoundFourState) MemberIndex() group.MemberIndex {
	return trfs.member.id
}
