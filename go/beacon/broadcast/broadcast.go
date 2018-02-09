package broadcast

import "github.com/dfinity/go-dfinity-crypto/bls"

// Message represents a message to communicate over a broadcast channel.
// TODO Combine with Raghav's work on protobuf messages.
type Message struct {
	sender    bls.ID
	receiver  bls.ID
	encrypted bool
	data      interface{}
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
	name      string
	recvChans []chan Message
}

func (channel *localChannel) Name() string {
	return channel.name
}

func (channel *localChannel) Send(message Message) bool {
	for _, recvChan := range channel.recvChans {
		recvChan <- message
	}

	return true
}

func (channel *localChannel) RecvChan() <-chan Message {
	newChan := make(chan Message)

	channel.recvChans = append(channel.recvChans, newChan)

	return newChan
}

// LocalChannel returns a Channel designed to mediate between local
// participants. It delivers all messages sent to the channel through its
// receive channels. RecvChan on a LocalChannel creates a new receive channel
// that is returned to the caller, so that all receive channels can receive
// the message.
func LocalChannel(name string) Channel {
	return &localChannel{name, make([]chan Message, 0)}
}
