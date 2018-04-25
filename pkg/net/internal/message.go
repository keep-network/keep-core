package internal

import "github.com/keep-network/keep-core/pkg/net"

// BasicMessage returns a struct-based trivial implementation of the net.Message
// interface for use by packages that don't need any frills.
func BasicMessage(
	networkSenderID net.ClientIdentifier,
	protocolSenderID net.ProtocolIdentifier,
	payload interface{},
) net.Message {
	return &basicMessage{networkSenderID,
		protocolSenderID,
		payload}
}

// basicMessage is a struct-based trivial implementation of the net.Message
// interface for use by packages that don't need any frills.
type basicMessage struct {
	networkSenderID  net.ClientIdentifier
	protocolSenderID net.ProtocolIdentifier
	payload          interface{}
}

func (m *basicMessage) NetworkSenderID() net.ClientIdentifier {
	return m.networkSenderID
}

func (m *basicMessage) ProtocolSenderID() net.ProtocolIdentifier {
	return m.protocolSenderID
}

func (m *basicMessage) Payload() interface{} {
	return m.payload
}
