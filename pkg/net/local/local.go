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

func (localIdentifier) ProviderName() string {
	return "local"
}

func (li localIdentifier) String() string {
	return string(li)
}

var channels map[string][]*localChannel

// Channel returns a BroadcastChannel designed to mediate between local
// participants. It delivers all messages sent to the channel through its
// receive channels. RecvChan on a LocalChannel creates a new receive channel
// that is returned to the caller, so that all receive channels can receive
// the message.
func Channel(name string) net.BroadcastChannel {
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
		name:                        name,
		identifier:                  &identifier,
		messageHandlersMutex:        sync.Mutex{},
		messageHandlers:             make([]net.HandleMessageFunc, 0),
		unmarshalersMutex:           sync.Mutex{},
		unmarshalersByType:          make(map[string]func() net.TaggedUnmarshaler, 0),
		identifiersMutex:            sync.Mutex{},
		transportToProtoIdentifiers: make(map[net.TransportIdentifier]net.ProtocolIdentifier),
		protoToTransportIdentifiers: make(map[net.ProtocolIdentifier]net.TransportIdentifier),
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
	name                        string
	identifier                  net.TransportIdentifier
	messageHandlersMutex        sync.Mutex
	messageHandlers             []net.HandleMessageFunc
	unmarshalersMutex           sync.Mutex
	unmarshalersByType          map[string]func() net.TaggedUnmarshaler
	identifiersMutex            sync.Mutex
	transportToProtoIdentifiers map[net.TransportIdentifier]net.ProtocolIdentifier
	protoToTransportIdentifiers map[net.ProtocolIdentifier]net.TransportIdentifier
}

func (channel *localChannel) Name() string {
	return channel.name
}

func doSend(
	channel *localChannel,
	recipient interface{},
	payload net.TaggedMarshaler,
) error {
	targetChannels := channels[channel.name]

	// If we have a recipient, filter `targetChannels` down to only the targeted
	// recipient (the recipient transport identifier is the same as the local
	// channel's identifier).
	var transportRecipient net.TransportIdentifier
	channel.identifiersMutex.Lock()
	if transportID, ok := recipient.(*localIdentifier); ok {
		transportRecipient = transportID
	} else if transportID, ok := channel.protoToTransportIdentifiers[recipient]; ok {
		transportRecipient = transportID
	}
	channel.identifiersMutex.Unlock()

	if transportRecipient != nil {
		potentialTargets := targetChannels
		targetChannels = make([]*localChannel, 0, 1)
		for _, targetChannel := range potentialTargets {
			if transportRecipient == targetChannel.identifier {
				targetChannels = append(targetChannels, targetChannel)
				break
			}
		}
	}

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

func (channel *localChannel) deliver(senderIdentifier net.TransportIdentifier, payload interface{}) {
	channel.messageHandlersMutex.Lock()
	snapshot := make([]net.HandleMessageFunc, len(channel.messageHandlers))
	copy(snapshot, channel.messageHandlers)
	channel.messageHandlersMutex.Unlock()

	channel.identifiersMutex.Lock()
	protocolIdentifier := channel.transportToProtoIdentifiers[senderIdentifier]
	channel.identifiersMutex.Unlock()

	message :=
		internal.BasicMessage(
			senderIdentifier,
			protocolIdentifier,
			payload,
		)

	go func() {
		for _, handler := range snapshot {
			handler(message)
		}
	}()
}

func (channel *localChannel) Send(message net.TaggedMarshaler) error {
	return doSend(channel, nil, message)
}

func (channel *localChannel) SendTo(
	recipient net.ProtocolIdentifier,
	message net.TaggedMarshaler) error {
	return doSend(channel, recipient, message)
}

func (channel *localChannel) Recv(handler net.HandleMessageFunc) error {
	channel.messageHandlersMutex.Lock()
	channel.messageHandlers = append(channel.messageHandlers, handler)
	channel.messageHandlersMutex.Unlock()

	return nil
}

func (channel *localChannel) RegisterIdentifier(
	transportIdentifier net.TransportIdentifier,
	protocolIdentifier net.ProtocolIdentifier,
) error {
	channel.identifiersMutex.Lock()
	defer channel.identifiersMutex.Unlock()

	if _, exists := channel.transportToProtoIdentifiers[transportIdentifier]; exists {
		return fmt.Errorf(
			"already have a protocol identifier associated with [%v]",
			transportIdentifier)
	}
	if _, exists := channel.protoToTransportIdentifiers[protocolIdentifier]; exists {
		return fmt.Errorf(
			"already have a transport identifier associated with [%v]",
			protocolIdentifier)
	}

	channel.transportToProtoIdentifiers[transportIdentifier] = protocolIdentifier
	channel.protoToTransportIdentifiers[protocolIdentifier] = transportIdentifier

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
