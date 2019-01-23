package internal

import "github.com/keep-network/keep-core/pkg/net"

// BasicMessage returns a struct-based trivial implementation of the net.Message
// interface for use by packages that don't need any frills.
func BasicMessage(
	protocolSenderID net.ProtocolIdentifier,
	payload interface{},
	messageType string,
) net.Message {
	return &basicMessage{
		protocolSenderID,
		payload,
		messageType,
	}
}

// basicMessage is a struct-based trivial implementation of the net.Message
// interface for use by packages that don't need any frills.
type basicMessage struct {
	protocolSenderID net.ProtocolIdentifier
	payload          interface{}
	messageType      string
}

func (m *basicMessage) ProtocolSenderID() net.ProtocolIdentifier {
	return m.protocolSenderID
}

func (m *basicMessage) Payload() interface{} {
	return m.payload
}

func (m *basicMessage) Type() string {
	return m.messageType
}
