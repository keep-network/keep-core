package local

import (
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
)

// Channel returns a BroadcastChannel designed to mediate between local
// participants. It delivers all messages sent to the channel through its
// receive channels. RecvChan on a LocalChannel creates a new receive channel
// that is returned to the caller, so that all receive channels can receive
// the message.
func Channel(name string) net.BroadcastChannel {
	return &localChannel{
		name,
		sync.Mutex{},
		make([]net.HandleMessageFunc, 0),
		sync.Mutex{},
		make(map[string]func() net.TaggedUnmarshaler, 0)}
}

type localChannel struct {
	name                 string
	messageHandlersMutex sync.Mutex
	messageHandlers      []net.HandleMessageFunc
	unmarshalersMutex    sync.Mutex
	unmarshalersByType   map[string]func() net.TaggedUnmarshaler
}

func (channel *localChannel) Name() string {
	return channel.name
}

func doSend(
	channel *localChannel,
	recipient *net.ClientIdentifier,
	message net.TaggedMarshaler,
) error {
	channel.messageHandlersMutex.Lock()
	snapshot := make([]net.HandleMessageFunc, len(channel.messageHandlers))
	copy(snapshot, channel.messageHandlers)
	channel.messageHandlersMutex.Unlock()

	bytes, err := message.Marshal()
	if err != nil {
		return err
	}

	unmarshaler, found := channel.unmarshalersByType[message.Type()]
	if !found {
		return fmt.Errorf("Couldn't find unmarshaler for type %s", message.Type())
	}

	unmarshaled := unmarshaler()
	err = unmarshaled.Unmarshal(bytes)
	if err != nil {
		return err
	}

	go func() {
		for _, messageHandler := range snapshot {
			messageHandler(unmarshaled) // TODO error handling?
		}
	}()

	return nil
}

func (channel *localChannel) Send(message net.TaggedMarshaler) error {
	return doSend(channel, nil, message)
}

func (channel *localChannel) SendTo(
	recipient net.ClientIdentifier,
	message net.TaggedMarshaler) error {
	return doSend(channel, &recipient, message)
}

func (channel *localChannel) Recv(handler net.HandleMessageFunc) error {
	channel.messageHandlersMutex.Lock()
	channel.messageHandlers = append(channel.messageHandlers, handler)
	channel.messageHandlersMutex.Unlock()

	return nil
}

func (channel *localChannel) RegisterUnmarshaler(
	unmarshaler func() net.TaggedUnmarshaler,
) (err error) {
	tpe := unmarshaler().Type()

	channel.unmarshalersMutex.Lock()
	_, exists := channel.unmarshalersByType[tpe]
	if exists {
		err = fmt.Errorf("type %s already has an associated unmarshaler", tpe)
	} else {
		channel.unmarshalersByType[tpe] = unmarshaler
	}
	channel.unmarshalersMutex.Unlock()
	return
}
