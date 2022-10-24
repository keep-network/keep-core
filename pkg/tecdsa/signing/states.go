package signing

import (
	"context"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/faststate"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"sync"
)

// ephemeralKeyPairGenerationState is the state during which members broadcast
// public ephemeral keys generated for other members of the group.
// `ephemeralPublicKeyMessage`s are valid in this state.
type ephemeralKeyPairGenerationState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *ephemeralKeyPairGeneratingMember
}

func (ekpgs *ephemeralKeyPairGenerationState) Initiate(ctx context.Context) error {
	ekpgs.action.run(func() error {
		message, err := ekpgs.member.generateEphemeralKeyPair()
		if err != nil {
			return err
		}

		if err := ekpgs.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (ekpgs *ephemeralKeyPairGenerationState) Receive(
	netMessage net.Message,
) error {
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
	messagingDone := len(receivedMessages[*ephemeralPublicKeyMessage](ekpgs)) ==
		len(ekpgs.member.group.OperatingMemberIDs())-1

	return ekpgs.action.isDone() && messagingDone
}

func (ekpgs *ephemeralKeyPairGenerationState) Next() (faststate.State, error) {
	if err := ekpgs.action.error(); err != nil {
		return nil, err
	}

	return &symmetricKeyGenerationState{
		BaseState: ekpgs.BaseState,
		action:    &stateAction{},
		channel:   ekpgs.channel,
		member:    ekpgs.member.initializeSymmetricKeyGeneration(),
	}, nil
}

func (ekpgs *ephemeralKeyPairGenerationState) MemberIndex() group.MemberIndex {
	return ekpgs.member.id
}

// symmetricKeyGenerationState is the state during which members compute
// symmetric keys from the previously exchanged ephemeral public keys.
// No messages are valid in this state.
type symmetricKeyGenerationState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *symmetricKeyGeneratingMember
}

func (skgs *symmetricKeyGenerationState) Initiate(ctx context.Context) error {
	skgs.action.run(func() error {
		return skgs.member.generateSymmetricKeys(
			receivedMessages[*ephemeralPublicKeyMessage](skgs),
		)
	})

	return nil
}

func (skgs *symmetricKeyGenerationState) Receive(net.Message) error {
	return nil
}

func (skgs *symmetricKeyGenerationState) CanTransition() bool {
	return skgs.action.isDone()
}

func (skgs *symmetricKeyGenerationState) Next() (faststate.State, error) {
	if err := skgs.action.error(); err != nil {
		return nil, err
	}

	return &tssRoundOneState{
		BaseState: skgs.BaseState,
		action:    &stateAction{},
		channel:   skgs.channel,
		member:    skgs.member.initializeTssRoundOne(),
	}, nil
}

func (skgs *symmetricKeyGenerationState) MemberIndex() group.MemberIndex {
	return skgs.member.id
}

// tssRoundOneState is the state during which members broadcast TSS
// round one messages.
// `tssRoundOneMessage`s are valid in this state.
type tssRoundOneState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *tssRoundOneMember
}

func (tros *tssRoundOneState) Initiate(ctx context.Context) error {
	tros.action.run(func() error {
		message, err := tros.member.tssRoundOne(ctx)
		if err != nil {
			return err
		}

		if err := tros.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

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
	messagingDone := len(receivedMessages[*tssRoundOneMessage](tros)) ==
		len(tros.member.group.OperatingMemberIDs())-1

	return tros.action.isDone() && messagingDone
}

func (tros *tssRoundOneState) Next() (faststate.State, error) {
	if err := tros.action.error(); err != nil {
		return nil, err
	}

	return &tssRoundTwoState{
		BaseState: tros.BaseState,
		action:    &stateAction{},
		channel:   tros.channel,
		member:    tros.member.initializeTssRoundTwo(),
	}, nil
}

func (tros *tssRoundOneState) MemberIndex() group.MemberIndex {
	return tros.member.id
}

// tssRoundTwoState is the state during which members broadcast TSS
// round two messages.
// `tssRoundTwoMessage`s are valid in this state.
type tssRoundTwoState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *tssRoundTwoMember
}

func (trts *tssRoundTwoState) Initiate(ctx context.Context) error {
	trts.action.run(func() error {
		message, err := trts.member.tssRoundTwo(
			ctx,
			receivedMessages[*tssRoundOneMessage](trts),
		)
		if err != nil {
			return err
		}

		if err := trts.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

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
	messagingDone := len(receivedMessages[*tssRoundTwoMessage](trts)) ==
		len(trts.member.group.OperatingMemberIDs())-1

	return trts.action.isDone() && messagingDone
}

func (trts *tssRoundTwoState) Next() (faststate.State, error) {
	if err := trts.action.error(); err != nil {
		return nil, err
	}

	return &tssRoundThreeState{
		BaseState: trts.BaseState,
		action:    &stateAction{},
		channel:   trts.channel,
		member:    trts.member.initializeTssRoundThree(),
	}, nil
}

func (trts *tssRoundTwoState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// tssRoundThreeState is the state during which members broadcast TSS
// round three messages.
// `tssRoundThreeMessage`s are valid in this state.
type tssRoundThreeState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *tssRoundThreeMember
}

func (trts *tssRoundThreeState) Initiate(ctx context.Context) error {
	trts.action.run(func() error {
		message, err := trts.member.tssRoundThree(
			ctx,
			receivedMessages[*tssRoundTwoMessage](trts),
		)
		if err != nil {
			return err
		}

		if err := trts.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

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
	messagingDone := len(receivedMessages[*tssRoundThreeMessage](trts)) ==
		len(trts.member.group.OperatingMemberIDs())-1

	return trts.action.isDone() && messagingDone
}

func (trts *tssRoundThreeState) Next() (faststate.State, error) {
	if err := trts.action.error(); err != nil {
		return nil, err
	}

	return &tssRoundFourState{
		BaseState: trts.BaseState,
		action:    &stateAction{},
		channel:   trts.channel,
		member:    trts.member.initializeTssRoundFour(),
	}, nil
}

func (trts *tssRoundThreeState) MemberIndex() group.MemberIndex {
	return trts.member.id
}

// tssRoundFourState is the state during which members broadcast TSS
// round four messages.
// `tssRoundFourMessage`s are valid in this state.
type tssRoundFourState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *tssRoundFourMember
}

func (trfs *tssRoundFourState) Initiate(ctx context.Context) error {
	trfs.action.run(func() error {
		message, err := trfs.member.tssRoundFour(
			ctx,
			receivedMessages[*tssRoundThreeMessage](trfs),
		)
		if err != nil {
			return err
		}

		if err := trfs.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

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
	messagingDone := len(receivedMessages[*tssRoundFourMessage](trfs)) ==
		len(trfs.member.group.OperatingMemberIDs())-1

	return trfs.action.isDone() && messagingDone
}

func (trfs *tssRoundFourState) Next() (faststate.State, error) {
	if err := trfs.action.error(); err != nil {
		return nil, err
	}

	return &tssRoundFiveState{
		BaseState: trfs.BaseState,
		action:    &stateAction{},
		channel:   trfs.channel,
		member:    trfs.member.initializeTssRoundFive(),
	}, nil
}

func (trfs *tssRoundFourState) MemberIndex() group.MemberIndex {
	return trfs.member.id
}

// tssRoundFiveState is the state during which members broadcast TSS
// round five messages.
// `tssRoundFiveMessage`s are valid in this state.
type tssRoundFiveState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *tssRoundFiveMember
}

func (trfs *tssRoundFiveState) Initiate(ctx context.Context) error {
	trfs.action.run(func() error {
		message, err := trfs.member.tssRoundFive(
			ctx,
			receivedMessages[*tssRoundFourMessage](trfs),
		)
		if err != nil {
			return err
		}

		if err := trfs.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

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
	messagingDone := len(receivedMessages[*tssRoundFiveMessage](trfs)) ==
		len(trfs.member.group.OperatingMemberIDs())-1

	return trfs.action.isDone() && messagingDone
}

func (trfs *tssRoundFiveState) Next() (faststate.State, error) {
	if err := trfs.action.error(); err != nil {
		return nil, err
	}

	return &tssRoundSixState{
		BaseState: trfs.BaseState,
		action:    &stateAction{},
		channel:   trfs.channel,
		member:    trfs.member.initializeTssRoundSix(),
	}, nil
}

func (trfs *tssRoundFiveState) MemberIndex() group.MemberIndex {
	return trfs.member.id
}

// tssRoundSixState is the state during which members broadcast TSS
// round six messages.
// `tssRoundSixMessage`s are valid in this state.
type tssRoundSixState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *tssRoundSixMember
}

func (trss *tssRoundSixState) Initiate(ctx context.Context) error {
	trss.action.run(func() error {
		message, err := trss.member.tssRoundSix(
			ctx,
			receivedMessages[*tssRoundFiveMessage](trss),
		)
		if err != nil {
			return err
		}

		if err := trss.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

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
	messagingDone := len(receivedMessages[*tssRoundSixMessage](trss)) ==
		len(trss.member.group.OperatingMemberIDs())-1

	return trss.action.isDone() && messagingDone
}

func (trss *tssRoundSixState) Next() (faststate.State, error) {
	if err := trss.action.error(); err != nil {
		return nil, err
	}

	return &tssRoundSevenState{
		BaseState: trss.BaseState,
		action:    &stateAction{},
		channel:   trss.channel,
		member:    trss.member.initializeTssRoundSeven(),
	}, nil
}

func (trss *tssRoundSixState) MemberIndex() group.MemberIndex {
	return trss.member.id
}

// tssRoundSevenState is the state during which members broadcast TSS
// round seven messages.
// `tssRoundSevenMessage`s are valid in this state.
type tssRoundSevenState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *tssRoundSevenMember
}

func (trss *tssRoundSevenState) Initiate(ctx context.Context) error {
	trss.action.run(func() error {
		message, err := trss.member.tssRoundSeven(
			ctx,
			receivedMessages[*tssRoundSixMessage](trss),
		)
		if err != nil {
			return err
		}

		if err := trss.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

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
	messagingDone := len(receivedMessages[*tssRoundSevenMessage](trss)) ==
		len(trss.member.group.OperatingMemberIDs())-1

	return trss.action.isDone() && messagingDone
}

func (trss *tssRoundSevenState) Next() (faststate.State, error) {
	if err := trss.action.error(); err != nil {
		return nil, err
	}

	return &tssRoundEightState{
		BaseState: trss.BaseState,
		action:    &stateAction{},
		channel:   trss.channel,
		member:    trss.member.initializeTssRoundEight(),
	}, nil
}

func (trss *tssRoundSevenState) MemberIndex() group.MemberIndex {
	return trss.member.id
}

// tssRoundEightState is the state during which members broadcast TSS
// round eight messages.
// `tssRoundEightMessage`s are valid in this state.
type tssRoundEightState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *tssRoundEightMember
}

func (tres *tssRoundEightState) Initiate(ctx context.Context) error {
	tres.action.run(func() error {
		message, err := tres.member.tssRoundEight(
			ctx,
			receivedMessages[*tssRoundSevenMessage](tres),
		)
		if err != nil {
			return err
		}

		if err := tres.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

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
	messagingDone := len(receivedMessages[*tssRoundEightMessage](tres)) ==
		len(tres.member.group.OperatingMemberIDs())-1

	return tres.action.isDone() && messagingDone
}

func (tres *tssRoundEightState) Next() (faststate.State, error) {
	if err := tres.action.error(); err != nil {
		return nil, err
	}

	return &tssRoundNineState{
		BaseState: tres.BaseState,
		action:    &stateAction{},
		channel:   tres.channel,
		member:    tres.member.initializeTssRoundNine(),
	}, nil
}

func (tres *tssRoundEightState) MemberIndex() group.MemberIndex {
	return tres.member.id
}

// tssRoundNineState is the state during which members broadcast TSS
// round nine messages.
// `tssRoundNineMessage`s are valid in this state.
type tssRoundNineState struct {
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *tssRoundNineMember
}

func (trns *tssRoundNineState) Initiate(ctx context.Context) error {
	trns.action.run(func() error {
		message, err := trns.member.tssRoundNine(
			ctx,
			receivedMessages[*tssRoundEightMessage](trns),
		)
		if err != nil {
			return err
		}

		if err := trns.channel.Send(ctx, message); err != nil {
			return err
		}

		return nil
	})

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
	messagingDone := len(receivedMessages[*tssRoundNineMessage](trns)) ==
		len(trns.member.group.OperatingMemberIDs())-1

	return trns.action.isDone() && messagingDone
}

func (trns *tssRoundNineState) Next() (faststate.State, error) {
	if err := trns.action.error(); err != nil {
		return nil, err
	}

	return &finalizationState{
		BaseState: trns.BaseState,
		action:    &stateAction{},
		channel:   trns.channel,
		member:    trns.member.initializeFinalization(),
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
	*faststate.BaseState

	action *stateAction

	channel net.BroadcastChannel
	member  *finalizingMember
}

func (fs *finalizationState) Initiate(ctx context.Context) error {
	fs.action.run(func() error {
		err := fs.member.tssFinalize(
			ctx,
			receivedMessages[*tssRoundNineMessage](fs),
		)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (fs *finalizationState) Receive(net.Message) error {
	return nil
}

func (fs *finalizationState) CanTransition() bool {
	return fs.action.isDone()
}

func (fs *finalizationState) Next() (faststate.State, error) {
	if err := fs.action.error(); err != nil {
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

// stateAction represents an asynchronous action performed in the given
// protocol state.
type stateAction struct {
	mutex   sync.RWMutex
	running bool
	done    bool
	err     error
}

// run triggers the action goroutine. Can be called only once. Successive
// calls do nothing.
func (sa *stateAction) run(actionFn func() error) {
	sa.mutex.Lock()
	defer sa.mutex.Unlock()

	if sa.running || sa.done {
		return
	}

	sa.running = true

	go func() {
		err := actionFn()

		sa.mutex.Lock()

		sa.running = false
		sa.done = true
		sa.err = err

		sa.mutex.Unlock()
	}()
}

// isDone returns whether the state's action is done.
func (sa *stateAction) isDone() bool {
	sa.mutex.RLock()
	defer sa.mutex.RUnlock()
	return sa.done
}

// error returns the state's action error if any. Calling this function makes
// sense only when the action is done, i.e. the isDone function returns true.
// After this function returns a non-nil error, successive calls return the
// same error.
func (sa *stateAction) error() error {
	sa.mutex.RLock()
	defer sa.mutex.RUnlock()
	return sa.err
}

// messageReceiverState is a type constraint that refers to a state which is
// supposed to receive network messages.
type messageReceiverState interface {
	GetAllReceivedMessages(messageType string) []net.Message
}

// receivedMessages returns all messages of type T that have been received
// and validated so far. Returned messages are deduplicated so there is a
// guarantee that only one message of the given type is returned for the
// given sender.
func receivedMessages[T message, S messageReceiverState](state S) []T {
	var template T

	payloads := make([]T, 0)
	for _, msg := range state.GetAllReceivedMessages(template.Type()) {
		payload, ok := msg.Payload().(T)
		if !ok {
			continue
		}

		payloads = append(payloads, payload)
	}

	return deduplicateBySender(payloads)
}

// deduplicateBySender removes duplicated items for the given sender.
// It always takes the first item that occurs for the given sender
// and ignores the subsequent ones.
func deduplicateBySender[T interface{ SenderID() group.MemberIndex }](
	list []T,
) []T {
	senders := make(map[group.MemberIndex]bool)
	result := make([]T, 0)

	for _, item := range list {
		if _, exists := senders[item.SenderID()]; !exists {
			senders[item.SenderID()] = true
			result = append(result, item)
		}
	}

	return result
}