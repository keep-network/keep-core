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
	tssRoundOneStateActiveBlocks = 11

	tssRoundTwoStateDelayBlocks  = 1
	tssRoundTwoStateActiveBlocks = 6

	tssRoundThreeStateDelayBlocks  = 1
	tssRoundThreeStateActiveBlocks = 5

	tssRoundFourStateDelayBlocks  = 1
	tssRoundFourStateActiveBlocks = 5

	tssRoundFiveStateDelayBlocks  = 1
	tssRoundFiveStateActiveBlocks = 5

	tssRoundSixStateDelayBlocks  = 1
	tssRoundSixStateActiveBlocks = 5

	tssRoundSevenStateDelayBlocks  = 1
	tssRoundSevenStateActiveBlocks = 5

	tssRoundEightStateDelayBlocks  = 1
	tssRoundEightStateActiveBlocks = 5

	tssRoundNineStateDelayBlocks  = 1
	tssRoundNineStateActiveBlocks = 5

	finalizationStateDelayBlocks  = 1
	finalizationStateActiveBlocks = 5
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
		tssRoundFourStateActiveBlocks +
		tssRoundFiveStateDelayBlocks +
		tssRoundFiveStateActiveBlocks +
		tssRoundSixStateDelayBlocks +
		tssRoundSixStateActiveBlocks +
		tssRoundSevenStateDelayBlocks +
		tssRoundSevenStateActiveBlocks +
		tssRoundEightStateDelayBlocks +
		tssRoundEightStateActiveBlocks +
		tssRoundNineStateDelayBlocks +
		tssRoundNineStateActiveBlocks +
		finalizationStateDelayBlocks +
		finalizationStateActiveBlocks
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

