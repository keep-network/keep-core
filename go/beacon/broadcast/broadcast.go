package broadcast

import (
	"sync"

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

// Channel represents a named broadcast channel. It allows consumers to send
// messages to the channel (via Send) and to access a low-level receive chan
// that furnishes messages sent onto the broadcast channel.
type Channel interface {
	Name() string

	Send(message Message) bool

	RecvChan() <-chan Message
}

type localChannel struct {
	name           string
	recvChansMutex sync.Mutex
	recvChans      []chan Message
}

func (channel *localChannel) Name() string {
	return channel.name
}

func (channel *localChannel) Send(message Message) bool {
	channel.recvChansMutex.Lock()
	go func(recvChans []chan Message) {
		for _, recvChan := range recvChans {
			recvChan <- message
		}
	}(channel.recvChans)
	channel.recvChansMutex.Unlock()

	return true
}

func (channel *localChannel) RecvChan() <-chan Message {
	newChan := make(chan Message)

	channel.recvChansMutex.Lock()
	channel.recvChans = append(channel.recvChans, newChan)
	channel.recvChansMutex.Unlock()

	return newChan
}

// LocalChannel returns a Channel designed to mediate between local
// participants. It delivers all messages sent to the channel through its
// receive channels. RecvChan on a LocalChannel creates a new receive channel
// that is returned to the caller, so that all receive channels can receive
// the message.
func LocalChannel(name string) Channel {
	return &localChannel{name, sync.Mutex{}, make([]chan Message, 0)}
}
