// Package local provides a local, non-networked implementation of the
// interfaces defined by the net package. It should largely be considered a
// sample implementation, and is not meant to be used at scale in any way.
package local

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/internal"
)

type localIdentifier string

func (li localIdentifier) String() string {
	return string(li)
}

var channelsMutex sync.Mutex
var channels map[string][]*localChannel

type localProvider struct {
	id localIdentifier
}

func (lp *localProvider) ID() net.TransportIdentifier {
	return lp.id
}

func (lp *localProvider) ChannelFor(name string) (net.BroadcastChannel, error) {
	return channel(name), nil
}

func (lp *localProvider) Type() string {
	return "local"
}

func (lp *localProvider) AddrStrings() []string {
	return make([]string, 0)
}

func (lp *localProvider) Peers() []string {
	return make([]string, 0)
}

// Connect returns a local instance of a net provider that does not go over the
// network.
func Connect() net.Provider {
	return &localProvider{
		id: localIdentifier(randomIdentifier()),
	}
}

// channel returns a BroadcastChannel designed to mediate between local
// participants. It delivers all messages sent to the channel through its
// receive channels. RecvChan on a LocalChannel creates a new receive channel
// that is returned to the caller, so that all receive channels can receive
// the message.
func channel(name string) net.BroadcastChannel {
	channelsMutex.Lock()
	defer channelsMutex.Unlock()
	if channels == nil {
		channels = make(map[string][]*localChannel)
	}

	localChannels, exists := channels[name]
	if !exists {
		localChannels = make([]*localChannel, 0)
		channels[name] = localChannels
	}

	identifier := localIdentifier(randomIdentifier())
	channel := &localChannel{
		name:                 name,
		identifier:           &identifier,
		messageHandlersMutex: sync.Mutex{},
		messageHandlers:      make([]net.HandleMessageFunc, 0),
		unmarshalersMutex:    sync.Mutex{},
		unmarshalersByType:   make(map[string]func() net.TaggedUnmarshaler, 0),
	}
	channels[name] = append(channels[name], channel)

	return channel
}

var letterRunes = [52]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y',
	'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func randomIdentifier() string {
	runes := make([]rune, 32)
	for i := range runes {
		runes[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(runes)
}

type localChannel struct {
	name                 string
	identifier           net.TransportIdentifier
	messageHandlersMutex sync.Mutex
	messageHandlers      []net.HandleMessageFunc
	unmarshalersMutex    sync.Mutex
	unmarshalersByType   map[string]func() net.TaggedUnmarshaler
}

func (lc *localChannel) Name() string {
	return lc.name
}

func doSend(channel *localChannel, payload net.TaggedMarshaler) error {
	channelsMutex.Lock()
	targetChannels := channels[channel.name]
	channelsMutex.Unlock()

	bytes, err := payload.Marshal()
	if err != nil {
		return err
	}

	unmarshaler, found := channel.unmarshalersByType[payload.Type()]
	if !found {
		return fmt.Errorf("Couldn't find unmarshaler for type %s", payload.Type())
	}

	unmarshaled := unmarshaler()
	err = unmarshaled.Unmarshal(bytes)
	if err != nil {
		return err
	}

	for _, targetChannel := range targetChannels {
		targetChannel.deliver(channel.identifier, unmarshaled) // TODO error handling?
	}

	return nil
}

func (lc *localChannel) deliver(transportIdentifier net.TransportIdentifier, payload interface{}) {
	lc.messageHandlersMutex.Lock()
	snapshot := make([]net.HandleMessageFunc, len(lc.messageHandlers))
	copy(snapshot, lc.messageHandlers)
	lc.messageHandlersMutex.Unlock()

	message :=
		internal.BasicMessage(
			transportIdentifier,
			payload,
			"local",
		)

	go func() {
		for _, handler := range snapshot {
			handler.Handler(message)
		}
	}()
}

func (lc *localChannel) Send(message net.TaggedMarshaler) error {
	return doSend(lc, message)
}

func (lc *localChannel) Recv(handler net.HandleMessageFunc) error {
	lc.messageHandlersMutex.Lock()
	lc.messageHandlers = append(lc.messageHandlers, handler)
	lc.messageHandlersMutex.Unlock()

	return nil
}

func (lc *localChannel) UnregisterRecv(handlerType string) error {
	lc.messageHandlersMutex.Lock()
	defer lc.messageHandlersMutex.Unlock()
	removedCount := 0
	for i, mh := range lc.messageHandlers {
		if mh.Type == handlerType {
			removedCount++
			lc.messageHandlers[i] = lc.messageHandlers[len(lc.messageHandlers)-removedCount]
		}
	}
	lc.messageHandlers = lc.messageHandlers[:len(lc.messageHandlers)-removedCount]

	return nil
}

func (lc *localChannel) RegisterUnmarshaler(
	unmarshaler func() net.TaggedUnmarshaler,
) (err error) {
	tpe := unmarshaler().Type()

	lc.unmarshalersMutex.Lock()
	_, exists := lc.unmarshalersByType[tpe]
	if exists {
		err = fmt.Errorf("type %s already has an associated unmarshaler", tpe)
	} else {
		lc.unmarshalersByType[tpe] = unmarshaler
	}
	lc.unmarshalersMutex.Unlock()
	return
}