func (ekpgs *ephemeralKeyPairGenerationState) Next() (state.SyncState, error) {
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

func (skgs *symmetricKeyGenerationState) Receive(net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) Next() (state.SyncState, error) {
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

func (tros *tssRoundOneState) Next() (state.SyncState, error) {
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

func (trts *tssRoundTwoState) Next() (state.SyncState, error) {
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

func (trts *tssRoundThreeState) Next() (state.SyncState, error) {
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
// `tssRoundFourMessage`s are valid in this state.
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

func (trfs *tssRoundFourState) Next() (state.SyncState, error) {
	err := <-trfs.outcomeChan
	if err != nil {
		return nil, err
	}

	return &tssRoundFiveState{
		channel:               trfs.channel,
		member:                trfs.member.initializeTssRoundFive(),
		outcomeChan:           make(chan error),
		previousPhaseMessages: trfs.phaseMessages,
	}, nil
}

func (trfs *tssRoundFourState) MemberIndex() group.MemberIndex {
	return trfs.member.id
}

// tssRoundFiveState is the state during which members broadcast TSS
// round five messages.
// `tssRoundFiveMessage`s are valid in this state.
type tssRoundFiveState struct {
	channel net.BroadcastChannel
	member  *tssRoundFiveMember

	outcomeChan chan error

	previousPhaseMessages []*tssRoundFourMessage

	phaseMessages []*tssRoundFiveMessage
}

func (trfs *tssRoundFiveState) DelayBlocks() uint64 {
	return tssRoundFiveStateDelayBlocks
}

func (trfs *tssRoundFiveState) ActiveBlocks() uint64 {
	return tssRoundFiveStateActiveBlocks
}

func (trfs *tssRoundFiveState) Initiate(ctx context.Context) error {
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
		message, err := trfs.member.tssRoundFive(ctx, trfs.previousPhaseMessages)
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

func (trfs *tssRoundFiveState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundFiveMessage:
		if trfs.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && trfs.member.sessionID == phaseMessage.sessionID {
			trfs.phaseMessages = append(trfs.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (trfs *tssRoundFiveState) Next() (state.SyncState, error) {
	err := <-trfs.outcomeChan
	if err != nil {
		return nil, err
	}

	return &tssRoundSixState{
		channel:               trfs.channel,
		member:                trfs.member.initializeTssRoundSix(),
		outcomeChan:           make(chan error),
		previousPhaseMessages: trfs.phaseMessages,
	}, nil
}

func (trfs *tssRoundFiveState) MemberIndex() group.MemberIndex {
	return trfs.member.id
}

// tssRoundSixState is the state during which members broadcast TSS
// round six messages.
// `tssRoundSixMessage`s are valid in this state.
type tssRoundSixState struct {
	channel net.BroadcastChannel
	member  *tssRoundSixMember

	outcomeChan chan error

	previousPhaseMessages []*tssRoundFiveMessage

	phaseMessages []*tssRoundSixMessage
}

func (trss *tssRoundSixState) DelayBlocks() uint64 {
	return tssRoundSixStateDelayBlocks
}

func (trss *tssRoundSixState) ActiveBlocks() uint64 {
	return tssRoundSixStateActiveBlocks
}

func (trss *tssRoundSixState) Initiate(ctx context.Context) error {
	trss.member.markInactiveMembers(trss.previousPhaseMessages)

	if len(trss.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(trss.member.group.InactiveMemberIDs())
	}

	// TSS computations can be time-consuming and can exceed the current
	// state's time window. The ctx parameter is scoped to the lifetime of
	// the current state so, it can be used as a timeout signal. However,
	// that ctx is cancelled upon state's end only after Initiate returns.
	// In order to make that working, Initiate must trigger the computations
	// in a separate goroutine and return before the end of the state.
	go func() {
		message, err := trss.member.tssRoundSix(ctx, trss.previousPhaseMessages)
		if err != nil {
			trss.outcomeChan <- err
			return
		}

		if err := trss.channel.Send(ctx, message); err != nil {
			trss.outcomeChan <- err
			return
		}

		close(trss.outcomeChan)
	}()

	return nil
}

func (trss *tssRoundSixState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundSixMessage:
		if trss.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && trss.member.sessionID == phaseMessage.sessionID {
			trss.phaseMessages = append(trss.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (trss *tssRoundSixState) Next() (state.SyncState, error) {
	err := <-trss.outcomeChan
	if err != nil {
		return nil, err
	}

	return &tssRoundSevenState{
		channel:               trss.channel,
		member:                trss.member.initializeTssRoundSeven(),
		outcomeChan:           make(chan error),
		previousPhaseMessages: trss.phaseMessages,
	}, nil
}

func (trss *tssRoundSixState) MemberIndex() group.MemberIndex {
	return trss.member.id
}

// tssRoundSevenState is the state during which members broadcast TSS
// round seven messages.
// `tssRoundSevenMessage`s are valid in this state.
type tssRoundSevenState struct {
	channel net.BroadcastChannel
	member  *tssRoundSevenMember

	outcomeChan chan error

	previousPhaseMessages []*tssRoundSixMessage

	phaseMessages []*tssRoundSevenMessage
}

func (trss *tssRoundSevenState) DelayBlocks() uint64 {
	return tssRoundSevenStateDelayBlocks
}

func (trss *tssRoundSevenState) ActiveBlocks() uint64 {
	return tssRoundSevenStateActiveBlocks
}

func (trss *tssRoundSevenState) Initiate(ctx context.Context) error {
	trss.member.markInactiveMembers(trss.previousPhaseMessages)

	if len(trss.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(trss.member.group.InactiveMemberIDs())
	}

	// TSS computations can be time-consuming and can exceed the current
	// state's time window. The ctx parameter is scoped to the lifetime of
	// the current state so, it can be used as a timeout signal. However,
	// that ctx is cancelled upon state's end only after Initiate returns.
	// In order to make that working, Initiate must trigger the computations
	// in a separate goroutine and return before the end of the state.
	go func() {
		message, err := trss.member.tssRoundSeven(ctx, trss.previousPhaseMessages)
		if err != nil {
			trss.outcomeChan <- err
			return
		}

		if err := trss.channel.Send(ctx, message); err != nil {
			trss.outcomeChan <- err
			return
		}

		close(trss.outcomeChan)
	}()

	return nil
}

func (trss *tssRoundSevenState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundSevenMessage:
		if trss.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && trss.member.sessionID == phaseMessage.sessionID {
			trss.phaseMessages = append(trss.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (trss *tssRoundSevenState) Next() (state.SyncState, error) {
	err := <-trss.outcomeChan
	if err != nil {
		return nil, err
	}

	return &tssRoundEightState{
		channel:               trss.channel,
		member:                trss.member.initializeTssRoundEight(),
		outcomeChan:           make(chan error),
		previousPhaseMessages: trss.phaseMessages,
	}, nil
}

func (trss *tssRoundSevenState) MemberIndex() group.MemberIndex {
	return trss.member.id
}

// tssRoundEightState is the state during which members broadcast TSS
// round eight messages.
// `tssRoundEightMessage`s are valid in this state.
type tssRoundEightState struct {
	channel net.BroadcastChannel
	member  *tssRoundEightMember

	outcomeChan chan error

	previousPhaseMessages []*tssRoundSevenMessage

	phaseMessages []*tssRoundEightMessage
}

func (tres *tssRoundEightState) DelayBlocks() uint64 {
	return tssRoundEightStateDelayBlocks
}

func (tres *tssRoundEightState) ActiveBlocks() uint64 {
	return tssRoundEightStateActiveBlocks
}

func (tres *tssRoundEightState) Initiate(ctx context.Context) error {
	tres.member.markInactiveMembers(tres.previousPhaseMessages)

	if len(tres.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(tres.member.group.InactiveMemberIDs())
	}

	// TSS computations can be time-consuming and can exceed the current
	// state's time window. The ctx parameter is scoped to the lifetime of
	// the current state so, it can be used as a timeout signal. However,
	// that ctx is cancelled upon state's end only after Initiate returns.
	// In order to make that working, Initiate must trigger the computations
	// in a separate goroutine and return before the end of the state.
	go func() {
		message, err := tres.member.tssRoundEight(ctx, tres.previousPhaseMessages)
		if err != nil {
			tres.outcomeChan <- err
			return
		}

		if err := tres.channel.Send(ctx, message); err != nil {
			tres.outcomeChan <- err
			return
		}

		close(tres.outcomeChan)
	}()

	return nil
}

func (tres *tssRoundEightState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundEightMessage:
		if tres.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && tres.member.sessionID == phaseMessage.sessionID {
			tres.phaseMessages = append(tres.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (tres *tssRoundEightState) Next() (state.SyncState, error) {
	err := <-tres.outcomeChan
	if err != nil {
		return nil, err
	}

	return &tssRoundNineState{
		channel:               tres.channel,
		member:                tres.member.initializeTssRoundNine(),
		outcomeChan:           make(chan error),
		previousPhaseMessages: tres.phaseMessages,
	}, nil
}

func (tres *tssRoundEightState) MemberIndex() group.MemberIndex {
	return tres.member.id
}

// tssRoundNineState is the state during which members broadcast TSS
// round nine messages.
// `tssRoundNineMessage`s are valid in this state.
type tssRoundNineState struct {
	channel net.BroadcastChannel
	member  *tssRoundNineMember

	outcomeChan chan error

	previousPhaseMessages []*tssRoundEightMessage

	phaseMessages []*tssRoundNineMessage
}

func (trns *tssRoundNineState) DelayBlocks() uint64 {
	return tssRoundNineStateDelayBlocks
}

func (trns *tssRoundNineState) ActiveBlocks() uint64 {
	return tssRoundNineStateActiveBlocks
}

func (trns *tssRoundNineState) Initiate(ctx context.Context) error {
	trns.member.markInactiveMembers(trns.previousPhaseMessages)

	if len(trns.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(trns.member.group.InactiveMemberIDs())
	}

	// TSS computations can be time-consuming and can exceed the current
	// state's time window. The ctx parameter is scoped to the lifetime of
	// the current state so, it can be used as a timeout signal. However,
	// that ctx is cancelled upon state's end only after Initiate returns.
	// In order to make that working, Initiate must trigger the computations
	// in a separate goroutine and return before the end of the state.
	go func() {
		message, err := trns.member.tssRoundNine(ctx, trns.previousPhaseMessages)
		if err != nil {
			trns.outcomeChan <- err
			return
		}

		if err := trns.channel.Send(ctx, message); err != nil {
			trns.outcomeChan <- err
			return
		}

		close(trns.outcomeChan)
	}()

	return nil
}

func (trns *tssRoundNineState) Receive(msg net.Message) error {
	switch phaseMessage := msg.Payload().(type) {
	case *tssRoundNineMessage:
		if trns.member.shouldAcceptMessage(
			phaseMessage.SenderID(),
			msg.SenderPublicKey(),
		) && trns.member.sessionID == phaseMessage.sessionID {
			trns.phaseMessages = append(trns.phaseMessages, phaseMessage)
		}
	}

	return nil
}

func (trns *tssRoundNineState) Next() (state.SyncState, error) {
	err := <-trns.outcomeChan
	if err != nil {
		return nil, err
	}

	return &finalizationState{
		channel:               trns.channel,
		member:                trns.member.initializeFinalization(),
		outcomeChan:           make(chan error),
		previousPhaseMessages: trns.phaseMessages,
	}, nil
}

func (trns *tssRoundNineState) MemberIndex() group.MemberIndex {
	return trns.member.id
}

// finalizationState is the last state of the signing protocol - in this state,
// signing is completed. No messages are valid in this state.
//
// State prepares a result that is returned to the caller.
type finalizationState struct {
	channel net.BroadcastChannel
	member  *finalizingMember

	outcomeChan chan error

	previousPhaseMessages []*tssRoundNineMessage
}

func (fs *finalizationState) DelayBlocks() uint64 {
	return finalizationStateDelayBlocks
}

func (fs *finalizationState) ActiveBlocks() uint64 {
	return finalizationStateActiveBlocks
}

func (fs *finalizationState) Initiate(ctx context.Context) error {
	fs.member.markInactiveMembers(fs.previousPhaseMessages)

	if len(fs.member.group.InactiveMemberIDs()) > 0 {
		return newInactiveMembersError(fs.member.group.InactiveMemberIDs())
	}

	// TSS computations can be time-consuming and can exceed the current
	// state's time window. The ctx parameter is scoped to the lifetime of
	// the current state so, it can be used as a timeout signal. However,
	// that ctx is cancelled upon state's end only after Initiate returns.
	// In order to make that working, Initiate must trigger the computations
	// in a separate goroutine and return before the end of the state.
	go func() {
		err := fs.member.tssFinalize(ctx, fs.previousPhaseMessages)
		if err != nil {
			fs.outcomeChan <- err
			return
		}

		close(fs.outcomeChan)
	}()

	return nil
}

func (fs *finalizationState) Receive(net.Message) error {
	return nil
}

func (fs *finalizationState) Next() (state.SyncState, error) {
	err := <-fs.outcomeChan
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (fs *finalizationState) MemberIndex() group.MemberIndex {
	return fs.member.id
}

func (fs *finalizationState) result() *Result {
	return fs.member.Result()
}
