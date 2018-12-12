package dkg2

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/net"
)

type keyGenerationState interface {
	activeBlocks() int

	initiate() error
	receive(msg net.Message) error
	nextState() (keyGenerationState, error)
}

// initializationState is the starting state of key generation; it waits for
// activePeriod and then enters joinState. No messages are valid in this state.
type initializationState struct {
	channel net.BroadcastChannel
	member  *gjkr.EphemeralKeyPairGeneratingMember
}

func (is *initializationState) activeBlocks() int { return 1 }

func (is *initializationState) initiate() error {
	return nil
}

func (is *initializationState) receive(msg net.Message) error {
	return fmt.Errorf("unexpected message for initialization state: [%#v]", msg)
}

func (is *initializationState) nextState() (keyGenerationState, error) {
	return &joinState{is.channel, is.member}, nil
}

// joinState is the state during which a member announces itself to the key
// generation broadcast channel to initiate the distributed protocol.
// `JoinMessage`s from other members are valid in this state.
type joinState struct {
	channel net.BroadcastChannel
	member  *gjkr.EphemeralKeyPairGeneratingMember
}

func (js *joinState) activeBlocks() int { return 1 }

func (js *joinState) initiate() error {
	return nil
}

func (js *joinState) receive(msg net.Message) error {
	switch joinMsg := msg.Payload().(type) {
	case *gjkr.JoinMessage:
		if err := js.channel.RegisterIdentifier(
			msg.TransportSenderID(),
			joinMsg.SenderID,
		); err != nil {
			return err
		}

		js.member.AddToGroup(joinMsg.SenderID)
	}
	return nil
}

func (js *joinState) nextState() (keyGenerationState, error) {
	return &ephemeralKeyPairGeneratingState{
		js.channel,
		js.member,
		make([]*gjkr.EphemeralPublicKeyMessage, 0),
	}, nil
}

// ephemeralKeyPairGeneratingState is the state during which members broadcast
// publish ephemeral keys generated for each other member in the group.
// `EphemeralPublicKeyMessage`s from other members are valid in this state.
type ephemeralKeyPairGeneratingState struct {
	channel       net.BroadcastChannel
	member        *gjkr.EphemeralKeyPairGeneratingMember
	phaseMessages []*gjkr.EphemeralPublicKeyMessage
}

func (ekpgs *ephemeralKeyPairGeneratingState) activeBlocks() int { return 1 }

func (ekpgs *ephemeralKeyPairGeneratingState) initiate() error {
	message, err := ekpgs.member.GenerateEphemeralKeyPair()
	if err != nil {
		return fmt.Errorf("ephemeral key generation phase failed [%v]", err)
	}

	if err := ekpgs.channel.Send(message); err != nil {
		return fmt.Errorf("ephemeral key generation phase failed [%v]", err)
	}
	return nil
}
func (ekpgs *ephemeralKeyPairGeneratingState) receive(msg net.Message) error {
	switch publicKeyMessage := msg.Payload().(type) {
	case *gjkr.EphemeralPublicKeyMessage:
		if senderID, ok := msg.ProtocolSenderID().(gjkr.MemberID); ok {
			if senderID == ekpgs.member.ID {
				return nil // ignore message from self
			}
			ekpgs.phaseMessages = append(ekpgs.phaseMessages, publicKeyMessage)
			return nil
		}

		return fmt.Errorf(
			"unknown protocol sender id type [%T]  [%v]",
			msg.ProtocolSenderID(),
			msg.TransportSenderID(),
		)
	}

	return fmt.Errorf(
		"unexpected message for ephemeral key generation state: [%#v]",
		msg,
	)
}
func (ekpgs *ephemeralKeyPairGeneratingState) nextState() (keyGenerationState, error) {
	return nil, nil
}
