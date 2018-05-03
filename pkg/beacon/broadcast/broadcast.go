package broadcast

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
)

// Message represents a message to communicate over a broadcast channel.
// TODO Combine with Raghav's work on protobuf messages.
type Message struct {
	Sender   bls.ID
	Receiver *bls.ID // pointer so it can be nil for broadcast messages
	Data     interface{}
}

// NewBroadcastMessage creates a new message from the given sender, carrying the
// given data payload, meant for broacast into a channel watched by others. The
// message is signed, but not encrypted.
func NewBroadcastMessage(sender bls.ID, data interface{}) Message {
	// FIXME Sign, will require private key...
	return Message{sender, nil, data}
}

// NewPrivateMessage creates a new private message from the given sender to the
// given receiver, carrying the given data payload, meant for broacast into a
// channel watched by others. The message is signed and encrypted.
func NewPrivateMessage(sender bls.ID, receiver bls.ID, data interface{}) Message {
	// FIXME Actually encrypt here... Will require a key, best taken from
	// FIXME chain...
	return Message{sender, &receiver, data}
}
