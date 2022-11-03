package signing

import (
	"context"
	"strconv"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

// ephemeralKeyPairGenerationState is the state during which members broadcast
// public ephemeral keys generated for other members of the group.
// `ephemeralPublicKeyMessage`s are valid in this state.
type ephemeralKeyPairGenerationState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *ephemeralKeyPairGeneratingMember
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

func (ekpgs *ephemeralKeyPairGenerationState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if ekpgs.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && ekpgs.member.sessionID == protocolMessage.SessionID() {
			ekpgs.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) CanTransition() bool {
	messagingDone := len(receivedMessages[*ephemeralPublicKeyMessage](ekpgs.BaseAsyncState)) ==
		len(ekpgs.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (ekpgs *ephemeralKeyPairGenerationState) Next() (state.AsyncState, error) {
	return &symmetricKeyGenerationState{
		BaseAsyncState: ekpgs.BaseAsyncState,
		channel:        ekpgs.channel,
		member:         ekpgs.member.initializeSymmetricKeyGeneration(),
	}, nil
}

func (ekpgs *ephemeralKeyPairGenerationState) MemberIndex() group.MemberIndex {
	return ekpgs.member.id
}

// symmetricKeyGenerationState is the state during which members compute
// symmetric keys from the previously exchanged ephemeral public keys.
// No messages are valid in this state.
type symmetricKeyGenerationState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *symmetricKeyGeneratingMember
}

func (skgs *symmetricKeyGenerationState) Initiate(ctx context.Context) error {
	return skgs.member.generateSymmetricKeys(
		receivedMessages[*ephemeralPublicKeyMessage](skgs.BaseAsyncState),
	)
}

func (skgs *symmetricKeyGenerationState) Receive(net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) CanTransition() bool {
	return true
}

func (skgs *symmetricKeyGenerationState) Next() (state.AsyncState, error) {
	return &tssRoundOneState{
		BaseAsyncState: skgs.BaseAsyncState,
		channel:        skgs.channel,
		member:         skgs.member.initializeTssRoundOne(),
	}, nil
}

func (skgs *symmetricKeyGenerationState) MemberIndex() group.MemberIndex {
	return skgs.member.id
}

// tssRoundOneState is the state during which members broadcast TSS
// round one messages.
// `tssRoundOneMessage`s are valid in this state.
type tssRoundOneState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundOneMember
}

func (tros *tssRoundOneState) Initiate(ctx context.Context) error {
	message, err := tros.member.tssRoundOne(ctx)
	if err != nil {
		return err
	}

	if err := tros.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (tros *tssRoundOneState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if tros.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && tros.member.sessionID == protocolMessage.SessionID() {
			tros.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (tros *tssRoundOneState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundOneCompositeMessage](tros.BaseAsyncState)) ==
		len(tros.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (tros *tssRoundOneState) Next() (state.AsyncState, error) {
	return &tssRoundTwoState{
		BaseAsyncState: tros.BaseAsyncState,
		channel:        tros.channel,
		member:         tros.member.initializeTssRoundTwo(),
	}, nil
}

func (tros *tssRoundOneState) MemberIndex() group.MemberIndex {
	return tros.member.id
}

// tssRoundTwoState is the state during which members broadcast TSS
// round two messages.
// `tssRoundTwoMessage`s are valid in this state.
type tssRoundTwoState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundTwoMember
}

func (trts *tssRoundTwoState) Initiate(ctx context.Context) error {
	message, err := trts.member.tssRoundTwo(
		ctx,
		receivedMessages[*tssRoundOneCompositeMessage](trts.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := trts.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trts *tssRoundTwoState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if trts.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && trts.member.sessionID == protocolMessage.SessionID() {
			trts.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (trts *tssRoundTwoState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundTwoCompositeMessage](trts.BaseAsyncState)) ==
		len(trts.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (trts *tssRoundTwoState) Next() (state.AsyncState, error) {
	return &tssRoundThreeState{
		BaseAsyncState: trts.BaseAsyncState,
		channel:        trts.channel,
		member:         trts.member.initializeTssRoundThree(),
	}, nil
}

func (trts *tssRoundTwoState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// tssRoundThreeState is the state during which members broadcast TSS
// round three messages.
// `tssRoundThreeMessage`s are valid in this state.
type tssRoundThreeState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundThreeMember
}

func (trts *tssRoundThreeState) Initiate(ctx context.Context) error {
	message, err := trts.member.tssRoundThree(
		ctx,
		receivedMessages[*tssRoundTwoCompositeMessage](trts.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := trts.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trts *tssRoundThreeState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if trts.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && trts.member.sessionID == protocolMessage.SessionID() {
			trts.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (trts *tssRoundThreeState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundThreeCompositeMessage](trts.BaseAsyncState)) ==
		len(trts.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (trts *tssRoundThreeState) Next() (state.AsyncState, error) {
	return &tssRoundFourState{
		BaseAsyncState: trts.BaseAsyncState,
		channel:        trts.channel,
		member:         trts.member.initializeTssRoundFour(),
	}, nil
}

func (trts *tssRoundThreeState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// tssRoundFourState is the state during which members broadcast TSS
// round four messages.
// `tssRoundFourMessage`s are valid in this state.
type tssRoundFourState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundFourMember
}

func (trfs *tssRoundFourState) Initiate(ctx context.Context) error {
	message, err := trfs.member.tssRoundFour(
		ctx,
		receivedMessages[*tssRoundThreeCompositeMessage](trfs.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := trfs.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trfs *tssRoundFourState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if trfs.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && trfs.member.sessionID == protocolMessage.SessionID() {
			trfs.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (trfs *tssRoundFourState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundFourCompositeMessage](trfs.BaseAsyncState)) ==
		len(trfs.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (trfs *tssRoundFourState) Next() (state.AsyncState, error) {
	return &tssRoundFiveState{
		BaseAsyncState: trfs.BaseAsyncState,
		channel:        trfs.channel,
		member:         trfs.member.initializeTssRoundFive(),
	}, nil
}

func (trfs *tssRoundFourState) MemberIndex() group.MemberIndex {
	return trfs.member.id
}

// tssRoundFiveState is the state during which members broadcast TSS
// round five messages.
// `tssRoundFiveMessage`s are valid in this state.
type tssRoundFiveState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundFiveMember
}

func (trfs *tssRoundFiveState) Initiate(ctx context.Context) error {
	message, err := trfs.member.tssRoundFive(
		ctx,
		receivedMessages[*tssRoundFourCompositeMessage](trfs.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := trfs.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trfs *tssRoundFiveState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if trfs.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && trfs.member.sessionID == protocolMessage.SessionID() {
			trfs.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (trfs *tssRoundFiveState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundFiveCompositeMessage](trfs.BaseAsyncState)) ==
		len(trfs.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (trfs *tssRoundFiveState) Next() (state.AsyncState, error) {
	return &tssRoundSixState{
		BaseAsyncState: trfs.BaseAsyncState,
		channel:        trfs.channel,
		member:         trfs.member.initializeTssRoundSix(),
	}, nil
}

func (trfs *tssRoundFiveState) MemberIndex() group.MemberIndex {
	return trfs.member.id
}

// tssRoundSixState is the state during which members broadcast TSS
// round six messages.
// `tssRoundSixMessage`s are valid in this state.
type tssRoundSixState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundSixMember
}

func (trss *tssRoundSixState) Initiate(ctx context.Context) error {
	message, err := trss.member.tssRoundSix(
		ctx,
		receivedMessages[*tssRoundFiveCompositeMessage](trss.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := trss.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trss *tssRoundSixState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if trss.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && trss.member.sessionID == protocolMessage.SessionID() {
			trss.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (trss *tssRoundSixState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundSixCompositeMessage](trss.BaseAsyncState)) ==
		len(trss.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (trss *tssRoundSixState) Next() (state.AsyncState, error) {
	return &tssRoundSevenState{
		BaseAsyncState: trss.BaseAsyncState,
		channel:        trss.channel,
		member:         trss.member.initializeTssRoundSeven(),
	}, nil
}

func (trss *tssRoundSixState) MemberIndex() group.MemberIndex {
	return trss.member.id
}

// tssRoundSevenState is the state during which members broadcast TSS
// round seven messages.
// `tssRoundSevenMessage`s are valid in this state.
type tssRoundSevenState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundSevenMember
}

func (trss *tssRoundSevenState) Initiate(ctx context.Context) error {
	message, err := trss.member.tssRoundSeven(
		ctx,
		receivedMessages[*tssRoundSixCompositeMessage](trss.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := trss.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trss *tssRoundSevenState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if trss.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && trss.member.sessionID == protocolMessage.SessionID() {
			trss.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (trss *tssRoundSevenState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundSevenCompositeMessage](trss.BaseAsyncState)) ==
		len(trss.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (trss *tssRoundSevenState) Next() (state.AsyncState, error) {
	return &tssRoundEightState{
		BaseAsyncState: trss.BaseAsyncState,
		channel:        trss.channel,
		member:         trss.member.initializeTssRoundEight(),
	}, nil
}

func (trss *tssRoundSevenState) MemberIndex() group.MemberIndex {
	return trss.member.id
}

// tssRoundEightState is the state during which members broadcast TSS
// round eight messages.
// `tssRoundEightMessage`s are valid in this state.
type tssRoundEightState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundEightMember
}

func (tres *tssRoundEightState) Initiate(ctx context.Context) error {
	message, err := tres.member.tssRoundEight(
		ctx,
		receivedMessages[*tssRoundSevenCompositeMessage](tres.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := tres.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (tres *tssRoundEightState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if tres.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && tres.member.sessionID == protocolMessage.SessionID() {
			tres.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (tres *tssRoundEightState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundEightCompositeMessage](tres.BaseAsyncState)) ==
		len(tres.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (tres *tssRoundEightState) Next() (state.AsyncState, error) {
	return &tssRoundNineState{
		BaseAsyncState: tres.BaseAsyncState,
		channel:        tres.channel,
		member:         tres.member.initializeTssRoundNine(),
	}, nil
}

func (tres *tssRoundEightState) MemberIndex() group.MemberIndex {
	return tres.member.id
}

// tssRoundNineState is the state during which members broadcast TSS
// round nine messages.
// `tssRoundNineMessage`s are valid in this state.
type tssRoundNineState struct {
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *tssRoundNineMember
}

func (trns *tssRoundNineState) Initiate(ctx context.Context) error {
	message, err := trns.member.tssRoundNine(
		ctx,
		receivedMessages[*tssRoundEightCompositeMessage](trns.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	if err := trns.channel.Send(ctx, message); err != nil {
		return err
	}

	return nil
}

func (trns *tssRoundNineState) Receive(netMessage net.Message) error {
	if protocolMessage, ok := netMessage.Payload().(message); ok {
		if trns.member.shouldAcceptMessage(
			protocolMessage.SenderID(),
			netMessage.SenderPublicKey(),
		) && trns.member.sessionID == protocolMessage.SessionID() {
			trns.ReceiveToHistory(netMessage)
		}
	}

	return nil
}

func (trns *tssRoundNineState) CanTransition() bool {
	messagingDone := len(receivedMessages[*tssRoundNineCompositeMessage](trns.BaseAsyncState)) ==
		len(trns.member.group.OperatingMemberIDs())-1

	return messagingDone
}

func (trns *tssRoundNineState) Next() (state.AsyncState, error) {
	return &finalizationState{
		BaseAsyncState: trns.BaseAsyncState,
		channel:        trns.channel,
		member:         trns.member.initializeFinalization(),
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
	*state.BaseAsyncState

	channel net.BroadcastChannel
	member  *finalizingMember
}

func (fs *finalizationState) Initiate(ctx context.Context) error {
	err := fs.member.tssFinalize(
		ctx,
		receivedMessages[*tssRoundNineCompositeMessage](fs.BaseAsyncState),
	)
	if err != nil {
		return err
	}

	return nil
}

func (fs *finalizationState) Receive(net.Message) error {
	return nil
}

func (fs *finalizationState) CanTransition() bool {
	return true
}

func (fs *finalizationState) Next() (state.AsyncState, error) {
	return nil, nil
}

func (fs *finalizationState) MemberIndex() group.MemberIndex {
	return fs.member.id
}

func (fs *finalizationState) result() *Result {
	return fs.member.Result()
}

// receivedMessages returns all messages of type T that have been received
// and validated so far. Returned messages are deduplicated so there is a
// guarantee that only one message of the given type is returned for the
// given sender.
func receivedMessages[T message](base *state.BaseAsyncState) []T {
	var messageTemplate T

	payloads := state.ExtractMessagesPayloads[T](base, messageTemplate.Type())

	return state.DeduplicateMessagesPayloads(
		payloads,
		func(message T) string {
			return strconv.Itoa(int(message.SenderID()))
		},
	)
}
